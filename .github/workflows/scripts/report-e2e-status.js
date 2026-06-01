export default async ({ github, context, core, process }) => {
  try {
    const prNumberRaw = process.env.PR_NUMBER;
    const prNumber = Number.parseInt(prNumberRaw, 10);
    if (!Number.isFinite(prNumber)) {
      core.setFailed(`Invalid PR_NUMBER: ${prNumberRaw}`);
      return;
    }
    const status = process.env.TEST_STATUS;
    const msg = status === 'success' ? 'End to End Tests Passed!' : 'End to End Tests Failed!';
    await github.rest.issues.createComment({
      issue_number: prNumber,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: `${msg} \n ${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`
    });
    core.info(`Successfully reported E2E test ${status} on PR #${prNumber}`);
  } catch (error) {
    core.setFailed(`Failed to create comment on release PR: ${error.message}`);
  }
};
