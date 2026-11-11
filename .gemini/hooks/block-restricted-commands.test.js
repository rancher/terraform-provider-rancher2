import assert from "node:assert";
import { spawn } from "node:child_process";
import test from "node:test";
import { fileURLToPath, URL } from "node:url";
import fs from "node:fs/promises";

function runHook(inputPayload) {
  return new Promise((resolve, reject) => {
    const hookPath = fileURLToPath(
      new URL("./block-restricted-commands.js", import.meta.url),
    );
    const child = spawn("node", [hookPath]);

    let stdout = "";
    let stderr = "";

    child.stdout.on("data", (data) => {
      stdout += data.toString();
    });

    child.stderr.on("data", (data) => {
      stderr += data.toString();
    });

    child.on("close", (code) => {
      resolve({ code, stdout, stderr });
    });

    child.on("error", (err) => {
      reject(err);
    });

    child.stdin.write(JSON.stringify(inputPayload));
    child.stdin.end();
  });
}

function parseJSON(stdout) {
  const lines = stdout.trim().split("\n");
  for (let i = lines.length - 1; i >= 0; i--) {
    try {
      return JSON.parse(lines[i]);
    } catch {
      // Continue scanning upward
    }
  }
  throw new Error(`Failed to parse JSON from stdout:\n${stdout}`);
}

async function loadTestVectors() {
  const filePath = new URL("./test-vectors.json", import.meta.url);
  const content = await fs.readFile(filePath, "utf-8");
  return JSON.parse(content);
}

test("block-restricted-commands.js deterministic guardrail rules", async (t) => {
  const vectors = await loadTestVectors();

  await t.test(
    "blocks blacklisted paths like /var/ or /usr/ in write_file",
    async () => {
      const payload = {
        tool_name: "write_file",
        tool_input: { file_path: "/usr/local/bin/some-script" },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(parsed.decision, "deny");
    },
  );

  await t.test("blocks id_rsa accesses", async () => {
    const payload = {
      tool_name: "read_file",
      tool_input: { file_path: "~/.ssh/id_rsa" },
    };
    const res = await runHook(payload);
    const parsed = parseJSON(res.stdout);
    assert.strictEqual(parsed.decision, "deny");
  });

  await t.test(
    "blocks destructive shell commands (rm, rmdir, shred, truncate, chmod, chown, mv)",
    async () => {
      for (const cmd of vectors.destructiveShellCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }
    },
  );

  await t.test(
    "blocks state-mutating git commands (commit, push, reset, checkout, rebase)",
    async () => {
      for (const cmd of vectors.mutatingGitCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }
    },
  );

  await t.test(
    "allows safe read-only git commands (status, diff, log, rev-parse, branch listing) and safe gh commands",
    async () => {
      for (const cmd of vectors.safeGitCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test("allows normal harmless commands", async () => {
    const payload = {
      tool_name: "run_shell_command",
      tool_input: { command: 'echo "hello"' },
    };
    const res = await runHook(payload);
    const parsed = parseJSON(res.stdout);
    assert.strictEqual(parsed.decision, "allow");
  });

  await t.test("blocks command substitution bypasses", async () => {
    for (const cmd of vectors.commandSubstitutionBypasses) {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: cmd },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(
        parsed.decision,
        "deny",
        `Expected "${cmd}" to be denied`,
      );
    }
  });

  await t.test("blocks backslash escape bypasses", async () => {
    for (const cmd of vectors.backslashEscapeBypasses) {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: cmd },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(
        parsed.decision,
        "deny",
        `Expected "${cmd}" to be denied`,
      );
    }
  });

  await t.test("blocks invoke_agent tool calls", async () => {
    const payload = {
      tool_name: "invoke_agent",
      tool_input: { agent_name: "quality_assurance", prompt: "test" },
    };
    const res = await runHook(payload);
    const parsed = parseJSON(res.stdout);
    assert.strictEqual(parsed.decision, "deny");
    assert.match(parsed.reason, /potentially destructive/);
  });

  await t.test(
    "outputs only valid JSON without unnecessary noise",
    async () => {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: 'echo "hello"' },
      };
      const res = await runHook(payload);
      const lines = res.stdout.trim().split("\n");
      assert.strictEqual(
        lines.length,
        1,
        "Expected exactly one line of output",
      );
      assert.doesNotThrow(
        () => JSON.parse(lines[0]),
        "Expected output to be valid JSON",
      );
    },
  );

  await t.test(
    "blocks access to security-controlled .gemini files (.gemini/hooks/block-restricted-commands.js, .gemini/settings.json, including traversal)",
    async () => {
      for (const file of vectors.restrictedFiles) {
        const readPayload = {
          tool_name: "read_file",
          tool_input: { file_path: file },
        };
        const readRes = await runHook(readPayload);
        const readParsed = parseJSON(readRes.stdout);
        assert.strictEqual(
          readParsed.decision,
          "deny",
          `Expected reading "${file}" to be denied`,
        );

        const writePayload = {
          tool_name: "write_file",
          tool_input: { file_path: file },
        };
        const writeRes = await runHook(writePayload);
        const writeParsed = parseJSON(writeRes.stdout);
        assert.strictEqual(
          writeParsed.decision,
          "deny",
          `Expected writing "${file}" to be denied`,
        );

        const cmdPayload = {
          tool_name: "run_shell_command",
          tool_input: { command: `cat ${file}` },
        };
        const cmdRes = await runHook(cmdPayload);
        const cmdParsed = parseJSON(cmdRes.stdout);
        assert.strictEqual(
          cmdParsed.decision,
          "deny",
          `Expected command "cat ${file}" to be denied`,
        );
      }
    },
  );

  await t.test(
    "allows access to legitimate non-gemini settings files like .vscode/settings.json",
    async () => {
      for (const file of vectors.allowedFiles) {
        const readPayload = {
          tool_name: "read_file",
          tool_input: { file_path: file },
        };
        const readRes = await runHook(readPayload);
        const readParsed = parseJSON(readRes.stdout);
        assert.strictEqual(
          readParsed.decision,
          "allow",
          `Expected reading "${file}" to be allowed`,
        );

        const writePayload = {
          tool_name: "write_file",
          tool_input: { file_path: file },
        };
        const writeRes = await runHook(writePayload);
        const writeParsed = parseJSON(writeRes.stdout);
        assert.strictEqual(
          writeParsed.decision,
          "allow",
          `Expected writing "${file}" to be allowed`,
        );

        const cmdPayload = {
          tool_name: "run_shell_command",
          tool_input: { command: `cat ${file}` },
        };
        const cmdRes = await runHook(cmdPayload);
        const cmdParsed = parseJSON(cmdRes.stdout);
        assert.strictEqual(
          cmdParsed.decision,
          "allow",
          `Expected command "cat ${file}" to be allowed`,
        );
      }
    },
  );

  await t.test("blocks shell interpreter payloads with -c or -e", async () => {
    for (const cmd of vectors.interpreterCmds) {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: cmd },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(
        parsed.decision,
        "deny",
        `Expected interpreter payload "${cmd}" to be denied`,
      );
    }
  });

  await t.test(
    "blocks wrapped and path-qualified destructive commands",
    async () => {
      for (const cmd of vectors.wrappedDestructiveCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }
    },
  );

  await t.test(
    "allows wrapper commands with harmless flags without false positive detection",
    async () => {
      for (const cmd of vectors.harmlessWrapperCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks wrapped and path-qualified mutating git commands",
    async () => {
      for (const cmd of vectors.wrappedMutatingGitCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }
    },
  );

  await t.test(
    "blocks commands wrapped in watch, su, doas, chroot, unshare, screen, and tmux",
    async () => {
      for (const cmd of vectors.wrappedBlockedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }

      for (const cmd of vectors.wrappedAllowedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks xargs wrapped destructive and mutating git commands across argument flags",
    async () => {
      for (const cmd of vectors.xargsBlockedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }

      for (const cmd of vectors.xargsAllowedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "allows wrapped or path-qualified safe read-only git commands",
    async () => {
      for (const cmd of vectors.safeReadGitCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks .env accesses while allowing .envrc without false positives",
    async () => {
      for (const envPath of vectors.blockedEnvPaths) {
        const readRes = await runHook({
          tool_name: "read_file",
          tool_input: { file_path: envPath },
        });
        assert.strictEqual(
          parseJSON(readRes.stdout).decision,
          "deny",
          `Expected read of "${envPath}" to be denied`,
        );

        const writeRes = await runHook({
          tool_name: "write_file",
          tool_input: { file_path: envPath },
        });
        assert.strictEqual(
          parseJSON(writeRes.stdout).decision,
          "deny",
          `Expected write to "${envPath}" to be denied`,
        );

        const cmdRes = await runHook({
          tool_name: "run_shell_command",
          tool_input: { command: `cat ${envPath}` },
        });
        assert.strictEqual(
          parseJSON(cmdRes.stdout).decision,
          "deny",
          `Expected command "cat ${envPath}" to be denied`,
        );

        const redirectRes = await runHook({
          tool_name: "run_shell_command",
          tool_input: { command: `echo secret >${envPath}` },
        });
        assert.strictEqual(
          parseJSON(redirectRes.stdout).decision,
          "deny",
          `Expected command "echo secret >${envPath}" to be denied`,
        );
      }

      for (const envrcPath of vectors.allowedEnvPaths) {
        const readRes = await runHook({
          tool_name: "read_file",
          tool_input: { file_path: envrcPath },
        });
        assert.strictEqual(
          parseJSON(readRes.stdout).decision,
          "allow",
          `Expected read of "${envrcPath}" to be allowed`,
        );

        const writeRes = await runHook({
          tool_name: "write_file",
          tool_input: { file_path: envrcPath },
        });
        assert.strictEqual(
          parseJSON(writeRes.stdout).decision,
          "allow",
          `Expected write to "${envrcPath}" to be allowed`,
        );

        const cmdRes = await runHook({
          tool_name: "run_shell_command",
          tool_input: { command: `cat ${envrcPath}` },
        });
        assert.strictEqual(
          parseJSON(cmdRes.stdout).decision,
          "allow",
          `Expected command "cat ${envrcPath}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks destructive and mutating commands wrapped in find (-delete, -exec, -execdir)",
    async () => {
      for (const cmd of vectors.findBlockedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected "${cmd}" to be denied`,
        );
      }

      for (const cmd of vectors.findAllowedCommands) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks web_fetch calls targeting GitHub pull requests or issues via url or prompt",
    async () => {
      for (const payload of vectors.blockedWebFetchPayloads) {
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected web_fetch payload ${JSON.stringify(payload)} to be denied`,
        );
      }

      for (const payload of vectors.allowedWebFetchPayloads) {
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected web_fetch payload ${JSON.stringify(payload)} to be allowed`,
        );
      }
    },
  );

  await t.test(
    "blocks flag-assigned inline options containing blacklisted file paths",
    async () => {
      for (const cmd of vectors.blockedFlagCmds) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected flag command "${cmd}" to be denied`,
        );
      }

      for (const cmd of vectors.allowedFlagCmds) {
        const payload = {
          tool_name: "run_shell_command",
          tool_input: { command: cmd },
        };
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "allow",
          `Expected flag command "${cmd}" to be allowed`,
        );
      }
    },
  );

  await t.test("blocks nested shell interpreter payloads", async () => {
    for (const cmd of vectors.nestedCmds) {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: cmd },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(
        parsed.decision,
        "deny",
        `Expected nested interpreter payload "${cmd}" to be denied`,
      );
    }
  });

  await t.test(
    "blocks env wrapper bypass attempts using absolute paths",
    async () => {
      const payload = {
        tool_name: "run_shell_command",
        tool_input: { command: "env /bin/rm -rf /" },
      };
      const res = await runHook(payload);
      const parsed = parseJSON(res.stdout);
      assert.strictEqual(
        parsed.decision,
        "deny",
        "Expected 'env /bin/rm' to be denied",
      );
    },
  );
});
