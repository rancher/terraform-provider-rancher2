import assert from "node:assert/strict";
import test from "node:test";

import {
  extractAddressedIssues,
  parseConventionalTitle,
} from "../lib/parser.js";

test("parseConventionalTitle", async (t) => {
  await t.test("parses standard conventional title with scope", () => {
    const res = parseConventionalTitle(
      "feat(cluster): add proxy configuration",
    );
    assert.strictEqual(res.type, "feat");
    assert.strictEqual(res.scope, "cluster");
    assert.strictEqual(res.breaking, false);
    assert.strictEqual(res.description, "add proxy configuration");
  });

  await t.test("parses title without scope", () => {
    const res = parseConventionalTitle("fix: handle nil pointer dereference");
    assert.strictEqual(res.type, "fix");
    assert.strictEqual(res.scope, null);
    assert.strictEqual(res.breaking, false);
    assert.strictEqual(res.description, "handle nil pointer dereference");
  });

  await t.test("parses breaking change with exclamation mark", () => {
    const res = parseConventionalTitle(
      "feat(auth)!: remove legacy v1 token endpoint",
    );
    assert.strictEqual(res.type, "feat");
    assert.strictEqual(res.scope, "auth");
    assert.strictEqual(res.breaking, true);
    assert.strictEqual(res.description, "remove legacy v1 token endpoint");
  });

  await t.test("handles non-conventional title gracefully", () => {
    const res = parseConventionalTitle(
      "Update README documentation and examples",
    );
    assert.strictEqual(res.type, null);
    assert.strictEqual(res.scope, null);
    assert.strictEqual(res.breaking, false);
    assert.strictEqual(
      res.description,
      "Update README documentation and examples",
    );
  });

  await t.test("handles null or empty title", () => {
    assert.deepStrictEqual(parseConventionalTitle(""), {
      type: null,
      scope: null,
      breaking: false,
      description: "",
    });
    assert.deepStrictEqual(parseConventionalTitle(null), {
      type: null,
      scope: null,
      breaking: false,
      description: "",
    });
  });
});

test("extractAddressedIssues", async (t) => {
  /**
   * Normalization Logic Rationale:
   * 1. Keywords: 'Addresses', 'Fixes', 'Closes', 'Resolves' are standard PR issue keywords.
   * 2. Patterns: Both short-hand (#123) and full GitHub issue URLs are captured.
   * 3. Normalization: URLs are stripped to extract only the numeric ID to maintain consistency with short-hand references.
   * 4. Deduplication & Sorting: Ensures the output is deterministic and unique, critical for backport tracking.
   * 5. Repo Scoping: Only issues within the target repository are resolved; external repository links are intentionally ignored.
   */
  await t.test("extracts issue from standard Addresses pattern", () => {
    const body =
      "This PR introduces cluster proxy settings.\n\n- Addresses: #1234\n";
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [1234]);
  });

  await t.test("handles unexpected input types safely", () => {
    assert.deepStrictEqual(extractAddressedIssues(null), []);
    assert.deepStrictEqual(extractAddressedIssues(undefined), []);
    assert.deepStrictEqual(extractAddressedIssues({}), []);
    assert.deepStrictEqual(extractAddressedIssues(123), []);
  });

  await t.test("extracts issues from complex nested lists", () => {
    const body = `
    - Feature:
      - Add proxy support
        - Fixes: #10
      - Improve performance
    - Addresses: #20
    `;
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [10, 20]);
  });

  await t.test("extracts issues from Markdown tables", () => {
    const body = `
    | ID | Description |
    |----|-------------|
    | 1  | Bug fix     |

    Addresses: #50
    | #60 | Refactor    |
    `;
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [50, 60]);
  });

  await t.test("extracts issues from complex URL variations", () => {
    const body = `
      - Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/100
      - Fixes: http://github.com/rancher/terraform-provider-rancher2/issues/200/
      - Closes: https://www.github.com/rancher/terraform-provider-rancher2/issues/300#issuecomment-1
      - Resolves: https://github.com/rancher/terraform-provider-rancher2/issues/400?foo=bar
    `;
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [100, 200, 300, 400]);
  });

  await t.test(
    "extracts multiple issues across multiple lines and keywords",
    () => {
      const body = `
      Some summary here.
      - Addresses: #100, #200
      Fixes: #300
      Closes #400
      Resolves: #500
    `;
      const issues = extractAddressedIssues(body);
      assert.deepStrictEqual(issues, [100, 200, 300, 400, 500]);
    },
  );

  await t.test("deduplicates issues and sorts ascending", () => {
    const body = "Addresses: #456\nCloses: #123\nResolves #456";
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [123, 456]);
  });

  await t.test("extracts issue from full GitHub issue URL", () => {
    const body =
      "This PR fixes a bug.\n\n- Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/123\n";
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, [123]);
  });

  await t.test("ignores malformed issue URLs lacking numeric IDs", () => {
    const body =
      "Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/\nFixes: https://github.com/rancher/terraform-provider-rancher2/pull/";
    const issues = extractAddressedIssues(body);
    assert.deepStrictEqual(issues, []);
  });

  await t.test(
    "extracts multiple issues mixing URLs and #IDs across lines",
    () => {
      const body = `
      Summary here.
      Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/100, #200
      Fixes: https://github.com/rancher/terraform-provider-rancher2/issues/300
      Closes: [Issue 400](https://github.com/rancher/terraform-provider-rancher2/issues/400)
    `;
      const issues = extractAddressedIssues(body);
      assert.deepStrictEqual(issues, [100, 200, 300, 400]);
    },
  );

  await t.test(
    "ignores issue URLs and references from external repositories",
    () => {
      const body = `
      Addresses https://github.com/rancher/rancher/issues/12345
      Fixes rancher/rancher#67890
      Closes: https://github.com/rancher/terraform-provider-rancher2/issues/555
      Resolves #777
    `;
      const issues = extractAddressedIssues(body);
      assert.deepStrictEqual(issues, [555, 777]);
    },
  );

  await t.test(
    "supports custom targetRepo parameter as string or object",
    () => {
      const body = `
      Addresses https://github.com/my-org/my-repo/issues/123
      Fixes my-org/my-repo#456
      Closes https://github.com/rancher/terraform-provider-rancher2/issues/789
    `;
      const issuesString = extractAddressedIssues(body, "my-org/my-repo");
      assert.deepStrictEqual(issuesString, [123, 456]);

      const issuesObj = extractAddressedIssues(body, {
        owner: "my-org",
        repo: "my-repo",
      });
      assert.deepStrictEqual(issuesObj, [123, 456]);
    },
  );

  await t.test(
    "strips HTML comments before extracting issues to ignore template boilerplate",
    () => {
      const bodyWithComments = `
      ## Summary
      This PR fixes the ingress sync logic.

      <!-- - Addresses: #123 -->
      <!-- Fixes: https://github.com/rancher/terraform-provider-rancher2/issues/999 -->
      <!--
      Multi-line template guidance:
      - Closes: #888
      - Resolves: #777
      -->

      - Addresses: #456
      Fixes: #789
    `;
      const issues = extractAddressedIssues(bodyWithComments);
      assert.deepStrictEqual(issues, [456, 789]);
    },
  );

  await t.test(
    "returns empty array if issue references only exist within HTML comments",
    () => {
      const templateBoilerplateOnly = `
      ## Issue References
      <!-- - Addresses: #100 -->
      <!-- Fixes: #200 -->
    `;
      assert.deepStrictEqual(
        extractAddressedIssues(templateBoilerplateOnly),
        [],
      );
    },
  );

  await t.test(
    "strips Markdown code blocks before extracting issues to ignore code snippets",
    () => {
      const bodyWithCode = `
      This PR fixes the issue.
      Addresses: #123

      Here is some log output:
      \`\`\`
      Fixes: #456
      Addresses: https://github.com/rancher/terraform-provider-rancher2/issues/789
      \`\`\`

      And inline \`Closes: #999\`
    `;
      const issues = extractAddressedIssues(bodyWithCode);
      assert.deepStrictEqual(issues, [123]);
    },
  );

  await t.test("returns empty array when no issues mentioned", () => {
    assert.deepStrictEqual(
      extractAddressedIssues("Just a normal description"),
      [],
    );
    assert.deepStrictEqual(extractAddressedIssues(""), []);
  });
});
