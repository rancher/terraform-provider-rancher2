import fs from "node:fs";
import path from "node:path";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

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
 * Checks whether the branch name is a maintenance release branch (e.g. release/v15).
 */
export function isReleaseBranch(branch) {
  return typeof branch === "string" && /^release\/v\d+$/.test(branch.trim());
}

/**
 * Parses and validates a semantic version tag (e.g. "v15.1.2" or "15.1.2").
 */
export function parseSemverTag(tag) {
  if (!tag || typeof tag !== "string") {
    throw new Error(`Invalid semver tag: '${tag}'`);
  }
  const cleanTag = tag.trim();
  const match = cleanTag.match(/^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$/);
  if (!match) {
    throw new Error(`Invalid semver tag: '${tag}'`);
  }
  const [, majorStr, minorStr, patchStr, prerelease] = match;
  const major = parseInt(majorStr, 10);
  const minor = parseInt(minorStr, 10);
  const patch = parseInt(patchStr, 10);
  const version = `${major}.${minor}.${patch}${prerelease ? `-${prerelease}` : ""}`;
  return {
    tag: cleanTag.startsWith("v") ? cleanTag : `v${cleanTag}`,
    version,
    major,
    minor,
    patch,
    prerelease: prerelease || null,
  };
}

/**
 * Updates .release-please-manifest.json content with the target version for root ".".
 */
export function updateManifestContent(rawJson, version) {
  let parsed = {};
  if (rawJson && rawJson.trim()) {
    parsed = JSON.parse(rawJson);
  }
  parsed["."] = version;
  return JSON.stringify(parsed, null, 2) + "\n";
}

/**
 * Calculates the next minor version (e.g. "15.1.2" -> { version: "15.2.0", tag: "v15.2.0" }).
 */
export function calculateNextMinorVersion(version) {
  const parsed = parseSemverTag(version);
  const nextMinorVersion = `${parsed.major}.${parsed.minor + 1}.0`;
  return {
    version: nextMinorVersion,
    tag: `v${nextMinorVersion}`,
  };
}

/**
 * Updates release-please-config.json with last-release-sha, release-as, and removes obsolete bootstrap-sha.
 */
export function updateConfigContent(rawJson, { lastReleaseSha, releaseAs }) {
  if (!lastReleaseSha) {
    throw new Error(
      "lastReleaseSha is required to update release-please-config.json",
    );
  }
  let parsed = {};
  if (rawJson && rawJson.trim()) {
    parsed = JSON.parse(rawJson);
  }
  parsed["last-release-sha"] = lastReleaseSha;
  if (releaseAs) {
    parsed["release-as"] = releaseAs;
  }
  if ("bootstrap-sha" in parsed) {
    delete parsed["bootstrap-sha"];
  }
  return JSON.stringify(parsed, null, 2) + "\n";
}

/**
 * Checks if synchronization is required. Returns false if already aligned.
 */
export function isSyncNeeded({
  manifestContent,
  configContent,
  targetVersion,
  releaseSha,
  releaseAs,
}) {
  let manifest = {};
  let config = {};
  try {
    if (manifestContent && manifestContent.trim()) {
      manifest = JSON.parse(manifestContent);
    }
  } catch {
    return true;
  }
  try {
    if (configContent && configContent.trim()) {
      config = JSON.parse(configContent);
    }
  } catch {
    return true;
  }

  const manifestVersion = manifest["."];
  const configReleaseSha = config["last-release-sha"];
  const hasBootstrapSha = "bootstrap-sha" in config;
  const configReleaseAs = config["release-as"];

  if (manifestVersion !== targetVersion) {
    return true;
  }
  if (configReleaseSha !== releaseSha) {
    return true;
  }
  if (hasBootstrapSha) {
    return true;
  }
  if (releaseAs && configReleaseAs !== releaseAs) {
    return true;
  }

  return false;
}

/**
 * Generates markdown body for the pull request.
 */
export function formatPrBody({
  branch,
  version,
  releaseTag,
  releaseSha,
  releaseAs,
  oldManifestVersion,
  oldReleaseSha,
  oldReleaseAs,
  removedBootstrapSha,
}) {
  const changes = [
    `- **Manifest Version**: \`${oldManifestVersion || "unknown"}\` -> \`${version}\``,
    `- **Last Release SHA**: \`${oldReleaseSha || "none"}\` -> \`${releaseSha}\``,
  ];

  if (releaseAs && oldReleaseAs !== releaseAs) {
    changes.push(
      `- **Release As Anchor**: \`${oldReleaseAs || "none"}\` -> \`${releaseAs}\``,
    );
  }

  if (removedBootstrapSha) {
    changes.push("- **Config Anchor**: Removed obsolete `bootstrap-sha`");
  }

  return [
    "## Automated Release Please Manifest & Config Synchronization",
    "",
    `This pull request aligns \`.release-please-manifest.json\` and \`release-please-config.json\` on \`${branch}\` with published release tag \`${releaseTag}\`.`,
    "",
    "### Changes",
    ...changes,
    "",
    "### Purpose",
    "Prevents Release Please from re-adding historical release notes or regressing changelogs following backport PRs.",
  ].join("\n");
}

/**
 * Resolves target release branch from environment or context.
 */
export function resolveTargetBranch({ context, process }) {
  const envBranch =
    process?.env?.TARGET_BRANCH ||
    process?.env?.INPUT_BRANCH ||
    context?.payload?.inputs?.branch ||
    context?.payload?.release?.target_commitish ||
    (context?.ref ? context.ref.replace(/^refs\/heads\//, "") : "");

  return envBranch ? envBranch.trim() : "";
}

/**
 * Resolves release tag from environment, context, or git tags.
 */
export async function resolveReleaseTag({
  context,
  process,
  targetBranch,
  gitExec = runGit,
}) {
  const manualInput =
    process?.env?.INPUT_VERSION || context?.payload?.inputs?.version || "";

  if (manualInput && manualInput.trim()) {
    let clean = manualInput.trim();
    if (!clean.startsWith("v")) {
      clean = `v${clean}`;
    }
    const parsed = parseSemverTag(clean);
    if (parsed.prerelease) {
      throw new Error(
        `Invalid release version '${manualInput}': manifest synchronization only targets full stable releases (e.g. '15.1.2'). Pre-release versions cannot be used as manifest anchors.`,
      );
    }
    return clean;
  }

  let tag =
    process?.env?.RELEASE_TAG ||
    context?.payload?.release?.tag_name ||
    "";

  if (tag) {
    tag = tag.trim();
    if (!tag.startsWith("v")) {
      tag = `v${tag}`;
    }
    return tag;
  }

  if (process?.env?.RELEASE_VERSION) {
    const v = process.env.RELEASE_VERSION.trim();
    return v.startsWith("v") ? v : `v${v}`;
  }

  // Auto-detect latest tag on target branch matching its major version
  if (targetBranch) {
    const match = targetBranch.match(/^release\/v(\d+)$/);
    const majorPrefix = match ? `v${match[1]}.` : "v";
    const candidates = [targetBranch, `origin/${targetBranch}`, "HEAD"];

    for (const ref of candidates) {
      try {
        const tagOutput = await gitExec([
          "tag",
          "--merged",
          ref,
          "--sort=-version:refname",
          `${majorPrefix}*`,
        ]);
        const tags = tagOutput.stdout
          .trim()
          .split("\n")
          .map((t) => t.trim())
          .filter(Boolean);
        const stableTags = tags.filter((t) => /^v?\d+\.\d+\.\d+$/.test(t));
        if (stableTags.length > 0) {
          return stableTags[0];
        }
      } catch {
        // try next candidate ref
      }
    }
  }

  throw new Error(
    `Could not resolve release tag from environment, workflow inputs, release payload, or git history for branch '${targetBranch}'.`,
  );
}

/**
 * Resolves the desired release-as version tag.
 * If explicitly provided via environment or inputs, normalizes it with a leading 'v'.
 * If not provided, computes the next minor version tag based on targetVersion.
 */
export function resolveReleaseAs({ context, process, targetVersion }) {
  const rawInput =
    process?.env?.RELEASE_AS ||
    process?.env?.INPUT_RELEASE_AS ||
    context?.payload?.inputs?.release_as ||
    "";

  if (rawInput && rawInput.trim()) {
    const clean = rawInput.trim();
    const tag = clean.startsWith("v") ? clean : `v${clean}`;
    parseSemverTag(tag);
    return tag;
  }

  if (targetVersion) {
    const { tag: nextMinorTag } = calculateNextMinorVersion(targetVersion);
    return nextMinorTag;
  }

  return "";
}

/**
 * Resolves release commit SHA for the target tag.
 */
export async function resolveReleaseSha({
  releaseTag,
  process,
  context,
  gitExec = runGit,
}) {
  const envSha =
    process?.env?.RELEASE_SHA ||
    process?.env?.INPUT_LAST_RELEASE_SHA ||
    context?.payload?.inputs?.last_release_sha;

  if (envSha && /^[0-9a-f]{40}$/i.test(envSha.trim())) {
    return envSha.trim();
  }

  try {
    const revOutput = await gitExec(["rev-parse", `${releaseTag}^{commit}`]);
    const sha = revOutput.stdout.trim();
    if (/^[0-9a-f]{40}$/i.test(sha)) {
      return sha;
    }
  } catch {
    // try rev-list fallback
    try {
      const revListOutput = await gitExec(["rev-list", "-n", "1", releaseTag]);
      const sha = revListOutput.stdout.trim();
      if (/^[0-9a-f]{40}$/i.test(sha)) {
        return sha;
      }
    } catch {
      throw new Error(
        `Release tag '${releaseTag}' does not exist in the repository or cannot be resolved to a commit SHA. ` +
        `The sync workflow aligns the manifest with an already published release tag (e.g. 'v15.1.2'). ` +
        `Leave the version input empty to auto-detect the latest published release tag.`,
      );
    }
  }

  throw new Error(`Could not resolve commit SHA for tag '${releaseTag}'.`);
}

/**
 * Main execution function for synchronizing release-please manifest and config.
 */
export async function runSyncManifest({
  github,
  context,
  core,
  process = globalThis.process,
  gitExec = runGit,
  fsModule = fs,
  workspaceDir = process?.cwd?.() || ".",
}) {
  const targetBranch = resolveTargetBranch({ context, process });
  if (!targetBranch) {
    throw new Error("Target branch is required but could not be resolved.");
  }

  if (!isReleaseBranch(targetBranch)) {
    core?.info?.(
      `Branch '${targetBranch}' is not a maintenance release branch (expected 'release/v*'). Skipping manifest synchronization.`,
    );
    return {
      skipped: true,
      reason: "not-release-branch",
      branch: targetBranch,
    };
  }

  core?.info?.(`Target maintenance branch: ${targetBranch}`);

  // Configure git authentication if merge token is available
  const mergeToken =
    process?.env?.GITHUB_MERGE_TOKEN ||
    process?.env?.GITHUB_TOKEN ||
    process?.env?.GH_TOKEN;

  if (mergeToken) {
    try {
      core?.setSecret?.(mergeToken);
      const b64Token = Buffer.from(`x-access-token:${mergeToken}`).toString(
        "base64",
      );
      core?.setSecret?.(b64Token);
      await gitExec([
        "config",
        "--local",
        "--replace-all",
        "http.https://github.com/.extraheader",
        `AUTHORIZATION: basic ${b64Token}`,
      ]);
    } catch (e) {
      core?.warning?.(
        `Failed to set git extraheader credentials: ${e.message}`,
      );
    }
  }

  // Ensure target branch is up-to-date locally before resolving tags and sha
  try {
    await gitExec(["fetch", "origin", `${targetBranch}:${targetBranch}`]);
  } catch {
    // If branch already checked out or fetch refspec fails, fetch origin targetBranch
    await gitExec(["fetch", "origin", targetBranch]);
  }
  await gitExec(["checkout", targetBranch]);

  const releaseTag = await resolveReleaseTag({
    context,
    process,
    targetBranch,
    gitExec,
  });
  const { version: targetVersion, prerelease } = parseSemverTag(releaseTag);

  if (prerelease) {
    core?.info?.(
      `Release tag '${releaseTag}' is a pre-release (${prerelease}). Release Please manifests only track full stable releases. Skipping manifest synchronization.`,
    );
    return {
      skipped: true,
      reason: "prerelease-tag",
      tag: releaseTag,
      branch: targetBranch,
    };
  }

  const releaseSha = await resolveReleaseSha({
    releaseTag,
    process,
    context,
    gitExec,
  });
  const releaseAs = resolveReleaseAs({
    context,
    process,
    targetVersion,
  });

  core?.info?.(`Release tag: ${releaseTag} (version: ${targetVersion})`);
  core?.info?.(`Release commit SHA: ${releaseSha}`);
  if (releaseAs) {
    core?.info?.(`Target release-as: ${releaseAs}`);
  }

  const manifestPath = path.join(workspaceDir, ".release-please-manifest.json");
  const configPath = path.join(workspaceDir, "release-please-config.json");

  let manifestContent = "";
  let configContent = "";

  if (fsModule.existsSync(manifestPath)) {
    manifestContent = fsModule.readFileSync(manifestPath, "utf8");
  } else {
    manifestContent = JSON.stringify({ ".": targetVersion }, null, 2) + "\n";
  }

  if (fsModule.existsSync(configPath)) {
    configContent = fsModule.readFileSync(configPath, "utf8");
  } else {
    throw new Error(`Configuration file '${configPath}' does not exist.`);
  }

  let oldManifest = {};
  let oldConfig = {};
  try {
    oldManifest = JSON.parse(manifestContent);
  } catch {
    // ignore
  }
  try {
    oldConfig = JSON.parse(configContent);
  } catch {
    // ignore
  }

  const oldManifestVersion = oldManifest["."];
  const oldReleaseSha = oldConfig["last-release-sha"];
  const oldReleaseAs = oldConfig["release-as"];
  const removedBootstrapSha = "bootstrap-sha" in oldConfig;

  const syncNeeded = isSyncNeeded({
    manifestContent,
    configContent,
    targetVersion,
    releaseSha,
    releaseAs,
  });

  if (!syncNeeded) {
    core?.info?.(
      `Branch '${targetBranch}' is already synchronized (manifest version ${targetVersion}, last-release-sha ${releaseSha}, release-as ${releaseAs}). No PR needed.`,
    );
    return {
      synced: true,
      prCreated: false,
      branch: targetBranch,
      version: targetVersion,
      sha: releaseSha,
      releaseAs,
    };
  }

  const newManifestContent = updateManifestContent(
    manifestContent,
    targetVersion,
  );
  const newConfigContent = updateConfigContent(configContent, {
    lastReleaseSha: releaseSha,
    releaseAs,
  });

  const cleanBranch = targetBranch.replace(/\//g, "-");
  const syncBranch = `chore/sync-release-please-${cleanBranch}`;

  const owner =
    context?.repo?.owner ||
    process?.env?.GITHUB_REPOSITORY?.split("/")[0] ||
    "rancher";
  const repo =
    context?.repo?.repo ||
    process?.env?.GITHUB_REPOSITORY?.split("/")[1] ||
    "terraform-provider-rancher2";

  // Check if an open PR already exists for this head and base
  let existingPr = null;
  if (github?.rest?.pulls?.list) {
    try {
      const openPrs = await github.rest.pulls.list({
        owner,
        repo,
        base: targetBranch,
        state: "open",
      });
      existingPr = openPrs.data?.find(
        (p) =>
          p.head?.ref === syncBranch ||
          p.head?.ref?.endsWith(`:${syncBranch}`) ||
          p.title?.includes(
            `chore(${targetBranch}): update release-please manifest`,
          ),
      );
    } catch (e) {
      core?.warning?.(`Failed to list existing pull requests: ${e.message}`);
    }
  }

  // Write updated files to disk
  fsModule.writeFileSync(manifestPath, newManifestContent, "utf8");
  fsModule.writeFileSync(configPath, newConfigContent, "utf8");

  // Configure author and commit
  await gitExec(["config", "user.name", "github-actions[bot]"]);
  await gitExec([
    "config",
    "user.email",
    "41898282+github-actions[bot]@users.noreply.github.com",
  ]);

  await gitExec(["checkout", "-B", syncBranch]);
  await gitExec([
    "add",
    ".release-please-manifest.json",
    "release-please-config.json",
  ]);

  const commitMsg = `chore(${targetBranch}): align release-please manifest with tag ${releaseTag}`;
  await gitExec(["commit", "-m", commitMsg]);
  await gitExec(["push", "-f", "origin", syncBranch]);

  core?.info?.(`Pushed synchronization commit to origin/${syncBranch}`);

  if (existingPr) {
    core?.info?.(
      `Open PR #${existingPr.number} already exists for ${syncBranch}. Branch updated with new commit.`,
    );
    return {
      synced: false,
      prCreated: false,
      prNumber: existingPr.number,
      prUrl: existingPr.html_url,
      branch: targetBranch,
    };
  }

  if (!github?.rest?.pulls?.create) {
    core?.info?.("GitHub client not provided. Skipping pull request creation.");
    return {
      synced: false,
      prCreated: false,
      branch: targetBranch,
      syncBranch,
    };
  }

  const prTitle = `chore(${targetBranch}): update release-please manifest to ${targetVersion}`;
  const prBody = formatPrBody({
    branch: targetBranch,
    version: targetVersion,
    releaseTag,
    releaseSha,
    releaseAs,
    oldManifestVersion,
    oldReleaseSha,
    oldReleaseAs,
    removedBootstrapSha,
  });

  core?.info?.(`Creating pull request targeting ${targetBranch}...`);
  const prResponse = await github.rest.pulls.create({
    owner,
    repo,
    title: prTitle,
    head: syncBranch,
    base: targetBranch,
    body: prBody,
  });

  const prNumber = prResponse.data.number;
  const prUrl = prResponse.data.html_url;
  core?.info?.(`Created PR #${prNumber}: ${prUrl}`);

  if (github?.rest?.issues?.addLabels) {
    try {
      await github.rest.issues.addLabels({
        owner,
        repo,
        issue_number: prNumber,
        labels: ["internal/pr-release-chore"],
      });
      core?.info?.(
        `Added 'internal/pr-release-chore' label to PR #${prNumber}`,
      );
    } catch (e) {
      core?.warning?.(`Failed to add labels to PR #${prNumber}: ${e.message}`);
    }
  }

  return {
    synced: false,
    prCreated: true,
    prNumber,
    prUrl,
    branch: targetBranch,
    version: targetVersion,
  };
}

export default async function run(args) {
  return await runSyncManifest(args);
}
