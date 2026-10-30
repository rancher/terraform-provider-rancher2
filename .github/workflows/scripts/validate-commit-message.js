import { execFileSync, execSync } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

export const PRODUCT_PREFIXES = [
  'fix: ',
  'feature: ',
  'feat: ',
  'refactor: ',
  'refactor!: ',
  'feature!: ',
  'feat!: ',
];

export const NON_PRODUCT_PREFIXES = [
  'dev: ', // representing changes to the development process (scripts/agentic framework changes)
  'env: ', // representing changes to the unified dev/test environment (flake or env variable changes)
  'docs: ', // representing changes only to documentation
  'ci: ', // representing changes only to workflows and workflow scripts
  'test: ', // representing changes only to tests
];

export const VALID_PREFIXES = [
  // end user product updates
  ...PRODUCT_PREFIXES,
  // merge and release (automated) updates
  'chore(main): ',
  'Merge branch ',
  // non product changes
  ...NON_PRODUCT_PREFIXES,
];

export function isProductPrefix(message) {
  return PRODUCT_PREFIXES.some((prefix) => message.startsWith(prefix));
}

export function isProductDirectory(file) {
  if (!file) {
    return false;
  }
  const normalized = file.trim().replace(/^\.\//, '').replace(/\\/g, '/');
  return (
    normalized.startsWith('rancher/') ||
    normalized.startsWith('rancher2/') ||
    normalized.startsWith('internal/') ||
    normalized === 'main.go' ||
    normalized === 'main_test.go'
  );
}

export function isProductFile(file) {
  if (!file) {
    return false;
  }
  const normalized = file.trim().replace(/^\.\//, '').replace(/\\/g, '/');
  return isProductDirectory(normalized) && !normalized.endsWith('_test.go');
}

export function validateEmpty(message) {
  if (!message || message.trim() === '') {
    throw new Error('...Empty commit message.');
  }
  return true;
}

export function validatePrefix(message) {
  const hasValidPrefix = VALID_PREFIXES.some((prefix) => message.startsWith(prefix));
  if (!hasValidPrefix) {
    const errorMsg = [
      '...Commit message does not start with the required prefix.',
      'Please use one of the following prefixes when making changes that change the product behavior: ',
      '"fix:", "feature:", "feat:", "refactor:", "refactor!:", "feature!:", or "feat!:".',
      'Please use one of the following prefixes when making changes that do not change the product behavior: ',
      '"dev:", "env:", "docs:", "ci:", or "test:".',
      'This enables release-please to automatically format release notes based on the commit message.',
      message,
    ].join('\n');
    throw new Error(errorMsg);
  }
  return true;
}

export function validateLength(message) {
  if (message.length > 70) {
    throw new Error(`...Commit message subject line should be less than 70 characters (this is a GitHub constraint), found ${message.length}.`);
  }
  return true;
}

export function validateSpelling(message, exec = execSync) {
  if (message.startsWith('Merge ') || message.startsWith('Merge branch ')) {
    return true;
  }
  try {
    exec('cspell stdin', { input: message, encoding: 'utf-8', stdio: ['pipe', 'pipe', 'pipe'] });
  } catch (error) {
    const output = (error.stdout || error.stderr || error.message || '').toString();
    const errorMsg = [
      `...Commit message contains spelling errors on: ^${output.trim()}$`,
      '...Also try updating the PR title.',
      '...If this is a mistake, add your word to the custom_words.txt file, it is case insensitive.',
    ].join('\n');
    throw new Error(errorMsg);
  }
  return true;
}

export function validateProductImpact(message, changedFiles = []) {
  if (message.startsWith('Merge branch ') || message.startsWith('Merge ')) {
    return true;
  }
  if (!changedFiles || changedFiles.length === 0) {
    return true;
  }

  const isProduct = isProductPrefix(message);
  const isTest = message.startsWith('test: ');
  const hasProductFiles = changedFiles.some(isProductFile);
  const hasProductDirectoryFiles = changedFiles.some(isProductDirectory);

  if (isTest) {
    if (hasProductFiles) {
      const errorMsg = [
        '...Commit message cannot use "test:" prefix because files in ./rancher or ./internal that are not tests were modified.',
        'The "test:" prefix is allowed only if all files are outside of end user product directories or only contain "_test.go" files.',
        'Please use one of the following product prefixes: "fix:", "feat:", "feature:", "refactor:", "refactor!:", "feature!:", or "feat!:".',
        message,
      ].join('\n');
      throw new Error(errorMsg);
    }
    return true;
  }

  if (hasProductFiles && !isProduct) {
    const errorMsg = [
      '...Commit message does not indicate a product change, but files in ./rancher or ./internal were modified.',
      'Please use one of the following product prefixes: "fix:", "feat:", "feature:", "refactor:", "refactor!:", "feature!:", or "feat!:".',
      message,
    ].join('\n');
    throw new Error(errorMsg);
  }

  if (!isProduct && hasProductDirectoryFiles) {
    const errorMsg = [
      '...Commit message does not indicate a product change, but files in ./rancher or ./internal were modified.',
      'Please use "test:" if modifying tests or one of the product prefixes for product code.',
      message,
    ].join('\n');
    throw new Error(errorMsg);
  }

  if (!hasProductFiles && isProduct) {
    const errorMsg = [
      '...Commit message cannot use a product prefix because no files in ./rancher or ./internal were modified.',
      'Please use one of the following non-product prefixes: "dev:", "env:", "docs:", "ci:", or "test:".',
      message,
    ].join('\n');
    throw new Error(errorMsg);
  }

  return true;
}

export function validateMessage(message, exec = execSync, changedFiles = undefined) {
  validateEmpty(message);
  validatePrefix(message);
  validateLength(message);
  validateSpelling(message, exec);
  if (changedFiles !== undefined) {
    validateProductImpact(message, changedFiles);
  }
  return true;
}

export function getCommitMessages(prNumber, exec = execFileSync) {
  if (!prNumber) {
    console.log('Notice: PR_NUMBER environment variable is not set.');
    console.log('Falling back to checking local commits against origin/main...');
    const output = exec('git', ['log', 'origin/main..HEAD', '--format=%s'], { encoding: 'utf-8' });
    return (output || '')
      .split('\n')
      .map((line) => line.replace(/\r$/, ''))
      .filter((line) => line !== '');
  }

  const output = exec(
    'gh',
    [
      'api',
      '--paginate',
      `repos/{owner}/{repo}/pulls/${prNumber}/commits`,
      '--jq',
      '.[].commit.message | split("\\n")[0]',
    ],
    { encoding: 'utf-8' }
  );
  return (output || '')
    .split('\n')
    .map((line) => line.replace(/\r$/, ''))
    .filter((line) => line !== '');
}

export function getCommits(prNumber, exec = execFileSync) {
  if (!prNumber) {
    console.log('Notice: PR_NUMBER environment variable is not set.');
    console.log('Falling back to checking local commits against origin/main...');
    const output = exec('git', ['log', 'origin/main..HEAD', '--format=%H%x09%s'], { encoding: 'utf-8' });
    return (output || '')
      .split('\n')
      .map((line) => line.replace(/\r$/, ''))
      .filter((line) => line.trim() !== '')
      .map((line) => {
        const tabIdx = line.indexOf('\t');
        if (tabIdx !== -1) {
          const firstPart = line.slice(0, tabIdx).trim();
          if (/^[0-9a-f]{40}$/i.test(firstPart)) {
            return {
              sha: firstPart,
              message: line.slice(tabIdx + 1),
            };
          }
        }
        return {
          sha: undefined,
          message: line,
        };
      });
  }

  const output = exec(
    'gh',
    [
      'api',
      '--paginate',
      `repos/{owner}/{repo}/pulls/${prNumber}/commits`,
      '--jq',
      '.[] | (.sha + "\\t" + (.commit.message | split("\\n")[0]))',
    ],
    { encoding: 'utf-8' }
  );

  return (output || '')
    .split('\n')
    .map((line) => line.replace(/\r$/, ''))
    .filter((line) => line.trim() !== '')
    .map((line) => {
      const tabIdx = line.indexOf('\t');
      if (tabIdx !== -1) {
        const firstPart = line.slice(0, tabIdx).trim();
        if (/^[0-9a-f]{40}$/i.test(firstPart)) {
          return {
            sha: firstPart,
            message: line.slice(tabIdx + 1),
          };
        }
      }
      return {
        sha: undefined,
        message: line,
      };
    });
}

export function getCommitChangedFiles(sha, prNumber, exec = execFileSync) {
  if (!sha) {
    return [];
  }
  if (!prNumber) {
    const output = exec('git', ['diff-tree', '--no-commit-id', '--name-only', '-r', '--root', sha], { encoding: 'utf-8' });
    return (output || '')
      .split('\n')
      .map((line) => line.replace(/\r$/, '').trim())
      .filter((line) => line !== '');
  }

  const output = exec(
    'gh',
    [
      'api',
      '--paginate',
      `repos/{owner}/{repo}/commits/${sha}`,
      '--jq',
      '.files[]?.filename // empty',
    ],
    { encoding: 'utf-8' }
  );
  return (output || '')
    .split('\n')
    .map((line) => line.replace(/\r$/, '').trim())
    .filter((line) => line !== '');
}

export function getChangedFiles(prNumber, exec = execFileSync) {
  if (!prNumber) {
    console.log('Notice: PR_NUMBER environment variable is not set.');
    console.log('Falling back to checking local changed files against origin/main...');
    const output = exec('git', ['diff', '--name-only', 'origin/main..HEAD'], { encoding: 'utf-8' });
    return (output || '')
      .split('\n')
      .map((line) => line.replace(/\r$/, '').trim())
      .filter((line) => line !== '');
  }

  const output = exec(
    'gh',
    [
      'api',
      '--paginate',
      `repos/{owner}/{repo}/pulls/${prNumber}/files`,
      '--jq',
      '.[].filename',
    ],
    { encoding: 'utf-8' }
  );
  return (output || '')
    .split('\n')
    .map((line) => line.replace(/\r$/, '').trim())
    .filter((line) => line !== '');
}

export async function main({
  prNumber = process.env.PR_NUMBER,
  exec = execFileSync,
  messageExec = execSync,
  commits,
  changedFiles,
} = {}) {
  const commitList = commits !== undefined ? commits : getCommits(prNumber, exec);

  console.log('Commit messages found:');
  console.log(commitList.map((c) => (typeof c === 'string' ? c : c.message)).join('\n'));

  if (commitList.length === 0) {
    console.log('No commit messages to check.');
    return;
  }

  for (const item of commitList) {
    const message = typeof item === 'string' ? item : item.message;
    const sha = typeof item === 'object' && item !== null ? item.sha : undefined;

    console.log(`checking message ^${message}$`);

    let filesForCommit;
    if (changedFiles !== undefined) {
      filesForCommit = changedFiles;
    } else if (sha) {
      filesForCommit = getCommitChangedFiles(sha, prNumber, exec);
    } else {
      filesForCommit = getChangedFiles(prNumber, exec);
    }

    validateMessage(message, messageExec, filesForCommit);
    console.log(`message ^${message}$ passed all checks`);
  }
}

export default async function run(params = {}) {
  const core = params.core;
  try {
    await main();
    if (core && typeof core.info === 'function') {
      core.info('Commit message validation passed.');
    }
  } catch (error) {
    if (core && typeof core.setFailed === 'function') {
      core.setFailed(error.message);
    } else {
      console.error(error.message);
      process.exit(1);
    }
  }
}

if (process.argv[1] && fileURLToPath(import.meta.url) === path.resolve(process.argv[1])) {
  main().catch((err) => {
    console.error(err.message);
    process.exit(1);
  });
}
