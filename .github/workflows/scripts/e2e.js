import { execSync } from 'child_process';

export default async ({ github, context, core, process }) => {
  const mode = process.env.SCRIPT_MODE;
  switch (mode) {
  case 'check-ip':
    return await runCheckIp({ github, context, core, process });
  case 'check-lock':
    return await runCheckLock({ github, context, core, process });
  case 'check-run':
    return await runCheckRun({ github, context, core, process });
  case 'clear-runner':
    return await runClearRunner({ github, context, core, process });
  case 'wait-for-e2e':
    return await runWaitForE2E({ github, context, core, process });
  case 'report-e2e-status':
    return await runReportE2EStatus({ github, context, core, process });
  default:
    throw new Error(`Unknown E2E script mode: ${mode}`);
  }
};

/**
 * check-ip: Checks if there are any IP collisions on active runners.
 */
async function runCheckIp({ github, core, process }) {
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
}

/**
 * check-lock: Checks if a test lock has been acquired by a run with a lower index.
 */
async function runCheckLock({ github, core, process }) {
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
}

/**
 * check-run: Validates that lock artifacts exist for all scheduled tests.
 */
async function runCheckRun({ github, core, process }) {
  const testNames = JSON.parse(process.env.ALL_TESTS_JSON);
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
      throw new Error(`No lock found for ${testName}, failing.`);
    }
    for (const lock of locks) {
      core.info(`Found lock ${lock.name}`);
    }
  }
}

/**
 * clear-runner: Cleans up heavy software packages on GitHub runners to free up disk space.
 */
async function runClearRunner({ core }) {
  const pathsToRemove = [
    '/usr/lib/jvm',
    '/usr/share/dotnet',
    '/usr/share/swift',
    '/usr/local/.ghcup',
    '/usr/local/julia*',
    '/usr/local/lib/android',
    '/usr/local/share/chromium',
    '/opt/microsoft',
    '/opt/google',
    '/opt/az',
    '/usr/local/share/powershell',
    '/opt/hostedtoolcache'
  ];

  try {
    const output = execSync(`df -h / --total | grep total | awk '{print $4}'`).toString();
    core.info(`Available disk space before cleanup: ${output}`);
  } catch (error) {
    throw new Error(`Failed running df to see disk space: ${error}`);
  }

  for (const path of pathsToRemove) {
    core.info(`Removing ${path}...`);
    try {
      execSync(`sudo bash -c "rm -rf ${path}"`);
    } catch (error) {
      throw new Error(`Failed to remove ${path}: ${error}`);
    }
  }

  core.info('Pruning Docker...');
  try {
    execSync(`docker system prune -af`);
    execSync(`docker builder prune -af`);
    execSync(`docker image prune -af`);
    execSync(`docker volume prune -af`);
  } catch(error) {
    throw new Error(`Failed pruning Docker: ${error}`);
  }

  core.info('Disk space after cleanup:');
  try {
    const output = execSync(`df -h / --total | grep total | awk '{print $4}'`).toString();
    core.info(`Available disk space after cleanup: ${output}`);
  } catch (error) {
    throw new Error(`Failed running df to see disk space: ${error}`);
  }
}

/**
 * wait-for-e2e: Creates an E2E testing notification comment on a release PR.
 */
async function runWaitForE2E({ github, context, core, process }) {
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
}

/**
 * report-e2e-status: Reports the passing or failing status of E2E tests on a release PR.
 */
async function runReportE2EStatus({ github, context, core, process }) {
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
}
