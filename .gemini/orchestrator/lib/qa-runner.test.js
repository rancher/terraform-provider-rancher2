import assert from "node:assert";
import { Buffer } from "node:buffer";
import test from "node:test";

import { runQAPipeline } from "./qa-runner.js";

test("runQAPipeline execution flow and concurrency rules", async (t) => {
  await t.test(
    "runs test suite and linter commands concurrently when fix is false",
    async () => {
      const executedCommands = [];
      const mockExec = async (cmd) => {
        executedCommands.push(cmd);
        return { stdout: `ok: ${cmd}`, stderr: "" };
      };

      const result = await runQAPipeline({ fix: false, execFn: mockExec });

      assert.strictEqual(result.success, true);
      assert.deepStrictEqual(executedCommands.sort(), [
        `${process.execPath} .gemini/skills/code-validation/scripts/lint.js all`,
        `${process.execPath} .gemini/skills/code-validation/scripts/test.js all`,
      ]);
    },
  );

  await t.test(
    "runs fixer sequentially before validation suites when fix is true",
    async () => {
      const callOrder = [];
      let fixerCompleted = false;

      const mockExec = async (cmd) => {
        callOrder.push(cmd);
        if (cmd.includes("--fix")) {
          fixerCompleted = true;
          return { stdout: "fixed files", stderr: "" };
        }
        assert.strictEqual(
          fixerCompleted,
          true,
          "Validation commands must run only after fixer completes",
        );
        return { stdout: `passed: ${cmd}`, stderr: "" };
      };

      const result = await runQAPipeline({ fix: true, execFn: mockExec });

      assert.strictEqual(result.success, true);
      assert.strictEqual(
        callOrder[0],
        `${process.execPath} .gemini/skills/code-validation/scripts/lint.js all --fix`,
      );
      assert.strictEqual(callOrder.length, 3);
      assert.ok(result.stdout.includes("fixed files"));
    },
  );

  await t.test("supports custom fixerCommand when fix is true", async () => {
    const callOrder = [];
    const mockExec = async (cmd) => {
      callOrder.push(cmd);
      return { stdout: `ran: ${cmd}`, stderr: "" };
    };

    const result = await runQAPipeline({
      fix: true,
      execFn: mockExec,
      fixerCommand: "npm run lint -- --fix",
      commands: ["npm test"],
    });

    assert.strictEqual(result.success, true);
    assert.strictEqual(callOrder[0], "npm run lint -- --fix");
    assert.strictEqual(callOrder[1], "npm test");
  });

  await t.test(
    "handles fixer failure and reports failure without throwing",
    async () => {
      const mockExec = async (cmd) => {
        if (cmd.includes("--fix")) {
          const err = new Error("Fixer failed");
          err.stdout = "";
          err.stderr = "Syntax error preventing automatic fix";
          throw err;
        }
        return { stdout: "ok", stderr: "" };
      };

      const result = await runQAPipeline({ fix: true, execFn: mockExec });

      assert.strictEqual(result.success, false);
      assert.ok(
        result.stderr.includes("Syntax error preventing automatic fix"),
      );
    },
  );

  await t.test(
    "safely string-casts Buffer or non-string stdout and stderr error properties",
    async () => {
      const mockExec = async (cmd) => {
        if (cmd.includes("--fix")) {
          const err = new Error("Process exited with code 1");
          err.stdout = Buffer.from("fix stdout buffer");
          err.stderr = Buffer.from("fix stderr buffer");
          throw err;
        }
        const err = new Error("Validation failed");
        err.stdout = 12345;
        err.stderr = null;
        throw err;
      };

      const result = await runQAPipeline({ fix: true, execFn: mockExec });

      assert.strictEqual(result.success, false);
      assert.ok(result.stdout.includes("fix stdout buffer"));
      assert.ok(result.stderr.includes("fix stderr buffer"));
      assert.ok(result.stderr.includes("Validation failed"));
    },
  );

  await t.test(
    "safely handles nullish or primitive errors without throwing TypeError",
    async () => {
      const mockExec = async (cmd) => {
        if (cmd.includes("--fix")) {
          throw null;
        }
        throw undefined;
      };

      const result = await runQAPipeline({ fix: true, execFn: mockExec });

      assert.strictEqual(result.success, false);
      assert.strictEqual(typeof result.stdout, "string");
      assert.strictEqual(typeof result.stderr, "string");
    },
  );

  await t.test(
    "passes repository root working directory to runner options",
    async () => {
      const capturedOptions = [];
      const mockExec = async (cmd, options) => {
        capturedOptions.push({ cmd, options });
        return { stdout: "ok", stderr: "" };
      };

      const customCwd = "/custom/repo/root";
      const result = await runQAPipeline({
        fix: true,
        execFn: mockExec,
        cwd: customCwd,
      });

      assert.strictEqual(result.success, true);
      assert.strictEqual(capturedOptions.length, 3);
      for (const call of capturedOptions) {
        assert.strictEqual(call.options.cwd, customCwd);
        assert.strictEqual(call.options.maxBuffer, 10 * 1024 * 1024);
      }
    },
  );
});
