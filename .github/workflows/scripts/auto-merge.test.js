import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import defaultHandler, {
  checkReviewerPermission,
  checkWorkflowModifications,
  evaluateReviews,
  resolveApprovalSignal,
  resolveTargetPr,
  runAutoMerge,
  runPreflight,
  validateApprovalSignal,
} from './auto-merge.js';

function createMockCore() {
  return {
    infoMessages: [],
    warningMessages: [],
    errorMessages: [],
    failedMessage: null,
    outputs: {},
    info(msg) {
      this.infoMessages.push(msg);
    },
    warning(msg) {
      this.warningMessages.push(msg);
    },
    error(msg) {
      this.errorMessages.push(msg);
    },
    setFailed(msg) {
      this.failedMessage = msg;
    },
    setOutput(key, val) {
      this.outputs[key] = val;
    },
  };
}

test('resolveTargetPr - resolves directly from context.payload.pull_request', async () => {
  const mockCore = createMockCore();
  const pr = { number: 10, title: 'Test PR' };
  const result = await resolveTargetPr({
    context: { payload: { pull_request: pr } },
    core: mockCore,
  });
  assert.strictEqual(result.number, 10);
});

test('resolveTargetPr - resolves via PR_NUMBER env var', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, title: 'From API' },
        }),
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'test' } },
    core: mockCore,
    process: { env: { PR_NUMBER: '42' } },
  });
  assert.strictEqual(result.number, 42);
});

test('resolveTargetPr - resolves via workflow_run.pull_requests', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, title: 'From Workflow Run' },
        }),
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          pull_requests: [{ number: 99 }],
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result.number, 99);
});

test('resolveTargetPr - resolves via commit SHA fallback', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 55, state: 'open', head: { ref: 'feature-branch', sha: 'commit-sha-123' } },
        ],
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-123',
          head_branch: 'feature-branch',
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result.number, 55);
});

test('resolveTargetPr - rejects fallback PR when head SHA does not match triggering SHA', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 55, state: 'open', head: { ref: 'feature-branch', sha: 'different-newer-sha' } },
        ],
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-123',
          head_branch: 'feature-branch',
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result, null);
});

test('resolveTargetPr - resolves via signal parameter when workflow_run.pull_requests is empty (fork PR)', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, title: 'From Signal Parameter' },
        }),
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          id: 1234,
          pull_requests: [],
          head_sha: 'commit-sha',
        },
      },
    },
    core: mockCore,
    process: { env: { PR_NUMBER: '' } },
    signal: { pr_number: 2588 },
  });
  assert.strictEqual(result.number, 2588);
});

test('resolveTargetPr - resolves via approval signal file on disk when signal parameter is omitted', async () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pr-signal-test-'));
  const signalPath = path.join(tempDir, 'approval-signal.json');
  fs.writeFileSync(signalPath, JSON.stringify({ pr_number: 777 }));

  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, title: 'From Signal File' },
        }),
      },
    },
  };

  try {
    const result = await resolveTargetPr({
      github: mockGithub,
      context: {
        repo: { owner: 'rancher', repo: 'test' },
        payload: {
          workflow_run: {
            id: 1234,
            pull_requests: [],
            head_sha: 'commit-sha',
          },
        },
      },
      core: mockCore,
      process: { env: { PR_NUMBER: '', APPROVAL_SIGNAL_PATH: signalPath } },
    });
    assert.strictEqual(result.number, 777);
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
});

test('resolveTargetPr - rejects fallback candidate when head.sha is missing or mismatched', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 88, state: 'open', head: { ref: 'fork-branch' } }, // missing head.sha
        ],
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-456',
          head_branch: 'fork-branch',
          pull_requests: [],
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result, null);
});

test('resolveTargetPr - rejects fallback candidate when fetched PR branch or repository mismatches', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 88, state: 'open', head: { ref: 'fork-branch', sha: 'commit-sha-456' } },
        ],
      },
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, state: 'open', head: { ref: 'other-branch', sha: 'commit-sha-456' } },
        }),
      },
    },
  };
  const result = await resolveTargetPr({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-456',
          head_branch: 'fork-branch',
          pull_requests: [],
        },
      },
    },
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result, null);
});

test('resolveTargetPr - rejects fallback candidate when fetched PR lacks head.sha or pulls.get fails', async () => {
  const mockCoreMissingSha = createMockCore();
  const mockGithubMissingSha = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 88, state: 'open', head: { ref: 'fork-branch', sha: 'commit-sha-456' } },
        ],
      },
      pulls: {
        get: async ({ pull_number }) => ({
          data: { number: pull_number, state: 'open', head: { ref: 'fork-branch' } }, // lacks sha
        }),
      },
    },
  };
  const resultMissingSha = await resolveTargetPr({
    github: mockGithubMissingSha,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-456',
          head_branch: 'fork-branch',
          pull_requests: [],
        },
      },
    },
    core: mockCoreMissingSha,
    process: { env: {} },
  });
  assert.strictEqual(resultMissingSha, null);

  const mockCoreThrow = createMockCore();
  const mockGithubThrow = {
    paginate: async (fn, params) => fn(params),
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => [
          { number: 88, state: 'open', head: { ref: 'fork-branch', sha: 'commit-sha-456' } },
        ],
      },
      pulls: {
        get: async () => {
          throw new Error('API failure');
        },
      },
    },
  };
  const resultThrow = await resolveTargetPr({
    github: mockGithubThrow,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: {
          head_sha: 'commit-sha-456',
          head_branch: 'fork-branch',
          pull_requests: [],
        },
      },
    },
    core: mockCoreThrow,
    process: { env: {} },
  });
  assert.strictEqual(resultThrow, null);
});

test('checkWorkflowModifications - flags changes in .github/', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [
      { filename: 'rancher2/resource_cluster.go' },
      { filename: '.github/workflows/pull_request.yaml' },
    ],
  };
  const res = await checkWorkflowModifications({
    github: mockGithub,
    owner: 'rancher',
    repo: 'test',
    prNumber: 1,
    core: mockCore,
  });
  assert.strictEqual(res.modifiesWorkflows, true);
  assert.deepStrictEqual(res.workflowFiles, ['.github/workflows/pull_request.yaml']);
});

test('checkWorkflowModifications - flags changes when previous_filename is in .github/ (rename protection)', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [
      { filename: 'scripts/moved_workflow.sh', previous_filename: '.github/workflows/old.yml' },
    ],
  };
  const res = await checkWorkflowModifications({
    github: mockGithub,
    owner: 'rancher',
    repo: 'test',
    prNumber: 1,
    core: mockCore,
  });
  assert.strictEqual(res.modifiesWorkflows, true);
  assert.deepStrictEqual(res.workflowFiles, ['.github/workflows/old.yml']);
});

test('checkWorkflowModifications - allows PRs without changes in .github/', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [
      { filename: 'rancher2/resource_cluster.go' },
      { filename: 'docs/resources/cluster.md' },
    ],
  };
  const res = await checkWorkflowModifications({
    github: mockGithub,
    owner: 'rancher',
    repo: 'test',
    prNumber: 1,
    core: mockCore,
  });
  assert.strictEqual(res.modifiesWorkflows, false);
});

test('evaluateReviews - handles CHANGES_REQUESTED and author exclusions', () => {
  const reviews = [
    { user: { login: 'reviewer1' }, state: 'CHANGES_REQUESTED' },
    { user: { login: 'reviewer2' }, state: 'APPROVED' },
    { user: { login: 'author' }, state: 'APPROVED' },
  ];
  const res = evaluateReviews(reviews, 'author');
  assert.strictEqual(res.hasChangesRequested, true);
  assert.deepStrictEqual(res.approvedReviewers, ['reviewer2']);
});

test('evaluateReviews - supersedes older reviews with latest state', () => {
  const reviews = [
    { user: { login: 'reviewer1' }, state: 'CHANGES_REQUESTED' },
    { user: { login: 'reviewer1' }, state: 'APPROVED' },
  ];
  const res = evaluateReviews(reviews, 'author');
  assert.strictEqual(res.hasChangesRequested, false);
  assert.deepStrictEqual(res.approvedReviewers, ['reviewer1']);
});

test('evaluateReviews - subsequent COMMENTED review does not overwrite CHANGES_REQUESTED', () => {
  const reviews = [
    { user: { login: 'reviewer1' }, state: 'CHANGES_REQUESTED' },
    { user: { login: 'reviewer1' }, state: 'COMMENTED' },
  ];
  const res = evaluateReviews(reviews, 'author');
  assert.strictEqual(res.hasChangesRequested, true);
  assert.deepStrictEqual(res.approvedReviewers, []);
});

test('evaluateReviews - subsequent COMMENTED review does not overwrite APPROVED', () => {
  const reviews = [
    { user: { login: 'reviewer1' }, state: 'APPROVED' },
    { user: { login: 'reviewer1' }, state: 'COMMENTED' },
  ];
  const res = evaluateReviews(reviews, 'author');
  assert.strictEqual(res.hasChangesRequested, false);
  assert.deepStrictEqual(res.approvedReviewers, ['reviewer1']);
});

test('evaluateReviews - DISMISSED review clears prior CHANGES_REQUESTED and APPROVED states', () => {
  const reviewsChangesRequested = [
    { user: { login: 'reviewer1' }, state: 'CHANGES_REQUESTED' },
    { user: { login: 'reviewer1' }, state: 'DISMISSED' },
  ];
  const res1 = evaluateReviews(reviewsChangesRequested, 'author');
  assert.strictEqual(res1.hasChangesRequested, false);
  assert.deepStrictEqual(res1.approvedReviewers, []);

  const reviewsApproved = [
    { user: { login: 'reviewer2' }, state: 'APPROVED' },
    { user: { login: 'reviewer2' }, state: 'DISMISSED' },
  ];
  const res2 = evaluateReviews(reviewsApproved, 'author');
  assert.strictEqual(res2.hasChangesRequested, false);
  assert.deepStrictEqual(res2.approvedReviewers, []);
});

test('evaluateReviews - case-insensitive user logins and author comparison', () => {
  const reviews = [
    { user: { login: 'ReviewerOne' }, state: 'APPROVED' },
    { user: { login: 'Author' }, state: 'APPROVED' },
  ];
  const res = evaluateReviews(reviews, 'author');
  assert.strictEqual(res.hasChangesRequested, false);
  assert.deepStrictEqual(res.approvedReviewers, ['ReviewerOne']);
});

test('evaluateReviews - excludes authenticated bot from approvedReviewers', () => {
  const reviews = [
    { user: { login: 'rancher-bot' }, state: 'APPROVED' },
    { user: { login: 'triage-reviewer' }, state: 'APPROVED' },
  ];
  const res = evaluateReviews(reviews, 'author', 'rancher-bot');
  assert.strictEqual(res.hasChangesRequested, false);
  assert.deepStrictEqual(res.approvedReviewers, ['triage-reviewer']);
});

test('checkReviewerPermission - authorizes triage, write, and admin', async () => {
  const mockCore = createMockCore();
  for (const role of ['triage', 'write', 'admin']) {
    const mockGithub = {
      rest: {
        repos: {
          getCollaboratorPermissionLevel: async () => ({
            data: { permission: role, role_name: role },
          }),
        },
      },
    };
    const res = await checkReviewerPermission({
      github: mockGithub,
      owner: 'rancher',
      repo: 'test',
      username: 'user1',
      core: mockCore,
    });
    assert.strictEqual(res.isAuthorized, true);
  }

  const mockGithubRead = {
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'read', role_name: 'read' },
        }),
      },
    },
  };
  const resRead = await checkReviewerPermission({
    github: mockGithubRead,
    owner: 'rancher',
    repo: 'test',
    username: 'userRead',
    core: mockCore,
  });
  assert.strictEqual(resRead.isAuthorized, false);
});

test('runAutoMerge - skips non-open PR', async () => {
  const mockCore = createMockCore();
  const mockGithub = {};
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: { number: 1, state: 'closed', base: { ref: 'main' } },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_not_open');
});

test('runAutoMerge - skips draft PR', async () => {
  const mockCore = createMockCore();
  const mockGithub = {};
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: { number: 1, state: 'open', draft: true, base: { ref: 'main' } },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_draft');
});

test('runAutoMerge - skips PR targeting branch other than main', async () => {
  const mockCore = createMockCore();
  const mockGithub = {};
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: { number: 1, state: 'open', draft: false, base: { ref: 'release/v15' } },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_not_main');
});

test('runAutoMerge - skips superseded PR when PR_HEAD_SHA mismatches', async () => {
  const mockCore = createMockCore();
  const mockGithub = {};
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'new-sha' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { PR_HEAD_SHA: 'old-sha' } },
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_superseded');
});

test('runAutoMerge - skips PR modifying .github/**', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [{ filename: '.github/workflows/auto-merge.yml' }],
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_workflow_files');
});

test('runAutoMerge - fails when GITHUB_MERGE_TOKEN is missing', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [{ filename: 'rancher2/config.go' }],
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });
  assert.strictEqual(result, false);
  assert.match(mockCore.failedMessage, /GITHUB_MERGE_TOKEN is required/);
});

test('runAutoMerge - fails when getOctokit is unavailable', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
      },
    },
    paginate: async () => [{ filename: 'rancher2/config.go' }],
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
  });
  assert.strictEqual(result, false);
  assert.match(mockCore.failedMessage, /getOctokit is unavailable/);
});

test('runAutoMerge - blocks when CHANGES_REQUESTED is active', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
        listReviews: async () => {},
      },
    },
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [
        { user: { login: 'reviewer' }, state: 'CHANGES_REQUESTED' },
        { user: { login: 'approver' }, state: 'APPROVED' },
      ];
    },
  };
  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
    },
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => elevatedGithub,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'blocked_changes_requested');
});

test('runAutoMerge - does not allow bot approval to self-qualify without human approval', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
        listReviews: async () => {},
      },
    },
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'rancher-bot' }, state: 'APPROVED' }];
    },
  };
  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
    },
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => elevatedGithub,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'no_approvals');
});

test('runAutoMerge - skips when reviewer lacks qualifying role', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
        listReviews: async () => {},
      },
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'read', role_name: 'read' },
        }),
      },
    },
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'guest' }, state: 'APPROVED' }];
    },
  };
  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
    },
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => elevatedGithub,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'no_qualifying_permission');
});

test('runAutoMerge - skips when qualifying review commit_id does not match PR head', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    rest: {
      pulls: {
        listFiles: async () => {},
        listReviews: async () => {},
      },
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'triage-reviewer' }, state: 'APPROVED', commit_id: 'old-sha' }];
    },
  };
  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
    },
  };
  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'new-sha' },
        user: { login: 'author' },
      },
    },
  };
  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => elevatedGithub,
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'skipped_superseded');
});

test('runAutoMerge - successfully submits bot approval and enables auto-merge for triage approver', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  let botReviewCreated = false;
  let autoMergeMutationCalled = false;
  let elevatedClientUsedForReview = false;
  let elevatedClientUsedForMutation = false;

  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
      pulls: {
        createReview: async ({ event, commit_id, body }) => {
          assert.strictEqual(event, 'APPROVE');
          assert.strictEqual(commit_id, 'sha1');
          assert.match(body, /@triage-engineer/);
          botReviewCreated = true;
          elevatedClientUsedForReview = true;
          return { data: { id: 123 } };
        },
      },
    },
    graphql: async (query, vars) => {
      assert.ok(query.includes('enablePullRequestAutoMerge'));
      assert.strictEqual(vars.pullRequestId, 'PR_node_id_123');
      autoMergeMutationCalled = true;
      elevatedClientUsedForMutation = true;
      return {
        enablePullRequestAutoMerge: {
          pullRequest: { id: 'PR_node_id_123' },
        },
      };
    },
  };

  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'triage-engineer' }, state: 'APPROVED', commit_id: 'sha1' }];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const mockGetOctokit = (token) => {
    assert.strictEqual(token, 'vault-merge-token-secret');
    return elevatedGithub;
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        node_id: 'PR_node_id_123',
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };

  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'vault-merge-token-secret' } },
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(result, true);
  assert.strictEqual(botReviewCreated, true);
  assert.strictEqual(elevatedClientUsedForReview, true);
  assert.strictEqual(autoMergeMutationCalled, true);
  assert.strictEqual(elevatedClientUsedForMutation, true);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'enabled');
});

test('runAutoMerge - handles already enabled auto-merge in GraphQL gracefully', async () => {
  const mockCore = createMockCore();
  let callCount = 0;

  const elevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
      pulls: {
        createReview: async () => ({ data: { id: 123 } }),
      },
    },
    graphql: async () => {
      throw new Error('Auto merge is already enabled for this pull request.');
    },
  };

  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'triage-engineer' }, state: 'APPROVED', commit_id: 'sha1' }];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const mockGetOctokit = (token) => {
    assert.strictEqual(token, 'vault-merge-token-secret');
    return elevatedGithub;
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        node_id: 'PR_node_id_123',
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };

  const result = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'vault-merge-token-secret' } },
    getOctokit: mockGetOctokit,
  });

  assert.strictEqual(result, true);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'already_enabled');
});

test('runAutoMerge - fails closed when authenticated bot login cannot be resolved', async () => {
  const mockCore = createMockCore();

  // Test case 1: getAuthenticated throws an error
  const failingElevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => {
          throw new Error('Bad credentials / token expired');
        },
      },
    },
  };

  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'write', role_name: 'write' },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      pull_request: {
        number: 1,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'sha1' },
        user: { login: 'author' },
      },
    },
  };

  const result1 = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => failingElevatedGithub,
  });

  assert.strictEqual(result1, false);
  assert.match(mockCore.failedMessage, /Failed to resolve authenticated bot login/);

  // Test case 2: getAuthenticated returns empty login
  const mockCore2 = createMockCore();
  const emptyLoginElevatedGithub = {
    rest: {
      users: {
        getAuthenticated: async () => ({
          data: { login: '' },
        }),
      },
    },
  };

  const result2 = await runAutoMerge({
    github: mockGithub,
    context: mockContext,
    core: mockCore2,
    process: { env: { GITHUB_MERGE_TOKEN: 'mock-token' } },
    getOctokit: () => emptyLoginElevatedGithub,
  });

  assert.strictEqual(result2, false);
  assert.match(mockCore2.failedMessage, /Could not resolve authenticated bot identity from GITHUB_MERGE_TOKEN/);
});

test('runPreflight - succeeds for qualifying PR with valid approval signal', async () => {
  const mockCore = createMockCore();
  let callCount = 0;

  const mockGithub = {
    paginate: async (_fn, params) => {
      if (params?.run_id) {
        // Artifacts lookup
        return [{ id: 101, name: 'auto-merge-approval-signal' }];
      }
      callCount++;
      if (callCount === 1) {
        // Changed files lookup
        return [{ filename: 'rancher2/resource_cluster.go' }];
      }
      // Reviews lookup
      return [{ id: 'rev-1', user: { login: 'triage-reviewer' }, state: 'APPROVED', commit_id: 'commit-sha' }];
    },
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [{ id: 101, name: 'auto-merge-approval-signal' }] },
        }),
        downloadArtifact: async () => ({
          data: {
            pr_number: 42,
            review_id: 'rev-1',
            review_user: 'triage-reviewer',
            review_commit_id: 'commit-sha',
            head_sha: 'commit-sha',
            review_state: 'APPROVED',
            base_ref: 'main',
          },
        }),
      },
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      workflow_run: {
        id: 9999,
        pull_requests: [{ number: 42, head: { sha: 'commit-sha' } }],
      },
      pull_request: {
        number: 42,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'commit-sha' },
        user: { login: 'author' },
      },
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { PR_NUMBER: '42', PR_HEAD_SHA: 'commit-sha' } },
  });

  assert.strictEqual(result, true);
  assert.strictEqual(mockCore.outputs.eligible, 'true');
  assert.strictEqual(mockCore.outputs.pr_number, '42');
  assert.strictEqual(mockCore.outputs.pr_head_sha, 'commit-sha');
  assert.strictEqual(mockCore.outputs.review_id, 'rev-1');
  assert.strictEqual(mockCore.outputs.review_user, 'triage-reviewer');
  assert.strictEqual(mockCore.outputs.review_commit_id, 'commit-sha');
});

test('runPreflight - fails when auto-merge-approval-signal artifact is missing in workflow_run', async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    paginate: async () => {
      return [{ name: 'other-unrelated-artifact' }];
    },
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [{ name: 'other-unrelated-artifact' }] },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      workflow_run: {
        id: 8888,
      },
      pull_request: {
        number: 42,
        state: 'open',
        draft: false,
        base: { ref: 'main' },
        head: { sha: 'commit-sha' },
        user: { login: 'author' },
      },
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes("lacks 'auto-merge-approval-signal' artifact")));
});

test('runPreflight - fails when target PR cannot be resolved', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      pulls: {
        get: async () => {
          throw new Error('Not found');
        },
      },
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {},
    },
    core: mockCore,
    process: { env: { PR_NUMBER: '99999' } },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.strictEqual(mockCore.failedMessage, 'Could not resolve target PR for preflight auto-merge.');
});

test('runPreflight - succeeds for fork PR when workflow_run.pull_requests is empty by reading PR from signal artifact', async () => {
  const mockCore = createMockCore();
  let callCount = 0;

  const mockGithub = {
    paginate: async (_fn, params) => {
      if (params?.run_id) {
        return [{ id: 101, name: 'auto-merge-approval-signal' }];
      }
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/resource_cluster.go' }];
      }
      return [{ id: 'rev-fork-1', user: { login: 'triage-reviewer' }, state: 'APPROVED', commit_id: 'fork-sha' }];
    },
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [{ id: 101, name: 'auto-merge-approval-signal' }] },
        }),
        downloadArtifact: async () => ({
          data: {
            pr_number: 2588,
            review_id: 'rev-fork-1',
            review_user: 'triage-reviewer',
            review_commit_id: 'fork-sha',
            head_sha: 'fork-sha',
            review_state: 'APPROVED',
            base_ref: 'main',
          },
        }),
      },
      pulls: {
        get: async ({ pull_number }) => ({
          data: {
            number: pull_number,
            state: 'open',
            draft: false,
            base: { ref: 'main' },
            head: { sha: 'fork-sha', ref: 'trivial-docs' },
            user: { login: 'fork-author' },
          },
        }),
      },
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'test' },
    payload: {
      workflow_run: {
        id: 9999,
        pull_requests: [], // Empty for fork PRs!
        head_sha: 'fork-sha',
        head_branch: 'trivial-docs',
      },
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: { PR_NUMBER: '', PR_HEAD_SHA: 'fork-sha', WORKFLOW_RUN_SHA: 'fork-sha' } },
  });

  assert.strictEqual(result, true);
  assert.strictEqual(mockCore.outputs.eligible, 'true');
  assert.strictEqual(mockCore.outputs.pr_number, '2588');
  assert.strictEqual(mockCore.outputs.pr_head_sha, 'fork-sha');
  assert.strictEqual(mockCore.outputs.review_id, 'rev-fork-1');
  assert.strictEqual(mockCore.outputs.review_user, 'triage-reviewer');
});

test('runPreflight - skips non-open or draft PR', async () => {
  const mockCoreClosed = createMockCore();
  const mockGithub = {
    paginate: async () => [],
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [] },
        }),
      },
    },
  };

  const resultClosed = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 10,
          state: 'closed',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha1' },
        },
      },
    },
    core: mockCoreClosed,
  });

  assert.strictEqual(resultClosed, false);
  assert.strictEqual(mockCoreClosed.outputs.eligible, 'false');

  const mockCoreDraft = createMockCore();
  const resultDraft = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 11,
          state: 'open',
          draft: true,
          base: { ref: 'main' },
          head: { sha: 'sha1' },
        },
      },
    },
    core: mockCoreDraft,
  });

  assert.strictEqual(resultDraft, false);
  assert.strictEqual(mockCoreDraft.outputs.eligible, 'false');
});

test('runPreflight - skips PR modifying .github/**', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async () => [{ filename: '.github/workflows/ci.yml' }],
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 12,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha1' },
        },
      },
    },
    core: mockCore,
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes('modifies files in .github/')));
});

test('runPreflight - skips PR with active CHANGES_REQUESTED', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [
        { user: { login: 'reviewer1' }, state: 'CHANGES_REQUESTED' },
        { user: { login: 'reviewer2' }, state: 'APPROVED' },
      ];
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 13,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha1' },
        },
      },
    },
    core: mockCore,
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes('active CHANGES_REQUESTED')));
});

test('runPreflight - skips PR when reviewer lacks qualifying role', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'read-only-user' }, state: 'APPROVED', commit_id: 'sha1' }];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'read', role_name: 'read' },
        }),
      },
    },
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 14,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha1' },
        },
      },
    },
    core: mockCore,
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes('does not have an approved review from a user with triage, write, or admin')));
});

test('default export - routes preflight mode and execute mode', async () => {
  const mockCorePreflight = createMockCore();
  const mockGithub = {
    paginate: async () => [{ filename: 'rancher2/config.go' }],
  };

  const preflightRes = await defaultHandler({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 20,
          state: 'closed',
          draft: false,
          base: { ref: 'main' },
        },
      },
    },
    core: mockCorePreflight,
    process: { env: { SCRIPT_MODE: 'preflight' } },
  });

  assert.strictEqual(preflightRes, false);
  assert.strictEqual(mockCorePreflight.outputs.eligible, 'false');

  const mockCoreExecute = createMockCore();
  const executeRes = await defaultHandler({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 21,
          state: 'closed',
          draft: false,
          base: { ref: 'main' },
        },
      },
    },
    core: mockCoreExecute,
    process: { env: { SCRIPT_MODE: 'execute' } },
  });

  assert.strictEqual(executeRes, false);
  assert.strictEqual(mockCoreExecute.outputs.auto_merge_status, 'skipped_not_open');
});

test('resolveApprovalSignal - resolves from pull_request_review event payload', async () => {
  const mockCore = createMockCore();
  const signal = await resolveApprovalSignal({
    context: {
      payload: {
        pull_request: { number: 10, head: { sha: 'sha-head' }, base: { ref: 'main' } },
        review: { id: 777, state: 'approved', user: { login: 'reviewer-alice' }, commit_id: 'sha-head' },
      },
    },
    core: mockCore,
  });

  assert.deepStrictEqual(signal, {
    pr_number: 10,
    review_id: 777,
    review_state: 'approved',
    review_user: 'reviewer-alice',
    review_commit_id: 'sha-head',
    head_sha: 'sha-head',
    base_ref: 'main',
  });
});

test('resolveApprovalSignal - resolves from APPROVAL_SIGNAL env var', async () => {
  const mockCore = createMockCore();
  const rawSignal = {
    pr_number: 15,
    review_id: 'rev-88',
    review_user: 'bob',
    review_commit_id: 'sha-15',
    head_sha: 'sha-15',
    review_state: 'APPROVED',
    base_ref: 'main',
  };

  const signal = await resolveApprovalSignal({
    core: mockCore,
    process: { env: { APPROVAL_SIGNAL: JSON.stringify(rawSignal) } },
  });

  assert.deepStrictEqual(signal, rawSignal);
});

test('resolveApprovalSignal - resolves from signal file on disk', async () => {
  const mockCore = createMockCore();
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'signal-test-'));
  const signalPath = path.join(tmpDir, 'approval-signal.json');
  const expectedSignal = {
    pr_number: 22,
    review_id: 'rev-99',
    review_user: 'carol',
    review_commit_id: 'sha-22',
    head_sha: 'sha-22',
    review_state: 'APPROVED',
    base_ref: 'main',
  };
  fs.writeFileSync(signalPath, JSON.stringify(expectedSignal));

  const signal = await resolveApprovalSignal({
    core: mockCore,
    process: { env: { APPROVAL_SIGNAL_PATH: signalPath } },
  });

  assert.deepStrictEqual(signal, expectedSignal);
  fs.rmSync(tmpDir, { recursive: true, force: true });
});

test('resolveApprovalSignal - resolves from TRIGGER_REVIEW_* env vars', async () => {
  const mockCore = createMockCore();
  const signal = await resolveApprovalSignal({
    core: mockCore,
    process: {
      env: {
        PR_NUMBER: '33',
        PR_HEAD_SHA: 'sha-33',
        TRIGGER_REVIEW_ID: 'rev-33',
        TRIGGER_REVIEW_USER: 'dan',
        TRIGGER_REVIEW_COMMIT_ID: 'sha-33',
      },
    },
  });

  assert.strictEqual(signal.pr_number, 33);
  assert.strictEqual(signal.review_id, 'rev-33');
  assert.strictEqual(signal.review_user, 'dan');
  assert.strictEqual(signal.review_commit_id, 'sha-33');
});

test('resolveApprovalSignal - resolves from downloadArtifact API', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [{ id: 456, name: 'auto-merge-approval-signal' }] },
        }),
        downloadArtifact: async () => ({
          data: {
            pr_number: 55,
            review_id: 'rev-55',
            review_user: 'erin',
            review_commit_id: 'sha-55',
            head_sha: 'sha-55',
            review_state: 'APPROVED',
            base_ref: 'main',
          },
        }),
      },
    },
  };

  const signal = await resolveApprovalSignal({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: { workflow_run: { id: 1234 } },
    },
    core: mockCore,
  });

  assert.strictEqual(signal.review_user, 'erin');
  assert.strictEqual(signal.review_id, 'rev-55');
});

test('validateApprovalSignal - validates structure and rejects invalid signals', () => {
  const mockCore = createMockCore();
  const validSignal = {
    pr_number: 100,
    review_id: 1,
    review_user: 'user1',
    review_commit_id: 'sha-100',
    head_sha: 'sha-100',
    review_state: 'APPROVED',
    base_ref: 'main',
  };
  const pr = { number: 100, head: { sha: 'sha-100' } };

  const validRes = validateApprovalSignal(validSignal, pr, mockCore);
  assert.strictEqual(validRes.isValid, true);
  assert.strictEqual(validRes.signal.review_user, 'user1');

  // Rejects missing fields
  const missingRes = validateApprovalSignal({ ...validSignal, review_user: '' }, pr, mockCore);
  assert.strictEqual(missingRes.isValid, false);
  assert.strictEqual(missingRes.reason, 'incomplete_signal');

  // Rejects non-APPROVED state
  const stateRes = validateApprovalSignal({ ...validSignal, review_state: 'CHANGES_REQUESTED' }, pr, mockCore);
  assert.strictEqual(stateRes.isValid, false);
  assert.strictEqual(stateRes.reason, 'not_approved_state');

  // Rejects base branch other than main
  const baseRes = validateApprovalSignal({ ...validSignal, base_ref: 'release/v1' }, pr, mockCore);
  assert.strictEqual(baseRes.isValid, false);
  assert.strictEqual(baseRes.reason, 'not_main_base');

  // Rejects mismatched PR number
  const prRes = validateApprovalSignal({ ...validSignal, pr_number: 999 }, pr, mockCore);
  assert.strictEqual(prRes.isValid, false);
  assert.strictEqual(prRes.reason, 'pr_mismatch');

  // Rejects stale commit SHA
  const staleRes = validateApprovalSignal({ ...validSignal, review_commit_id: 'old-sha' }, pr, mockCore);
  assert.strictEqual(staleRes.isValid, false);
  assert.strictEqual(staleRes.reason, 'stale_signal');
});

test('runPreflight - rejects when triggering review commit does not match PR head', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async () => {
      // Changed files then reviews
      return [
        { id: 'rev-stale', user: { login: 'triage-user' }, state: 'APPROVED', commit_id: 'old-commit' },
      ];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const signal = {
    pr_number: 70,
    review_id: 'rev-stale',
    review_user: 'triage-user',
    review_commit_id: 'sha-current',
    head_sha: 'sha-current',
    review_state: 'APPROVED',
    base_ref: 'main',
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 70,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha-current' },
          user: { login: 'author-user' },
        },
      },
    },
    core: mockCore,
    process: { env: { APPROVAL_SIGNAL: JSON.stringify(signal) } },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes('Skipping stale review')));
});

test('runPreflight - rejects when triggering review is from PR author', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async () => [],
  };

  const signal = {
    pr_number: 80,
    review_id: 'rev-author',
    review_user: 'author-user',
    review_commit_id: 'sha-80',
    head_sha: 'sha-80',
    review_state: 'APPROVED',
    base_ref: 'main',
  };

  const result = await runPreflight({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        pull_request: {
          number: 80,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha-80' },
          user: { login: 'author-user' },
        },
      },
    },
    core: mockCore,
    process: { env: { APPROVAL_SIGNAL: JSON.stringify(signal) } },
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.eligible, 'false');
  assert.ok(mockCore.infoMessages.some(m => m.includes('Self-approvals are not permitted')));
});

test('runAutoMerge - binds to triggering review signal and rejects read-only user even with existing write approval on PR', async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    paginate: async () => {
      return [
        { id: 'rev-write', user: { login: 'write-user' }, state: 'APPROVED', commit_id: 'sha-90' },
        { id: 'rev-read', user: { login: 'read-contributor' }, state: 'APPROVED', commit_id: 'sha-90' },
      ];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async ({ username }) => {
          if (username === 'read-contributor') {
            return { data: { permission: 'read', role_name: 'read' } };
          }
          return { data: { permission: 'write', role_name: 'write' } };
        },
      },
    },
  };

  // Signal was produced by read-contributor's review
  const signal = {
    pr_number: 90,
    review_id: 'rev-read',
    review_user: 'read-contributor',
    review_commit_id: 'sha-90',
    head_sha: 'sha-90',
    review_state: 'APPROVED',
    base_ref: 'main',
  };

  const mockElevatedGithub = {
    rest: {
      users: { getAuthenticated: async () => ({ data: { login: 'rancher-bot' } }) },
    },
  };

  const result = await runAutoMerge({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: { id: 9991 },
        pull_request: {
          number: 90,
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha-90' },
          user: { login: 'author-user' },
        },
      },
    },
    core: mockCore,
    process: {
      env: {
        GITHUB_MERGE_TOKEN: 'mock-token',
        APPROVAL_SIGNAL: JSON.stringify(signal),
      },
    },
    getOctokit: () => mockElevatedGithub,
  });

  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'no_qualifying_permission');
  assert.ok(mockCore.infoMessages.some(m => m.includes('does not have triage, write, or admin permissions')));
});

test('runAutoMerge - binds to qualifying triage signal and succeeds', async () => {
  const mockCore = createMockCore();
  let reviewCreated = false;
  let autoMergeEnabled = false;

  const mockGithub = {
    paginate: async () => [
      { id: 'rev-triage', user: { login: 'triage-member' }, state: 'APPROVED', commit_id: 'sha-95' },
    ],
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
    },
  };

  const signal = {
    pr_number: 95,
    review_id: 'rev-triage',
    review_user: 'triage-member',
    review_commit_id: 'sha-95',
    head_sha: 'sha-95',
    review_state: 'APPROVED',
    base_ref: 'main',
  };

  const mockElevatedGithub = {
    rest: {
      users: { getAuthenticated: async () => ({ data: { login: 'rancher-bot' } }) },
      pulls: {
        createReview: async ({ body }) => {
          assert.match(body, /Automated approval triggered by @triage-member \(role: triage\)/);
          reviewCreated = true;
          return { data: { id: 888 } };
        },
      },
    },
    graphql: async (query) => {
      assert.match(query, /enablePullRequestAutoMerge/);
      autoMergeEnabled = true;
      return { enablePullRequestAutoMerge: {} };
    },
  };

  const result = await runAutoMerge({
    github: mockGithub,
    context: {
      repo: { owner: 'rancher', repo: 'test' },
      payload: {
        workflow_run: { id: 9992 },
        pull_request: {
          number: 95,
          node_id: 'PR_node_95',
          state: 'open',
          draft: false,
          base: { ref: 'main' },
          head: { sha: 'sha-95' },
          user: { login: 'author-user' },
        },
      },
    },
    core: mockCore,
    process: {
      env: {
        GITHUB_MERGE_TOKEN: 'mock-token',
        APPROVAL_SIGNAL: JSON.stringify(signal),
      },
    },
    getOctokit: () => mockElevatedGithub,
  });

  assert.strictEqual(result, true);
  assert.strictEqual(reviewCreated, true);
  assert.strictEqual(autoMergeEnabled, true);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'enabled');
});
