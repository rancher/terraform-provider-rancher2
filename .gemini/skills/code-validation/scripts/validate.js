#!/usr/bin/env node
import path from "node:path";
import { fileURLToPath } from "node:url";
import { runFormat } from "./format.js";
import { runLint, SUPPORTED_MODES as LINT_MODES } from "./lint.js";
import { runTest, SUPPORTED_MODES as TEST_MODES } from "./test.js";

/**
 * Prints help and usage information for validate.js.
 */
export function printHelp() {
  console.log(`Usage: validate.js [options]

Unified code validation harness executing formatting, static analysis, and testing.

Options:
  -f, --fix              Auto-format code before running validation suites
  -a, --all              Run all linter checks and all test suites (lint: all, test: all)
  -l, --lint-mode <mode> Specify static analysis mode (default: 'all')
  -t, --test-mode <mode> Specify test suite mode (default: 'unit')
  --skip-lint            Skip static analysis / linting step
  --skip-test            Skip test execution step
  -h, --help             Show this help message

Supported Lint Modes:
  ${LINT_MODES.join(", ")}

Supported Test Modes:
  ${TEST_MODES.join(", ")}
`);
}

/**
 * Parses arguments for validate.js.
 * @param {Array<string>} rawArgs
 * @returns {{help?: boolean, fix?: boolean, all?: boolean, lintMode?: string, testMode?: string, skipLint?: boolean, skipTest?: boolean, error?: string}}
 */
export function parseValidateArgs(rawArgs = []) {
  if (rawArgs.includes("-h") || rawArgs.includes("--help")) {
    return { help: true };
  }

  let fix = false;
  let all = false;
  let lintMode = "all";
  let testMode = "unit";
  let skipLint = false;
  let skipTest = false;

  for (let i = 0; i < rawArgs.length; i++) {
    const arg = rawArgs[i];
    if (arg === "-f" || arg === "--fix") {
      fix = true;
    } else if (arg === "-a" || arg === "--all") {
      all = true;
      lintMode = "all";
      testMode = "all";
    } else if (arg === "--skip-lint") {
      skipLint = true;
    } else if (arg === "--skip-test") {
      skipTest = true;
    } else if (arg === "-l" || arg === "--lint-mode") {
      if (i + 1 < rawArgs.length && !rawArgs[i + 1].startsWith("-")) {
        lintMode = rawArgs[++i];
      } else {
        return { error: "Option --lint-mode requires a value." };
      }
    } else if (arg.startsWith("--lint-mode=")) {
      lintMode = arg.slice("--lint-mode=".length) || "all";
    } else if (arg === "-t" || arg === "--test-mode") {
      if (i + 1 < rawArgs.length && !rawArgs[i + 1].startsWith("-")) {
        testMode = rawArgs[++i];
      } else {
        return { error: "Option --test-mode requires a value." };
      }
    } else if (arg.startsWith("--test-mode=")) {
      testMode = arg.slice("--test-mode=".length) || "unit";
    } else if (arg.startsWith("-")) {
      return { error: `Unknown option: '${arg}'` };
    } else {
      return { error: `Unexpected argument: '${arg}'` };
    }
  }

  lintMode = lintMode || "all";
  testMode = testMode || "unit";

  if (!skipLint && !LINT_MODES.includes(lintMode)) {
    return {
      error: `Unknown lint mode: '${lintMode}'. Supported modes: ${LINT_MODES.join(", ")}`,
    };
  }

  if (!skipTest && !TEST_MODES.includes(testMode)) {
    return {
      error: `Unknown test mode: '${testMode}'. Supported modes: ${TEST_MODES.join(", ")}`,
    };
  }

  return {
    fix,
    all,
    lintMode,
    testMode,
    skipLint,
    skipTest,
  };
}

/**
 * Runs the complete code validation pipeline asynchronously.
 *
 * @param {object} [options]
 * @param {boolean} [options.fix=false] - Whether to auto-format before running validation
 * @param {string} [options.lintMode='all'] - The linting mode to execute
 * @param {string} [options.testMode='unit'] - The test mode to execute
 * @param {boolean} [options.skipLint=false] - Skip static analysis
 * @param {boolean} [options.skipTest=false] - Skip tests
 * @param {string} [options.cwd] - Working directory (defaults to repo root)
 * @param {'inherit'|'pipe'} [options.stdio='inherit'] - Child process stdio configuration
 * @param {Function} [options.commandRunner] - Custom command runner function
 * @returns {Promise<{success: boolean, steps: {format?: object, lint?: object, test?: object}}>}
 */
export async function runValidate({
  fix = false,
  lintMode = "all",
  testMode = "unit",
  skipLint = false,
  skipTest = false,
  cwd = null,
  stdio = "inherit",
  commandRunner = null,
} = {}) {
  const steps = {};
  let overallSuccess = true;

  if (fix) {
    if (stdio === "inherit") {
      console.log("--- Step 1: Auto-Formatting Codebase ---");
    }
    const formatOpts = { target: "all", check: false, cwd, stdio };
    if (commandRunner) {
      formatOpts.commandRunner = commandRunner;
    }
    const formatResult = await runFormat(formatOpts);
    steps.format = formatResult;
    if (!formatResult.success) {
      overallSuccess = false;
    }
  }

  if (!skipLint) {
    if (stdio === "inherit") {
      console.log(
        `--- Step ${fix ? "2" : "1"}: Static Analysis [mode: ${lintMode}] ---`,
      );
    }
    const lintOpts = { mode: lintMode, fix: false, cwd, stdio };
    if (commandRunner) {
      lintOpts.commandRunner = commandRunner;
    }
    const lintResult = await runLint(lintOpts);
    steps.lint = lintResult;
    if (!lintResult.success) {
      overallSuccess = false;
    }
  }

  if (!skipTest) {
    const testStepNum = (fix ? 1 : 0) + (!skipLint ? 1 : 0) + 1;
    if (stdio === "inherit") {
      console.log(
        `--- Step ${testStepNum}: Test Suite [mode: ${testMode}] ---`,
      );
    }
    const testOpts = { mode: testMode, cwd, stdio };
    if (commandRunner) {
      testOpts.commandRunner = commandRunner;
    }
    const testResult = await runTest(testOpts);
    steps.test = testResult;
    if (!testResult.success) {
      overallSuccess = false;
    }
  }

  return {
    success: overallSuccess,
    steps,
  };
}

/**
 * Formats a clean summary table of validation results.
 * @param {object} steps
 * @param {boolean} overallSuccess
 * @returns {string}
 */
export function formatSummary(steps, overallSuccess) {
  const lines = [
    "==================================================",
    "Code Validation Pipeline Summary",
    "==================================================",
  ];

  if (steps.format) {
    lines.push(
      `  Format:   ${steps.format.success ? "✅ PASSED" : "❌ FAILED"}`,
    );
  } else {
    lines.push("  Format:   ⚪ SKIPPED");
  }

  if (steps.lint) {
    lines.push(
      `  Lint:     ${steps.lint.success ? "✅ PASSED" : "❌ FAILED"} [mode: ${steps.lint.mode}]`,
    );
  } else {
    lines.push("  Lint:     ⚪ SKIPPED");
  }

  if (steps.test) {
    lines.push(
      `  Tests:    ${steps.test.success ? "✅ PASSED" : "❌ FAILED"} [mode: ${steps.test.mode}]`,
    );
  } else {
    lines.push("  Tests:    ⚪ SKIPPED");
  }

  lines.push("--------------------------------------------------");
  lines.push(
    `Result:     ${overallSuccess ? "✅ ALL CHECKS PASSED" : "❌ VALIDATION FAILED"}`,
  );
  lines.push("==================================================");

  return lines.join("\n");
}

async function main() {
  const args = process.argv.slice(2);
  const parsed = parseValidateArgs(args);

  if (parsed.help) {
    printHelp();
    process.exit(0);
  }

  if (parsed.error) {
    console.error(`❌ Error: ${parsed.error}`);
    printHelp();
    process.exit(1);
  }

  const { fix, lintMode, testMode, skipLint, skipTest } = parsed;

  console.log("🚀 Starting code validation pipeline...");
  const { success, steps } = await runValidate({
    fix,
    lintMode,
    testMode,
    skipLint,
    skipTest,
    stdio: "inherit",
  });

  console.log("\n" + formatSummary(steps, success));

  process.exit(success ? 0 : 1);
}

if (
  process.argv[1] &&
  path.resolve(process.argv[1]) === path.resolve(fileURLToPath(import.meta.url))
) {
  main().catch((err) => {
    console.error(
      `❌ Fatal error in validation pipeline: ${err?.message || err}`,
    );
    process.exit(1);
  });
}
