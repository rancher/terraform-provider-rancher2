export default async ({ github, context, core, process }) => {
  const mode = process.env.SCRIPT_MODE;
  switch (mode) {
  case 'check-maintainer':
    return await runCheckMaintainer({ github, context, core, process });
  case 'rc-notify':
    return await runRcNotify({ github, context, core, process });
  case 'publish-release':
    return await runPublishRelease({ github, context, core, process });
  case 'tracking-issue':
    return await runTrackingIssue({ github, context, core, process });
  default:
    throw new Error(`Unknown release script mode: ${mode}`);
  }
};

/**
 * check-maintainer: Checks if the user triggering the workflow is an authorized maintainer.
 */
async function runCheckMaintainer({ context, core, process }) {
  let maintainers = ["matttrach"];
  
  if (process.env.MAINTAINERS && process.env.MAINTAINERS !== "undefined") {
    try {
      maintainers = JSON.parse(process.env.MAINTAINERS);
    } catch (e) {
      core.info(`problem parsing maintainers, trying again: ${e.message}`);
      maintainers = process.env.MAINTAINERS.split(',').map(m => m.trim());
    }
  }

  const isMaintainer = maintainers.includes(context.actor);
  core.info(`Checking if '${context.actor}' is an authorized maintainer: ${isMaintainer}`);
  
  return isMaintainer;
}

/**
 * rc-notify: Sends notifications about release candidates.
 */
async function runRcNotify({ github, context, core, process }) {
  let tagName =
    process.env.TAG ||
    process.env.TAG_NAME ||
    context.payload.release?.tag_name;
  let branchLabel =
    process.env.BRANCH ||
    process.env.BRANCH_LABEL ||
    context.payload.release?.target_commitish;

  if (!tagName || !branchLabel) {
    core.setFailed('tagName and branchLabel must be provided via env (TAG/BRANCH) or release payload.');
    return;
  }

  const owner = "rancher";
  const repo = "terraform-provider-rancher2";

  if (!tagName.toLowerCase().includes('rc')) {
    core.info(`Tag "${tagName}" does not appear to be an RC. Skipping notification.`);
    return;
  }

  const isValidBranch = /^release\/v\d{1,2}$/.test(branchLabel);
  if (!isValidBranch) {
    throw new Error(`Target branch label "${branchLabel}" is invalid. It must start with "release/v" and end with exactly one or two digits.`);
  }

  core.info(`RC Detected: ${tagName}`);
  core.info(`Searching for open issues with labels: "${branchLabel}", "internal/backport", and "internal/merged"`);

  const issues = await github.paginate(github.rest.search.issuesAndPullRequests, {
    q: `repo:${owner}/${repo} is:issue is:open label:${branchLabel} label:internal/backport label:internal/merged`
  });

  if (issues.length === 0) {
    core.info('No matching issues found. Exiting.');
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
      core.info(`Commented on issue #${issue.number}`);
      commentedCount++;
    } catch (error) {
      core.setFailed(`Failed to comment on issue #${issue.number}: ${error.message}`);
    }
  }
  
  core.info(`Success! Notified ${commentedCount} issues.`);
}

/**
 * publish-release: Publishes draft GitHub releases.
 */
async function runPublishRelease({ github, context, core, process }) {
  try {
    const version = process.env.VERSION;
    const tag = version.startsWith('v') ? version : `v${version}`;

    const releases = await github.paginate(github.rest.repos.listReleases, {
      owner: context.repo.owner,
      repo: context.repo.repo,
    });

    const release = releases.find(r => r.tag_name === tag);
    if (!release) {
      return core.setFailed(`Could not find release for tag ${tag}`);
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
    const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);

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
