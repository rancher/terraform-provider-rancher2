import { execSync } from 'child_process';
export default async ({ github, core, process }) => {
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const mergeCommitSha = process.env.MERGE_COMMIT_SHA;
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  let response; // used to hold all github responses

  // https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28#list-pull-requests-associated-with-a-commit
  try {
    response = await github.rest.repos.listPullRequestsAssociatedWithCommit({
      owner,
      repo,
      commit_sha: mergeCommitSha
    });
  } catch (error) {
    core.setFailed(`Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`);
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

  // https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-issues-and-pull-requests
  core.info(`Searching for 'internal/tracking' issue linked to PR #${pr.number}`);
  try {
    response = await github.request('GET /search/issues', {
      q: `is:issue state:open label:"internal/tracking" repo:${owner}/${repo} in:body #${pr.number}`,
      advanced_search: true,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    core.setFailed(`Failed to search for internal/tracking issue for PR #${pr.number}: ${error.message}`);
  }
  const searchResults = response.data;
  if (searchResults.total_count === 0) {
    core.info(`No 'internal/tracking' issue found for PR #${pr.number}. Exiting.`);
    return;
  }
  const trackingIssue = searchResults.items[0];
  core.info(`Found tracking issue: #${trackingIssue.number}`);

  // https://docs.github.com/en/rest/issues/sub-issues?apiVersion=2022-11-28#add-sub-issue
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
    core.setFailed(`Failed to fetch sub-issues for tracking issue #${trackingIssue.number}: ${error.message}`);
  }
  const subIssues = response.data;
  if (subIssues.length === 0) {
    core.info(`No sub-issues found for issue #${trackingIssue.number}. Exiting.`);
    return;
  }
  core.info(`Found ${subIssues.length} sub-issues.`);

  for (const subIssue of subIssues) {
    const subIssueNumber = subIssue.number;
    core.info(`Processing sub-issue #${subIssueNumber}...`);

    // Find the release label directly on the sub-issue object
    const releaseLabel = subIssue.labels.find(label => label.name.startsWith('release/v'));
    if (!releaseLabel) {
      core.warning(`Sub-issue #${subIssueNumber} has no 'release/v...' label. Skipping.`);
      continue;
    }
    const targetBranch = releaseLabel.name;
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
      core.setFailed(`Failed to create and push branch ${newBranchName}: ${error.message}`);
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
          `Please add the proper milestone to this PR`,
          `Copied from main PR:`,
          `${pr.body}`
        ].join("\n\n")
      });
    } catch (error) {
      core.setFailed(`Failed to create pull request for branch ${newBranchName}: ${error.message}`);
    }
    const newPR = response.data;
    const prNumber = newPR.number;
    try {
      await github.rest.issues.addAssignees({
        owner,
        repo,
        issue_number: prNumber,
        assignees: assignees
      });
    } catch (error) {
      core.setFailed(`Failed to assign PR #${prNumber}: ${error.message}`);
    }
  }
};
