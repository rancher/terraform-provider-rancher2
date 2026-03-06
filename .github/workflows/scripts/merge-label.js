export default async ({ github, context, core, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
  // https://github.com/actions/toolkit/tree/main/packages/core
  // https://docs.github.com/en/actions/reference/workflows-and-actions/contexts

  // this script should find the issue associated with this PR and add a label "internal/merged" to that issue
  const owner = "rancher";
  const repo = "terraform-provider-rancher2";
  const pr = context.payload.pull_request;
  const prNumber = pr.number;
  let response; // used to hold all github responses
  let issue;

  // https://docs.github.io/en/rest/search/search?apiVersion=2022-11-28#search-issues-and-pull-requests
  core.info(`Searching for 'internal/backport' issue linked to PR #${prNumber}`);
  try {
    response = await github.rest.search.issuesAndPullRequests({
      q: `repo:${owner}/${repo} is:issue state:open label:"internal/backport" in:body #${prNumber}`,
    });
  } catch (error) {
    core.setFailed(`Failed to search for internal/backport issue for PR #${prNumber}: ${error.message}`);
  }

  const searchResults = response.data;
  if (searchResults.total_count === 0) {
    core.setFailed(`No 'internal/tracking' issue found for PR #${prNumber}.`);
  }

  issue = searchResults.items[0];
  core.info(`Found tracking issue: #${issue.number}. Adding 'internal/merged' label.`);

  try {
    await github.rest.issues.addLabels({
      owner: owner,
      repo: repo,
      issue_number: issue.number,
      labels: ["internal/merged"]
    });
  } catch (error) {
    core.setFailed(`Failed to add label to issue #${issue.number}: ${error.message}`);
  }
};
