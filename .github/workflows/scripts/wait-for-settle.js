/**
 * Allows GitHub's API slightly more time to index merge commits and retrieve the PR list.
 */
export default async ({ github, context, core, process }) => {
  const owner = context.repo.owner;
  const repo = context.repo.repo;
  
  // Handle input from either manual dispatch or push
  const mergeCommitSha = process.env.MERGE_COMMIT_SHA || context.payload.head_commit?.id;
  
  if (!mergeCommitSha) {
    core.setFailed("No merge commit SHA found in environment or context payload.");
    return;
  }

  // wait 10 seconds to allow GitHub to index the commit and associated PRs
  await new Promise(resolve => setTimeout(resolve, 10000));

  try {
    await github.paginate(github.rest.repos.listPullRequestsAssociatedWithCommit, {
      owner,
      repo,
      commit_sha: mergeCommitSha
    });
  } catch (error) {
    core.setFailed(`Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`);
  }
  
  // Set output for next steps
  core.setOutput('merge_commit_sha', mergeCommitSha);
};
