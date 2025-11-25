async ({ github, context }) => {
  const release = context.payload.release;
  const tagName = release.tag_name;
  
  if (!tagName.toLowerCase().includes('rc')) {
    console.log(`Tag "${tagName}" does not appear to be an RC. Skipping notification.`);
    return; 
  }
  const branchLabel = release.target_commitish;
  
  console.log(`RC Detected: ${tagName}`);
  console.log(`Searching for open issues with label: "${branchLabel}"`);

  const issues = await github.paginate(github.rest.search.issuesAndPullRequests, {
    q: `repo:${context.repo.owner}/${context.repo.repo} is:issue is:open label:${branchLabel} -label:internal/user -label:internal/tracking`
  });

  if (issues.length === 0) {
    console.log('No matching issues found. Exiting.');
    return;
  }

  const releaseUrl = release.html_url;
  const commentBody = `New Release Candidate Available for Validation: [${tagName}](${releaseUrl})\n\n`;

  let commentedCount = 0;
  for (const issue of issues) {
    if (issue.pull_request) continue; // don't inform PRs

    try {
      await github.rest.issues.createComment({
        owner: context.repo.owner,
        repo: context.repo.repo,
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
