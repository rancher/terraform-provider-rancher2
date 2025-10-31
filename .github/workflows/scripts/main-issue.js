async ({ github, core, context, process }) => {
  const repo = context.repo.repo;
  const owner = context.repo.owner;
  const pr = context.payload.pull_request;
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  let response; // used to hold all github responses

  let newLabels = ['internal/tracking'];
  const releaseLabel = pr.labels.find(label => label.name.startsWith('release/v'));
  if (releaseLabel) {
    newLabels.push(releaseLabel);
  }

  // Create the tracking issue
  // https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28#create-an-issue
  // Note: issues can't have teams assigned to them and our default token can't retrieve org level team members
  try {
    response = await github.rest.issues.create({
      owner: owner,
      repo: repo,
      title: pr.title,
      body:  `This is the main issue tracking #${pr.number} \n\n` +
        `Please add labels indicating the release versions eg. 'release/v13' \n\n` +
        `Please add comments for user issues which this issue addresses. \n\n` +
        `Description copied from PR: \n${pr.body}`,
      labels: newLabels,
      assignees: assignees
    });
  } catch (error) {
    core.setFailed(`Failed to create main issue: ${error.message}`);
  }
  const newIssue = response.data;
  if (releaseLabel) {
    // if release label detected, then add appropriate sub-issues
    const parentIssue = newIssue.data;
    const parentIssueTitle = parentIssue.title;
    const parentIssueNumber = parentIssue.number;
    // Note: can't get terraform-maintainers team, the default token can't access org level objects
    // Create the sub-issue
    try {
      response = await github.rest.issues.create({
        owner: owner,
        repo: repo,
        title: `[Backport][${releaseLabel.name}] ${parentIssueTitle}`,
        body:  `Backport #${pr.number} to ${releaseLabel.name} for #${parentIssueNumber}\n\n` +
          `Copied from PR: \n${pr.body}`,
        labels: [releaseLabel],
        assignees: assignees
      });
    } catch (error) {
      core.setFailed(`Failed to create backport issue: ${error.message}`);
    }
    const newSubIssue = response.data;
    const subIssueId = newSubIssue.data.id;
    // Attach the sub-issue to the parent using API request
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
      core.setFailed(`Failed to link backport issue to main issue: ${error.message}`);
    }
  }
};
