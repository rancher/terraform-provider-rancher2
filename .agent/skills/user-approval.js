#!/usr/bin/env node
//
// Skill: user-approval.js
// Description: Securely prompts for developer approval and verifies user-approval status using diff-hash tieing.
//              Conforms to repository standards and rules.
// Usage:
//   Write Mode:   node .agent/skills/user-approval.js --approve
//   Verify Mode:  node .agent/skills/user-approval.js --verify
//   Prompt Mode:  node .agent/skills/user-approval.js [message] [defaultOption=N]

import { execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import path from 'path';

// Refuse to run rather than fall back to /tmp — a world-writable directory undermines the
// security properties of this approval gate (any local user/process could pre-create or swap
// the approval file there).
if (!process.env.HOME) {
  console.error("Error: HOME environment variable is not set. Refusing to run rather than write approval state to a shared/world-writable location.");
  process.exit(1);
}

const TARGET_DIR = path.resolve(process.env.HOME, '.gemini/tmp/terraform-provider-rancher2');
const APPROVAL_FILE = path.join(TARGET_DIR, 'user-approval.json');

// Helper to calculate active local diff hash securely using SHA-256
function calculateSHA256() {
  try {
    const diff = execSync('git diff HEAD', { stdio: ['ignore', 'pipe', 'ignore'] }).toString();
    return crypto.createHash('sha256').update(diff).digest('hex');
  } catch {
    return null;
  }
}

// ------------------------------------------------------------------------------
// VERIFY OPERATION
// ------------------------------------------------------------------------------
function verifyApproval() {
  console.log("Verifying developer manual IDE review approval status...");

  if (!fs.existsSync(APPROVAL_FILE)) {
    console.error("Error: Developer manual IDE review approval not found!");
    console.error("       In accordance with Gate 2 (IDE & Commit Gate) of 'development-process.md',");
    console.error("       you MUST request developer approval first: @user-approval");
    process.exit(1);
  }

  // Reject symbolic links to prevent symlink bypasses
  const stat = fs.lstatSync(APPROVAL_FILE);
  if (stat.isSymbolicLink()) {
    console.error("Error: Developer approval file is a symbolic link (Prohibited).");
    process.exit(1);
  }

  // Verify file ownership matches current user UID natively
  if (process.getuid) {
    const fileUid = fs.statSync(APPROVAL_FILE).uid;
    const currentUid = process.getuid();
    if (fileUid !== currentUid) {
      console.error(`Error: Developer approval file is not owned by the current user (UID: ${currentUid}, Owner: ${fileUid}).`);
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
    if (activeHash === null) {
      console.error("Error: Failed to calculate the active diff hash (is this a git repository?).");
      process.exit(1);
    }
    if (content.diff_hash !== activeHash) {
      console.error("Error: Local changes have been modified since your last manual developer approval!");
      console.error(`       Approved SHA-256 hash: ${content.diff_hash}`);
      console.error(`       Current active SHA-256 hash: ${activeHash}`);
      console.error("       Please request developer approval again on your latest changes.");
      process.exit(1);
    }

    console.log(`✅ Developer visual IDE review approval verified! (SHA-256 Hash: ${activeHash})`);
    process.exit(0);
  } catch (err) {
    console.error("Error: Failed to parse developer approval file:", err.message);
    process.exit(1);
  }
}

// ------------------------------------------------------------------------------
// APPROVE OPERATION
// ------------------------------------------------------------------------------
function writeApproval() {
  const activeHash = calculateSHA256();
  if (!activeHash) {
    console.error("Error: Failed to calculate active diff hash.");
    process.exit(1);
  }

  const approvalData = {
    status: 'approved',
    diff_hash: activeHash,
    timestamp: new Date().toISOString()
  };

  try {
    fs.mkdirSync(TARGET_DIR, { recursive: true });

    // Write to a freshly-created temp file, then atomically rename it into place, rather than
    // unlinking the existing file/symlink and writing separately. That leaves a TOCTOU window
    // where a symlink recreated at APPROVAL_FILE between the unlink and the write could
    // redirect it to an unintended path. fs.renameSync (rename(2) on the same filesystem)
    // replaces whatever is at the destination atomically without dereferencing an existing
    // symlink there.
    const tmpFile = path.join(TARGET_DIR, `.user-approval.${crypto.randomBytes(8).toString('hex')}.tmp`);
    try {
      fs.writeFileSync(tmpFile, JSON.stringify(approvalData, null, 2), { mode: 0o600 });
      fs.renameSync(tmpFile, APPROVAL_FILE);
    } catch (err) {
      try { fs.unlinkSync(tmpFile); } catch { /* best-effort cleanup */ }
      throw err;
    }
    console.log(`✅ Developer approval successfully recorded and tied to diff hash: ${activeHash}`);
    process.exit(0);
  } catch (err) {
    console.error("Error: Failed to write developer approval file:", err.message);
    process.exit(1);
  }
}

// ------------------------------------------------------------------------------
// PROMPT OPERATION
// ------------------------------------------------------------------------------
function promptApproval(message, defaultOption) {
  console.log("============================================================");
  console.log("🚨 HOOK GATEWAY: DEVELOPER CONFIRMATION REQUIRED");
  console.log("============================================================");
  console.log(message);
  console.log("=============================================================");

  const isDefaultYes = defaultOption.toLowerCase() === 'y' || defaultOption.toLowerCase() === 'yes';
  const optionPrompt = isDefaultYes ? "[Y/n]" : "[y/N]";

  let response = defaultOption;

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
  } catch {
    // Fallback gracefully in headless/piped environments without throwing
    console.warn(`Non-interactive terminal detected. Fallback to default: ${defaultOption}`);
    response = defaultOption;
  }

  if (response.toLowerCase() === 'y' || response.toLowerCase() === 'yes') {
    writeApproval();
  } else {
    console.error("❌ Action aborted by developer.");
    process.exit(1);
  }
}

// ==============================================================================
// MAIN ENTRY
// ==============================================================================
function main() {
  const arg = process.argv[2] || '';

  if (arg === '-h' || arg === '--help') {
    console.log("Usage: node user-approval.js [--verify | --approve | message]");
    process.exit(0);
  }

  if (arg === '--verify') {
    verifyApproval();
  }

  if (arg === '--approve') {
    writeApproval();
  }

  const message = arg || "Do you approve GPG-signing, committing, and pushing these changes?";
  const defaultOption = process.argv[3] || "N";
  promptApproval(message, defaultOption);
}

main();
