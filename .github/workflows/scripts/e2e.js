import { execSync, execFile } from "child_process";
import { promisify } from "util";

const execFileAsync = promisify(execFile);

export async function defaultRunCommand(file, args, opts = {}) {
  return await execFileAsync(file, args, opts);
}

export default async ({
  github,
  context,
  core,
  process,
  runCommand,
}) => {
  const mode = process.env.SCRIPT_MODE;
  switch (mode) {
  case "filter-pr":
    return await runFilterPr({ github, context, core, process });
  case "e2e-status":
    return await runE2EStatus({ github, context, core, process });
  case "clear-runner":
    return await runClearRunner({ github, context, core, process });
  case "wait-for-e2e":
    return await runWaitForE2E({ github, context, core, process });
  case "report-e2e-status":
    return await runReportE2EStatus({ github, context, core, process });
  case "find-completed-e2e-runs":
    return await runFindCompletedE2eRuns({ github, context, core, process });
  case "execute-cleanup":
    return await runExecuteCleanup({ core, process, runCommand });
  default:
    throw new Error(`Unknown E2E script mode: ${mode}`);
  }
};

/**
 * clear-runner: Cleans up heavy software packages on GitHub runners to free up disk space.
 */
async function runClearRunner({ core }) {
  const pathsToRemove = [
    "/usr/lib/jvm",
    "/usr/share/dotnet",
    "/usr/share/swift",
    "/usr/local/.ghcup",
    "/usr/local/julia*",
    "/usr/local/lib/android",
    "/usr/local/share/chromium",
    "/opt/microsoft",
    "/opt/google",
    "/opt/az",
    "/usr/local/share/powershell",
    "/opt/hostedtoolcache",
  ];

  try {
    const output = execSync(
      `df -h / --total | grep total | awk '{print $4}'`,
    ).toString();
    core.info(`Available disk space before cleanup: ${output}`);
  } catch (error) {
    throw new Error(`Failed running df to see disk space: ${error}`);
  }

  for (const path of pathsToRemove) {
    core.info(`Removing ${path}...`);
    try {
      execSync(`sudo bash -c "rm -rf ${path}"`);
    } catch (error) {
      throw new Error(`Failed to remove ${path}: ${error}`);
    }
  }

  core.info("Pruning Docker...");
  try {
    execSync(`docker system prune -af`);
    execSync(`docker builder prune -af`);
    execSync(`docker image prune -af`);
    execSync(`docker volume prune -af`);
  } catch (error) {
    throw new Error(`Failed pruning Docker: ${error}`);
  }

  core.info("Disk space after cleanup:");
  try {
    const output = execSync(
      `df -h / --total | grep total | awk '{print $4}'`,
    ).toString();
    core.info(`Available disk space after cleanup: ${output}`);
  } catch (error) {
    throw new Error(`Failed running df to see disk space: ${error}`);
  }
}

/**
 * wait-for-e2e: Creates an E2E testing notification comment on a release PR.
 */
async function runWaitForE2E({ github, context, core, process }) {
  try {
    const prNumberRaw = process.env.PR_NUMBER;
    const prNumber = Number.parseInt(prNumberRaw, 10);
    if (!Number.isFinite(prNumber)) {
      core.setFailed(`Invalid PR_NUMBER: ${prNumberRaw}`);
      return;
    }
    await github.rest.issues.createComment({
      issue_number: prNumber,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: `Please make sure e2e tests pass before merging this PR! \n ${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`,
    });
    core.info(`Successfully commented on PR #${prNumber}`);
  } catch (error) {
    core.setFailed(`Failed to create comment on release PR: ${error.message}`);
  }
}

/**
 * report-e2e-status: Reports the passing or failing status of E2E tests on a release PR.
 */
async function runReportE2EStatus({ github, context, core, process }) {
  try {
    const prNumberRaw = process.env.PR_NUMBER;
    const prNumber = Number.parseInt(prNumberRaw, 10);
    if (!Number.isFinite(prNumber)) {
      core.setFailed(`Invalid PR_NUMBER: ${prNumberRaw}`);
      return;
    }
    const status = process.env.TEST_STATUS;
    const msg =
      status === "success"
        ? "End to End Tests Passed!"
        : "End to End Tests Failed!";
    await github.rest.issues.createComment({
      issue_number: prNumber,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: `${msg} \n ${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID}`,
    });
    core.info(`Successfully reported E2E test ${status} on PR #${prNumber}`);
  } catch (error) {
    core.setFailed(`Failed to create comment on release PR: ${error.message}`);
  }
}

/**
 * shouldRunE2E: Evaluates if a PR is an authorized, automated PR targeting a release branch.
 */
function shouldRunE2E(pr, core = { info: () => {}, warning: () => {} }) {
  if (!pr) {
    core.warning("No pull request payload found.");
    return false;
  }

  const baseRepo = pr.base?.repo?.full_name;
  const headRepo = pr.head?.repo?.full_name;
  if (!baseRepo || !headRepo) {
    core.warning(
      "Missing repository information in PR payload. Skipping E2E tests.",
    );
    return false;
  }

  const isFork = Boolean(
    pr.head?.repo?.fork || baseRepo.toLowerCase() !== headRepo.toLowerCase(),
  );
  if (isFork) {
    core.info(
      `PR is from an external fork (${headRepo} -> ${baseRepo}). Skipping E2E tests for security.`,
    );
    return false;
  }

  const rawBaseRef = pr.base?.ref || "";
  const cleanBaseRef = rawBaseRef.replace(/^refs\/heads\//, "");
  if (!/^release\/v\d+$/.test(cleanBaseRef)) {
    core.info(
      `PR target branch "${rawBaseRef}" is not a release branch. Skipping E2E tests.`,
    );
    return false;
  }

  const rawHeadRef = pr.head?.ref || "";
  const cleanHeadRef = rawHeadRef.replace(/^refs\/heads\//, "");
  const userType = pr.user?.type;
  const userLogin = pr.user?.login;
  const normalizedLogin =
    typeof userLogin === "string" ? userLogin.toLowerCase() : "";
  const authorizedBots = [
    "eio-github-actions-power-enhancer[bot]",
    "app/eio-github-actions-power-enhancer",
  ];
  const isBot =
    (userType === "Bot" ||
      normalizedLogin.endsWith("[bot]") ||
      normalizedLogin.startsWith("app/")) &&
    authorizedBots.includes(normalizedLogin);

  const isBackport =
    isBot && /^backport-\d+-(?:release-)?v\d+$/.test(cleanHeadRef);
  const isReleasePlease =
    isBot &&
    /^release-please--branches--release(?:--|\/)v\d+$/.test(cleanHeadRef);

  if (isBackport || isReleasePlease) {
    core.info(
      `PR #${pr.number} qualifies for E2E tests (Type: ${isBackport ? "backport" : "release-please"}, Head: ${cleanHeadRef}, Author: ${userLogin}).`,
    );
    return true;
  }

  core.info(
    `PR #${pr.number} is not an authorized automated PR (Author: ${userLogin}, Type: ${userType}, Head: ${cleanHeadRef}). Skipping E2E tests.`,
  );
  return false;
}

/**
 * resolveTargetPr: Resolves the target pull request object from payload, PR number, or workflow_run commit fallback.
 */
async function resolveTargetPr({ github, context, core, process }) {
  if (context?.payload?.pull_request) {
    return context.payload.pull_request;
  }

  let prNumber = null;
  if (process.env.PR_NUMBER) {
    const parsed = Number.parseInt(process.env.PR_NUMBER, 10);
    if (Number.isFinite(parsed) && parsed > 0) {
      prNumber = parsed;
    }
  }

  const parentRun = context?.payload?.workflow_run;

  if (!prNumber && parentRun?.pull_requests?.length > 0) {
    const parentPrNum = parentRun.pull_requests[0].number;
    if (Number.isFinite(parentPrNum) && parentPrNum > 0) {
      prNumber = parentPrNum;
    }
  }

  if (prNumber) {
    core.info(`Fetching PR #${prNumber} from API...`);
    try {
      const response = await github.rest.pulls.get({
        owner: context.repo.owner,
        repo: context.repo.repo,
        pull_number: prNumber,
      });
      return response.data;
    } catch (error) {
      core.setFailed(`Failed to fetch PR #${prNumber}: ${error.message}`);
      return null;
    }
  }

  // Fallback for listener runs triggered by fork PRs where workflow_run.pull_requests is empty
  const headSha = parentRun?.head_sha || process.env.WORKFLOW_RUN_SHA;
  if (
    headSha &&
    typeof github?.rest?.repos?.listPullRequestsAssociatedWithCommit ===
      "function"
  ) {
    core.info(
      `Attempting fallback PR resolution for commit SHA: ${headSha}...`,
    );
    try {
      let associatedPrs = [];
      if (typeof github.paginate === "function") {
        associatedPrs = await github.paginate(
          github.rest.repos.listPullRequestsAssociatedWithCommit,
          {
            owner: context.repo.owner,
            repo: context.repo.repo,
            commit_sha: headSha,
          },
        );
      } else {
        const res =
          await github.rest.repos.listPullRequestsAssociatedWithCommit({
            owner: context.repo.owner,
            repo: context.repo.repo,
            commit_sha: headSha,
          });
        associatedPrs = res?.data || res || [];
      }

      // Match PR by head branch and repo owner if available, strictly requiring open PRs
      const parentBranch = parentRun?.head_branch
        ? parentRun.head_branch.replace(/^refs\/heads\//, "")
        : null;
      const matched = associatedPrs.find((p) => {
        if (p.state && p.state !== "open") {
          return false;
        }
        const prHeadRef = p.head?.ref
          ? p.head.ref.replace(/^refs\/heads\//, "")
          : null;
        if (parentBranch && prHeadRef !== parentBranch) {
          return false;
        }
        if (
          parentRun?.head_repository?.full_name &&
          p.head?.repo?.full_name?.toLowerCase() !==
            parentRun.head_repository.full_name.toLowerCase()
        ) {
          return false;
        }
        return true;
      });

      if (matched) {
        core.info(`Resolved PR #${matched.number} via commit lookup.`);
        if (matched.number && typeof github?.rest?.pulls?.get === "function") {
          try {
            const fullPr = await github.rest.pulls.get({
              owner: context.repo.owner,
              repo: context.repo.repo,
              pull_number: matched.number,
            });
            if (fullPr?.data) {
              return fullPr.data;
            }
          } catch (err) {
            core.info(
              `Could not fetch full details for PR #${matched.number}, using commit-matched PR object: ${err.message}`,
            );
          }
        }
        return matched;
      }
    } catch (error) {
      core.warning(
        `Error searching PRs associated with commit ${headSha}: ${error.message}`,
      );
    }
  }

  return null;
}

/**
 * filter-pr: Evaluates PR qualification and sets the 'should_run' step output.
 */
async function runFilterPr({ github, context, core, process }) {
  if (context?.eventName === "merge_group") {
    core.info(
      "Triggered by merge_group event. E2E tests are not executed in merge queue (M2). Skipping E2E tests safely.",
    );
    core.setOutput("should_run", "false");
    core.setOutput(
      "sha",
      context?.payload?.merge_group?.head_sha || context?.sha || "",
    );
    core.setOutput("is_fork", "false");
    return true;
  }

  const pr = await resolveTargetPr({ github, context, core, process });
  const parentRun = context?.payload?.workflow_run;

  if (!pr) {
    const isExternalRun =
      parentRun?.head_repository?.fork === true ||
      (parentRun?.head_repository?.full_name &&
        parentRun.head_repository.full_name.toLowerCase() !==
          `${context?.repo?.owner}/${context?.repo?.repo}`.toLowerCase());

    if (isExternalRun) {
      core.info(
        "Triggered by an external fork repository without an identifiable open PR. Skipping E2E tests safely.",
      );
      core.setOutput("should_run", "false");
      core.setOutput(
        "sha",
        parentRun?.head_sha || process.env.WORKFLOW_RUN_SHA || "",
      );
      core.setOutput("is_fork", "true");
      return false;
    }

    core.setFailed(
      "Missing pull_request in context payload and PR_NUMBER env var.",
    );
    core.setOutput("should_run", "false");
    core.setOutput("sha", "");
    core.setOutput("is_fork", "false");
    return false;
  }

  const baseRepo = pr.base?.repo?.full_name;
  const headRepo = pr.head?.repo?.full_name;
  const isFork = Boolean(
    pr.head?.repo?.fork ||
    !pr.head?.repo ||
    (baseRepo && headRepo && baseRepo.toLowerCase() !== headRepo.toLowerCase()),
  );
  core.setOutput("is_fork", isFork ? "true" : "false");

  if (!pr.head?.sha) {
    core.setFailed(
      `Target PR #${pr.number || "unknown"} is missing a verified head commit SHA.`,
    );
    core.setOutput("should_run", "false");
    core.setOutput("sha", "");
    return false;
  }

  const triggerPrHeadSha =
    process.env.PR_HEAD_SHA || parentRun?.pull_requests?.[0]?.head?.sha;
  if (triggerPrHeadSha && triggerPrHeadSha !== pr.head.sha) {
    core.info(
      `Workflow run PR head commit (${triggerPrHeadSha}) is superseded by current PR head (${pr.head.sha}). Skipping stale run.`,
    );
    core.setOutput("should_run", "false");
    core.setOutput("sha", "");
    return false;
  }

  const shouldRun = shouldRunE2E(pr, core);
  core.setOutput("should_run", shouldRun ? "true" : "false");
  core.setOutput("sha", pr.head.sha);

  return shouldRun;
}

/**
 * evaluateE2EStatus: Evaluates the combined statuses of the workflow jobs for Branch Protection.
 */
function evaluateE2EStatus(
  {
    filterResult,
    setupResult = "success",
    shouldRun,
    testResult,
    cleanupResult,
  },
  core = { info: () => {}, setFailed: () => {} },
) {
  if (filterResult !== "success") {
    core.setFailed(
      `Filter job did not complete successfully (Result: ${filterResult}).`,
    );
    return false;
  }

  if (shouldRun !== "true") {
    core.info(
      "E2E tests were safely skipped (PR does not qualify). Status check passing.",
    );
    return true;
  }

  if (setupResult !== "success" && setupResult !== "skipped") {
    core.setFailed(
      `Setup job did not complete successfully (Result: ${setupResult}).`,
    );
    return false;
  }

  if (
    testResult === "success" &&
    (cleanupResult === "success" || cleanupResult === "skipped")
  ) {
    core.info("All E2E tests passed successfully!");
    return true;
  }

  core.setFailed(
    `E2E tests failed or incomplete (Test: ${testResult}, Cleanup: ${cleanupResult}).`,
  );
  return false;
}

/**
 * e2e-status: Evaluates workflow outcome and determines status check pass/fail.
 * Posts the final status directly to the PR's commit SHA to enforce branch protection.
 */
async function runE2EStatus({ github, context, core, process }) {
  const filterResult = process.env.FILTER_RESULT;
  const setupResult = process.env.SETUP_RESULT;
  const shouldRun = process.env.SHOULD_RUN;
  const testResult = process.env.TEST_RESULT;
  const cleanupResult = process.env.CLEANUP_RESULT;
  const sha = process.env.SHA;
  const isFork = process.env.IS_FORK === "true";

  const passed = evaluateE2EStatus(
    {
      filterResult,
      setupResult,
      shouldRun,
      testResult,
      cleanupResult,
    },
    core,
  );

  const skipCommitStatus = String(process.env.SKIP_COMMIT_STATUS || "").trim().toLowerCase() === "true";

  if (skipCommitStatus) {
    core.info("Skipping commit status creation due to SKIP_COMMIT_STATUS flag.");
  } else if (
    sha &&
    github?.rest?.repos?.createCommitStatus
  ) {
    core.info(`Posting commit status to SHA: ${sha}`);
    let description;
    if (!passed) {
      description =
        filterResult !== "success"
          ? "E2E test filter job failed"
          : "E2E tests failed or incomplete";
    } else if (shouldRun !== "true") {
      description = "E2E tests safely skipped (PR does not qualify)";
    } else {
      description = "All E2E tests passed successfully";
    }

    const runId = context.runId || process.env.GITHUB_RUN_ID;
    try {
      await github.rest.repos.createCommitStatus({
        owner: context.repo.owner,
        repo: context.repo.repo,
        sha: sha,
        state: passed ? "success" : "failure",
        context: "E2E Tests / executor",
        description,
        target_url: `https://github.com/${context.repo.owner}/${context.repo.repo}/actions/runs/${runId}`,
      });
    } catch (error) {
      if (isFork) {
        core.info(
          `Unable to post commit status to fork PR commit (${sha}): ${error.message}`,
        );
      } else {
        core.setFailed(
          `Failed to post commit status to SHA ${sha}: ${error.message}`,
        );
        return false;
      }
    }
  }

  return passed;
}

/**
 * find-completed-e2e-runs: Discovers E2E workflow runs completed since the last cleanup run.
 */
export async function runFindCompletedE2eRuns({
  github,
  context,
  core,
  process = globalThis.process,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  // Check if explicit run_id input was provided
  const explicitRunId = (
    process?.env?.EXPLICIT_RUN_ID ||
    context?.payload?.inputs?.run_id ||
    ""
  ).trim();

  if (explicitRunId) {
    if (core && typeof core.info === "function") {
      core.info(`Explicit E2E run ID provided: ${explicitRunId}`);
    }
    const parsedId = parseInt(explicitRunId, 10);
    const runIds = !isNaN(parsedId) ? [parsedId] : [];
    if (core && typeof core.setOutput === "function") {
      core.setOutput("run_ids", JSON.stringify(runIds));
      core.setOutput("count", String(runIds.length));
    }
    return { runIds, count: runIds.length };
  }

  // Determine 'since' timestamp
  let sinceTimestamp = (
    process?.env?.SINCE ||
    context?.payload?.inputs?.since ||
    ""
  ).trim();

  if (
    !sinceTimestamp &&
    typeof github?.rest?.actions?.listWorkflowRuns === "function"
  ) {
    try {
      const { data: cleanupRuns } = await github.rest.actions.listWorkflowRuns({
        owner,
        repo,
        workflow_id: "cleanup-e2e.yml",
        status: "completed",
        per_page: 5,
      });

      const currentRunId =
        context?.runId || parseInt(process?.env?.GITHUB_RUN_ID || "0", 10);
      const prevRun =
        cleanupRuns?.workflow_runs?.find(
          (r) => r.id !== currentRunId && r.conclusion === "success",
        ) ||
        cleanupRuns?.workflow_runs?.find((r) => r.id !== currentRunId);

      if (prevRun?.created_at) {
        sinceTimestamp = prevRun.created_at;
        if (core && typeof core.info === "function") {
          core.info(
            `Last cleanup workflow run #${prevRun.id} completed at ${sinceTimestamp}`,
          );
        }
      }
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(
          `Could not determine previous cleanup workflow run timestamp: ${err.message}`,
        );
      }
    }
  }

  // Fallback to 24 hours ago if no previous run timestamp was found
  if (!sinceTimestamp) {
    const twentyFourHoursAgo = new Date(
      Date.now() - 24 * 60 * 60 * 1000,
    ).toISOString();
    sinceTimestamp = twentyFourHoursAgo;
    if (core && typeof core.info === "function") {
      core.info(
        `No previous cleanup run found; defaulting to runs completed since 24h ago: ${sinceTimestamp}`,
      );
    }
  }

  const sinceDate = new Date(sinceTimestamp);

  // Fetch completed E2E workflow runs using pagination
  let allRuns = [];
  let completedRuns = [];
  try {
    allRuns =
      typeof github?.paginate === "function"
        ? await github.paginate(github.rest.actions.listWorkflowRuns, {
          owner,
          repo,
          workflow_id: "e2e.yml",
          status: "completed",
          per_page: 100,
        })
        : (
          await github.rest.actions.listWorkflowRuns({
            owner,
            repo,
            workflow_id: "e2e.yml",
            status: "completed",
            per_page: 100,
          })
        )?.data?.workflow_runs || [];

    completedRuns = allRuns.filter((r) => {
      if ((r.status || "").toLowerCase() !== "completed") return false;
      const completedAt = r.updated_at || r.created_at;
      return new Date(completedAt) >= sinceDate;
    });
  } catch (err) {
    throw new Error(`Failed to list E2E workflow runs: ${err.message}`);
  }

  const runIds = completedRuns.map((r) => r.id);
  if (core && typeof core.info === "function") {
    core.info(
      `Found ${runIds.length} E2E run(s) completed since ${sinceTimestamp}: ${runIds.join(", ") || "(none)"}`,
    );
  }

  if (core && typeof core.setOutput === "function") {
    core.setOutput("run_ids", JSON.stringify(runIds));
    core.setOutput("count", String(runIds.length));
  }

  return { runIds, count: runIds.length, since: sinceTimestamp };
}

export const DEFAULT_E2E_TEST_NAMES = [
  "TestOneBasic",
  "TestThreeBasic",
  "TestProductionBasic",
];

export function resolveCleanupIdentifiers(identifier, envTests = "") {
  let testNames = DEFAULT_E2E_TEST_NAMES;
  if (envTests) {
    try {
      const parsed = JSON.parse(envTests);
      if (Array.isArray(parsed) && parsed.length > 0) {
        testNames = parsed;
      }
    } catch {}
  }

  const clean = (identifier || "").trim();
  if (!clean) return [];

  // If identifier already has a test suffix, return it directly
  if (testNames.some((t) => clean.endsWith(`-${t}`))) {
    return [clean];
  }

  // Otherwise, return base identifier plus all suffixed test identifiers
  return [clean, ...testNames.map((t) => `${clean}-${t}`)];
}

/**
 * execute-cleanup: Runs cleanup.sh for a specified run ID and all matrix test suffixed identifiers.
 */
export async function runExecuteCleanup({
  core,
  process = globalThis.process,
  runCommand = defaultRunCommand,
}) {
  const baseIdentifier = (process?.env?.IDENTIFIER || "").trim();
  if (!baseIdentifier || !/^[a-zA-Z0-9_-]+$/.test(baseIdentifier)) {
    throw new Error(
      `Invalid or missing cleanup identifier: "${baseIdentifier}". Identifier must contain only alphanumeric characters, dashes, or underscores.`,
    );
  }

  const targetIdentifiers = resolveCleanupIdentifiers(
    baseIdentifier,
    process?.env?.ALL_TESTS_JSON,
  );

  if (core && typeof core.info === "function") {
    core.info(
      `Starting scheduled AWS leftovers cleanup for identifier: ${baseIdentifier} (targets: ${targetIdentifiers.join(", ")})`,
    );
  }

  const runner =
    typeof runCommand === "function" ? runCommand : defaultRunCommand;

  for (const id of targetIdentifiers) {
    if (core && typeof core.info === "function") {
      core.info(`Executing cleanup for target identifier: ${id}`);
    }
    await runner(".github/workflows/scripts/nix-run.sh", ["./cleanup.sh", id], {
      env: {
        ...process.env,
        AWS_MAX_ATTEMPTS: "100",
        AWS_RETRY_MODE: "adaptive",
        IDENTIFIER: id,
      },
    });
  }

  if (core && typeof core.info === "function") {
    core.info(
      `Successfully executed cleanup for ${targetIdentifiers.length} target identifier(s) related to ${baseIdentifier}.`,
    );
  }

  return {
    identifier: baseIdentifier,
    targetIdentifiers,
    success: true,
  };
}

export {
  runClearRunner,
  runWaitForE2E,
  runReportE2EStatus,
  resolveTargetPr,
  runFilterPr,
  shouldRunE2E,
  runE2EStatus,
  evaluateE2EStatus,
};
