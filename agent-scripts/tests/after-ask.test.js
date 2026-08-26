import { execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import assert from 'node:assert';
import test from 'node:test';
import os from 'os';
import path from 'path';
import { handlePlanApproval } from '../after-ask.js';

test('after-ask.js: Touch ID gate signing unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  // Resolve tempHome inside homedir so that ssh-keygen verification won't fail due to world-writable /tmp permissions
  const tempHome = path.resolve(os.homedir(), `.gemini/tmp/gemini-after-ask-test-${uniqueId}`);
  const tempTmpDir = path.resolve(tempHome, '.gemini/tmp/terraform-provider-rancher2');

  // Set HOME environment variable so that os.homedir() returns tempHome in subsequent gating checks
  process.env.HOME = tempHome;

  fs.mkdirSync(tempTmpDir, { recursive: true });

  t.after(() => {
    fs.rmSync(tempHome, { recursive: true, force: true });
  });

  await t.test('handlePlanApproval signs Gate 1 with SSH keys', () => {
    const keysDir = path.join(tempHome, '.gemini');
    fs.mkdirSync(keysDir, { recursive: true });
    const privKeyFile = path.join(keysDir, 'ssh-key');
    const pubKeyFile = path.join(keysDir, 'ssh-key.pub');

    // Generate real passwordless SSH key pair
    execSync(`ssh-keygen -t ed25519 -C "gemini" -N "" -f "${privKeyFile}"`, { stdio: 'ignore' });

    const sessionPlansDir = path.join(tempTmpDir, 'session1/plans');
    fs.mkdirSync(sessionPlansDir, { recursive: true });
    const planFile = path.join(sessionPlansDir, 'PR404-Resolution.md');

    // Ensure the plan text satisfies validation: has checklist, tests, gates, framework, doc checks
    const compliantText =
      `# PR404 Resolution\n` +
      `- [ ] Task 1: Create a script.\n` +
      `- [ ] Run comprehensive tests and linters.\n` +
      `- [ ] Maintain the Agentic Framework if improvements are found.\n` +
      `- [ ] Enforce standard quality gates.\n` +
      `- [ ] Update documentation to describe the changes.`;
    fs.writeFileSync(planFile, compliantText);

    const result = handlePlanApproval(
      tempTmpDir,
      pubKeyFile,
      'Do you approve the plan?\n```markdown\n# Bootstrap Test\n```',
    );

    assert.strictEqual(result.status, 'approved');
    const planApprovalFile = path.join(tempTmpDir, 'plan-approval.json');
    assert.strictEqual(fs.existsSync(planApprovalFile), true);
    const approval = JSON.parse(fs.readFileSync(planApprovalFile, 'utf-8'));
    assert.strictEqual(approval.status, 'approved');
  });
});
