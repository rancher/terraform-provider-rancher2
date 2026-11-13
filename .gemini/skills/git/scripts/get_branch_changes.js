#!/usr/bin/env node
import { simpleGit } from "simple-git";
import { PREFERRED_DEFAULT_BRANCHES } from "./create-branch.js";

/**
 * Gets diff between current branch and base branch (main/master).
 * @param {object} [params]
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} [params.baseBranch]
 * @param {Array<string>} [params.files]
 * @param {boolean} [params.statOnly]
 * @param {boolean} [params.nameOnly]
 * @returns {Promise<string>}
 */
export async function getBranchChanges({
  git = simpleGit({ config: ["core.quotepath=false"] }),
  baseBranch,
  files = [],
  statOnly = false,
  nameOnly = false,
} = {}) {
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

  const extraFlags = [];
  if (statOnly) {
    extraFlags.push("--stat");
  }
  if (nameOnly) {
    extraFlags.push("--name-only");
  }

  const fileArgs =
    files && files.length > 0 ? files.filter((f) => f !== "--") : [];
  if (git.raw) {
    try {
      const rawDiff = await git.raw([
        "-c",
        "core.quotepath=false",
        "diff",
        "--no-ext-diff",
        "--no-color",
        ...extraFlags,
        mergeBase,
        "--",
        ...fileArgs,
      ]);
      if (typeof rawDiff === "string") {
        return rawDiff;
      }
    } catch {
      // Fallback to git.diff if raw diff fails
    }
  }
  if (git.diff) {
    // Fallback `git.diff` strictly appends to `git diff` and cannot accept global `-c` flags inline.
    // It relies entirely on the injected `git` instance having been configured with `core.quotepath=false`
    // at instantiation (e.g. `simpleGit({ config: ['core.quotepath=false'] })`) to prevent octal-escaped paths.
    return await git.diff([
      "--no-ext-diff",
      "--no-color",
      ...extraFlags,
      mergeBase,
      "--",
      ...fileArgs,
    ]);
  }
  return "";
}

/**
 * Parses command-line arguments for get_branch_changes.js.
 * Supports:
 * - [--stat|--name-only]
 * - [-b base_branch] or [--base base_branch] or [--base=base_branch]
 * - [base_branch] (only for standard branches like main, master, dev, trunk or remote tracking refs)
 * - [file_path...]
 * @param {Array<string>} [args]
 * @returns {{baseBranch?: string, files: Array<string>, statOnly: boolean, nameOnly: boolean}}
 */
export function parseBranchChangesArgs(args = []) {
  let baseBranch;
  const files = [];
  let statOnly = false;
  let nameOnly = false;
  let separatorReached = false;

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (separatorReached) {
      files.push(arg);
      continue;
    }
    if (arg === "--") {
      separatorReached = true;
      continue;
    }
    if (arg === "--stat") {
      statOnly = true;
    } else if (arg === "--name-only") {
      nameOnly = true;
    } else if (arg === "-b" || arg === "--base") {
      if (i + 1 < args.length) {
        baseBranch = args[++i];
      } else {
        throw new Error(`Flag '${arg}' requires a value.`);
      }
    } else if (arg.startsWith("-b=") || arg.startsWith("--base=")) {
      baseBranch = arg.slice(arg.indexOf("=") + 1);
    } else if (!baseBranch && !arg.startsWith("-")) {
      if (
        PREFERRED_DEFAULT_BRANCHES.includes(arg) ||
        arg.startsWith("origin/") ||
        arg.startsWith("upstream/") ||
        arg.startsWith("remotes/") ||
        arg.startsWith("refs/")
      ) {
        baseBranch = arg;
      } else {
        files.push(arg);
      }
    } else if (!arg.startsWith("-")) {
      files.push(arg);
    }
  }

  return { baseBranch, files, statOnly, nameOnly };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_branch_changes.js [--stat | --name-only] [-b base_branch] [file_path...]

Outputs git diff between current branch and the merge-base with the base branch.

Options:
  -b, --base <base>  Specify base branch (default: main/master/trunk/dev)
  --stat             Show diffstat summary instead of full diff
  --name-only        Show only names of changed files
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
    const { baseBranch, files, statOnly, nameOnly } = parseBranchChangesArgs(
      process.argv.slice(2),
    );

    const diff = await getBranchChanges({
      baseBranch,
      files,
      statOnly,
      nameOnly,
    });
    if (!diff.trim()) {
      if (!nameOnly) {
        console.log("No branch changes found relative to base branch.");
      }
      return;
    }
    process.stdout.write(diff.endsWith("\n") ? diff : diff + "\n");
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
