export default async ({ github, context, core, process }) => {
  // Context for this script
  // https://github.com/actions/github-script?tab=readme-ov-file#this-action
  // https://octokit.github.io/rest.js/v22/#custom-requests replace octokit with github in the examples

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
};
