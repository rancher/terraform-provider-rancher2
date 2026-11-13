#!/usr/bin/env node
import { simpleGit } from "simple-git";

export const PREFERRED_DEFAULT_BRANCHES = ["main", "master", "trunk", "dev"];

/**
 * Creates a new branch, syncing the base branch from upstream first.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} params.branchName
 * @param {string} [params.baseBranch]
 * @returns {Promise<{branchName: string, baseBranch: string, syncedRemote: string}>}
 */
export async function createBranch({
  git = simpleGit(),
  branchName,
  baseBranch,
} = {}) {
  if (!branchName) {
    throw new Error("Branch name cannot be empty.");
  }

  if (branchName.startsWith("-")) {
    throw new Error("Branch name cannot start with a hyphen.");
  }

  // Determine base branch if not provided
  let base = baseBranch;
  if (!base) {
    if (typeof git.raw === "function") {
      try {
        const remoteHead = await git.raw([
          "symbolic-ref",
          "--short",
          "refs/remotes/origin/HEAD",
        ]);
        if (remoteHead && remoteHead.trim()) {
          base = remoteHead.trim().replace(/^origin\//, "");
        }
      } catch {
        // ignore and fallback to local branch detection
      }
    }
    if (!base) {
      const branches = await git.branchLocal();
      base =
        PREFERRED_DEFAULT_BRANCHES.find((b) => branches.all.includes(b)) ||
        "main";
    }
  }

  // Detect upstream remote
  const remotes = await git.getRemotes();
  const upstreamRemote = remotes.some((r) => r.name === "upstream")
    ? "upstream"
    : "origin";

  const remoteRef =
    base.startsWith(`${upstreamRemote}/`) ||
    base.startsWith("origin/") ||
    base.startsWith("upstream/")
      ? base.replace(
        /^(?:origin|upstream|remotes\/upstream|remotes\/origin)\//,
        "",
      )
      : base;

  console.log(
    `Syncing base branch '${remoteRef}' from remote '${upstreamRemote}'...`,
  );
  try {
    await git.fetch(upstreamRemote, remoteRef);
    const current = (await git.branch()).current;
    if (current === remoteRef) {
      await git.merge(["--ff-only", `${upstreamRemote}/${remoteRef}`]);
    } else {
      await git.fetch(upstreamRemote, `${remoteRef}:${remoteRef}`);
    }
    console.log(
      `Base branch '${remoteRef}' synced with '${upstreamRemote}/${remoteRef}'.`,
    );
  } catch (err) {
    throw new Error(
      `Failed to sync base branch '${remoteRef}' from '${upstreamRemote}': ${err.message}`,
      {
        cause: err,
      },
    );
  }

  // Verify branch does not already exist
  const existingBranches = await git.branchLocal();
  if (existingBranches.all.includes(branchName)) {
    throw new Error(`Branch '${branchName}' already exists.`);
  }

  console.log(
    `Creating and switching to branch '${branchName}' based on '${base}'...`,
  );
  await git.checkoutBranch(branchName, base);

  return { branchName, baseBranch: base, syncedRemote: upstreamRemote };
}

/**
 * Parses command-line arguments for create-branch.js.
 * Supports:
 * - <branch-name> [base-branch]
 * - <branch-name> --base <base-branch> or -b <base-branch>
 * @param {Array<string>} [args]
 * @returns {{branchName: string, baseBranch?: string}}
 */
export function parseCreateBranchArgs(args = []) {
  let branchName = "";
  let baseBranch = undefined;

  let i = 0;
  while (i < args.length) {
    const arg = args[i];
    if (arg === "--base" || arg === "-b") {
      if (i + 1 >= args.length) {
        throw new Error(`Flag '${arg}' requires a value.`);
      }
      baseBranch = args[i + 1];
      i += 2;
    } else if (arg.startsWith("--base=")) {
      baseBranch = arg.slice(7);
      i++;
    } else if (arg.startsWith("-b=")) {
      baseBranch = arg.slice(3);
      i++;
    } else if (!branchName && !arg.startsWith("-")) {
      branchName = arg;
      i++;
    } else if (!baseBranch && !arg.startsWith("-")) {
      baseBranch = arg;
      i++;
    } else {
      i++;
    }
  }

  return { branchName, baseBranch };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: create-branch.js <branch-name> [--base <base-branch>]

Creates and checks out a new branch, synchronizing base from upstream or origin.

Arguments:
  branch-name                  Name of the new branch to create

Options:
  -b, --base <base-branch>     Base branch to branch off from (default: auto-detected default branch)
  -h, --help                   Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  let parsed;
  try {
    parsed = parseCreateBranchArgs(process.argv.slice(2));
  } catch (err) {
    console.error(`❌ Error: ${err.message}`);
    process.exit(1);
  }
  const { branchName, baseBranch } = parsed;

  if (!branchName) {
    console.error(
      "Usage: create-branch.js <branch-name> [--base <base-branch>]",
    );
    process.exit(1);
  }

  try {
    const result = await createBranch({ branchName, baseBranch });
    console.log(
      `✅ Successfully created and switched to branch '${result.branchName}'.`,
    );
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
