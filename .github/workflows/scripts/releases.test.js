import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import defaultHandler, {
  computeNextRcTag,
  runCheckMaintainer,
  runRcNotify,
  runTrackingIssue,
  runPrepareRc,
  runVerifyTag,
  validateReleaseTag,
  runDispatchTag,
  runPrepareRegistryManifest,
} from "./releases.js";

test("runCheckMaintainer - authorizes maintainer from MAINTAINERS env", async () => {
  let loggedInfo = null;
  const mockCore = {
    info: (msg) => {
      loggedInfo = msg;
    },
  };
  const mockContext = { actor: "alice" };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(["alice", "bob"]),
    },
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, true);
  assert.match(loggedInfo, /alice.*authorized.*true/);
});

test("runCheckMaintainer - does not authorize bot accounts unless explicitly allowed", async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: "rancher-eio[bot]" };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(["matt"]),
    },
  };

  let isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, false);

  mockProcess.env.MAINTAINERS = JSON.stringify(["matt", "rancher-eio[bot]"]);
  isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, true);
});

test("runCheckMaintainer - rejects unauthorized actors", async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: "unauthorized-user" };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(["matt"]),
    },
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, false);
});

test("validateReleaseTag - accepts valid stable release tag", () => {
  const result = validateReleaseTag({ tag: "v1.2.3" });
  assert.deepStrictEqual(result, { valid: true, tag: "v1.2.3", isRc: false });
});

test("validateReleaseTag - accepts valid RC tag with canonical release branch", () => {
  const result = validateReleaseTag({
    tag: "v16.0.0-rc.1",
    branch: "release/v16",
  });
  assert.deepStrictEqual(result, {
    valid: true,
    tag: "v16.0.0-rc.1",
    isRc: true,
  });
});

test("validateReleaseTag - normalizes branch prefix on RC releases", () => {
  const result = validateReleaseTag({
    tag: "v15.2.0-rc.3",
    branch: "refs/heads/origin/release/v15/",
  });
  assert.deepStrictEqual(result, {
    valid: true,
    tag: "v15.2.0-rc.3",
    isRc: true,
  });
});

test("validateReleaseTag - throws when tag does not start with v", () => {
  assert.throws(() => validateReleaseTag({ tag: "1.2.3" }), {
    message: /Tag must start with 'v'/,
  });
});

test("validateReleaseTag - throws when RC tag lacks branch input", () => {
  assert.throws(() => validateReleaseTag({ tag: "v16.0.0-rc.1" }), {
    message: /'branch' input is required for RC releases/,
  });
});

test("validateReleaseTag - throws when RC branch is not a canonical release branch", () => {
  assert.throws(
    () => validateReleaseTag({ tag: "v16.0.0-rc.1", branch: "main" }),
    { message: /RC releases must target a 'release\/v<major>' branch/ },
  );
  assert.throws(
    () => validateReleaseTag({ tag: "v16.0.0-rc.1", branch: "feature-branch" }),
    { message: /RC releases must target a 'release\/v<major>' branch/ },
  );
});

test("runVerifyTag - logs info and succeeds for valid tag", async () => {
  let loggedInfo = null;
  const mockCore = {
    info: (msg) => {
      loggedInfo = msg;
    },
    setFailed: () => {},
  };

  const result = await runVerifyTag({
    core: mockCore,
    process: { env: { TAG: "v1.2.3" } },
  });

  assert.strictEqual(result, true);
  assert.match(loggedInfo, /Tag 'v1.2.3' is valid for release/);
});

test("runVerifyTag - calls core.setFailed and throws on invalid tag", async () => {
  let failedMessage = null;
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    },
  };

  await assert.rejects(
    async () => {
      await runVerifyTag({
        core: mockCore,
        process: { env: { TAG: "invalid-tag" } },
      });
    },
    { message: /Tag must start with 'v'/ },
  );

  assert.match(failedMessage, /Tag must start with 'v'/);
});

test("computeNextRcTag - increments existing RC tag", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.2" },
      { name: "v15.2.0-rc.1" },
      { name: "v15.2.0-rc.2" },
      { name: "v15.2.0-rc.3" },
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.4");
});

test("computeNextRcTag - starts new RC cycle as rc.1 when latest is full release", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.1.1" },
      { name: "v15.1.2" },
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.1.3-rc.1");
});

test("computeNextRcTag - uses version from release-please-config.json release-as if present", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.1.1" },
      { name: "v15.1.2" },
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === "release-please-config.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ "release-as": "v15.3.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.3.0-rc.1");
});

test("computeNextRcTag - uses version from .release-please-manifest.json if higher than latest full release", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.1.1" },
      { name: "v15.1.2" },
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ ".": "15.2.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.1");
});

test("computeNextRcTag - ignores manifest with mismatched major version", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ ".": "14.0.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.1.1-rc.1");
});

test("computeNextRcTag - advances active RC when manifest version equals latest full release", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.2" }, { name: "v15.1.3-rc.1" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ ".": "15.1.2" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.1.3-rc.2");
});

test("computeNextRcTag - normalizes branch ref prefix", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "refs/heads/release/v15",
  });

  assert.strictEqual(tag, "v15.1.1-rc.1");
});

test("computeNextRcTag - brand new release branch with no tags", async () => {
  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v16",
  });

  assert.strictEqual(tag, "v16.0.0-rc.1");
});

test("computeNextRcTag - rejects invalid branch name", async () => {
  const mockGithub = {
    paginate: async () => [],
  };

  await assert.rejects(
    async () => {
      await computeNextRcTag({
        github: mockGithub,
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        branch: "main",
      });
    },
    { message: /does not match expected pattern/ },
  );
});

test("runPrepareRc - skips RC release and outputs should_release=false when commit is branch initialization", async () => {
  const outputs = {};
  let loggedInfo = null;
  const mockCore = {
    info: (msg) => {
      loggedInfo = msg;
    },
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const result = await runPrepareRc({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      payload: {
        head_commit: {
          message: "fix(release/v16): initialize release branch configuration",
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, { should_release: false, tag: null });
  assert.strictEqual(outputs.should_release, "false");
  assert.match(
    loggedInfo,
    /release branch initialization or maintenance; skipping RC release/,
  );
});

test("runPrepareRc - skips RC release when commit message is branch initialization with PR number", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const result = await runPrepareRc({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      payload: {
        head_commit: {
          message:
            "fix(release/v16): initialize release branch configuration (#2710)",
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, { should_release: false, tag: null });
  assert.strictEqual(outputs.should_release, "false");
});

test("runPrepareRc - skips RC release when commit message is manifest sync or release chore", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const resultManifest = await runPrepareRc({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      payload: {
        head_commit: {
          message:
            "chore(release/v16): align release-please manifest with tag v16.0.0",
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(resultManifest, { should_release: false, tag: null });
  assert.strictEqual(outputs.should_release, "false");

  const resultReleaseChore = await runPrepareRc({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v15",
      payload: {
        head_commit: {
          message: "chore(release/v15): release 15.2.0 (#2705)",
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(resultReleaseChore, {
    should_release: false,
    tag: null,
  });
  assert.strictEqual(outputs.should_release, "false");
});

test("runPrepareRc - computes next RC tag and outputs should_release=true for regular backport commit", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const mockGithub = {
    paginate: async () => [{ name: "v16.0.0-rc.1" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const gitCalls = [];
  const mockRunGit = async (args) => {
    gitCalls.push(args);
    if (args[0] === "ls-remote") {
      return { stdout: "", stderr: "" };
    }
    return { stdout: "", stderr: "" };
  };

  const result = await runPrepareRc({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      payload: {
        head_commit: {
          message: "fix: update provider logic (#2709)",
        },
      },
    },
    core: mockCore,
    process: { env: { SHA: "explicit-backport-sha" } },
    runGit: mockRunGit,
  });

  assert.deepStrictEqual(result, {
    should_release: true,
    tag: "v16.0.0-rc.2",
  });
  assert.strictEqual(outputs.should_release, "true");
  assert.strictEqual(outputs.tag, "v16.0.0-rc.2");
  assert.ok(
    gitCalls.some(
      (c) =>
        c[0] === "tag" &&
        c[1] === "-a" &&
        c[2] === "v16.0.0-rc.2" &&
        c[3] === "-m" &&
        c[4] === "Release Candidate v16.0.0-rc.2" &&
        c[5] === "explicit-backport-sha",
    ),
  );
  assert.ok(
    gitCalls.some(
      (c) => c[0] === "push" && c[1] === "origin" && c[2] === "v16.0.0-rc.2",
    ),
  );
});

test("runPrepareRc - skips creating and pushing tag if tag already exists on remote origin with matching SHA", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const mockGithub = {
    paginate: async () => [{ name: "v16.0.0-rc.1" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const gitCalls = [];
  const mockRunGit = async (args) => {
    gitCalls.push(args);
    if (args[0] === "ls-remote") {
      return {
        stdout: "abcd1234\trefs/tags/v16.0.0-rc.2\n",
        stderr: "",
      };
    }
    return { stdout: "", stderr: "" };
  };

  const result = await runPrepareRc({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      payload: {
        head_commit: {
          message: "fix: update provider logic (#2709)",
        },
      },
    },
    core: mockCore,
    process: { env: { SHA: "abcd1234" } },
    runGit: mockRunGit,
  });

  assert.deepStrictEqual(result, {
    should_release: true,
    tag: "v16.0.0-rc.2",
  });
  assert.strictEqual(outputs.tag, "v16.0.0-rc.2");
  // Verify tag -a and push were NOT called
  assert.strictEqual(
    gitCalls.some((c) => c[0] === "tag"),
    false,
  );
  assert.strictEqual(
    gitCalls.some((c) => c[0] === "push"),
    false,
  );
});

test("runPrepareRc - throws error when remote tag points to different SHA than target SHA", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v16.0.0-rc.1" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const mockRunGit = async (args) => {
    if (args[0] === "ls-remote") {
      return {
        stdout: "different-remote-sha\trefs/tags/v16.0.0-rc.2\n",
        stderr: "",
      };
    }
    return { stdout: "", stderr: "" };
  };

  await assert.rejects(
    async () => {
      await runPrepareRc({
        github: mockGithub,
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          ref: "refs/heads/release/v16",
          payload: {
            head_commit: {
              message: "fix: update provider logic (#2709)",
            },
          },
        },
        core: { info: () => {}, warning: () => {}, setOutput: () => {} },
        process: { env: { SHA: "my-target-sha" } },
        runGit: mockRunGit,
      });
    },
    {
      message:
        /Tag 'v16.0.0-rc.2' already exists on remote pointing to different-remote-sha, which does not match requested target SHA my-target-sha/,
    },
  );
});

test("runPrepareRc - resolves commit message via getCommit when not in env or context payload", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        getCommit: async ({ ref }) => {
          assert.strictEqual(ref, "commit-sha-123");
          return {
            data: {
              commit: {
                message:
                  "fix(release/v16): initialize release branch configuration",
              },
            },
          };
        },
      },
    },
  };

  const result = await runPrepareRc({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v16",
      sha: "commit-sha-123",
      payload: {},
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, { should_release: false, tag: null });
  assert.strictEqual(outputs.should_release, "false");
});

test("runPrepareRc - throws when target branch is missing", async () => {
  await assert.rejects(
    async () => {
      await runPrepareRc({
        github: {},
        context: {},
        core: {},
        process: { env: {} },
      });
    },
    { message: /Target branch must be provided/ },
  );
});

test("runPrepareRc - normalizes branch name with refs/heads/ or origin/ prefix and trailing slash", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const result = await runPrepareRc({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/origin/release/v15/",
      payload: {
        head_commit: {
          message: "fix: backport fix for token (#2700)",
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    should_release: true,
    tag: "v15.1.1-rc.1",
  });
  assert.strictEqual(outputs.tag, "v15.1.1-rc.1");
  assert.strictEqual(outputs.should_release, "true");
});

test("runPrepareRc - skips RC release when COMMIT_MESSAGE is passed via env", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const result = await runPrepareRc({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v15",
    },
    core: mockCore,
    process: {
      env: {
        COMMIT_MESSAGE:
          "fix(release/v15): initialize release branch configuration",
      },
    },
  });

  assert.deepStrictEqual(result, { should_release: false, tag: null });
  assert.strictEqual(outputs.should_release, "false");
});

test("runPrepareRc - handles getCommit returning commit without message gracefully", async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getCommit: async () => ({
          data: { commit: {} },
        }),
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const result = await runPrepareRc({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      ref: "refs/heads/release/v15",
      sha: "sha-empty-msg",
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    should_release: true,
    tag: "v15.1.1-rc.1",
  });
  assert.strictEqual(outputs.should_release, "true");
  assert.strictEqual(outputs.tag, "v15.1.1-rc.1");
});

test("computeNextRcTag - falls back to direct listTags when github.paginate is undefined", async () => {
  const mockGithub = {
    rest: {
      repos: {
        listTags: async () => ({
          data: [{ name: "v15.1.0" }, { name: "v15.2.0-rc.1" }],
        }),
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.2");
});

test("computeNextRcTag - uses version from release-please-config.json packages release-as", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === "release-please-config.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({
                    packages: { ".": { "release-as": "15.4.0" } },
                  }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.4.0-rc.1");
});

test("computeNextRcTag - handles trailing slash on branch name", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "refs/heads/release/v15/",
  });

  assert.strictEqual(tag, "v15.1.1-rc.1");
});

test("runRcNotify - skips notification if tag is not an RC", async () => {
  const loggedInfo = [];
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {},
  };
  await runRcNotify({
    github: {},
    context: {},
    core: mockCore,
    process: { env: { TAG: "v15.2.0", BRANCH: "release/v15" } },
  });
  assert.ok(loggedInfo.some((m) => m.includes("does not appear to be an RC")));
});

test("runRcNotify - fails when tag or branch is missing", async () => {
  let failedMessage = null;
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    },
  };
  await runRcNotify({
    github: {},
    context: {},
    core: mockCore,
    process: { env: {} },
  });
  assert.match(failedMessage, /tagName and branchLabel must be provided/);
});

test("runRcNotify - throws when branch format is invalid", async () => {
  const mockCore = {
    info: () => {},
    setFailed: () => {},
  };
  await assert.rejects(
    async () => {
      await runRcNotify({
        github: {},
        context: {},
        core: mockCore,
        process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "main" } },
      });
    },
    { message: /Target branch label "main" is invalid/ },
  );
});

test("runRcNotify - handles empty matching issues gracefully", async () => {
  const loggedInfo = [];
  const mockGithub = {
    paginate: async () => [],
    rest: {
      search: {
        issuesAndPullRequests: () => {},
      },
    },
  };
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {},
  };
  await runRcNotify({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "release/v15" } },
  });
  assert.ok(loggedInfo.some((m) => m.includes("No matching issues found")));
});

test("runRcNotify - comments on matching issues with release url", async () => {
  const loggedInfo = [];
  const createdComments = [];
  const mockGithub = {
    paginate: async (endpoint, params) => {
      assert.match(params.q, /label:"release\/v15"/);
      assert.match(params.q, /label:"internal\/backport"/);
      assert.match(params.q, /label:"internal\/merged"/);
      return [{ number: 101 }, { number: 102 }];
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {},
      },
      issues: {
        createComment: async (params) => {
          createdComments.push(params);
        },
      },
    },
  };
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {},
  };
  await runRcNotify({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "release/v15" } },
  });
  assert.strictEqual(createdComments.length, 2);
  assert.strictEqual(createdComments[0].issue_number, 101);
  assert.match(createdComments[0].body, /v15\.2\.0-rc\.1/);
  assert.strictEqual(createdComments[1].issue_number, 102);
  assert.ok(loggedInfo.some((m) => m.includes("Notified 2 issues")));
});

test("runRcNotify - handles error on comment creation", async () => {
  let failedMessage = null;
  const mockGithub = {
    paginate: async () => [{ number: 999 }],
    rest: {
      search: { issuesAndPullRequests: () => {} },
      issues: {
        createComment: async () => {
          throw new Error("API Rate Limit");
        },
      },
    },
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    },
  };
  await runRcNotify({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "release/v15" } },
  });
  assert.match(
    failedMessage,
    /Failed to comment on issue #999: API Rate Limit/,
  );
});

test("runRcNotify - strips trailing slash from branch name and succeeds", async () => {
  const createdComments = [];
  const mockGithub = {
    paginate: async (endpoint, params) => {
      assert.match(params.q, /label:"release\/v15"/);
      return [{ number: 201 }];
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {},
      },
      issues: {
        createComment: async (params) => {
          createdComments.push(params);
        },
      },
    },
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      assert.fail(`Unexpected failure: ${msg}`);
    },
  };
  await runRcNotify({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "release/v15/" } },
  });
  assert.strictEqual(createdComments.length, 1);
  assert.strictEqual(createdComments[0].issue_number, 201);
});

test("computeNextRcTag - increments RC when existing tags use un-dotted format like -rc1, -rc2", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.2.0-rc1" },
      { name: "v15.2.0-rc2" },
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.3");
});

test("computeNextRcTag - correctly increments double-digit RC tags numerically", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.2.0-rc.1" },
      { name: "v15.2.0-rc.2" },
      { name: "v15.2.0-rc.9" },
      { name: "v15.2.0-rc.10" },
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.11");
});

test("computeNextRcTag - selects highest active RC regardless of tag ordering", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.4-rc.1" },
      { name: "v15.1.3" },
      { name: "v15.2.0-rc.2" },
      { name: "v15.2.0-rc.1" },
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.3");
});

test("runCheckMaintainer - does not allow substring match when MAINTAINERS is JSON string or array", async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: "super" };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify("superuser"),
    },
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, false);
});

test("computeNextRcTag - falls back to first object value in manifest if dot key is absent", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ "terraform-provider-rancher2": "15.5.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.5.0-rc.1");
});

test("runCheckMaintainer - authorizes maintainers case-insensitively", async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: "Alice" };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(["alice", "BOB"]),
    },
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(isAuth, true);
});

test("computeNextRcTag - uses version from release-please-config.json packages with arbitrary package key", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === "release-please-config.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({
                    packages: {
                      "terraform-provider-rancher2": { "release-as": "15.6.0" },
                    },
                  }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.6.0-rc.1");
});

test("runRcNotify - fails gracefully when search API throws an error", async () => {
  let failedMessage = null;
  const mockGithub = {
    paginate: async () => {
      throw new Error("GitHub Search API rate limit exceeded");
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {},
      },
    },
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    },
  };
  await runRcNotify({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.1", BRANCH: "release/v15" } },
  });
  assert.match(
    failedMessage,
    /Failed to search issues for RC notification: GitHub Search API rate limit exceeded/,
  );
});

test("computeNextRcTag - falls back to manifestVersion if configVersion has mismatched major version", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === "release-please-config.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ "release-as": "v14.9.0" }),
                ).toString("base64"),
              },
            };
          }
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ ".": "15.3.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.3.0-rc.1");
});

test("computeNextRcTag - does not regress version when active RC cycle exceeds manifest targetBase", async () => {
  const mockGithub = {
    paginate: async () => [
      { name: "v15.1.0" },
      { name: "v15.3.0-rc.1" },
      { name: "v15.3.0-rc.2" },
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === ".release-please-manifest.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ ".": "15.2.0" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.3.0-rc.3");
});

test("computeNextRcTag - parses release-as with pre-release suffix cleanly", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === "release-please-config.json") {
            return {
              data: {
                content: Buffer.from(
                  JSON.stringify({ "release-as": "v15.4.0-rc.1" }),
                ).toString("base64"),
              },
            };
          }
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.4.0-rc.1");
});

test("runCheckMaintainer - rejects empty or missing actor", async () => {
  const mockCore = { info: () => {} };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(["alice", ""]),
    },
  };

  const isAuthEmpty = await runCheckMaintainer({
    context: { actor: "" },
    core: mockCore,
    process: mockProcess,
  });
  assert.strictEqual(isAuthEmpty, false);

  const isAuthMissing = await runCheckMaintainer({
    context: {},
    core: mockCore,
    process: mockProcess,
  });
  assert.strictEqual(isAuthMissing, false);
});

test("computeNextRcTag - correctly parses case-insensitive uppercase RC tags", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }, { name: "v15.2.0-RC.1" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    branch: "release/v15",
  });

  assert.strictEqual(tag, "v15.2.0-rc.2");
});

test("runCheckMaintainer - safely handles null JSON string in MAINTAINERS", async () => {
  const mockCore = { info: () => {} };
  const isAuth = await runCheckMaintainer({
    context: { actor: "null" },
    core: mockCore,
    process: { env: { MAINTAINERS: "null" } },
  });

  assert.strictEqual(isAuth, false);
});

test("default export - routes check-maintainer mode", async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: "alice" };
  const mockProcess = {
    env: {
      SCRIPT_MODE: "check-maintainer",
      MAINTAINERS: JSON.stringify(["alice"]),
    },
  };

  const result = await defaultHandler({
    github: {},
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(result, true);
});

test("default export - routes verify-tag mode", async () => {
  const mockCore = {
    info: () => {},
    setFailed: () => {},
  };

  const result = await defaultHandler({
    core: mockCore,
    process: {
      env: { SCRIPT_MODE: "verify-tag", TAG: "v1.2.3" },
    },
  });

  assert.strictEqual(result, true);
});

test("default export - routes prepare-rc mode", async () => {
  const mockGithub = {
    paginate: async () => [{ name: "v15.1.0" }],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error("Not found");
        },
      },
    },
  };
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: () => {},
  };
  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    ref: "refs/heads/release/v15",
    payload: {
      head_commit: {
        message: "fix: backport bugfix (#2701)",
      },
    },
  };

  const result = await defaultHandler({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: {
      env: { SCRIPT_MODE: "prepare-rc" },
    },
  });

  assert.deepStrictEqual(result, {
    should_release: true,
    tag: "v15.1.1-rc.1",
  });
});

test("default export - routes rc-notify mode", async () => {
  const loggedInfo = [];
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {},
  };

  await defaultHandler({
    github: {},
    context: {},
    core: mockCore,
    process: {
      env: { SCRIPT_MODE: "rc-notify", TAG: "v15.2.0", BRANCH: "release/v15" },
    },
  });

  assert.ok(loggedInfo.some((m) => m.includes("does not appear to be an RC")));
});

test("runDispatchTag - resolves target SHA from explicit SHA and creates/pushes tag", async () => {
  const gitCalls = [];
  const mockRunGit = async (args) => {
    gitCalls.push(args);
    if (args[0] === "rev-parse") {
      return { stdout: "target-sha-1234\n", stderr: "" };
    }
    if (args[0] === "ls-remote") {
      return { stdout: "", stderr: "" };
    }
    return { stdout: "", stderr: "" };
  };

  const outputs = {};
  const mockCore = {
    info: () => {},
    setOutput: (k, v) => {
      outputs[k] = v;
    },
  };

  const result = await runDispatchTag({
    context: { actor: "testactor" },
    core: mockCore,
    process: {
      env: { TAG: "v1.2.3", SHA: "explicit-sha" },
      cwd: () => "/workspace",
    },
    runGit: mockRunGit,
  });

  assert.deepStrictEqual(result, {
    tag: "v1.2.3",
    targetSha: "target-sha-1234",
  });
  assert.strictEqual(outputs.target_sha, "target-sha-1234");
  assert.ok(
    gitCalls.some(
      (c) =>
        c[0] === "tag" &&
        c[1] === "v1.2.3" &&
        c[2] === "-m" &&
        c[3] === "Release v1.2.3" &&
        c[4] === "target-sha-1234",
    ),
  );
  assert.ok(
    gitCalls.some(
      (c) => c[0] === "push" && c[1] === "origin" && c[2] === "v1.2.3",
    ),
  );
});

test("runDispatchTag - checks remote tag existence and succeeds if remote tag matches target SHA", async () => {
  const gitCalls = [];
  let loggedInfo = null;
  const mockRunGit = async (args) => {
    gitCalls.push(args);
    if (args[0] === "rev-parse") {
      return { stdout: "target-sha-1234\n", stderr: "" };
    }
    if (args[0] === "ls-remote") {
      return { stdout: "target-sha-1234\trefs/tags/v1.2.3^{}\n", stderr: "" };
    }
    return { stdout: "", stderr: "" };
  };

  const mockCore = {
    info: (msg) => {
      loggedInfo = msg;
    },
    setOutput: () => {},
  };

  const result = await runDispatchTag({
    context: {},
    core: mockCore,
    process: {
      env: { TAG: "v1.2.3", SHA: "explicit-sha" },
      cwd: () => "/workspace",
    },
    runGit: mockRunGit,
  });

  assert.deepStrictEqual(result, {
    tag: "v1.2.3",
    targetSha: "target-sha-1234",
  });
  // Verify tag creation and push were NOT called
  assert.strictEqual(
    gitCalls.some((c) => c[0] === "tag"),
    false,
  );
  assert.strictEqual(
    gitCalls.some((c) => c[0] === "push"),
    false,
  );
  assert.match(
    loggedInfo,
    /already exists on remote origin pointing to matching target SHA/,
  );
});

test("runDispatchTag - throws error if remote tag exists but points to different SHA", async () => {
  const mockRunGit = async (args) => {
    if (args[0] === "rev-parse") {
      return { stdout: "target-sha-1234\n", stderr: "" };
    }
    if (args[0] === "ls-remote") {
      return {
        stdout: "different-sha-9999\trefs/tags/v1.2.3^{}\n",
        stderr: "",
      };
    }
    return { stdout: "", stderr: "" };
  };

  await assert.rejects(
    async () => {
      await runDispatchTag({
        context: {},
        core: { info: () => {} },
        process: {
          env: { TAG: "v1.2.3", SHA: "explicit-sha" },
          cwd: () => "/workspace",
        },
        runGit: mockRunGit,
      });
    },
    { message: /which does not match requested target SHA target-sha-1234/ },
  );
});

test("runDispatchTag - throws error when TAG is missing", async () => {
  await assert.rejects(
    async () => {
      await runDispatchTag({
        context: {},
        core: {},
        process: { env: {} },
      });
    },
    { message: /Tag input is required for dispatch-tag/ },
  );
});

test("runPrepareRegistryManifest - creates default terraform-registry-manifest.json if missing", async () => {
  const writtenFiles = {};
  const mockFs = {
    existsSync: () => false,
    promises: {
      writeFile: async (p, content) => {
        writtenFiles[p] = content;
      },
    },
  };

  const gitCalls = [];
  const mockRunGit = async (args) => {
    gitCalls.push(args);
    if (args[0] === "tag") {
      return { stdout: "v1.2.3\nv1.2.2\nv1.2.1\n" };
    }
    return { stdout: "" };
  };

  const mockCore = {
    info: () => {},
    warning: () => {},
  };

  const result = await runPrepareRegistryManifest({
    core: mockCore,
    process: {
      env: { TAG: "v1.2.3", GITHUB_WORKSPACE: "/workspace" },
      cwd: () => "/workspace",
    },
    runGit: mockRunGit,
    fsModule: mockFs,
  });

  assert.strictEqual(result.tag, "v1.2.3");
  const manifestKey = Object.keys(writtenFiles).find((k) =>
    k.endsWith("terraform-registry-manifest.json"),
  );
  assert.ok(manifestKey, "Expected manifest to be written");
  const parsed = JSON.parse(writtenFiles[manifestKey]);
  assert.strictEqual(parsed.version, 1);
  assert.deepStrictEqual(parsed.metadata.protocol_versions, [
    "4.0",
    "5.0",
    "6.0",
  ]);
  // Verify extraneous tags deleted
  assert.ok(
    gitCalls.some(
      (c) =>
        c[0] === "tag" &&
        c[1] === "-d" &&
        c.includes("v1.2.2") &&
        c.includes("v1.2.1") &&
        !c.includes("v1.2.3"),
    ),
  );
});

test("default export - routes dispatch-tag and prepare-manifest modes", async () => {
  const mockCore = {
    info: () => {},
    setOutput: () => {},
  };

  const mockRunGit = async (args) => {
    if (args[0] === "rev-parse") {
      return { stdout: "target-sha-1234\n", stderr: "" };
    }
    return { stdout: "", stderr: "" };
  };

  const resultDispatch = await defaultHandler({
    context: {},
    core: mockCore,
    process: {
      env: { SCRIPT_MODE: "dispatch-tag", TAG: "v1.2.3", SHA: "some-sha" },
      cwd: () => "/workspace",
    },
    runGit: mockRunGit,
  });
  assert.strictEqual(resultDispatch.tag, "v1.2.3");

  const mockFs = {
    existsSync: () => true,
    promises: { writeFile: async () => {} },
  };
  const resultManifest = await defaultHandler({
    core: mockCore,
    process: {
      env: {
        SCRIPT_MODE: "prepare-manifest",
        TAG: "v1.2.3",
        GITHUB_WORKSPACE: "/workspace",
      },
      cwd: () => "/workspace",
    },
    fsModule: mockFs,
    runGit: async () => ({ stdout: "v1.2.3\n" }),
  });
  assert.strictEqual(resultManifest.tag, "v1.2.3");
});

test("default export - throws when SCRIPT_MODE is unknown or undefined", async () => {
  await assert.rejects(
    async () => {
      await defaultHandler({
        process: { env: { SCRIPT_MODE: "unsupported-mode" } },
      });
    },
    { message: /Unknown release script mode: unsupported-mode/ },
  );

  await assert.rejects(
    async () => {
      await defaultHandler({
        process: { env: {} },
      });
    },
    { message: /Unknown release script mode: undefined/ },
  );
});

function createMockCore() {
  return {
    infoMessages: [],
    warningMessages: [],
    failedMessage: null,
    info(msg) {
      this.infoMessages.push(msg);
    },
    warning(msg) {
      this.warningMessages.push(msg);
    },
    setFailed(msg) {
      this.failedMessage = msg;
    },
  };
}

test("runTrackingIssue - successfully converts referenced issue into tracking issue", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [
          { name: "main" },
          { name: "release/v14" },
          { name: "release/v15" },
          { name: "release/v9" },
        ];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 42,
            title: "Fix cluster driver bug",
            html_url:
              "https://github.com/rancher/terraform-provider-rancher2/pull/42",
            body: "Some text\n- Addresses: #100\nMore text",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
      },
      search: {
        issuesAndPullRequests: async () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 100);
          return {
            data: {
              number: 100,
              title: "Cluster driver issue",
              labels: [],
              state: "open",
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          assert.strictEqual(issue_number, 100);
          labelsAdded.push(labels);
        },
        createComment: async ({ issue_number, body }) => {
          assert.strictEqual(issue_number, 100);
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  // First step: internal/tracking and internal/user
  assert.deepStrictEqual(labelsAdded[0], [
    "internal/tracking",
    "internal/user",
  ]);
  // Second step: tracking comments (PR linking + release instructions)
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #42 -->/);
  assert.match(commentsCreated[0], /This is the tracking issue for PR #42/);
  assert.match(
    commentsCreated[1],
    /A label has been added for the latest release branch/,
  );
  // Third step: latest release branch label added last
  assert.deepStrictEqual(labelsAdded[1], ["release/v15"]);
});

test("runTrackingIssue - respects existing internal/user label", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 43,
            title: "Add feature",
            body: "- Addresses: #101",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 101,
            title: "User issue",
            labels: [{ name: "internal/user" }],
            state: "open",
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  // internal/user was already present, so only internal/tracking should be added
  assert.deepStrictEqual(labelsAdded[0], ["internal/tracking"]);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #43 -->/);
  assert.match(
    commentsCreated[1],
    /A label has been added for the latest release branch/,
  );
  assert.deepStrictEqual(labelsAdded[1], ["release/v15"]);
});

test("runTrackingIssue - skips issue already marked internal/tracking", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 44,
            body: "- Addresses: #102",
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 102,
            labels: [{ name: "internal/tracking" }, { name: "internal/user" }],
            state: "open",
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("already has 'internal/tracking' label"),
    ),
  );
});

test("runTrackingIssue - skips PR without valid Addresses line", async () => {
  const mockCore = createMockCore();

  let getIssueCalled = false;
  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 45,
            body: "Just a PR description with no issue reference",
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => {
          getIssueCalled = true;
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(getIssueCalled, false);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("does not contain valid issue references"),
    ),
  );
});

test("runTrackingIssue - fails when no release branches exist", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: { listBranches: async () => {} },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, "No branches found");
});

test("runTrackingIssue - respects target release branch from PR labels over latest branch", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }, { name: "release/v14" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 50,
            title: "Fix issue for v14",
            body: "- Addresses: #200",
            labels: [{ name: "release/v14" }],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            title: "Issue 200",
            labels: [],
            state: "open",
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], [
    "internal/tracking",
    "internal/user",
  ]);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #50 -->/);
  assert.match(
    commentsCreated[1],
    /A label has been added for the latest release branch/,
  );
  assert.deepStrictEqual(labelsAdded[1], ["release/v14"]);
});

test("runTrackingIssue - skips referenced issue that is closed", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 51,
            body: "- Addresses: #201",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 201,
            state: "closed",
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "Issue #201 referenced in PR #51 is already closed. Skipping.",
      ),
    ),
  );
});

test("runTrackingIssue - skips referenced item that is a pull request", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 52,
            body: "- Addresses: #202",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 202,
            state: "open",
            pull_request: {},
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "Issue #202 referenced in PR #52 is a pull request, not an issue. Skipping.",
      ),
    ),
  );
});

test("runTrackingIssue - handles multiple referenced issues where one already has internal/tracking", async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 53,
            body: "- Addresses: #100\n- Addresses: #101",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 100) {
            return {
              data: {
                number: 100,
                labels: [{ name: "internal/tracking" }],
                state: "open",
              },
            };
          }
          return {
            data: {
              number: 101,
              labels: [],
              state: "open",
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          labelsAdded.push({ issue_number, labels });
        },
        createComment: async ({ issue_number, body }) => {
          commentsCreated.push({ issue_number, body });
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("Issue #100 already has 'internal/tracking' label. Skipping."),
    ),
  );
  assert.strictEqual(labelsAdded.length, 2);
  assert.deepStrictEqual(labelsAdded[0], {
    issue_number: 101,
    labels: ["internal/tracking", "internal/user"],
  });
  assert.deepStrictEqual(labelsAdded[1], {
    issue_number: 101,
    labels: ["release/v15"],
  });
  assert.strictEqual(commentsCreated.length, 2);
  assert.strictEqual(commentsCreated[0].issue_number, 101);
  assert.strictEqual(commentsCreated[1].issue_number, 101);
  assert.match(commentsCreated[0].body, /<!-- tracking-pr: #53 -->/);
  assert.match(
    commentsCreated[1].body,
    /A label has been added for the latest release branch/,
  );
});

test("run via default export with SCRIPT_MODE=tracking-issue", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      SCRIPT_MODE: "tracking-issue",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v16" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  await defaultHandler({
    github: mockGithub,
    context: {},
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("Latest release branch detected: release/v16"),
    ),
  );
});

test("runTrackingIssue - fails when branches exist but none match release/v pattern", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async () => [
      { name: "main" },
      { name: "feature/foo" },
      { name: "release/v-custom" },
    ],
    rest: {
      repos: { listBranches: async () => {} },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, "No release branches found");
});

test("runTrackingIssue - logs warning and continues when fetching issue fails", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 60,
            body: "- Addresses: #300",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => {
          throw new Error("Not Found (404)");
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes("Could not fetch issue #300 for PR #60: Not Found (404)"),
    ),
  );
});

test("runTrackingIssue - records error and fails when addLabels throws", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 61,
            html_url:
              "https://github.com/rancher/terraform-provider-rancher2/pull/61",
            body: "- Addresses: #301",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 301,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {
          throw new Error("Secondary rate limit reached");
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.ok(mockCore.failedMessage);
  assert.ok(
    mockCore.failedMessage.includes("Failed to process some pull requests"),
  );
  assert.ok(mockCore.failedMessage.includes("Secondary rate limit reached"));
});

test("runTrackingIssue - handles PR description with inline code backticks", async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 62,
            body: "- Addresses: `#302`",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 302,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], [
    "internal/tracking",
    "internal/user",
  ]);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #62 -->/);
  assert.match(commentsCreated[0], /This is the tracking issue for PR #62/);
  assert.match(
    commentsCreated[1],
    /A label has been added for the latest release branch/,
  );
  assert.deepStrictEqual(labelsAdded[1], ["release/v15"]);
});

test("runTrackingIssue - picks highest release branch when multiple are on PR", async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v16" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 63,
            body: "- Addresses: #303",
            labels: [
              { name: "release/v13" },
              { name: "release/v15" },
              { name: "release/v14" },
            ],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 303,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async () => {},
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[1], ["release/v15"]);
});

test("runTrackingIssue - does not add release branch label if already present on issue", async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 64,
            body: "- Addresses: #304",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 304,
            state: "open",
            labels: [{ name: "release/v15" }],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 1);
  assert.deepStrictEqual(labelsAdded[0], [
    "internal/tracking",
    "internal/user",
  ]);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #64 -->/);
  assert.match(
    commentsCreated[1],
    /A label has been added for the latest release branch/,
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Release branch label 'release/v15' already present on issue #304",
      ),
    ),
  );
});

test("runTrackingIssue - handles malformed label structures gracefully", async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 65,
            body: "- Addresses: #305",
            labels: [
              null,
              undefined,
              { invalid: true },
              { name: "release/v15" },
            ],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 305,
            state: "open",
            labels: [null, { name: "internal/user" }],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async () => {},
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], ["internal/tracking"]);
  assert.deepStrictEqual(labelsAdded[1], ["release/v15"]);
});

test("runTrackingIssue - dispatches backport-issues.yml when adopting tracking issue", async () => {
  const mockCore = createMockCore();
  let dispatchedWorkflow = null;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 10,
            title: "Add feature",
            body: "- Addresses: #200",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 200,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async (params) => {
          dispatchedWorkflow = params;
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(dispatchedWorkflow, {
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    workflow_id: "backport-issues.yml",
    ref: "main",
    inputs: {
      issue_number: "200",
      release_label: "release/v15",
      pr_number: "10",
    },
  });
});

test("runTrackingIssue - records error and fails when workflow dispatch fails", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 11,
            title: "Fix something",
            body: "- Addresses: #201",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 201,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {
          throw new Error("API Rate Limit or Network Error");
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.match(
    mockCore.failedMessage,
    /Failed to dispatch 'backport-issues.yml' for issue #201: API Rate Limit or Network Error/,
  );
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "Failed to dispatch 'backport-issues.yml' for issue #201: API Rate Limit or Network Error",
      ),
    ),
  );
});

test("runTrackingIssue - uses GITHUB_MERGE_TOKEN and getOctokit when available", async () => {
  const mockCore = createMockCore();
  let tokenUsed = null;
  let customDispatchCalled = false;
  let customDispatchParams = null;

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: "custom-vault-token-123",
    },
  };

  const mockGetOctokit = (token) => {
    tokenUsed = token;
    return {
      rest: {
        actions: {
          createWorkflowDispatch: async (params) => {
            customDispatchCalled = true;
            customDispatchParams = params;
          },
        },
      },
    };
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 12,
            title: "Update provider",
            body: "- Addresses: #202",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 202,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {
          throw new Error(
            "Should not be called when custom octokit is provided",
          );
        },
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(tokenUsed, "custom-vault-token-123");
  assert.strictEqual(customDispatchCalled, true);
  assert.deepStrictEqual(customDispatchParams, {
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    workflow_id: "backport-issues.yml",
    ref: "main",
    inputs: {
      issue_number: "202",
      release_label: "release/v15",
      pr_number: "12",
    },
  });
});

test("runTrackingIssue - falls back to github.getOctokit when getOctokit argument is undefined", async () => {
  const mockCore = createMockCore();
  let tokenUsed = null;
  let customDispatchCalled = false;
  let customDispatchParams = null;

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: "token-via-github-octokit",
    },
  };

  const mockGithub = {
    getOctokit: (token) => {
      tokenUsed = token;
      return {
        rest: {
          actions: {
            createWorkflowDispatch: async (params) => {
              customDispatchCalled = true;
              customDispatchParams = params;
            },
          },
        },
      };
    },
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 13,
            title: "Fix issue",
            body: "- Addresses: #203",
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 203,
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(tokenUsed, "token-via-github-octokit");
  assert.strictEqual(customDispatchCalled, true);
  assert.deepStrictEqual(customDispatchParams, {
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    workflow_id: "backport-issues.yml",
    ref: "main",
    inputs: {
      issue_number: "203",
      release_label: "release/v15",
      pr_number: "13",
    },
  });
});

test("runTrackingIssue - targets specific PR when PR_NUMBER env variable is provided", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let searchCalled = false;
  let subIssueCreated = false;

  const mockProcess = {
    env: {
      PR_NUMBER: "42",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Specific PR",
              body: "- Addresses: #301",
              state: "open",
              draft: false,
              base: { ref: "main" },
              labels: [],
            },
          };
        },
      },
      issues: {
        get: async () => ({
          data: {
            number: 301,
            title: "Fix issue 301",
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 999, number: 302 } };
        },
      },
      actions: {
        createWorkflowDispatch: async () => {},
      },
    },
    request: async () => ({ data: [] }),
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(searchCalled, false);
  assert.strictEqual(subIssueCreated, true);
});

test("runTrackingIssue - falls back to search when targeted PR is ineligible", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let searchCalled = false;

  const mockProcess = {
    env: {
      PR_NUMBER: "43",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Draft PR with internal/ignore",
              body: "- Addresses: #301",
              state: "open",
              draft: true,
              base: { ref: "main" },
              labels: [{ name: "internal/ignore" }],
            },
          };
        },
      },
      actions: {
        createWorkflowDispatch: async () => {},
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(searchCalled, true);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("Target PR #43 is not eligible"),
    ),
  );
});

test("runTrackingIssue - falls back to runner github client when dispatchGithub throws 403", async () => {
  const mockCore = createMockCore();
  let runnerDispatchCalled = false;
  let customDispatchCalled = false;

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: "scoped-pull-requests-token",
    },
  };

  const mockGetOctokit = () => ({
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          customDispatchCalled = true;
          const err = new Error("Resource not accessible by integration");
          err.status = 403;
          throw err;
        },
      },
    },
  });

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      return [
        {
          number: 55,
          title: "PR with token fallback",
          body: "- Addresses: #401",
          labels: [],
        },
      ];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 401,
            title: "Issue 401",
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
        create: async () => ({ data: { id: 888, number: 402 } }),
      },
      actions: {
        createWorkflowDispatch: async () => {
          runnerDispatchCalled = true;
        },
      },
    },
    request: async () => ({ data: [] }),
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(customDispatchCalled, true);
  assert.strictEqual(runnerDispatchCalled, true);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Successfully dispatched 'backport-issues.yml' using fallback client",
      ),
    ),
  );
});

test("runTrackingIssue - does not fall back to runner client when dispatch error is not 403", async () => {
  const mockCore = createMockCore();
  let runnerDispatchCalled = false;
  let customDispatchCalled = false;

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: "scoped-pull-requests-token",
    },
  };

  const mockGetOctokit = () => ({
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          customDispatchCalled = true;
          const err = new Error("Gateway Timeout");
          err.status = 504;
          throw err;
        },
      },
    },
  });

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      return [
        {
          number: 56,
          title: "PR with non-403 error",
          body: "- Addresses: #403",
          labels: [],
        },
      ];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 403,
            title: "Issue 403",
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
        create: async () => ({ data: { id: 889, number: 404 } }),
      },
      actions: {
        createWorkflowDispatch: async () => {
          runnerDispatchCalled = true;
        },
      },
    },
    request: async () => ({ data: [] }),
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(customDispatchCalled, true);
  assert.strictEqual(runnerDispatchCalled, false);
  assert.ok(
    mockCore.failedMessage?.includes(
      "Failed to dispatch 'backport-issues.yml'",
    ),
  );
});

test("runTrackingIssue - resolves PR from commit SHA when workflow_run pull_requests is empty", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let commitLookupCalled = false;
  let searchCalled = false;
  let subIssueCreatedBeforeLabel = false;
  let releaseLabelAdded = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (
        method === mockGithub.rest.repos.listPullRequestsAssociatedWithCommit
      ) {
        commitLookupCalled = true;
        return [
          {
            number: 88,
            state: "open",
            base: { ref: "main" },
            labels: [],
          },
        ];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
        listPullRequestsAssociatedWithCommit: async () => {},
      },
      search: { issuesAndPullRequests: async () => {} },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Fork PR",
              body: "- Addresses: #500",
              state: "open",
              draft: false,
              base: { ref: "main" },
              labels: [],
            },
          };
        },
      },
      issues: {
        get: async () => ({
          data: {
            number: 500,
            title: "Issue 500",
            state: "open",
            labels: [],
          },
        }),
        create: async () => {
          if (!releaseLabelAdded) {
            subIssueCreatedBeforeLabel = true;
          }
          return { data: { id: 501, number: 501 } };
        },
        addLabels: async ({ labels }) => {
          if (labels.includes("release/v15")) {
            releaseLabelAdded = true;
          }
        },
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {},
      },
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return { data: [] };
      }
      return { data: {} };
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-sha-fork-88",
        pull_requests: [],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(commitLookupCalled, true);
  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(searchCalled, false);
  assert.strictEqual(subIssueCreatedBeforeLabel, true);
  assert.strictEqual(releaseLabelAdded, true);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("Resolved PR #88 from commit SHA commit-sha-fork-88"),
    ),
  );
});

test("runTrackingIssue - resolves PR from PR signal artifact file on disk when workflow_run pull_requests is empty", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let subIssueCreatedBeforeLabel = false;
  let releaseLabelAdded = false;

  const tmpSignalDir = fs.mkdtempSync(
    path.join(os.tmpdir(), "pr-signal-test-"),
  );
  const signalFilePath = path.join(tmpSignalDir, "pr-signal.json");
  fs.writeFileSync(
    signalFilePath,
    JSON.stringify({ pr_number: 99, head_sha: "signal-head-sha" }),
    "utf8",
  );

  try {
    const mockGithub = {
      paginate: async (method) => {
        if (method === mockGithub.rest.repos.listBranches) {
          return [{ name: "main" }, { name: "release/v15" }];
        }
        return [];
      },
      rest: {
        repos: {
          listBranches: async () => {},
        },
        pulls: {
          get: async ({ pull_number }) => {
            pullGetCalled = true;
            return {
              data: {
                number: pull_number,
                title: "PR with Signal",
                body: "- Addresses: #502",
                state: "open",
                draft: false,
                base: { ref: "main" },
                labels: [],
              },
            };
          },
        },
        issues: {
          get: async () => ({
            data: {
              number: 502,
              title: "Issue 502",
              state: "open",
              labels: [],
            },
          }),
          create: async () => {
            if (!releaseLabelAdded) {
              subIssueCreatedBeforeLabel = true;
            }
            return { data: { id: 503, number: 503 } };
          },
          addLabels: async ({ labels }) => {
            if (labels.includes("release/v15")) {
              releaseLabelAdded = true;
            }
          },
          createComment: async () => {},
        },
        actions: {
          createWorkflowDispatch: async () => {},
        },
      },
      request: async (endpoint) => {
        if (endpoint.includes("/sub_issues")) {
          return { data: [] };
        }
        return { data: {} };
      },
    };

    const mockContext = {
      eventName: "workflow_run",
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        workflow_run: {
          head_sha: "signal-head-sha",
          pull_requests: [],
        },
      },
    };

    const mockProcess = {
      env: {
        PR_SIGNAL_PATH: signalFilePath,
      },
      cwd: () => tmpSignalDir,
    };

    await runTrackingIssue({
      github: mockGithub,
      context: mockContext,
      core: mockCore,
      process: mockProcess,
    });

    assert.strictEqual(pullGetCalled, true);
    assert.strictEqual(subIssueCreatedBeforeLabel, true);
    assert.strictEqual(releaseLabelAdded, true);
    assert.ok(
      mockCore.infoMessages.some((m) =>
        m.includes("Resolved PR #99 from PR signal artifact"),
      ),
    );
  } finally {
    fs.rmSync(tmpSignalDir, { recursive: true, force: true });
  }
});

test("runTrackingIssue - ignores PR signal artifact file when head_sha does not match workflow_run head_sha", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;

  const tmpSignalDir = fs.mkdtempSync(
    path.join(os.tmpdir(), "pr-signal-test-"),
  );
  const signalFilePath = path.join(tmpSignalDir, "pr-signal.json");
  fs.writeFileSync(
    signalFilePath,
    JSON.stringify({ pr_number: 99, head_sha: "spoofed-sha" }),
    "utf8",
  );

  try {
    const mockGithub = {
      paginate: async (method) => {
        if (method === mockGithub.rest.repos.listBranches) {
          return [{ name: "main" }, { name: "release/v15" }];
        }
        return [];
      },
      rest: {
        repos: {
          listBranches: async () => {},
        },
        pulls: {
          get: async () => {
            pullGetCalled = true;
            return { data: {} };
          },
        },
      },
    };

    const mockContext = {
      eventName: "workflow_run",
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        workflow_run: {
          head_sha: "actual-sha",
          pull_requests: [],
        },
      },
    };

    const mockProcess = {
      env: {
        PR_SIGNAL_PATH: signalFilePath,
      },
      cwd: () => tmpSignalDir,
    };

    await runTrackingIssue({
      github: mockGithub,
      context: mockContext,
      core: mockCore,
      process: mockProcess,
    });

    assert.strictEqual(pullGetCalled, false);
    assert.ok(
      mockCore.warningMessages.some((m) =>
        m.includes(
          "PR signal head_sha 'spoofed-sha' does not match workflow_run head_sha 'actual-sha'",
        ),
      ),
    );
    assert.ok(
      mockCore.infoMessages.some((m) =>
        m.includes(
          "Could not identify a single PR from workflow_run event. Skipping adoption",
        ),
      ),
    );
  } finally {
    fs.rmSync(tmpSignalDir, { recursive: true, force: true });
  }
});

test("runTrackingIssue - ignores PR signal artifact file when head_sha is missing from artifact", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;

  const tmpSignalDir = fs.mkdtempSync(
    path.join(os.tmpdir(), "pr-signal-test-"),
  );
  const signalFilePath = path.join(tmpSignalDir, "pr-signal.json");
  fs.writeFileSync(signalFilePath, JSON.stringify({ pr_number: 99 }), "utf8");

  try {
    const mockGithub = {
      paginate: async (method) => {
        if (method === mockGithub.rest.repos.listBranches) {
          return [{ name: "main" }, { name: "release/v15" }];
        }
        return [];
      },
      rest: {
        repos: {
          listBranches: async () => {},
        },
        pulls: {
          get: async () => {
            pullGetCalled = true;
            return { data: {} };
          },
        },
      },
    };

    const mockContext = {
      eventName: "workflow_run",
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        workflow_run: {
          head_sha: "actual-sha",
          pull_requests: [],
        },
      },
    };

    const mockProcess = {
      env: {
        PR_SIGNAL_PATH: signalFilePath,
      },
      cwd: () => tmpSignalDir,
    };

    await runTrackingIssue({
      github: mockGithub,
      context: mockContext,
      core: mockCore,
      process: mockProcess,
    });

    assert.strictEqual(pullGetCalled, false);
    assert.ok(
      mockCore.warningMessages.some((m) =>
        m.includes(
          "PR signal head_sha '' does not match workflow_run head_sha 'actual-sha'",
        ),
      ),
    );
  } finally {
    fs.rmSync(tmpSignalDir, { recursive: true, force: true });
  }
});

test("runTrackingIssue - ignores PR signal artifact file when head_branch does not match workflow_run head_branch", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;

  const tmpSignalDir = fs.mkdtempSync(
    path.join(os.tmpdir(), "pr-signal-test-"),
  );
  const signalFilePath = path.join(tmpSignalDir, "pr-signal.json");
  fs.writeFileSync(
    signalFilePath,
    JSON.stringify({
      pr_number: 99,
      head_sha: "actual-sha",
      head_branch: "spoofed-branch",
      head_repo: "rancher/terraform-provider-rancher2",
    }),
    "utf8",
  );

  try {
    const mockGithub = {
      paginate: async (method) => {
        if (method === mockGithub.rest.repos.listBranches) {
          return [{ name: "main" }, { name: "release/v15" }];
        }
        return [];
      },
      rest: {
        repos: {
          listBranches: async () => {},
        },
        pulls: {
          get: async () => {
            pullGetCalled = true;
            return { data: {} };
          },
        },
      },
    };

    const mockContext = {
      eventName: "workflow_run",
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        workflow_run: {
          head_sha: "actual-sha",
          head_branch: "actual-branch",
          head_repository: { full_name: "rancher/terraform-provider-rancher2" },
          pull_requests: [],
        },
      },
    };

    const mockProcess = {
      env: {
        PR_SIGNAL_PATH: signalFilePath,
      },
      cwd: () => tmpSignalDir,
    };

    await runTrackingIssue({
      github: mockGithub,
      context: mockContext,
      core: mockCore,
      process: mockProcess,
    });

    assert.strictEqual(pullGetCalled, false);
    assert.ok(
      mockCore.warningMessages.some((m) =>
        m.includes(
          "PR signal head_branch 'spoofed-branch' does not match workflow_run head_branch 'actual-branch'",
        ),
      ),
    );
    assert.ok(
      mockCore.infoMessages.some((m) =>
        m.includes(
          "Could not identify a single PR from workflow_run event. Skipping adoption",
        ),
      ),
    );
  } finally {
    fs.rmSync(tmpSignalDir, { recursive: true, force: true });
  }
});

test("runTrackingIssue - ignores PR signal artifact file when head_repo does not match workflow_run head_repository", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;

  const tmpSignalDir = fs.mkdtempSync(
    path.join(os.tmpdir(), "pr-signal-test-"),
  );
  const signalFilePath = path.join(tmpSignalDir, "pr-signal.json");
  fs.writeFileSync(
    signalFilePath,
    JSON.stringify({
      pr_number: 99,
      head_sha: "actual-sha",
      head_branch: "actual-branch",
      head_repo: "attacker/terraform-provider-rancher2",
    }),
    "utf8",
  );

  try {
    const mockGithub = {
      paginate: async (method) => {
        if (method === mockGithub.rest.repos.listBranches) {
          return [{ name: "main" }, { name: "release/v15" }];
        }
        return [];
      },
      rest: {
        repos: {
          listBranches: async () => {},
        },
        pulls: {
          get: async () => {
            pullGetCalled = true;
            return { data: {} };
          },
        },
      },
    };

    const mockContext = {
      eventName: "workflow_run",
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        workflow_run: {
          head_sha: "actual-sha",
          head_branch: "actual-branch",
          head_repository: { full_name: "rancher/terraform-provider-rancher2" },
          pull_requests: [],
        },
      },
    };

    const mockProcess = {
      env: {
        PR_SIGNAL_PATH: signalFilePath,
      },
      cwd: () => tmpSignalDir,
    };

    await runTrackingIssue({
      github: mockGithub,
      context: mockContext,
      core: mockCore,
      process: mockProcess,
    });

    assert.strictEqual(pullGetCalled, false);
    assert.ok(
      mockCore.warningMessages.some((m) =>
        m.includes(
          "PR signal head_repo 'attacker/terraform-provider-rancher2' does not match workflow_run head_repo 'rancher/terraform-provider-rancher2'",
        ),
      ),
    );
    assert.ok(
      mockCore.infoMessages.some((m) =>
        m.includes(
          "Could not identify a single PR from workflow_run event. Skipping adoption",
        ),
      ),
    );
  } finally {
    fs.rmSync(tmpSignalDir, { recursive: true, force: true });
  }
});

test("runTrackingIssue - rejects PR adoption when fetched PR head branch does not match workflow_run head branch", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let subIssueCreated = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
      },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Mismatched Head Branch PR",
              body: "- Addresses: #502",
              state: "open",
              draft: false,
              base: { ref: "main" },
              head: {
                sha: "expected-head-sha",
                ref: "mismatched-branch",
                repo: { full_name: "rancher/terraform-provider-rancher2" },
              },
              labels: [],
            },
          };
        },
      },
      issues: {
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 503, number: 503 } };
        },
      },
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "expected-head-sha",
        head_branch: "expected-branch",
        head_repository: { full_name: "rancher/terraform-provider-rancher2" },
        pull_requests: [{ number: 99 }],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(subIssueCreated, false);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "Target PR #99 head branch 'mismatched-branch' does not match workflow_run head branch 'expected-branch'",
      ),
    ),
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Skipping tracking adoption for workflow_run PR #99 with mismatched head branch",
      ),
    ),
  );
});

test("runTrackingIssue - rejects PR adoption when fetched PR head repository does not match workflow_run head repository", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let subIssueCreated = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
      },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Mismatched Head Repo PR",
              body: "- Addresses: #502",
              state: "open",
              draft: false,
              base: { ref: "main" },
              head: {
                sha: "expected-head-sha",
                ref: "expected-branch",
                repo: {
                  full_name: "untrusted-fork/terraform-provider-rancher2",
                },
              },
              labels: [],
            },
          };
        },
      },
      issues: {
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 503, number: 503 } };
        },
      },
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "expected-head-sha",
        head_branch: "expected-branch",
        head_repository: { full_name: "rancher/terraform-provider-rancher2" },
        pull_requests: [{ number: 99 }],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(subIssueCreated, false);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "Target PR #99 head repository 'untrusted-fork/terraform-provider-rancher2' does not match workflow_run head repository 'rancher/terraform-provider-rancher2'",
      ),
    ),
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Skipping tracking adoption for workflow_run PR #99 with mismatched head repository",
      ),
    ),
  );
});

test("runTrackingIssue - rejects PR adoption when fetched PR head SHA does not match workflow_run head SHA", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let subIssueCreated = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
      },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Mismatched Head SHA PR",
              body: "- Addresses: #502",
              state: "open",
              draft: false,
              base: { ref: "main" },
              head: { sha: "different-head-sha" },
              labels: [],
            },
          };
        },
      },
      issues: {
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 503, number: 503 } };
        },
      },
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "expected-head-sha",
        pull_requests: [{ number: 99 }],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(subIssueCreated, false);
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes(
        "PR #99 head SHA 'different-head-sha' does not match workflow_run head SHA 'expected-head-sha'",
      ),
    ),
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Skipping tracking adoption for workflow_run PR #99 with mismatched head SHA",
      ),
    ),
  );
});

test("runTrackingIssue - skips search candidate PR when head_sha does not match workflow_run head_sha", async () => {
  const mockCore = createMockCore();
  let pullCalls = 0;
  let subIssueCreated = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (
        method === mockGithub.rest.repos.listPullRequestsAssociatedWithCommit
      ) {
        return [];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
        listPullRequestsAssociatedWithCommit: async () => {},
      },
      search: {
        issuesAndPullRequests: async () => ({
          data: {
            items: [
              { number: 101, state: "open", pull_request: {} },
              { number: 102, state: "open", pull_request: {} },
            ],
          },
        }),
      },
      pulls: {
        get: async ({ pull_number }) => {
          pullCalls++;
          if (pull_number === 101) {
            return {
              data: {
                number: 101,
                title: "Wrong SHA PR",
                state: "open",
                draft: false,
                base: { ref: "main" },
                head: { sha: "wrong-sha" },
                labels: [],
              },
            };
          }
          return {
            data: {
              number: 102,
              title: "Correct SHA PR",
              body: "- Addresses: #504",
              state: "open",
              draft: false,
              base: { ref: "main" },
              head: { sha: "matching-sha" },
              labels: [],
            },
          };
        },
      },
      issues: {
        get: async () => ({
          data: {
            number: 504,
            title: "Issue 504",
            state: "open",
            labels: [],
          },
        }),
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 505, number: 505 } };
        },
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {},
      },
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return { data: [] };
      }
      return { data: {} };
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "matching-sha",
        pull_requests: [],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(subIssueCreated, true);
  assert.ok(pullCalls >= 2);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Search candidate PR #101 head SHA 'wrong-sha' does not match commit SHA 'matching-sha'",
      ),
    ),
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Resolved PR #102 from search API for commit SHA matching-sha",
      ),
    ),
  );
});

test("runTrackingIssue - resolves PR from search API fallback when commit lookup returns empty", async () => {
  const mockCore = createMockCore();
  let pullGetCalled = false;
  let searchApiCalled = false;
  let subIssueCreatedBeforeLabel = false;
  let releaseLabelAdded = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (
        method === mockGithub.rest.repos.listPullRequestsAssociatedWithCommit
      ) {
        return [];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
        listPullRequestsAssociatedWithCommit: async () => {},
      },
      search: {
        issuesAndPullRequests: async () => {
          searchApiCalled = true;
          return {
            data: {
              items: [
                {
                  number: 101,
                  state: "open",
                  pull_request: {},
                },
              ],
            },
          };
        },
      },
      pulls: {
        get: async ({ pull_number }) => {
          pullGetCalled = true;
          return {
            data: {
              number: pull_number,
              title: "Search fallback PR",
              body: "- Addresses: #504",
              state: "open",
              draft: false,
              base: { ref: "main" },
              labels: [],
            },
          };
        },
      },
      issues: {
        get: async () => ({
          data: {
            number: 504,
            title: "Issue 504",
            state: "open",
            labels: [],
          },
        }),
        create: async () => {
          if (!releaseLabelAdded) {
            subIssueCreatedBeforeLabel = true;
          }
          return { data: { id: 505, number: 505 } };
        },
        addLabels: async ({ labels }) => {
          if (labels.includes("release/v15")) {
            releaseLabelAdded = true;
          }
        },
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {},
      },
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return { data: [] };
      }
      return { data: {} };
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "search-head-sha",
        pull_requests: [],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(searchApiCalled, true);
  assert.strictEqual(pullGetCalled, true);
  assert.strictEqual(subIssueCreatedBeforeLabel, true);
  assert.strictEqual(releaseLabelAdded, true);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Resolved PR #101 from search API for commit SHA search-head-sha",
      ),
    ),
  );
});

test("runTrackingIssue - skips adoption without broad search when workflow_run has no resolvable PR", async () => {
  const mockCore = createMockCore();
  let searchCalled = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      if (
        method === mockGithub.rest.repos.listPullRequestsAssociatedWithCommit
      ) {
        return [];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
        listPullRequestsAssociatedWithCommit: async () => {},
      },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "unknown-sha",
        pull_requests: [],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(searchCalled, false);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Could not identify a single PR from workflow_run event. Skipping adoption to prevent broad PR search.",
      ),
    ),
  );
});

test("runTrackingIssue - skips adoption without broad search when workflow_run PR is ineligible", async () => {
  const mockCore = createMockCore();
  let searchCalled = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      pulls: {
        get: async ({ pull_number }) => ({
          data: {
            number: pull_number,
            title: "Closed PR",
            state: "closed",
            base: { ref: "main" },
            labels: [],
          },
        }),
      },
    },
  };

  const mockContext = {
    eventName: "workflow_run",
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "closed-sha",
        pull_requests: [{ number: 99 }],
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });

  assert.strictEqual(searchCalled, false);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Skipping tracking adoption for ineligible workflow_run PR #99.",
      ),
    ),
  );
});

test("runTrackingIssue - fails and skips search when PR_NUMBER is non-numeric", async () => {
  const mockCore = createMockCore();
  let searchCalled = false;

  const mockProcess = {
    env: {
      PR_NUMBER: "abc",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(searchCalled, false);
  assert.ok(mockCore.failedMessage?.includes('Invalid PR_NUMBER: "abc"'));
});

test("runTrackingIssue - fails and skips search when PR_NUMBER contains trailing non-numeric characters", async () => {
  const mockCore = createMockCore();
  let searchCalled = false;

  const mockProcess = {
    env: {
      PR_NUMBER: "42abc",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(searchCalled, false);
  assert.ok(mockCore.failedMessage?.includes('Invalid PR_NUMBER: "42abc"'));
});

test("runTrackingIssue - fails and skips search when PR_NUMBER is zero or negative", async () => {
  const mockCore = createMockCore();
  let searchCalled = false;

  const mockProcess = {
    env: {
      PR_NUMBER: "0",
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: "main" }, { name: "release/v15" }];
      }
      searchCalled = true;
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(searchCalled, false);
  assert.ok(mockCore.failedMessage?.includes('Invalid PR_NUMBER: "0"'));
});
