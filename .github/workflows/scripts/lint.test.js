import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import {
  runTerraform,
  runActionlint,
  getShellcheckFiles,
  runShellcheck,
  getJsFiles,
  runNodeCheck,
  getTestFiles,
  runNodeTest,
  runEslint,
  runGitleaks,
  runAll,
  main,
} from './lint.js';

describe('lint.js', () => {
  describe('runTerraform', () => {
    it('executes terraform fmt and tflint commands', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      runTerraform(mockExec);
      assert.deepStrictEqual(commands, [
        'terraform fmt -check -recursive -diff',
        'tflint --recursive',
      ]);
    });

    it('throws when terraform fmt fails', () => {
      const mockExec = (cmd) => {
        if (cmd.startsWith('terraform fmt')) {
          throw new Error('Format failed');
        }
      };
      assert.throws(() => runTerraform(mockExec), /Format failed/);
    });
  });

  describe('runActionlint', () => {
    it('executes actionlint command', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      runActionlint(mockExec);
      assert.deepStrictEqual(commands, ['actionlint']);
    });
  });

  describe('getShellcheckFiles & runShellcheck', () => {
    it('finds script files starting with #!', () => {
      const files = getShellcheckFiles('.');
      assert.ok(Array.isArray(files));
      assert.ok(files.some((f) => f.endsWith('.sh') || f.includes('bin/')));
    });

    it('executes shellcheck on found files', () => {
      const checkedCalls = [];
      const mockFs = {
        readdirSync: () => [
          { name: 'test.sh', isDirectory: () => false, isFile: () => true },
          { name: 'node_modules', isDirectory: () => true, isFile: () => false },
        ],
        openSync: () => 1,
        readSync: (fd, buf) => {
          buf.write('#!');
        },
        closeSync: () => {},
        existsSync: () => true,
      };
      const mockExec = (cmd, args) => {
        checkedCalls.push({ cmd, args });
      };

      runShellcheck(mockExec, mockFs);
      assert.strictEqual(checkedCalls.length, 1);
      assert.strictEqual(checkedCalls[0].cmd, 'shellcheck');
      assert.deepStrictEqual(checkedCalls[0].args, ['-x', 'test.sh']);
    });

    it('passes filenames with quotes and shell metacharacters safely as argv elements without injection', () => {
      const checkedCalls = [];
      const trickyFile = 'my $script; "evil".sh';
      const mockFs = {
        readdirSync: () => [
          { name: trickyFile, isDirectory: () => false, isFile: () => true },
        ],
        openSync: () => 1,
        readSync: (fd, buf) => {
          buf.write('#!');
        },
        closeSync: () => {},
        existsSync: () => true,
      };
      const mockExec = (cmd, args) => {
        checkedCalls.push({ cmd, args });
      };

      runShellcheck(mockExec, mockFs);
      assert.strictEqual(checkedCalls.length, 1);
      assert.strictEqual(checkedCalls[0].cmd, 'shellcheck');
      assert.deepStrictEqual(checkedCalls[0].args, ['-x', trickyFile]);
    });

    it('handles scenario when no shell scripts exist', () => {
      const mockFs = {
        readdirSync: () => [],
      };
      const mockExec = () => {
        throw new Error('Should not be called');
      };
      assert.doesNotThrow(() => runShellcheck(mockExec, mockFs));
    });
  });

  describe('getJsFiles & runNodeCheck', () => {
    it('finds js files and excludes node_modules', () => {
      const files = getJsFiles('.');
      assert.ok(files.length > 0);
      assert.ok(files.every((f) => !f.includes('node_modules')));
      assert.ok(files.every((f) => f.endsWith('.js')));
    });

    it('runs node --check on all discovered js files', () => {
      const calls = [];
      const mockFs = {
        readdirSync: () => [
          { name: 'script.js', isDirectory: () => false, isFile: () => true },
        ],
        existsSync: () => true,
      };
      const mockExec = (cmd, args) => {
        calls.push({ cmd, args });
      };
      runNodeCheck(mockExec, mockFs);
      assert.strictEqual(calls.length, 1);
      assert.strictEqual(calls[0].cmd, 'node');
      assert.deepStrictEqual(calls[0].args, ['--check', 'script.js']);
    });
  });

  describe('getTestFiles & runNodeTest', () => {
    it('finds test files ending in .test.js or .test.mjs', () => {
      const files = getTestFiles('.');
      assert.ok(files.length > 0);
      assert.ok(files.every((f) => f.endsWith('.test.js') || f.endsWith('.test.mjs')));
    });

    it('executes node --test on discovered test files', () => {
      const calls = [];
      const mockFs = {
        readdirSync: () => [
          { name: 'sample.test.js', isDirectory: () => false, isFile: () => true },
        ],
      };
      const mockExec = (cmd, args) => {
        calls.push({ cmd, args });
      };
      runNodeTest(mockExec, mockFs);
      assert.strictEqual(calls.length, 1);
      assert.strictEqual(calls[0].cmd, 'node');
      assert.deepStrictEqual(calls[0].args, ['--test', 'sample.test.js']);
    });
  });

  describe('runEslint', () => {
    it('executes npm ci and eslint .', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      runEslint(mockExec);
      assert.deepStrictEqual(commands, ['npm ci', 'eslint .']);
    });
  });

  describe('runGitleaks', () => {
    it('executes secret scanning commands', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      runGitleaks(mockExec);
      assert.deepStrictEqual(commands, [
        'gitleaks detect --no-banner -v --no-git',
        'gitleaks detect --no-banner -v',
      ]);
    });
  });

  describe('runAll and main dispatch', () => {
    it('runAll executes all validation steps', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      const mockFs = {
        readdirSync: () => [],
      };
      runAll(mockExec, mockFs);
      assert.ok(commands.includes('terraform fmt -check -recursive -diff'));
      assert.ok(commands.includes('actionlint'));
      assert.ok(commands.includes('npm ci'));
      assert.ok(commands.includes('eslint .'));
      assert.ok(commands.includes('gitleaks detect --no-banner -v --no-git'));
    });

    it('main dispatches correctly for specific modes', () => {
      const commands = [];
      const mockExec = (cmd) => {
        commands.push(cmd);
      };
      const mockFs = {
        readdirSync: () => [],
      };

      main('terraform', mockExec, mockFs);
      assert.ok(commands.includes('terraform fmt -check -recursive -diff'));

      main('actionlint', mockExec, mockFs);
      assert.ok(commands.includes('actionlint'));

      main('eslint', mockExec, mockFs);
      assert.ok(commands.includes('eslint .'));
    });
  });
});
