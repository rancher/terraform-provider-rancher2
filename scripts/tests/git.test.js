import assert from "node:assert/strict";
import test from "node:test";

import {
  getChangedFilesForCommit,
  getChangedFilesFromCommit,
  isValidCommitSha,
  parseDiffTreeOutput,
  validateCommitSha,
} from "../lib/git.js";

test("parseDiffTreeOutput", async (t) => {
  await t.test("parses newline-delimited output", () => {
    const raw = "rancher2/resource_cluster.go\ndocs/resources/cluster.md\n\n";
    const files = parseDiffTreeOutput(raw);
    assert.deepStrictEqual(files, [
      "rancher2/resource_cluster.go",
      "docs/resources/cluster.md",
    ]);
  });

  await t.test("handles empty or null output", () => {
    assert.deepStrictEqual(parseDiffTreeOutput(""), []);
    assert.deepStrictEqual(parseDiffTreeOutput(null), []);
  });
});

test("getChangedFilesFromCommit", async (t) => {
  await t.test("returns file array when git succeeds", async () => {
    const mockExec = async (cmd, args) => {
      assert.strictEqual(cmd, "git");
      assert.deepStrictEqual(args, [
        "diff-tree",
        "--no-commit-id",
        "--name-only",
        "-r",
        "--",
        "abc1234",
      ]);
      return { stdout: "rancher2/provider.go\nCHANGELOG.md\n", stderr: "" };
    };

    const files = await getChangedFilesFromCommit("abc1234", {
      execFn: mockExec,
    });
    assert.deepStrictEqual(files, ["rancher2/provider.go", "CHANGELOG.md"]);
  });

  await t.test("returns null when git fails or commit not found", async () => {
    const mockExec = async () => {
      throw new Error("fatal: bad object");
    };

    const files = await getChangedFilesFromCommit("nonexistent", {
      execFn: mockExec,
    });
    assert.strictEqual(files, null);
  });

  await t.test(
    "returns null without executing git on null, undefined, empty, whitespace, or malformed commit SHA",
    async () => {
      let execCalled = false;
      const mockExec = async () => {
        execCalled = true;
        return { stdout: "file.txt\n", stderr: "" };
      };

      const invalidShas = [
        null,
        undefined,
        "",
        "   ",
        "\t\n",
        "abc 123",
        "abc;rm -rf",
        "--diff-filter=M",
        "abc$def",
        "abc\x00def",
      ];

      for (const sha of invalidShas) {
        execCalled = false;
        const res1 = await getChangedFilesFromCommit(sha, {
          execFn: mockExec,
        });
        assert.strictEqual(
          res1,
          null,
          `Expected null for getChangedFilesFromCommit(${JSON.stringify(sha)})`,
        );
        assert.strictEqual(
          execCalled,
          false,
          `execFn should not be called for ${JSON.stringify(sha)}`,
        );

        const res2 = await getChangedFilesForCommit(sha, {
          execFn: mockExec,
        });
        assert.strictEqual(
          res2,
          null,
          `Expected null for getChangedFilesForCommit(${JSON.stringify(sha)})`,
        );
        assert.strictEqual(
          execCalled,
          false,
          `execFn should not be called for ${JSON.stringify(sha)}`,
        );
      }
    },
  );
});

test("isValidCommitSha and validateCommitSha", async (t) => {
  await t.test("returns true for valid commit SHAs", () => {
    assert.strictEqual(isValidCommitSha("abc1234"), true);
    assert.strictEqual(
      isValidCommitSha("da39a3ee5e6b4b0d3255bfef95601890afd80709"),
      true,
    );
    assert.strictEqual(isValidCommitSha(" 1234567890abcdef "), true);
  });

  await t.test("returns false for invalid commit SHAs", () => {
    assert.strictEqual(isValidCommitSha(null), false);
    assert.strictEqual(isValidCommitSha(undefined), false);
    assert.strictEqual(isValidCommitSha(""), false);
    assert.strictEqual(isValidCommitSha("   "), false);
    assert.strictEqual(isValidCommitSha("abc 123"), false);
    assert.strictEqual(isValidCommitSha("--flag"), false);
    assert.strictEqual(isValidCommitSha("abc$123"), false);
    assert.strictEqual(isValidCommitSha("abc;rm -rf"), false);
  });

  await t.test(
    "throws Error with rejected value when throwOnError is true or via validateCommitSha",
    () => {
      assert.throws(
        () => isValidCommitSha("bad;sha", true),
        (err) => {
          assert.ok(err instanceof Error);
          assert.strictEqual(err.message, "Invalid commit SHA: bad;sha");
          return true;
        },
      );

      assert.throws(
        () => validateCommitSha("--diff-filter=M"),
        (err) => {
          assert.ok(err instanceof Error);
          assert.strictEqual(
            err.message,
            "Invalid commit SHA: --diff-filter=M",
          );
          return true;
        },
      );

      assert.doesNotThrow(() => validateCommitSha("abc1234"));
    },
  );
});
