#!/usr/bin/env node
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { getRepoRoot, runCommand } from "./utils.js";

export const SUPPORTED_MODES = [
  "unit",
  "compile",
  "acc",
  "acc-relay",
  "workflow-scripts",
  "agent-scripts",
  "all",
];

export const MODE_TO_MAKE_TARGET = {
  unit: "unit",
  compile: "compile",
  acc: "acc",
  "acc-relay": "testaccrelay",
  all: "test-all",
};

/**
 * Recursively locates test files matching test conventions.
 * @param {string} dir
 * @returns {Promise<Array<string>>}
 */
export async function findTestFiles(dir) {
  const testFiles = [];
  try {
    const entries = await fs.readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        if (entry.name !== "node_modules" && entry.name !== ".git") {
          const subFiles = await findTestFiles(fullPath);
          testFiles.push(...subFiles);
        }
      } else if (
        entry.isFile() &&
        (entry.name.endsWith(".test.js") ||
          entry.name.endsWith(".test.mjs") ||
          entry.name.endsWith(".test.ts") ||
          entry.name.endsWith(".spec.js") ||
          entry.name.endsWith(".spec.mjs") ||
          entry.name.endsWith(".spec.ts"))
      ) {
        testFiles.push(fullPath);
      }
    }
  } catch {
    // Ignore unreadable or missing directories
  }
  return testFiles;
}

/**
 * Prints help and usage information for test.js.
 */
export function printHelp() {
  console.log(`Usage: test.js [mode]
   or: test.js [-m|--mode <mode>]

Executes test suites and verification harnesses via standard interfaces.

Modes:
  unit              Run unit tests (default, executes 'make unit')
  compile           Run compile check on provider and tests (executes 'make compile')
  acc               Run acceptance tests (executes 'make acc')
  acc-relay         Run AWS Test Relay acceptance tests (executes 'make testaccrelay')
  workflow-scripts  Run workflow script unit tests via Node test runner
  agent-scripts     Run agent script unit tests via Node test runner
  all               Run full test suite across all languages (executes 'make test-all')

Options:
  -m, --mode        Explicitly specify test mode
  -h, --help        Show this help message
`);
}

/**
 * Parses arguments for test.js.
 * @param {Array<string>} rawArgs
 * @returns {{help?: boolean, mode?: string, error?: string}}
 */
export function parseTestArgs(rawArgs = []) {
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    return { help: true };
  }

  let mode = null;

  for (let i = 0; i < rawArgs.length; i++) {
    const arg = rawArgs[i];
    if (arg === "-m" || arg === "--mode") {
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

  mode = mode || "unit";

  if (!SUPPORTED_MODES.includes(mode)) {
    return {
      error: `Unknown test mode: '${mode}'. Supported modes: ${SUPPORTED_MODES.join(", ")}`,
    };
  }

  return { mode };
}

/**
 * Executes the test harness asynchronously using standard GNUmakefile targets.
 *
 * @param {object} [options]
 * @param {string} [options.mode='unit'] - The test mode to execute
 * @param {string} [options.cwd] - Working directory (defaults to repo root)
 * @param {'inherit'|'pipe'} [options.stdio='inherit'] - Child process stdio configuration
 * @param {Function} [options.commandRunner=runCommand] - Custom command runner function
 * @returns {Promise<{success: boolean, code: number, stdout: string, stderr: string, mode: string}>}
 */
export async function runTest({
  mode = "unit",
  cwd = null,
  stdio = "inherit",
  commandRunner = runCommand,
} = {}) {
  const repoRoot = await getRepoRoot();
  const workingDir = cwd || repoRoot;

  if (mode === "agent-scripts") {
    const candidateDirs = [
      path.join(repoRoot, ".gemini/orchestrator"),
      path.join(repoRoot, "agent-scripts"),
      path.join(repoRoot, ".gemini/hooks"),
      path.join(repoRoot, ".gemini/skills"),
    ];
    const testFiles = [];
    for (const dir of candidateDirs) {
      const files = await findTestFiles(dir);
      testFiles.push(...files);
    }
    if (testFiles.length === 0) {
      return {
        success: false,
        code: 1,
        stdout: "",
        stderr: "No agent script test files found",
        mode,
      };
    }
    const result = await commandRunner(
      process.execPath,
      ["--test", ...testFiles],
      { cwd: workingDir, stdio },
    );
    return { ...result, mode };
  }

  if (mode === "workflow-scripts") {
    const candidateDirs = [
      path.join(repoRoot, ".github/workflows/scripts/tests"),
      path.join(repoRoot, ".github/workflows/scripts"),
    ];
    let testFiles = [];
    for (const dir of candidateDirs) {
      try {
        const files = await findTestFiles(dir);
        if (files.length > 0) {
          testFiles = files;
          break;
        }
      } catch {
        // Continue searching fallback directories
      }
    }
    if (testFiles.length === 0) {
      return {
        success: false,
        code: 1,
        stdout: "",
        stderr: "No workflow script test files found",
        mode,
      };
    }
    const testEnv = { ...process.env };
    if (testEnv.GIT_CONFIG_GLOBAL === "/dev/null") {
      testEnv.GIT_CONFIG_GLOBAL = path.join(os.tmpdir(), ".gitconfig");
    }
    const result = await commandRunner(
      process.execPath,
      ["--test", ...testFiles],
      { cwd: workingDir, stdio, env: testEnv },
    );
    return { ...result, mode };
  }

  const makeTarget = MODE_TO_MAKE_TARGET[mode] || mode;
  let result = await commandRunner("make", [makeTarget], {
    cwd: workingDir,
    stdio,
  });

  // Graceful fallback to core C1 targets if extended target is not defined in host Makefile
  const checkResult = await commandRunner("make", ["-n", makeTarget], {
    cwd: workingDir,
    stdio: "pipe",
  });
  if (
    !result.success &&
    !checkResult.success &&
    typeof checkResult.stderr === "string" &&
    checkResult.stderr.includes("No rule to make target")
  ) {
    if (mode === "compile") {
      result = await commandRunner("make", ["build"], {
        cwd: workingDir,
        stdio,
      });
    } else if (mode === "acc-relay") {
      result = await commandRunner("make", ["acc"], { cwd: workingDir, stdio });
    } else if (mode === "all") {
      result = await commandRunner("make", ["test"], {
        cwd: workingDir,
        stdio,
      });
    }
  }

  return {
    ...result,
    mode,
  };
}

async function main() {
  const args = process.argv.slice(2);
  const parsed = parseTestArgs(args);

  if (parsed.help) {
    printHelp();
    process.exit(0);
  }

  if (parsed.error) {
    console.error(`❌ Error: ${parsed.error}`);
    printHelp();
    process.exit(1);
  }

  const { mode } = parsed;
  console.log(`▶️ Running test suite [mode: ${mode}]...`);

  const result = await runTest({ mode, stdio: "inherit" });

  if (result.success) {
    console.log(`✅ Test suite passed [mode: ${mode}]`);
    process.exit(0);
  } else {
    console.error(
      `❌ Test suite failed [mode: ${mode}] (exit code: ${result.code})`,
    );
    process.exit(result.code);
  }
}

if (
  process.argv[1] &&
  path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))
) {
  main().catch((err) => {
    console.error(`❌ Fatal error in test runner: ${err?.message || err}`);
    process.exit(1);
  });
}
