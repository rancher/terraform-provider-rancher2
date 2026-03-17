export default async ({ github, core, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples
  // https://github.com/actions/toolkit/tree/main/packages/core
  // https://docs.github.com/en/actions/reference/workflows-and-actions/contexts
  try {
    const repo = "terraform-provider-rancher2";
    const owner = "rancher";
    // Note: unable to dynamically retrieve team members and unable to assign a team to an issue
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
        q: `repo:${owner}/${repo} is:pr state:open base:main -draft:true -label:internal/pr-tracked -label:internal/pr-backport -label:"autorelease: pending" -label:"autorelease: tagged"`
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
          await github.rest.issues.addLabels({
            owner: owner,
            repo: repo,
            issue_number: pr.number,
            labels: ["internal/pr-tracked"]
          });
          core.info(`Tracking issue already exists for PR #${pr.number}. Skipping.`);
          continue;
        }

        // Create the tracking issue
        // https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28#create-an-issue
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

        // add appropriate sub-issues for either release label or latest release branch
        const parentIssue = newIssue;
        const parentIssueTitle = parentIssue.title;
        const parentIssueNumber = parentIssue.number;
        // Create the sub-issue
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
        // Attach the sub-issue to the parent using API request
        await github.request('POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues', {
          owner: owner,
          repo: repo,
          issue_number: parentIssueNumber,
          sub_issue_id: subIssueId,
          headers: {
            'X-GitHub-Api-Version': '2022-11-28'
          }
        });

        await github.rest.issues.addLabels({
          owner: owner,
          repo: repo,
          issue_number: pr.number,
          labels: ["internal/pr-tracked"]
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
};
