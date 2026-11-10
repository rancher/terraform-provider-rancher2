import assert from "node:assert";
import test from "node:test";
import { getPRComments, sys } from "./github-context.js";
import { getRepoRoot } from "./utils.js";

test("github-context - getPRComments", async (t) => {
  await t.test(
    "executes read_comments.js with arguments and returns formatted comments",
    async () => {
      let capturedCmd = null;
      let capturedArgs = null;
      let capturedOptions = null;

      t.mock.method(sys, "execFileAsync", async (cmd, args, options) => {
        capturedCmd = cmd;
        capturedArgs = args;
        capturedOptions = options;
        return {
          stdout: "=== Comments for PR #448 ===\nSome mock comment text\n",
        };
      });

      const comments = await getPRComments(448);
      assert.strictEqual(typeof comments, "string");
      assert.match(comments, /Comments for PR #448/);
      assert.strictEqual(capturedCmd, process.execPath);
      const repoRoot = await getRepoRoot();
      assert.ok(capturedArgs[0].startsWith(repoRoot));
      assert.ok(capturedArgs[0].endsWith("read_comments.js"));
      assert.strictEqual(capturedArgs[1], "448");
      assert.strictEqual(capturedOptions.cwd, repoRoot);
      assert.strictEqual(capturedOptions.maxBuffer, 10 * 1024 * 1024);
    },
  );

  await t.test(
    "falls back to 'current' when prTarget is omitted or null",
    async () => {
      let capturedArgs = null;
      t.mock.method(sys, "execFileAsync", async (cmd, args) => {
        capturedArgs = args;
        return {
          stdout: "=== Comments for PR current ===\nMock comment\n",
        };
      });

      await getPRComments();
      assert.strictEqual(capturedArgs[1], "current");

      await getPRComments(null);
      assert.strictEqual(capturedArgs[1], "current");
    },
  );

  await t.test(
    "propagates error message when child process fails",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with code 1");
        err.stderr = "API rate limit exceeded";
        throw err;
      });

      await assert.rejects(
        async () => {
          await getPRComments("448");
        },
        {
          message:
            /Failed to fetch PR comments programmatically: API rate limit exceeded/,
        },
      );
    },
  );

  await t.test(
    "falls back to stdout diagnostics when stderr is empty or whitespace",
    async () => {
      t.mock.method(sys, "execFileAsync", async () => {
        const err = new Error("Process exited with code 1");
        err.stderr = "   \n\t  ";
        err.stdout = "Diagnostic output from read_comments.js";
        throw err;
      });

      await assert.rejects(
        async () => {
          await getPRComments("448");
        },
        {
          message:
            /Failed to fetch PR comments programmatically: Diagnostic output from read_comments\.js/,
        },
      );
    },
  );
});
