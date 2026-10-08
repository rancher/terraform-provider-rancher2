import assert from 'node:assert/strict';
import test from 'node:test';
import defaultHandler, {
  computeNextRcTag,
  runCheckMaintainer,
  runRcNotify,
  runTrackingIssue,
  runTriggerRcRelease,
} from './releases.js';

test('runCheckMaintainer - authorizes maintainer from MAINTAINERS env', async () => {
  let loggedInfo = null;
  const mockCore = {
    info: (msg) => { loggedInfo = msg; }
  };
  const mockContext = { actor: 'alice' };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(['alice', 'bob'])
    }
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, true);
  assert.match(loggedInfo, /alice.*authorized.*true/);
});

test('runCheckMaintainer - does not authorize bot accounts unless explicitly allowed', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'rancher-eio[bot]' };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(['matt'])
    }
  };

  let isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, false);

  mockProcess.env.MAINTAINERS = JSON.stringify(['matt', 'rancher-eio[bot]']);
  isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, true);
});

test('runCheckMaintainer - rejects unauthorized actors', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'unauthorized-user' };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(['matt'])
    }
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, false);
});

test('computeNextRcTag - increments existing RC tag', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.2' },
      { name: 'v15.2.0-rc.1' },
      { name: 'v15.2.0-rc.2' },
      { name: 'v15.2.0-rc.3' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.4');
});

test('computeNextRcTag - starts new RC cycle as rc.1 when latest is full release', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.1.1' },
      { name: 'v15.1.2' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.1.3-rc.1');
});

test('computeNextRcTag - uses version from release-please-config.json release-as if present', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.1.1' },
      { name: 'v15.1.2' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === 'release-please-config.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ 'release-as': 'v15.3.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.3.0-rc.1');
});

test('computeNextRcTag - uses version from .release-please-manifest.json if higher than latest full release', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.1.1' },
      { name: 'v15.1.2' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ '.': '15.2.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.1');
});

test('computeNextRcTag - ignores manifest with mismatched major version', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ '.': '14.0.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.1.1-rc.1');
});

test('computeNextRcTag - advances active RC when manifest version equals latest full release', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.2' },
      { name: 'v15.1.3-rc.1' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ '.': '15.1.2' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.1.3-rc.2');
});

test('computeNextRcTag - normalizes branch ref prefix', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'refs/heads/release/v15'
  });

  assert.strictEqual(tag, 'v15.1.1-rc.1');
});

test('computeNextRcTag - brand new release branch with no tags', async () => {
  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v16'
  });

  assert.strictEqual(tag, 'v16.0.0-rc.1');
});

test('computeNextRcTag - rejects invalid branch name', async () => {
  const mockGithub = {
    paginate: async () => []
  };

  await assert.rejects(
    async () => {
      await computeNextRcTag({
        github: mockGithub,
        owner: 'rancher',
        repo: 'terraform-provider-rancher2',
        branch: 'main'
      });
    },
    { message: /does not match expected pattern/ }
  );
});

test('runTriggerRcRelease - dispatches workflow with computed tag and explicit sha', async () => {
  let dispatchedParams = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      },
      actions: {
        createWorkflowDispatch: async (params) => {
          dispatchedParams = params;
        }
      }
    }
  };

  let failedMessage = null;
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    }
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        base: { ref: 'refs/heads/release/v15' },
        merge_commit_sha: 'sha123456789'
      }
    }
  };

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: 'test-token'
    }
  };

  const mockGetOctokit = () => mockGithub;

  await runTriggerRcRelease({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit
  });

  assert.strictEqual(failedMessage, null);
  assert.notStrictEqual(dispatchedParams, null);
  assert.strictEqual(dispatchedParams.owner, 'rancher');
  assert.strictEqual(dispatchedParams.repo, 'terraform-provider-rancher2');
  assert.strictEqual(dispatchedParams.workflow_id, 'rc-release.yml');
  assert.strictEqual(dispatchedParams.ref, 'main');
  assert.deepStrictEqual(dispatchedParams.inputs, {
    branch: 'release/v15',
    sha: 'sha123456789',
    tag: 'v15.1.1-rc.1'
  });
});

test('runTriggerRcRelease - succeeds without sha and dispatches with sha omitted', async () => {
  let dispatchedParams = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      },
      actions: {
        createWorkflowDispatch: async (params) => {
          dispatchedParams = params;
        }
      }
    }
  };

  let failedMessage = null;
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    }
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        base: { ref: 'release/v15' }
      }
    }
  };

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: 'test-token'
    }
  };

  const mockGetOctokit = () => mockGithub;

  await runTriggerRcRelease({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit
  });

  assert.strictEqual(failedMessage, null);
  assert.notStrictEqual(dispatchedParams, null);
  assert.deepStrictEqual(dispatchedParams.inputs, {
    branch: 'release/v15',
    tag: 'v15.1.1-rc.1'
  });
});

test('runTriggerRcRelease - warns when GITHUB_MERGE_TOKEN is not provided', async () => {
  let warningMessage = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      },
      actions: {
        createWorkflowDispatch: async () => {}
      }
    }
  };

  const mockCore = {
    info: () => {},
    warning: (msg) => { warningMessage = msg; },
    setFailed: () => {}
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        base: { ref: 'release/v15' }
      }
    }
  };

  await runTriggerRcRelease({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} }
  });

  assert.match(warningMessage, /GITHUB_MERGE_TOKEN not provided/);
});

test('runTriggerRcRelease - fails gracefully when branch is missing', async () => {
  let failedMessage = null;
  const mockCore = {
    info: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    }
  };

  await runTriggerRcRelease({
    github: {},
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }, payload: {} },
    core: mockCore,
    process: { env: {} }
  });

  assert.match(failedMessage, /Target branch must be provided/);
});

test('runTriggerRcRelease - catches errors gracefully and sets failure', async () => {
  let failedMessage = null;
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: (msg) => {
      failedMessage = msg;
    }
  };

  const mockGithub = {
    paginate: async () => {
      throw new Error('API Rate Limit Exceeded');
    },
    rest: {
      repos: {
        listTags: () => {}
      }
    }
  };

  await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        pull_request: {
          base: { ref: 'release/v15' },
          merge_commit_sha: 'sha123'
        }
      }
    },
    core: mockCore,
    process: { env: {} }
  });

  assert.match(failedMessage, /Failed to trigger RC release workflow: API Rate Limit Exceeded/);
});

test('runTriggerRcRelease - skips RC release when workflow run published a full release', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.run_id, 12345);
      return [
        { name: 'Generate Full Release', conclusion: 'success' }
      ];
    },
    rest: {
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 12345,
          head_branch: 'release/v15',
          head_sha: 'abc12345'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'abc12345', WORKFLOW_RUN_ID: '12345' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - skips RC release when commit is from release-please PR', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async (fn, params) => {
      if (params.run_id) {
        return [{ name: 'Generate Full Release', conclusion: 'skipped' }];
      }
      assert.strictEqual(params.commit_sha, 'commit-rp-123');
      return [
        {
          number: 2506,
          state: 'closed',
          merged_at: '2026-10-06T15:00:00Z',
          head: { ref: 'release-please--branches--release/v15' }
        }
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 54321,
          head_branch: 'release/v15',
          head_sha: 'commit-rp-123'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'commit-rp-123', WORKFLOW_RUN_ID: '54321' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /is a release-please PR.*skipping RC release/);
});

test('runTriggerRcRelease - skips RC release when direct PR payload is a release-please PR', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        pull_request: {
          number: 2506,
          head: { ref: 'release-please--branches--release/v15' }
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /is a release-please PR.*skipping RC release/);
});

test('runTriggerRcRelease - skips RC release when multiple PRs associated and one is release-please', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async (fn, params) => {
      if (params.commit_sha) {
        return [
          {
            number: 100,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            head: { ref: 'some-feature-branch' }
          },
          {
            number: 101,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            head: { ref: 'release-please--branches--release/v15' }
          }
        ];
      }
      return [];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'shared-commit-sha' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /is a release-please PR.*skipping RC release/);
});

test('runTriggerRcRelease - skips RC release when workflow run job matches full release case-insensitively', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'generate full release', conclusion: 'success' }
    ],
    rest: {
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 9999,
          head_branch: 'release/v15',
          head_sha: 'sha-9999'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-9999', WORKFLOW_RUN_ID: '9999' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - skips RC release when workflow run job matches legacy publish job name', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'publish', conclusion: 'success' }
    ],
    rest: {
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 8888,
          head_branch: 'release/v13',
          head_sha: 'sha-8888'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v13', SHA: 'sha-8888', WORKFLOW_RUN_ID: '8888' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - skips RC release when workflow run job matches release job name', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async () => [
      { name: 'release', conclusion: 'success' }
    ],
    rest: {
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 7777,
          head_branch: 'release/v14',
          head_sha: 'sha-7777'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v14', SHA: 'sha-7777', WORKFLOW_RUN_ID: '7777' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - skips RC release when target branch starts with release-please', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release-please--branches--release/v15' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /is a release-please branch; skipping RC release/);
});

test('runTriggerRcRelease - skips RC release when workflow_run.pull_requests contains a release-please PR without commit lookup', async () => {
  let dispatched = false;
  let loggedInfo = null;
  let commitLookupCalled = false;
  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {
          commitLookupCalled = true;
        }
      },
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 7777,
          head_branch: 'release/v15',
          head_sha: 'sha-7777',
          pull_requests: [
            {
              number: 3000,
              head: { ref: 'release-please--branches--release/v15' }
            }
          ]
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-7777', WORKFLOW_RUN_ID: '7777' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.strictEqual(commitLookupCalled, false);
  assert.match(loggedInfo, /is a release-please PR.*skipping RC release/);
});

test('runTriggerRcRelease - matches release-please PR targeting the specific branch when multiple PRs associated', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async (fn, params) => {
      if (params.commit_sha) {
        return [
          {
            number: 201,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            head: { ref: 'release-please--branches--release/v14' },
            base: { ref: 'release/v14' }
          },
          {
            number: 202,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            head: { ref: 'release-please--branches--release/v15' },
            base: { ref: 'release/v15' }
          }
        ];
      }
      return [];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-branch-match' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #202 is a release-please PR/);
});

test('runTriggerRcRelease - skips RC release when PR has autorelease label even with custom branch name', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    paginate: async (fn, params) => {
      if (params.commit_sha) {
        return [
          {
            number: 501,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            head: { ref: 'custom-automation-branch' },
            labels: [{ name: 'autorelease: pending' }]
          }
        ];
      }
      return [];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-autorelease-label' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #501 is a release-please PR/);
});

test('runTriggerRcRelease - falls back to non-paginated listJobsForWorkflowRun when github.paginate is undefined', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    rest: {
      actions: {
        listJobsForWorkflowRun: async ({ run_id }) => {
          assert.strictEqual(run_id, 4444);
          return {
            data: {
              jobs: [{ name: 'Generate Full Release', conclusion: 'success' }]
            }
          };
        },
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 4444,
          head_branch: 'release/v15',
          head_sha: 'sha-4444'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-4444', WORKFLOW_RUN_ID: '4444' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - falls back to non-paginated listPullRequestsAssociatedWithCommit when github.paginate is undefined', async () => {
  let dispatched = false;
  let loggedInfo = null;
  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async ({ commit_sha }) => {
          assert.strictEqual(commit_sha, 'sha-fallback-no-paginate');
          return {
            data: [
              {
                number: 601,
                state: 'closed',
                merged_at: '2026-10-06T12:00:00Z',
                head: { ref: 'release-please--branches--release/v15' }
              }
            ]
          };
        }
      },
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-fallback-no-paginate' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #601 is a release-please PR/);
});

test('runTriggerRcRelease - skips RC release when PR title matches chore release pattern', async () => {
  let dispatched = false;
  let loggedInfo = '';

  const mockGithub = {
    paginate: async () => [
      {
        number: 701,
        title: 'chore(release/v15): release 15.2.0',
        state: 'closed',
        merged_at: '2026-10-06T12:00:00Z',
        head: { ref: 'custom-release-branch' },
        base: { ref: 'refs/heads/release/v15' }
      }
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-title-release' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #701 is a release-please PR/);
});

test('runTriggerRcRelease - matches candidate release-please PR when base.ref has refs/heads/ prefix', async () => {
  let dispatched = false;
  let loggedInfo = '';

  const mockGithub = {
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 12345,
          pull_requests: [
            {
              number: 702,
              head: { ref: 'release-please--branches--release/v15' },
              base: { ref: 'refs/heads/release/v15' }
            }
          ]
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #702 is a release-please PR/);
});

test('runTriggerRcRelease - falls back to pulls.get when commit lookup throws and workflow_run.pull_requests has candidate', async () => {
  let dispatched = false;
  let loggedInfo = '';

  const mockGithub = {
    paginate: async () => {
      throw new Error('API Rate Limit (403)');
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 801);
          return {
            data: {
              number: 801,
              title: 'chore(release/v15): release 15.2.0',
              head: { ref: 'custom-pr-branch' }
            }
          };
        }
      },
      actions: {
        listJobsForWorkflowRun: async () => ({ data: { jobs: [] } }),
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 55555,
          pull_requests: [
            {
              number: 801,
              base: { ref: 'release/v15' }
            }
          ]
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sha-pr-fallback' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #801 is a release-please PR/);
});

test('runTriggerRcRelease - safely handles object { jobs: [...] } from non-paginated listJobsForWorkflowRun', async () => {
  let dispatched = false;
  let loggedInfo = '';

  const mockGithub = {
    rest: {
      actions: {
        listJobsForWorkflowRun: async () => ({
          data: {
            jobs: [
              { name: 'Generate Full Release', conclusion: 'success' }
            ]
          }
        })
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 66666,
          head_branch: 'release/v15'
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', WORKFLOW_RUN_ID: '66666' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /attempted a full release/);
});

test('runTriggerRcRelease - resolves branch and sha from workflow_run payload when process.env lacks them', async () => {
  let dispatchedBranch = null;
  let dispatchedSha = null;
  let dispatchedTag = null;

  const mockGithub = {
    paginate: async (fn, params) => {
      if (params?.run_id) return [];
      if (params?.commit_sha) return [];
      return [
        { name: 'v15.1.0' },
        { name: 'v15.2.0-rc.1' }
      ];
    },
    rest: {
      actions: {
        listJobsForWorkflowRun: () => {},
        createWorkflowDispatch: async ({ inputs }) => {
          dispatchedBranch = inputs.branch;
          dispatchedSha = inputs.sha;
          dispatchedTag = inputs.tag;
        }
      },
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 77777,
          head_branch: 'release/v15',
          head_sha: 'sha-from-payload'
        }
      }
    },
    core: mockCore,
    process: { env: {} }
  });

  assert.strictEqual(result, 'v15.2.0-rc.2');
  assert.strictEqual(dispatchedBranch, 'release/v15');
  assert.strictEqual(dispatchedSha, 'sha-from-payload');
  assert.strictEqual(dispatchedTag, 'v15.2.0-rc.2');
});

test('runTriggerRcRelease - matches release-please PR when head.ref has refs/heads/ prefix', async () => {
  let dispatched = false;
  let loggedInfo = '';

  const mockGithub = {
    rest: {
      actions: {
        createWorkflowDispatch: async () => {
          dispatched = true;
        }
      }
    }
  };

  const mockCore = {
    info: (msg) => { loggedInfo = msg; },
    warning: () => {},
    setFailed: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 12345,
          pull_requests: [
            {
              number: 703,
              head: { ref: 'refs/heads/release-please--branches--release/v15' },
              base: { ref: 'release/v15' }
            }
          ]
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15' } }
  });

  assert.strictEqual(result, null);
  assert.strictEqual(dispatched, false);
  assert.match(loggedInfo, /PR #703 is a release-please PR/);
});

test('computeNextRcTag - falls back to direct listTags when github.paginate is undefined', async () => {
  const mockGithub = {
    rest: {
      repos: {
        listTags: async () => ({
          data: [
            { name: 'v15.1.0' },
            { name: 'v15.2.0-rc.1' }
          ]
        }),
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.2');
});


test('computeNextRcTag - uses version from release-please-config.json packages release-as', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === 'release-please-config.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ packages: { '.': { 'release-as': '15.4.0' } } })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.4.0-rc.1');
});

test('computeNextRcTag - handles trailing slash on branch name', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'refs/heads/release/v15/'
  });

  assert.strictEqual(tag, 'v15.1.1-rc.1');
});

test('runRcNotify - skips notification if tag is not an RC', async () => {
  const loggedInfo = [];
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {}
  };
  await runRcNotify({
    github: {},
    context: {},
    core: mockCore,
    process: { env: { TAG: 'v15.2.0', BRANCH: 'release/v15' } }
  });
  assert.ok(loggedInfo.some(m => m.includes('does not appear to be an RC')));
});

test('runRcNotify - fails when tag or branch is missing', async () => {
  let failedMessage = null;
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMessage = msg; }
  };
  await runRcNotify({
    github: {},
    context: {},
    core: mockCore,
    process: { env: {} }
  });
  assert.match(failedMessage, /tagName and branchLabel must be provided/);
});

test('runRcNotify - throws when branch format is invalid', async () => {
  const mockCore = {
    info: () => {},
    setFailed: () => {}
  };
  await assert.rejects(
    async () => {
      await runRcNotify({
        github: {},
        context: {},
        core: mockCore,
        process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'main' } }
      });
    },
    { message: /Target branch label "main" is invalid/ }
  );
});

test('runRcNotify - handles empty matching issues gracefully', async () => {
  const loggedInfo = [];
  const mockGithub = {
    paginate: async () => [],
    rest: {
      search: {
        issuesAndPullRequests: () => {}
      }
    }
  };
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {}
  };
  await runRcNotify({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'release/v15' } }
  });
  assert.ok(loggedInfo.some(m => m.includes('No matching issues found')));
});

test('runRcNotify - comments on matching issues with release url', async () => {
  const loggedInfo = [];
  const createdComments = [];
  const mockGithub = {
    paginate: async (endpoint, params) => {
      assert.match(params.q, /label:"release\/v15"/);
      assert.match(params.q, /label:"internal\/backport"/);
      assert.match(params.q, /label:"internal\/merged"/);
      return [{ number: 101 }, { number: 102 }];
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {}
      },
      issues: {
        createComment: async (params) => {
          createdComments.push(params);
        }
      }
    }
  };
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {}
  };
  await runRcNotify({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'release/v15' } }
  });
  assert.strictEqual(createdComments.length, 2);
  assert.strictEqual(createdComments[0].issue_number, 101);
  assert.match(createdComments[0].body, /v15\.2\.0-rc\.1/);
  assert.strictEqual(createdComments[1].issue_number, 102);
  assert.ok(loggedInfo.some(m => m.includes('Notified 2 issues')));
});

test('runRcNotify - handles error on comment creation', async () => {
  let failedMessage = null;
  const mockGithub = {
    paginate: async () => [{ number: 999 }],
    rest: {
      search: { issuesAndPullRequests: () => {} },
      issues: {
        createComment: async () => {
          throw new Error('API Rate Limit');
        }
      }
    }
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMessage = msg; }
  };
  await runRcNotify({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'release/v15' } }
  });
  assert.match(failedMessage, /Failed to comment on issue #999: API Rate Limit/);
});

test('runRcNotify - strips trailing slash from branch name and succeeds', async () => {
  const createdComments = [];
  const mockGithub = {
    paginate: async (endpoint, params) => {
      assert.match(params.q, /label:"release\/v15"/);
      return [{ number: 201 }];
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {}
      },
      issues: {
        createComment: async (params) => {
          createdComments.push(params);
        }
      }
    }
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { assert.fail(`Unexpected failure: ${msg}`); }
  };
  await runRcNotify({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'release/v15/' } }
  });
  assert.strictEqual(createdComments.length, 1);
  assert.strictEqual(createdComments[0].issue_number, 201);
});

test('computeNextRcTag - increments RC when existing tags use un-dotted format like -rc1, -rc2', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.2.0-rc1' },
      { name: 'v15.2.0-rc2' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.3');
});

test('computeNextRcTag - correctly increments double-digit RC tags numerically', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.2.0-rc.1' },
      { name: 'v15.2.0-rc.2' },
      { name: 'v15.2.0-rc.9' },
      { name: 'v15.2.0-rc.10' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.11');
});

test('computeNextRcTag - selects highest active RC regardless of tag ordering', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.4-rc.1' },
      { name: 'v15.1.3' },
      { name: 'v15.2.0-rc.2' },
      { name: 'v15.2.0-rc.1' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.3');
});

test('runCheckMaintainer - does not allow substring match when MAINTAINERS is JSON string or array', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'super' };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify('superuser')
    }
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, false);
});

test('computeNextRcTag - falls back to first object value in manifest if dot key is absent', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ 'terraform-provider-rancher2': '15.5.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.5.0-rc.1');
});

test('runTriggerRcRelease - sets step output tag and returns tag on success', async () => {
  const outputs = {};
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: (key, val) => {
      outputs[key] = val;
    }
  };

  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      },
      actions: {
        createWorkflowDispatch: async () => {}
      }
    }
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        base: { ref: 'release/v15' }
      }
    }
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'token' } },
    getOctokit: () => mockGithub
  });

  assert.strictEqual(result, 'v15.1.1-rc.1');
  assert.strictEqual(outputs.tag, 'v15.1.1-rc.1');
});

test('runCheckMaintainer - authorizes maintainers case-insensitively', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'Alice' };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(['alice', 'BOB'])
    }
  };

  const isAuth = await runCheckMaintainer({
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(isAuth, true);
});

test('computeNextRcTag - uses version from release-please-config.json packages with arbitrary package key', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === 'release-please-config.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ packages: { 'terraform-provider-rancher2': { 'release-as': '15.6.0' } } })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.6.0-rc.1');
});

test('runRcNotify - fails gracefully when search API throws an error', async () => {
  let failedMessage = null;
  const mockGithub = {
    paginate: async () => {
      throw new Error('GitHub Search API rate limit exceeded');
    },
    rest: {
      search: {
        issuesAndPullRequests: () => {}
      }
    }
  };
  const mockCore = {
    info: () => {},
    setFailed: (msg) => { failedMessage = msg; }
  };
  await runRcNotify({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: { env: { TAG: 'v15.2.0-rc.1', BRANCH: 'release/v15' } }
  });
  assert.match(failedMessage, /Failed to search issues for RC notification: GitHub Search API rate limit exceeded/);
});

test('computeNextRcTag - falls back to manifestVersion if configVersion has mismatched major version', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === 'release-please-config.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ 'release-as': 'v14.9.0' })).toString('base64')
              }
            };
          }
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ '.': '15.3.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.3.0-rc.1');
});

test('computeNextRcTag - does not regress version when active RC cycle exceeds manifest targetBase', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.3.0-rc.1' },
      { name: 'v15.3.0-rc.2' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === '.release-please-manifest.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ '.': '15.2.0' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.3.0-rc.3');
});

test('computeNextRcTag - parses release-as with pre-release suffix cleanly', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' }
    ],
    rest: {
      repos: {
        getContent: async ({ path }) => {
          if (path === 'release-please-config.json') {
            return {
              data: {
                content: Buffer.from(JSON.stringify({ 'release-as': 'v15.4.0-rc.1' })).toString('base64')
              }
            };
          }
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.4.0-rc.1');
});

test('runCheckMaintainer - rejects empty or missing actor', async () => {
  const mockCore = { info: () => {} };
  const mockProcess = {
    env: {
      MAINTAINERS: JSON.stringify(['alice', ''])
    }
  };

  const isAuthEmpty = await runCheckMaintainer({
    context: { actor: '' },
    core: mockCore,
    process: mockProcess
  });
  assert.strictEqual(isAuthEmpty, false);

  const isAuthMissing = await runCheckMaintainer({
    context: {},
    core: mockCore,
    process: mockProcess
  });
  assert.strictEqual(isAuthMissing, false);
});

test('runTriggerRcRelease - re-throws error when core.setFailed is not available', async () => {
  const mockGithub = {
    paginate: async () => {
      throw new Error('API failure');
    },
    rest: {
      repos: { listTags: () => {} }
    }
  };

  await assert.rejects(
    async () => {
      await runTriggerRcRelease({
        github: mockGithub,
        context: {
          repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
          payload: {
            pull_request: {
              base: { ref: 'release/v15' },
              merge_commit_sha: 'sha123'
            }
          }
        },
        core: {},
        process: { env: {} }
      });
    },
    { message: /API failure/ }
  );
});

test('computeNextRcTag - correctly parses case-insensitive uppercase RC tags', async () => {
  const mockGithub = {
    paginate: async () => [
      { name: 'v15.1.0' },
      { name: 'v15.2.0-RC.1' }
    ],
    rest: {
      repos: {
        getContent: async () => {
          throw new Error('Not found');
        }
      }
    }
  };

  const tag = await computeNextRcTag({
    github: mockGithub,
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    branch: 'release/v15'
  });

  assert.strictEqual(tag, 'v15.2.0-rc.2');
});

test('runCheckMaintainer - safely handles null JSON string in MAINTAINERS', async () => {
  const mockCore = { info: () => {} };
  const isAuth = await runCheckMaintainer({
    context: { actor: 'null' },
    core: mockCore,
    process: { env: { MAINTAINERS: 'null' } }
  });

  assert.strictEqual(isAuth, false);
});

test('default export - routes check-maintainer mode', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'alice' };
  const mockProcess = {
    env: {
      SCRIPT_MODE: 'check-maintainer',
      MAINTAINERS: JSON.stringify(['alice'])
    }
  };

  const result = await defaultHandler({
    github: {},
    context: mockContext,
    core: mockCore,
    process: mockProcess
  });

  assert.strictEqual(result, true);
});

test('default export - routes trigger-rc-release mode', async () => {
  const mockGithub = {
    paginate: async () => [{ name: 'v15.1.0' }],
    rest: {
      repos: { getContent: async () => { throw new Error('Not found'); } },
      actions: { createWorkflowDispatch: async () => {} }
    }
  };
  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: () => {}
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: { base: { ref: 'release/v15' } }
    }
  };

  const result = await defaultHandler({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { SCRIPT_MODE: 'trigger-rc-release', GITHUB_MERGE_TOKEN: 'token' } },
    getOctokit: () => mockGithub
  });

  assert.strictEqual(result, 'v15.1.1-rc.1');
});

test('default export - routes rc-notify mode', async () => {
  const loggedInfo = [];
  const mockCore = {
    info: (msg) => loggedInfo.push(msg),
    setFailed: () => {}
  };

  await defaultHandler({
    github: {},
    context: {},
    core: mockCore,
    process: { env: { SCRIPT_MODE: 'rc-notify', TAG: 'v15.2.0', BRANCH: 'release/v15' } }
  });

  assert.ok(loggedInfo.some(m => m.includes('does not appear to be an RC')));
});

test('default export - throws when SCRIPT_MODE is unknown or undefined', async () => {
  await assert.rejects(
    async () => {
      await defaultHandler({
        process: { env: { SCRIPT_MODE: 'unsupported-mode' } }
      });
    },
    { message: /Unknown release script mode: unsupported-mode/ }
  );

  await assert.rejects(
    async () => {
      await defaultHandler({
        process: { env: {} }
      });
    },
    { message: /Unknown release script mode: undefined/ }
  );
});

function createMockCore() {
  return {
    infoMessages: [],
    warningMessages: [],
    failedMessage: null,
    info(msg) {
      this.infoMessages.push(msg);
    },
    warning(msg) {
      this.warningMessages.push(msg);
    },
    setFailed(msg) {
      this.failedMessage = msg;
    },
  };
}

test('runTrackingIssue - successfully converts referenced issue into tracking issue', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [
          { name: 'main' },
          { name: 'release/v14' },
          { name: 'release/v15' },
          { name: 'release/v9' },
        ];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 42,
            title: 'Fix cluster driver bug',
            html_url: 'https://github.com/rancher/terraform-provider-rancher2/pull/42',
            body: 'Some text\n- Addresses: #100\nMore text',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: {
        listBranches: async () => {},
      },
      search: {
        issuesAndPullRequests: async () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 100);
          return {
            data: {
              number: 100,
              title: 'Cluster driver issue',
              labels: [],
              state: 'open',
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          assert.strictEqual(issue_number, 100);
          labelsAdded.push(labels);
        },
        createComment: async ({ issue_number, body }) => {
          assert.strictEqual(issue_number, 100);
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  // First step: internal/tracking and internal/user
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking', 'internal/user']);
  // Second step: tracking comments (PR linking + release instructions)
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #42 -->/);
  assert.match(commentsCreated[0], /This is the tracking issue for PR #42/);
  assert.match(commentsCreated[1], /A label has been added for the latest release branch/);
  // Third step: latest release branch label added last
  assert.deepStrictEqual(labelsAdded[1], ['release/v15']);
});

test('runTrackingIssue - respects existing internal/user label', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 43,
            title: 'Add feature',
            body: '- Addresses: #101',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 101,
            title: 'User issue',
            labels: [{ name: 'internal/user' }],
            state: 'open',
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  // internal/user was already present, so only internal/tracking should be added
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking']);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #43 -->/);
  assert.match(commentsCreated[1], /A label has been added for the latest release branch/);
  assert.deepStrictEqual(labelsAdded[1], ['release/v15']);
});

test('runTrackingIssue - skips issue already marked internal/tracking', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 44,
            body: '- Addresses: #102',
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 102,
            labels: [{ name: 'internal/tracking' }, { name: 'internal/user' }],
            state: 'open',
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(mockCore.infoMessages.some(m => m.includes('already has \'internal/tracking\' label')));
});

test('runTrackingIssue - skips PR without valid Addresses line', async () => {
  const mockCore = createMockCore();

  let getIssueCalled = false;
  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 45,
            body: 'Just a PR description with no issue reference',
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => {
          getIssueCalled = true;
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(getIssueCalled, false);
  assert.ok(mockCore.infoMessages.some(m => m.includes('does not contain valid issue references')));
});

test('runTrackingIssue - fails when no release branches exist', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: { listBranches: async () => {} },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, 'No branches found');
});

test('runTrackingIssue - respects target release branch from PR labels over latest branch', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }, { name: 'release/v14' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 50,
            title: 'Fix issue for v14',
            body: '- Addresses: #200',
            labels: [{ name: 'release/v14' }],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            title: 'Issue 200',
            labels: [],
            state: 'open',
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking', 'internal/user']);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #50 -->/);
  assert.match(commentsCreated[1], /A label has been added for the latest release branch/);
  assert.deepStrictEqual(labelsAdded[1], ['release/v14']);
});

test('runTrackingIssue - skips referenced issue that is closed', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 51,
            body: '- Addresses: #201',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 201,
            state: 'closed',
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(mockCore.warningMessages.some(m => m.includes('Issue #201 referenced in PR #51 is already closed. Skipping.')));
});

test('runTrackingIssue - skips referenced item that is a pull request', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 52,
            body: '- Addresses: #202',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 202,
            state: 'open',
            pull_request: {},
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 0);
  assert.strictEqual(commentsCreated.length, 0);
  assert.ok(mockCore.warningMessages.some(m => m.includes('Issue #202 referenced in PR #52 is a pull request, not an issue. Skipping.')));
});

test('runTrackingIssue - handles multiple referenced issues where one already has internal/tracking', async () => {
  const mockCore = createMockCore();

  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 53,
            body: '- Addresses: #100\n- Addresses: #101',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 100) {
            return {
              data: {
                number: 100,
                labels: [{ name: 'internal/tracking' }],
                state: 'open',
              },
            };
          }
          return {
            data: {
              number: 101,
              labels: [],
              state: 'open',
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          labelsAdded.push({ issue_number, labels });
        },
        createComment: async ({ issue_number, body }) => {
          commentsCreated.push({ issue_number, body });
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes("Issue #100 already has 'internal/tracking' label. Skipping.")));
  assert.strictEqual(labelsAdded.length, 2);
  assert.deepStrictEqual(labelsAdded[0], { issue_number: 101, labels: ['internal/tracking', 'internal/user'] });
  assert.deepStrictEqual(labelsAdded[1], { issue_number: 101, labels: ['release/v15'] });
  assert.strictEqual(commentsCreated.length, 2);
  assert.strictEqual(commentsCreated[0].issue_number, 101);
  assert.strictEqual(commentsCreated[1].issue_number, 101);
  assert.match(commentsCreated[0].body, /<!-- tracking-pr: #53 -->/);
  assert.match(commentsCreated[1].body, /A label has been added for the latest release branch/);
});

test('run via default export with SCRIPT_MODE=tracking-issue', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      SCRIPT_MODE: 'tracking-issue',
    },
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v16' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
    },
  };

  await defaultHandler({ github: mockGithub, context: {}, core: mockCore, process: mockProcess });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Latest release branch detected: release/v16')));
});

test('runTrackingIssue - fails when branches exist but none match release/v pattern', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async () => [{ name: 'main' }, { name: 'feature/foo' }, { name: 'release/v-custom' }],
    rest: {
      repos: { listBranches: async () => {} },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, 'No release branches found');
});

test('runTrackingIssue - logs warning and continues when fetching issue fails', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 60,
            body: '- Addresses: #300',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => {
          throw new Error('Not Found (404)');
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.warningMessages.some(m => m.includes('Could not fetch issue #300 for PR #60: Not Found (404)')));
});

test('runTrackingIssue - records error and fails when addLabels throws', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 61,
            html_url: 'https://github.com/rancher/terraform-provider-rancher2/pull/61',
            body: '- Addresses: #301',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 301,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async () => {
          throw new Error('Secondary rate limit reached');
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.ok(mockCore.failedMessage);
  assert.ok(mockCore.failedMessage.includes('Failed to process some pull requests'));
  assert.ok(mockCore.failedMessage.includes('Secondary rate limit reached'));
});

test('runTrackingIssue - handles PR description with inline code backticks', async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 62,
            body: '- Addresses: `#302`',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 302,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking', 'internal/user']);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #62 -->/);
  assert.match(commentsCreated[0], /This is the tracking issue for PR #62/);
  assert.match(commentsCreated[1], /A label has been added for the latest release branch/);
  assert.deepStrictEqual(labelsAdded[1], ['release/v15']);
});

test('runTrackingIssue - picks highest release branch when multiple are on PR', async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v16' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 63,
            body: '- Addresses: #303',
            labels: [{ name: 'release/v13' }, { name: 'release/v15' }, { name: 'release/v14' }],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 303,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async () => {},
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[1], ['release/v15']);
});

test('runTrackingIssue - does not add release branch label if already present on issue', async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];
  const commentsCreated = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 64,
            body: '- Addresses: #304',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 304,
            state: 'open',
            labels: [{ name: 'release/v15' }],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async ({ body }) => {
          commentsCreated.push(body);
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded.length, 1);
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking', 'internal/user']);
  assert.strictEqual(commentsCreated.length, 2);
  assert.match(commentsCreated[0], /<!-- tracking-pr: #64 -->/);
  assert.match(commentsCreated[1], /A label has been added for the latest release branch/);
  assert.ok(mockCore.infoMessages.some(m => m.includes("Release branch label 'release/v15' already present on issue #304")));
});

test('runTrackingIssue - handles malformed label structures gracefully', async () => {
  const mockCore = createMockCore();
  const labelsAdded = [];

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 65,
            body: '- Addresses: #305',
            labels: [null, undefined, { invalid: true }, { name: 'release/v15' }],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 305,
            state: 'open',
            labels: [null, { name: 'internal/user' }],
          },
        }),
        addLabels: async ({ labels }) => {
          labelsAdded.push(labels);
        },
        createComment: async () => {},
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(labelsAdded[0], ['internal/tracking']);
  assert.deepStrictEqual(labelsAdded[1], ['release/v15']);
});

test('runTriggerRcRelease - handles null and sparse entries in jobs and pull_requests safely', async () => {
  let dispatched = false;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.actions.listJobsForWorkflowRun) {
        return [null, undefined, { name: 'Run Tests', conclusion: 'success' }];
      }
      if (method === mockGithub.rest.repos.listTags) {
        return [{ name: 'v15.2.0-rc.1' }];
      }
      return [];
    },
    rest: {
      actions: {
        listJobsForWorkflowRun: async () => {},
        createWorkflowDispatch: async ({ workflow_id, ref, inputs }) => {
          dispatched = true;
          assert.strictEqual(workflow_id, 'rc-release.yml');
          assert.strictEqual(ref, 'main');
          assert.strictEqual(inputs.tag, 'v15.2.0-rc.2');
        }
      },
      repos: {
        listTags: async () => {},
        listPullRequestsAssociatedWithCommit: async () => {},
        getContent: async () => { throw new Error('Not found'); }
      }
    }
  };

  const mockCore = {
    info: () => {},
    warning: () => {},
    setFailed: () => {},
    setOutput: () => {}
  };

  const result = await runTriggerRcRelease({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
      payload: {
        workflow_run: {
          id: 77777,
          head_sha: 'sparse-sha',
          head_branch: 'release/v15',
          pull_requests: [null, undefined]
        }
      }
    },
    core: mockCore,
    process: { env: { BRANCH: 'release/v15', SHA: 'sparse-sha', WORKFLOW_RUN_ID: '77777' } }
  });

  assert.strictEqual(result, 'v15.2.0-rc.2');
  assert.strictEqual(dispatched, true);
});

test('runTrackingIssue - dispatches backport-issues.yml when adopting tracking issue', async () => {
  const mockCore = createMockCore();
  let dispatchedWorkflow = null;

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'main' }, { name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 10,
            title: 'Add feature',
            body: '- Addresses: #200',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 200,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async (params) => {
          dispatchedWorkflow = params;
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(dispatchedWorkflow, {
    owner: 'rancher',
    repo: 'terraform-provider-rancher2',
    workflow_id: 'backport-issues.yml',
    ref: 'main',
    inputs: {
      issue_number: '200',
      release_label: 'release/v15',
      pr_number: '10',
    },
  });
});

test('runTrackingIssue - warns and continues when workflow dispatch fails', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'main' }, { name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 11,
            title: 'Fix something',
            body: '- Addresses: #201',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 201,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {
          throw new Error('API Rate Limit or Network Error');
        },
      },
    },
  };

  await runTrackingIssue({ github: mockGithub, core: mockCore });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.warningMessages.some(m => m.includes("Failed to dispatch 'backport-issues.yml' for issue #201: API Rate Limit or Network Error")));
});

test('runTrackingIssue - uses GITHUB_MERGE_TOKEN and getOctokit when available', async () => {
  const mockCore = createMockCore();
  let tokenUsed = null;
  let customDispatchCalled = false;

  const mockProcess = {
    env: {
      GITHUB_MERGE_TOKEN: 'custom-vault-token-123',
    },
  };

  const mockGetOctokit = (token) => {
    tokenUsed = token;
    return {
      rest: {
        actions: {
          createWorkflowDispatch: async () => {
            customDispatchCalled = true;
          },
        },
      },
    };
  };

  const mockGithub = {
    paginate: async (method) => {
      if (method === mockGithub.rest.repos.listBranches) {
        return [{ name: 'main' }, { name: 'release/v15' }];
      }
      if (method === mockGithub.rest.search.issuesAndPullRequests) {
        return [
          {
            number: 12,
            title: 'Update provider',
            body: '- Addresses: #202',
            labels: [],
          },
        ];
      }
      return [];
    },
    rest: {
      repos: { listBranches: async () => {} },
      search: { issuesAndPullRequests: async () => {} },
      issues: {
        get: async () => ({
          data: {
            number: 202,
            state: 'open',
            labels: [],
          },
        }),
        addLabels: async () => {},
        createComment: async () => {},
      },
      actions: {
        createWorkflowDispatch: async () => {
          throw new Error('Should not be called when custom octokit is provided');
        },
      },
    },
  };

  await runTrackingIssue({
    github: mockGithub,
    core: mockCore,
    process: mockProcess,
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(tokenUsed, 'custom-vault-token-123');
  assert.strictEqual(customDispatchCalled, true);
});
