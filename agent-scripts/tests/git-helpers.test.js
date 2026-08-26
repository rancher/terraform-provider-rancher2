import { execSync } from 'child_process';
import crypto from 'crypto';
import fs from 'fs';
import assert from 'node:assert';
import test from 'node:test';
import os from 'os';
import path from 'path';
import {
  calculateFileHash,
  checkAndRevokeStaleGates,
  findLatestActivePlan,
  verifyPlanGate,
  verifyReviewGate,
  verifyTestGate,
} from '../gating.js';
import {
  calculateSha256,
  getFileOwnerUid,
  getForkOwner,
  getPhaseState,
  runCommitPushHelper,
  runCreatePrHelper,
  runPhaseManager,
  savePhaseState,
  verifyPushSafety,
} from '../git-helpers.js';
import { checkActiveBlueprint, validatePlanContent } from '../planning.js';

test('git-helpers.js: master consolidated helper unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  const tempHome = path.resolve(os.homedir(), `.gemini/tmp/gemini-helpers-test-${uniqueId}`);
  process.env.HOME = tempHome;
  const tempTmpDir = path.resolve(tempHome, '.gemini/tmp/terraform-provider-rancher2');

  fs.mkdirSync(tempTmpDir, { recursive: true });
  execSync('git init', { cwd: tempHome, stdio: 'ignore' });
  fs.mkdirSync(path.join(tempHome, 'docs/development'), { recursive: true });
  fs.writeFileSync(path.join(tempHome, 'docs/development/.gitkeep'), '');
  execSync('git add -A && git -c user.email=test@test.com -c user.name=test commit -m init', {
    cwd: tempHome,
    stdio: 'ignore',
  });

  t.after(() => {
    fs.rmSync(tempHome, { recursive: true, force: true });
  });

  // --- HASH & UTILS TESTS ---
  await t.test('calculateSha256 calculates sha256 correctly', () => {
    const raw = 'my test data';
    const expected = crypto.createHash('sha256').update(raw).digest('hex');
    assert.strictEqual(calculateSha256(raw), expected);
  });

  await t.test('calculateFileHash hashes file content correctly', () => {
    const file = path.join(tempHome, 'test.txt');
    fs.writeFileSync(file, 'hello world');
    const hash = calculateFileHash(file);
    const expected = crypto.createHash('sha256').update('hello world').digest('hex');
    assert.strictEqual(hash, expected);
  });

  await t.test('getFileOwnerUid retrieves ownership correctly', () => {
    const filePath = path.join(tempHome, 'test.txt');
    fs.writeFileSync(filePath, 'hello');
    const expectedUid = fs.statSync(filePath).uid;
    assert.strictEqual(getFileOwnerUid(filePath), expectedUid);
  });

  await t.test('getFileOwnerUid returns null for non-existent files', () => {
    assert.strictEqual(getFileOwnerUid(path.join(tempHome, 'nonexistent.txt')), null);
  });

  // --- BLUEPRINTS & PLANNING TESTS ---
  await t.test('checkActiveBlueprint returns false when working tree is clean', () => {
    const hasPlan = checkActiveBlueprint(tempHome);
    assert.strictEqual(hasPlan, false);
  });

  await t.test('checkActiveBlueprint returns true when a blueprint is added or modified', () => {
    fs.mkdirSync(path.join(tempHome, 'docs/development'), { recursive: true });
    fs.writeFileSync(path.join(tempHome, 'docs/development/TestBlueprint.md'), '# Test Blueprint');
    const hasPlan = checkActiveBlueprint(tempHome);
    assert.strictEqual(hasPlan, true);
  });

  await t.test('validatePlanContent rejects structurally incomplete plans', () => {
    const planPath = path.join(tempHome, 'incomplete-plan.md');
    fs.writeFileSync(planPath, '# Incomplete Plan\n- Create something.');

    const validation = validatePlanContent(planPath);
    assert.strictEqual(validation.valid, false);
    assert.strictEqual(validation.errors.length, 5); // missing all 5 required elements
  });

  await t.test('validatePlanContent approves compliant plans', () => {
    const planPath = path.join(tempHome, 'complete-plan.md');
    const compliantText =
      `# Complete Plan\n` +
      `- [ ] Task 1: Create a script.\n` +
      `- [ ] Run comprehensive tests and linters.\n` +
      `- [ ] Maintain the Agentic Framework if improvements are found.\n` +
      `- [ ] Enforce standard quality gates.\n` +
      `- [ ] Update documentation to describe the changes.`;
    fs.writeFileSync(planPath, compliantText);

    const validation = validatePlanContent(planPath);
    assert.strictEqual(validation.valid, true);
  });

  // --- GATING & PR TESTS ---
  await t.test('findLatestActivePlan returns latest plan file', () => {
    const session1 = path.join(tempTmpDir, 'session1/plans');
    const session2 = path.join(tempTmpDir, 'session2/plans');
    fs.mkdirSync(session1, { recursive: true });
    fs.mkdirSync(session2, { recursive: true });

    const plan1 = path.join(session1, 'Plan1.md');
    const plan2 = path.join(session2, 'Plan2.md');

    fs.writeFileSync(plan1, '# Plan 1');
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

  await t.test('verifyTestGate verifies testing approvals correctly', () => {
    const stateFile = path.join(tempTmpDir, 'phase-state.json');
    fs.writeFileSync(
      stateFile,
      JSON.stringify({
        currentPhase: 'review',
        tested_plan_hash: 'plan123',
        tested_diff_hash: 'diff123',
      }),
    );

    const isValid = verifyTestGate(tempTmpDir, 'plan123', 'diff123');
    assert.strictEqual(isValid, true);

    const isInvalidDiff = verifyTestGate(tempTmpDir, 'plan123', 'differentdiff');
    assert.strictEqual(isInvalidDiff, false);
  });

  await t.test('verifyReviewGate verifies review approvals correctly', () => {
    const reviewApprovalFile = path.join(tempTmpDir, 'review-approval.json');
    fs.writeFileSync(
      reviewApprovalFile,
      JSON.stringify({
        status: 'approved',
        plan_hash: 'plan123',
        diff_hash: 'diff123',
      }),
    );

    const isValid = verifyReviewGate(tempTmpDir, 'diff123', 'plan123');
    assert.strictEqual(isValid, true);

    const isInvalidPlan = verifyReviewGate(tempTmpDir, 'diff123', 'differentplan');
    assert.strictEqual(isInvalidPlan, false);
  });

  await t.test('checkAndRevokeStaleGates unlinks stale approvals when hashes mismatch', () => {
    const stateFile = path.join(tempTmpDir, 'phase-state.json');
    const reviewApprovalFile = path.join(tempTmpDir, 'review-approval.json');

    fs.writeFileSync(
      stateFile,
      JSON.stringify({ currentPhase: 'review', tested_plan_hash: 'plan123', tested_diff_hash: 'stalediff' }),
    );
    fs.writeFileSync(
      reviewApprovalFile,
      JSON.stringify({ status: 'approved', plan_hash: 'plan123', diff_hash: 'stalediff' }),
    );

    const revoked = checkAndRevokeStaleGates(tempTmpDir, 'currentdiff', 'plan123');
    assert.strictEqual(revoked, true);
    assert.strictEqual(fs.existsSync(reviewApprovalFile), false);

    const content = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
    assert.strictEqual(content.tested_diff_hash, '');
    assert.strictEqual(content.tested_plan_hash, '');
  });

  await t.test('getForkOwner parses origin remote url correctly', () => {
    execSync('git remote add origin https://github.com/my-fork/terraform-provider-rancher2.git', { cwd: tempHome });
    const owner = getForkOwner(tempHome);
    assert.strictEqual(owner, 'my-fork');
  });

  // --- SECURITY & PUSH TESTS ---
  await t.test('verifyPushSafety allows safe remotes and blocks unsafe remotes', () => {
    verifyPushSafety('origin', tempHome);

    execSync('git remote add rancherremote https://github.com/rancher/terraform-provider-rancher2', { cwd: tempHome });
    assert.throws(() => {
      verifyPushSafety('rancherremote', tempHome);
    }, /CRITICAL SECURITY ERROR: Unsafe push prevented!/);
  });

  // --- CLI RUNNERS TESTS ---
  await t.test('runCommitPushHelper rejects execution if message is missing', () => {
    assert.throws(() => {
      runCommitPushHelper(['-f'], tempHome);
    }, /Commit message is required/);
  });

  await t.test('runCommitPushHelper rejects execution if unknown option is provided', () => {
    assert.throws(() => {
      runCommitPushHelper(['-m', 'feat: clean commit', '--invalid-option'], tempHome);
    }, /Unknown argument/);
  });

  await t.test('runCreatePrHelper rejects execution if unknown option is provided', () => {
    assert.throws(() => {
      runCreatePrHelper(['--invalid-pr-option'], tempHome);
    }, /Unknown parameter/);
  });

  // --- PHASE STATE TESTS ---
  await t.test('getPhaseState returns research default if file is missing', () => {
    const nonExistentFile = path.join(tempHome, 'nonexistent-state.json');
    const state = getPhaseState(nonExistentFile);
    assert.strictEqual(state.currentPhase, 'research');
  });

  await t.test('savePhaseState and getPhaseState correctly write and read state', () => {
    const stateFile = path.join(tempHome, 'phase-state.json');
    const targetState = { currentPhase: 'implement' };
    savePhaseState(targetState, stateFile);

    const loaded = getPhaseState(stateFile);
    assert.strictEqual(loaded.currentPhase, 'implement');
  });

  await t.test('runPhaseManager rejects invalid phase names', () => {
    assert.throws(() => {
      runPhaseManager(['invalid_phase_name'], tempHome);
    }, /Invalid phase:/);
  });
});
