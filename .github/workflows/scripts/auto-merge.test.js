import assert from 'node:assert/strict';
import test from 'node:test';
import {
  checkReviewerPermission,
  checkWorkflowModifications,
  evaluateReviews,
  resolveTargetPr,
  runAutoMerge,
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
          { number: 55, state: 'open', head: { ref: 'feature-branch' } },
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
    { user: { login: 'AuthorUser' }, state: 'APPROVED' },
  ];
  const res = evaluateReviews(reviews, 'authoruser');
  assert.strictEqual(res.hasChangesRequested, false);
  assert.deepStrictEqual(res.approvedReviewers, ['ReviewerOne']);
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
        // listFiles
        return [{ filename: 'rancher2/config.go' }];
      }
      // listReviews
      return [
        { user: { login: 'reviewer' }, state: 'CHANGES_REQUESTED' },
        { user: { login: 'approver' }, state: 'APPROVED' },
      ];
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
    process: { env: {} },
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'blocked_changes_requested');
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
    process: { env: {} },
  });
  assert.strictEqual(result, false);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'no_qualifying_permission');
});

test('runAutoMerge - successfully submits bot approval and enables auto-merge for triage approver', async () => {
  const mockCore = createMockCore();
  let callCount = 0;
  let botReviewCreated = false;
  let autoMergeMutationCalled = false;

  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'triage-engineer' }, state: 'APPROVED' }];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
      users: {
        getAuthenticated: async () => ({
          data: { login: 'rancher-bot' },
        }),
      },
      pulls: {
        createReview: async ({ event, commit_id }) => {
          assert.strictEqual(event, 'APPROVE');
          assert.strictEqual(commit_id, 'sha1');
          botReviewCreated = true;
          return { data: { id: 123 } };
        },
      },
    },
    graphql: async (query, vars) => {
      assert.ok(query.includes('enablePullRequestAutoMerge'));
      assert.strictEqual(vars.pullRequestId, 'PR_node_id_123');
      autoMergeMutationCalled = true;
      return {
        enablePullRequestAutoMerge: {
          pullRequest: { id: 'PR_node_id_123' },
        },
      };
    },
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
    process: { env: {} },
  });

  assert.strictEqual(result, true);
  assert.strictEqual(botReviewCreated, true);
  assert.strictEqual(autoMergeMutationCalled, true);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'enabled');
});

test('runAutoMerge - handles already enabled auto-merge in GraphQL gracefully', async () => {
  const mockCore = createMockCore();
  let callCount = 0;

  const mockGithub = {
    paginate: async () => {
      callCount++;
      if (callCount === 1) {
        return [{ filename: 'rancher2/config.go' }];
      }
      return [{ user: { login: 'triage-engineer' }, state: 'APPROVED' }];
    },
    rest: {
      repos: {
        getCollaboratorPermissionLevel: async () => ({
          data: { permission: 'triage', role_name: 'triage' },
        }),
      },
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
    process: { env: {} },
  });

  assert.strictEqual(result, true);
  assert.strictEqual(mockCore.outputs.auto_merge_status, 'already_enabled');
});
