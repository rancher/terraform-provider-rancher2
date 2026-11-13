import { execFile, spawn } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";
import { Buffer } from "node:buffer";

const execFileAsync = promisify(execFile);

/**
 * Resolves the git repository root directory asynchronously.
 * Falls back to directory traversal or process.cwd().
 * @returns {Promise<string>} Repository root absolute path
 */
export async function getRepoRoot() {
  try {
    const scriptDir = path.dirname(fileURLToPath(import.meta.url));
    const { stdout } = await execFileAsync(
      "git",
      ["rev-parse", "--show-toplevel"],
      { cwd: scriptDir },
    );
    if (stdout && stdout.trim()) {
      return stdout.trim();
    }
  } catch {
    // Fall back to upward directory traversal
  }

  let current = path.dirname(fileURLToPath(import.meta.url));
  while (current !== path.dirname(current)) {
    try {
      await fs.stat(path.join(current, ".git"));
      return current;
    } catch {
      current = path.dirname(current);
    }
  }
  return process.cwd();
}

/**
 * Checks if a file or path exists asynchronously.
 * @param {string} targetPath - Path to check
 * @returns {Promise<boolean>}
 */
export async function fileExists(targetPath) {
  try {
    await fs.access(targetPath);
    return true;
  } catch {
    return false;
  }
}

/**
 * Runs a command asynchronously using spawn, returning execution details.
 *
 * @param {string} command - The command or executable to run
 * @param {Array<string>} args - Arguments to pass to the executable
 * @param {object} [options]
 * @param {string} [options.cwd] - Working directory
 * @param {'inherit'|'pipe'} [options.stdio='inherit'] - Child process stdio configuration
 * @param {object} [options.env] - Additional environment variables
 * @returns {Promise<{success: boolean, code: number, stdout: string, stderr: string}>}
 */
export function runCommand(
  command,
  args,
  { cwd, stdio = "inherit", env = {} } = {},
) {
  return new Promise((resolve) => {
    const child = spawn(command, args, {
      cwd,
      stdio,
      env: {
        ...process.env,
        NO_COLOR: process.env.NO_COLOR || "1",
        ...env,
      },
    });

    const stdoutChunks = [];
    const stderrChunks = [];

    if (child.stdout) {
      child.stdout.on("data", (data) => {
        stdoutChunks.push(Buffer.isBuffer(data) ? data : Buffer.from(data));
      });
    }
    if (child.stderr) {
      child.stderr.on("data", (data) => {
        stderrChunks.push(Buffer.isBuffer(data) ? data : Buffer.from(data));
      });
    }

    let isResolved = false;

    child.once("close", (code) => {
      if (isResolved) {
        return;
      }
      isResolved = true;
      resolve({
        success: code === 0,
        code: code ?? 1,
        stdout: Buffer.concat(stdoutChunks).toString("utf8"),
        stderr: Buffer.concat(stderrChunks).toString("utf8"),
      });
    });

    child.once("error", (err) => {
      if (isResolved) {
        return;
      }
      isResolved = true;
      resolve({
        success: false,
        code: 1,
        stdout: Buffer.concat(stdoutChunks).toString("utf8"),
        stderr: err?.message || String(err),
      });
    });
  });
}
