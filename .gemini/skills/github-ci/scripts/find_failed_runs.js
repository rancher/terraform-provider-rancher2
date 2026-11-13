#!/usr/bin/env node
import { getCurrentBranch, getOctokit, getRepoDetails } from "./auth.js";

/**
 * Finds failed workflow runs.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {string} [params.branch]
 * @param {number} [params.limit=10]
 * @param {string} [params.status]
 * @returns {Promise<Array>} Failed runs
 */
const LIFECYCLE_STATUSES = new Set([
  "completed",
  "in_progress",
  "queued",
  "requested",
  "waiting",
  "pending",
]);

export async function findFailedRuns({
  octokit,
  owner,
  repo,
  branch,
  limit = 10,
  status,
}) {
  const failedRuns = [];
  const perPage = Math.min(Math.max(limit, 30), 100);
  let page = 1;

  const apiStatus =
    status && LIFECYCLE_STATUSES.has(status) ? status : "completed";

  while (failedRuns.length < limit && page <= 15) {
    const req = {
      owner,
      repo,
      status: apiStatus,
      per_page: perPage,
      page,
    };
    if (branch) {
      req.branch = branch;
    }

    const { data } = await octokit.rest.actions.listWorkflowRunsForRepo(req);
    const runs = data.workflow_runs || [];
    if (runs.length === 0) {
      break;
    }

    for (const run of runs) {
      const isMatch = status
        ? run.conclusion === status || run.status === status
        : run.conclusion === "failure" ||
          run.conclusion === "timed_out" ||
          run.conclusion === "startup_failure";
      if (isMatch) {
        failedRuns.push(run);
        if (failedRuns.length >= limit) {
          break;
        }
      }
    }

    if (runs.length < perPage) {
      break;
    }
    page++;
  }

  return failedRuns;
}

/**
 * Parses command-line arguments for workflow run queries.
 * Supports:
 * - [limit] (e.g. 20)
 * - [-n, limit] or [--limit, limit]
 * - [-n=limit] or [--limit=limit]
 * - [-b, branch] or [--branch, branch]
 * - [-b=branch] or [--branch=branch]
 * - [-s, status] or [--status, status]
 * - [-s=status] or [--status=status]
 * - [branch] (e.g. 'main', 'all', 'current')
 * - [branch, limit]
 * - [branch, -n, limit]
 * @param {Array<string>} [args]
 * @returns {{branchTarget: string, limit: number, status?: string}}
 */
export function parseRunArgs(args = []) {
  let limit = 10;
  let branchTarget = "current";
  let status = undefined;

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if ((arg === "-n" || arg === "--limit") && i + 1 < args.length) {
      const parsed = parseInt(args[++i], 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (arg.startsWith("-n=") || arg.startsWith("--limit=")) {
      const parsed = parseInt(arg.slice(arg.indexOf("=") + 1), 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if ((arg === "-b" || arg === "--branch") && i + 1 < args.length) {
      branchTarget = args[++i];
    } else if (arg.startsWith("-b=") || arg.startsWith("--branch=")) {
      branchTarget = arg.slice(arg.indexOf("=") + 1);
    } else if ((arg === "-s" || arg === "--status") && i + 1 < args.length) {
      status = args[++i];
    } else if (arg.startsWith("-s=") || arg.startsWith("--status=")) {
      status = arg.slice(arg.indexOf("=") + 1);
    } else if (/^\d+$/.test(arg)) {
      const parsed = parseInt(arg, 10);
      if (!isNaN(parsed) && parsed > 0) {
        limit = parsed;
      }
    } else if (arg && !arg.startsWith("-")) {
      branchTarget = arg;
    }
  }

  return { branchTarget, limit, status };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: find_failed_runs.js [branch] [limit]
   or: find_failed_runs.js [-b branch] [-n limit] [-s status]

Arguments:
  branch                       Target branch (e.g. 'main', 'all', 'current') (default: 'current')
  limit                        Maximum number of runs to list (default: 10)

Options:
  -b, --branch <name>          Filter workflow runs by branch
  -n, --limit <num>            Maximum number of runs to list (default: 10)
  -s, --status <status>        Filter by run status or conclusion (e.g. 'completed', 'in_progress', 'failure')
  -h, --help                   Show this help message

Environment Variables:
  GH_REPO / GITHUB_REPOSITORY  Target repository in 'owner/repo' format (overrides git remote detection)
  GITHUB_TOKEN / GH_TOKEN      GitHub authentication token (falls back to 'gh auth token')
`);
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("-h") || args.includes("--help")) {
    printHelp();
    process.exit(0);
  }
  const { branchTarget, limit, status } = parseRunArgs(args);

  let branch;
  if (branchTarget === "all") {
    branch = undefined;
  } else if (branchTarget === "current") {
    branch = await getCurrentBranch();
  } else {
    branch = branchTarget;
  }

  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  try {
    const runs = await findFailedRuns({
      octokit,
      owner,
      repo,
      branch,
      limit,
      status,
    });
    const statusNote = status ? `, Status: ${status}` : "";
    console.log(
      `=== Failed Workflow Runs for ${owner}/${repo} (Branch: ${branch || "all"}${statusNote}) ===\n`,
    );
    if (runs.length === 0) {
      console.log("No failed workflow runs found.");
      return;
    }

    for (const r of runs) {
      console.log(`X [${r.id}] ${r.name} (${r.head_branch})`);
      console.log(
        `   Conclusion: ${r.conclusion} | Event: ${r.event} | Triggered: ${r.created_at}`,
      );
      console.log(`   URL: ${r.html_url}`);
    }
  } catch (err) {
    console.error(`❌ Error finding failed workflow runs: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
