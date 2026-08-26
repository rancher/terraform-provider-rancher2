#!/usr/bin/env node
//
// Skill: user-approval.js
// Description: Securely prompts for developer approval and verifies user-approval status using diff-hash tieing.
//              Conforms to repository standards and rules.
// Usage:
//   Write Mode:   node .gemini/skills/user-approval.js --approve
//   Verify Mode:  node .gemini/skills/user-approval.js --verify
//   Prompt Mode:  node .gemini/skills/user-approval.js [message] [defaultOption=N]

import fs from 'fs';
import path from 'path';
import os from 'os';
import crypto from 'crypto';
import { execSync } from 'child_process';
import { Buffer } from 'buffer';

let repoName;
try {
  const topLevel = execSync('git rev-parse --show-toplevel', { stdio: ['ignore', 'pipe', 'ignore'] })
    .toString()
    .trim();
  repoName = path.basename(topLevel);
} catch (err) {
  console.error(err.message || err);
  repoName = path.basename(process.cwd()) || 'generic-repo';
}

const TARGET_DIR = path.resolve(os.homedir(), '.gemini/tmp', repoName);
const APPROVAL_FILE = path.join(TARGET_DIR, 'user-approval.json');
const PLAN_APPROVAL_FILE = path.join(TARGET_DIR, 'plan-approval.json');
const PLAN_CHALLENGE_FILE = path.join(TARGET_DIR, 'plan-approval.challenge');

// Helper to calculate active local diff hash securely using SHA-256
function calculateSHA256() {
  try {
    const diff = execSync('git diff HEAD', { stdio: ['ignore', 'pipe', 'ignore'] }).toString();
    return crypto.createHash('sha256').update(diff).digest('hex');
  } catch (err) {
    console.error(err.message || err);
    return null;
  }
}

function calculateFileHash(filePath) {
  try {
    const content = fs.readFileSync(filePath);
    return crypto.createHash('sha256').update(content).digest('hex');
  } catch (err) {
    console.error(err.message || err);
    return null;
  }
}

function findLatestActivePlan() {
  try {
    const activeSessions = fs.readdirSync(TARGET_DIR);
    const planFiles = [];

    for (const session of activeSessions) {
      const plansPath = path.join(TARGET_DIR, session, 'plans');
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
    console.error(err.message || err);
    return null;
  }
}

function verifyPlan() {
  console.log('Verifying developer plan approval status...');

  if (!fs.existsSync(PLAN_APPROVAL_FILE) || !fs.existsSync(PLAN_CHALLENGE_FILE)) {
    console.error('Error: Plan approval signature not found!');
    process.exit(1);
  }

  const stat = fs.lstatSync(PLAN_APPROVAL_FILE);
  if (stat.isSymbolicLink()) {
    console.error('Error: Plan approval file is a symbolic link (Prohibited).');
    process.exit(1);
  }

  if (process.getuid) {
    const fileUid = fs.statSync(PLAN_APPROVAL_FILE).uid;
    const currentUid = process.getuid();
    if (fileUid !== currentUid) {
      console.error('Error: Plan approval file is not owned by the current user.');
      process.exit(1);
    }
  }

  try {
    const content = JSON.parse(fs.readFileSync(PLAN_APPROVAL_FILE, 'utf-8'));
    const challenge = JSON.parse(fs.readFileSync(PLAN_CHALLENGE_FILE, 'utf-8'));

    if (content.status !== 'approved') {
      console.error('Error: Plan approval status is not approved.');
      process.exit(1);
    }

    const token = content.challenge_token;
    if (!token) {
      console.error('Error: Plan approval missing challenge token.');
      process.exit(1);
    }

    const calculatedHash = crypto.createHash('sha256').update(token).digest('hex');
    if (calculatedHash !== challenge.challenge_hash) {
      console.error('Error: Plan approval challenge validation failed.');
      process.exit(1);
    }

    const activePlan = findLatestActivePlan();
    if (!activePlan) {
      console.error('Error: No active plan found in session memory.');
      process.exit(1);
    }

    const currentPlanHash = calculateFileHash(activePlan);
    if (content.plan_hash !== currentPlanHash) {
      console.error('Error: The active plan has been modified since approval!');
      process.exit(1);
    }

    console.log(`✅ Plan approval verified! (Plan Hash: ${currentPlanHash})`);
    process.exit(0);
  } catch (err) {
    console.error('Error: Failed to verify plan approval:', err.message);
    process.exit(1);
  }
}

// ------------------------------------------------------------------------------
// VERIFY OPERATION
// ------------------------------------------------------------------------------
function verifyApproval() {
  console.log('Verifying developer manual IDE review approval status...');

  if (!fs.existsSync(APPROVAL_FILE)) {
    console.error('Error: Developer manual IDE review approval not found!');
    console.error(
      "       In accordance with Gate 2 (IDE & Commit Gate) of 'docs/development/AgenticFramework/DevelopmentProcess.md',",
    );
    console.error('       you MUST request developer approval first: @user-approval');
    process.exit(1);
  }

  // Reject symbolic links to prevent symlink bypasses
  const stat = fs.lstatSync(APPROVAL_FILE);
  if (stat.isSymbolicLink()) {
    console.error('Error: Developer approval file is a symbolic link (Prohibited).');
    process.exit(1);
  }

  // Verify file ownership matches current user UID natively
  if (process.getuid) {
    const fileUid = fs.statSync(APPROVAL_FILE).uid;
    const currentUid = process.getuid();
    if (fileUid !== currentUid) {
      console.error(
        `Error: Developer approval file is not owned by the current user (UID: ${currentUid}, Owner: ${fileUid}).`,
      );
      process.exit(1);
    }
  }

  try {
    const content = JSON.parse(fs.readFileSync(APPROVAL_FILE, 'utf-8'));
    if (content.status !== 'approved') {
      console.error(`Error: Developer approval status is '${content.status}' (not approved).`);
      process.exit(1);
    }

    const activeHash = calculateSHA256();
    if (content.diff_hash !== activeHash) {
      console.error('Error: Local changes have been modified since your last manual developer approval!');
      console.error(`       Approved SHA-256 hash: ${content.diff_hash}`);
      console.error(`       Current active SHA-256 hash: ${activeHash}`);
      console.error('       Please request developer approval again on your latest changes.');
      process.exit(1);
    }

    console.log(`✅ Developer visual IDE review approval verified! (SHA-256 Hash: ${activeHash})`);
    process.exit(0);
  } catch (err) {
    console.error('Error: Failed to parse developer approval file:', err.message);
    process.exit(1);
  }
}

// ------------------------------------------------------------------------------
// APPROVE OPERATION
// ------------------------------------------------------------------------------
function writeApproval() {
  const activeHash = calculateSHA256();
  if (!activeHash) {
    console.error('Error: Failed to calculate active diff hash.');
    process.exit(1);
  }

  const approvalData = {
    status: 'approved',
    diff_hash: activeHash,
    timestamp: new Date().toISOString(),
  };

  try {
    fs.mkdirSync(TARGET_DIR, { recursive: true });

    // Delete any existing file or symlink to prevent symlink overwrites
    try {
      fs.unlinkSync(APPROVAL_FILE);
    } catch (err) {
      if (err.code !== 'ENOENT') {
        throw err;
      }
    }

    // Write file securely with highly restrictive 0600 permissions
    fs.writeFileSync(APPROVAL_FILE, JSON.stringify(approvalData, null, 2), { mode: 0o600 });
    console.log(`✅ Developer approval successfully recorded and tied to diff hash: ${activeHash}`);
    process.exit(0);
  } catch (err) {
    console.error('Error: Failed to write developer approval file:', err.message);
    process.exit(1);
  }
}

// ------------------------------------------------------------------------------
// PROMPT OPERATION
// ------------------------------------------------------------------------------
function promptApproval(message, defaultOption) {
  console.log('============================================================');
  console.log('🚨 HOOK GATEWAY: DEVELOPER CONFIRMATION REQUIRED');
  console.log('============================================================');
  console.log(message);
  console.log('=============================================================');

  const isDefaultYes = defaultOption.toLowerCase() === 'y' || defaultOption.toLowerCase() === 'yes';
  const optionPrompt = isDefaultYes ? '[Y/n]' : '[y/N]';

  let response;

  try {
    // Open direct TTY read/write streams
    const ttyRead = fs.openSync('/dev/tty', 'r');
    const ttyWrite = fs.openSync('/dev/tty', 'w');

    fs.writeSync(ttyWrite, `Confirm ${optionPrompt}: `);

    const buffer = Buffer.alloc(1024);
    const bytesRead = fs.readSync(ttyRead, buffer, 0, 1024, null);

    fs.closeSync(ttyRead);
    fs.closeSync(ttyWrite);

    response = buffer.toString('utf8', 0, bytesRead).trim();
    if (!response) {
      response = defaultOption;
    }
  } catch (err) {
    console.error(err.message || err);
    // Fallback gracefully in headless/piped environments
    console.warn(`Non-interactive terminal detected. Fallback to default: ${defaultOption}`);
    response = defaultOption;
  }

  if (response.toLowerCase() === 'y' || response.toLowerCase() === 'yes') {
    writeApproval();
  } else {
    console.error('❌ Action aborted by developer.');
    process.exit(1);
  }
}

// ==============================================================================
// MAIN ENTRY
// ==============================================================================
function main() {
  const arg = process.argv[2] || '';

  if (arg === '-h' || arg === '--help') {
    console.log('Usage: node user-approval.js [--verify | --approve | message]');
    process.exit(0);
  }

  if (arg === '--verify') {
    verifyApproval();
  }

  if (arg === '--verify-plan') {
    verifyPlan();
  }

  if (arg === '--approve') {
    writeApproval();
  }

  const message = arg || 'Do you approve GPG-signing, committing, and pushing these changes?';
  const defaultOption = process.argv[3] || 'N';
  promptApproval(message, defaultOption);
}

main();
