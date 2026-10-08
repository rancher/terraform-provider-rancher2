import { validatePrDescription } from './validate-pr-description.js';

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
    return await runTrackingIssue({ github, context, core, process, getOctokit });
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
 * Filters and sorts release branches in descending order (e.g. release/v15 before release/v14).
 *
 * @param {Array<string|object>} items List of branch objects or names.
 * @returns {string[]} Sorted release branch names.
 */
function extractReleaseBranches(items) {
  return (items || [])
    .map(b => (typeof b === 'string' ? b : b?.name))
    .filter(name => typeof name === 'string' && /^release\/v\d+$/.test(name))
    .sort((a, b) => {
      const versionA = parseInt(a.replace('release/v', ''), 10);
      const versionB = parseInt(b.replace('release/v', ''), 10);
      return versionB - versionA;
    });
}

/**
 * tracking-issue: Automatically converts referenced issues into tracking issues for open pull requests.
 */
export async function runTrackingIssue({ github, context, core, process = globalThis.process, getOctokit }) {
  try {
    const repo = context?.repo?.repo || "terraform-provider-rancher2";
    const owner = context?.repo?.owner || "rancher";

    let latestReleaseBranch = "";
    const branches = await github.paginate(github.rest.repos.listBranches, {
      owner,
      repo,
    });

    if (branches.length === 0) {
      core.setFailed('No branches found');
      return;
    }

    const releaseBranches = extractReleaseBranches(branches);

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
        let issueNumbers;
        try {
          issueNumbers = validatePrDescription(pr.body);
        } catch (error) {
          core.info(`PR #${pr.number} does not contain valid issue references (${error.message}). Skipping.`);
          continue;
        }

        const releaseLabels = extractReleaseBranches(pr.labels);
        const targetReleaseBranch = releaseLabels.length > 0 ? releaseLabels[0] : latestReleaseBranch;

        const uniqueIssues = [...new Set(issueNumbers)];

        for (const issueNumStr of uniqueIssues) {
          const issueNumber = parseInt(issueNumStr, 10);
          let issue;
          try {
            ({ data: issue } = await github.rest.issues.get({
              owner,
              repo,
              issue_number: issueNumber,
            }));
          } catch (error) {
            core.warning(`Could not fetch issue #${issueNumber} for PR #${pr.number}: ${error.message}`);
            continue;
          }

          if (issue.pull_request) {
            core.warning(`Issue #${issueNumber} referenced in PR #${pr.number} is a pull request, not an issue. Skipping.`);
            continue;
          }

          if (issue.state !== 'open') {
            core.warning(`Issue #${issueNumber} referenced in PR #${pr.number} is already closed. Skipping.`);
            continue;
          }

          const currentLabels = (issue.labels || []).map(l => (typeof l === 'string' ? l : l?.name));

          if (currentLabels.includes('internal/tracking')) {
            core.info(`Issue #${issueNumber} already has 'internal/tracking' label. Skipping.`);
            continue;
          }

          const initialLabelsToAdd = ['internal/tracking'];
          if (!currentLabels.includes('internal/user')) {
            initialLabelsToAdd.push('internal/user');
          }

          await github.rest.issues.addLabels({
            owner,
            repo,
            issue_number: issueNumber,
            labels: initialLabelsToAdd,
          });
          core.info(`Added labels [${initialLabelsToAdd.join(', ')}] to issue #${issueNumber}`);

          const linkCommentBody = `<!-- tracking-pr: #${pr.number} -->\nThis is the tracking issue for PR #${pr.number}`;
          await github.rest.issues.createComment({
            owner,
            repo,
            issue_number: issueNumber,
            body: linkCommentBody,
          });
          core.info(`Added PR linking comment for PR #${pr.number} to issue #${issueNumber}`);

          const instructionCommentBody = `A label has been added for the latest release branch, if you need this change to go to any other release branches, please add labels for the branches you need this to go to. Please don't skip branches, eg. if you need something added to release/v13 and the latest is release/v15, you must also add release/v14. Once a label is added a sub-issue will be generated to facilitate the backport, these sub-issues must be in place before the PR is merged or automatic backports won't happen.`;
          await github.rest.issues.createComment({
            owner,
            repo,
            issue_number: issueNumber,
            body: instructionCommentBody,
          });
          core.info(`Added release branch instructions comment to issue #${issueNumber}`);

          if (!currentLabels.includes(targetReleaseBranch)) {
            await github.rest.issues.addLabels({
              owner,
              repo,
              issue_number: issueNumber,
              labels: [targetReleaseBranch],
            });
            core.info(`Added release branch label '${targetReleaseBranch}' to issue #${issueNumber}`);
          } else {
            core.info(`Release branch label '${targetReleaseBranch}' already present on issue #${issueNumber}`);
          }

          const mergeToken = process?.env?.GITHUB_MERGE_TOKEN ? process.env.GITHUB_MERGE_TOKEN.trim() : undefined;
          const octokitFactory = (typeof getOctokit === 'function')
            ? getOctokit
            : (typeof github?.getOctokit === 'function' ? github.getOctokit.bind(github) : undefined);
          const dispatchGithub = (mergeToken && typeof octokitFactory === 'function') ? octokitFactory(mergeToken) : github;

          if (typeof dispatchGithub?.rest?.actions?.createWorkflowDispatch === 'function') {
            try {
              await dispatchGithub.rest.actions.createWorkflowDispatch({
                owner,
                repo,
                workflow_id: 'backport-issues.yml',
                ref: 'main',
                inputs: {
                  issue_number: String(issueNumber),
                  release_label: targetReleaseBranch,
                  pr_number: String(pr.number),
                },
              });
              core.info(`Dispatched 'backport-issues.yml' for issue #${issueNumber} with label '${targetReleaseBranch}'`);
            } catch (dispatchError) {
              core.warning(`Failed to dispatch 'backport-issues.yml' for issue #${issueNumber}: ${dispatchError.message}`);
            }
          }
        }
      } catch (error) {
        errors.push(`Failed to process PR [${pr.number}](${pr.html_url || ''}): ${error.message}`);
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
  const tags = typeof github?.paginate === 'function'
    ? await github.paginate(github.rest?.repos?.listTags, {
      owner,
      repo,
      per_page: 100
    })
    : (await github.rest?.repos?.listTags({ owner, repo, per_page: 100 }))?.data || [];

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

    const normalizeBranch = b => (b || '').toString().trim().replace(/^refs\/heads\//, '').replace(/^origin\//, '').replace(/\/+$/, '');
    const rawBranch = (process?.env?.BRANCH || context?.payload?.workflow_run?.head_branch || pr?.base?.ref || '').trim();
    const branch = normalizeBranch(rawBranch);
    const sha = (process?.env?.SHA || context?.payload?.workflow_run?.head_sha || pr?.merge_commit_sha || '').trim();

    if (!branch) {
      if (core && typeof core.setFailed === 'function') {
        core.setFailed('Target branch must be provided via env (BRANCH) or PR payload.');
      }
      return;
    }

    if (branch.startsWith('release-please')) {
      if (core && typeof core.info === 'function') {
        core.info(`Branch ${branch} is a release-please branch; skipping RC release.`);
      }
      return null;
    }

    if (!sha && core && typeof core.info === 'function') {
      core.info(`No commit SHA provided; rc-release.yml will default to HEAD of branch ${branch}.`);
    }

    const isReleasePleasePr = (p) => {
      if (!p) return false;
      const headRef = normalizeBranch(p?.head?.ref || '');
      if (headRef.startsWith('release-please')) return true;
      if (Array.isArray(p?.labels) && p.labels.some(l => (typeof l === 'string' ? l : l?.name)?.startsWith('autorelease'))) return true;
      if (typeof p?.title === 'string' && /^chore(?:\([^)]*\))?:\s*release/i.test(p.title)) return true;
      return false;
    };

    const workflowRunId = (process?.env?.WORKFLOW_RUN_ID || context?.payload?.workflow_run?.id || '').toString().trim();
    if (workflowRunId && (typeof github?.paginate === 'function' || typeof github?.rest?.actions?.listJobsForWorkflowRun === 'function')) {
      const runId = parseInt(workflowRunId, 10);
      if (!isNaN(runId)) {
        try {
          const jobs = typeof github?.paginate === 'function'
            ? await github.paginate(github.rest?.actions?.listJobsForWorkflowRun, {
              owner,
              repo,
              run_id: runId
            })
            : (await github.rest?.actions?.listJobsForWorkflowRun({ owner, repo, run_id: runId }))?.data?.jobs || [];
          const jobList = Array.isArray(jobs) ? jobs : (Array.isArray(jobs?.jobs) ? jobs.jobs : []);
          const fullReleaseJob = jobList.find(j => j?.name === 'Generate Full Release' || j?.name?.toLowerCase().includes('full release') || j?.name?.toLowerCase() === 'publish' || j?.name?.toLowerCase() === 'release');
          if (fullReleaseJob && fullReleaseJob.conclusion !== 'skipped') {
            if (core && typeof core.info === 'function') {
              core.info(`Workflow run #${workflowRunId} attempted a full release (conclusion: ${fullReleaseJob.conclusion}); skipping RC release.`);
            }
            return null;
          }
        } catch (err) {
          if (core && typeof core.warning === 'function') {
            core.warning(`Could not check workflow run jobs for #${workflowRunId}: ${err.message}`);
          }
        }
      }
    }

    let releasePleasePr = isReleasePleasePr(pr) ? pr : null;
    if (!releasePleasePr && Array.isArray(context?.payload?.workflow_run?.pull_requests)) {
      releasePleasePr = context.payload.workflow_run.pull_requests.find(p => isReleasePleasePr(p) && (!p?.base?.ref || !branch || normalizeBranch(p?.base?.ref) === branch))
        || context.payload.workflow_run.pull_requests.find(isReleasePleasePr);
    }
    if (!releasePleasePr && sha && (typeof github?.paginate === 'function' || typeof github?.rest?.repos?.listPullRequestsAssociatedWithCommit === 'function')) {
      try {
        const pulls = typeof github?.paginate === 'function'
          ? await github.paginate(github.rest?.repos?.listPullRequestsAssociatedWithCommit, {
            owner,
            repo,
            commit_sha: sha
          })
          : (await github.rest?.repos?.listPullRequestsAssociatedWithCommit({ owner, repo, commit_sha: sha }))?.data || [];
        releasePleasePr = pulls.find(p => isReleasePleasePr(p) && (!p?.base?.ref || !branch || normalizeBranch(p?.base?.ref) === branch))
          || pulls.find(isReleasePleasePr);
      } catch (err) {
        if (core && typeof core.warning === 'function') {
          core.warning(`Could not check associated pull requests for SHA ${sha}: ${err.message}`);
        }
      }
    }

    if (!releasePleasePr && Array.isArray(context?.payload?.workflow_run?.pull_requests) && github?.rest?.pulls?.get) {
      const candidates = context.payload.workflow_run.pull_requests.filter(p => p && (!branch || !p.base?.ref || normalizeBranch(p.base?.ref) === branch));
      const pool = candidates.length > 0 ? candidates : context.payload.workflow_run.pull_requests;
      for (const candidate of pool) {
        if (candidate?.number) {
          try {
            const { data: fullPr } = await github.rest.pulls.get({
              owner,
              repo,
              pull_number: candidate.number
            });
            if (isReleasePleasePr(fullPr)) {
              releasePleasePr = fullPr;
              break;
            }
          } catch (err) {
            if (core && typeof core.warning === 'function') {
              core.warning(`Could not fetch PR #${candidate.number} details: ${err.message}`);
            }
          }
        }
      }
    }

    if (releasePleasePr) {
      const headRef = releasePleasePr.head?.ref || '';
      const prIdentifier = releasePleasePr.number ? `PR #${releasePleasePr.number}` : 'Triggering PR';
      if (core && typeof core.info === 'function') {
        core.info(`${prIdentifier} is a release-please PR (${headRef || 'autorelease'}); skipping RC release.`);
      }
      return null;
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
