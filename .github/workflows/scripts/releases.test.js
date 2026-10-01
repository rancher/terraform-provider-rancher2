import assert from 'node:assert/strict';
import test from 'node:test';
import defaultHandler, { computeNextRcTag, runCheckMaintainer, runRcNotify, runTriggerRcRelease } from './releases.js';

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

test('runCheckMaintainer - authorizes bot accounts ending in [bot]', async () => {
  const mockCore = { info: () => {} };
  const mockContext = { actor: 'rancher-eio[bot]' };
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


