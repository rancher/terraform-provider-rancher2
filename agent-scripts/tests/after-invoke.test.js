import crypto from 'crypto';
import fs from 'fs';
import assert from 'node:assert';
import test from 'node:test';
import path from 'path';
import { saveReport, verifyTestReport } from '../after-invoke.js';

test('after-invoke.js: subagent report parsing unit tests', async (t) => {
  const uniqueId = crypto.randomBytes(8).toString('hex');
  const tempHome = path.resolve(`/tmp/gemini-after-invoke-test-${uniqueId}`);
  const tempTmpDir = path.resolve(tempHome, '.gemini/tmp/terraform-provider-rancher2');

  fs.mkdirSync(tempTmpDir, { recursive: true });

  t.after(() => {
    fs.rmSync(tempHome, { recursive: true, force: true });
  });

  await t.test('saveReport writes report markdown correctly', () => {
    const logsDir = path.join(tempTmpDir, 'logs');
    saveReport('testing-agent', 'Standard Report Content', logsDir);

    const reportFile = path.join(logsDir, 'testing-agent_report.md');
    assert.strictEqual(fs.existsSync(reportFile), true);
    assert.strictEqual(fs.readFileSync(reportFile, 'utf-8'), 'Standard Report Content');
  });

  await t.test('verifyTestReport signs Gate 2 when report indicates SUCCESS', () => {
    const stateFile = path.join(tempTmpDir, 'phase-state.json');
    const result = verifyTestReport(
      'TEST RUN status: 🟢 SUCCESS - All tests and linting passed.',
      'diff123',
      'plan123',
      stateFile,
    );

    assert.strictEqual(result.status, 'approved');
    assert.strictEqual(fs.existsSync(stateFile), true);
    const content = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
    assert.strictEqual(content.tested_diff_hash, 'diff123');
    assert.strictEqual(content.tested_plan_hash, 'plan123');
  });

  await t.test('verifyTestReport revokes Gate 2 when report indicates FAILURE', () => {
    const stateFile = path.join(tempTmpDir, 'phase-state.json');
    // Pre-populate with approved signature
    fs.writeFileSync(stateFile, JSON.stringify({ tested_diff_hash: 'diff123' }));

    const result = verifyTestReport('TEST RUN status: 🔴 FAILED', 'diff123', 'plan123', stateFile);

    assert.strictEqual(result.status, 'rejected');
    const content = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
    assert.strictEqual(content.tested_diff_hash, '');
    assert.strictEqual(content.tested_plan_hash, '');
  });
});
