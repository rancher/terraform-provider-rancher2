import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import {
  getChangedFiles,
  getCommitMessages,
  isProductFile,
  isProductPrefix,
  main,
  NON_PRODUCT_PREFIXES,
  PRODUCT_PREFIXES,
  VALID_PREFIXES,
  validateEmpty,
  validateLength,
  validateMessage,
  validatePrefix,
  validateProductImpact,
  validateSpelling,
} from './validate-commit-message.js';

describe('validate-commit-message', () => {
  describe('PRODUCT_PREFIXES and NON_PRODUCT_PREFIXES', () => {
    it('contains expected prefixes partitioned by impact', () => {
      assert.ok(PRODUCT_PREFIXES.includes('fix: '));
      assert.ok(PRODUCT_PREFIXES.includes('feat: '));
      assert.ok(PRODUCT_PREFIXES.includes('feature: '));
      assert.ok(PRODUCT_PREFIXES.includes('refactor: '));
      assert.ok(PRODUCT_PREFIXES.includes('refactor!: '));
      assert.ok(NON_PRODUCT_PREFIXES.includes('dev: '));
      assert.ok(NON_PRODUCT_PREFIXES.includes('env: '));
      assert.ok(NON_PRODUCT_PREFIXES.includes('docs: '));
      assert.ok(NON_PRODUCT_PREFIXES.includes('ci: '));
      assert.ok(NON_PRODUCT_PREFIXES.includes('test: '));
    });
  });

  describe('VALID_PREFIXES', () => {
    it('contains expected conventional commit prefixes', () => {
      assert.ok(VALID_PREFIXES.includes('fix: '));
      assert.ok(VALID_PREFIXES.includes('feat: '));
      assert.ok(VALID_PREFIXES.includes('feature: '));
      assert.ok(VALID_PREFIXES.includes('refactor: '));
      assert.ok(VALID_PREFIXES.includes('refactor!: '));
      assert.ok(VALID_PREFIXES.includes('feature!: '));
      assert.ok(VALID_PREFIXES.includes('feat!: '));
      assert.ok(VALID_PREFIXES.includes('chore(main): '));
      assert.ok(VALID_PREFIXES.includes('Merge branch '));
      assert.ok(VALID_PREFIXES.includes('dev: '));
      assert.ok(VALID_PREFIXES.includes('env: '));
      assert.ok(VALID_PREFIXES.includes('docs: '));
      assert.ok(VALID_PREFIXES.includes('ci: '));
      assert.ok(VALID_PREFIXES.includes('test: '));
    });
  });

  describe('isProductFile', () => {
    it('identifies files in rancher or internal directories as product changes', () => {
      assert.strictEqual(isProductFile('rancher/resource_cluster.go'), true);
      assert.strictEqual(isProductFile('./rancher/resource_cluster.go'), true);
      assert.strictEqual(isProductFile('rancher2/resource_cluster.go'), true);
      assert.strictEqual(isProductFile('./rancher2/resource_cluster.go'), true);
      assert.strictEqual(isProductFile('internal/provider/provider.go'), true);
      assert.strictEqual(isProductFile('./internal/provider/provider.go'), true);
    });

    it('identifies files outside rancher or internal as non-product changes', () => {
      assert.strictEqual(isProductFile('docs/resources/cluster.md'), false);
      assert.strictEqual(isProductFile('.github/workflows/lint.yml'), false);
      assert.strictEqual(isProductFile('examples/use-cases/one/main.tf'), false);
      assert.strictEqual(isProductFile('rancher_test.go'), false);
      assert.strictEqual(isProductFile(null), false);
    });
  });

  describe('isProductPrefix', () => {
    it('identifies product prefixes correctly', () => {
      assert.strictEqual(isProductPrefix('fix: resolve issue'), true);
      assert.strictEqual(isProductPrefix('feat: new capability'), true);
      assert.strictEqual(isProductPrefix('feature: new capability'), true);
      assert.strictEqual(isProductPrefix('refactor: simplify code'), true);
      assert.strictEqual(isProductPrefix('refactor!: breaking rewrite'), true);
    });

    it('returns false for non-product prefixes', () => {
      assert.strictEqual(isProductPrefix('docs: update readme'), false);
      assert.strictEqual(isProductPrefix('ci: update workflow'), false);
      assert.strictEqual(isProductPrefix('test: add unit test'), false);
      assert.strictEqual(isProductPrefix('dev: update scripts'), false);
      assert.strictEqual(isProductPrefix('env: update nix flake'), false);
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
      assert.strictEqual(validatePrefix('refactor: rewrite logic'), true);
      assert.strictEqual(validatePrefix('refactor!: breaking change'), true);
      assert.strictEqual(validatePrefix('feature!: breaking change'), true);
      assert.strictEqual(validatePrefix('feat!: breaking change'), true);
      assert.strictEqual(validatePrefix('chore(main): update dependencies'), true);
      assert.strictEqual(validatePrefix('Merge branch main into dev'), true);
      assert.strictEqual(validatePrefix('dev: update agent scripts'), true);
      assert.strictEqual(validatePrefix('env: update nix flake'), true);
      assert.strictEqual(validatePrefix('docs: update readme'), true);
      assert.strictEqual(validatePrefix('ci: update workflow'), true);
      assert.strictEqual(validatePrefix('test: add unit test'), true);
    });

    it('rejects disallowed prefixes', () => {
      assert.throws(
        () => validatePrefix('style: format whitespace'),
        /\.\.\.Commit message does not start with the required prefix\./
      );
      assert.throws(
        () => validatePrefix('random: unsupported commit'),
        /\.\.\.Commit message does not start with the required prefix\./
      );
      assert.throws(
        () => validatePrefix(' feat: leading space prefix'),
        /\.\.\.Commit message does not start with the required prefix\./
      );
    });
  });

  describe('validateLength', () => {
    it('accepts messages 70 characters or fewer', () => {
      assert.strictEqual(validateLength('feat: short subject line'), true);
      const exact70 = 'feat: ' + 'a'.repeat(64);
      assert.strictEqual(exact70.length, 70);
      assert.strictEqual(validateLength(exact70), true);
    });

    it('rejects messages longer than 70 characters', () => {
      const over70 = 'feat: ' + 'a'.repeat(65);
      assert.strictEqual(over70.length, 71);
      assert.throws(
        () => validateLength(over70),
        /\.\.\.Commit message subject line should be less than 70 characters \(this is a GitHub constraint\), found 71\./
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

  describe('validateProductImpact', () => {
    it('accepts product commit prefix when product files are changed', () => {
      assert.strictEqual(
        validateProductImpact('feat: add cluster proxy', ['rancher/resource_cluster.go']),
        true
      );
      assert.strictEqual(
        validateProductImpact('fix: fix token parsing', ['internal/auth/token.go']),
        true
      );
    });

    it('rejects non-product commit prefix when product files are changed', () => {
      assert.throws(
        () => validateProductImpact('docs: update docstrings', ['rancher/resource_cluster.go']),
        /\.\.\.Commit message does not indicate a product change, but files in \.\/rancher or \.\/internal were modified\./
      );
    });

    it('accepts non-product commit prefix when only non-product files are changed', () => {
      assert.strictEqual(
        validateProductImpact('docs: update readme', ['docs/resources/cluster.md']),
        true
      );
    });

    it('rejects product commit prefix when only non-product files are changed', () => {
      assert.throws(
        () => validateProductImpact('feat: add feature', ['docs/resources/cluster.md']),
        /\.\.\.Commit message cannot use a product prefix because no files in \.\/rancher or \.\/internal were modified\./
      );
    });

    it('skips validation for merge commits or empty changed files', () => {
      assert.strictEqual(validateProductImpact('Merge branch main into dev', ['rancher/foo.go']), true);
      assert.strictEqual(validateProductImpact('feat: something', []), true);
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

    it('enforces product impact when changedFiles are passed', () => {
      const mockExec = () => '';
      assert.strictEqual(
        validateMessage('feat: add cluster proxy', mockExec, ['rancher/resource_cluster.go']),
        true
      );
      assert.throws(
        () => validateMessage('docs: update readme', mockExec, ['rancher/resource_cluster.go']),
        /\.\.\.Commit message does not indicate a product change/
      );
    });
  });

  describe('getChangedFiles', () => {
    it('falls back to git diff when prNumber is not provided', () => {
      const mockExec = (cmd, args) => {
        assert.strictEqual(cmd, 'git');
        assert.deepStrictEqual(args, ['diff', '--name-only', 'origin/main..HEAD']);
        return 'rancher/resource_cluster.go\ndocs/cluster.md\n';
      };
      const files = getChangedFiles(undefined, mockExec);
      assert.deepStrictEqual(files, ['rancher/resource_cluster.go', 'docs/cluster.md']);
    });

    it('calls gh api when prNumber is provided', () => {
      const mockExec = (cmd, args) => {
        assert.strictEqual(cmd, 'gh');
        assert.deepStrictEqual(args, [
          'api',
          '--paginate',
          'repos/{owner}/{repo}/pulls/42/files',
          '--jq',
          '.[].filename',
        ]);
        return 'internal/provider/provider.go\n';
      };
      const files = getChangedFiles('42', mockExec);
      assert.deepStrictEqual(files, ['internal/provider/provider.go']);
    });

    it('propagates git diff failures without swallowing them', () => {
      const mockExec = () => {
        throw new Error('fatal: ambiguous argument origin/main');
      };
      assert.throws(() => getChangedFiles(undefined, mockExec), /fatal: ambiguous argument/);
    });
  });

  describe('getCommitMessages', () => {
    it('falls back to git log when prNumber is not provided', () => {
      const mockExec = (cmd, args) => {
        assert.strictEqual(cmd, 'git');
        assert.deepStrictEqual(args, ['log', 'origin/main..HEAD', '--format=%s']);
        return 'feat: commit one\nfix: commit two\n';
      };
      const messages = getCommitMessages(undefined, mockExec);
      assert.deepStrictEqual(messages, ['feat: commit one', 'fix: commit two']);
    });

    it('propagates git log failures without swallowing them', () => {
      const mockExec = () => {
        throw new Error('fatal: ambiguous argument origin/main: unknown revision');
      };
      assert.throws(
        () => getCommitMessages(undefined, mockExec),
        /fatal: ambiguous argument/
      );
    });

    it('calls gh api when prNumber is provided', () => {
      const mockExec = (cmd, args) => {
        assert.strictEqual(cmd, 'gh');
        assert.deepStrictEqual(args, [
          'api',
          '--paginate',
          'repos/{owner}/{repo}/pulls/42/commits',
          '--jq',
          '.[].commit.message | split("\\n")[0]',
        ]);
        return 'feat: api commit one\nfix: api commit two\n';
      };
      const messages = getCommitMessages('42', mockExec);
      assert.deepStrictEqual(messages, ['feat: api commit one', 'fix: api commit two']);
    });

    it('prevents shell injection when malicious prNumber is supplied', () => {
      const calls = [];
      const maliciousPr = '42"; rm -rf /; echo "pwned';
      const mockExec = (cmd, args) => {
        calls.push({ cmd, args });
        return 'feat: test commit\n';
      };
      const messages = getCommitMessages(maliciousPr, mockExec);
      assert.strictEqual(calls.length, 1);
      assert.strictEqual(calls[0].cmd, 'gh');
      assert.strictEqual(calls[0].args[2], `repos/{owner}/{repo}/pulls/${maliciousPr}/commits`);
      assert.deepStrictEqual(messages, ['feat: test commit']);
    });

    it('preserves commit text with leading whitespace without trimming', () => {
      const mockExecGit = () => ' feat: leading whitespace commit\n';
      assert.deepStrictEqual(getCommitMessages(undefined, mockExecGit), [' feat: leading whitespace commit']);

      const mockExecGh = () => ' feat: api leading whitespace commit\n';
      assert.deepStrictEqual(getCommitMessages('42', mockExecGh), [' feat: api leading whitespace commit']);
    });
  });

  describe('main', () => {
    it('runs validation loop across all commits', async () => {
      const mockExec = (cmd, args) => {
        if (cmd === 'gh') {
          if (args && args.some((arg) => typeof arg === 'string' && arg.includes('/files'))) {
            return 'rancher/resource_cluster.go\n';
          }
          return 'feat: first commit\nfix: second commit\n';
        }
        return '';
      };
      await assert.doesNotReject(async () => {
        await main({ prNumber: '100', exec: mockExec, messageExec: () => '' });
      });
    });

    it('handles empty commit list without error', async () => {
      const mockExec = () => '';
      await assert.doesNotReject(async () => {
        await main({ prNumber: '100', exec: mockExec, messageExec: () => '' });
      });
    });

    it('rejects in main when commit prefix does not match file impact', async () => {
      const mockExec = (cmd, args) => {
        if (cmd === 'gh') {
          if (args && args.some((arg) => typeof arg === 'string' && arg.includes('/files'))) {
            return 'docs/README.md\n';
          }
          return 'feat: docs update\n';
        }
        return '';
      };
      await assert.rejects(
        async () => {
          await main({ prNumber: '100', exec: mockExec, messageExec: () => '' });
        },
        /\.\.\.Commit message cannot use a product prefix/
      );
    });
  });
});
