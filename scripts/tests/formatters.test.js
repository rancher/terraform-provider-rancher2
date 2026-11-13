import assert from "node:assert/strict";
import test from "node:test";

import { formatJSON, formatMarkdown, formatTable } from "../lib/formatters.js";

const samplePRs = [
  {
    number: 101,
    title: "feat(cluster): add proxy configuration",
    author: "alice",
    url: "https://github.com/rancher/terraform-provider-rancher2/pull/101",
    conventional_type: "feat",
    conventional_scope: "cluster",
    addressed_issues: [1001],
    changed_product_files: ["rancher2/resource_cluster.go"],
  },
  {
    number: 102,
    title: "fix: nil pointer check on delete",
    author: "bob",
    url: "https://github.com/rancher/terraform-provider-rancher2/pull/102",
    conventional_type: "fix",
    conventional_scope: null,
    addressed_issues: [],
    changed_product_files: ["rancher2/provider.go"],
  },
  {
    number: 103,
    title: "chore: update dependencies",
    author: "carol",
    url: "https://github.com/rancher/terraform-provider-rancher2/pull/103",
    conventional_type: "chore",
    conventional_scope: null,
    addressed_issues: [],
    changed_product_files: ["docs/index.md"],
  },
];

test("formatJSON", () => {
  const jsonStr = formatJSON(samplePRs);
  const parsed = JSON.parse(jsonStr);
  assert.strictEqual(parsed.length, 3);
  assert.strictEqual(parsed[0].number, 101);
});

test("formatMarkdown", async (t) => {
  await t.test("generates grouped sections by conventional commit type", () => {
    const md = formatMarkdown(samplePRs);
    assert.ok(md.includes("### Features"));
    assert.ok(md.includes("### Bug Fixes"));
    assert.ok(md.includes("### Other Changes"));
    assert.ok(
      md.includes(
        "- feat(cluster): add proxy configuration ([#101](https://github.com/rancher/terraform-provider-rancher2/pull/101)) (Addresses #1001) by @alice",
      ),
    );
    assert.ok(
      md.includes(
        "- fix: nil pointer check on delete ([#102](https://github.com/rancher/terraform-provider-rancher2/pull/102)) by @bob",
      ),
    );
  });

  await t.test("handles empty list gracefully", () => {
    const md = formatMarkdown([]);
    assert.ok(md.includes("No merged pull requests found"));
  });
});

test("formatTable", async (t) => {
  await t.test("generates markdown table headers and rows", () => {
    const tbl = formatTable(samplePRs);
    assert.ok(tbl.includes("| PR   | Type"));
    assert.ok(tbl.includes("| #101 | feat"));
    assert.ok(tbl.includes("| #102 | fix"));
  });

  await t.test("escapes pipe characters in pull request titles", () => {
    const prWithPipe = [
      {
        number: 201,
        title: "feat(cluster): proxy | v2 configuration",
        author: "dave",
        url: "https://github.com/rancher/terraform-provider-rancher2/pull/201",
        conventional_type: "feat",
        conventional_scope: "cluster",
        addressed_issues: [1002],
        changed_product_files: ["rancher2/resource_cluster.go"],
      },
    ];
    const tbl = formatTable(prWithPipe);
    assert.ok(tbl.includes("proxy \\| v2 configuration"));

    // Verify row column count matches header column count
    const lines = tbl.trim().split("\n");
    assert.strictEqual(lines.length, 3); // header, divider, row
    const headerCols = lines[0].split("|").length;
    // Check that there are 6 data columns (surrounded by leading and trailing pipes = 8 segments)
    assert.strictEqual(
      lines[2].split("\\|").join("").split("|").length,
      headerCols,
    );
  });

  await t.test("handles empty list gracefully", () => {
    const tbl = formatTable([]);
    assert.ok(tbl.includes("No merged pull requests found"));
  });
});
