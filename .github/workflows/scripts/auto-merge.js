/**
 * auto-merge.js: Gating and automation logic for the Automated PR Merge System on main.
 * Validates PR eligibility, reviewer permissions (triage/write/admin), file paths, and changes requested,
 * then submits a bot approval and enables native auto-merge via GraphQL.
 */

import fs from 'node:fs';
import path from 'node:path';

/**
 * resolveApprovalSignal: Resolves triggering approval signal from payload, env vars, signal file, or Actions API.
 */
export async function resolveApprovalSignal({ github, context, core, process = globalThis.process }) {
  // 1. Direct pull_request_review event payload
  if (context?.payload?.review && context?.payload?.pull_request) {
    return {
      pr_number: context.payload.pull_request.number,
      review_id: context.payload.review.id,
      review_state: context.payload.review.state,
      review_user: context.payload.review.user?.login,
      review_commit_id: context.payload.review.commit_id,
      head_sha: context.payload.pull_request.head?.sha,
      base_ref: context.payload.pull_request.base?.ref,
    };
  }

  // 2. Direct JSON string passed via APPROVAL_SIGNAL env var
  if (process?.env?.APPROVAL_SIGNAL) {
    try {
      const parsed = typeof process.env.APPROVAL_SIGNAL === 'string'
        ? JSON.parse(process.env.APPROVAL_SIGNAL)
        : process.env.APPROVAL_SIGNAL;
      if (parsed && typeof parsed === 'object') {
        return parsed;
      }
    } catch (err) {
      core.warning(`Failed to parse APPROVAL_SIGNAL env var: ${err.message}`);
    }
  }

  // 3. Signal file on disk (from actions/download-artifact step)
  const candidatePaths = [];
  if (process?.env?.APPROVAL_SIGNAL_PATH) {
    candidatePaths.push(process.env.APPROVAL_SIGNAL_PATH);
    if (process?.env?.GITHUB_WORKSPACE) {
      candidatePaths.push(path.resolve(process.env.GITHUB_WORKSPACE, process.env.APPROVAL_SIGNAL_PATH));
    }
  }
  const workspace = process?.env?.GITHUB_WORKSPACE || process?.cwd?.() || '.';
  candidatePaths.push(
    path.resolve(workspace, 'signal/approval-signal.json'),
    path.resolve(workspace, 'approval-signal.json'),
    'signal/approval-signal.json',
    'approval-signal.json'
  );

  for (const filePath of candidatePaths) {
    try {
      if (fs.existsSync(filePath)) {
        const raw = fs.readFileSync(filePath, 'utf8');
        const parsed = JSON.parse(raw);
        if (parsed && typeof parsed === 'object') {
          return parsed;
        }
      }
    } catch (err) {
      core.warning(`Failed to read approval signal file at ${filePath}: ${err.message}`);
    }
  }

  // 4. Fallback to TRIGGER_REVIEW_* env variables if set from previous job outputs
  if (process?.env?.TRIGGER_REVIEW_ID && process?.env?.TRIGGER_REVIEW_USER) {
    return {
      pr_number: process.env.PR_NUMBER ? Number.parseInt(process.env.PR_NUMBER, 10) : undefined,
      review_id: process.env.TRIGGER_REVIEW_ID,
      review_user: process.env.TRIGGER_REVIEW_USER,
      review_commit_id: process.env.TRIGGER_REVIEW_COMMIT_ID || process.env.PR_HEAD_SHA,
      head_sha: process.env.PR_HEAD_SHA || process.env.TRIGGER_REVIEW_COMMIT_ID,
      review_state: 'APPROVED',
      base_ref: 'main',
    };
  }

  // 5. Direct API download fallback if in workflow_run context
  const parentRun = context?.payload?.workflow_run;
  const owner = context?.repo?.owner || process?.env?.GITHUB_REPOSITORY?.split('/')[0];
  const repo = context?.repo?.repo || process?.env?.GITHUB_REPOSITORY?.split('/')[1];

  if (parentRun?.id && owner && repo && typeof github?.rest?.actions?.downloadArtifact === 'function') {
    try {
      let artifacts = [];
      if (typeof github.paginate === 'function' && typeof github?.rest?.actions?.listWorkflowRunArtifacts === 'function') {
        artifacts = await github.paginate(github.rest.actions.listWorkflowRunArtifacts, {
          owner,
          repo,
          run_id: parentRun.id,
        });
      } else if (typeof github?.rest?.actions?.listWorkflowRunArtifacts === 'function') {
        const res = await github.rest.actions.listWorkflowRunArtifacts({
          owner,
          repo,
          run_id: parentRun.id,
        });
        artifacts = res?.data?.artifacts || res?.artifacts || [];
      }

      const signalArtifact = artifacts.find(a => a.name === 'auto-merge-approval-signal');
      if (signalArtifact) {
        const downloadRes = await github.rest.actions.downloadArtifact({
          owner,
          repo,
          artifact_id: signalArtifact.id,
          archive_format: 'zip',
        });
        if (downloadRes?.data) {
          if (typeof downloadRes.data === 'string') {
            try { return JSON.parse(downloadRes.data); } catch { /* ignore */ }
          } else if (Buffer.isBuffer(downloadRes.data)) {
            try { return JSON.parse(downloadRes.data.toString('utf8')); } catch { /* ignore */ }
          } else if (typeof downloadRes.data === 'object') {
            return downloadRes.data;
          }
        }
      }
    } catch (apiErr) {
      core.warning(`API download of approval signal failed: ${apiErr.message}`);
    }
  }

  return null;
}

/**
 * validateApprovalSignal: Validates the integrity and freshness of the approval signal against the PR.
 */
export function validateApprovalSignal(signal, pr, core) {
  if (!signal || typeof signal !== 'object') {
    core.info('Missing or invalid approval signal.');
    return { isValid: false, reason: 'missing_signal' };
  }

  const reviewUser = signal.review_user ? String(signal.review_user).trim() : '';
  const reviewId = signal.review_id !== undefined && signal.review_id !== null ? String(signal.review_id).trim() : '';
  const reviewCommitId = signal.review_commit_id ? String(signal.review_commit_id).trim() : '';
  const signalHeadSha = signal.head_sha ? String(signal.head_sha).trim() : '';
  const reviewState = signal.review_state ? String(signal.review_state).trim().toUpperCase() : '';
  const baseRef = signal.base_ref ? String(signal.base_ref).trim().replace(/^refs\/heads\//, '') : '';
  const prNum = signal.pr_number !== undefined && signal.pr_number !== null ? Number.parseInt(signal.pr_number, 10) : null;

  if (!reviewUser || !reviewId || !reviewCommitId || !signalHeadSha) {
    core.info('Approval signal lacks required fields (review_user, review_id, review_commit_id, head_sha).');
    return { isValid: false, reason: 'incomplete_signal' };
  }

  if (reviewState !== 'APPROVED') {
    core.info(`Approval signal review_state is "${reviewState}", expected "APPROVED".`);
    return { isValid: false, reason: 'not_approved_state' };
  }

  if (baseRef && baseRef !== 'main') {
    core.info(`Approval signal base_ref is "${baseRef}", expected "main".`);
    return { isValid: false, reason: 'not_main_base' };
  }

  if (prNum && pr && pr.number && prNum !== pr.number) {
    core.info(`Approval signal PR #${prNum} does not match target PR #${pr.number}.`);
    return { isValid: false, reason: 'pr_mismatch' };
  }

  if (pr?.head?.sha && (reviewCommitId !== pr.head.sha || signalHeadSha !== pr.head.sha)) {
    core.info(`Approval signal commit (${reviewCommitId}) or head (${signalHeadSha}) does not match current PR head (${pr.head.sha}). Skipping stale signal.`);
    return { isValid: false, reason: 'stale_signal' };
  }

  return {
    isValid: true,
    signal: {
      pr_number: prNum,
      review_id: reviewId,
      review_user: reviewUser,
      review_commit_id: reviewCommitId,
      head_sha: signalHeadSha,
      review_state: reviewState,
      base_ref: baseRef,
    },
  };
}

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

  const headSha = parentRun?.head_sha || process?.env?.WORKFLOW_RUN_SHA || process?.env?.PR_HEAD_SHA;
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
        // Reject fallback PR matches with stale head SHA
        if (!p.head?.sha || p.head.sha !== headSha) {
          return false;
        }
        return true;
      });

      if (matched) {
        core.info(`Resolved PR #${matched.number} via commit search fallback.`);
        if (matched.number && typeof github?.rest?.pulls?.get === 'function') {
          try {
            const fullPr = await github.rest.pulls.get({
              owner,
              repo,
              pull_number: matched.number,
            });
            if (fullPr?.data) {
              if (fullPr.data.head?.sha && fullPr.data.head.sha !== headSha) {
                core.info(`PR #${matched.number} head (${fullPr.data.head.sha}) has changed since commit ${headSha}. Rejecting stale match.`);
                return null;
              }
              return fullPr.data;
            }
          } catch (err) {
            core.info(`Could not fetch full details for PR #${matched.number}, using commit-matched PR object: ${err.message}`);
          }
        }
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
 * Validates both filename and previous_filename to protect against workflow renaming.
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

  const workflowFiles = files.filter(f => {
    const filenameMatches = typeof f?.filename === 'string' && f.filename.startsWith('.github/');
    const prevMatches = typeof f?.previous_filename === 'string' && f.previous_filename.startsWith('.github/');
    return filenameMatches || prevMatches;
  });

  const matchedPaths = [];
  for (const f of workflowFiles) {
    if (typeof f?.filename === 'string' && f.filename.startsWith('.github/')) {
      matchedPaths.push(f.filename);
    }
    if (typeof f?.previous_filename === 'string' && f.previous_filename.startsWith('.github/')) {
      matchedPaths.push(f.previous_filename);
    }
  }

  return {
    modifiesWorkflows: workflowFiles.length > 0,
    workflowFiles: [...new Set(matchedPaths)],
  };
}

/**
 * evaluateReviews: Determines the latest review state for each reviewer.
 * Excludes author and authenticated bot reviews from approved human candidates.
 */
export function evaluateReviews(reviews, authorLogin, botLogin) {
  const DECISION_STATES = new Set(['APPROVED', 'CHANGES_REQUESTED', 'DISMISSED']);
  const normalizedAuthor = typeof authorLogin === 'string' ? authorLogin.toLowerCase() : null;
  const normalizedBot = typeof botLogin === 'string' ? botLogin.toLowerCase() : null;
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
      const userLower = r.user.login.toLowerCase();
      const isAuthor = normalizedAuthor && userLower === normalizedAuthor;
      const isBot = (normalizedBot && userLower === normalizedBot)
        || userLower.endsWith('[bot]')
        || r.user?.type === 'Bot';
      return isApproved && !isAuthor && !isBot;
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

  const triggerPrHeadSha = process?.env?.PR_HEAD_SHA || process?.env?.WORKFLOW_RUN_SHA;
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

  let botLogin = null;
  try {
    if (typeof prGithub?.rest?.users?.getAuthenticated === 'function') {
      const { data: authUser } = await prGithub.rest.users.getAuthenticated();
      if (authUser?.login && typeof authUser.login === 'string') {
        botLogin = authUser.login.trim();
      }
    }
  } catch (error) {
    core.setFailed(`Failed to resolve authenticated bot login: ${error.message}`);
    return false;
  }

  if (!botLogin) {
    core.setFailed('Could not resolve authenticated bot identity from GITHUB_MERGE_TOKEN. Failing closed.');
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
  const { hasChangesRequested, approvedReviewers, latestReviews } = evaluateReviews(reviews, authorLogin, botLogin);

  if (hasChangesRequested) {
    core.info(`PR #${prNumber} has one or more active CHANGES_REQUESTED reviews. Skipping auto-merge.`);
    core.setOutput('auto_merge_status', 'blocked_changes_requested');
    return false;
  }

  // Gate 4: Verify Reviewer Permissions and Bind to Triggering Signal
  const parentRun = context?.payload?.workflow_run;
  const signal = await resolveApprovalSignal({ github, context, core, process });

  let qualifyingReviewer = null;

  if (signal) {
    const { isValid, reason, signal: validatedSignal } = validateApprovalSignal(signal, pr, core);
    if (!isValid) {
      core.setFailed(`Triggering approval signal failed validation (${reason}). Failing closed.`);
      return false;
    }

    if (authorLogin && validatedSignal.review_user.toLowerCase() === authorLogin.toLowerCase()) {
      core.info(`Reviewer @${validatedSignal.review_user} is the author of PR #${prNumber}. Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'skipped_author');
      return false;
    }

    if (validatedSignal.review_user.toLowerCase() === botLogin.toLowerCase()) {
      core.info(`Reviewer @${validatedSignal.review_user} is the authenticated bot. Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'skipped_bot');
      return false;
    }

    const { isAuthorized, roleName, permission } = await checkReviewerPermission({
      github,
      owner,
      repo,
      username: validatedSignal.review_user,
      core,
    });

    if (!isAuthorized) {
      core.info(`Triggering reviewer @${validatedSignal.review_user} does not have triage, write, or admin permissions (role: ${roleName || permission}). Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'no_qualifying_permission');
      return false;
    }

    const triggerReview = latestReviews.find(
      r => r.user?.login?.toLowerCase() === validatedSignal.review_user.toLowerCase()
    );

    if (!triggerReview || (typeof triggerReview.state === 'string' && triggerReview.state.toUpperCase() !== 'APPROVED')) {
      core.info(`Triggering review by @${validatedSignal.review_user} is no longer in APPROVED state. Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'no_approvals');
      return false;
    }

    if (triggerReview.commit_id && pr.head?.sha && triggerReview.commit_id !== pr.head.sha) {
      core.info(`Approval by @${validatedSignal.review_user} was submitted on commit ${triggerReview.commit_id}, but current PR head is ${pr.head.sha}. Skipping stale review.`);
      core.setOutput('auto_merge_status', 'skipped_superseded');
      return false;
    }

    if (triggerReview.id && validatedSignal.review_id && String(triggerReview.id) !== String(validatedSignal.review_id)) {
      core.info(`Latest review ID (${triggerReview.id}) for @${validatedSignal.review_user} does not match triggering review ID (${validatedSignal.review_id}). Skipping superseded review.`);
      core.setOutput('auto_merge_status', 'skipped_superseded');
      return false;
    }

    qualifyingReviewer = {
      username: validatedSignal.review_user,
      role: roleName || permission,
      review: triggerReview,
    };
  } else if (parentRun?.id) {
    core.setFailed(`Workflow run ${parentRun.id} lacks valid 'auto-merge-approval-signal'. Failing closed.`);
    return false;
  } else {
    if (approvedReviewers.length === 0) {
      core.info(`PR #${prNumber} has no approved reviews from qualifying non-author, non-bot users. Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'no_approvals');
      return false;
    }

    for (const username of approvedReviewers) {
      const { isAuthorized, roleName, permission } = await checkReviewerPermission({
        github,
        owner,
        repo,
        username,
        core,
      });
      if (isAuthorized) {
        const qualifyingReview = latestReviews.find(
          r => r.user?.login?.toLowerCase() === username.toLowerCase()
        );
        qualifyingReviewer = {
          username,
          role: roleName || permission,
          review: qualifyingReview,
        };
        break;
      }
    }

    if (!qualifyingReviewer) {
      core.info(`PR #${prNumber} does not have an approved review from a user with triage, write, or admin permissions. Skipping auto-merge.`);
      core.setOutput('auto_merge_status', 'no_qualifying_permission');
      return false;
    }

    if (qualifyingReviewer.review?.commit_id && pr.head?.sha && qualifyingReviewer.review.commit_id !== pr.head.sha) {
      core.info(`Approval by @${qualifyingReviewer.username} was submitted on commit ${qualifyingReviewer.review.commit_id}, but current PR head is ${pr.head.sha}. Skipping stale review.`);
      core.setOutput('auto_merge_status', 'skipped_superseded');
      return false;
    }
  }

  core.info(`PR #${prNumber} qualified for auto-merge via sign-off from @${qualifyingReviewer.username} (role: ${qualifyingReviewer.role}).`);

  // Gate 5: Check if Bot already approved the current head commit
  const alreadyApprovedByBot = botLogin && reviews.some(
    r => r.user?.login?.toLowerCase() === botLogin.toLowerCase() && r.commit_id === pr.head.sha && r.state === 'APPROVED'
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
          const currentReviews = await (typeof github.paginate === 'function' && typeof github.rest?.pulls?.listReviews === 'function'
            ? github.paginate(github.rest.pulls.listReviews, {
              owner,
              repo,
              pull_number: prNumber,
            })
            : (github.rest?.pulls?.listReviews
              ? (await github.rest.pulls.listReviews({ owner, repo, pull_number: prNumber }))?.data || []
              : []));
          const { latestReviews: freshReviews } = evaluateReviews(currentReviews, authorLogin, botLogin);
          approvalConfirmed = freshReviews.some(
            r => r.user?.login?.toLowerCase() === botLogin.toLowerCase() && r.commit_id === pr.head.sha && r.state === 'APPROVED'
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

/**
 * runPreflight: Unprivileged preflight validation to determine PR eligibility before retrieving secrets.
 */
export async function runPreflight({ github, context, core, process = globalThis.process }) {
  const owner = context?.repo?.owner || process?.env?.GITHUB_REPOSITORY?.split('/')[0];
  const repo = context?.repo?.repo || process?.env?.GITHUB_REPOSITORY?.split('/')[1];

  if (!owner || !repo) {
    core.setFailed('Missing repository owner and repo in context.');
    core.setOutput('eligible', 'false');
    return false;
  }

  const parentRun = context?.payload?.workflow_run;

  const pr = await resolveTargetPr({ github, context, core, process });
  if (!pr) {
    core.info('Could not resolve target PR for preflight auto-merge. Skipping.');
    core.setOutput('eligible', 'false');
    return false;
  }

  const prNumber = pr.number;
  core.info(`Running preflight auto-merge validation for PR #${prNumber}...`);

  // Gate 1: State, Draft, and Target Branch Checks
  if (pr.state !== 'open') {
    core.info(`PR #${prNumber} is not open (state: "${pr.state}"). Skipping auto-merge.`);
    core.setOutput('eligible', 'false');
    return false;
  }

  if (pr.draft) {
    core.info(`PR #${prNumber} is a draft. Skipping auto-merge.`);
    core.setOutput('eligible', 'false');
    return false;
  }

  const baseRef = pr.base?.ref ? pr.base.ref.replace(/^refs\/heads\//, '') : '';
  if (baseRef !== 'main') {
    core.info(`PR #${prNumber} targets branch "${baseRef}", not "main". Skipping auto-merge.`);
    core.setOutput('eligible', 'false');
    return false;
  }

  const triggerPrHeadSha = process?.env?.PR_HEAD_SHA || process?.env?.WORKFLOW_RUN_SHA;
  if (triggerPrHeadSha && pr.head?.sha && triggerPrHeadSha !== pr.head.sha) {
    core.info(`PR #${prNumber} head commit (${triggerPrHeadSha}) is superseded by current head (${pr.head.sha}). Skipping stale run.`);
    core.setOutput('eligible', 'false');
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
    core.setOutput('eligible', 'false');
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
    core.setOutput('eligible', 'false');
    return false;
  }

  const authorLogin = pr.user?.login;
  const { hasChangesRequested, approvedReviewers, latestReviews } = evaluateReviews(reviews, authorLogin, null);

  if (hasChangesRequested) {
    core.info(`PR #${prNumber} has one or more active CHANGES_REQUESTED reviews. Skipping auto-merge.`);
    core.setOutput('eligible', 'false');
    return false;
  }

  // Gate 4: Download, read, and validate Approval Signal
  const signal = await resolveApprovalSignal({ github, context, core, process });
  if (parentRun?.id && !signal) {
    core.info(`Workflow run ${parentRun.id} lacks 'auto-merge-approval-signal' artifact (listener was skipped or event was not an approved review on main). Skipping auto-merge.`);
    core.setOutput('eligible', 'false');
    return false;
  }

  let qualifyingReviewer = null;

  if (signal) {
    const { isValid, reason, signal: validatedSignal } = validateApprovalSignal(signal, pr, core);
    if (!isValid) {
      core.info(`Approval signal failed validation (${reason}). Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    if (authorLogin && validatedSignal.review_user.toLowerCase() === authorLogin.toLowerCase()) {
      core.info(`Triggering reviewer @${validatedSignal.review_user} is the author of PR #${prNumber}. Self-approvals are not permitted. Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    const { isAuthorized, roleName, permission } = await checkReviewerPermission({
      github,
      owner,
      repo,
      username: validatedSignal.review_user,
      core,
    });

    if (!isAuthorized) {
      core.info(`Triggering reviewer @${validatedSignal.review_user} does not have an approved review from a user with triage, write, or admin permissions (role: ${roleName || permission}). Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    const triggerReview = latestReviews.find(
      r => r.user?.login?.toLowerCase() === validatedSignal.review_user.toLowerCase()
    );

    if (!triggerReview || (typeof triggerReview.state === 'string' && triggerReview.state.toUpperCase() !== 'APPROVED')) {
      core.info(`Review by triggering reviewer @${validatedSignal.review_user} is not in APPROVED state. Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    if (triggerReview.commit_id && pr.head?.sha && triggerReview.commit_id !== pr.head.sha) {
      core.info(`Approval by @${validatedSignal.review_user} was submitted on commit ${triggerReview.commit_id}, but current PR head is ${pr.head.sha}. Skipping stale review.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    if (triggerReview.id && validatedSignal.review_id && String(triggerReview.id) !== String(validatedSignal.review_id)) {
      core.info(`Latest review ID (${triggerReview.id}) for @${validatedSignal.review_user} does not match triggering review ID (${validatedSignal.review_id}). Skipping superseded review.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    qualifyingReviewer = {
      username: validatedSignal.review_user,
      role: roleName || permission,
      review: triggerReview,
      signal: validatedSignal,
    };
  } else {
    // Direct execution fallback when not triggered via workflow_run signal
    if (approvedReviewers.length === 0) {
      core.info(`PR #${prNumber} has no approved reviews from qualifying non-author, non-bot users. Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    for (const username of approvedReviewers) {
      const { isAuthorized, roleName, permission } = await checkReviewerPermission({
        github,
        owner,
        repo,
        username,
        core,
      });
      if (isAuthorized) {
        const qualifyingReview = latestReviews.find(
          r => r.user?.login?.toLowerCase() === username.toLowerCase()
        );
        qualifyingReviewer = {
          username,
          role: roleName || permission,
          review: qualifyingReview,
        };
        break;
      }
    }

    if (!qualifyingReviewer) {
      core.info(`PR #${prNumber} does not have an approved review from a user with triage, write, or admin permissions. Skipping auto-merge.`);
      core.setOutput('eligible', 'false');
      return false;
    }

    if (qualifyingReviewer.review?.commit_id && pr.head?.sha && qualifyingReviewer.review.commit_id !== pr.head.sha) {
      core.info(`Approval by @${qualifyingReviewer.username} was submitted on commit ${qualifyingReviewer.review.commit_id}, but current PR head is ${pr.head.sha}. Skipping stale review.`);
      core.setOutput('eligible', 'false');
      return false;
    }
  }

  core.info(`Preflight validation succeeded: PR #${prNumber} is eligible for auto-merge via sign-off from @${qualifyingReviewer.username} (role: ${qualifyingReviewer.role}).`);
  core.setOutput('eligible', 'true');
  core.setOutput('pr_number', String(pr.number));
  core.setOutput('pr_head_sha', pr.head?.sha || '');
  if (qualifyingReviewer.signal) {
    core.setOutput('review_id', String(qualifyingReviewer.signal.review_id));
    core.setOutput('review_user', qualifyingReviewer.signal.review_user);
    core.setOutput('review_commit_id', qualifyingReviewer.signal.review_commit_id);
  }
  return true;
}

export default async ({ github, context, core, process = globalThis.process, getOctokit }) => {
  const mode = process?.env?.SCRIPT_MODE;
  if (mode === 'preflight') {
    return await runPreflight({ github, context, core, process });
  }
  return await runAutoMerge({ github, context, core, process, getOctokit });
};
