import { promptIdContext, TerminalQuotaError } from "@google/gemini-cli-core";
import { GeminiCliAgent, tool, z } from "@google/gemini-cli-sdk";
import { Buffer } from "node:buffer";
import fs from "node:fs";
import fsPromises from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import {
  bumpTurnCapacity,
  extractIntent,
  getMaxTurnsForModel,
  promptTurnBudgetExhaustion,
  requestHandoffSummary,
  syncSessionTools,
  TurnTracker,
} from "./turn-accounting.js";
import {
  getRepoRoot,
  loadWorkspaceSkillsContext,
  renderProseCard,
  renderToolCard,
  renderTopicBanner,
} from "./utils.js";

const MAX_SILENT_RETRY_DELAY_MS = 300000; // 5 minutes

// Intercept retryDelayMs of TerminalQuotaError to prevent massive retry hangs (> 5 minutes)
try {
  if (TerminalQuotaError && TerminalQuotaError.prototype) {
    Object.defineProperty(TerminalQuotaError.prototype, "name", {
      get() {
        return this._name;
      },
      set(val) {
        this._name = val;
        Object.defineProperty(this, "retryDelayMs", {
          get() {
            return this._retryDelayMs;
          },
          set(delayVal) {
            this._retryDelayMs = delayVal;
            if (
              delayVal !== undefined &&
              delayVal > MAX_SILENT_RETRY_DELAY_MS
            ) {
              if (typeof this.message === "string") {
                this.message = this.message.replace(
                  /exhausted your capacity|capacity exceeded|MODEL_CAPACITY_EXHAUSTED/gi,
                  "exhausted capacity (immediate fallback)",
                );
              }
              if (typeof this._reason === "string") {
                this._reason = this._reason.replace(
                  /MODEL_CAPACITY_EXHAUSTED|MODEL_CAPACITY_EXCEEDED/g,
                  "MODEL_CAPACITY_EXHAUSTED_IMMEDIATE_FALLBACK",
                );
              }
            }
          },
          configurable: true,
          enumerable: true,
        });

        Object.defineProperty(this, "reason", {
          get() {
            return this._reason;
          },
          set(reasonVal) {
            if (
              reasonVal !== undefined &&
              this.retryDelayMs > MAX_SILENT_RETRY_DELAY_MS
            ) {
              this._reason = reasonVal.replace(
                /MODEL_CAPACITY_EXHAUSTED|MODEL_CAPACITY_EXCEEDED/g,
                "MODEL_CAPACITY_EXHAUSTED_IMMEDIATE_FALLBACK",
              );
            } else {
              this._reason = reasonVal;
            }
          },
          configurable: true,
          enumerable: true,
        });
      },
      configurable: true,
      enumerable: true,
    });
  }
} catch (err) {
  console.debug(
    `[DEBUG] Failed to patch TerminalQuotaError prototype: ${err.message}`,
  );
}

let initPromise = null;
let repoRoot = null;
let logStream = null;
let isConsoleIntercepted = false;
let originalStdoutWrite = null;
let originalStderrWrite = null;
let activeAbortHandler = null;

let MODEL_PRO = "gemini-3.1-pro-preview";
let MODEL_FLASH = "gemini-3.8-flash";
let MODEL_FLASH_LITE = "gemini-3.1-flash-lite";
let MODEL_HIERARCHY = [MODEL_PRO, MODEL_FLASH, MODEL_FLASH_LITE];

/**
 * Determines whether an error is an AbortError from an aborted stream or controller.
 * @param {any} err
 * @returns {boolean}
 */
export function isAbortError(err) {
  return (
    err?.name === "AbortError" ||
    err?.type === "aborted" ||
    err?.code === "ABORT_ERR"
  );
}

function writeLogAsync(msg) {
  if (
    logStream &&
    logStream.writable &&
    !logStream.writableEnded &&
    !logStream.destroyed &&
    !logStream.errored
  ) {
    logStream.write(msg + "\n");
  }
}

/**
 * Flushes and closes the active debug log stream if open.
 * @returns {Promise<void>}
 */
export async function flushLogs() {
  if (isConsoleIntercepted) {
    if (originalStdoutWrite) {
      process.stdout.write = originalStdoutWrite;
    }
    if (originalStderrWrite) {
      process.stderr.write = originalStderrWrite;
    }
    originalStdoutWrite = null;
    originalStderrWrite = null;
    isConsoleIntercepted = false;
  }
  if (activeAbortHandler) {
    process.removeListener("uncaughtException", activeAbortHandler);
    activeAbortHandler = null;
  }
  if (!logStream) {
    return;
  }
  const stream = logStream;
  logStream = null;

  if (stream.destroyed || stream.errored || stream.writableEnded) {
    initPromise = null;
    return;
  }

  return new Promise((resolve) => {
    const timer = setTimeout(() => {
      initPromise = null;
      resolve();
    }, 1000);
    try {
      stream.end(() => {
        clearTimeout(timer);
        initPromise = null;
        resolve();
      });
    } catch {
      clearTimeout(timer);
      initPromise = null;
      resolve();
    }
  });
}

function shouldRedirectLog(msg) {
  if (typeof msg !== "string") {
    return false;
  }
  return (
    msg.includes("[DEBUG]") ||
    msg.includes("[PolicyEngine.check]") ||
    msg.includes("[Routing]") ||
    msg.includes("[TopicTool]") ||
    msg.includes("Experiments loaded") ||
    msg.includes("Loading ignore patterns") ||
    msg.includes("Ripgrep is not available") ||
    /Tool with name ['"].*?['"] (?:already registered|has already been)/i.test(
      msg,
    ) ||
    msg.includes("GrepLogic:") ||
    msg.includes("Loaded cached credentials")
  );
}

function setupConsoleIntercept() {
  if (isConsoleIntercepted) {
    return;
  }
  isConsoleIntercepted = true;
  originalStdoutWrite = process.stdout.write;
  originalStderrWrite = process.stderr.write;

  function createInterceptor(originalWrite) {
    return function (chunk, encoding, callback) {
      const str = chunk ? chunk.toString() : "";

      if (
        logStream &&
        logStream.writable &&
        !logStream.writableEnded &&
        !logStream.destroyed &&
        !logStream.errored
      ) {
        if (typeof encoding === "string") {
          logStream.write(chunk, encoding);
        } else {
          logStream.write(chunk);
        }
      }

      if (shouldRedirectLog(str)) {
        if (typeof encoding === "function") {
          encoding();
        } else if (typeof callback === "function") {
          callback();
        }
        return true;
      } else {
        return originalWrite.call(this, chunk, encoding, callback);
      }
    };
  }

  process.stdout.write = createInterceptor(originalStdoutWrite);
  process.stderr.write = createInterceptor(originalStderrWrite);
}

/**
 * Initializes configuration state and logging streams asynchronously.
 * Can be called explicitly or lazily when runAgentSession executes.
 * @returns {Promise<{repoRoot: string, models: {pro: string, flash: string, flash_lite: string}}>}
 */
export async function initializeAgentRunner() {
  if (initPromise) {
    return initPromise;
  }
  initPromise = (async () => {
    repoRoot = await getRepoRoot();
    const logDir = path.join(repoRoot, ".gemini", "tmp");
    await fsPromises.mkdir(logDir, { recursive: true });
    const logPath = path.join(logDir, "orchestrator.log");
    logStream = fs.createWriteStream(logPath, { flags: "a" });
    logStream.on("error", (err) => {
      console.debug(
        `[DEBUG] Debug log stream error encountered: ${err.message}`,
      );
    });
    setupConsoleIntercept();

    // Prevent unhandled AbortError from crashing Node process when Gemini SDK stream is aborted
    if (!activeAbortHandler) {
      activeAbortHandler = (err, origin) => {
        if (isAbortError(err)) {
          return;
        }
        process.removeListener("uncaughtException", activeAbortHandler);
        activeAbortHandler = null;
        if (process.listenerCount("uncaughtException") > 0) {
          process.emit("uncaughtException", err, origin);
        } else {
          throw err;
        }
      };
      process.on("uncaughtException", activeAbortHandler);
    }

    let modelConfig;
    try {
      const settingsStr = await fsPromises.readFile(
        path.join(repoRoot, ".gemini/settings.json"),
        "utf8",
      );
      const settingsObj = JSON.parse(settingsStr);
      modelConfig = settingsObj.models || {};
    } catch (err) {
      console.debug(
        `[DEBUG] Could not load .gemini/settings.json models: ${err.message}. Using defaults.`,
      );
      modelConfig = {};
    }

    MODEL_PRO = modelConfig.pro || "gemini-3.1-pro-preview";
    MODEL_FLASH = modelConfig.flash || "gemini-3.8-flash";
    MODEL_FLASH_LITE = modelConfig.flash_lite || "gemini-3.1-flash-lite";
    MODEL_HIERARCHY = [MODEL_PRO, MODEL_FLASH, MODEL_FLASH_LITE];

    return {
      repoRoot,
      logPath,
      models: {
        pro: MODEL_PRO,
        flash: MODEL_FLASH,
        flash_lite: MODEL_FLASH_LITE,
      },
    };
  })().catch((err) => {
    initPromise = null;
    throw err;
  });

  return initPromise;
}

function getModelFallbackSequence(requestedModel) {
  let modelName = requestedModel;
  if (!modelName || modelName === "auto-gemini-3") {
    modelName = MODEL_PRO;
  }
  const index = MODEL_HIERARCHY.indexOf(modelName);
  if (index === -1) {
    return [requestedModel, ...MODEL_HIERARCHY];
  }
  return MODEL_HIERARCHY.slice(index);
}

const STANDARD_STRING_PROPS = ["message", "name", "reason", "code", "status"];
const NESTED_KEYS = ["error", "cause", "details", "errors"];
const IGNORED_KEYS = new Set([
  "config",
  "request",
  "response",
  "socket",
  "req",
  "res",
  "stack",
  "options",
  "params",
  "agent",
  "session",
  "runner",
  "client",
  ...STANDARD_STRING_PROPS,
  ...NESTED_KEYS,
]);

const MAX_TURNS_REGEX =
  /\b(?:(?:maximum|max)[\s_-]+turns?|turn[\s_-]+limits?|(?:max[-_]?turns?|turn[-_]?limits?)\w*)\b/i;

function safeGet(obj, prop) {
  if (obj === null || obj === undefined) {
    return undefined;
  }
  try {
    return obj[prop];
  } catch (err) {
    void err;
    return undefined;
  }
}

function testObjectForMaxTurns(target, depth = 0, seen = new Set()) {
  if (!target || depth > 4) {
    return false;
  }
  if (typeof target === "string") {
    return MAX_TURNS_REGEX.test(target);
  }
  if (typeof target !== "object") {
    return MAX_TURNS_REGEX.test(String(target));
  }
  if (seen.has(target)) {
    return false;
  }
  seen.add(target);

  if (Buffer.isBuffer(target)) {
    if (target.length <= 64 * 1024) {
      const text = target.toString("utf8");
      try {
        const parsed = JSON.parse(text);
        if (typeof parsed === "object" && parsed !== null) {
          return testObjectForMaxTurns(parsed, depth + 1, seen);
        }
      } catch (parseErr) {
        void parseErr;
      }
      return MAX_TURNS_REGEX.test(text);
    }
    return false;
  }

  if (
    typeof ArrayBuffer !== "undefined" &&
    (target instanceof ArrayBuffer || ArrayBuffer.isView(target))
  ) {
    return false;
  }

  if (Array.isArray(target)) {
    for (const item of target) {
      if (testObjectForMaxTurns(item, depth + 1, seen)) {
        return true;
      }
    }
    return false;
  }

  for (const prop of STANDARD_STRING_PROPS) {
    const val = safeGet(target, prop);
    if (val !== null && val !== undefined) {
      if (typeof val !== "object") {
        if (MAX_TURNS_REGEX.test(String(val))) {
          return true;
        }
      } else if (testObjectForMaxTurns(val, depth + 1, seen)) {
        return true;
      }
    }
  }

  for (const key of NESTED_KEYS) {
    const val = safeGet(target, key);
    if (val !== null && val !== undefined) {
      if (typeof val !== "object") {
        if (MAX_TURNS_REGEX.test(String(val))) {
          return true;
        }
      } else if (testObjectForMaxTurns(val, depth + 1, seen)) {
        return true;
      }
    }
  }

  const responseObj = safeGet(target, "response");
  if (responseObj && typeof responseObj === "object") {
    const data = safeGet(responseObj, "data");
    if (data !== null && data !== undefined) {
      if (testObjectForMaxTurns(data, depth + 1, seen)) {
        return true;
      }
    }
    const body = safeGet(responseObj, "body");
    if (body !== null && body !== undefined && body !== data) {
      if (testObjectForMaxTurns(body, depth + 1, seen)) {
        return true;
      }
    }
    const statusText = safeGet(responseObj, "statusText");
    if (statusText !== null && statusText !== undefined) {
      if (typeof statusText !== "object") {
        if (MAX_TURNS_REGEX.test(String(statusText))) {
          return true;
        }
      } else if (testObjectForMaxTurns(statusText, depth + 1, seen)) {
        return true;
      }
    }
  } else if (
    responseObj !== null &&
    responseObj !== undefined &&
    typeof responseObj !== "object"
  ) {
    if (MAX_TURNS_REGEX.test(String(responseObj))) {
      return true;
    }
  }

  let keys = [];
  try {
    keys = Object.keys(target);
  } catch (err) {
    void err;
  }

  for (const key of keys) {
    if (IGNORED_KEYS.has(key)) {
      continue;
    }
    const val = safeGet(target, key);
    if (
      val !== false &&
      val !== null &&
      val !== undefined &&
      val !== 0 &&
      val !== "" &&
      MAX_TURNS_REGEX.test(key)
    ) {
      return true;
    }
    if (val !== null && val !== undefined) {
      if (typeof val !== "object") {
        if (MAX_TURNS_REGEX.test(String(val))) {
          return true;
        }
      } else if (testObjectForMaxTurns(val, depth + 1, seen)) {
        return true;
      }
    }
  }

  return false;
}

export function isMaxTurnsError(err) {
  if (!err) {
    return false;
  }
  if (typeof err === "string") {
    return MAX_TURNS_REGEX.test(err);
  }
  if (typeof err === "object") {
    const directMessage = safeGet(err, "message");
    if (
      typeof directMessage === "string" &&
      MAX_TURNS_REGEX.test(directMessage)
    ) {
      return true;
    }
  }
  try {
    return testObjectForMaxTurns(err);
  } catch (extractErr) {
    void extractErr;
    return false;
  }
}

// Re-export turn accounting helpers for backward compatibility
export {
  bumpTurnCapacity,
  getMaxTurnsForModel,
  promptTurnBudgetExhaustion,
  requestHandoffSummary,
  syncSessionTools,
  TurnTracker,
};

function isQuotaError(err) {
  if (!err) {
    return false;
  }
  const message = String(err.message || err).toLowerCase();
  const name = String(err.name || "").toLowerCase();
  return (
    name.includes("quota") ||
    message.includes("quota") ||
    message.includes("resource_exhausted") ||
    message.includes("429") ||
    message.includes("rate limit") ||
    message.includes("capacity") ||
    message.includes("exhausted") ||
    message.includes("exceeded")
  );
}

function truncateDeep(val, maxLen = 2000, seen = new Set(), depth = 0) {
  if (depth > 20) {
    return "[Max Depth Exceeded]";
  }
  if (typeof val === "string") {
    if (val.length > maxLen) {
      return (
        val.substring(0, maxLen) +
        `... [Truncated, total length: ${val.length} characters]`
      );
    }
    return val;
  }
  if (Buffer.isBuffer(val)) {
    return "[Buffer]";
  }
  if (Array.isArray(val)) {
    if (seen.has(val)) {
      return "[Duplicate Reference]";
    }
    seen.add(val);
    const res = val
      .slice(0, 100)
      .map((item) => truncateDeep(item, maxLen, seen, depth + 1));
    if (val.length > 100) {
      res.push(`... [Truncated ${val.length - 100} remaining items]`);
    }
    seen.delete(val);
    return res;
  }
  if (typeof val === "object" && val !== null) {
    if (seen.has(val)) {
      return "[Duplicate Reference]";
    }
    seen.add(val);
    const res = {};
    if (val instanceof Error) {
      res.name = val.name;
      res.message =
        typeof val.message === "string"
          ? truncateDeep(val.message, maxLen, seen, depth + 1)
          : String(val.message);
      if (typeof val.stack === "string") {
        res.stack = truncateDeep(val.stack, maxLen, seen, depth + 1);
      }
    }
    const entries = Object.entries(val);
    for (const [k, v] of entries.slice(0, 100)) {
      res[k] = truncateDeep(v, maxLen, seen, depth + 1);
    }
    if (entries.length > 100) {
      res["..."] = `[Truncated ${entries.length - 100} remaining properties]`;
    }
    seen.delete(val);
    return res;
  }
  return val;
}

/**
 * Runs a Gemini Cli Agent session with specified parameters, automatically managing models,
 * fallbacks, budgeting, sandboxing, and error intercepts.
 *
 * @param {Object} options
 * @param {string} options.initialPrompt - The user objective or prompt
 * @param {string} [options.systemInstructions] - Custom system instructions for the agent
 * @param {string} [options.requestedModel] - Model to start with
 * @param {Array<string>} [options.blockTools] - List of tools to unregister/block
 * @param {Array<Object>} [options.customTools] - Array of custom SDK tools to register
 * @param {boolean} [options.isolate] - Whether to isolate the agent session
 * @param {boolean} [options.standalone] - Whether to auto-flush logs upon session completion
 * @param {number} [options.maxTurns] - Maximum turn capacity override for the session
 * @returns {Promise<string>} Accumulated text output from the agent
 */
export async function runAgentSession({
  initialPrompt,
  systemInstructions = "",
  requestedModel,
  blockTools = [],
  allowedTools = null,
  customTools = [],
  isolate = false,
  standalone = false,
  rl = null,
  ioOptions = {},
  maxTurns: inputMaxTurns,
  maxTurnsOverride,
  validateIntent = false,
  requireTopic = !isolate,
  streamIdleTimeoutMs = 180000,
}) {
  await initializeAgentRunner();

  try {
    const startingModel = requestedModel || MODEL_FLASH;
    const fallbackSequence = getModelFallbackSequence(startingModel);
    const effectiveMaxTurnsOverride = maxTurnsOverride ?? inputMaxTurns;
    const isNumeric =
      typeof effectiveMaxTurnsOverride === "number" ||
      (typeof effectiveMaxTurnsOverride === "string" &&
        effectiveMaxTurnsOverride.trim() !== "");
    const parsedTurnsOverride = isNumeric
      ? Number(effectiveMaxTurnsOverride)
      : NaN;
    const hasTurnOverride =
      Number.isInteger(parsedTurnsOverride) && parsedTurnsOverride > 0;

    if (
      effectiveMaxTurnsOverride !== undefined &&
      effectiveMaxTurnsOverride !== null &&
      !hasTurnOverride
    ) {
      console.warn(
        `⚠️ Invalid maxTurns override "${effectiveMaxTurnsOverride}". Expected a positive integer. Falling back to default model turn budget.`,
      );
    }

    if (requireTopic && isolate) {
      throw new Error(
        "Cannot enable requireTopic when isolate is true; isolated sessions do not have access to the update_topic tool.",
      );
    }

    if (requireTopic && blockTools.includes("update_topic")) {
      throw new Error(
        'Cannot block "update_topic" tool when requireTopic is enabled; topics are required for session progress and user visibility.',
      );
    }

    for (let i = 0; i < fallbackSequence.length; i++) {
      const currentModel = fallbackSequence[i];
      const maxTurns = hasTurnOverride
        ? parsedTurnsOverride
        : getMaxTurnsForModel(currentModel, {
          pro: MODEL_PRO,
          flash: MODEL_FLASH,
          flash_lite: MODEL_FLASH_LITE,
        });
      console.log(
        `\n[Initializing Gemini SDK Agentic Session] (Model: ${currentModel}, Max Turns: ${maxTurns})...`,
      );

      const topicInstructions = requireTopic
        ? `\n\n🎯 TOPIC & OBJECTIVE INSTRUCTIONS:
- You MUST call update_topic before executing any tools (on your first tool call) with a descriptive "title", "summary", and a granular "strategic_intent" (your concrete objective).
- Update the topic whenever you transition phases (e.g., Research -> Strategy -> Execution -> Validation) or whenever your strategic intent shifts.
- Keep your strategic_intent granular and actionable (e.g., "Empirically reproduce failure in parser.test.js", "Refactor buffer allocation in chunk.go").
- Do NOT output "Intent: <reason>" before each tool call.`
        : "";

      const intentInstructions = validateIntent
        ? `\n\n⚠️ TOOL USAGE INSTRUCTIONS: Immediately before calling a tool, you MUST output a single short sentence starting with "Intent:" explaining why you are calling it.`
        : "";

      let skillsContext = "";
      if (!isolate) {
        try {
          skillsContext = await loadWorkspaceSkillsContext(repoRoot);
        } catch {
          skillsContext = "";
        }
      }

      const hasSkillsAlready =
        typeof systemInstructions === "string" &&
        (systemInstructions.includes("<available_skills>") ||
          systemInstructions.includes("Skills-First Mandate"));

      const modelInstructions =
        (systemInstructions || "You are a highly capable agentic assistant.") +
        (hasSkillsAlready ? "" : skillsContext) +
        topicInstructions +
        intentInstructions +
        `\n\n⚠️ IMPORTANT TURN BUDGET: You are allowed a MAXIMUM of ${maxTurns} turns/iterations for this entire run. Conduct yourself efficiently, use tools in parallel, avoid unnecessary turns, and complete your task before reaching this limit.`;

      const noOpTool = tool(
        {
          name: "no_op",
          description: "A dummy tool that does nothing.",
          inputSchema: z.object({
            dummy: z.string().optional().describe("Ignored parameter"),
          }),
        },
        async () => "No operation performed.",
      );

      const updateTopicTool = tool(
        {
          name: "update_topic",
          description:
            "Manages your narrative flow and establishes your objective. Call this before executing tools and when transitioning phases.",
          inputSchema: z.object({
            strategic_intent: z
              .string()
              .optional()
              .describe("A statement of your immediate granular intent."),
            title: z
              .string()
              .optional()
              .describe("The title of the new topic or chapter."),
            summary: z
              .string()
              .optional()
              .describe("A summary of current progress and next steps."),
          }),
        },
        async ({ strategic_intent, title } = {}) => {
          return `Topic updated to: ${title || strategic_intent || "default"}`;
        },
      );

      const hasTopicTool = customTools.some(
        (t) => (t?.name || t?.definition?.name) === "update_topic",
      );
      const finalTools = isolate
        ? [noOpTool]
        : !hasTopicTool
          ? [...customTools, updateTopicTool]
          : customTools;
      const controller = new globalThis.AbortController();
      let accumulatedText = "";
      let turnTracker = null;

      try {
        const agent = new GeminiCliAgent({
          model: currentModel,
          max_turns: maxTurns,
          instructions: modelInstructions,
          tools: finalTools,
        });

        const session = agent.session();
        const rootDir = repoRoot || (await getRepoRoot());
        const repoName = path.basename(rootDir) || "workspace";
        const projectTempDir = path.join(os.homedir(), ".gemini/tmp", repoName);
        await fsPromises.mkdir(projectTempDir, { recursive: true });
        session.config.getWorkspaceContext().addDirectory(projectTempDir);
        await session.initialize();

        // Register GrepSearch tool alias if grep_search is available in registry
        const builtinGrep = session.config.toolRegistry.getTool("grep_search");
        if (builtinGrep && typeof builtinGrep.constructor === "function") {
          const GrepToolClass = builtinGrep.constructor;
          class GrepSearchTool extends GrepToolClass {
            constructor(config) {
              super(config);
              this.name = "GrepSearch";
              this.displayName = "GrepSearch";
            }
            getSchema(modelId) {
              const schema = super.getSchema ? super.getSchema(modelId) : {};
              return {
                ...schema,
                name: "GrepSearch",
                description:
                  "Searches for a regular expression pattern within file contents across the workspace. Alias for grep_search.",
              };
            }
          }
          try {
            session.config.toolRegistry.registerTool(
              new GrepSearchTool(session.config),
            );
            syncSessionTools(session);
          } catch (err) {
            console.debug(
              `[DEBUG] Failed to register GrepSearch alias: ${err.message}`,
            );
          }
        }

        const allowedSet = Array.isArray(allowedTools)
          ? new Set(allowedTools)
          : null;
        if (
          allowedSet &&
          (allowedSet.has("grep_search") || allowedSet.has("GrepSearch"))
        ) {
          allowedSet.add("grep_search");
          allowedSet.add("GrepSearch");
        }

        const blockedSet = new Set(blockTools || []);
        if (blockedSet.has("grep_search") || blockedSet.has("GrepSearch")) {
          blockedSet.add("grep_search");
          blockedSet.add("GrepSearch");
        }

        if (isolate) {
          const registeredTools = session.config.toolRegistry.getAllToolNames();
          for (const toolName of registeredTools) {
            if (
              toolName !== "no_op" &&
              (!allowedSet || !allowedSet.has(toolName))
            ) {
              session.config.toolRegistry.unregisterTool(toolName);
            }
          }
        } else {
          // Enforce custom sandboxing by unregistering specified tools
          const registeredTools = session.config.toolRegistry.getAllToolNames();
          for (const toolName of registeredTools) {
            if (
              allowedSet &&
              !allowedSet.has(toolName) &&
              toolName !== "no_op"
            ) {
              session.config.toolRegistry.unregisterTool(toolName);
            } else if (blockedSet.has(toolName)) {
              session.config.toolRegistry.unregisterTool(toolName);
            }
          }
          // Always unregister invoke_agent to prevent unmonitored delegation
          session.config.toolRegistry.unregisterTool("invoke_agent");
        }

        // Synchronize advertised session declarations with updated registry state
        syncSessionTools(session);

        let currentTurnText = "";
        let unrenderedProse = "";
        const pendingToolCalls = [];

        const flushProse = () => {
          if (unrenderedProse && unrenderedProse.trim() !== "") {
            console.log(renderProseCard(unrenderedProse));
            unrenderedProse = "";
          }
        };

        const renderToolOutput = (
          toolName,
          args,
          result,
          status = "success",
        ) => {
          flushProse();
          if (toolName === "update_topic" && status !== "error") {
            console.log(
              renderTopicBanner({
                topic: args?.title,
                objective: args?.strategic_intent,
                summary: args?.summary,
              }),
            );
          } else {
            console.log(
              renderToolCard({
                toolName,
                args,
                result,
                status,
              }),
            );
          }
        };

        turnTracker = new TurnTracker({
          maxTurns,
          controller,
          session,
          logger: writeLogAsync,
          validateIntent,
          requireTopic,
          getTurnText: () => currentTurnText,
          onToolExecution: ({ toolName, callId, args, result, error }) => {
            const argsStr = JSON.stringify(args || {});
            const pending =
              pendingToolCalls.find(
                (p) =>
                  !p.rendered &&
                  ((callId && p.callId && p.callId === callId) ||
                    (p.toolName === toolName &&
                      JSON.stringify(p.args || {}) === argsStr)),
              ) ||
              pendingToolCalls.find(
                (p) => !p.rendered && p.toolName === toolName,
              );

            if (pending) {
              pending.rendered = true;
            }
            renderToolOutput(
              toolName,
              args,
              result,
              error ? "error" : "success",
            );
          },
        });

        const runResult = await promptIdContext.run(session.id, async () => {
          let currentPrompt = initialPrompt;

          while (!controller.signal.aborted) {
            currentTurnText = "";
            let streamHaltedByBudget = false;
            let heartbeatInterval = null;
            let idleTimeoutTimer = null;
            let streamTimedOut = false;

            const resetIdleTimeout = () => {
              if (idleTimeoutTimer) {
                clearTimeout(idleTimeoutTimer);
                idleTimeoutTimer = null;
              }
              if (streamIdleTimeoutMs > 0 && !controller.signal.aborted) {
                idleTimeoutTimer = setTimeout(() => {
                  streamTimedOut = true;
                  const timeoutSec = Math.round(streamIdleTimeoutMs / 1000);
                  const errMsg = `[Stream Timeout] No response received from ${currentModel || "model"} for ${timeoutSec}s. Stream timed out.`;
                  writeLogAsync(errMsg);
                  console.warn(`\n⚠️ ${errMsg}`);
                  controller.abort(new Error(errMsg));
                }, streamIdleTimeoutMs);
                if (idleTimeoutTimer.unref) {
                  idleTimeoutTimer.unref();
                }
              }
            };

            try {
              const streamStartTime = Date.now();
              let firstChunkReceived = false;

              heartbeatInterval = setInterval(() => {
                if (!firstChunkReceived && !controller.signal.aborted) {
                  const elapsedSec = Math.round(
                    (Date.now() - streamStartTime) / 1000,
                  );
                  console.log(
                    `⏳ [Thinking] Waiting for ${currentModel || "model"} response (${elapsedSec}s elapsed)...`,
                  );
                }
              }, 15000);
              if (heartbeatInterval.unref) {
                heartbeatInterval.unref();
              }

              resetIdleTimeout();

              const stream = session.sendStream(
                currentPrompt,
                controller.signal,
              );

              for await (const chunk of stream) {
                if (streamTimedOut) {
                  const timeoutSec = Math.round(streamIdleTimeoutMs / 1000);
                  throw new Error(
                    `Stream idle timeout exceeded (${timeoutSec}s). No response from model.`,
                  );
                }

                resetIdleTimeout();

                if (!firstChunkReceived) {
                  firstChunkReceived = true;
                  if (heartbeatInterval) {
                    clearInterval(heartbeatInterval);
                    heartbeatInterval = null;
                  }
                }
                if (
                  chunk.type === "error" ||
                  chunk.type === "invalid_stream" ||
                  chunk.type === "agent_execution_blocked"
                ) {
                  if (isMaxTurnsError(chunk.value)) {
                    writeLogAsync(
                      `[Turn Budget Reached] SDK reported max turns reached via chunk error.`,
                    );
                    turnTracker.turnsExecuted = turnTracker.currentTurnBudget;
                    streamHaltedByBudget = true;
                    break;
                  }
                  throw new Error(
                    `Agent execution failed: ${chunk.type}. Details: ${JSON.stringify(chunk.value || "")}`,
                  );
                }
                if (chunk.type === "content") {
                  const shouldContinue = turnTracker.onContent();
                  if (!shouldContinue) {
                    streamHaltedByBudget = true;
                    break;
                  }
                  const text = chunk.value || "";
                  accumulatedText += text;
                  currentTurnText += text;
                  unrenderedProse += text;
                } else if (chunk.type === "tool_call_request") {
                  flushProse();
                  const toolCall = chunk.value;
                  const toolName = toolCall.name;
                  if (toolName === "invoke_agent") {
                    throw new Error("Sub-agent delegation is blocked.");
                  }

                  let args = toolCall.args;
                  if (typeof args === "string") {
                    try {
                      args = JSON.parse(args);
                    } catch (err) {
                      console.debug(
                        `[DEBUG] Failed to parse tool arguments string as JSON: ${err.message}`,
                      );
                    }
                  }

                  let formattedRawArgs;
                  if (typeof args === "object" && args !== null) {
                    formattedRawArgs = JSON.stringify(
                      truncateDeep(args, 2000),
                      null,
                      2,
                    );
                  } else {
                    formattedRawArgs = String(args);
                  }

                  writeLogAsync(
                    `\n[Tool Call]: ${toolName}\nArguments:\n${formattedRawArgs}`,
                  );

                  let shouldContinue = true;
                  if (!turnTracker.hasSessionHooks) {
                    shouldContinue = turnTracker.onToolCall(toolName);
                  } else {
                    shouldContinue =
                      turnTracker.count <= turnTracker.currentTurnBudget;
                  }
                  if (!shouldContinue) {
                    streamHaltedByBudget = true;
                    break;
                  }

                  const reason =
                    extractIntent(currentTurnText) ||
                    turnTracker.lastDeclaredIntent;
                  if (reason) {
                    turnTracker.lastDeclaredIntent = reason;
                  }

                  const callId = toolCall.callId || toolCall.id;
                  pendingToolCalls.push({
                    callId,
                    toolName,
                    args,
                    rendered: false,
                  });
                } else if (chunk.type === "tool_call_result") {
                  currentTurnText = "";
                  if (!turnTracker.hasSessionHooks) {
                    turnTracker.onToolResult();
                  }
                  try {
                    const safeValue = truncateDeep(chunk.value, 2000);
                    writeLogAsync(
                      `\n[Tool Result]: ${JSON.stringify(safeValue, null, 2)}`,
                    );
                  } catch (err) {
                    console.debug(
                      `[DEBUG] Failed to log tool call result: ${err.message}`,
                    );
                  }
                  const resultCallId =
                    chunk.callId ||
                    chunk.id ||
                    chunk.value?.callId ||
                    chunk.value?.id;
                  let pendingIndex = -1;
                  if (resultCallId) {
                    pendingIndex = pendingToolCalls.findIndex(
                      (p) => p.callId === resultCallId,
                    );
                  }
                  const pending =
                    pendingIndex !== -1
                      ? pendingToolCalls.splice(pendingIndex, 1)[0]
                      : pendingToolCalls.shift();

                  if (pending && !pending.rendered) {
                    pending.rendered = true;
                    const isError = Boolean(
                      chunk.value &&
                      typeof chunk.value === "object" &&
                      chunk.value.error === true,
                    );
                    const status = isError ? "error" : "success";
                    renderToolOutput(
                      pending.toolName,
                      pending.args,
                      chunk.value,
                      status,
                    );
                  }
                }
              }

              flushProse();
              if (streamTimedOut) {
                const timeoutSec = Math.round(streamIdleTimeoutMs / 1000);
                throw new Error(
                  `Stream idle timeout exceeded (${timeoutSec}s). No response from model.`,
                );
              }
              for (const pending of pendingToolCalls) {
                if (!pending.rendered) {
                  pending.rendered = true;
                  renderToolOutput(
                    pending.toolName,
                    pending.args,
                    null,
                    "error",
                  );
                }
              }
              pendingToolCalls.length = 0;
            } catch (err) {
              flushProse();
              if (streamTimedOut) {
                const timeoutSec = Math.round(streamIdleTimeoutMs / 1000);
                throw new Error(
                  `Stream idle timeout exceeded (${timeoutSec}s). No response from model.`,
                  {
                    cause: err,
                  },
                );
              }
              if (isMaxTurnsError(err)) {
                writeLogAsync(
                  `[Turn Budget Reached] SDK threw max turns error during stream: ${err.message}`,
                );
                turnTracker.turnsExecuted = turnTracker.currentTurnBudget;
                streamHaltedByBudget = true;
              } else {
                throw err;
              }
            } finally {
              if (heartbeatInterval) {
                clearInterval(heartbeatInterval);
              }
              if (idleTimeoutTimer) {
                clearTimeout(idleTimeoutTimer);
              }
            }

            if (streamHaltedByBudget || turnTracker.isBudgetExhausted) {
              writeLogAsync(
                `[Turn Limit Reached] Tracker detected budget limit (${turnTracker.count}/${turnTracker.currentTurnBudget}).`,
              );
              if (hasTurnOverride) {
                writeLogAsync(
                  `[Turn Limit Reached] Override limit enforced as a hard stop.`,
                );
                return accumulatedText;
              }

              const summary = await requestHandoffSummary(
                agent,
                session,
                writeLogAsync,
                controller,
              );
              if (summary) {
                accumulatedText += "\n\n" + summary;
              }

              const decision = await promptTurnBudgetExhaustion(
                turnTracker.count,
                turnTracker.currentTurnBudget,
                maxTurns,
                controller,
                { rl, ...ioOptions },
                writeLogAsync,
              );
              if (!decision.shouldContinue) {
                return accumulatedText;
              }

              turnTracker.currentTurnBudget = decision.newBudget;
              bumpTurnCapacity(agent, session, maxTurns);
              syncSessionTools(session);
              turnTracker.resetStepLatch();
              currentPrompt =
                "Please continue working to complete the task based on your handoff checkpoint summary.";
              continue;
            }

            break;
          }
          return accumulatedText;
        });

        if (controller.signal.aborted) {
          writeLogAsync(`[Session Stopped] Agent session aborted cleanly.`);
          return accumulatedText;
        }

        return runResult ?? accumulatedText;
      } catch (err) {
        if (isQuotaError(err) && i < fallbackSequence.length - 1) {
          console.warn(
            `⚠️ Model ${currentModel} hit quota limit. Retrying with lesser model ${fallbackSequence[i + 1]}...`,
          );
          continue;
        }
        throw err;
      } finally {
        turnTracker?.dispose();
      }
    }
    throw new Error("All models in fallback sequence failed.");
  } finally {
    if (standalone) {
      await flushLogs();
    }
  }
}
