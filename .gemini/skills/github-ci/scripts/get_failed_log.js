#!/usr/bin/env node
import { Buffer } from "node:buffer";
import { getCurrentBranch, getOctokit, getRepoDetails } from "./auth.js";
import { findFailedRuns } from "./find_failed_runs.js";

/**
 * Strips ANSI escape codes from terminal output.
 * @param {string} str
 * @returns {string}
 */
export function stripAnsi(str) {
   
  return str ? str.replace(/\x1b\[[0-9;]*[a-zA-Z]/g, "") : "";
}

/**
 * Retrieves and parses logs from failed jobs in a workflow run.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {number} params.runId
 * @returns {Promise<{run: object, failedJobs: Array<{job: object, errorLines: Array<string>}>}>}
 */
export async function getFailedLog({ octokit, owner, repo, runId }) {
  const { data: run } = await octokit.rest.actions.getWorkflowRun({
    owner,
    repo,
    run_id: runId,
  });

  // List all jobs for the workflow run across all pages
  const jobs = octokit.paginate
    ? await octokit.paginate(octokit.rest.actions.listJobsForWorkflowRun, {
      owner,
      repo,
      run_id: runId,
      per_page: 100,
    })
    : (
      await octokit.rest.actions.listJobsForWorkflowRun({
        owner,
        repo,
        run_id: runId,
      })
    ).data.jobs || [];

  const failedJobs = [];
  const failureRegex =
    /\b(error|failed|failure|fatal|panic|exit code [1-9]|exit status [1-9]|undefined|cannot find|timed out)\b|##\[error\]/i;

  for (const job of jobs) {
    const isFailed =
      job.conclusion === "failure" ||
      job.conclusion === "timed_out" ||
      job.conclusion === "startup_failure" ||
      job.status === "in_progress";
    if (isFailed) {
      let logText = "";
      try {
        const logRes = await octokit.rest.actions.downloadJobLogsForWorkflowRun(
          {
            owner,
            repo,
            job_id: job.id,
          },
        );
        if (Buffer.isBuffer(logRes.data)) {
          logText = logRes.data.toString("utf8");
        } else if (
          logRes.data instanceof ArrayBuffer ||
          ArrayBuffer.isView(logRes.data)
        ) {
          logText = Buffer.from(logRes.data).toString("utf8");
        } else {
          logText =
            typeof logRes.data === "string" ? logRes.data : String(logRes.data);
        }
      } catch {
        // Logs might not be available
      }

      const lines = logText.split("\n");
      const errorLines = lines
        .map((l) => stripAnsi(l).trim())
        .filter((l) => l.length > 0 && failureRegex.test(l))
        .slice(-30); // Last 30 error lines

      const nonEmpties = lines
        .map((l) => stripAnsi(l).trim())
        .filter((l) => l.length > 0);
      const tailLines = nonEmpties.slice(-15);

      failedJobs.push({
        job,
        errorLines,
        tailLines,
      });
    }
  }

  return { run, failedJobs };
}

/**
 * Parses CLI arguments for get_failed_log.js.
 * Supports:
 * - [run_id] (e.g. 12345 or #12345)
 * - [-r run_id] or [--run run_id] or [--run-id run_id] or [flag=val]
 * - [-b branch] or [--branch branch] or [flag=val]
 * - [branch] (e.g. 'main', 'current', 'all')
 * @param {Array<string>} [args]
 * @returns {{runId: number|null, targetBranch?: string}}
 */
export function parseFailedLogArgs(args = []) {
  let runId = null;
  let targetBranch = undefined;

  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if ((a === "-b" || a === "--branch") && i + 1 < args.length) {
      targetBranch = args[++i];
    } else if (a.startsWith("-b=") || a.startsWith("--branch=")) {
      targetBranch = a.slice(a.indexOf("=") + 1);
    } else if (
      (a === "-r" || a === "--run" || a === "--run-id") &&
      i + 1 < args.length
    ) {
      const parsed = parseInt(args[++i].replace(/^#/, ""), 10);
      if (!isNaN(parsed)) {
        runId = parsed;
      }
    } else if (
      a.startsWith("-r=") ||
      a.startsWith("--run=") ||
      a.startsWith("--run-id=")
    ) {
      const parsed = parseInt(
        a.slice(a.indexOf("=") + 1).replace(/^#/, ""),
        10,
      );
      if (!isNaN(parsed)) {
        runId = parsed;
      }
    } else if (/^#?\d+$/.test(a)) {
      runId = parseInt(a.replace(/^#/, ""), 10);
    } else if (a !== "current" && !a.startsWith("-")) {
      targetBranch = a;
    }
  }

  return { runId, targetBranch };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: get_failed_log.js [run_id] [branch]
   or: get_failed_log.js [-r run_id] [-b branch]

Arguments:
  run_id                       Workflow run ID to inspect (default: latest failed run)
  branch                       Target branch to query for failed runs (e.g. 'main', 'current', 'all') (default: 'current')

Options:
  -r, --run, --run-id <id>     Workflow run ID to inspect
  -b, --branch <name>          Branch name to query for failed runs
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
  const { runId: parsedRunId, targetBranch } = parseFailedLogArgs(args);
  let runId = parsedRunId;

  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  if (!runId || isNaN(runId)) {
    let branch = targetBranch;
    if (!branch || branch === "current") {
      branch = await getCurrentBranch();
    } else if (branch === "all") {
      branch = undefined;
    }
    console.log(`Finding latest failed run for branch '${branch || "all"}'...`);
    const failedRuns = await findFailedRuns({
      octokit,
      owner,
      repo,
      branch,
      limit: 1,
    });
    if (failedRuns.length === 0) {
      console.log("No failed workflow runs found.");
      return;
    }
    runId = failedRuns[0].id;
  }

  console.log(`==================================================`);
  console.log(
    `Inspecting failed workflow run #${runId} in ${owner}/${repo}...`,
  );
  console.log(`==================================================\n`);

  try {
    const { run, failedJobs } = await getFailedLog({
      octokit,
      owner,
      repo,
      runId,
    });
    console.log(`Workflow:   ${run.name} (${run.event})`);
    console.log(`Status:     ${run.status} / ${run.conclusion}`);
    console.log(
      `Branch:     ${run.head_branch} (Commit: ${run.head_sha?.slice(0, 7)})`,
    );
    console.log(`URL:        ${run.html_url}`);

    if (failedJobs.length === 0) {
      console.log("\nNo specific failed jobs found or logs are still pending.");
      return;
    }

    console.log(
      `\n=== Failed Job Log Analysis (${failedJobs.length} job(s)) ===`,
    );
    for (const item of failedJobs) {
      console.log(`\n❌ Job: ${item.job.name} (ID: ${item.job.id})`);
      const failedSteps =
        item.job.steps?.filter(
          (s) =>
            s.conclusion === "failure" ||
            s.conclusion === "timed_out" ||
            s.conclusion === "startup_failure",
        ) || [];
      for (const step of failedSteps) {
        console.log(
          `   Failed Step: ${step.name} (Step #${step.number}, conclusion: ${step.conclusion})`,
        );
      }
      if (item.errorLines.length === 0) {
        if (item.tailLines && item.tailLines.length > 0) {
          console.log("   Recent job output lines before exit:");
          for (const line of item.tailLines) {
            console.log(`   > ${line}`);
          }
        } else {
          console.log(
            "   (No explicit error lines matched, check full logs via download_logs.js)",
          );
        }
      } else {
        console.log("   Relevant failure log lines:");
        for (const line of item.errorLines) {
          console.log(`   > ${line}`);
        }
      }
    }
  } catch (err) {
    console.error(`❌ Error inspecting failed run: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
