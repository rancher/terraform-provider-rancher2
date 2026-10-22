/**
 * auto-merge.js: Gating and automation logic for the Automated PR Merge System on main.
 * Validates PR eligibility, reviewer permissions (triage/write/admin), file paths, and changes requested,
 * then submits a bot approval and enables native auto-merge via GraphQL.
 */

/**
 * resolveTargetPr: Resolves the target pull request object from payload, PR number, or workflow_run commit fallback.
 */
export async function resolveTargetPr({ github, context, core, process = globalThis.process }) {
  if (context?.payload?.pull_request) {
    return context.payload.pull_request;
  }

  const owner = context?.repo?.owner || process?.env?.GITHUB_REPOSITORY?.split('/')[0];
  const repo = context?.repo?.repo || process?.env?.GITHUB_REPOSITORY?.split('/')[1];

  let prNumber = null;
  if (process?.env?.PR_NUMBER) {
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
    core.info(`Fetching PR #${prNumber} from GitHub API...`);
    try {
      const response = await github.rest.pulls.get({
        owner,
        repo,
        pull_number: prNumber,
      });
      return response.data;
    } catch (error) {
      core.setFailed(`Failed to fetch PR #${prNumber}: ${error.message}`);
      return null;
    }
  }

  const headSha = parentRun?.head_sha || process?.env?.WORKFLOW_RUN_SHA;
  if (headSha && typeof github?.rest?.repos?.listPullRequestsAssociatedWithCommit === 'function') {
    core.info(`Attempting fallback PR resolution for commit SHA: ${headSha}...`);
    try {
      let associatedPrs = [];
      if (typeof github.paginate === 'function') {
        associatedPrs = await github.paginate(
          github.rest.repos.listPullRequestsAssociatedWithCommit,
          {
            owner,
            repo,
            commit_sha: headSha,
          }
        );
      } else {
        const res = await github.rest.repos.listPullRequestsAssociatedWithCommit({
          owner,
          repo,
          commit_sha: headSha,
        });
        associatedPrs = res?.data || res || [];
      }

      const parentBranch = parentRun?.head_branch ? parentRun.head_branch.replace(/^refs\/heads\//, '') : null;
      const matched = associatedPrs.find(p => {
        if (p.state && p.state !== 'open') {
          return false;
        }
        const prHeadRef = p.head?.ref ? p.head.ref.replace(/^refs\/heads\//, '') : null;
        if (parentBranch && prHeadRef !== parentBranch) {
          return false;
        }
        if (parentRun?.head_repository?.full_name && p.head?.repo?.full_name?.toLowerCase() !== parentRun.head_repository.full_name.toLowerCase()) {
          return false;
        }
        return true;
      });

      if (matched) {
        core.info(`Resolved PR #${matched.number} via commit search fallback.`);
        return matched;
      }
    } catch (error) {
      core.warning(`Commit fallback failed to list PRs: ${error.message}`);
    }
  }

  return null;
}

/**
 * checkWorkflowModifications: Checks if any files modified in the PR reside under .github/.
 */
export async function checkWorkflowModifications({ github, owner, repo, prNumber, core }) {
  let files = [];
  try {
    if (typeof github?.paginate === 'function' && typeof github?.rest?.pulls?.listFiles === 'function') {
      files = await github.paginate(github.rest.pulls.listFiles, {
        owner,
        repo,
        pull_number: prNumber,
      });
    } else if (typeof github?.rest?.pulls?.listFiles === 'function') {
      const res = await github.rest.pulls.listFiles({
        owner,
        repo,
        pull_number: prNumber,
      });
      files = res?.data || res || [];
    } else if (typeof github?.paginate === 'function') {
      files = await github.paginate();
    }
  } catch (error) {
    core.setFailed(`Failed to list files for PR #${prNumber}: ${error.message}`);
    throw error;
  }

  const workflowFiles = files.filter(f => f.filename && f.filename.startsWith('.github/'));
  return {
    modifiesWorkflows: workflowFiles.length > 0,
    workflowFiles: workflowFiles.map(f => f.filename),
  };
}

/**
 * evaluateReviews: Determines the latest review state for each reviewer.
 */
export function evaluateReviews(reviews, authorLogin) {
  const DECISION_STATES = new Set(['APPROVED', 'CHANGES_REQUESTED', 'DISMISSED']);
  const normalizedAuthor = typeof authorLogin === 'string' ? authorLogin.toLowerCase() : null;
  const latestByUser = new Map();

  for (const review of reviews || []) {
    if (!review?.user?.login) {
      continue;
    }
    const state = typeof review.state === 'string' ? review.state.toUpperCase() : '';
    if (!DECISION_STATES.has(state)) {
      continue;
    }

    const userKey = review.user.login.toLowerCase();
    if (state === 'DISMISSED') {
      latestByUser.delete(userKey);
    } else {
      latestByUser.set(userKey, review);
    }
  }

  const latestReviews = Array.from(latestByUser.values());
  const hasChangesRequested = latestReviews.some(
    r => (typeof r.state === 'string' ? r.state.toUpperCase() : '') === 'CHANGES_REQUESTED'
  );
  const approvedReviewers = latestReviews
    .filter(r => {
      const isApproved = (typeof r.state === 'string' ? r.state.toUpperCase() : '') === 'APPROVED';
      const isAuthor = normalizedAuthor && r.user.login.toLowerCase() === normalizedAuthor;
      return isApproved && !isAuthor;
    })
    .map(r => r.user.login);

  return {
    hasChangesRequested,
    approvedReviewers,
    latestReviews,
  };
}

/**
 * checkReviewerPermission: Checks if a user has triage, write, or admin permission in the repository.
 */
export async function checkReviewerPermission({ github, owner, repo, username, core }) {
  const allowedRoles = ['triage', 'write', 'admin'];
  try {
    const res = await github.rest.repos.getCollaboratorPermissionLevel({
      owner,
      repo,
      username,
    });
    const permission = res?.data?.permission;
    const roleName = res?.data?.role_name || permission;
    core.info(`User @${username} has repository permission "${permission}" (role: "${roleName}").`);

    const isAuthorized = allowedRoles.includes(permission) || allowedRoles.includes(roleName);
    return {
      isAuthorized,
      permission,
      roleName,
    };
  } catch (error) {
    core.warning(`Failed to retrieve collaborator permission for @${username}: ${error.message}`);
    return {
      isAuthorized: false,
      permission: null,
      roleName: null,
      error: error.message,
    };
  }
}

/**
 * runAutoMerge: Main orchestration entrypoint for validating and enabling auto-merge on qualifying PRs.
 */
export async function runAutoMerge({ github, context, core, process = globalThis.process, getOctokit }) {
  const owner = context?.repo?.owner || process?.env?.GITHUB_REPOSITORY?.split('/')[0];
  const repo = context?.repo?.repo || process?.env?.GITHUB_REPOSITORY?.split('/')[1];

  if (!owner || !repo) {
    core.setFailed('Missing repository owner and repo in context.');
    return false;
  }

  const pr = await resolveTargetPr({ github, context, core, process });
  if (!pr) {
    core.setFailed('Could not resolve target PR for auto-merge.');
    return false;
  }

  const prNumber = pr.number;
  core.info(`Evaluating PR #${prNumber} for auto-merge...`);

  // Gate 1: State, Draft, and Target Branch Checks
  if (pr.state !== 'open') {
    core.info(`PR #${prNumber} is not open (state: "${pr.state}"). Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'skipped_not_open');
    return false;
  }

  if (pr.draft) {
    core.info(`PR #${prNumber} is a draft. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'skipped_draft');
    return false;
  }

  const baseRef = pr.base?.ref ? pr.base.ref.replace(/^refs\/heads\//, '') : '';
  if (baseRef !== 'main') {
    core.info(`PR #${prNumber} targets branch "${baseRef}", not "main". Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'skipped_not_main');
    return false;
  }

  const triggerPrHeadSha = process?.env?.PR_HEAD_SHA;
  if (triggerPrHeadSha && pr.head?.sha && triggerPrHeadSha !== pr.head.sha) {
    core.info(`PR #${prNumber} head commit (${triggerPrHeadSha}) is superseded by current head (${pr.head.sha}). Skipping stale run.`);
    core.setOutput('auto_merge_status', 'skipped_superseded');
    return false;
  }

  // Gate 2: Exclude modifications to .github/**
  const { modifiesWorkflows, workflowFiles } = await checkWorkflowModifications({
    github,
    owner,
    repo,
    prNumber,
    core,
  });

  if (modifiesWorkflows) {
    core.info(`PR #${prNumber} modifies files in .github/ (${workflowFiles.join(', ')}). Workflow changes require human write approval. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'skipped_workflow_files');
    return false;
  }

  // Gate 3: Reviews and Approvals Checks
  let reviews = [];
  try {
    if (typeof github?.paginate === 'function' && typeof github?.rest?.pulls?.listReviews === 'function') {
      reviews = await github.paginate(github.rest.pulls.listReviews, {
        owner,
        repo,
        pull_number: prNumber,
      });
    } else if (typeof github?.rest?.pulls?.listReviews === 'function') {
      const res = await github.rest.pulls.listReviews({
        owner,
        repo,
        pull_number: prNumber,
      });
      reviews = res?.data || res || [];
    } else if (typeof github?.paginate === 'function') {
      reviews = await github.paginate();
    }
  } catch (error) {
    core.setFailed(`Failed to list reviews for PR #${prNumber}: ${error.message}`);
    return false;
  }

  const authorLogin = pr.user?.login;
  const { hasChangesRequested, approvedReviewers } = evaluateReviews(reviews, authorLogin);

  if (hasChangesRequested) {
    core.info(`PR #${prNumber} has one or more active CHANGES_REQUESTED reviews. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'blocked_changes_requested');
    return false;
  }

  if (approvedReviewers.length === 0) {
    core.info(`PR #${prNumber} has no approved reviews from non-author users. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'no_approvals');
    return false;
  }

  // Gate 4: Verify Reviewer Permissions
  let qualifyingReviewer = null;
  for (const username of approvedReviewers) {
    const { isAuthorized, roleName, permission } = await checkReviewerPermission({
      github,
      owner,
      repo,
      username,
      core,
    });
    if (isAuthorized) {
      qualifyingReviewer = { username, role: roleName || permission };
      break;
    }
  }

  if (!qualifyingReviewer) {
    core.info(`PR #${prNumber} does not have an approved review from a user with triage, write, or admin permissions. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'no_qualifying_permission');
    return false;
  }

  core.info(`PR #${prNumber} qualified for auto-merge via sign-off from @${qualifyingReviewer.username} (role: ${qualifyingReviewer.role}).`);

  // Initialize elevated Octokit instance with GITHUB_MERGE_TOKEN
  const mergeToken = process?.env?.GITHUB_MERGE_TOKEN ? process.env.GITHUB_MERGE_TOKEN.trim() : undefined;
  if (!mergeToken) {
    core.setFailed('GITHUB_MERGE_TOKEN is required to approve PRs and enable auto-merge.');
    return false;
  }
  const octokitFactory = (typeof getOctokit === 'function')
    ? getOctokit
    : (typeof github?.getOctokit === 'function' ? github.getOctokit.bind(github) : undefined);
  if (typeof octokitFactory !== 'function') {
    core.setFailed('Cannot initialize the merge-token client: getOctokit is unavailable.');
    return false;
  }
  const prGithub = octokitFactory(mergeToken);

  // Gate 5: Check if Bot already approved the current head commit
  let botLogin = null;
  try {
    if (typeof prGithub.rest?.users?.getAuthenticated === 'function') {
      const { data: authUser } = await prGithub.rest.users.getAuthenticated();
      botLogin = authUser?.login;
    }
  } catch (error) {
    core.info(`Could not resolve authenticated bot login: ${error.message}`);
  }

  const alreadyApprovedByBot = botLogin && reviews.some(
    r => r.user?.login === botLogin && r.commit_id === pr.head.sha && r.state === 'APPROVED'
  );

  if (!alreadyApprovedByBot) {
    core.info(`Submitting bot approval on PR #${prNumber} for commit ${pr.head.sha}...`);
    try {
      await prGithub.rest.pulls.createReview({
        owner,
        repo,
        pull_number: prNumber,
        commit_id: pr.head.sha,
        event: 'APPROVE',
        body: `Automated approval triggered by @${qualifyingReviewer.username} (role: ${qualifyingReviewer.role}).`,
      });
      core.info(`Successfully submitted bot approval on PR #${prNumber}.`);
    } catch (error) {
      let approvalConfirmed = false;
      if (botLogin) {
        try {
          const currentReviews = await github.paginate(github.rest.pulls.listReviews, {
            owner,
            repo,
            pull_number: prNumber,
          });
          const { latestReviews } = evaluateReviews(currentReviews, authorLogin);
          approvalConfirmed = latestReviews.some(
            r => r.user.login === botLogin && r.commit_id === pr.head.sha && r.state === 'APPROVED'
          );
        } catch (verificationError) {
          core.warning(`Could not verify bot approval on PR #${prNumber}: ${verificationError.message}`);
        }
      }
      if (!approvalConfirmed) {
        core.setFailed(`Failed to submit bot approval on PR #${prNumber}: ${error.message}`);
        return false;
      }
      core.warning(`Bot approval on PR #${prNumber} is already present for commit ${pr.head.sha}: ${error.message}`);
    }
  } else {
    core.info(`Bot (@${botLogin}) has already approved commit ${pr.head.sha} on PR #${prNumber}.`);
  }

  // Gate 6: Enable Native Auto-Merge via GraphQL
  core.info(`Enabling native auto-merge for PR #${prNumber} (node_id: ${pr.node_id})...`);
  try {
    await prGithub.graphql(`
      mutation EnableAutoMerge($pullRequestId: ID!) {
        enablePullRequestAutoMerge(input: {
          pullRequestId: $pullRequestId,
          mergeMethod: SQUASH
        }) {
          pullRequest {
            id
            autoMergeRequest {
              enabledAt
            }
          }
        }
      }
    `, { pullRequestId: pr.node_id });

    core.info(`Successfully enabled native auto-merge for PR #${prNumber}.`);
    core.setOutput('auto_merge_status', 'enabled');
    return true;
  } catch (error) {
    if (error.message && (error.message.includes('already enabled') || error.message.includes('Auto merge is already enabled'))) {
      core.info(`Native auto-merge is already enabled for PR #${prNumber}.`);
      core.setOutput('auto_merge_status', 'already_enabled');
      return true;
    }

    core.setFailed(`Failed to enable native auto-merge for PR #${prNumber}: ${error.message}`);
    return false;
  }
}

export default async ({ github, context, core, process = globalThis.process, getOctokit }) => {
  return await runAutoMerge({ github, context, core, process, getOctokit });
};
