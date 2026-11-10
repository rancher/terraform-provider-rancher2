#!/usr/bin/env node
import path from "node:path";
import { simpleGit } from "simple-git";
import { isExcludedByAIExclude } from "./read-ignored-file.js";

/**
 * Stages files safely, ensuring .aiexclude patterns are respected.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {Array<string>} [params.files]
 * @param {boolean} [params.all]
 * @param {boolean} [params.includeUntracked]
 * @returns {Promise<Array<string>>} List of staged files
 */
export async function stageFiles({
  git = simpleGit(),
  files = [],
  all = false,
  includeUntracked = false,
}) {
  const root = (await git.revparse(["--show-toplevel"])).trim();

  let targets = [...files];
  if (all || includeUntracked) {
    const args = ["-c", "core.quotepath=false", "ls-files"];
    if (all) {
      args.push("-m", "-d");
    }
    if (includeUntracked) {
      args.push("-o", "--exclude-standard");
    }
    const statusOutput = await git.raw(args);
    const listed = statusOutput
      .split("\n")
      .map((l) => l.trim().replace(/^"(.*)"$/, "$1"))
      .filter(Boolean);
    targets = [...new Set([...targets, ...listed])];
    if (targets.length === 0) {
      return [];
    }
  }

  if (targets.length === 0) {
    throw new Error("No files specified to stage.");
  }

  await Promise.all(
    targets.map(async (file) => {
      const fullPath = path.isAbsolute(file) ? file : path.resolve(root, file);
      const isBlocked = await isExcludedByAIExclude(fullPath, root);
      if (isBlocked) {
        throw new Error(
          `File '${file}' is blocked by .aiexclude policy. Staging aborted.`,
        );
      }
    }),
  );

  await git.add(targets);
  return targets;
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: stage-files.js <file1> [file2...] | --all | --include-untracked

Stages files safely in the git index while enforcing .aiexclude policy checks.

Options:
  -A, --all              Stage modified and deleted files
  -u, --include-untracked Include untracked files
  -h, --help             Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  const dashDashIdx = rawArgs.indexOf("--");
  const optArgs = dashDashIdx !== -1 ? rawArgs.slice(0, dashDashIdx) : rawArgs;
  const postDashArgs =
    dashDashIdx !== -1 ? rawArgs.slice(dashDashIdx + 1) : [];

  if (optArgs.includes("-h") || optArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }

  if (rawArgs.length === 0) {
    console.error(
      "Usage: stage-files.js <file1> [file2...] | --all | --include-untracked",
    );
    process.exit(1);
  }

  const all = optArgs.includes("--all") || optArgs.includes("-A");
  const includeUntracked =
    optArgs.includes("--include-untracked") || optArgs.includes("-u");
  const preDashFiles = optArgs.filter(
    (a) =>
      a !== "--all" &&
      a !== "-A" &&
      a !== "--include-untracked" &&
      a !== "-u",
  );
  const files = [...preDashFiles, ...postDashArgs];

  try {
    const staged = await stageFiles({ files, all, includeUntracked });
    if (staged.length === 0) {
      console.log("No files to stage.");
      return;
    }
    console.log(`✅ Successfully staged: ${staged.join(", ")}`);
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
