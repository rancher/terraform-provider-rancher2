import { execFile } from "node:child_process";
import path from "node:path";
import { promisify } from "node:util";
import { getRepoRoot } from "./utils.js";

export const sys = {
  execFileAsync: promisify(execFile),
};

/**
 * Programmatically retrieves PR comments using the git-pr skill script.
 * @param {string|number} [prTarget] - Optional PR number or "current"
 * @returns {Promise<string>} The formatted PR comments
 */
export async function getPRComments(prTarget = "current") {
  const target = prTarget || "current";
  const repoRoot = await getRepoRoot();
  const scriptPath = path.join(
    repoRoot,
    ".gemini/skills/github-pr/scripts/read_comments.js",
  );
  try {
    const { stdout } = await sys.execFileAsync(
      process.execPath,
      [scriptPath, String(target)],
      {
        cwd: repoRoot,
        maxBuffer: 10 * 1024 * 1024,
        timeout: 60000,
      },
    );
    return stdout;
  } catch (err) {
    const stderrMsg = err.stderr?.trim() || err.stdout?.trim() || err.message;
    throw new Error(
      `Failed to fetch PR comments programmatically: ${stderrMsg}`,
      { cause: err },
    );
  }
}
