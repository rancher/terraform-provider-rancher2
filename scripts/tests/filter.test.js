import assert from "node:assert/strict";
import test from "node:test";

import {
  extractChangesPath,
  filterProductFiles,
  filterPullRequest,
  hasProductChanges,
  isAuthorExcluded,
  isDateInRange,
  isProductPath,
  matchesLabels,
  matchesMilestone,
} from "../lib/filter.js";

test("isProductPath and hasProductChanges", async (t) => {
  await t.test("identifies rancher2/ files as default product paths", () => {
    assert.strictEqual(
      isProductPath("rancher2/resource_rancher2_cluster_v2.go"),
      true,
    );
    assert.strictEqual(isProductPath("./rancher2/provider.go"), true);
    assert.strictEqual(isProductPath("docs/resources/cluster_v2.md"), false);
    assert.strictEqual(
      isProductPath("website/docs/r/cluster.html.markdown"),
      false,
    );
  });

  await t.test("supports custom allowed prefixes when provided", () => {
    assert.strictEqual(
      isProductPath("docs/resources/cluster_v2.md", ["rancher2/", "docs/"]),
      true,
    );
  });

  await t.test("rejects testing and non-product paths", () => {
    assert.strictEqual(isProductPath("examples/use-cases/cluster.tf"), false);
    assert.strictEqual(isProductPath("test/main.tf"), false);
    assert.strictEqual(isProductPath(".github/workflows/e2e.yml"), false);
    assert.strictEqual(isProductPath(".gemini/skills/git/SKILL.md"), false);
    assert.strictEqual(isProductPath("GNUmakefile"), false);
    assert.strictEqual(isProductPath("flake.nix"), false);
  });

  await t.test("filterProductFiles returns only qualifying paths", () => {
    const files = [
      "rancher2/resource_cluster.go",
      "examples/quickstart.tf",
      "docs/index.md",
      "test/setup.go",
      "flake.nix",
    ];
    const filtered = filterProductFiles(files);
    assert.deepStrictEqual(filtered, ["rancher2/resource_cluster.go"]);
    assert.strictEqual(hasProductChanges(files), true);
  });

  await t.test(
    "hasProductChanges returns false when no product files present",
    () => {
      const files = [
        "examples/quickstart.tf",
        ".github/workflows/lint.yml",
        "flake.nix",
      ];
      assert.strictEqual(hasProductChanges(files), false);
      assert.deepStrictEqual(filterProductFiles(files), []);
    },
  );
});

test("extractChangesPath", async (t) => {
  await t.test("extracts and orders product directory prefixes", () => {
    const files = [
      "docs/resources/cluster.md",
      "rancher2/resource_cluster.go",
      "rancher2/provider.go",
    ];
    assert.strictEqual(extractChangesPath(files), "rancher2/, docs/");
  });

  await t.test("extracts non-product directories and root files", () => {
    const files = [
      ".github/workflows/e2e.yml",
      "CONTRIBUTING.md",
      "examples/test.tf",
    ];
    assert.strictEqual(
      extractChangesPath(files),
      ".github/, CONTRIBUTING.md, examples/",
    );
  });

  await t.test("handles empty or invalid inputs gracefully", () => {
    assert.strictEqual(extractChangesPath([]), "-");
    assert.strictEqual(extractChangesPath(null), "-");
  });
});

test("isAuthorExcluded", async (t) => {
  await t.test("matches author case-insensitively", () => {
    assert.strictEqual(
      isAuthorExcluded("matt-trachier", ["Matt-Trachier"]),
      true,
    );
    assert.strictEqual(
      isAuthorExcluded("Matt-Trachier", ["matt-trachier", "other"]),
      true,
    );
    assert.strictEqual(
      isAuthorExcluded("contributor", ["matt-trachier"]),
      false,
    );
  });

  await t.test("handles empty lists or missing authors gracefully", () => {
    assert.strictEqual(isAuthorExcluded("", ["someone"]), false);
    assert.strictEqual(isAuthorExcluded("alice", []), false);
    assert.strictEqual(isAuthorExcluded(null, ["someone"]), false);
  });
});

test("isDateInRange", async (t) => {
  const date = "2026-06-15T12:00:00Z";

  await t.test("validates within range", () => {
    assert.strictEqual(
      isDateInRange(date, {
        since: "2026-06-01T00:00:00Z",
        until: "2026-06-30T00:00:00Z",
      }),
      true,
    );
  });

  await t.test("rejects before since", () => {
    assert.strictEqual(
      isDateInRange(date, { since: "2026-07-01T00:00:00Z" }),
      false,
    );
  });

  await t.test("rejects after until", () => {
    assert.strictEqual(
      isDateInRange(date, { until: "2026-05-01T00:00:00Z" }),
      false,
    );
  });

  await t.test("passes with no boundaries", () => {
    assert.strictEqual(isDateInRange(date), true);
  });
});

test("matchesLabels and matchesMilestone", async (t) => {
  await t.test("matches labels case-insensitively", () => {
    const labels = [{ name: "Area/Cluster" }, { name: "release/v15" }];
    assert.strictEqual(matchesLabels(labels, ["area/cluster"]), true);
    assert.strictEqual(matchesLabels(labels, ["missing-label"]), false);
    assert.strictEqual(matchesLabels(labels, []), true);
  });

  await t.test("matches milestone case-insensitively", () => {
    const milestone = { title: "v15.3.0" };
    assert.strictEqual(matchesMilestone(milestone, "V15.3.0"), true);
    assert.strictEqual(matchesMilestone(milestone, "v15.4.0"), false);
    assert.strictEqual(matchesMilestone(null, "v15.3.0"), false);
    assert.strictEqual(matchesMilestone(milestone, null), true);
  });
});

test("filterPullRequest", async (t) => {
  const pr = {
    number: 100,
    merged_at: "2026-06-15T12:00:00Z",
    user: { login: "contributor" },
    labels: [{ name: "enhancement" }],
    milestone: { title: "v15.3.0" },
  };

  await t.test(
    "qualifies when PR touches product files and matches criteria",
    () => {
      const res = filterPullRequest(pr, {
        filePaths: ["rancher2/resource_cluster.go", "examples/cluster.tf"],
      });
      assert.strictEqual(res.qualifies, true);
      assert.deepStrictEqual(res.productFiles, [
        "rancher2/resource_cluster.go",
      ]);
    },
  );

  await t.test("rejects when PR is not merged", () => {
    const res = filterPullRequest(
      { ...pr, merged_at: null },
      {
        filePaths: ["rancher2/resource_cluster.go"],
      },
    );
    assert.strictEqual(res.qualifies, false);
    assert.strictEqual(res.reason, "PR is not merged");
  });

  await t.test("rejects when author is excluded", () => {
    const res = filterPullRequest(pr, {
      filePaths: ["rancher2/resource_cluster.go"],
      excludedAuthors: ["Contributor"],
    });
    assert.strictEqual(res.qualifies, false);
    assert.strictEqual(res.reason, "Author 'contributor' is excluded");
  });

  await t.test(
    "rejects when only non-product files modified and allPaths is false",
    () => {
      const res = filterPullRequest(pr, {
        filePaths: ["examples/cluster.tf", ".github/workflows/ci.yml"],
      });
      assert.strictEqual(res.qualifies, false);
      assert.strictEqual(res.reason, "No product files modified");
    },
  );

  await t.test("qualifies non-product changes when allPaths is true", () => {
    const res = filterPullRequest(pr, {
      filePaths: ["examples/cluster.tf", ".github/workflows/ci.yml"],
      allPaths: true,
    });
    assert.strictEqual(res.qualifies, true);
    assert.deepStrictEqual(res.productFiles, [
      "examples/cluster.tf",
      ".github/workflows/ci.yml",
    ]);
  });
});
