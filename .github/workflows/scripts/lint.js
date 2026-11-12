import { execSync, execFileSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import fs from "node:fs";
import path from "node:path";

export function runTerraform(exec = execSync) {
  console.log("==> Linting Terraform files...");
  try {
    exec("terraform fmt -check -recursive -diff", { stdio: "inherit" });
  } catch (error) {
    console.error(
      "terraform format error, please run 'terraform fmt -recursive'",
    );
    throw error;
  }
  exec("tflint --recursive", { stdio: "inherit" });
}

export function runActionlint(exec = execSync) {
  console.log("==> Linting GitHub workflows...");
  exec("actionlint", { stdio: "inherit" });
}

export function getShellcheckFiles(dir = ".", fsModule = fs) {
  const ignoredDirs = new Set([
    ".git",
    ".terraform",
    ".agent",
    "bin",
    "node_modules",
    ".github-tools",
  ]);
  const files = [];

  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (ignoredDirs.has(entry.name)) {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (entry.isFile()) {
        if (entry.name.endsWith(".md") || entry.name.endsWith(".js")) {
          continue;
        }
        try {
          const fd = fsModule.openSync(fullPath, "r");
          const buffer = Buffer.alloc(2);
          fsModule.readSync(fd, buffer, 0, 2, 0);
          fsModule.closeSync(fd);
          if (buffer.toString("utf-8") === "#!") {
            files.push(fullPath);
          }
        } catch {
          // Ignore unreadable files
        }
      }
    }
  }

  scan(dir);
  return files;
}

export function runShellcheck(exec = execFileSync, fsModule = fs) {
  console.log("==> Running shellcheck...");
  const files = getShellcheckFiles(".", fsModule);

  if (files.length === 0) {
    console.log("No shell scripts found to check.");
    return;
  }

  for (const file of files) {
    if (fsModule.existsSync(file)) {
      console.log(`Checking ${file}...`);
      exec("shellcheck", ["-x", file], { stdio: "inherit" });
    }
  }
}

export function getJsFiles(dir = ".", fsModule = fs) {
  const files = [];
  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (
        entry.name === "node_modules" ||
        entry.name === ".git" ||
        entry.name === ".github-tools"
      ) {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (entry.isFile() && entry.name.endsWith(".js")) {
        files.push(fullPath);
      }
    }
  }
  scan(dir);
  return files;
}

export function runNodeCheck(exec = execFileSync, fsModule = fs) {
  console.log("==> Running Node Syntax Check...");
  const files = getJsFiles(".", fsModule);

  if (files.length === 0) {
    console.log("No Node files found to check.");
    return;
  }

  for (const file of files) {
    if (fsModule.existsSync(file)) {
      console.log(`Checking ${file}...`);
      exec("node", ["--check", file], { stdio: "inherit" });
    }
  }
}

export function getTestFiles(dir = ".", fsModule = fs) {
  const files = [];
  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (
        entry.name === "node_modules" ||
        entry.name === ".git" ||
        entry.name === ".github-tools"
      ) {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (
        entry.isFile() &&
        (entry.name.endsWith(".test.js") || entry.name.endsWith(".test.mjs"))
      ) {
        files.push(fullPath);
      }
    }
  }
  scan(dir);
  return files;
}

export function runNodeTest(exec = execFileSync, fsModule = fs) {
  console.log("==> Running Node Unit Tests...");
  if (fsModule.existsSync?.(".gemini/package.json")) {
    console.log("==> Installing .gemini dependencies...");
    exec("npm", ["ci", "--prefix", ".gemini"], { stdio: "inherit" });
    if (fsModule.existsSync?.(".gemini/orchestrator/install-sdk.js")) {
      console.log("==> Installing @google/gemini-cli-sdk...");
      exec("node", [".gemini/orchestrator/install-sdk.js"], {
        stdio: "inherit",
      });
    }
  }
  const files = getTestFiles(".", fsModule);

  if (files.length === 0) {
    console.log("No test files found to run.");
    return;
  }

  exec("node", ["--test", ...files], { stdio: "inherit" });
}

export function runEslint(exec = execSync, fsModule = fs) {
  console.log("==> Running ESLint...");
  exec("npm ci", { stdio: "inherit" });
  const configArg =
    !fsModule.existsSync("eslint.config.js") &&
    fsModule.existsSync(".github-tools/eslint.config.js")
      ? "--config .github-tools/eslint.config.js "
      : "";
  exec(`eslint ${configArg}.`, { stdio: "inherit" });
}

export function runGitleaks(exec = execSync) {
  console.log("==> Scanning for secrets with gitleaks...");
  exec("gitleaks detect --no-banner -v --no-git", { stdio: "inherit" });
  exec("gitleaks detect --no-banner -v", { stdio: "inherit" });
}

export function compareSemver(v1, v2) {
  const p1 = String(v1).replace(/^v/, "").split("-")[0].split(".").map(Number);
  const p2 = String(v2).replace(/^v/, "").split("-")[0].split(".").map(Number);
  for (let i = 0; i < 3; i++) {
    const num1 = p1[i] || 0;
    const num2 = p2[i] || 0;
    if (num1 > num2) return 1;
    if (num1 < num2) return -1;
  }
  return 0;
}

export function checkManifestDowngrade({
  baseBranch,
  baseManifestContent,
  headManifestContent,
}) {
  if (!baseBranch || !/^release\/v\d+$/.test(baseBranch.trim())) {
    return { valid: true, skipped: true, reason: "not-release-branch" };
  }
  if (!baseManifestContent || !headManifestContent) {
    return { valid: true, skipped: true, reason: "missing-manifest" };
  }

  let baseManifest = {};
  let headManifest = {};
  try {
    baseManifest = JSON.parse(baseManifestContent);
    headManifest = JSON.parse(headManifestContent);
  } catch (e) {
    throw new Error(`Failed to parse manifest JSON: ${e.message}`);
  }

  const baseVersion = baseManifest["."];
  const headVersion = headManifest["."];

  if (!baseVersion || !headVersion) {
    return { valid: true, skipped: true, reason: "missing-version" };
  }

  if (compareSemver(headVersion, baseVersion) < 0) {
    const error = `❌ Manifest version downgrade detected: PR proposes '${headVersion}' but base branch is '${baseVersion}'. Backport PRs must not downgrade release manifests.`;
    return { valid: false, error, headVersion, baseVersion };
  }

  return { valid: true, headVersion, baseVersion };
}

export function runManifestCheck(
  exec = execFileSync,
  fsModule = fs,
  env = process.env,
) {
  const baseBranch = (env.GITHUB_BASE_REF || env.BASE_BRANCH || "").trim();
  if (!baseBranch || !/^release\/v\d+$/.test(baseBranch)) {
    console.log(
      "==> Skipping manifest downgrade check (not targeting release/v* branch)...",
    );
    return;
  }

  console.log(
    `==> Checking manifest version downgrade against base branch '${baseBranch}'...`,
  );

  const manifestPath = ".release-please-manifest.json";
  if (!fsModule.existsSync(manifestPath)) {
    console.log("No .release-please-manifest.json found. Skipping.");
    return;
  }

  const headManifestContent = fsModule.readFileSync(manifestPath, "utf8");

  let baseManifestContent = "";
  try {
    baseManifestContent = exec(
      "git",
      ["show", `origin/${baseBranch}:${manifestPath}`],
      { encoding: "utf8", stdio: ["pipe", "pipe", "pipe"] },
    );
  } catch {
    try {
      baseManifestContent = exec(
        "git",
        ["show", `${baseBranch}:${manifestPath}`],
        { encoding: "utf8", stdio: ["pipe", "pipe", "pipe"] },
      );
    } catch {
      console.log(
        `Could not load ${manifestPath} from base branch '${baseBranch}'. Skipping.`,
      );
      return;
    }
  }

  const result = checkManifestDowngrade({
    baseBranch,
    baseManifestContent,
    headManifestContent,
  });

  if (!result.valid) {
    console.error(result.error);
    throw new Error(result.error);
  }

  console.log(
    `✅ Manifest check passed: version '${result.headVersion}' >= '${result.baseVersion}' on ${baseBranch}.`,
  );
}

export function runAll(
  exec = execSync,
  fsModule = fs,
  fileExec = execFileSync,
) {
  runTerraform(exec);
  runActionlint(exec);
  runShellcheck(fileExec, fsModule);
  runNodeCheck(fileExec, fsModule);
  runNodeTest(fileExec, fsModule);
  runEslint(exec);
  runGitleaks(exec);
  runManifestCheck(fileExec, fsModule);
}

export function main(
  mode = process.argv[2] || "all",
  exec = execSync,
  fsModule = fs,
  fileExec = execFileSync,
) {
  const activeFileExec =
    fileExec !== execFileSync
      ? fileExec
      : exec !== execSync
        ? exec
        : execFileSync;

  switch (mode) {
  case "terraform":
    runTerraform(exec);
    break;
  case "actionlint":
    runActionlint(exec);
    break;
  case "shellcheck":
    runShellcheck(activeFileExec, fsModule);
    break;
  case "node-check":
    runNodeCheck(activeFileExec, fsModule);
    break;
  case "node-test":
    runNodeTest(activeFileExec, fsModule);
    break;
  case "eslint":
    runEslint(exec);
    break;
  case "gitleaks":
    runGitleaks(exec);
    break;
  case "manifest":
    runManifestCheck(activeFileExec, fsModule);
    break;
  case "all":
    runAll(exec, fsModule, activeFileExec);
    break;
  default:
    console.error(`Error: Unknown lint mode: ${mode}`);
    console.error(
      "Usage: lint.js [terraform|actionlint|shellcheck|node-check|node-test|eslint|gitleaks|manifest|all]",
    );
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  fileURLToPath(import.meta.url) === path.resolve(process.argv[1])
) {
  main();
}
