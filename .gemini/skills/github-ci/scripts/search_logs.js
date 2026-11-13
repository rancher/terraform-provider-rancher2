#!/usr/bin/env node
import { createReadStream } from "node:fs";
import fs from "node:fs/promises";
import path from "node:path";
import readline from "node:readline";
import { logError } from "../../../orchestrator/lib/logger.js";
import { DEFAULT_PATTERN } from "./constants.js";

/**
 * Strips ANSI escape codes from terminal output.
 * @param {string} str
 * @returns {string}
 */
export function stripAnsi(str) {
  return str ? str.replace(/\x1b\[[0-9;]*[a-zA-Z]/g, "") : "";
}

/**
 * Recursively searches a directory for regex matches across text files.
 * @param {object} params
 * @param {string} params.logsDir
 * @param {RegExp|string} [params.pattern]
 * @param {number} [params.limit=100]
 * @returns {Promise<Array<{file: string, line: number, text: string}>>}
 */
export async function searchLogs({
  logsDir,
  pattern = DEFAULT_PATTERN,
  limit = 100,
}) {
  let regex;
  if (pattern instanceof RegExp) {
    regex = new RegExp(pattern.source, pattern.flags.replace(/[gy]/g, ""));
  } else if (typeof pattern === "string") {
    const slashMatch = pattern.match(/^\/(.+)\/([a-z]*)$/i);
    if (slashMatch) {
      try {
        regex = new RegExp(slashMatch[1], slashMatch[2].replace(/[gy]/g, ""));
      } catch {
        regex = new RegExp(
          slashMatch[1].replace(/[.*+?^${}()|[\]\\]/g, "\\$&"),
          "i",
        );
      }
    } else {
      try {
        regex = new RegExp(pattern, "i");
      } catch {
        regex = new RegExp(pattern.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"), "i");
      }
    }
  } else {
    regex = DEFAULT_PATTERN;
  }
  const results = [];

  async function walk(dir) {
    if (limit && results.length >= limit) {
      return;
    }
    let entries;
    try {
      entries = await fs.readdir(dir, { withFileTypes: true });
    } catch (err) {
      logError(`Error listing directory ${dir}: ${err.message}`);
      return;
    }
    for (const entry of entries) {
      if (limit && results.length >= limit) {
        return;
      }
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        await walk(fullPath);
      } else if (entry.isFile()) {
        try {
          const fileStream = createReadStream(fullPath, { encoding: "utf8" });
          const rl = readline.createInterface({
            input: fileStream,
            crlfDelay: Infinity,
          });

          let lineNum = 0;
          for await (const line of rl) {
            lineNum++;
            if (limit && results.length >= limit) {
              rl.close();
              fileStream.destroy();
              break;
            }
            const cleanLine = stripAnsi(line);
            if (regex.test(cleanLine) || regex.test(line)) {
              results.push({
                file: path.relative(logsDir, fullPath),
                line: lineNum,
                text: cleanLine.trim(),
              });
            }
          }
        } catch (err) {
          // Log errors for better observability, but continue walking to avoid total failure
          logError(`Error processing file ${fullPath}: ${err.message}`);
        }
      }
    }
  }

  await walk(logsDir);
  return results;
}

/**
 * Parses command-line arguments for search_logs.js.
 * Supports:
 * - <logs_directory> [pattern]
 * - [-n limit] or [--limit limit] or [flag=limit]
 * @param {Array<string>} [args]
 * @returns {{logsDir?: string, pattern?: string, limit: number}}
 */
export function parseSearchArgs(args = []) {
  let limit = 100;
  let logsDir = undefined;
  let pattern = undefined;

  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if ((a === "-n" || a === "--limit") && i + 1 < args.length) {
      const parsed = parseInt(args[++i], 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (a.startsWith("-n=") || a.startsWith("--limit=")) {
      const parsed = parseInt(a.slice(a.indexOf("=") + 1), 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (!logsDir && !a.startsWith("-")) {
      logsDir = a;
    } else if (!pattern && !a.startsWith("-")) {
      pattern = a;
    }
  }

  return { logsDir, pattern, limit };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: search_logs.js <logs_directory> [regex_pattern] [-n limit]

Arguments:
  <logs_directory>    Path to directory containing downloaded log files
  [regex_pattern]     Regex pattern to search for (default: error/fail/failure/exception/exit status/fatal/panic)

Options:
  -n, --limit <num>   Maximum number of matching lines to display (default: 100)
  -h, --help          Show this help message
`);
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("-h") || args.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const { logsDir, pattern: customPattern, limit } = parseSearchArgs(args);

  if (!logsDir) {
    logError(
      "Usage: search_logs.js <logs_directory> [regex_pattern] [-n limit]",
    );
    process.exit(1);
  }

  try {
    const stat = await fs.stat(logsDir);
    if (!stat.isDirectory()) {
      logError(`'${logsDir}' is not a directory.`);
      process.exit(1);
    }
  } catch {
    logError(`Directory '${logsDir}' not found.`);
    process.exit(1);
  }

  const pattern = customPattern || DEFAULT_PATTERN;
  console.log(
    `Searching for pattern ${pattern} in ${logsDir}...\n----------------------------------------`,
  );

  try {
    const matches = await searchLogs({ logsDir, pattern, limit });
    if (matches.length === 0) {
      console.log(`No matches found for pattern: ${pattern}`);
      return;
    }

    console.log(`Found ${matches.length} matching lines:`);
    for (const m of matches) {
      console.log(`[${m.file}:${m.line}] ${m.text}`);
    }
  } catch (err) {
    logError(`Error searching logs: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
