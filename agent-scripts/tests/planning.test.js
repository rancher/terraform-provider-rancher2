import test from 'node:test';
import assert from 'node:assert';
import fs from 'fs';
import path from 'path';
import crypto from 'crypto';
import { execSync } from 'child_process';
import { checkActivePlan } from '../planning.js';

test('planning.js: checkActivePlan unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  const tempHome = path.resolve(`/tmp/gemini-planning-test-${uniqueId}`);
  const tempTmpDir = path.resolve(tempHome, '.gemini/tmp');

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

  await t.test('returns false when working tree is clean (no modified plan)', () => {
    const hasPlan = checkActivePlan(tempHome);
    assert.strictEqual(hasPlan, false);
  });

  await t.test('returns true when a plan is added or modified under docs/development/', () => {
    fs.writeFileSync(path.join(tempHome, 'docs/development/TestPlan.md'), '# Test Plan');
    const hasPlan = checkActivePlan(tempHome);
    assert.strictEqual(hasPlan, true);
  });
});
