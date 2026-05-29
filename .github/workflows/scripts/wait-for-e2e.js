export default async ({ github, context, core, process }) => {
  try {
    const prNumber = parseInt(process.env.PR_NUMBER, 10);
    await github.rest.issues.createComment({
      issue_number: prNumber,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: `Please make sure e2e tests pass before merging this PR! \n ${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`
    });
    core.info(`Successfully commented on PR #${prNumber}`);
  } catch (error) {
    core.setFailed(`Failed to create comment on release PR: ${error.message}`);
  }
};
