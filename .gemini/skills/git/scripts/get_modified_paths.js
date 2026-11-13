#!/usr/bin/env node
import { simpleGit } from "simple-git";

/**
 * Parses git status porcelain output and returns an array of modified/untracked file paths.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @returns {Promise<Array<string>>}
 */
export async function getModifiedPaths({ git = simpleGit() } = {}) {
  const porcelain = await git.raw([
    "-c",
    "core.quotepath=false",
    "status",
    "--porcelain",
  ]);
  const lines = porcelain.split("\n");
  const paths = [];

  for (const line of lines) {
    if (!line || line.length < 3) {
      continue;
    }
    let filePath = line.slice(3);
    if (filePath.includes(" -> ")) {
      filePath = filePath.split(" -> ")[1];
    }
    filePath = filePath.trim();
    if (filePath.startsWith('"') && filePath.endsWith('"')) {
      filePath = filePath.slice(1, -1);
    }
    if (filePath) {
      paths.push(filePath);
    }
  }

  return paths;
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_modified_paths.js

Lists modified and untracked file paths from git status porcelain.

Options:
  -h, --help   Show this help message
`);
}

async function main() {
  if (process.argv.includes("-h") || process.argv.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  try {
    const paths = await getModifiedPaths();
    for (const p of paths) {
      console.log(p);
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
