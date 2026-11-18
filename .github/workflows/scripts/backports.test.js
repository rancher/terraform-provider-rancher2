import assert from "node:assert/strict";
import test from "node:test";
import {
  runBackportIssues,
  runBackportPr,
  runMergeLabel,
  runCloseTrackingIssue,
  runReleaseDone,
  transitionProjectItemsStatus,
  updateIssueProjectStatus,
} from "./backports.js";

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

test("runBackportPr - finds tracking issue and resolves via PR body addressed issue in O(1)", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "abc1234",
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
              title: "Fix issue",
              body: "Fixes bug\n- Addresses: #100",
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
            },
          ],
        }),
      },
    },
    request: async (endpoint, params) => {
      if (endpoint === "GET /search/issues") {
        searchEndpointCalled = true;
        // Verify valid GitHub issue search query syntax
        assert.ok(params.q.includes("42"));
        assert.ok(params.q.includes('label:"internal/tracking"'));
        assert.strictEqual(params.q.includes("tracking-pr"), false);
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 100,
                title: "User Bug Report",
                body: "Bug description",
              },
            ],
          },
        };
      }
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        subIssuesEndpointCalled = true;
        assert.strictEqual(params.issue_number, 100);
        return { data: [] }; // Empty sub-issues terminates execution cleanly
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(searchEndpointCalled, true);
  assert.strictEqual(subIssuesEndpointCalled, true);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found tracking issue: #100")),
  );
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("No sub-issues found for issue #100. Exiting."),
    ),
  );
});

test("runBackportPr - disambiguates multiple search results using hidden comment header", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "jkl3456",
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 45,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              // No Addresses in body, requiring comment inspection
              body: "Generic PR description",
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
        return [
          {
            user: { login: "github-actions[bot]" },
            body: "Mentioned #45 here casually",
          },
        ];
      }
      if (params.issue_number === 401) {
        return [
          {
            user: { login: "github-actions[bot]" },
            body: "<!-- tracking-pr: #45 -->\nThis is the tracking issue for PR #45",
          },
        ];
      }
      return [];
    },
    request: async (endpoint, params) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 2,
            items: [
              { number: 400, body: "Not the tracking issue" },
              { number: 401, body: "Real tracking issue" },
            ],
          },
        };
      }
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        assert.strictEqual(params.issue_number, 401);
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found tracking issue: #401")),
  );
});

test("runBackportPr - ignores search results without matching tracking header when multiple exist", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "ghi9012",
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 44,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR description without Addresses",
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async () => [
      { body: "Hey, is this related to #44 or is it independent?" },
    ],
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 2,
            items: [
              { number: 300, body: "Unrelated issue 1" },
              { number: 301, body: "Unrelated issue 2" },
            ],
          },
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "No verified 'internal/tracking' issue found for PR #44. Exiting.",
      ),
    ),
  );
});

test("runBackportPr - ignores tracking comment when author is not github-actions[bot]", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "botcheck1",
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 44,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR description without Addresses",
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
        user: { login: "regular-user" },
        body: "<!-- tracking-pr: #44 -->\nThis looks like a tracking issue but is user spoofed",
      },
    ],
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [{ number: 500, body: "User issue" }],
          },
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "No verified 'internal/tracking' issue found for PR #44. Exiting.",
      ),
    ),
  );
});

test("runBackportPr - ignores false positive issue when search returns single issue that merely references PR number", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "ghi9012",
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 44,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR description without Addresses",
            },
          ],
        }),
      },
      issues: {
        listComments: async () => {},
      },
    },
    paginate: async () => [
      { body: "Hey, is this related to #44 or is it independent?" },
    ],
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 300,
                title: "Unrelated issue",
                body: "Some description without tracking header",
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
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "No verified 'internal/tracking' issue found for PR #44. Exiting.",
      ),
    ),
  );
});

test("runBackportPr - resolves matching issue when PR description addresses multiple issues", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "multi123",
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
              title: "Multi-issue fix",
              body: "<!-- comment -->\n- Addresses: #101\n- Addresses: #102",
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
            },
          ],
        }),
      },
    },
    request: async (endpoint, params) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 102,
                title: "Second Issue",
                body: "Bug description",
              },
            ],
          },
        };
      }
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        subIssuesEndpointCalled = true;
        assert.strictEqual(params.issue_number, 102);
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(subIssuesEndpointCalled, true);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found tracking issue: #102")),
  );
});

test("runBackportPr - exits when search results are empty", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      MERGE_COMMIT_SHA: "mno7890",
    },
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 46,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
            },
          ],
        }),
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
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
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("No 'internal/tracking' issue found for PR #46. Exiting."),
    ),
  );
});

test("runBackportIssues - creates sub-issue and handles quoted PR numbers", async () => {
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
              body: "Fix description",
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
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return { data: [] };
      }
      if (
        endpoint ===
        "POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        subIssueLinkedParams = params;
        return { data: {} };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      label: { name: "release/v15" },
      issue: {
        number: 500,
        title: "Parent Tracking Issue",
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
  assert.strictEqual(
    createdIssueParams.title,
    "[release/v15] Parent Tracking Issue",
  );
  assert.deepStrictEqual(createdIssueParams.labels, [
    "release/v15",
    "internal/backport",
  ]);
  assert.strictEqual(createdIssueParams.assignees, undefined);
  assert.match(
    createdIssueParams.body,
    /Backport #47 to release\/v15 for #500/,
  );
  assert.strictEqual(subIssueLinkedParams.issue_number, 500);
  assert.strictEqual(subIssueLinkedParams.sub_issue_id, 9999);
});

test("runBackportIssues - creates backport issue without assignees", async () => {
  const mockCore = createMockCore();
  let createdIssueParams = null;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 500) {
            return {
              data: {
                id: 100,
                number: 500,
                title: "Parent Tracking Issue",
                pull_request: null,
              },
            };
          }
          if (issue_number === 47) {
            return {
              data: {
                number: 47,
                body: "PR body content",
              },
            };
          }
          throw new Error("Not found");
        },
        create: async (params) => {
          createdIssueParams = params;
          return { data: { id: 9999, number: 601 } };
        },
      },
    },
    request: async () => ({
      data: {
        id: 9999,
        sub_issues_summary: { total: 0 },
      },
    }),
  };

  await runBackportIssues({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {
        label: { name: "release/v15" },
        issue: { number: 500 },
      },
    },
    core: mockCore,
    process: {
      env: {
        PR: "47",
        TERRAFORM_MAINTAINERS: '["maintainer1", "maintainer2"]',
      },
    },
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(createdIssueParams.assignees, undefined);
});

test("runBackportIssues - throws on invalid PR number", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "not-a-number",
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      label: { name: "release/v15" },
      issue: { number: 500, title: "Title" },
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
    { message: /Invalid PR number: not-a-number/ },
  );
});

test("runBackportIssues - handles workflow_dispatch inputs via env and fetches parent issue", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "600",
      RELEASE_LABEL: "release/v15",
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
                title: "Add support for feature X",
                body: "Detailed PR description for feature X",
              },
            };
          }
          if (issue_number === 600) {
            return {
              data: {
                number: 600,
                title: "Parent Issue From Env",
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
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return { data: [] };
      }
      if (
        endpoint ===
        "POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        subIssueLinkedParams = params;
        return { status: 201 };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {},
  };

  await runBackportIssues({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: mockProcess,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(
    createdIssueParams.title,
    "[release/v15] Parent Issue From Env",
  );
  assert.deepStrictEqual(createdIssueParams.labels, [
    "release/v15",
    "internal/backport",
  ]);
  assert.strictEqual(createdIssueParams.assignees, undefined);
  assert.match(
    createdIssueParams.body,
    /Backport #48 to release\/v15 for #600/,
  );
  assert.strictEqual(subIssueLinkedParams.issue_number, 600);
  assert.strictEqual(subIssueLinkedParams.sub_issue_id, 8888);
});

test("runBackportIssues - throws when release label is missing", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "49",
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      issue: { number: 700, title: "Title" },
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
    {
      message:
        /Release label must be provided via env \(RELEASE_LABEL\) or label event payload/,
    },
  );
});

test("runBackportIssues - throws when parent issue retrieval fails", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "49",
      ISSUE_NUMBER: "999",
      RELEASE_LABEL: "release/v15",
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 999) {
            throw new Error("Not Found");
          }
          return { data: { number: 49, title: "PR" } };
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
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
    { message: /Failed to retrieve parent issue #999: Not Found/ },
  );
});

test("runBackportIssues - skips creation when sub-issue for release branch already exists (idempotent)", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "600",
      RELEASE_LABEL: "release/v15",
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
                title: "Existing Parent Issue",
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
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return {
          data: [
            {
              number: 601,
              title: "[release/v15] Existing Parent Issue",
              labels: [{ name: "release/v15" }, { name: "internal/backport" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
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
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes(
        "Sub-issue for release label 'release/v15' already exists on tracking issue #600. Skipping creation.",
      ),
    ),
  );
});

test('runBackportIssues - sanitizes leading "#" from ISSUE_NUMBER and PR env inputs', async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "#48",
      ISSUE_NUMBER: "#600",
      RELEASE_LABEL: "release/v15",
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
            return { data: { number: 600, title: "Parent Issue With Hash" } };
          }
          if (issue_number === 48) {
            fetchedPrNumber = issue_number;
            return {
              data: { number: 48, title: "PR With Hash", body: "Body" },
            };
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
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return { data: [] };
      }
      if (
        endpoint ===
        "POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return { status: 201 };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
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
  assert.strictEqual(
    createdIssueParams.title,
    "[release/v15] Parent Issue With Hash",
  );
});

test("trackingRegex in backport-issues - accurately extracts PR number from hidden comment header", () => {
  const trackingRegex = /<!--\s*tracking-pr:\s*#?(\d+)\s*-->/i;

  const standardHeader =
    "<!-- tracking-pr: #42 -->\nThis is the tracking issue for PR #42";
  const match1 = standardHeader.match(trackingRegex);
  assert.ok(match1);
  assert.strictEqual(parseInt(match1[1], 10), 42);

  const withoutHash =
    "<!-- tracking-pr: 99 -->\nThis is the tracking issue for PR #99";
  const match2 = withoutHash.match(trackingRegex);
  assert.ok(match2);
  assert.strictEqual(parseInt(match2[1], 10), 99);

  const extraSpaces = "<!--    tracking-pr:   #12345   -->";
  const match3 = extraSpaces.match(trackingRegex);
  assert.ok(match3);
  assert.strictEqual(parseInt(match3[1], 10), 12345);

  const userComment = "Here is a regular user comment mentioning #42 casually.";
  assert.strictEqual(userComment.match(trackingRegex), null);
});

test("runMergeLabel - adds internal/merged label to issues referenced in PR payload", async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 123);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #123 on release branch",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - resolves PR from SHA via listPullRequestsAssociatedWithCommit in workflow_run context", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.commit_sha, "commit-sha-789");
      return [
        {
          number: 555,
          state: "closed",
          merged_at: "2026-10-06T12:00:00Z",
          body: "Resolves #321",
        },
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 321);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-sha-789",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: {} },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 321);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - skips unmerged PR when resolving via SHA", async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    paginate: async (fn, params) => {
      assert.strictEqual(params.commit_sha, "unmerged-sha");
      return [
        {
          number: 888,
          state: "open",
          merged_at: null,
          body: "Addresses #999",
        },
      ];
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      issues: {
        addLabels: async () => {
          addLabelsCalled = true;
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "unmerged-sha",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: {} },
  });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(
    core.infoMessages.some((msg) =>
      msg.includes("No pull request found for merge-label; skipping."),
    ),
  );
});

test("runMergeLabel - fails and continues when issue fetch throws error", async () => {
  const mockGithub = {
    rest: {
      issues: {
        get: async () => {
          throw new Error("Not found (404)");
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #404",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.ok(
    core.failedMessage !== null &&
      core.failedMessage.includes("Could not process issue #404"),
  );
});

test("runMergeLabel - ignores issue references inside HTML comments", async () => {
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
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #123 on release branch <!-- tracking-pr: #999 -->",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(fetchedIssues, [123]);
  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - ignores issue references inside fenced code blocks", async () => {
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
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #123 on release branch\n```\n#999 should be ignored\n```\n~~~bash\n#888 also ignored\n~~~",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(fetchedIssues, [123]);
  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 123);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - skips adding internal/merged if already present on issue", async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: {
            pull_request: null,
            labels: [
              { name: "internal/backport" },
              { name: "internal/merged" },
            ],
          },
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #123 on release branch",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(
    core.infoMessages.some((msg) =>
      msg.includes("already has 'internal/merged' label; skipping."),
    ),
  );
});

test("runMergeLabel - prioritizes PR matching target branch when multiple PRs associated with commit SHA", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      {
        number: 100,
        state: "closed",
        merged_at: "2026-10-06T12:00:00Z",
        base: { ref: "main" },
        body: "Resolves #111",
      },
      {
        number: 200,
        state: "closed",
        merged_at: "2026-10-06T12:05:00Z",
        base: { ref: "release/v15" },
        body: "Resolves #222",
      },
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 222);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "shared-sha-123",
        head_branch: "release/v15",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 222);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - falls back to fetching PR by number from workflow_run.pull_requests when commit lookup fails", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => {
      throw new Error("API Rate Limit (403)");
    },
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 888);
          return {
            data: {
              number: 888,
              state: "closed",
              merged_at: "2026-10-06T12:00:00Z",
              body: "Resolves #999",
            },
          };
        },
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 999);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-fallback-sha",
        head_branch: "release/v15",
        pull_requests: [
          {
            number: 888,
            base: { ref: "release/v15" },
          },
        ],
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 999);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - normalizes branch prefix (refs/heads/) when matching associated PR", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      {
        number: 300,
        state: "closed",
        merged_at: "2026-10-06T12:00:00Z",
        base: { ref: "release/v15" },
        body: "Resolves #444",
      },
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 444);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-prefixed-branch",
        head_branch: "refs/heads/release/v15",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "refs/heads/release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 444);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - checks subsequent pull_requests candidates if earlier ones are unmerged", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: () => {},
      },
      pulls: {
        get: async ({ pull_number }) => {
          if (pull_number === 101) {
            return {
              data: {
                number: 101,
                state: "open",
                merged_at: null,
                body: "Fixes #555",
              },
            };
          }
          if (pull_number === 102) {
            return {
              data: {
                number: 102,
                state: "closed",
                merged_at: "2026-10-06T12:00:00Z",
                body: "Fixes #666",
              },
            };
          }
          throw new Error("Unexpected PR number");
        },
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 666);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-multi-pr",
        head_branch: "release/v15",
        pull_requests: [
          { number: 101, base: { ref: "release/v15" } },
          { number: 102, base: { ref: "release/v15" } },
        ],
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 666);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - logs info and skips when issue lacks internal/backport label", async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: {
            pull_request: null,
            labels: [{ name: "enhancement" }],
          },
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Fixes #789 on release branch",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(
    core.infoMessages.some((msg) =>
      msg.includes("does not have 'internal/backport' label; skipping."),
    ),
  );
});

test("runMergeLabel - falls back to direct listPullRequestsAssociatedWithCommit when github.paginate is undefined", async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async ({ commit_sha }) => {
          assert.strictEqual(commit_sha, "sha-no-paginate-backport");
          return {
            data: [
              {
                number: 777,
                state: "closed",
                merged_at: "2026-10-06T12:00:00Z",
                base: { ref: "release/v15" },
                body: "Resolves #888",
              },
            ],
          };
        },
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 888);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "sha-no-paginate-backport",
        head_branch: "release/v15",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 888);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - normalizes candidate base.ref prefix when matching PR from workflow_run.pull_requests", async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          assert.strictEqual(pull_number, 555);
          return {
            data: {
              number: 555,
              state: "closed",
              merged_at: "2026-10-06T12:00:00Z",
              base: { ref: "refs/heads/release/v15" },
              body: "Fixes #777",
            },
          };
        },
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 777);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_branch: "release/v15",
        pull_requests: [
          { number: 555, base: { ref: "refs/heads/release/v15" } },
        ],
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { BRANCH: "refs/heads/release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 777);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - matches both issue URL and shorthand #issue with word boundaries", async () => {
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
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Addresses: #101\nAlso resolves https://github.com/rancher/terraform-provider-rancher2/issues/202 and ignores token#303 and https://github.com/other/repo/issues/404",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.deepStrictEqual(
    fetchedIssues.sort((a, b) => a - b),
    [101, 202],
  );
  assert.strictEqual(addedLabels.length, 2);
  assert.deepStrictEqual(
    addedLabels.map((a) => a.issue_number).sort((a, b) => a - b),
    [101, 202],
  );
});

test("runMergeLabel - skips and logs when referenced item is a pull request", async () => {
  let addLabelsCalled = false;
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            pull_request: {
              url: `https://api.github.com/repos/rancher/terraform-provider-rancher2/pulls/${issue_number}`,
            },
            labels: [{ name: "internal/backport" }],
          },
        }),
        addLabels: async () => {
          addLabelsCalled = true;
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 456,
        body: "Addresses: #101",
      },
    },
  };

  await runMergeLabel({ github: mockGithub, context, core });

  assert.strictEqual(addLabelsCalled, false);
  assert.ok(
    core.infoMessages.some((msg) =>
      msg.includes("#101 is a pull request, not an issue; skipping."),
    ),
  );
});

test("runMergeLabel - uses candidate.body directly from workflow_run.pull_requests without pulls.get API call", async () => {
  const addedLabels = [];
  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 505);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_branch: "release/v15",
        pull_requests: [
          {
            number: 999,
            state: "closed",
            merged_at: "2026-10-06T12:00:00Z",
            base: { ref: "release/v15" },
            body: "Fixes #505",
          },
        ],
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: {} },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 505);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runMergeLabel - handles null and sparse entries in pulls and pull_requests safely", async () => {
  const addedLabels = [];
  const mockGithub = {
    paginate: async () => [
      null,
      undefined,
      {
        number: 888,
        state: "closed",
        merged_at: "2026-10-06T12:00:00Z",
        base: { ref: "release/v15" },
        body: "Closes #707",
      },
    ],
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => {},
      },
      issues: {
        get: async ({ issue_number }) => {
          assert.strictEqual(issue_number, 707);
          return {
            data: {
              pull_request: null,
              labels: [{ name: "internal/backport" }],
            },
          };
        },
        addLabels: async (params) => {
          addedLabels.push(params);
        },
      },
    },
  };

  const core = createMockCore();
  const context = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      workflow_run: {
        head_sha: "commit-sparse-sha",
        head_branch: "release/v15",
        pull_requests: [null, undefined],
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context,
    core,
    process: { env: { SHA: "commit-sparse-sha", BRANCH: "release/v15" } },
  });

  assert.strictEqual(addedLabels.length, 1);
  assert.strictEqual(addedLabels[0].issue_number, 707);
  assert.deepStrictEqual(addedLabels[0].labels, ["internal/merged"]);
});

test("runBackportIssues - throws when parent issue is a pull request", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "500",
      RELEASE_LABEL: "release/v15",
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
                title: "This is actually a pull request",
                pull_request: {
                  url: "https://api.github.com/repos/rancher/terraform-provider-rancher2/pulls/500",
                },
              },
            };
          }
          return { data: { number: issue_number, title: "PR" } };
        },
      },
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
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
    { message: /Issue #500 is a pull request, not an issue/ },
  );
});

test("runBackportIssues - throws when issue number is non-positive", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "0",
      RELEASE_LABEL: "release/v15",
    },
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: {},
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: {},
        },
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Invalid issue number: 0/ },
  );
});

test("runBackportIssues - throws when PR number is non-positive", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "-5",
      ISSUE_NUMBER: "600",
      RELEASE_LABEL: "release/v15",
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: { number: 600, title: "Valid Parent Issue" },
        }),
      },
    },
    request: async () => ({ data: [] }),
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: mockGithub,
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: {},
        },
        core: mockCore,
        process: mockProcess,
      });
    },
    { message: /Invalid PR number: -5/ },
  );
});

test("runBackportIssues - throws when release label format is invalid", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "600",
      RELEASE_LABEL: "invalid-label",
    },
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: {},
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: {},
        },
        core: mockCore,
        process: mockProcess,
      });
    },
    {
      message:
        /Invalid release label: "invalid-label". It must match the pattern release\/v<major>/,
    },
  );
});

test("runBackportIssues - throws when issue number is malformed string or decimal", async () => {
  const mockCore = createMockCore();

  for (const malformed of ["600abc", "1.5", "not-a-number"]) {
    await assert.rejects(
      async () => {
        await runBackportIssues({
          github: {},
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
            payload: {},
          },
          core: mockCore,
          process: {
            env: {
              PR: "48",
              ISSUE_NUMBER: malformed,
              RELEASE_LABEL: "release/v15",
            },
          },
        });
      },
      { message: new RegExp(`Invalid issue number: ${malformed}`) },
    );
  }
});

test("runBackportIssues - throws when PR number is malformed string or decimal", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: { number: 600, title: "Valid Parent Issue" },
        }),
      },
    },
    request: async () => ({ data: [] }),
  };

  for (const malformed of ["48abc", "1.5"]) {
    await assert.rejects(
      async () => {
        await runBackportIssues({
          github: mockGithub,
          context: {
            repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
            payload: {},
          },
          core: mockCore,
          process: {
            env: {
              PR: malformed,
              ISSUE_NUMBER: "600",
              RELEASE_LABEL: "release/v15",
            },
          },
        });
      },
      { message: new RegExp(`Invalid PR number: ${malformed}`) },
    );
  }
});

test("runBackportIssues - throws when sub-issues lookup throws an API error", async () => {
  const mockCore = createMockCore();
  const mockProcess = {
    env: {
      PR: "48",
      ISSUE_NUMBER: "600",
      RELEASE_LABEL: "release/v15",
    },
  };

  const mockGithub = {
    rest: {
      issues: {
        get: async () => ({
          data: { number: 600, title: "Valid Parent Issue" },
        }),
      },
    },
    request: async (endpoint) => {
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        throw new Error("API rate limit exceeded");
      }
      return { data: [] };
    },
  };

  await assert.rejects(
    async () => {
      await runBackportIssues({
        github: mockGithub,
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: {},
        },
        core: mockCore,
        process: mockProcess,
      });
    },
    {
      message:
        /Failed to check existing sub-issues for #600: API rate limit exceeded/,
    },
  );
});

test("extract-pr workflow logic - validates explicit pr_number inputs and tracking headers correctly", () => {
  const parsePrInput = (rawInput) => {
    const rawPrInput = String(rawInput ?? "").trim();
    if (rawPrInput !== "") {
      const cleanedPr = rawPrInput.replace(/^#/, "");
      if (/^[1-9]\d*$/.test(cleanedPr)) {
        return { success: true, pr: parseInt(cleanedPr, 10) };
      }
      return {
        success: false,
        error: `Invalid pr_number input provided: ${rawInput}`,
      };
    }
    return { success: true, pr: null };
  };

  assert.deepStrictEqual(parsePrInput("48"), { success: true, pr: 48 });
  assert.deepStrictEqual(parsePrInput("#48"), { success: true, pr: 48 });
  assert.deepStrictEqual(parsePrInput("  #102  "), { success: true, pr: 102 });
  assert.deepStrictEqual(parsePrInput(""), { success: true, pr: null });
  assert.deepStrictEqual(parsePrInput("   "), { success: true, pr: null });
  assert.deepStrictEqual(parsePrInput(undefined), { success: true, pr: null });

  const invalidAlpha = parsePrInput("abc");
  assert.strictEqual(invalidAlpha.success, false);
  assert.match(invalidAlpha.error, /Invalid pr_number input provided: abc/);

  const invalidAlphaSuffix = parsePrInput("#48abc");
  assert.strictEqual(invalidAlphaSuffix.success, false);
  assert.match(
    invalidAlphaSuffix.error,
    /Invalid pr_number input provided: #48abc/,
  );

  const invalidFloat = parsePrInput("1.5");
  assert.strictEqual(invalidFloat.success, false);
  assert.match(invalidFloat.error, /Invalid pr_number input provided: 1.5/);

  const invalidNegative = parsePrInput("-10");
  assert.strictEqual(invalidNegative.success, false);
  assert.match(invalidNegative.error, /Invalid pr_number input provided: -10/);

  const invalidZero = parsePrInput("0");
  assert.strictEqual(invalidZero.success, false);
  assert.match(invalidZero.error, /Invalid pr_number input provided: 0/);

  const parseIssueInput = ({ issueNumberPayload, issueNumberInput }) => {
    const rawInputIssue = issueNumberInput
      ? String(issueNumberInput).trim().replace(/^#/, "")
      : undefined;
    const rawIssueNumber =
      issueNumberPayload !== undefined
        ? String(issueNumberPayload)
        : rawInputIssue;
    if (!rawIssueNumber || !/^[1-9]\d*$/.test(rawIssueNumber)) {
      return {
        success: false,
        error:
          "Could not determine valid issue number from event payload or workflow inputs.",
      };
    }
    return { success: true, issueNumber: parseInt(rawIssueNumber, 10) };
  };

  assert.deepStrictEqual(parseIssueInput({ issueNumberPayload: 500 }), {
    success: true,
    issueNumber: 500,
  });
  assert.deepStrictEqual(parseIssueInput({ issueNumberInput: "600" }), {
    success: true,
    issueNumber: 600,
  });
  assert.deepStrictEqual(parseIssueInput({ issueNumberInput: "#600" }), {
    success: true,
    issueNumber: 600,
  });
  assert.strictEqual(
    parseIssueInput({ issueNumberInput: "#600abc" }).success,
    false,
  );
  assert.strictEqual(
    parseIssueInput({ issueNumberInput: "1.5" }).success,
    false,
  );
  assert.strictEqual(parseIssueInput({ issueNumberInput: "0" }).success, false);
  assert.strictEqual(parseIssueInput({}).success, false);
});

test("runBackportPr - adopts referenced issue on demand when search returns 0 results", async () => {
  const mockCore = createMockCore();
  let labelsAdded = false;
  let subIssuesFetched = false;
  let gitCommands = [];

  const mockGitExec = async (args) => {
    gitCommands.push(args.join(" "));
    return { stdout: "", stderr: "" };
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 47,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "Fix bug\n\n- Addresses: #701",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({
          data: { number: 471, node_id: "PR_node_471" },
        }),
      },
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            title: "Bug report 701",
            state: "open",
            labels: [{ name: "release/v15" }],
          },
        }),
        addLabels: async () => {
          labelsAdded = true;
        },
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      if (endpoint.includes("/sub_issues")) {
        subIssuesFetched = true;
        return {
          data: [
            {
              number: 702,
              title: "[release/v15] Bug report 701",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "abc1234" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(labelsAdded, true);
  assert.strictEqual(subIssuesFetched, true);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found 1 sub-issues.")),
  );
});

test("runBackportPr - generates sub-issue on demand when tracking issue has 0 sub-issues", async () => {
  const mockCore = createMockCore();
  let subIssueCreated = false;
  let subIssueCallCount = 0;
  let gitCommands = [];

  const mockGitExec = async (args) => {
    gitCommands.push(args.join(" "));
    return { stdout: "", stderr: "" };
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 48,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "Feature PR\n\n- Addresses: #801",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({
          data: { number: 481, node_id: "PR_node_481" },
        }),
      },
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            title: "Feature issue 801",
            state: "open",
            labels: [{ name: "internal/tracking" }, { name: "release/v15" }],
          },
        }),
        create: async () => {
          subIssueCreated = true;
          return { data: { id: 802, number: 802 } };
        },
        addLabels: async () => {},
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 801,
                title: "Feature issue 801",
                state: "open",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v15" },
                ],
              },
            ],
          },
        };
      }
      if (
        endpoint ===
        "POST /repos/rancher/terraform-provider-rancher2/issues/801/sub_issues"
      ) {
        return { data: {} };
      }
      if (endpoint.includes("/sub_issues")) {
        subIssueCallCount++;
        if (subIssueCallCount <= 2) {
          // First check in runBackportPr and second check in runBackportIssues return 0 sub-issues
          return { data: [] };
        }
        // Third check after on-demand generation returns the created sub-issue
        return {
          data: [
            {
              number: 802,
              title: "[release/v15] Feature issue 801",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "def5678" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(subIssueCreated, true);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found 1 sub-issues.")),
  );
});

test("runBackportPr - iterates multiple addressed numbers until an open non-PR issue is found", async () => {
  const mockCore = createMockCore();
  let adoptedIssueNumber = null;
  let gitCommands = [];

  const mockGitExec = async (args) => {
    gitCommands.push(args.join(" "));
    return { stdout: "", stderr: "" };
  };

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 49,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "Multi-addressed PR\n\n- Addresses: #901\n- Addresses: #902",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({
          data: { number: 491, node_id: "PR_node_491" },
        }),
      },
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 901) {
            // First reference is closed
            return {
              data: {
                number: 901,
                title: "Closed issue 901",
                state: "closed",
                labels: [{ name: "release/v15" }],
              },
            };
          }
          // Second reference is open
          return {
            data: {
              number: 902,
              title: "Open issue 902",
              state: "open",
              labels: [{ name: "release/v15" }],
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          if (labels.includes("internal/tracking")) {
            adoptedIssueNumber = issue_number;
          }
        },
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            {
              number: 903,
              title: "[release/v15] Sub issue for 902",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "ghi7890" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.strictEqual(adoptedIssueNumber, 902);
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found tracking issue: #902")),
  );
});

test("runBackportPr - creates backport children for every release label on tracking issue", async () => {
  const mockCore = createMockCore();
  const createdSubIssueBranches = [];
  let gitCommands = [];

  const mockGitExec = async (args) => {
    gitCommands.push(args.join(" "));
    return { stdout: "", stderr: "" };
  };

  let subIssueFetchCount = 0;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 50,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "Multi-release backport PR\n\n- Addresses: #1001",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({
          data: { number: 501, node_id: "PR_node_501" },
        }),
      },
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            title: "Issue 1001",
            state: "open",
            labels: [
              { name: "internal/tracking" },
              { name: "release/v14" },
              { name: "release/v15" },
            ],
          },
        }),
        create: async ({ body }) => {
          if (body?.includes("release/v14")) {
            createdSubIssueBranches.push("release/v14");
            return { data: { id: 1002, number: 1002 } };
          }
          if (body?.includes("release/v15")) {
            createdSubIssueBranches.push("release/v15");
            return { data: { id: 1003, number: 1003 } };
          }
          return { data: { id: 9999, number: 9999 } };
        },
        addLabels: async () => {},
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 1001,
                title: "Issue 1001",
                state: "open",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v14" },
                  { name: "release/v15" },
                ],
              },
            ],
          },
        };
      }
      if (endpoint.startsWith("POST ")) {
        return { data: {} };
      }
      if (endpoint.includes("/sub_issues")) {
        subIssueFetchCount++;
        // Initial check in runBackportPr and checks in runBackportIssues
        if (subIssueFetchCount <= 3) {
          return { data: [] };
        }
        // Final refetched call returns both created sub-issues
        return {
          data: [
            {
              number: 1002,
              title: "[release/v14] Issue 1001",
              labels: [{ name: "release/v14" }],
            },
            {
              number: 1003,
              title: "[release/v15] Issue 1001",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "jkl0123" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(createdSubIssueBranches.includes("release/v14"));
  assert.ok(createdSubIssueBranches.includes("release/v15"));
  assert.ok(
    mockCore.infoMessages.some((m) => m.includes("Found 2 sub-issues.")),
  );
});

test("runBackportPr - throws and propagates API errors from direct issue lookup (non-404)", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 51,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1101",
            },
          ],
        }),
      },
      issues: {
        get: async () => {
          const err = new Error("Internal Server Error");
          err.status = 500;
          throw err;
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "err123" } },
      gitExec: mockGitExec,
    });
  }, /Direct issue lookup for #1101 failed with API error: Internal Server Error/);
});

test("runBackportPr - continues search when direct issue lookup returns 404 and adopts subsequent valid issue", async () => {
  const mockCore = createMockCore();
  let adoptedNumber = null;
  const mockGitExec = async () => ({ stdout: "", stderr: "" });

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 52,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1101\n- Addresses: #1102",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({ data: { number: 521, node_id: "PR_521" } }),
      },
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 1101) {
            const err = new Error("Not Found");
            err.status = 404;
            throw err;
          }
          return {
            data: {
              number: 1102,
              title: "Valid Issue",
              state: "open",
              labels: [{ name: "release/v15" }],
            },
          };
        },
        addLabels: async ({ issue_number, labels }) => {
          if (labels.includes("internal/tracking")) {
            adoptedNumber = issue_number;
          }
        },
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            {
              number: 1103,
              title: "[release/v15] Valid Issue",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "search404" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(adoptedNumber, 1102);
  assert.ok(
    mockCore.infoMessages.some((m) =>
      m.includes("Referenced issue #1101 was not found (404)"),
    ),
  );
});

test("runBackportPr - throws and propagates error when adding tracking labels fails during direct lookup", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 53,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1103",
            },
          ],
        }),
      },
      issues: {
        get: async () => ({
          data: {
            number: 1103,
            title: "Found Issue",
            state: "open",
            labels: [],
          },
        }),
        addLabels: async () => {
          throw new Error("Permission denied on label write");
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "labelerr" } },
      gitExec: mockGitExec,
    });
  }, /Failed to add tracking labels to issue #1103: Permission denied on label write/);
});

test("runBackportPr - throws and propagates error when on-demand sub-issue generation fails", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 54,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1104",
            },
          ],
        }),
      },
      issues: {
        get: async () => ({
          data: {
            number: 1104,
            title: "Issue 1104",
            state: "open",
            labels: [{ name: "internal/tracking" }, { name: "release/v15" }],
          },
        }),
        create: async () => {
          throw new Error("API Rate limit exceeded on issue create");
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 1104,
                title: "Issue 1104",
                state: "open",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v15" },
                ],
              },
            ],
          },
        };
      }
      if (endpoint.includes("/sub_issues")) {
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "subcreateerr" } },
      gitExec: mockGitExec,
    });
  }, /Failed to create on-demand sub-issue\(s\) for #1104/);
});

test("runBackportPr - throws and propagates error when refetching sub-issues fails", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });
  let subIssueFetchCount = 0;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 55,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1105",
            },
          ],
        }),
      },
      issues: {
        get: async () => ({
          data: {
            number: 1105,
            title: "Issue 1105",
            state: "open",
            labels: [{ name: "internal/tracking" }, { name: "release/v15" }],
          },
        }),
        create: async () => ({ data: { id: 1106, number: 1106 } }),
        addLabels: async () => {},
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 1105,
                title: "Issue 1105",
                state: "open",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v15" },
                ],
              },
            ],
          },
        };
      }
      if (endpoint.startsWith("POST ")) {
        return { data: {} };
      }
      if (endpoint.includes("/sub_issues")) {
        subIssueFetchCount++;
        if (subIssueFetchCount <= 2) {
          return { data: [] };
        }
        // Refetch call fails
        throw new Error("Refetch network timeout");
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "refetcherr" } },
      gitExec: mockGitExec,
    });
  }, /Failed to refetch sub-issues for #1105: Refetch network timeout/);
});

test("runBackportPr - throws when release branch discovery fails or finds no branches", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 56,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1106",
            },
          ],
        }),
        listBranches: async () => {
          throw new Error("GitHub API 500 on listBranches");
        },
      },
      issues: {
        get: async () => ({
          data: {
            number: 1106,
            title: "Issue 1106",
            state: "open",
            labels: [{ name: "internal/tracking" }],
          },
        }),
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 1106,
                title: "Issue 1106",
                state: "open",
                labels: [{ name: "internal/tracking" }],
              },
            ],
          },
        };
      }
      if (endpoint.includes("/sub_issues")) {
        return { data: [] };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "branchdiscoveryfail" } },
      gitExec: mockGitExec,
    });
  }, /Failed to determine target release branch for #1106: GitHub API 500 on listBranches/);
});

test("runBackportPr - throws and fails when refetched sub-issues list is missing one of the requested release labels", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });
  let subIssueFetchCount = 0;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 60,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #1201",
            },
          ],
        }),
      },
      issues: {
        get: async () => ({
          data: {
            number: 1201,
            title: "Issue 1201",
            state: "open",
            labels: [
              { name: "internal/tracking" },
              { name: "release/v14" },
              { name: "release/v15" },
            ],
          },
        }),
        create: async () => ({ data: { id: 1202, number: 1202 } }),
        addLabels: async () => {},
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            total_count: 1,
            items: [
              {
                number: 1201,
                title: "Issue 1201",
                state: "open",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v14" },
                  { name: "release/v15" },
                ],
              },
            ],
          },
        };
      }
      if (endpoint.startsWith("POST ")) {
        return { data: {} };
      }
      if (endpoint.includes("/sub_issues")) {
        subIssueFetchCount++;
        // Initial check in runBackportPr and checks in runBackportIssues
        if (subIssueFetchCount <= 3) {
          return { data: [] };
        }
        // Refetched call only returns release/v14, missing release/v15
        return {
          data: [
            {
              number: 1202,
              title: "[release/v14] Issue 1201",
              labels: [{ name: "release/v14" }],
            },
          ],
        };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await assert.rejects(async () => {
    await runBackportPr({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { MERGE_COMMIT_SHA: "missingsublabel" } },
      gitExec: mockGitExec,
    });
  }, /Failed to verify created sub-issues on demand for #1201: missing sub-issue\(s\) for label\(s\) \[release\/v15\]/);
});

test("runBackportPr - adopts referenced issue on demand when PR body uses trailing period (- Addresses: #1301.)", async () => {
  const mockCore = createMockCore();
  const gitCommands = [];
  const mockGitExec = async (args) => {
    gitCommands.push(args.join(" "));
    return { stdout: "", stderr: "" };
  };

  const addedLabels = [];
  let subIssuesEndpointCalled = false;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 61,
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "Context\n- Addresses: #1301.\n* Addresses: `#1302`.\nFooter",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({ data: { number: 611, node_id: "PR_611" } }),
      },
      issues: {
        get: async ({ issue_number }) => {
          if (issue_number === 1301) {
            return {
              data: {
                number: 1301,
                title: "Trailing Period Issue",
                state: "open",
                labels: [{ name: "release/v15" }],
              },
            };
          }
          throw new Error("Not found");
        },
        create: async () => ({ data: { id: 1303, number: 1303 } }),
        addLabels: async ({ issue_number, labels }) => {
          addedLabels.push({ issue_number, labels });
        },
      },
    },
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return { data: { total_count: 0, items: [] } };
      }
      if (endpoint.startsWith("POST ")) {
        return { data: {} };
      }
      if (endpoint.includes("/sub_issues")) {
        if (!subIssuesEndpointCalled) {
          subIssuesEndpointCalled = true;
          return { data: [] };
        }
        return {
          data: [
            {
              number: 1303,
              title: "[release/v15] Trailing Period Issue",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      if (endpoint === "POST /graphql") {
        return { data: {} };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "trailperiodsha" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(
    mockCore.infoMessages.some((msg) =>
      msg.includes("Attempting direct lookup for referenced issue #1301..."),
    ),
  );
  assert.ok(
    addedLabels.some(
      (l) => l.issue_number === 1301 && l.labels.includes("internal/tracking"),
    ),
  );
});

test("runBackportPr - posts summary comment to tracking issue when backport PRs are created", async () => {
  const mockCore = createMockCore();
  const mockGitExec = async () => ({ stdout: "", stderr: "" });
  let createdComment = null;

  const mockGithub = {
    rest: {
      repos: {
        listPullRequestsAssociatedWithCommit: async () => ({
          data: [
            {
              number: 42,
              title: "fix: something",
              base: { ref: "main" },
              merged_at: "2026-10-01T00:00:00Z",
              body: "PR body\n\n- Addresses: #2000",
            },
          ],
        }),
      },
      pulls: {
        create: async () => ({
          data: {
            number: 1042,
            node_id: "PR_node_id",
            html_url:
              "https://github.com/rancher/terraform-provider-rancher2/pull/1042",
          },
        }),
      },
      issues: {
        addLabels: async () => ({}),
        addAssignees: async () => ({}),
        createComment: async ({ issue_number, body }) => {
          createdComment = { issue_number, body };
          return { data: {} };
        },
      },
    },
    graphql: async () => ({}),
    request: async (endpoint) => {
      if (endpoint === "GET /search/issues") {
        return {
          data: {
            items: [
              {
                number: 2000,
                title: "Tracking issue 2000",
                labels: [
                  { name: "internal/tracking" },
                  { name: "release/v15" },
                ],
                body: "<!-- tracking-pr: #42 -->",
                pull_request: null,
              },
            ],
          },
        };
      }
      if (
        endpoint ===
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues"
      ) {
        return {
          data: [
            {
              number: 2001,
              title: "[release/v15] Sub issue",
              labels: [{ name: "release/v15" }],
            },
          ],
        };
      }
      if (endpoint === "POST /graphql") {
        return { data: {} };
      }
      throw new Error(`Unexpected endpoint: ${endpoint}`);
    },
  };

  await runBackportPr({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { MERGE_COMMIT_SHA: "sha_summary_test" } },
    gitExec: mockGitExec,
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(createdComment !== null);
  assert.strictEqual(createdComment.issue_number, 2000);
  assert.ok(createdComment.body.includes("Automated Backport PRs Created"));
  assert.ok(createdComment.body.includes("release/v15"));
  assert.ok(createdComment.body.includes("#1042"));
  assert.ok(createdComment.body.includes("/checks"));
});

test("runMergeLabel - posts progression comment to tracking issue linking release-please PR", async () => {
  const mockCore = createMockCore();
  let createdComment = null;

  const mockGithub = {
    rest: {
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            pull_request: null,
            labels: [{ name: "internal/backport" }],
          },
        }),
        addLabels: async () => ({}),
        createComment: async ({ issue_number, body }) => {
          createdComment = { issue_number, body };
          return { data: {} };
        },
      },
      pulls: {
        list: async () => ({
          data: [
            {
              number: 3000,
              title: "chore(release/v15): release 15.1.0",
              head: { ref: "release-please--branches--release/v15" },
              labels: [{ name: "autorelease: pending" }],
            },
          ],
        }),
      },
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 2626,
        base: { ref: "release/v15" },
        body: "This PR cherry-picks changes.\n\nAddresses #2623 for #2614",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.ok(createdComment !== null);
  assert.strictEqual(createdComment.issue_number, 2614);
  assert.ok(createdComment.body.includes("Backport Merged: `release/v15`"));
  assert.ok(createdComment.body.includes("#2626"));
  assert.ok(createdComment.body.includes("#3000"));
  assert.ok(createdComment.body.includes("/pull/3000/checks"));
});

test("updateIssueProjectStatus - updates Status field to Ready for QA", async () => {
  const mockCore = createMockCore();
  let mutationVariables = null;

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getIssueProjectItems")) {
        return {
          repository: {
            issue: {
              id: "issue-node-id",
              projectItems: {
                nodes: [
                  {
                    id: "item-123",
                    project: {
                      id: "project-456",
                      title: "Rancher Project Board",
                      fields: {
                        nodes: [
                          {
                            id: "field-status",
                            name: "Status",
                            options: [
                              { id: "opt-backlog", name: "Backlog" },
                              { id: "opt-qa", name: "Ready for QA" },
                            ],
                          },
                        ],
                      },
                    },
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        mutationVariables = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      throw new Error(`Unexpected query: ${query}`);
    },
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    targetStatus: "Ready for QA",
    core: mockCore,
  });

  assert.strictEqual(updated, true);
  assert.deepStrictEqual(mutationVariables, {
    projectId: "project-456",
    itemId: "item-123",
    fieldId: "field-status",
    optionId: "opt-qa",
  });
});

test("updateIssueProjectStatus - adds issue directly to Project 73 when not already present in project items", async () => {
  const mockCore = createMockCore();
  let addVars = null;
  let updateVars = null;

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getIssueProjectItems")) {
        return {
          repository: {
            issue: {
              id: "issue-node-777",
              projectItems: { nodes: [] },
            },
          },
        };
      }
      if (query.includes("getProjectByNumber")) {
        assert.strictEqual(vars.projectNumber, 73);
        return {
          organization: {
            projectV2: {
              id: "p73-id",
              number: 73,
              title: "Rancher Board #73",
              fields: {
                nodes: [
                  {
                    id: "f-status-73",
                    name: "Status",
                    options: [{ id: "opt-qa-73", name: "Ready for QA" }],
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("addItemToProject")) {
        addVars = vars;
        return {
          addProjectV2ItemById: {
            item: { id: "item-node-888" },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        updateVars = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    targetStatus: "Ready for QA",
    projectNumber: 73,
    core: mockCore,
  });

  assert.strictEqual(updated, true);
  assert.deepStrictEqual(addVars, {
    projectId: "p73-id",
    contentId: "issue-node-777",
  });
  assert.deepStrictEqual(updateVars, {
    projectId: "p73-id",
    itemId: "item-node-888",
    fieldId: "f-status-73",
    optionId: "opt-qa-73",
  });
});

test("updateIssueProjectStatus - safely handles issue not associated with project items", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "issue-node-id",
          projectItems: { nodes: [] },
        },
      },
    }),
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    core: mockCore,
  });

  assert.strictEqual(updated, false);
});

test("updateIssueProjectStatus - paginates projectItems across multiple pages", async () => {
  const mockCore = createMockCore();
  const updatedItemIds = [];

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getIssueProjectItems")) {
        if (!vars.cursor) {
          return {
            repository: {
              issue: {
                id: "issue-node-page-1",
                projectItems: {
                  pageInfo: { hasNextPage: true, endCursor: "cursor-item-1" },
                  nodes: [
                    {
                      id: "item-p1",
                      project: {
                        id: "proj-1",
                        title: "Project 1",
                        fields: {
                          nodes: [
                            {
                              id: "f1",
                              name: "Status",
                              options: [{ id: "opt-qa", name: "Ready for QA" }],
                            },
                          ],
                        },
                      },
                    },
                  ],
                },
              },
            },
          };
        }
        return {
          repository: {
            issue: {
              id: "issue-node-page-1",
              projectItems: {
                pageInfo: { hasNextPage: false, endCursor: null },
                nodes: [
                  {
                    id: "item-p2",
                    project: {
                      id: "proj-2",
                      title: "Project 2",
                      fields: {
                        nodes: [
                          {
                            id: "f2",
                            name: "Status",
                            options: [{ id: "opt-qa", name: "Ready for QA" }],
                          },
                        ],
                      },
                    },
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        updatedItemIds.push(vars.itemId);
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    targetStatus: "Ready for QA",
    core: mockCore,
  });

  assert.strictEqual(updated, true);
  assert.deepStrictEqual(updatedItemIds, ["item-p1", "item-p2"]);
});

test("updateIssueProjectStatus - safely handles project missing Status field or matching option", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "issue-node-id",
          projectItems: {
            nodes: [
              {
                id: "item-123",
                project: {
                  id: "project-456",
                  title: "Other Board",
                  fields: {
                    nodes: [
                      {
                        id: "field-status",
                        name: "Status",
                        options: [{ id: "opt-backlog", name: "Backlog" }],
                      },
                    ],
                  },
                },
              },
            ],
          },
        },
      },
    }),
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    targetStatus: "Ready for QA",
    core: mockCore,
  });

  assert.strictEqual(updated, false);
});

test("updateIssueProjectStatus - handles GraphQL failure gracefully without throwing", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    graphql: async () => {
      throw new Error("GraphQL permission denied");
    },
  };

  const updated = await updateIssueProjectStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    issueNumber: 1234,
    core: mockCore,
  });

  assert.strictEqual(updated, false);
  assert.ok(
    mockCore.warningMessages.some((msg) =>
      msg.includes("Failed to update project status for issue #1234"),
    ),
  );
});

test("runMergeLabel - updates project status to Ready for QA when labeling issue", async () => {
  const mockCore = createMockCore();
  let mutationVars = null;

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getIssueProjectItems")) {
        return {
          repository: {
            issue: {
              id: "node-101",
              projectItems: {
                nodes: [
                  {
                    id: "item-101",
                    project: {
                      id: "proj-101",
                      title: "Rancher QA",
                      fields: {
                        nodes: [
                          {
                            id: "f-status",
                            name: "Status",
                            options: [{ id: "opt-qa", name: "Ready for QA" }],
                          },
                        ],
                      },
                    },
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        mutationVars = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
    rest: {
      issues: {
        get: async ({ issue_number }) => ({
          data: {
            number: issue_number,
            pull_request: null,
            labels: [{ name: "internal/backport" }],
          },
        }),
        addLabels: async () => ({}),
      },
    },
  };

  const mockContext = {
    repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    payload: {
      pull_request: {
        number: 2626,
        base: { ref: "release/v15" },
        body: "Addresses #2623",
      },
    },
  };

  await runMergeLabel({
    github: mockGithub,
    context: mockContext,
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(mockCore.failedMessage, null);
  assert.deepStrictEqual(mutationVars, {
    projectId: "proj-101",
    itemId: "item-101",
    fieldId: "f-status",
    optionId: "opt-qa",
  });
});

test("runCloseTrackingIssue - closes tracking issue and sets status to Ready For Release when all sub-issues are closed and milestone is set", async () => {
  const mockCore = createMockCore();
  let updatedIssueParams = null;
  let createdComment = null;
  let statusMutationVars = null;

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-100",
              number: 100,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: { number: 3, title: "v16.0.0" },
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      if (query.includes("getIssueProjectItems")) {
        return {
          repository: {
            issue: {
              id: "parent-50",
              projectItems: {
                nodes: [
                  {
                    id: "p-item-50",
                    project: {
                      id: "p-proj-1",
                      title: "Project Board",
                      fields: {
                        nodes: [
                          {
                            id: "p-field-status",
                            name: "Status",
                            options: [
                              { id: "opt-release", name: "Ready For Release" },
                            ],
                          },
                        ],
                      },
                    },
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        statusMutationVars = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            { number: 100, state: "closed" },
            { number: 101, state: "closed" },
          ],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async (params) => {
          updatedIssueParams = params;
          return { data: {} };
        },
        createComment: async (params) => {
          createdComment = params;
          return { data: {} };
        },
        listComments: async () => ({ data: [] }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, { closed: true, trackingIssueNumber: 50 });
  assert.strictEqual(updatedIssueParams.issue_number, 50);
  assert.strictEqual(updatedIssueParams.state, "closed");
  assert.strictEqual(updatedIssueParams.state_reason, "completed");
  assert.ok(
    createdComment.body.includes("Automatically closing tracking issue #50"),
  );
  assert.deepStrictEqual(statusMutationVars, {
    projectId: "p-proj-1",
    itemId: "p-item-50",
    fieldId: "p-field-status",
    optionId: "opt-release",
  });
});

test("runCloseTrackingIssue - adopts milestone from sub-issue if tracking issue lacks milestone", async () => {
  const mockCore = createMockCore();
  let assignedMilestone = null;
  let closedIssue = false;

  const mockGithub = {
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-100",
              number: 100,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: null,
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      return { repository: { issue: { projectItems: { nodes: [] } } } };
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            {
              number: 100,
              state: "closed",
              milestone: { number: 8, title: "v15.2.0" },
            },
          ],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async (params) => {
          if (params.milestone) {
            assignedMilestone = params.milestone;
          }
          if (params.state === "closed") {
            closedIssue = true;
          }
          return { data: {} };
        },
        createComment: async () => ({ data: {} }),
        listComments: async () => ({ data: [] }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, { closed: true, trackingIssueNumber: 50 });
  assert.strictEqual(assignedMilestone, 8);
  assert.strictEqual(closedIssue, true);
});

test("runCloseTrackingIssue - halts closure and posts warning comment if milestone assignment fails", async () => {
  const mockCore = createMockCore();
  let warningComment = null;
  let closeAttempted = false;

  const mockGithub = {
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-100",
              number: 100,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: null,
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      return {};
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            {
              number: 100,
              state: "closed",
              milestone: { number: 8, title: "v15.2.0" },
            },
          ],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async (params) => {
          if (params.milestone) {
            throw new Error("Milestone assignment permission denied");
          }
          if (params.state === "closed") {
            closeAttempted = true;
          }
          return { data: {} };
        },
        createComment: async (params) => {
          warningComment = params;
          return { data: {} };
        },
        listComments: async () => ({ data: [] }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    closed: false,
    trackingIssueNumber: 50,
    reason: "missing_milestone",
  });
  assert.strictEqual(closeAttempted, false);
  assert.ok(warningComment);
  assert.ok(
    warningComment.body.includes(
      "cannot be automatically closed because it has not been assigned to a milestone",
    ),
  );
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes("Failed to assign milestone #8 to tracking issue #50"),
    ),
  );
});

test("runCloseTrackingIssue - paginates sub_issues and detects open sub-issue on later page", async () => {
  const mockCore = createMockCore();
  let closeAttempted = false;

  const mockGithub = {
    paginate: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return [
          { number: 101, state: "closed" },
          { number: 102, state: "closed" },
          { number: 103, state: "open" },
        ];
      }
      return [];
    },
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-101",
              number: 101,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: { number: 1, title: "v16.0.0" },
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      return {};
    },
    rest: {
      issues: {
        update: async () => {
          closeAttempted = true;
          return { data: {} };
        },
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 101 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    closed: false,
    trackingIssueNumber: 50,
    reason: "open_sub_issues",
  });
  assert.strictEqual(closeAttempted, false);
});

test("runCloseTrackingIssue - skips closure and posts warning comment when milestone is completely missing", async () => {
  const mockCore = createMockCore();
  let warningComment = null;
  let closeAttempted = false;

  const mockGithub = {
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-100",
              number: 100,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: null,
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      return {};
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [{ number: 100, state: "closed", milestone: null }],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async (params) => {
          if (params.state === "closed") {
            closeAttempted = true;
          }
          return { data: {} };
        },
        createComment: async (params) => {
          warningComment = params;
          return { data: {} };
        },
        listComments: async () => ({ data: [] }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    closed: false,
    trackingIssueNumber: 50,
    reason: "missing_milestone",
  });
  assert.strictEqual(closeAttempted, false);
  assert.ok(warningComment);
  assert.ok(
    warningComment.body.includes(
      "cannot be automatically closed because it has not been assigned to a milestone",
    ),
  );
});

test("runCloseTrackingIssue - skips closure when any sub-issue remains open", async () => {
  const mockCore = createMockCore();
  let closeAttempted = false;

  const mockGithub = {
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "sub-100",
              number: 100,
              parent: {
                id: "parent-50",
                number: 50,
                title: "Tracking Issue #50",
                state: "OPEN",
                milestone: { number: 1, title: "v16.0.0" },
                labels: { nodes: [{ name: "internal/tracking" }] },
              },
            },
          },
        };
      }
      return {};
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [
            { number: 100, state: "closed" },
            { number: 101, state: "open" },
          ],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async () => {
          closeAttempted = true;
          return { data: {} };
        },
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    closed: false,
    trackingIssueNumber: 50,
    reason: "open_sub_issues",
  });
  assert.strictEqual(closeAttempted, false);
});

test("runCloseTrackingIssue - skips when parent is already closed", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "sub-100",
          number: 100,
          parent: {
            id: "parent-50",
            number: 50,
            state: "CLOSED",
            labels: { nodes: [{ name: "internal/tracking" }] },
          },
        },
      },
    }),
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, null);
});

test("runCloseTrackingIssue - skips when issue is not tracking and has no tracking parent", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "sub-100",
          number: 100,
          parent: null,
          labels: { nodes: [{ name: "bug" }] },
        },
      },
    }),
    rest: {
      issues: {
        get: async () => ({
          data: { number: 100, labels: [{ name: "bug" }] },
        }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.strictEqual(result, null);
});

test("default export - routes close-tracking-issue mode", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          number: 200,
          state: "CLOSED",
          labels: { nodes: [{ name: "internal/tracking" }] },
        },
      },
    }),
  };

  const result = await import("./backports.js").then((m) =>
    m.default({
      github: mockGithub,
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
        payload: { issue: { number: 200 } },
      },
      core: mockCore,
      process: { env: { SCRIPT_MODE: "close-tracking-issue" } },
    }),
  );

  assert.strictEqual(result, null);
});

test("runReleaseDone - skips transition when release lookup throws error", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => {
          throw new Error("GitHub API 500 Internal Error");
        },
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { HEAD_SHA: "lookup-fail-sha" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "release_lookup_failed",
    error: "GitHub API 500 Internal Error",
  });
  assert.ok(
    mockCore.warningMessages.some((m) =>
      m.includes("Could not fetch releases: GitHub API 500 Internal Error"),
    ),
  );
});

test("runReleaseDone - skips transition when tag cannot be verified", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => ({ data: [] }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "unverified_tag",
  });
});

test("runReleaseDone - matches release by headSha and skips when matching release is an RC", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => ({
          data: [
            {
              tag_name: "v16.0.0",
              prerelease: false,
              target_commitish: "old-sha",
            },
            {
              tag_name: "v16.0.1-rc.1",
              prerelease: true,
              target_commitish: "rc-sha-999",
            },
          ],
        }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { HEAD_SHA: "rc-sha-999" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "prerelease",
    tag: "v16.0.1-rc.1",
  });
});

test("runReleaseDone - skips transition when HEAD_SHA does not match any release", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => ({
          data: [
            {
              tag_name: "v16.0.0",
              prerelease: false,
              target_commitish: "other-sha",
            },
          ],
        }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { HEAD_SHA: "unmatched-sha-123" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "unverified_tag",
    headSha: "unmatched-sha-123",
  });
});

test("runReleaseDone - verifies explicit TAG against getReleaseByTag and skips if not found", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        getReleaseByTag: async () => {
          const err = new Error("Not Found");
          err.status = 404;
          throw err;
        },
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v16.0.0-typo" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "release_not_found",
    tag: "v16.0.0-typo",
    error: "Not Found",
  });
});

test("runReleaseDone - verifies explicit TAG against getReleaseByTag and skips if prerelease", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        getReleaseByTag: async ({ tag }) => ({
          data: { tag_name: tag, prerelease: true, draft: false },
        }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v16.0.0-alpha.1" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "prerelease",
    tag: "v16.0.0-alpha.1",
  });
});

test("runReleaseDone - verifies explicit TAG against getReleaseByTag and skips if draft", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        getReleaseByTag: async ({ tag }) => ({
          data: { tag_name: tag, prerelease: false, draft: true },
        }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v16.0.0" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "draft",
    tag: "v16.0.0",
  });
});

test("runReleaseDone - matches release by headSha and skips when matching release is draft", async () => {
  const mockCore = createMockCore();
  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => ({
          data: [
            {
              tag_name: "v16.0.0",
              prerelease: false,
              draft: true,
              target_commitish: "release-sha-draft",
            },
          ],
        }),
      },
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { HEAD_SHA: "release-sha-draft" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "draft",
    tag: "v16.0.0",
  });
});

test("runReleaseDone - skips transition when TAG env contains rc", async () => {
  const mockCore = createMockCore();

  const result = await runReleaseDone({
    github: {},
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { TAG: "v15.2.0-rc.2" } },
  });

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "prerelease",
    tag: "v15.2.0-rc.2",
  });
});

test("runReleaseDone - transitions project items from Ready for Release to Done for full release", async () => {
  const mockCore = createMockCore();
  let mutationVars = null;

  const mockGithub = {
    rest: {
      repos: {
        listReleases: async () => ({
          data: [
            {
              tag_name: "v16.0.0",
              prerelease: false,
              target_commitish: "release-sha-123",
            },
          ],
        }),
      },
    },
    graphql: async (query, vars) => {
      if (query.includes("getRepoProjects")) {
        return {
          repository: {
            projectsV2: {
              nodes: [
                {
                  id: "proj-1",
                  title: "Release Board",
                  fields: {
                    nodes: [
                      {
                        id: "field-status",
                        name: "Status",
                        options: [
                          { id: "opt-ready", name: "Ready for Release" },
                          { id: "opt-done", name: "Done" },
                        ],
                      },
                    ],
                  },
                },
              ],
            },
          },
        };
      }
      if (query.includes("getOrgProjects")) {
        return { organization: { projectsV2: { nodes: [] } } };
      }
      if (query.includes("getProjectItems")) {
        return {
          node: {
            items: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: "item-55",
                  fieldValues: {
                    nodes: [
                      {
                        field: { name: "Status" },
                        name: "Ready for Release",
                        optionId: "opt-ready",
                      },
                    ],
                  },
                  content: { number: 450, title: "Tracking Issue #450" },
                },
              ],
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        mutationVars = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const result = await runReleaseDone({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
    },
    core: mockCore,
    process: { env: { HEAD_SHA: "release-sha-123" } },
  });

  assert.deepStrictEqual(result, { transitionedCount: 1 });
  assert.deepStrictEqual(mutationVars, {
    projectId: "proj-1",
    itemId: "item-55",
    fieldId: "field-status",
    optionId: "opt-done",
  });
});

test("transitionProjectItemsStatus - handles missing GraphQL client gracefully", async () => {
  const mockCore = createMockCore();
  const result = await transitionProjectItemsStatus({
    github: {},
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    core: mockCore,
  });

  assert.deepStrictEqual(result, { transitionedCount: 0 });
});

test("default export - routes release-done mode", async () => {
  const mockCore = createMockCore();
  const result = await import("./backports.js").then((m) =>
    m.default({
      github: {},
      context: {
        repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      },
      core: mockCore,
      process: { env: { SCRIPT_MODE: "release-done", TAG: "v1.2.3-rc.1" } },
    }),
  );

  assert.deepStrictEqual(result, {
    skipped: true,
    reason: "prerelease",
    tag: "v1.2.3-rc.1",
  });
});

test("runCloseTrackingIssue - handles direct tracking issue number in ISSUE_NUMBER", async () => {
  const mockCore = createMockCore();
  let updatedIssue = null;

  const mockGithub = {
    graphql: async (query) => {
      if (query.includes("getIssueParent")) {
        return {
          repository: {
            issue: {
              id: "direct-50",
              number: 50,
              state: "OPEN",
              milestone: { number: 4, title: "v16.0.0" },
              labels: { nodes: [{ name: "internal/tracking" }] },
              parent: null,
            },
          },
        };
      }
      return { repository: { issue: { projectItems: { nodes: [] } } } };
    },
    request: async (endpoint) => {
      if (endpoint.includes("/sub_issues")) {
        return {
          data: [{ number: 101, state: "closed" }],
        };
      }
      return { data: [] };
    },
    rest: {
      issues: {
        update: async (params) => {
          updatedIssue = params;
          return { data: {} };
        },
        createComment: async () => ({ data: {} }),
        listComments: async () => ({ data: [] }),
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: {},
    },
    core: mockCore,
    process: { env: { ISSUE_NUMBER: "50" } },
  });

  assert.deepStrictEqual(result, { closed: true, trackingIssueNumber: 50 });
  assert.strictEqual(updatedIssue.issue_number, 50);
  assert.strictEqual(updatedIssue.state, "closed");
});

test("runCloseTrackingIssue - deduplicates warning comment if already posted", async () => {
  const mockCore = createMockCore();
  let commentCreated = false;

  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "sub-100",
          number: 100,
          parent: {
            id: "parent-50",
            number: 50,
            state: "OPEN",
            milestone: null,
            labels: { nodes: [{ name: "internal/tracking" }] },
          },
        },
      },
    }),
    request: async () => ({
      data: [{ number: 100, state: "closed", milestone: null }],
    }),
    rest: {
      issues: {
        listComments: async () => ({
          data: [
            {
              body: "All sub-issues for this tracking issue have been resolved and closed. However, this tracking issue cannot be automatically closed because it has not been assigned to a milestone.",
            },
          ],
        }),
        createComment: async () => {
          commentCreated = true;
          return { data: {} };
        },
      },
    },
  };

  const result = await runCloseTrackingIssue({
    github: mockGithub,
    context: {
      repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
      payload: { issue: { number: 100 } },
    },
    core: mockCore,
    process: { env: {} },
  });

  assert.deepStrictEqual(result, {
    closed: false,
    trackingIssueNumber: 50,
    reason: "missing_milestone",
  });
  assert.strictEqual(commentCreated, false);
});

test("runCloseTrackingIssue - throws when sub_issues fetch throws an error", async () => {
  const mockCore = createMockCore();

  const mockGithub = {
    graphql: async () => ({
      repository: {
        issue: {
          id: "sub-100",
          number: 100,
          parent: {
            id: "parent-50",
            number: 50,
            state: "OPEN",
            labels: { nodes: [{ name: "internal/tracking" }] },
          },
        },
      },
    }),
    request: async () => {
      throw new Error("GitHub API rate limit exceeded");
    },
  };

  await assert.rejects(
    async () => {
      await runCloseTrackingIssue({
        github: mockGithub,
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: { issue: { number: 100 } },
        },
        core: mockCore,
        process: { env: {} },
      });
    },
    { message: /Failed to fetch sub-issues for tracking issue #50/ },
  );
});

test("runCloseTrackingIssue - throws when ISSUE_NUMBER has invalid non-numeric format", async () => {
  const mockCore = createMockCore();

  await assert.rejects(
    async () => {
      await runCloseTrackingIssue({
        github: {},
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: {},
        },
        core: mockCore,
        process: { env: { ISSUE_NUMBER: "50abc" } },
      });
    },
    { message: /Invalid issue number: 50abc/ },
  );

  await assert.rejects(
    async () => {
      await runCloseTrackingIssue({
        github: {},
        context: {
          repo: { owner: "rancher", repo: "terraform-provider-rancher2" },
          payload: { issue: { number: -5 } },
        },
        core: mockCore,
        process: { env: {} },
      });
    },
    { message: /Invalid issue number: -5/ },
  );
});

test("transitionProjectItemsStatus - handles multi-page items pagination", async () => {
  const mockCore = createMockCore();
  const transitionedItemIds = [];

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getRepoProjects")) {
        return {
          repository: {
            projectsV2: {
              nodes: [
                {
                  id: "proj-paged",
                  title: "Paged Board",
                  fields: {
                    nodes: [
                      {
                        id: "f-status",
                        name: "Status",
                        options: [
                          { id: "opt-ready", name: "Ready for Release" },
                          { id: "opt-done", name: "Done" },
                        ],
                      },
                    ],
                  },
                },
              ],
            },
          },
        };
      }
      if (query.includes("getOrgProjects")) {
        return { organization: { projectsV2: { nodes: [] } } };
      }
      if (query.includes("getProjectItems")) {
        if (!vars.after) {
          // Page 1
          return {
            node: {
              items: {
                pageInfo: { hasNextPage: true, endCursor: "cursor-page-2" },
                nodes: [
                  {
                    id: "item-p1",
                    fieldValues: {
                      nodes: [
                        {
                          field: { name: "Status" },
                          optionId: "opt-ready",
                          name: "Ready for Release",
                        },
                      ],
                    },
                    content: { number: 111 },
                  },
                ],
              },
            },
          };
        }
        // Page 2
        return {
          node: {
            items: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: "item-p2",
                  fieldValues: {
                    nodes: [
                      {
                        field: { name: "Status" },
                        optionId: "opt-ready",
                        name: "Ready for Release",
                      },
                    ],
                  },
                  content: { number: 222 },
                },
              ],
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        transitionedItemIds.push(vars.itemId);
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const result = await transitionProjectItemsStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    core: mockCore,
  });

  assert.deepStrictEqual(result, { transitionedCount: 2 });
  assert.deepStrictEqual(transitionedItemIds, ["item-p1", "item-p2"]);
});

test("transitionProjectItemsStatus - handles project query error gracefully and continues with other projects", async () => {
  const mockCore = createMockCore();
  const transitionedItemIds = [];

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getRepoProjects")) {
        return {
          repository: {
            projectsV2: {
              nodes: [
                {
                  id: "proj-error",
                  title: "Broken Board",
                  fields: {
                    nodes: [
                      {
                        id: "f-status",
                        name: "Status",
                        options: [
                          { id: "opt-ready", name: "Ready for Release" },
                          { id: "opt-done", name: "Done" },
                        ],
                      },
                    ],
                  },
                },
                {
                  id: "proj-working",
                  title: "Working Board",
                  fields: {
                    nodes: [
                      {
                        id: "f-status",
                        name: "Status",
                        options: [
                          { id: "opt-ready", name: "Ready for Release" },
                          { id: "opt-done", name: "Done" },
                        ],
                      },
                    ],
                  },
                },
              ],
            },
          },
        };
      }
      if (query.includes("getOrgProjects")) {
        return { organization: { projectsV2: { nodes: [] } } };
      }
      if (query.includes("getProjectItems")) {
        if (vars.projectId === "proj-error") {
          throw new Error("Project access token expired");
        }
        return {
          node: {
            items: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: "item-working-1",
                  fieldValues: {
                    nodes: [
                      {
                        field: { name: "Status" },
                        optionId: "opt-ready",
                        name: "Ready for Release",
                      },
                    ],
                  },
                  content: { number: 333 },
                },
              ],
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        transitionedItemIds.push(vars.itemId);
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const result = await transitionProjectItemsStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    core: mockCore,
  });

  assert.deepStrictEqual(result, { transitionedCount: 1 });
  assert.deepStrictEqual(transitionedItemIds, ["item-working-1"]);
  assert.ok(
    mockCore.warningMessages.some((msg) =>
      msg.includes(
        'Failed to query/transition items for project "Broken Board"',
      ),
    ),
  );
});

test("transitionProjectItemsStatus - queries Project 73 directly by number and transitions items", async () => {
  const mockCore = createMockCore();
  let mutationVars = null;

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getTargetProject")) {
        assert.strictEqual(vars.projectNumber, 73);
        assert.strictEqual(vars.projectOwner, "rancher");
        return {
          organization: {
            projectV2: {
              id: "p73-id",
              number: 73,
              title: "Rancher Project #73",
              fields: {
                nodes: [
                  {
                    id: "f-status",
                    name: "Status",
                    options: [
                      { id: "opt-ready-73", name: "Ready for Release" },
                      { id: "opt-done-73", name: "Done" },
                    ],
                  },
                ],
              },
            },
          },
        };
      }
      if (query.includes("getRepoProjects")) {
        return { repository: { projectsV2: { nodes: [] } } };
      }
      if (query.includes("getProjectItems")) {
        return {
          node: {
            items: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: "item-73",
                  fieldValues: {
                    nodes: [
                      {
                        field: { name: "Status" },
                        optionId: "opt-ready-73",
                        name: "Ready for Release",
                      },
                    ],
                  },
                  content: { number: 888 },
                },
              ],
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        mutationVars = vars;
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const result = await transitionProjectItemsStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    projectNumber: 73,
    core: mockCore,
  });

  assert.deepStrictEqual(result, { transitionedCount: 1 });
  assert.deepStrictEqual(mutationVars, {
    projectId: "p73-id",
    itemId: "item-73",
    fieldId: "f-status",
    optionId: "opt-done-73",
  });
});

test("transitionProjectItemsStatus - paginates repository projects across multiple pages", async () => {
  const mockCore = createMockCore();
  const transitionedItemIds = [];

  const mockGithub = {
    graphql: async (query, vars) => {
      if (query.includes("getTargetProject")) {
        return { organization: { projectV2: null } };
      }
      if (query.includes("getRepoProjects")) {
        if (!vars.cursor) {
          return {
            repository: {
              projectsV2: {
                pageInfo: { hasNextPage: true, endCursor: "cursor-repo-1" },
                nodes: [
                  {
                    id: "repo-proj-1",
                    title: "Repo Board 1",
                    fields: {
                      nodes: [
                        {
                          id: "f1",
                          name: "Status",
                          options: [
                            { id: "opt-ready-1", name: "Ready for Release" },
                            { id: "opt-done-1", name: "Done" },
                          ],
                        },
                      ],
                    },
                  },
                ],
              },
            },
          };
        }
        return {
          repository: {
            projectsV2: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: "repo-proj-2",
                  title: "Repo Board 2",
                  fields: {
                    nodes: [
                      {
                        id: "f2",
                        name: "Status",
                        options: [
                          { id: "opt-ready-2", name: "Ready for Release" },
                          { id: "opt-done-2", name: "Done" },
                        ],
                      },
                    ],
                  },
                },
              ],
            },
          },
        };
      }
      if (query.includes("getOrgProjects")) {
        return { organization: { projectsV2: { nodes: [] } } };
      }
      if (query.includes("getProjectItems")) {
        const optionId =
          vars.projectId === "repo-proj-1" ? "opt-ready-1" : "opt-ready-2";
        const itemId =
          vars.projectId === "repo-proj-1" ? "item-repo-1" : "item-repo-2";
        return {
          node: {
            items: {
              pageInfo: { hasNextPage: false, endCursor: null },
              nodes: [
                {
                  id: itemId,
                  fieldValues: {
                    nodes: [
                      {
                        field: { name: "Status" },
                        optionId,
                        name: "Ready for Release",
                      },
                    ],
                  },
                  content: { number: 99 },
                },
              ],
            },
          },
        };
      }
      if (query.includes("updateItemStatus")) {
        transitionedItemIds.push(vars.itemId);
        return { projectV2Item: { id: vars.itemId } };
      }
      return {};
    },
  };

  const result = await transitionProjectItemsStatus({
    github: mockGithub,
    owner: "rancher",
    repo: "terraform-provider-rancher2",
    core: mockCore,
  });

  assert.deepStrictEqual(result, { transitionedCount: 2 });
  assert.deepStrictEqual(transitionedItemIds, ["item-repo-1", "item-repo-2"]);
});
