import { execFileSync, execSync } from 'child_process';
import fs from 'fs';
import path from 'path';

/**
 * Validates a shell command against Rancher git remote operations, manual git commits/pushes, and branch draft PR checks.
 * @returns {object} - { decision: 'allow'|'deny', reason?: string, systemMessage?: string }
 */
export function verifyGitCommand(command, cwd) {
  const trimmedCmd = command.trim();

  // Strip leading env var assignments (e.g. KEY=value or KEY="value" or KEY='value') and optional sudo
  let commandClean = trimmedCmd;
  while (true) {
    const next = commandClean.replace(/^[A-Za-z_][A-Za-z0-9_]*=(?:'[^']*'|"[^"]*"|\S+)\s+/, '');
    if (next === commandClean) {
      break;
    }
    commandClean = next;
  }

  // Anti-Bypass Guardrail: Unconditionally deny any manual execution of enforcer hook scripts inside .gemini/hooks/ or .claude/hooks/
  const isExecutingHooksManually =
    trimmedCmd.includes('.gemini/hooks/') ||
    trimmedCmd.includes('.gemini/hooks') ||
    trimmedCmd.includes('.claude/hooks/') ||
    trimmedCmd.includes('.claude/hooks') ||
    trimmedCmd.includes('agent-scripts/');
  if (isExecutingHooksManually) {
    return {
      decision: 'deny',
      reason:
        '🔒 Security Policy Violation: Manual execution of enforcer hook or agent scripts is strictly prohibited.\n\n' +
        'These scripts are part of the secure system pipeline and must only be executed automatically by the Gemini CLI lifecycle.',
      systemMessage: '🔒 Security Block: Manual execution of secure scripts is prohibited.',
    };
  }

  // Anti-Bypass Guardrail: Unconditionally deny any manual writing, editing, or spoofing of any gate approval/challenge JSON/age files
  const isManipulatingApproval =
    /\b(echo|cat|touch|rm|mv|cp|write|tee|vim|vi|nano|printf|sed|awk)\b.*\b(plan-approval|test-approval|review-approval|user-approval)\.(json|challenge|age|sig)\b|>>?[^>]*\b(plan-approval|test-approval|review-approval|user-approval)\.(json|challenge|age|sig)\b/.test(
      commandClean,
    );
  if (isManipulatingApproval) {
    return {
      decision: 'deny',
      reason:
        'Security Policy Violation: Manually writing, editing, or spoofing any planning, testing, review, or commit gate approval files is strictly prohibited.\n\n' +
        'Gating approval files must ONLY be generated automatically and securely by our pipeline hooks and sub-agents.\n\n' +
        'To proceed:\n' +
        '1. Comply strictly with our gated sequence (Plan -> Test -> Review -> Commit).\n' +
        '2. Use the proper tools (biometric Touches or sub-agent runs) to obtain valid signatures.\n' +
        '3. Never attempt to manually create, edit, or spoof any gate approval or challenge files.',
      systemMessage: '🔒 Security Block: Direct manipulation of approval files is prohibited.',
    };
  }

  // Check if we are attempting to switch branches while current PR is in Draft mode (Phase 6, Step 18 / Phase 7, Step 20)
  let isBranchSwitch = false;
  if (/\bgit\s+switch\b/.test(commandClean)) {
    isBranchSwitch = true;
  } else if (/\bgit\s+checkout\b/.test(commandClean)) {
    const hasDoubleDash = commandClean.includes(' -- ');
    if (hasDoubleDash) {
      isBranchSwitch = false;
    } else {
      const parts = commandClean.split(/\s+/).filter((p) => p !== 'git' && p !== 'checkout');
      const nonFlagParts = parts.filter((p) => !p.startsWith('-') || p === '-');

      if (nonFlagParts.length > 0) {
        const target = nonFlagParts[0];
        if (target === '-') {
          isBranchSwitch = true;
        } else {
          const resolvedTarget = path.resolve(cwd || process.cwd(), target);
          if (!fs.existsSync(resolvedTarget)) {
            isBranchSwitch = true;
          }
        }
      } else {
        isBranchSwitch = false;
      }
    }
  }

  if (isBranchSwitch) {
    try {
      const currentBranch = execSync('git branch --show-current', {
        cwd: cwd || process.cwd(),
        stdio: ['ignore', 'pipe', 'ignore'],
      })
        .toString()
        .trim();

      if (currentBranch && currentBranch !== 'main') {
        let prStatusOutput;
        try {
          prStatusOutput = execFileSync('gh', ['pr', 'view', currentBranch, '--json', 'isDraft,number'], {
            cwd: cwd || process.cwd(),
            stdio: ['ignore', 'pipe', 'pipe'],
          })
            .toString()
            .trim();
        } catch (execErr) {
          const stderr = execErr.stderr ? execErr.stderr.toString() : '';
          if (stderr.includes('no pull requests found')) {
            prStatusOutput = '';
          } else {
            throw execErr;
          }
        }

        if (prStatusOutput) {
          const prInfo = JSON.parse(prStatusOutput);
          if (prInfo.isDraft === true) {
            return {
              decision: 'deny',
              reason:
                `Security Policy Violation: Moving to a new PR or branch is prohibited while the current branch PR (#${prInfo.number}) is still in Draft mode.\n\n` +
                `In accordance with Phase 6, Step 18 (Convert to Ready) and Phase 7, Step 20 (Proceed to Next Layer) of 'docs/development/AgenticFramework/DevelopmentProcess.md', you MUST first graduate the current PR from Draft to Ready-for-Review before checking out 'main' or switching tasks.\n\n` +
                `To proceed:\n` +
                `1. Complete all iteration reviews and obtain local sign-off.\n` +
                `2. Convert the draft PR to Ready-for-Review (Phase 6, Step 18) using: \`gh pr ready ${prInfo.number}\` (or the create-pr.sh skill).\n` +
                `3. Once the PR is marked as ready for review on GitHub, you will be authorized to switch branches (Phase 7, Step 20).`,
              systemMessage: `🔒 Security Block: Current PR #${prInfo.number} is in Draft mode. Please comply with Phase 6, Step 18 of docs/development/AgenticFramework/DevelopmentProcess.md.`,
            };
          }
        }
      }
    } catch (err) {
      console.error('Failed to verify branch PR status:', err);
      return {
        decision: 'deny',
        reason:
          'Security Policy Violation: Failed to verify draft PR status on GitHub. To prevent branch state divergence, operations are blocked until status can be verified.',
        systemMessage: '🔒 Security Block: Branch PR verification failed.',
      };
    }
  }

  // Check for unauthorized git commit or push operations
  const isCommitOrPush = /\bgit\s+(commit|push)\b/.test(commandClean);
  if (isCommitOrPush) {
    return {
      decision: 'deny',
      reason:
        `Security Policy Violation: Direct manual git commit and push commands are strictly prohibited in this repository.\n\n` +
        `In accordance with Phase 6, Step 15 (Authorized Commit & Push) of 'docs/development/AgenticFramework/DevelopmentProcess.md', you MUST use our custom, secure, and synchronized commit-and-push skill to perform commits and pushes.\n\n` +
        `This skill automatically validates file count limits, verifies proactive code reviews, synchronizes with upstream main, pulls/fetches from your fork remote, and executes GPG/SSH cryptographically signed and signed-off commits.\n\n` +
        `To proceed:\n` +
        `1. Stage your changes cleanly: \`git add <files>...\`\n` +
        `2. Execute the commit-push skill in the chat: \`.gemini/skills/commit-push.sh -m "your commit message"\`\n` +
        `3. Provide manual confirmation and allow the skill to execute fully.`,
      systemMessage:
        '🔒 Security Block: Direct git commit/push is blocked. Please run \'.gemini/skills/commit-push.sh -m "..."\' to commit.',
    };
  }

  // Check if it is a git command and performs a remote-interacting operation
  const isGitCmd = /^(?:sudo\s+)?git\b/.test(commandClean);
  const isRemoteOp = /\b(push|pull|fetch|clone|remote)\b/.test(commandClean);

  if (isGitCmd && isRemoteOp) {
    const hasRancherRef = /rancher/i.test(trimmedCmd.replace(/block-restricted-commands\.js/g, ''));
    if (hasRancherRef) {
      return {
        decision: 'deny',
        reason:
          'Security Policy Violation: Git command contains references to Rancher remote/URLs, which is strictly blocked.',
        systemMessage: '🔒 Security Block: Prohibited remote/URL reference detected.',
      };
    }

    try {
      const remotesOutput = execSync('git remote -v', {
        cwd: cwd || process.cwd(),
        stdio: ['ignore', 'pipe', 'ignore'],
      }).toString();

      if (/rancher/i.test(remotesOutput)) {
        return {
          decision: 'deny',
          reason:
            'Security Policy Violation: Operations (push, pull, fetch, remote) targeting Rancher-owned remotes are strictly blocked.',
          systemMessage: '🔒 Security Block: Git remote operation against a Rancher remote is prohibited.',
        };
      }
    } catch (err) {
      console.error('Failed to check git remote safety:', err);
      return {
        decision: 'deny',
        reason: 'Security Policy Violation: Failed to check git remote safety configuration.',
        systemMessage: '🔒 Security Block: Remote safety verification failed.',
      };
    }
  }

  return { decision: 'allow' };
}
