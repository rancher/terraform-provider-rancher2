import { stdin as input, stdout as output } from "node:process";
import * as readline from "node:readline/promises";

/**
 * Returns the maximum turns allowed for a given model tier.
 *
 * @param {string} modelName - The name of the model
 * @param {Object} [models] - Optional dynamically loaded model configuration
 * @returns {number} The maximum number of turns
 */
export function getMaxTurnsForModel(modelName, models = null) {
  if (models) {
    if (modelName === models.pro) {
      return 10;
    }
    if (modelName === models.flash_lite) {
      return 2;
    }
    if (modelName === models.flash) {
      return 5;
    }
  }

  if (typeof modelName === "string") {
    if (modelName === "gemini-3.1-pro-preview" || modelName.includes("pro")) {
      return 10;
    }
    if (
      modelName === "gemini-3.1-flash-lite" ||
      modelName.includes("flash-lite")
    ) {
      return 2;
    }
    if (
      modelName === "gemini-3.8-flash" ||
      modelName === "gemini-3.5-flash" ||
      modelName.includes("flash")
    ) {
      return 5;
    }
  }

  return 5;
}

/**
 * Extracts declared intent from streamed turn text.
 * Returns empty string if no valid "Intent:" declaration is present.
 *
 * @param {string} text - The turn text to inspect
 * @returns {string} The normalized intent reason, or empty string
 */
export function extractIntent(text) {
  if (typeof text !== "string") {
    return "";
  }
  const trimmed = text.trim();
  if (!trimmed) {
    return "";
  }
  const match = trimmed.match(
    /(?:^|\n)\s*(?:(?:[-*>]|#{1,6})+\s+)?[_*]{0,2}Intent[_*]{0,2}:[_*]{0,2}\s*([\s\S]+)$/i,
  );
  if (!match || !match[1]) {
    return "";
  }
  const reason = match[1]
    .trim()
    .replace(/^[*_]+/, "")
    .replace(/[*_]+$/, "")
    .trim();
  return reason.replace(/^(to\s+)/i, "");
}

/**
 * Safely increments maxTurns across agent and session configurations.
 *
 * @param {Object} [agent] - Agent instance
 * @param {Object} [session] - Session instance
 * @param {number} delta - Number of turns to add
 */
export function bumpTurnCapacity(agent, session, delta) {
  if (agent) {
    if (typeof agent.maxTurns === "number") {
      agent.maxTurns += delta;
    }
    if (typeof agent.max_turns === "number") {
      agent.max_turns += delta;
    }
    if (agent.config && typeof agent.config.maxTurns === "number") {
      agent.config.maxTurns += delta;
    }
    if (agent.config && typeof agent.config.max_turns === "number") {
      agent.config.max_turns += delta;
    }
  }
  if (session && session.config) {
    if (typeof session.config.maxTurns === "number") {
      session.config.maxTurns += delta;
    }
    if (typeof session.config.max_turns === "number") {
      session.config.max_turns += delta;
    }
  }
}

/**
 * Prompts the user via readline when turn budget is exhausted to either continue or stop.
 * Defaults to extending budget when user presses Enter, or stopping on Ctrl+C/stop.
 *
 * @param {number} turnsExecuted - Total turns executed so far
 * @param {number} currentTurnBudget - Current turn budget limit
 * @param {number} maxTurns - The turn increment for this model
 * @param {AbortController} [controller] - The abort controller for the session
 * @param {Object} [ioOptions] - Custom I/O streams for testing
 * @param {Function} [logger] - Async logger function
 * @returns {Promise<{shouldContinue: boolean, newBudget: number}>}
 */
export async function promptTurnBudgetExhaustion(
  turnsExecuted,
  currentTurnBudget,
  maxTurns,
  controller,
  ioOptions = {},
  logger = () => {},
) {
  logger(
    `[Turn Limit Reached] Turns executed: ${turnsExecuted}/${currentTurnBudget}. Pausing for user decision.`,
  );
  console.log(
    `\n\n⚠️ Turn budget reached (${turnsExecuted}/${currentTurnBudget} turns).`,
  );

  const rlInput = ioOptions.input || input;
  const rlOutput = ioOptions.output || output;

  const isInteractive = Boolean(
    ioOptions.rl || ioOptions.input || rlInput.isTTY,
  );
  if (!isInteractive) {
    logger(
      `[DEBUG] Non-interactive environment detected before turn limit prompt. Defaulting to stop.`,
    );
    console.log(
      "\n🛑 Non-interactive environment detected. Gracefully aborting agent session.\n",
    );
    controller?.abort();
    return { shouldContinue: false, newBudget: currentTurnBudget };
  }

  const ownsRl = !ioOptions.rl;
  const rl =
    ioOptions.rl ||
    readline.createInterface({ input: rlInput, output: rlOutput });
  try {
    let answer = null;
    let attempts = 0;
    const isContinue = (val) => val === "" || val === "continue";
    const isStop = (val) => val === "stop" || val === "ctrl+c" || val === "^c";

    while (answer === null && attempts < 5) {
      attempts++;
      let response;
      try {
        response = await rl.question(
          `👉 Hit enter to continue with an additional ${maxTurns} turns, or type ctrl+c to stop the script: `,
        );
      } catch (err) {
        logger(`[DEBUG] Readline question error: ${err.message}`);
        answer = "stop";
        break;
      }
      if (response === undefined || response === null) {
        answer = "stop";
        break;
      }
      const trimmed = response.trim().toLowerCase();
      if (isContinue(trimmed)) {
        answer = "continue";
        break;
      }
      if (isStop(trimmed)) {
        answer = "stop";
        break;
      }

      if (!rlInput.isTTY && !ioOptions.input && !ioOptions.rl) {
        logger(
          `[DEBUG] Non-interactive environment detected at turn limit prompt. Defaulting to stop.`,
        );
        answer = "stop";
        break;
      }
      console.log(
        "Please hit enter to continue or type ctrl+c / 'stop' to stop.",
      );
    }

    if (answer !== "continue" && answer !== "stop") {
      answer = "stop";
    }

    if (answer === "continue") {
      logger(`[Turn Limit Reached] User opted to continue session.`);
      console.log(`\n▶️ Resuming / continuing agent workflow...\n`);
      return { shouldContinue: true, newBudget: currentTurnBudget + maxTurns };
    } else {
      logger(`[Turn Limit Aborted] User opted to stop session.`);
      console.log("\n🛑 Gracefully aborting agent session at user request.\n");
      if (controller) {
        controller.abort();
      }
      return { shouldContinue: false, newBudget: currentTurnBudget };
    }
  } finally {
    if (ownsRl) {
      rl.close();
    }
    // Do not call input.resume() to prevent leaving stdin in flowing mode without listeners
  }
}

/**
 * Ensures session client and chat tool configurations are synchronized with toolRegistry.
 * Resets cached model ID to prevent stale tool state on resume.
 *
 * @param {Object} session - The active agent session
 */
export function syncSessionTools(session) {
  if (!session) {
    return;
  }
  try {
    const client = session.client || session.config?.geminiClient;
    const registry = session.config?.toolRegistry;
    if (client) {
      client.lastUsedModelId = undefined;
      if (registry && typeof registry.getFunctionDeclarations === "function") {
        const modelId =
          typeof session.config?.getModel === "function"
            ? session.config.getModel()
            : undefined;
        const declarations = registry.getFunctionDeclarations(modelId);
        const tools =
          Array.isArray(declarations) && declarations.length > 0
            ? [{ functionDeclarations: declarations }]
            : [];
        if (typeof client.getChat === "function") {
          const chat = client.getChat();
          if (chat && typeof chat.setTools === "function") {
            chat.setTools(tools);
          }
        }
      }
    }
  } catch (err) {
    console.debug(
      `[DEBUG] syncSessionTools defensive sync encountered error: ${err.message}`,
    );
  }
}

/**
 * Bumps the agent session turn limit and requests a structured handoff summary.
 *
 * @param {Object} agent - The GeminiCliAgent instance
 * @param {Object} session - The active agent session instance
 * @param {Function} [logger] - Async logger hook
 * @param {AbortController} [controller] - The abort controller for the session
 * @returns {Promise<string>} The handoff summary produced by the agent
 */
export async function requestHandoffSummary(
  agent,
  session,
  logger = () => {},
  controller = null,
) {
  logger("[Turn Limit Reached] Requesting handoff summary from agent...");
  console.log(
    "\n⏳ Turn limit reached. Asking agent for a handoff checkpoint summary...\n",
  );

  // Increase turn capacity on the agent and session configurations to allow the summary prompt
  bumpTurnCapacity(agent, session, 2);

  const handoffPrompt =
    "⚠️ TURN LIMIT REACHED. Do not invoke any further tools. Output a structured markdown checkpoint summary with: 1) What has been completed so far, 2) Current state of modified files, and 3) The exact remaining steps to complete the task.";

  // Disarm tool execution non-destructively for the handoff turn so tool declarations stay intact
  const registry = session?.config?.toolRegistry;
  const restoreFns = [];
  if (registry?.getAllToolNames && registry?.getTool) {
    for (const name of registry.getAllToolNames()) {
      const tool = registry.getTool(name);
      if (tool && typeof tool.createInvocation === "function") {
        const origCreate =
          tool._originalCreateInvocation || tool.createInvocation;
        tool._originalCreateInvocation = origCreate;
        tool.createInvocation = function () {
          return {
            execute: async () => ({
              llmContent:
                "Tool execution is disabled during handoff checkpoint summary generation. Please output the requested markdown checkpoint summary.",
              returnDisplay: "Tool execution disabled during handoff summary.",
            }),
            getDescription: () => `Disarmed: ${name}`,
            getDisplayTitle: () => `Disarmed: ${name}`,
          };
        };
        restoreFns.push(() => {
          tool.createInvocation = origCreate;
          delete tool._originalCreateInvocation;
        });
      } else if (tool && typeof tool.execute === "function") {
        const origExec = tool._originalExecute || tool.execute;
        tool._originalExecute = origExec;
        tool.execute = async () => ({
          llmContent:
            "Tool execution is disabled during handoff checkpoint summary generation. Please output the requested markdown checkpoint summary.",
          returnDisplay: "Tool execution disabled during handoff summary.",
        });
        restoreFns.push(() => {
          tool.execute = origExec;
          delete tool._originalExecute;
        });
      }
    }
  }

  let summary = "";
  try {
    const handoffStream = session.sendStream(handoffPrompt, controller?.signal);
    for await (const chunk of handoffStream) {
      if (chunk.type === "content") {
        const text = chunk.value || "";
        process.stdout.write(text);
        summary += text;
      }
    }
  } catch (err) {
    logger(`[DEBUG] Failed to obtain handoff summary: ${err.message}`);
  } finally {
    for (const restore of restoreFns) {
      try {
        restore();
      } catch (err) {
        console.debug(
          `[DEBUG] Failed to restore disarmed tool: ${err.message}`,
        );
      }
    }
    syncSessionTools(session);
  }

  return summary.trim();
}

/**
 * State tracker for agent stream turns and budget exhaustion.
 * Supports native session event hooks as well as stream-level fallback inspection.
 */
export class TurnTracker {
  /**
   * @param {Object} options
   * @param {number} options.maxTurns - Increment turn budget per model
   * @param {AbortController} options.controller - Controller used to abort execution
   * @param {Object} [options.session] - The active GeminiCliAgent session
   * @param {Function} [options.logger] - Logger hook for recording state changes
   * @param {boolean} [options.validateIntent=false] - Whether to enforce intent declaration before tools
   * @param {boolean} [options.requireTopic=false] - Whether to enforce calling update_topic with granular objective before executing tools on tool-using turns
   * @param {Function} [options.onToolExecution=null] - Hook called upon tool completion with execution details
   * @param {string[]} [options.exemptTools=['no_op']] - Tool names exempt from intent and topic enforcement
   * @param {Function} [options.getTurnText=null] - Accessor for currently streamed turn text
   * @param {number} [options.maxIntentRetries=2] - Maximum free intent-correction retries per turn
   */
  constructor({
    maxTurns,
    controller,
    session = null,
    logger = () => {},
    validateIntent = false,
    requireTopic = false,
    onToolExecution = null,
    exemptTools = ["no_op"],
    getTurnText = null,
    maxIntentRetries = 2,
  }) {
    this.maxTurns = maxTurns;
    this.currentTurnBudget = maxTurns;
    this.turnsExecuted = 0;
    this.turnCountedForModelStep = false;
    this.controller = controller;
    this.session = session;
    this.logger = logger;
    this.validateIntent = validateIntent;
    this.requireTopic = requireTopic;
    this.onToolExecution = onToolExecution;
    this.exemptTools = exemptTools;
    this.getTurnText = getTurnText;
    this.maxIntentRetries = maxIntentRetries;
    this.intentRetries = 0;
    this.lastDeclaredIntent = "";
    this.hasEstablishedTopic = false;
    this.activeTopic = null;
    this._hasSessionEvents = false;
    this._cleanupFns = [];

    if (this.session?.events && typeof this.session.events.on === "function") {
      this._attachSessionHooks(this.session.events);
    }
    if (this.session?.config?.toolRegistry) {
      this._attachToolRegistryHooks(this.session.config.toolRegistry);
    }
  }

  /**
   * Connects to native SDK session event emitter.
   * @private
   */
  _attachSessionHooks(events) {
    this._hasSessionEvents = true;
    const onToolCall = (call) => {
      const name = call?.name || "unknown_tool";
      this.onToolCall(name);
    };
    const onToolResult = () => {
      this.onToolResult();
    };

    events.on("tool_call", onToolCall);
    events.on("tool_result", onToolResult);

    this._cleanupFns.push(() => {
      events.removeListener?.("tool_call", onToolCall);
      events.removeListener?.("tool_result", onToolResult);
      this._hasSessionEvents = false;
    });
  }

  /**
   * Wraps tool executions in the SDK tool registry to reliably reset the step latch upon completion.
   * @private
   */
  _attachToolRegistryHooks(registry) {
    if (!registry || typeof registry.getAllToolNames !== "function") {
      return;
    }
    const wrapTool = (tool) => {
      const toolName = tool?.name || tool?.definition?.name || "tool";
      const tracker = this;

      const createWrappedExecute = (
        origExec,
        invocationArgs = null,
        invocation = null,
        capturedCallId = null,
      ) => {
        return async function (...execArgs) {
          const callId =
            capturedCallId ||
            execArgs[0]?.callId ||
            execArgs[0]?.id ||
            execArgs[1]?.callId ||
            execArgs[1]?.id ||
            invocation?.callId ||
            invocation?.id;
          try {
            let pendingTopic = null;
            const toolArgs =
              invocationArgs !== null &&
              invocationArgs !== undefined &&
              typeof invocationArgs === "object" &&
              Object.keys(invocationArgs).length > 0
                ? invocationArgs
                : invocation?.params ||
                  invocation?.input ||
                  (execArgs[0] && typeof execArgs[0] === "object"
                    ? execArgs[0]
                    : {});

            if (toolName === "update_topic") {
              const title = toolArgs.title || "";
              const objective =
                toolArgs.strategic_intent ||
                toolArgs.objective ||
                toolArgs.summary ||
                "";
              const summary = toolArgs.summary || "";
              if (
                tracker.requireTopic &&
                (!objective ||
                  typeof objective !== "string" ||
                  !objective.trim())
              ) {
                tracker.refundTurn();
                return {
                  llmContent: `Topic update rejected: You must provide a granular "strategic_intent" describing your specific objective (e.g., 'Empirically reproduce test failure in parser.test.js').`,
                  returnDisplay: `Topic update rejected: Missing granular strategic_intent.`,
                  error: true,
                };
              }
              pendingTopic = {
                title,
                objective: objective.trim(),
                summary: summary.trim(),
              };
            } else if (
              tracker.requireTopic &&
              !tracker.hasEstablishedTopic &&
              !tracker.exemptTools.includes(toolName)
            ) {
              tracker.refundTurn();
              return {
                llmContent: `Tool call deferred: You must establish an objective first by calling update_topic with a "title", "summary", and a granular "strategic_intent" before executing ${toolName}.`,
                returnDisplay: `Tool call deferred: Please call update_topic first.`,
                error: true,
              };
            }

            if (
              tracker.validateIntent &&
              !tracker.exemptTools.includes(toolName)
            ) {
              const turnText = tracker.getTurnText ? tracker.getTurnText() : "";
              const declared = extractIntent(turnText);
              if (declared) {
                tracker.lastDeclaredIntent = declared;
              }
              const intent = declared || tracker.lastDeclaredIntent;
              if (!intent) {
                tracker.refundTurn();
                return {
                  llmContent: `Tool call refused: You must declare an "Intent: <reason>" immediately before calling ${toolName}. Please output a single sentence starting with "Intent:" explaining why you are calling this tool, then retry.`,
                  returnDisplay: `Tool call refused: Missing Intent declaration for ${toolName}.`,
                  error: true,
                };
              }
            }
            tracker.intentRetries = 0;
            let result;
            let executionError = false;
            try {
              result = await origExec.apply(this, execArgs);
              if (
                pendingTopic &&
                !(result && typeof result === "object" && result.error === true)
              ) {
                tracker.hasEstablishedTopic = true;
                tracker.activeTopic = pendingTopic;
              }
            } catch (execErr) {
              result = execErr.message;
              executionError = true;
              throw execErr;
            } finally {
              if (typeof tracker.onToolExecution === "function") {
                try {
                  tracker.onToolExecution({
                    toolName,
                    callId,
                    args: toolArgs,
                    result,
                    error:
                      executionError ||
                      (result &&
                        typeof result === "object" &&
                        result.error === true),
                  });
                } catch (e) {
                  console.debug(
                    `[DEBUG] onToolExecution hook failed: ${e.message}`,
                  );
                }
              }
            }
            return result;
          } finally {
            tracker.onToolResult();
          }
        };
      };

      if (
        tool &&
        tool.definition &&
        typeof tool.definition.action === "function" &&
        !tool.definition._turnAccountingWrapped
      ) {
        const origAction = tool.definition.action;
        tool.definition.action = createWrappedExecute(
          origAction,
          null,
          tool.definition,
        );
        tool.definition._turnAccountingWrapped = true;
        this._cleanupFns.push(() => {
          tool.definition.action = origAction;
          delete tool.definition._turnAccountingWrapped;
        });
      } else if (
        tool &&
        typeof tool.action === "function" &&
        !tool._turnAccountingWrapped
      ) {
        const origAction = tool.action;
        tool.action = createWrappedExecute(origAction, null, tool);
        tool._turnAccountingWrapped = true;
        this._cleanupFns.push(() => {
          tool.action = origAction;
          delete tool._turnAccountingWrapped;
        });
      } else if (
        tool &&
        typeof tool.createInvocation === "function" &&
        !tool._turnAccountingWrapped
      ) {
        const origCreate = tool.createInvocation;
        tool.createInvocation = function (...createArgs) {
          const invocation = origCreate.apply(this, createArgs);
          const invocationCallId =
            createArgs[1]?.callId ||
            createArgs[1]?.id ||
            createArgs[0]?.callId ||
            createArgs[0]?.id ||
            invocation?.callId ||
            invocation?.id;
          if (invocation && typeof invocation.execute === "function") {
            invocation.execute = createWrappedExecute(
              invocation.execute,
              createArgs[0],
              invocation,
              invocationCallId,
            );
          }
          return invocation;
        };
        tool._turnAccountingWrapped = true;
        this._cleanupFns.push(() => {
          tool.createInvocation = origCreate;
          delete tool._turnAccountingWrapped;
        });
      } else if (
        tool &&
        typeof tool.execute === "function" &&
        !tool._turnAccountingWrapped
      ) {
        const origExec = tool.execute;
        tool.execute = createWrappedExecute(origExec, null, tool);
        tool._turnAccountingWrapped = true;
        this._cleanupFns.push(() => {
          tool.execute = origExec;
          delete tool._turnAccountingWrapped;
        });
      }
    };

    for (const name of registry.getAllToolNames()) {
      wrapTool(registry.getTool(name));
    }

    if (Array.isArray(this.session?.tools)) {
      for (const t of this.session.tools) {
        wrapTool(t);
      }
    }

    if (typeof registry.registerTool === "function") {
      const origRegister = registry.registerTool;
      registry.registerTool = function (...args) {
        for (const arg of args) {
          wrapTool(arg);
        }
        return origRegister.apply(this, args);
      };
      this._cleanupFns.push(() => {
        registry.registerTool = origRegister;
      });
    }
  }

  /**
   * Detaches native event hooks to prevent listener leaks.
   */
  dispose() {
    for (const cleanup of this._cleanupFns) {
      try {
        cleanup();
      } catch (err) {
        this.logger(
          `[DEBUG] Error during turn tracker cleanup: ${err.message}`,
        );
      }
    }
    this._cleanupFns = [];
  }

  /**
   * Whether native session event hooks are currently active.
   */
  get hasSessionHooks() {
    return Boolean(this._hasSessionEvents);
  }

  /**
   * Refunds the turn count for an intercepted/refused tool call.
   * Decrements turnsExecuted if within retry limit.
   * @returns {boolean} Whether the refund was granted
   */
  refundTurn() {
    this.lastDeclaredIntent = "";
    if (this.turnsExecuted > 0 && this.turnCountedForModelStep) {
      if (this.intentRetries < this.maxIntentRetries) {
        this.intentRetries++;
        this.turnsExecuted--;
        this.logger(
          `[Turn Refunded] Tool call refused due to missing intent (retry ${this.intentRetries}/${this.maxIntentRetries}). Turns executed refunded to: ${this.turnsExecuted}/${this.currentTurnBudget}`,
        );
        return true;
      }
      this.logger(
        `[Turn Budget] Intent retries exceeded maximum limit (${this.maxIntentRetries}). Turn not refunded.`,
      );
    }
    return false;
  }

  /**
   * Resets the turn step latch upon receiving a tool result.
   */
  onToolResult() {
    this.turnCountedForModelStep = false;
  }

  /**
   * Handles content streaming tokens.
   * @returns {boolean} true if execution should continue, false if aborted.
   */
  onContent() {
    if (!this.turnCountedForModelStep) {
      this.turnCountedForModelStep = true;
      this.lastDeclaredIntent = "";
      this.turnsExecuted++;
      this.logger(
        `[Turn ${this.turnsExecuted}/${this.currentTurnBudget}] Agent streaming content response.`,
      );
    }
    return this._checkBudget();
  }

  /**
   * Handles a tool call request chunk.
   * @param {string} toolName - Name of the requested tool
   * @returns {boolean} true if execution should continue, false if aborted.
   */
  onToolCall(toolName) {
    if (!this.turnCountedForModelStep) {
      this.turnCountedForModelStep = true;
      this.lastDeclaredIntent = "";
      this.turnsExecuted++;
      this.logger(
        `[Turn ${this.turnsExecuted}/${this.currentTurnBudget}] Tool call request: ${toolName}`,
      );
    }
    return this._checkBudget();
  }

  /**
   * Evaluates turn budget and allows content streaming through the current turn.
   * @private
   * @returns {boolean} false if budget is exceeded, true otherwise
   */
  _checkBudget() {
    return this.turnsExecuted <= this.currentTurnBudget;
  }

  /**
   * Explicitly resets the step latch across stream prompt restarts.
   */
  resetStepLatch() {
    this.turnCountedForModelStep = false;
    this.lastDeclaredIntent = "";
  }

  /**
   * Whether the budget limit has been reached.
   */
  get isBudgetExhausted() {
    return this.turnsExecuted > this.currentTurnBudget;
  }

  /**
   * Current summary count of turns taken.
   */
  get count() {
    return this.turnsExecuted;
  }
}
