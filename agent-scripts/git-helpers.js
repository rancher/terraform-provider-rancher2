import { execFileSync, execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import os from 'os';
import path from 'path';
import { findLatestActivePlan, verifyPlanGate, verifyReviewGate, verifyTestGate } from './gating.js';
import { validatePlanContent } from './planning.js';

export function calculateSha256(data) {
  return crypto.createHash('sha256').update(data).digest('hex');
}

export function getFileOwnerUid(filePath) {
  try {
    return fs.statSync(filePath).uid;
  } catch (err) {
    console.error(err.message || err);
    return null;
  }
}

export function checkDefunctBranch(branch, cwd) {
  if (branch === 'main') {
    return;
  }
  try {
    const prStatus = execFileSync(
      'gh',
      ['pr', 'view', branch, '--json', 'state,number', '--template', '{{.state}} {{.number}}'],
      {
        cwd: cwd || process.cwd(),
        env: { ...process.env, GITHUB_TOKEN: '' },
        stdio: ['ignore', 'pipe', 'ignore'],
      },
    )
      .toString()
      .trim();
    if (prStatus) {
      const [state, number] = prStatus.split(' ');
      if (state === 'MERGED') {
        throw new Error(`The current branch '${branch}' already has a merged Pull Request (#${number}) on GitHub.`);
      }
    }
  } catch (err) {
    if (err.message.includes('already has a merged Pull Request')) {
      throw err;
    }
    console.error(err.message || err);
    // If gh command fails safely skip defunct check
  }
}

export function verifyRemoteAncestry(branch, remoteName = 'origin', cwd) {
  console.log('Verifying local branch ancestry is fully up-to-date with remote fork...');
  try {
    execFileSync('git', ['fetch', '-q', remoteName, branch], {
      cwd: cwd || process.cwd(),
      stdio: 'ignore',
    });
  } catch (err) {
    console.error(`🔒 Hook Info: Fetch skipped: ${err.message || err}`);
    console.log(
      `--> [FETCH SKIPPED] Tracking reference on remote '${remoteName}/${branch}' does not exist yet. Safe to proceed.`,
    );
    return;
  }

  try {
    const localSha = execFileSync('git', ['rev-parse', 'HEAD'], { cwd: cwd || process.cwd() })
      .toString()
      .trim();
    const remoteSha = execFileSync('git', ['rev-parse', `${remoteName}/${branch}`], { cwd: cwd || process.cwd() })
      .toString()
      .trim();

    if (!remoteSha) {
      console.log('--> [RESOLVE SKIPPED] No remote tracking branch SHA found. Safe to proceed.');
      return;
    }

    if (localSha === remoteSha) {
      console.log('✅ Local branch is identical to remote fork tracking reference.');
      return;
    }

    try {
      execFileSync('git', ['merge-base', '--is-ancestor', remoteSha, localSha], { cwd: cwd || process.cwd() });
      console.log('✅ Local branch contains all remote tracking changes (Fast-Forward ancestry confirmed).');
      return;
    } catch (err) {
      console.error(err.message || err);
      throw new Error(
        `CRITICAL FAILURE: LOCAL BRANCH IS OUT-OF-SYNC WITH REMOTE FORK! Remote tracking reference '${remoteName}/${branch}' (${remoteSha}) has changes that are not present locally.`,
        { cause: err },
      );
    }
  } catch (err) {
    if (err.message.includes('LOCAL BRANCH IS OUT-OF-SYNC')) {
      throw err;
    }
    console.error(err.message || err);
  }
}

export function verifyStagingLimits(override, cwd) {
  let maxAllowed = 5;
  if (override) {
    const parsed = parseInt(override, 10);
    if (isNaN(parsed) || parsed < 0) {
      throw new Error(`COMMIT_LIMIT_OVERRIDE must be a positive integer, got: '${override}'`);
    }
    maxAllowed = parsed;
    console.log(`--> [OVERRIDE] Using custom staged file limit from COMMIT_LIMIT_OVERRIDE: ${maxAllowed}`);
  }

  const stagedOutput = execSync('git diff --cached --name-only', { cwd: cwd || process.cwd() })
    .toString()
    .trim();
  const stagedCount = stagedOutput ? stagedOutput.split('\n').length : 0;

  if (stagedCount === 0) {
    if (
      fs.existsSync(path.join(cwd || process.cwd(), '.git/MERGE_HEAD')) ||
      fs.existsSync(path.join(cwd || process.cwd(), '.git/CHERRY_PICK_HEAD')) ||
      fs.existsSync(path.join(cwd || process.cwd(), '.git/REBASE_HEAD'))
    ) {
      console.log(
        '--> [MERGE STATE] Active merge/rebase/cherry-pick in progress. Allowing 0 staged files to create merge commit.',
      );
      return;
    }
    throw new Error('No changes are currently staged for commit. Please stage your changes first.');
  }

  if (stagedCount > maxAllowed) {
    throw new Error(
      `Committing too much code at once is prohibited (${stagedCount} files staged; max allowed is ${maxAllowed}).`,
    );
  }
}

export function verifyProactiveReview(agentStateDir, cwd) {
  let targetDir = agentStateDir;
  if (!targetDir) {
    const topLevel = execSync('git rev-parse --show-toplevel', { cwd: cwd || process.cwd() })
      .toString()
      .trim();
    const repoName = path.basename(topLevel);
    targetDir = path.resolve(os.homedir(), '.gemini/tmp', repoName);
  }
  const reviewFile = path.join(targetDir, 'review-approval.json');

  console.log('Verifying proactive review approval status...');

  if (!fs.existsSync(reviewFile)) {
    throw new Error('Proactive review approval file not found!');
  }

  if (fs.lstatSync(reviewFile).isSymbolicLink()) {
    throw new Error('Proactive review approval file is a symbolic link (Prohibited).');
  }

  const fileUid = getFileOwnerUid(reviewFile);
  if (fileUid === null) {
    throw new Error('Could not determine owner UID for proactive review approval file.');
  }

  const currentUid = os.userInfo().uid;
  if (fileUid !== currentUid) {
    throw new Error(
      `Proactive review approval file is not owned by the current user (UID: ${currentUid}, Owner: ${fileUid}).`,
    );
  }

  const activeDiff = execSync('git diff HEAD', { cwd: cwd || process.cwd() }).toString();
  const activeHash = calculateSha256(activeDiff);

  try {
    const content = JSON.parse(fs.readFileSync(reviewFile, 'utf-8'));
    if (content.status !== 'approved') {
      throw new Error(`Proactive review approval status is '${content.status}' (not approved).`);
    }
    if (content.diff_hash !== activeHash) {
      throw new Error(
        `Local changes have been modified since your last proactive review! Approved: ${content.diff_hash}, Current: ${activeHash}`,
      );
    }
  } catch (err) {
    throw new Error(`Failed to parse proactive review approval file: ${err.message}`, { cause: err });
  }

  console.log(`✅ Proactive review approval verified! (SHA-256 Hash: ${activeHash})`);
}

export function executeCommit(commitMsg, branch, cwd) {
  console.log('Staging changes...');
  execFileSync('git', ['add', '-A'], { cwd: cwd || process.cwd() });

  console.log(`Creating conventional signed commit on branch '${branch}'...`);
  try {
    execFileSync('git', ['commit', '-S', '-s', '-m', commitMsg], {
      cwd: cwd || process.cwd(),
      stdio: 'inherit',
    });
    console.log('✅ Conventional GPG/SSH-signed commit successfully created!');
  } catch (err) {
    console.error(err);
    throw new Error('GPG/SSH COMMIT SIGNATURE FAILURE! Commit signature operation failed or was cancelled.', { cause: err });
  }
}

export function verifyPushSafety(remoteName, cwd) {
  try {
    const url = execFileSync('git', ['remote', 'get-url', remoteName], {
      cwd: cwd || process.cwd(),
      stdio: ['ignore', 'pipe', 'ignore'],
    })
      .toString()
      .trim();
    if (/[/:](rancher|rancherlabs)[/:]/.test(url)) {
      throw new Error(
        `CRITICAL SECURITY ERROR: Unsafe push prevented! Remote '${remoteName}' points to Rancher-owned repository: ${url}`,
      );
    }
  } catch (err) {
    if (err.message.includes('Unsafe push prevented')) {
      throw err;
    }
    throw new Error(`Remote '${remoteName}' has no configured URL.`, { cause: err });
  }
}

export function executePush(remoteName, branch, forcePush = false, cwd) {
  verifyPushSafety(remoteName, cwd);

  if (forcePush) {
    console.log(`Safely force-pushing branch '${branch}' to '${remoteName}' with lease...`);
    try {
      execFileSync('git', ['push', '-u', remoteName, branch, '--force-with-lease'], {
        cwd: cwd || process.cwd(),
        stdio: 'inherit',
      });
    } catch (err) {
      console.error(err.message || err);
      throw new Error('Remote force-push with lease failed.', { cause: err });
    }
  } else {
    console.log(`Pushing branch '${branch}' to '${remoteName}'...`);
    try {
      execFileSync('git', ['push', '-u', remoteName, branch], {
        cwd: cwd || process.cwd(),
        stdio: 'inherit',
      });
    } catch (err) {
      console.error(err.message || err);
      throw new Error('Remote push failed.', { cause: err });
    }
  }

  console.log(`✅ Changes successfully pushed to remote '${remoteName}/${branch}'!`);
}

// ==============================================================================
// GITHUB PR HELPERS (merged from pr-helper.js)
// ==============================================================================

export function getForkOwner(cwd) {
  try {
    const originUrl = execFileSync('git', ['remote', 'get-url', 'origin'], {
      cwd: cwd || process.cwd(),
      stdio: ['ignore', 'pipe', 'ignore'],
    })
      .toString()
      .trim();

    const match = originUrl.match(/github\.com[:/]([^/]+)\//);
    if (match) {
      return match[1];
    }
    throw new Error(`Could not parse fork owner from origin URL: ${originUrl}`);
  } catch (err) {
    throw new Error(`Failed to retrieve configured origin remote URL: ${err.message}`, { cause: err });
  }
}

export function createPullRequest(title, body, base = 'main', draftFlag = '', cwd) {
  const branch = execFileSync('git', ['branch', '--show-current'], { cwd: cwd || process.cwd() })
    .toString()
    .trim();
  const forkOwner = getForkOwner(cwd);

  try {
    const existingPr = execFileSync('gh', ['pr', 'list', '--head', branch, '--json', 'url', '--jq', '.[0].url'], {
      cwd: cwd || process.cwd(),
      env: { ...process.env, GITHUB_TOKEN: '' },
      stdio: ['ignore', 'pipe', 'ignore'],
    })
      .toString()
      .trim();

    if (existingPr) {
      console.log(`✅ Pull Request already exists for branch '${branch}': ${existingPr}`);
      return;
    }
  } catch (err) {
    console.error(err.message || err);
    // Safely ignore if gh command fails
  }

  let upstreamRepo;
  try {
    const originUrl = execFileSync('git', ['remote', 'get-url', 'origin'], {
      cwd: cwd || process.cwd(),
      stdio: ['ignore', 'pipe', 'ignore'],
    })
      .toString()
      .trim();

    const match = originUrl.match(/github\.com[:/]([^/]+)\/([^/]+)(?:\.git)?/);
    if (match) {
      const repoNameClean = match[2].replace(/\.git$/, '');
      upstreamRepo = `rancher/${repoNameClean}`;
    } else {
      const topLevel = execFileSync('git', ['rev-parse', '--show-toplevel'], { cwd: cwd || process.cwd() })
        .toString()
        .trim();
      const repoName = path.basename(topLevel);
      upstreamRepo = `rancher/${repoName}`;
    }
  } catch (err) {
    console.error(err.message || err);
    const repoName = path.basename(cwd || process.cwd()) || 'generic-repo';
    upstreamRepo = `rancher/${repoName}`;
  }

  console.log(`Creating Pull Request for branch '${branch}' on fork '${forkOwner}' to upstream '${upstreamRepo}'...`);

  try {
    execFileSync('gh', ['repo', 'set-default', upstreamRepo], {
      cwd: cwd || process.cwd(),
      env: { ...process.env, GITHUB_TOKEN: '' },
      stdio: 'ignore',
    });

    const prArgs = [
      'pr',
      'create',
      '--repo',
      upstreamRepo,
      '--base',
      base,
      '--head',
      `${forkOwner}:${branch}`,
      '--title',
      title,
      '--body',
      body,
    ];
    if (draftFlag) {
      prArgs.push('--draft');
    }

    execFileSync('gh', prArgs, {
      cwd: cwd || process.cwd(),
      stdio: 'inherit',
      env: { ...process.env, GITHUB_TOKEN: '' },
    });
  } catch (err) {
    throw new Error(`Failed to create Pull Request: ${err.message}`, { cause: err });
  }
}

export function graduatePullRequest(target, cwd) {
  const branch = execFileSync('git', ['branch', '--show-current'], { cwd: cwd || process.cwd() })
    .toString()
    .trim();

  try {
    if (!target) {
      console.log(`Graduating draft pull request for the current branch '${branch}'...`);
      execFileSync('gh', ['pr', 'ready'], {
        cwd: cwd || process.cwd(),
        env: { ...process.env, GITHUB_TOKEN: '' },
        stdio: 'inherit',
      });
    } else {
      console.log(`Graduating draft pull request for target '${target}'...`);
      execFileSync('gh', ['pr', 'ready', target], {
        cwd: cwd || process.cwd(),
        env: { ...process.env, GITHUB_TOKEN: '' },
        stdio: 'inherit',
      });
    }
    console.log('✅ Pull Request successfully graduated to ready for review!');
  } catch (err) {
    throw new Error(`Failed to graduate Pull Request: ${err.message}`, { cause: err });
  }
}

// ==============================================================================
// PHASE STATE MACHINE HELPERS (merged from phase-manager.js)
// ==============================================================================

export function getPhaseState(stateFile) {
  if (!fs.existsSync(stateFile)) {
    return { currentPhase: 'research' };
  }
  try {
    return JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
  } catch (err) {
    console.error(err.message || err);
    return { currentPhase: 'research' };
  }
}

export function savePhaseState(state, stateFile) {
  fs.mkdirSync(path.dirname(stateFile), { recursive: true });
  fs.writeFileSync(stateFile, JSON.stringify(state, null, 2));
}

export function runPhaseManager(args, cwd) {
  const topLevel = execSync('git rev-parse --show-toplevel', { cwd: cwd || process.cwd() })
    .toString()
    .trim();
  const repoName = path.basename(topLevel);
  const homeDir = os.homedir();
  const targetDir = path.resolve(homeDir, '.gemini/tmp', repoName);
  const stateFile = path.join(targetDir, 'phase-state.json');

  if (args.length === 0) {
    const { currentPhase } = getPhaseState(stateFile);
    console.log(`📍 Current development phase: ${currentPhase.toUpperCase()}`);
    return;
  }

  const targetPhase = args[0].toLowerCase();
  const allowedPhases = ['research', 'plan', 'implement', 'test', 'review', 'commit'];
  if (!allowedPhases.includes(targetPhase)) {
    throw new Error(`Invalid phase: ${targetPhase}. Allowed: ${allowedPhases.join(', ')}`);
  }

  const { currentPhase } = getPhaseState(stateFile);
  console.log(`🔄 Attempting transition: ${currentPhase.toUpperCase()} -> ${targetPhase.toUpperCase()}`);

  if (targetPhase === currentPhase) {
    console.log(`ℹ️ Already in phase: ${targetPhase.toUpperCase()}`);
    return;
  }

  // --- RESEARCH -> PLAN ---
  if (currentPhase === 'research' && targetPhase === 'plan') {
    const isDirty = execFileSync('git', ['status', '--porcelain'], { cwd: cwd || process.cwd() })
      .toString()
      .trim();
    if (isDirty) {
      throw new Error(
        'Transition Blocked: You have uncommitted changes in your workspace. Please commit or stash them before transitioning phases to prevent data loss.',
      );
    }
    console.log('🧹 Exiting RESEARCH: Workspace is clean. Proceeding to PLAN phase.');
  }

  // --- PLAN -> IMPLEMENT ---
  if (targetPhase === 'implement') {
    const planHash = verifyPlanGate(targetDir);
    if (!planHash) {
      throw new Error('Transition Blocked: Missing or invalid plan cryptographic approval (Gate 1).');
    }

    const isDirty = execFileSync('git', ['status', '--porcelain'], { cwd: cwd || process.cwd() })
      .toString()
      .trim();
    if (isDirty) {
      throw new Error(
        'Transition Blocked: You have uncommitted changes in your workspace. Please commit or stash them before transitioning phases to prevent data loss.',
      );
    }
    console.log('🧹 Exiting PLAN: Workspace is clean. Proceeding to IMPLEMENT phase.');

    const activePlan = findLatestActivePlan(targetDir);
    if (activePlan) {
      // Programmatic verification of Plan structural content
      const validation = validatePlanContent(activePlan);
      if (!validation.valid) {
        throw new Error(`Transition Blocked: Active plan has invalid structure! ${validation.errors.join(', ')}`);
      }

      const planName = path
        .basename(activePlan, '.md')
        .toLowerCase()
        .replace(/[^a-z0-9_-]/g, '-');
      const branchName = `feature/${planName}`;

      const currentBranch = execFileSync('git', ['branch', '--show-current'], { cwd: cwd || process.cwd() })
        .toString()
        .trim();
      if (currentBranch !== branchName) {
        console.log(`🌿 Switching to feature branch: ${branchName}`);
        try {
          const branchesOutput = execFileSync('git', ['branch', '--list', '--format=%(refname:short)'], {
            cwd: cwd || process.cwd(),
          }).toString();
          const branches = branchesOutput
            .split('\n')
            .map((b) => b.trim())
            .filter(Boolean);
          if (branches.includes(branchName)) {
            execFileSync('git', ['switch', branchName], { cwd: cwd || process.cwd() });
          } else {
            execFileSync('git', ['checkout', '-b', branchName], { cwd: cwd || process.cwd() });
          }
        } catch (err) {
          throw new Error(`Failed to switch to branch ${branchName}: ${err.message}`, { cause: err });
        }
      }
    }
  }

  // Validation checks for other phase transitions
  if (targetPhase === 'test') {
    const planHash = verifyPlanGate(targetDir);
    if (!planHash) {
      throw new Error('Transition Blocked: Missing or invalid plan cryptographic approval (Gate 1).');
    }
  }

  if (targetPhase === 'review') {
    const planHash = verifyPlanGate(targetDir);
    if (!planHash) {
      throw new Error('Transition Blocked: Missing or invalid plan cryptographic approval (Gate 1).');
    }

    const diffHash = execSync('git diff HEAD', { cwd: cwd || process.cwd() }).toString();
    const diffHashClean = calculateSha256(diffHash);
    const testPassed = verifyTestGate(targetDir, planHash, diffHashClean);
    if (!testPassed) {
      throw new Error('Transition Blocked: Missing or invalid test approval (Gate 2). Run and pass all tests first.');
    }
  }

  if (targetPhase === 'commit') {
    const planHash = verifyPlanGate(targetDir);
    if (!planHash) {
      throw new Error('Transition Blocked: Missing or invalid plan cryptographic approval (Gate 1).');
    }

    const diffHash = execSync('git diff HEAD', { cwd: cwd || process.cwd() }).toString();
    const diffHashClean = calculateSha256(diffHash);
    const testPassed = verifyTestGate(targetDir, planHash, diffHashClean);
    if (!testPassed) {
      throw new Error('Transition Blocked: Missing or invalid test approval (Gate 2).');
    }

    const reviewPassed = verifyReviewGate(targetDir, diffHashClean, planHash);
    if (!reviewPassed) {
      throw new Error(
        'Transition Blocked: Missing or invalid review approval (Gate 3). Get code review sign-off first.',
      );
    }
  }

  // Authorized
  savePhaseState({ currentPhase: targetPhase }, stateFile);
  console.log(`✅ Transitioned successfully to phase: ${targetPhase.toUpperCase()}`);
}

// ==============================================================================
// RUNNERS (executable CLI entrypoints)
// ==============================================================================

export function runCommitPushHelper(args, cwd) {
  let commitMsg = '';
  let forcePush = false;

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '-h' || arg === '--help') {
      console.log(`Usage: commit-push-helper.js [options] -m "COMMIT_MESSAGE"

Programmatically commit and push local changes with GPG/SSH signature, sign-off, and fork synchronization.

Options:
  -h, --help            Show this message and exit.
  -m MESSAGE            The conventional commit message (Required).
  -f, --force           Bypass remote ancestry check and perform safe force-push with lease.`);
      return;
    } else if (arg === '-f' || arg === '--force') {
      forcePush = true;
    } else if (arg === '-m') {
      if (i + 1 >= args.length || args[i + 1].startsWith('-')) {
        throw new Error('Error: -m option requires a non-empty commit message argument.');
      }
      commitMsg = args[++i];
    } else {
      throw new Error(`Error: Unknown argument '${arg}'`);
    }
  }

  if (!commitMsg) {
    throw new Error('Error: Commit message is required. Specify using -m "message".');
  }

  const activeBranch = execSync('git branch --show-current', { cwd: cwd || process.cwd() })
    .toString()
    .trim();

  // 1. Defunct branch check
  checkDefunctBranch(activeBranch, cwd);

  // 2. Proactive review validation
  verifyProactiveReview(null, cwd);

  // 3. Stage & limit checks
  execSync('git add -A', { cwd: cwd || process.cwd() });
  verifyStagingLimits(process.env.COMMIT_LIMIT_OVERRIDE, cwd);

  // 4. Remote ancestry validation
  if (!forcePush) {
    verifyRemoteAncestry(activeBranch, 'origin', cwd);
  }

  // 5. Commit
  executeCommit(commitMsg, activeBranch, cwd);

  // 6. Push
  executePush('origin', activeBranch, forcePush, cwd);
}

export function runCreatePrHelper(args, cwd) {
  let title = '';
  let body = '';
  let base = 'main';
  let draftFlag = '';
  let action = 'create';
  let readyTarget = '';

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '-h' || arg === '--help') {
      console.log(`Usage: create-pr-helper.js [options]

Safely creates or graduates a pull request from your local fork branch to the upstream repository.

Options:
  --title TITLE        The title of the pull request (Required for creation).
  --body BODY          The markdown description body of the pull request (Required for creation).
  --base BASE          The target upstream branch (default: main).
  --draft              Create the pull request as a draft.
  --ready [TARGET]     Graduate a draft pull request to ready-for-review (accepts PR number, branch, or URL; defaults to current branch).`);
      return;
    } else if (arg === '--title') {
      if (i + 1 >= args.length || args[i + 1].startsWith('-')) {
        throw new Error('Error: --title requires an argument.');
      }
      title = args[++i];
    } else if (arg === '--body') {
      if (i + 1 >= args.length || args[i + 1].startsWith('-')) {
        throw new Error('Error: --body requires an argument.');
      }
      body = args[++i];
    } else if (arg === '--base') {
      if (i + 1 >= args.length || args[i + 1].startsWith('-')) {
        throw new Error('Error: --base requires an argument.');
      }
      base = args[++i];
    } else if (arg === '--draft') {
      draftFlag = '--draft';
    } else if (arg === '--ready') {
      action = 'ready';
      if (i + 1 < args.length && !args[i + 1].startsWith('-')) {
        readyTarget = args[++i];
      }
    } else {
      throw new Error(`Unknown parameter: ${arg}`);
    }
  }

  if (action === 'ready') {
    graduatePullRequest(readyTarget, cwd);
    return;
  }

  if (!title) {
    try {
      title = execSync('git log -1 --pretty=%s', { cwd: cwd || process.cwd() })
        .toString()
        .trim();
    } catch (err) {
      console.error(`🔒 Hook Warning: Failed to retrieve commit title: ${err.message || err}`);
      title = 'chore: automated development pull request';
    }
  }

  if (!body) {
    try {
      body = execSync('git log -1 --pretty=%b', { cwd: cwd || process.cwd() })
        .toString()
        .trim();
    } catch (err) {
      console.error(`🔒 Hook Warning: Failed to retrieve commit body: ${err.message || err}`);
      body = 'Automated Draft Pull Request generated by the Agentic Framework.';
    }
  }

  createPullRequest(title, body, base, draftFlag, cwd);
}
