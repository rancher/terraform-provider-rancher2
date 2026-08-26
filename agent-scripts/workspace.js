import { execSync } from 'child_process';
import os from 'os';
import path from 'path';

/**
 * Resolves the target temporary directory for the current repository workspace.
 * @param {string} [cwd] - Optional current working directory, defaults to process.cwd()
 * @returns {string} The absolute path to the target temporary directory.
 */
export function resolveTargetDir(cwd = process.cwd()) {
  const homeDir = os.homedir();
  let repoName;
  try {
    const topLevel = execSync('git rev-parse --show-toplevel', {
      cwd: cwd,
      stdio: ['ignore', 'pipe', 'ignore'],
    })
      .toString()
      .trim();
    repoName = path.basename(topLevel);
  } catch (err) {
    console.error(`🔒 Hook Info: Failed to resolve git toplevel directory: ${err.message || err}`);
    repoName = path.basename(cwd) || 'generic-repo';
  }
  return path.resolve(homeDir, '.gemini/tmp', repoName);
}
