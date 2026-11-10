#!/usr/bin/env node
import { simpleGit } from "simple-git";

/**
 * Unstages files from the git index safely.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {Array<string>} [params.files]
 * @param {boolean} [params.all]
 * @returns {Promise<Array<string>>} List of unstaged files
 */
export async function unstageFiles({
  git = simpleGit(),
  files = [],
  all = false,
}) {
  if (all) {
    let rawStaged = "";
    try {
      rawStaged = await git.raw([
        "-c",
        "core.quotepath=false",
        "diff",
        "--name-only",
        "--cached",
      ]);
    } catch {
      try {
        const stageOutput = await git.raw([
          "-c",
          "core.quotepath=false",
          "ls-files",
          "--stage",
        ]);
        rawStaged = stageOutput
          .split("\n")
          .map((l) => {
            const tabIdx = l.indexOf("\t");
            return tabIdx !== -1
              ? l.slice(tabIdx + 1).replace(/^"(.*)"$/, "$1")
              : "";
          })
          .filter(Boolean)
          .join("\n");
      } catch {
        // Fall through
      }
    }
    const stagedFiles = rawStaged
      .split("\n")
      .map((l) => l.trim().replace(/^"(.*)"$/, "$1"))
      .filter(Boolean);

    if (stagedFiles.length === 0) {
      return [];
    }

    try {
      await git.reset(["HEAD", "--"]);
    } catch {
      await git.raw(["rm", "--cached", "-r", "."]);
    }
    return stagedFiles;
  }

  if (!files || files.length === 0) {
    throw new Error("No files specified to unstage.");
  }

  try {
    await git.reset(["HEAD", "--", ...files]);
  } catch {
    await git.raw(["rm", "--cached", "--", ...files]);
  }
  return files;
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: unstage-files.js <file1> [file2...] | --all

Unstages specified files or all staged changes from the git index safely.

Options:
  -A, --all    Unstage all staged files from the index
  -h, --help   Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }

  const args = process.argv.slice(2);
  if (args.length === 0) {
    console.error("Usage: unstage-files.js <file1> [file2...] | --all");
    process.exit(1);
  }

  const all = args.includes("--all") || args.includes("-A");
  const files = args.filter((a) => a !== "--all" && a !== "-A" && a !== "--");

  try {
    const unstaged = await unstageFiles({ files, all });
    if (unstaged.length === 0) {
      console.log("No staged files to unstage.");
      return;
    }
    console.log(`✅ Successfully unstaged: ${unstaged.join(", ")}`);
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
