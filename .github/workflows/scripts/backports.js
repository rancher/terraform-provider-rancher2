import { execFile } from 'node:child_process';
import { promisify } from 'node:util';

const execFileAsync = promisify(execFile);

/**
 * Executes a git command asynchronously and extracts detailed stderr on failure.
 */
async function runGit(args, options = {}) {
  try {
    return await execFileAsync('git', args, { encoding: 'utf8', ...options });
  } catch (error) {
    const rawDetails = error.stderr?.trim() || error.stdout?.trim() || error.message;
    const details = rawDetails
      ?.replace(/authorization:\s*basic\s+[^\s'"]+/gi, 'AUTHORIZATION: basic [REDACTED]')
      ?.replace(/https:\/\/[^@\s]+@github\.com/gi, 'https://[REDACTED]@github.com');
    const sanitizedArgs = args.map(arg =>
      (arg.toLowerCase().includes('authorization:') ? '[REDACTED]' : arg)
    );
    const sanitizedCommand = `git ${sanitizedArgs.join(' ')}`;
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
    core.warning(`Could not parse TERRAFORM_MAINTAINERS: ${err.message}. Defaulting to no assignees.`);
    return [];
  }
}

export default async ({ github, context, core, process = globalThis.process, getOctokit }) => {
  const mode = process?.env?.SCRIPT_MODE;
  switch (mode) {
  case 'wait-for-settle':
    return await runWaitForSettle({ github, context, core, process });
  case 'backport-pr':
    return await runBackportPr({ github, context, core, process, getOctokit });
  case 'backport-issues':
    return await runBackportIssues({ github, context, core, process });
  case 'merge-label':
    return await runMergeLabel({ github, context, core });
  default:
    throw new Error(`Unknown backport script mode: ${mode}`);
  }
};

/**
 * wait-for-settle: Allows GitHub's API time to index merge commits and retrieve the PR list with exponential back-off.
 */
async function runWaitForSettle({ github, context, core, process = globalThis.process }) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";
  
  // Handle input from either manual dispatch or push
  const mergeCommitSha = process?.env?.MERGE_COMMIT_SHA || context?.payload?.head_commit?.id;
  
  if (!mergeCommitSha) {
    core.setFailed("No merge commit SHA found in environment or context payload.");
    return;
  }

  const maxAttempts = 5;
  let delayMs = 3000;

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const prs = await github.paginate(github.rest.repos.listPullRequestsAssociatedWithCommit, {
        owner,
        repo,
        commit_sha: mergeCommitSha
      });

      const mergedMainPr = prs?.find(p => p.base?.ref === 'main' && p.merged_at);
      if (mergedMainPr) {
        core.info(`Found merged PR #${mergedMainPr.number} associated with commit ${mergeCommitSha} targeting main on attempt ${attempt}.`);
        core.setOutput('merge_commit_sha', mergeCommitSha);
        return;
      }

      if (attempt < maxAttempts) {
        core.info(`Attempt ${attempt}/${maxAttempts}: Merged PR targeting main for commit ${mergeCommitSha} not found yet. Retrying in ${delayMs / 1000}s...`);
        await new Promise(resolve => setTimeout(resolve, delayMs));
        delayMs *= 2;
      } else {
        core.setFailed(`No merged PR targeting main associated with commit ${mergeCommitSha} found after ${maxAttempts} attempts.`);
        return;
      }
    } catch (error) {
      if (attempt < maxAttempts) {
        core.warning(`Attempt ${attempt}/${maxAttempts} failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}. Retrying in ${delayMs / 1000}s...`);
        await new Promise(resolve => setTimeout(resolve, delayMs));
        delayMs *= 2;
      } else {
        core.setFailed(`Failed to retrieve PRs associated with commit ${mergeCommitSha} after ${maxAttempts} attempts: ${error.message}`);
        return;
      }
    }
  }
}

/**
 * backport-pr: Cherry-picks a commit to the appropriate backport branches and creates PRs.
 */
async function runBackportPr({ github, context, core, process = globalThis.process, getOctokit }) {
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
  const mergeToken = process?.env?.GITHUB_MERGE_TOKEN ? process.env.GITHUB_MERGE_TOKEN.trim() : undefined;
  const octokitFactory = (typeof getOctokit === 'function')
    ? getOctokit
    : (typeof github?.getOctokit === 'function' ? github.getOctokit.bind(github) : undefined);
  const prGithub = (mergeToken && typeof octokitFactory === 'function') ? octokitFactory(mergeToken) : github;
  let response;

  try {
    response = await github.rest.repos.listPullRequestsAssociatedWithCommit({
      owner,
      repo,
      commit_sha: mergeCommitSha
    });
  } catch (error) {
    throw new Error(`Failed to retrieve PRs associated with commit ${mergeCommitSha}: ${error.message}`);
  }
  const associatedPrs = response.data;
  if (associatedPrs.length === 0) {
    core.info(`No PRs associated with commit ${mergeCommitSha}. Exiting.`);
    return;
  }

  const pr = associatedPrs.find(p => p.base.ref === 'main' && p.merged_at);
  if (!pr) {
    core.info(`No merged PR found for commit ${mergeCommitSha}.`);
    return;
  }
  core.info(`Found associated PR: #${pr.number}`);

  core.info(`Searching for 'internal/tracking' issue linked to PR #${pr.number}`);
  try {
    response = await github.request('GET /search/issues', {
      q: `repo:${owner}/${repo} is:issue state:open label:"internal/tracking" in:body #${pr.number}`,
      advanced_search: true,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    throw new Error(`Failed to search for internal/tracking issue for PR #${pr.number}: ${error.message}`);
  }
  const searchResults = response.data;
  if (searchResults.total_count === 0) {
    core.info(`No 'internal/tracking' issue found for PR #${pr.number}. Exiting.`);
    return;
  }
  const trackingIssue = searchResults.items[0];
  core.info(`Found tracking issue: #${trackingIssue.number}`);

  core.info(`Fetching sub-issues for tracking issue #${trackingIssue.number}`);
  try {
    response = await github.request('GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues', {
      owner: owner,
      repo: repo,
      issue_number: trackingIssue.number,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    throw new Error(`Failed to fetch sub-issues for tracking issue #${trackingIssue.number}: ${error.message}`);
  }
  const subIssues = response.data;
  core.info(`Sub-issues data: ${JSON.stringify(subIssues)}`);
  if (!Array.isArray(subIssues)) {
    core.warning(`Unexpected sub-issues data format: ${JSON.stringify(subIssues)}`);
    return;
  }
  if (subIssues.length === 0) {
    core.info(`No sub-issues found for issue #${trackingIssue.number}. Exiting.`);
    return;
  }
  core.info(`Found ${subIssues.length} sub-issues.`);

  const failedBackports = [];

  try {
    try {
      await runGit(['config', '--local', 'user.name', 'github-actions[bot]']);
      await runGit(['config', '--local', 'user.email', 'github-actions[bot]@users.noreply.github.com']);
      if (mergeToken) {
        try {
          const includeConfig = await runGit(['config', '--local', '--name-only', '--get-regexp', '^includeif\\.gitdir:']);
          const includeKeys = [...new Set(includeConfig.stdout.trim().split('\n').map(k => k.trim()).filter(Boolean))];
          for (const key of includeKeys) {
            await runGit(['config', '--local', '--unset-all', key]);
          }
        } catch {
          // ignore if no includeif.gitdir keys exist
        }

        try {
          const headerConfig = await runGit(['config', '--local', '--name-only', '--get-regexp', 'http\\..*extraheader']);
          const headerKeys = [...new Set(headerConfig.stdout.trim().split('\n').map(k => k.trim()).filter(Boolean))];
          for (const key of headerKeys) {
            await runGit(['config', '--local', '--unset-all', key]);
          }
        } catch {
          // ignore if no extraheader keys exist
        }

        core.setSecret(mergeToken);
        const b64Token = Buffer.from(`x-access-token:${mergeToken}`).toString('base64');
        core.setSecret(b64Token);
        await runGit(['config', '--local', '--replace-all', 'http.https://github.com/.extraheader', `AUTHORIZATION: basic ${b64Token}`]);
      }
    } catch (error) {
      throw new Error(`Failed to configure git credentials: ${error.message}`);
    }

    for (const subIssue of subIssues) {
      const subIssueNumber = subIssue.number;
      core.info(`Processing sub-issue #${subIssueNumber}...`);

      const releaseLabel = subIssue.labels?.find(label => label.name?.startsWith('release/v'));

      if (!releaseLabel) {
        core.warning(`Sub-issue #${subIssueNumber} has no 'release/v...' label. Skipping.`);
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

      core.info(`Processing sub-issue #${subIssueNumber} for target branch: ${targetBranch}`);
      const newBranchName = `backport-${pr.number}-${targetBranch.replace(/\//g, '-')}`;

      try {
        await runGit(['fetch', 'origin', targetBranch]);
        await runGit(['checkout', '-B', newBranchName, `origin/${targetBranch}`]);
        await runGit(['cherry-pick', '--allow-empty', '-x', mergeCommitSha, '-X', 'theirs']);
        await runGit(['push', 'origin', newBranchName]);

        core.info(`Creating pull request for branch ${newBranchName} targeting ${targetBranch}...`);
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
            `${pr.body}`
          ].join("\n\n")
        });

        const newPR = response.data;
        core.info(`Created backport PR data: ${JSON.stringify(newPR)}`);
        const prNumber = newPR.number;

        try {
          await prGithub.graphql(`
            mutation EnableAutoMerge($pullRequestId: ID!) {
              enablePullRequestAutoMerge(input: {
                pullRequestId: $pullRequestId,
                mergeMethod: SQUASH
              }) {
                pullRequest { id }
              }
            }
          `, { pullRequestId: newPR.node_id });
          core.info(`Auto-merge enabled for PR #${prNumber}`);
        } catch (error) {
          core.warning(`Failed to enable auto-merge on PR #${prNumber}: ${error.message}`);
        }

        if (assignees.length > 0) {
          try {
            await github.rest.issues.addAssignees({
              owner,
              repo,
              issue_number: prNumber,
              assignees: assignees
            });
          } catch (error) {
            throw new Error(`Failed to assign PR #${prNumber}: ${error.message}`);
          }
        }

        try {
          await github.rest.issues.addLabels({
            owner,
            repo,
            issue_number: prNumber,
            labels: ["internal/pr-backport", targetBranch]
          });
        } catch (error) {
          throw new Error(`Failed to add backport label to PR #${prNumber}: ${error.message}`);
        }
      } catch (error) {
        core.error(`Failed to backport to ${targetBranch} for sub-issue #${subIssueNumber}: ${error.message}`);
        failedBackports.push({ subIssueNumber, targetBranch, error: error.message });
        try {
          await runGit(['cherry-pick', '--abort']);
        } catch {
          // ignore if cherry-pick was not active
        }
        try {
          await runGit(['checkout', 'main', '-f']);
        } catch {
          // ignore fallback checkout failure
        }
      }
    }
  } finally {
    if (mergeToken) {
      try {
        await runGit(['config', '--local', '--unset-all', 'http.https://github.com/.extraheader']);
      } catch {
        // ignore if unsetting extraheader fails or key was not present
      }
    }
  }

  if (failedBackports.length > 0) {
    const errorSummary = failedBackports
      .map(f => `Sub-issue #${f.subIssueNumber} (${f.targetBranch}): ${f.error}`)
      .join('\n');
    throw new Error(`Failed to complete ${failedBackports.length} backport(s):\n${errorSummary}`);
  }
}

/**
 * backport-issues: Creates sub-issues for tracking backports.
 */
async function runBackportIssues({ github, context, core, process = globalThis.process }) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";
  const releaseLabel = context.payload.label.name;
  const parentIssue = context.payload.issue;
  const parentIssueTitle = parentIssue.title;
  const parentIssueNumber = parentIssue.number;
  const assignees = parseMaintainers(core, process?.env?.TERRAFORM_MAINTAINERS);
  const extractedPrNumber = parseInt(process?.env?.PR, 10);
  if (isNaN(extractedPrNumber)) {
    throw new Error(`Invalid PR number: ${process?.env?.PR}`);
  }
  let response;

  try {
    response = await github.rest.issues.get({
      owner: owner,
      repo: repo,
      issue_number: extractedPrNumber
    });
  } catch (error) {
    throw new Error(`Failed to retrieve PR #${extractedPrNumber}: ${error.message}`);
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
        `${pr.body}`
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
  core.info(`New backport issue data: ${JSON.stringify(newIssue)}`);
  const subIssueId = newIssue.id;

  try {
    await github.request('POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues', {
      owner: owner,
      repo: repo,
      issue_number: parentIssueNumber,
      sub_issue_id: subIssueId,
      headers: {
        'X-GitHub-Api-Version': '2022-11-28'
      }
    });
  } catch (error) {
    throw new Error(`Failed to link backport issue to tracking issue: ${error.message}`);
  }
}

/**
 * merge-label: Adds internal/merged label to issues referenced in a merged PR body.
 */
async function runMergeLabel({ github, context, core }) {
  const owner = context?.repo?.owner || "rancher";
  const repo = context?.repo?.repo || "terraform-provider-rancher2";
  const pr = context.payload.pull_request;

  const issueRegex = /#(\d+)/g;
  const prBody = pr.body ?? "";
  const matches = prBody.matchAll(issueRegex);
  const issueNumbers = [...new Set(Array.from(matches, m => parseInt(m[1], 10)))];

  core.info(`Found issue numbers in PR body: ${issueNumbers}`);

  for (const issueNumber of issueNumbers) {
    try {
      const { data: issueData } = await github.rest.issues.get({
        owner,
        repo,
        issue_number: issueNumber,
      });

      if (!issueData.pull_request && issueData.labels.some(l => l.name === 'internal/backport')) {
        core.info(`Adding 'internal/merged' label to issue #${issueNumber}`);
        await github.rest.issues.addLabels({
          owner,
          repo,
          issue_number: issueNumber,
          labels: ["internal/merged"]
        });
      }
    } catch (error) {
      core.setFailed(`Could not process issue #${issueNumber}: ${error.message}`);
    }
  }
}
