import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import path from 'node:path';

const SCRIPT_PATH = path.resolve('.github/workflows/scripts/goreleaser.sh');

describe('goreleaser.sh', () => {
  it('passes bash syntax check (bash -n)', () => {
    assert.doesNotThrow(() => {
      execFileSync('bash', ['-n', SCRIPT_PATH], { stdio: 'pipe' });
    });
  });

  it('fails if GPG_PASSPHRASE is empty', () => {
    assert.throws(
      () => {
        execFileSync('bash', [SCRIPT_PATH], {
          env: { PATH: process.env.PATH },
          stdio: 'pipe',
        });
      },
      (error) => {
        const stderr = error.stderr?.toString() || '';
        return stderr.includes('Error: gpg passphrase empty');
      }
    );
  });

  it('fails if GPG_KEY_ID is empty', () => {
    assert.throws(
      () => {
        execFileSync('bash', [SCRIPT_PATH], {
          env: { PATH: process.env.PATH, GPG_PASSPHRASE: 'secret' },
          stdio: 'pipe',
        });
      },
      (error) => {
        const stderr = error.stderr?.toString() || '';
        return stderr.includes('Error: key id empty');
      }
    );
  });

  it('fails if GPG_KEY is empty', () => {
    assert.throws(
      () => {
        execFileSync('bash', [SCRIPT_PATH], {
          env: {
            PATH: process.env.PATH,
            GPG_PASSPHRASE: 'secret',
            GPG_KEY_ID: '12345678',
          },
          stdio: 'pipe',
        });
      },
      (error) => {
        const stderr = error.stderr?.toString() || '';
        return stderr.includes('Error: key contents empty');
      }
    );
  });
});
