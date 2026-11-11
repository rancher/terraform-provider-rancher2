import assert from "node:assert";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import test from "node:test";

import {
  parseFormatArgs,
  runFormat,
  SUPPORTED_TARGETS,
} from "../scripts/format.js";
import {
  parseLintArgs,
  runLint,
  SUPPORTED_MODES as LINT_MODES,
} from "../scripts/lint.js";
import {
  findTestFiles,
  parseTestArgs,
  runTest,
  SUPPORTED_MODES as TEST_MODES,
  withIsolatedGitConfig,
} from "../scripts/test.js";
import { getRepoRoot, runCommand } from "../scripts/utils.js";
import {
  formatSummary,
  parseValidateArgs,
  runValidate,
} from "../scripts/validate.js";

test("code-validation skill - utils.js", async (t) => {
  await t.test(
    "getRepoRoot resolves repository root with package.json",
    async () => {
      const root = await getRepoRoot();
      assert.ok(typeof root === "string" && root.length > 0);
      const pkgStat = await fs.stat(path.join(root, "package.json"));
      assert.strictEqual(pkgStat.isFile(), true);
    },
  );

  await t.test(
    "runCommand executes command with stdio pipe and returns stdout",
    async () => {
      const res = await runCommand(
        "node",
        ["-e", 'process.stdout.write("hello world")'],
        {
          stdio: "pipe",
        },
      );
      assert.strictEqual(res.success, true);
      assert.strictEqual(res.code, 0);
      assert.strictEqual(res.stdout, "hello world");
    },
  );

  await t.test(
    "runCommand captures output when stdio is specified as an array",
    async () => {
      const res = await runCommand(
        "node",
        ["-e", 'process.stdout.write("out"); process.stderr.write("err")'],
        {
          stdio: ["ignore", "pipe", "pipe"],
        },
      );
      assert.strictEqual(res.success, true);
      assert.strictEqual(res.code, 0);
      assert.strictEqual(res.stdout, "out");
      assert.strictEqual(res.stderr, "err");
    },
  );

  await t.test(
    "runCommand handles non-zero exit codes gracefully without throwing",
    async () => {
      const res = await runCommand("node", ["-e", "process.exit(42)"], {
        stdio: "pipe",
      });
      assert.strictEqual(res.success, false);
      assert.strictEqual(res.code, 42);
    },
  );

  await t.test(
    "runCommand handles spawn error on non-existent command",
    async () => {
      const res = await runCommand("non_existent_command_123456789", [], {
        stdio: "pipe",
      });
      assert.strictEqual(res.success, false);
      assert.strictEqual(res.code, 1);
      assert.ok(res.stderr.length > 0);
    },
  );

  await t.test(
    "runCommand sets NO_COLOR to 1 even when empty string is provided",
    async () => {
      const origNoColor = process.env.NO_COLOR;
      process.env.NO_COLOR = "";
      try {
        const res = await runCommand(
          "node",
          ["-e", 'process.stdout.write(process.env.NO_COLOR || "")'],
          {
            stdio: "pipe",
          },
        );
        assert.strictEqual(res.success, true);
        assert.strictEqual(res.stdout, "1");
      } finally {
        if (origNoColor === undefined) {
          delete process.env.NO_COLOR;
        } else {
          process.env.NO_COLOR = origNoColor;
        }
      }
    },
  );
});

test("code-validation skill - parseLintArgs", async (t) => {
  await t.test("defaults to all mode without fix", () => {
    const res = parseLintArgs([]);
    assert.strictEqual(res.mode, "all");
    assert.strictEqual(res.fix, false);
  });

  await t.test("handles help flags", () => {
    assert.deepStrictEqual(parseLintArgs(["-h"]), { help: true });
    assert.deepStrictEqual(parseLintArgs(["--help"]), { help: true });
  });

  await t.test("parses positional mode argument", () => {
    const res = parseLintArgs(["eslint"]);
    assert.strictEqual(res.mode, "eslint");
    assert.strictEqual(res.fix, false);
  });

  await t.test("parses -f and --fix flags", () => {
    const res1 = parseLintArgs(["-f"]);
    assert.strictEqual(res1.fix, true);
    assert.strictEqual(res1.mode, "all");

    const res2 = parseLintArgs(["shfmt", "--fix"]);
    assert.strictEqual(res2.mode, "shfmt");
    assert.strictEqual(res2.fix, true);
  });

  await t.test("parses -m and --mode flags", () => {
    const res1 = parseLintArgs(["-m", "golangci-lint"]);
    assert.strictEqual(res1.mode, "golangci-lint");

    const res2 = parseLintArgs(["--mode", "cspell"]);
    assert.strictEqual(res2.mode, "cspell");

    const res3 = parseLintArgs(["--mode=prettier"]);
    assert.strictEqual(res3.mode, "prettier");
  });

  await t.test("returns error for missing --mode value", () => {
    const res = parseLintArgs(["-m"]);
    assert.ok(res.error?.includes("requires a value"));
  });

  await t.test(
    "returns error when flag immediately follows -m or --mode",
    () => {
      const res1 = parseLintArgs(["-m", "--fix"]);
      assert.ok(res1.error?.includes("requires a value"));
      const res2 = parseLintArgs(["--mode", "-f"]);
      assert.ok(res2.error?.includes("requires a value"));
    },
  );

  await t.test("returns error for unsupported mode", () => {
    const res = parseLintArgs(["invalid-mode"]);
    assert.ok(res.error?.includes("Unknown lint mode"));
  });

  await t.test("all supported modes are accepted", () => {
    for (const mode of LINT_MODES) {
      const res = parseLintArgs([mode]);
      assert.strictEqual(res.mode, mode);
    }
  });

  await t.test("returns error for unknown options", () => {
    const res = parseLintArgs(["--unknown-flag"]);
    assert.ok(res.error?.includes("Unknown option"));
  });
});

test("code-validation skill - parseTestArgs", async (t) => {
  await t.test("defaults to unit mode", () => {
    const res = parseTestArgs([]);
    assert.strictEqual(res.mode, "unit");
  });

  await t.test("handles help flags", () => {
    assert.deepStrictEqual(parseTestArgs(["-h"]), { help: true });
    assert.deepStrictEqual(parseTestArgs(["--help"]), { help: true });
  });

  await t.test("parses positional mode argument", () => {
    const res = parseTestArgs(["all"]);
    assert.strictEqual(res.mode, "all");
  });

  await t.test("parses -m and --mode flags", () => {
    const res1 = parseTestArgs(["-m", "workflow-scripts"]);
    assert.strictEqual(res1.mode, "workflow-scripts");

    const res2 = parseTestArgs(["--mode", "agent-scripts"]);
    assert.strictEqual(res2.mode, "agent-scripts");

    const res3 = parseTestArgs(["--mode=compile"]);
    assert.strictEqual(res3.mode, "compile");
  });

  await t.test("returns error for missing --mode value", () => {
    const res = parseTestArgs(["--mode"]);
    assert.ok(res.error?.includes("requires a value"));
  });

  await t.test(
    "returns error when flag immediately follows -m or --mode",
    () => {
      const res1 = parseTestArgs(["-m", "--all"]);
      assert.ok(res1.error?.includes("requires a value"));
      const res2 = parseTestArgs(["--mode", "-m"]);
      assert.ok(res2.error?.includes("requires a value"));
    },
  );

  await t.test("returns error for unsupported test mode", () => {
    const res = parseTestArgs(["non-existent-test"]);
    assert.ok(res.error?.includes("Unknown test mode"));
  });

  await t.test("all supported test modes are accepted", () => {
    for (const mode of TEST_MODES) {
      const res = parseTestArgs([mode]);
      assert.strictEqual(res.mode, mode);
    }
  });

  await t.test("returns error for unknown options", () => {
    const res = parseTestArgs(["--unknown-flag"]);
    assert.ok(res.error?.includes("Unknown option"));
  });
});

test("code-validation skill - findTestFiles", async (t) => {
  await t.test(
    "returns discovered test files in deterministic sorted order",
    async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "find-test-files-"),
      );
      try {
        await fs.writeFile(path.join(tmpDir, "z.test.js"), "");
        await fs.writeFile(path.join(tmpDir, "a.test.js"), "");
        await fs.writeFile(path.join(tmpDir, "m.spec.js"), "");
        await fs.writeFile(path.join(tmpDir, "not-a-test.txt"), "");

        const files = await findTestFiles(tmpDir);
        assert.strictEqual(files.length, 3);
        assert.deepStrictEqual(files, [
          path.join(tmpDir, "a.test.js"),
          path.join(tmpDir, "m.spec.js"),
          path.join(tmpDir, "z.test.js"),
        ]);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    },
  );
});

test("code-validation skill - parseFormatArgs", async (t) => {
  await t.test("defaults to all target without check", () => {
    const res = parseFormatArgs([]);
    assert.strictEqual(res.target, "all");
    assert.strictEqual(res.check, false);
  });

  await t.test("handles help flags", () => {
    assert.deepStrictEqual(parseFormatArgs(["-h"]), { help: true });
    assert.deepStrictEqual(parseFormatArgs(["--help"]), { help: true });
  });

  await t.test("parses positional target argument", () => {
    const res = parseFormatArgs(["go"]);
    assert.strictEqual(res.target, "go");
    assert.strictEqual(res.check, false);
  });

  await t.test("parses -c and --check flags", () => {
    const res1 = parseFormatArgs(["-c"]);
    assert.strictEqual(res1.check, true);
    assert.strictEqual(res1.target, "all");

    const res2 = parseFormatArgs(["prettier", "--check"]);
    assert.strictEqual(res2.target, "prettier");
    assert.strictEqual(res2.check, true);
  });

  await t.test("parses -t and --target flags", () => {
    const res1 = parseFormatArgs(["-t", "shfmt"]);
    assert.strictEqual(res1.target, "shfmt");

    const res2 = parseFormatArgs(["--target", "golangci-lint"]);
    assert.strictEqual(res2.target, "golangci-lint");

    const res3 = parseFormatArgs(["--target=gofmt"]);
    assert.strictEqual(res3.target, "gofmt");
  });

  await t.test("returns error for missing --target value", () => {
    const res = parseFormatArgs(["-t"]);
    assert.ok(res.error?.includes("requires a value"));
  });

  await t.test(
    "returns error when flag immediately follows -t or --target",
    () => {
      const res1 = parseFormatArgs(["-t", "--check"]);
      assert.ok(res1.error?.includes("requires a value"));
      const res2 = parseFormatArgs(["--target", "-c"]);
      assert.ok(res2.error?.includes("requires a value"));
    },
  );

  await t.test("returns error for unsupported format target", () => {
    const res = parseFormatArgs(["rustfmt"]);
    assert.ok(res.error?.includes("Unknown format target"));
  });

  await t.test("all supported format targets are accepted", () => {
    for (const target of SUPPORTED_TARGETS) {
      const res = parseFormatArgs([target]);
      assert.strictEqual(res.target, target);
    }
  });

  await t.test("returns error for unknown options", () => {
    const res = parseFormatArgs(["--unknown-flag"]);
    assert.ok(res.error?.includes("Unknown option"));
  });
});

test("code-validation skill - parseValidateArgs", async (t) => {
  await t.test("defaults to lint all and test unit without fix", () => {
    const res = parseValidateArgs([]);
    assert.strictEqual(res.fix, false);
    assert.strictEqual(res.all, false);
    assert.strictEqual(res.lintMode, "all");
    assert.strictEqual(res.testMode, "unit");
    assert.strictEqual(res.skipLint, false);
    assert.strictEqual(res.skipTest, false);
  });

  await t.test("handles help flags", () => {
    assert.deepStrictEqual(parseValidateArgs(["-h"]), { help: true });
    assert.deepStrictEqual(parseValidateArgs(["--help"]), { help: true });
  });

  await t.test("parses -a and --all flags", () => {
    const res = parseValidateArgs(["--all"]);
    assert.strictEqual(res.all, true);
    assert.strictEqual(res.lintMode, "all");
    assert.strictEqual(res.testMode, "all");
  });

  await t.test("parses -f and --fix flags", () => {
    const res = parseValidateArgs(["-f"]);
    assert.strictEqual(res.fix, true);
  });

  await t.test("parses custom lint and test modes", () => {
    const res = parseValidateArgs(["-l", "eslint", "-t", "workflow-scripts"]);
    assert.strictEqual(res.lintMode, "eslint");
    assert.strictEqual(res.testMode, "workflow-scripts");
  });

  await t.test("parses skip flags", () => {
    const res = parseValidateArgs(["--skip-lint", "--skip-test"]);
    assert.strictEqual(res.skipLint, true);
    assert.strictEqual(res.skipTest, true);
  });

  await t.test("returns error for invalid lint mode", () => {
    const res = parseValidateArgs(["-l", "non-existent"]);
    assert.ok(res.error?.includes("Unknown lint mode"));
  });

  await t.test("returns error for invalid test mode", () => {
    const res = parseValidateArgs(["-t", "non-existent"]);
    assert.ok(res.error?.includes("Unknown test mode"));
  });

  await t.test(
    "returns error when flag immediately follows -l or -t in parseValidateArgs",
    () => {
      const res1 = parseValidateArgs(["-l", "--all"]);
      assert.ok(res1.error?.includes("requires a value"));
      const res2 = parseValidateArgs(["-t", "--skip-lint"]);
      assert.ok(res2.error?.includes("requires a value"));
    },
  );

  await t.test(
    "falls back gracefully when --lint-mode= or --test-mode= has empty value",
    () => {
      const res = parseValidateArgs(["--lint-mode=", "--test-mode="]);
      assert.strictEqual(res.lintMode, "all");
      assert.strictEqual(res.testMode, "unit");
    },
  );

  await t.test(
    "falls back gracefully when empty strings are passed positionally for -l and -t",
    () => {
      const res = parseValidateArgs(["-l", "", "-t", ""]);
      assert.strictEqual(res.lintMode, "all");
      assert.strictEqual(res.testMode, "unit");
    },
  );

  await t.test("returns error for unknown options", () => {
    const res = parseValidateArgs(["--unknown-flag"]);
    assert.ok(res.error?.includes("Unknown option"));
  });

  await t.test("returns error for unexpected positional arguments", () => {
    const res = parseValidateArgs(["eslint"]);
    assert.ok(res.error?.includes("Unexpected argument"));
  });
});

test("code-validation skill - formatSummary", async (t) => {
  await t.test("formats passing summary", () => {
    const steps = {
      format: { success: true },
      lint: { success: true, mode: "all" },
      test: { success: true, mode: "unit" },
    };
    const summary = formatSummary(steps, true);
    assert.ok(summary.includes("Format:   ✅ PASSED"));
    assert.ok(summary.includes("Lint:     ✅ PASSED [mode: all]"));
    assert.ok(summary.includes("Tests:    ✅ PASSED [mode: unit]"));
    assert.ok(summary.includes("✅ ALL CHECKS PASSED"));
  });

  await t.test("formats failing summary with skipped steps", () => {
    const steps = {
      lint: { success: false, mode: "eslint" },
    };
    const summary = formatSummary(steps, false);
    assert.ok(summary.includes("Format:   ⚪ SKIPPED"));
    assert.ok(summary.includes("Lint:     ❌ FAILED [mode: eslint]"));
    assert.ok(summary.includes("Tests:    ⚪ SKIPPED"));
    assert.ok(summary.includes("❌ VALIDATION FAILED"));
  });
});

test("code-validation skill - execution integration", async (t) => {
  await t.test(
    "runValidate executes with skipped steps successfully",
    async () => {
      const res = await runValidate({
        skipLint: true,
        skipTest: true,
        fix: false,
        stdio: "pipe",
      });
      assert.strictEqual(res.success, true);
      assert.strictEqual(res.steps.format, undefined);
      assert.strictEqual(res.steps.lint, undefined);
      assert.strictEqual(res.steps.test, undefined);
    },
  );

  await t.test(
    "runValidate executes format step when fix is true with mock commandRunner",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args) => {
        executed.push({ cmd, args });
        return { success: true, code: 0, stdout: "ok", stderr: "" };
      };
      const res = await runValidate({
        fix: true,
        skipLint: true,
        skipTest: true,
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.success, true);
      assert.ok(res.steps.format !== undefined);
      assert.strictEqual(res.steps.format.target, "all");
      assert.strictEqual(executed.length, 4);
    },
  );

  await t.test(
    "runValidate reports success: false when a sub-step fails with mock commandRunner",
    async () => {
      const mockRunner = async () => ({
        success: false,
        code: 2,
        stdout: "",
        stderr: "Lint error",
      });
      const res = await runValidate({
        skipLint: false,
        skipTest: true,
        lintMode: "eslint",
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.success, false);
      assert.strictEqual(res.steps.lint.success, false);
      assert.strictEqual(res.steps.lint.code, 2);
    },
  );

  await t.test(
    "runLint delegates to commandRunner with proper args",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args) => {
        executed.push({ cmd, args });
        return { success: true, code: 0, stdout: "cspell passed", stderr: "" };
      };
      const res = await runLint({ mode: "cspell", commandRunner: mockRunner });
      assert.strictEqual(res.mode, "cspell");
      assert.strictEqual(res.success, true);
      assert.strictEqual(res.stdout, "cspell passed");
      assert.strictEqual(executed[0].cmd, "cspell");
      assert.deepStrictEqual(executed[0].args, [
        "lint",
        "--no-progress",
        "**/*",
      ]);
    },
  );

  await t.test("runLint mode all executes make lint", async () => {
    const executed = [];
    const mockRunner = async (cmd, args) => {
      executed.push({ cmd, args });
      return { success: true, code: 0, stdout: "lint clean", stderr: "" };
    };
    const res = await runLint({ mode: "all", commandRunner: mockRunner });
    assert.strictEqual(res.mode, "all");
    assert.strictEqual(res.success, true);
    assert.strictEqual(executed[0].cmd, "make");
    assert.deepStrictEqual(executed[0].args, ["lint"]);
  });

  await t.test("runTest delegates to make unit without recursion", async () => {
    const executed = [];
    const mockRunner = async (cmd, args) => {
      executed.push({ cmd, args });
      return { success: true, code: 0, stdout: "tests passed", stderr: "" };
    };
    const res = await runTest({ mode: "unit", commandRunner: mockRunner });
    assert.strictEqual(res.mode, "unit");
    assert.strictEqual(res.success, true);
    assert.strictEqual(res.stdout, "tests passed");
    assert.strictEqual(executed[0].cmd, "make");
    assert.deepStrictEqual(executed[0].args, ["unit"]);
  });

  await t.test("runTest mode acc delegates to make acc", async () => {
    const executed = [];
    const mockRunner = async (cmd, args) => {
      executed.push({ cmd, args });
      return { success: true, code: 0, stdout: "acc tests passed", stderr: "" };
    };
    const res = await runTest({ mode: "acc", commandRunner: mockRunner });
    assert.strictEqual(res.mode, "acc");
    assert.strictEqual(res.success, true);
    assert.strictEqual(executed[0].cmd, "make");
    assert.deepStrictEqual(executed[0].args, ["acc"]);
  });

  await t.test(
    "runTest mode workflow-scripts isolates GIT_CONFIG_GLOBAL to unique temp file when /dev/null",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const executed = [];
        const mockRunner = async (cmd, args, opts) => {
          executed.push({ cmd, args, opts });
          return {
            success: true,
            code: 0,
            stdout: "workflow tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "workflow-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "workflow-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(executed[0].opts?.env?.GIT_CONFIG_GLOBAL);
        assert.notStrictEqual(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          "/dev/null",
        );
        assert.match(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          /gemini-gitconfig-[^/]+\/\.gitconfig$/,
        );
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
      }
    },
  );

  await t.test(
    "runTest mode workflow-scripts unlinks temporary .gitconfig if created by runner",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      let tempPath = null;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const mockRunner = async (cmd, args, opts) => {
          tempPath = opts?.env?.GIT_CONFIG_GLOBAL;
          if (tempPath) {
            await fs.writeFile(tempPath, "[user]\n  name = Test\n", "utf8");
          }
          return {
            success: true,
            code: 0,
            stdout: "workflow tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "workflow-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "workflow-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(tempPath);
        await assert.rejects(async () => fs.stat(tempPath));
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
        if (tempPath) {
          try {
            await fs.unlink(tempPath);
          } catch {
            // Ignore
          }
        }
      }
    },
  );

  await t.test(
    "runTest mode agent-scripts isolates GIT_CONFIG_GLOBAL to unique temp file when /dev/null",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const executed = [];
        const mockRunner = async (cmd, args, opts) => {
          executed.push({ cmd, args, opts });
          return {
            success: true,
            code: 0,
            stdout: "agent tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "agent-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "agent-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(executed[0].opts?.env?.GIT_CONFIG_GLOBAL);
        assert.notStrictEqual(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          "/dev/null",
        );
        assert.match(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          /gemini-gitconfig-[^/]+\/\.gitconfig$/,
        );
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
      }
    },
  );

  await t.test(
    "runTest mode agent-scripts unlinks temporary .gitconfig if created by runner",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      let tempPath = null;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const mockRunner = async (cmd, args, opts) => {
          tempPath = opts?.env?.GIT_CONFIG_GLOBAL;
          if (tempPath) {
            await fs.writeFile(tempPath, "[user]\n  name = Test\n", "utf8");
          }
          return {
            success: true,
            code: 0,
            stdout: "agent tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "agent-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "agent-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(tempPath);
        await assert.rejects(async () => fs.stat(tempPath));
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
        if (tempPath) {
          try {
            await fs.unlink(tempPath);
          } catch {
            // Ignore
          }
        }
      }
    },
  );

  await t.test(
    "runTest mode release-scripts executes test files directly via process.execPath",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args) => {
        executed.push({ cmd, args });
        return { success: true, code: 0, stdout: "tests passed", stderr: "" };
      };
      const res = await runTest({
        mode: "release-scripts",
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.mode, "release-scripts");
      assert.strictEqual(res.success, true);
      assert.strictEqual(executed[0].cmd, process.execPath);
      assert.strictEqual(executed[0].args[0], "--test");
      assert.ok(executed[0].args.length > 1);
      assert.ok(executed[0].args.some((f) => f.includes("scripts")));
    },
  );

  await t.test(
    "runTest mode release-scripts isolates GIT_CONFIG_GLOBAL to unique temp file when /dev/null",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const executed = [];
        const mockRunner = async (cmd, args, opts) => {
          executed.push({ cmd, args, opts });
          return {
            success: true,
            code: 0,
            stdout: "release tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "release-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "release-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(executed[0].opts?.env?.GIT_CONFIG_GLOBAL);
        assert.notStrictEqual(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          "/dev/null",
        );
        assert.match(
          executed[0].opts.env.GIT_CONFIG_GLOBAL,
          /gemini-gitconfig-[^/]+\/\.gitconfig$/,
        );
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
      }
    },
  );

  await t.test(
    "runTest mode release-scripts unlinks temporary .gitconfig if created by runner",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      let tempPath = null;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const mockRunner = async (cmd, args, opts) => {
          tempPath = opts?.env?.GIT_CONFIG_GLOBAL;
          if (tempPath) {
            await fs.writeFile(tempPath, "[user]\n  name = Test\n", "utf8");
          }
          return {
            success: true,
            code: 0,
            stdout: "release tests passed",
            stderr: "",
          };
        };
        const res = await runTest({
          mode: "release-scripts",
          commandRunner: mockRunner,
        });
        assert.strictEqual(res.mode, "release-scripts");
        assert.strictEqual(res.success, true);
        assert.ok(tempPath);
        await assert.rejects(async () => fs.stat(tempPath));
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
        if (tempPath) {
          try {
            await fs.unlink(tempPath);
          } catch {
            // Ignore
          }
        }
      }
    },
  );

  await t.test(
    "withIsolatedGitConfig preserves environment when GIT_CONFIG_GLOBAL is not /dev/null",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/custom/path/.gitconfig";
        let capturedEnv = null;
        const result = await withIsolatedGitConfig(async (testEnv) => {
          capturedEnv = testEnv;
          return "custom-success";
        });
        assert.strictEqual(result, "custom-success");
        assert.strictEqual(
          capturedEnv?.GIT_CONFIG_GLOBAL,
          "/custom/path/.gitconfig",
        );
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
      }
    },
  );

  await t.test(
    "withIsolatedGitConfig creates isolated temp directory and cleans up on completion",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      let tempConfigPath = null;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        const result = await withIsolatedGitConfig(async (testEnv) => {
          tempConfigPath = testEnv.GIT_CONFIG_GLOBAL;
          assert.match(tempConfigPath, /gemini-gitconfig-[^/]+\/\.gitconfig$/);
          const initialStat = await fs.stat(tempConfigPath);
          assert.strictEqual(initialStat.isFile(), true);
          assert.strictEqual(await fs.readFile(tempConfigPath, "utf8"), "");
          await fs.writeFile(tempConfigPath, "[user]\n  name = Test\n", "utf8");
          const stat = await fs.stat(tempConfigPath);
          assert.strictEqual(stat.isFile(), true);
          return "isolated-ok";
        });
        assert.strictEqual(result, "isolated-ok");
        assert.ok(tempConfigPath);
        await assert.rejects(async () => fs.stat(tempConfigPath));
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
        if (tempConfigPath) {
          try {
            await fs.rm(path.dirname(tempConfigPath), {
              recursive: true,
              force: true,
            });
          } catch {
            // Ignore
          }
        }
      }
    },
  );

  await t.test(
    "withIsolatedGitConfig cleans up temporary directory even if callback throws",
    async () => {
      const origGitConfig = process.env.GIT_CONFIG_GLOBAL;
      let tempConfigPath = null;
      try {
        process.env.GIT_CONFIG_GLOBAL = "/dev/null";
        await assert.rejects(
          async () => {
            await withIsolatedGitConfig(async (testEnv) => {
              tempConfigPath = testEnv.GIT_CONFIG_GLOBAL;
              await fs.writeFile(
                tempConfigPath,
                "[user]\n  name = Test\n",
                "utf8",
              );
              throw new Error("Callback failure");
            });
          },
          { message: "Callback failure" },
        );
        assert.ok(tempConfigPath);
        await assert.rejects(async () => fs.stat(tempConfigPath));
      } finally {
        if (origGitConfig !== undefined) {
          process.env.GIT_CONFIG_GLOBAL = origGitConfig;
        } else {
          delete process.env.GIT_CONFIG_GLOBAL;
        }
        if (tempConfigPath) {
          try {
            await fs.rm(path.dirname(tempConfigPath), {
              recursive: true,
              force: true,
            });
          } catch {
            // Ignore
          }
        }
      }
    },
  );

  await t.test(
    "runFormat normalizes gofmt target to go with mock commandRunner",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args) => {
        executed.push({ cmd, args });
        return { success: true, code: 0, stdout: "", stderr: "" };
      };
      const res = await runFormat({
        target: "gofmt",
        check: true,
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.target, "gofmt");
      assert.strictEqual(res.results.length, 1);
      assert.strictEqual(res.results[0].target, "go");
      assert.strictEqual(executed[0].cmd, "gofmt");
      assert.deepStrictEqual(executed[0].args, ["-l", "."]);
    },
  );

  await t.test(
    "runFormat executes prettier directly without npx wrapper",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args) => {
        executed.push({ cmd, args });
        return { success: true, code: 0, stdout: "clean", stderr: "" };
      };
      const res = await runFormat({
        target: "prettier",
        check: true,
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.target, "prettier");
      assert.strictEqual(res.success, true);
      assert.strictEqual(executed.length, 1);
      assert.strictEqual(executed[0].cmd, "prettier");
      assert.deepStrictEqual(executed[0].args, ["--check", "."]);
    },
  );

  await t.test("runFormat executes shfmt with commandRunner", async () => {
    const executed = [];
    const mockRunner = async (cmd, args) => {
      executed.push({ cmd, args });
      return { success: true, code: 0, stdout: "clean", stderr: "" };
    };
    const res = await runFormat({
      target: "shfmt",
      check: true,
      commandRunner: mockRunner,
    });
    assert.strictEqual(res.target, "shfmt");
    assert.strictEqual(res.success, true);
    assert.strictEqual(executed.length, 1);
  });

  await t.test(
    "runLint resolves script from getRepoRoot while passing custom cwd to commandRunner",
    async () => {
      const executed = [];
      const mockRunner = async (cmd, args, opts) => {
        executed.push({ cmd, args, opts });
        return { success: true, code: 0, stdout: "ok", stderr: "" };
      };
      const customCwd = "/custom/fixture/dir";
      const res = await runLint({
        mode: "all",
        cwd: customCwd,
        commandRunner: mockRunner,
      });
      assert.strictEqual(res.success, true);
      assert.strictEqual(executed.length, 1);
      assert.strictEqual(executed[0].cmd, "make");
      assert.deepStrictEqual(executed[0].args, ["lint"]);
      assert.strictEqual(executed[0].opts.cwd, customCwd);
    },
  );
});
