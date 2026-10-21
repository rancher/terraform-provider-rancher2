import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

export const VALID_PREFIXES = [
  'fix: ',
  'feature: ',
  'feat: ',
  'refactor!: ',
  'feature!: ',
  'feat!: ',
  'chore(main): ',
  'Merge branch ',
];

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
      'Please use one of the following prefixes: "fix:", "feature:", "feat:", "refactor!:", "feature!:", or "feat!:".',
      'This enables release-please to automatically format release notes based on the commit message.',
      message,
    ].join('\n');
    throw new Error(errorMsg);
  }
  return true;
}

export function validateLength(message) {
  if (message.length > 100) {
    throw new Error(`...Commit message subject line should be less than 100 characters, found ${message.length}.`);
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

export function validateMessage(message, exec = execSync) {
  validateEmpty(message);
  validatePrefix(message);
  validateLength(message);
  validateSpelling(message, exec);
  return true;
}

export function getCommitMessages(prNumber, exec = execSync) {
  if (!prNumber) {
    console.log('Notice: PR_NUMBER environment variable is not set.');
    console.log('Falling back to checking local commits against origin/main...');
    try {
      const output = exec('git log origin/main..HEAD --format=%s', { encoding: 'utf-8' });
      return (output || '')
        .split('\n')
        .map((line) => line.replace(/\r$/, ''))
        .filter((line) => line !== '');
    } catch (err) {
      console.warn('Warning: Failed to fetch commits from origin/main..HEAD:', err.message);
      return [];
    }
  }

  const output = exec(`gh api --paginate "repos/{owner}/{repo}/pulls/${prNumber}/commits" --jq '.[].commit.message | split("\\n")[0]'`, {
    encoding: 'utf-8',
  });
  return (output || '')
    .split('\n')
    .map((line) => line.replace(/\r$/, ''))
    .filter((line) => line !== '');
}

export async function main({ prNumber = process.env.PR_NUMBER, exec = execSync } = {}) {
  const commitMessages = getCommitMessages(prNumber, exec);

  console.log('Commit messages found:');
  console.log(commitMessages.join('\n'));

  if (commitMessages.length === 0) {
    console.log('No commit messages to check.');
    return;
  }

  for (const message of commitMessages) {
    console.log(`checking message ^${message}$`);
    validateMessage(message, exec);
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
