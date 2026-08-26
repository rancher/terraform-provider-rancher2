import { execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import assert from 'node:assert';
import test from 'node:test';
import os from 'os';
import path from 'path';
import { calculateFileHash, findLatestActivePlan, verifyPlanGate } from '../gating.js';

test('gating.js: verification unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  // Resolve tempHome inside homedir so that ssh-keygen verification won't fail due to world-writable /tmp permissions
  const tempHome = path.resolve(os.homedir(), `.gemini/tmp/gemini-gating-test-${uniqueId}`);
  const tempTmpDir = path.resolve(tempHome, '.gemini/tmp/terraform-provider-rancher2');

  // Set HOME environment variable so that os.homedir() returns tempHome in subsequent gating checks
  process.env.HOME = tempHome;

  fs.mkdirSync(tempTmpDir, { recursive: true });
  execSync('git init', { cwd: tempHome, stdio: 'ignore' });

  t.after(() => {
    fs.rmSync(tempHome, { recursive: true, force: true });
  });

  await t.test('calculateFileHash hashes file content correctly', () => {
    const file = path.join(tempHome, 'test.txt');
    fs.writeFileSync(file, 'hello world');
    const hash = calculateFileHash(file);
    const expected = crypto.createHash('sha256').update('hello world').digest('hex');
    assert.strictEqual(hash, expected);
  });

  await t.test('findLatestActivePlan returns latest plan file', () => {
    const session1 = path.join(tempTmpDir, 'session1/plans');
    const session2 = path.join(tempTmpDir, 'session2/plans');
    fs.mkdirSync(session1, { recursive: true });
    fs.mkdirSync(session2, { recursive: true });

    const plan1 = path.join(session1, 'Plan1.md');
    const plan2 = path.join(session2, 'Plan2.md');

    fs.writeFileSync(plan1, '# Plan 1');
    // Ensure separate mtimes
    execSync('touch -d "2 hours ago" ' + plan1);
    fs.writeFileSync(plan2, '# Plan 2');

    const latest = findLatestActivePlan(tempTmpDir);
    assert.strictEqual(latest, plan2);
  });

  await t.test('verifyPlanGate checks signatures correctly', () => {
    // Create keys dir and real SSH key pair
    const keysDir = path.join(tempHome, '.gemini');
    fs.mkdirSync(keysDir, { recursive: true });
    const privKeyFile = path.join(keysDir, 'ssh-key');
    execSync(`ssh-keygen -t ed25519 -C "gemini" -N "" -f "${privKeyFile}"`, { stdio: 'ignore' });

    // Create session and plan
    const session1 = path.join(tempTmpDir, 'session1/plans');
    fs.mkdirSync(session1, { recursive: true });
    const plan1 = path.join(session1, 'Plan1.md');
    fs.writeFileSync(plan1, '# Plan 1');
    const planHash = crypto.createHash('sha256').update('# Plan 1').digest('hex');

    // Create plan approval and sign it natively
    const approvalFile = path.join(tempTmpDir, 'plan-approval.json');
    fs.writeFileSync(
      approvalFile,
      JSON.stringify({
        status: 'approved',
        plan_file: plan1,
        plan_hash: planHash,
      }),
    );
    execSync(`ssh-keygen -Y sign -f "${privKeyFile}" -n gemini "${approvalFile}"`, { stdio: 'ignore' });

    const verifiedHash = verifyPlanGate(tempTmpDir);
    assert.strictEqual(verifiedHash, planHash);
  });
});
