import assert from 'node:assert/strict';
import test from 'node:test';
import defaultHandler, {
  evaluateE2EStatus,
  resolveTargetPr,
  runCheckIp,
  runCheckLock,
  runCheckRun,
  runE2EStatus,
  runFilterPr,
  shouldRunE2E,
} from './e2e.js';

test('shouldRunE2E - qualifies automated backport PR from github-actions[bot]', () => {
  const pr = {
    number: 101,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - qualifies automated backport PR from custom bot user', () => {
  const pr = {
    number: 102,
    base: {
      ref: 'release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-456-release-v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'rancher-merge-bot[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - qualifies automated release-please PR from bot', () => {
  const pr = {
    number: 201,
    base: {
      ref: 'release/v3',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'release-please--branches--release--v3',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - rejects PR from an external fork', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 301,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'attacker-fork/terraform-provider-rancher2', fork: true },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr, mockCore);
  assert.strictEqual(result, false);
  assert.match(loggedInfo, /external fork/);
});

test('shouldRunE2E - rejects manual PR from human user', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 302,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'developer-jane',
      type: 'User',
    },
  };

  const result = shouldRunE2E(pr, mockCore);
  assert.strictEqual(result, false);
  assert.match(loggedInfo, /not an authorized automated PR/);
});

test('shouldRunE2E - rejects PR targeting main branch', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 303,
    base: {
      ref: 'main',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr, mockCore);
  assert.strictEqual(result, false);
  assert.match(loggedInfo, /not a release branch/);
});

test('shouldRunE2E - rejects bot PR with non-qualifying head branch', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 304,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'renovate/upgrade-deps',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'renovate[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr, mockCore);
  assert.strictEqual(result, false);
  assert.match(loggedInfo, /not an authorized automated PR/);
});

test('shouldRunE2E - returns false on null or empty payload', () => {
  let warned = false;
  const mockCore = {
    info: () => {},
    warning: () => { warned = true; },
  };

  assert.strictEqual(shouldRunE2E(null, mockCore), false);
  assert.strictEqual(warned, true);
});

test('runFilterPr - sets step output should_run to true for qualifying PR', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    payload: {
      pull_request: {
        number: 401,
        base: {
          ref: 'release/v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
        },
        head: {
          ref: 'backport-999-release-v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
          sha: '1234567890abcdef',
        },
        user: {
          login: 'github-actions[bot]',
          type: 'Bot',
        },
      },
    },
  };

  const result = await runFilterPr({ context: mockContext, core: mockCore, process: { env: {} } });
  assert.strictEqual(result, true);
  assert.strictEqual(outputs['should_run'], 'true');
  assert.strictEqual(outputs['sha'], '1234567890abcdef');
});

test('runFilterPr - sets step output should_run to false when PR missing from payload and no PR_NUMBER', async () => {
  let failedMsg = '';
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = { payload: {} };
  const result = await runFilterPr({ context: mockContext, core: mockCore, process: { env: {} } });
  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.match(failedMsg, /Missing pull_request in context payload and PR_NUMBER env var/);
});

test('runFilterPr - fetches PR dynamically via API when payload is missing', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'repo' },
    payload: {}, // Missing PR payload
  };

  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 402);
          return {
            data: {
              number: 402,
              base: {
                ref: 'release/v2',
                repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
              },
              head: {
                ref: 'backport-999-release-v2',
                repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
                sha: 'abcdef1234567890',
              },
              user: {
                login: 'github-actions[bot]',
                type: 'Bot',
              },
            },
          };
        },
      },
    },
  };

  const result = await runFilterPr({ 
    github: mockGithub, 
    context: mockContext, 
    core: mockCore, 
    process: { env: { PR_NUMBER: '402' } } 
  });
  
  assert.strictEqual(result, true);
  assert.strictEqual(outputs['should_run'], 'true');
  assert.strictEqual(outputs['sha'], 'abcdef1234567890');
});

test('evaluateE2EStatus - passes when tests were safely skipped', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    setFailed: () => {},
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'false',
    testResult: 'skipped',
    auditResult: 'skipped',
    cleanupResult: 'skipped',
  }, mockCore);

  assert.strictEqual(result, true);
  assert.match(loggedInfo, /safely skipped/);
});

test('evaluateE2EStatus - passes when all test jobs succeeded', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    setFailed: () => {},
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'true',
    testResult: 'success',
    auditResult: 'success',
    cleanupResult: 'success',
  }, mockCore);

  assert.strictEqual(result, true);
  assert.match(loggedInfo, /passed successfully/);
});

test('evaluateE2EStatus - fails when filter job failed', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'failure',
    shouldRun: 'true',
    testResult: 'success',
    auditResult: 'success',
    cleanupResult: 'success',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /Filter job did not complete successfully/);
});

test('evaluateE2EStatus - fails when test job failed', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'true',
    testResult: 'failure',
    auditResult: 'success',
    cleanupResult: 'success',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /E2E tests failed or incomplete/);
});

test('evaluateE2EStatus - fails when audit job failed', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'true',
    testResult: 'success',
    auditResult: 'failure',
    cleanupResult: 'success',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /E2E tests failed or incomplete/);
});

test('evaluateE2EStatus - fails when cleanup job failed', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'true',
    testResult: 'success',
    auditResult: 'success',
    cleanupResult: 'failure',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /E2E tests failed or incomplete/);
});

test('runE2EStatus - evaluates statuses from process environment', async () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    setFailed: () => {},
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'true',
      TEST_RESULT: 'success',
      AUDIT_RESULT: 'success',
      CLEANUP_RESULT: 'success',
    },
  };

  const result = await runE2EStatus({ core: mockCore, process: mockProcess });
  assert.strictEqual(result, true);
  assert.match(loggedInfo, /passed successfully/);
});

test('defaultHandler - routes filter-pr and e2e-status modes', async () => {
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: () => {},
  };

  const mockContext = {
    payload: {
      pull_request: {
        base: { ref: 'release/v2', repo: { full_name: 'rancher/repo', fork: false } },
        head: { ref: 'backport-1-release-v2', repo: { full_name: 'rancher/repo', fork: false }, sha: 'mockheadsha123' },
        user: { login: 'github-actions[bot]', type: 'Bot' },
      },
    },
  };

  const filterResult = await defaultHandler({
    github: {},
    context: mockContext,
    core: mockCore,
    process: { env: { SCRIPT_MODE: 'filter-pr' } },
  });
  assert.strictEqual(filterResult, true);

  const statusResult = await defaultHandler({
    github: {},
    context: mockContext,
    core: mockCore,
    process: {
      env: {
        SCRIPT_MODE: 'e2e-status',
        FILTER_RESULT: 'success',
        SHOULD_RUN: 'false',
      },
    },
  });
  assert.strictEqual(statusResult, true);

  await assert.rejects(
    async () => {
      await defaultHandler({
        github: {},
        context: mockContext,
        core: mockCore,
        process: { env: { SCRIPT_MODE: 'unknown-mode' } },
      });
    },
    /Unknown E2E script mode: unknown-mode/,
  );
});

test('shouldRunE2E - qualifies release-please PR when branch uses slash notation (release-please--branches--release/v15)', () => {
  const pr = {
    number: 501,
    base: {
      ref: 'release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'release-please--branches--release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - normalizes refs/heads/ prefix on head and base refs', () => {
  const pr = {
    number: 502,
    base: {
      ref: 'refs/heads/release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'refs/heads/backport-123-release-v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - qualifies bot account ending in [bot] even if user.type is not Bot', () => {
  const pr = {
    number: 503,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-101-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'rancher-merge-bot[bot]',
      type: 'User',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - rejects PR when base.repo or head.repo is missing', () => {
  let warned = false;
  const mockCore = {
    info: () => {},
    warning: () => { warned = true; },
  };

  const prMissingHeadRepo = {
    number: 504,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-101-release-v2',
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  assert.strictEqual(shouldRunE2E(prMissingHeadRepo, mockCore), false);
  assert.strictEqual(warned, true);
});

test('shouldRunE2E - case-insensitively recognizes matching repository names', () => {
  const pr = {
    number: 505,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'Rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-101-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('evaluateE2EStatus - fails when cleanup job was cancelled', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    shouldRun: 'true',
    testResult: 'success',
    auditResult: 'success',
    cleanupResult: 'cancelled',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /E2E tests failed or incomplete/);
});

test('runCheckIp - detects IP collision using context fallback and outputs status', async () => {
  let statusOutput = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (key, val) => {
      if (key === 'status') statusOutput = val;
    },
  };

  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: {
            artifacts: [
              { name: 'ip-10.0.0.1-0' },
              { name: 'ip-10.0.0.1-1' },
            ],
          },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 999,
  };

  const mockProcess = {
    env: {
      IP: '10.0.0.1',
      INDEX: '1',
    },
  };

  await runCheckIp({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(statusOutput, 'collision');
});

test('runCheckLock - detects lock using context fallback and outputs status', async () => {
  let statusOutput = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setOutput: (key, val) => {
      if (key === 'status') statusOutput = val;
    },
  };

  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: {
            artifacts: [
              { name: 'lock-TestOneBasic-0' },
            ],
          },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 999,
  };

  const mockProcess = {
    env: {
      TEST_NAME: 'TestOneBasic',
      INDEX: '1',
    },
  };

  await runCheckLock({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(statusOutput, 'locked');
});

test('runCheckRun - verifies lock artifacts exist using context fallback and throws on missing', async () => {
  const mockCore = {
    info: () => {},
  };

  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: {
            artifacts: [
              { name: 'lock-TestOneBasic-0' },
            ],
          },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 999,
  };

  // Succeeded when all scheduled tests have a lock
  await runCheckRun({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: {
      env: {
        ALL_TESTS_JSON: JSON.stringify(['TestOneBasic']),
      },
    },
  });

  // Throws when a scheduled test is missing a lock
  await assert.rejects(
    async () => {
      await runCheckRun({
        github: mockGithub,
        context: mockContext,
        core: mockCore,
        process: {
          env: {
            ALL_TESTS_JSON: JSON.stringify(['TestOneBasic', 'TestMissingBasic']),
          },
        },
      });
    },
    /No lock found for TestMissingBasic, failing\./,
  );
});

test('evaluateE2EStatus - fails when setup job failed', () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    setupResult: 'failure',
    shouldRun: 'true',
    testResult: 'skipped',
    auditResult: 'skipped',
    cleanupResult: 'skipped',
  }, mockCore);

  assert.strictEqual(result, false);
  assert.match(failedMsg, /Setup job did not complete successfully/);
});

test('evaluateE2EStatus - passes when setup job succeeded and shouldRun is false (clean skip)', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    setFailed: () => {},
  };

  const result = evaluateE2EStatus({
    filterResult: 'success',
    setupResult: 'success',
    shouldRun: 'false',
    testResult: 'skipped',
    auditResult: 'skipped',
    cleanupResult: 'skipped',
  }, mockCore);

  assert.strictEqual(result, true);
  assert.match(loggedInfo, /safely skipped/);
});

test('shouldRunE2E - rejects PR when base.ref has no digits after release/v (e.g. release/v or release/various)', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr1 = {
    number: 601,
    base: {
      ref: 'release/v',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr1, mockCore), false);
  assert.match(loggedInfo, /not a release branch/);

  const pr2 = {
    number: 602,
    base: {
      ref: 'release/various',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr2, mockCore), false);
  assert.match(loggedInfo, /not a release branch/);
});

test('shouldRunE2E - rejects PR when base.ref has trailing characters or minor versions (e.g. release/v15-malicious, release/v15.1)', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr1 = {
    number: 604,
    base: {
      ref: 'release/v15-malicious',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr1, mockCore), false);
  assert.match(loggedInfo, /not a release branch/);

  const pr2 = {
    number: 605,
    base: {
      ref: 'release/v15.1',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'release-please--branches--release--v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr2, mockCore), false);
  assert.match(loggedInfo, /not a release branch/);
});

test('shouldRunE2E - rejects bot PR when head.ref has trailing characters or malicious suffix', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr1 = {
    number: 606,
    base: {
      ref: 'release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v15-malicious',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr1, mockCore), false);
  assert.match(loggedInfo, /not an authorized automated PR/);

  const pr2 = {
    number: 607,
    base: {
      ref: 'release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'release-please--branches--release--v15-malicious',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr2, mockCore), false);
  assert.match(loggedInfo, /not an authorized automated PR/);

  const pr3 = {
    number: 608,
    base: {
      ref: 'release/v15',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'release-please--branches--release/v15.1',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };
  assert.strictEqual(shouldRunE2E(pr3, mockCore), false);
  assert.match(loggedInfo, /not an authorized automated PR/);
});

test('shouldRunE2E - rejects bot PR when head.repo.fork is true even if repo full_names match', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 603,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: true },
    },
    user: { login: 'github-actions[bot]', type: 'Bot' },
  };

  assert.strictEqual(shouldRunE2E(pr, mockCore), false);
  assert.match(loggedInfo, /external fork/);
});

test('runCheckIp - paginates when github.paginate is defined', async () => {
  let statusOutput = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      if (key === 'status') statusOutput = val;
    },
  };

  const mockGithub = {
    paginate: async () => [
      { name: 'ip-10.0.0.1-0' },
      { name: 'ip-10.0.0.1-1' },
    ],
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 999,
  };

  const mockProcess = {
    env: {
      IP: '10.0.0.1',
      INDEX: '1',
    },
  };

  await runCheckIp({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(statusOutput, 'collision');
});

test('runCheckLock - paginates when github.paginate is defined', async () => {
  let statusOutput = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      if (key === 'status') statusOutput = val;
    },
  };

  const mockGithub = {
    paginate: async () => [
      { name: 'lock-TestOneBasic-0' },
    ],
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 999,
  };

  const mockProcess = {
    env: {
      TEST_NAME: 'TestOneBasic',
      INDEX: '1',
    },
  };

  await runCheckLock({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(statusOutput, 'locked');
});

test('runCheckIp - sets failed when API call throws an error', async () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
    setOutput: () => {},
  };

  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => { throw new Error('API Rate Limit'); },
      },
    },
  };

  await runCheckIp({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'repo' }, runId: 1 },
    core: mockCore,
    process: { env: { IP: '1.2.3.4', INDEX: '1' } },
  });

  assert.match(failedMsg, /Failed to list workflow run artifacts: API Rate Limit/);
});

test('shouldRunE2E - qualifies automated backport PR with short branch format (backport-2033-v14)', () => {
  const pr = {
    number: 701,
    base: {
      ref: 'release/v14',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-2033-v14',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('shouldRunE2E - qualifies bot account with mixed-case login', () => {
  const pr = {
    number: 702,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'GitHub-Actions[bot]',
      type: 'User',
    },
  };

  const result = shouldRunE2E(pr);
  assert.strictEqual(result, true);
});

test('runCheckRun - sets core.setFailed when lock is missing', async () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: {
            artifacts: [],
          },
        }),
      },
    },
  };

  await assert.rejects(
    async () => {
      await runCheckRun({
        github: mockGithub,
        context: { repo: { owner: 'rancher', repo: 'repo' }, runId: 1 },
        core: mockCore,
        process: {
          env: {
            ALL_TESTS_JSON: JSON.stringify(['TestOneBasic']),
          },
        },
      });
    },
    /No lock found for TestOneBasic, failing\./,
  );

  assert.strictEqual(failedMsg, 'No lock found for TestOneBasic, failing.');
});


test('shouldRunE2E - rejects unrelated bot account', () => {
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
  };

  const pr = {
    number: 801,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    user: {
      login: 'some-malicious-app[bot]',
      type: 'Bot',
    },
  };

  const result = shouldRunE2E(pr, mockCore);
  assert.strictEqual(result, false);
  assert.match(loggedInfo, /not an authorized automated PR/);
});

test('runE2EStatus - posts commit status to API when SHA is provided', async () => {
  let loggedInfo = '';
  let commitStatusCalled = false;
  let commitStatusArgs = {};
  
  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    setFailed: () => {},
  };

  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async (args) => {
          commitStatusCalled = true;
          commitStatusArgs = args;
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 9999,
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'true',
      TEST_RESULT: 'success',
      AUDIT_RESULT: 'success',
      CLEANUP_RESULT: 'success',
      SHA: '1234567890abcdef',
    },
  };

  const result = await runE2EStatus({ github: mockGithub, context: mockContext, core: mockCore, process: mockProcess });
  
  assert.strictEqual(result, true);
  assert.match(loggedInfo, /Posting commit status to SHA: 1234567890abcdef/);
  assert.strictEqual(commitStatusCalled, true);
  assert.strictEqual(commitStatusArgs.sha, '1234567890abcdef');
  assert.strictEqual(commitStatusArgs.state, 'success');
  assert.strictEqual(commitStatusArgs.context, 'E2E Tests / executor');
});

test('defaultHandler - forwards github and context to runE2EStatus and creates status check', async () => {
  let createdStatus = null;
  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async (args) => {
          createdStatus = args;
        },
      },
    },
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 12345,
  };
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
  };

  const result = await defaultHandler({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: {
      env: {
        SCRIPT_MODE: 'e2e-status',
        FILTER_RESULT: 'success',
        SETUP_RESULT: 'success',
        SHOULD_RUN: 'true',
        TEST_RESULT: 'success',
        AUDIT_RESULT: 'success',
        CLEANUP_RESULT: 'success',
        SHA: 'feedbeef123',
      },
    },
  });

  assert.strictEqual(result, true);
  assert.ok(createdStatus);
  assert.strictEqual(createdStatus.sha, 'feedbeef123');
  assert.strictEqual(createdStatus.state, 'success');
});

test('runFilterPr - always outputs sha even when PR does not qualify (should_run=false)', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  // Human PR targeting release branch
  const mockContext = {
    payload: {
      pull_request: {
        number: 888,
        base: {
          ref: 'release/v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
        },
        head: {
          ref: 'feature-human',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
          sha: 'humanheadcommitsha',
        },
        user: {
          login: 'developer-bob',
          type: 'User',
        },
      },
    },
  };

  const result = await runFilterPr({
    github: {},
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], 'humanheadcommitsha');
  assert.strictEqual(outputs['is_fork'], 'false');
});

test('runFilterPr - detects superseded commit when PR_HEAD_SHA does not match pr.head.sha', async () => {
  const outputs = {};
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo += msg; },
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    payload: {
      pull_request: {
        number: 889,
        base: {
          ref: 'release/v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
        },
        head: {
          ref: 'backport-1-release-v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
          sha: 'newestheadsha999',
        },
        user: {
          login: 'github-actions[bot]',
          type: 'Bot',
        },
      },
    },
  };

  const result = await runFilterPr({
    github: {},
    context: mockContext,
    core: mockCore,
    process: {
      env: {
        PR_HEAD_SHA: 'olderstalesha111',
      },
    },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], '');
  assert.match(loggedInfo, /superseded by current PR head/);
});

test('runFilterPr - does not mark run as superseded when PR_HEAD_SHA matches pr.head.sha', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    payload: {
      pull_request: {
        number: 890,
        base: {
          ref: 'release/v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
        },
        head: {
          ref: 'backport-1-release-v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
          sha: 'matchingheadsha888',
        },
        user: {
          login: 'github-actions[bot]',
          type: 'Bot',
        },
      },
    },
  };

  const result = await runFilterPr({
    github: {},
    context: mockContext,
    core: mockCore,
    process: {
      env: {
        PR_HEAD_SHA: 'matchingheadsha888',
      },
    },
  });

  assert.strictEqual(result, true);
  assert.strictEqual(outputs['should_run'], 'true');
  assert.strictEqual(outputs['sha'], 'matchingheadsha888');
});

test('runFilterPr - detects superseded commit when parentRun.pull_requests[0].head.sha does not match pr.head.sha', async () => {
  const outputs = {};
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo += msg; },
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        pull_requests: [
          {
            number: 891,
            head: {
              ref: 'backport-1-release-v2',
              sha: 'staleheadsha333',
            },
          },
        ],
      },
    },
  };

  const mockGithub = {
    rest: {
      pulls: {
        get: async () => ({
          data: {
            number: 891,
            base: {
              ref: 'release/v2',
              repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
            },
            head: {
              ref: 'backport-1-release-v2',
              repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
              sha: 'newestprheadsha444',
            },
            user: {
              login: 'github-actions[bot]',
              type: 'Bot',
            },
          },
        }),
      },
    },
  };

  const result = await runFilterPr({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], '');
  assert.match(loggedInfo, /superseded by current PR head/);
});

test('runFilterPr - safe fork fallback does not mark as superseded when pull_requests is empty and PR_HEAD_SHA is undefined', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    payload: {
      pull_request: {
        number: 892,
        base: {
          ref: 'release/v2',
          repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
        },
        head: {
          ref: 'patch-fork',
          repo: { full_name: 'external-contributor/terraform-provider-rancher2', fork: true },
          sha: 'forkprheadsha555',
        },
        user: {
          login: 'external-user',
          type: 'User',
        },
      },
    },
  };

  const result = await runFilterPr({
    github: {},
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], 'forkprheadsha555');
  assert.strictEqual(outputs['is_fork'], 'true');
});

test('runFilterPr - resolves fork PR via commit search fallback when pull_requests is empty in workflow_run', async () => {
  const outputs = {};
  let loggedInfo = '';
  const mockCore = {
    info: (msg) => { loggedInfo += msg; },
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'forkcommitsha123',
        head_branch: 'patch-1',
        head_repository: {
          full_name: 'contributor/terraform-provider-rancher2',
        },
        pull_requests: [], // Empty due to GHA fork limitation
      },
    },
  };

  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.commit_sha, 'forkcommitsha123');
      return [
        {
          number: 991,
          base: {
            ref: 'release/v2',
            repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
          },
          head: {
            ref: 'patch-1',
            repo: { full_name: 'contributor/terraform-provider-rancher2', fork: true },
            sha: 'forkcommitsha123',
          },
          user: {
            login: 'contributor',
            type: 'User',
          },
        },
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => {},
      },
    },
  };

  const result = await runFilterPr({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], 'forkcommitsha123');
  assert.strictEqual(outputs['is_fork'], 'true');
  assert.match(loggedInfo, /Resolved PR #991 via commit lookup/);
  assert.match(loggedInfo, /external fork/);
});

test('runFilterPr - safely skips external fork workflow_run without open PR without failing', async () => {
  const outputs = {};
  let loggedInfo = '';
  let failed = false;
  const mockCore = {
    info: (msg) => { loggedInfo += msg; },
    warning: () => {},
    setFailed: () => { failed = true; },
    setOutput: (key, val) => {
      outputs[key] = val;
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'unmatchedforksha',
        head_branch: 'custom-branch',
        head_repository: {
          full_name: 'external-user/terraform-provider-rancher2',
        },
        pull_requests: [],
      },
    },
  };

  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({ data: [] }),
      },
    },
  };

  const result = await runFilterPr({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(failed, false);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], 'unmatchedforksha');
  assert.strictEqual(outputs['is_fork'], 'true');
  assert.match(loggedInfo, /Triggered by an external fork repository without an identifiable open PR/);
});

test('runE2EStatus - sets core.setFailed and returns false when createCommitStatus fails on non-fork PR', async () => {
  let failedMsg = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
  };

  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async () => {
          throw new Error('API 500 Internal Error');
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 9999,
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'true',
      TEST_RESULT: 'success',
      AUDIT_RESULT: 'success',
      CLEANUP_RESULT: 'success',
      SHA: 'abc1234',
      IS_FORK: 'false',
    },
  };

  const result = await runE2EStatus({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(result, false);
  assert.match(failedMsg, /Failed to post commit status to SHA abc1234: API 500 Internal Error/);
});

test('runE2EStatus - handles createCommitStatus failure on fork PR gracefully without setFailed', async () => {
  let loggedInfo = '';
  let failed = false;
  const mockCore = {
    info: (msg) => { loggedInfo += msg; },
    warning: () => {},
    setFailed: () => { failed = true; },
  };

  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async () => {
          const err = new Error('Validation Failed');
          err.status = 422;
          throw err;
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 9999,
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'false',
      SHA: 'forksha456',
      IS_FORK: 'true',
    },
  };

  const result = await runE2EStatus({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(result, true);
  assert.strictEqual(failed, false);
  assert.match(loggedInfo, /Unable to post commit status to fork PR commit/);
});

test('runE2EStatus - posts skip description when shouldRun is not true', async () => {
  let createdStatus = null;
  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async (args) => {
          createdStatus = args;
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 9999,
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'false',
      SHA: 'skipsha789',
      IS_FORK: 'false',
    },
  };

  const result = await runE2EStatus({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, setFailed: () => {} },
    process: mockProcess,
  });

  assert.strictEqual(result, true);
  assert.ok(createdStatus);
  assert.strictEqual(createdStatus.state, 'success');
  assert.strictEqual(createdStatus.description, 'E2E tests safely skipped (PR does not qualify)');
});

test('resolveTargetPr - returns pull_request directly from context payload', async () => {
  const mockPr = { number: 456, head: { sha: 'directsha123' } };
  const pr = await resolveTargetPr({
    context: { payload: { pull_request: mockPr } },
    core: { info: () => {}, setFailed: () => {} },
    process: { env: {} },
  });
  assert.strictEqual(pr, mockPr);
});

test('resolveTargetPr - fetches PR using PR_NUMBER from process env', async () => {
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 555);
          return { data: { number: 555, head: { sha: 'sha555' } } };
        },
      },
    },
  };
  const pr = await resolveTargetPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }, payload: {} },
    core: { info: () => {}, setFailed: () => {} },
    process: { env: { PR_NUMBER: '555' } },
  });
  assert.strictEqual(pr.number, 555);
});

test('resolveTargetPr - uses workflow_run.pull_requests[0].number fallback', async () => {
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 777);
          return { data: { number: 777, head: { sha: 'sha777' } } };
        },
      },
    },
  };
  const pr = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          pull_requests: [{ number: 777 }],
        },
      },
    },
    core: { info: () => {}, setFailed: () => {} },
    process: { env: {} },
  });
  assert.strictEqual(pr.number, 777);
});

test('resolveTargetPr - returns null and sets failed when pulls.get throws error', async () => {
  let failedMsg = '';
  const mockGithub = {
    rest: {
      pulls: {
        get: async () => {
          throw new Error('PR not found');
        },
      },
    },
  };
  const pr = await resolveTargetPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }, payload: {} },
    core: { info: () => {}, setFailed: (msg) => { failedMsg = msg; } },
    process: { env: { PR_NUMBER: '999' } },
  });
  assert.strictEqual(pr, null);
  assert.match(failedMsg, /Failed to fetch PR #999: PR not found/);
});

test('runE2EStatus - posts filter failure description when filter job fails', async () => {
  let createdStatus = null;
  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async (args) => {
          createdStatus = args;
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    runId: 9999,
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'failure',
      SHOULD_RUN: 'false',
      SHA: 'failsha123',
      IS_FORK: 'false',
    },
  };

  const result = await runE2EStatus({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, setFailed: () => {} },
    process: mockProcess,
  });

  assert.strictEqual(result, false);
  assert.ok(createdStatus);
  assert.strictEqual(createdStatus.state, 'failure');
  assert.strictEqual(createdStatus.description, 'E2E test filter job failed');
});

test('resolveTargetPr - fetches full PR object when match is found via commit lookup', async () => {
  let fullPrFetched = false;
  const fullPrData = {
    number: 992,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'patch-2',
      repo: { full_name: 'contributor/terraform-provider-rancher2', fork: true },
      sha: 'commitsha992',
    },
    user: {
      login: 'contributor',
      type: 'User',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          {
            number: 992,
            state: 'open',
            head: { ref: 'patch-2', repo: { full_name: 'contributor/terraform-provider-rancher2' } },
          },
        ],
      },
      pulls: {
        get: async ({ pull_number }) => {
          if (pull_number === 992) {
            fullPrFetched = true;
            return { data: fullPrData };
          }
          throw new Error('Not found');
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commitsha992',
        head_branch: 'patch-2',
        head_repository: { full_name: 'contributor/terraform-provider-rancher2' },
        pull_requests: [],
      },
    },
  };

  const pr = await resolveTargetPr({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, warning: () => {}, setFailed: () => {} },
    process: { env: {} },
  });

  assert.strictEqual(fullPrFetched, true);
  assert.strictEqual(pr.number, 992);
  assert.strictEqual(pr.head.sha, 'commitsha992');
});

test('resolveTargetPr - fails closed and returns null when associated PRs do not match branch or repository', async () => {
  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          {
            number: 101,
            state: 'open',
            head: { ref: 'other-branch', repo: { full_name: 'unrelated/terraform-provider-rancher2' } },
          },
          {
            number: 102,
            state: 'closed',
            head: { ref: 'patch-3', repo: { full_name: 'contributor/terraform-provider-rancher2' } },
          },
        ],
      },
      pulls: {
        get: async () => {
          throw new Error('Should not be called');
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commitsha103',
        head_branch: 'patch-3',
        head_repository: { full_name: 'contributor/terraform-provider-rancher2' },
        pull_requests: [],
      },
    },
  };

  const pr = await resolveTargetPr({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, warning: () => {}, setFailed: () => {} },
    process: { env: {} },
  });

  assert.strictEqual(pr, null);
});

test('runFilterPr - fails closed and sets core.setFailed when resolved PR lacks head.sha', async () => {
  let failedMsg = '';
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
    setOutput: (key, val) => { outputs[key] = val; },
  };

  const mockPr = {
    number: 888,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'backport-123-release-v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
      // head.sha is intentionally missing/undefined
    },
    user: {
      login: 'github-actions[bot]',
      type: 'Bot',
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: mockPr,
    },
  };

  const result = await runFilterPr({
    github: {},
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.match(failedMsg, /missing a verified head commit SHA/);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], '');
});

test('resolveTargetPr - matches head branch even when one has refs/heads/ prefix', async () => {
  const fullPrData = {
    number: 993,
    base: {
      ref: 'release/v2',
      repo: { full_name: 'rancher/terraform-provider-rancher2', fork: false },
    },
    head: {
      ref: 'refs/heads/patch-normalized',
      repo: { full_name: 'contributor/terraform-provider-rancher2', fork: true },
      sha: 'commitsha993',
    },
    user: {
      login: 'contributor',
      type: 'User',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          {
            number: 993,
            state: 'open',
            head: { ref: 'refs/heads/patch-normalized', repo: { full_name: 'contributor/terraform-provider-rancher2' } },
          },
        ],
      },
      pulls: {
        get: async () => ({ data: fullPrData }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commitsha993',
        head_branch: 'patch-normalized', // without refs/heads/
        head_repository: { full_name: 'contributor/terraform-provider-rancher2' },
        pull_requests: [],
      },
    },
  };

  const pr = await resolveTargetPr({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, warning: () => {}, setFailed: () => {} },
    process: { env: {} },
  });

  assert.ok(pr);
  assert.strictEqual(pr.number, 993);
  assert.strictEqual(pr.head.sha, 'commitsha993');
});

test('runFilterPr - sets sha to empty string and is_fork to false when PR cannot be resolved for internal run', async () => {
  let failedMsg = '';
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
    setOutput: (key, val) => { outputs[key] = val; },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'missingprsha',
        head_branch: 'main',
        head_repository: {
          full_name: 'rancher/terraform-provider-rancher2',
          fork: false,
        },
        pull_requests: [],
      },
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({ data: [] }),
      },
    },
  };

  const result = await runFilterPr({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.match(failedMsg, /Missing pull_request in context payload and PR_NUMBER env var/);
  assert.strictEqual(outputs['should_run'], 'false');
  assert.strictEqual(outputs['sha'], '');
  assert.strictEqual(outputs['is_fork'], 'false');
});

test('runE2EStatus - uses process.env.GITHUB_RUN_ID when context.runId is missing', async () => {
  let createdStatus = null;
  const mockGithub = {
    rest: {
      repos: {
        createCommitStatus: async (args) => {
          createdStatus = args;
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    // runId intentionally missing
  };

  const mockProcess = {
    env: {
      FILTER_RESULT: 'success',
      SHOULD_RUN: 'true',
      TEST_RESULT: 'success',
      AUDIT_RESULT: 'success',
      CLEANUP_RESULT: 'success',
      SHA: 'runidsha987',
      IS_FORK: 'false',
      GITHUB_RUN_ID: '778899',
    },
  };

  const result = await runE2EStatus({
    github: mockGithub,
    context: mockContext,
    core: { info: () => {}, setFailed: () => {} },
    process: mockProcess,
  });

  assert.strictEqual(result, true);
  assert.ok(createdStatus);
  assert.strictEqual(createdStatus.target_url, 'https://github.com/rancher/terraform-provider-rancher2/actions/runs/778899');
});
