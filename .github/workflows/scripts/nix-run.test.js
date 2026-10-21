import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';

const SCRIPT_PATH = path.resolve('.github/workflows/scripts/nix-run.sh');

describe('nix-run.sh', () => {
  it('passes bash syntax check (bash -n)', () => {
    assert.doesNotThrow(() => {
      execFileSync('bash', ['-n', SCRIPT_PATH], { stdio: 'pipe' });
    });
  });

  it('executes wrapped command successfully when nix runner is available', () => {
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'nix-mock-'));
    const mockSudo = path.join(tmpDir, 'sudo');
    fs.writeFileSync(
      mockSudo,
      `#!/usr/bin/env bash
while [[ $# -gt 0 ]]; do
  if [[ "$1" == "--command" ]]; then
    shift
    exec "$@"
  fi
  shift
done
`,
      { mode: 0o755 }
    );

    try {
      const output = execFileSync('bash', [SCRIPT_PATH, 'echo', 'runner active'], {
        env: { ...process.env, PATH: `${tmpDir}:${process.env.PATH}` },
        stdio: 'pipe',
        encoding: 'utf-8',
      });
      assert.ok(output.includes('runner active'));
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });

  it('fails deterministically when wrapped command exits with non-zero code', () => {
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'nix-mock-'));
    const mockSudo = path.join(tmpDir, 'sudo');
    fs.writeFileSync(
      mockSudo,
      `#!/usr/bin/env bash
while [[ $# -gt 0 ]]; do
  if [[ "$1" == "--command" ]]; then
    shift
    exec "$@"
  fi
  shift
done
`,
      { mode: 0o755 }
    );

    try {
      assert.throws(
        () => {
          execFileSync('bash', [SCRIPT_PATH, 'bash', '-c', 'exit 42'], {
            env: { ...process.env, PATH: `${tmpDir}:${process.env.PATH}` },
            stdio: 'pipe',
          });
        },
        (error) => {
          return error.status === 42;
        }
      );
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });

  it('fails safely when runner tools fail or are unavailable', () => {
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'nix-mock-'));
    const mockSudo = path.join(tmpDir, 'sudo');
    fs.writeFileSync(
      mockSudo,
      `#!/usr/bin/env bash
echo "sudo: command not found" >&2
exit 127
`,
      { mode: 0o755 }
    );

    try {
      assert.throws(
        () => {
          execFileSync('bash', [SCRIPT_PATH, 'echo', 'hello'], {
            env: { ...process.env, PATH: `${tmpDir}:${process.env.PATH}` },
            stdio: 'pipe',
          });
        },
        (error) => {
          return error.status !== 0;
        }
      );
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });
});
