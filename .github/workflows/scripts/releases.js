export default async ({ github, context, core, process = globalThis.process, getOctokit }) => {
  const mode = process?.env?.SCRIPT_MODE;
  switch (mode) {
  case 'check-maintainer':
    return await runCheckMaintainer({ github, context, core, process });
  case 'rc-notify':
    return await runRcNotify({ github, context, core, process });
  case 'publish-release':
    return await runPublishRelease({ github, context, core, process });
  case 'tracking-issue':
    return await runTrackingIssue({ github, context, core, process });
  case 'trigger-rc-release':
    return await runTriggerRcRelease({ github, context, core, process, getOctokit });
  default:
    throw new Error(`Unknown release script mode: ${mode}`);
  }
};

/**
 * check-maintainer: Checks if the user triggering the workflow is an authorized maintainer.
 */
export async function runCheckMaintainer({ context, core, process = globalThis.process }) {
  let maintainers = [];
  
  if (process?.env?.MAINTAINERS && process.env.MAINTAINERS !== "undefined") {
    try {
      const parsed = JSON.parse(process.env.MAINTAINERS);
      maintainers = Array.isArray(parsed) ? parsed : (parsed ? [parsed] : []);
    } catch (e) {
      core?.info?.(`problem parsing maintainers, trying again: ${e.message}`);
      maintainers = process.env.MAINTAINERS.split(',').map(m => m.trim()).filter(Boolean);
    }
  }

  const rawActor = context?.actor || '';
  const actor = rawActor.trim().toLowerCase();
  if (!actor) {
    core?.info?.('Checking if actor is an authorized maintainer: false (empty actor)');
    return false;
  }
  const normalizedMaintainers = maintainers
    .map(m => String(m).trim().toLowerCase())
    .filter(Boolean);
  const isMaintainer = Boolean(normalizedMaintainers.includes(actor));
  core?.info?.(`Checking if '${rawActor}' is an authorized maintainer: ${isMaintainer}`);
  
  return isMaintainer;
}

/**
 * rc-notify: Sends notifications about release candidates.
 */
export async function runRcNotify({ github, context, core, process = globalThis.process }) {
  let tagName =
    process?.env?.TAG ||
    process?.env?.TAG_NAME ||
    context?.payload?.release?.tag_name;
  let rawBranch =
    process?.env?.BRANCH ||
    process?.env?.BRANCH_LABEL ||
    context?.payload?.release?.target_commitish ||
    '';
  let branchLabel = rawBranch.trim().replace(/^refs\/heads\//, '').replace(/^origin\//, '').replace(/\/+$/, '');

  if (!tagName || !branchLabel) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed('tagName and branchLabel must be provided via env (TAG/BRANCH) or release payload.');
    }
    return;
  }

  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  if (!tagName.toLowerCase().includes('rc')) {
    core?.info?.(`Tag "${tagName}" does not appear to be an RC. Skipping notification.`);
    return;
  }

  const isValidBranch = /^release\/v\d+$/.test(branchLabel);
  if (!isValidBranch) {
    throw new Error(`Target branch label "${branchLabel}" is invalid. It must start with "release/v" followed by the major version number.`);
  }

  core?.info?.(`RC Detected: ${tagName}`);
  core?.info?.(`Searching for open issues with labels: "${branchLabel}", "internal/backport", and "internal/merged"`);

  let issues;
  try {
    issues = await github.paginate(github.rest.search.issuesAndPullRequests, {
      q: `repo:${owner}/${repo} is:issue is:open label:"${branchLabel}" label:"internal/backport" label:"internal/merged"`
    });
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(`Failed to search issues for RC notification: ${error.message}`);
    }
    return;
  }

  if (issues.length === 0) {
    core?.info?.('No matching issues found. Exiting.');
    return;
  }

  const releaseUrl = `https://github.com/${owner}/${repo}/releases/tag/${tagName}`;
  const commentBody = `New Release Candidate Available for Validation: [${tagName}](${releaseUrl})\n\n`;

  let commentedCount = 0;
  for (const issue of issues) {
    try {
      await github.rest.issues.createComment({
        owner: owner,
        repo: repo,
        issue_number: issue.number,
        body: commentBody
      });
      core?.info?.(`Commented on issue #${issue.number}`);
      commentedCount++;
    } catch (error) {
      if (core && typeof core.setFailed === 'function') {
        core.setFailed(`Failed to comment on issue #${issue.number}: ${error.message}`);
      }
    }
  }
  
  if (commentedCount === issues.length) {
    core?.info?.(`Success! Notified ${commentedCount} issues.`);
  } else {
    core?.info?.(`Notified ${commentedCount} of ${issues.length} issues.`);
  }
}

/**
 * publish-release: Publishes draft GitHub releases.
 */
async function runPublishRelease({ github, context, core, process }) {
  try {
    const version = process?.env?.VERSION;
    if (!version) {
      return core.setFailed('VERSION environment variable is not defined.');
    }

    const releases = await github.paginate(github.rest.repos.listReleases, {
      owner: context.repo.owner,
      repo: context.repo.repo,
    });

    let release = releases.find(r => r.tag_name === version);
    let tag = version;

    if (!release) {
      const fallbackTag = version.startsWith('v') ? version.slice(1) : `v${version}`;
      release = releases.find(r => r.tag_name === fallbackTag);
      if (release) {
        tag = fallbackTag;
      }
    }

    if (!release) {
      return core.setFailed(`Could not find release for tag "${version}"`);
    }

    if (release.draft) {
      core.info(`Publishing release ID ${release.id} for tag ${tag}`);
      await github.rest.repos.updateRelease({
        owner: context.repo.owner,
        repo: context.repo.repo,
        release_id: release.id,
        draft: false
      });
    } else {
      core.info(`Release for tag ${tag} is already published.`);
    }
  } catch (error) {
    core.setFailed(`Failed to publish release: ${error.message}`);
  }
}

/**
 * tracking-issue: Automatically creates tracking and backport issues for open pull requests.
 */
async function runTrackingIssue({ github, core, process }) {
  try {
    const repo = "terraform-provider-rancher2";
    const owner = "rancher";
    let assignees = [];
    if (process?.env?.TERRAFORM_MAINTAINERS && process.env.TERRAFORM_MAINTAINERS !== "undefined") {
      try {
        const parsed = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
        assignees = Array.isArray(parsed) ? parsed : [parsed];
      } catch (err) {
        if (core && typeof core.warning === 'function') {
          core.warning(`Could not parse TERRAFORM_MAINTAINERS: ${err.message}. Defaulting to no assignees.`);
        }
        assignees = process.env.TERRAFORM_MAINTAINERS.split(',').map(m => m.trim()).filter(Boolean);
      }
    }

    let latestReleaseBranch = "";
    const branches = await github.paginate(github.rest.repos.listBranches,{
      owner,
      repo,
    });

    if (branches.length === 0) {
      core.setFailed('No branches found');
      return;
    }

    const releaseBranches = branches
      .map(b => b.name)
      .filter(name => name.startsWith('release/v'))
      .sort((a, b) => {
        const versionA = parseInt(a.replace('release/v', ''), 10);
        const versionB = parseInt(b.replace('release/v', ''), 10);
        return versionB - versionA;
      });

    if (releaseBranches.length > 0) {
      latestReleaseBranch = releaseBranches[0];
      core.info(`Latest release branch detected: ${latestReleaseBranch}`);
    } else {
      core.setFailed('No release branches found');
      return;
    }

    let pulls;
    try {
      pulls = await github.paginate(github.rest.search.issuesAndPullRequests, {
        q: `repo:${owner}/${repo} is:pr state:open base:main -draft:true -label:internal/ignore -label:internal/pr-backport -label:"autorelease: pending" -label:"autorelease: tagged"`
      });
    } catch (error) {
      throw new Error(`Failed to retrieve pull requests for tracking issue: ${error.message}`);
    }

    const errors = [];
    for (const pr of pulls) {
      try {
        let response;
        let newLabels = ['internal/tracking'];
        let releaseName = "";

        const releaseLabels = pr.labels
          .filter(label => label.name.startsWith('release/v'))
          .sort((a, b) => {
            const versionA = parseInt(a.name.replace('release/v', ''), 10);
            const versionB = parseInt(b.name.replace('release/v', ''), 10);
            return versionB - versionA;
          });
        const latestReleaseLabel = (releaseLabels.length > 0) ? releaseLabels[0].name : null;

        if (latestReleaseLabel) {
          newLabels.push(latestReleaseLabel);
          releaseName = latestReleaseLabel;
        } else {
          newLabels.push(latestReleaseBranch);
          releaseName = latestReleaseBranch;
        }

        const existingIssues = await github.paginate(github.rest.search.issuesAndPullRequests, {
          q: `repo:${owner}/${repo} is:issue is:open label:internal/tracking in:body #${pr.number}`
        });

        if (existingIssues.length > 0) {
          core.info(`Tracking issue already exists for PR #${pr.number}. Skipping.`);
          continue;
        }

        response = await github.rest.issues.create({
          owner: owner,
          repo:  repo,
          title: pr.title,
          body:  `This is the tracking issue for #${pr.number} \n\n` +
            `Please add labels indicating the release versions eg. '${releaseName}' \n\n` +
            `Please add comments for user issues which this issue addresses. \n\n` +
            `Description copied from PR: \n${pr.body ?? ''}`,
          labels: newLabels,
          assignees: assignees
        });

        const newIssue = response.data;
        core.info(`Created tracking issue #${newIssue.number}: ${newIssue.html_url}`);

        const parentIssue = newIssue;
        const parentIssueTitle = parentIssue.title;
        const parentIssueNumber = parentIssue.number;
        
        response = await github.rest.issues.create({
          owner: owner,
          repo: repo,
          title: `[${releaseName}] ${parentIssueTitle}`,
          body:  `Backport #${pr.number} to ${releaseName} for #${parentIssueNumber}\n\n` +
            `Please add this issue to the proper milestone.\n` +
            `Copied from PR: \n${pr.body ?? ''}`,
          labels: [releaseName, "internal/backport"],
          assignees: assignees
        });
        const newSubIssue = response.data;
        core.info(`Created backport issue #${newSubIssue.number}: ${newSubIssue.html_url}`);
        const subIssueId = newSubIssue.id;
        
        await github.request('POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues', {
          owner: owner,
          repo: repo,
          issue_number: parentIssueNumber,
          sub_issue_id: subIssueId,
          headers: {
            'X-GitHub-Api-Version': '2022-11-28'
          }
        });
      } catch (error) {
        errors.push(`Failed to process PR [${pr.number}](${pr.html_url}): ${error.message}`);
      }
    }

    if (errors.length > 0) {
      core.setFailed(`Failed to process some pull requests:\n- ${errors.join('\n- ')}`);
    }
  } catch (error) {
    core.setFailed(`Script failed with error: ${error.message}`);
  }
}

/**
 * Computes the next Release Candidate tag for a given release branch.
 */
export async function computeNextRcTag({ github, owner, repo, branch, core }) {
  const branchName = (branch || '').trim().replace(/^refs\/heads\//, '').replace(/^origin\//, '').replace(/\/+$/, '');
  const match = branchName.match(/^release\/v(\d+)$/);
  if (!match) {
    throw new Error(`Branch '${branch}' does not match expected pattern release/v<major>`);
  }
  const major = parseInt(match[1], 10);
  const prefix = `v${major}.`;

  // Fetch all tags
  const tags = await github.paginate(github.rest.repos.listTags, {
    owner,
    repo,
    per_page: 100
  });

  const branchTags = tags
    .map(t => t?.name)
    .filter(name => typeof name === 'string' && name.startsWith(prefix));

  // Try to read explicit version from release-please-config.json ('release-as') or .release-please-manifest.json
  let configVersion = null;
  try {
    const response = await github.rest.repos.getContent({
      owner,
      repo,
      path: 'release-please-config.json',
      ref: branchName
    });
    if (response.data && response.data.content) {
      const content = Buffer.from(response.data.content, 'base64').toString('utf8');
      const config = JSON.parse(content);
      const releaseAs = config['release-as'] || config.packages?.['.']?.['release-as'] || (config.packages && Object.values(config.packages)[0]?.['release-as']);
      if (releaseAs) {
        const releaseAsStr = String(releaseAs).trim();
        configVersion = releaseAsStr.startsWith('v') ? releaseAsStr : `v${releaseAsStr}`;
      }
    }
  } catch (err) {
    if (core && typeof core.info === 'function') {
      core.info(`Could not read release-please-config.json from ${branchName}: ${err.message}.`);
    }
  }

  let manifestVersion = null;
  try {
    const response = await github.rest.repos.getContent({
      owner,
      repo,
      path: '.release-please-manifest.json',
      ref: branchName
    });
    if (response.data && response.data.content) {
      const content = Buffer.from(response.data.content, 'base64').toString('utf8');
      const manifest = JSON.parse(content);
      const rawManifestVal = manifest['.'] || Object.values(manifest)[0];
      if (rawManifestVal) {
        const manifestStr = String(rawManifestVal).trim();
        manifestVersion = manifestStr.startsWith('v') ? manifestStr : `v${manifestStr}`;
      }
    }
  } catch (err) {
    if (core && typeof core.info === 'function') {
      core.info(`Could not read .release-please-manifest.json from ${branchName}: ${err.message}. Relying on tags.`);
    }
  }

  // Parse existing tags for this major
  const parsedTags = [];
  const tagRegex = /^v(\d+)\.(\d+)\.(\d+)(?:-rc\.?(\d+))?$/i;

  for (const tag of branchTags) {
    const m = tag.match(tagRegex);
    if (m) {
      const tagMajor = parseInt(m[1], 10);
      if (tagMajor === major) {
        parsedTags.push({
          raw: tag,
          major: tagMajor,
          minor: parseInt(m[2], 10),
          patch: parseInt(m[3], 10),
          rc: m[4] !== undefined ? parseInt(m[4], 10) : null,
          isRc: m[4] !== undefined
        });
      }
    }
  }

  // Find all full releases
  const fullReleases = parsedTags.filter(t => !t.isRc);
  fullReleases.sort((a, b) => {
    if (a.minor !== b.minor) return b.minor - a.minor;
    return b.patch - a.patch;
  });
  const latestFull = fullReleases[0] || null;

  // Find all RCs and sort descending by minor, patch, then rc number
  const rcReleases = parsedTags.filter(t => t.isRc);
  rcReleases.sort((a, b) => {
    if (a.minor !== b.minor) return b.minor - a.minor;
    if (a.patch !== b.patch) return b.patch - a.patch;
    return (b.rc ?? 0) - (a.rc ?? 0);
  });

  // Determine candidate target base version:
  // Preference given to explicit release-as from config, then manifest version
  let targetBase = null;
  for (const candidate of [configVersion, manifestVersion]) {
    if (candidate) {
      const m = candidate.match(/^v(\d+)\.(\d+)\.(\d+)(?:-.*)?$/);
      if (m && parseInt(m[1], 10) === major) {
        targetBase = {
          major,
          minor: parseInt(m[2], 10),
          patch: parseInt(m[3], 10)
        };
        break;
      }
    }
  }

  const highestRc = rcReleases.find(r =>
    !latestFull || r.minor > latestFull.minor || (r.minor === latestFull.minor && r.patch > latestFull.patch)
  ) || null;

  // If no manifest version or manifest version is <= latestFull, check active RCs or bump patch
  if (!targetBase || (latestFull && (targetBase.minor < latestFull.minor || (targetBase.minor === latestFull.minor && targetBase.patch <= latestFull.patch)))) {
    if (highestRc) {
      targetBase = { major, minor: highestRc.minor, patch: highestRc.patch };
    } else if (latestFull) {
      targetBase = { major, minor: latestFull.minor, patch: latestFull.patch + 1 };
    } else {
      targetBase = { major, minor: 0, patch: 0 };
    }
  } else if (highestRc && (highestRc.minor > targetBase.minor || (highestRc.minor === targetBase.minor && highestRc.patch > targetBase.patch))) {
    // Prevent regressing version if an active RC cycle already exists higher than targetBase
    targetBase = { major, minor: highestRc.minor, patch: highestRc.patch };
  }

  // Find existing RCs for targetBase
  const matchingRcs = rcReleases.filter(r => r.minor === targetBase.minor && r.patch === targetBase.patch);
  let nextRcNum = 1;
  if (matchingRcs.length > 0) {
    const validRcs = matchingRcs
      .map(r => r.rc)
      .filter(rc => typeof rc === 'number' && !isNaN(rc));
    const maxRc = validRcs.length > 0 ? Math.max(...validRcs) : 0;
    nextRcNum = maxRc + 1;
  }

  return `v${targetBase.major}.${targetBase.minor}.${targetBase.patch}-rc.${nextRcNum}`;
}

/**
 * trigger-rc-release: Computes the next RC tag and dispatches rc-release.yml on main.
 */
export async function runTriggerRcRelease({ github, context, core, process = globalThis.process, getOctokit }) {
  try {
    const owner = context?.repo?.owner || "rancher";
    const repo = context?.repo?.repo || "terraform-provider-rancher2";
    const pr = context?.payload?.pull_request;

    const rawBranch = (process?.env?.BRANCH || pr?.base?.ref || '').trim();
    const branch = rawBranch.replace(/^refs\/heads\//, '').replace(/^origin\//, '').replace(/\/+$/, '');
    const sha = (process?.env?.SHA || pr?.merge_commit_sha || '').trim();

    if (!branch) {
      if (core && typeof core.setFailed === 'function') {
        core.setFailed('Target branch must be provided via env (BRANCH) or PR payload.');
      }
      return;
    }

    if (!sha && core && typeof core.info === 'function') {
      core.info(`No commit SHA provided; rc-release.yml will default to HEAD of branch ${branch}.`);
    }

    const mergeToken = process?.env?.GITHUB_MERGE_TOKEN ? process.env.GITHUB_MERGE_TOKEN.trim() : undefined;
    if (!mergeToken && core && typeof core.warning === 'function') {
      core.warning('GITHUB_MERGE_TOKEN not provided. Workflow dispatch triggered by the default GITHUB_TOKEN may not trigger downstream workflow runs.');
    }

    const octokitFactory = (typeof getOctokit === 'function')
      ? getOctokit
      : (typeof github?.getOctokit === 'function' ? github.getOctokit.bind(github) : undefined);
    const dispatchGithub = (mergeToken && typeof octokitFactory === 'function') ? octokitFactory(mergeToken) : github;

    if (core && typeof core.info === 'function') {
      core.info(`Computing next RC tag for branch: ${branch}...`);
    }
    const tag = await computeNextRcTag({ github: dispatchGithub || github, owner, repo, branch, core });
    if (core && typeof core.info === 'function') {
      core.info(`Computed next RC tag: ${tag}`);
    }

    if (core && typeof core.setOutput === 'function') {
      core.setOutput('tag', tag);
    }

    if (core && typeof core.info === 'function') {
      core.info(`Dispatching 'rc-release.yml' on 'main' with branch=${branch}, sha=${sha || '(HEAD)'}, tag=${tag}...`);
    }
    const inputs = {
      branch,
      tag
    };
    if (sha) {
      inputs.sha = sha;
    }

    await dispatchGithub.rest.actions.createWorkflowDispatch({
      owner,
      repo,
      workflow_id: 'rc-release.yml',
      ref: 'main',
      inputs
    });

    if (core && typeof core.info === 'function') {
      core.info(`Successfully dispatched 'rc-release.yml' for ${tag}!`);
    }
    return tag;
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(`Failed to trigger RC release workflow: ${error.message}`);
    } else {
      throw error;
    }
  }
}
