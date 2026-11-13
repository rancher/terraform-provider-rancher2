#!/usr/bin/env node
import path from "node:path";
import { simpleGit } from "simple-git";

/**
 * Retrieves the contents of a deleted or modified file from Git history.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} params.filePath
 * @param {string} [params.commit] - Specific commit SHA or revision
 * @returns {Promise<string>}
 */
export async function getDeletedFile({ git = simpleGit(), filePath, commit }) {
  if (!filePath) {
    throw new Error("Missing file path.");
  }

  let relPath = filePath.replace(/^\.\//, "");
  if (git.revparse) {
    try {
      const root = (await git.revparse(["--show-toplevel"])).trim();
      const fullPath = path.isAbsolute(filePath)
        ? filePath
        : path.resolve(root, filePath);
      const computed = path.relative(root, fullPath).replace(/\\/g, "/");
      if (!computed.startsWith("..") && !path.isAbsolute(computed)) {
        relPath = computed;
      }
    } catch {
      // Ignore revparse failure
    }
  }

  if (commit) {
    try {
      return await git.show([`${commit}:${relPath}`]);
    } catch {
      throw new Error(
        `Could not retrieve contents for '${filePath}' at commit '${commit}'.`,
      );
    }
  }

  // Find last commit where file was modified
  const log = await git.log({ file: relPath, maxCount: 1 });
  if (log.all.length === 0) {
    throw new Error(`File '${filePath}' was never found in git history.`);
  }

  const lastCommit = log.all[0].hash;

  try {
    const content = await git.show([`${lastCommit}:${relPath}`]);
    return content;
  } catch {
    // If not found at last commit (e.g. was deleted in that commit), check parent commit
    try {
      const content = await git.show([`${lastCommit}~1:${relPath}`]);
      return content;
    } catch {
      throw new Error(
        `Could not retrieve contents for '${filePath}' from history.`,
      );
    }
  }
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_deleted_file.js <file_path> [commit_or_rev]

Retrieves contents of a deleted or modified file from git history.

Arguments:
  file_path           Relative or absolute path to target file
  commit_or_rev       Optional commit SHA or git revision (default: auto-detected last commit)

Options:
  -h, --help          Show this help message
`);
}

async function main() {
  if (process.argv.includes("-h") || process.argv.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const filePath = process.argv[2];
  const commit = process.argv[3];
  if (!filePath) {
    console.error("Usage: get_deleted_file.js <file_path> [commit_or_rev]");
    process.exit(1);
  }

  try {
    const content = await getDeletedFile({ filePath, commit });
    process.stdout.write(content.endsWith("\n") ? content : content + "\n");
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
