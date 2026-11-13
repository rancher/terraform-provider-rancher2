#!/usr/bin/env node
import { simpleGit } from "simple-git";

/**
 * Parses raw git ls-files --stage output into structured objects.
 * Format: <mode> <hash> <stage>\t<path>
 * @param {string} rawOutput - Output of `git ls-files --stage`
 * @returns {Array<{mode: string, hash: string, stage: number, path: string}>}
 */
export function parseStageOutput(rawOutput) {
  if (!rawOutput || !rawOutput.trim()) {
    return [];
  }

  const entries = [];
  const lines = rawOutput.split("\n");

  for (const line of lines) {
    if (!line.trim()) {
      continue;
    }
    const tabIdx = line.indexOf("\t");
    if (tabIdx === -1) {
      continue;
    }

    const metadata = line.slice(0, tabIdx).trim().split(/\s+/);
    let filePath = line.slice(tabIdx + 1).trim();

    if (filePath.startsWith('"') && filePath.endsWith('"')) {
      filePath = filePath.slice(1, -1);
    }

    if (metadata.length >= 3) {
      entries.push({
        mode: metadata[0],
        hash: metadata[1],
        stage: parseInt(metadata[2], 10),
        path: filePath,
      });
    }
  }

  return entries;
}

/**
 * Retrieves the list of currently staged files with stage details and status.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @returns {Promise<{entries: Array<{mode: string, hash: string, stage: number, path: string, status?: string}>, paths: Array<string>, hasConflicts: boolean}>}
 */
export async function getStagedFiles({ git = simpleGit() } = {}) {
  const rawStage = await git.raw([
    "-c",
    "core.quotepath=false",
    "ls-files",
    "--stage",
  ]);
  const entries = parseStageOutput(rawStage);

  // Retrieve cached diff status (A, M, D, R, etc.)
  const statusMap = new Map();
  try {
    const rawDiff = await git.raw([
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--name-status",
      "--cached",
    ]);
    if (rawDiff.trim()) {
      for (const line of rawDiff.split("\n")) {
        if (!line.trim()) {
          continue;
        }
        const parts = line.split("\t");
        if (parts.length >= 2) {
          const statusCode = parts[0].trim();
          const targetFile = parts[parts.length - 1]
            .trim()
            .replace(/^"(.*)"$/, "$1");
          statusMap.set(targetFile, statusCode);
        }
      }
    }
  } catch {
    // Ignore error if diff fails
  }

  let hasConflicts = false;
  const stagedEntries = [];

  for (const entry of entries) {
    if (statusMap.has(entry.path)) {
      entry.status = statusMap.get(entry.path);
    }
    if (entry.stage > 0) {
      hasConflicts = true;
    }
    if (entry.status || entry.stage > 0) {
      stagedEntries.push(entry);
    }
  }

  // Also include any deleted files that appear in statusMap ('D') but are no longer in ls-files stage 0
  const existingPaths = new Set(stagedEntries.map((e) => e.path));
  for (const [filePath, status] of statusMap.entries()) {
    if (!existingPaths.has(filePath)) {
      stagedEntries.push({
        mode: "000000",
        hash: "0000000",
        stage: 0,
        path: filePath,
        status,
      });
      existingPaths.add(filePath);
    }
  }

  const paths = [...new Set(stagedEntries.map((e) => e.path))];
  return { entries: stagedEntries, paths, hasConflicts };
}

/**
 * Parses command-line arguments for get_staged_files.js.
 * Supports:
 * - [--name-only] or [-n]
 * @param {Array<string>} [args]
 * @returns {{nameOnly: boolean}}
 */
export function parseStageArgs(args = []) {
  const nameOnly = args.includes("--name-only") || args.includes("-n");
  return { nameOnly };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_staged_files.js [-n | --name-only]

Retrieves the list of currently staged files in index with stage details and conflicts.

Options:
  -n, --name-only  Output only file paths
  -h, --help       Show this help message
`);
}

async function main() {
  if (process.argv.includes("-h") || process.argv.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const { nameOnly } = parseStageArgs(process.argv.slice(2));
  try {
    const { entries, paths, hasConflicts } = await getStagedFiles();
    if (nameOnly) {
      for (const p of paths) {
        console.log(p);
      }
      return;
    }

    console.log(`=== Staged Files in Index (${paths.length} file(s)) ===`);

    if (entries.length === 0) {
      console.log("(No staged changes)");
      return;
    }

    for (const e of entries) {
      const statusBadge = e.status ? `[${e.status}] ` : "";
      const conflictBadge =
        e.stage > 0 ? ` ⚠️ CONFLICT (stage ${e.stage})` : "";
      console.log(
        `${statusBadge}${e.path} (mode: ${e.mode}, hash: ${e.hash.slice(0, 7)})${conflictBadge}`,
      );
    }

    if (hasConflicts) {
      console.warn("\n⚠️ Warning: Merge conflicts detected in index stages!");
    }
  } catch (err) {
    console.error(`❌ Error retrieving staged files: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
