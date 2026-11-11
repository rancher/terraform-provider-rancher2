import { spawn, exec, execFile } from "node:child_process";
import path from "node:path";
import { promisify } from "node:util";
import {
  fileAffectsProduct,
  getProductRoots,
  validateCommitTitle,
} from "../../skills/git/scripts/commit.js";
import { logError } from "./logger.js";
import { getRepoRoot } from "./utils.js";

export const sys = {
  execAsync: promisify(exec),
  execFileAsync: promisify(execFile),
  spawnAsync: (command, args, options = {}) => {
    return new Promise((resolve, reject) => {
      const spawnOpts = { ...options };
      const timeoutMs = spawnOpts.timeout || 30000;
      if (!spawnOpts.signal && typeof AbortSignal !== "undefined" && AbortSignal.timeout) {
        spawnOpts.signal = AbortSignal.timeout(timeoutMs);
      }
      
      const child = spawn(command, args, spawnOpts);
      let stdout = "";
      let stderr = "";
      if (child.stdout) {
        child.stdout.setEncoding("utf8");
        child.stdout.on("data", (data) => (stdout += data));
      }
      if (child.stderr) {
        child.stderr.setEncoding("utf8");
        child.stderr.on("data", (data) => (stderr += data));
      }
      child.on("close", (code, signal) => {
        if (code !== 0) {
          const isTimeout = spawnOpts.signal && spawnOpts.signal.aborted && spawnOpts.signal.reason?.name === "TimeoutError";
          const isAborted = spawnOpts.signal && spawnOpts.signal.aborted;
          const errorMsg = isTimeout 
            ? `Command timed out after ${timeoutMs}ms: ${stderr}`
            : isAborted
              ? `Command aborted: ${stderr}`
              : `Command failed with code ${code} (signal ${signal}): ${stderr}`;
          const error = new Error(errorMsg);
          error.code = code;
          error.signal = signal;
          error.stdout = stdout;
          error.stderr = stderr;
          reject(error);
        } else {
          resolve({ stdout, stderr });
        }
      });
      child.on("error", (err) => {
        // AbortError from signal is handled, but we reject it properly
        if (err.name === "AbortError") {
          const isTimeout = spawnOpts.signal && spawnOpts.signal.aborted && spawnOpts.signal.reason?.name === "TimeoutError";
          const error = new Error(isTimeout ? `Command timed out after ${timeoutMs}ms` : "Command aborted");
          error.stdout = stdout;
          error.stderr = stderr;
          reject(error);
        } else {
          reject(err);
        }
      });
    });
  },
};

/**
 * Gets the current unstaged and staged git diff.
 * @returns {Promise<string>} The git diff
 */
export async function getGitDiff() {
  const repoRoot = await getRepoRoot();
  try {
    const { stdout } = await sys.spawnAsync(
      "git",
      ["diff", "--no-ext-diff", "HEAD"],
      {
        cwd: repoRoot,
      },
    );
    return stdout;
  } catch (err) {
    if (err.stderr) logError(`Error executing git diff: ${err.stderr}`);
    throw err;
  }
}

/**
 * Checks if the cached (staged) or active changes affect product files.
 * @returns {Promise<boolean>}
 */
export async function affectsProductFiles() {
  const repoRoot = await getRepoRoot();
  let diffFiles = "";
  let untrackedFiles = "";
  try {
    const [{ stdout: dFiles }, { stdout: uFiles }] = await Promise.all([
      sys.spawnAsync(
        "git",
        ["diff", "-z", "--no-ext-diff", "--name-only", "HEAD"],
        {
          cwd: repoRoot,
        },
      ),
      sys.spawnAsync(
        "git",
        ["ls-files", "--others", "--exclude-standard", "-z"],
        {
          cwd: repoRoot,
        },
      ),
    ]);
    diffFiles = dFiles;
    untrackedFiles = uFiles;
  } catch (err) {
    if (err.stderr)
      logError(`Error executing git diff/ls-files: ${err.stderr}`);
    throw err;
  }
  const changedFiles = [...diffFiles.split("\0"), ...untrackedFiles.split("\0")]
    .map((f) => f.trim())
    .filter(Boolean);
  const roots = await getProductRoots();
  return changedFiles.some((f) => fileAffectsProduct(f, roots));
}

/**
 * Validates a candidate commit message against the repository rules.
 * @param {string} message - The commit message to validate
 * @returns {Promise<{valid: boolean, reason?: string}>}
 */
export async function validateMessage(message) {
  const isProduct = await affectsProductFiles();
  const title = (message || "").split(/\r?\n/)[0].trim();
  const result = validateCommitTitle(title, isProduct, false);
  return result;
}

/**
 * Stages all changes and executes a signed, validated git commit.
 * @param {string} commitMessage - The message to commit with
 * @returns {Promise<{stdout: string, stderr: string}>} The commit result
 */
export async function stageAndCommit(commitMessage) {
  const repoRoot = await getRepoRoot();

  // Prevent commits directly to the main branch
  let branchOutput = "";
  try {
    const { stdout } = await sys.spawnAsync(
      "git",
      ["rev-parse", "--abbrev-ref", "HEAD"],
      {
        cwd: repoRoot,
      },
    );
    branchOutput = stdout;
  } catch (err) {
    if (err.stderr) logError(`Error executing git rev-parse: ${err.stderr}`);
    throw err;
  }
  if (["main", "master"].includes(branchOutput.trim())) {
    throw new Error(
      'Commits directly to the "main" or "master" branch are not allowed.',
    );
  }

  // Get explicitly modified and untracked files
  const scriptPath = path.join(
    repoRoot,
    ".gemini/skills/git/scripts/get_modified_paths.js",
  );
  const { stdout: pathsOutput } = await sys.execFileAsync(
    process.execPath,
    [scriptPath],
    {
      cwd: repoRoot,
    },
  );
  const modifiedPaths = pathsOutput
    .split("\n")
    .map((p) => p.trim())
    .filter(Boolean);

  if (modifiedPaths.length === 0) {
    throw new Error("No modified or untracked paths found to commit.");
  }

  // Verify the message before staging or committing to prevent leaving index in a mutated state
  const validation = await validateMessage(commitMessage);
  if (!validation.valid) {
    throw new Error(`Commit message validation failed: ${validation.reason}`);
  }

  // Stage files in chunks to prevent E2BIG (Argument list too long) OS errors
  const chunkSize = 100;
  const stageScriptPath = path.join(
    repoRoot,
    ".gemini/skills/git/scripts/stage-files.js",
  );
  for (let i = 0; i < modifiedPaths.length; i += chunkSize) {
    const chunk = modifiedPaths.slice(i, i + chunkSize);
    await sys.execFileAsync(
      process.execPath,
      [stageScriptPath, "--", ...chunk],
      {
        cwd: repoRoot,
      },
    );
  }

  // Perform signed and signed-off commit
  const commitScriptPath = path.join(
    repoRoot,
    ".gemini/skills/git/scripts/commit.js",
  );
  const { stdout, stderr } = await sys.execFileAsync(
    process.execPath,
    [commitScriptPath, "-m", commitMessage],
    {
      cwd: repoRoot,
    },
  );
  return { stdout, stderr };
}

/**
 * Strips volatile diff metadata (index hashes and chunk header line numbers)
 * while preserving diff content whitespace and hunk location context,
 * so that cryptographic diff hashing is resilient to line number shifts.
 * @param {string} diff - The unified git diff string
 * @returns {string} The sanitized diff string
 */
export function stripDiffMetadata(diff) {
  if (typeof diff !== "string" || !diff) {
    return "";
  }
  return diff
    .split(/\r?\n/)
    .filter((line) => !line.startsWith("index "))
    .map((line) =>
      line.replace(/^@@\s+-\d+(?:,\d+)?\s+\+\d+(?:,\d+)?\s+@@/, "@@ @@"),
    )
    .join("\n");
}
