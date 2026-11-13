#!/usr/bin/env node
import { getCurrentBranch, getOctokit, getRepoDetails } from "./auth.js";
import { parseRunArgs } from "./find_failed_runs.js";

/**
 * Lists workflow runs from GitHub Actions.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {string} [params.branch]
 * @param {number} [params.limit=10]
 * @param {string} [params.status]
 * @returns {Promise<Array>} List of workflow runs
 */
export async function listRuns({
  octokit,
  owner,
  repo,
  branch,
  limit = 10,
  status,
}) {
  const req = {
    owner,
    repo,
    per_page: Math.min(limit, 100),
  };
  if (branch) {
    req.branch = branch;
  }
  if (status) {
    req.status = status;
  }

  const { data } = await octokit.rest.actions.listWorkflowRunsForRepo(req);
  return data.workflow_runs || [];
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: list_runs.js [branch] [limit]
   or: list_runs.js [-b branch] [-n limit] [-s status]

Arguments:
  branch                       Target branch (e.g. 'main', 'all', 'current') (default: 'current')
  limit                        Maximum number of runs to list (default: 10)

Options:
  -b, --branch <name>          Filter workflow runs by branch
  -n, --limit <num>            Maximum number of runs to list (default: 10)
  -s, --status <status>        Filter by run status (e.g. 'completed', 'in_progress', 'queued')
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
    const runs = await listRuns({
      octokit,
      owner,
      repo,
      branch,
      limit,
      status,
    });
    const statusNote = status ? `, Status: ${status}` : "";
    console.log(
      `=== Recent Workflow Runs for ${owner}/${repo} (Branch: ${branch || "all"}${statusNote}) ===\n`,
    );
    if (runs.length === 0) {
      console.log("No workflow runs found.");
      return;
    }

    for (const r of runs) {
      const state = (r.conclusion || r.status).toUpperCase();
      const symbol =
        state === "SUCCESS" ? "✓" : state === "FAILURE" ? "X" : "•";
      console.log(
        `${symbol} [${r.id}] ${r.name} - ${state} (${r.head_branch})`,
      );
      console.log(
        `   Event: ${r.event} | Triggered: ${r.created_at} | URL: ${r.html_url}`,
      );
    }
  } catch (err) {
    console.error(`❌ Error listing workflow runs: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
