import assert from "node:assert/strict";
import test from "node:test";

import {
  buildHeaders,
  fetchBranchPRs,
  fetchClosedPRs,
  fetchIssueDetails,
  fetchPRFiles,
  fetchPullRequest,
  fetchReleaseBranches,
  fetchTagDate,
  fetchTrackingIssuesMap,
  searchPullRequests,
} from "../lib/github.js";

test("buildHeaders", () => {
  const withoutToken = buildHeaders();
  assert.strictEqual(withoutToken.Accept, "application/vnd.github+json");
  assert.strictEqual(withoutToken.Authorization, undefined);

  const withToken = buildHeaders("my-token");
  assert.strictEqual(withToken.Authorization, "Bearer my-token");
});

test("fetchClosedPRs pagination and boundaries", async (t) => {
  await t.test("fetches single page when items less than limit", async () => {
    const mockFetch = async (url) => {
      assert.ok(url.includes("page=1"));
      return {
        ok: true,
        json: async () => [
          { number: 1, title: "PR 1", updated_at: "2026-06-10T00:00:00Z" },
          { number: 2, title: "PR 2", updated_at: "2026-06-09T00:00:00Z" },
        ],
      };
    };

    const res = await fetchClosedPRs({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      limit: 10,
      fetchFn: mockFetch,
    });
    assert.strictEqual(res.length, 2);
    assert.strictEqual(res[0].number, 1);
  });

  await t.test("paginates until limit is reached", async () => {
    let callCount = 0;
    const mockFetch = async () => {
      callCount++;
      // Return 100 items on page 1, 100 items on page 2
      return {
        ok: true,
        json: async () =>
          Array.from({ length: 100 }, (_, i) => ({
            number: (callCount - 1) * 100 + i + 1,
            title: `PR ${(callCount - 1) * 100 + i + 1}`,
          })),
      };
    };

    const res = await fetchClosedPRs({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      limit: 150,
      fetchFn: mockFetch,
    });
    assert.strictEqual(res.length, 150);
    assert.strictEqual(callCount, 2);
  });

  await t.test("stops early when since boundary is exceeded", async () => {
    let callCount = 0;
    const mockFetch = async () => {
      callCount++;
      return {
        ok: true,
        json: async () => [
          { number: 1, updated_at: "2026-06-15T00:00:00Z" },
          { number: 2, updated_at: "2026-05-01T00:00:00Z" }, // older than since
        ],
      };
    };

    const res = await fetchClosedPRs({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      since: "2026-06-01T00:00:00Z",
      limit: 50,
      fetchFn: mockFetch,
    });
    assert.strictEqual(res.length, 2);
    assert.strictEqual(callCount, 1);
  });

  await t.test("throws error on API failure", async () => {
    const mockFetch = async () => ({
      ok: false,
      status: 403,
      statusText: "Rate limit exceeded",
      text: async () => "Rate limit exceeded",
    });

    await assert.rejects(
      async () =>
        fetchClosedPRs({
          owner: "rancher",
          repo: "terraform-provider-rancher2",
          fetchFn: mockFetch,
        }),
      { message: /GitHub API error \(403\): Rate limit exceeded/ },
    );
  });
});

test("fetchPRFiles", async (t) => {
  await t.test("extracts filenames from response", async () => {
    const mockFetch = async () => ({
      ok: true,
      json: async () => [
        { filename: "rancher2/resource_cluster.go" },
        { filename: "docs/resources/cluster.md" },
      ],
    });

    const files = await fetchPRFiles({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      prNumber: 123,
      fetchFn: mockFetch,
    });
    assert.deepStrictEqual(files, [
      "rancher2/resource_cluster.go",
      "docs/resources/cluster.md",
    ]);
  });
});

test("fetchTagDate", async (t) => {
  await t.test("returns published_at from releases API", async () => {
    const mockFetch = async (url) => {
      if (url.includes("/releases/tags/v5.2.0")) {
        return {
          ok: true,
          json: async () => ({
            published_at: "2026-05-10T14:00:00Z",
          }),
        };
      }
      return { ok: false };
    };

    const date = await fetchTagDate({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      tag: "v5.2.0",
      fetchFn: mockFetch,
    });
    assert.strictEqual(date, "2026-05-10T14:00:00Z");
  });

  await t.test("returns null if tag not found", async () => {
    const mockFetch = async () => ({ ok: false });
    const date = await fetchTagDate({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      tag: "v99.99.99",
      fetchFn: mockFetch,
    });
    assert.strictEqual(date, null);
  });
});

test("searchPullRequests", async (t) => {
  await t.test("executes search query and extracts items", async () => {
    const mockFetch = async (url) => {
      assert.ok(url.includes("search/issues"));
      return {
        ok: true,
        json: async () => ({
          total_count: 1,
          items: [{ number: 999, title: "backport PR" }],
        }),
      };
    };

    const items = await searchPullRequests({
      query:
        "repo:rancher/terraform-provider-rancher2 is:pr head:backport-123-",
      fetchFn: mockFetch,
    });
    assert.strictEqual(items.length, 1);
    assert.strictEqual(items[0].number, 999);
  });
});

test("fetchPullRequest", async (t) => {
  await t.test("fetches single PR object by number", async () => {
    const mockFetch = async (url) => {
      assert.ok(url.endsWith("/pulls/123"));
      return {
        ok: true,
        json: async () => ({
          number: 123,
          title: "Test PR",
          user: { login: "test" },
        }),
      };
    };

    const pr = await fetchPullRequest({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      prNumber: 123,
      fetchFn: mockFetch,
    });
    assert.strictEqual(pr.number, 123);
    assert.strictEqual(pr.user.login, "test");
  });
});

test("fetchBranchPRs", async (t) => {
  await t.test(
    "fetches PRs targeting a base branch via Core REST API",
    async () => {
      const mockFetch = async (url) => {
        assert.ok(url.includes("base=release%2Fv15"));
        return {
          ok: true,
          json: async () => [
            { number: 2529, head: { ref: "backport-2514-release-v15" } },
          ],
        };
      };

      const prs = await fetchBranchPRs({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        baseBranch: "release/v15",
        fetchFn: mockFetch,
      });
      assert.strictEqual(prs.length, 1);
      assert.strictEqual(prs[0].number, 2529);
    },
  );
});

test("fetchReleaseBranches", async (t) => {
  await t.test("filters and sorts release branches descending", async () => {
    const mockFetch = async () => ({
      ok: true,
      json: async () => [
        { name: "main" },
        { name: "release/v9" },
        { name: "release/v13" },
        { name: "release/v15" },
        { name: "release/v14" },
        { name: "feature/test" },
      ],
    });

    const branches = await fetchReleaseBranches({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      fetchFn: mockFetch,
    });
    assert.deepStrictEqual(branches, [
      "release/v15",
      "release/v14",
      "release/v13",
      "release/v9",
    ]);
  });
});

test("fetchTrackingIssuesMap", async (t) => {
  await t.test("maps PR numbers to tracking issue numbers", async () => {
    const mockFetch = async () => ({
      ok: true,
      json: async () => [
        {
          number: 2466,
          body: "This is the tracking issue for #2465",
        },
        {
          number: 2312,
          body: "<!-- tracking-pr: #2311 -->\nThis is the tracking issue for PR #2311",
        },
      ],
    });

    const map = await fetchTrackingIssuesMap({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      fetchFn: mockFetch,
    });
    assert.deepStrictEqual(map.get(2465), { number: 2466, state: "open" });
    assert.deepStrictEqual(map.get(2311), { number: 2312, state: "open" });
    assert.strictEqual(map.issueStates.get(2466), "open");
  });
});

test("fetchIssueDetails", async (t) => {
  await t.test("fetches closed issue with closed_by actor", async () => {
    const mockFetch = async () => ({
      ok: true,
      json: async () => ({
        number: 2313,
        state: "closed",
        closed_by: { login: "markusewalker" },
        assignee: { login: "HarrisonWAffel" },
        html_url:
          "https://github.com/rancher/terraform-provider-rancher2/issues/2313",
      }),
    });

    const res = await fetchIssueDetails({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      issueNumber: 2313,
      fetchFn: mockFetch,
    });
    assert.strictEqual(res.number, 2313);
    assert.strictEqual(res.state, "closed");
    assert.strictEqual(res.person, "markusewalker");
  });

  await t.test(
    "fetches closed issue with null closed_by falling back to assignee or unknown",
    async () => {
      const mockFetchWithAssignee = async () => ({
        ok: true,
        json: async () => ({
          number: 2314,
          state: "closed",
          closed_by: null,
          assignee: { login: "HarrisonWAffel" },
          html_url:
            "https://github.com/rancher/terraform-provider-rancher2/issues/2314",
        }),
      });

      const res1 = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 2314,
        fetchFn: mockFetchWithAssignee,
      });
      assert.strictEqual(res1.number, 2314);
      assert.strictEqual(res1.state, "closed");
      assert.strictEqual(res1.person, "HarrisonWAffel");

      const mockFetchWithoutAssignee = async () => ({
        ok: true,
        json: async () => ({
          number: 2315,
          state: "closed",
          closed_by: null,
          assignee: null,
          html_url:
            "https://github.com/rancher/terraform-provider-rancher2/issues/2315",
        }),
      });

      const res2 = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 2315,
        fetchFn: mockFetchWithoutAssignee,
      });
      assert.strictEqual(res2.number, 2315);
      assert.strictEqual(res2.state, "closed");
      assert.strictEqual(res2.person, "unknown");
    },
  );

  await t.test("fetches open issue with assignee or unassigned", async () => {
    const mockFetch = async () => ({
      ok: true,
      json: async () => ({
        number: 2518,
        state: "open",
        closed_by: null,
        assignee: null,
        html_url:
          "https://github.com/rancher/terraform-provider-rancher2/issues/2518",
      }),
    });

    const res = await fetchIssueDetails({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      issueNumber: 2518,
      fetchFn: mockFetch,
    });
    assert.strictEqual(res.number, 2518);
    assert.strictEqual(res.state, "open");
    assert.strictEqual(res.person, "unassigned");
  });

  await t.test(
    "fetches open issue falling back to assignees array when assignee is null",
    async () => {
      const mockFetch = async () => ({
        ok: true,
        json: async () => ({
          number: 2519,
          state: "open",
          closed_by: null,
          assignee: null,
          assignees: [{ login: "qa-lead" }, { login: "dev-lead" }],
          html_url:
            "https://github.com/rancher/terraform-provider-rancher2/issues/2519",
        }),
      });

      const res = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 2519,
        fetchFn: mockFetch,
      });
      assert.strictEqual(res.number, 2519);
      assert.strictEqual(res.state, "open");
      assert.strictEqual(res.person, "qa-lead");
    },
  );

  await t.test(
    "falls back to unknown state and person on non-ok or error",
    async () => {
      const mockFetch404 = async () => ({
        ok: false,
        status: 404,
        json: async () => ({}),
      });

      const res1 = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 9999,
        fetchFn: mockFetch404,
      });
      assert.strictEqual(res1.number, 9999);
      assert.strictEqual(res1.state, "unknown");
      assert.strictEqual(res1.person, "unknown");

      const mockFetchThrow = async () => {
        throw new Error("Network timeout");
      };

      const res2 = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 8888,
        fetchFn: mockFetchThrow,
      });
      assert.strictEqual(res2.number, 8888);
      assert.strictEqual(res2.state, "unknown");
      assert.strictEqual(res2.person, "unknown");

      const mockFetchSyntaxError = async () => ({
        ok: true,
        json: async () => {
          throw new SyntaxError("Unexpected token");
        },
      });

      const res3 = await fetchIssueDetails({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
        issueNumber: 7777,
        fetchFn: mockFetchSyntaxError,
      });
      assert.strictEqual(res3.number, 7777);
      assert.strictEqual(res3.state, "unknown");
      assert.strictEqual(res3.person, "unknown");
    },
  );
});
