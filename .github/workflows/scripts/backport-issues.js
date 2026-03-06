async ({ github, context, core, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples

  const owner = "rancher";
  const repo =  "terraform-provider-rancher2";
  const releaseLabel = context.payload.label.name;
  const parentIssue = context.payload.issue;
  const parentIssueTitle = parentIssue.title;
  const parentIssueNumber = parentIssue.number;
  // Note: unable to dynamically retrieve team members and unable to assign a team to an issue
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  const extractedPrNumber = JSON.parse(process.env.PR);
  let response; // used to hold all github responses

  // Retrieve the PR to get its data
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

  // Create the sub-issue
  try {
    response = await github.rest.issues.create({
      owner: owner,
      repo: repo,
      title: `[${releaseLabel}] ${parentIssueTitle}`,
      body:  [
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

  // Attach the sub-issue to the parent, use REST API because there isn't a github-script API yet.
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
};
