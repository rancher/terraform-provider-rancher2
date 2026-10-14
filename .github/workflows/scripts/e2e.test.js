import assert from 'node:assert/strict';
import test from 'node:test';
import defaultHandler, {
  evaluateE2EStatus,
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
  let outputKey = '';
  let outputValue = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputKey = key;
      outputValue = val;
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
        },
        user: {
          login: 'github-actions[bot]',
          type: 'Bot',
        },
      },
    },
  };

  const result = await runFilterPr({ context: mockContext, core: mockCore });
  assert.strictEqual(result, true);
  assert.strictEqual(outputKey, 'should_run');
  assert.strictEqual(outputValue, 'true');
});

test('runFilterPr - sets step output should_run to false when PR missing from payload', async () => {
  let failedMsg = '';
  let outputKey = '';
  let outputValue = '';
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => { failedMsg = msg; },
    setOutput: (key, val) => {
      outputKey = key;
      outputValue = val;
    },
  };

  const mockContext = { payload: {} };
  const result = await runFilterPr({ context: mockContext, core: mockCore });
  assert.strictEqual(result, false);
  assert.strictEqual(outputKey, 'should_run');
  assert.strictEqual(outputValue, 'false');
  assert.match(failedMsg, /Missing pull_request/);
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
        head: { ref: 'backport-1-release-v2', repo: { full_name: 'rancher/repo', fork: false } },
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

