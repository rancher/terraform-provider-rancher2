export default async ({ github, context, core, process }) => {
  try {
    const prNumberRaw = process.env.PR_NUMBER;
    const prNumber = Number.parseInt(prNumberRaw, 10);
    if (!Number.isFinite(prNumber)) {
      core.setFailed(`Invalid PR_NUMBER: ${prNumberRaw}`);
      return;
    }
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
