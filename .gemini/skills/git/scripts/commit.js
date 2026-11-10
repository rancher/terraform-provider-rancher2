#!/usr/bin/env node
import { execFile } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { simpleGit } from "simple-git";

const CONVENTIONAL_COMMIT_REGEX =
  /^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert|env|dev)(?:\([^)]+\))?(!?): .+/i;

export const DEFAULT_PRODUCT_ROOTS = [
  "internal/",
  "src/",
  "lib/",
  "pkg/",
  "cmd/",
];

/**
 * Loads configured product path roots from .gemini/settings.json, falling back to DEFAULT_PRODUCT_ROOTS.
 * @param {string} [repoRoot] - Repository root directory
 * @returns {Promise<Array<string>>}
 */
export async function getProductRoots(repoRoot = null) {
  try {
    const root = repoRoot || process.cwd();
    const settingsPath = path.join(root, ".gemini/settings.json");
    const content = await fs.readFile(settingsPath, "utf8");
    const parsed = JSON.parse(content);
    if (
      Array.isArray(parsed?.project?.productPaths) &&
      parsed.project.productPaths.length > 0
    ) {
      return parsed.project.productPaths;
    }
  } catch {
    // Fallback if settings file is missing or unparseable
  }
  return DEFAULT_PRODUCT_ROOTS;
}

/**
 * Checks whether a relative file path belongs to product code.
 * @param {string} filePath - Path to evaluate
 * @param {Array<string>} [productRoots] - Product directory roots
 * @returns {boolean}
 */
export function fileAffectsProduct(
  filePath,
  productRoots = DEFAULT_PRODUCT_ROOTS,
) {
  const normalized = filePath.replace(/^[./]+/, "");
  for (const root of productRoots) {
    const cleanRoot = root.endsWith("/") ? root : `${root}/`;
    if (normalized.startsWith(cleanRoot)) {
      return true;
    }
  }
  return false;
}

/**
 * Runs cspell check on text asynchronously.
 * @param {string} text - Text to check
 * @param {string} [repoRoot] - Repository root directory
 * @returns {Promise<{valid: boolean, misspelled?: string, reason?: string, skipped?: boolean}>}
 */
export function checkSpelling(text, repoRoot = process.cwd()) {
  return new Promise((resolve) => {
    const child = execFile(
      "cspell",
      ["stdin", "--quiet", "--words-only"],
      { cwd: repoRoot },
      (err, stdout) => {
        if (err && !stdout) {
          // cspell not found or non-spelling error; skip gracefully
          return resolve({ valid: true, skipped: true });
        }
        const words = (stdout || "").trim();
        if (words) {
          const list = words.replace(/\r?\n/g, ", ");
          return resolve({
            valid: false,
            misspelled: list,
            reason: `Commit message contains spelling errors on: ${list}`,
          });
        }
        resolve({ valid: true });
      },
    );

    child.on("error", () => {
      resolve({ valid: true, skipped: true });
    });

    if (child.stdin) {
      child.stdin.end(text);
    }
  });
}

/**
 * Validates a candidate commit message against conventional commit format and repository boundaries.
 * @param {object} params
 * @param {string} params.message
 * @param {boolean} params.affectsProduct
 * @param {boolean} [params.isMerge]
 * @returns {{valid: boolean, reason?: string}}
 */
export function validateCommitMessage({
  message,
  affectsProduct,
  isMerge = false,
}) {
  if (!message || !message.trim()) {
    return { valid: false, reason: "Missing commit message." };
  }

  const subject = message.split("\n")[0].trim().replace(/\r$/, "");
  if (isMerge || /^Merge\s+/i.test(subject)) {
    return { valid: true };
  }

  if (!CONVENTIONAL_COMMIT_REGEX.test(subject)) {
    return {
      valid: false,
      reason: `Commit message title does not follow Conventional Commits format: '${subject}'. Expected 'type: description' or 'type(scope): description'.`,
    };
  }

  if (subject.length > 72) {
    return {
      valid: false,
      reason: `Commit subject line exceeds 72 characters (${subject.length}): '${subject}'`,
    };
  }

  const match = subject.match(/^([a-z]+)(?:\([^)]+\))?(!?):/i);
  const type = match[1].toLowerCase();
  const hasBreakingIndicator =
    match[2] === "!" || /(^|\s)BREAKING[ -]CHANGE:/m.test(message);

  if (!affectsProduct) {
    if (type === "feat" || type === "refactor" || hasBreakingIndicator) {
      return {
        valid: false,
        reason:
          "Non-product change (outside 'internal/') must NOT use 'feat', 'refactor', or '!' breaking-change indicators.",
      };
    }
  }

  return { valid: true };
}

/**
 * Validates a single commit title against Conventional Commits.
 * @param {string} title - Commit title
 * @param {boolean} affectsProduct - Whether commit touches product code
 * @param {boolean} [isMerge=false] - Whether commit is a merge commit
 * @returns {{valid: boolean, reason?: string, isMerge?: boolean}}
 */
export function validateCommitTitle(title, affectsProduct, isMerge = false) {
  return validateCommitMessage({ message: title, affectsProduct, isMerge });
}

/**
 * Asynchronously validates all commit checks including spell-checking.
 * @param {object} params
 * @param {string} params.message
 * @param {boolean} params.affectsProduct
 * @param {Function} [params.checkSpellingFn]
 * @returns {Promise<{valid: boolean, reason?: string}>}
 */
export async function validateCommit({
  message,
  affectsProduct,
  checkSpellingFn = checkSpelling,
}) {
  const syncResult = validateCommitMessage({ message, affectsProduct });
  if (!syncResult.valid) {
    return syncResult;
  }

  const subject = message.split("\n")[0].trim();
  const isMerge = /^Merge\s+/i.test(subject);

  if (!isMerge && checkSpellingFn) {
    const spellResult = await checkSpellingFn(subject);
    if (!spellResult.valid) {
      return { valid: false, reason: spellResult.reason };
    }
  }

  return { valid: true };
}

/**
 * Creates a validated commit, blocking main/master branch commits and enforcing all checks.
 * @param {object} params
 * @param {import('simple-git').SimpleGit} [params.git]
 * @param {string} params.message
 * @param {Function} [params.checkSpellingFn]
 * @returns {Promise<object>} Commit summary
 */
export async function createCommit({
  git = simpleGit(),
  message,
  checkSpellingFn = checkSpelling,
}) {
  const currentBranch = (await git.branch()).current;
  if (!currentBranch || currentBranch === "HEAD") {
    throw new Error("Cannot commit on a detached HEAD.");
  }

  if (currentBranch === "main" || currentBranch === "master") {
    throw new Error(
      `Commits directly to the '${currentBranch}' branch are not allowed.`,
    );
  }

  const diffSummary = await git.diffSummary(["--no-ext-diff", "--cached"]);
  if (diffSummary.files.length === 0) {
    const unstaged = await git.diffSummary(["--no-ext-diff"]);
    if (unstaged.files.length === 0) {
      throw new Error(
        "No changes to commit (working tree clean and staging area empty).",
      );
    } else {
      throw new Error(
        "No staged changes found. Please stage changes with stage-files.js before committing.",
      );
    }
  }

  const roots = await getProductRoots();
  const affectsProduct = diffSummary.files.some((f) =>
    fileAffectsProduct(f.file, roots),
  );
  const validation = await validateCommit({
    message,
    affectsProduct,
    checkSpellingFn,
  });
  if (!validation.valid) {
    throw new Error(validation.reason);
  }

  const result = await git.commit(message, undefined, { "-s": null });
  return { ...result, branch: result.branch || currentBranch };
}

/**
 * Parses commit message arguments from CLI tokens.
 * Multiple -m or --message flags are joined with double newlines (\n\n) as paragraphs.
 * Bare unflagged arguments are joined with a space.
 * @param {Array<string>} args
 * @returns {string} Combined commit message
 */
export function parseCommitArgs(args) {
  if (!args || args.length === 0) {
    return "";
  }

  const paragraphs = [];
  const positionalWords = [];
  let hasMessageFlag = false;

  for (let i = 0; i < args.length; i++) {
    if ((args[i] === "-m" || args[i] === "--message") && i + 1 < args.length) {
      hasMessageFlag = true;
      paragraphs.push(args[i + 1]);
      i++;
    } else if (args[i].startsWith("-m=") || args[i].startsWith("--message=")) {
      hasMessageFlag = true;
      paragraphs.push(args[i].slice(args[i].indexOf("=") + 1));
    } else {
      positionalWords.push(args[i]);
    }
  }

  if (hasMessageFlag) {
    return paragraphs.join("\n\n").trim();
  }

  return positionalWords.join(" ").trim();
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: commit.js [-m <message>] <message>

Creates a validated conventional commit with spellcheck and SemVer branch safety.

Options:
  -m, --message <text>   Specify commit message text (supports multiple flags)
  -h, --help             Show this help message
`);
}

async function main() {
  const rawArgs = process.argv.slice(2);
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    printHelp();
    process.exit(0);
  }

  const args = process.argv.slice(2);
  const message = parseCommitArgs(args);

  if (!message) {
    console.error("Usage: commit.js [-m <message>] <message>");
    process.exit(1);
  }

  try {
    const summary = await createCommit({ message });
    console.log(
      `✅ Commit successfully created on branch '${summary.branch}'! (${summary.commit})`,
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
