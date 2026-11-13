#!/usr/bin/env node
import fs from "node:fs/promises";
import { parseArgs } from "node:util";

import { getAuthToken } from "./lib/auth.js";
import {
  buildBackportSearchQuery,
  correlateBackports,
  extractAttachedIssueNumber,
  filterOpenQAIssues,
  formatBackportGaps,
  formatBackportMatrix,
  formatBackportTable,
} from "./lib/backports.js";
import { formatJSON } from "./lib/formatters.js";
import {
  fetchBranchPRs,
  fetchIssueDetails,
  fetchIssueState,
  fetchPullRequest,
  fetchReleaseBranches,
  fetchTrackingIssuesMap,
  searchPullRequests,
} from "./lib/github.js";
import { extractAddressedIssues } from "./lib/parser.js";

export function printHelp() {
  console.log(`Usage: get-backport-prs.js [options]

Cross-references main pull requests with automated backport PRs across maintenance branches.
Consumes output from get-merged-prs.js via stdin or --input file.

Options:
  --input <path>               Path to JSON file containing main PRs
  --pr <number>                Specify main PR number directly (repeatable or comma-separated)
  --expected-branches <names>  Comma-separated list of expected maintenance branches (e.g. release/v15,release/v14)
  --open-qa-issues-only        List only rows that have QA issues that are open (alias: --open-issues-only)
  --format <table|matrix|gap|json> Output format: table, matrix, gap, or json (default: json)
  --output <path>              Write output directly to file rather than stdout
  --owner <owner>              GitHub repository owner (default: rancher)
  --repo <repo>                GitHub repository name (default: terraform-provider-rancher2)
  -h, --help                   Show this help message
`);
}

/**
 * Normalizes list option.
 *
 * @param {string|Array<string>|undefined} opt
 * @returns {Array<string>}
 */
export function normalizeListOption(opt) {
  if (!opt) {
    return [];
  }
  const rawList = Array.isArray(opt) ? opt : [opt];
  return rawList
    .flatMap((item) => item.split(","))
    .map((s) => s.trim())
    .filter(Boolean);
}

/**
 * Reads all data from standard input if piped.
 *
 * @returns {Promise<string>}
 */
export async function readStdin() {
  if (process.stdin.isTTY) {
    return "";
  }
  const chunks = [];
  for await (const chunk of process.stdin) {
    chunks.push(chunk);
  }
  return Buffer.concat(chunks).toString("utf8").trim();
}

/**
 * Main execution pipeline for get-backport-prs.
 *
 * @param {Array<string>} [args=process.argv.slice(2)]
 */
export async function runBackportPRs(args = process.argv.slice(2)) {
  const { values } = parseArgs({
    args,
    options: {
      input: { type: "string" },
      pr: { type: "string", multiple: true },
      "expected-branches": { type: "string" },
      "open-qa-issues-only": { type: "boolean", default: false },
      "open-issues-only": { type: "boolean", default: false },
      "only-open-issues": { type: "boolean", default: false },
      format: { type: "string", default: "json" },
      output: { type: "string" },
      owner: { type: "string", default: "rancher" },
      repo: { type: "string", default: "terraform-provider-rancher2" },
      help: { type: "boolean", short: "h" },
    },
    allowPositionals: true,
  });

  if (values.help) {
    printHelp();
    return;
  }

  const owner = values.owner;
  const repo = values.repo;
  const format = (values.format || "json").toLowerCase();
  const expectedBranches = normalizeListOption(values["expected-branches"]);
  const openQAIssuesOnly =
    Boolean(values["open-qa-issues-only"]) ||
    Boolean(values["open-issues-only"]) ||
    Boolean(values["only-open-issues"]);
  const token = await getAuthToken();

  // Ingest main PR array
  let mainPRs = [];

  if (values.input) {
    const fileContent = await fs.readFile(values.input, "utf8");
    mainPRs = JSON.parse(fileContent);
  } else if (values.pr) {
    const prNums = normalizeListOption(values.pr)
      .map((n) => parseInt(n, 10))
      .filter(Boolean);
    for (const num of prNums) {
      try {
        const prData = await fetchPullRequest({
          owner,
          repo,
          prNumber: num,
          token,
        });
        mainPRs.push({
          number: num,
          title: prData.title || `PR #${num}`,
          author:
            prData.user?.login || prData.author?.login || prData.author || "",
          state: prData.merged_at ? "merged" : prData.state || "open",
          target_branch: prData.base?.ref || "main",
          url:
            prData.html_url ||
            `https://github.com/${owner}/${repo}/pull/${num}`,
          addressed_issues: extractAddressedIssues(
            prData.body || "",
            `${owner}/${repo}`,
          ),
        });
      } catch {
        mainPRs.push({
          number: num,
          title: `PR #${num}`,
          author: "",
          state: "open",
          target_branch: "main",
          url: `https://github.com/${owner}/${repo}/pull/${num}`,
          addressed_issues: [],
        });
      }
    }
  } else {
    const stdinContent = await readStdin();
    if (stdinContent) {
      mainPRs = JSON.parse(stdinContent);
    }
  }

  if (!Array.isArray(mainPRs)) {
    throw new Error(
      "No main PRs provided. Pipe JSON from get-merged-prs.js, provide --input <file>, or specify --pr <number>.",
    );
  }

  if (mainPRs.length === 0) {
    const emptyMsg =
      format === "table" || format === "markdown" || format === "md"
        ? "No pull requests to display in table.\n"
        : "[]\n";
    if (values.output) {
      await fs.writeFile(values.output, emptyMsg, "utf8");
    } else {
      process.stdout.write(emptyMsg);
    }
    return;
  }

  const mainNumbers = mainPRs.map((p) => p.number).filter(Boolean);

  // Discover backport PRs via Core REST API (5,000 req/hr rate limit)
  let targetBranches = expectedBranches;
  if (targetBranches.length === 0) {
    const allReleaseBranches = await fetchReleaseBranches({
      owner,
      repo,
      token,
    });
    targetBranches = allReleaseBranches.slice(0, 5);
  }

  const branchPRLists = await Promise.all(
    targetBranches.map((branch) =>
      fetchBranchPRs({ owner, repo, baseBranch: branch, token }).catch(
        (err) => {
          console.warn(
            `⚠️ Error fetching PRs for branch ${branch}:`,
            err.message,
          );
          return [];
        },
      ),
    ),
  );
  const discoveredBackports = branchPRLists.flat();

  // If a main PR has no backports found via branch search and only 1-2 PRs queried, try search fallback
  if (discoveredBackports.length === 0 && mainNumbers.length <= 2) {
    for (const num of mainNumbers) {
      try {
        const query = buildBackportSearchQuery({ owner, repo, prNumber: num });
        const items = await searchPullRequests({ query, token });
        discoveredBackports.push(...items);
      } catch {
        // Search API might be rate-limited; ignore fallback errors
      }
    }
  }

  // Resolve tracking issues for main PRs
  const trackingMap = await fetchTrackingIssuesMap({ owner, repo, token });
  const issuesNeedingState = new Set();

  for (const pr of mainPRs) {
    const trackingInfo = trackingMap.get(pr.number);
    const existingIssues = Array.isArray(pr.addressed_issues)
      ? pr.addressed_issues
      : [];
    const normalizedList = [];

    if (trackingInfo) {
      normalizedList.push(
        typeof trackingInfo === "object"
          ? trackingInfo
          : { number: trackingInfo, state: "open" },
      );
    }

    for (const item of existingIssues) {
      const num = typeof item === "object" ? item.number : item;
      const state =
        typeof item === "object"
          ? item.state
          : trackingMap.issueStates?.get(num) || null;
      if (!normalizedList.some((iss) => iss.number === num)) {
        normalizedList.push({ number: num, state });
        if (!state) {
          issuesNeedingState.add(num);
        }
      }
    }
    pr.addressed_issues = normalizedList;
  }

  // Lookup state for any remaining issues without state
  if (issuesNeedingState.size > 0) {
    const unknownList = Array.from(issuesNeedingState);
    const resolvedStates = await Promise.all(
      unknownList.map(async (num) => {
        const state = await fetchIssueState({
          owner,
          repo,
          issueNumber: num,
          token,
        });
        return { num, state };
      }),
    );
    const stateLookup = new Map(resolvedStates.map((r) => [r.num, r.state]));
    for (const pr of mainPRs) {
      if (Array.isArray(pr.addressed_issues)) {
        for (const iss of pr.addressed_issues) {
          if (!iss.state && stateLookup.has(iss.number)) {
            iss.state = stateLookup.get(iss.number);
          }
        }
      }
    }
  }

  // Resolve attached issues for discovered backport PRs
  const bpIssueNumbers = new Set();
  for (const bpr of discoveredBackports) {
    const issNum = extractAttachedIssueNumber(bpr.body);
    if (issNum) {
      bpr.attached_issue_number = issNum;
      bpIssueNumbers.add(issNum);
    }
  }

  if (bpIssueNumbers.size > 0) {
    const issueList = Array.from(bpIssueNumbers);
    const resolvedIssues = await Promise.all(
      issueList.map((issNum) =>
        fetchIssueDetails({ owner, repo, issueNumber: issNum, token }).catch(
          () => ({
            number: issNum,
            state: "unknown",
            person: "unknown",
            url: `https://github.com/${owner}/${repo}/issues/${issNum}`,
          }),
        ),
      ),
    );
    const bpIssueDetailsMap = new Map(
      resolvedIssues.map((item) => [item.number, item]),
    );

    for (const bpr of discoveredBackports) {
      if (
        bpr.attached_issue_number &&
        bpIssueDetailsMap.has(bpr.attached_issue_number)
      ) {
        bpr.qa_issue = bpr.attached_issue = bpIssueDetailsMap.get(
          bpr.attached_issue_number,
        );
      }
    }
  }

  // Correlate backports
  let enrichedPRs = correlateBackports(mainPRs, discoveredBackports, {
    expectedBranches,
  });

  if (openQAIssuesOnly) {
    enrichedPRs = filterOpenQAIssues(enrichedPRs);
  }

  // Format output
  let outputData = "";
  if (format === "table" || format === "markdown" || format === "md") {
    outputData = formatBackportTable(enrichedPRs, { owner, repo });
  } else if (format === "matrix") {
    outputData = formatBackportMatrix(enrichedPRs, expectedBranches);
  } else if (format === "gap" || format === "gaps") {
    outputData = formatBackportGaps(enrichedPRs);
  } else {
    outputData = formatJSON(enrichedPRs) + "\n";
  }

  if (values.output) {
    await fs.writeFile(values.output, outputData, "utf8");
  } else {
    process.stdout.write(outputData);
  }
}

async function main() {
  try {
    await runBackportPRs();
  } catch (err) {
    console.error(`❌ Error: ${err?.stack || err?.message || err}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
