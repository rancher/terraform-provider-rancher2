import fs from "node:fs/promises";
import path from "node:path";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

/**
 * Standard caller stubs pointing strictly to @main.
 */
export const CALLER_STUBS = {
  ".github/workflows/pull_request.yaml": `name: Pull Request

on:
  pull_request:
  merge_group:

permissions: {}

jobs:
  pull-request:
    permissions:
      contents: read
      pull-requests: read
    uses: rancher/terraform-provider-rancher2/.github/workflows/pull_request.yaml@main
    secrets: inherit
`,
  ".github/workflows/release.yml": `name: Release Please

on:
  push:
    branches:
      - 'release/v*'

permissions: {}

jobs:
  release:
    permissions:
      contents: write
      id-token: write
      issues: write
      pull-requests: write
      actions: read
    uses: rancher/terraform-provider-rancher2/.github/workflows/release-please.yml@main
    secrets: inherit
`,
  ".github/workflows/e2e-pr-listener.yml": `name: PR Executor Listener
# This unprivileged workflow acts as a middleman to securely trigger e2e tests.

on:
  pull_request:
    branches:
      - 'release/v*'
    types: [opened, synchronize, reopened]

permissions: {}

jobs:
  emit-event:
    name: 'Emit E2E Workflow Run Event'
    runs-on: ubuntu-latest
    timeout-minutes: 5
    permissions:
      contents: read
    steps:
      - name: 'Execute PR Listener'
        run: echo "PR listener executed. The completion of this workflow will trigger the privileged E2E Tests (Release PRs) workflow safely from the base branch via the workflow_run event."
`,
  ".github/workflows/sync-release-manifest.yml": `name: "Sync Release Please Manifest Caller"

on:
  workflow_dispatch:
    inputs:
      branch:
        description: "Target release branch"
        required: true
        type: string
      version:
        description: "Last published full release version"
        required: false
        type: string
      release_as:
        description: "Target upcoming release version for release-as"
        required: false
        type: string

permissions: {}

jobs:
  sync:
    permissions:
      contents: write
      pull-requests: write
      id-token: write
    uses: rancher/terraform-provider-rancher2/.github/workflows/sync-release-manifest.yml@main
    with:
      branch: \${{ inputs.branch }}
      version: \${{ inputs.version }}
      release_as: \${{ inputs.release_as }}
    secrets: inherit
`,
};

/**
 * Allowed stub filenames inside .github/workflows.
 */
export const ALLOWED_STUB_FILENAMES = new Set([
  "pull_request.yaml",
  "release.yml",
  "e2e-pr-listener.yml",
  "sync-release-manifest.yml",
]);

/**
 * Executes a git command asynchronously with error sanitization.
 */
export async function runGit(args, options = {}) {
  try {
    return await execFileAsync("git", args, { encoding: "utf8", ...options });
  } catch (error) {
    const rawDetails =
      error.stderr?.trim() || error.stdout?.trim() || error.message;
    const details = rawDetails
      ?.replace(
        /authorization:\s*basic\s+[^\s'"]+/gi,
        "AUTHORIZATION: basic [REDACTED]",
      )
      ?.replace(
        /https:\/\/[^@\s]+@github\.com/gi,
        "https://[REDACTED]@github.com",
      );
    const sanitizedArgs = args.map((arg) =>
      arg.toLowerCase().includes("authorization:") ? "[REDACTED]" : arg,
    );
    const sanitizedCommand = `git ${sanitizedArgs.join(" ")}`;
    throw new Error(`'${sanitizedCommand}' failed: ${details}`);
  }
}

/**
 * Validates and parses release branch input.
 * Strictly accepts canonical 'release/v<version>' (e.g. 'release/v16').
 * Chomps surrounding whitespace; rejects shorthands like 'v16', '16', 'release/16', leading zeros, or arbitrary branch names.
 */
export function parseReleaseBranch(input) {
  if (!input || typeof input !== "string") {
    throw new Error(
      `Invalid branch format: '${input}'. Branch must strictly match 'release/v<version>' (e.g. 'release/v16'). Shorthands like 'v16' or '16' are not permitted.`,
    );
  }

  const clean = input.trim();
  const match = clean.match(/^release\/v([1-9]\d*)$/);
  if (!match) {
    throw new Error(
      `Invalid branch format: '${clean}'. Branch must strictly match 'release/v<version>' (e.g. 'release/v16'). Shorthands like 'v16' or '16' and leading zeros are not permitted.`,
    );
  }

  const versionString = match[1];
  const versionNumber = parseInt(versionString, 10);

  return {
    branch: `release/v${versionString}`,
    versionNumber,
    versionString,
    releaseAs: `v${versionNumber}.0.0`,
    manifestVersion: `${versionNumber}.0.0`,
  };
}

/**
 * Resolves the commit SHA where the release branch was created from origin/main.
 * Computes 'git merge-base origin/main origin/<branch>' with fallback to 'git rev-parse origin/<branch>'.
 * Throws if the SHA cannot be resolved.
 */
export async function resolveBranchCreationSha(
  git,
  branch,
  baseRef = "origin/main",
) {
  try {
    const mergeBaseResult = await git([
      "merge-base",
      baseRef,
      `origin/${branch}`,
    ]);
    const sha = mergeBaseResult?.stdout?.trim() || "";
    if (/^[0-9a-f]{40}$/i.test(sha)) {
      return sha;
    }
  } catch {
    // If merge-base fails, fallback to rev-parse on origin/<branch>
  }

  try {
    const revResult = await git(["rev-parse", `origin/${branch}`]);
    const sha = revResult?.stdout?.trim() || "";
    if (/^[0-9a-f]{40}$/i.test(sha)) {
      return sha;
    }
  } catch {
    // ignore
  }

  throw new Error(
    `Failed to resolve branch creation SHA for '${branch}' against '${baseRef}'.`,
  );
}

/**
 * Resolves the commit SHA of the tag with the latest release history.
 * Strictly ignores pre-release / RC tags (e.g. -rc, -rc.1) and selects the highest stable semver tag prior to maxMajor (e.g. v15.3.0 for v16).
 */
export async function resolveLatestReleaseTagSha(git, maxMajor = null) {
  try {
    const tagResult = await git(["tag", "-l", "v*"]);
    const tags = (tagResult?.stdout || "")
      .trim()
      .split("\n")
      .map((t) => t.trim())
      .filter(Boolean);

    // Strictly filter out any -rc or pre-release tags, accepting only stable releases (e.g. v15.3.0)
    const stableTags = tags
      .filter(
        (t) =>
          !t.includes("-rc") && !t.includes("-") && /^v?\d+\.\d+\.\d+$/.test(t),
      )
      .map((t) => {
        const match = t.match(/^v?(\d+)\.(\d+)\.(\d+)$/);
        return {
          tag: t,
          major: parseInt(match[1], 10),
          minor: parseInt(match[2], 10),
          patch: parseInt(match[3], 10),
        };
      })
      .filter((item) =>
        maxMajor !== null && maxMajor !== undefined
          ? item.major < maxMajor
          : true,
      )
      .sort((a, b) => {
        if (b.major !== a.major) return b.major - a.major;
        if (b.minor !== a.minor) return b.minor - a.minor;
        return b.patch - a.patch;
      });

    if (stableTags.length > 0) {
      const latestTag = stableTags[0].tag;
      const revResult = await git(["rev-list", "-n", "1", latestTag]);
      const sha = revResult?.stdout?.trim() || "";
      if (/^[0-9a-f]{40}$/i.test(sha)) {
        return sha;
      }
    }
  } catch {
    // Return null if tags cannot be queried or resolved
  }
  return null;
}

/**
 * Builds release-please-config.json content with release-as version and cleans obsolete sha keys.
 */
export function buildReleasePleaseConfig(
  existingConfigOrJson,
  releaseAs,
  lastReleaseSha = null,
) {
  let config = {};
  if (typeof existingConfigOrJson === "string") {
    if (existingConfigOrJson.trim()) {
      try {
        config = JSON.parse(existingConfigOrJson);
      } catch (err) {
        throw new Error(
          `Failed to parse existing release-please-config.json: ${err.message}`,
        );
      }
    }
  } else if (existingConfigOrJson && typeof existingConfigOrJson === "object") {
    config = JSON.parse(JSON.stringify(existingConfigOrJson));
  }

  config["release-as"] = releaseAs;
  if (!("always-update" in config)) config["always-update"] = true;
  if (!("commit-search-depth" in config)) config["commit-search-depth"] = 100;
  if (!("prerelease" in config)) config["prerelease"] = true;
  if (!("draft" in config)) config["draft"] = true;
  if (!("force-tag-creation" in config)) config["force-tag-creation"] = true;
  if (!("include-v-in-tag" in config)) config["include-v-in-tag"] = true;
  if (!("include-component-in-tag" in config))
    config["include-component-in-tag"] = false;

  if (!config.packages || typeof config.packages !== "object") {
    config.packages = {
      ".": {
        "release-type": "go",
      },
    };
  }

  if (!("$schema" in config)) {
    config["$schema"] =
      "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json";
  }

  delete config["bootstrap-sha"];
  if (lastReleaseSha) {
    config["last-release-sha"] = lastReleaseSha;
  } else {
    delete config["last-release-sha"];
  }

  return JSON.stringify(config, null, 2) + "\n";
}

/**
 * Builds .release-please-manifest.json content.
 */
export function buildReleasePleaseManifest(manifestVersion) {
  const manifest = {
    ".": manifestVersion,
  };
  return JSON.stringify(manifest, null, 2) + "\n";
}

/**
 * Returns standard caller stubs mapping.
 */
export function buildWorkflowStubs() {
  return { ...CALLER_STUBS };
}

/**
 * Removes any non-stub workflows from the workflows directory.
 */
export async function cleanupNonStubWorkflows(workflowsDir) {
  try {
    const entries = await fs.readdir(workflowsDir);
    for (const entry of entries) {
      if (entry.endsWith(".yml") || entry.endsWith(".yaml")) {
        if (!ALLOWED_STUB_FILENAMES.has(entry)) {
          await fs.unlink(path.join(workflowsDir, entry));
        }
      }
    }
  } catch (err) {
    if (err.code !== "ENOENT") {
      throw err;
    }
  }
}

/**
 * Generates pull request description body.
 */
export function formatPrBody({ branch, releaseAs, manifestVersion }) {
  return `### Release Branch Initialization: \`${branch}\`

This automated pull request initializes configuration and caller stubs for the new release branch \`${branch}\`.

#### Summary of Changes:
- **Release Please Config**: Initialized \`release-please-config.json\` with \`release-as: ${releaseAs}\` and draft release enabled.
- **Release Please Manifest**: Set \`.release-please-manifest.json\` root version to \`${manifestVersion}\`.
- **Workflow Caller Stubs**: Replaced standalone workflows with thin caller stubs delegating execution strictly to \`@main\`:
  - \`.github/workflows/pull_request.yaml\`
  - \`.github/workflows/release.yml\`
  - \`.github/workflows/sync-release-manifest.yml\`
  - \`.github/workflows/e2e-pr-listener.yml\` (triggers privileged E2E workflows)

---
*Created automatically by [setup-release-branch.yml](https://github.com/rancher/terraform-provider-rancher2/actions/workflows/setup-release-branch.yml).*`;
}

/**
 * Main release branch creation and setup orchestrator.
 */
export async function runSetupReleaseBranch({
  github,
  context,
  core,
  process: proc = process,
  gitExec,
  workDir,
} = {}) {
  const branchInput =
    proc?.env?.INPUT_BRANCH || context?.payload?.inputs?.branch;
  const parsed = parseReleaseBranch(branchInput);

  const owner =
    context?.repo?.owner ||
    proc?.env?.GITHUB_REPOSITORY?.split("/")[0] ||
    "rancher";
  const repo =
    context?.repo?.repo ||
    proc?.env?.GITHUB_REPOSITORY?.split("/")[1] ||
    "terraform-provider-rancher2";

  const rootDir = workDir || proc?.cwd?.() || process.cwd();
  const git =
    gitExec || ((args, opts) => runGit(args, { cwd: rootDir, ...opts }));

  core?.info?.(
    `Configuring git identity and preparing release branch '${parsed.branch}'...`,
  );
  await git(["config", "user.name", "github-actions[bot]"]);
  await git([
    "config",
    "user.email",
    "41898282+github-actions[bot]@users.noreply.github.com",
  ]);

  // Ensure origin/main and tags are fetched
  await git(["fetch", "origin", "main", "--tags"]);

  // Check if target release branch already exists on remote origin
  const lsRemoteResult = await git([
    "ls-remote",
    "--heads",
    "origin",
    `refs/heads/${parsed.branch}`,
  ]);
  const branchExistsOnOrigin = Boolean(
    lsRemoteResult?.stdout &&
    lsRemoteResult.stdout.includes(`refs/heads/${parsed.branch}`),
  );

  if (branchExistsOnOrigin) {
    throw new Error(
      `Release branch '${parsed.branch}' already exists on remote origin. Release branch setup workflow cannot be run on existing branches; configurations must be updated manually.`,
    );
  }

  // Pre-calculate lastReleaseSha before pushing the remote branch so failures don't leave an unconfigured branch
  let lastReleaseSha = null;
  const envSha = proc?.env?.LAST_RELEASE_SHA?.trim();
  if (envSha && /^[0-9a-f]{40}$/i.test(envSha)) {
    lastReleaseSha = envSha;
  } else {
    // Determine the branching point from origin/main. For new branches, this is the current origin/main.
    lastReleaseSha = await resolveBranchCreationSha(git, "main", "origin/main");
  }

  if (!lastReleaseSha || !/^[0-9a-f]{40}$/i.test(lastReleaseSha)) {
    throw new Error(
      `Failed to resolve required last-release-sha anchor for release branch '${parsed.branch}'.`,
    );
  }

  core?.info?.(
    `Release branch '${parsed.branch}' does not exist on origin. Creating strictly from 'origin/main'...`,
  );
  // Invariant: New release branch must be branched strictly and exclusively from origin/main
  await git(["push", "origin", `origin/main:refs/heads/${parsed.branch}`]);

  // Fetch the release branch from origin
  await git(["fetch", "origin", parsed.branch]);

  // Checkout setup chore branch based on origin/<branch>
  const choreBranch = `chore/setup-release-v${parsed.versionNumber}`;
  core?.info?.(
    `Checking out setup branch '${choreBranch}' based on 'origin/${parsed.branch}'...`,
  );
  await git(["checkout", "-B", choreBranch, `origin/${parsed.branch}`]);

  // 1. Configure release-please-config.json
  const configPath = path.join(rootDir, "release-please-config.json");
  let existingConfig = "";
  try {
    existingConfig = await fs.readFile(configPath, "utf8");
  } catch (err) {
    if (err.code !== "ENOENT") throw err;
  }

  core?.info?.(
    `Resolved release anchor SHA '${lastReleaseSha}' for release-please-config.json`,
  );
  const updatedConfig = buildReleasePleaseConfig(
    existingConfig,
    parsed.releaseAs,
    lastReleaseSha,
  );
  await fs.writeFile(configPath, updatedConfig, "utf8");

  // 2. Configure .release-please-manifest.json
  const manifestPath = path.join(rootDir, ".release-please-manifest.json");
  const updatedManifest = buildReleasePleaseManifest(parsed.manifestVersion);
  await fs.writeFile(manifestPath, updatedManifest, "utf8");

  // 3. Configure workflows and caller stubs
  const workflowsDir = path.join(rootDir, ".github", "workflows");
  await fs.mkdir(workflowsDir, { recursive: true });
  await cleanupNonStubWorkflows(workflowsDir);

  const stubs = buildWorkflowStubs();
  for (const [relPath, content] of Object.entries(stubs)) {
    const fullPath = path.join(rootDir, relPath);
    await fs.mkdir(path.dirname(fullPath), { recursive: true });
    await fs.writeFile(fullPath, content, "utf8");
  }

  // 4. Git add and commit
  await git([
    "add",
    "-A",
    ".github/workflows",
    "release-please-config.json",
    ".release-please-manifest.json",
  ]);

  const statusResult = await git(["status", "--porcelain"]);
  const hasChanges = Boolean(
    statusResult?.stdout && statusResult.stdout.trim().length > 0,
  );

  let choreBranchPushed = hasChanges;

  if (hasChanges) {
    core?.info?.(`Committing and pushing changes to '${choreBranch}'...`);
    await git([
      "commit",
      "-m",
      `chore(${parsed.branch}): initialize release branch configuration`,
    ]);
    await git(["push", "-u", "origin", choreBranch, "--force"]);
    choreBranchPushed = true;
  } else {
    core?.info?.(
      `No file modifications needed; setup files already match target configuration.`,
    );
    // Check if chore branch exists on remote origin and is ahead of the release branch
    const remoteChoreResult = await git([
      "ls-remote",
      "--heads",
      "origin",
      `refs/heads/${choreBranch}`,
    ]);
    const choreExists = Boolean(
      remoteChoreResult?.stdout &&
      remoteChoreResult.stdout.includes(`refs/heads/${choreBranch}`),
    );

    if (choreExists) {
      try {
        await git(["fetch", "origin", choreBranch]);
        const revCountResult = await git([
          "rev-list",
          "--count",
          `origin/${parsed.branch}..origin/${choreBranch}`,
        ]);
        const aheadCount = parseInt(revCountResult?.stdout?.trim() || "0", 10);
        choreBranchPushed = aheadCount > 0;
      } catch (fetchErr) {
        core?.warning?.(
          `Failed to inspect remote '${choreBranch}' commits: ${fetchErr.message}`,
        );
        choreBranchPushed = false;
      }
    } else {
      choreBranchPushed = false;
    }
  }

  // If there are no local changes and the remote chore branch is not ahead of the release branch,
  // the release branch is already fully configured and up to date. No configuration PR is needed.
  if (!choreBranchPushed) {
    core?.info?.(
      `Release branch '${parsed.branch}' already has target configuration applied and no chore branch changes are pending. No setup PR needed.`,
    );
    return {
      branch: parsed.branch,
      versionNumber: parsed.versionNumber,
      choreBranch,
      prNumber: null,
      prUrl: null,
      prCreated: false,
    };
  }

  // 5. Open pull request targeting release branch
  const prTitle = `chore(${parsed.branch}): initialize release branch configuration`;
  const prBody = formatPrBody(parsed);

  let prNumber = null;
  let prUrl = null;
  let prCreated = false;

  if (github?.rest?.pulls) {
    const existingPrs = await github.rest.pulls.list({
      owner,
      repo,
      head: `${owner}:${choreBranch}`,
      base: parsed.branch,
      state: "open",
    });

    if (existingPrs?.data && existingPrs.data.length > 0) {
      const existingPr = existingPrs.data[0];
      prNumber = existingPr.number;
      prUrl = existingPr.html_url;
      core?.info?.(
        `Open configuration PR already exists for '${choreBranch}': PR #${prNumber} (${prUrl})`,
      );
    } else {
      core?.info?.(
        `Opening pull request targeting '${parsed.branch}' from '${choreBranch}'...`,
      );
      const createdPr = await github.rest.pulls.create({
        owner,
        repo,
        title: prTitle,
        head: choreBranch,
        base: parsed.branch,
        body: prBody,
      });

      prNumber = createdPr.data.number;
      prUrl = createdPr.data.html_url;
      prCreated = true;
      core?.info?.(`Successfully opened PR #${prNumber}: ${prUrl}`);
    }

    if (prNumber && github?.rest?.issues?.addLabels) {
      try {
        await github.rest.issues.addLabels({
          owner,
          repo,
          issue_number: prNumber,
          labels: ["internal/pr-release-chore"],
        });
        core?.info?.(
          `Ensured 'internal/pr-release-chore' label on PR #${prNumber}`,
        );
      } catch (labelErr) {
        core?.warning?.(
          `Failed to add label 'internal/pr-release-chore': ${labelErr.message}`,
        );
      }
    }
  }

  return {
    branch: parsed.branch,
    versionNumber: parsed.versionNumber,
    choreBranch,
    prNumber,
    prUrl,
    prCreated,
  };
}

export default async function run(args) {
  try {
    return await runSetupReleaseBranch(args);
  } catch (err) {
    args?.core?.setFailed?.(err.message);
    throw err;
  }
}
