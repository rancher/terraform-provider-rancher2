export default async ({ process, github, core }) => {
  const testNames = JSON.parse(process.env.ALL_TEST_JSON);
  const owner = process.env.OWNER;
  const repo = process.env.REPO;
  const runId = process.env.RUN_ID;

  core.info(`Checking for lock files.`);

  const { data: { artifacts } } = await github.rest.actions.listWorkflowRunArtifacts({
    owner: owner,
    repo: repo,
    run_id: runId,
  });

  for (const testName of testNames){
    core.info(`Checking lock for test: ${testName}`);
    const prefix = `lock-${testName}-`;
    const locks = artifacts.filter(a => a.name.startsWith(prefix));
    if (locks.length == 0) {
      core.setFailed(`No lock found for ${testName}, failing.`);
    }
    for (const lock of locks) {
      core.info(`Found lock ${lock.name}`);
    }
  }
};
