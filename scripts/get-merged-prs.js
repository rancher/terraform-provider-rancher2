#!/usr/bin/env node
import fs from "node:fs/promises";
import { parseArgs } from "node:util";

import { getAuthToken, getAuthenticatedUser } from "./lib/auth.js";
import {
  DEFAULT_PRODUCT_PATHS,
  extractChangesPath,
  filterProductFiles,
  filterPullRequest,
} from "./lib/filter.js";
import { formatBackportTable } from "./lib/backports.js";
import { formatJSON, formatMarkdown, formatTable } from "./lib/formatters.js";
import { getChangedFilesFromCommit } from "./lib/git.js";
import {
  fetchClosedPRs,
  fetchPRFiles,
  fetchTagDate,
  fetchTrackingIssuesMap,
} from "./lib/github.js";
import {
  extractAddressedIssues,
  parseConventionalTitle,
} from "./lib/parser.js";

export function printHelp() {
  console.log(`Usage: get-merged-prs.js [options]

Retrieves pull requests merged into main and outputs structured metadata for release management.

Options:
  --exclude-me                 Exclude pull requests authored by authenticated user
  --no-exclude-me              Include pull requests authored by authenticated user (default)
  --exclude-author <login>     Exclude PRs authored by specific user (repeatable or comma-separated)
  --include-paths <paths>      Comma-separated list of product path prefixes (default: rancher2)
  --all-paths                  Disable path-based filtering (retrieve all PRs regardless of files)
  --since <tag|date>           Filter PRs merged after specific release tag or ISO date
  --until <tag|date>           Filter PRs merged before specific release tag or ISO date
  --milestone <title>          Restrict to PRs associated with milestone
  --label <name>               Restrict to PRs containing label (repeatable or comma-separated)
  --limit <n>                  Maximum PRs to return (default: 100)
  --format <json|markdown|table> Output format (default: json)
  --output <path>              Write output directly to file rather than stdout
  --owner <owner>              GitHub repository owner (default: rancher)
  --repo <repo>                GitHub repository name (default: terraform-provider-rancher2)
  -h, --help                   Show this help message
`);
}

/**
 * Normalizes comma-separated or array options into a flat array of clean strings.
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
 * Main execution pipeline for get-merged-prs.
 *
 * @param {Array<string>} [args=process.argv.slice(2)]
 */
export async function runMergedPRs(args = process.argv.slice(2)) {
  const { values } = parseArgs({
    args,
    options: {
      "exclude-me": { type: "boolean", default: false },
      "no-exclude-me": { type: "boolean" },
      "exclude-author": { type: "string", multiple: true },
      "include-paths": { type: "string" },
      "all-paths": { type: "boolean", default: false },
      since: { type: "string" },
      until: { type: "string" },
      milestone: { type: "string" },
      label: { type: "string", multiple: true },
      limit: { type: "string", default: "100" },
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
  const limit = parseInt(values.limit || "100", 10);
  const format = (values.format || "json").toLowerCase();
  const allPaths = Boolean(values["all-paths"]);
  const excludeMe = Boolean(values["exclude-me"]) && !values["no-exclude-me"];

  // Path boundaries
  const includePaths = values["include-paths"]
    ? normalizeListOption(values["include-paths"]).map((p) =>
      p.endsWith("/") ? p : `${p}/`,
    )
    : DEFAULT_PRODUCT_PATHS;

  // Resolve token and authenticated user
  const token = await getAuthToken();
  let me = null;
  if (excludeMe) {
    if (!token) {
      console.error(
        "Error: --exclude-me requires an authenticated user via GITHUB_TOKEN or 'gh auth login'",
      );
      process.exit(1);
    }
    me = await getAuthenticatedUser({ token });
    if (!me) {
      console.error(
        "Error: --exclude-me requires an authenticated user with a valid GitHub token",
      );
      process.exit(1);
    }
  }
  const excludedAuthors = normalizeListOption(values["exclude-author"]);

  if (excludeMe && me) {
    if (!excludedAuthors.some((a) => a.toLowerCase() === me.toLowerCase())) {
      excludedAuthors.push(me);
    }
  }

  // Resolve since/until tag dates if needed
  let since = values.since ? values.since.trim() : undefined;
  let until = values.until ? values.until.trim() : undefined;

  if (since && isNaN(new Date(since).getTime())) {
    const tagDate = await fetchTagDate({ owner, repo, tag: since, token });
    if (!tagDate) {
      throw new Error(`Invalid tag or date: ${since}`);
    }
    since = tagDate;
  }

  if (until && isNaN(new Date(until).getTime())) {
    const tagDate = await fetchTagDate({ owner, repo, tag: until, token });
    if (!tagDate) {
      throw new Error(`Invalid tag or date: ${until}`);
    }
    until = tagDate;
  }

  const requiredLabels = normalizeListOption(values.label);
  const milestone = values.milestone;

  // Fetch closed PRs
  const rawPRs = await fetchClosedPRs({
    owner,
    repo,
    base: "main",
    limit: limit * 2, // Fetch buffer to allow for filtered PRs
    since,
    token,
  });

  const trackingMap = await fetchTrackingIssuesMap({ owner, repo, token });
  const qualifiedPRs = [];

  for (const pr of rawPRs) {
    if (qualifiedPRs.length >= limit) {
      break;
    }

    if (!pr.merged_at) {
      continue;
    }

    // Resolve changed files: local git first, then API fallback
    let filePaths = [];
    if (pr.merge_commit_sha) {
      filePaths = await getChangedFilesFromCommit(pr.merge_commit_sha);
    }
    if (!filePaths || filePaths.length === 0) {
      filePaths = await fetchPRFiles({
        owner,
        repo,
        prNumber: pr.number,
        token,
      });
    }

    const evalResult = filterPullRequest(pr, {
      filePaths,
      includePaths,
      allPaths,
      excludedAuthors,
      since,
      until,
      labels: requiredLabels,
      milestone,
    });

    if (!evalResult.qualifies) {
      continue;
    }

    const conventional = parseConventionalTitle(pr.title);
    const rawIssueNums = extractAddressedIssues(
      pr.body || "",
      `${owner}/${repo}`,
    );
    const trackingInfo = trackingMap.get(pr.number);
    const addressedIssues = [];

    if (trackingInfo) {
      addressedIssues.push(
        typeof trackingInfo === "object"
          ? trackingInfo
          : { number: trackingInfo, state: "open" },
      );
    }

    for (const num of rawIssueNums) {
      if (!addressedIssues.some((iss) => iss.number === num)) {
        const state = trackingMap.issueStates?.get(num) || null;
        addressedIssues.push({ number: num, state });
      }
    }

    qualifiedPRs.push({
      number: pr.number,
      title: pr.title,
      author: pr.user?.login || pr.author?.login || pr.author || "",
      merged_at: pr.merged_at,
      merge_commit_sha: pr.merge_commit_sha || null,
      labels: Array.isArray(pr.labels)
        ? pr.labels.map((l) => (typeof l === "string" ? l : l.name || ""))
        : [],
      milestone: pr.milestone?.title || null,
      url:
        pr.html_url || `https://github.com/${owner}/${repo}/pull/${pr.number}`,
      state: "merged",
      target_branch: pr.base?.ref || "main",
      conventional_type: conventional.type,
      conventional_scope: conventional.scope,
      addressed_issues: addressedIssues,
      changes_path: extractChangesPath(filePaths),
      changed_product_files: allPaths
        ? filePaths
        : filterProductFiles(filePaths, includePaths),
    });
  }

  // Format output
  let outputData = "";
  if (format === "markdown" || format === "md") {
    outputData = formatMarkdown(qualifiedPRs);
  } else if (format === "tracking") {
    outputData = formatBackportTable(qualifiedPRs, { owner, repo });
  } else if (format === "table") {
    outputData = formatTable(qualifiedPRs);
  } else {
    outputData = formatJSON(qualifiedPRs) + "\n";
  }

  if (values.output) {
    await fs.writeFile(values.output, outputData, "utf8");
  } else {
    process.stdout.write(outputData);
  }
}

async function main() {
  try {
    await runMergedPRs();
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
