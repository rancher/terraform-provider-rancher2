import assert from "node:assert/strict";
import test from "node:test";

import {
  buildBackportSearchQuery,
  correlateBackports,
  extractAttachedIssueNumber,
  extractBackportInfo,
  filterOpenBackportIssues,
  filterOpenQAIssues,
  formatBackportGaps,
  formatBackportMatrix,
  formatBackportTable,
  parseBackportBranch,
} from "../lib/backports.js";

test("parseBackportBranch", async (t) => {
  await t.test("parses standard backport branch name", () => {
    const res = parseBackportBranch("backport-2633-release-v15");
    assert.deepStrictEqual(res, {
      mainPrNumber: 2633,
      branchSuffix: "release-v15",
      targetBranch: "release/v15",
    });
  });

  await t.test("parses custom branch target", () => {
    const res = parseBackportBranch("backport-1234-maintenance");
    assert.deepStrictEqual(res, {
      mainPrNumber: 1234,
      branchSuffix: "maintenance",
      targetBranch: "maintenance",
    });
  });

  await t.test("returns null for non-backport branch", () => {
    assert.strictEqual(parseBackportBranch("feature/cluster-proxy"), null);
    assert.strictEqual(parseBackportBranch(""), null);
    assert.strictEqual(parseBackportBranch(null), null);
  });
});

test("buildBackportSearchQuery", async (t) => {
  await t.test("builds single PR search query with prNumber", () => {
    const query = buildBackportSearchQuery({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      prNumber: 2633,
    });
    assert.strictEqual(
      query,
      "repo:rancher/terraform-provider-rancher2 is:pr head:backport-2633-",
    );
  });

  await t.test("builds single PR search query with prNumbers array", () => {
    const query = buildBackportSearchQuery({
      owner: "rancher",
      repo: "terraform-provider-rancher2",
      prNumbers: [2514],
    });
    assert.strictEqual(
      query,
      "repo:rancher/terraform-provider-rancher2 is:pr head:backport-2514-",
    );
  });

  await t.test(
    "builds general backport query when no PR number provided",
    () => {
      const query = buildBackportSearchQuery({
        owner: "rancher",
        repo: "terraform-provider-rancher2",
      });
      assert.strictEqual(
        query,
        "repo:rancher/terraform-provider-rancher2 is:pr head:backport-",
      );
    },
  );
});

test("extractBackportInfo", async (t) => {
  await t.test("extracts info from full PR object with head.ref", () => {
    const fullPR = {
      number: 2529,
      head: { ref: "backport-2514-release-v15" },
      base: { ref: "release/v15" },
      merged_at: "2026-09-29T16:50:43Z",
      html_url:
        "https://github.com/rancher/terraform-provider-rancher2/pull/2529",
    };
    const info = extractBackportInfo(fullPR);
    assert.strictEqual(info.mainPrNumber, 2514);
    assert.strictEqual(info.targetBranch, "release/v15");
    assert.strictEqual(info.prNumber, 2529);
    assert.strictEqual(info.state, "merged");
  });

  await t.test(
    "extracts info from /search/issues item without head.ref",
    () => {
      const searchItem = {
        number: 2529,
        body: "This pull request cherry-picks the changes from #2514 into release/v15\n\nAddresses #2518",
        labels: [{ name: "release/v15" }, { name: "internal/pr-backport" }],
        pull_request: {
          merged_at: "2026-09-29T16:50:43Z",
        },
        html_url:
          "https://github.com/rancher/terraform-provider-rancher2/pull/2529",
      };
      const info = extractBackportInfo(searchItem);
      assert.strictEqual(info.mainPrNumber, 2514);
      assert.strictEqual(info.targetBranch, "release/v15");
      assert.strictEqual(info.prNumber, 2529);
      assert.strictEqual(info.state, "merged");
      assert.strictEqual(info.mergedAt, "2026-09-29T16:50:43Z");
    },
  );

  await t.test("extracts info from manual backport with PR URL in body", () => {
    const manualPR = {
      number: 2497,
      head: { ref: "backport-priv-registry-improvements" },
      base: { ref: "release/v15" },
      body: "This is a backport of https://github.com/rancher/terraform-provider-rancher2/pull/2311 for v15.",
      merged_at: "2026-09-21T16:46:42Z",
      html_url:
        "https://github.com/rancher/terraform-provider-rancher2/pull/2497",
    };
    const info = extractBackportInfo(manualPR);
    assert.strictEqual(info.mainPrNumber, 2311);
    assert.strictEqual(info.targetBranch, "release/v15");
    assert.strictEqual(info.prNumber, 2497);
    assert.strictEqual(info.state, "merged");
  });

  await t.test(
    "extracts info from manual backport with #PR shorthand in body",
    () => {
      const manualPR = {
        number: 3001,
        base: { ref: "release/v14" },
        body: "Backport of #2200 to maintenance branch.",
        merged_at: "2026-08-10T12:00:00Z",
        html_url:
          "https://github.com/rancher/terraform-provider-rancher2/pull/3001",
      };
      const info = extractBackportInfo(manualPR);
      assert.strictEqual(info.mainPrNumber, 2200);
      assert.strictEqual(info.targetBranch, "release/v14");
      assert.strictEqual(info.prNumber, 3001);
      assert.strictEqual(info.state, "merged");
    },
  );

  await t.test("extracts info from backport reference in title", () => {
    const titlePR = {
      number: 3002,
      base: { ref: "release/v15" },
      title: "[Backport #2100] feat: cluster feature",
      body: "Manual cherry-pick without body reference.",
      merged_at: "2026-08-11T12:00:00Z",
      html_url:
        "https://github.com/rancher/terraform-provider-rancher2/pull/3002",
    };
    const info = extractBackportInfo(titlePR);
    assert.strictEqual(info.mainPrNumber, 2100);
    assert.strictEqual(info.targetBranch, "release/v15");
    assert.strictEqual(info.prNumber, 3002);
  });

  await t.test("returns null for non-backport item", () => {
    assert.strictEqual(extractBackportInfo(null), null);
    assert.strictEqual(
      extractBackportInfo({ number: 999, body: "Normal PR" }),
      null,
    );
  });
});

test("correlateBackports", async (t) => {
  const mainPRs = [
    { number: 100, title: "feat: main pr 100" },
    { number: 200, title: "fix: main pr 200" },
  ];

  const rawBackports = [
    {
      number: 105,
      head: { ref: "backport-100-release-v15" },
      base: { ref: "release/v15" },
      state: "closed",
      merged_at: "2026-06-15T15:00:00Z",
      url: "https://github.com/rancher/terraform-provider-rancher2/pull/105",
    },
    {
      number: 205,
      head: { ref: "backport-200-release-v15" },
      base: { ref: "release/v15" },
      state: "open",
      merged_at: null,
      url: "https://github.com/rancher/terraform-provider-rancher2/pull/205",
    },
  ];

  await t.test(
    "correlates backports and computes completed and missing branches",
    () => {
      const enriched = correlateBackports(mainPRs, rawBackports, {
        expectedBranches: ["release/v15", "release/v14"],
      });

      assert.strictEqual(enriched.length, 2);

      // PR 100: release/v15 is merged, release/v14 is missing
      const pr100 = enriched.find((p) => p.number === 100);
      assert.strictEqual(pr100.backports.length, 1);
      assert.strictEqual(pr100.backports[0].state, "merged");
      assert.deepStrictEqual(pr100.backport_status.completed_branches, [
        "release/v15",
      ]);
      assert.deepStrictEqual(pr100.backport_status.missing_branches, [
        "release/v14",
      ]);

      // PR 200: release/v15 is open (not merged), both are missing
      const pr200 = enriched.find((p) => p.number === 200);
      assert.strictEqual(pr200.backports[0].state, "open");
      assert.deepStrictEqual(pr200.backport_status.completed_branches, []);
      assert.deepStrictEqual(pr200.backport_status.missing_branches, [
        "release/v15",
        "release/v14",
      ]);
    },
  );
});

test("formatBackportMatrix and formatBackportGaps", async (t) => {
  const enrichedPRs = [
    {
      number: 100,
      title: "feat: proxy support",
      url: "https://github.com/rancher/terraform-provider-rancher2/pull/100",
      backports: [
        {
          target_branch: "release/v15",
          pr_number: 105,
          state: "merged",
          url: "https://github.com/rancher/terraform-provider-rancher2/pull/105",
        },
      ],
      backport_status: {
        expected_branches: ["release/v15", "release/v14"],
        completed_branches: ["release/v15"],
        missing_branches: ["release/v14"],
      },
    },
  ];

  await t.test("formatBackportMatrix renders markdown table", () => {
    const matrix = formatBackportMatrix(enrichedPRs, [
      "release/v15",
      "release/v14",
    ]);
    assert.ok(matrix.includes("Main PR"));
    assert.ok(matrix.includes("release/v15"));
    assert.ok(matrix.includes("release/v14"));
    assert.ok(
      matrix.includes(
        "✅ [#105](https://github.com/rancher/terraform-provider-rancher2/pull/105) (Merged)",
      ),
    );
    assert.ok(matrix.includes("❌ Missing"));
  });

  await t.test("formatBackportGaps renders report for missing branches", () => {
    const gaps = formatBackportGaps(enrichedPRs);
    assert.ok(gaps.includes("# Backport Gap Report"));
    assert.ok(gaps.includes("**Missing Target Branches**: release/v14"));
  });
});

test("extractAttachedIssueNumber", async (t) => {
  await t.test(
    "extracts issue number from Addresses #issue for #tracking",
    () => {
      const body =
        "This pull request cherry-picks the changes from #2514 into release/v15\n\nAddresses #2518 for #2517";
      assert.strictEqual(extractAttachedIssueNumber(body), 2518);
    },
  );

  await t.test("extracts issue number from Addresses with issue URL", () => {
    const body =
      "- Addresses # https://github.com/rancher/terraform-provider-rancher2/issues/2313";
    assert.strictEqual(extractAttachedIssueNumber(body), 2313);
  });

  await t.test("returns null when no attached issue mentioned", () => {
    assert.strictEqual(
      extractAttachedIssueNumber("Just a regular PR description"),
      null,
    );
    assert.strictEqual(extractAttachedIssueNumber(""), null);
    assert.strictEqual(extractAttachedIssueNumber(null), null);
  });
});

test("filterOpenQAIssues", async (t) => {
  const prWithOpenIssue = {
    number: 1,
    backports: [
      {
        pr_number: 10,
        qa_issue: { number: 100, state: "open" },
      },
    ],
  };

  const prWithClosedIssue = {
    number: 2,
    backports: [
      {
        pr_number: 20,
        qa_issue: { number: 200, state: "closed" },
      },
    ],
  };

  const prWithoutIssue = {
    number: 3,
    backports: [
      {
        pr_number: 30,
        qa_issue: null,
      },
    ],
  };

  const prNoBackports = {
    number: 4,
    backports: [],
  };

  await t.test("keeps only PRs having at least one open QA issue", () => {
    const filtered = filterOpenQAIssues([
      prWithOpenIssue,
      prWithClosedIssue,
      prWithoutIssue,
      prNoBackports,
    ]);
    assert.strictEqual(filtered.length, 1);
    assert.strictEqual(filtered[0].number, 1);
  });

  await t.test("filterOpenBackportIssues alias functions identically", () => {
    const filtered = filterOpenBackportIssues([
      prWithOpenIssue,
      prWithClosedIssue,
    ]);
    assert.strictEqual(filtered.length, 1);
    assert.strictEqual(filtered[0].number, 1);
  });

  await t.test("returns empty array for invalid or empty inputs", () => {
    assert.deepStrictEqual(filterOpenQAIssues([]), []);
    assert.deepStrictEqual(filterOpenQAIssues(null), []);
  });
});

test("formatBackportTable", async (t) => {
  const sample = [
    {
      number: 123,
      title: "feat: cluster proxy support",
      author: "Jim",
      state: "merged",
      target_branch: "main",
      url: "https://github.com/rancher/terraform-provider-rancher2/pull/123",
      addressed_issues: [{ number: 456, state: "open" }],
      changed_product_files: ["rancher2/cluster.go", "docs/cluster.md"],
      backports: [
        {
          pr_number: 789,
          state: "merged",
          target_branch: "release/v15",
          url: "https://github.com/rancher/terraform-provider-rancher2/pull/789",
          qa_issue: {
            number: 7890,
            state: "closed",
            person: "Karol",
            url: "https://github.com/rancher/terraform-provider-rancher2/issues/7890",
          },
        },
        {
          pr_number: 109,
          state: "merged",
          target_branch: "release/v14",
          url: "https://github.com/rancher/terraform-provider-rancher2/pull/109",
          qa_issue: {
            number: 1090,
            state: "open",
            person: "unassigned",
            url: "https://github.com/rancher/terraform-provider-rancher2/issues/1090",
          },
        },
      ],
    },
  ];

  await t.test(
    "renders table with exact columns: Title, PR#, Author, Tracking Issue #, Backport PRs, QA Issues, ChangesPath",
    () => {
      const table = formatBackportTable(sample);
      assert.ok(table.includes("Title"));
      assert.ok(table.includes("PR#"));
      assert.ok(table.includes("Author"));
      assert.ok(table.includes("Tracking Issue #"));
      assert.ok(table.includes("Backport PRs"));
      assert.ok(table.includes("QA Issues"));
      assert.ok(table.includes("ChangesPath"));
      assert.ok(table.includes("feat: cluster proxy support"));
      assert.ok(
        table.includes(
          "[123 - merged main](https://github.com/rancher/terraform-provider-rancher2/pull/123)",
        ),
      );
      assert.ok(table.includes("Jim"));
      assert.ok(
        table.includes(
          "[456 - open](https://github.com/rancher/terraform-provider-rancher2/issues/456)",
        ),
      );
      assert.ok(
        table.includes(
          "[789 - merged release/v15](https://github.com/rancher/terraform-provider-rancher2/pull/789), [109 - merged release/v14](https://github.com/rancher/terraform-provider-rancher2/pull/109)",
        ),
      );
      assert.ok(
        table.includes(
          "[7890 - closed - Karol](https://github.com/rancher/terraform-provider-rancher2/issues/7890), [1090 - open - unassigned](https://github.com/rancher/terraform-provider-rancher2/issues/1090)",
        ),
      );
      assert.ok(table.includes("rancher2/, docs/"));
    },
  );

  await t.test("handles PR without issues or backports gracefully", () => {
    const sampleNoIssues = [
      {
        number: 300,
        author: "Alice",
        state: "open",
        target_branch: "main",
        url: "https://github.com/rancher/terraform-provider-rancher2/pull/300",
        addressed_issues: [],
        backports: [],
      },
    ];
    const table = formatBackportTable(sampleNoIssues);
    assert.ok(
      table.includes(
        "[300 - open main](https://github.com/rancher/terraform-provider-rancher2/pull/300)",
      ),
    );
    assert.ok(table.includes("Alice"));
    assert.ok(table.includes("-"));
  });

  await t.test("escapes pipe characters in PR titles", () => {
    const sampleWithPipe = [
      {
        number: 400,
        title: "feat(cluster): proxy | v2 configuration",
        author: "Alice",
        state: "merged",
        target_branch: "main",
        url: "https://github.com/rancher/terraform-provider-rancher2/pull/400",
        addressed_issues: [],
        backports: [],
      },
    ];
    const table = formatBackportTable(sampleWithPipe);
    assert.ok(table.includes("feat(cluster): proxy \\| v2 configuration"));
  });

  await t.test("handles empty array gracefully", () => {
    const table = formatBackportTable([]);
    assert.ok(table.includes("No pull requests to display in table."));
  });
});
