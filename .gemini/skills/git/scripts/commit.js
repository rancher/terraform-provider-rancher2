#!/usr/bin/env node
import { execFile } from "node:child_process";
import nodeFs from "node:fs";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { simpleGit } from "simple-git";
import {
  isProductDirectory,
  isProductFile,
  isProductPrefix,
  NON_PRODUCT_PREFIXES,
  PRODUCT_PREFIXES,
  VALID_PREFIXES,
  validateEmpty,
  validateLength,
  validatePrefix,
  validateProductImpact,
} from "../../../../.github/workflows/scripts/validate-commit-message.js";

export {
  isProductDirectory,
  isProductFile,
  isProductPrefix,
  NON_PRODUCT_PREFIXES,
  PRODUCT_PREFIXES,
  VALID_PREFIXES,
};

export const DEFAULT_PRODUCT_ROOTS = ["rancher/", "rancher2/", "internal/"];

/**
 * Resolves the user's global Git configuration file path if present.
 * Respects explicit custom GIT_CONFIG_GLOBAL environments (like test isolation),
 * falling back to user's config when unset or pointing to /dev/null.
 * @param {object} [options]
 * @param {string} [options.userConfigPath] - Override path to check
 * @param {object} [options.fsModule] - File system module to use
 * @param {object} [options.env] - Environment variables object
 * @returns {string|null} Resolved config path or null if absent
 */
export function resolveGlobalGitConfig({
  userConfigPath = path.join(os.homedir(), ".gitconfig"),
  fsModule = nodeFs,
  env = process.env,
} = {}) {
  if (env.GIT_CONFIG_GLOBAL && env.GIT_CONFIG_GLOBAL !== "/dev/null") {
    return env.GIT_CONFIG_GLOBAL;
  }
  try {
    if (fsModule.existsSync(userConfigPath)) {
      return userConfigPath;
    }
  } catch {
    // Ignore access or stat errors
  }
  return null;
}

/**
 * Creates an isolated simpleGit client with scoped access to the global Git config.
 * @param {object} [options]
 * @param {string|null} [options.globalConfigPath] - Optional global config path
 * @param {string} [options.baseDir] - Repository base directory
 * @returns {import('simple-git').SimpleGit}
 */
export function createGitClient({
  globalConfigPath = resolveGlobalGitConfig(),
  baseDir = process.cwd(),
} = {}) {
  const options = {
    baseDir,
    unsafe: { allowUnsafeConfigPaths: true },
    allowEnvironment: ["GIT_CONFIG_GLOBAL"],
  };
  const git = simpleGit(options);
  if (globalConfigPath) {
    git.env("GIT_CONFIG_GLOBAL", globalConfigPath);
  }
  return git;
}

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
  if (isProductFile(filePath)) {
    return true;
  }
  const normalized = filePath.replace(/^[./]+/, "");
  if (normalized.endsWith("_test.go")) {
    return false;
  }
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
 * Enforces strict alignment with CI validate-commit-message rules.
 * @param {object} params
 * @param {string} params.message
 * @param {boolean} [params.affectsProduct]
 * @param {Array<string>} [params.changedFiles]
 * @param {boolean} [params.isMerge]
 * @returns {{valid: boolean, reason?: string}}
 */
export function validateCommitMessage({
  message,
  affectsProduct,
  changedFiles,
  isMerge = false,
  productRoots = DEFAULT_PRODUCT_ROOTS,
}) {
  if (!message || !message.trim()) {
    return { valid: false, reason: "Missing commit message." };
  }

  const subject = message.split("\n")[0].trim().replace(/\r$/, "");
  if (isMerge || /^Merge\s+/i.test(subject)) {
    return { valid: true };
  }

  try {
    validateEmpty(subject);
    validatePrefix(subject);
    validateLength(subject);

    const hasProductChanges =
      changedFiles !== undefined
        ? changedFiles.some((f) => fileAffectsProduct(f, productRoots))
        : Boolean(affectsProduct);

    const files = hasProductChanges
      ? ["rancher2/provider.go"]
      : ["docs/README.md"];

    validateProductImpact(subject, files);
    return { valid: true };
  } catch (err) {
    return { valid: false, reason: err.message };
  }
}

/**
 * Validates a single commit title against Conventional Commits and CI rules.
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
 * @param {boolean} [params.affectsProduct]
 * @param {Array<string>} [params.changedFiles]
 * @param {Array<string>} [params.productRoots]
 * @param {Function} [params.checkSpellingFn]
 * @returns {Promise<{valid: boolean, reason?: string}>}
 */
export async function validateCommit({
  message,
  affectsProduct,
  changedFiles,
  productRoots = DEFAULT_PRODUCT_ROOTS,
  checkSpellingFn = checkSpelling,
}) {
  const syncResult = validateCommitMessage({
    message,
    affectsProduct,
    changedFiles,
    productRoots,
  });
  if (!syncResult.valid) {
    return syncResult;
  }

  const subject = message.split("\n")[0].trim().replace(/\r$/, "");
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
 * @param {string|null} [params.globalConfigPath]
 * @param {Function} [params.resolveGlobalConfigFn]
 * @returns {Promise<object>} Commit summary
 */
export async function createCommit({
  git = null,
  message,
  checkSpellingFn = checkSpelling,
  globalConfigPath = undefined,
  resolveGlobalConfigFn = resolveGlobalGitConfig,
}) {
  const gitClient =
    git ||
    createGitClient({
      globalConfigPath:
        globalConfigPath !== undefined
          ? globalConfigPath
          : resolveGlobalConfigFn(),
    });

  const currentBranch = (await gitClient.branch()).current;
  if (!currentBranch || currentBranch === "HEAD") {
    throw new Error("Cannot commit on a detached HEAD.");
  }

  if (currentBranch === "main" || currentBranch === "master") {
    throw new Error(
      `Commits directly to the '${currentBranch}' branch are not allowed.`,
    );
  }

  const diffSummary = await gitClient.diffSummary([
    "--no-ext-diff",
    "--cached",
  ]);
  if (diffSummary.files.length === 0) {
    const unstaged = await gitClient.diffSummary(["--no-ext-diff"]);
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

  const changedFiles = diffSummary.files.map((f) => f.file);
  const roots = await getProductRoots();
  const affectsProduct = diffSummary.files.some((f) =>
    fileAffectsProduct(f.file, roots),
  );
  const validation = await validateCommit({
    message,
    affectsProduct,
    changedFiles,
    productRoots: roots,
    checkSpellingFn,
  });
  if (!validation.valid) {
    throw new Error(validation.reason);
  }

  const result = await gitClient.commit(message, undefined, {
    "-s": null,
    "-S": null,
  });
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
