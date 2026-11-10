#!/usr/bin/env node
import path from "node:path";
import { fileURLToPath } from "node:url";
import { getRepoRoot, runCommand } from "./utils.js";

export const SUPPORTED_MODES = [
  "all",
  "terraform",
  "actionlint",
  "eslint",
  "shellcheck",
  "shfmt",
  "prettier",
  "markdownlint",
  "golangci-lint",
  "tests",
  "gitleaks",
  "cspell",
];

/**
 * Prints help and usage information for lint.js.
 */
export function printHelp() {
  console.log(`Usage: lint.js [mode] [-f|--fix]
   or: lint.js [-m|--mode <mode>] [-f|--fix]

Executes static analysis and formatting verification suites via standard interfaces.

Modes:
  all             Execute all verification suites via 'make lint' (default)
  terraform       Lint Terraform configurations
  actionlint      Lint GitHub Actions workflow files
  eslint          Lint JavaScript/MJS codebase files
  shellcheck      Lint shell script code
  shfmt           Verify shell script formatting
  prettier        Verify formatting of JS, JSON, and MD files
  markdownlint    Validate Markdown semantic structure
  golangci-lint   Check Go source code against static analysis rules
  tests           Run Go test-specific linters
  gitleaks        Scan repository for hardcoded secrets
  cspell          Verify spelling across codebase files

Options:
  -f, --fix       Auto-format code where supported (Prettier, shfmt, golangci-lint)
  -m, --mode      Explicitly specify linting mode
  -h, --help      Show this help message
`);
}

/**
 * Parses arguments for lint.js.
 * @param {Array<string>} rawArgs
 * @returns {{help?: boolean, mode?: string, fix?: boolean, error?: string}}
 */
export function parseLintArgs(rawArgs = []) {
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    return { help: true };
  }

  let fix = false;
  let mode = null;

  for (let i = 0; i < rawArgs.length; i++) {
    const arg = rawArgs[i];
    if (arg === "-f" || arg === "--fix") {
      fix = true;
    } else if (arg === "-m" || arg === "--mode") {
      if (i + 1 < rawArgs.length && !rawArgs[i + 1].startsWith("-")) {
        mode = rawArgs[++i];
      } else {
        return { error: "Option --mode requires a value." };
      }
    } else if (arg.startsWith("--mode=")) {
      mode = arg.slice("--mode=".length);
    } else if (arg.startsWith("-")) {
      return { error: `Unknown option: '${arg}'` };
    } else if (!arg.startsWith("-")) {
      if (mode === null) {
        mode = arg;
      } else {
        return { error: `Unexpected argument: '${arg}'` };
      }
    }
  }

  mode = mode || "all";

  if (!SUPPORTED_MODES.includes(mode)) {
    return {
      error: `Unknown lint mode: '${mode}'. Supported modes: ${SUPPORTED_MODES.join(", ")}`,
    };
  }

  return { mode, fix };
}

/**
 * Executes the linting harness asynchronously.
 *
 * @param {object} [options]
 * @param {string} [options.mode='all'] - The linting mode to execute
 * @param {boolean} [options.fix=false] - Whether to apply fixes where supported
 * @param {string} [options.cwd] - Working directory (defaults to repo root)
 * @param {'inherit'|'pipe'} [options.stdio='inherit'] - Child process stdio configuration
 * @param {Function} [options.commandRunner=runCommand] - Custom command runner function
 * @returns {Promise<{success: boolean, code: number, stdout: string, stderr: string, mode: string, fix: boolean}>}
 */
export async function runLint({
  mode = "all",
  fix = false,
  cwd = null,
  stdio = "inherit",
  commandRunner = runCommand,
} = {}) {
  const repoRoot = await getRepoRoot();
  const workingDir = cwd || repoRoot;

  if (mode === "all") {
    if (fix) {
      const fmtRes = await commandRunner("make", ["fmt"], {
        cwd: workingDir,
        stdio,
      });
      const lintRes = await commandRunner("make", ["lint"], {
        cwd: workingDir,
        stdio,
      });
      return {
        success: fmtRes.success && lintRes.success,
        code: lintRes.code || fmtRes.code || 0,
        stdout: [fmtRes.stdout, lintRes.stdout].filter(Boolean).join("\n"),
        stderr: [fmtRes.stderr, lintRes.stderr].filter(Boolean).join("\n"),
        mode,
        fix,
      };
    }
    const result = await commandRunner("make", ["lint"], {
      cwd: workingDir,
      stdio,
    });
    return { ...result, mode, fix };
  }

  let result;
  switch (mode) {
  case "golangci-lint":
    result = await commandRunner(
      "golangci-lint",
      fix ? ["run", "--fix"] : ["run"],
      { cwd: workingDir, stdio },
    );
    break;
  case "eslint":
    result = await commandRunner(
      "npx",
      [
        "eslint",
        ".github/workflows/scripts/",
        ".gemini/hooks/",
        ".gemini/orchestrator/",
        ".gemini/skills/",
      ],
      { cwd: workingDir, stdio },
    );
    break;
  case "prettier":
    result = await commandRunner(
      "prettier",
      fix ? ["--write", "."] : ["--check", "."],
      { cwd: workingDir, stdio },
    );
    break;
  case "shfmt":
    result = await commandRunner("shfmt", fix ? ["-w", "."] : ["-d", "."], {
      cwd: workingDir,
      stdio,
    });
    break;
  case "cspell":
    result = await commandRunner(
      "cspell",
      ["lint", "--no-progress", "**/*"],
      { cwd: workingDir, stdio },
    );
    break;
  case "gitleaks":
    result = await commandRunner(
      "gitleaks",
      ["detect", "--no-banner", "-v", "--no-git"],
      { cwd: workingDir, stdio },
    );
    break;
  case "markdownlint":
    result = await commandRunner(
      "markdownlint",
      ["**/*.md", "--ignore", "node_modules"],
      { cwd: workingDir, stdio },
    );
    break;
  case "actionlint":
    result = await commandRunner("actionlint", [], {
      cwd: workingDir,
      stdio,
    });
    break;
  case "terraform":
    result = await commandRunner("terraform", ["fmt", fix ? "." : "-check"], {
      cwd: workingDir,
      stdio,
    });
    break;
  default:
    result = await commandRunner(
      "make",
      [mode === "all" ? "lint" : `lint-${mode}`],
      { cwd: workingDir, stdio },
    );
    break;
  }

  return {
    ...result,
    mode,
    fix,
  };
}

async function main() {
  const args = process.argv.slice(2);
  const parsed = parseLintArgs(args);

  if (parsed.help) {
    printHelp();
    process.exit(0);
  }

  if (parsed.error) {
    console.error(`❌ Error: ${parsed.error}`);
    printHelp();
    process.exit(1);
  }

  const { mode, fix } = parsed;
  console.log(
    `▶️ Running static analysis [mode: ${mode}${fix ? ", fix: true" : ""}]...`,
  );

  const result = await runLint({ mode, fix, stdio: "inherit" });

  if (result.success) {
    console.log(
      `✅ Static analysis passed [mode: ${mode}${fix ? ", fix: true" : ""}]`,
    );
    process.exit(0);
  } else {
    console.error(
      `❌ Static analysis failed [mode: ${mode}] (exit code: ${result.code})`,
    );
    process.exit(result.code);
  }
}

if (
  process.argv[1] &&
  path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))
) {
  main().catch((err) => {
    console.error(`❌ Fatal error in lint runner: ${err?.message || err}`);
    process.exit(1);
  });
}
