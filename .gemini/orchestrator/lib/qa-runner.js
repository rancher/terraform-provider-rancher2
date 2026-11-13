import { exec } from "node:child_process";
import { promisify } from "node:util";
import { getRepoRoot } from "./utils.js";

const execAsync = promisify(exec);

/**
 * Runs the workspace tests and linters programmatically.
 * When `fix: true` is requested, the fixer is run sequentially before validation
 * to ensure determinism and prevent test runner race conditions.
 *
 * @param {object} [options]
 * @param {boolean} [options.fix=false] - Whether to allow linters to automatically fix files
 * @param {Function} [options.execFn=null] - Custom exec function for dependency injection in testing
 * @param {Array<string>} [options.commands] - Test and linter commands to run
 * @param {string} [options.fixerCommand] - Fixer command to run when fix is true
 * @param {string} [options.cwd=null] - Working directory for test and linter execution
 * @returns {Promise<{success: boolean, stdout: string, stderr: string}>}
 */
export async function runQAPipeline({
  fix = false,
  execFn = null,
  commands = [
    `${process.execPath} .gemini/skills/code-validation/scripts/lint.js all`,
    `${process.execPath} .gemini/skills/code-validation/scripts/test.js all`,
  ],
  fixerCommand = `${process.execPath} .gemini/skills/code-validation/scripts/lint.js all --fix`,
  cwd = null,
} = {}) {
  const runner = execFn || execAsync;
  const repoRoot = cwd || (await getRepoRoot());

  let fixResult = null;
  if (fix) {
    try {
      const { stdout, stderr } = await runner(fixerCommand, {
        cwd: repoRoot,
        maxBuffer: 10 * 1024 * 1024,
        timeout: 600000,
      });
      fixResult = {
        success: true,
        stdout: stdout && stdout.trim() ? stdout.trim() : "",
        stderr: stderr && stderr.trim() ? stderr.trim() : "",
      };
    } catch (err) {
      const outStr = String(err?.stdout ?? "").trim();
      const errStr = String(err?.stderr ?? "").trim();
      fixResult = {
        success: false,
        stdout: outStr,
        stderr: errStr || String(err?.message || err || "").trim(),
      };
    }
  }

  const results = await Promise.all(
    commands.map(async (cmd) => {
      try {
        const { stdout, stderr } = await runner(cmd, {
          cwd: repoRoot,
          maxBuffer: 10 * 1024 * 1024,
          timeout: 600000, // 10 minutes timeout
        });
        return {
          success: true,
          stdout: stdout && stdout.trim() ? stdout.trim() : "",
          stderr: stderr && stderr.trim() ? stderr.trim() : "",
        };
      } catch (err) {
        const outStr = String(err?.stdout ?? "").trim();
        const errStr = String(err?.stderr ?? "").trim();
        return {
          success: false,
          stdout: outStr,
          stderr: errStr || String(err?.message || err || "").trim(),
        };
      }
    }),
  );

  const allSuccess =
    (fixResult ? fixResult.success : true) && results.every((r) => r.success);
  const stdoutChunks = [
    fixResult?.stdout,
    ...results.map((r) => r.stdout),
  ].filter(Boolean);
  const stderrChunks = [
    fixResult?.stderr,
    ...results.map((r) => r.stderr),
  ].filter(Boolean);

  return {
    success: allSuccess,
    stdout: stdoutChunks.join("\n"),
    stderr: stderrChunks.join("\n"),
  };
}
