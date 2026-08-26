#!/usr/bin/env node

import fs from 'fs';
import os from 'os';
import path from 'path';
import { handlePlanApproval } from '../../agent-scripts/after-ask.js';
import { findLatestActivePlan, verifyPlanGate } from '../../agent-scripts/gating.js';
import { validatePlanContent } from '../../agent-scripts/planning.js';
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

function clearPrePlanFlag(targetDir) {
  const requirePlanModeFile = path.join(targetDir, 'require-plan-mode.flag');
  if (fs.existsSync(requirePlanModeFile)) {
    try {
      fs.unlinkSync(requirePlanModeFile);
    } catch (err) {
      console.warn(`Warning: Failed to delete require-plan-mode.flag. Error: ${err.message || err}`);
    }
  }

  // Force Plan mode by setting the active planning flag
  fs.writeFileSync(path.join(targetDir, 'in-plan-mode.flag'), 'true', 'utf-8');

  console.log(
    JSON.stringify({
      decision: 'allow',
      systemMessage:
        '✨ You have successfully entered Plan Mode. All tools are now unlocked for planning. 👉 ACTION REQUIRED: Draft your plan and then use `ask_user` to request approval.',
    }),
  );
  process.exit(0);
}

function verifyExitPlanMode(inputData, targetDir) {
  // BeforeTool hook for exit_plan_mode
  if (inputData.tool_name !== 'exit_plan_mode') {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not exit_plan_mode.' }));
    process.exit(0);
  }

  const planHash = verifyPlanGate(targetDir);
  if (!planHash) {
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '🔒 Security Policy Violation: You cannot exit Plan Mode until the user has cryptographically approved the plan.\n\n' +
          'Please present the plan to the user using the `ask_user` tool first and ask for their cryptographic approval. Only after the user approves the plan will you be permitted to exit plan mode.',
        systemMessage: '🔒 Security Block: User cryptographic plan approval required before exiting plan mode. 👉 ACTION REQUIRED: Use ask_user to get plan approval.',
      }),
    );
    process.exit(0);
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: verifyExitPlanMode check complete, execution allowed.' }));
  process.exit(0);
}

function clearPlanModeFlag(inputData, targetDir) {
  // AfterTool hook for exit_plan_mode
  if (inputData.tool_name !== 'exit_plan_mode') {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not exit_plan_mode.' }));
    process.exit(0);
  }

  const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
  if (fs.existsSync(inPlanModeFile)) {
    try {
      fs.unlinkSync(inPlanModeFile);
    } catch (err) {
      console.warn(`Warning: Failed to delete in-plan-mode.flag. Error: ${err.message || err}`);
    }
  }

  console.log(
    JSON.stringify({
      decision: 'allow',
      systemMessage: '✅ Exited Plan Mode. Implementation phase successfully unlocked! 👉 ACTION REQUIRED: Proceed immediately to Implement your plan, then move to the Review Phase by invoking the review_agent.',
    }),
  );
  process.exit(0);
}

function beforeAskUserPlan(inputData, targetDir) {
  const { tool_name, tool_input } = inputData;

  if (tool_name !== 'ask_user' || !tool_input) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not ask_user.' }));
    process.exit(0);
  }

  const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
  if (!fs.existsSync(inPlanModeFile)) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, currently not in plan mode.' }));
    process.exit(0);
  }

  // Verify the plan is valid before allowing ask_user to prompt the user
  const activePlan = findLatestActivePlan(targetDir);
  if (!activePlan) {
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '🔒 Security Policy Violation: Active plan file not found in session directory!\n\n' +
          'Please write your plan file as a markdown document under plans/ first before calling `ask_user`.',
        systemMessage: '🔒 Security Block: Active plan file not found.',
      }),
    );
    process.exit(0);
  }

  const validation = validatePlanContent(activePlan);
  if (!validation.valid) {
    const errorsList = validation.errors.map((err) => `  - ${err}`).join('\n');
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '❌ Plan Format Validation Failure!\n\n' +
          'The proposed plan has invalid structure and violates repository standards:\n' +
          errorsList +
          '\n\n' +
          '👉 ACTION REQUIRED: You must rewrite the plan file at:\n' +
          `   ${activePlan}\n` +
          'to satisfy all repository requirements (include markdown checklist - [ ], comprehensive tests, quality gates, agentic framework maintenance, and documentation updates) before you can ask the user for approval.',
        systemMessage: '🔒 Security Block: Plan format validation failed.',
      }),
    );
    process.exit(0);
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: beforeAskUserPlan check complete, execution allowed.' }));
  process.exit(0);
}

function askUserPlanProof(inputData, targetDir) {
  const { tool_name, tool_input, tool_response } = inputData;

  if (tool_name !== 'ask_user' || !tool_response) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not ask_user or response is missing.' }));
    process.exit(0);
  }

  const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
  if (!fs.existsSync(inPlanModeFile)) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, currently not in plan mode.' }));
    process.exit(0);
  }

  const safeToolInput = JSON.stringify(tool_input);
  const isPlanAsk =
    safeToolInput.includes('plan') || safeToolInput.includes('blueprint') || safeToolInput.includes('Planning');

  if (!isPlanAsk) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, not a plan approval request.' }));
    process.exit(0);
  }

  let answerText;
  try {
    let res = typeof tool_response === 'string' ? JSON.parse(tool_response) : tool_response;
    if (res.output && typeof res.output === 'string') {
      try {
        res = JSON.parse(res.output);
      } catch (err) {
        console.error(err.message || err);
      }
    }
    if (res.answers) {
      answerText = Object.values(res.answers)[0] || '';
    } else if (res.llmContent) {
      const parsed = JSON.parse(res.llmContent);
      if (parsed && parsed.answers) {
        answerText = Object.values(parsed.answers)[0] || '';
      } else {
        answerText = Object.values(parsed)[0] || '';
      }
    } else {
      answerText = Object.values(res)[0] || '';
    }
  } catch (err) {
    console.error(err.message || err);
    answerText = tool_response.llmContent || JSON.stringify(tool_response);
  }

  const isApproved =
    answerText.toLowerCase() === 'yes' ||
    answerText.toLowerCase() === 'y' ||
    answerText.toLowerCase() === 'approve' ||
    answerText.toLowerCase() === 'approve plan' ||
    answerText.toLowerCase() === 'looks good';

  if (!isApproved) {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: User did not approve plan.' }));
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

  let promptText;
  try {
    const question = tool_input.questions && tool_input.questions[0];
    promptText = (question && question.question) || JSON.stringify(tool_input);
  } catch (err) {
    console.error(err.message || err);
    promptText = JSON.stringify(tool_input);
  }

  const result = handlePlanApproval(targetDir, sshPubKeyFile, promptText);

  console.log(
    JSON.stringify({
      decision: 'allow',
      systemMessage: result.systemMessage + ' You may now call `exit_plan_mode`.',
    }),
  );
  process.exit(0);
}

function prePlanPhaseInterruption(inputData, targetDir) {
  const requirePlanModeFile = path.join(targetDir, 'require-plan-mode.flag');
  if (fs.existsSync(requirePlanModeFile)) {
    if (inputData.tool_name !== 'enter_plan_mode') {
      console.log(
        JSON.stringify({
          decision: 'deny',
          reason:
            '🔒 Security Policy Violation: A plan has not been accepted. All tools are blocked until you enter the Plan Phase.\n\n' +
            'Please call the `enter_plan_mode` tool to transition the workflow into the Plan Phase before utilizing other tools.',
          systemMessage: '🔒 Security Block: Call the `enter_plan_mode` tool first.',
        }),
      );
      process.exit(0);
    }
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: prePlanPhaseInterruption check complete, execution allowed.' }));
  process.exit(0);
}

function verifyGateArtifactProtection(inputData) {
  const { tool_name, tool_input } = inputData;

  if (
    tool_name === 'write_file' ||
    tool_name === 'replace' ||
    tool_name === 'edit_file' ||
    tool_name === 'create_file'
  ) {
    if (tool_input) {
      const filePath = tool_input.file_path || tool_input.path || '';
      const fileName = path.basename(filePath);

      const isApprovalFile =
        /^(plan-approval|test-approval|review-approval|user-approval)\.(json|challenge|age|sig)$/.test(fileName) ||
        fileName.endsWith('-approval.json') ||
        fileName.endsWith('.sig');

      if (isApprovalFile) {
        console.log(
          JSON.stringify({
            decision: 'deny',
            reason:
              '🔒 Security Policy Violation: Direct creation or modification of gate approvals or signature files is strictly prohibited.\n\n' +
              'Gating approval files must ONLY be generated automatically and securely by our pipeline hooks and sub-agents.',
            systemMessage: '🔒 Security Block: Direct manipulation of approval files is prohibited.',
          }),
        );
        process.exit(0);
      }
    }
  }
}

function main() {
  let inputData;
  try {
    inputData = JSON.parse(fs.readFileSync(0, 'utf-8'));
  } catch (err) {
    console.error('Failed to parse stdin JSON:', err);
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Failed to parse input, allowing execution by default.' }));
    process.exit(0);
  }

  // Enforce Gate Artifact Tamper Protection
  verifyGateArtifactProtection(inputData);

  const targetDir = resolveTargetDir();
  if (!fs.existsSync(targetDir)) {
    fs.mkdirSync(targetDir, { recursive: true });
  }

  const args = process.argv.slice(2);
  if (args.includes('--enter-proof')) {
    clearPrePlanFlag(targetDir);
  } else if (args.includes('--verify-exit')) {
    verifyExitPlanMode(inputData, targetDir);
  } else if (args.includes('--clear-plan-mode')) {
    clearPlanModeFlag(inputData, targetDir);
  } else if (args.includes('--before-ask-proof')) {
    beforeAskUserPlan(inputData, targetDir);
  } else if (args.includes('--ask-proof')) {
    askUserPlanProof(inputData, targetDir);
  } else {
    prePlanPhaseInterruption(inputData, targetDir);
  }
}

main();
