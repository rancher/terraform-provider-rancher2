import { execSync } from 'child_process';

export default async ({ github, context, core, process }) => {
  const mode = process.env.SCRIPT_MODE;
  switch (mode) {
  case 'wait-for-settle':
    return await runWaitForSettle({ github, context, core, process });
  case 'backport-pr':
    return await runBackportPr({ github, context, core, process });
  case 'backport-issues':
    return await runBackportIssues({ github, context, core, process });
  case 'merge-label':
    return await runMergeLabel({ github, context, core, process });
  default:
    throw new Error(`Unknown backport script mode: ${mode}`);
  }
};

/**
 * wait-for-settle: Allows GitHub's API slightly more time to index merge commits and retrieve the PR list.
 */
async function runWaitForSettle({ github, context, core, process }) {
  const owner = context.repo.owner;
  const repo = context.repo.repo;
  
  // Handle input from either manual dispatch or push
  const mergeCommitSha = process.env.MERGE_COMMIT_SHA || context.payload.head_commit?.id;
  
  if (!mergeCommitSha) {
    core.setFailed("No merge commit SHA found in environment or context payload.");
    return;
  }

  // wait 10 seconds to allow GitHub to index the commit and associated PRs
  await new Promise(resolve => setTimeout(resolve, 10000));

  try {
    await github.paginate(github.rest.repos.listPullRequestsAssociatedWithCommit, {
      owner,
      repo,
      commit_sha: mergeCommitSha
    });
  } catch (error) {
    core.setFailed(`Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`);
  }
  
  // Set output for next steps
  core.setOutput('merge_commit_sha', mergeCommitSha);
}

/**
 * backport-pr: Cherry-picks a commit to the appropriate backport branches and creates PRs.
 */
async function runBackportPr({ github, core, process }) {
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const mergeCommitSha = process.env.MERGE_COMMIT_SHA;
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  let response;

  try {
    response = await github.rest.repos.listPullRequestsAssociatedWithCommit({
      owner,
      repo,
      commit_sha: mergeCommitSha
    });
  } catch (error) {
    throw new Error(`Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`);
  }
  const associatedPrs = response.data;
  if (associatedPrs.length === 0) {
    core.info(`No PRs associated with commit ${mergeCommitSha}. Exiting.`);
    return;
  }

  const pr = associatedPrs.find(p => p.base.ref === 'main' && p.merged_at);
  if (!pr) {
    core.info(`No merged PR found for commit ${mergeCommitSha}.`);
    return;
  }
  core.info(`Found associated PR: #${pr.number}`);

  core.info(`Searching for 'internal/tracking' issue linked to PR #${pr.number}`);
  try {
    response = await github.request('GET /search/issues', {
      q: `repo:${owner}/${repo} is:issue state:open label:"internal/tracking" in:body #${pr.number}`,
      advanced_search: true,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    throw new Error(`Failed to search for internal/tracking issue for PR #${pr.number}: ${error.message}`);
  }
  const searchResults = response.data;
  if (searchResults.total_count === 0) {
    core.info(`No 'internal/tracking' issue found for PR #${pr.number}. Exiting.`);
    return;
  }
  const trackingIssue = searchResults.items[0];
  core.info(`Found tracking issue: #${trackingIssue.number}`);

  core.info(`Fetching sub-issues for tracking issue #${trackingIssue.number}`);
  try {
    response = await github.request('GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues', {
      owner: owner,
      repo: repo,
      issue_number: trackingIssue.number,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    throw new Error(`Failed to fetch sub-issues for tracking issue #${trackingIssue.number}: ${error.message}`);
  }
  const subIssues = response.data;
  core.info(`Sub-issues data: ${JSON.stringify(subIssues)}`);
  if (!Array.isArray(subIssues)) {
    core.warning(`Unexpected sub-issues data format: ${JSON.stringify(subIssues)}`);
    return;
  }
  if (subIssues.length === 0) {
    core.info(`No sub-issues found for issue #${trackingIssue.number}. Exiting.`);
    return;
  }
  core.info(`Found ${subIssues.length} sub-issues.`);

  for (const subIssue of subIssues) {
    const subIssueNumber = subIssue.number;
    core.info(`Processing sub-issue #${subIssueNumber}...`);

    const releaseLabel = subIssue.labels.find(label => label.name.startsWith('release/v'));

    if (!releaseLabel) {
      core.warning(`Sub-issue #${subIssueNumber} has no 'release/v...' label. Skipping.`);
      continue;
    }

    const targetBranch = releaseLabel.name;
    const isValidBranch = /^release\/v\d{1,2}$/.test(targetBranch);

    if (!isValidBranch) {
      throw new Error(`Target branch label "${targetBranch}" is invalid. It must start with "release/v" and end with exactly one or two digits.`);
    }

    core.info(`Processing sub-issue #${subIssueNumber} for target branch: ${targetBranch}`);
    const newBranchName = `backport-${pr.number}-${targetBranch.replace(/\//g, '-')}`;
    try {
      execSync(`git config user.name "github-actions[bot]"`);
      execSync(`git config user.email "github-actions[bot]@users.noreply.github.com"`);
      execSync(`git fetch origin ${targetBranch}`);
      execSync(`git checkout -b ${newBranchName} origin/${targetBranch}`);
      execSync(`git cherry-pick --allow-empty -x ${mergeCommitSha} -X theirs`);
      execSync(`git push origin ${newBranchName}`);
    } catch (error) {
      throw new Error(`Failed to create and push branch ${newBranchName}: ${error.message}`);
    }

    core.info(`Creating pull request for branch ${newBranchName} targeting ${targetBranch}...`);
    try {
      response = await github.rest.pulls.create({
        owner,
        repo,
        title: pr.title,
        head: newBranchName,
        base: targetBranch,
        body: [
          `This pull request cherry-picks the changes from #${pr.number} into ${targetBranch}`,
          `Addresses #${subIssueNumber} for #${trackingIssue.number}`,
          `**WARNING!**: to avoid having to resolve merge conflicts this PR is generated with 'git cherry-pick -X theirs'.`,
          `Please make sure to carefully inspect this PR so that you don't accidentally revert anything!`,
          `Copied from main PR:`,
          `${pr.body}`
        ].join("\n\n")
      });
    } catch (error) {
      throw new Error(`Failed to create pull request for branch ${newBranchName}: ${error.message}`);
    }
    const newPR = response.data;
    core.info(`Created backport PR data: ${JSON.stringify(newPR)}`);
    const prNumber = newPR.number;
    try {
      await github.rest.issues.addAssignees({
        owner,
        repo,
        issue_number: prNumber,
        assignees: assignees
      });
    } catch (error) {
      throw new Error(`Failed to assign PR #${prNumber}: ${error.message}`);
    }
    try {
      await github.rest.issues.addLabels({
        owner,
        repo,
        issue_number: prNumber,
        labels: ["internal/pr-backport", targetBranch]
      });
    } catch (error) {
      throw new Error(`Failed to add backport label to PR #${prNumber}: ${error.message}`);
    }
  }
}

/**
 * backport-issues: Creates sub-issues for tracking backports.
 */
async function runBackportIssues({ github, context, core, process }) {
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const releaseLabel = context.payload.label.name;
  const parentIssue = context.payload.issue;
  const parentIssueTitle = parentIssue.title;
  const parentIssueNumber = parentIssue.number;
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  const extractedPrNumber = JSON.parse(process.env.PR);
  let response;

  try {
    response = await github.rest.issues.get({
      owner: owner,
      repo: repo,
      issue_number: extractedPrNumber
    });
  } catch (error) {
    throw new Error(`Failed to retrieve PR #${extractedPrNumber}: ${error.message}`);
  }
  const pr = response.data;
  core.info(`PR data: ${JSON.stringify(pr)}`);
  const prNumber = pr.number;

  try {
    response = await github.rest.issues.create({
      owner: owner,
      repo: repo,
      title: `[${releaseLabel}] ${parentIssueTitle}`,
      body: [
        `Backport #${prNumber} to ${releaseLabel} for #${parentIssueNumber}`,
        `Please add this issue to the proper milestone.`,
        `Copied from PR:`,
        `${pr.body}`
      ].join("\n\n"),
      labels: [releaseLabel, "internal/backport"],
      assignees: assignees
    });
  } catch (error) {
    throw new Error(`Failed to create backport issue: ${error.message}`);
  }
  const newIssue = response.data;
  core.info(`New backport issue data: ${JSON.stringify(newIssue)}`);
  const subIssueId = newIssue.id;

  try {
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
    throw new Error(`Failed to link backport issue to tracking issue: ${error.message}`);
  }
}

/**
 * merge-label: Adds internal/merged label to issues referenced in a merged PR body.
 */
async function runMergeLabel({ github, context, core }) {
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const pr = context.payload.pull_request;

  const issueRegex = /#(\d+)/g;
  const prBody = pr.body ?? "";
  const matches = prBody.matchAll(issueRegex);
  const issueNumbers = Array.from(matches, m => parseInt(m[1]));

  core.info(`Found issue numbers in PR body: ${issueNumbers}`);

  for (const issueNumber of issueNumbers) {
    try {
      const { data: issueData } = await github.rest.issues.get({
        owner,
        repo,
        issue_number: issueNumber,
      });

      if (!issueData.pull_request && issueData.labels.some(l => l.name === 'internal/backport')) {
        core.info(`Adding 'internal/merged' label to issue #${issueNumber}`);
        await github.rest.issues.addLabels({
          owner,
          repo,
          issue_number: issueNumber,
          labels: ["internal/merged"]
        });
      }
    } catch (error) {
      core.setFailed(`Could not process issue #${issueNumber}: ${error.message}`);
    }
  }
}
