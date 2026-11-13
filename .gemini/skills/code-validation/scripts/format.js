#!/usr/bin/env node
import path from "node:path";
import { fileURLToPath } from "node:url";
import { getRepoRoot, runCommand } from "./utils.js";

export const SUPPORTED_TARGETS = [
  "all",
  "go",
  "gofmt",
  "prettier",
  "shfmt",
  "golangci-lint",
];

/**
 * Prints help and usage information for format.js.
 */
export function printHelp() {
  console.log(`Usage: format.js [target] [--check]
   or: format.js [-t|--target <target>] [-c|--check]

Executes automatic code formatters across the codebase.

Targets:
  all             Format all file types (Go, JS/TS/JSON/MD, shell scripts) (default)
  go, gofmt       Format Go code using gofmt
  prettier        Format JS, JSON, and Markdown files using Prettier
  shfmt           Format shell scripts using shfmt
  golangci-lint   Apply automatic fixes via golangci-lint

Options:
  -c, --check     Verify code formatting without modifying files (dry-run)
  -t, --target    Explicitly specify formatting target
  -h, --help      Show this help message
`);
}

/**
 * Parses arguments for format.js.
 * @param {Array<string>} rawArgs
 * @returns {{help?: boolean, target?: string, check?: boolean, error?: string}}
 */
export function parseFormatArgs(rawArgs = []) {
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    return { help: true };
  }

  let check = false;
  let target = null;

  for (let i = 0; i < rawArgs.length; i++) {
    const arg = rawArgs[i];
    if (arg === "-c" || arg === "--check") {
      check = true;
    } else if (arg === "-t" || arg === "--target") {
      if (i + 1 < rawArgs.length && !rawArgs[i + 1].startsWith("-")) {
        target = rawArgs[++i];
      } else {
        return { error: "Option --target requires a value." };
      }
    } else if (arg.startsWith("--target=")) {
      target = arg.slice("--target=".length);
    } else if (arg.startsWith("-")) {
      return { error: `Unknown option: '${arg}'` };
    } else if (!arg.startsWith("-")) {
      if (target === null) {
        target = arg;
      } else {
        return { error: `Unexpected argument: '${arg}'` };
      }
    }
  }

  target = target || "all";

  if (!SUPPORTED_TARGETS.includes(target)) {
    return {
      error: `Unknown format target: '${target}'. Supported targets: ${SUPPORTED_TARGETS.join(", ")}`,
    };
  }

  return { target, check };
}

/**
 * Executes code formatting tools asynchronously.
 *
 * @param {object} [options]
 * @param {string} [options.target='all'] - The formatting target to execute
 * @param {boolean} [options.check=false] - Whether to verify without writing modifications
 * @param {string} [options.cwd] - Working directory (defaults to repo root)
 * @param {'inherit'|'pipe'} [options.stdio='inherit'] - Child process stdio configuration
 * @param {Function} [options.commandRunner=runCommand] - Custom command runner function
 * @returns {Promise<{success: boolean, code: number, results: Array<{target: string, success: boolean, code: number}>, target: string, check: boolean}>}
 */
export async function runFormat({
  target = "all",
  check = false,
  cwd = null,
  stdio = "inherit",
  commandRunner = runCommand,
} = {}) {
  if (target !== "all" && !SUPPORTED_TARGETS.includes(target)) {
    return {
      success: false,
      code: 1,
      results: [
        {
          target,
          success: false,
          code: 1,
          stderr: `Unsupported target: ${target}`,
        },
      ],
      target,
      check,
    };
  }

  const repoRoot = await getRepoRoot();
  const workingDir = cwd || repoRoot;

  const normalized = target === "gofmt" ? "go" : target;
  const toolsToRun =
    normalized === "all"
      ? ["go", "prettier", "shfmt", "golangci-lint"]
      : [normalized];

  const results = [];
  let overallSuccess = true;

  for (const tool of toolsToRun) {
    let res;
    if (tool === "go") {
      if (check) {
        res = await commandRunner("gofmt", ["-l", "."], {
          cwd: workingDir,
          stdio: "pipe",
        });
        if (res.success && res.stdout.trim().length > 0) {
          res.success = false;
          res.code = 1;
        }
        if (stdio === "inherit") {
          const out = res.stdout ? res.stdout.trim() : "";
          if (out) {
            console.log(out);
          }
          const err = res.stderr ? res.stderr.trim() : "";
          if (err) {
            console.error(err);
          }
        }
      } else {
        res = await commandRunner("gofmt", ["-s", "-w", "-e", "."], {
          cwd: workingDir,
          stdio,
        });
      }
    } else if (tool === "prettier") {
      const prettierArgs = check ? ["--check", "."] : ["--write", "."];
      res = await commandRunner("prettier", prettierArgs, {
        cwd: workingDir,
        stdio,
      });
    } else if (tool === "shfmt") {
      const shfmtArgs = check
        ? ["-i", "2", "-ci", "-bn", "-d", "."]
        : ["-i", "2", "-ci", "-bn", "-w", "."];
      res = await commandRunner("shfmt", shfmtArgs, { cwd: workingDir, stdio });
    } else if (tool === "golangci-lint") {
      const lintArgs = check ? ["run"] : ["run", "--fix"];
      res = await commandRunner("golangci-lint", lintArgs, {
        cwd: workingDir,
        stdio,
      });
    }

    const toolSuccess = res.success;
    results.push({
      target: tool,
      success: toolSuccess,
      code: res.code,
      stdout: res.stdout,
      stderr: res.stderr,
    });

    if (!toolSuccess) {
      overallSuccess = false;
    }
  }

  return {
    success: overallSuccess,
    code: overallSuccess ? 0 : 1,
    results,
    target,
    check,
  };
}

async function main() {
  const args = process.argv.slice(2);
  const parsed = parseFormatArgs(args);

  if (parsed.help) {
    printHelp();
    process.exit(0);
  }

  if (parsed.error) {
    console.error(`❌ Error: ${parsed.error}`);
    printHelp();
    process.exit(1);
  }

  const { target, check } = parsed;
  const actionName = check ? "Verifying format of" : "Formatting";
  console.log(`▶️ ${actionName} [target: ${target}]...`);

  const result = await runFormat({ target, check, stdio: "inherit" });

  if (result.success) {
    console.log(
      `✅ Formatting ${check ? "verification passed" : "completed successfully"} [target: ${target}]`,
    );
    process.exit(0);
  } else {
    console.error(
      `❌ Formatting ${check ? "verification failed" : "failed"} [target: ${target}]`,
    );
    process.exit(result.code);
  }
}

if (
  process.argv[1] &&
  path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))
) {
  main().catch((err) => {
    console.error(`❌ Fatal error in format runner: ${err?.message || err}`);
    process.exit(1);
  });
}
