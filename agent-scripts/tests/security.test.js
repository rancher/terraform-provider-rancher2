import test from 'node:test';
import assert from 'node:assert';
import fs from 'fs';
import path from 'path';
import crypto from 'crypto';
import { execSync } from 'child_process';
import { verifyGitCommand } from '../security.js';

test('security.js: verifyGitCommand unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  const tempHome = path.resolve(`/tmp/gemini-security-test-${uniqueId}`);
  fs.mkdirSync(tempHome, { recursive: true });

  execSync('git init', { cwd: tempHome, stdio: 'ignore' });

  t.after(() => {
    fs.rmSync(tempHome, { recursive: true, force: true });
  });

  await t.test('allows normal non-git shell commands', () => {
    const result = verifyGitCommand('echo "hello"', tempHome);
    assert.strictEqual(result.decision, 'allow');
  });

  await t.test('denies manual execution of enforcer hooks and agent-scripts', () => {
    const resultGemini = verifyGitCommand('node .gemini/hooks/block-secrets.js', tempHome);
    assert.strictEqual(resultGemini.decision, 'deny');
    assert.ok(
      resultGemini.reason.includes('Manual execution of enforcer hook or agent scripts is strictly prohibited'),
    );

    const resultClaude = verifyGitCommand('node .claude/hooks/block-direct-git.js', tempHome);
    assert.strictEqual(resultClaude.decision, 'deny');

    const resultAgentScripts = verifyGitCommand('sh agent-scripts/git-utils.sh', tempHome);
    assert.strictEqual(resultAgentScripts.decision, 'deny');
  });

  await t.test('denies direct manipulation of gate approval and challenge files', () => {
    const resultEcho = verifyGitCommand('echo "approved" > review-approval.json', tempHome);
    assert.strictEqual(resultEcho.decision, 'deny');
    assert.ok(
      resultEcho.reason.includes(
        'Manually writing, editing, or spoofing any planning, testing, review, or commit gate approval files is strictly prohibited',
      ),
    );

    const resultPlan = verifyGitCommand('cat /tmp/spoof > plan-approval.challenge', tempHome);
    assert.strictEqual(resultPlan.decision, 'deny');

    const resultTest = verifyGitCommand('touch test-approval.json', tempHome);
    assert.strictEqual(resultTest.decision, 'deny');

    const resultCommit = verifyGitCommand('rm user-approval.age', tempHome);
    assert.strictEqual(resultCommit.decision, 'deny');
  });

  await t.test('denies direct git commit or git push', () => {
    const resultCommit = verifyGitCommand('git commit -m "feat: bypass"', tempHome);
    assert.strictEqual(resultCommit.decision, 'deny');
    assert.ok(resultCommit.reason.includes('Direct manual git commit and push commands are strictly prohibited'));

    const resultPush = verifyGitCommand('git push origin main', tempHome);
    assert.strictEqual(resultPush.decision, 'deny');
  });
});
