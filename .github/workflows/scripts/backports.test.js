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

test('runBackportPr - finds tracking issue and resolves via PR body addressed issue in O(1)', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'abc1234',
    },
  };

  let searchEndpointCalled = false;
  let subIssuesEndpointCalled = false;

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
    },
    request: async (endpoint, params) => {
      if (endpoint === 'GET /search/issues') {
        searchEndpointCalled = true;
        // Verify valid GitHub issue search query syntax
        assert.ok(params.q.includes('42'));
        assert.ok(params.q.includes('label:"internal/tracking"'));
        assert.strictEqual(params.q.includes('tracking-pr'), false);
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
        subIssuesEndpointCalled = true;
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
  assert.strictEqual(searchEndpointCalled, true);
  assert.strictEqual(subIssuesEndpointCalled, true);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Found tracking issue: #100')));
  assert.ok(mockCore.infoMessages.some(m => m.includes('No sub-issues found for issue #100. Exiting.')));
});

test('runBackportPr - disambiguates multiple search results using hidden comment header', async () => {
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
              // No Addresses in body, requiring comment inspection
              body: 'Generic PR description',
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
        return [{ user: { login: 'github-actions[bot]' }, body: 'Mentioned #45 here casually' }];
      }
      if (params.issue_number === 401) {
        return [{
          user: { login: 'github-actions[bot]' },
          body: '<!-- tracking-pr: #45 -->\nThis is the tracking issue for PR #45',
        }];
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
              { number: 401, body: 'Real tracking issue' },
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

test('runBackportPr - ignores search results without matching tracking header when multiple exist', async () => {
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
              body: 'PR description without Addresses',
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
            total_count: 2,
            items: [
              { number: 300, body: 'Unrelated issue 1' },
              { number: 301, body: 'Unrelated issue 2' },
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

test('runBackportPr - ignores tracking comment when author is not github-actions[bot]', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'botcheck1',
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
              body: 'PR description without Addresses',
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async () => [
      {
        user: { login: 'regular-user' },
        body: '<!-- tracking-pr: #44 -->\nThis looks like a tracking issue but is user spoofed',
      },
    ],
    request: async (endpoint) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 1,
            items: [
              { number: 500, body: 'User issue' },
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
              body: 'PR description without Addresses',
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

test('runBackportPr - resolves matching issue when PR description addresses multiple issues', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: 'multi123',
    },
  };

  let subIssuesEndpointCalled = false;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 47,
              title: 'Multi-issue fix',
              body: '<!-- comment -->\n- Addresses: #101\n- Addresses: #102',
              base: { ref: 'main' },
              merged_at: '2026-10-01T00:00:00Z',
            },
          ],
        }),
      },
    },
    request: async (endpoint, params) => {
      if (endpoint === 'GET /search/issues') {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 102,
                title: 'Second Issue',
                body: 'Bug description',
              },
            ],
          },
        };
      }
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        subIssuesEndpointCalled = true;
        assert.strictEqual(params.issue_number, 102);
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
  assert.strictEqual(subIssuesEndpointCalled, true);
  assert.ok(mockCore.infoMessages.some(m => m.includes('Found tracking issue: #102')));
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

test('trackingRegex in backport-issues - accurately extracts PR number from hidden comment header', () => {
  const trackingRegex = /<!--\s*tracking-pr:\s*#?(\d+)\s*-->/i;

  const standardHeader = '<!-- tracking-pr: #42 -->\nThis is the tracking issue for PR #42';
  const match1 = standardHeader.match(trackingRegex);
  assert.ok(match1);
  assert.strictEqual(Number(match1[1]), 42);

  const withoutHash = '<!-- tracking-pr: 99 -->\nThis is the tracking issue for PR #99';
  const match2 = withoutHash.match(trackingRegex);
  assert.ok(match2);
  assert.strictEqual(Number(match2[1]), 99);

  const extraSpaces = '<!--    tracking-pr:   #12345   -->';
  const match3 = extraSpaces.match(trackingRegex);
  assert.ok(match3);
  assert.strictEqual(Number(match3[1]), 12345);

  const userComment = 'Here is a regular user comment mentioning #42 casually.';
  assert.strictEqual(userComment.match(trackingRegex), null);
});
