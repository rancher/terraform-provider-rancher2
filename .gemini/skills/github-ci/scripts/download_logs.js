#!/usr/bin/env node
import { Buffer } from "node:buffer";
import crypto from "node:crypto";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { logError } from "../../../orchestrator/lib/logger.js";
import { findFailedRuns } from "./find_failed_runs.js";
import { getCurrentBranch, getOctokit, getRepoDetails } from "./auth.js";
import { parseRunId } from "./constants.js";

/**
 * Downloads logs for a workflow run and saves them locally.
 * @param {object} params
 * @param {import('octokit').Octokit} params.octokit
 * @param {string} params.owner
 * @param {string} params.repo
 * @param {number} params.runId
 * @param {string} params.outputDir
 * @param {boolean} [params.failedOnly=false]
 * @returns {Promise<{outputDir: string, jobsCount: number}>}
 */
export async function downloadLogs({
  octokit,
  owner,
  repo,
  runId,
  outputDir,
  failedOnly = false,
}) {
  const validRunId = parseRunId(runId);
  if (!validRunId) {
    throw new Error(`Invalid workflow run ID: ${runId}`);
  }

  outputDir = outputDir || path.join(os.tmpdir(), `ci-logs-${validRunId}`);
  await fs.mkdir(outputDir, { recursive: true });

  // List all jobs for the workflow run across all pages
  const jobs = octokit.paginate
    ? await octokit.paginate(octokit.rest.actions.listJobsForWorkflowRun, {
      owner,
      repo,
      run_id: validRunId,
      per_page: 100,
    })
    : (
      await octokit.rest.actions.listJobsForWorkflowRun({
        owner,
        repo,
        run_id: validRunId,
        per_page: 100,
      })
    ).data.jobs || [];

  const targetJobs = failedOnly
    ? jobs.filter(
      (j) =>
        j.conclusion === "failure" ||
          j.conclusion === "timed_out" ||
          j.conclusion === "startup_failure" ||
          j.conclusion === "cancelled" ||
          j.status === "in_progress",
    )
    : jobs.filter((j) => j.status !== "queued" && j.conclusion !== "skipped");

  let savedCount = 0;

  for (const job of targetJobs) {
    let logText;
    let filePath;
    try {
      const response = await octokit.rest.actions.downloadJobLogsForWorkflowRun(
        {
          owner,
          repo,
          job_id: job.id,
        },
      );

      if (Buffer.isBuffer(response.data)) {
        logText = response.data.toString("utf8");
      } else if (
        response.data instanceof ArrayBuffer ||
        ArrayBuffer.isView(response.data)
      ) {
        logText = Buffer.from(response.data).toString("utf8");
      } else {
        logText =
          typeof response.data === "string"
            ? response.data
            : response.data != null
              ? String(response.data)
              : "";
      }
      const rawJobName = job.name || `job_${job.id}`;
      const safeJobName = rawJobName
        .replace(/[^a-zA-Z0-9._-]/g, "_")
        .slice(0, 100);
      const hash = crypto
        .createHash("sha256")
        .update(rawJobName)
        .digest("hex")
        .slice(0, 8);
      filePath = path.join(
        outputDir,
        `job-${job.id}-${safeJobName}-${hash}.log`,
      );
    } catch (err) {
      if (err.status === 401 || err.status === 403) {
        throw err;
      }
      if (err.status !== 404 && err.status !== 410) {
        logError(`Error downloading logs for job ${job.id}: ${err.message}`);
      }
      // Some jobs might not have logs available yet
      continue;
    }

    await fs.writeFile(filePath, logText, "utf8");
    savedCount++;
  }

  return { outputDir, jobsCount: savedCount };
}

/**
 * Parses command-line arguments for download_logs.js.
 * Supports:
 * - [run_id] [output_dir]
 * - [output_dir]
 * - [run_id]
 * - ['current', output_dir]
 * - [-o, output_dir] or [--output, output_dir] or [--output=dir]
 * - [--failed-only] or [-f]
 * @param {Array<string>} [args]
 * @returns {{runId: number|null, outputDir: string|null, failedOnly: boolean}}
 */
export function parseDownloadArgs(args = []) {
  let runId = null;
  let outputDir = null;
  let failedOnly = false;

  const positional = [];
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === "--failed-only" || arg === "-f") {
      failedOnly = true;
    } else if ((arg === "-o" || arg === "--output") && i + 1 < args.length) {
      outputDir = args[++i];
    } else if (arg.startsWith("--output=") || arg.startsWith("-o=")) {
      outputDir = arg.slice(arg.indexOf("=") + 1);
    } else if (
      (arg === "-r" || arg === "--run" || arg === "--run-id") &&
      i + 1 < args.length
    ) {
      runId = parseRunId(args[++i]) ?? runId;
    } else if (
      arg.startsWith("-r=") ||
      arg.startsWith("--run=") ||
      arg.startsWith("--run-id=")
    ) {
      runId = parseRunId(arg.slice(arg.indexOf("=") + 1)) ?? runId;
    } else if (arg && !arg.startsWith("-")) {
      positional.push(arg);
    }
  }

  if (positional.length > 0) {
    const first = positional[0];
    const parsedFirst = parseRunId(first);
    if (parsedFirst !== null) {
      runId = runId ?? parsedFirst;
      if (positional.length > 1 && !outputDir) {
        outputDir = positional[1];
      }
    } else if (first === "current") {
      runId = null;
      if (positional.length > 1 && !outputDir) {
        outputDir = positional[1];
      }
    } else {
      if (!outputDir) {
        outputDir = first;
      }
      if (positional.length > 1) {
        const parsedSecond = parseRunId(positional[1]);
        if (parsedSecond !== null) {
          runId = runId ?? parsedSecond;
        }
      }
    }
  }

  if (runId && !outputDir) {
    outputDir = path.join(os.tmpdir(), `ci-logs-${runId}`);
  }

  return { runId, outputDir, failedOnly };
}

/**
 * Prints help and usage information.
 */
export function printHelp() {
  console.log(`Usage: download_logs.js [run_id] [output_dir] [--failed-only]
   or: download_logs.js [-o output_dir] [-f]

Arguments:
  run_id                       Workflow run ID (default: latest failed run on current branch)
  output_dir                   Directory path to save downloaded logs (default: /tmp/ci-logs-<run_id>)

Options:
  -o, --output <dir>           Output directory path for saved logs
  -f, --failed-only            Download logs only for failed/timed_out jobs
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
  const {
    runId: parsedRunId,
    outputDir: parsedOutputDir,
    failedOnly,
  } = parseDownloadArgs(args);
  let runId = parsedRunId;
  let outputDir = parsedOutputDir;

  const { owner, repo } = await getRepoDetails();
  const octokit = await getOctokit();

  if (!runId) {
    const branch = await getCurrentBranch();
    console.log(`Finding latest failed run for branch '${branch || "all"}'...`);
    const failedRuns = await findFailedRuns({
      octokit,
      owner,
      repo,
      branch,
      limit: 1,
    });
    if (failedRuns.length === 0) {
      console.log("No failed workflow runs found to download logs for.");
      return;
    }
    runId = failedRuns[0].id;
    if (!outputDir) {
      outputDir = path.join(os.tmpdir(), `ci-logs-${runId}`);
    }
  }

  try {
    const scopeDesc = failedOnly ? "failed job logs" : "job logs";
    console.log(
      `Downloading ${scopeDesc} for run #${runId} to ${outputDir}...`,
    );
    const result = await downloadLogs({
      octokit,
      owner,
      repo,
      runId,
      outputDir,
      failedOnly,
    });
    console.log(
      `✅ Successfully saved ${result.jobsCount} job log(s) to: ${result.outputDir}`,
    );
  } catch (err) {
    logError(`Error downloading logs: ${err.message}`);
    process.exit(1);
  }
}

if (
  process.argv[1] &&
  import.meta.url.endsWith(process.argv[1].replace(/^.*\//, ""))
) {
  main();
}
