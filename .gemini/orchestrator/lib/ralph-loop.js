import { Buffer } from "node:buffer";
import { execFile } from "node:child_process";
import crypto from "node:crypto";
import fs from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";

import { isAbortError, runAgentSession } from "./agent-runner.js";
import { stripDiffMetadata } from "./git-release.js";
import { runQAPipeline } from "./qa-runner.js";
import {
  exists,
  getRepoRoot,
  loadWorkspaceSkillsContext,
  parseJSONFromText,
} from "./utils.js";

const execFileAsync = promisify(execFile);

const BINARY_EXTENSIONS = new Set([
  ".png",
  ".jpg",
  ".jpeg",
  ".gif",
  ".ico",
  ".webp",
  ".bmp",
  ".tiff",
  ".pdf",
  ".zip",
  ".tar",
  ".gz",
  ".tgz",
  ".bz2",
  ".xz",
  ".7z",
  ".bin",
  ".exe",
  ".dll",
  ".so",
  ".dylib",
  ".class",
  ".jar",
  ".wasm",
  ".iso",
  ".dmg",
  ".pkg",
  ".deb",
  ".rpm",
  ".pyc",
  ".pyo",
  ".pyd",
  ".o",
  ".obj",
  ".a",
  ".lib",
  ".out",
  ".sqlite",
  ".sqlite3",
  ".db",
  ".woff",
  ".woff2",
  ".ttf",
  ".mp4",
  ".mov",
  ".avi",
  ".mkv",
  ".mp3",
  ".wav",
  ".pack",
  ".idx",
]);

/**
 * Git empty tree hash (SHA-1 hash of an empty tree object).
 */
export const GIT_EMPTY_TREE_HASH = "4b825dc642cb6eb9a060e54bf8d69288fbee4904";

/**
 * Captures repository state from git and workspace artifacts.
 *
 * Architectural Note on Skills-First Mandate:
 * While mutating operations (such as staging via `stage-files.js`) strictly invoke skill scripts,
 * `captureGitState` executes read-only git status and diff commands directly via `runner` because:
 * 1. Rapid Loop Performance: Eliminates the overhead of repeatedly spawning Node.js runtime instances
 *    to execute CLI wrappers during tight inner-loop iterations.
 * 2. Dependency Injection & Testability: Accepts a pluggable `execFn` runner enabling hermetic unit testing
 *    without spawning physical processes or touching the real filesystem.
 * 3. Low-Level Control: Directly passes necessary low-level flags (`core.quotepath=false`, `--no-ext-diff`,
 *    `--no-color`, `maxBuffer`) and handles unborn HEAD fallback chains (`GIT_EMPTY_TREE_HASH`).
 *
 * @param {Object} [options]
 * @param {string} [options.repoRoot] - Root directory of the repository
 * @param {string} [options.planPath] - Path to the active plan file
 * @param {boolean} [options.stageBeforeIteration=true] - Whether to stage all changes before reading diff
 * @param {Function} [options.execFn] - Custom child_process exec function for testing
 * @param {number} [options.maxFileSizeBytes=30720] - Maximum file size in bytes for a single file to include in context
 * @param {number} [options.maxTotalContextBytes=153600] - Maximum total bytes across all modified files to include in context
 * @returns {Promise<{status: string, diff: string, diffHash: string, plan: string, changedFiles: string[], modifiedFilesContext: string}>}
 */
export async function captureGitState({
  repoRoot = null,
  planPath = null,
  stageBeforeIteration = true,
  execFn = null,
  maxFileSizeBytes = 30 * 1024,
  maxTotalContextBytes = 150 * 1024,
} = {}) {
  const root = repoRoot || (await getRepoRoot());
  const runner = execFn || execFileAsync;

  if (stageBeforeIteration) {
    const stageScript = path.join(
      root,
      ".gemini/skills/git/scripts/stage-files.js",
    );
    await runner("node", [stageScript, "--all", "--include-untracked"], {
      cwd: root,
      maxBuffer: 10 * 1024 * 1024,
    });
  }

  // 1. Get git porcelain status
  let status;
  try {
    const { stdout } = await runner(
      "git",
      ["-c", "core.quotepath=false", "status", "--porcelain"],
      {
        cwd: root,
        maxBuffer: 10 * 1024 * 1024,
      },
    );
    status = stdout;
  } catch (err) {
    console.debug(`[DEBUG] Failed to get git status: ${err.message}`);
    status = "";
  }

  // 2. Get git diff
  let diff;
  try {
    const { stdout } = await runner(
      "git",
      [
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        "HEAD",
      ],
      {
        cwd: root,
        maxBuffer: 10 * 1024 * 1024,
      },
    );
    diff = stdout;
  } catch {
    // If HEAD is unborn or error occurs, fallback to diff against the index (--cached)
    try {
      const { stdout } = await runner(
        "git",
        [
          "-c",
          "core.quotepath=false",
          "diff",
          "--no-ext-diff",
          "--no-color",
          "--cached",
        ],
        {
          cwd: root,
          maxBuffer: 10 * 1024 * 1024,
        },
      );
      diff = stdout;
    } catch {
      // Further fallback: compare staged index against empty tree hash
      try {
        console.debug(
          `[DEBUG] Comparing staged index against empty tree hash (${GIT_EMPTY_TREE_HASH}) as fallback.`,
        );
        const { stdout } = await runner(
          "git",
          [
            "-c",
            "core.quotepath=false",
            "diff",
            "--no-ext-diff",
            "--no-color",
            "--cached",
            GIT_EMPTY_TREE_HASH,
          ],
          {
            cwd: root,
            maxBuffer: 10 * 1024 * 1024,
          },
        );
        diff = stdout;
      } catch (err) {
        throw new Error(
          `All git diff strategies failed to capture repository diff: ${err?.message || err}`,
          {
            cause: err,
          },
        );
      }
    }
  }

  const sanitizedDiff = stripDiffMetadata(diff);
  const diffHash = crypto
    .createHash("sha256")
    .update(sanitizedDiff || "")
    .digest("hex");

  // Spillover diff artifact if diff is unusually large (> 50KB)
  const maxDiffInlineChars = 50 * 1024;
  let promptDiff = diff;
  if (diff && diff.length > maxDiffInlineChars) {
    const tmpDir = path.join(root, ".gemini/tmp");
    const artifactPath = path.join(tmpDir, "staged-review.diff");
    try {
      await fs.mkdir(tmpDir, { recursive: true });
      await fs.writeFile(artifactPath, diff, "utf8");
      const head = diff.slice(0, 10000);
      const tail = diff.slice(-10000);
      const relPath = path.relative(root, artifactPath);
      promptDiff = `${head}\n\n... [Git diff truncated (${diff.length} chars). Full diff saved to '${relPath}'. Inspect with read_file using start_line and end_line] ...\n\n${tail}`;
    } catch (err) {
      console.debug(
        `[DEBUG] Failed to write diff spillover artifact: ${err.message}`,
      );
    }
  }

  // 3. Read active plan file if present
  let plan = "";
  const resolvedPlanPath = planPath || path.join(root, "plans/current.md");
  if (await exists(resolvedPlanPath)) {
    try {
      plan = await fs.readFile(resolvedPlanPath, "utf8");
    } catch (err) {
      console.debug(
        `[DEBUG] Could not read plan file ${resolvedPlanPath}: ${err.message}`,
      );
      plan = "";
    }
  }

  // 4. Retrieve modified file paths
  let changedFiles;
  try {
    const { stdout } = await runner(
      "git",
      [
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        "--name-only",
        "HEAD",
      ],
      {
        cwd: root,
        maxBuffer: 10 * 1024 * 1024,
      },
    );
    changedFiles = stdout
      .split("\n")
      .map((f) => f.trim())
      .filter(Boolean);
  } catch {
    // If HEAD is unborn, fallback to --cached to list staged files
    try {
      const { stdout } = await runner(
        "git",
        [
          "-c",
          "core.quotepath=false",
          "diff",
          "--no-ext-diff",
          "--no-color",
          "--name-only",
          "--cached",
        ],
        {
          cwd: root,
          maxBuffer: 10 * 1024 * 1024,
        },
      );
      changedFiles = stdout
        .split("\n")
        .map((f) => f.trim())
        .filter(Boolean);
    } catch {
      // Further fallback: compare staged index against empty tree hash
      try {
        const { stdout } = await runner(
          "git",
          [
            "-c",
            "core.quotepath=false",
            "diff",
            "--no-ext-diff",
            "--no-color",
            "--name-only",
            "--cached",
            GIT_EMPTY_TREE_HASH,
          ],
          {
            cwd: root,
            maxBuffer: 10 * 1024 * 1024,
          },
        );
        changedFiles = stdout
          .split("\n")
          .map((f) => f.trim())
          .filter(Boolean);
      } catch {
        changedFiles = [];
      }
    }
  }

  // 5. Read contents of modified files for contextual visibility (capped to prevent token blowup)
  let totalContextBytes = 0;
  const fileBlocks = [];
  const maxOmittedPlaceholders = 50;
  let omittedPlaceholderCount = 0;
  let remainingOmittedCount = 0;

  const addOmittedPlaceholder = (relPath, reasonText) => {
    if (omittedPlaceholderCount < maxOmittedPlaceholders) {
      fileBlocks.push(`\n<file path="${relPath}">\n${reasonText}\n</file>\n`);
      omittedPlaceholderCount++;
      return true;
    }
    remainingOmittedCount++;
    return false;
  };

  const fileEntries = await Promise.all(
    changedFiles.map(async (relPath) => {
      if (
        relPath.startsWith(".git/") ||
        relPath === ".git" ||
        relPath.startsWith(`.git${path.sep}`)
      ) {
        console.debug(`[DEBUG] Skipping internal git path: ${relPath}`);
        return null;
      }
      const fullPath = path.join(root, relPath);
      if (!fullPath.startsWith(root + path.sep) && fullPath !== root) {
        console.debug(
          `[DEBUG] Skipping path escaping repository root: ${relPath}`,
        );
        return null;
      }
      const fileStats = await fs.lstat(fullPath).catch((err) => {
        if (err.code !== "ENOENT") {
          console.debug(
            `[DEBUG] Could not lstat path ${relPath}: ${err.message}`,
          );
        }
        return null;
      });
      if (!fileStats) {
        return null;
      }
      if (fileStats.isSymbolicLink()) {
        try {
          const resolvedPath = await fs.realpath(fullPath);
          if (
            !resolvedPath.startsWith(root + path.sep) &&
            resolvedPath !== root
          ) {
            console.debug(
              `[DEBUG] Skipping symlink resolving outside repository root: ${relPath} -> ${resolvedPath}`,
            );
            return null;
          }
          const targetStats = await fs.stat(resolvedPath).catch(() => null);
          if (targetStats && targetStats.isFile()) {
            return { relPath, fullPath: resolvedPath, stats: targetStats };
          }
        } catch (err) {
          console.debug(
            `[DEBUG] Could not resolve symlink ${relPath}: ${err.message}`,
          );
          return null;
        }
        return null;
      }
      if (fileStats.isFile()) {
        return { relPath, fullPath, stats: fileStats };
      }
      return null;
    }),
  );

  for (let i = 0; i < fileEntries.length; i++) {
    const entry = fileEntries[i];
    if (!entry) {
      continue;
    }
    const { relPath, fullPath, stats } = entry;
    try {
      if (totalContextBytes >= maxTotalContextBytes) {
        if (omittedPlaceholderCount < maxOmittedPlaceholders) {
          addOmittedPlaceholder(
            relPath,
            `[File content omitted: total modified context limit (${maxTotalContextBytes} bytes) reached. Inspect '${relPath}' with read_file using start_line and end_line if needed.]`,
          );
        } else {
          const remainingCount = fileEntries.length - i;
          fileBlocks.push(
            `\n[... and ${remainingCount} additional modified file(s) omitted due to total context limit]\n`,
          );
          break;
        }
        continue;
      }

      if (stats.size > maxFileSizeBytes) {
        addOmittedPlaceholder(
          relPath,
          `[File content omitted: size (${stats.size} bytes) exceeds limit of ${maxFileSizeBytes} bytes. Inspect '${relPath}' with read_file using start_line and end_line if needed.]`,
        );
        continue;
      }

      const ext = path.extname(relPath).toLowerCase();
      if (BINARY_EXTENSIONS.has(ext)) {
        addOmittedPlaceholder(
          relPath,
          `[File content omitted: binary file (${ext})]`,
        );
        continue;
      }

      // Check for null bytes (0x00) within the first 8000 bytes using open/read to detect binary content without buffering entire file
      let fileHandle = null;
      let isBinary = false;
      const checkBuffer = Buffer.alloc(Math.min(8000, stats.size));
      try {
        fileHandle = await fs.open(fullPath, "r");
        const { bytesRead } = await fileHandle.read(
          checkBuffer,
          0,
          checkBuffer.length,
          0,
        );
        if (checkBuffer.subarray(0, bytesRead).includes(0x00)) {
          isBinary = true;
        }
      } finally {
        if (fileHandle) {
          await fileHandle.close().catch(() => {});
        }
      }

      if (isBinary) {
        const extDesc = ext || "no extension";
        addOmittedPlaceholder(
          relPath,
          `[File content omitted: binary file detected (${extDesc}, ${stats.size} bytes)]`,
        );
        continue;
      }

      const buffer = await fs.readFile(fullPath);

      if (totalContextBytes + buffer.length > maxTotalContextBytes) {
        if (omittedPlaceholderCount < maxOmittedPlaceholders) {
          addOmittedPlaceholder(
            relPath,
            `[File content omitted: total modified context limit (${maxTotalContextBytes} bytes) exceeded. Inspect '${relPath}' with read_file using start_line and end_line if needed.]`,
          );
        } else {
          const remainingCount = fileEntries.length - i;
          fileBlocks.push(
            `\n[... and ${remainingCount} additional modified file(s) omitted due to total context limit]\n`,
          );
          totalContextBytes = maxTotalContextBytes;
          break;
        }
        totalContextBytes = maxTotalContextBytes;
        continue;
      }

      totalContextBytes += buffer.length;
      const content = buffer.toString("utf8");
      fileBlocks.push(`\n<file path="${relPath}">\n${content}\n</file>\n`);
    } catch (err) {
      console.debug(
        `[DEBUG] Could not read modified file ${relPath}: ${err.message}`,
      );
    }
  }

  if (remainingOmittedCount > 0) {
    fileBlocks.push(
      `\n[... and ${remainingOmittedCount} additional omitted file placeholder(s) omitted due to placeholder limit (${maxOmittedPlaceholders})]\n`,
    );
  }

  const modifiedFilesContext = fileBlocks.join("");

  return {
    status,
    diff: promptDiff,
    rawDiff: diff,
    diffHash,
    plan,
    changedFiles,
    modifiedFilesContext: modifiedFilesContext.trim(),
  };
}

/**
 * Spills oversized validation outputs (stdout / stderr) to disk artifacts in .gemini/tmp/
 * and provides head/tail previews with artifact file references.
 *
 * @param {Object} validationResult - Result object from validator function
 * @param {string} [repoRoot] - Root path of the repository
 * @param {number} [maxInlineChars=3000] - Threshold beyond which outputs spill over
 * @returns {Promise<Object>} Formatted validation result object
 */
export async function formatValidationSpillover(
  validationResult,
  repoRoot = null,
  maxInlineChars = 3000,
) {
  if (!validationResult || validationResult.success) {
    return validationResult;
  }

  const root = repoRoot || (await getRepoRoot());
  const tmpDir = path.join(root, ".gemini/tmp");

  const spilloverField = async (val, filename) => {
    if (!val) {
      return val;
    }
    const str = String(val);
    if (str.length <= maxInlineChars) {
      return str;
    }
    const artifactPath = path.join(tmpDir, filename);
    try {
      await fs.mkdir(tmpDir, { recursive: true });
      await fs.writeFile(artifactPath, str, "utf8");
      const halfLimit = Math.floor(maxInlineChars / 2);
      const head = str.slice(0, halfLimit);
      const tail = str.slice(-halfLimit);
      const relPath = path.relative(root, artifactPath);
      return `${head}\n\n... [Output truncated (${str.length} chars). Full log saved to '${relPath}'. Inspect with read_file using start_line and end_line] ...\n\n${tail}`;
    } catch (err) {
      console.debug(
        `[DEBUG] Failed to write validation artifact ${filename}: ${err.message}`,
      );
      return `... [truncated, showing last ${maxInlineChars} chars]\n${str.slice(-maxInlineChars)}`;
    }
  };

  const formattedStdout = await spilloverField(
    validationResult.stdout,
    "validation-stdout.log",
  );
  const formattedStderr = await spilloverField(
    validationResult.stderr ||
      validationResult.error ||
      validationResult.sessionError,
    "validation-stderr.log",
  );

  return {
    ...validationResult,
    stdout: formattedStdout,
    stderr: formattedStderr,
    rawStdout: validationResult.stdout,
    rawStderr: validationResult.stderr,
  };
}

/**
 * Detects whether repository or plan state has changed between two snapshots.
 *
 * @param {Object} before - State snapshot before agent execution
 * @param {Object} after - State snapshot after agent execution
 * @returns {{hasChanged: boolean, diffChanged: boolean, statusChanged: boolean, planChanged: boolean}}
 */
export function detectGitStateChange(before, after) {
  if (!before || !after) {
    return {
      hasChanged: true,
      diffChanged: true,
      statusChanged: true,
      planChanged: true,
    };
  }
  const diffChanged = before.diffHash !== after.diffHash;
  const statusChanged =
    (before.status || "").trimEnd() !== (after.status || "").trimEnd();
  const planChanged = (before.plan || "") !== (after.plan || "");
  const hasChanged = diffChanged || statusChanged || planChanged;
  return {
    hasChanged,
    diffChanged,
    statusChanged,
    planChanged,
  };
}

/**
 * Formats validation commands into a grammatically clean description string.
 *
 * @param {Array<string>|string} [validationCommands]
 * @returns {string} Formatted validation commands string
 */
export function formatValidationCommands(validationCommands) {
  if (Array.isArray(validationCommands) && validationCommands.length > 0) {
    const quoted = validationCommands.map((c) => `'${c}'`);
    if (quoted.length === 1) {
      return quoted[0];
    }
    if (quoted.length === 2) {
      return `${quoted[0]} and ${quoted[1]}`;
    }
    return `${quoted.slice(0, -1).join(", ")}, and ${quoted[quoted.length - 1]}`;
  }
  if (
    typeof validationCommands === "string" &&
    validationCommands.trim().length > 0
  ) {
    return validationCommands.trim();
  }
  return "";
}

/**
 * Builds the default review prompt with negative feedback injection from previous validation failures.
 *
 * @param {Object} context
 * @param {number} context.iteration - Current iteration number (1-based)
 * @param {number} context.maxIterations - Total allowable iterations
 * @param {Object} context.stateBefore - State snapshot captured before agent run
 * @param {Object|null} [context.lastValidation] - Previous iteration validation result
 * @param {string} [context.activePlan] - Active plan text
 * @returns {string} Formatted review prompt
 */
export function buildReviewPrompt({
  iteration,
  maxIterations,
  stateBefore,
  lastValidation = null,
  activePlan = "",
  validationCommands = [
    "node .gemini/skills/code-validation/scripts/lint.js all",
    "node .gemini/skills/code-validation/scripts/test.js all",
  ],
  skillsContext = "",
} = {}) {
  const diff = stateBefore?.diff || "";
  const modifiedFiles = stateBefore?.modifiedFilesContext || "";
  const plan = activePlan || stateBefore?.plan || "";
  const validationCmdsText = formatValidationCommands(validationCommands);
  const validationClause = validationCmdsText
    ? `- Ensure all tests and linters pass cleanly (e.g., ${validationCmdsText}).`
    : `- Ensure all tests and linters pass cleanly.`;

  let validationBlock = "";
  if (lastValidation && !lastValidation.success) {
    if (lastValidation.isSessionError || lastValidation.sessionError) {
      validationBlock = `
⚠️ AGENT SESSION ERROR FEEDBACK:
An unexpected execution or infrastructure error occurred during the previous session:

<session_error>
${lastValidation.sessionError || lastValidation.stderr || lastValidation.error}
</session_error>
`;
    } else {
      validationBlock = `
⚠️ CRITICAL VALIDATION FEEDBACK (Negative Backpressure):
The deterministic automated checks (tests / linters) failed in the previous iteration.
You MUST resolve and fix the following failures before making any other changes:

<validation_failure>
STDOUT:
${lastValidation.stdout || "(no stdout)"}

STDERR / ERRORS:
${lastValidation.stderr || lastValidation.error || "(no stderr)"}
</validation_failure>
`;
    }
  } else if (lastValidation && lastValidation.success) {
    validationBlock = `
<validation_status>
🟢 All workspace tests and linters passed cleanly in the previous iteration.
</validation_status>
`;
  }

  return `Review iteration ${iteration}/${maxIterations}.

Please review the staged code and the current plan (in plans/current.md), then implement any suggestions or fixes you have that are within the scope of the plan.

On your first turn, call update_topic with a descriptive "title" and "summary" (e.g. title: "Ralph Review Loop - Iteration ${iteration}", summary: "Review staged code and active plan, implement surgical improvements, and verify tests/linters").
${validationBlock}
Instructions & Constraints:
- Do NOT enter plan mode (do not call enter_plan_mode).
- Do NOT invoke subagents (do not call invoke_agent).
- Do NOT commit changes to git.
- Comply with the Skills-First Mandate: use installed skills in .gemini/skills/ (such as 'node .gemini/skills/git/scripts/get_diff_to_review.js' or 'node .gemini/skills/git/scripts/stage-files.js') for repository operations instead of raw mutating commands.
- The active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.
- Please add any decisions you make to the plan (in plans/current.md) so that further review cycles are refining the vision rather than overwriting it.
${validationClause}
- If you don't have any suggestions and all tests pass cleanly, please state that the code looks good and make no file modifications.
${skillsContext ? `\n${skillsContext}\n` : ""}
<active_plan>
${plan}
</active_plan>

<git_diff>
${diff}
</git_diff>

<modified_files>
${modifiedFiles}
</modified_files>`;
}

/**
 * Builds the prompt for the thinking review agent (Pro model).
 * Instructs the agent to strictly review staged code, rate each suggestion (critical/medium/minor),
 * and always output at least 5 suggestions without modifying any files.
 *
 * @param {Object} context
 * @returns {string} Formatted reviewer prompt
 */
export function buildReviewerPrompt({
  iteration = 1,
  maxIterations = 15,
  stateBefore = {},
  lastValidation = null,
  activePlan = "",
  skillsContext = "",
} = {}) {
  const diff = stateBefore?.diff || "";
  const modifiedFiles = stateBefore?.modifiedFilesContext || "";
  const plan = activePlan || stateBefore?.plan || "";

  let validationBlock = "";
  if (lastValidation && !lastValidation.success) {
    if (lastValidation.isSessionError || lastValidation.sessionError) {
      validationBlock = `
⚠️ PREVIOUS SESSION ERROR:
An unexpected infrastructure or runtime error occurred during the previous iteration:

<session_error>
${lastValidation.sessionError || lastValidation.stderr || lastValidation.error}
</session_error>
`;
    } else {
      validationBlock = `
⚠️ PREVIOUS VALIDATION FAILURE (Deterministic Backpressure):
Workspace tests or linters failed in the previous iteration:

<validation_failure>
STDOUT:
${lastValidation.stdout || "(no stdout)"}

STDERR / ERRORS:
${lastValidation.stderr || lastValidation.error || "(no stderr)"}
</validation_failure>
Prioritize suggestions that directly address and resolve these failures.
`;
    }
  }

  return `Review iteration ${iteration}/${maxIterations}.

You are an expert principal software engineer performing a strict, constructive pull request (PR) review on the staged git diff against the approved plan (in plans/current.md).
Your task is ONLY to review the code and output a numbered list of actionable suggestions within the scope of the plan.
${validationBlock}
Review Focus & Principles:
1. Analyze for Impact: Focus strictly on significant issues:
   - Bugs & Logic Flaws: Edge cases, off-by-one errors, state mutation risks, unhandled null/undefined values.
   - Security: Secret leakage, injection vulnerabilities, untrusted input handling.
   - Performance & Resources: Memory leaks, inefficient queries/loops, unnecessary allocations, concurrency/deadlock issues.
   - Architecture & Missing Tests: Idiomatic style, error propagation, missing unit test scenarios.
2. Minimize Pedantic Noise: Ignore minor stylistic preferences, import sorting, or formatting unless they violate critical language conventions or obscure readability.
3. Be Concrete: For every finding, cite the exact file name and line reference from the diff, explain the risk, and provide a concrete replacement snippet or fix rationale.

Instructions & Strict Constraints:
1. Your output must ONLY be a numbered list of suggestions, OR an empty list [] if no suggestions exist.
2. Every suggestion MUST be explicitly categorized with a severity rating of either "critical", "medium", or "minor". Do NOT add markdown formatting (such as asterisks or bolding, e.g. **[critical]** or [**critical**]) around or inside the severity tags or brackets:
   - [critical]: Must-fix bugs, data loss, security vulnerabilities, or regression risks.
   - [medium]: Missing tests, error handling gaps, resource leaks, or architectural inconsistencies.
   - [minor]: Clean-up recommendations, maintainability improvements, or non-blocking edge cases.
3. Goal: Strive for a minimum of 5 suggestions per review. If the core code is solid, look for minor or trivial enhancements such as test coverage, documentation, type/error assertions, or maintainability.
4. Escape Hatch: If you cannot find 5 suggestions even after considering minor or trivial enhancements, provide whatever valid suggestions you identified. If the staged code meets all requirements and absolutely no enhancements or changes are warranted, output: []
5. Do NOT make file modifications, call mutating tools, or commit code.
6. Consider the repository's installed skills (in .gemini/skills/) and the Skills-First Mandate when reviewing code and proposing suggestions.

Format each suggestion clearly as:
<number>. [<severity>] <Summary>: <File:line reference> - <Explanation of risk and recommended fix>

Example format:
1. [critical] Fix resource leak: agent-scripts/lib/ralph-loop.js:120 - Ensure stream is closed in finally block...
2. [medium] Add missing unit test: internal/provider/resource_test.go:45 - Add test case for empty directory input...
3. [minor] Variable naming: agent-scripts/lib/utils.js:80 - Rename tmpVar to descriptive identifier...
${skillsContext ? `\n${skillsContext}\n` : ""}
<active_plan>
${plan}
</active_plan>

<git_diff>
${diff}
</git_diff>

<modified_files>
${modifiedFiles}
</modified_files>`;
}

/**
 * Builds the counting prompt for the flash-lite model to parse suggestion severity counts.
 *
 * @param {string} suggestionsText - Raw suggestions text from the reviewer agent
 * @returns {string} Formatted tally prompt
 */
export function buildTallyPrompt(suggestionsText = "") {
  return `You are a strict text analysis and counting model.
Read the following list of review suggestions and count the number of suggestions rated as "critical", "medium", and "minor".

Output ONLY a raw JSON object (no markdown code blocks, backticks, or other text):
{
  "critical": <integer>,
  "medium": <integer>,
  "minor": <integer>
}

Review Suggestions to evaluate:
${suggestionsText}`;
}

/**
 * Builds the implementation prompt for the Flash model to apply approved suggestions.
 *
 * @param {Object} context
 * @returns {string} Formatted implementation prompt
 */
export function buildImplementationPrompt({
  suggestionsText = "",
  tally = { critical: 0, medium: 0, minor: 0 },
  activePlan = "",
  stateBefore = {},
  lastValidation = null,
  iteration = 1,
  maxIterations = 15,
  validationCommands = [
    "node .gemini/skills/code-validation/scripts/lint.js all",
    "node .gemini/skills/code-validation/scripts/test.js all",
  ],
  skillsContext = "",
} = {}) {
  const diff = stateBefore?.diff || "";
  const modifiedFiles = stateBefore?.modifiedFilesContext || "";
  const plan = activePlan || stateBefore?.plan || "";
  const validationCmdsText = formatValidationCommands(validationCommands);
  const validationClause = validationCmdsText
    ? `- Ensure all tests and linters pass cleanly (e.g., ${validationCmdsText}).`
    : `- Ensure all tests and linters pass cleanly.`;

  let validationBlock = "";
  if (lastValidation && !lastValidation.success) {
    if (lastValidation.isSessionError || lastValidation.sessionError) {
      validationBlock = `
⚠️ PREVIOUS SESSION ERROR:
An unexpected infrastructure or runtime error occurred during the previous iteration:

<session_error>
${lastValidation.sessionError || lastValidation.stderr || lastValidation.error}
</session_error>
`;
    } else {
      validationBlock = `
⚠️ PREVIOUS VALIDATION FAILURE (Deterministic Backpressure):
The automated workspace checks failed in the previous iteration:

<validation_failure>
STDOUT:
${lastValidation.stdout || "(no stdout)"}

STDERR / ERRORS:
${lastValidation.stderr || lastValidation.error || "(no stderr)"}
</validation_failure>
You must resolve these failures while implementing the suggestions.
`;
    }
  }

  return `Review iteration ${iteration}/${maxIterations}.

The thinking review agent identified the following suggestions (${tally.critical} critical, ${tally.medium} medium, ${tally.minor} minor):

<suggestions>
${suggestionsText}
</suggestions>
${validationBlock}
Instructions & Constraints:
- On your first turn, call update_topic with a descriptive "title" and "summary" (e.g. title: "Implement Review Suggestions - Iteration ${iteration}", summary: "Implement critical and medium suggestions from review, update plan, and verify tests/linters").
- Since there are critical and/or medium suggestions, you must implement ALL suggestions that are within scope of the plan.
- You have access to GrepSearch (and grep_search) to search across workspace file contents.
- Do NOT enter plan mode (do not call enter_plan_mode).
- Do NOT invoke subagents (do not call invoke_agent).
- Comply strictly with the Skills-First Mandate: use installed skills in .gemini/skills/ (such as 'node .gemini/skills/git/scripts/get_diff_to_review.js' or 'node .gemini/skills/git/scripts/stage-files.js') for repository operations instead of raw mutating commands.
- The active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.
- Formulate a clear implementation plan for addressing the review suggestions and append it to 'plans/current.md'.
- CRITICAL PLAN INTEGRITY RULE: Do NOT overwrite or wipe 'plans/current.md'. It represents all changes and historical context across the entire PR. You must only append new sections (such as review remediation plans and decisions) or edit/refine specific checklist items to preserve all PR context.
- Do NOT commit changes to git.
${validationClause}
${skillsContext ? `\n${skillsContext}\n` : ""}
<active_plan>
${plan}
</active_plan>

<git_diff>
${diff}
</git_diff>

<modified_files>
${modifiedFiles}
</modified_files>`;
}

/**
 * Parses individual numbered or bulleted suggestions from reviewer text into discrete suggestion objects.
 *
 * @param {string} [text=''] - Raw reviewer text containing suggestions
 * @returns {Array<{id: number, severity: string, text: string, raw: string}>}
 */
export function parseIndividualSuggestions(text = "") {
  if (typeof text !== "string" || !text.trim()) {
    return [];
  }
  const lines = text.replace(/\r/g, "").split("\n");
  const items = [];
  let currentItem = null;
  let insideCodeBlock = false;
  let codeFenceMarker = null;

  // Matches lines like: 1. [critical] ..., 2. [medium] ..., * [minor] ..., 1. Title ...
  // Tolerates leading whitespace for indented list items.
  const headerRegex =
    /^\s*(?:[*-]|\d+[.)])\s+(?:\[\s*([a-z]+)\s*\]|\(\s*([a-z]+)\s*\))?\s*(.*)$/i;

  for (const line of lines) {
    const trimmed = line.trim();
    if (!currentItem && !trimmed) {
      continue;
    }
    const isBacktickFence = trimmed.startsWith("```");
    const isTildeFence = trimmed.startsWith("~~~");
    const isCodeFence = insideCodeBlock
      ? codeFenceMarker
        ? trimmed.startsWith(codeFenceMarker)
        : isBacktickFence || isTildeFence
      : isBacktickFence || isTildeFence;
    const wasInsideCodeBlock = insideCodeBlock;
    if (isCodeFence) {
      const marker = isBacktickFence ? "```" : "~~~";
      const isSingleLineFence =
        !insideCodeBlock &&
        trimmed.length >= 6 &&
        trimmed.slice(3).includes(marker);
      if (!isSingleLineFence) {
        insideCodeBlock = !insideCodeBlock;
        codeFenceMarker = insideCodeBlock ? marker : null;
      }
    }
    const isInCodeBlock = wasInsideCodeBlock || insideCodeBlock || isCodeFence;
    let match = !isInCodeBlock ? line.match(headerRegex) : null;
    if (match && /^\s+/.test(line) && !match[1] && !match[2] && currentItem) {
      match = null;
    }
    if (!match && insideCodeBlock) {
      const candidateMatch = line.match(headerRegex);
      if (candidateMatch) {
        const rawSeverity = (
          candidateMatch[1] ||
          candidateMatch[2] ||
          ""
        ).toLowerCase();
        if (["critical", "medium", "minor"].includes(rawSeverity)) {
          insideCodeBlock = false;
          codeFenceMarker = null;
          match = candidateMatch;
        }
      }
    }
    if (match) {
      if (currentItem) {
        currentItem.text = currentItem.text.replace(/^\n+|\s+$/g, "");
        currentItem.raw = currentItem.raw.trim();
        items.push(currentItem);
      }
      const rawSeverity = (match[1] || match[2] || "").toLowerCase();
      const severity = ["critical", "medium", "minor"].includes(rawSeverity)
        ? rawSeverity
        : "unclassified";
      const body = match[3] ? match[3].trim() : "";
      currentItem = {
        id: items.length + 1,
        severity,
        text: body,
        raw: trimmed,
      };
    } else if (currentItem) {
      if (!trimmed) {
        if (currentItem.text && !currentItem.text.endsWith("\n")) {
          currentItem.text += "\n";
        }
        if (currentItem.raw && !currentItem.raw.endsWith("\n")) {
          currentItem.raw += "\n";
        }
      } else {
        const content = line.trimEnd();
        currentItem.text = currentItem.text
          ? `${currentItem.text}\n${content}`
          : content;
        currentItem.raw = currentItem.raw
          ? `${currentItem.raw}\n${line.trimEnd()}`
          : line.trimEnd();
      }
    }
  }

  if (currentItem) {
    currentItem.text = currentItem.text.replace(/^\n+|\s+$/g, "");
    currentItem.raw = currentItem.raw.trim();
    items.push(currentItem);
  }

  return items;
}

/**
 * Extracts completed suggestion ID(s) from implementer output based on tags and completion verbs.
 *
 * @param {string} [agentOutput='']
 * @param {Array<{id: number}>} [remainingSuggestions=[]]
 * @returns {number[]} Array of completed suggestion IDs, or empty array if none detected
 */
export function extractCompletedSuggestionIds(
  agentOutput = "",
  remainingSuggestions = [],
) {
  if (typeof agentOutput !== "string") {
    return [];
  }

  const matchedIds = new Set();

  // 1. Tag match: <completed_suggestion id="1,2"> or <completed_suggestion>1, 2</completed_suggestion>
  const tagMatches = agentOutput.matchAll(
    /<completed_suggestion[^>]*id=["']?([\d\s,]+)["']?[^>]*>/gi,
  );
  for (const match of tagMatches) {
    const ids = match[1]
      .split(",")
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));
    ids.forEach((id) => {
      if (remainingSuggestions.some((s) => s.id === id)) {
        matchedIds.add(id);
      }
    });
  }

  const tagBodyMatches = agentOutput.matchAll(
    /<completed_suggestion[^>]*>\s*(?:(?:suggestions?\s*)?(?:#\s*|id\b\s*:?\s*)|suggestions?\s+)?([\d\s,]+)\b.*?<\/completed_suggestion>/gis,
  );
  for (const match of tagBodyMatches) {
    const ids = match[1]
      .split(",")
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));
    ids.forEach((id) => {
      if (remainingSuggestions.some((s) => s.id === id)) {
        matchedIds.add(id);
      }
    });
  }

  if (matchedIds.size > 0) {
    return Array.from(matchedIds);
  }

  // 2. Conversational match: action verbs ("completed suggestions 1, 2", "suggestion #3 resolved")
  const phraseMatches = [
    ...agentOutput.matchAll(
      /(?:completed?|completing|implemented?|implementing|resolved?|resolving|addressed?|addressing|fixed?|fixing)\s+suggestions?\s*(?:#|id\s*:?\s*)?([\d\s,]+)/gi,
    ),
    ...agentOutput.matchAll(
      /suggestions?\s*(?:#|id\s*)?([\d\s,]+)\s+(?:completed?|completing|implemented?|implementing|resolved?|resolving|addressed?|addressing|fixed?|fixing)/gi,
    ),
  ];

  for (const match of phraseMatches) {
    const ids = match[1]
      .split(",")
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));
    ids.forEach((id) => {
      if (remainingSuggestions.some((s) => s.id === id)) {
        matchedIds.add(id);
      }
    });
  }

  if (matchedIds.size > 0) {
    return Array.from(matchedIds);
  }

  return [];
}

/**
 * Builds prompt for batch suggestion iteration within the implementation Ralph loop.
 *
 * @param {Object} options
 * @returns {string} Formatted prompt
 */
export function buildBatchSuggestionImplementationPrompt({
  remainingSuggestions = [],
  activePlan = "",
  stateBefore = {},
  lastValidation = null,
  currentSuggestionIndex = 1,
  totalSuggestions = 1,
  validationCommands = [
    "node .gemini/skills/code-validation/scripts/lint.js all",
    "node .gemini/skills/code-validation/scripts/test.js all",
  ],
  skillsContext = "",
} = {}) {
  const diff = stateBefore?.diff || "";
  const modifiedFiles = stateBefore?.modifiedFilesContext || "";
  const plan = activePlan || stateBefore?.plan || "";
  const validationCmdsText = formatValidationCommands(validationCommands);
  const validationClause = validationCmdsText
    ? `- Ensure all tests and linters pass cleanly (e.g., ${validationCmdsText}).`
    : `- Ensure all tests and linters pass cleanly.`;

  let validationBlock = "";
  if (lastValidation && !lastValidation.success) {
    const maxValidationChars = 2000;
    const formatOutput = (text, fallback) => {
      if (!text) {
        return fallback;
      }
      const str = String(text);
      if (str.length <= maxValidationChars) {
        return str;
      }
      return `... [truncated, showing last ${maxValidationChars} chars]\n${str.slice(-maxValidationChars)}`;
    };

    const stdoutText = formatOutput(lastValidation.stdout, "(no stdout)");
    const stderrText = formatOutput(
      lastValidation.stderr || lastValidation.error,
      "(no stderr)",
    );

    validationBlock = `
⚠️ PREVIOUS VALIDATION FAILURE (Deterministic Backpressure):
The automated workspace checks failed in the previous iteration.
You MUST resolve and fix the following failures while addressing the suggestions:

<validation_failure>
STDOUT:
${stdoutText}

STDERR / ERRORS:
${stderrText}
</validation_failure>
Prioritize resolving these failures cleanly.
`;
  }

  const formattedRemaining = remainingSuggestions
    .map(
      (s, idx) => `${idx + 1}. [${s.severity}] (Original #${s.id}) ${s.text}`,
    )
    .join("\n");

  return `Implementation Phase - Step ${currentSuggestionIndex} of ${totalSuggestions}.

There are ${remainingSuggestions.length} remaining suggestion(s) to implement:

<remaining_suggestions>
${formattedRemaining}
</remaining_suggestions>
${validationBlock}
Instructions & Constraints:
- On your first turn, call update_topic with a descriptive "title" and "summary" (e.g. title: "Implement Suggestions ${currentSuggestionIndex}-${Math.min(currentSuggestionIndex + 2, totalSuggestions)}/${totalSuggestions}").
- You may implement UP TO 3 SUGGESTIONS in this single turn if they are related or straightforward. Do not overwhelm the turn budget.
- Focus surgically on those suggestions and keep changes minimal and precise.
- You have access to GrepSearch (and grep_search) to search across workspace file contents.
- Do NOT enter plan mode (do not call enter_plan_mode).
- Do NOT invoke subagents (do not call invoke_agent).
- Comply strictly with the Skills-First Mandate: use installed skills in .gemini/skills/ (such as 'node .gemini/skills/git/scripts/get_diff_to_review.js' or 'node .gemini/skills/git/scripts/stage-files.js') for repository operations instead of raw mutating commands.
- The active plan from 'plans/current.md' is already embedded below in <active_plan>. Do NOT waste turns reading it from disk.
- If this is Step 1 (or iteration 1), formulate a clear implementation plan for addressing the suggestions and append it to 'plans/current.md'.
- CRITICAL PLAN INTEGRITY RULE: Do NOT overwrite or wipe 'plans/current.md'. It represents all changes and historical context across the entire PR. You must only append new sections (such as review remediation plans and decisions) or edit/refine specific checklist items to preserve all PR context.
- Do NOT commit changes to git.
${validationClause}
- When finished, include a tag indicating WHICH suggestions were implemented (comma-separated if multiple):
  <completed_suggestion id="${remainingSuggestions[0]?.id || 1}">Brief summary of the fixes...</completed_suggestion>
${skillsContext ? `\n${skillsContext}\n` : ""}
<active_plan>
${plan}
</active_plan>

<git_diff>
${diff}
</git_diff>

<modified_files>
${modifiedFiles}
</modified_files>`;
}

/**
 * Parses counts of "critical", "medium", and "minor" suggestions from model output or raw text.
 *
 * @param {string} [primaryText=''] - Output from flash-lite model or JSON string
 * @param {string} [fallbackText=''] - Raw reviewer text used as fallback
 * @param {Function|null} [logger=null] - Optional logger function for debug messages (defaults to no-op)
 * @returns {{critical: number, medium: number, minor: number, total: number}}
 */
export function parseSuggestionCounts(
  primaryText = "",
  fallbackText = "",
  logger = null,
) {
  const logDebug = typeof logger === "function" ? logger : () => {};
  const candidates = [primaryText, fallbackText].filter(Boolean);

  for (const text of candidates) {
    if (typeof text !== "string" || text.trim().length === 0) {
      continue;
    }

    if (text.trim() === "[]") {
      return { critical: 0, medium: 0, minor: 0, total: 0 };
    }

    // 1. Try parsing JSON block or object via parseJSONFromText
    try {
      const rawParsed = parseJSONFromText(text, null);
      if (Array.isArray(rawParsed) && rawParsed.length === 0) {
        return { critical: 0, medium: 0, minor: 0, total: 0 };
      }
      if (
        rawParsed &&
        typeof rawParsed === "object" &&
        !Array.isArray(rawParsed)
      ) {
        // Normalize keys to lower-case to defend against model casing variations
        const parsed = Object.fromEntries(
          Object.entries(rawParsed).map(([key, val]) => [
            key.toLowerCase().trim(),
            val,
          ]),
        );
        const hasCritical =
          parsed.critical !== undefined && parsed.critical !== null;
        const hasMedium = parsed.medium !== undefined && parsed.medium !== null;
        const hasMinor = parsed.minor !== undefined && parsed.minor !== null;

        if (hasCritical && hasMedium && hasMinor) {
          const criticalParsed = parseInt(parsed.critical, 10);
          const mediumParsed = parseInt(parsed.medium, 10);
          const minorParsed = parseInt(parsed.minor, 10);

          if (
            !Number.isNaN(criticalParsed) &&
            !Number.isNaN(mediumParsed) &&
            !Number.isNaN(minorParsed)
          ) {
            const critical = Math.max(0, criticalParsed);
            const medium = Math.max(0, mediumParsed);
            const minor = Math.max(0, minorParsed);
            return {
              critical,
              medium,
              minor,
              total: critical + medium + minor,
            };
          }
        }

        logDebug(
          `[DEBUG] Incomplete JSON schema in parseSuggestionCounts (must include all numeric keys: critical, medium, minor): ${JSON.stringify(parsed)}, falling back to heuristics.`,
        );
      } else if (rawParsed === null) {
        logDebug(
          `[DEBUG] parseJSONFromText returned null for suggestions count text: "${text.slice(0, 50)}", falling back to heuristics.`,
        );
      }
    } catch (err) {
      const errMsg =
        (err.message || "").length > 120
          ? `${(err.message || "").slice(0, 120)}...`
          : err.message;
      logDebug(
        `[DEBUG] Exception in parseSuggestionCounts JSON parse, falling back to heuristics: ${errMsg}`,
      );
    }

    // 2. Check for explicit text counts like "5 critical, 2 medium, 1 minor" or "critical: 5, medium: 2, minor: 1"
    // Also supports pluralized forms like "5 criticals, 2 mediums, 1 minor"
    // Use global matchAll to extract the final occurrence, defending against conversational drifting/preambles
    const extractCount = (cat) => {
      const matches = Array.from(
        text.matchAll(
          new RegExp(
            `(?:(\\d+)\\s+${cat}s?|${cat}s?\\s*[:=-]?\\s*(\\d+))`,
            "gi",
          ),
        ),
      );
      if (matches.length === 0) {
        return null;
      }
      const last = matches[matches.length - 1];
      const numStr = last[1] || last[2];
      return numStr !== undefined ? parseInt(numStr, 10) : null;
    };

    const critCount = extractCount("critical");
    const medCount = extractCount("medium");
    const minCount = extractCount("minor");

    if (critCount !== null || medCount !== null || minCount !== null) {
      const critical = Math.max(0, critCount ?? 0);
      const medium = Math.max(0, medCount ?? 0);
      const minor = Math.max(0, minCount ?? 0);
      return {
        critical,
        medium,
        minor,
        total: critical + medium + minor,
      };
    }
  }

  // 3. Fallback: scan text for bracketed tag occurrences anchored to list structures (e.g. "1. [critical]" or "- [critical]" or "* [critical]")
  // Strip code backticks, bold/italic formatting, and tag-wrapping markers while preserving markdown bullet characters like "* [critical]"
  const normalizedText = (fallbackText || primaryText || "")
    .replace(/`([^`\n]+)`/g, "$1")
    .replace(/\*\*(.*?)\*\*/g, "$1")
    .replace(/(\[[^\]\n]+\])[*`_]/g, "$1")
    .replace(/[*`_](\[[^\]\n]+\])/g, "$1");
  const severities = ["critical", "medium", "minor"];
  const counts = { critical: 0, medium: 0, minor: 0 };

  let foundList = false;
  for (const sev of severities) {
    const listMatches =
      normalizedText.match(
        new RegExp(`^\\s*(?:\\d+[.)]|[-*])\\s*\\[\\s*${sev}\\s*\\]`, "gim"),
      ) || [];
    counts[sev] = listMatches.length;
    if (counts[sev] > 0) {
      foundList = true;
    }
  }
  if (foundList) {
    counts.total = counts.critical + counts.medium + counts.minor;
    return counts;
  }

  // Also check line-anchored bracket tags without numbers (e.g. "[critical]")
  let foundLine = false;
  for (const sev of severities) {
    const lineMatches =
      normalizedText.match(new RegExp(`^\\s*\\[\\s*${sev}\\s*\\]`, "gim")) ||
      [];
    counts[sev] = lineMatches.length;
    if (counts[sev] > 0) {
      foundLine = true;
    }
  }
  if (foundLine) {
    counts.total = counts.critical + counts.medium + counts.minor;
    return counts;
  }

  // 4. Bracketed or parenthesized severity tag fallback: check for [critical] or (critical) anywhere in the text
  let foundTag = false;
  for (const sev of severities) {
    const tagMatches =
      normalizedText.match(new RegExp(`[\\[(]\\s*${sev}\\s*[\\])]`, "gim")) ||
      [];
    counts[sev] = tagMatches.length;
    if (counts[sev] > 0) {
      foundTag = true;
    }
  }
  if (foundTag) {
    counts.total = counts.critical + counts.medium + counts.minor;
    return counts;
  }

  // If no structural tags or explicit counts were matched, return zero counts to prevent false positives
  return { critical: 0, medium: 0, minor: 0, total: 0 };
}

/**
 * Invokes the flash-lite model to read suggestions and return parsed category counts.
 *
 * @param {Object} options
 * @param {string} options.suggestionsText - Raw reviewer output
 * @param {Function} [options.sessionRunner=runAgentSession] - Session runner function
 * @param {string} [options.model='gemini-3.1-flash-lite'] - Model to use
 * @param {Object|null} [options.rl=null] - Readline interface
 * @param {Function} [options.logger=console.log] - Logger
 * @returns {Promise<{critical: number, medium: number, minor: number, total: number}>}
 */
export async function tallySuggestionsWithFlashLite({
  suggestionsText = "",
  sessionRunner = runAgentSession,
  model = "gemini-3.1-flash-lite",
  rl = null,
  logger: inputLogger = console.log,
} = {}) {
  const logger = typeof inputLogger === "function" ? inputLogger : console.log;
  const prompt = buildTallyPrompt(suggestionsText);

  let rawOutput = "";
  try {
    rawOutput = await sessionRunner({
      initialPrompt: prompt,
      systemInstructions:
        'You are a strict text analysis and counting agent. Output ONLY a raw, valid JSON object with numeric keys "critical", "medium", and "minor" without any markdown formatting, backticks, code fences, or explanatory text. Do NOT call any tools.',
      requestedModel: model,
      maxTurnsOverride: 3,
      requireTopic: false,
      isolate: true,
      rl,
    });
  } catch (err) {
    logger(
      `[DEBUG] Flash-lite tally session error: ${String(err.message || err)}. Falling back to deterministic parsing.`,
    );
  }

  return parseSuggestionCounts(rawOutput, suggestionsText, logger);
}

const DEFAULT_REVIEW_SYSTEM_INSTRUCTIONS = `You are in the REVIEW LOOP phase. Review the staged code and the approved plan documented in 'plans/current.md', make surgical improvements within scope if necessary, update the plan with any decisions made, and ensure all tests and linters pass. If no changes or improvements are needed, make no modifications.`;

/**
 * Runs a multi-purpose Ralph Loop (Wiggum Loop).
 * Enforces fresh context per iteration, externalized state in filesystem/git,
 * deterministic automated backpressure, negative feedback injection, and explicit exit criteria.
 *
 * @param {Object} [options]
 * @param {Function} [options.sessionRunner=runAgentSession] - Async function executing an agent session
 * @param {Function|null} [options.validator=runQAPipeline] - Async function checking code validity
 * @param {Function|null} [options.iterationRunner=null] - Custom multi-agent per-iteration coordinator
 * @param {Function|string} [options.buildPrompt=buildReviewPrompt] - Prompt builder or prompt string
 * @param {Function|string} [options.systemInstructions] - System instructions string or builder
 * @param {string|null} [options.model=null] - Model to use for the session
 * @param {number} [options.maxIterations=15] - Maximum allowable outer loop iterations
 * @param {number} [options.maxTurnsPerSession=20] - Turn capacity per agent session
 * @param {number} [options.noProgressLimit=3] - Maximum consecutive iterations without progress before aborting
 * @param {number} [options.maxFileSizeBytes=30720] - Maximum file size in bytes per file
 * @param {number} [options.maxTotalContextBytes=153600] - Maximum total bytes across all modified files
 * @param {string|null} [options.repoRoot=null] - Repository root
 * @param {string|null} [options.planPath=null] - Path to plan file
 * @param {boolean} [options.requireTopic=true] - Whether to require topic updates
 * @param {boolean} [options.validateIntent=false] - Whether to require intent declarations
 * @param {Array<string>} [options.blockTools] - Tools to unregister/block
 * @param {Array<any>} [options.customTools=[]] - Custom tools to pass to agent session
 * @param {Object|null} [options.rl=null] - Readline interface for interactive prompts
 * @param {boolean} [options.stageBeforeIteration=true] - Whether to run git add -A before diffing
 * @param {boolean} [options.requireInitialDiff=false] - If true, exits immediately when initial diff is empty
 * @param {Function|null} [options.checkCompletion=null] - Custom completion predicate
 * @param {Function|null} [options.captureState=null] - Custom state capturer
 * @param {Function|null} [options.detectChange=null] - Custom change detector
 * @param {Function} [options.logger=console.log] - Output logger
 * @param {Function|null} [options.onIterationStart=null] - Lifecycle hook before iteration starts
 * @param {Function|null} [options.onIterationEnd=null] - Lifecycle hook after iteration ends
 * @returns {Promise<{success: boolean, completed: boolean, iterations: number, reason: string, finalState: Object, validationResult: Object, history: Array<Object>}>}
 */
export async function runRalphLoop(options = {}) {
  const {
    sessionRunner = runAgentSession,
    validator = runQAPipeline,
    iterationRunner = null,
    buildPrompt = buildReviewPrompt,
    systemInstructions = DEFAULT_REVIEW_SYSTEM_INSTRUCTIONS,
    validationCommands = [
      "node .gemini/skills/code-validation/scripts/lint.js all",
      "node .gemini/skills/code-validation/scripts/test.js all",
    ],
    model = null,
    maxIterations = 15,
    maxTurnsPerSession = 20,
    noProgressLimit = 3,
    maxFileSizeBytes = 30 * 1024,
    maxTotalContextBytes = 150 * 1024,
    repoRoot = null,
    planPath = null,
    requireTopic = true,
    validateIntent = false,
    blockTools = ["invoke_agent", "enter_plan_mode"],
    customTools = [],
    rl = null,
    stageBeforeIteration = true,
    requireInitialDiff = false,
    checkCompletion = null,
    captureState = null,
    detectChange = null,
    recaptureStateAfterValidation = false,
    logger: inputLogger = console.log,
    onIterationStart = null,
    onIterationEnd = null,
  } = options;

  const logger = typeof inputLogger === "function" ? inputLogger : console.log;
  const root = repoRoot || (await getRepoRoot());
  const resolvedPlanPath = planPath || path.join(root, "plans/current.md");
  const captureFn = captureState || captureGitState;
  const detectFn = detectChange || detectGitStateChange;

  let iteration = 0;
  let consecutiveNoProgress = 0;
  let lastValidation = null;
  let lastAgentOutput = "";
  let isCompleted = false;
  let completionReason = "";
  const history = [];
  const failingDiffHashCounts = new Map();

  // Capture baseline state
  const initialBaseline = await captureFn({
    repoRoot: root,
    planPath: resolvedPlanPath,
    stageBeforeIteration,
    maxFileSizeBytes,
    maxTotalContextBytes,
  });

  if (
    requireInitialDiff &&
    (!initialBaseline.diff || initialBaseline.diff.trim() === "")
  ) {
    let emptyDiffValidation = { success: true };
    if (typeof validator === "function") {
      emptyDiffValidation = await validator({
        iteration: 0,
        maxIterations,
        stateBefore: initialBaseline,
        stateAfter: initialBaseline,
        stateChanged: false,
        repoRoot: root,
        planPath: resolvedPlanPath,
      });
      if (!emptyDiffValidation?.success) {
        lastValidation = await formatValidationSpillover(
          emptyDiffValidation,
          root,
        );
        logger(
          "\n⚠️ No staged changes found to review, but deterministic workspace validation failed. Entering loop to resolve issues...",
        );
      }
    }

    if (emptyDiffValidation?.success) {
      logger("\n✅ No staged changes found to review. Ralph loop complete.");
      return {
        success: true,
        completed: true,
        iterations: 0,
        reason: "No staged changes found to review.",
        finalState: initialBaseline,
        validationResult: emptyDiffValidation,
        history: [],
      };
    }
  }

  let currentStateBefore = initialBaseline;
  let sessionError = null;

  while (iteration < maxIterations && !isCompleted) {
    iteration++;
    sessionError = null;
    logger(`\n🔄 [Ralph Loop] Iteration ${iteration}/${maxIterations}...`);

    if (iteration > 1) {
      currentStateBefore = await captureFn({
        repoRoot: root,
        planPath: resolvedPlanPath,
        stageBeforeIteration,
        maxFileSizeBytes,
        maxTotalContextBytes,
      });
    }

    if (typeof onIterationStart === "function") {
      try {
        await onIterationStart({
          iteration,
          maxIterations,
          stateBefore: currentStateBefore,
          lastValidation,
        });
      } catch (err) {
        logger(`[DEBUG] onIterationStart hook error: ${err.message}`);
      }
    }

    let agentOutput;
    let iterationMeta;

    if (typeof iterationRunner === "function") {
      try {
        iterationMeta = await iterationRunner({
          iteration,
          maxIterations,
          stateBefore: currentStateBefore,
          lastValidation,
          repoRoot: root,
          planPath: resolvedPlanPath,
          sessionRunner,
          rl,
          logger,
        });
        agentOutput = iterationMeta?.agentOutput || "";
        if (iterationMeta?.error) {
          sessionError = iterationMeta.error;
        }
        lastAgentOutput = agentOutput;
      } catch (err) {
        if (isAbortError(err)) {
          throw err;
        }
        logger(
          `❌ [Ralph Loop] Iteration runner threw an error: ${err.message}`,
        );
        sessionError = err;
        agentOutput = `Error: ${err.message}`;
        lastAgentOutput = agentOutput;
      }
    } else {
      let skillsContext;
      try {
        skillsContext = await loadWorkspaceSkillsContext(root);
      } catch {
        skillsContext = "";
      }

      // Build the prompt for this fresh-context session
      const promptText =
        typeof buildPrompt === "function"
          ? await buildPrompt({
            iteration,
            maxIterations,
            stateBefore: currentStateBefore,
            lastValidation,
            lastAgentOutput,
            repoRoot: root,
            planPath: resolvedPlanPath,
            activePlan: currentStateBefore.plan,
            validationCommands,
            skillsContext,
          })
          : String(buildPrompt || "");

      // Resolve system instructions
      let currentInstructions =
        typeof systemInstructions === "function"
          ? await systemInstructions({
            iteration,
            maxIterations,
            stateBefore: currentStateBefore,
            lastValidation,
            skillsContext,
          })
          : String(systemInstructions);

      if (
        skillsContext &&
        !currentInstructions.includes("<available_skills>") &&
        !currentInstructions.includes("Skills-First Mandate")
      ) {
        currentInstructions += `\n\n${skillsContext}`;
      }

      // Run thinking agent with a clean context window
      try {
        agentOutput = await sessionRunner({
          initialPrompt: promptText,
          systemInstructions: currentInstructions,
          requestedModel: model,
          maxTurnsOverride: maxTurnsPerSession,
          requireTopic,
          validateIntent,
          blockTools,
          customTools,
          rl,
        });
      } catch (err) {
        if (isAbortError(err)) {
          throw err;
        }
        logger(`❌ [Ralph Loop] Agent session threw an error: ${err.message}`);
        sessionError = err;
        agentOutput = `Error: ${err.message}`;
      }
      lastAgentOutput = agentOutput;
    }

    // Capture state after agent execution
    let stateAfter = await captureFn({
      repoRoot: root,
      planPath: resolvedPlanPath,
      stageBeforeIteration,
      maxFileSizeBytes,
      maxTotalContextBytes,
    });

    const changeResult = detectFn(currentStateBefore, stateAfter);
    let stateChanged =
      typeof changeResult === "boolean"
        ? changeResult
        : Boolean(changeResult?.hasChanged);

    // Deterministic Backpressure: run machine validation if no session error
    let validationResult = sessionError
      ? {
        success: false,
        isSessionError: true,
        sessionError: sessionError.message || String(sessionError),
        stdout: "",
        stderr: sessionError.message || String(sessionError),
        error: sessionError.message || String(sessionError),
      }
      : { success: true };
    if (!sessionError && typeof validator === "function") {
      logger(`⚙️  Running deterministic validation pipeline...`);
      try {
        validationResult = await validator({
          iteration,
          repoRoot: root,
          state: stateAfter,
          stateBefore: currentStateBefore,
          stateChanged,
        });
      } catch (err) {
        validationResult = {
          success: false,
          stdout: "",
          stderr: err.message || String(err),
          error: err.message || String(err),
        };
      }

      if (recaptureStateAfterValidation) {
        const postValidationState = await captureFn({
          repoRoot: root,
          planPath: resolvedPlanPath,
          stageBeforeIteration,
          maxFileSizeBytes,
          maxTotalContextBytes,
        });

        const postValidationChange = detectFn(
          currentStateBefore,
          postValidationState,
        );
        const postValidationChanged =
          typeof postValidationChange === "boolean"
            ? postValidationChange
            : Boolean(postValidationChange?.hasChanged);

        stateAfter = postValidationState;
        stateChanged = postValidationChanged;
      }
    }

    if (!validationResult.success) {
      validationResult = await formatValidationSpillover(
        validationResult,
        root,
      );
    }
    lastValidation = validationResult;

    // Evaluate loop exit criteria
    if (iterationMeta?.forceComplete) {
      if (validationResult.success) {
        isCompleted = true;
        completionReason =
          iterationMeta.completionReason || "Completed by iteration runner.";
      } else {
        logger(
          `⚠️ Iteration runner requested early completion, but deterministic validation failed. Retrying with error feedback...`,
        );
      }
    }

    if (!isCompleted) {
      if (sessionError) {
        consecutiveNoProgress++;
        logger(
          `⚠️ Agent session encountered an error (No-progress: ${consecutiveNoProgress}/${noProgressLimit}): ${sessionError.message}`,
        );
      } else if (typeof checkCompletion === "function") {
        const customCheck = await checkCompletion({
          iteration,
          maxIterations,
          stateBefore: currentStateBefore,
          stateAfter,
          stateChanged,
          validationResult,
          agentOutput,
        });
        if (customCheck?.isComplete) {
          isCompleted = true;
          completionReason =
            customCheck.reason || "Custom completion check satisfied.";
        }
      } else {
        // Default Ralph Loop exit criteria
        if (!validationResult.success) {
          if (!stateChanged) {
            consecutiveNoProgress++;
            logger(
              `⚠️ Agent made no changes, but deterministic validation failed (No-progress: ${consecutiveNoProgress}/${noProgressLimit}).`,
            );
          } else {
            // Check if this failing diff hash has been encountered previously
            const diffHash = stateAfter?.diffHash || "";
            const occurrences = (failingDiffHashCounts.get(diffHash) || 0) + 1;
            failingDiffHashCounts.set(diffHash, occurrences);

            if (occurrences > 1) {
              consecutiveNoProgress++;
              logger(
                `⚠️ Detected recurring failing diff state (Hash repeated ${occurrences} times; No-progress: ${consecutiveNoProgress}/${noProgressLimit}).`,
              );
            } else {
              consecutiveNoProgress = 0;
              logger(
                `📝 Agent made changes, but deterministic validation failed. Retrying with error feedback...`,
              );
            }
          }
        } else {
          // Validation passed cleanly
          failingDiffHashCounts.clear();
          if (!stateChanged) {
            if (iterationMeta?.shouldImplement) {
              consecutiveNoProgress++;
              logger(
                `⚠️ Implementation was requested for suggestions, but no repository changes were made (No-progress: ${consecutiveNoProgress}/${noProgressLimit}).`,
              );
            } else {
              isCompleted = true;
              completionReason =
                "No further changes suggested and deterministic validation passed cleanly.";
              logger(
                `\n✅ Convergence achieved: No further changes suggested and deterministic checks passed.`,
              );
            }
          } else {
            consecutiveNoProgress = 0;
            if (
              agentOutput &&
              /<promise>\s*COMPLETE\s*<\/promise>/i.test(agentOutput)
            ) {
              isCompleted = true;
              completionReason =
                "Agent signaled completion promise and deterministic validation passed.";
              logger(
                `\n✅ Agent signaled completion promise and all deterministic checks passed.`,
              );
            } else {
              logger(
                `\n📝 Agent made valid improvements. Looping again to verify stability...`,
              );
            }
          }
        }
      }
    }

    // Circuit breaker: abort if no progress limit is reached
    if (!isCompleted && consecutiveNoProgress >= noProgressLimit) {
      if (sessionError) {
        completionReason = `Agent session failed with error for ${consecutiveNoProgress} consecutive iteration(s): ${sessionError.message}`;
      } else {
        completionReason = `No progress made for ${consecutiveNoProgress} consecutive iterations while validation failed.`;
      }
      logger(`\n🛑 Halting Ralph Loop: ${completionReason}`);
      break;
    }

    history.push({
      iteration,
      stateChanged,
      validationSuccess: validationResult.success,
      sessionError: sessionError ? sessionError.message : null,
      completionReason,
      isCompleted,
    });

    if (typeof onIterationEnd === "function") {
      try {
        await onIterationEnd({
          iteration,
          maxIterations,
          stateChanged,
          validationResult,
          isCompleted,
          reason: completionReason,
        });
      } catch (err) {
        logger(`[DEBUG] onIterationEnd hook error: ${err.message}`);
      }
    }
  }

  if (!isCompleted && iteration >= maxIterations && !completionReason) {
    completionReason = `Reached maximum iterations (${maxIterations}) without resolving all issues.`;
    logger(`\n⚠️ ${completionReason}`);
  }

  const finalCapture = await captureFn({
    repoRoot: root,
    planPath: resolvedPlanPath,
    stageBeforeIteration: false,
  });

  const overallSuccess =
    isCompleted && (!lastValidation || lastValidation.success) && !sessionError;

  return {
    success: overallSuccess,
    completed: isCompleted,
    iterations: iteration,
    reason: completionReason,
    finalState: finalCapture,
    validationResult: lastValidation,
    history,
  };
}

/**
 * Calculates scaled turns for implementation sessions based on suggestion volume.
 *
 * @param {number} [baseTurns=20] - Base turn capacity
 * @param {number} [suggestionCount=5] - Number of suggestions to address
 * @param {number} [maxCeiling=60] - Maximum turn ceiling
 * @returns {number} Scaled turns
 */
export function calculateScaledTurns(
  baseTurns = 20,
  suggestionCount = 5,
  maxCeiling = 60,
) {
  const count = Math.max(0, suggestionCount ?? 5);
  const effectiveBase = Math.min(baseTurns, maxCeiling);
  const scaled = Math.round(effectiveBase * Math.max(1, count / 5));
  return Math.min(maxCeiling, Math.max(effectiveBase, scaled));
}

export const DEFAULT_IMPLEMENTER_INSTRUCTIONS =
  "You are an expert engineer in the IMPLEMENTATION phase of the review loop. Pick and implement UP TO 3 SUGGESTIONS from the remaining list. Ensure all tests and linters pass cleanly.";

/**
 * Executes an inner Ralph Loop for the implementation phase of the review loop.
 * Iterates through suggestions one by one with fresh context per iteration,
 * running automated validation backpressure, removing completed suggestions,
 * and staging the code before returning to the outer review loop.
 *
 * @param {Object} [options]
 * @returns {Promise<{success: boolean, completed: boolean, totalSuggestions: number, implementedCount: number, iterations: number}>}
 */
export async function runImplementationRalphLoop(options = {}) {
  const {
    suggestions = [],
    suggestionsText = "",
    sessionRunner = runAgentSession,
    model = "gemini-3.8-flash",
    proModel = "gemini-3.1-pro-preview",
    proImplementationIterations = 1,
    validator = () => runQAPipeline({ fix: false }),
    validationCommands = [
      "node .gemini/skills/code-validation/scripts/lint.js all",
      "node .gemini/skills/code-validation/scripts/test.js all",
    ],
    repoRoot = null,
    planPath = null,
    skillsContext = "",
    captureState = null,
    detectChange = null,
    stageFn = null,
    rl = null,
    logger: inputLogger = console.log,
    maxTurnsPerSession = 20,
    maxIterations = 15,
    customTools = [],
    blockTools = ["invoke_agent", "enter_plan_mode"],
    implementerInstructions = DEFAULT_IMPLEMENTER_INSTRUCTIONS,
  } = options;

  const logger = typeof inputLogger === "function" ? inputLogger : console.log;
  const root = repoRoot || (await getRepoRoot());
  const resolvedPlanPath = planPath || path.join(root, "plans/current.md");
  const captureFn = captureState || captureGitState;
  const detectFn = detectChange || detectGitStateChange;

  // Parse suggestions into a list if not already parsed
  const parsedItems =
    suggestions.length > 0
      ? [...suggestions]
      : parseIndividualSuggestions(suggestionsText);
  const remaining =
    parsedItems.length > 0
      ? parsedItems
      : typeof suggestionsText === "string" && suggestionsText.trim()
        ? [
          {
            id: 1,
            severity: "critical",
            text: suggestionsText.trim(),
            raw: suggestionsText.trim(),
          },
        ]
        : [];

  const totalSuggestions = remaining.length;
  if (totalSuggestions === 0) {
    logger(
      `\nℹ️ [Inner Implementation Ralph Loop] No suggestions provided to implement.`,
    );
    return {
      success: true,
      completed: true,
      totalSuggestions: 0,
      implementedCount: 0,
      iterations: 0,
    };
  }

  let implementedCount = 0;
  let innerIteration = 0;
  let lastValidation = null;

  logger(
    `\n🔄 [Inner Implementation Ralph Loop] Initiating incremental implementation of ${totalSuggestions} suggestion(s)...`,
  );

  // maxIterations regular attempts + 1 final exhaustion fallback attempt using the Pro thinking model
  const maxTotalAttempts = maxIterations + 1;

  while (remaining.length > 0 && innerIteration < maxTotalAttempts) {
    innerIteration++;

    // If we've hit the max iteration limit without fully resolving all suggestions,
    // we use the final innerIteration as a last-ditch effort with the pro thinking model
    // before giving up.
    const isExhaustionFallback = innerIteration > maxIterations;
    if (isExhaustionFallback) {
      logger(
        `\n⚠️ [Inner Loop] Iteration ceiling reached. Escalating remaining ${remaining.length} suggestion(s) to Thinking Model (${proModel}) for final fallback attempt...`,
      );
    }

    const currentSuggestion = remaining[0];
    const hasCritical = currentSuggestion?.severity === "critical";
    const isProTier = innerIteration <= proImplementationIterations;
    const currentModel =
      isProTier || hasCritical || isExhaustionFallback ? proModel : model;
    const maxBatchSize = 3;
    const batchSize = Math.min(remaining.length, maxBatchSize);
    const batchedSuggestions = remaining.slice(0, batchSize);
    const batchSummary = batchedSuggestions.map((s) => `#${s.id}`).join(", ");

    logger(
      `\n🛠️ [Inner Loop] Iteration ${Math.min(innerIteration, maxIterations)}/${maxIterations}: Waiting for implementation of suggestion(s) ${batchSummary} (e.g., ID #${currentSuggestion.id} [${currentSuggestion.severity}]: "${currentSuggestion.text.slice(0, 80)}") [${implementedCount}/${totalSuggestions} completed]...`,
    );

    // Capture per-attempt baseline state immediately before this turn
    const attemptStateBefore = await captureFn({
      repoRoot: root,
      planPath: resolvedPlanPath,
      stageBeforeIteration: false,
    });

    const prompt = buildBatchSuggestionImplementationPrompt({
      remainingSuggestions: remaining,
      activePlan: attemptStateBefore.plan,
      stateBefore: attemptStateBefore,
      lastValidation,
      currentSuggestionIndex: implementedCount + 1,
      totalSuggestions,
      validationCommands,
      skillsContext,
    });

    const instructions = implementerInstructions;

    let agentOutput;
    let sessionError = null;

    try {
      agentOutput = await sessionRunner({
        initialPrompt: instructions + "\n\n" + prompt,
        requestedModel: currentModel,
        maxTurnsOverride: calculateScaledTurns(
          maxTurnsPerSession,
          remaining.length,
        ),
        requireTopic: true,
        blockTools,
        customTools,
        rl,
      });
    } catch (err) {
      if (isAbortError(err)) {
        throw err;
      }
      logger(`❌ [Inner Loop] Implementer session threw error: ${err.message}`);
      sessionError = err;
      agentOutput = `Error: ${err.message}`;
    }

    const attemptedIds = extractCompletedSuggestionIds(
      agentOutput,
      remaining,
      logger,
    );

    // Capture state after implementer execution
    const stateAfter = await captureFn({
      repoRoot: root,
      planPath: resolvedPlanPath,
      stageBeforeIteration: false,
    });
    const changeResult = detectFn(attemptStateBefore, stateAfter);
    const stateChanged =
      typeof changeResult === "boolean"
        ? changeResult
        : Boolean(changeResult?.hasChanged);

    // Deterministic Backpressure: validate after each suggestion
    let validationResult = sessionError
      ? { success: false, stdout: "", stderr: sessionError.message }
      : { success: true };

    if (!sessionError && typeof validator === "function") {
      logger(
        `⚙️  Running deterministic validation pipeline for suggestion ID(s) ${attemptedIds.join(", ")}...`,
      );
      try {
        validationResult = await validator({
          iteration: innerIteration,
          repoRoot: root,
          state: stateAfter,
          stateBefore: attemptStateBefore,
          stateChanged,
        });
      } catch (err) {
        validationResult = {
          success: false,
          stdout: "",
          stderr: err.message || String(err),
        };
      }
    }

    if (!validationResult.success) {
      validationResult = await formatValidationSpillover(
        validationResult,
        root,
      );
    }

    lastValidation = validationResult;

    if (validationResult.success) {
      if (attemptedIds.length > 0) {
        for (const attemptedId of attemptedIds) {
          const completedIdx = remaining.findIndex((s) => s.id === attemptedId);
          if (completedIdx !== -1) {
            const removed = remaining.splice(completedIdx, 1)[0];
            implementedCount++;
            logger(
              `✅ [Inner Loop] Successfully implemented and verified suggestion #${removed.id} [${removed.severity}]. (${remaining.length} remaining)`,
            );
          }
        }
      } else if (stateChanged && remaining.length > 0) {
        // Fallback only if the repository state actually changed during this attempt and validation passed
        const removed = remaining.shift();
        implementedCount++;
        logger(
          `✅ [Inner Loop] Verified repository changes and deterministic validation for suggestion #${removed.id} [${removed.severity}]. (${remaining.length} remaining)`,
        );
      } else {
        logger(
          `⚠️ [Inner Loop] Implementer emitted no completion ID and made no code changes. Suggestion #${currentSuggestion.id} remains unaddressed.`,
        );
      }
    } else {
      logger(
        `⚠️ [Inner Loop] Validation failed for suggestion ID(s) ${attemptedIds.join(", ")}. Retrying with automated backpressure...`,
      );
    }
  }

  // Only stage code if all suggestions were successfully implemented and last validation passed
  if (remaining.length === 0 && (!lastValidation || lastValidation.success)) {
    logger(`\n📦 [Inner Loop] Staging all implemented changes...`);
    if (typeof stageFn === "function") {
      await stageFn({ repoRoot: root });
    } else {
      try {
        const stageScript = path.join(
          root,
          ".gemini/skills/git/scripts/stage-files.js",
        );
        await execFileAsync("node", [stageScript, "--all"], { cwd: root });
        logger(`✅ Successfully staged code via git skill stage-files.js.`);
      } catch (err) {
        logger(`⚠️ Could not stage files via skill: ${err.message}`);
      }
    }
  } else {
    logger(
      `\n⚠️ [Inner Loop] Skipping staging: loop finished with ${remaining.length} remaining suggestion(s) and validation success=${lastValidation?.success ?? "n/a"}.`,
    );
  }

  return {
    success: remaining.length === 0,
    completed: remaining.length === 0,
    totalSuggestions,
    implementedCount,
    iterations: innerIteration,
  };
}

export const DEFAULT_REVIEWER_INSTRUCTIONS =
  "You are a senior technical reviewer. Review the staged code against the active plan in plans/current.md. Strive for a minimum of 5 suggestions with severity ratings [critical], [medium], or [minor] (including minor or trivial enhancements). If after striving for 5 you cannot find 5, provide whatever valid suggestions you identified. If the staged code meets all requirements and no changes or enhancements are needed, output: []. Do NOT modify files or call tools.";

/**
 * Specialized high-level wrapper executing a Multi-Agent Review Ralph Loop.
 * 1. Thinking model (Pro) reviews staged code and generates at least 5 rated suggestions (critical/medium/minor).
 * 2. Flash-Lite model reads and tallies the counts of each category.
 * 3. Decision gate: if critical > 0 or medium > 0, Flash model implements all suggestions; if only minor, implementation is skipped and convergence is achieved.
 * 4. Deterministic validation pipeline runs after implementation to verify test/linter compliance.
 *
 * @param {Object} [options]
 * @param {string} [options.proModel='gemini-3.1-pro-preview'] - Thinking model for reviewing
 * @param {string} [options.flashModel='gemini-3.8-flash'] - Flash model for implementing suggestions
 * @param {string} [options.flashLiteModel='gemini-3.1-flash-lite'] - Flash-Lite model for tallying counts
 * @param {Object} [options.models] - Models dictionary containing { pro, flash, flash_lite }
 * @returns {Promise<{success: boolean, completed: boolean, iterations: number, reason: string, finalState: Object, validationResult: Object, history: Array<Object>}>}
 */
export async function runReviewRalphLoop(options = {}) {
  const {
    proModel = "gemini-3.1-pro-preview",
    flashModel = "gemini-3.8-flash",
    flashLiteModel = "gemini-3.1-flash-lite",
    models: modelsConfig = null,
    sessionRunner = runAgentSession,
    validator = runQAPipeline,
    repoRoot = null,
    planPath = null,
    customTools = [],
    stageFn = null,
    rl = null,
    logger: inputLogger = console.log,
    maxIterations = 15,
    maxTurnsPerSession = 20,
    ...restOptions
  } = options;

  const logger = typeof inputLogger === "function" ? inputLogger : console.log;

  const resolvedProModel = modelsConfig?.pro || proModel;
  const resolvedFlashModel = modelsConfig?.flash || options.model || flashModel;
  const resolvedFlashLiteModel =
    modelsConfig?.flash_lite || modelsConfig?.flashLite || flashLiteModel;
  const proReviewIterations = options.proReviewIterations ?? 3;
  const proImplementationIterations = options.proImplementationIterations ?? 1;

  const root = repoRoot || (await getRepoRoot());
  const resolvedPlanPath = planPath || path.join(root, "plans/current.md");
  let skillsContext = "";
  try {
    skillsContext = await loadWorkspaceSkillsContext(root);
  } catch {
    skillsContext = "";
  }

  const reviewIterationRunner = async ({
    iteration,
    maxIterations: totalIterations,
    stateBefore,
    lastValidation,
    logger: iterLogger,
  }) => {
    const activeLogger = typeof iterLogger === "function" ? iterLogger : logger;
    const activeReviewModel =
      iteration <= proReviewIterations ? resolvedProModel : resolvedFlashModel;
    activeLogger(
      `\n🤔 [Review Loop] Step 1: Reviewer (${activeReviewModel}) evaluating staged code (Iteration ${iteration}/${totalIterations})...`,
    );

    const reviewerPrompt = buildReviewerPrompt({
      iteration,
      maxIterations: totalIterations,
      stateBefore,
      lastValidation,
      activePlan: stateBefore.plan,
      skillsContext,
    });

    const reviewerInstructions = DEFAULT_REVIEWER_INSTRUCTIONS;

    // 1. Review staged code (Pro model for first 3 iterations, Flash model for remaining)
    const suggestionsText = await sessionRunner({
      initialPrompt: reviewerInstructions + "\n\n" + reviewerPrompt,
      requestedModel: activeReviewModel,
      maxTurnsOverride: 1,
      requireTopic: false,
      isolate: true,
      allowedTools: [],
      rl,
    });

    const trimmed =
      typeof suggestionsText === "string" ? suggestionsText.trim() : "";

    // If reviewer returned an empty list '[]' or text explicitly stating 0 suggestions / no suggestions
    const isCleanApproval =
      trimmed === "[]" ||
      trimmed === "```json\n[]\n```" ||
      trimmed === "```\n[]\n```" ||
      /^(?:\[\]|no\s+(?:critical|medium|actionable\s+)?suggestions(?:\s+found)?\.?|none\.?|pristine\.?)$/i.test(
        trimmed,
      );

    if (isCleanApproval) {
      const tally = { critical: 0, medium: 0, minor: 0, total: 0 };
      const completionReason =
        "Reviewer reported 0 suggestions ([]). Code review approved cleanly.";
      activeLogger(`\n✅ ${completionReason}`);
      return {
        agentOutput: suggestionsText,
        shouldImplement: false,
        tally,
        suggestionsText,
        forceComplete: true,
        completionReason,
      };
    }

    if (!trimmed) {
      const errMsg =
        "Reviewer model returned an empty response. Expected a numbered list of suggestions or [].";
      activeLogger(`\n⚠️ ${errMsg} Retrying review iteration...`);
      return {
        agentOutput: "",
        error: new Error(errMsg),
        shouldImplement: false,
      };
    }

    activeLogger(
      `\n📝 [Review Loop] Step 2: Flash-Lite (${resolvedFlashLiteModel}) tallying suggestion categories...`,
    );

    // 2. Flash-Lite model counts suggestions by category
    const tally = await tallySuggestionsWithFlashLite({
      suggestionsText,
      sessionRunner,
      model: resolvedFlashLiteModel,
      rl,
      logger: activeLogger,
    });

    activeLogger(
      `📊 Suggestion Tally: ${tally.critical} critical, ${tally.medium} medium, ${tally.minor} minor (Total: ${tally.total})`,
    );

    // 3. Decision Gate:
    // If there are any critical or medium suggestions we implement all suggestions.
    // If only minor suggestions exist (critical and medium are 0), we do not need to resolve them,
    // skipping implementation to achieve convergence.
    const shouldImplement =
      (tally.critical || 0) > 0 || (tally.medium || 0) > 0;

    if (!shouldImplement) {
      if ((tally.total || 0) === 0) {
        const errMsg =
          "Reviewer returned unparseable output without categorized suggestions. Retrying review iteration...";
        activeLogger(`\n⚠️ ${errMsg}`);
        return {
          agentOutput: suggestionsText,
          error: new Error(errMsg),
          shouldImplement: false,
        };
      }

      const completionReason = `Only minor suggestions found (${tally.minor} minor, 0 critical, 0 medium). No implementation needed.`;
      activeLogger(
        `\n✅ No critical or medium suggestions found (only ${tally.minor} minor suggestions). Skipping implementation per review policy.`,
      );
      return {
        agentOutput: suggestionsText,
        shouldImplement: false,
        tally,
        suggestionsText,
        forceComplete: true,
        completionReason,
      };
    }

    activeLogger(
      `\n🛠️ [Review Loop] Step 3: Critical or medium suggestions detected. Initiating inner implementation Ralph Loop...`,
    );

    // 4. Run the inner Ralph Loop to implement suggestions (Pro model for iteration 1, Flash model for subsequent retries)
    const innerResult = await runImplementationRalphLoop({
      suggestionsText,
      sessionRunner,
      model: resolvedFlashModel,
      proModel: resolvedProModel,
      proImplementationIterations,
      validator,
      validationCommands: restOptions.validationCommands,
      repoRoot: root,
      planPath: resolvedPlanPath,
      skillsContext,
      captureState: restOptions.captureState || captureGitState,
      detectChange: restOptions.detectChange,
      stageFn: stageFn || restOptions.stageFn,
      rl,
      logger: activeLogger,
      maxTurnsPerSession,
      maxIterations: Math.max(15, tally.total * 3),
      customTools: restOptions.customTools || customTools,
      blockTools: restOptions.blockTools || ["invoke_agent", "enter_plan_mode"],
      implementerInstructions:
        restOptions.implementerInstructions || DEFAULT_IMPLEMENTER_INSTRUCTIONS,
    });

    const isSuccess = Boolean(innerResult?.success);
    const failureMsg = !isSuccess
      ? `Inner implementation loop failed: only ${innerResult.implementedCount}/${innerResult.totalSuggestions} suggestions implemented.`
      : null;

    return {
      agentOutput: `Inner Implementation Loop completed: ${innerResult.implementedCount}/${innerResult.totalSuggestions} suggestions implemented across ${innerResult.iterations} iteration(s).`,
      shouldImplement: true,
      tally,
      suggestionsText,
      innerResult,
      ...(failureMsg ? { error: new Error(failureMsg) } : {}),
    };
  };

  const requireInitialDiff =
    options.requireInitialDiff !== undefined
      ? options.requireInitialDiff
      : true;

  return runRalphLoop({
    iterationRunner: reviewIterationRunner,
    validator,
    repoRoot,
    planPath,
    customTools,
    rl,
    logger,
    maxIterations,
    maxTurnsPerSession,
    ...restOptions,
    requireInitialDiff,
  });
}
