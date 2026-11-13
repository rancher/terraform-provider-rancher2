/**
 * Output formatters for merged pull requests (JSON, Markdown, and tabular).
 */

/**
 * Formats PR array as pretty JSON string.
 *
 * @param {Array<object>} prs
 * @param {object} [options]
 * @param {number} [options.indent=2]
 * @returns {string}
 */
export function formatJSON(prs, { indent = 2 } = {}) {
  return JSON.stringify(prs, null, indent);
}

const TYPE_SECTION_HEADERS = {
  feat: "### Features",
  fix: "### Bug Fixes",
  docs: "### Documentation",
  perf: "### Performance Improvements",
  refactor: "### Refactoring",
};

/**
 * Formats PR array as a categorized Markdown changelog.
 *
 * @param {Array<object>} prs
 * @returns {string}
 */
export function formatMarkdown(prs) {
  if (!Array.isArray(prs) || prs.length === 0) {
    return "No merged pull requests found matching criteria.\n";
  }

  const grouped = {
    feat: [],
    fix: [],
    docs: [],
    perf: [],
    refactor: [],
    other: [],
  };

  for (const pr of prs) {
    const type = pr.conventional_type
      ? pr.conventional_type.toLowerCase()
      : null;
    if (type && grouped[type]) {
      grouped[type].push(pr);
    } else {
      grouped.other.push(pr);
    }
  }

  const sections = [];

  for (const [type, header] of Object.entries(TYPE_SECTION_HEADERS)) {
    const list = grouped[type];
    if (list && list.length > 0) {
      const items = list.map((pr) => formatMarkdownItem(pr)).join("\n");
      sections.push(`${header}\n\n${items}`);
    }
  }

  if (grouped.other.length > 0) {
    const items = grouped.other.map((pr) => formatMarkdownItem(pr)).join("\n");
    sections.push(`### Other Changes\n\n${items}`);
  }

  return sections.join("\n\n") + "\n";
}

function formatMarkdownItem(pr) {
  const prLink = pr.url ? `[#${pr.number}](${pr.url})` : `#${pr.number}`;
  const authorSuffix = pr.author ? ` by @${pr.author}` : "";
  const issuesSuffix =
    Array.isArray(pr.addressed_issues) && pr.addressed_issues.length > 0
      ? ` (Addresses #${pr.addressed_issues.join(", #")})`
      : "";
  return `- ${pr.title} (${prLink})${issuesSuffix}${authorSuffix}`;
}

/**
 * Formats PR array as an aligned text / markdown table.
 *
 * @param {Array<object>} prs
 * @returns {string}
 */
export function formatTable(prs) {
  if (!Array.isArray(prs) || prs.length === 0) {
    return "No merged pull requests found matching criteria.\n";
  }

  const headers = ["PR", "Type", "Title", "Author", "Issues", "Product Files"];
  const rows = prs.map((pr) => {
    const rawTitle = pr.title || "";
    const title =
      rawTitle.length > 50 ? `${rawTitle.slice(0, 47)}...` : rawTitle;
    return [
      `#${pr.number}`,
      pr.conventional_type || "-",
      title,
      pr.author || "-",
      Array.isArray(pr.addressed_issues) && pr.addressed_issues.length > 0
        ? pr.addressed_issues.map((i) => `#${i}`).join(", ")
        : "-",
      String(
        Array.isArray(pr.changed_product_files)
          ? pr.changed_product_files.length
          : 0,
      ),
    ].map((cell) => cell.replaceAll("|", "\\|"));
  });

  const colWidths = headers.map((h, i) => {
    const maxRow = Math.max(...rows.map((r) => r[i].length));
    return Math.max(h.length, maxRow);
  });

  const formatRow = (cols) =>
    "| " + cols.map((c, i) => c.padEnd(colWidths[i])).join(" | ") + " |";

  const divider = "| " + colWidths.map((w) => "-".repeat(w)).join(" | ") + " |";

  return (
    [formatRow(headers), divider, ...rows.map(formatRow)].join("\n") + "\n"
  );
}
