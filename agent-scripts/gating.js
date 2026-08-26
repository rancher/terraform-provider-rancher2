import { execFileSync, execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import os from 'os';
import path from 'path';

// Calculate SHA-256 hash of a file's content
export function calculateFileHash(filePath) {
  try {
    const content = fs.readFileSync(filePath);
    return crypto.createHash('sha256').update(content).digest('hex');
  } catch (err) {
    console.error('🔒 Hook Debug: calculateFileHash failed:', err.message || err);
    return null;
  }
}

// Calculate active local diff hash securely (staged + unstaged combined, relative to main on feature branches)
export function calculateDiffHash() {
  try {
    const currentBranch = execSync('git branch --show-current', { stdio: ['ignore', 'pipe', 'ignore'] })
      .toString()
      .trim();
    const diffCmdArgs = currentBranch !== 'main' && currentBranch !== '' ? ['diff', 'main'] : ['diff', 'HEAD'];
    const diff = execFileSync('git', diffCmdArgs, { stdio: ['ignore', 'pipe', 'ignore'] }).toString();
    return crypto.createHash('sha256').update(diff).digest('hex');
  } catch (err) {
    console.error('🔒 Hook Debug: calculateDiffHash failed:', err.message || err);
    return null;
  }
}

// Automatically scans the Gemini tmp directories to find the latest active plan file
export function findLatestActivePlan(targetDir) {
  try {
    const activeSessions = fs.readdirSync(targetDir);
    const planFiles = [];

    for (const session of activeSessions) {
      const plansPath = path.join(targetDir, session, 'plans');
      if (fs.existsSync(plansPath) && fs.statSync(plansPath).isDirectory()) {
        const files = fs.readdirSync(plansPath);
        for (const file of files) {
          if (file.endsWith('.md')) {
            const filePath = path.join(plansPath, file);
            planFiles.push({
              path: filePath,
              mtime: fs.statSync(filePath).mtimeMs,
            });
          }
        }
      }
    }

    if (planFiles.length === 0) {
      return null;
    }

    planFiles.sort((a, b) => b.mtime - a.mtime);
    return planFiles[0].path;
  } catch (err) {
    console.error('🔒 Hook Debug: findLatestActivePlan failed:', err.message || err);
    return null;
  }
}

// Verify Gate 1: Plan Gate and return plan_hash
export function verifyPlanGate(targetDir) {
  const planApprovalFile = path.join(targetDir, 'plan-approval.json');
  const sigFile = path.join(targetDir, 'plan-approval.json.sig');
  const pubKeyFile = path.resolve(os.homedir(), '.gemini/ssh-key.pub');

  if (!fs.existsSync(planApprovalFile) || !fs.existsSync(sigFile) || !fs.existsSync(pubKeyFile)) {
    return null;
  }

  try {
    const allowedSignersFile = path.join(targetDir, 'allowed_signers');
    const pubKeyContent = fs.readFileSync(pubKeyFile, 'utf-8').trim();
    fs.writeFileSync(allowedSignersFile, `gemini ${pubKeyContent}`);

    execFileSync(
      'ssh-keygen',
      ['-Y', 'verify', '-f', allowedSignersFile, '-I', 'gemini', '-n', 'gemini', '-s', sigFile],
      {
        input: fs.readFileSync(planApprovalFile),
        stdio: ['pipe', 'ignore', 'ignore'],
      },
    );

    const content = JSON.parse(fs.readFileSync(planApprovalFile, 'utf-8'));

    if (content.status !== 'approved') {
      return null;
    }

    const activePlan = findLatestActivePlan(targetDir);
    if (!activePlan) {
      return null;
    }

    const currentPlanHash = calculateFileHash(activePlan);
    if (content.plan_hash !== currentPlanHash) {
      return null;
    }

    return content.plan_hash;
  } catch (err) {
    console.error('🔒 Hook Debug: verifyPlanGate failed:', err.message || err);
    return null;
  }
}

// Verify Gate 2: Test Gate
export function verifyTestGate(targetDir, expectedPlanHash, activeDiffHash) {
  const stateFile = path.join(targetDir, 'phase-state.json');

  if (!fs.existsSync(stateFile)) {
    return false;
  }

  try {
    const content = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
    if (content.tested_plan_hash !== expectedPlanHash) {
      return false;
    }

    if (content.tested_diff_hash !== activeDiffHash) {
      return false;
    }

    return true;
  } catch (err) {
    console.error('🔒 Hook Debug: verifyTestGate failed:', err.message || err);
    return false;
  }
}

// Verify Gate 3: Review Gate
export function verifyReviewGate(targetDir, expectedDiffHash, expectedPlanHash) {
  const reviewApprovalFile = path.join(targetDir, 'review-approval.json');

  if (!fs.existsSync(reviewApprovalFile)) {
    return false;
  }

  try {
    const content = JSON.parse(fs.readFileSync(reviewApprovalFile, 'utf-8'));
    if (content.status !== 'approved') {
      return false;
    }

    if (content.plan_hash !== expectedPlanHash) {
      return false;
    }
    if (content.diff_hash !== expectedDiffHash) {
      return false;
    }

    return true;
  } catch (err) {
    console.error('🔒 Hook Debug: verifyReviewGate failed:', err.message || err);
    return false;
  }
}

// Check and actively revoke stale signatures (Gate 2) if the diff hash has changed
export function checkAndRevokeStaleGates(targetDir, activeDiffHash, expectedPlanHash) {
  const reviewApprovalFile = path.join(targetDir, 'review-approval.json');

  let hasRevoked = false;

  // Check review approval
  if (fs.existsSync(reviewApprovalFile)) {
    try {
      const content = JSON.parse(fs.readFileSync(reviewApprovalFile, 'utf-8'));
      if (
        activeDiffHash &&
        expectedPlanHash &&
        (content.diff_hash !== activeDiffHash || content.plan_hash !== expectedPlanHash)
      ) {
        fs.unlinkSync(reviewApprovalFile);
        console.error(
          '❌ Active Gate Revocation: Stale review signature deleted because workspace changes were modified since your last review!',
        );
        hasRevoked = true;
      }
    } catch (err) {
      console.warn(`🔒 Hook Warning: Failed to parse review approval JSON: ${err.message}`);
      // If unparsable, delete it
      try {
        fs.unlinkSync(reviewApprovalFile);
      /* eslint-disable-next-line no-shadow */
      } catch (err) {
        console.error(`🔒 Hook Warning: Failed to delete unparsable review approval: ${err.message}`);
      }
      hasRevoked = true;
    }
  }

  // Also reset tested_diff_hash in phase-state.json on change to enforce re-testing
  const stateFile = path.join(targetDir, 'phase-state.json');
  if (fs.existsSync(stateFile)) {
    try {
      const content = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
      if (
        activeDiffHash &&
        expectedPlanHash &&
        (content.tested_diff_hash !== activeDiffHash || content.tested_plan_hash !== expectedPlanHash)
      ) {
        content.tested_diff_hash = '';
        content.tested_plan_hash = '';
        fs.writeFileSync(stateFile, JSON.stringify(content, null, 2));
      }
    } catch (err) {
      console.error(`🔒 Hook Warning: Failed to parse phase-state JSON in revoke: ${err.message}`);
    }
  }

  return hasRevoked;
}
