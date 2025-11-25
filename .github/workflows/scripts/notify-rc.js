async ({ github, process }) => {
  const tagName = process.env.TAG;
  const branchLabel = process.env.BRANCH;
  const owner = github.repo.owner;
  const repo = github.repo.repo;
  
  if (!tagName.toLowerCase().includes('rc')) {
    console.log(`Tag "${tagName}" does not appear to be an RC. Skipping notification.`);
    return; 
  }
  
  console.log(`RC Detected: ${tagName}`);
  console.log(`Searching for open issues with label: "${branchLabel}"`);

  const issues = await github.paginate(github.rest.search.issuesAndPullRequests, {
    q: `repo:${owner}/${repo} is:issue is:open label:${branchLabel} -label:internal/user -label:internal/tracking`
  });

  if (issues.length === 0) {
    console.log('No matching issues found. Exiting.');
    return;
  }

  const releaseUrl = `https://github.com/${owner}/${repo}/releases/tag/${tagName}`;
  const commentBody = `New Release Candidate Available for Validation: [${tagName}](${releaseUrl})\n\n`;

  let commentedCount = 0;
  for (const issue of issues) {
    if (issue.pull_request) continue; // don't inform PRs

    try {
      await github.rest.issues.createComment({
        owner: owner,
        repo: repo,
        issue_number: issue.number,
        body: commentBody
      });
      console.log(`Commented on issue #${issue.number}`);
      commentedCount++;
    } catch (error) {
      console.error(`Failed to comment on issue #${issue.number}: ${error.message}`);
    }
  }
  
  console.log(`Success! Notified ${commentedCount} issues.`);
};
