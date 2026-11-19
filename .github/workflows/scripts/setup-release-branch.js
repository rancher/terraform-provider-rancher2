import fs from "node:fs/promises";
import path from "node:path";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

/**
 * Generates standard caller stubs pointing strictly to @main for a target repository.
 */
export function buildWorkflowStubs(
  repository,
  proc = typeof globalThis !== "undefined" ? globalThis.process : undefined,
) {
  const repo = (
    typeof repository === "string"
      ? repository
      : proc?.env?.GITHUB_REPOSITORY || "rancher/terraform-provider-rancher2"
  ).trim();
  return {
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
    uses: ${repo}/.github/workflows/pull_request.yaml@main
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
    uses: ${repo}/.github/workflows/release-please.yml@main
    secrets: inherit
`,
    ".github/workflows/e2e.yml": `name: "E2E Tests"

on:
  pull_request:
    branches:
      - 'release/v*'
    types: [opened, synchronize, reopened]
  merge_group:
    branches:
      - 'release/v*'
  workflow_call:

permissions: {}

jobs:
  e2e:
    permissions:
      actions: read
      contents: read
      id-token: write
      pull-requests: read
    uses: ${repo}/.github/workflows/e2e.yml@main
    secrets: inherit
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
    uses: ${repo}/.github/workflows/sync-release-manifest.yml@main
    with:
      branch: \${{ inputs.branch }}
      version: \${{ inputs.version }}
      release_as: \${{ inputs.release_as }}
    secrets: inherit
`,
  };
}

/**
 * Standard default caller stubs pointing strictly to upstream @main.
 */
export const CALLER_STUBS = buildWorkflowStubs(
  "rancher/terraform-provider-rancher2",
);

/**
 * Allowed stub filenames inside .github/workflows.
 */
export const ALLOWED_STUB_FILENAMES = new Set([
  "pull_request.yaml",
  "release.yml",
  "e2e.yml",
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

  if (!config || typeof config !== "object" || Array.isArray(config)) {
    config = {};
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

  if (
    !config.packages ||
    typeof config.packages !== "object" ||
    Array.isArray(config.packages)
  ) {
    config.packages = {};
  }
  config.packages["."] = config.packages["."] || {
    "release-type": "go",
  };

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
 * Removes any non-stub workflows from the workflows directory.
 */
export async function cleanupNonStubWorkflows(workflowsDir) {
  try {
    const entries = await fs.readdir(workflowsDir, { withFileTypes: true });
    for (const entry of entries) {
      if (
        entry.isFile() &&
        (entry.name.endsWith(".yml") || entry.name.endsWith(".yaml"))
      ) {
        if (!ALLOWED_STUB_FILENAMES.has(entry.name)) {
          await fs.unlink(path.join(workflowsDir, entry.name));
        }
      }
    }
    await fs.rm(path.join(workflowsDir, "scripts"), {
      recursive: true,
      force: true,
    });
  } catch (err) {
    if (err.code !== "ENOENT") {
      throw err;
    }
  }
}

/**
 * Generates pull request description body.
 */
export function formatPrBody({
  branch,
  releaseAs,
  manifestVersion,
  isUpdate = false,
}) {
  const actionText = isUpdate
    ? `updates configuration and caller stubs on release branch \`${branch}\``
    : `initializes configuration and caller stubs for the new release branch \`${branch}\``;

  return `### Release Branch ${isUpdate ? "Update" : "Initialization"}: \`${branch}\`

This automated pull request ${actionText}.

#### Summary of Changes:
- **Release Please Config**: Ensured \`release-please-config.json\` is configured with \`release-as: ${releaseAs}\` and both draft and prerelease are enabled.
- **Release Please Manifest**: Ensured \`.release-please-manifest.json\` root version is properly set to \`${manifestVersion}\`.
- **Workflow Caller Stubs**: Ensured standalone workflows are replaced with thin caller stubs delegating execution strictly to \`@main\`:
  - \`.github/workflows/pull_request.yaml\`
  - \`.github/workflows/release.yml\`
  - \`.github/workflows/sync-release-manifest.yml\`
  - \`.github/workflows/e2e.yml\` (delegates E2E test execution to @main)
- **Release Marker**: Configured \`rancher2/release.md\` with \`${branch}\`.

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
  const repoSlug =
    proc?.env?.GITHUB_REPOSITORY ||
    (owner && repo
      ? `${owner}/${repo}`
      : "rancher/terraform-provider-rancher2");

  const rootDir = workDir || proc?.cwd?.() || process.cwd();
  const git =
    gitExec || ((args, opts) => runGit(args, { cwd: rootDir, ...opts }));

  core?.info?.(
    `Configuring git identity and preparing release branch '${parsed.branch}'...`,
  );
  if (proc?.env?.GITHUB_ACTIONS === "true") {
    await git(["config", "user.name", "github-actions[bot]"]);
    await git([
      "config",
      "user.email",
      "41898282+github-actions[bot]@users.noreply.github.com",
    ]);
  }

  // Ensure origin/main and tags are fetched
  await git([
    "fetch",
    "origin",
    "+refs/heads/main:refs/remotes/origin/main",
    "--tags",
  ]);

  // Check if target release branch already exists on remote origin
  const lsRemoteResult = await git([
    "ls-remote",
    "--heads",
    "origin",
    `refs/heads/${parsed.branch}`,
  ]);
  const branchExistsOnOrigin = Boolean(
    lsRemoteResult?.stdout &&
    lsRemoteResult.stdout
      .trim()
      .split(/\r?\n/)
      .some((line) => line.trim().endsWith(`refs/heads/${parsed.branch}`)),
  );

  let defaultLastReleaseSha = null;
  const envSha = proc?.env?.LAST_RELEASE_SHA?.trim();
  if (envSha && /^[0-9a-f]{40}$/i.test(envSha)) {
    defaultLastReleaseSha = envSha;
  }

  if (!branchExistsOnOrigin) {
    // Pre-calculate lastReleaseSha before pushing the remote branch so failures don't leave an unconfigured branch
    if (!defaultLastReleaseSha) {
      defaultLastReleaseSha = await resolveBranchCreationSha(
        git,
        "main",
        "origin/main",
      );
    }

    if (
      !defaultLastReleaseSha ||
      !/^[0-9a-f]{40}$/i.test(defaultLastReleaseSha)
    ) {
      throw new Error(
        `Failed to resolve required last-release-sha anchor for release branch '${parsed.branch}'.`,
      );
    }

    core?.info?.(
      `Release branch '${parsed.branch}' does not exist on origin. Creating strictly from 'origin/main'...`,
    );
    // Invariant: New release branch must be branched strictly and exclusively from origin/main
    if (proc?.env?.GITHUB_ACTIONS === "true") {
      await git(["push", "origin", `origin/main:refs/heads/${parsed.branch}`]);
    }
  } else {
    core?.info?.(
      `Release branch '${parsed.branch}' already exists on origin. Updating existing release branch configuration...`,
    );
  }

  // Fetch the release branch from origin with explicit refspec to ensure local remote-tracking branch is updated
  if (branchExistsOnOrigin || proc?.env?.GITHUB_ACTIONS === "true") {
    await git([
      "fetch",
      "origin",
      `+refs/heads/${parsed.branch}:refs/remotes/origin/${parsed.branch}`,
    ]);
  }

  // Checkout setup chore branch based on origin/<branch>
  const choreBranch = `chore/setup-release-v${parsed.versionNumber}`;
  const baseRef =
    branchExistsOnOrigin || proc?.env?.GITHUB_ACTIONS === "true"
      ? `origin/${parsed.branch}`
      : "origin/main";
  core?.info?.(
    `Checking out setup branch '${choreBranch}' based on '${baseRef}'...`,
  );
  await git(["checkout", "-B", choreBranch, baseRef]);
  await git(["reset", "--hard", baseRef]);

  // 1. Configure .release-please-manifest.json
  const manifestPath = path.join(rootDir, ".release-please-manifest.json");
  let existingManifest = "";
  try {
    existingManifest = await fs.readFile(manifestPath, "utf8");
  } catch (err) {
    if (err.code !== "ENOENT") throw err;
  }

  let updatedManifest;
  let preservedVersion = null;
  if (branchExistsOnOrigin && existingManifest) {
    try {
      const parsedManifest = JSON.parse(existingManifest);
      if (
        typeof parsedManifest["."] === "string" &&
        parsedManifest["."].trim().length > 0
      ) {
        preservedVersion = parsedManifest["."].trim();
        core?.info?.(
          `Preserving existing manifest version '${preservedVersion}' on '${parsed.branch}'`,
        );
        updatedManifest = existingManifest;
      }
    } catch {
      // ignore
    }
  }

  if (!updatedManifest) {
    updatedManifest = buildReleasePleaseManifest(parsed.manifestVersion);
  }
  await fs.writeFile(manifestPath, updatedManifest, "utf8");

  // 2. Configure release-please-config.json
  const configPath = path.join(rootDir, "release-please-config.json");
  let existingConfig = "";
  try {
    existingConfig = await fs.readFile(configPath, "utf8");
  } catch (err) {
    if (err.code !== "ENOENT") throw err;
  }

  let targetLastReleaseSha = defaultLastReleaseSha;
  let targetReleaseAs = parsed.releaseAs;
  if (existingConfig) {
    try {
      const parsedConfig = JSON.parse(existingConfig);
      if (
        parsedConfig &&
        typeof parsedConfig === "object" &&
        !Array.isArray(parsedConfig)
      ) {
        if (
          !targetLastReleaseSha &&
          typeof parsedConfig["last-release-sha"] === "string" &&
          /^[0-9a-f]{40}$/i.test(parsedConfig["last-release-sha"].trim())
        ) {
          targetLastReleaseSha = parsedConfig["last-release-sha"].trim();
        }
        if (branchExistsOnOrigin) {
          if (parsedConfig["release-as"]) {
            targetReleaseAs = parsedConfig["release-as"];
            core?.info?.(
              `Preserving existing release-as '${targetReleaseAs}' on '${parsed.branch}'`,
            );
          } else if (preservedVersion) {
            targetReleaseAs = preservedVersion.startsWith("v")
              ? preservedVersion
              : `v${preservedVersion}`;
            core?.info?.(
              `Deriving release-as '${targetReleaseAs}' from preserved manifest version '${preservedVersion}' on '${parsed.branch}'`,
            );
          }
        }
      }
    } catch {
      // ignore
    }
  } else if (branchExistsOnOrigin && preservedVersion) {
    targetReleaseAs = preservedVersion.startsWith("v")
      ? preservedVersion
      : `v${preservedVersion}`;
    core?.info?.(
      `Deriving release-as '${targetReleaseAs}' from preserved manifest version '${preservedVersion}' on '${parsed.branch}'`,
    );
  }

  if (!targetLastReleaseSha) {
    targetLastReleaseSha = await resolveBranchCreationSha(
      git,
      parsed.branch,
      "origin/main",
    );
  }

  if (targetLastReleaseSha) {
    core?.info?.(
      `Resolved release anchor SHA '${targetLastReleaseSha}' for release-please-config.json`,
    );
  }

  const updatedConfig = buildReleasePleaseConfig(
    existingConfig,
    targetReleaseAs,
    targetLastReleaseSha,
  );
  await fs.writeFile(configPath, updatedConfig, "utf8");

  // 3. Configure workflows and caller stubs
  const workflowsDir = path.join(rootDir, ".github", "workflows");
  await fs.mkdir(workflowsDir, { recursive: true });
  await cleanupNonStubWorkflows(workflowsDir);

  const stubs = buildWorkflowStubs(repoSlug);
  for (const [relPath, content] of Object.entries(stubs)) {
    const fullPath = path.join(rootDir, relPath);
    await fs.mkdir(path.dirname(fullPath), { recursive: true });
    await fs.writeFile(fullPath, content, "utf8");
  }

  // 4. Configure rancher2/release.md
  const releaseMdPath = path.join(rootDir, "rancher2", "release.md");
  await fs.mkdir(path.dirname(releaseMdPath), { recursive: true });
  await fs.writeFile(releaseMdPath, `${parsed.branch}\n`, "utf8");

  // 5. Git add and commit
  await git([
    "add",
    "-A",
    ".github/workflows",
    "release-please-config.json",
    ".release-please-manifest.json",
    "rancher2/release.md",
  ]);

  const diffResult = await git(["diff", "--cached", "--name-only"]);
  const hasChanges = Boolean(
    diffResult?.stdout && diffResult.stdout.trim().length > 0,
  );

  let choreBranchPushed = hasChanges;
  const prAction = branchExistsOnOrigin ? "update" : "initialize";

  if (hasChanges) {
    core?.info?.(`Committing and pushing changes to '${choreBranch}'...`);
    await git([
      "commit",
      "-m",
      `fix(${parsed.branch}): ${prAction} release branch configuration`,
    ]);
    if (proc?.env?.GITHUB_ACTIONS === "true") {
      await git(["push", "-u", "origin", choreBranch, "--force"]);
      choreBranchPushed = true;
    } else {
      core?.info?.(
        "Local environment detected: skipping push and PR creation.",
      );
      choreBranchPushed = false;
    }
  } else {
    core?.info?.(
      `No file modifications needed; setup files already match target configuration.`,
    );
    choreBranchPushed = false;
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

  // 6. Open pull request targeting release branch
  const prTitle = `fix(${parsed.branch}): ${prAction} release branch configuration`;
  const prBody = formatPrBody({
    ...parsed,
    releaseAs: targetReleaseAs,
    manifestVersion: preservedVersion || parsed.manifestVersion,
    isUpdate: branchExistsOnOrigin,
  });

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

      if (createdPr?.data?.number) {
        prNumber = createdPr.data.number;
        prUrl = createdPr.data.html_url;
        prCreated = true;
        core?.info?.(`Successfully opened PR #${prNumber}: ${prUrl}`);
      } else {
        core?.warning?.(
          `Pull request creation succeeded but response did not contain a valid PR number.`,
        );
      }
    }

    if (prNumber && typeof github?.rest?.issues?.addLabels === "function") {
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
