#!/usr/bin/env node
import fs from 'fs';
import path from 'path';
import { saveReport } from '../../agent-scripts/after-invoke.js';
import { calculateDiffHash, findLatestActivePlan, verifyPlanGate } from '../../agent-scripts/gating.js';
import { runPreReviewTests } from '../../agent-scripts/testing.js';
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

const TARGET_DIR = resolveTargetDir();
const LOGS_DIR = path.join(TARGET_DIR, 'logs');
const REVIEW_APPROVAL_FILE = path.join(TARGET_DIR, 'review-approval.json');

function revokeReviewState() {
  try {
    if (fs.existsSync(REVIEW_APPROVAL_FILE)) {
      fs.unlinkSync(REVIEW_APPROVAL_FILE);
    }
  } catch (err) {
    console.warn(`Warning: Failed to revoke review state. Error: ${err.message || err}`);
  }
  const flagFile = path.join(TARGET_DIR, 'require-ask-user.flag');
  try {
    if (fs.existsSync(flagFile)) {
      fs.unlinkSync(flagFile);
    }
  } catch (err) {
    console.warn(`Warning: Failed to delete require-ask-user.flag. Error: ${err.message || err}`);
  }
}

function verifyPlanning() {
  const planHash = verifyPlanGate(TARGET_DIR);
  if (!planHash) {
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason:
          '🔒 Security Policy Violation: You cannot execute pre-review testing because Gate 1 (Planning Gate) is missing or invalid!\n\n' +
          'Please obtain planning approval from the developer first by executing `exit_plan_mode`.',
        systemMessage: '🔒 Security Block: Gate 1 must be approved before testing/review.',
      }),
    );
    process.exit(0);
  }
}

function preReviewTesting(tool_input) {
  const result = runPreReviewTests();
  if (result.success) {
    const modifiedToolInput = tool_input;
    const activePlan = findLatestActivePlan(TARGET_DIR);
    if (activePlan && fs.existsSync(activePlan)) {
      const planContent = fs.readFileSync(activePlan, 'utf-8');
      modifiedToolInput.prompt = (modifiedToolInput.prompt || '') + '\n\n### ACTIVE PLAN CONTEXT ###\n' + planContent;
    }

    console.log(
      JSON.stringify({
        decision: 'allow',
        tool_input: modifiedToolInput,
        systemMessage: '🟢 Pre-Review Testing Passed. Starting review agent with active plan context.',
      }),
    );
    process.exit(0);
  } else {
    revokeReviewState();
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason: `Pre-review testing failed. Please fix the following issues before invoking the review agent:\n\n${result.failureOutput}`,
        systemMessage: '🔒 Review blocked: Automated tests failed. 👉 ACTION REQUIRED: Please fix the failing tests and re-run pre-review testing.',
      }),
    );
    process.exit(0);
  }
}

function afterInvoke(inputData) {
  const { tool_name, tool_input, tool_response } = inputData;

  if (tool_name !== 'invoke_agent' || !tool_input || tool_input.agent_name !== 'review_agent') {
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not invoke_agent or agent is not review_agent.' }));
    process.exit(0);
  }

  if (!tool_response || !tool_response.llmContent) {
    console.error('🔒 Hook Error: Sub-agent response is missing, empty, or unparsable.');
    revokeReviewState();
    console.log(JSON.stringify({ 
      decision: 'allow',
      systemMessage: '🔒 Review failed: The review_agent returned an empty report. 👉 ACTION REQUIRED: Please re-run the review_agent.'
    }));
    process.exit(0);
  }

  let report = '';
  if (Array.isArray(tool_response.llmContent)) {
    report = tool_response.llmContent.map((item) => item.text || '').join('\n');
  } else if (typeof tool_response.llmContent === 'string') {
    report = tool_response.llmContent;
  }

  if (!report || report.trim() === '') {
    console.error('🔒 Hook Error: Sub-agent report is empty or unparsable.');
    revokeReviewState();
    console.log(JSON.stringify({ 
      decision: 'allow',
      systemMessage: '🔒 Review failed: The review_agent returned an empty report. 👉 ACTION REQUIRED: Please re-run the review_agent.'
    }));
    process.exit(0);
  }

  saveReport('review_agent', report, LOGS_DIR);

  const planHash = verifyPlanGate(TARGET_DIR);
  if (!planHash) {
    revokeReviewState();
    console.log(
      JSON.stringify({
        decision: 'allow',
        systemMessage: `🔒 Hook Notification: Sub-agent review_agent finished but no signature was written because Gate 1 (Planning Gate) is missing or invalid.`,
      }),
    );
    process.exit(0);
  }

  const reportLower = report.toLowerCase();
  const requiredTopics = [
    'pass 1',
    'pass 2',
    'pass 3',
    'security',
    'standard',
    'performance',
    'logic',
    'error handling',
    'concurrency',
    'edge cases',
    'maintainability',
    'testability',
    'commit title',
    'commit message',
  ];
  const missingTopics = requiredTopics.filter((topic) => !reportLower.includes(topic));

  if (missingTopics.length > 0) {
    revokeReviewState();
    console.log(
      JSON.stringify({
        decision: 'allow',
        systemMessage: `⚠️ Review Verification Failed: The review agent's report is incomplete. It is missing explicit checks for: ${missingTopics.join(', ')}.\n\n👉 ACTION REQUIRED: Please explicitly instruct the review agent to perform these checks and re-run the review_agent to proceed.`,
      }),
    );
    process.exit(0);
  }

  const hasCheckedPasses =
    /- \[[xX]\] Pass 1/i.test(report) &&
    /- \[[xX]\] Pass 2/i.test(report) &&
    /- \[[xX]\] Pass 3/i.test(report) &&
    /- \[[xX]\] Pass 4/i.test(report);

  if (!hasCheckedPasses) {
    revokeReviewState();
    console.log(
      JSON.stringify({
        decision: 'allow',
        systemMessage: '⚠️ Review Verification Failed: The review agent\'s passes are incomplete or unchecked.\n\nAll 4 sequential passes must be checked as complete (e.g. - [x] Pass 1, - [x] Pass 2, etc.) in the report checklist to proceed.',
      }),
    );
    process.exit(0);
  }

  const hasCleanMarker = /0 comments\/findings|0 findings/i.test(report);
  const isPerfect = hasCleanMarker || reportLower.includes('perfect') || reportLower.includes('approved');
  const hasComments = reportLower.includes('finding') || reportLower.includes('issue') || reportLower.includes('suggestion') || reportLower.includes('comment');

  if (!isPerfect) {
    revokeReviewState();
    let sysMsg = '⚠️ Review Verification Failed: The review agent did not approve the changes but provided no explicit comments. 👉 ACTION REQUIRED: Please re-run the review_agent.';
    if (hasComments) {
      sysMsg = '⚠️ Review Verification Failed: The review agent found issues. 👉 ACTION REQUIRED: Please review the comments and findings from the review_agent, implement the necessary fixes, and then re-run the review_agent.';
    }
    
    console.log(
      JSON.stringify({
        decision: 'allow',
        systemMessage: sysMsg,
      }),
    );
    process.exit(0);
  }

  const diffHash = calculateDiffHash();
  if (planHash && diffHash) {
    try {
      fs.unlinkSync(REVIEW_APPROVAL_FILE);
    } catch (err) {
      if (err.code !== 'ENOENT') {
        console.warn(`Warning: Failed to unlink review approval file. Error: ${err.message || err}`);
      }
    }

    fs.writeFileSync(
      REVIEW_APPROVAL_FILE,
      JSON.stringify(
        {
          status: 'approved',
          plan_hash: planHash,
          diff_hash: diffHash,
          timestamp: new Date().toISOString(),
        },
        null,
        2,
      ),
      { mode: 0o400 },
    );

    fs.writeFileSync(path.join(TARGET_DIR, 'require-ask-user.flag'), 'true', 'utf-8');
  }

  console.log(
    JSON.stringify({
      decision: 'allow',
      systemMessage: '🟢 Gate 2 (Review) Cryptographically Signed. Multi-pass review successful. 👉 ACTION REQUIRED: You must now move to the Commit Phase by calling the ask_user tool to request commit approval.',
    }),
  );
  process.exit(0);
}

function main() {
  let inputData;
  try {
    inputData = JSON.parse(fs.readFileSync(0, 'utf-8'));
  } catch (err) {
    console.error('Failed to parse stdin JSON in 03-review-phase:', err.message || err);
    console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Failed to parse input, allowing execution by default.' }));
    process.exit(0);
  }

  const args = process.argv.slice(2);

  if (args.includes('--after-invoke')) {
    afterInvoke(inputData);
  } else {
    const { tool_name, tool_input } = inputData;

    // Proceed only if the target is the review agent being invoked
    if (tool_name !== 'invoke_agent' || !tool_input || tool_input.agent_name !== 'review_agent') {
      console.log(JSON.stringify({ decision: 'allow', systemMessage: '🔒 Hook Notification: Execution allowed, tool is not invoke_agent or agent is not review_agent.' }));
      process.exit(0);
    }

    verifyPlanning();
    preReviewTesting(tool_input);
  }
}

main();
