#!/usr/bin/env node
import { simpleGit } from "simple-git";

/**
 * Switches to an existing local or remote branch safely.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} params.branchName
 * @returns {Promise<{branchName: string}>}
 */
export async function switchBranch({ git = simpleGit(), branchName } = {}) {
  if (!branchName) {
    throw new Error("Branch name cannot be empty.");
  }

  if (branchName.startsWith("-")) {
    throw new Error("Branch name cannot start with a hyphen.");
  }

  const cleanBranch = branchName.replace(
    /^(?:remotes\/)?(?:upstream|origin)\//,
    "",
  );

  const branches = await git.branch();
  const hasLocal = branches.all.includes(cleanBranch);
  const hasOrigin =
    branches.all.includes(`remotes/origin/${cleanBranch}`) ||
    branches.all.includes(`origin/${cleanBranch}`);
  const hasUpstream =
    branches.all.includes(`remotes/upstream/${cleanBranch}`) ||
    branches.all.includes(`upstream/${cleanBranch}`);

  const branchExists =
    hasLocal || hasOrigin || hasUpstream || branches.all.includes(branchName);

  if (!branchExists) {
    throw new Error(`Branch '${branchName}' does not exist.`);
  }

  if (!hasLocal && !hasOrigin && hasUpstream) {
    try {
      await git.checkout([
        "-b",
        cleanBranch,
        "--track",
        `upstream/${cleanBranch}`,
      ]);
    } catch {
      await git.checkout(cleanBranch);
    }
  } else {
    await git.checkout(cleanBranch);
  }
  return { branchName: cleanBranch };
}

/**
 * Gets the current active local branch name.
 * @param {object} [params]
 * @param {import('simple-git').SimpleGit} [params.git]
 * @returns {Promise<string>}
 */
export async function getCurrentBranch({ git = simpleGit() } = {}) {
  const branchSummary = await git.branch();
  return branchSummary.current || "";
}

/**
 * Checks whether a branch exists locally or on tracked remotes.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} params.branchName
 * @returns {Promise<boolean>}
 */
export async function checkBranchExists({
  git = simpleGit(),
  branchName,
} = {}) {
  if (!branchName) {
    return false;
  }
  const cleanBranch = branchName
    .replace(/^refs\/(?:heads|remotes)\//, "")
    .replace(/^(?:remotes\/)?(?:upstream|origin)\//, "");
  try {
    const branches = await git.branch();
    return (
      branches.all.includes(cleanBranch) ||
      branches.all.includes(`remotes/origin/${cleanBranch}`) ||
      branches.all.includes(`origin/${cleanBranch}`) ||
      branches.all.includes(`remotes/upstream/${cleanBranch}`) ||
      branches.all.includes(`upstream/${cleanBranch}`) ||
      branches.all.includes(branchName)
    );
  } catch (err) {
    console.debug(
      `[DEBUG] Error checking branch existence for '${branchName}': ${err.message}`,
    );
    return false;
  }
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: switch-branch.js <branch-name> [--current] [--check <branch-name>]

Switches branches safely with tracking support, or inspects branch state.

Arguments:
  branch-name          Target branch name to check out

Options:
  --current            Print the name of the currently active branch
  --check, --exists    Check if specified branch exists (exits 0 if yes, 1 if no)
  -h, --help           Show this help message
`);
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("-h") || args.includes("--help")) {
    printHelp();
    process.exit(0);
  }

  if (args.includes("--current")) {
    try {
      const branch = await getCurrentBranch();
      console.log(branch);
      process.exit(0);
    } catch (err) {
      console.error(`❌ Error: ${err.message}`);
      process.exit(1);
    }
  }

  const checkIdx = args.findIndex((a) => a === "--check" || a === "--exists");
  if (checkIdx !== -1) {
    const target =
      args[checkIdx + 1] ||
      args.find((a) => a !== "--check" && a !== "--exists");
    if (!target) {
      console.error("Usage: switch-branch.js --check <branch-name>");
      process.exit(1);
    }
    try {
      const exists = await checkBranchExists({ branchName: target });
      if (exists) {
        console.log(`Branch '${target}' exists.`);
        process.exit(0);
      } else {
        console.error(`Branch '${target}' does not exist.`);
        process.exit(1);
      }
    } catch (err) {
      console.error(`❌ Error: ${err.message}`);
      process.exit(1);
    }
  }

  const branchName = args[0];
  if (!branchName) {
    console.error(
      "Usage: switch-branch.js <branch-name> [--current] [--check <branch-name>]",
    );
    process.exit(1);
  }

  try {
    const result = await switchBranch({ branchName });
    console.log(`✅ Successfully switched to branch '${result.branchName}'.`);
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
