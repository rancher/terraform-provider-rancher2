import { exec, execFile } from "node:child_process";
import path from "node:path";
import { promisify } from "node:util";
import {
  fileAffectsProduct,
  getProductRoots,
  validateCommitTitle,
} from "../../skills/git/scripts/commit.js";
import { getRepoRoot } from "./utils.js";

export const sys = {
  execAsync: promisify(exec),
  execFileAsync: promisify(execFile),
};

/**
 * Gets the current unstaged and staged git diff.
 * @returns {Promise<string>} The git diff
 */
export async function getGitDiff() {
  const repoRoot = await getRepoRoot();
  const { stdout } = await sys.execAsync("git diff --no-ext-diff HEAD", {
    cwd: repoRoot,
    maxBuffer: 10 * 1024 * 1024,
  });
  return stdout;
}

/**
 * Checks if the cached (staged) or active changes affect product files.
 * @returns {Promise<boolean>}
 */
export async function affectsProductFiles() {
  const repoRoot = await getRepoRoot();
  const [{ stdout: diffFiles }, { stdout: untrackedFiles }] = await Promise.all([
    sys.execAsync("git diff -z --no-ext-diff --name-only HEAD", {
      cwd: repoRoot,
      maxBuffer: 10 * 1024 * 1024,
    }),
    sys.execAsync("git ls-files --others --exclude-standard -z", {
      cwd: repoRoot,
      maxBuffer: 10 * 1024 * 1024,
    }),
  ]);
  const changedFiles = [
    ...diffFiles.split("\0"),
    ...untrackedFiles.split("\0"),
  ]
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
  const { stdout: branchOutput } = await sys.execAsync(
    "git rev-parse --abbrev-ref HEAD",
    {
      cwd: repoRoot,
      maxBuffer: 10 * 1024 * 1024,
    },
  );
  if (["main", "master"].includes(branchOutput.trim())) {
    throw new Error('Commits directly to the "main" or "master" branch are not allowed.');
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
      maxBuffer: 10 * 1024 * 1024,
    },
  );
  const modifiedPaths = pathsOutput
    .split("\n")
    .map((p) => p.trim())
    .filter(Boolean);

  if (modifiedPaths.length === 0) {
    throw new Error("No modified or untracked paths found to commit.");
  }

  // Stage files in chunks to prevent E2BIG (Argument list too long) OS errors
  const chunkSize = 100;
  const stageScriptPath = path.join(
    repoRoot,
    ".gemini/skills/git/scripts/stage-files.js",
  );
  for (let i = 0; i < modifiedPaths.length; i += chunkSize) {
    const chunk = modifiedPaths.slice(i, i + chunkSize);
    await sys.execFileAsync(process.execPath, [stageScriptPath, "--", ...chunk], {
      cwd: repoRoot,
      maxBuffer: 10 * 1024 * 1024,
    });
  }

  // Verify the message before committing
  const validation = await validateMessage(commitMessage);
  if (!validation.valid) {
    throw new Error(`Commit message validation failed: ${validation.reason}`);
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
      maxBuffer: 10 * 1024 * 1024,
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
