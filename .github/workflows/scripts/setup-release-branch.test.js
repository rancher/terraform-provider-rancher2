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
  resolveBranchCreationSha,
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

    it("chomps surrounding whitespace cleanly", () => {
      const resPadded = parseReleaseBranch("  release/v16  ");
      assert.strictEqual(resPadded.branch, "release/v16");
      assert.strictEqual(resPadded.versionNumber, 16);

      const resLeading = parseReleaseBranch(" release/v16");
      assert.strictEqual(resLeading.branch, "release/v16");

      const resTrailing = parseReleaseBranch("release/v16 ");
      assert.strictEqual(resTrailing.branch, "release/v16");

      const resNewline = parseReleaseBranch("\n release/v99 \t");
      assert.strictEqual(resNewline.branch, "release/v99");
      assert.strictEqual(resNewline.versionNumber, 99);
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

    it("sets last-release-sha when provided and removes bootstrap-sha", () => {
      const initial = {
        "bootstrap-sha": "abc1234",
        packages: { ".": { "release-type": "go" } },
      };
      const updatedStr = buildReleasePleaseConfig(
        JSON.stringify(initial),
        "v16.0.0",
        "2c006a36c2833ff9f07eb4bdcc370fcd201d15e8",
      );
      const parsed = JSON.parse(updatedStr);
      assert.strictEqual(parsed["release-as"], "v16.0.0");
      assert.strictEqual("bootstrap-sha" in parsed, false);
      assert.strictEqual(
        parsed["last-release-sha"],
        "2c006a36c2833ff9f07eb4bdcc370fcd201d15e8",
      );
    });

    it("throws when existing config string is invalid JSON", () => {
      assert.throws(
        () => buildReleasePleaseConfig("{ invalid json ", "v16.0.0"),
        /Failed to parse existing release-please-config\.json/,
      );
    });

    it("resets invalid packages array to default packages map", () => {
      const initial = {
        packages: ["invalid-array"],
      };
      const updatedStr = buildReleasePleaseConfig(
        JSON.stringify(initial),
        "v16.0.0",
      );
      const parsed = JSON.parse(updatedStr);
      assert.deepStrictEqual(parsed.packages, {
        ".": { "release-type": "go" },
      });
    });

    it("handles primitive JSON types by defaulting to empty object", () => {
      const updatedStr = buildReleasePleaseConfig("null", "v16.0.0");
      const parsed = JSON.parse(updatedStr);
      assert.strictEqual(parsed["release-as"], "v16.0.0");
    });
  });

  describe("resolveBranchCreationSha", () => {
    it("resolves the merge-base SHA between origin/main and the release branch", async () => {
      const mockGit = async (args) => {
        if (
          args[0] === "merge-base" &&
          args[1] === "origin/main" &&
          args[2] === "origin/release/v16"
        ) {
          return {
            stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
            stderr: "",
          };
        }
        return { stdout: "", stderr: "" };
      };

      const sha = await resolveBranchCreationSha(mockGit, "release/v16");
      assert.strictEqual(sha, "96b10c94caf2dad853185325e586f1d9e7ed719e");
    });

    it("falls back to rev-parse on origin/<branch> when merge-base fails", async () => {
      const mockGit = async (args) => {
        if (args[0] === "merge-base") {
          throw new Error("merge-base failed");
        }
        if (args[0] === "rev-parse" && args[1] === "origin/release/v16") {
          return {
            stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
            stderr: "",
          };
        }
        return { stdout: "", stderr: "" };
      };

      const sha = await resolveBranchCreationSha(mockGit, "release/v16");
      assert.strictEqual(sha, "96b10c94caf2dad853185325e586f1d9e7ed719e");
    });

    it("throws when git lookups fail or return invalid SHA", async () => {
      const mockGit = async () => {
        throw new Error("git failure");
      };

      await assert.rejects(
        () => resolveBranchCreationSha(mockGit, "release/v16"),
        /Failed to resolve branch creation SHA for 'release\/v16' against 'origin\/main'/,
      );

      const mockGitInvalid = async () => ({
        stdout: "invalid-sha\n",
        stderr: "",
      });
      await assert.rejects(
        () => resolveBranchCreationSha(mockGitInvalid, "release/v16"),
        /Failed to resolve branch creation SHA for 'release\/v16' against 'origin\/main'/,
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

      const e2eStub = stubs[".github/workflows/e2e.yml"];
      assert.match(e2eStub, /name: "E2E Tests"/);
      assert.match(e2eStub, /branches:\s+- 'release\/v\*'/);
      assert.match(e2eStub, /types: \[opened, synchronize, reopened\]/);
      assert.match(e2eStub, /merge_group:/);
      assert.match(e2eStub, /id-token: write/);
      assert.match(e2eStub, /pull-requests: read/);
      assert.match(
        e2eStub,
        /uses: rancher\/terraform-provider-rancher2\/\.github\/workflows\/e2e\.yml@main/,
      );

      assert.strictEqual(ALLOWED_STUB_FILENAMES.size, 4);
      assert.ok(ALLOWED_STUB_FILENAMES.has("pull_request.yaml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("release.yml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("sync-release-manifest.yml"));
      assert.ok(ALLOWED_STUB_FILENAMES.has("e2e.yml"));
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
        const scriptsDir = path.join(tmpDir, "scripts");
        await fs.mkdir(scriptsDir, { recursive: true });
        await fs.writeFile(
          path.join(scriptsDir, "legacy.js"),
          "console.log(1)",
          "utf8",
        );

        await cleanupNonStubWorkflows(tmpDir);

        const files = await fs.readdir(tmpDir);
        assert.ok(files.includes("pull_request.yaml"));
        assert.ok(files.includes("release.yml"));
        assert.ok(files.includes("not-a-workflow.txt"));
        assert.ok(!files.includes("arbitrary-ci.yml"));
        assert.ok(!files.includes("legacy-release.yaml"));
        assert.ok(!files.includes("scripts"));
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
      assert.match(body, /release\.md/);
      assert.match(body, /both draft and prerelease are enabled/);
    });

    it("formats update description and reflects preserved manifest version", () => {
      const body = formatPrBody({
        branch: "release/v15",
        versionNumber: 15,
        releaseAs: "v15.0.0",
        manifestVersion: "15.2.0",
        isUpdate: true,
      });
      assert.match(body, /### Release Branch Update: `release\/v15`/);
      assert.match(
        body,
        /updates configuration and caller stubs on release branch `release\/v15`/,
      );
      assert.match(body, /root version is properly set to `15\.2\.0`/);
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
          if (args[0] === "diff" && args[1] === "--cached") {
            return { stdout: "release-please-config.json\n", stderr: "" };
          }
          if (args[0] === "merge-base") {
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          }
          if (args[0] === "tag") {
            return { stdout: "v15.3.0\n", stderr: "" };
          }
          if (args[0] === "rev-list" && args[3] === "v15.3.0") {
            return {
              stdout: "2c006a36c2833ff9f07eb4bdcc370fcd201d15e8\n",
              stderr: "",
            };
          }
          return { stdout: "", stderr: "" };
        };

        let createdPrArgs = null;
        const mockPulls = {
          list: async () => ({ data: [] }),
          create: async (args) => {
            createdPrArgs = args;
            return {
              data: {
                number: 42,
                html_url: "https://github.com/rancher/test/pull/42",
              },
            };
          },
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
          process: {
            env: { INPUT_BRANCH: "release/v16", GITHUB_ACTIONS: "true" },
            cwd: () => tmpDir,
          },
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

        // Verify git fetch calls use explicit refspecs to guarantee local remote-tracking branches are updated
        const fetchMainCall = gitCalls.find(
          (call) =>
            call[0] === "fetch" &&
            call[1] === "origin" &&
            call[2] === "+refs/heads/main:refs/remotes/origin/main",
        );
        assert.ok(
          fetchMainCall,
          "Expected git fetch for main to use explicit refspec '+refs/heads/main:refs/remotes/origin/main'",
        );

        const fetchBranchCall = gitCalls.find(
          (call) =>
            call[0] === "fetch" &&
            call[1] === "origin" &&
            call[2] ===
              "+refs/heads/release/v16:refs/remotes/origin/release/v16",
        );
        assert.ok(
          fetchBranchCall,
          "Expected git fetch for release branch to use explicit refspec '+refs/heads/release/v16:refs/remotes/origin/release/v16'",
        );

        // Verify chore branch was checked out and reset --hard to guarantee pristine state
        const checkoutIndex = gitCalls.findIndex(
          (call) =>
            call[0] === "checkout" &&
            call[1] === "-B" &&
            call[2] === "chore/setup-release-v16" &&
            call[3] === "origin/release/v16",
        );
        assert.ok(
          checkoutIndex !== -1,
          "Expected checkout -B chore/setup-release-v16 origin/release/v16",
        );
        const resetCall = gitCalls[checkoutIndex + 1];
        assert.deepStrictEqual(
          resetCall,
          ["reset", "--hard", "origin/release/v16"],
          "Expected reset --hard origin/release/v16 immediately after checkout",
        );

        // Verify files were created in tmpDir
        const configRaw = await fs.readFile(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        );
        const configJson = JSON.parse(configRaw);
        assert.strictEqual(configJson["release-as"], "v16.0.0");
        assert.strictEqual(
          configJson["last-release-sha"],
          "96b10c94caf2dad853185325e586f1d9e7ed719e",
        );

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

        // Verify PR title and commit message used fix prefix
        assert.strictEqual(
          createdPrArgs?.title,
          "fix(release/v16): initialize release branch configuration",
        );
        const commitCall = gitCalls.find(
          (call) =>
            call[0] === "commit" &&
            call[1] === "-m" &&
            call[2] ===
              "fix(release/v16): initialize release branch configuration",
        );
        assert.ok(commitCall, "Expected fix commit message for release setup");

        // Verify rancher2/release.md was created
        const releaseMdRaw = await fs.readFile(
          path.join(tmpDir, "rancher2", "release.md"),
          "utf8",
        );
        assert.strictEqual(releaseMdRaw, "release/v16\n");
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("idempotently updates existing release branch on remote origin when changes are needed", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-existing-update-"),
      );
      try {
        await fs.writeFile(
          path.join(tmpDir, ".release-please-manifest.json"),
          JSON.stringify({ ".": "16.2.0" }, null, 2),
          "utf8",
        );
        await fs.writeFile(
          path.join(tmpDir, "release-please-config.json"),
          JSON.stringify({ "release-as": "v16.2.0" }, null, 2),
          "utf8",
        );

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
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "diff" && args[1] === "--cached") {
            return {
              stdout: ".github/workflows/pull_request.yaml\n",
              stderr: "",
            };
          }
          if (args[0] === "merge-base" || args[0] === "rev-parse") {
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          }
          return { stdout: "", stderr: "" };
        };

        let createdPrArgs = null;
        const mockGithub = {
          rest: {
            pulls: {
              list: async () => ({ data: [] }),
              create: async (args) => {
                createdPrArgs = args;
                return {
                  data: {
                    number: 202,
                    html_url:
                      "https://github.com/rancher/terraform-provider-rancher2/pull/202",
                  },
                };
              },
            },
            issues: {
              addLabels: async () => ({}),
            },
          },
        };

        const result = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: { INPUT_BRANCH: "release/v16", GITHUB_ACTIONS: "true" },
            cwd: () => tmpDir,
          },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(result.prCreated, true);
        assert.strictEqual(result.prNumber, 202);
        assert.strictEqual(
          createdPrArgs?.title,
          "fix(release/v16): update release branch configuration",
        );
        assert.ok(
          createdPrArgs?.body?.includes(
            "root version is properly set to `16.2.0`",
          ),
          "Expected preserved manifest version in PR body",
        );
        assert.ok(
          createdPrArgs?.body?.includes("release-as: v16.2.0"),
          "Expected preserved release-as version in PR body",
        );
        assert.ok(
          !createdPrArgs?.body?.includes("16.0.0"),
          "PR body should not state 16.0.0 when manifest was preserved",
        );

        // Verify release-please-config.json preserved release-as
        const updatedConfigRaw = await fs.readFile(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        );
        const updatedConfig = JSON.parse(updatedConfigRaw);
        assert.strictEqual(updatedConfig["release-as"], "v16.2.0");

        // Verify it did not attempt to push origin/main to release/v16
        const pushedNewBranch = gitCalls.some(
          (call) =>
            call[0] === "push" &&
            call.includes("origin/main:refs/heads/release/v16"),
        );
        assert.strictEqual(
          pushedNewBranch,
          false,
          "Should not push branch creation when branch already exists",
        );

        // Verify chore branch was checked out and reset --hard to guarantee pristine state
        const checkoutIndex = gitCalls.findIndex(
          (call) =>
            call[0] === "checkout" &&
            call[1] === "-B" &&
            call[2] === "chore/setup-release-v16" &&
            call[3] === "origin/release/v16",
        );
        assert.ok(
          checkoutIndex !== -1,
          "Expected checkout -B chore/setup-release-v16 origin/release/v16",
        );
        const resetCall = gitCalls[checkoutIndex + 1];
        assert.deepStrictEqual(
          resetCall,
          ["reset", "--hard", "origin/release/v16"],
          "Expected reset --hard origin/release/v16 immediately after checkout",
        );

        // Verify commit message used update
        const commitCall = gitCalls.find(
          (call) =>
            call[0] === "commit" &&
            call[1] === "-m" &&
            call[2] === "fix(release/v16): update release branch configuration",
        );
        assert.ok(
          commitCall,
          "Expected update commit message for existing branch",
        );
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("idempotently skips PR creation when existing release branch already has target configuration", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-existing-noop-"),
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
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "diff" && args[1] === "--cached") {
            return { stdout: "", stderr: "" };
          }
          if (args[0] === "merge-base" || args[0] === "rev-parse") {
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          }
          return { stdout: "", stderr: "" };
        };

        const mockGithub = {
          rest: {
            pulls: {
              list: async () => {
                throw new Error("Should not be called when no changes");
              },
              create: async () => {
                throw new Error("Should not be called when no changes");
              },
            },
          },
        };

        const result = await runSetupReleaseBranch({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: { INPUT_BRANCH: "release/v16", GITHUB_ACTIONS: "true" },
            cwd: () => tmpDir,
          },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        assert.strictEqual(result.prCreated, false);
        assert.strictEqual(result.prNumber, null);
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }
    });

    it("normalizes valid LAST_RELEASE_SHA and falls back to branch creation lookup when malformed", async () => {
      const tmpDir = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-env-sha-"),
      );
      try {
        const mockGit = async (args) => {
          if (args[0] === "ls-remote") return { stdout: "", stderr: "" };
          if (args[0] === "diff" && args[1] === "--cached")
            return { stdout: "config\n", stderr: "" };
          if (args[0] === "merge-base") {
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          }
          if (args[0] === "tag") return { stdout: "v15.3.0\n", stderr: "" };
          if (args[0] === "rev-list" && args[3] === "v15.3.0") {
            return {
              stdout: "2c006a36c2833ff9f07eb4bdcc370fcd201d15e8\n",
              stderr: "",
            };
          }
          return { stdout: "", stderr: "" };
        };

        const mockGithub = {
          rest: {
            pulls: {
              list: async () => ({ data: [] }),
              create: async () => ({
                data: { number: 10, html_url: "https://pr/10" },
              }),
            },
            issues: { addLabels: async () => ({}) },
          },
        };

        // Test with padded valid SHA
        await runSetupReleaseBranch({
          github: mockGithub,
          context: { repo: { owner: "rancher", repo: "test" } },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: {
              INPUT_BRANCH: "release/v16",
              LAST_RELEASE_SHA: "  abcdefabcdefabcdefabcdefabcdefabcdefabcd  ",
            },
            cwd: () => tmpDir,
          },
          gitExec: mockGit,
          workDir: tmpDir,
        });

        const configRaw = await fs.readFile(
          path.join(tmpDir, "release-please-config.json"),
          "utf8",
        );
        const configJson = JSON.parse(configRaw);
        assert.strictEqual(
          configJson["last-release-sha"],
          "abcdefabcdefabcdefabcdefabcdefabcdefabcd",
        );
      } finally {
        await fs.rm(tmpDir, { recursive: true, force: true });
      }

      const tmpDirFallback = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-env-fallback-"),
      );
      try {
        const mockGit = async (args) => {
          if (args[0] === "ls-remote") return { stdout: "", stderr: "" };
          if (args[0] === "diff" && args[1] === "--cached")
            return { stdout: "config\n", stderr: "" };
          if (args[0] === "merge-base") {
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          }
          if (args[0] === "tag") return { stdout: "v15.3.0\n", stderr: "" };
          if (args[0] === "rev-list" && args[3] === "v15.3.0") {
            return {
              stdout: "2c006a36c2833ff9f07eb4bdcc370fcd201d15e8\n",
              stderr: "",
            };
          }
          return { stdout: "", stderr: "" };
        };

        const mockGithub = {
          rest: {
            pulls: {
              list: async () => ({ data: [] }),
              create: async () => ({
                data: { number: 10, html_url: "https://pr/10" },
              }),
            },
            issues: { addLabels: async () => ({}) },
          },
        };

        // Test with malformed SHA: falls back to branch creation lookup
        await runSetupReleaseBranch({
          github: mockGithub,
          context: { repo: { owner: "rancher", repo: "test" } },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: {
              INPUT_BRANCH: "release/v16",
              LAST_RELEASE_SHA: "invalid-short-sha",
            },
            cwd: () => tmpDirFallback,
          },
          gitExec: mockGit,
          workDir: tmpDirFallback,
        });

        const configRawFallback = await fs.readFile(
          path.join(tmpDirFallback, "release-please-config.json"),
          "utf8",
        );
        const configJsonFallback = JSON.parse(configRawFallback);
        assert.strictEqual(
          configJsonFallback["last-release-sha"],
          "96b10c94caf2dad853185325e586f1d9e7ed719e",
        );
      } finally {
        await fs.rm(tmpDirFallback, { recursive: true, force: true });
      }
    });

    it("guards git identity configuration and remote push operations behind GITHUB_ACTIONS === 'true'", async () => {
      const tmpDirLocal = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-local-config-"),
      );
      const tmpDirCi = await fs.mkdtemp(
        path.join(os.tmpdir(), "setup-release-branch-ci-config-"),
      );
      try {
        const localGitCalls = [];
        const mockLocalGit = async (args) => {
          localGitCalls.push(args);
          if (args[0] === "ls-remote") return { stdout: "", stderr: "" };
          if (args[0] === "diff" && args[1] === "--cached")
            return { stdout: "config\n", stderr: "" };
          if (args[0] === "merge-base")
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          return { stdout: "", stderr: "" };
        };

        let localPullsCreateCalls = 0;
        const mockLocalGithub = {
          rest: {
            pulls: {
              list: async () => ({ data: [] }),
              create: async () => {
                localPullsCreateCalls++;
                return { data: { number: 1, html_url: "https://pr/1" } };
              },
            },
            issues: { addLabels: async () => ({}) },
          },
        };

        // 1. Local environment (GITHUB_ACTIONS is undefined or false)
        const localResult = await runSetupReleaseBranch({
          github: mockLocalGithub,
          context: { repo: { owner: "rancher", repo: "test" } },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: { INPUT_BRANCH: "release/v16", GITHUB_ACTIONS: "false" },
            cwd: () => tmpDirLocal,
          },
          gitExec: mockLocalGit,
          workDir: tmpDirLocal,
        });

        assert.strictEqual(
          localPullsCreateCalls,
          0,
          "Expected no pulls.create calls when GITHUB_ACTIONS !== 'true'",
        );
        assert.strictEqual(
          localResult.prCreated,
          false,
          "Expected prCreated to be false when GITHUB_ACTIONS !== 'true'",
        );
        assert.strictEqual(
          localResult.prNumber,
          null,
          "Expected prNumber to be null when GITHUB_ACTIONS !== 'true'",
        );

        const hasLocalConfig = localGitCalls.some(
          (call) =>
            call[0] === "config" &&
            (call[1] === "user.name" || call[1] === "user.email"),
        );
        assert.strictEqual(
          hasLocalConfig,
          false,
          "Expected no git config user.name/user.email calls when GITHUB_ACTIONS !== 'true'",
        );

        const hasLocalPush = localGitCalls.some((call) => call[0] === "push");
        assert.strictEqual(
          hasLocalPush,
          false,
          "Expected no git push calls when GITHUB_ACTIONS !== 'true'",
        );

        const hasLocalBranchFetch = localGitCalls.some(
          (call) =>
            call[0] === "fetch" &&
            call[1] === "origin" &&
            call[2]?.includes("refs/heads/release/v16"),
        );
        assert.strictEqual(
          hasLocalBranchFetch,
          false,
          "Expected no fetch of non-existent release branch when GITHUB_ACTIONS !== 'true' and branch not on origin",
        );

        const localCheckoutCall = localGitCalls.find(
          (call) =>
            call[0] === "checkout" &&
            call[1] === "-B" &&
            call[2] === "chore/setup-release-v16",
        );
        assert.strictEqual(
          localCheckoutCall[3],
          "origin/main",
          "Expected checkout baseRef to fall back to origin/main when branch was not pushed to origin",
        );

        // 2. CI environment (GITHUB_ACTIONS === "true")
        const ciGitCalls = [];
        const mockCiGit = async (args) => {
          ciGitCalls.push(args);
          if (args[0] === "ls-remote") return { stdout: "", stderr: "" };
          if (args[0] === "diff" && args[1] === "--cached")
            return { stdout: "config\n", stderr: "" };
          if (args[0] === "merge-base")
            return {
              stdout: "96b10c94caf2dad853185325e586f1d9e7ed719e\n",
              stderr: "",
            };
          return { stdout: "", stderr: "" };
        };

        let ciPullsCreateCalls = 0;
        const mockCiGithub = {
          rest: {
            pulls: {
              list: async () => ({ data: [] }),
              create: async () => {
                ciPullsCreateCalls++;
                return { data: { number: 1, html_url: "https://pr/1" } };
              },
            },
            issues: { addLabels: async () => ({}) },
          },
        };

        const ciResult = await runSetupReleaseBranch({
          github: mockCiGithub,
          context: { repo: { owner: "rancher", repo: "test" } },
          core: { info: () => {}, warning: () => {}, setFailed: () => {} },
          process: {
            env: { INPUT_BRANCH: "release/v16", GITHUB_ACTIONS: "true" },
            cwd: () => tmpDirCi,
          },
          gitExec: mockCiGit,
          workDir: tmpDirCi,
        });

        assert.strictEqual(
          ciPullsCreateCalls,
          1,
          "Expected pulls.create to be called when GITHUB_ACTIONS === 'true'",
        );
        assert.strictEqual(
          ciResult.prCreated,
          true,
          "Expected prCreated to be true when GITHUB_ACTIONS === 'true'",
        );
        assert.strictEqual(
          ciResult.prNumber,
          1,
          "Expected prNumber to be 1 when GITHUB_ACTIONS === 'true'",
        );

        const userNameCall = ciGitCalls.find(
          (call) => call[0] === "config" && call[1] === "user.name",
        );
        assert.deepStrictEqual(userNameCall, [
          "config",
          "user.name",
          "github-actions[bot]",
        ]);

        const userEmailCall = ciGitCalls.find(
          (call) => call[0] === "config" && call[1] === "user.email",
        );
        assert.deepStrictEqual(userEmailCall, [
          "config",
          "user.email",
          "41898282+github-actions[bot]@users.noreply.github.com",
        ]);

        const pushBranchCall = ciGitCalls.find(
          (call) =>
            call[0] === "push" &&
            call[1] === "origin" &&
            call[2] === "origin/main:refs/heads/release/v16",
        );
        assert.ok(
          pushBranchCall,
          "Expected push origin origin/main:refs/heads/release/v16 when GITHUB_ACTIONS === 'true'",
        );

        const pushChoreCall = ciGitCalls.find(
          (call) =>
            call[0] === "push" &&
            call[1] === "-u" &&
            call[2] === "origin" &&
            call[3] === "chore/setup-release-v16" &&
            call[4] === "--force",
        );
        assert.ok(
          pushChoreCall,
          "Expected push -u origin chore/setup-release-v16 --force when GITHUB_ACTIONS === 'true'",
        );
      } finally {
        await fs.rm(tmpDirLocal, { recursive: true, force: true });
        await fs.rm(tmpDirCi, { recursive: true, force: true });
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
