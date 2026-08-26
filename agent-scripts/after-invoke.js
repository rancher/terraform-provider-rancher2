import fs from 'fs';
import path from 'path';
import { execSync } from 'child_process';

/**
 * Saves a sub-agent execution report to disk.
 */
export function saveReport(agentName, report, logsDir) {
  try {
    fs.mkdirSync(logsDir, { recursive: true });
    const reportFile = path.join(logsDir, `${agentName}_report.md`);
    try {
      fs.unlinkSync(reportFile);
    } catch (err) {
      console.error(err.message || err);
    }
    fs.writeFileSync(reportFile, report, { mode: 0o600 });
  } catch (err) {
    console.error(`🔒 Hook Error: Failed to write sub-agent report for ${agentName}:`, err.message);
  }
}

/**
 * Verifies a test sub-agent report and writes the Gate 2 signature if successful.
 */
export function verifyTestReport(report, diffHash, planHash, stateFile) {
  const isSuccess = report.includes('TEST RUN status: 🟢 SUCCESS');

  if (isSuccess) {
    try {
      let state = { currentPhase: 'review' };
      if (fs.existsSync(stateFile)) {
        try {
          state = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
        } catch (err) {
          console.error(`🔒 Hook Warning: Failed to parse state JSON in verifyTestReport: ${err.message}`);
        }
      }
      state.tested_diff_hash = diffHash;
      state.tested_plan_hash = planHash;
      fs.writeFileSync(stateFile, JSON.stringify(state, null, 2));
      return {
        status: 'approved',
        systemMessage: '✅ Gate 2 Approved: Testing report verified.',
      };
    } catch (err) {
      console.error('🔒 Hook Error: Failed to write tested diff_hash:', err.message);
      return { status: 'error', error: err.message };
    }
  } else {
    try {
      if (fs.existsSync(stateFile)) {
        const state = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
        state.tested_diff_hash = '';
        state.tested_plan_hash = '';
        fs.writeFileSync(stateFile, JSON.stringify(state, null, 2));
      }
    } catch (err) {
      console.error(`🔒 Hook Warning: Failed to revoke tested state: ${err.message}`);
    }
    return {
      status: 'rejected',
      systemMessage: '❌ Gate 2 Rejected: Testing failures reported.',
    };
  }
}

/**
 * Verifies a review sub-agent report and writes the Gate 3 signature if successful.
 */
export function verifyReviewReport(report, diffHash, planHash, reviewApprovalFile, testApprovalFile) {
  const isSuccess = report.includes('PR Review status: 🟢 PERFECT - 0 findings.');

  if (isSuccess) {
    // Hook enforces Gate 2 must also be valid! (Review requires Tests to be passed)
    if (!fs.existsSync(testApprovalFile)) {
      return {
        status: 'gated',
        systemMessage:
          '🔒 Hook Notification: Review agent completed with 0 findings, but Gate 3 (Review Gate) cannot be signed because Gate 2 (Testing Gate) is missing!',
      };
    }

    try {
      const testContent = JSON.parse(fs.readFileSync(testApprovalFile, 'utf-8'));
      if (testContent.diff_hash !== diffHash || testContent.plan_hash !== planHash) {
        return {
          status: 'gated',
          systemMessage:
            '🔒 Hook Notification: Review agent completed with 0 findings, but Gate 3 cannot be signed because the current diff/plan does not match Gate 2 (Testing Gate).',
        };
      }

      const approvalData = {
        status: 'approved',
        message: 'PR Review status: 🟢 PERFECT - 0 findings.',
        commit_sha: execSync('git rev-parse HEAD 2>/dev/null || echo "unknown"', {
          stdio: ['ignore', 'pipe', 'ignore'],
        })
          .toString()
          .trim(),
        diff_hash: diffHash,
        plan_hash: planHash,
        timestamp: new Date().toISOString(),
      };

      try {
        fs.unlinkSync(reviewApprovalFile);
      } catch (err) {
        console.error(err.message || err);
      }
      fs.writeFileSync(reviewApprovalFile, JSON.stringify(approvalData, null, 2), { mode: 0o600 });
      return {
        status: 'approved',
        systemMessage:
          '✅ Gate 3 Approved: Review sub-agent report verified. Gate 3 signature successfully written and chained!',
      };
    } catch (err) {
      console.error('🔒 Hook Error: Failed to write Gate 3 signature:', err.message);
      return { status: 'error', error: err.message };
    }
  } else {
    // Self-Healing: Revoke existing signature if review failed
    try {
      fs.unlinkSync(reviewApprovalFile);
    } catch (err) {
      console.error(err.message || err);
    }
    return {
      status: 'rejected',
      systemMessage: '❌ Gate 3 Rejected: Review sub-agent reported violations. Gate 3 signature revoked/missing.',
    };
  }
}

/**
 * Programmatically validates that a testing report contains all of our strict verification requirements.
 * @param {string} report - The testing report markdown content
 * @returns {object} - { valid: boolean, errors: string[] }
 */
export function validateTestContent(report) {
  const errors = [];

  // 1. Linters check
  const lintMatch = /linter|lint|eslint|golangci-lint|tflint|static check/i.test(report);
  if (!lintMatch) {
    errors.push('Testing report must explicitly confirm successful static analysis and linter execution.');
  }

  // 2. Test execution check
  const testMatch = /unit test|test suite|go test|run_tests/i.test(report);
  if (!testMatch) {
    errors.push(
      'Testing report must explicitly confirm execution of the active unit tests or comprehensive test suites.',
    );
  }

  // 3. Status Check
  const statusMatch = /status:\s*🟢\s*SUCCESS|success|pass/i.test(report);
  if (!statusMatch) {
    errors.push('Testing report must contain an explicit success status declaration.');
  }

  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Programmatically validates that a review report contains all of our strict verification requirements.
 * @param {string} report - The review report markdown content
 * @returns {object} - { valid: boolean, errors: string[] }
 */
export function validateReviewContent(report) {
  const errors = [];

  // 1. Security Check
  const securityMatch = /security|vulnerab|secret/i.test(report);
  if (!securityMatch) {
    errors.push('Review report must explicitly confirm verification of security and credential safeguards.');
  }

  // 2. Coding Standards Check
  const standardsMatch = /standards|coding\s+standards|rules|conventions/i.test(report);
  if (!standardsMatch) {
    errors.push('Review report must explicitly confirm verification of repository coding standards.');
  }

  // 3. Spelling & Wording Check
  const spellingMatch = /spelling|wording|typo|discrepancy/i.test(report);
  if (!spellingMatch) {
    errors.push(
      'Review report must explicitly confirm verification of spelling and documentation/wording consistency.',
    );
  }

  // 4. Automation Audit Check
  const automationMatch = /automation\s+audit|automat(e|ed|ion)/i.test(report);
  if (!automationMatch) {
    errors.push('Review report must explicitly confirm conducting an automation audit of checked items.');
  }

  // 5. Approval Status Check
  const approvalMatch = /approved|approval|perfect|pass/i.test(report);
  if (!approvalMatch) {
    errors.push('Review report must contain an explicit approval status declaration.');
  }

  return {
    valid: errors.length === 0,
    errors,
  };
}
