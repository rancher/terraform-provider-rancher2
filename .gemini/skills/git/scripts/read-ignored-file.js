#!/usr/bin/env node
import { execFile } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";
import { simpleGit } from "simple-git";

const execFileAsync = promisify(execFile);

/**
 * Checks if a given path is excluded by .aiexclude policy.
 * @param {string} fullPath - Full absolute path to file
 * @param {string} repoRoot - Root of git repository
 * @returns {Promise<boolean>} True if blocked by .aiexclude
 */
export async function isExcludedByAIExclude(fullPath, repoRoot) {
  const aiexcludePath = path.join(repoRoot, ".aiexclude");
  try {
    await fs.access(aiexcludePath);
  } catch (err) {
    if (err.code === "ENOENT") {
      return false; // .aiexclude file does not exist
    }
    throw new Error(`Failed to access .aiexclude policy: ${err.message}`, {
      cause: err,
    });
  }

  const relPath = path.relative(repoRoot, fullPath);
  try {
    const { stdout } = await execFileAsync("git", [
      "-C",
      repoRoot,
      "-c",
      "core.quotepath=false",
      "ls-files",
      "-o",
      "-c",
      "-i",
      "--exclude-from=.aiexclude",
      "--",
      relPath,
    ]);
    const lines = stdout
      .split("\n")
      .map((l) => l.trim().replace(/^"(.*)"$/, "$1"))
      .filter(Boolean);
    return lines.includes(relPath) || lines.length > 0;
  } catch (err) {
    throw new Error(
      `Failed to evaluate .aiexclude policy for '${relPath}': ${err.message}`,
      { cause: err },
    );
  }
}

/**
 * Reads a file that may be ignored by .gitignore, enforcing .aiexclude rules.
 * @param {object} params
 * @param {string} params.targetPath
 * @param {number} [params.startLine]
 * @param {number} [params.endLine]
 * @param {string} [params.repoRoot]
 * @returns {Promise<string>} Content of requested file or line range
 */
export async function readIgnoredFile({
  targetPath,
  startLine,
  endLine,
  repoRoot,
}) {
  const root =
    repoRoot || (await simpleGit().revparse(["--show-toplevel"])).trim();
  const fullPath = path.isAbsolute(targetPath)
    ? targetPath
    : path.resolve(root, targetPath);

  const relativeFromRoot = path.relative(root, fullPath);
  if (relativeFromRoot.startsWith("..") || path.isAbsolute(relativeFromRoot)) {
    throw new Error(
      `Cannot read file outside repository root: '${targetPath}'.`,
    );
  }

  try {
    const stat = await fs.stat(fullPath);
    if (stat.isDirectory()) {
      throw new Error(`Path '${targetPath}' is a directory, not a file.`);
    }

    const realPath = await fs.realpath(fullPath);
    const realRoot = await fs.realpath(root);
    const realRel = path.relative(realRoot, realPath);
    if (realRel.startsWith("..") || path.isAbsolute(realRel)) {
      throw new Error(
        `Cannot read file outside repository root: '${targetPath}'.`,
      );
    }
  } catch (err) {
    if (err.code === "ENOENT") {
      throw new Error(`File '${targetPath}' not found.`, { cause: err });
    }
    throw err;
  }

  const isBlocked = await isExcludedByAIExclude(fullPath, root);
  if (isBlocked) {
    throw new Error(
      `Access to '${targetPath}' is blocked by .aiexclude policy.`,
    );
  }

  const rawContent = await fs.readFile(fullPath, "utf8");
  if (startLine === undefined && endLine === undefined) {
    return rawContent;
  }

  const parsedStart =
    startLine !== undefined && startLine !== null
      ? Number(startLine)
      : undefined;
  const parsedEnd =
    endLine !== undefined && endLine !== null ? Number(endLine) : undefined;

  if (parsedStart !== undefined && isNaN(parsedStart)) {
    throw new Error(`Invalid startLine: '${startLine}'. Must be a number.`);
  }
  if (parsedEnd !== undefined && isNaN(parsedEnd)) {
    throw new Error(`Invalid endLine: '${endLine}'. Must be a number.`);
  }

  if (parsedStart && parsedEnd && parsedStart > parsedEnd) {
    throw new Error(
      `startLine (${parsedStart}) cannot be greater than endLine (${parsedEnd}).`,
    );
  }

  const lines = rawContent.split("\n");
  const start = parsedStart ? Math.max(1, parsedStart) : 1;
  const end = parsedEnd ? Math.min(lines.length, parsedEnd) : lines.length;

  return lines.slice(start - 1, end).join("\n");
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: read-ignored-file.js <file_path> [start_line] [end_line]

Reads a file that may be ignored by .gitignore, enforcing .aiexclude rules.

Arguments:
  file_path    Path to the target file
  start_line   Starting line number (1-based, optional)
  end_line     Ending line number (1-based, optional)

Options:
  -h, --help   Show this help message
`);
}

async function main() {
  if (process.argv.includes("-h") || process.argv.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const targetPath = process.argv[2];
  const startArg = process.argv[3];
  const endArg = process.argv[4];

  if (!targetPath) {
    console.error(
      "Usage: read-ignored-file.js <file_path> [start_line] [end_line]",
    );
    process.exit(1);
  }

  const startLine = startArg ? parseInt(startArg, 10) : undefined;
  const endLine = endArg ? parseInt(endArg, 10) : undefined;

  try {
    const content = await readIgnoredFile({ targetPath, startLine, endLine });
    process.stdout.write(content.endsWith("\n") ? content : content + "\n");
  } catch (err) {
    console.error(`❌ Error: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
