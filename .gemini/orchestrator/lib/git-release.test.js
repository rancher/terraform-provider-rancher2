import assert from "node:assert";
import test from "node:test";
import {
  affectsProductFiles,
  getGitDiff,
  stageAndCommit,
  stripDiffMetadata,
  sys,
  validateMessage,
} from "./git-release.js";
import { getRepoRoot } from "./utils.js";

test("stripDiffMetadata", async (t) => {
  await t.test(
    "filters index lines and chunk header @@ lines with line number shifts",
    () => {
      const diff1 = [
        "diff --git a/file.txt b/file.txt",
        "index 1234567..89abcdef 100644",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -10,3 +10,4 @@",
        " context line",
        "+added line",
        " context line 2",
      ].join("\n");

      const diff2 = [
        "diff --git a/file.txt b/file.txt",
        "index fedcba9..7654321 100644",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -45,3 +45,4 @@",
        " context line",
        "+added line",
        " context line 2",
      ].join("\n");

      const stripped1 = stripDiffMetadata(diff1);
      const stripped2 = stripDiffMetadata(diff2);

      assert.strictEqual(stripped1, stripped2);
      assert.doesNotMatch(stripped1, /^index /m);
      assert.match(stripped1, /^@@ @@$/m);
    },
  );

  await t.test(
    "retains hunk location context while ignoring line number shifts",
    () => {
      const diffFoo1 = [
        "diff --git a/file.txt b/file.txt",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -10,3 +10,4 @@ func Foo()",
        " context line",
        "+added line",
        " context line 2",
      ].join("\n");

      const diffFoo2 = [
        "diff --git a/file.txt b/file.txt",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -80,3 +80,4 @@ func Foo()",
        " context line",
        "+added line",
        " context line 2",
      ].join("\n");

      const diffBar = [
        "diff --git a/file.txt b/file.txt",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -10,3 +10,4 @@ func Bar()",
        " context line",
        "+added line",
        " context line 2",
      ].join("\n");

      assert.strictEqual(
        stripDiffMetadata(diffFoo1),
        stripDiffMetadata(diffFoo2),
      );
      assert.notStrictEqual(
        stripDiffMetadata(diffFoo1),
        stripDiffMetadata(diffBar),
      );
    },
  );

  await t.test(
    "preserves whitespace in diff content lines and trailing newlines",
    () => {
      const diffWithTrailingSpace = [
        "diff --git a/file.txt b/file.txt",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -1,2 +1,2 @@",
        "+added line with spaces   ",
        " context line",
        "",
      ].join("\n");

      const diffWithoutTrailingSpace = [
        "diff --git a/file.txt b/file.txt",
        "--- a/file.txt",
        "+++ b/file.txt",
        "@@ -1,2 +1,2 @@",
        "+added line with spaces",
        " context line",
        "",
      ].join("\n");

      const strippedWith = stripDiffMetadata(diffWithTrailingSpace);
      const strippedWithout = stripDiffMetadata(diffWithoutTrailingSpace);

      assert.match(strippedWith, /\+added line with spaces {3}\n/);
      assert.notStrictEqual(strippedWith, strippedWithout);
      assert.ok(strippedWith.endsWith("\n"));
    },
  );
});

test("getGitDiff and affectsProductFiles", async (t) => {
  await t.test("getGitDiff returns a string and passes repoRoot cwd", async () => {
    let capturedOptions = null;
    t.mock.method(sys, "execAsync", async (cmd, options) => {
      capturedOptions = options;
      return { stdout: "diff --git a/test b/test\n" };
    });

    const diff = await getGitDiff();
    assert.strictEqual(typeof diff, "string");
    const repoRoot = await getRepoRoot();
    assert.strictEqual(capturedOptions?.cwd, repoRoot);
    assert.strictEqual(capturedOptions?.maxBuffer, 10 * 1024 * 1024);
  });

  await t.test("affectsProductFiles returns a boolean and passes repoRoot cwd", async () => {
    let capturedOptions = null;
    t.mock.method(sys, "execAsync", async (cmd, options) => {
      capturedOptions = options;
      return { stdout: "test/file.go\0" };
    });

    const isProduct = await affectsProductFiles();
    assert.strictEqual(typeof isProduct, "boolean");
    const repoRoot = await getRepoRoot();
    assert.strictEqual(capturedOptions?.cwd, repoRoot);
    assert.strictEqual(capturedOptions?.maxBuffer, 10 * 1024 * 1024);
  });

  await t.test("affectsProductFiles detects untracked product files when tracked diff is empty", async () => {
    t.mock.method(sys, "execAsync", async (cmd) => {
      if (cmd.includes("ls-files")) {
        return { stdout: "internal/resource_rancher2_new.go\0" };
      }
      return { stdout: "" };
    });

    const isProduct = await affectsProductFiles();
    assert.strictEqual(isProduct, true);
  });
});

test("stageAndCommit execution and process.execPath", async (t) => {
  await t.test("throws error when attempting to commit on main branch", async () => {
    t.mock.method(sys, "execAsync", async (cmd) => {
      if (cmd.includes("rev-parse")) {
        return { stdout: "main\n" };
      }
      return { stdout: "" };
    });

    await assert.rejects(
      async () => {
        await stageAndCommit("docs: update readme");
      },
      {
        message: 'Commits directly to the "main" or "master" branch are not allowed.',
      },
    );
  });

  await t.test("throws error when no modified or untracked paths are found", async () => {
    t.mock.method(sys, "execAsync", async (cmd) => {
      if (cmd.includes("rev-parse")) {
        return { stdout: "feature-branch\n" };
      }
      return { stdout: "" };
    });

    t.mock.method(sys, "execFileAsync", async () => {
      return { stdout: "" };
    });

    await assert.rejects(
      async () => {
        await stageAndCommit("docs: update readme");
      },
      {
        message: "No modified or untracked paths found to commit.",
      },
    );
  });

  await t.test(
    "executes get_modified_paths.js, stage-files.js, and commit.js with process.execPath",
    async () => {
      const execFileCalls = [];
      const execAsyncCalls = [];

      t.mock.method(sys, "execAsync", async (cmd, options) => {
        execAsyncCalls.push({ cmd, options });
        if (cmd.includes("rev-parse")) {
          return { stdout: "feature-branch\n" };
        }
        if (cmd.includes("diff -z")) {
          return { stdout: "docs/readme.md\0" };
        }
        return { stdout: "" };
      });

      t.mock.method(sys, "execFileAsync", async (file, args, options) => {
        execFileCalls.push({ file, args, options });
        const script = args[0] || "";
        if (script.includes("get_modified_paths.js")) {
          return { stdout: "docs/readme.md\n" };
        }
        if (script.includes("stage-files.js")) {
          return { stdout: "staged docs/readme.md\n" };
        }
        if (script.includes("commit.js")) {
          return {
            stdout: "[feature-branch 1234567] docs: update readme\n",
            stderr: "",
          };
        }
        return { stdout: "", stderr: "" };
      });

      const res = await stageAndCommit("docs: update readme");
      assert.ok(res.stdout.includes("[feature-branch 1234567]"));

      const repoRoot = await getRepoRoot();
      assert.strictEqual(execAsyncCalls.length >= 1, true);
      assert.strictEqual(execAsyncCalls[0].options?.cwd, repoRoot);
      assert.strictEqual(execFileCalls.length, 3);
      assert.strictEqual(execFileCalls[0].file, process.execPath);
      assert.ok(execFileCalls[0].args[0].endsWith("get_modified_paths.js"));
      assert.strictEqual(execFileCalls[0].options?.cwd, repoRoot);

      assert.strictEqual(execFileCalls[1].file, process.execPath);
      assert.ok(execFileCalls[1].args[0].endsWith("stage-files.js"));
      assert.strictEqual(execFileCalls[1].args[1], "--");
      assert.strictEqual(execFileCalls[1].args[2], "docs/readme.md");
      assert.strictEqual(execFileCalls[1].options?.cwd, repoRoot);

      assert.strictEqual(execFileCalls[2].file, process.execPath);
      assert.ok(execFileCalls[2].args[0].endsWith("commit.js"));
      assert.strictEqual(execFileCalls[2].options?.cwd, repoRoot);
      assert.deepStrictEqual(execFileCalls[2].args.slice(1), [
        "-m",
        "docs: update readme",
      ]);
    },
  );

  await t.test(
    "safely passes hyphen-prefixed paths with end-of-options delimiter to stage-files.js",
    async () => {
      const execFileCalls = [];

      t.mock.method(sys, "execAsync", async (cmd) => {
        if (cmd.includes("rev-parse")) {
          return { stdout: "feature-branch\n" };
        }
        if (cmd.includes("diff -z")) {
          return { stdout: "-f\0" };
        }
        return { stdout: "" };
      });

      t.mock.method(sys, "execFileAsync", async (file, args, options) => {
        execFileCalls.push({ file, args, options });
        const script = args[0] || "";
        if (script.includes("get_modified_paths.js")) {
          return { stdout: "-f\n-file.txt\n" };
        }
        if (script.includes("stage-files.js")) {
          return { stdout: "staged -f, -file.txt\n" };
        }
        if (script.includes("commit.js")) {
          return {
            stdout: "[feature-branch 1234567] docs: update readme\n",
            stderr: "",
          };
        }
        return { stdout: "", stderr: "" };
      });

      const res = await stageAndCommit("docs: update readme");
      assert.ok(res.stdout.includes("[feature-branch 1234567]"));

      const stageCall = execFileCalls.find((c) =>
        c.args[0].endsWith("stage-files.js"),
      );
      assert.ok(stageCall);
      assert.strictEqual(stageCall.args[1], "--");
      assert.strictEqual(stageCall.args[2], "-f");
      assert.strictEqual(stageCall.args[3], "-file.txt");
    },
  );
});

test("validateMessage", async (t) => {
  await t.test("strips commit message body and validates title", async () => {
    t.mock.method(sys, "execAsync", async (cmd) => {
      if (cmd.includes("diff -z")) {
        return { stdout: "docs/readme.md\0" };
      }
      return { stdout: "" };
    });

    const multiLineMessage = [
      "docs: update documentation for skills",
      "",
      "Detailed explanation of why skills were modernized.",
      "Signed-off-by: Test Developer <test@example.com>",
    ].join("\n");

    const result = await validateMessage(multiLineMessage);
    assert.strictEqual(result.valid, true);
  });
});

