export default async ({ process, github, core }) => {
  const ip = process.env.IP;
  const index = parseInt(process.env.INDEX);
  const owner = process.env.OWNER;
  const repo = process.env.REPO;
  const runId = process.env.RUN_ID;
    
  core.info(`Checking collisions for IP: ${ip} (My Index: ${index})`);
  
  const { data: { artifacts } } = await github.rest.actions.listWorkflowRunArtifacts({
    owner: owner,
    repo: repo,
    run_id: runId,
  });
  
  const prefix = `ip-${ip}-`;
  const conflicts = artifacts.filter(a => a.name.startsWith(prefix));
  
  let status = 'clean';

  for (const artifact of conflicts) {
    // format: ip-1.2.3.4-0
    const parts = artifact.name.split('-');
    const otherIndex = parseInt(parts[parts.length - 1]);
    
    core.info(`Found data: ${artifact.name} (Index: ${otherIndex})`);

    if (!isNaN(otherIndex) && otherIndex < index) {
      core.warning(`Index ${otherIndex} beat us to IP ${ip}.`);
      status = 'collision';
      break; 
    }
  }
  core.info(`Final Status: ${status}`);
  core.setOutput('status', status);
};
