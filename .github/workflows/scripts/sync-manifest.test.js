import assert from "node:assert/strict";
import test from "node:test";
import {
  isReleaseBranch,
  parseSemverTag,
  calculateNextMinorVersion,
  resolveReleaseAs,
  updateManifestContent,
  updateConfigContent,
  isSyncNeeded,
  formatPrBody,
  resolveTargetBranch,
  resolveReleaseTag,
  resolveReleaseSha,
  runSyncManifest,
} from "./sync-manifest.js";

test("isReleaseBranch - valid release branches", () => {
  assert.strictEqual(isReleaseBranch("release/v15"), true);
  assert.strictEqual(isReleaseBranch("release/v16"), true);
  assert.strictEqual(isReleaseBranch("release/v1"), true);
  assert.strictEqual(isReleaseBranch("release/v99"), true);
});

test("isReleaseBranch - invalid branches", () => {
  assert.strictEqual(isReleaseBranch("main"), false);
  assert.strictEqual(isReleaseBranch("master"), false);
  assert.strictEqual(isReleaseBranch("v15"), false);
  assert.strictEqual(isReleaseBranch("release/15"), false);
  assert.strictEqual(isReleaseBranch("release/v15.1"), false);
  assert.strictEqual(isReleaseBranch("release/v15-fix"), false);
  assert.strictEqual(isReleaseBranch(""), false);
  assert.strictEqual(isReleaseBranch(null), false);
  assert.strictEqual(isReleaseBranch(undefined), false);
});

test("parseSemverTag - valid standard tags", () => {
  const res1 = parseSemverTag("v15.1.2");
  assert.deepStrictEqual(res1, {
    tag: "v15.1.2",
    version: "15.1.2",
    major: 15,
    minor: 1,
    patch: 2,
    prerelease: null,
  });

  const res2 = parseSemverTag("15.1.2");
  assert.deepStrictEqual(res2, {
    tag: "v15.1.2",
    version: "15.1.2",
    major: 15,
    minor: 1,
    patch: 2,
    prerelease: null,
  });

  const res3 = parseSemverTag("v1.0.0-rc.1");
  assert.deepStrictEqual(res3, {
    tag: "v1.0.0-rc.1",
    version: "1.0.0-rc.1",
    major: 1,
    minor: 0,
    patch: 0,
    prerelease: "rc.1",
  });
});

test("parseSemverTag - invalid tags throw error", () => {
  assert.throws(() => parseSemverTag("v15"), /Invalid semver tag/);
  assert.throws(() => parseSemverTag("invalid"), /Invalid semver tag/);
  assert.throws(() => parseSemverTag(""), /Invalid semver tag/);
  assert.throws(() => parseSemverTag(null), /Invalid semver tag/);
});

test("updateManifestContent - updates version and formatting", () => {
  const initial = JSON.stringify({ ".": "15.0.0" }, null, 2) + "\n";
  const updated = updateManifestContent(initial, "15.1.2");
  assert.strictEqual(
    updated,
    JSON.stringify({ ".": "15.1.2" }, null, 2) + "\n",
  );

  const fromEmpty = updateManifestContent("", "16.0.0");
  assert.strictEqual(
    fromEmpty,
    JSON.stringify({ ".": "16.0.0" }, null, 2) + "\n",
  );
});

test("calculateNextMinorVersion - increments minor and resets patch", () => {
  assert.deepStrictEqual(calculateNextMinorVersion("15.1.2"), {
    version: "15.2.0",
    tag: "v15.2.0",
  });
  assert.deepStrictEqual(calculateNextMinorVersion("v14.0.0"), {
    version: "14.1.0",
    tag: "v14.1.0",
  });
});

test("resolveReleaseAs - resolves explicit input or computes next minor", () => {
  // Explicit input without v
  assert.strictEqual(
    resolveReleaseAs({ process: { env: { RELEASE_AS: "15.3.0" } } }),
    "v15.3.0",
  );
  // Explicit input with v
  assert.strictEqual(
    resolveReleaseAs({ process: { env: { RELEASE_AS: "v15.3.0" } } }),
    "v15.3.0",
  );
  // From context payload
  assert.strictEqual(
    resolveReleaseAs({
      context: { payload: { inputs: { release_as: "15.4.0" } } },
    }),
    "v15.4.0",
  );
  // Computed next minor from targetVersion
  assert.strictEqual(
    resolveReleaseAs({ targetVersion: "15.1.2" }),
    "v15.2.0",
  );
});

test("updateConfigContent - updates last-release-sha, release-as, and removes bootstrap-sha", () => {
  const initial =
    JSON.stringify(
      {
        "bootstrap-sha": "ab77bd40ab77bd40ab77bd40ab77bd40ab77bd40",
        "release-as": "v15.0.0",
        "release-type": "go",
        packages: { ".": {} },
      },
      null,
      2,
    ) + "\n";

  const updated = updateConfigContent(initial, {
    lastReleaseSha: "1111222233334444555566667777888899990000",
    releaseAs: "v15.3.0",
  });

  const parsed = JSON.parse(updated);
  assert.strictEqual(
    parsed["last-release-sha"],
    "1111222233334444555566667777888899990000",
  );
  assert.strictEqual(parsed["release-as"], "v15.3.0");
  assert.strictEqual("bootstrap-sha" in parsed, false);
  assert.strictEqual(parsed["release-type"], "go");
});

test("updateConfigContent - throws if lastReleaseSha is missing", () => {
  assert.throws(
    () => updateConfigContent("{}", {}),
    /lastReleaseSha is required/,
  );
});

test("isSyncNeeded - identifies whether synchronization is required", () => {
  const manifestSynced = JSON.stringify({ ".": "15.1.2" });
  const configSynced = JSON.stringify({
    "last-release-sha": "1111222233334444555566667777888899990000",
  });

  // Fully aligned
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: manifestSynced,
      configContent: configSynced,
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
    }),
    false,
  );

  // Manifest version mismatch
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: JSON.stringify({ ".": "15.0.0" }),
      configContent: configSynced,
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
    }),
    true,
  );

  // SHA mismatch
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: manifestSynced,
      configContent: JSON.stringify({
        "last-release-sha": "oldshaoldshaoldshaoldshaoldshaoldsha",
      }),
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
    }),
    true,
  );

  // Obsolete bootstrap-sha present
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: manifestSynced,
      configContent: JSON.stringify({
        "last-release-sha": "1111222233334444555566667777888899990000",
        "bootstrap-sha": "oldbootstrap",
      }),
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
    }),
    true,
  );

  // release-as mismatch
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: manifestSynced,
      configContent: JSON.stringify({
        "last-release-sha": "1111222233334444555566667777888899990000",
        "release-as": "v15.2.0",
      }),
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
      releaseAs: "v15.3.0",
    }),
    true,
  );

  // Malformed JSON returns true
  assert.strictEqual(
    isSyncNeeded({
      manifestContent: "not-json",
      configContent: configSynced,
      targetVersion: "15.1.2",
      releaseSha: "1111222233334444555566667777888899990000",
    }),
    true,
  );
});

test("formatPrBody - constructs detailed markdown", () => {
  const body = formatPrBody({
    branch: "release/v15",
    version: "15.1.2",
    releaseTag: "v15.1.2",
    releaseSha: "1111222233334444555566667777888899990000",
    oldManifestVersion: "15.0.0",
    oldReleaseSha: "oldsha",
    removedBootstrapSha: true,
  });

  assert.match(body, /release\/v15/);
  assert.match(body, /15\.1\.2/);
  assert.match(body, /1111222233334444555566667777888899990000/);
  assert.match(body, /Removed obsolete `bootstrap-sha`/);
});

test("resolveTargetBranch - resolves from various sources", () => {
  assert.strictEqual(
    resolveTargetBranch({ process: { env: { TARGET_BRANCH: "release/v15" } } }),
    "release/v15",
  );
  assert.strictEqual(
    resolveTargetBranch({ process: { env: { INPUT_BRANCH: "release/v16" } } }),
    "release/v16",
  );
  assert.strictEqual(
    resolveTargetBranch({
      context: { payload: { inputs: { branch: "release/v14" } } },
    }),
    "release/v14",
  );
  assert.strictEqual(
    resolveTargetBranch({
      context: { payload: { release: { target_commitish: "release/v13" } } },
    }),
    "release/v13",
  );
  assert.strictEqual(
    resolveTargetBranch({ context: { ref: "refs/heads/release/v12" } }),
    "release/v12",
  );
});

test("resolveReleaseTag - resolves tag from inputs, env, or git tags", async () => {
  // From env RELEASE_TAG
  const tag1 = await resolveReleaseTag({
    process: { env: { RELEASE_TAG: "v15.1.2" } },
  });
  assert.strictEqual(tag1, "v15.1.2");

  // From env without v
  const tag2 = await resolveReleaseTag({
    process: { env: { RELEASE_TAG: "15.1.2" } },
  });
  assert.strictEqual(tag2, "v15.1.2");

  // From RELEASE_VERSION
  const tag3 = await resolveReleaseTag({
    process: { env: { RELEASE_VERSION: "15.1.2" } },
  });
  assert.strictEqual(tag3, "v15.1.2");

  // From git tag list
  const mockGit = async (args) => {
    if (args[0] === "tag") {
      return { stdout: "v15.1.2\nv15.1.1\nv15.1.0\n" };
    }
    return { stdout: "" };
  };
  const tag4 = await resolveReleaseTag({
    targetBranch: "release/v15",
    gitExec: mockGit,
  });
  assert.strictEqual(tag4, "v15.1.2");

  // Filters out -rc.* pre-release tags and selects latest stable release tag
  const mockGitWithRc = async (args) => {
    if (args[0] === "tag") {
      return {
        stdout: "v15.3.0-rc.0\nv15.2.0-rc.4\nv15.1.2\nv15.1.1\nv15.1.0\n",
      };
    }
    return { stdout: "" };
  };
  const tag5 = await resolveReleaseTag({
    targetBranch: "release/v15",
    gitExec: mockGitWithRc,
  });
  assert.strictEqual(tag5, "v15.1.2");

  // Rejects manual inputs with a pre-release version
  await assert.rejects(
    () =>
      resolveReleaseTag({
        process: { env: { INPUT_VERSION: "15.3.0-rc.0" } },
      }),
    /manifest synchronization only targets full stable releases/,
  );

  // Rejects branches where only pre-release tags exist
  const mockGitOnlyRc = async (args) => {
    if (args[0] === "tag") {
      return { stdout: "v15.3.0-rc.0\nv15.2.0-rc.4\n" };
    }
    return { stdout: "" };
  };
  await assert.rejects(
    () =>
      resolveReleaseTag({
        targetBranch: "release/v15",
        gitExec: mockGitOnlyRc,
      }),
    /Could not resolve release tag from environment/,
  );
});

test("resolveReleaseSha - resolves commit SHA from env or git", async () => {
  const sha = "1111222233334444555566667777888899990000";

  // From env
  const sha1 = await resolveReleaseSha({
    releaseTag: "v15.1.2",
    process: { env: { RELEASE_SHA: sha } },
  });
  assert.strictEqual(sha1, sha);

  // From git rev-parse
  const mockGit = async (args) => {
    if (args[0] === "rev-parse") {
      return { stdout: `${sha}\n` };
    }
    return { stdout: "" };
  };
  const sha2 = await resolveReleaseSha({
    releaseTag: "v15.1.2",
    gitExec: mockGit,
  });
  assert.strictEqual(sha2, sha);

  // Fails with clear diagnostic when tag does not exist
  const mockGitFail = async () => {
    throw new Error("unknown revision");
  };
  await assert.rejects(
    () => resolveReleaseSha({ releaseTag: "v15.3.0", gitExec: mockGitFail }),
    /Release tag 'v15\.3\.0' does not exist in the repository/,
  );
});

test("runSyncManifest - skips non-release branches cleanly", async () => {
  const logged = [];
  const core = {
    info: (msg) => logged.push(msg),
    warning: () => {},
    error: () => {},
  };
  const result = await runSyncManifest({
    context: { ref: "refs/heads/main" },
    core,
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "not-release-branch",
    branch: "main",
  });
  assert.ok(
    logged.some((m) => m.includes("is not a maintenance release branch")),
  );
});

test("runSyncManifest - skips pre-release tags cleanly without creating PR", async () => {
  const logged = [];
  const core = {
    info: (msg) => logged.push(msg),
    warning: () => {},
    error: () => {},
  };
  const mockGit = async () => ({ stdout: "" });

  const result = await runSyncManifest({
    process: {
      env: {
        TARGET_BRANCH: "release/v15",
        RELEASE_TAG: "v15.3.0-rc.0",
      },
    },
    core,
    gitExec: mockGit,
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "prerelease-tag",
    tag: "v15.3.0-rc.0",
    branch: "release/v15",
  });
  assert.ok(logged.some((m) => m.includes("is a pre-release")));
});

test("runSyncManifest - detects already synchronized branch and exits without PR", async () => {
  const targetBranch = "release/v15";
  const releaseTag = "v15.1.2";
  const releaseSha = "1111222233334444555566667777888899990000";

  const files = {
    ".release-please-manifest.json": JSON.stringify({ ".": "15.1.2" }, null, 2),
    "release-please-config.json": JSON.stringify(
      { "last-release-sha": releaseSha, "release-as": "v15.2.0" },
      null,
      2,
    ),
  };

  const mockFs = {
    existsSync: (p) => Object.keys(files).some((k) => p.endsWith(k)),
    readFileSync: (p) => {
      const key = Object.keys(files).find((k) => p.endsWith(k));
      return files[key];
    },
    writeFileSync: () => {
      assert.fail("writeFileSync should not be called when already synced");
    },
  };

  const gitCalls = [];
  const mockGit = async (args) => {
    gitCalls.push(args);
    return { stdout: "" };
  };

  const logged = [];
  const core = {
    info: (msg) => logged.push(msg),
    warning: () => {},
    error: () => {},
  };

  const result = await runSyncManifest({
    process: {
      env: {
        TARGET_BRANCH: targetBranch,
        RELEASE_TAG: releaseTag,
        RELEASE_SHA: releaseSha,
      },
    },
    core,
    gitExec: mockGit,
    fsModule: mockFs,
  });

  assert.strictEqual(result.synced, true);
  assert.strictEqual(result.prCreated, false);
  assert.strictEqual(result.version, "15.1.2");
  assert.ok(logged.some((m) => m.includes("is already synchronized")));
});

test("runSyncManifest - performs synchronization, commits, pushes and opens PR", async () => {
  const targetBranch = "release/v15";
  const releaseTag = "v15.1.2";
  const releaseSha = "1111222233334444555566667777888899990000";

  const writtenFiles = {};
  const initialFiles = {
    ".release-please-manifest.json": JSON.stringify({ ".": "15.0.0" }, null, 2),
    "release-please-config.json": JSON.stringify(
      {
        "bootstrap-sha": "oldbootstrapsha",
        "last-release-sha": "oldsha",
      },
      null,
      2,
    ),
  };

  const mockFs = {
    existsSync: (p) => Object.keys(initialFiles).some((k) => p.endsWith(k)),
    readFileSync: (p) => {
      const key = Object.keys(initialFiles).find((k) => p.endsWith(k));
      return writtenFiles[key] || initialFiles[key];
    },
    writeFileSync: (p, content) => {
      const key = Object.keys(initialFiles).find((k) => p.endsWith(k));
      writtenFiles[key] = content;
    },
  };

  const gitCalls = [];
  const mockGit = async (args) => {
    gitCalls.push(args);
    return { stdout: "" };
  };

  const prCreated = [];
  const labelsAdded = [];
  const mockGithub = {
    rest: {
      pulls: {
        list: async () => ({ data: [] }),
        create: async (params) => {
          prCreated.push(params);
          return {
            data: {
              number: 100,
              html_url: "https://github.com/rancher/pr/100",
            },
          };
        },
      },
      issues: {
        addLabels: async (params) => {
          labelsAdded.push(params);
          return { data: [] };
        },
      },
    },
  };

  const logged = [];
  const secrets = [];
  const core = {
    info: (msg) => logged.push(msg),
    warning: () => {},
    error: () => {},
    setSecret: (s) => secrets.push(s),
  };

  const result = await runSyncManifest({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    process: {
      env: {
        TARGET_BRANCH: targetBranch,
        RELEASE_TAG: releaseTag,
        RELEASE_SHA: releaseSha,
        GITHUB_MERGE_TOKEN: "secret-token",
      },
    },
    core,
    gitExec: mockGit,
    fsModule: mockFs,
  });

  assert.strictEqual(result.prCreated, true);
  assert.strictEqual(result.prNumber, 100);

  // Verify file writes
  const updatedManifest = JSON.parse(
    writtenFiles[".release-please-manifest.json"],
  );
  assert.strictEqual(updatedManifest["."], "15.1.2");

  const updatedConfig = JSON.parse(writtenFiles["release-please-config.json"]);
  assert.strictEqual(updatedConfig["last-release-sha"], releaseSha);
  assert.strictEqual(updatedConfig["release-as"], "v15.2.0");
  assert.strictEqual("bootstrap-sha" in updatedConfig, false);

  // Verify git operations
  assert.ok(
    gitCalls.some(
      (c) =>
        c[0] === "checkout" &&
        c[1] === "-B" &&
        c[2] === "chore/sync-release-please-release-v15",
    ),
  );
  assert.ok(gitCalls.some((c) => c[0] === "commit"));
  assert.ok(
    gitCalls.some(
      (c) =>
        c[0] === "push" &&
        c[1] === "-f" &&
        c[3] === "chore/sync-release-please-release-v15",
    ),
  );

  // Verify PR details
  assert.strictEqual(prCreated.length, 1);
  assert.strictEqual(prCreated[0].base, "release/v15");
  assert.strictEqual(
    prCreated[0].head,
    "chore/sync-release-please-release-v15",
  );
  assert.strictEqual(
    prCreated[0].title,
    "chore(release/v15): update release-please manifest to 15.1.2",
  );

  // Verify label added
  assert.strictEqual(labelsAdded.length, 1);
  assert.strictEqual(labelsAdded[0].issue_number, 100);
  assert.deepStrictEqual(labelsAdded[0].labels, ["internal/pr-release-chore"]);
});

test("runSyncManifest - updates branch without duplicate PR when open PR already exists", async () => {
  const targetBranch = "release/v15";
  const releaseTag = "v15.1.2";
  const releaseSha = "1111222233334444555566667777888899990000";

  const initialFiles = {
    ".release-please-manifest.json": JSON.stringify({ ".": "15.0.0" }, null, 2),
    "release-please-config.json": JSON.stringify(
      { "last-release-sha": "oldsha" },
      null,
      2,
    ),
  };

  const writtenFiles = {};
  const mockFs = {
    existsSync: () => true,
    readFileSync: (p) => {
      const key = Object.keys(initialFiles).find((k) => p.endsWith(k));
      return writtenFiles[key] || initialFiles[key];
    },
    writeFileSync: (p, content) => {
      const key = Object.keys(initialFiles).find((k) => p.endsWith(k));
      writtenFiles[key] = content;
    },
  };

  const mockGit = async () => ({ stdout: "" });

  const mockGithub = {
    rest: {
      pulls: {
        list: async () => ({
          data: [
            {
              number: 200,
              head: { ref: "chore/sync-release-please-release-v15" },
              base: { ref: "release/v15" },
              html_url: "https://github.com/rancher/pr/200",
            },
          ],
        }),
        create: async () => {
          assert.fail(
            "pulls.create should not be called when PR already exists",
          );
        },
      },
    },
  };

  const logged = [];
  const core = {
    info: (msg) => logged.push(msg),
    warning: () => {},
    error: () => {},
  };

  const result = await runSyncManifest({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    process: {
      env: {
        TARGET_BRANCH: targetBranch,
        RELEASE_TAG: releaseTag,
        RELEASE_SHA: releaseSha,
      },
    },
    core,
    gitExec: mockGit,
    fsModule: mockFs,
  });

  assert.strictEqual(result.prCreated, false);
  assert.strictEqual(result.prNumber, 200);
  assert.ok(logged.some((m) => m.includes("Open PR #200 already exists")));
});
