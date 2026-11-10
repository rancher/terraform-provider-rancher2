#!/usr/bin/env node
import { simpleGit } from "simple-git";
import { PREFERRED_DEFAULT_BRANCHES } from "./create-branch.js";

/**
 * Lists modified files and branch diff status vs main/master.
 * @param {object} [params]
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} [params.baseBranch]
 * @returns {Promise<{statusText: string, branchDiffText: string, baseBranch: string}>}
 */
export async function getFilesChanged({
  git = simpleGit({ config: ["core.quotepath=false"] }),
  baseBranch,
} = {}) {
  const status = await git.raw([
    "-c",
    "core.quotepath=false",
    "status",
    "--short",
  ]);

  let base = baseBranch;
  if (!base) {
    const branches = git.branch ? await git.branch() : await git.branchLocal();
    for (const pref of PREFERRED_DEFAULT_BRANCHES) {
      if (branches.all.includes(pref)) {
        base = pref;
        break;
      }
    }
    if (!base) {
      for (const pref of PREFERRED_DEFAULT_BRANCHES) {
        if (
          branches.all.includes(`remotes/upstream/${pref}`) ||
          branches.all.includes(`upstream/${pref}`) ||
          branches.all.includes(`refs/remotes/upstream/${pref}`)
        ) {
          base = `upstream/${pref}`;
          break;
        }
        if (
          branches.all.includes(`remotes/origin/${pref}`) ||
          branches.all.includes(`origin/${pref}`) ||
          branches.all.includes(`refs/remotes/origin/${pref}`)
        ) {
          base = `origin/${pref}`;
          break;
        }
      }
    }
    if (!base) {
      base = "main";
    }
  }

  let mergeBase;
  try {
    mergeBase = (await git.raw(["merge-base", base, "HEAD"])).trim();
  } catch {
    try {
      await git.raw(["rev-parse", "--verify", base]);
      mergeBase = base;
    } catch {
      throw new Error(`Base reference '${base}' is invalid or unresolvable.`);
    }
  }
  if (!mergeBase) {
    mergeBase = base;
  }

  let branchDiff = "";
  if (mergeBase) {
    branchDiff = await git.raw([
      "-c",
      "core.quotepath=false",
      "diff",
      "--no-ext-diff",
      "--name-status",
      mergeBase,
    ]);
  }

  return { statusText: status, branchDiffText: branchDiff, baseBranch: base };
}

/**
 * Parses command-line arguments for get_files_changed.js.
 * Supports:
 * - [-b base_branch] or [--base base_branch] or [--base=base_branch]
 * - [base_branch] (only for standard branches like main, master, dev, trunk or remote tracking refs)
 * @param {Array<string>} [args]
 * @returns {{baseBranch?: string}}
 */
export function parseFilesChangedArgs(args = []) {
  let baseBranch = undefined;
  let separatorReached = false;
  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if (separatorReached) {
      continue;
    }
    if (a === "--") {
      separatorReached = true;
      continue;
    }
    if (a === "-b" || a === "--base") {
      if (i + 1 < args.length) {
        baseBranch = args[++i];
      } else {
        throw new Error(`Flag '${a}' requires a value.`);
      }
    } else if (a.startsWith("-b=") || a.startsWith("--base=")) {
      baseBranch = a.slice(a.indexOf("=") + 1);
    } else if (!baseBranch && a && !a.startsWith("-")) {
      if (
        PREFERRED_DEFAULT_BRANCHES.includes(a) ||
        a.startsWith("origin/") ||
        a.startsWith("upstream/") ||
        a.startsWith("remotes/") ||
        a.startsWith("refs/")
      ) {
        baseBranch = a;
      }
    }
  }
  return { baseBranch };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_files_changed.js [-b base_branch]

Lists modified files in the working tree and branch diff status vs base.

Options:
  -b, --base <base>  Specify base branch (default: main/master/trunk/dev)
  -h, --help         Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  try {
    const { baseBranch } = parseFilesChangedArgs(process.argv.slice(2));
    const {
      statusText,
      branchDiffText,
      baseBranch: detectedBase,
    } = await getFilesChanged({ baseBranch });
    console.log("=== Current Working Tree Status ===");
    console.log(statusText.trim() ? statusText : "(Working tree clean)");
    console.log(
      `\n=== Files Changed in Current Branch vs ${detectedBase || "base"} ===`,
    );
    console.log(
      branchDiffText.trim()
        ? branchDiffText
        : "(No changes relative to base branch)",
    );
  } catch (err) {
    console.error(`❌ Error: ${err?.message || err}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
