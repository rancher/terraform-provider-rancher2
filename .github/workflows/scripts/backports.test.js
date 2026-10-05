import assert from 'node:assert/strict';
import test from 'node:test';
import { runBackportIssues, runBackportPr } from './backports.js';

function createMockCore() {
  return {
    infoMessages: [],
    warningMessages: [],
    errorMessages: [],
    failedMessage: null,
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
  };
}

test('runBackportPr - finds tracking issue when tracking text is in comments', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'abc1234',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 42,
              title: 'Fix issue',
              body: 'Fixes bug\n- Addresses: #100',
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async (method, params) => {
      if (method === mockGithub.rest.issues.listComments && params.issue_number === 100) {
        return [
          {
            body: 'This is the tracking issue for PR #42, a label has been added for the latest release branch',
          },
        ];
      }
      return [];
    },
    request: async (endpoint, params) => {
      if (endpoint === 'GET /search/issues') {
        assert.ok(params.q.includes('#42'));
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 100,
                title: 'User Bug Report',
                body: 'Bug description',
              },
            ],
          },
        };
      }
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        assert.strictEqual(params.issue_number, 100);
        return { data: [] }; // Empty sub-issues terminates execution cleanly
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Found tracking issue: #100')));
  assert.ok(mockCore.infoMessages.some(m => m.includes('No sub-issues found for issue #100. Exiting.')));
});

test('runBackportPr - finds tracking issue when tracking text is in body (legacy format)', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'def5678',
    },
  };

  let listCommentsCalled = false;
  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 43,
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {
          listCommentsCalled = true;
          return { data: [] };
        },
      },
    },
    paginate: async () => [],
    request: async (endpoint, params) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 200,
                body: 'This is the tracking issue for #43\n\nPlease add labels...',
              },
            ],
          },
        };
      }
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        assert.strictEqual(params.issue_number, 200);
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(listCommentsCalled, false);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Found tracking issue: #200')));
});

test('runBackportPr - ignores false positive issue when search returns single issue that merely references PR number', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'ghi9012',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 44,
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async () => [
      { body: 'Hey, is this related to #44 or is it independent?' },
    ],
    request: async (endpoint) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 300,
                title: 'Unrelated issue',
                body: 'Some description without tracking header',
              },
            ],
          },
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes("No verified 'internal/tracking' issue found for PR #44. Exiting.")));
});

test('runBackportPr - disambiguates correct tracking issue when multiple search results exist', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'jkl3456',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 45,
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async (method, params) => {
      if (params.issue_number === 400) {
        return [{ body: 'Mentioned #45 here casually' }];
      }
      if (params.issue_number === 401) {
        return [{ body: 'This is the tracking issue for PR #45, a label has been added' }];
      }
      return [];
    },
    request: async (endpoint, params) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 2,
            items: [
              { number: 400, body: 'Not the tracking issue' },
              { number: 401, body: 'Real tracking issue body' },
            ],
          },
        };
      }
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        assert.strictEqual(params.issue_number, 401);
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Found tracking issue: #401')));
});

test('runBackportPr - exits when search results are empty', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'mno7890',
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 46,
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
    },
    request: async (endpoint) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 0,
            items: [],
          },
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' } },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(mockCore.infoMessages.some(m => m.includes("No 'internal/tracking' issue found for PR #46. Exiting.")));
});

test('runBackportIssues - creates sub-issue and handles quoted PR numbers', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '"47"',
      TERRAFORM_MAINTAINERS: '["testmaintainer"]',
    },
  };

  let createdIssueParams = null;
  let subIssueLinkedParams = null;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 47);
          return {
            data: {
              number: 47,
              body: 'Fix description',
            },
          };
        },
        create: async (params) => {
          createdIssueParams = params;
          return {
            data: {
              id: 9999,
              number: 501,
            },
          };
        },
      },
    },
    request: async (endpoint, params) => {
      if (endpoint === 'POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        subIssueLinkedParams = params;
        return { data: {} };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      label: { name: 'release/v15' },
      issue: {
        number: 500,
        title: 'Parent Tracking Issue',
      },
    },
  };

  await runBackportIssues({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(createdIssueParams.title, '[release/v15] Parent Tracking Issue');
  assert.deepStrictEqual(createdIssueParams.labels, ['release/v15', 'internal/backport']);
  assert.deepStrictEqual(createdIssueParams.assignees, ['testmaintainer']);
  assert.match(createdIssueParams.body, /Backport #47 to release\/v15 for #500/);
  assert.strictEqual(subIssueLinkedParams.issue_number, 500);
  assert.strictEqual(subIssueLinkedParams.sub_issue_id, 9999);
});

test('runBackportIssues - throws on invalid PR number', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: 'not-a-number',
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      label: { name: 'release/v15' },
      issue: { number: 500, title: 'Title' },
    },
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: {},
        context: mockContext,
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Invalid PR number: not-a-number/ }
  );
});
