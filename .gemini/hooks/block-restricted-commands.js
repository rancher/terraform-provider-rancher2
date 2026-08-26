#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { verifyGitCommand } from '../../agent-scripts/security.js';

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
      decision: 'deny',
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
      decision: 'deny',
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
      decision: 'deny',
      systemMessage: `${introLog}\n${errMsg}`
    }) + '\n');
    hasLogged = true;
  }
  process.exit(1);
});

function main() {
  let inputData;
  try {
    inputData = JSON.parse(fs.readFileSync(0, 'utf-8'));
  } catch (err) {
    console.error('Failed to parse stdin JSON in block-restricted-commands:', err.message || err);
    console.log(
      JSON.stringify({
        decision: 'deny',
        reason: '🔒 Security Policy Violation: Failed to parse input payload in block-restricted-commands hook.',
        systemMessage: '🔒 Security Block: Hook input parsing failed, execution denied.',
      }),
    );
    process.exit(0);
  }

  const { tool_name, tool_input, cwd } = inputData;

  if (tool_name === 'run_shell_command' && tool_input && tool_input.command) {
    const result = verifyGitCommand(tool_input.command, cwd || process.cwd());
    if (result && result.decision === 'deny') {
      console.log(
        JSON.stringify({
          decision: 'deny',
          reason: result.reason || 'Command execution blocked by security policy.',
          systemMessage: '🔒 Security Block: Restricted shell command denied.',
        }),
      );
      process.exit(0);
    }
  }

  const fileModificationTools = ['write_file', 'replace', 'edit_file', 'create_file'];
  if (fileModificationTools.includes(tool_name) && tool_input) {
    const targetPath = tool_input.file_path || tool_input.path || '';
    if (targetPath.endsWith('eslint.config.mjs')) {
      console.log(
        JSON.stringify({
          decision: 'deny',
          reason: 'Direct modification of eslint.config.mjs is restricted. If you need to change linting rules, you must use the ask_user tool to present the proposed changes and request that the developer apply them manually.',
          systemMessage: '🔒 Security Block: Modifying ESLint configuration is denied.',
        }),
      );
      process.exit(0);
    }
  }

  let allowedMessage = 'execution allowed.';
  if (tool_name === 'run_shell_command' && tool_input && tool_input.command) {
    const cmdStr = tool_input.command.length > 100 ? tool_input.command.substring(0, 97) + '...' : tool_input.command;
    allowedMessage = `\`${cmdStr}\` command execution allowed.`;
  } else if (tool_name) {
    allowedMessage = `\`${tool_name}\` execution allowed.`;
  }

  console.log(JSON.stringify({ decision: 'allow', systemMessage: `🔒 Hook Notification: ${allowedMessage}` }));
  process.exit(0);
}

main();
