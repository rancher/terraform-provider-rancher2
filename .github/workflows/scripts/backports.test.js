import assert from 'node:assert/strict';
import test from 'node:test';
import { runBackportIssues, runBackportPr, runMergeLabel } from './backports.js';

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

test('runBackportIssues - handles workflow_dispatch inputs via env and fetches parent issue', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '48',
      ISSUE_NUMBER: '600',
      RELEASE_LABEL: 'release/v15',
      TERRAFORM_MAINTAINERS: '["testmaintainer"]',
    },
  };

  let createdIssueParams = null;
  let subIssueLinkedParams = null;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 48) {
            return {
              data: {
                number: 48,
                title: 'Add support for feature X',
                body: 'Detailed PR description for feature X',
              },
            };
          }
          if (issue_number === 600) {
            return {
              data: {
                number: 600,
                title: 'Parent Issue From Env',
              },
            };
          }
          throw new Error(`Unexpected issue number: ${issue_number}`);
        },
        create: async (params) => {
          createdIssueParams = params;
          return {
            data: {
              id: 8888,
              number: 601,
            },
          };
        },
      },
    },
    request: async (endpoint, params) => {
      if (endpoint === 'POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        subIssueLinkedParams = params;
        return { status: 201 };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {},
  };

  await runBackportIssues({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(createdIssueParams.title, '[release/v15] Parent Issue From Env');
  assert.deepStrictEqual(createdIssueParams.labels, ['release/v15', 'internal/backport']);
  assert.deepStrictEqual(createdIssueParams.assignees, ['testmaintainer']);
  assert.match(createdIssueParams.body, /Backport #48 to release\/v15 for #600/);
  assert.strictEqual(subIssueLinkedParams.issue_number, 600);
  assert.strictEqual(subIssueLinkedParams.sub_issue_id, 8888);
});

test('runBackportIssues - throws when release label is missing', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '49',
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      issue: { number: 700, title: 'Title' },
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
    { message: /Release label must be provided via env \(RELEASE_LABEL\) or label event payload/ }
  );
});

test('runBackportIssues - throws when parent issue retrieval fails', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '49',
      ISSUE_NUMBER: '999',
      RELEASE_LABEL: 'release/v15',
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 999) {
            throw new Error('Not Found');
          }
          return { data: { number: 49, title: 'PR' } };
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {},
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: mockGithub,
        context: mockContext,
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Failed to retrieve parent issue #999: Not Found/ }
  );
});

test('runBackportIssues - skips creation when sub-issue for release branch already exists (idempotent)', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '48',
      ISSUE_NUMBER: '600',
      RELEASE_LABEL: 'release/v15',
    },
  };

  let createIssueCalled = false;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 600) {
            return {
              data: {
                number: 600,
                title: 'Existing Parent Issue',
              },
            };
          }
          throw new Error(`Unexpected issue number: ${issue_number}`);
        },
        create: async () => {
          createIssueCalled = true;
          return { data: { id: 9999, number: 601 } };
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        return {
          data: [
            {
              number: 601,
              title: '[release/v15] Existing Parent Issue',
              labels: [{ name: 'release/v15' }, { name: 'internal/backport' }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {},
  };

  await runBackportIssues({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(createIssueCalled, false);
  assert.ok(mockCore.infoMessages.some(m => m.includes("Sub-issue for release label 'release/v15' already exists on tracking issue #600. Skipping creation.")));
});

test('runBackportIssues - sanitizes leading "#" from ISSUE_NUMBER and PR env inputs', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '#48',
      ISSUE_NUMBER: '#600',
      RELEASE_LABEL: 'release/v15',
    },
  };

  let fetchedIssueNumber = null;
  let fetchedPrNumber = null;
  let createdIssueParams = null;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 600) {
            fetchedIssueNumber = issue_number;
            return { data: { number: 600, title: 'Parent Issue With Hash' } };
          }
          if (issue_number === 48) {
            fetchedPrNumber = issue_number;
            return { data: { number: 48, title: 'PR With Hash', body: 'Body' } };
          }
          throw new Error(`Unexpected issue number: ${issue_number}`);
        },
        create: async (params) => {
          createdIssueParams = params;
          return { data: { id: 7777, number: 605 } };
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === 'GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        return { data: [] };
      }
      if (endpoint === 'POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues') {
        return { status: 201 };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {},
  };

  await runBackportIssues({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(fetchedIssueNumber, 600);
  assert.strictEqual(fetchedPrNumber, 48);
  assert.strictEqual(createdIssueParams.title, '[release/v15] Parent Issue With Hash');
});

test('trackingRegex in backport-issues - accurately extracts PR number from hidden comment header', () => {
  const trackingRegex = /<!--\s*tracking-pr:\s*#?(\d+)\s*-->/i;

  const standardHeader = '<!-- tracking-pr: #42 -->\nThis is the tracking issue for PR #42';
  const match1 = standardHeader.match(trackingRegex);
  assert.ok(match1);
  assert.strictEqual(parseInt(match1[1], 10), 42);

  const withoutHash = '<!-- tracking-pr: 99 -->\nThis is the tracking issue for PR #99';
  const match2 = withoutHash.match(trackingRegex);
  assert.ok(match2);
  assert.strictEqual(parseInt(match2[1], 10), 99);

  const extraSpaces = '<!--    tracking-pr:   #12345   -->';
  const match3 = extraSpaces.match(trackingRegex);
  assert.ok(match3);
  assert.strictEqual(parseInt(match3[1], 10), 12345);

  const userComment = 'Here is a regular user comment mentioning #42 casually.';
  assert.strictEqual(userComment.match(trackingRegex), null);
});

test('runMergeLabel - adds internal/merged label to issues referenced in PR payload', async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 123);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #123 on release branch'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - resolves PR from SHA via listPullRequestsAssociatedWithCommit in workflow_run context', async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.commit_sha, 'commit-sha-789');
      return [
        {
          number: 555,
          state: 'closed',
          merged_at: '2026-10-06T12:00:00Z',
          body: 'Resolves #321'
        }
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 321);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commit-sha-789'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: {} } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 321);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - skips unmerged PR when resolving via SHA', async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.commit_sha, 'unmerged-sha');
      return [
        {
          number: 888,
          state: 'open',
          merged_at: null,
          body: 'Addresses #999'
        }
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      issues: {
        addLabels: async () => {
          addLabelsCalled = true;
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'unmerged-sha'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: {} } });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(core.infoMessages.some(msg => msg.includes('No pull request found for merge-label; skipping.')));
});

test('runMergeLabel - fails and continues when issue fetch throws error', async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          throw new Error('Not found (404)');
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #404'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.ok(core.failedMessage !== null && core.failedMessage.includes('Could not process issue #404'));
});

test('runMergeLabel - ignores issue references inside HTML comments', async () => {
  const addedLabels = [];
  const fetchedIssues = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          fetchedIssues.push(issue_number);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #123 on release branch <!-- tracking-pr: #999 -->'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(fetchedIssues, [123]);
  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - ignores issue references inside fenced code blocks', async () => {
  const addedLabels = [];
  const fetchedIssues = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          fetchedIssues.push(issue_number);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #123 on release branch\n```\n#999 should be ignored\n```\n~~~bash\n#888 also ignored\n~~~'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(fetchedIssues, [123]);
  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - skips adding internal/merged if already present on issue', async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: {
            pull_request: null,
            labels: [{ name: 'internal/backport' }, { name: 'internal/merged' }]
          }
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #123 on release branch'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(core.infoMessages.some(msg => msg.includes("already has 'internal/merged' label; skipping.")));
});

test('runMergeLabel - prioritizes PR matching target branch when multiple PRs associated with commit SHA', async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      {
        number: 100,
        state: 'closed',
        merged_at: '2026-10-06T12:00:00Z',
        base: { ref: 'main' },
        body: 'Resolves #111'
      },
      {
        number: 200,
        state: 'closed',
        merged_at: '2026-10-06T12:05:00Z',
        base: { ref: 'release/v15' },
        body: 'Resolves #222'
      }
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 222);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'shared-sha-123',
        head_branch: 'release/v15'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 222);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - falls back to fetching PR by number from workflow_run.pull_requests when commit lookup fails', async () => {
  const addedLabels = [];
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
          assert.strictEqual(pull_number, 888);
          return {
            data: {
              number: 888,
              state: 'closed',
              merged_at: '2026-10-06T12:00:00Z',
              body: 'Resolves #999'
            }
          };
        }
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 999);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commit-fallback-sha',
        head_branch: 'release/v15',
        pull_requests: [
          {
            number: 888,
            base: { ref: 'release/v15' }
          }
        ]
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 999);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - normalizes branch prefix (refs/heads/) when matching associated PR', async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      {
        number: 300,
        state: 'closed',
        merged_at: '2026-10-06T12:00:00Z',
        base: { ref: 'release/v15' },
        body: 'Resolves #444'
      }
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 444);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commit-prefixed-branch',
        head_branch: 'refs/heads/release/v15'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'refs/heads/release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 444);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - checks subsequent pull_requests candidates if earlier ones are unmerged', async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {}
      },
      pulls: {
        get: async ({ pull_number }) => {
          if (pull_number === 101) {
            return {
              data: {
                number: 101,
                state: 'open',
                merged_at: null,
                body: 'Fixes #555'
              }
            };
          }
          if (pull_number === 102) {
            return {
              data: {
                number: 102,
                state: 'closed',
                merged_at: '2026-10-06T12:00:00Z',
                body: 'Fixes #666'
              }
            };
          }
          throw new Error('Unexpected PR number');
        }
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 666);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commit-multi-pr',
        head_branch: 'release/v15',
        pull_requests: [
          { number: 101, base: { ref: 'release/v15' } },
          { number: 102, base: { ref: 'release/v15' } }
        ]
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 666);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - logs info and skips when issue lacks internal/backport label', async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: {
            pull_request: null,
            labels: [{ name: 'enhancement' }]
          }
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Fixes #789 on release branch'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(core.infoMessages.some(msg => msg.includes("does not have 'internal/backport' label; skipping.")));
});

test('runMergeLabel - falls back to direct listPullRequestsAssociatedWithCommit when github.paginate is undefined', async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async ({ commit_sha }) => {
          assert.strictEqual(commit_sha, 'sha-no-paginate-backport');
          return {
            data: [
              {
                number: 777,
                state: 'closed',
                merged_at: '2026-10-06T12:00:00Z',
                base: { ref: 'release/v15' },
                body: 'Resolves #888'
              }
            ]
          };
        }
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 888);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'sha-no-paginate-backport',
        head_branch: 'release/v15'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 888);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - normalizes candidate base.ref prefix when matching PR from workflow_run.pull_requests', async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 555);
          return {
            data: {
              number: 555,
              state: 'closed',
              merged_at: '2026-10-06T12:00:00Z',
              base: { ref: 'refs/heads/release/v15' },
              body: 'Fixes #777'
            }
          };
        }
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 777);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_branch: 'release/v15',
        pull_requests: [
          { number: 555, base: { ref: 'refs/heads/release/v15' } }
        ]
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { BRANCH: 'refs/heads/release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 777);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - matches both issue URL and shorthand #issue with word boundaries', async () => {
  const addedLabels = [];
  const fetchedIssues = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          fetchedIssues.push(issue_number);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Addresses: #101\nAlso resolves https://github.com/rancher/terraform-provider-rancher2/issues/202 and ignores token#303 and https://github.com/other/repo/issues/404'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(fetchedIssues.sort((a, b) => a - b), [101, 202]);
  assert.strictEqual(addedLabels.length, 2);
  assert.deepStrictEqual(addedLabels.map(a => a.issue_number).sort((a, b) => a - b), [101, 202]);
});

test('runMergeLabel - skips and logs when referenced item is a pull request', async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            pull_request: { url: `https://api.github.com/repos/rancher/terraform-provider-rancher2/pulls/${issue_number}` },
            labels: [{ name: 'internal/backport' }]
          }
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      pull_request: {
        number: 456,
        body: 'Addresses: #101'
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(core.infoMessages.some(msg => msg.includes('#101 is a pull request, not an issue; skipping.')));
});

test('runMergeLabel - uses candidate.body directly from workflow_run.pull_requests without pulls.get API call', async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 505);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_branch: 'release/v15',
        pull_requests: [
          {
            number: 999,
            state: 'closed',
            merged_at: '2026-10-06T12:00:00Z',
            base: { ref: 'release/v15' },
            body: 'Fixes #505'
          }
        ]
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: {} } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 505);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});

test('runMergeLabel - handles null and sparse entries in pulls and pull_requests safely', async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      null,
      undefined,
      {
        number: 888,
        state: 'closed',
        merged_at: '2026-10-06T12:00:00Z',
        base: { ref: 'release/v15' },
        body: 'Closes #707'
      }
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => {}
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 707);
          return {
            data: {
              pull_request: null,
              labels: [{ name: 'internal/backport' }]
            }
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        }
      }
    }
  };

  const core = createMockCore();
  const context = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {
      workflow_run: {
        head_sha: 'commit-sparse-sha',
        head_branch: 'release/v15',
        pull_requests: [null, undefined]
      }
    }
  };

  await runMergeLabel({ github: mockGithub, context, core, process: { env: { SHA: 'commit-sparse-sha', BRANCH: 'release/v15' } } });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 707);
  assert.deepStrictEqual(addedLabels[0].labels, ['internal/merged']);
});




test('runBackportIssues - throws when parent issue is a pull request', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '48',
      ISSUE_NUMBER: '500',
      RELEASE_LABEL: 'release/v15',
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 500) {
            return {
              data: {
                number: 500,
                title: 'This is actually a pull request',
                pull_request: { url: 'https://api.github.com/repos/rancher/terraform-provider-rancher2/pulls/500' },
              },
            };
          }
          return { data: { number: issue_number, title: 'PR' } };
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' },
    payload: {},
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: mockGithub,
        context: mockContext,
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Issue #500 is a pull request, not an issue/ }
  );
});

test('runBackportIssues - throws when issue number is non-positive', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '48',
      ISSUE_NUMBER: '0',
      RELEASE_LABEL: 'release/v15',
    },
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: {},
        context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }, payload: {} },
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Invalid issue number: 0/ }
  );
});

test('runBackportIssues - throws when PR number is non-positive', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: '-5',
      ISSUE_NUMBER: '600',
      RELEASE_LABEL: 'release/v15',
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: { number: 600, title: 'Valid Parent Issue' },
        }),
      },
    },
    request: async () => ({ data: [] }),
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: mockGithub,
        context: { repo: { owner: 'rancher', repo: 'terraform-provider-rancher2' }, payload: {} },
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Invalid PR number: -5/ }
  );
});

test('extract-pr workflow logic - validates explicit pr_number inputs and tracking headers correctly', () => {
  const parsePrInput = (rawInput) => {
    const rawPrInput = String(rawInput ?? '').trim();
    if (rawPrInput !== '') {
      const cleanedPr = rawPrInput.replace(/^#/, '');
      const parsedPr = parseInt(cleanedPr, 10);
      if (!isNaN(parsedPr) && parsedPr > 0) {
        return { success: true, pr: parsedPr };
      }
      return { success: false, error: `Invalid pr_number input provided: ${rawInput}` };
    }
    return { success: true, pr: null };
  };

  assert.deepStrictEqual(parsePrInput('48'), { success: true, pr: 48 });
  assert.deepStrictEqual(parsePrInput('#48'), { success: true, pr: 48 });
  assert.deepStrictEqual(parsePrInput('  #102  '), { success: true, pr: 102 });
  assert.deepStrictEqual(parsePrInput(''), { success: true, pr: null });
  assert.deepStrictEqual(parsePrInput('   '), { success: true, pr: null });
  assert.deepStrictEqual(parsePrInput(undefined), { success: true, pr: null });

  const invalidAlpha = parsePrInput('abc');
  assert.strictEqual(invalidAlpha.success, false);
  assert.match(invalidAlpha.error, /Invalid pr_number input provided: abc/);

  const invalidNegative = parsePrInput('-10');
  assert.strictEqual(invalidNegative.success, false);
  assert.match(invalidNegative.error, /Invalid pr_number input provided: -10/);

  const invalidZero = parsePrInput('0');
  assert.strictEqual(invalidZero.success, false);
  assert.match(invalidZero.error, /Invalid pr_number input provided: 0/);
});
