export default async ({ process, github, core }) => {
  const testName = process.env.TEST_NAME;
  const index = parseInt(process.env.INDEX);
  const owner = process.env.OWNER;
  const repo = process.env.REPO;
  const runId = process.env.RUN_ID;
    
  core.info(`Checking lock for test: ${testName} (My Index: ${index})`);
  
  const { data: { artifacts } } = await github.rest.actions.listWorkflowRunArtifacts({
    owner: owner,
    repo: repo,
    run_id: runId,
  });
  
  const prefix = `lock-${testName}-`;
  const conflicts = artifacts.filter(a => a.name.startsWith(prefix));
  
  let status = 'clean';

  for (const artifact of conflicts) {
    // format: lock-TestOneBasic-0
    const parts = artifact.name.split('-');
    const otherIndex = parseInt(parts[parts.length - 1]);
    
    core.info(`Found data: ${artifact.name} (Index: ${otherIndex})`);

    if (!isNaN(otherIndex) && otherIndex < index) {
      core.warning(`Index ${otherIndex} beat us to ${testName}.`);
      status = 'locked';
      break; 
    }
  }
  core.info(`Final Status: ${status}`);
  core.setOutput('status', status);
};
