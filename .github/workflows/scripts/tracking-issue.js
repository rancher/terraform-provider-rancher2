export default async ({ github, core, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
  // https://github.com/actions/toolkit/tree/main/packages/core
  // https://docs.github.com/en/actions/reference/workflows-and-actions/contexts

  const repo = "terraform-provider-rancher2";
  const owner = "rancher";
  const assignees = JSON.parse(process.env.TERRAFORM_MAINTAINERS);
  let response; // used to hold all github responses

  let latestReleaseBranch = "";
  try {
    const branches = await github.paginate(github.rest.repos.listBranches,{
      owner,
      repo,
      protected: true,
    });

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
      throw new Error('No release branches found');
    }
  } catch (error) {
    throw new Error(`Failed to find latest release branch: ${error.message}`);
  }

  let pulls;
  try {
    pulls = await github.paginate(github.rest.search.issuesAndPullRequests, {
      q: `repo:${owner}/${repo} is:pr state:open base:main -draft:true -label:internal/pr-tracked -label:internal/pr-backport -label:"autorelease: pending" -label:"autorelease: tagged"`
    });
  } catch (error) {
    throw new Error(`Failed to retrieve PRs: ${error.message}`);
  }

  for (const pr of pulls) {
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

    try {
      const existingIssues = await github.paginate(github.rest.search.issuesAndPullRequests, {
        q: `repo:${owner}/${repo} is:issue is:open label:internal/tracking in:body #${pr.number}`
      });
      if (existingIssues.length > 0) {
        try {
          await github.rest.issues.addLabels({
            owner: owner,
            repo: repo,
            issue_number: pr.number,
            labels: ["internal/pr-tracked"]
          });
        } catch (error) {
          throw new Error(`Failed to add tracking label to PR #${pr.number}: ${error.message}`);
        }
        core.info(`Tracking issue already exists for PR #${pr.number}. Skipping.`);
        continue;
      }
    } catch (error) {
      core.warning(`Failed to check for existing tracking issue for PR #${pr.number}: ${error.message}`);
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
          `Please add labels indicating the release versions eg. '${releaseName}' \n\n` +
          `Please add comments for user issues which this issue addresses. \n\n` +
          `Description copied from PR: \n${pr.body ?? ''}`,
        labels: newLabels,
        assignees: assignees
      });
    } catch (error) {
      throw new Error(`Failed to create tracking issue: ${error.message}`);
    }

    const newIssue = response.data;
    core.info(`New tracking issue data: ${JSON.stringify(newIssue)}`);

    // add appropriate sub-issues for either release label or latest release branch
    const parentIssue = newIssue;
    const parentIssueTitle = parentIssue.title;
    const parentIssueNumber = parentIssue.number;
    // Note: can't get terraform-maintainers team, the default token can't access org level objects
    // Create the sub-issue
    try {
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
    } catch (error) {
      throw new Error(`Failed to create backport issue: ${error.message}`);
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
      throw new Error(`Failed to link backport issue to tracking issue: ${error.message}`);
    }

    try {
      await github.rest.issues.addLabels({
        owner: owner,
        repo: repo,
        issue_number: pr.number,
        labels: ["internal/pr-tracked"]
      });
    } catch (error) {
      throw new Error(`Failed to add label to PR #${pr.number}: ${error.message}`);
    }
  }
};
