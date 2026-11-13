import toml from "@iarna/toml";
import { execFile } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

let cachedRepoRoot = null;
let cachedSkillsContext = null;
let cachedSkillsFingerprint = "";
let lastSkillsFingerprintTime = 0;
const SKILLS_FINGERPRINT_TTL_MS = 2000;

/**
 * Resets the cached skills context (useful for testing).
 */
export function resetSkillsContextCache() {
  cachedSkillsContext = null;
  cachedSkillsFingerprint = "";
  lastSkillsFingerprintTime = 0;
}

/**
 * Recursively calculates a fingerprint across all files in a skills directory.
 * Safely guards against circular symlink recursion and avoids call stack overflow.
 * Captures file paths, file mtimes, and sizes so additions, deletions, renames,
 * and modifications of files invalidate the cache. Empty-directory additions
 * and removals do not alter the fingerprint because directory entries themselves
 * are not hashed.
 * @param {string} dir
 * @param {string} [baseDir=dir]
 * @param {Set<string>} [visited=new Set()]
 * @returns {Promise<string>}
 */
export async function getSkillsFingerprint(
  dir,
  baseDir = dir,
  visited = new Set(),
) {
  const parts = [];
  try {
    const realDir = await fs.realpath(dir).catch(() => path.resolve(dir));
    if (visited.has(realDir)) {
      return "";
    }
    visited.add(realDir);

    const entries = await fs.readdir(dir, { withFileTypes: true });
    entries.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));

    const results = await Promise.all(
      entries.map(async (entry) => {
        if (entry.isSymbolicLink()) {
          return null;
        }
        const fullPath = path.join(dir, entry.name);
        const relPath = path.relative(baseDir, fullPath);
        try {
          if (entry.isDirectory()) {
            return await getSkillsFingerprint(fullPath, baseDir, visited);
          } else if (entry.isFile()) {
            const stats = await fs.stat(fullPath);
            return `f:${relPath}:${stats.mtimeMs}:${stats.size}`;
          }
        } catch {
          // ignore missing/unreadable entry
        }
        return null;
      }),
    );

    for (const res of results) {
      if (res) {
        parts.push(res);
      }
    }
  } catch {
    // ignore readdir error
  }
  return parts.join("|");
}

/**
 * Loads and formats workspace skills from .gemini/skills/ and .gemini/system.md.
 * Automatically checks directory and file mtimes to invalidate stale cache.
 * @param {string} [customRepoRoot=null]
 * @returns {Promise<string>} Formatted markdown block of available skills and mandate
 */
export async function loadWorkspaceSkillsContext(customRepoRoot = null) {
  if (
    cachedSkillsContext &&
    !customRepoRoot &&
    Date.now() - lastSkillsFingerprintTime < SKILLS_FINGERPRINT_TTL_MS
  ) {
    return cachedSkillsContext;
  }

  const root = customRepoRoot || (await getRepoRoot());
  const skillsDir = path.join(root, ".gemini/skills");
  const systemMdPath = path.join(root, ".gemini/system.md");

  let currentFingerprint = "";
  let hasSystemMd = false;

  if (!customRepoRoot) {
    try {
      currentFingerprint = await getSkillsFingerprint(skillsDir);
    } catch (err) {
      console.debug(
        "[DEBUG] Failed to get skills fingerprint:",
        err?.message || String(err),
      );
    }
  }

  try {
    const stats = await fs.stat(systemMdPath);
    if (!customRepoRoot) {
      currentFingerprint += `|system.md:${stats.mtimeMs}:${stats.size}`;
    }
    hasSystemMd = stats.isFile();
  } catch {
    // ignore missing system.md
  }

  if (!customRepoRoot) {
    lastSkillsFingerprintTime = Date.now();
  }

  if (
    cachedSkillsContext &&
    !customRepoRoot &&
    currentFingerprint === cachedSkillsFingerprint
  ) {
    return cachedSkillsContext;
  }

  let skillsBlock = "";
  try {
    const entries = await fs.readdir(skillsDir, { withFileTypes: true });
    const skillDirs = entries
      .filter((e) => e.isDirectory())
      .map((e) => e.name)
      .sort();

    const skills = (
      await Promise.all(
        skillDirs.map(async (dirName) => {
          const skillMdPath = path.join(skillsDir, dirName, "SKILL.md");
          if (await exists(skillMdPath)) {
            try {
              const content = await fs.readFile(skillMdPath, "utf8");
              const nameMatch = content.match(
                /(?:^name:\s*|<name>)(.+?)(?:<\/name>|$)/m,
              );
              const descMatch = content.match(
                /(?:^description:\s*(?:[>|]\s*\n)?|<description>)([\s\S]+?)(?:<\/description>|\n---|\n[a-zA-Z0-9_-]+:\s*|$)/m,
              );
              const name = nameMatch ? nameMatch[1].trim() : dirName;
              const description = descMatch ? descMatch[1].trim() : "";

              // Discover scripts
              const scriptsDir = path.join(skillsDir, dirName, "scripts");
              let scriptFiles = [];
              if (await exists(scriptsDir)) {
                try {
                  const scriptEntries = await fs.readdir(scriptsDir, {
                    withFileTypes: true,
                  });
                  scriptFiles = scriptEntries
                    .filter(
                      (e) =>
                        e.isFile() &&
                        (e.name.endsWith(".js") ||
                          e.name.endsWith(".mjs") ||
                          e.name.endsWith(".cjs") ||
                          e.name.endsWith(".sh") ||
                          e.name.endsWith(".py")),
                    )
                    .map((e) => e.name)
                    .sort()
                    .map((f) => {
                      if (f.endsWith(".sh")) {
                        return `bash .gemini/skills/${dirName}/scripts/${f}`;
                      }
                      if (f.endsWith(".py")) {
                        return `python3 .gemini/skills/${dirName}/scripts/${f}`;
                      }
                      return `node .gemini/skills/${dirName}/scripts/${f}`;
                    });
                } catch {
                  scriptFiles = [];
                }
              }

              return {
                name,
                description,
                location: `.gemini/skills/${dirName}/SKILL.md`,
                scripts: scriptFiles,
              };
            } catch {
              // ignore error reading single skill
            }
          }
          return null;
        }),
      )
    ).filter(Boolean);

    if (skills.length > 0) {
      skillsBlock =
        `<available_skills>\n` +
        skills
          .map(
            (s) =>
              `  <skill>\n    <name>${s.name}</name>\n    <description>${s.description}</description>\n    <location>${s.location}</location>\n    <scripts>\n${s.scripts.map((sc) => `      - ${sc}`).join("\n")}\n    </scripts>\n  </skill>`,
          )
          .join("\n") +
        `\n</available_skills>`;
    }
  } catch {
    skillsBlock = "";
  }

  // Load Skills-First Mandate from .gemini/system.md if present
  let mandateText = "";
  if (hasSystemMd) {
    try {
      const systemContent = await fs.readFile(systemMdPath, "utf8");
      const mandateMatch = systemContent.match(
        /## (?:[^\n]*?)Skills-First Mandate[\s\S]*?(?=\n---|\n## |$)/,
      );
      if (mandateMatch) {
        mandateText = mandateMatch[0].trim();
      }
    } catch {
      mandateText = "";
    }
  }

  if (!mandateText) {
    console.warn(
      "⚠️ WARNING: .gemini/system.md is missing or lacks the Skills-First Mandate block. Using hardcoded fallback which may drift out of sync.",
    );
    mandateText = `## 🛠️ Skills-First Mandate\nTo guarantee safety and compliance with repository guardrails, you MUST prioritize the installed skills in .gemini/skills/ (such as git, github-ci, and github-pr) for repository operations instead of raw mutating commands.`;
  }

  const result = `\n\n# Available Workspace Skills\n\n${skillsBlock}\n\n${mandateText}\n`;
  if (!customRepoRoot) {
    cachedSkillsContext = result;
    cachedSkillsFingerprint = currentFingerprint;
    lastSkillsFingerprintTime = Date.now();
  }
  return result;
}

/**
 * Gets the repository root path dynamically.
 * @returns {Promise<string>}
 */
export async function getRepoRoot() {
  if (cachedRepoRoot) {
    return cachedRepoRoot;
  }
  try {
    const { stdout } = await execFileAsync("git", [
      "rev-parse",
      "--show-toplevel",
    ]);
    cachedRepoRoot = stdout.trim();
    return cachedRepoRoot;
  } catch (err) {
    console.debug(
      `[DEBUG] Failed to resolve git repo root via execFile: ${err?.message || err}. Falling back to cwd.`,
    );
    cachedRepoRoot = process.cwd();
    return cachedRepoRoot;
  }
}

/**
 * Checks if a file exists.
 * @param {string} filePath
 * @returns {Promise<boolean>}
 */
export async function exists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

/**
 * Loads custom agent instructions from .gemini/agents/<agentName>.toml
 * @param {string} agentName
 * @returns {Promise<Object>}
 */
export async function loadAgentInstructions(agentName) {
  const repoRoot = await getRepoRoot();
  const filePath = path.join(repoRoot, `.gemini/agents/${agentName}.toml`);
  try {
    const content = await fs.readFile(filePath, "utf8");
    const parsed = toml.parse(content);
    return {
      instructions: parsed.instructions || "",
      model: parsed.model,
      max_turns: parsed.max_turns,
      temperature: parsed.temperature,
    };
  } catch (err) {
    console.warn(
      `⚠️ Could not load custom instructions for agent ${agentName}: ${err?.message || err}`,
    );
    return { instructions: "" };
  }
}

/**
 * Extracts balanced top-level JSON objects from mixed text.
 * @param {string} str
 * @returns {Array<Object>}
 */
function extractTopLevelJSONObjects(str) {
  const objects = [];
  const stack = [];
  let startIndex = -1;
  let inString = false;
  let escape = false;

  for (let i = 0; i < str.length; i++) {
    const char = str[i];
    if (inString) {
      if (escape) {
        escape = false;
      } else if (char === "\\") {
        escape = true;
      } else if (char === '"') {
        inString = false;
      } else if (char === "\n") {
        inString = false;
        escape = false;
      }
      continue;
    }

    if (char === '"') {
      inString = true;
    } else if (char === "{" || char === "[") {
      if (stack.length === 0) {
        startIndex = i;
      }
      stack.push(char);
    } else if (char === "}" || char === "]") {
      if (stack.length > 0) {
        const top = stack[stack.length - 1];
        if ((top === "{" && char === "}") || (top === "[" && char === "]")) {
          stack.pop();
          if (stack.length === 0 && startIndex !== -1) {
            const candidate = str.slice(startIndex, i + 1);
            try {
              const parsed = JSON.parse(candidate);
              if (parsed && typeof parsed === "object") {
                objects.push(parsed);
              }
            } catch (err) {
              console.debug(
                `[DEBUG] Skipping malformed JSON candidate: ${err?.message || err}`,
              );
            }
            startIndex = -1;
          }
        } else {
          startIndex = -1;
          stack.length = 0;
        }
      }
    }
  }
  return objects;
}

/**
 * Normalizes any parsed JSON object or array into a standard implementation plan structure.
 * @param {Object|Array} raw
 * @param {string} [fallbackObjective]
 * @returns {Object}
 */
export function normalizePlanObject(raw, fallbackObjective = "") {
  const obj = Array.isArray(raw) ? { plan: raw } : { ...raw };

  let title =
    typeof obj.title === "string" && obj.title.trim() ? obj.title.trim() : "";
  let objective =
    typeof obj.objective === "string" && obj.objective.trim()
      ? obj.objective.trim()
      : "";
  const description =
    typeof obj.description === "string" && obj.description.trim()
      ? obj.description.trim()
      : "";
  const name =
    typeof obj.name === "string" && obj.name.trim() ? obj.name.trim() : "";

  let firstPlanItemDesc = "";
  if (
    Array.isArray(obj.plan) &&
    obj.plan.length > 0 &&
    typeof obj.plan[0] === "object" &&
    obj.plan[0]
  ) {
    const p0 = obj.plan[0];
    const filePrefix = p0.file ? `[${p0.file}] ` : "";
    if (p0.description) {
      firstPlanItemDesc = p0.description;
    } else if (p0.title) {
      firstPlanItemDesc = p0.title;
    } else if (p0.task) {
      firstPlanItemDesc = `${filePrefix}${p0.task}`;
    } else if (Array.isArray(p0.steps) && p0.steps.length > 0) {
      firstPlanItemDesc = `${filePrefix}${p0.steps[0]}`;
    }
  }

  if (!title) {
    title =
      objective ||
      name ||
      description ||
      firstPlanItemDesc ||
      fallbackObjective ||
      "Implementation Plan";
  }
  if (!objective) {
    objective =
      fallbackObjective ||
      description ||
      title ||
      firstPlanItemDesc ||
      "Implementation Plan";
  }

  let implementationTasks = [];
  if (Array.isArray(obj.implementation_tasks)) {
    implementationTasks = obj.implementation_tasks.map((t) =>
      typeof t === "object" && t
        ? t.description || t.task || JSON.stringify(t)
        : String(t),
    );
  } else if (Array.isArray(obj.tasks)) {
    implementationTasks = obj.tasks.map((t) =>
      typeof t === "object" && t
        ? t.description || t.task || JSON.stringify(t)
        : String(t),
    );
  } else if (Array.isArray(obj.steps)) {
    implementationTasks = obj.steps.map((s) =>
      typeof s === "object" && s
        ? s.description || s.step || JSON.stringify(s)
        : String(s),
    );
  } else if (Array.isArray(obj.plan)) {
    for (const item of obj.plan) {
      if (typeof item === "string") {
        implementationTasks.push(item);
      } else if (item && typeof item === "object") {
        const filePrefix = item.file ? `[${item.file}] ` : "";
        if (Array.isArray(item.steps) && item.steps.length > 0) {
          for (const step of item.steps) {
            implementationTasks.push(`${filePrefix}${step}`);
          }
        } else if (item.description) {
          implementationTasks.push(`${filePrefix}${item.description}`);
        } else if (item.task) {
          implementationTasks.push(`${filePrefix}${item.task}`);
        } else if (item.title) {
          implementationTasks.push(`${filePrefix}${item.title}`);
        }
      }
    }
  }

  if (implementationTasks.length === 0) {
    if (typeof obj.plan === "string" && obj.plan.trim()) {
      implementationTasks.push(obj.plan.trim());
    } else {
      implementationTasks.push("Execute implementation according to plan.");
    }
  }

  let inScope = [];
  let outOfScope = [];
  if (obj.scope_boundaries && typeof obj.scope_boundaries === "object") {
    if (Array.isArray(obj.scope_boundaries.in_scope)) {
      inScope = obj.scope_boundaries.in_scope.map(String);
    }
    if (Array.isArray(obj.scope_boundaries.out_of_scope)) {
      outOfScope = obj.scope_boundaries.out_of_scope.map(String);
    }
  }
  if (inScope.length === 0) {
    if (Array.isArray(obj.in_scope)) {
      inScope = obj.in_scope.map(String);
    } else if (Array.isArray(obj.plan)) {
      const files = obj.plan
        .map((p) => (p && typeof p === "object" ? p.file : null))
        .filter(Boolean);
      if (files.length > 0) {
        inScope = files.map((f) => `Modify ${f}`);
      }
    }
  }
  if (inScope.length === 0) {
    inScope = ["Implement requested tasks as planned"];
  }

  if (outOfScope.length === 0) {
    if (Array.isArray(obj.out_of_scope)) {
      outOfScope = obj.out_of_scope.map(String);
    } else {
      outOfScope = [
        "Unrelated refactoring or changes outside the defined scope",
      ];
    }
  }

  let exitCriteria;
  if (Array.isArray(obj.exit_criteria)) {
    exitCriteria = obj.exit_criteria.map(String);
  } else if (
    typeof obj.exit_criteria === "string" &&
    obj.exit_criteria.trim()
  ) {
    exitCriteria = [obj.exit_criteria.trim()];
  } else {
    exitCriteria = [
      "All implementation tasks completed successfully.",
      "Tests and linters pass cleanly.",
    ];
  }

  return {
    title,
    objective,
    scope_boundaries: {
      in_scope: inScope,
      out_of_scope: outOfScope,
    },
    exit_criteria: exitCriteria,
    implementation_tasks: implementationTasks,
  };
}

/**
 * Parses JSON blocks out of raw text.
 * @param {string} text
 * @param {string} [fallbackType]
 * @returns {Object|null}
 */
export function parseJSONFromText(text, fallbackType = "qa") {
  const safeText = typeof text === "string" ? text : "";
  const normalizeResult = (obj) => {
    if (!obj || typeof obj !== "object") {
      return obj;
    }
    if (fallbackType === "qa") {
      if (Array.isArray(obj.findings) && !obj.approval_status) {
        obj.approval_status =
          obj.findings.length === 0 ? "APPROVED" : "UNAPPROVED";
      }
    } else if (fallbackType === "plan") {
      return normalizePlanObject(obj);
    }
    return obj;
  };

  const matches = [
    ...safeText.matchAll(/```(?:[a-zA-Z0-9_-]+)?\s*?\n?([\s\S]*?)\n?\s*```/gi),
  ];
  for (let i = matches.length - 1; i >= 0; i--) {
    const candidate = matches[i][1].trim();
    try {
      const parsed = JSON.parse(candidate);
      if (parsed && typeof parsed === "object") {
        return normalizeResult(parsed);
      }
    } catch (err) {
      console.debug(
        `[DEBUG] Code block candidate did not contain valid JSON: ${err?.message || err}`,
      );
    }
  }
  try {
    const clean = safeText.trim();
    if (clean) {
      const parsed = JSON.parse(clean);
      if (parsed && typeof parsed === "object") {
        return normalizeResult(parsed);
      }
    }
  } catch (err) {
    console.debug(`[DEBUG] Raw text was not pure JSON: ${err?.message || err}`);
  }

  const unadorned = extractTopLevelJSONObjects(safeText);
  if (unadorned.length > 0) {
    if (fallbackType === "plan") {
      const validPlans = unadorned.filter((item) => {
        if (Array.isArray(item)) {
          return (
            item.length > 0 && typeof item[0] === "object" && item[0] !== null
          );
        }
        return (
          typeof item === "object" &&
          item !== null &&
          (item.plan ||
            item.tasks ||
            item.steps ||
            item.title ||
            item.implementation_tasks)
        );
      });
      if (validPlans.length > 0) {
        return normalizeResult(validPlans[validPlans.length - 1]);
      }
    } else {
      return normalizeResult(unadorned[unadorned.length - 1]);
    }
  }

  console.warn(
    `⚠️ JSON parsing failed: no valid JSON block found. Failing closed.`,
  );
  if (fallbackType === "qa") {
    return {
      approval_status: "UNAPPROVED",
      findings: [
        {
          file: "unknown",
          line_numbers: [],
          narrative: `Malformed output: agent did not produce valid JSON findings.\nRaw response: ${safeText.trim()}`,
        },
      ],
    };
  }
  return null;
}

/**
 * Formats and saves an implementation plan to plans/current.md
 * @param {string|Object} textOrObj - Raw JSON plan string or parsed plan object
 * @param {string} [fallbackObjective] - Fallback objective if plan has no explicit objective
 * @returns {Promise<boolean>}
 */
export async function savePlanFromJSON(textOrObj, fallbackObjective = "") {
  const planJSON =
    typeof textOrObj === "object" && textOrObj !== null
      ? normalizePlanObject(textOrObj, fallbackObjective)
      : parseJSONFromText(textOrObj, "plan");
  if (planJSON && planJSON.title) {
    if (fallbackObjective) {
      if (!planJSON.objective || planJSON.objective === "Implementation Plan") {
        planJSON.objective = fallbackObjective;
      }
      if (!planJSON.title || planJSON.title === "Implementation Plan") {
        planJSON.title = fallbackObjective;
      }
    }
    const mdPlan = [
      `# IMPLEMENTATION PLAN: ${planJSON.title}`,
      "",
      `## 🎯 OBJECTIVE`,
      planJSON.objective || "",
      "",
      `## 🚧 SCOPE BOUNDARIES`,
      `**IN SCOPE:**`,
      ...(planJSON.scope_boundaries?.in_scope || []).map((i) => `- ${i}`),
      "",
      `**OUT OF SCOPE (Do NOT attempt):**`,
      ...(planJSON.scope_boundaries?.out_of_scope || []).map((i) => `- ${i}`),
      "",
      `## ✅ EXIT CRITERIA (Definition of Done)`,
      ...(planJSON.exit_criteria || []).map((i) => `- [ ] ${i}`),
      "",
      `## 🛠️ IMPLEMENTATION TASKS`,
      ...(planJSON.implementation_tasks || []).map((i) => `- [ ] ${i}`),
      "",
    ].join("\n");
    const repoRoot = await getRepoRoot();
    await fs.mkdir(path.join(repoRoot, "plans"), { recursive: true });
    await fs.writeFile(path.join(repoRoot, "plans/current.md"), mdPlan, "utf8");
    return true;
  }
  return false;
}

/**
 * Unwraps Zod wrapper types (effects, default, optional, nullable).
 * @param {import('zod').ZodTypeAny} s
 * @returns {import('zod').ZodTypeAny}
 */
function unwrapZodSchema(s) {
  if (s?._def?.schema) {
    return unwrapZodSchema(s._def.schema);
  }
  if (s?._def?.innerType) {
    return unwrapZodSchema(s._def.innerType);
  }
  return s;
}

/**
 * Generates a simplified object representation of a Zod schema structure.
 * @param {import('zod').ZodTypeAny} zodType
 * @returns {Object|Array|string}
 */
function describeZodType(zodType) {
  if (!zodType) {
    return "unknown";
  }
  const inner = unwrapZodSchema(zodType);
  if (inner.shape) {
    const shapeObj =
      typeof inner.shape === "function" ? inner.shape() : inner.shape;
    const out = {};
    for (const [k, v] of Object.entries(shapeObj)) {
      out[k] = describeZodType(v);
    }
    return out;
  }
  if (
    inner._def?.typeName === "ZodArray" ||
    (inner._def?.type && !inner._def?.values)
  ) {
    return [describeZodType(inner._def.type)];
  }
  if (inner._def?.values && Array.isArray(inner._def.values)) {
    return inner._def.values.join(" | ");
  }
  if (inner._def?.value !== undefined) {
    return JSON.stringify(inner._def.value);
  }
  if (inner._def?.options && Array.isArray(inner._def.options)) {
    return inner._def.options.map(describeZodType).join(" | ");
  }
  return (
    inner.description ||
    inner._def?.typeName?.replace("Zod", "") ||
    typeof inner
  );
}

/**
 * Validates and potentially manipulates agent outputs to fit a given Zod schema.
 * Orchestrates up to 4 validation agent sessions across 5 total validation attempts.
 * @param {string|Object} rawInput - The raw output text or pre-parsed object from another agent
 * @param {import('zod').ZodSchema} schema - The Zod schema to validate against
 * @param {string|Function|Object} [requestedModelOrRunner='gemini-3.8-flash'] - The model name, session runner function, or options object
 * @param {Function} [sessionRunner=null] - The session runner function to use for validator agent calls
 * @returns {Promise<Object|null>} The parsed and validated JSON object, or null on failure
 */
export async function validateAgentOutput(
  rawInput,
  schema,
  requestedModelOrRunner = "gemini-3.8-flash",
  sessionRunner = null,
) {
  if (!schema || typeof schema.parse !== "function") {
    throw new TypeError(
      "validateAgentOutput requires a valid Zod schema with a .parse() method.",
    );
  }

  let requestedModel = "gemini-3.8-flash";
  let runner = sessionRunner;

  if (typeof requestedModelOrRunner === "function") {
    runner = requestedModelOrRunner;
  } else if (typeof requestedModelOrRunner === "string") {
    requestedModel = requestedModelOrRunner;
  } else if (
    requestedModelOrRunner &&
    typeof requestedModelOrRunner === "object"
  ) {
    if (requestedModelOrRunner.requestedModel) {
      requestedModel = requestedModelOrRunner.requestedModel;
    }
    if (typeof requestedModelOrRunner.sessionRunner === "function") {
      runner = requestedModelOrRunner.sessionRunner;
    }
  }

  let currentData = rawInput;
  const rawTextRepresentation =
    typeof rawInput === "object" && rawInput !== null
      ? JSON.stringify(rawInput, null, 2)
      : String(rawInput);

  for (let i = 0; i < 5; i++) {
    try {
      let parsedObj;
      if (typeof currentData === "object" && currentData !== null) {
        parsedObj = currentData;
      } else {
        parsedObj = parseJSONFromText(currentData, null);
        if (!parsedObj) {
          parsedObj = JSON.parse(currentData);
        }
      }

      const validatedData = schema.parse(parsedObj);
      return validatedData;
    } catch (err) {
      console.log(
        `\n⚠️ Agent Output Validation failed (Try ${i + 1}/5): ${err?.message || err}. Retrying via Validator Agent...`,
      );
      if (i === 4) {
        console.error(
          `❌ Failed to validate agent output against schema after 5 tries. Last error: ${err?.message || err}`,
        );
        return null;
      }

      let schemaDetails = schema?.description || "";
      if (!schemaDetails) {
        try {
          const targetSchema = unwrapZodSchema(schema);
          const described = describeZodType(targetSchema);
          if (described && described !== "unknown") {
            schemaDetails = JSON.stringify(described, null, 2);
          }
        } catch (shapeErr) {
          console.debug(
            `[DEBUG] Failed to inspect schema shape: ${shapeErr?.message || shapeErr}`,
          );
          schemaDetails = "";
        }
      }

      const schemaContext = schemaDetails
        ? `Expected Schema Structure:\n${schemaDetails}\n\n`
        : "";

      const currentTextStr =
        typeof currentData === "object" && currentData !== null
          ? JSON.stringify(currentData, null, 2)
          : String(currentData);

      const dataComparison =
        i > 0 && currentTextStr !== rawTextRepresentation
          ? `Original Raw Input:\n${rawTextRepresentation}\n\nLatest Attempt That Failed:\n${currentTextStr}`
          : `Original Data:\n${currentTextStr}`;

      const prompt = `The following JSON output failed schema validation:
Error:
${err?.message || err}

${schemaContext}${dataComparison}

Please re-format the original data to strictly fit the required schema without losing or changing the semantic data. Output ONLY valid JSON matching the schema, with no markdown formatting or extra text.`;

      const systemInstructions = `You are a strict data formatter and schema validator agent. Your only job is to transform the provided data so it successfully parses against the target JSON schema. Do not change the meaning of the data. Return pure JSON block.`;

      try {
        if (!runner) {
          const runnerMod = await import("./agent-runner.js");
          runner = runnerMod.runAgentSession;
        }
        currentData = await runner({
          initialPrompt: prompt,
          systemInstructions: systemInstructions,
          requestedModel,
          maxTurns: 5,
          isolate: true,
        });
      } catch (sessionErr) {
        console.error(
          `❌ Validator Agent session failed: ${sessionErr?.message || sessionErr}`,
        );
        return null;
      }
    }
  }

  return null;
}

/**
 * Strips ANSI escape sequences, OSC sequences, and unsafe terminal control characters (including C0 and C1 controls).
 * @param {string} str
 * @returns {string}
 */
export function stripAnsi(str) {
  if (typeof str !== "string") {
    return String(str ?? "");
  }
  const esc = String.fromCharCode(27);
  const bell = String.fromCharCode(7);
  const c1Csi = String.fromCharCode(0x9b);
  const c1Osc = String.fromCharCode(0x9d);
  const c1St = String.fromCharCode(0x9c);

  const oscPattern = new RegExp(
    `(?:${esc}\\]|${c1Osc})[^${bell}${esc}${c1St}]*(?:${bell}|${esc}\\\\|${c1St})`,
    "g",
  );
  const csiPattern = new RegExp(`(?:${esc}\\[|${c1Csi})[0-?]*[ -/]*[@-~]`, "g");
  const esc2Pattern = new RegExp(`${esc}[@-Z\\\\-_]`, "g");
  const ctrlPattern = new RegExp(
    `[${String.fromCharCode(0)}-${String.fromCharCode(8)}${String.fromCharCode(11)}${String.fromCharCode(12)}${String.fromCharCode(14)}-${String.fromCharCode(31)}${String.fromCharCode(127)}-${String.fromCharCode(159)}]`,
    "g",
  );

  return str
    .replace(oscPattern, "")
    .replace(csiPattern, "")
    .replace(esc2Pattern, "")
    .replace(ctrlPattern, "");
}

/**
 * Wraps a long line while preserving leading indentation.
 * @param {string} line
 * @param {number} maxWidth
 * @returns {string[]}
 */
export function wrapLine(line, maxWidth) {
  if (typeof line !== "string" || !line) {
    return [line || ""];
  }
  const safeMaxWidth = Math.max(
    Number.isFinite(maxWidth) ? Math.floor(maxWidth) : 0,
    1,
  );
  if (line.length <= safeMaxWidth) {
    return [line];
  }
  const indentMatch = line.match(/^(\s*)/);
  const fullIndent = indentMatch ? indentMatch[1] : "";
  const indent = fullIndent.length < safeMaxWidth ? fullIndent : "";
  const trimmed = line.slice(fullIndent.length);
  const rawWords = trimmed.split(/\s+/).filter(Boolean);
  if (rawWords.length === 0) {
    const chunks = [];
    for (let i = 0; i < line.length; i += safeMaxWidth) {
      chunks.push(line.slice(i, i + safeMaxWidth));
    }
    return chunks.length > 0 ? chunks : [line];
  }

  const availableWidth = Math.max(safeMaxWidth - indent.length, 1);
  const words = [];
  for (const w of rawWords) {
    if (w.length <= availableWidth) {
      words.push(w);
    } else {
      for (let i = 0; i < w.length; i += availableWidth) {
        words.push(w.slice(i, i + availableWidth));
      }
    }
  }

  const wrapped = [];
  let current = indent;
  for (const word of words) {
    if (current === indent) {
      current += word;
    } else if (current.length + 1 + word.length <= safeMaxWidth) {
      current += " " + word;
    } else {
      wrapped.push(current);
      current = indent + word;
    }
  }
  if (current && current !== indent) {
    wrapped.push(current);
  }
  return wrapped.length > 0 ? wrapped : [line];
}

/**
 * Renders multi-line text inside a unicode box.
 * If text is single-line, it returns the original text unmodified.
 * @param {string} text
 * @param {number} [columns]
 * @param {'square'|'round'} [style='square']
 * @returns {string}
 */
export function renderBox(
  text,
  columns = process.stdout.columns,
  style = "square",
) {
  if (typeof text !== "string") {
    return String(text ?? "");
  }
  const termWidth = columns ? Math.max(columns - 6, 40) : 0;
  const rawLines = text.replace(/\r/g, "").split("\n");
  if (rawLines.length <= 1) {
    return text;
  }
  const lines = [];
  for (const rLine of rawLines) {
    if (termWidth > 0 && stripAnsi(rLine).length > termWidth) {
      lines.push(...wrapLine(rLine, termWidth));
    } else {
      lines.push(rLine);
    }
  }
  const maxLen = Math.max(...lines.map((l) => stripAnsi(l).length));
  const isRound = style === "round";
  const topLeft = isRound ? "╭─" : "┌─";
  const topRight = isRound ? "─╮" : "─┐";
  const bottomLeft = isRound ? "╰─" : "└─";
  const bottomRight = isRound ? "─╯" : "─┘";
  const top = `${topLeft}${"─".repeat(maxLen)}${topRight}`;
  const bottom = `${bottomLeft}${"─".repeat(maxLen)}${bottomRight}`;
  const middle = lines
    .map((l) => {
      const pad = maxLen - stripAnsi(l).length;
      return `│ ${l}${" ".repeat(pad)} │`;
    })
    .join("\n");
  return `${top}\n${middle}\n${bottom}`;
}

/**
 * Formats a tool name and arguments into a concise Gemini CLI style action header.
 * @param {string} toolName
 * @param {Object|string} [args]
 * @returns {string}
 */
export function formatToolAction(toolName, args = {}) {
  if (!toolName) {
    return "Tool";
  }
  let toolArgs = args;
  if (typeof args === "string") {
    try {
      toolArgs = JSON.parse(args);
    } catch {
      toolArgs = {};
    }
  }
  if (!toolArgs || typeof toolArgs !== "object") {
    toolArgs = {};
  }

  const normalizedToolName = String(toolName).toLowerCase();
  switch (normalizedToolName) {
  case "run_shell_command":
    return `Shell ${toolArgs.command || ""}`.trim();
  case "read_file":
    return `ReadFile ${toolArgs.file_path || ""}`.trim();
  case "write_file":
    return `WriteFile ${toolArgs.file_path || ""}`.trim();
  case "replace":
    return `Replace ${toolArgs.file_path || ""}`.trim();
  case "grep_search":
  case "grepsearch":
    return `GrepSearch ${toolArgs.pattern ? `'${toolArgs.pattern}'` : ""}`.trim();
  case "glob":
    return `Glob ${toolArgs.pattern ? `'${toolArgs.pattern}'` : ""}`.trim();
  case "list_directory":
    return `ListDirectory ${toolArgs.dir_path || "."}`.trim();
  case "ask_user":
    return "AskUser";
  case "update_topic":
    return `UpdateTopic ${toolArgs.title || toolArgs.strategic_intent || ""}`.trim();
  default: {
    const pascal = toolName
      .split("_")
      .map((part) =>
        part.length > 0 ? part[0].toUpperCase() + part.slice(1) : "",
      )
      .join("");
    const firstArg =
        Object.keys(toolArgs).length > 0
          ? String(Object.values(toolArgs)[0])
          : "";
    return firstArg ? `${pascal} ${firstArg}`.trim() : pascal;
  }
  }
}

/**
 * Renders a Gemini CLI style rounded tool execution card.
 * @param {Object} options
 * @param {string} options.toolName
 * @param {Object|string} [options.args]
 * @param {'success'|'running'|'error'|'info'} [options.status='success']
 * @param {any} [options.result]
 * @param {string} [options.agentType]
 * @param {number} [columns]
 * @returns {string}
 */
export function renderToolCard(
  { toolName, args, status = "success", result, agentType } = {},
  columns = process.stdout.columns,
) {
  const width =
    typeof columns === "number" && columns > 0
      ? columns
      : process.stdout.columns || 100;
  const boxWidth = Math.max(Math.min(width, 160), 40);
  const innerWidth = boxWidth - 4;

  const icons = {
    success: "✓  ",
    running: "⚡  ",
    error: "❌  ",
    info: "ℹ  ",
  };
  const icon = icons[status] || icons.success;
  const actionText = stripAnsi(formatToolAction(toolName, args)).replace(
    /[\r\n]+/g,
    " ",
  );
  const prefix = agentType
    ? `[${stripAnsi(agentType).replace(/[\r\n]+/g, " ")}] `
    : "";
  const header = `${icon}${prefix}${actionText}`;

  const lines = [];
  if (stripAnsi(header).length > innerWidth) {
    lines.push(...wrapLine(header, innerWidth));
  } else {
    lines.push(header);
  }

  let formattedResult = "";
  if (result !== undefined && result !== null) {
    if (typeof result === "string") {
      formattedResult = result;
    } else if (typeof result === "object") {
      if (typeof result.returnDisplay === "string") {
        formattedResult = result.returnDisplay;
      } else if (typeof result.llmContent === "string") {
        formattedResult = result.llmContent;
      } else {
        try {
          formattedResult = JSON.stringify(result, null, 2);
        } catch {
          formattedResult = String(result);
        }
      }
    } else {
      formattedResult = String(result);
    }
  }

  const MAX_RESULT_CHARS = 20000;
  let safeResult = formattedResult;
  if (
    typeof formattedResult === "string" &&
    formattedResult.length > MAX_RESULT_CHARS
  ) {
    safeResult = `${formattedResult.slice(0, MAX_RESULT_CHARS)}\n... [truncated remaining ${formattedResult.length - MAX_RESULT_CHARS} characters] ...`;
  }

  const trimmedResult = stripAnsi(safeResult)
    .replace(/\r\n/g, "\n")
    .replace(/\r/g, "")
    .trim();
  if (trimmedResult.length > 0) {
    lines.push(""); // Empty separator line matching Gemini CLI format
    const rawResultLines = trimmedResult.split("\n");
    const MAX_VISIBLE_LINES = 20;
    let visibleLines = rawResultLines;
    if (rawResultLines.length > MAX_VISIBLE_LINES) {
      const hidden = rawResultLines.length - MAX_VISIBLE_LINES;
      const tail = rawResultLines.slice(
        rawResultLines.length - MAX_VISIBLE_LINES,
      );
      visibleLines = [
        `... first ${hidden} lines hidden (Ctrl+O to show) ...`,
        ...tail,
      ];
    }
    for (const rLine of visibleLines) {
      if (stripAnsi(rLine).length > innerWidth) {
        lines.push(...wrapLine(rLine, innerWidth));
      } else {
        lines.push(rLine);
      }
    }
  }

  const top = `╭─${"─".repeat(innerWidth)}─╮`;
  const bottom = `╰─${"─".repeat(innerWidth)}─╯`;
  const middle = lines
    .map((l) => {
      const stripped = stripAnsi(l);
      const pad = Math.max(innerWidth - stripped.length, 0);
      return `│ ${l}${" ".repeat(pad)} │`;
    })
    .join("\n");

  return `${top}\n${middle}\n${bottom}`;
}

/**
 * Renders a prominent Gemini CLI style topic banner with title, objective, and summary.
 * @param {Object} options
 * @param {string} [options.topic] - Phase or topic title
 * @param {string} [options.objective] - Granular strategic intent
 * @param {string} [options.summary] - Contextual summary of previous work and new goals
 * @param {number} [columns]
 * @returns {string}
 */
export function renderTopicBanner(
  { topic, objective, summary } = {},
  columns = process.stdout.columns,
) {
  const width =
    typeof columns === "number" && columns > 0
      ? columns
      : process.stdout.columns || 100;
  const boxWidth = Math.max(Math.min(width, 160), 40);
  const innerWidth = boxWidth - 4;

  const rawLines = [];
  const normalizeText = (text) =>
    typeof text === "string"
      ? stripAnsi(text).replace(/\r\n/g, "\n").replace(/\r/g, "").trim()
      : "";

  const cleanTopic = normalizeText(topic);
  const cleanObjective = normalizeText(objective);
  const cleanSummary = normalizeText(summary);

  if (cleanTopic) {
    const topicLines = cleanTopic
      .split(/\r?\n/)
      .map((l) => l.trim())
      .filter(Boolean);
    if (topicLines.length > 0) {
      rawLines.push(`📂 Topic: ${topicLines[0]}`);
      for (let i = 1; i < topicLines.length; i++) {
        rawLines.push(`   ${topicLines[i]}`);
      }
    }
  }
  if (cleanObjective) {
    const objLines = cleanObjective
      .split(/\r?\n/)
      .map((l) => l.trim())
      .filter(Boolean);
    if (objLines.length > 0) {
      rawLines.push(`🎯 Objective: ${objLines[0]}`);
      for (let i = 1; i < objLines.length; i++) {
        rawLines.push(`   ${objLines[i]}`);
      }
    }
  }
  if (cleanSummary) {
    const sumLines = cleanSummary
      .split(/\r?\n/)
      .map((l) => l.trim())
      .filter(Boolean);
    if (sumLines.length > 0) {
      rawLines.push(`📝 Summary: ${sumLines[0]}`);
      for (let i = 1; i < sumLines.length; i++) {
        rawLines.push(`   ${sumLines[i]}`);
      }
    }
  }

  if (rawLines.length === 0) {
    rawLines.push("📂 Topic: (Active Goal)");
  }

  const lines = [];
  for (const rLine of rawLines) {
    if (stripAnsi(rLine).length > innerWidth) {
      lines.push(...wrapLine(rLine, innerWidth));
    } else {
      lines.push(rLine);
    }
  }

  const top = `╭─${"─".repeat(innerWidth)}─╮`;
  const bottom = `╰─${"─".repeat(innerWidth)}─╯`;
  const middle = lines
    .map((l) => {
      const stripped = stripAnsi(l);
      const pad = Math.max(innerWidth - stripped.length, 0);
      return `│ ${l}${" ".repeat(pad)} │`;
    })
    .join("\n");

  return `${top}\n${middle}\n${bottom}`;
}

/**
 * Renders prose and commentary text inside a rounded box aligned with tool cards and topic banners.
 * @param {string} text
 * @param {number} [columns]
 * @returns {string}
 */
export function renderProseCard(text, columns = process.stdout.columns) {
  if (typeof text !== "string") {
    return "";
  }
  const cleanText = stripAnsi(text)
    .replace(/\r\n/g, "\n")
    .replace(/\r/g, "")
    .trim();
  if (!cleanText) {
    return "";
  }

  const width =
    typeof columns === "number" && columns > 0
      ? columns
      : process.stdout.columns || 100;
  const boxWidth = Math.max(Math.min(width, 160), 40);
  const innerWidth = boxWidth - 4;

  const rawLines = cleanText.split("\n");
  const lines = [];
  for (const rLine of rawLines) {
    const sanitized = rLine.replace(/\t/g, "  ");
    if (sanitized.length > innerWidth) {
      lines.push(...wrapLine(sanitized, innerWidth));
    } else {
      lines.push(sanitized);
    }
  }

  const top = `╭─${"─".repeat(innerWidth)}─╮`;
  const bottom = `╰─${"─".repeat(innerWidth)}─╯`;
  const middle = lines
    .map((l) => {
      const pad = Math.max(innerWidth - l.length, 0);
      return `│ ${l}${" ".repeat(pad)} │`;
    })
    .join("\n");

  return `${top}\n${middle}\n${bottom}`;
}
