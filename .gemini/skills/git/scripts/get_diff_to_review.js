#!/usr/bin/env node
import { simpleGit } from "simple-git";

/**
 * Parses command-line arguments for get_diff_to_review.js.
 * Extracts flags (--staged/--cached, --unstaged/--working-tree, --stat, --name-only, --head, --raw) and file paths.
 * @param {Array<string>} [args]
 * @returns {{files: Array<string>, stagedOnly: boolean, unstagedOnly: boolean, statOnly: boolean, nameOnly: boolean, head: boolean, raw: boolean}}
 */
export function parseDiffArgs(args = []) {
  const files = [];
  let stagedOnly = false;
  let unstagedOnly = false;
  let statOnly = false;
  let nameOnly = false;
  let head = false;
  let raw = false;
  let separatorReached = false;

  for (const arg of args) {
    if (separatorReached) {
      files.push(arg);
      continue;
    }
    if (arg === "--") {
      separatorReached = true;
      continue;
    }
    if (arg === "--staged" || arg === "--cached") {
      stagedOnly = true;
    } else if (arg === "--unstaged" || arg === "--working-tree") {
      unstagedOnly = true;
    } else if (arg === "--stat") {
      statOnly = true;
    } else if (arg === "--name-only") {
      nameOnly = true;
    } else if (arg === "--head") {
      head = true;
    } else if (arg === "--raw") {
      raw = true;
    } else if (arg && !arg.startsWith("-")) {
      files.push(arg);
    }
  }

  if (stagedOnly && unstagedOnly) {
    stagedOnly = false;
    unstagedOnly = false;
  }

  return { files, stagedOnly, unstagedOnly, statOnly, nameOnly, head, raw };
}

/**
 * Gets staged and unstaged diffs for review.
 * @param {object} [params]
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {Array<string>} [params.files] - Optional list of file paths to inspect
 * @param {boolean} [params.stagedOnly] - Retrieve only staged changes
 * @param {boolean} [params.unstagedOnly] - Retrieve only unstaged changes
 * @param {boolean} [params.statOnly] - Retrieve diffstat summary instead of full diff
 * @param {boolean} [params.nameOnly] - Retrieve only changed file names
 * @param {boolean} [params.head] - Retrieve full diff against HEAD
 * @returns {Promise<{staged: string, unstaged: string, head?: string}>}
 */
export async function getDiffToReview({
  git = simpleGit(),
  files = [],
  stagedOnly = false,
  unstagedOnly = false,
  statOnly = false,
  nameOnly = false,
  head = false,
} = {}) {
  const extraFlags = [];
  if (statOnly) {
    extraFlags.push("--stat");
  }
  if (nameOnly) {
    extraFlags.push("--name-only");
  }

  const fileArgs = files && files.length > 0 ? files : [];

  let headDiff = "";
  if (head) {
    try {
      headDiff = await git.raw([
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        ...extraFlags,
        "HEAD",
        "--",
        ...fileArgs,
      ]);
    } catch (err) {
      const errMsg = err?.message || String(err);
      if (
        /ambiguous argument ['"]?HEAD['"]?/i.test(errMsg) ||
        /unknown revision or path not in the working tree/i.test(errMsg) ||
        /bad revision ['"]?HEAD['"]?/i.test(errMsg)
      ) {
        headDiff = "";
      } else {
        throw err;
      }
    }
  }

  let staged = "";
  if (!unstagedOnly && !head) {
    staged = await git.raw([
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--no-color",
      ...extraFlags,
      "--staged",
      "--",
      ...fileArgs,
    ]);
  }

  let unstaged = "";
  if (!stagedOnly && !head) {
    unstaged = await git.raw([
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--no-color",
      ...extraFlags,
      "--",
      ...fileArgs,
    ]);
  }

  return { staged, unstaged, head: headDiff };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_diff_to_review.js [--staged] [--unstaged] [--stat] [--name-only] [--head] [--raw] [files...]

Gets formatted staged and unstaged diffs for review.

Options:
  --staged, --cached       Show only staged changes
  --unstaged, --working-tree Show only unstaged changes
  --stat                   Show diffstat summary instead of full diff
  --name-only              Show changed file names only
  --head                   Show complete diff against HEAD
  --raw                    Output raw diff stream without decorative headers
  -h, --help               Show this help message
`);
}

async function main() {
  if (process.argv.includes("-h") || process.argv.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const { files, stagedOnly, unstagedOnly, statOnly, nameOnly, head, raw } =
    parseDiffArgs(process.argv.slice(2));
  try {
    const {
      staged,
      unstaged,
      head: headDiff,
    } = await getDiffToReview({
      files,
      stagedOnly,
      unstagedOnly,
      statOnly,
      nameOnly,
      head,
    });
    if (raw) {
      if (head) {
        process.stdout.write(headDiff || "");
      } else if (stagedOnly) {
        process.stdout.write(staged);
      } else if (unstagedOnly) {
        process.stdout.write(unstaged);
      } else {
        // Concatenate staged and unstaged diffs (both already terminate with newlines natively)
        process.stdout.write([staged, unstaged].filter(Boolean).join(""));
      }
      return;
    }
    if (head) {
      console.log("=== Changes relative to HEAD ===");
      console.log(headDiff.trim() ? headDiff : "(No changes relative to HEAD)");
      return;
    }
    if (!unstagedOnly) {
      console.log("=== Staged Changes (to be committed) ===");
      console.log(staged.trim() ? staged : "(No staged changes)");
    }
    if (!stagedOnly && !unstagedOnly) {
      console.log("\n=== Unstaged Changes (working tree) ===");
      console.log(unstaged.trim() ? unstaged : "(No unstaged changes)");
    } else if (unstagedOnly) {
      console.log("=== Unstaged Changes (working tree) ===");
      console.log(unstaged.trim() ? unstaged : "(No unstaged changes)");
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
