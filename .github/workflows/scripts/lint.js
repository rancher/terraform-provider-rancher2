import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import fs from 'node:fs';
import path from 'node:path';

export function runTerraform(exec = execSync) {
  console.log('==> Linting Terraform files...');
  try {
    exec('terraform fmt -check -recursive -diff', { stdio: 'inherit' });
  } catch (error) {
    console.error("terraform format error, please run 'terraform fmt -recursive'");
    throw error;
  }
  exec('tflint --recursive', { stdio: 'inherit' });
}

export function runActionlint(exec = execSync) {
  console.log('==> Linting GitHub workflows...');
  exec('actionlint', { stdio: 'inherit' });
}

export function getShellcheckFiles(dir = '.', fsModule = fs) {
  const ignoredDirs = new Set(['.git', '.terraform', '.agent', 'bin', 'node_modules']);
  const files = [];

  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (ignoredDirs.has(entry.name)) {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (entry.isFile()) {
        if (entry.name.endsWith('.md') || entry.name.endsWith('.js')) {
          continue;
        }
        try {
          const fd = fsModule.openSync(fullPath, 'r');
          const buffer = Buffer.alloc(2);
          fsModule.readSync(fd, buffer, 0, 2, 0);
          fsModule.closeSync(fd);
          if (buffer.toString('utf-8') === '#!') {
            files.push(fullPath);
          }
        } catch {
          // Ignore unreadable files
        }
      }
    }
  }

  scan(dir);
  return files;
}

export function runShellcheck(exec = execSync, fsModule = fs) {
  console.log('==> Running shellcheck...');
  const files = getShellcheckFiles('.', fsModule);

  if (files.length === 0) {
    console.log('No shell scripts found to check.');
    return;
  }

  for (const file of files) {
    if (fsModule.existsSync(file)) {
      console.log(`Checking ${file}...`);
      exec(`shellcheck -x "${file}"`, { stdio: 'inherit' });
    }
  }
}

export function getJsFiles(dir = '.', fsModule = fs) {
  const files = [];
  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (entry.name === 'node_modules' || entry.name === '.git') {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (entry.isFile() && entry.name.endsWith('.js')) {
        files.push(fullPath);
      }
    }
  }
  scan(dir);
  return files;
}

export function runNodeCheck(exec = execSync, fsModule = fs) {
  console.log('==> Running Node Syntax Check...');
  const files = getJsFiles('.', fsModule);

  if (files.length === 0) {
    console.log('No Node files found to check.');
    return;
  }

  for (const file of files) {
    if (fsModule.existsSync(file)) {
      console.log(`Checking ${file}...`);
      exec(`node --check "${file}"`, { stdio: 'inherit' });
    }
  }
}

export function getTestFiles(dir = '.', fsModule = fs) {
  const files = [];
  function scan(currentDir) {
    const entries = fsModule.readdirSync(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      if (entry.name === 'node_modules' || entry.name === '.git') {
        continue;
      }
      const fullPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        scan(fullPath);
      } else if (
        entry.isFile() &&
        (entry.name.endsWith('.test.js') || entry.name.endsWith('.test.mjs'))
      ) {
        files.push(fullPath);
      }
    }
  }
  scan(dir);
  return files;
}

export function runNodeTest(exec = execSync, fsModule = fs) {
  console.log('==> Running Node Unit Tests...');
  const files = getTestFiles('.', fsModule);

  if (files.length === 0) {
    console.log('No test files found to run.');
    return;
  }

  const quotedFiles = files.map((f) => `"${f}"`).join(' ');
  exec(`node --test ${quotedFiles}`, { stdio: 'inherit' });
}

export function runEslint(exec = execSync) {
  console.log('==> Running ESLint...');
  exec('npm ci', { stdio: 'inherit' });
  exec('eslint .', { stdio: 'inherit' });
}

export function runGitleaks(exec = execSync) {
  console.log('==> Scanning for secrets with gitleaks...');
  exec('gitleaks detect --no-banner -v --no-git', { stdio: 'inherit' });
  exec('gitleaks detect --no-banner -v', { stdio: 'inherit' });
}

export function runAll(exec = execSync, fsModule = fs) {
  runTerraform(exec);
  runActionlint(exec);
  runShellcheck(exec, fsModule);
  runNodeCheck(exec, fsModule);
  runNodeTest(exec, fsModule);
  runEslint(exec);
  runGitleaks(exec);
}

export function main(mode = process.argv[2] || 'all', exec = execSync, fsModule = fs) {
  switch (mode) {
  case 'terraform':
    runTerraform(exec);
    break;
  case 'actionlint':
    runActionlint(exec);
    break;
  case 'shellcheck':
    runShellcheck(exec, fsModule);
    break;
  case 'node-check':
    runNodeCheck(exec, fsModule);
    break;
  case 'node-test':
    runNodeTest(exec, fsModule);
    break;
  case 'eslint':
    runEslint(exec);
    break;
  case 'gitleaks':
    runGitleaks(exec);
    break;
  case 'all':
    runAll(exec, fsModule);
    break;
  default:
    console.error(`Error: Unknown lint mode: ${mode}`);
    console.error(
      'Usage: lint.js [terraform|actionlint|shellcheck|node-check|node-test|eslint|gitleaks|all]'
    );
    process.exit(1);
  }
}

if (process.argv[1] && fileURLToPath(import.meta.url) === path.resolve(process.argv[1])) {
  main();
}
