#!/usr/bin/env node

import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';
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

function main() {
  // Log diagnostics to stderr to comply with the silence rule on stdout
  console.error('Loading session-start workspace context...');

  let combinedContext = `###############################################################################
#                           CRITICAL AGENT MANDATES                                 #
#                                                                                   #
# 1. YOU MUST FOLLOW THE DEVELOPMENT PROCESS IN 'docs/development/AgenticFramework.md'. #
# 2. YOU ARE STRICTLY FORBIDDEN FROM EXECUTING ANY COMMIT OR PUSH COMMANDS.         #
#    COMMITS AND PUSHES ARE SOLELY MANAGED OUT-OF-BAND BY SYSTEM HOOKS.             #
# 3. SOURCE EDITS ARE BLOCKED UNTIL PLAN APPROVAL IS GRANTED.                       #
#    ALL ACTIVE TASK CHECKLISTS RESIDE STRICTLY INSIDE THE PLAN.                    #
# 4. WE ENFORCE A GATED 4-PHASE LIFECYCLE (Plan, Implement, Review, Commit).        #
#    YOU MUST TRANSITION THROUGH THESE PHASES SEQUENTIALLY WITHOUT SKIPPING.        #
#    UPON COMPLETING A PHASE, IMMEDIATELY PROCEED TO THE NEXT.                      #
#                                                                                   #
# FAILURE TO COMPLY WILL TRIGGER SECURITY BLOCKS AND PROCESS TERMINATION.           #
###############################################################################\n\n`;

  const inNixShell = !!process.env.IN_NIX_SHELL;
  if (inNixShell) {
    combinedContext += '✅ NIX ENVIRONMENT: Session is securely running inside a hermetic Nix shell.\n\n';
    console.error('Nix shell environment verified.');
  } else {
    combinedContext +=
      '⚠️ NIX ENVIRONMENT WARNING: Session is NOT running inside a Nix shell. Standard dependencies may be missing. Advise the developer to run `nix develop`.\n\n';
    console.error('Warning: Not running in a Nix shell.');
  }

  const frameworkDocPath = 'docs/development/AgenticFramework.md';

  if (fs.existsSync(frameworkDocPath)) {
    combinedContext += '# Context from docs/development/AgenticFramework.md\n\n';
    combinedContext += fs.readFileSync(frameworkDocPath, 'utf-8');
    combinedContext += '\n\n';
    console.error(`Loaded ${frameworkDocPath}`);
  } else {
    console.error(`Warning: ${frameworkDocPath} not found`);
  }

  combinedContext += `###############################################################################
#                           IMMEDIATE ACTION REQUIRED                               #
#                                                                                   #
# You must immediately enter PLAN MODE as your first action in this session.        #
# Evaluate the user's initial request and draft a step-by-step imperative plan.     #
#                                                                                   #
# PLAN FORMAT REQUIREMENTS:                                                         #
# Your plan MUST include a markdown checklist (using "- [ ]") that covers:          #
# 1. The specific implementation tasks.                                             #
# 2. Running comprehensive tests and linters.                                       #
# 3. Maintaining the agentic framework if improvements or bugs are found.           #
# 4. Enforcing standard quality gates.                                              #
# 5. Updating documentation to describe the changes.                                #
#                                                                                   #
# REVIEW AGENT INSTRUCTIONS:                                                        #
# When invoking the review_agent, you must explicitly instruct it to check for:     #
# security, coding standards, spelling/wording, and an automation audit.            #
# It must also be instructed to provide a Commit Title and Commit Message,          #
# and output an explicit approval status.                                           #
###############################################################################\n`;

  combinedContext += `\n👉 ACTION REQUIRED: You must call the \`enter_plan_mode\` tool to formally enter the Plan Phase before utilizing other tools or modifying code.\n`;

  // Enforce enter_plan_mode on startup by writing a flag file
  try {
    const targetDir = resolveTargetDir();
    fs.mkdirSync(targetDir, { recursive: true });
    fs.writeFileSync(path.join(targetDir, 'require-plan-mode.flag'), 'true', 'utf-8');

    // Clean up stale flags from previous sessions to ensure a clean slate
    const inPlanModeFile = path.join(targetDir, 'in-plan-mode.flag');
    if (fs.existsSync(inPlanModeFile)) {
      fs.unlinkSync(inPlanModeFile);
    }
    const requireAskUserFile = path.join(targetDir, 'require-ask-user.flag');
    if (fs.existsSync(requireAskUserFile)) {
      fs.unlinkSync(requireAskUserFile);
    }
  } catch (err) {
    console.error(
      `Warning: Failed to create require-plan-mode.flag. This can be safely ignored. Error: ${err.message || err}`,
    );
  }

  // Enforce OS-level read-only permissions on AI exclude files to prevent agent tampering
  try {
    let repoRoot = process.cwd();
    try {
      repoRoot = execSync('git rev-parse --show-toplevel', { stdio: 'pipe' }).toString().trim();
    } catch (err) {
      console.error(`Warning: Failed to determine git repo root: ${err.message}`);
    }

    const excludeFiles = ['.aiexclude', '.claudeignore'];
    for (const file of excludeFiles) {
      const filePath = path.join(repoRoot, file);
      if (fs.existsSync(filePath)) {
        fs.chmodSync(filePath, 0o400);
      }
    }
  } catch (err) {
    console.error(`Warning: Failed to set read-only permissions on exclude files. Error: ${err.message || err}`);
  }

  // Output clean JSON structure to stdout
  const output = {
    hookSpecificOutput: {
      additionalContext: combinedContext,
    },
    systemMessage: `✨ Workspace context injected. ${
      inNixShell ? '[Nix Shell: Active]' : '[Nix Shell: Inactive]'
    } 👉 ACTION REQUIRED: Enter Plan Mode via 02-plan-phase.js.`,
  };

  console.log(JSON.stringify(output, null, 2));
  process.exit(0);
}

main();
