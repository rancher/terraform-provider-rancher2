#!/usr/bin/env node

import fs from 'fs';
import os from 'os';
import path from 'path';
import { handleCommitApproval } from '../../agent-scripts/after-ask.js';
import {
  calculateDiffHash,
  checkAndRevokeStaleGates,
  verifyPlanGate,
  verifyReviewGate,
} from '../../agent-scripts/gating.js';
import { resolveTargetDir } from '../../agent-scripts/workspace.js';

const hookName = path.basename(process.argv[1]);
const isStartup = hookName === '01-startup-context.js';
const introLog = `🔒 Hook: ${hookName} - ${isStartup ? 'Loading startup context...' : 'Loading hook context...'}`;
console.error(introLog);

const originalLog = console.log;
let hasLogged = false;

console.log = function (msg) {
  if (hasLogged) {
    return;
  }
  try {
    const parsed = JSON.parse(msg);
    if (parsed.systemMessage) {
      console.error(parsed.systemMessage);
    }
    const exitLog = `🔒 Hook: ${hookName} - ${isStartup ? 'context successfully loaded.' : 'Hook successfully loaded.'}`;
    console.error(exitLog);

    const msgs = [introLog];
    if (parsed.systemMessage) {
      msgs.push(parsed.systemMessage);
    }
    msgs.push(exitLog);
    parsed.systemMessage = msgs.join('\n');

    if (!parsed.decision && !isStartup) {
      parsed.decision = 'allow';
    }

    originalLog(JSON.stringify(parsed, null, 2));
    hasLogged = true;
  } catch (err) {
    console.error(err.message || err);
    originalLog(msg);
  }
};

process.on('exit', (code) => {
  if (!hasLogged) {
    const exitMsg = `🔒 Hook Error (${hookName}): Silent early exit detected with code ${code}.`;
    console.error(exitMsg);
    process.stdout.write(JSON.stringify({
      decision: 'allow',
      systemMessage: `${introLog}\n${exitMsg}`
    }) + '\n');
    hasLogged = true;
  }
});

process.on('uncaughtException', (err) => {
  const errMsg = `🔒 Hook Error (${hookName}): Unhandled exception - ${err.message || err}`;
  console.error(errMsg);
  if (!hasLogged) {
    process.stdout.write(JSON.stringify({
      decision: 'allow',
      systemMessage: `${introLog}\n${errMsg}`
    }) + '\n');
    hasLogged = true;
  }
  process.exit(1);
});

process.on('unhandledRejection', (reason) => {
  const errMsg = `🔒 Hook Error (${hookName}): Unhandled promise rejection - ${reason.message || reason}`;
  console.error(errMsg);
  if (!hasLogged) {
    process.stdout.write(JSON.stringify({
      decision: 'allow',
      systemMessage: `${introLog}\n${errMsg}`
    }) + '\n');
    hasLogged = true;
  }
  process.exit(1);
});

function revokeReviewState(targetDir) {
  const reviewApprovalFile = path.join(targetDir, 'review-approval.json');
  try {
    if (fs.existsSync(reviewApprovalFile)) {
      fs.unlinkSync(reviewApprovalFile);
      console.error('❌ Gate 2 (Review) Revoked: User rejected the commit. Review approval has been deleted.');
    }
  } catch (err) {
    console.warn(`Warning: Failed to revoke review state. Error: ${err.message || err}`);
  }
}

function preCommitPhaseInterruption(inputData, targetDir) {
  const flagFile = path.join(targetDir, 'require-ask-user.flag');

  if (fs.existsSync(flagFile)) {
    if (inputData.tool_name !== 'ask_user') {
      console.log(
        JSON.stringify({
          decision: 'deny',
          reason:
            '🔒 Security Policy Violation: The review phase has completed successfully. All tools are blocked until you present the changes to the user for commit approval.\n\n' +
            'Please call the `ask_user` tool to request commit approval.',
          systemMessage: '🔒 Security Block: Call the `ask_user` tool to proceed to the Commit Phase. 👉 ACTION REQUIRED: Present the changes to the user for commit approval.',
        }),
      );
      process.exit(0);
    }

    // Present the suggested commit message from the review agent
    let suggestedCommitMessage = 'chore: automated development commit';
    try {
      const reviewApprovalFile = path.join(targetDir, 'review-approval.json');
      if (fs.existsSync(reviewApprovalFile)) {
        const approvalData = JSON.parse(fs.readFileSync(reviewApprovalFile, 'utf-8'));
        if (approvalData.suggested_commit_message) {
          suggestedCommitMessage = approvalData.suggested_commit_message;
        }
      }
    } catch (err) {
      console.warn('Failed to retrieve suggested commit message from review approval:', err.message);
    }

    const modifiedInput = inputData.tool_input || {};
    const reviewContext = `\n\n### 🔍 AUTOMATED REVIEW COMPLETE 🔍 ###\nThe review agent has verified the changes and formulated the following commit message:\n\nCommit Message: \`${suggestedCommitMessage}\`\n\nPlease review the code changes in your IDE. Do you approve these changes for commit? (Yes/No)`;

    // Strip any raw "Commit Message:" directives the main agent might have formulated to avoid collision/parse issues
    const replaceCommitMsg = (str) =>
      typeof str === 'string' ? str.replace(/Commit Message/gi, 'Proposed Message') : str;

    if (modifiedInput.questions && Array.isArray(modifiedInput.questions) && modifiedInput.questions.length > 0) {
      modifiedInput.questions[0].question = replaceCommitMsg(modifiedInput.questions[0].question) + reviewContext;
    } else if (modifiedInput.question !== undefined) {
      modifiedInput.question = replaceCommitMsg(modifiedInput.question) + reviewContext;
    } else if (modifiedInput.prompt !== undefined) {
      modifiedInput.prompt = replaceCommitMsg(modifiedInput.prompt) + reviewContext;
    } else {
      modifiedInput.question = reviewContext;
    }

    console.log(
      JSON.stringify({
        decision: 'allow',
        tool_input: modifiedInput,
        systemMessage: '🟢 Pre-Commit Phase: Appended review agent commit message to the user prompt.',
      }),
    );
    process.exit(0);
  }
}

function beforeAskUser(inputData, targetDir) {
  const { tool_name, tool_input } = inputData;

  if (tool_name !== 'ask_user' || !tool_input) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not ask_user.' }));
    process.exit(0);
  }

  const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
  if (fs.existsSync(inPlanModeFile)) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, currently in plan mode (deferring to plan phase hook).' }));
    process.exit(0);
  }

  const safeToolInput = JSON.stringify(tool_input);
  const isCommitAsk =
    safeToolInput.includes('approve these changes for commit') ||
    safeToolInput.includes('Commit Message:') ||
    safeToolInput.includes('Gate 3');

  if (!isCommitAsk) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, not a commit approval request.' }));
    process.exit(0);
  }

  const planHash = verifyPlanGate(targetDir);
  if (!planHash) {
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '🔒 Security Policy Violation: You cannot ask for Developer Commit Approval (Gate 3) because Gate 1 (Planning Gate) is missing or invalid!\n\n' +
          'Please obtain planning approval from the developer first.',
        systemMessage: '🔒 Security Block: Gate 1 must be approved before commit.',
      }),
    );
    process.exit(0);
  }

  const diffHash = calculateDiffHash();

  checkAndRevokeStaleGates(targetDir, diffHash, planHash);

  const reviewPassed = verifyReviewGate(targetDir, diffHash, planHash);
  if (!reviewPassed) {
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '🔒 Security Policy Violation: You cannot ask for Developer Commit Approval (Gate 3) because the Review prerequisite is missing or invalid!\n\n' +
          'In accordance with our zero-trust pipeline, you MUST successfully run the Review Subagent first:\n' +
          '   `invoke_agent(agent_name="review_agent", prompt="Please review my current changes.")`',
        systemMessage: '🔒 Security Block: Review must be approved before commit.',
      }),
    );
    process.exit(0);
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: beforeAskUser check complete, execution allowed.' }));
  process.exit(0);
}

function afterAskUser(inputData, targetDir) {
  const { tool_name, tool_input, tool_response } = inputData;

  if (tool_name !== 'ask_user' || !tool_input || !tool_response) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not ask_user or response is missing.' }));
    process.exit(0);
  }

  const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
  if (fs.existsSync(inPlanModeFile)) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, currently in plan mode.' }));
    process.exit(0);
  }

  const FLAG_FILE = path.join(targetDir, 'require-ask-user.flag');
  if (fs.existsSync(FLAG_FILE)) {
    try {
      fs.unlinkSync(FLAG_FILE);
    } catch (err) {
      console.warn(`Warning: Failed to delete require-ask-user.flag. Error: ${err.message || err}`);
    }
  }

  const safeToolInput = JSON.stringify(tool_input);
  const isCommitAsk =
    safeToolInput.includes('approve these changes for commit') ||
    safeToolInput.includes('Commit Message:') ||
    safeToolInput.includes('Gate 3');

  if (!isCommitAsk) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, not a commit approval request.' }));
    process.exit(0);
  }

  let answerText = '';
  try {
    console.error(`🔒 Hook Debug: Raw tool_response: ${JSON.stringify(tool_response)}`);
    let res = typeof tool_response === 'string' ? JSON.parse(tool_response) : tool_response;
    if (res && res.output && typeof res.output === 'string') {
      try {
        res = JSON.parse(res.output);
      } catch (err) {
        console.error(`🔒 Hook Debug: Failed to parse res.output: ${err.message}`);
      }
    }
    console.error(`🔒 Hook Debug: Parsed res: ${JSON.stringify(res)}`);
    if (res && res.answers) {
      answerText = Object.values(res.answers)[0] || '';
    } else if (res && res.llmContent) {
      const parsed = JSON.parse(res.llmContent);
      if (parsed && parsed.answers) {
        answerText = Object.values(parsed.answers)[0] || '';
      } else {
        answerText = Object.values(parsed)[0] || '';
      }
    } else if (res) {
      answerText = Object.values(res)[0] || '';
    }
  } catch (err) {
    console.error(`🔒 Hook Debug: Catch block triggered: ${err.message}`);
    answerText = (tool_response && tool_response.llmContent) || JSON.stringify(tool_response);
  }

  const safeAnswerText = String(answerText || '');
  console.error(`🔒 Hook Debug: safeAnswerText: "${safeAnswerText}"`);
  const isApproved =
    safeAnswerText.toLowerCase() === 'yes' ||
    safeAnswerText.toLowerCase() === 'y' ||
    safeAnswerText.toLowerCase() === 'approve' ||
    safeAnswerText.toLowerCase() === 'approve plan' ||
    safeAnswerText.toLowerCase() === 'approve commit' ||
    safeAnswerText.toLowerCase() === 'looks good';

  console.error(`🔒 Hook Debug: isApproved: ${isApproved}`);

  if (!isApproved) {
    if (isCommitAsk) {
      revokeReviewState(targetDir);
    }
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: User rejected commit approval.' }));
    process.exit(0);
  }

  const homeDir = os.homedir();
  const sshPubKeyFile = path.resolve(homeDir, '.gemini/ssh-key.pub');
  const sshPrivKeyFile = sshPubKeyFile.replace(/\.pub$/, '');

  let keyExists = false;
  let keyErrorMsg = '';
  try {
    fs.accessSync(sshPubKeyFile, fs.constants.R_OK);

    let hasPrivKey = false;
    try {
      fs.accessSync(sshPrivKeyFile, fs.constants.R_OK);
      hasPrivKey = true;
    } catch (err) {
      console.error(`🔒 Hook Debug: Private key access check failed: ${err.message}`);
    }

    if (hasPrivKey || process.env.SSH_AUTH_SOCK) {
      keyExists = true;
    } else {
      keyErrorMsg = ' (Private key not found on disk and SSH agent is not running).';
    }
  } catch (err) {
    if (err.code === 'EACCES' || err.code === 'EPERM') {
      keyErrorMsg = ' (Permission denied! Please check read permissions for the public key).';
    } else if (err.code !== 'ENOENT') {
      keyErrorMsg = ` (Error checking public key: ${err.message}).`;
    }
  }

  if (!keyExists) {
    console.log(
      JSON.stringify({
        decision: 'allow',
        systemMessage: `🔒 Hook Notification: Cryptographic signing skipped because your SSH key material is not accessible at ~/.gemini/ssh-key.pub${keyErrorMsg} Please ensure your public key is linked here and your SSH agent is active.`,
      }),
    );
    process.exit(0);
  }

  const promptText = tool_input.questions && tool_input.questions[0] ? tool_input.questions[0].question : '';

  if (isCommitAsk) {
    const planHash = verifyPlanGate(targetDir);
    const diffHash = calculateDiffHash();
    checkAndRevokeStaleGates(targetDir, diffHash, planHash);

    const reviewPassed = verifyReviewGate(targetDir, diffHash, planHash);
    if (!reviewPassed) {
      console.log(
        JSON.stringify({
          decision: 'allow',
          systemMessage:
            '🔒 Security Block: Commit aborted. Review approval is missing or was invalidated by recent file changes. Please run the review agent again.',
        }),
      );
      process.exit(0);
    }

    handleCommitApproval(targetDir, sshPubKeyFile, promptText);
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: afterAskUser check complete, execution allowed.' }));
  process.exit(0);
}

function main() {
  let inputData;
  try {
    inputData = JSON.parse(fs.readFileSync(0, 'utf-8'));
  } catch (err) {
    console.error('Failed to parse stdin JSON in 04-commit-phase:', err.message || err);
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason: '🔒 Security Policy Violation: Failed to parse input payload in 04-commit-phase hook.',
        systemMessage: '🔒 Security Block: Hook input parsing failed, execution denied.',
      }),
    );
    process.exit(0);
  }

  const targetDir = resolveTargetDir();
  const args = process.argv.slice(2);

  if (args.includes('--before-ask')) {
    beforeAskUser(inputData, targetDir);
  } else if (args.includes('--after-ask')) {
    afterAskUser(inputData, targetDir);
  } else {
    preCommitPhaseInterruption(inputData, targetDir);
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Pre-commit phase check complete, execution allowed.' }));
    process.exit(0);
  }
}

main();
