import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

/**
 * Parses git diff-tree stdout into an array of clean file paths.
 *
 * @param {string} stdout
 * @returns {Array<string>}
 */
export function parseDiffTreeOutput(stdout) {
  if (!stdout || typeof stdout !== "string") {
    return [];
  }
  return stdout
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean);
}

/**
 * Validates whether a commit SHA is a non-empty alphanumeric string.
 *
 * @param {string} commitSha
 * @param {boolean} [throwOnError=false] - If true, throws an Error with the rejected value instead of returning false
 * @returns {boolean}
 */
export function isValidCommitSha(commitSha, throwOnError = false) {
  const valid = Boolean(
    commitSha &&
      typeof commitSha === "string" &&
      /^[a-zA-Z0-9]+$/.test(commitSha.trim()),
  );
  if (!valid && throwOnError) {
    throw new Error(`Invalid commit SHA: ${commitSha}`);
  }
  return valid;
}

/**
 * Validates a commit SHA and throws an Error with the rejected value if it is not valid.
 *
 * @param {string} commitSha
 * @throws {Error} If commitSha is not a valid commit SHA
 */
export function validateCommitSha(commitSha) {
  isValidCommitSha(commitSha, true);
}

/**
 * Inspects changed files for a given commit SHA locally via git diff-tree.
 * Returns null if the commit SHA is not found locally or command fails.
 *
 * @param {string} commitSha
 * @param {object} [options]
 * @param {string} [options.cwd=process.cwd()]
 * @param {Function} [options.execFn] - For testing injection
 * @returns {Promise<Array<string>|null>}
 */
export async function getChangedFilesFromCommit(
  commitSha,
  { cwd = process.cwd(), execFn = execFileAsync } = {},
) {
  if (!isValidCommitSha(commitSha)) {
    return null;
  }

  try {
    const { stdout } = await execFn(
      "git",
      [
        "diff-tree",
        "--no-commit-id",
        "--name-only",
        "-r",
        "--",
        commitSha.trim(),
      ],
      { cwd },
    );
    return parseDiffTreeOutput(stdout);
  } catch {
    return null;
  }
}

export const getChangedFilesForCommit = getChangedFilesFromCommit;
