async ({ github, context, core, process }) => {
  const labelName = context.payload.label.name;
  const parentIssue = context.payload.issue;
  const parentIssueTitle = parentIssue.title;
  const parentIssueNumber = parentIssue.number;
  const repo = context.repo.repo;
  const owner = context.repo.owner;
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
    core.setFailed(`Failed to retrieve PR #${extractedPrNumber}: ${error.message}`);
  }
  const pr = response.data;
  core.info(`PR data: ${JSON.stringify(pr)}`);
  const prNumber = pr.number;

  // Note: can't get terraform-maintainers team, the default token can't access org level objects
  // Create the sub-issue
  try {
    response = await github.rest.issues.create({
      owner: owner,
      repo: repo,
      title: `[${labelName}] ${parentIssueTitle}`,
      body:  [
        `Backport #${prNumber} to ${labelName} for #${parentIssueNumber}`,
        `Copied from PR:`,
        `${pr.body}`
      ].join("\n\n"),
      labels: [labelName],
      assignees: assignees
    });
  } catch (error) {
    core.setFailed(`Failed to create backport issue: ${error.message}`);
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
    core.setFailed(`Failed to link backport issue to tracking issue: ${error.message}`);
  }
};
