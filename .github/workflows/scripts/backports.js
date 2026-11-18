import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

/**
 * Executes a git command asynchronously and extracts detailed stderr on failure.
 */
async function runGit(args, options = {}) {
  try {
    return await execFileAsync("git", args, { encoding: "utf8", ...options });
  } catch (error) {
    const rawDetails =
      error.stderr?.trim() || error.stdout?.trim() || error.message;
    const details = rawDetails
      ?.replace(
        /authorization:\s*basic\s+[^\s'"]+/gi,
        "AUTHORIZATION: basic [REDACTED]",
      )
      ?.replace(
        /https:\/\/[^@\s]+@github\.com/gi,
        "https://[REDACTED]@github.com",
      );
    const sanitizedArgs = args.map((arg) =>
      arg.toLowerCase().includes("authorization:") ? "[REDACTED]" : arg,
    );
    const sanitizedCommand = `git ${sanitizedArgs.join(" ")}`;
    throw new Error(`'${sanitizedCommand}' failed: ${details}`);
  }
}

/**
 * Safely parses the TERRAFORM_MAINTAINERS environment variable.
 */
function parseMaintainers(core, envVar) {
  if (!envVar) {
    return [];
  }
  try {
    const parsed = JSON.parse(envVar);
    return Array.isArray(parsed) ? parsed : [];
  } catch (err) {
    if (core && typeof core.warning === "function") {
      core.warning(
        `Could not parse TERRAFORM_MAINTAINERS: ${err.message}. Defaulting to no assignees.`,
      );
    }
    return [];
  }
}

export default async ({
  github,
  context,
  core,
  process = globalThis.process,
  getOctokit,
}) => {
  const mode = process?.env?.SCRIPT_MODE;
  switch (mode) {
  case "wait-for-settle":
    return await runWaitForSettle({ github, context, core, process });
  case "backport-pr":
    return await runBackportPr({
      github,
      context,
      core,
      process,
      getOctokit,
    });
  case "backport-issues":
    return await runBackportIssues({ github, context, core, process });
  case "merge-label":
    return await runMergeLabel({
      github,
      context,
      core,
      process,
      getOctokit,
    });
  case "close-tracking-issue":
    return await runCloseTrackingIssue({
      github,
      context,
      core,
      process,
      getOctokit,
    });
  case "release-done":
    return await runReleaseDone({
      github,
      context,
      core,
      process,
      getOctokit,
    });
  default:
    throw new Error(`Unknown backport script mode: ${mode}`);
  }
};

/**
 * wait-for-settle: Allows GitHub's API time to index merge commits and retrieve the PR list with exponential back-off.
 */
async function runWaitForSettle({
  github,
  context,
  core,
  process = globalThis.process,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  // Handle input from either manual dispatch or push
  const mergeCommitSha =
    process?.env?.MERGE_COMMIT_SHA || context?.payload?.head_commit?.id;

  if (!mergeCommitSha) {
    core.setFailed(
      "No merge commit SHA found in environment or context payload.",
    );
    return;
  }

  const maxAttempts = 5;
  let delayMs = 3000;

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const prs = await github.paginate(
        github.rest.repos.listPullRequestsAssociatedWithCommit,
        {
          owner,
          repo,
          commit_sha: mergeCommitSha,
        },
      );

      const mergedMainPr = prs?.find(
        (p) => p.base?.ref === "main" && p.merged_at,
      );
      if (mergedMainPr) {
        core.info(
          `Found merged PR #${mergedMainPr.number} associated with commit ${mergeCommitSha} targeting main on attempt ${attempt}.`,
        );
        core.setOutput("merge_commit_sha", mergeCommitSha);
        return;
      }

      if (attempt < maxAttempts) {
        core.info(
          `Attempt ${attempt}/${maxAttempts}: Merged PR targeting main for commit ${mergeCommitSha} not found yet. Retrying in ${delayMs / 1000}s...`,
        );
        await new Promise((resolve) => setTimeout(resolve, delayMs));
        delayMs *= 2;
      } else {
        core.setFailed(
          `No merged PR targeting main associated with commit ${mergeCommitSha} found after ${maxAttempts} attempts.`,
        );
        return;
      }
    } catch (error) {
      if (attempt < maxAttempts) {
        core.warning(
          `Attempt ${attempt}/${maxAttempts} failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}. Retrying in ${delayMs / 1000}s...`,
        );
        await new Promise((resolve) => setTimeout(resolve, delayMs));
        delayMs *= 2;
      } else {
        core.setFailed(
          `Failed to retrieve PRs associated with commit ${mergeCommitSha} after ${maxAttempts} attempts: ${error.message}`,
        );
        return;
      }
    }
  }
}

/**
 * backport-pr: Cherry-picks a commit to the appropriate backport branches and creates PRs.
 */
export async function runBackportPr({
  github,
  context,
  core,
  process = globalThis.process,
  getOctokit,
  gitExec = runGit,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";
  const mergeCommitSha = process?.env?.MERGE_COMMIT_SHA;
  const assignees = parseMaintainers(core, process?.env?.TERRAFORM_MAINTAINERS);

  // GITHUB_MERGE_TOKEN permissions: contents: write, metadata: read, pull_requests: write, workflows: write
  // Using an injected client factory (getOctokit) to instantiate a dedicated `prGithub` client is the
  // recommended approach here:
  // 1. Scope isolation: Overriding the action runner's default `github-token` with GITHUB_MERGE_TOKEN breaks
  //    standard issue operations (like GET /search/issues) because GITHUB_MERGE_TOKEN lacks 'issues' permissions.
  // 2. Instance safety: Octokit instance methods do not support re-scoping via `.defaults(...)` (which is a static constructor).
  // 3. Principle of least privilege: `github` retains the runner's default permissions for issues/labels/assignees,
  //    while `prGithub` is scoped specifically to operations requiring elevated PR/auto-merge privileges.
  const mergeToken = process?.env?.GITHUB_MERGE_TOKEN
    ? process.env.GITHUB_MERGE_TOKEN.trim()
    : undefined;
  const octokitFactory =
    typeof getOctokit === "function"
      ? getOctokit
      : typeof github?.getOctokit === "function"
        ? github.getOctokit.bind(github)
        : undefined;
  const prGithub =
    mergeToken && typeof octokitFactory === "function"
      ? octokitFactory(mergeToken)
      : github;
  let response;

  try {
    response = await github.rest.repos.listPullRequestsAssociatedWithCommit({
      owner,
      repo,
      commit_sha: mergeCommitSha,
    });
  } catch (error) {
    throw new Error(
      `Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`,
    );
  }
  const associatedPrs = response.data;
  if (associatedPrs.length === 0) {
    core.info(`No PRs associated with commit ${mergeCommitSha}. Exiting.`);
    return;
  }

  const pr = associatedPrs.find((p) => p.base.ref === "main" && p.merged_at);
  if (!pr) {
    core.info(`No merged PR found for commit ${mergeCommitSha}.`);
    return;
  }
  core.info(`Found associated PR: #${pr.number}`);

  core.info(
    `Searching for 'internal/tracking' issue linked to PR #${pr.number}`,
  );
  try {
    response = await github.request("GET /search/issues", {
      q: `repo:${owner}/${repo} is:issue state:open label:"internal/tracking" ${pr.number}`,
      advanced_search: true,
      headers: {
        "X-GitHub-Api-Version": "2022-11-28",
      },
    });
  } catch (error) {
    throw new Error(
      `Failed to search for internal/tracking issue for PR #${pr.number}: ${error.message}`,
    );
  }
  const searchResults = response.data;
  const cleanBody = pr.body
    ? pr.body
      .replace(/<!--[\s\S]*?-->/g, "")
      .replace(/(~{3,}|`{3,})[\s\S]*?\1/g, "")
    : "";
  const addressedNumbers = [];
  for (const line of cleanBody.split("\n")) {
    const match = line
      .trim()
      .match(/^[-*]\s+Addresses:\s+(`?)#([0-9]+)\1\.?$/i);
    if (match) {
      const num = parseInt(match[2], 10);
      if (!addressedNumbers.includes(num)) {
        addressedNumbers.push(num);
      }
    }
  }

  let trackingIssue = null;
  if (searchResults?.items && searchResults.items.length > 0) {
    if (addressedNumbers.length > 0) {
      trackingIssue = searchResults.items.find((item) =>
        addressedNumbers.includes(item.number),
      );
    }

    if (!trackingIssue) {
      const trackingRegex = new RegExp(
        `<!--\\s*tracking-pr:\\s*#?${pr.number}\\s*-->`,
        "i",
      );
      for (const item of searchResults.items) {
        try {
          const comments = await github.paginate(
            github.rest.issues.listComments,
            {
              owner,
              repo,
              issue_number: item.number,
            },
          );
          if (
            comments.some(
              (c) =>
                c.user?.login === "github-actions[bot]" &&
                c.body &&
                trackingRegex.test(c.body),
            )
          ) {
            trackingIssue = item;
            break;
          }
        } catch (error) {
          core.info(
            `Could not retrieve comments for issue #${item.number}: ${error.message}`,
          );
        }
      }
    }
  }

  // Fallback: If no tracking issue found via search, check if PR body references an open issue directly
  if (!trackingIssue && addressedNumbers.length > 0) {
    for (const fallbackNum of addressedNumbers) {
      core.info(
        `No tracking issue found in search for PR #${pr.number}. Attempting direct lookup for referenced issue #${fallbackNum}...`,
      );
      let directIssue;
      try {
        const directIssueRes = await github.rest.issues.get({
          owner,
          repo,
          issue_number: fallbackNum,
        });
        directIssue = directIssueRes?.data;
      } catch (directErr) {
        const isNotFound =
          directErr?.status === 404 || directErr?.response?.status === 404;
        if (isNotFound) {
          core.info(
            `Referenced issue #${fallbackNum} was not found (404). Continuing search...`,
          );
          continue;
        }
        throw new Error(
          `Direct issue lookup for #${fallbackNum} failed with API error: ${directErr.message}`,
        );
      }

      if (
        directIssue &&
        !directIssue.pull_request &&
        directIssue.state === "open"
      ) {
        const currentLabels = (directIssue.labels || []).map((l) =>
          typeof l === "string" ? l : l?.name,
        );
        const toAdd = [];
        if (!currentLabels.includes("internal/tracking"))
          toAdd.push("internal/tracking");
        if (!currentLabels.includes("internal/user"))
          toAdd.push("internal/user");
        if (
          toAdd.length > 0 &&
          typeof github?.rest?.issues?.addLabels === "function"
        ) {
          try {
            await github.rest.issues.addLabels({
              owner,
              repo,
              issue_number: fallbackNum,
              labels: toAdd,
            });
            core.info(
              `Added labels [${toAdd.join(", ")}] to issue #${fallbackNum}`,
            );
          } catch (labelErr) {
            throw new Error(
              `Failed to add tracking labels to issue #${fallbackNum}: ${labelErr.message}`,
            );
          }
        }
        trackingIssue = directIssue;
        break;
      } else {
        core.info(
          `Referenced item #${fallbackNum} is not an open issue (pull_request=${Boolean(directIssue?.pull_request)}, state=${directIssue?.state}).`,
        );
      }
    }
  }

  if (!trackingIssue) {
    if (!searchResults?.items || searchResults.items.length === 0) {
      core.info(
        `No 'internal/tracking' issue found for PR #${pr.number}. Exiting.`,
      );
    } else {
      core.info(
        `No verified 'internal/tracking' issue found for PR #${pr.number}. Exiting.`,
      );
    }
    return;
  }
  core.info(`Found tracking issue: #${trackingIssue.number}`);

  core.info(`Fetching sub-issues for tracking issue #${trackingIssue.number}`);
  try {
    response = await github.request(
      "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
      {
        owner: owner,
        repo: repo,
        issue_number: trackingIssue.number,
        headers: {
          "X-GitHub-Api-Version": "2022-11-28",
        },
      },
    );
  } catch (error) {
    throw new Error(
      `Failed to fetch sub-issues for tracking issue #${trackingIssue.number}: ${error.message}`,
    );
  }
  let subIssues = response.data;
  core.info(`Sub-issues data: ${JSON.stringify(subIssues)}`);
  if (!Array.isArray(subIssues)) {
    core.warning(
      `Unexpected sub-issues data format: ${JSON.stringify(subIssues)}`,
    );
    return;
  }

  const currentLabels = (trackingIssue.labels || []).map((l) =>
    typeof l === "string" ? l : l?.name,
  );
  let releaseLabels = currentLabels.filter((l) => /^release\/v\d+$/.test(l));
  if (
    releaseLabels.length === 0 &&
    typeof github?.rest?.repos?.listBranches === "function"
  ) {
    try {
      const branches =
        typeof github.paginate === "function"
          ? await github.paginate(github.rest.repos.listBranches, {
            owner,
            repo,
          })
          : (await github.rest.repos.listBranches({ owner, repo }))?.data || [];
      const releaseBranches = branches
        .map((b) => b.name)
        .filter((name) => /^release\/v\d+$/.test(name))
        .sort((a, b) => {
          const aVer = parseInt(a.replace("release/v", ""), 10);
          const bVer = parseInt(b.replace("release/v", ""), 10);
          return bVer - aVer;
        });
      if (!releaseBranches || releaseBranches.length === 0) {
        throw new Error(
          "No release branches matching release/v<major> found in repository.",
        );
      }
      releaseLabels = [releaseBranches[0]];
    } catch (branchErr) {
      core.error(
        `Failed to discover release branches for #${trackingIssue.number}: ${branchErr.message}`,
      );
      throw new Error(
        `Failed to determine target release branch for #${trackingIssue.number}: ${branchErr.message}`,
      );
    }
  }

  // Find any release labels that do not already have a sub-issue
  const existingSubIssueLabels = new Set(
    subIssues
      .flatMap((sub) =>
        (sub.labels || []).map((l) => (typeof l === "string" ? l : l?.name)),
      )
      .filter((l) => /^release\/v\d+$/.test(l)),
  );
  const missingReleaseLabels = releaseLabels.filter(
    (l) => !existingSubIssueLabels.has(l),
  );

  if (missingReleaseLabels.length > 0) {
    core.info(
      `Missing sub-issues for release labels [${missingReleaseLabels.join(", ")}] on tracking issue #${trackingIssue.number}. Attempting on-demand sub-issue generation...`,
    );
    const creationErrors = [];
    for (const releaseLabel of missingReleaseLabels) {
      try {
        await runBackportIssues({
          github,
          context,
          core,
          process: {
            ...process,
            env: {
              ...process?.env,
              RELEASE_LABEL: releaseLabel,
              ISSUE_NUMBER: String(trackingIssue.number),
              PR: String(pr.number),
            },
          },
        });
        core.info(
          `Created sub-issue on demand for #${trackingIssue.number} on branch '${releaseLabel}'`,
        );
      } catch (createErr) {
        core.warning(
          `On-demand sub-issue generation failed for ${releaseLabel}: ${createErr.message}`,
        );
        creationErrors.push(`${releaseLabel}: ${createErr.message}`);
      }
    }

    try {
      const refetched = await github.request(
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
        {
          owner: owner,
          repo: repo,
          issue_number: trackingIssue.number,
          headers: {
            "X-GitHub-Api-Version": "2022-11-28",
          },
        },
      );
      if (Array.isArray(refetched?.data)) {
        subIssues = refetched.data;
        core.info(
          `Successfully refetched sub-issues on demand for #${trackingIssue.number} (found ${subIssues.length})`,
        );
      }
    } catch (refetchErr) {
      throw new Error(
        `Failed to refetch sub-issues for #${trackingIssue.number}: ${refetchErr.message}`,
      );
    }

    if (creationErrors.length > 0) {
      throw new Error(
        `Failed to create on-demand sub-issue(s) for #${trackingIssue.number}:\n- ${creationErrors.join("\n- ")}`,
      );
    }

    if (subIssues.length === 0) {
      throw new Error(
        `Failed to verify created sub-issues on demand for #${trackingIssue.number}: sub-issue list is empty.`,
      );
    }

    const refetchedLabels = new Set(
      subIssues
        .flatMap((sub) =>
          (sub.labels || []).map((l) => (typeof l === "string" ? l : l?.name)),
        )
        .filter((l) => typeof l === "string" && /^release\/v\d+$/.test(l)),
    );
    const unverifiedReleaseLabels = missingReleaseLabels.filter(
      (l) => !refetchedLabels.has(l),
    );
    if (unverifiedReleaseLabels.length > 0) {
      throw new Error(
        `Failed to verify created sub-issues on demand for #${trackingIssue.number}: missing sub-issue(s) for label(s) [${unverifiedReleaseLabels.join(", ")}].`,
      );
    }
  }

  if (subIssues.length === 0) {
    core.info(
      `No sub-issues found for issue #${trackingIssue.number}. Exiting.`,
    );
    return;
  }
  core.info(`Found ${subIssues.length} sub-issues.`);

  const failedBackports = [];
  const createdBackportPrs = [];

  try {
    try {
      await gitExec(["config", "--local", "user.name", "github-actions[bot]"]);
      await gitExec([
        "config",
        "--local",
        "user.email",
        "github-actions[bot]@users.noreply.github.com",
      ]);
      if (mergeToken) {
        try {
          const includeConfig = await gitExec([
            "config",
            "--local",
            "--name-only",
            "--get-regexp",
            "^includeif\\.gitdir:",
          ]);
          const includeKeys = [
            ...new Set(
              includeConfig.stdout
                .trim()
                .split("\n")
                .map((k) => k.trim())
                .filter(Boolean),
            ),
          ];
          for (const key of includeKeys) {
            await gitExec(["config", "--local", "--unset-all", key]);
          }
        } catch {
          // ignore if no includeif.gitdir keys exist
        }

        try {
          const headerConfig = await gitExec([
            "config",
            "--local",
            "--name-only",
            "--get-regexp",
            "http\\..*extraheader",
          ]);
          const headerKeys = [
            ...new Set(
              headerConfig.stdout
                .trim()
                .split("\n")
                .map((k) => k.trim())
                .filter(Boolean),
            ),
          ];
          for (const key of headerKeys) {
            await gitExec(["config", "--local", "--unset-all", key]);
          }
        } catch {
          // ignore if no extraheader keys exist
        }

        core.setSecret(mergeToken);
        const b64Token = Buffer.from(`x-access-token:${mergeToken}`).toString(
          "base64",
        );
        core.setSecret(b64Token);
        await gitExec([
          "config",
          "--local",
          "--replace-all",
          "http.https://github.com/.extraheader",
          `AUTHORIZATION: basic ${b64Token}`,
        ]);
      }
    } catch (error) {
      throw new Error(`Failed to configure git credentials: ${error.message}`);
    }

    for (const subIssue of subIssues) {
      const subIssueNumber = subIssue.number;
      core.info(`Processing sub-issue #${subIssueNumber}...`);

      const releaseLabel = subIssue.labels?.find((label) =>
        label.name?.startsWith("release/v"),
      );

      if (!releaseLabel) {
        core.warning(
          `Sub-issue #${subIssueNumber} has no 'release/v...' label. Skipping.`,
        );
        continue;
      }

      const targetBranch = releaseLabel.name;
      const isValidBranch = /^release\/v\d{1,2}$/.test(targetBranch);

      if (!isValidBranch) {
        const msg = `Target branch label "${targetBranch}" on sub-issue #${subIssueNumber} is invalid. It must start with "release/v" and end with exactly one or two digits.`;
        core.error(msg);
        failedBackports.push({ subIssueNumber, targetBranch, error: msg });
        continue;
      }

      core.info(
        `Processing sub-issue #${subIssueNumber} for target branch: ${targetBranch}`,
      );
      const newBranchName = `backport-${pr.number}-${targetBranch.replace(/\//g, "-")}`;

      try {
        await gitExec(["fetch", "origin", targetBranch]);
        await gitExec([
          "checkout",
          "-B",
          newBranchName,
          `origin/${targetBranch}`,
        ]);
        await gitExec([
          "cherry-pick",
          "--allow-empty",
          "-x",
          mergeCommitSha,
          "-X",
          "theirs",
        ]);
        await gitExec(["push", "origin", newBranchName]);

        core.info(
          `Creating pull request for branch ${newBranchName} targeting ${targetBranch}...`,
        );
        response = await prGithub.rest.pulls.create({
          owner,
          repo,
          title: pr.title,
          head: newBranchName,
          base: targetBranch,
          body: [
            `This pull request cherry-picks the changes from #${pr.number} into ${targetBranch}`,
            `Addresses #${subIssueNumber} for #${trackingIssue.number}`,
            `**WARNING!**: to avoid having to resolve merge conflicts this PR is generated with 'git cherry-pick -X theirs'.`,
            `Please make sure to carefully inspect this PR so that you don't accidentally revert anything!`,
            `Copied from main PR:`,
            `${pr.body}`,
          ].join("\n\n"),
        });

        const newPR = response.data;
        core.info(`Created backport PR data: ${JSON.stringify(newPR)}`);
        const prNumber = newPR.number;
        createdBackportPrs.push({ prNumber, targetBranch, subIssueNumber });

        try {
          await prGithub.graphql(
            `
            mutation EnableAutoMerge($pullRequestId: ID!) {
              enablePullRequestAutoMerge(input: {
                pullRequestId: $pullRequestId,
                mergeMethod: SQUASH
              }) {
                pullRequest { id }
              }
            }
          `,
            { pullRequestId: newPR.node_id },
          );
          core.info(`Auto-merge enabled for PR #${prNumber}`);
        } catch (error) {
          core.warning(
            `Failed to enable auto-merge on PR #${prNumber}: ${error.message}`,
          );
        }

        if (assignees.length > 0) {
          try {
            await github.rest.issues.addAssignees({
              owner,
              repo,
              issue_number: prNumber,
              assignees: assignees,
            });
          } catch (error) {
            throw new Error(
              `Failed to assign PR #${prNumber}: ${error.message}`,
            );
          }
        }

        try {
          await github.rest.issues.addLabels({
            owner,
            repo,
            issue_number: prNumber,
            labels: ["internal/pr-backport", targetBranch],
          });
        } catch (error) {
          throw new Error(
            `Failed to add backport label to PR #${prNumber}: ${error.message}`,
          );
        }
      } catch (error) {
        core.error(
          `Failed to backport to ${targetBranch} for sub-issue #${subIssueNumber}: ${error.message}`,
        );
        failedBackports.push({
          subIssueNumber,
          targetBranch,
          error: error.message,
        });
        try {
          await gitExec(["cherry-pick", "--abort"]);
        } catch {
          // ignore if cherry-pick was not active
        }
        try {
          await gitExec(["checkout", "main", "-f"]);
        } catch {
          // ignore fallback checkout failure
        }
      }
    }
  } finally {
    if (mergeToken) {
      try {
        await gitExec([
          "config",
          "--local",
          "--unset-all",
          "http.https://github.com/.extraheader",
        ]);
      } catch {
        // ignore if unsetting extraheader fails or key was not present
      }
    }
  }

  if (
    createdBackportPrs.length > 0 &&
    trackingIssue?.number &&
    github?.rest?.issues?.createComment
  ) {
    const tableHeader =
      "| Target Branch | Backport PR | CI Checks |\n| --- | --- | --- |";
    const tableRows = createdBackportPrs
      .map((b) => {
        const prLink = `[#${b.prNumber}](https://github.com/${owner}/${repo}/pull/${b.prNumber})`;
        const checksLink = `[View CI Checks](https://github.com/${owner}/${repo}/pull/${b.prNumber}/checks)`;
        return `| \`${b.targetBranch}\` | ${prLink} | ${checksLink} |`;
      })
      .join("\n");

    const commentBody = `### Automated Backport PRs Created\n\n${tableHeader}\n${tableRows}`;
    try {
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number: trackingIssue.number,
        body: commentBody,
      });
      core.info(
        `Successfully posted backport summary comment to tracking issue #${trackingIssue.number}`,
      );
    } catch (commentErr) {
      core.warning(
        `Failed to post backport comment to tracking issue #${trackingIssue.number}: ${commentErr.message}`,
      );
    }
  }

  if (failedBackports.length > 0) {
    const errorSummary = failedBackports
      .map(
        (f) => `Sub-issue #${f.subIssueNumber} (${f.targetBranch}): ${f.error}`,
      )
      .join("\n");
    throw new Error(
      `Failed to complete ${failedBackports.length} backport(s):\n${errorSummary}`,
    );
  }
}

/**
 * backport-issues: Creates sub-issues for tracking backports.
 */
export async function runBackportIssues({
  github,
  context,
  core,
  process = globalThis.process,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  const envLabel = process?.env?.RELEASE_LABEL
    ? process.env.RELEASE_LABEL.trim()
    : "";
  const releaseLabel = envLabel || context?.payload?.label?.name;
  if (!releaseLabel) {
    throw new Error(
      "Release label must be provided via env (RELEASE_LABEL) or label event payload.",
    );
  }
  if (!/^release\/v\d+$/.test(releaseLabel)) {
    throw new Error(
      `Invalid release label: "${releaseLabel}". It must match the pattern release/v<major> (e.g. release/v15).`,
    );
  }

  let parentIssue = context?.payload?.issue;
  const envIssue = process?.env?.ISSUE_NUMBER
    ? String(process.env.ISSUE_NUMBER)
      .replace(/^["']|["']$/g, "")
      .trim()
      .replace(/^#/, "")
    : "";
  const rawIssueNumber =
    envIssue ||
    (parentIssue?.number !== undefined ? String(parentIssue.number) : "");
  if (!rawIssueNumber || !/^[1-9]\d*$/.test(rawIssueNumber)) {
    throw new Error(
      `Invalid issue number: ${process?.env?.ISSUE_NUMBER ?? parentIssue?.number}`,
    );
  }
  const parentIssueNumber = parseInt(rawIssueNumber, 10);

  if (!parentIssue || !parentIssue.title) {
    try {
      const issueRes = await github.rest.issues.get({
        owner,
        repo,
        issue_number: parentIssueNumber,
      });
      parentIssue = issueRes.data;
    } catch (error) {
      throw new Error(
        `Failed to retrieve parent issue #${parentIssueNumber}: ${error.message}`,
      );
    }
  }

  if (parentIssue?.pull_request) {
    throw new Error(
      `Issue #${parentIssueNumber} is a pull request, not an issue.`,
    );
  }

  const parentIssueTitle = parentIssue.title;

  // Check if a sub-issue for this release branch already exists on the parent issue to ensure idempotency
  if (typeof github?.request === "function") {
    try {
      const existingSubIssuesRes = await github.request(
        "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
        {
          owner,
          repo,
          issue_number: parentIssueNumber,
          headers: {
            "X-GitHub-Api-Version": "2022-11-28",
          },
        },
      );
      const existingSubIssues = existingSubIssuesRes?.data;
      if (Array.isArray(existingSubIssues)) {
        const alreadyExists = existingSubIssues.some(
          (si) =>
            si?.labels?.some(
              (l) => (typeof l === "string" ? l : l?.name) === releaseLabel,
            ) ||
            (si?.title && si.title.startsWith(`[${releaseLabel}]`)),
        );
        if (alreadyExists) {
          if (core && typeof core.info === "function") {
            core.info(
              `Sub-issue for release label '${releaseLabel}' already exists on tracking issue #${parentIssueNumber}. Skipping creation.`,
            );
          }
          return;
        }
      }
    } catch (checkError) {
      throw new Error(
        `Failed to check existing sub-issues for #${parentIssueNumber}: ${checkError.message}`,
      );
    }
  }

  const isExcludedMaintainer = (m) => {
    if (typeof m !== "string") return false;
    const lower = m.toLowerCase().trim();
    return (
      lower === "matt-trachier" ||
      lower === "matt.trachier" ||
      lower === "mtrachier"
    );
  };
  const assignees = parseMaintainers(
    core,
    process?.env?.TERRAFORM_MAINTAINERS,
  ).filter((a) => !isExcludedMaintainer(a));
  const rawPr = process?.env?.PR
    ? String(process.env.PR)
      .replace(/^["']|["']$/g, "")
      .trim()
      .replace(/^#/, "")
    : "";
  if (!rawPr || !/^[1-9]\d*$/.test(rawPr)) {
    throw new Error(`Invalid PR number: ${process?.env?.PR}`);
  }
  const extractedPrNumber = parseInt(rawPr, 10);
  let response;

  try {
    response = await github.rest.issues.get({
      owner: owner,
      repo: repo,
      issue_number: extractedPrNumber,
    });
  } catch (error) {
    throw new Error(
      `Failed to retrieve PR #${extractedPrNumber}: ${error.message}`,
    );
  }
  const pr = response.data;
  const prNumber = pr.number;

  try {
    const issueParams = {
      owner: owner,
      repo: repo,
      title: `[${releaseLabel}] ${parentIssueTitle}`,
      body: [
        `Backport #${prNumber} to ${releaseLabel} for #${parentIssueNumber}`,
        `Please add this issue to the proper milestone.`,
        `Copied from PR:`,
        `${pr.body}`,
      ].join("\n\n"),
      labels: [releaseLabel, "internal/backport"],
    };
    if (assignees.length > 0) {
      issueParams.assignees = assignees;
    }
    response = await github.rest.issues.create(issueParams);
  } catch (error) {
    throw new Error(`Failed to create backport issue: ${error.message}`);
  }
  const newIssue = response.data;
  if (core && typeof core.info === "function") {
    core.info(`New backport issue data: ${JSON.stringify(newIssue)}`);
  }
  const subIssueId = newIssue.id;

  try {
    await github.request(
      "POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
      {
        owner: owner,
        repo: repo,
        issue_number: parentIssueNumber,
        sub_issue_id: subIssueId,
        headers: {
          "X-GitHub-Api-Version": "2022-11-28",
        },
      },
    );
  } catch (error) {
    throw new Error(
      `Failed to link backport issue to tracking issue: ${error.message}`,
    );
  }
}

export const DEFAULT_PROJECT_NUMBER = 73;
export const DEFAULT_PROJECT_OWNER = "rancher";

/**
 * Updates an issue's single-select Status field to targetStatus (e.g. "Ready for QA") in any ProjectV2 it belongs to.
 */
export async function updateIssueProjectStatus({
  github,
  owner,
  repo,
  issueNumber,
  targetStatus = "Ready for QA",
  projectNumber = DEFAULT_PROJECT_NUMBER,
  core,
}) {
  if (!github || typeof github.graphql !== "function") {
    if (core && typeof core.info === "function") {
      core.info(
        `GraphQL client unavailable; skipping project status update for #${issueNumber}.`,
      );
    }
    return false;
  }

  try {
    let issueId = null;
    const projectItems = [];
    let hasNextPage = true;
    let cursor = null;

    while (hasNextPage) {
      const query = `
        query getIssueProjectItems($owner: String!, $repo: String!, $issueNumber: Int!, $cursor: String) {
          repository(owner: $owner, name: $repo) {
            issue(number: $issueNumber) {
              id
              projectItems(first: 50, after: $cursor) {
                pageInfo {
                  hasNextPage
                  endCursor
                }
                nodes {
                  id
                  project {
                    id
                    number
                    title
                    fields(first: 50) {
                      nodes {
                        ... on ProjectV2SingleSelectField {
                          id
                          name
                          options {
                            id
                            name
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      `;

      const data = await github.graphql(query, {
        owner,
        repo,
        issueNumber,
        cursor,
      });

      if (!issueId) {
        issueId = data?.repository?.issue?.id;
      }

      const connection = data?.repository?.issue?.projectItems;
      const nodes = connection?.nodes?.filter(Boolean) || [];
      projectItems.push(...nodes);

      hasNextPage = Boolean(connection?.pageInfo?.hasNextPage);
      cursor = connection?.pageInfo?.endCursor || null;
    }

    let updatedCount = 0;

    if (projectItems.length > 0) {
      for (const item of projectItems) {
        const project = item.project;
        if (!project) continue;

        const fields = project.fields?.nodes?.filter(Boolean) || [];
        const statusField = fields.find(
          (f) =>
            typeof f?.name === "string" && f.name.toLowerCase() === "status",
        );

        if (!statusField) {
          if (core && typeof core.info === "function") {
            core.info(
              `Project "${project.title}" has no single-select 'Status' field.`,
            );
          }
          continue;
        }

        const targetOption = statusField.options?.find(
          (o) =>
            typeof o?.name === "string" &&
            o.name.trim().toLowerCase() === targetStatus.trim().toLowerCase(),
        );

        if (!targetOption) {
          if (core && typeof core.info === "function") {
            core.info(
              `Project "${project.title}" Status field has no option matching "${targetStatus}".`,
            );
          }
          continue;
        }

        const mutation = `
          mutation updateItemStatus($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
            updateProjectV2ItemFieldValue(
              input: {
                projectId: $projectId
                itemId: $itemId
                fieldId: $fieldId
                value: {
                  singleSelectOptionId: $optionId
                }
              }
            ) {
              projectV2Item {
                id
              }
            }
          }
        `;

        await github.graphql(mutation, {
          projectId: project.id,
          itemId: item.id,
          fieldId: statusField.id,
          optionId: targetOption.id,
        });

        if (core && typeof core.info === "function") {
          core.info(
            `Updated status of issue #${issueNumber} in project "${project.title}" to "${targetOption.name}".`,
          );
        }
        updatedCount++;
      }
    } else if (issueId && projectNumber) {
      // If issue is not yet associated with any project items, attempt to add it directly to Project 73
      try {
        const projectOwner = owner || DEFAULT_PROJECT_OWNER;
        const projectQuery = `
          query getProjectByNumber($projectOwner: String!, $projectNumber: Int!) {
            organization(login: $projectOwner) {
              projectV2(number: $projectNumber) {
                id
                number
                title
                fields(first: 50) {
                  nodes {
                    ... on ProjectV2SingleSelectField {
                      id
                      name
                      options {
                        id
                        name
                      }
                    }
                  }
                }
              }
            }
          }
        `;

        const projectData = await github.graphql(projectQuery, {
          projectOwner,
          projectNumber,
        });
        const project = projectData?.organization?.projectV2;

        if (project) {
          const fields = project.fields?.nodes?.filter(Boolean) || [];
          const statusField = fields.find(
            (f) =>
              typeof f?.name === "string" && f.name.toLowerCase() === "status",
          );

          const targetOption = statusField?.options?.find(
            (o) =>
              typeof o?.name === "string" &&
              o.name.trim().toLowerCase() === targetStatus.trim().toLowerCase(),
          );

          if (statusField && targetOption) {
            const addMutation = `
              mutation addItemToProject($projectId: ID!, $contentId: ID!) {
                addProjectV2ItemById(input: { projectId: $projectId, contentId: $contentId }) {
                  item {
                    id
                  }
                }
              }
            `;
            const addResult = await github.graphql(addMutation, {
              projectId: project.id,
              contentId: issueId,
            });
            const newItemId = addResult?.addProjectV2ItemById?.item?.id;

            if (newItemId) {
              const updateMutation = `
                mutation updateItemStatus($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
                  updateProjectV2ItemFieldValue(
                    input: {
                      projectId: $projectId
                      itemId: $itemId
                      fieldId: $fieldId
                      value: {
                        singleSelectOptionId: $optionId
                      }
                    }
                  ) {
                    projectV2Item {
                      id
                    }
                  }
                }
              `;
              await github.graphql(updateMutation, {
                projectId: project.id,
                itemId: newItemId,
                fieldId: statusField.id,
                optionId: targetOption.id,
              });

              if (core && typeof core.info === "function") {
                core.info(
                  `Added issue #${issueNumber} to project "${project.title}" (#${project.number}) and set status to "${targetOption.name}".`,
                );
              }
              updatedCount++;
            }
          }
        }
      } catch (addErr) {
        if (core && typeof core.info === "function") {
          core.info(
            `Could not add issue #${issueNumber} to project #${projectNumber}: ${addErr.message}`,
          );
        }
      }
    }

    if (updatedCount === 0 && projectItems.length === 0) {
      if (core && typeof core.info === "function") {
        core.info(
          `Issue #${issueNumber} is not associated with any GitHub Project items.`,
        );
      }
      return false;
    }

    return updatedCount > 0;
  } catch (err) {
    if (core && typeof core.warning === "function") {
      core.warning(
        `Failed to update project status for issue #${issueNumber}: ${err.message}`,
      );
    }
    return false;
  }
}

/**
 * merge-label: Adds internal/merged label to issues referenced in a merged PR body.
 */
export async function runMergeLabel({
  github,
  context,
  core,
  process = globalThis.process,
  getOctokit,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";
  let pr = context?.payload?.pull_request;

  const issuesToken = process?.env?.GITHUB_ISSUES_TOKEN
    ? process.env.GITHUB_ISSUES_TOKEN.trim()
    : undefined;
  const octokitFactory =
    typeof getOctokit === "function"
      ? getOctokit
      : typeof github?.getOctokit === "function"
        ? github.getOctokit.bind(github)
        : undefined;
  const issuesGithub =
    issuesToken && typeof octokitFactory === "function"
      ? octokitFactory(issuesToken)
      : github;

  const normalizeBranch = (b) =>
    (b || "")
      .toString()
      .trim()
      .replace(/^refs\/heads\//, "")
      .replace(/^origin\//, "")
      .replace(/\/+$/, "");

  if (!pr) {
    const sha = (
      process?.env?.SHA ||
      context?.payload?.workflow_run?.head_sha ||
      ""
    ).trim();
    if (
      sha &&
      (typeof github?.paginate === "function" ||
        typeof github?.rest?.repos?.listPullRequestsAssociatedWithCommit ===
          "function")
    ) {
      try {
        const pulls =
          typeof github?.paginate === "function"
            ? await github.paginate(
              github.rest?.repos?.listPullRequestsAssociatedWithCommit,
              {
                owner,
                repo,
                commit_sha: sha,
              },
            )
            : (
              await github.rest?.repos?.listPullRequestsAssociatedWithCommit({
                owner,
                repo,
                commit_sha: sha,
              })
            )?.data || [];
        const rawBranch = (
          process?.env?.BRANCH ||
          context?.payload?.workflow_run?.head_branch ||
          ""
        ).trim();
        const branch = normalizeBranch(rawBranch);
        pr =
          pulls.find(
            (p) =>
              p?.state === "closed" &&
              Boolean(p?.merged_at) &&
              (!branch || normalizeBranch(p?.base?.ref) === branch),
          ) ||
          pulls.find((p) => p?.state === "closed" && Boolean(p?.merged_at));
      } catch (err) {
        if (core && typeof core.warning === "function") {
          core.warning(
            `Could not fetch PR associated with SHA ${sha}: ${err.message}`,
          );
        }
      }
    }
  }

  if (!pr && Array.isArray(context?.payload?.workflow_run?.pull_requests)) {
    const rawBranch = (
      process?.env?.BRANCH ||
      context?.payload?.workflow_run?.head_branch ||
      ""
    ).trim();
    const branch = normalizeBranch(rawBranch);
    const candidates = context.payload.workflow_run.pull_requests.filter(
      (p) =>
        p &&
        (!branch || !p.base?.ref || normalizeBranch(p.base?.ref) === branch),
    );
    const pool =
      candidates.length > 0
        ? candidates
        : context.payload.workflow_run.pull_requests;
    for (const candidate of pool) {
      if (candidate?.number) {
        if (
          candidate.body &&
          candidate.state === "closed" &&
          Boolean(candidate.merged_at)
        ) {
          pr = candidate;
          break;
        }
        if (github?.rest?.pulls?.get) {
          try {
            const { data: fullPr } = await github.rest.pulls.get({
              owner,
              repo,
              pull_number: candidate.number,
            });
            if (fullPr?.state === "closed" && Boolean(fullPr?.merged_at)) {
              pr = fullPr;
              break;
            }
          } catch (err) {
            if (core && typeof core.warning === "function") {
              core.warning(
                `Could not fetch PR #${candidate.number} details: ${err.message}`,
              );
            }
          }
        }
      }
    }
  }

  if (!pr) {
    if (core && typeof core.info === "function") {
      core.info("No pull request found for merge-label; skipping.");
    }
    return;
  }

  const rawBody = pr.body ?? "";
  const cleanBody = rawBody
    .replace(/<!--[\s\S]*?-->/g, "")
    .replace(/(~{3,}|`{3,})[\s\S]*?\1/g, "");

  const issueRegex = /(?:^|[^\w])#(\d+)\b/g;
  const hashMatches = cleanBody.matchAll(issueRegex);
  const hashNumbers = Array.from(hashMatches, (m) => parseInt(m[1], 10));

  const escapeRegex = (s) => (s || "").replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const issueUrlRegex = new RegExp(
    `https?:\\/\\/github\\.com\\/${escapeRegex(owner)}\\/${escapeRegex(repo)}\\/issues\\/(\\d+)\\b`,
    "gi",
  );
  const urlMatches = cleanBody.matchAll(issueUrlRegex);
  const urlNumbers = Array.from(urlMatches, (m) => parseInt(m[1], 10));

  const issueNumbers = [...new Set([...hashNumbers, ...urlNumbers])];

  if (core && typeof core.info === "function") {
    core.info(`Found issue numbers in PR body: ${issueNumbers}`);
  }

  for (const issueNumber of issueNumbers) {
    try {
      const { data: issueData } = await github.rest.issues.get({
        owner,
        repo,
        issue_number: issueNumber,
      });

      const hasBackport = issueData.labels?.some(
        (l) => (typeof l === "string" ? l : l?.name) === "internal/backport",
      );
      const alreadyMerged = issueData.labels?.some(
        (l) => (typeof l === "string" ? l : l?.name) === "internal/merged",
      );

      if (!issueData.pull_request && hasBackport && !alreadyMerged) {
        if (core && typeof core.info === "function") {
          core.info(`Adding 'internal/merged' label to issue #${issueNumber}`);
        }
        await github.rest.issues.addLabels({
          owner,
          repo,
          issue_number: issueNumber,
          labels: ["internal/merged"],
        });
        await updateIssueProjectStatus({
          github: issuesGithub || github,
          owner,
          repo,
          issueNumber,
          targetStatus: "Ready for QA",
          core,
        });
      } else if (!issueData.pull_request && hasBackport && alreadyMerged) {
        if (core && typeof core.info === "function") {
          core.info(
            `Issue #${issueNumber} already has 'internal/merged' label; skipping.`,
          );
        }
        await updateIssueProjectStatus({
          github: issuesGithub || github,
          owner,
          repo,
          issueNumber,
          targetStatus: "Ready for QA",
          core,
        });
      } else if (!issueData.pull_request && !hasBackport) {
        if (core && typeof core.info === "function") {
          core.info(
            `Issue #${issueNumber} does not have 'internal/backport' label; skipping.`,
          );
        }
      } else if (issueData.pull_request) {
        if (core && typeof core.info === "function") {
          core.info(
            `Issue #${issueNumber} is a pull request, not an issue; skipping.`,
          );
        }
      }
    } catch (error) {
      if (core && typeof core.setFailed === "function") {
        core.setFailed(
          `Could not process issue #${issueNumber}: ${error.message}`,
        );
      }
    }
  }

  const baseBranch = pr.base?.ref;
  const trackingMatch = cleanBody.match(/Addresses\s+#?\d+\s+for\s+#?(\d+)/i);
  const trackingIssueNumber = trackingMatch
    ? parseInt(trackingMatch[1], 10)
    : null;

  if (
    trackingIssueNumber &&
    baseBranch &&
    github?.rest?.issues?.createComment
  ) {
    try {
      const pulls =
        typeof github?.paginate === "function"
          ? await github.paginate(github.rest.pulls.list, {
            owner,
            repo,
            state: "open",
            base: baseBranch,
          })
          : (
            await github.rest.pulls.list({
              owner,
              repo,
              state: "open",
              base: baseBranch,
            })
          )?.data || [];

      const isReleasePlease = (p) => {
        const headRef = (p?.head?.ref || "").replace(/^refs\/heads\//, "");
        return (
          headRef.startsWith("release-please--branches--") ||
          (Array.isArray(p?.labels) &&
            p.labels.some((l) =>
              (typeof l === "string" ? l : l?.name)?.startsWith("autorelease"),
            )) ||
          (typeof p?.title === "string" &&
            /^chore(?:\([^)]*\))?:\s*release/i.test(p.title))
        );
      };

      const releasePleasePr = pulls.find(isReleasePlease);
      let commentBody = `### Backport Merged: \`${baseBranch}\`\n\nBackport PR #${pr.number} has merged into \`${baseBranch}\`.`;
      if (releasePleasePr) {
        const rpPrLink = `[#${releasePleasePr.number}](https://github.com/${owner}/${repo}/pull/${releasePleasePr.number})`;
        const rpChecksLink = `[View CI Checks](https://github.com/${owner}/${repo}/pull/${releasePleasePr.number}/checks)`;
        commentBody += `\n\nRelease staging is progressing in Release-Please PR ${rpPrLink} (${rpChecksLink}).`;
      }

      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number: trackingIssueNumber,
        body: commentBody,
      });
      if (core && typeof core.info === "function") {
        core.info(
          `Posted release progression comment to tracking issue #${trackingIssueNumber}`,
        );
      }
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(
          `Failed to post release comment to tracking issue #${trackingIssueNumber}: ${err.message}`,
        );
      }
    }
  }
}

/**
 * close-tracking-issue: Automatically closes a tracking issue when all its sub-issues are closed,
 * ensuring the tracking issue is assigned to a milestone before closure.
 */
export async function runCloseTrackingIssue({
  github,
  context,
  core,
  process = globalThis.process,
  getOctokit,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  const rawIssueNumber =
    process?.env?.ISSUE_NUMBER ||
    context?.payload?.inputs?.issue_number ||
    context?.payload?.issue?.number;

  if (!rawIssueNumber) {
    if (core && typeof core.info === "function") {
      core.info(
        "No issue number provided in payload or environment; skipping.",
      );
    }
    return null;
  }

  const issueNumber = parseInt(String(rawIssueNumber).trim(), 10);
  if (isNaN(issueNumber)) {
    throw new Error(`Invalid issue number: ${rawIssueNumber}`);
  }

  const issuesToken = process?.env?.GITHUB_ISSUES_TOKEN
    ? process.env.GITHUB_ISSUES_TOKEN.trim()
    : undefined;
  const octokitFactory =
    typeof getOctokit === "function"
      ? getOctokit
      : typeof github?.getOctokit === "function"
        ? github.getOctokit.bind(github)
        : undefined;
  const issuesGithub =
    issuesToken && typeof octokitFactory === "function"
      ? octokitFactory(issuesToken)
      : github;

  // Resolve the tracking issue: either issueNumber itself or its parent
  let trackingIssueNumber = null;
  let trackingIssueData = null;

  // Try GraphQL query to find parent or issue details
  if (typeof issuesGithub?.graphql === "function") {
    try {
      const parentQuery = `
        query getIssueParent($owner: String!, $repo: String!, $issueNumber: Int!) {
          repository(owner: $owner, name: $repo) {
            issue(number: $issueNumber) {
              id
              number
              state
              milestone {
                id
                number
                title
              }
              labels(first: 20) {
                nodes {
                  name
                }
              }
              parent {
                id
                number
                title
                state
                milestone {
                  id
                  number
                  title
                }
                labels(first: 20) {
                  nodes {
                    name
                  }
                }
              }
            }
          }
        }
      `;
      const gqlResult = await issuesGithub.graphql(parentQuery, {
        owner,
        repo,
        issueNumber,
      });
      const issueNode = gqlResult?.repository?.issue;
      if (issueNode?.parent) {
        trackingIssueNumber = issueNode.parent.number;
        trackingIssueData = {
          number: issueNode.parent.number,
          state: issueNode.parent.state,
          milestone: issueNode.parent.milestone,
          labels: issueNode.parent.labels?.nodes || [],
        };
      } else if (
        issueNode?.labels?.nodes?.some((l) => l.name === "internal/tracking")
      ) {
        trackingIssueNumber = issueNode.number;
        trackingIssueData = {
          number: issueNode.number,
          state: issueNode.state,
          milestone: issueNode.milestone,
          labels: issueNode.labels?.nodes || [],
        };
      }
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(
          `GraphQL issue lookup for #${issueNumber} failed: ${err.message}`,
        );
      }
    }
  }

  // Fallback to REST if tracking issue wasn't resolved via GraphQL
  if (!trackingIssueNumber) {
    try {
      const { data: issueRest } = await issuesGithub.rest.issues.get({
        owner,
        repo,
        issue_number: issueNumber,
      });
      if (
        Array.isArray(issueRest.labels) &&
        issueRest.labels.some((l) => (l?.name || l) === "internal/tracking")
      ) {
        trackingIssueNumber = issueRest.number;
        trackingIssueData = issueRest;
      }
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(
          `REST issue lookup for #${issueNumber} failed: ${err.message}`,
        );
      }
    }
  }

  if (!trackingIssueNumber) {
    if (core && typeof core.info === "function") {
      core.info(
        `Issue #${issueNumber} is neither an internal/tracking issue nor a sub-issue of a tracking issue; skipping.`,
      );
    }
    return null;
  }

  // Ensure trackingIssueData is populated
  if (!trackingIssueData) {
    const { data: fetchedParent } = await issuesGithub.rest.issues.get({
      owner,
      repo,
      issue_number: trackingIssueNumber,
    });
    trackingIssueData = fetchedParent;
  }

  const hasTrackingLabel = Array.isArray(trackingIssueData.labels)
    ? trackingIssueData.labels.some(
      (l) => (l?.name || l) === "internal/tracking",
    )
    : false;
  if (!hasTrackingLabel) {
    if (core && typeof core.info === "function") {
      core.info(
        `Issue #${trackingIssueNumber} does not have 'internal/tracking' label; skipping.`,
      );
    }
    return null;
  }

  const isClosed =
    (trackingIssueData.state || "").toString().toLowerCase() === "closed";
  if (isClosed) {
    if (core && typeof core.info === "function") {
      core.info(
        `Tracking issue #${trackingIssueNumber} is already closed; skipping.`,
      );
    }
    return null;
  }

  // Fetch all sub-issues of the tracking issue
  let subIssues = [];
  try {
    subIssues =
      typeof issuesGithub?.paginate === "function"
        ? await issuesGithub.paginate(
          "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
          {
            owner,
            repo,
            issue_number: trackingIssueNumber,
            per_page: 100,
            headers: {
              "X-GitHub-Api-Version": "2022-11-28",
            },
          },
        )
        : (
          await issuesGithub.request(
            "GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues",
            {
              owner,
              repo,
              issue_number: trackingIssueNumber,
              per_page: 100,
              headers: {
                "X-GitHub-Api-Version": "2022-11-28",
              },
            },
          )
        )?.data || [];
  } catch (err) {
    throw new Error(
      `Failed to fetch sub-issues for tracking issue #${trackingIssueNumber}: ${err.message}`,
    );
  }

  if (subIssues.length === 0) {
    if (core && typeof core.info === "function") {
      core.info(
        `Tracking issue #${trackingIssueNumber} has no sub-issues; skipping closure.`,
      );
    }
    return null;
  }

  const openSubIssues = subIssues.filter(
    (sub) => (sub?.state || "").toString().toLowerCase() !== "closed",
  );

  if (openSubIssues.length > 0) {
    const openList = openSubIssues.map((s) => `#${s.number}`).join(", ");
    if (core && typeof core.info === "function") {
      core.info(
        `Tracking issue #${trackingIssueNumber} still has open sub-issues (${openList}); skipping closure.`,
      );
    }
    return {
      closed: false,
      trackingIssueNumber,
      reason: "open_sub_issues",
    };
  }

  // All sub-issues are closed. Validate milestone requirement before closing tracking issue.
  let milestoneNumber = trackingIssueData.milestone?.number;

  // If parent tracking issue lacks a milestone, attempt to adopt from sub-issues
  if (!milestoneNumber) {
    const subWithMilestone = subIssues.find(
      (sub) => sub?.milestone && sub.milestone.number,
    );
    if (subWithMilestone) {
      milestoneNumber = subWithMilestone.milestone.number;
      try {
        await issuesGithub.rest.issues.update({
          owner,
          repo,
          issue_number: trackingIssueNumber,
          milestone: milestoneNumber,
        });
        if (core && typeof core.info === "function") {
          core.info(
            `Assigned milestone "${subWithMilestone.milestone.title}" (#${milestoneNumber}) to tracking issue #${trackingIssueNumber}.`,
          );
        }
      } catch (err) {
        milestoneNumber = null;
        if (core && typeof core.warning === "function") {
          core.warning(
            `Failed to assign milestone #${subWithMilestone.milestone.number} to tracking issue #${trackingIssueNumber}: ${err.message}`,
          );
        }
      }
    }
  }

  if (!milestoneNumber) {
    if (core && typeof core.warning === "function") {
      core.warning(
        `All sub-issues for tracking issue #${trackingIssueNumber} are closed, but the tracking issue has no milestone assigned. Skipping closure to prevent tracking issue from being lost.`,
      );
    }

    // Check if warning comment was already posted to avoid comment spam
    try {
      const existingComments = await issuesGithub.rest.issues.listComments({
        owner,
        repo,
        issue_number: trackingIssueNumber,
        per_page: 20,
      });
      const alreadyWarned = Array.isArray(existingComments?.data)
        ? existingComments.data.some((c) =>
          (c?.body || "").includes(
            "cannot be automatically closed because it has not been assigned to a milestone",
          ),
        )
        : false;

      if (!alreadyWarned) {
        await issuesGithub.rest.issues.createComment({
          owner,
          repo,
          issue_number: trackingIssueNumber,
          body: `All sub-issues for this tracking issue have been resolved and closed. However, this tracking issue cannot be automatically closed because it has not been assigned to a milestone. Please assign a milestone before closing.`,
        });
      }
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(
          `Failed to post milestone warning comment to #${trackingIssueNumber}: ${err.message}`,
        );
      }
    }

    return {
      closed: false,
      trackingIssueNumber,
      reason: "missing_milestone",
    };
  }

  // All sub-issues are closed and milestone is set: post closure comment and close tracking issue
  try {
    await issuesGithub.rest.issues.createComment({
      owner,
      repo,
      issue_number: trackingIssueNumber,
      body: `All sub-issues for this tracking issue have been resolved and closed. Automatically closing tracking issue #${trackingIssueNumber}.`,
    });
  } catch (err) {
    if (core && typeof core.warning === "function") {
      core.warning(
        `Failed to post closure comment to #${trackingIssueNumber}: ${err.message}`,
      );
    }
  }

  await issuesGithub.rest.issues.update({
    owner,
    repo,
    issue_number: trackingIssueNumber,
    state: "closed",
    state_reason: "completed",
  });

  if (core && typeof core.info === "function") {
    core.info(`Successfully closed tracking issue #${trackingIssueNumber}.`);
  }

  // Progress GitHub Project V2 status to "Ready For Release"
  await updateIssueProjectStatus({
    github: issuesGithub || github,
    owner,
    repo,
    issueNumber: trackingIssueNumber,
    targetStatus: "Ready For Release",
    core,
  });

  return { closed: true, trackingIssueNumber };
}

/**
 * Transitions items in GitHub Projects V2 from fromStatus to toStatus.
 */
export async function transitionProjectItemsStatus({
  github,
  owner,
  repo,
  fromStatus = "Ready for Release",
  toStatus = "Done",
  projectNumber = DEFAULT_PROJECT_NUMBER,
  core,
}) {
  if (!github || typeof github.graphql !== "function") {
    if (core && typeof core.info === "function") {
      core.info(
        "GraphQL client unavailable; skipping project items status transition.",
      );
    }
    return { transitionedCount: 0 };
  }

  let projects = [];

  // Directly query the target project (e.g. Project 73)
  if (projectNumber) {
    try {
      const projectOwner = owner || DEFAULT_PROJECT_OWNER;
      const targetProjectQuery = `
        query getTargetProject($projectOwner: String!, $projectNumber: Int!) {
          organization(login: $projectOwner) {
            projectV2(number: $projectNumber) {
              id
              number
              title
              fields(first: 50) {
                nodes {
                  ... on ProjectV2SingleSelectField {
                    id
                    name
                    options {
                      id
                      name
                    }
                  }
                }
              }
            }
          }
        }
      `;
      const targetData = await github.graphql(targetProjectQuery, {
        projectOwner,
        projectNumber,
      });
      const targetProject = targetData?.organization?.projectV2;
      if (targetProject) {
        projects.push(targetProject);
        if (core && typeof core.info === "function") {
          core.info(
            `Found target project "${targetProject.title}" (#${targetProject.number}).`,
          );
        }
      }
    } catch (err) {
      if (core && typeof core.info === "function") {
        core.info(
          `Could not fetch target project #${projectNumber} directly: ${err.message}`,
        );
      }
    }
  }

  // Also query repository projects with pagination
  try {
    let hasNextPage = true;
    let cursor = null;

    while (hasNextPage) {
      const repoProjectsQuery = `
        query getRepoProjects($owner: String!, $repo: String!, $cursor: String) {
          repository(owner: $owner, name: $repo) {
            projectsV2(first: 50, after: $cursor) {
              pageInfo {
                hasNextPage
                endCursor
              }
              nodes {
                id
                number
                title
                fields(first: 50) {
                  nodes {
                    ... on ProjectV2SingleSelectField {
                      id
                      name
                      options {
                        id
                        name
                      }
                    }
                  }
                }
              }
            }
          }
        }
      `;
      const repoData = await github.graphql(repoProjectsQuery, {
        owner,
        repo,
        cursor,
      });
      const connection = repoData?.repository?.projectsV2;
      const repoProjects = connection?.nodes?.filter(Boolean) || [];
      for (const rp of repoProjects) {
        if (!projects.some((p) => p.id === rp.id)) {
          projects.push(rp);
        }
      }

      hasNextPage = Boolean(connection?.pageInfo?.hasNextPage);
      cursor = connection?.pageInfo?.endCursor || null;
    }
  } catch (err) {
    if (core && typeof core.warning === "function") {
      core.warning(`Could not fetch repository projects: ${err.message}`);
    }
  }

  // Query organization projects as additional fallback with pagination
  try {
    let hasNextPage = true;
    let cursor = null;

    while (hasNextPage) {
      const orgProjectsQuery = `
        query getOrgProjects($owner: String!, $cursor: String) {
          organization(login: $owner) {
            projectsV2(first: 50, after: $cursor) {
              pageInfo {
                hasNextPage
                endCursor
              }
              nodes {
                id
                number
                title
                fields(first: 50) {
                  nodes {
                    ... on ProjectV2SingleSelectField {
                      id
                      name
                      options {
                        id
                        name
                      }
                    }
                  }
                }
              }
            }
          }
        }
      `;
      const orgData = await github.graphql(orgProjectsQuery, {
        owner,
        cursor,
      });
      const connection = orgData?.organization?.projectsV2;
      const orgProjects = connection?.nodes?.filter(Boolean) || [];
      for (const op of orgProjects) {
        if (!projects.some((p) => p.id === op.id)) {
          projects.push(op);
        }
      }

      hasNextPage = Boolean(connection?.pageInfo?.hasNextPage);
      cursor = connection?.pageInfo?.endCursor || null;
    }
  } catch (err) {
    if (core && typeof core.info === "function") {
      core.info(
        `Could not fetch organization projects (may lack org scope): ${err.message}`,
      );
    }
  }

  if (projects.length === 0) {
    if (core && typeof core.info === "function") {
      core.info("No GitHub Projects found to transition.");
    }
    return { transitionedCount: 0 };
  }

  let transitionedCount = 0;

  for (const project of projects) {
    const fields = project.fields?.nodes?.filter(Boolean) || [];
    const statusField = fields.find(
      (f) => typeof f?.name === "string" && f.name.toLowerCase() === "status",
    );

    if (!statusField) continue;

    const fromOption = statusField.options?.find(
      (o) =>
        typeof o?.name === "string" &&
        o.name.trim().toLowerCase() === fromStatus.trim().toLowerCase(),
    );

    const toOption = statusField.options?.find(
      (o) =>
        typeof o?.name === "string" &&
        o.name.trim().toLowerCase() === toStatus.trim().toLowerCase(),
    );

    if (!fromOption || !toOption) continue;

    // Fetch items in project
    let hasNextPage = true;
    let endCursor = null;

    while (hasNextPage) {
      try {
        const itemsQuery = `
          query getProjectItems($projectId: ID!, $after: String) {
            node(id: $projectId) {
              ... on ProjectV2 {
                items(first: 100, after: $after) {
                  pageInfo {
                    hasNextPage
                    endCursor
                  }
                  nodes {
                    id
                    fieldValues(first: 20) {
                      nodes {
                        ... on ProjectV2ItemFieldSingleSelectValue {
                          field {
                            ... on ProjectV2SingleSelectField {
                              name
                            }
                          }
                          name
                          optionId
                        }
                      }
                    }
                    content {
                      ... on Issue {
                        number
                        title
                      }
                    }
                  }
                }
              }
            }
          }
        `;

        const itemsData = await github.graphql(itemsQuery, {
          projectId: project.id,
          after: endCursor,
        });

        const projectNode = itemsData?.node;
        const pageInfo = projectNode?.items?.pageInfo;
        const items = projectNode?.items?.nodes?.filter(Boolean) || [];

        for (const item of items) {
          const fieldValues = item.fieldValues?.nodes?.filter(Boolean) || [];
          const statusVal = fieldValues.find(
            (fv) =>
              fv?.field?.name?.toLowerCase() === "status" ||
              fv?.name?.trim().toLowerCase() ===
                fromStatus.trim().toLowerCase() ||
              fv?.optionId === fromOption.id,
          );

          if (
            statusVal &&
            (statusVal.optionId === fromOption.id ||
              statusVal.name?.trim().toLowerCase() ===
                fromStatus.trim().toLowerCase())
          ) {
            const mutation = `
              mutation updateItemStatus($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
                updateProjectV2ItemFieldValue(
                  input: {
                    projectId: $projectId
                    itemId: $itemId
                    fieldId: $fieldId
                    value: {
                      singleSelectOptionId: $optionId
                    }
                  }
                ) {
                  projectV2Item {
                    id
                  }
                }
              }
            `;

            await github.graphql(mutation, {
              projectId: project.id,
              itemId: item.id,
              fieldId: statusField.id,
              optionId: toOption.id,
            });

            const issueIdentifier = item.content?.number
              ? `#${item.content.number}`
              : item.id;
            if (core && typeof core.info === "function") {
              core.info(
                `Transitioned issue ${issueIdentifier} in project "${project.title}" from "${fromOption.name}" to "${toOption.name}".`,
              );
            }
            transitionedCount++;
          }
        }

        hasNextPage = Boolean(pageInfo?.hasNextPage);
        endCursor = pageInfo?.endCursor || null;
      } catch (err) {
        if (core && typeof core.warning === "function") {
          core.warning(
            `Failed to query/transition items for project "${project.title}": ${err.message}`,
          );
        }
        break;
      }
    }
  }

  if (core && typeof core.info === "function") {
    core.info(
      `Transitioned ${transitionedCount} item(s) from "${fromStatus}" to "${toStatus}".`,
    );
  }

  return { transitionedCount };
}

/**
 * release-done: Transitions all project items in "Ready for Release" to "Done" when a full release completes.
 */
export async function runReleaseDone({
  github,
  context,
  core,
  process = globalThis.process,
  getOctokit,
}) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";

  const issuesToken = process?.env?.GITHUB_ISSUES_TOKEN
    ? process.env.GITHUB_ISSUES_TOKEN.trim()
    : undefined;
  const octokitFactory =
    typeof getOctokit === "function"
      ? getOctokit
      : typeof github?.getOctokit === "function"
        ? github.getOctokit.bind(github)
        : undefined;
  const issuesGithub =
    issuesToken && typeof octokitFactory === "function"
      ? octokitFactory(issuesToken)
      : github;

  const headSha = (
    process?.env?.HEAD_SHA ||
    context?.payload?.workflow_run?.head_sha ||
    ""
  ).trim();

  // Determine release tag
  let tag = (process?.env?.TAG || "").trim();

  if (tag) {
    // Verify explicit tag against GitHub releases API to ensure it exists and is not a prerelease
    if (issuesGithub?.rest?.repos?.getReleaseByTag) {
      try {
        const { data: release } = await issuesGithub.rest.repos.getReleaseByTag({
          owner,
          repo,
          tag,
        });
        if (release.prerelease || /rc/i.test(release.tag_name || tag)) {
          if (core && typeof core.info === "function") {
            core.info(
              `Release "${tag}" is a pre-release/RC; skipping transition to Done.`,
            );
          }
          return { skipped: true, reason: "prerelease", tag };
        }
      } catch (err) {
        if (core && typeof core.warning === "function") {
          core.warning(
            `Specified tag "${tag}" could not be verified as an existing release: ${err.message}`,
          );
        }
        return {
          skipped: true,
          reason: "release_not_found",
          tag,
          error: err.message,
        };
      }
    } else if (/rc/i.test(tag)) {
      return { skipped: true, reason: "prerelease", tag };
    }
  } else if (headSha && issuesGithub?.rest?.repos?.listReleases) {
    try {
      const { data: releases } = await issuesGithub.rest.repos.listReleases({
        owner,
        repo,
        per_page: 20,
      });

      const matchingRelease = releases?.find(
        (r) =>
          r.target_commitish === headSha ||
          r.tag_name === headSha,
      );

      if (!matchingRelease) {
        if (core && typeof core.info === "function") {
          core.info(
            `No release found matching HEAD_SHA "${headSha}"; skipping transition to Done.`,
          );
        }
        return { skipped: true, reason: "unverified_tag", headSha };
      }

      if (
        matchingRelease.prerelease ||
        /rc/i.test(matchingRelease.tag_name || "")
      ) {
        if (core && typeof core.info === "function") {
          core.info(
            `Release ${matchingRelease.tag_name} is a pre-release/RC; skipping transition to Done.`,
          );
        }
        return {
          skipped: true,
          reason: "prerelease",
          tag: matchingRelease.tag_name,
        };
      }
      tag = matchingRelease.tag_name;
    } catch (err) {
      if (core && typeof core.warning === "function") {
        core.warning(`Could not fetch releases: ${err.message}`);
      }
      return {
        skipped: true,
        reason: "release_lookup_failed",
        error: err.message,
      };
    }
  } else {
    if (core && typeof core.info === "function") {
      core.info(
        "Neither verified TAG nor HEAD_SHA provided; skipping transition to Done.",
      );
    }
    return { skipped: true, reason: "unverified_tag" };
  }

  if (!tag) {
    if (core && typeof core.info === "function") {
      core.info(
        "Could not verify completed release tag; skipping transition to Done.",
      );
    }
    return { skipped: true, reason: "unverified_tag" };
  }

  if (/rc/i.test(tag)) {
    if (core && typeof core.info === "function") {
      core.info(`Tag "${tag}" is an RC; skipping transition to Done.`);
    }
    return { skipped: true, reason: "prerelease", tag };
  }

  if (core && typeof core.info === "function") {
    core.info(
      `Full release verified (tag: ${tag}). Transitioning 'Ready for Release' issues to 'Done'...`,
    );
  }

  const projectNumber = parseInt(
    process?.env?.PROJECT_NUMBER || String(DEFAULT_PROJECT_NUMBER),
    10,
  );

  return await transitionProjectItemsStatus({
    github: issuesGithub,
    owner,
    repo,
    fromStatus: "Ready for Release",
    toStatus: "Done",
    projectNumber: !isNaN(projectNumber) ? projectNumber : DEFAULT_PROJECT_NUMBER,
    core,
  });
}
