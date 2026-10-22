import assert from 'node:assert/strict';
import test from 'node:test';
import defaultHandler, {
  checkReviewerPermission,
  checkWorkflowModifications,
  evaluateReviews,
  resolveTargetPr,
  runAutoMerge,
  runPreflight,
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
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        // Artifacts lookup
        return [{ name: 'auto-merge-approval-signal' }];
      }
      if (callCount === 2) {
        // Changed files lookup
        return [{ filename: 'rancher2/resource_cluster.go' }];
      }
      // Reviews lookup
      return [{ user: { login: 'triage-reviewer' }, state: 'APPROVED', commit_id: 'commit-sha' }];
    },
    rest: {
      actions: {
        listWorkflowRunArtifacts: async () => ({
          data: { artifacts: [{ name: 'auto-merge-approval-signal' }] },
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
