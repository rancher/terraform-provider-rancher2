import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import path from 'node:path';

const SCRIPT_PATH = path.resolve('.github/workflows/scripts/nix-run.sh');

describe('nix-run.sh', () => {
  it('passes bash syntax check (bash -n)', () => {
    assert.doesNotThrow(() => {
      execFileSync('bash', ['-n', SCRIPT_PATH], { stdio: 'pipe' });
    });
  });

  it('fails safely when sudo or nix is missing outside the Nix container', () => {
    assert.throws(
      () => {
        execFileSync('bash', [SCRIPT_PATH, 'echo', 'hello'], {
          stdio: 'pipe',
        });
      },
      (error) => {
        // Exits with code > 0 on host system where suse/sudo/nix develop cannot drop privileges
        return error.status !== 0;
      }
    );
  });
});
