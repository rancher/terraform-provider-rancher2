export default async ({ github, core, context, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
  // https://github.com/actions/toolkit/tree/main/packages/core
  // https://docs.github.com/en/actions/reference/workflows-and-actions/contexts

  const repo = "terraform-provider-rancher2";
  const owner = "rancher";
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  let response; // used to hold all github responses
  let pulls;
  
  try {
    pulls = await github.paginate(github.rest.search.issuesAndPullRequests, {
      q: `repo:${owner}/${repo} is:pr state:open base:main -draft:true -label:internal/user -label:internal/tracking -label:internal/pr-tracked -label:internal/pr-backported -label:internal/backport -label:"autorelease: pending" -label:"autorelease: tagged"`
    });
  } catch (error) {
    // setFailed exits
    core.setFailed(`Failed to retrieve PRs: ${error.message}`);
  }

  for (const pr of pulls) {
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
        repo:  repo,
        title: pr.title,
        body:  `This is the tracking issue for #${pr.number} \n\n` +
          `Please add labels indicating the release versions eg. 'release/v13' \n\n` +
          `Please add comments for user issues which this issue addresses. \n\n` +
          `Description copied from PR: \n${pr.body}`,
        labels: newLabels,
        assignees: assignees
      });
    } catch (error) {
      core.setFailed(`Failed to create tracking issue: ${error.message}`);
    }

    const newIssue = response.data;
    core.info(`New tracking issue data: ${JSON.stringify(newIssue)}`);
    if (releaseLabel) {
      // if release label detected, then add appropriate sub-issues
      const parentIssue = newIssue;
      const parentIssueTitle = parentIssue.title;
      const parentIssueNumber = parentIssue.number;
      // Note: can't get terraform-maintainers team, the default token can't access org level objects
      // Create the sub-issue
      try {
        response = await github.rest.issues.create({
          owner: owner,
          repo: repo,
          title: `[${releaseLabel.name}] ${parentIssueTitle}`,
          body:  `Backport #${pr.number} to ${releaseLabel.name} for #${parentIssueNumber}\n\n` +
            `Copied from PR: \n${pr.body}`,
          labels: [releaseLabel],
          assignees: assignees
        });
      } catch (error) {
        core.setFailed(`Failed to create backport issue: ${error.message}`);
      }
      const newSubIssue = response.data;
      core.info(`New backport issue data: ${JSON.stringify(newSubIssue)}`);
      const subIssueId = newSubIssue.id;
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
        core.setFailed(`Failed to link backport issue to tracking issue: ${error.message}`);
      }
    }

    try {
      await github.rest.issues.addLabels({
        owner: owner,
        repo: repo,
        issue_number: pr.number,
        labels: ["internal/pr-tracked"]
      });
    } catch (error) {
      core.setFailed(`Failed to add label to PR #${pr.number}: ${error.message}`);
    }
  }
};
