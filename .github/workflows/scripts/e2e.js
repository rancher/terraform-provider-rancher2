import { execSync } from 'child_process';

export default async ({ github, context, core, process }) => {
  const mode = process.env.SCRIPT_MODE;
  switch (mode) {
  case 'filter-pr':
    return await runFilterPr({ github, context, core, process });
  case 'e2e-status':
    return await runE2EStatus({ github, context, core, process });
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
 * getWorkflowArtifacts: Retrieves artifacts for a run using pagination if available, falling back to direct REST call.
 */
async function getWorkflowArtifacts({ github, owner, repo, runId }) {
  if (typeof github?.paginate === 'function') {
    return await github.paginate(github?.rest?.actions?.listWorkflowRunArtifacts, {
      owner,
      repo,
      run_id: runId,
    });
  }
  const response = await github?.rest?.actions?.listWorkflowRunArtifacts({
    owner,
    repo,
    run_id: runId,
  });
  return response?.data?.artifacts || response?.artifacts || [];
}

/**
 * check-ip: Checks if there are any IP collisions on active runners.
 */
async function runCheckIp({ github, context, core, process }) {
  const ip = process.env.IP;
  const index = parseInt(process.env.INDEX, 10);
  const owner = process.env.OWNER || context?.repo?.owner;
  const repo = process.env.REPO || context?.repo?.repo;
  const runId = process.env.RUN_ID || context?.runId;
    
  core.info(`Checking collisions for IP: ${ip} (My Index: ${index})`);
  
  let artifacts;
  try {
    artifacts = await getWorkflowArtifacts({ github, owner, repo, runId });
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(`Failed to list workflow run artifacts: ${error.message}`);
    }
    return;
  }
  
  const prefix = `ip-${ip}-`;
  const conflicts = artifacts.filter(a => a.name.startsWith(prefix));
  
  let status = 'clean';

  for (const artifact of conflicts) {
    const parts = artifact.name.split('-');
    const otherIndex = parseInt(parts[parts.length - 1], 10);
    
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
async function runCheckLock({ github, context, core, process }) {
  const testName = process.env.TEST_NAME;
  const index = parseInt(process.env.INDEX, 10);
  const owner = process.env.OWNER || context?.repo?.owner;
  const repo = process.env.REPO || context?.repo?.repo;
  const runId = process.env.RUN_ID || context?.runId;
    
  core.info(`Checking lock for test: ${testName} (My Index: ${index})`);
  
  let artifacts;
  try {
    artifacts = await getWorkflowArtifacts({ github, owner, repo, runId });
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(`Failed to list workflow run artifacts: ${error.message}`);
    }
    return;
  }
  
  const prefix = `lock-${testName}-`;
  const conflicts = artifacts.filter(a => a.name.startsWith(prefix));
  
  let status = 'clean';

  for (const artifact of conflicts) {
    const parts = artifact.name.split('-');
    const otherIndex = parseInt(parts[parts.length - 1], 10);
    
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
async function runCheckRun({ github, context, core, process }) {
  const testNames = JSON.parse(process.env.ALL_TESTS_JSON);
  const owner = process.env.OWNER || context?.repo?.owner;
  const repo = process.env.REPO || context?.repo?.repo;
  const runId = process.env.RUN_ID || context?.runId;

  core.info(`Checking for lock files.`);

  let artifacts;
  try {
    artifacts = await getWorkflowArtifacts({ github, owner, repo, runId });
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(`Failed to list workflow run artifacts: ${error.message}`);
    }
    throw error;
  }

  for (const testName of testNames){
    core.info(`Checking lock for test: ${testName}`);
    const prefix = `lock-${testName}-`;
    const locks = artifacts.filter(a => a.name.startsWith(prefix));
    if (locks.length == 0) {
      if (core && typeof core.setFailed === 'function') {
        core.setFailed(`No lock found for ${testName}, failing.`);
      }
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

/**
 * shouldRunE2E: Evaluates if a PR is an authorized, automated PR targeting a release branch.
 */
function shouldRunE2E(pr, core = { info: () => {}, warning: () => {} }) {
  if (!pr) {
    core.warning('No pull request payload found.');
    return false;
  }

  const baseRepo = pr.base?.repo?.full_name;
  const headRepo = pr.head?.repo?.full_name;
  if (!baseRepo || !headRepo) {
    core.warning('Missing repository information in PR payload. Skipping E2E tests.');
    return false;
  }

  const isFork = Boolean(pr.head?.repo?.fork || baseRepo.toLowerCase() !== headRepo.toLowerCase());
  if (isFork) {
    core.info(`PR is from an external fork (${headRepo} -> ${baseRepo}). Skipping E2E tests for security.`);
    return false;
  }

  const rawBaseRef = pr.base?.ref || '';
  const cleanBaseRef = rawBaseRef.replace(/^refs\/heads\//, '');
  if (!/^release\/v\d+.*$/.test(cleanBaseRef)) {
    core.info(`PR target branch "${rawBaseRef}" is not a release branch. Skipping E2E tests.`);
    return false;
  }

  const rawHeadRef = pr.head?.ref || '';
  const cleanHeadRef = rawHeadRef.replace(/^refs\/heads\//, '');
  const userType = pr.user?.type;
  const userLogin = pr.user?.login;
  const normalizedLogin = typeof userLogin === 'string' ? userLogin.toLowerCase() : '';
  const authorizedBots = ['github-actions[bot]', 'rancher-merge-bot[bot]'];
  const isBot = authorizedBots.includes(normalizedLogin);

  const isBackport = isBot && /^backport-\d+-(?:release-)?v\d+.*$/.test(cleanHeadRef);
  const isReleasePlease = isBot && /^release-please--branches--release(?:--|\/)v\d+.*$/.test(cleanHeadRef);

  if (isBackport || isReleasePlease) {
    core.info(`PR #${pr.number} qualifies for E2E tests (Type: ${isBackport ? 'backport' : 'release-please'}, Head: ${cleanHeadRef}, Author: ${userLogin}).`);
    return true;
  }

  core.info(`PR #${pr.number} is not an authorized automated PR (Author: ${userLogin}, Type: ${userType}, Head: ${cleanHeadRef}). Skipping E2E tests.`);
  return false;
}

/**
 * resolveTargetPr: Resolves the target pull request object from payload, PR number, or workflow_run commit fallback.
 */
async function resolveTargetPr({ github, context, core, process }) {
  if (context?.payload?.pull_request) {
    return context.payload.pull_request;
  }

  let prNumber = null;
  if (process.env.PR_NUMBER) {
    const parsed = Number.parseInt(process.env.PR_NUMBER, 10);
    if (Number.isFinite(parsed) && parsed > 0) {
      prNumber = parsed;
    }
  }

  const parentRun = context?.payload?.workflow_run;

  if (!prNumber && parentRun?.pull_requests?.length > 0) {
    const parentPrNum = parentRun.pull_requests[0].number;
    if (Number.isFinite(parentPrNum) && parentPrNum > 0) {
      prNumber = parentPrNum;
    }
  }

  if (prNumber) {
    core.info(`Fetching PR #${prNumber} from API...`);
    try {
      const response = await github.rest.pulls.get({
        owner: context.repo.owner,
        repo: context.repo.repo,
        pull_number: prNumber,
      });
      return response.data;
    } catch (error) {
      core.setFailed(`Failed to fetch PR #${prNumber}: ${error.message}`);
      return null;
    }
  }

  // Fallback for listener runs triggered by fork PRs where workflow_run.pull_requests is empty
  const headSha = parentRun?.head_sha || process.env.WORKFLOW_RUN_SHA;
  if (headSha && typeof github?.rest?.repos?.listPullRequestsAssociatedWithCommit === 'function') {
    core.info(`Attempting fallback PR resolution for commit SHA: ${headSha}...`);
    try {
      let associatedPrs = [];
      if (typeof github.paginate === 'function') {
        associatedPrs = await github.paginate(
          github.rest.repos.listPullRequestsAssociatedWithCommit,
          {
            owner: context.repo.owner,
            repo: context.repo.repo,
            commit_sha: headSha,
          }
        );
      } else {
        const res = await github.rest.repos.listPullRequestsAssociatedWithCommit({
          owner: context.repo.owner,
          repo: context.repo.repo,
          commit_sha: headSha,
        });
        associatedPrs = res?.data || res || [];
      }

      // Match PR by head branch and repo owner if available, prioritizing open PRs
      const matched = associatedPrs.find(p => {
        if (p.state && p.state !== 'open') {
          return false;
        }
        if (parentRun?.head_branch && p.head?.ref !== parentRun.head_branch) {
          return false;
        }
        if (parentRun?.head_repository?.full_name && p.head?.repo?.full_name?.toLowerCase() !== parentRun.head_repository.full_name.toLowerCase()) {
          return false;
        }
        return true;
      }) || associatedPrs.find(p => p.state === 'open') || associatedPrs[0];

      if (matched) {
        core.info(`Resolved PR #${matched.number} via commit lookup.`);
        if (matched.number && typeof github?.rest?.pulls?.get === 'function') {
          try {
            const fullPr = await github.rest.pulls.get({
              owner: context.repo.owner,
              repo: context.repo.repo,
              pull_number: matched.number,
            });
            if (fullPr?.data) {
              return fullPr.data;
            }
          } catch (err) {
            core.info(`Could not fetch full details for PR #${matched.number}, using commit-matched PR object: ${err.message}`);
          }
        }
        return matched;
      }
    } catch (error) {
      core.warning(`Error searching PRs associated with commit ${headSha}: ${error.message}`);
    }
  }

  return null;
}

/**
 * filter-pr: Evaluates PR qualification and sets the 'should_run' step output.
 */
async function runFilterPr({ github, context, core, process }) {
  const pr = await resolveTargetPr({ github, context, core, process });
  const parentRun = context?.payload?.workflow_run;

  if (!pr) {
    const isExternalRun = (parentRun?.head_repository?.fork === true) ||
      (parentRun?.head_repository?.full_name &&
        parentRun.head_repository.full_name.toLowerCase() !== `${context?.repo?.owner}/${context?.repo?.repo}`.toLowerCase());

    if (isExternalRun) {
      core.info('Triggered by an external fork repository without an identifiable open PR. Skipping E2E tests safely.');
      core.setOutput('should_run', 'false');
      core.setOutput('sha', parentRun?.head_sha || process.env.WORKFLOW_RUN_SHA || '');
      core.setOutput('is_fork', 'true');
      return false;
    }

    core.setFailed('Missing pull_request in context payload and PR_NUMBER env var.');
    core.setOutput('should_run', 'false');
    return false;
  }

  const baseRepo = pr.base?.repo?.full_name;
  const headRepo = pr.head?.repo?.full_name;
  const isFork = Boolean(
    pr.head?.repo?.fork ||
    !pr.head?.repo ||
    (baseRepo && headRepo && baseRepo.toLowerCase() !== headRepo.toLowerCase())
  );
  core.setOutput('is_fork', isFork ? 'true' : 'false');

  const workflowRunSha = process.env.WORKFLOW_RUN_SHA || parentRun?.head_sha;
  if (workflowRunSha && pr.head?.sha && workflowRunSha !== pr.head.sha) {
    core.info(`Workflow run commit (${workflowRunSha}) is superseded by current PR head (${pr.head.sha}). Skipping stale run.`);
    core.setOutput('should_run', 'false');
    core.setOutput('sha', '');
    return false;
  }

  const shouldRun = shouldRunE2E(pr, core);
  core.setOutput('should_run', shouldRun ? 'true' : 'false');
  core.setOutput('sha', pr.head?.sha || '');

  return shouldRun;
}

/**
 * evaluateE2EStatus: Evaluates the combined statuses of the workflow jobs for Branch Protection.
 */
function evaluateE2EStatus({ filterResult, setupResult = 'success', shouldRun, testResult, auditResult, cleanupResult }, core = { info: () => {}, setFailed: () => {} }) {
  if (filterResult !== 'success') {
    core.setFailed(`Filter job did not complete successfully (Result: ${filterResult}).`);
    return false;
  }

  if (shouldRun !== 'true') {
    core.info('E2E tests were safely skipped (PR does not qualify). Status check passing.');
    return true;
  }

  if (setupResult !== 'success' && setupResult !== 'skipped') {
    core.setFailed(`Setup job did not complete successfully (Result: ${setupResult}).`);
    return false;
  }

  if (testResult === 'success' && auditResult === 'success' && (cleanupResult === 'success' || cleanupResult === 'skipped')) {
    core.info('All E2E tests and audits passed successfully!');
    return true;
  }

  core.setFailed(`E2E tests failed or incomplete (Test: ${testResult}, Audit: ${auditResult}, Cleanup: ${cleanupResult}).`);
  return false;
}

/**
 * e2e-status: Evaluates workflow outcome and determines status check pass/fail.
 * Posts the final status directly to the PR's commit SHA to enforce branch protection.
 */
async function runE2EStatus({ github, context, core, process }) {
  const filterResult = process.env.FILTER_RESULT;
  const setupResult = process.env.SETUP_RESULT;
  const shouldRun = process.env.SHOULD_RUN;
  const testResult = process.env.TEST_RESULT;
  const auditResult = process.env.AUDIT_RESULT;
  const cleanupResult = process.env.CLEANUP_RESULT;
  const sha = process.env.SHA;
  const isFork = process.env.IS_FORK === 'true';

  const passed = evaluateE2EStatus({
    filterResult,
    setupResult,
    shouldRun,
    testResult,
    auditResult,
    cleanupResult,
  }, core);

  if (sha && github?.rest?.repos?.createCommitStatus) {
    core.info(`Posting commit status to SHA: ${sha}`);
    let description;
    if (!passed) {
      description = filterResult !== 'success'
        ? 'E2E test filter job failed'
        : 'E2E tests failed or incomplete';
    } else if (shouldRun !== 'true') {
      description = 'E2E tests safely skipped (PR does not qualify)';
    } else {
      description = 'All E2E tests and audits passed successfully';
    }

    try {
      await github.rest.repos.createCommitStatus({
        owner: context.repo.owner,
        repo: context.repo.repo,
        sha: sha,
        state: passed ? 'success' : 'failure',
        context: 'E2E Tests / executor',
        description,
        target_url: `https://github.com/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.runId}`,
      });
    } catch (error) {
      if (isFork) {
        core.info(`Unable to post commit status to fork PR commit (${sha}): ${error.message}`);
      } else {
        core.setFailed(`Failed to post commit status to SHA ${sha}: ${error.message}`);
        return false;
      }
    }
  }

  return passed;
}

export {
  getWorkflowArtifacts,
  runCheckIp,
  runCheckLock,
  runCheckRun,
  runClearRunner,
  runWaitForE2E,
  runReportE2EStatus,
  resolveTargetPr,
  runFilterPr,
  shouldRunE2E,
  runE2EStatus,
  evaluateE2EStatus,
};
