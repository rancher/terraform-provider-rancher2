import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import {
  VALID_PREFIXES,
  validateEmpty,
  validatePrefix,
  validateLength,
  validateSpelling,
  validateMessage,
  getCommitMessages,
  main,
} from './validate-commit-message.js';

describe('validate-commit-message', () => {
  describe('VALID_PREFIXES', () => {
    it('contains expected conventional commit prefixes', () => {
      assert.ok(VALID_PREFIXES.includes('fix: '));
      assert.ok(VALID_PREFIXES.includes('feat: '));
      assert.ok(VALID_PREFIXES.includes('feature: '));
      assert.ok(VALID_PREFIXES.includes('refactor!: '));
      assert.ok(VALID_PREFIXES.includes('feature!: '));
      assert.ok(VALID_PREFIXES.includes('feat!: '));
      assert.ok(VALID_PREFIXES.includes('chore(main): '));
      assert.ok(VALID_PREFIXES.includes('Merge branch '));
    });
  });

  describe('validateEmpty', () => {
    it('throws error for empty or whitespace messages', () => {
      assert.throws(() => validateEmpty(''), /...Empty commit message\./);
      assert.throws(() => validateEmpty('   '), /...Empty commit message\./);
      assert.throws(() => validateEmpty(null), /...Empty commit message\./);
    });

    it('returns true for non-empty message', () => {
      assert.strictEqual(validateEmpty('fix: bug'), true);
    });
  });

  describe('validatePrefix', () => {
    it('accepts allowed prefixes', () => {
      assert.strictEqual(validatePrefix('fix: resolve issue'), true);
      assert.strictEqual(validatePrefix('feat: new capability'), true);
      assert.strictEqual(validatePrefix('feature: new capability'), true);
      assert.strictEqual(validatePrefix('refactor!: breaking change'), true);
      assert.strictEqual(validatePrefix('feature!: breaking change'), true);
      assert.strictEqual(validatePrefix('feat!: breaking change'), true);
      assert.strictEqual(validatePrefix('chore(main): update dependencies'), true);
      assert.strictEqual(validatePrefix('Merge branch main into dev'), true);
    });

    it('rejects disallowed prefixes', () => {
      assert.throws(
        () => validatePrefix('docs: update readme'),
        /\.\.\.Commit message does not start with the required prefix\./
      );
      assert.throws(
        () => validatePrefix('random: unsupported commit'),
        /\.\.\.Commit message does not start with the required prefix\./
      );
    });
  });

  describe('validateLength', () => {
    it('accepts messages 100 characters or fewer', () => {
      assert.strictEqual(validateLength('feat: short subject line'), true);
      const exact100 = 'feat: ' + 'a'.repeat(94);
      assert.strictEqual(exact100.length, 100);
      assert.strictEqual(validateLength(exact100), true);
    });

    it('rejects messages longer than 100 characters', () => {
      const over100 = 'feat: ' + 'a'.repeat(95);
      assert.strictEqual(over100.length, 101);
      assert.throws(
        () => validateLength(over100),
        /\.\.\.Commit message subject line should be less than 100 characters, found 101\./
      );
    });
  });

  describe('validateSpelling', () => {
    it('skips spell check for merge commits', () => {
      const mockExec = () => {
        throw new Error('Should not be called for merge commit');
      };
      assert.strictEqual(validateSpelling('Merge branch foo into bar', mockExec), true);
      assert.strictEqual(validateSpelling('Merge pull request #123', mockExec), true);
    });

    it('passes when cspell exits cleanly', () => {
      let executedCommand = null;
      let executedInput = null;
      const mockExec = (cmd, opts) => {
        executedCommand = cmd;
        executedInput = opts?.input;
        return '';
      };
      assert.strictEqual(validateSpelling('feat: proper spelling', mockExec), true);
      assert.strictEqual(executedCommand, 'cspell stdin');
      assert.strictEqual(executedInput, 'feat: proper spelling');
    });

    it('throws error when cspell detects spelling error', () => {
      const mockExec = () => {
        const err = new Error('Command failed: cspell stdin');
        err.stdout = 'Unknown word (invaliddword)';
        throw err;
      };
      assert.throws(
        () => validateSpelling('feat: invaliddword included', mockExec),
        /\.\.\.Commit message contains spelling errors on: \^Unknown word \(invaliddword\)\$/
      );
    });
  });

  describe('validateMessage', () => {
    it('validates a completely valid commit message', () => {
      const mockExec = () => '';
      assert.strictEqual(validateMessage('feat: add support for cluster proxy', mockExec), true);
    });

    it('fails if prefix is invalid', () => {
      const mockExec = () => '';
      assert.throws(() => validateMessage('invalid: prefix message', mockExec));
    });
  });

  describe('getCommitMessages', () => {
    it('falls back to git log when prNumber is not provided', () => {
      const mockExec = (cmd) => {
        assert.ok(cmd.includes('git log origin/main..HEAD'));
        return 'feat: commit one\nfix: commit two\n';
      };
      const messages = getCommitMessages(undefined, mockExec);
      assert.deepStrictEqual(messages, ['feat: commit one', 'fix: commit two']);
    });

    it('calls gh api when prNumber is provided', () => {
      const mockExec = (cmd) => {
        assert.ok(cmd.includes('pulls/42/commits'));
        return 'feat: api commit one\nfix: api commit two\n';
      };
      const messages = getCommitMessages('42', mockExec);
      assert.deepStrictEqual(messages, ['feat: api commit one', 'fix: api commit two']);
    });
  });

  describe('main', () => {
    it('runs validation loop across all commits', async () => {
      const mockExec = (cmd) => {
        if (cmd.startsWith('gh api')) {
          return 'feat: first commit\nfix: second commit\n';
        }
        if (cmd.startsWith('cspell')) {
          return '';
        }
        return '';
      };
      await assert.doesNotReject(async () => {
        await main({ prNumber: '100', exec: mockExec });
      });
    });

    it('handles empty commit list without error', async () => {
      const mockExec = () => '';
      await assert.doesNotReject(async () => {
        await main({ prNumber: '100', exec: mockExec });
      });
    });
  });
});
