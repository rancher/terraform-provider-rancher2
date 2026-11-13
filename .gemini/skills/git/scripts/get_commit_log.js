#!/usr/bin/env node
import { simpleGit } from "simple-git";
/**
 * Gets concise commit log for current branch or specific file/branch.
 * @param {object} [params]
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {number} [params.limit=5]
 * @param {string} [params.file]
 * @param {string} [params.branch]
 * @returns {Promise<Array<{hash: string, message: string, date: string}>>}
 */
export async function getCommitLog({
  git = simpleGit(),
  limit = 5,
  file,
  branch,
} = {}) {
  try {
    const options = { maxCount: limit };
    if (file) {
      options.file = file;
    }
    if (branch) {
      options.from = branch;
    }
    const log = await git.log(options);
    return log.all.map((c) => ({
      hash: c.hash.slice(0, 7),
      message: c.message,
      date: c.date,
    }));
  } catch {
    return [];
  }
}

/**
 * Parses full CLI options for get_commit_log.js.
 * Supports:
 * - [limit]
 * - [-n, limit] or [--limit, limit]
 * - [-f, file] or [--file, file]
 * - [-b, branch] or [--branch, branch]
 * - [file_path] or [branch]
 * @param {Array<string>} [args]
 * @returns {{limit: number, file?: string, branch?: string}}
 */
export function parseCommitLogOptions(args = []) {
  let limit = 5;
  let file = undefined;
  let branch = undefined;

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if ((arg === "-n" || arg === "--limit") && i + 1 < args.length) {
      const parsed = parseInt(args[++i], 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (arg.startsWith("-n=") || arg.startsWith("--limit=")) {
      const parsed = parseInt(arg.slice(arg.indexOf("=") + 1), 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if ((arg === "-f" || arg === "--file") && i + 1 < args.length) {
      file = args[++i];
    } else if (arg.startsWith("-f=") || arg.startsWith("--file=")) {
      file = arg.slice(arg.indexOf("=") + 1);
    } else if ((arg === "-b" || arg === "--branch") && i + 1 < args.length) {
      branch = args[++i];
    } else if (arg.startsWith("-b=") || arg.startsWith("--branch=")) {
      branch = arg.slice(arg.indexOf("=") + 1);
    } else if (/^\d+$/.test(arg)) {
      const parsed = parseInt(arg, 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (!file && (arg.includes("/") || arg.includes("."))) {
      file = arg;
    } else if (!branch && !arg.startsWith("-")) {
      branch = arg;
    }
  }

  return { limit, file, branch };
}

/**
 * Parses CLI arguments for get_commit_log.js.
 * Supports:
 * - [limit]
 * - [-n, limit]
 * - [--limit, limit]
 * - [-n=limit]
 * - [--limit=limit]
 * @param {Array<string>} args
 * @returns {number}
 */
export function parseCommitLogArgs(args = []) {
  return parseCommitLogOptions(args).limit;
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_commit_log.js [limit] [-n limit] [-f file] [-b branch]

Outputs a concise commit log for the current branch, target branch, or specific file.

Options:
  -n, --limit <count>  Number of commits to display (default: 5)
  -f, --file <path>    Filter commits affecting specific file
  -b, --branch <name>  Show commit history for specific branch
  -h, --help           Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }

  const { limit, file, branch } = parseCommitLogOptions(process.argv.slice(2));

  try {
    const commits = await getCommitLog({ limit, file, branch });
    if (commits.length === 0) {
      console.log("No commits found on current branch.");
      return;
    }

    const scopeDesc = file
      ? ` for file '${file}'`
      : branch
        ? ` for branch '${branch}'`
        : "";
    console.log(`=== Last ${commits.length} Commit(s)${scopeDesc} ===`);
    for (const c of commits) {
      console.log(`${c.hash} ${c.message}`);
    }
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
