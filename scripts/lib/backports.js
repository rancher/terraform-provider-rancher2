/**
 * Backport branch parsing, discovery, and correlation engine.
 */

import { extractChangesPath } from "./filter.js";

/**
 * Parses main PR number and target branch from a backport head branch reference.
 * E.g., "backport-2633-release-v15" -> { mainPrNumber: 2633, branchSuffix: "release-v15", targetBranch: "release/v15" }
 *
 * @param {string} headRef
 * @returns {{ mainPrNumber: number, branchSuffix: string, targetBranch: string } | null}
 */
export function parseBackportBranch(headRef) {
  if (!headRef || typeof headRef !== "string") {
    return null;
  }

  const match = headRef.trim().match(/^backport-(\d+)-(.+)$/);
  if (!match) {
    return null;
  }

  const mainPrNumber = parseInt(match[1], 10);
  const branchSuffix = match[2];

  // Reconstruct standard release branch name (e.g. "release-v15" -> "release/v15")
  let targetBranch = branchSuffix;
  if (branchSuffix.startsWith("release-")) {
    targetBranch = `release/${branchSuffix.slice("release-".length)}`;
  }

  return {
    mainPrNumber,
    branchSuffix,
    targetBranch,
  };
}

/**
 * Builds search query string for a main PR number.
 * E.g., "repo:rancher/terraform-provider-rancher2 is:pr head:backport-123-"
 * Note: GitHub Search qualifiers (such as head:) do not support boolean OR grouping.
 *
 * @param {object} options
 * @param {string} options.owner
 * @param {string} options.repo
 * @param {number|Array<number>} [options.prNumber]
 * @param {Array<number>} [options.prNumbers]
 * @returns {string}
 */
export function buildBackportSearchQuery({ owner, repo, prNumber, prNumbers }) {
  const num =
    prNumber ||
    (Array.isArray(prNumbers) && prNumbers.length > 0 ? prNumbers[0] : null);
  if (!num) {
    return `repo:${owner}/${repo} is:pr head:backport-`;
  }

  return `repo:${owner}/${repo} is:pr head:backport-${num}-`;
}

/**
 * Extracts attached backport issue number from a pull request description.
 *
 * @param {string} body
 * @returns {number|null}
 */
export function extractAttachedIssueNumber(body) {
  if (!body || typeof body !== "string") {
    return null;
  }
  // Pattern 1: Addresses #<issue> for #<tracking> (supporting URLs or #)
  const mFor = body.match(
    /Addresses\s+(?:https?:\/\/[^\s\)]+\/issues\/|#)(\d+)\s+for\s+(?:https?:\/\/[^\s\)]+\/issues\/|#)?(\d+)/i,
  );
  if (mFor) {
    return parseInt(mFor[1], 10);
  }
  // Pattern 2: Addresses / Fixes / Closes with URL or #
  const mGen = body.match(
    /(?:Addresses|Fixes|Closes|Resolves):?\s*#?\s*(?:https?:\/\/[^\s\)]+\/issues\/|#)(\d+)/i,
  );
  if (mGen) {
    return parseInt(mGen[1], 10);
  }
  return null;
}

/**
 * Extracts correlation information from a backport PR (supporting both full PR objects and /search/issues objects).
 *
 * @param {object} bpr
 * @returns {{ mainPrNumber: number, targetBranch: string, prNumber: number, state: string, mergedAt: string | null, mergeCommitSha: string | null, url: string, attachedIssueNumber: number | null } | null}
 */
export function extractBackportInfo(bpr) {
  if (!bpr) {
    return null;
  }

  let mainPrNumber = null;
  let targetBranch = bpr.base?.ref || null;

  // 1. Try head.ref (present in full PR objects)
  if (bpr.head && bpr.head.ref) {
    const parsed = parseBackportBranch(bpr.head.ref);
    if (parsed) {
      mainPrNumber = parsed.mainPrNumber;
      if (!targetBranch) {
        targetBranch = parsed.targetBranch;
      }
    }
  }

  // 2. Try body signatures (automated cherry-pick or manual backport reference)
  if (!mainPrNumber && bpr.body) {
    const match = bpr.body.match(
      /cherry-picks the changes from #(\d+)\b(?:\s+into\s+([^\s\r\n]+))?/i,
    );
    if (match) {
      mainPrNumber = parseInt(match[1], 10);
      if (!targetBranch && match[2]) {
        targetBranch = match[2].trim();
      }
    }

    if (!mainPrNumber) {
      const manualMatch = bpr.body.match(
        /(?:backport(?:ed|ing)?(?:\s+is|\s+of|\s+from)?|cherry-pick(?:ed)?(?:\s+from)?)\s+(?:https?:\/\/[^\s\)]+\/(?:pull|pulls)\/|#)(\d+)\b/i,
      );
      if (manualMatch) {
        mainPrNumber = parseInt(manualMatch[1], 10);
      }
    }
  }

  // 3. Try title signatures (e.g. "[Backport #2311]" or "backport #2311")
  if (!mainPrNumber && bpr.title) {
    const titleMatch = bpr.title.match(
      /(?:backport(?:ed|ing)?(?:\s+is|\s+of|\s+from|\s+for)?|cherry-pick(?:ed)?(?:\s+from)?)\s+(?:https?:\/\/[^\s\)]+\/(?:pull|pulls)\/|#)(\d+)\b/i,
    );
    if (titleMatch) {
      mainPrNumber = parseInt(titleMatch[1], 10);
    }
  }

  // 4. Fallback target branch from labels (e.g. release/v15)
  if (!targetBranch && Array.isArray(bpr.labels)) {
    const branchLabel = bpr.labels.find((l) => {
      const name = typeof l === "string" ? l : l.name || "";
      return name.startsWith("release/");
    });
    if (branchLabel) {
      targetBranch =
        typeof branchLabel === "string" ? branchLabel : branchLabel.name;
    }
  }

  if (!mainPrNumber) {
    return null;
  }

  const mergedAt = bpr.merged_at || bpr.pull_request?.merged_at || null;
  let state = bpr.state || "open";
  if (mergedAt || bpr.merged) {
    state = "merged";
  }

  const attachedIssueNumber = extractAttachedIssueNumber(bpr.body);

  return {
    mainPrNumber,
    targetBranch: targetBranch || "unknown",
    prNumber: bpr.number,
    state,
    mergedAt,
    mergeCommitSha: bpr.merge_commit_sha || null,
    url: bpr.html_url || bpr.url || "",
    attachedIssueNumber,
  };
}

/**
 * Correlates raw backport pull request records with main pull request objects.
 *
 * @param {Array<object>} mainPRs
 * @param {Array<object>} backportPRs
 * @param {object} [options]
 * @param {Array<string>} [options.expectedBranches=[]]
 * @returns {Array<object>} Enriched main PR objects
 */
export function correlateBackports(
  mainPRs,
  backportPRs,
  { expectedBranches = [] } = {},
) {
  if (!Array.isArray(mainPRs)) {
    return [];
  }

  // Index backports by main PR number
  const backportMap = new Map();

  for (const bpr of backportPRs || []) {
    const info = extractBackportInfo(bpr);
    if (!info) {
      continue;
    }

    const record = {
      target_branch: info.targetBranch,
      pr_number: info.prNumber,
      state: info.state,
      merged_at: info.mergedAt,
      merge_commit_sha: info.mergeCommitSha,
      url: info.url,
      attached_issue_number: info.attachedIssueNumber || null,
      attached_issue: bpr.qa_issue || bpr.attached_issue || null,
      qa_issue: bpr.qa_issue || bpr.attached_issue || null,
    };

    if (!backportMap.has(info.mainPrNumber)) {
      backportMap.set(info.mainPrNumber, []);
    }

    const existing = backportMap.get(info.mainPrNumber);
    if (!existing.some((r) => r.pr_number === info.prNumber)) {
      existing.push(record);
    }
  }

  return mainPRs.map((mainPr) => {
    const list = backportMap.get(mainPr.number) || [];
    const completedBranches = list
      .filter((b) => b.state === "merged")
      .map((b) => b.target_branch);

    const missingBranches = expectedBranches.filter(
      (eb) => !completedBranches.includes(eb),
    );

    return {
      ...mainPr,
      backports: list,
      backport_status: {
        expected_branches: expectedBranches,
        completed_branches: completedBranches,
        missing_branches: missingBranches,
      },
    };
  });
}

/**
 * Filters enriched PRs to only those that have at least one open QA issue.
 *
 * @param {Array<object>} enrichedPRs
 * @returns {Array<object>}
 */
export function filterOpenQAIssues(enrichedPRs) {
  if (!Array.isArray(enrichedPRs)) {
    return [];
  }
  return enrichedPRs.filter((pr) => {
    const qaIssues = (pr.backports || [])
      .map((bp) => bp.qa_issue || bp.attached_issue)
      .filter(Boolean);
    return qaIssues.some((iss) => iss.state === "open");
  });
}

// Alias for backwards compatibility
export const filterOpenBackportIssues = filterOpenQAIssues;

/**
 * Formats enriched PRs as a Markdown backport status matrix.
 *
 * @param {Array<object>} enrichedPRs
 * @param {Array<string>} [targetBranches]
 * @returns {string}
 */
export function formatBackportMatrix(enrichedPRs, targetBranches = []) {
  if (!Array.isArray(enrichedPRs) || enrichedPRs.length === 0) {
    return "No pull requests to display in backport matrix.\n";
  }

  // Determine all branches present across all PRs if not provided
  let branches = targetBranches;
  if (!branches || branches.length === 0) {
    const branchSet = new Set();
    for (const pr of enrichedPRs) {
      for (const bp of pr.backports || []) {
        if (bp.target_branch) {
          branchSet.add(bp.target_branch);
        }
      }
      for (const exp of pr.backport_status?.expected_branches || []) {
        branchSet.add(exp);
      }
    }
    branches = Array.from(branchSet);
  }

  const headers = ["Main PR", "Title", ...branches];
  const rows = enrichedPRs.map((pr) => {
    const prCell = pr.url ? `[#${pr.number}](${pr.url})` : `#${pr.number}`;
    const titleCell =
      pr.title.length > 40 ? `${pr.title.slice(0, 37)}...` : pr.title;

    const branchCells = branches.map((branch) => {
      const match = (pr.backports || []).find(
        (b) => b.target_branch === branch,
      );
      if (!match) {
        return "❌ Missing";
      }
      const link = match.url
        ? `[#${match.pr_number}](${match.url})`
        : `#${match.pr_number}`;
      if (match.state === "merged") {
        return `✅ ${link} (Merged)`;
      }
      if (match.state === "open") {
        return `⏳ ${link} (Open)`;
      }
      return `⚠️ ${link} (${match.state})`;
    });

    return [prCell, titleCell, ...branchCells].map((cell) =>
      cell.replaceAll("|", "\\|"),
    );
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

/**
 * Formats enriched PRs as a markdown tracking table with columns:
 * | PR# | Author | Tracking Issue # | Backport PRs |
 *
 * @param {Array<object>} enrichedPRs
 * @param {object} [options]
 * @param {string} [options.owner="rancher"]
 * @param {string} [options.repo="terraform-provider-rancher2"]
 * @returns {string}
 */
export function formatBackportTable(
  enrichedPRs,
  { owner = "rancher", repo = "terraform-provider-rancher2" } = {},
) {
  if (!Array.isArray(enrichedPRs) || enrichedPRs.length === 0) {
    return "No pull requests to display in table.\n";
  }

  const headers = [
    "Title",
    "PR#",
    "Author",
    "Tracking Issue #",
    "Backport PRs",
    "QA Issues",
    "ChangesPath",
  ];
  const rows = enrichedPRs.map((pr) => {
    // 1. Title
    const titleCell = pr.title || "-";

    // 2. PR#
    const prUrl =
      pr.url || `https://github.com/${owner}/${repo}/pull/${pr.number}`;
    const mainState = pr.state || (pr.merged_at ? "merged" : "open");
    const mainBranch = pr.target_branch || pr.base?.ref || "main";
    const prCell = `[${pr.number} - ${mainState} ${mainBranch}](${prUrl})`;

    // 3. Author
    const authorCell = pr.author || "-";

    // 4. Tracking Issue #
    let trackingCell = "-";
    if (Array.isArray(pr.addressed_issues) && pr.addressed_issues.length > 0) {
      trackingCell = pr.addressed_issues
        .map((issue) => {
          const issueNum = typeof issue === "object" ? issue.number : issue;
          const issueState = typeof issue === "object" ? issue.state : null;
          const linkText = issueState
            ? `${issueNum} - ${issueState}`
            : String(issueNum);
          return `[${linkText}](https://github.com/${owner}/${repo}/issues/${issueNum})`;
        })
        .join(", ");
    }

    // 5. Backport PRs
    let backportsCell = "-";
    if (Array.isArray(pr.backports) && pr.backports.length > 0) {
      backportsCell = pr.backports
        .map((bp) => {
          const bpUrl =
            bp.url ||
            `https://github.com/${owner}/${repo}/pull/${bp.pr_number}`;
          const bpState = bp.state || (bp.merged_at ? "merged" : "open");
          const bpBranch = bp.target_branch || "unknown";
          return `[${bp.pr_number} - ${bpState} ${bpBranch}](${bpUrl})`;
        })
        .join(", ");
    }

    // 6. QA Issues
    let qaIssuesCell = "-";
    const qaIssues = (pr.backports || [])
      .map((bp) => bp.qa_issue || bp.attached_issue)
      .filter(Boolean);

    const uniqueQAIssues = [];
    for (const iss of qaIssues) {
      if (!uniqueQAIssues.some((u) => u.number === iss.number)) {
        uniqueQAIssues.push(iss);
      }
    }

    if (uniqueQAIssues.length > 0) {
      qaIssuesCell = uniqueQAIssues
        .map((iss) => {
          const person = iss.person || "unassigned";
          const linkText = `${iss.number} - ${iss.state} - ${person}`;
          const url =
            iss.url ||
            `https://github.com/${owner}/${repo}/issues/${iss.number}`;
          return `[${linkText}](${url})`;
        })
        .join(", ");
    }

    // 7. ChangesPath
    const files = pr.changed_product_files || pr.files || [];
    const changesPathCell = pr.changes_path || extractChangesPath(files);

    return [
      titleCell,
      prCell,
      authorCell,
      trackingCell,
      backportsCell,
      qaIssuesCell,
      changesPathCell,
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

/**
 * Formats a Gap Report highlighting missing or incomplete backports.
 *
 * @param {Array<object>} enrichedPRs
 * @returns {string}
 */
export function formatBackportGaps(enrichedPRs) {
  if (!Array.isArray(enrichedPRs) || enrichedPRs.length === 0) {
    return "No pull requests evaluated.\n";
  }

  const gaps = [];

  for (const pr of enrichedPRs) {
    const missing = pr.backport_status?.missing_branches || [];
    const openBackports = (pr.backports || []).filter(
      (b) => b.state === "open",
    );

    if (missing.length > 0 || openBackports.length > 0) {
      gaps.push({ pr, missing, openBackports });
    }
  }

  if (gaps.length === 0) {
    return "✅ All pull requests have completed backports across all expected branches.\n";
  }

  const lines = ["# Backport Gap Report\n"];

  for (const { pr, missing, openBackports } of gaps) {
    lines.push(`### Main PR #${pr.number}: ${pr.title}`);
    if (pr.url) {
      lines.push(`URL: ${pr.url}`);
    }
    if (missing.length > 0) {
      lines.push(`- **Missing Target Branches**: ${missing.join(", ")}`);
    }
    if (openBackports.length > 0) {
      const openLinks = openBackports
        .map((b) => `${b.target_branch} (${b.url || `#${b.pr_number}`})`)
        .join(", ");
      lines.push(`- **Open Backports Pending Merge**: ${openLinks}`);
    }
    lines.push("");
  }

  return lines.join("\n");
}
