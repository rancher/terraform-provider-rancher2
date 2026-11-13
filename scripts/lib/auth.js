import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

/**
 * Resolves an active GitHub authentication token.
 * Precedence:
 * 1. env.GITHUB_TOKEN or env.GH_TOKEN
 * 2. Fallback to `gh auth token` via GitHub CLI
 *
 * @param {object} [options]
 * @param {object} [options.env=process.env]
 * @param {Function} [options.execFn] - For testing injection
 * @returns {Promise<string>}
 */
export async function getAuthToken({
  env = process.env,
  execFn = execFileAsync,
} = {}) {
  const envToken = env.GITHUB_TOKEN || env.GH_TOKEN;
  if (envToken && envToken.trim()) {
    return envToken.trim();
  }

  try {
    const { stdout } = await execFn("gh", ["auth", "token"], {
      env: { ...env, GITHUB_TOKEN: "", GH_TOKEN: "" },
    });
    if (stdout && stdout.trim()) {
      return stdout.trim();
    }
  } catch {
    // gh CLI might not be installed or logged in
  }

  return "";
}

/**
 * Resolves the authenticated user's login from GitHub API.
 *
 * @param {object} options
 * @param {string} options.token
 * @param {Function} [options.fetchFn=globalThis.fetch]
 * @returns {Promise<string|null>} Login username or null if unauthenticated
 */
export async function getAuthenticatedUser({
  token,
  fetchFn = globalThis.fetch,
} = {}) {
  if (!token) {
    return null;
  }

  try {
    const res = await fetchFn("https://api.github.com/user", {
      headers: {
        Accept: "application/vnd.github+json",
        Authorization: `Bearer ${token}`,
        "User-Agent": "rancher-terraform-provider-release-scripts",
      },
    });

    if (!res.ok) {
      return null;
    }

    const data = await res.json();
    return data && data.login ? data.login : null;
  } catch {
    return null;
  }
}
