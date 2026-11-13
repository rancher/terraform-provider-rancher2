import assert from "node:assert/strict";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { describe, it } from "node:test";
import run, {
  ALLOWED_STUB_FILENAMES,
  buildReleasePleaseConfig,
  buildReleasePleaseManifest,
  buildWorkflowStubs,
  cleanupNonStubWorkflows,
  formatPrBody,
  parseReleaseBranch,
  runGit,
  runSetupReleaseBranch,
} from "./setup-release-branch.js";

describe("setup-release-branch.js", () => {
  describe("parseReleaseBranch", () => {
    it("parses valid canonical release branch names", () => {
      const res16 = parseReleaseBranch("release/v16");
      assert.deepStrictEqual(res16, {
        branch: "release/v16",
        versionNumber: 16,
        versionString: "16",
        releaseAs: "v16.0.0",
        manifestVersion: "16.0.0",
      });

      const res15 = parseReleaseBranch("release/v15");
      assert.deepStrictEqual(res15, {
        branch: "release/v15",
        versionNumber: 15,
        versionString: "15",
        releaseAs: "v15.0.0",
        manifestVersion: "15.0.0",
      });

      const res1 = parseReleaseBranch("release/v1");
      assert.deepStrictEqual(res1, {
        branch: "release/v1",
        versionNumber: 1,
        versionString: "1",
        releaseAs: "v1.0.0",
        manifestVersion: "1.0.0",
      });

      const res100 = parseReleaseBranch("release/v100");
      assert.deepStrictEqual(res100, {
        branch: "release/v100",
        versionNumber: 100,
        versionString: "100",
        releaseAs: "v100.0.0",
        manifestVersion: "100.0.0",
      });
    });

    it("rejects leading-zero versions (e.g. release/v016, release/v01, release/v00, release/v0)", () => {
      assert.throws(
        () => parseReleaseBranch("release/v016"),
        /leading zeros.*are not permitted/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v01"),
        /leading zeros.*are not permitted/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v00"),
        /leading zeros.*are not permitted/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v0"),
        /leading zeros.*are not permitted/,
      );
    });

    it("rejects surrounding whitespace strictly", () => {
      assert.throws(
        () => parseReleaseBranch("  release/v16  "),
        /surrounding whitespace are not permitted/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v16 "),
        /surrounding whitespace are not permitted/,
      );
      assert.throws(
        () => parseReleaseBranch(" release/v16"),
        /surrounding whitespace are not permitted/,
      );
    });

    it("rejects shorthand branch names (v16, 16, release/16)", () => {
      assert.throws(
        () => parseReleaseBranch("v16"),
        /Branch must strictly match 'release\/v<version>'/,
      );
      assert.throws(
        () => parseReleaseBranch("16"),
        /Branch must strictly match 'release\/v<version>'/,
      );
      assert.throws(
        () => parseReleaseBranch("release/16"),
        /Branch must strictly match 'release\/v<version>'/,
      );
    });

    it("rejects non-canonical or arbitrary branch names", () => {
      assert.throws(() => parseReleaseBranch("main"), /Invalid branch format/);
      assert.throws(
        () => parseReleaseBranch("feature/release"),
        /Invalid branch format/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v"),
        /Invalid branch format/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v16-rc1"),
        /Invalid branch format/,
      );
      assert.throws(
        () => parseReleaseBranch("release/v16/test"),
        /Invalid branch format/,
      );
      assert.throws(() => parseReleaseBranch(""), /Invalid branch format/);
      assert.throws(() => parseReleaseBranch(null), /Invalid branch format/);
      assert.throws(
        () => parseReleaseBranch(undefined),
        /Invalid branch format/,
      );
    });
  });

  describe("buildReleasePleaseConfig", () => {
    it("creates standard config from empty/null input", () => {
      const configStr = buildReleasePleaseConfig(null, "v16.0.0");
      const parsed = JSON.parse(configStr);
      assert.strictEqual(parsed["release-as"], "v16.0.0");
      assert.strictEqual(parsed["always-update"], true);
      assert.strictEqual(parsed["commit-search-depth"], 100);
      assert.strictEqual(parsed["prerelease"], true);
      assert.strictEqual(parsed["draft"], true);
      assert.strictEqual(parsed["force-tag-creation"], true);
      assert.strictEqual(parsed["include-v-in-tag"], true);
      assert.deepStrictEqual(parsed["packages"], {
        ".": { "release-type": "go" },
      });
      assert.ok(parsed["$schema"]);
      assert.strictEqual(configStr.endsWith("\n"), true);
    });

    it("updates existing config and strips obsolete sha keys", () => {
      const initial = {
        "release-as": "v15.2.0",
        "bootstrap-sha": "abc1234",
        "last-release-sha": "def5678",
        packages: { ".": { "release-type": "go" } },
      };
      const updatedStr = buildReleasePleaseConfig(
        JSON.stringify(initial),
        "v16.0.0",
      );
      const parsed = JSON.parse(updatedStr);
      assert.strictEqual(parsed["release-as"], "v16.0.0");
      assert.strictEqual("bootstrap-sha" in parsed, false);
      assert.strictEqual("last-release-sha" in parsed, false);
      assert.strictEqual(parsed["draft"], true);
    });

    it("throws when existing config string is invalid JSON", () => {
      assert.throws(
        () => buildReleasePleaseConfig("{ invalid json ", "v16.0.0"),
        /Failed to parse existing release-please-config\.json/,
      );
    });
  });

  describe("buildReleasePleaseManifest", () => {
    it("generates correct manifest JSON with trailing newline", () => {
      const manifestStr = buildReleasePleaseManifest("16.0.0");
      const parsed = JSON.parse(manifestStr);
      assert.deepStrictEqual(parsed, { ".": "16.0.0" });
      assert.strictEqual(manifestStr.endsWith("\n"), true);
    });
  });

  describe("buildWorkflowStubs & ALLOWED_STUB_FILENAMES", () => {
    it("returns expected caller stubs pointing strictly to @main", () => {
      const stubs = buildWorkflowStubs();
      assert.strictEqual(Object.keys(stubs).length, 4);

      const prStub = stubs[".github/workflows/pull_request.yaml"];
      assert.match(
        prStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/pull_request\.yaml@main/,
      );

      const releaseStub = stubs[".github/workflows/release.yml"];
      assert.match(
        releaseStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/release-please\.yml@main/,
      );

      const syncStub = stubs[".github/workflows/sync-release-manifest.yml"];
      assert.match(
        syncStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/sync-release-manifest\.yml@main/,
      );

      const listenerStub = stubs[".github/workflows/e2e-pr-listener.yml"];
      assert.match(listenerStub, /name: PR Executor Listener/);
      assert.match(listenerStub, /branches:\s+- 'release\/v\*'/);

      assert.strictEqual(ALLOWED_STUB_FILENAMES.size, 4);
      assert.ok(ALLOWED_STUB_FILENAMES.has("pull_request.yaml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("release.yml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("sync-release-manifest.yml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("e2e-pr-listener.yml"));
    });
  });

  describe("cleanupNonStubWorkflows", () => {
    it("deletes arbitrary non-stub workflows and preserves allowed stubs", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "cleanup-workflows-test-"),
      );
      try {
        await fs.writeFile(
          path.join(tmpDir, "pull_request.yaml"),
          "stub",
          "utf8",
        );
        await fs.writeFile(path.join(tmpDir, "release.yml"), "stub", "utf8");
        await fs.writeFile(path.join(tmpDir, "arbitrary-ci.yml"), "ci", "utf8");
        await fs.writeFile(
          path.join(tmpDir, "legacy-release.yaml"),
          "legacy",
          "utf8",
        );
        await fs.writeFile(
          path.join(tmpDir, "not-a-workflow.txt"),
          "text",
          "utf8",
        );

        await cleanupNonStubWorkflows(tmpDir);

        const files = await fs.readdir(tmpDir);
        assert.ok(files.includes("pull_request.yaml"));
        assert.ok(files.includes("release.yml"));
        assert.ok(files.includes("not-a-workflow.txt"));
        assert.ok(!files.includes("arbitrary-ci.yml"));
        assert.ok(!files.includes("legacy-release.yaml"));
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("handles non-existent directory safely without error", async () => {
      const nonExistent = path.join(
        os.tmpdir(),
        "non-existent-dir-" + Date.now(),
      );
      await assert.doesNotReject(() => cleanupNonStubWorkflows(nonExistent));
    });
  });

  describe("formatPrBody", () => {
    it("includes branch, release-as, and manifest version details", () => {
      const body = formatPrBody({
        branch: "release/v16",
        versionNumber: 16,
        releaseAs: "v16.0.0",
        manifestVersion: "16.0.0",
      });
      assert.match(body, /release\/v16/);
      assert.match(body, /v16\.0\.0/);
      assert.match(body, /16\.0\.0/);
      assert.match(body, /pull_request\.yaml/);
      assert.match(body, /release\.yml/);
    });
  });

  describe("runGit error sanitization", () => {
    it("sanitizes authorization headers and tokens in error messages", async () => {
      await assert.rejects(
        () =>
          runGit([
            "subcommand-that-does-not-exist",
            "Authorization: basic secrettoken12345",
          ]),
        (err) => {
          assert.ok(!err.message.includes("secrettoken12345"));
          assert.ok(err.message.includes("[REDACTED]"));
          return true;
        },
      );
    });
  });

  describe("runSetupReleaseBranch", () => {
    it("strictly creates release branch from origin/main when not present on remote", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-run-"),
      );
      try {
        const gitCalls = [];
        const mockGit = async (args) => {
          gitCalls.push(args);
          if (args[0] === "ls-remote") {
            // Simulate branch does not exist on remote
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "status") {
            return { stdout: "M release-please-config.json\n", stderr: "" };
          }
          return { stdout: "", stderr: "" };
        };

        const mockPulls = {
          list: async () => ({ data: [] }),
          create: async () => ({
            data: {
              number: 42,
              html_url: "https://github.com/rancher/test/pull/42",
            },
          }),
        };

        const addedLabels = [];
        const mockIssues = {
          addLabels: async (params) => {
            addedLabels.push(params);
            return { data: [] };
          },
        };

        const mockGithub = {
          rest: {
            pulls: mockPulls,
            issues: mockIssues,
          },
        };

        const infoLogs = [];
        const mockCore = {
          info: (msg) => infoLogs.push(msg),
          warning: () => {},
          setFailed: () => {},
        };

        const res = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: mockCore,
          process: { env: { INPUT_BRANCH: "release/v16" }, cwd: () => tmpDir },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(res.branch, "release/v16");
        assert.strictEqual(res.versionNumber, 16);
        assert.strictEqual(res.choreBranch, "chore/setup-release-v16");
        assert.strictEqual(res.prNumber, 42);
        assert.strictEqual(res.prCreated, true);

        // Verify invariant: branch was created strictly from origin/main
        const pushCall = gitCalls.find(
          (call) =>
            call[0] === "push" &&
            call[1] === "origin" &&
            call[2] === "origin/main:refs/heads/release/v16",
        );
        assert.ok(
          pushCall,
          "Expected release branch to be created strictly and exclusively from 'origin/main'",
        );

        // Verify files were created in tmpDir
        const configRaw = await fs.readFile(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        );
        const configJson = JSON.parse(configRaw);
        assert.strictEqual(configJson["release-as"], "v16.0.0");

        const manifestRaw = await fs.readFile(
          path.join(tmpDir, ".release-please-manifest.json"),
          "utf8",
        );
        const manifestJson = JSON.parse(manifestRaw);
        assert.strictEqual(manifestJson["."], "16.0.0");

        const prStub = await fs.readFile(
          path.join(tmpDir, ".github", "workflows", "pull_request.yaml"),
          "utf8",
        );
        assert.match(prStub, /uses:.*@main/);

        // Verify label was added
        assert.strictEqual(addedLabels.length, 1);
        assert.deepStrictEqual(addedLabels[0].labels, [
          "internal/pr-release-chore",
        ]);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("skips origin/main push if release branch already exists on remote and ensures labels on rerun", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-existing-"),
      );
      try {
        const gitCalls = [];
        const mockGit = async (args) => {
          gitCalls.push(args);
          if (args[0] === "ls-remote") {
            if (args[3] === "refs/heads/release/v16") {
              return {
                stdout: "abcd1234refs/heads/release/v16\n",
                stderr: "",
              };
            }
            if (args[3] === "refs/heads/chore/setup-release-v16") {
              return {
                stdout: "ef01234refs/heads/chore/setup-release-v16\n",
                stderr: "",
              };
            }
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "status") {
            return { stdout: "", stderr: "" }; // No changes
          }
          if (args[0] === "rev-list") {
            return { stdout: "1\n", stderr: "" }; // 1 commit ahead
          }
          return { stdout: "", stderr: "" };
        };

        const addedLabels = [];
        const mockGithub = {
          rest: {
            pulls: {
              list: async () => ({
                data: [
                  {
                    number: 99,
                    html_url: "https://github.com/rancher/test/pull/99",
                  },
                ],
              }),
              create: async () => {
                throw new Error("Should not be called when PR already exists");
              },
            },
            issues: {
              addLabels: async (params) => {
                addedLabels.push(params);
                return { data: [] };
              },
            },
          },
        };

        const res = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: { env: { INPUT_BRANCH: "release/v16" }, cwd: () => tmpDir },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(res.prNumber, 99);
        assert.strictEqual(res.prCreated, false);

        // Verify label was ensured on existing PR during rerun
        assert.strictEqual(addedLabels.length, 1);
        assert.strictEqual(addedLabels[0].issue_number, 99);
        assert.deepStrictEqual(addedLabels[0].labels, [
          "internal/pr-release-chore",
        ]);

        // Verify origin/main was NOT pushed to release/v16 since it already existed
        const pushNewBranchCall = gitCalls.find(
          (call) =>
            call[0] === "push" &&
            call[2] === "origin/main:refs/heads/release/v16",
        );
        assert.strictEqual(pushNewBranchCall, undefined);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("handles no-op cleanly without creating PR when branch is already configured and no chore branch exists", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-noop-"),
      );
      try {
        const mockGit = async (args) => {
          if (args[0] === "ls-remote") {
            if (args[3] === "refs/heads/release/v16") {
              return {
                stdout: "abcd1234refs/heads/release/v16\n",
                stderr: "",
              };
            }
            // Chore branch does not exist on remote
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "status") {
            return { stdout: "", stderr: "" }; // No changes
          }
          return { stdout: "", stderr: "" };
        };

        const mockGithub = {
          rest: {
            pulls: {
              list: async () => {
                throw new Error("Should not call pulls.list on no-op");
              },
              create: async () => {
                throw new Error("Should not call pulls.create on no-op");
              },
            },
          },
        };

        const res = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: { env: { INPUT_BRANCH: "release/v16" }, cwd: () => tmpDir },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(res.prNumber, null);
        assert.strictEqual(res.prCreated, false);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("handles no-op cleanly when remote chore branch exists but is not ahead of release branch", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-zero-ahead-"),
      );
      try {
        const mockGit = async (args) => {
          if (args[0] === "ls-remote") {
            if (args[3] === "refs/heads/release/v16") {
              return {
                stdout: "abcd1234refs/heads/release/v16\n",
                stderr: "",
              };
            }
            if (args[3] === "refs/heads/chore/setup-release-v16") {
              return {
                stdout: "abcd1234refs/heads/chore/setup-release-v16\n",
                stderr: "",
              };
            }
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "status") {
            return { stdout: "", stderr: "" }; // No changes
          }
          if (args[0] === "rev-list") {
            return { stdout: "0\n", stderr: "" }; // 0 commits ahead
          }
          return { stdout: "", stderr: "" };
        };

        const mockGithub = {
          rest: {
            pulls: {
              list: async () => {
                throw new Error(
                  "Should not call pulls.list when 0 commits ahead",
                );
              },
              create: async () => {
                throw new Error(
                  "Should not call pulls.create when 0 commits ahead",
                );
              },
            },
          },
        };

        const res = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: { env: { INPUT_BRANCH: "release/v16" }, cwd: () => tmpDir },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(res.prNumber, null);
        assert.strictEqual(res.prCreated, false);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("rejects invalid branch names and calls core.setFailed in default export", async () => {
      let failedMessage = null;
      const mockCore = {
        info: () => {},
        warning: () => {},
        setFailed: (msg) => {
          failedMessage = msg;
        },
      };

      await assert.rejects(
        () =>
          run({
            core: mockCore,
            process: { env: { INPUT_BRANCH: "v16" } },
          }),
        /Branch must strictly match 'release\/v<version>'/,
      );

      assert.match(
        failedMessage,
        /Branch must strictly match 'release\/v<version>'/,
      );
    });
  });
});
