import assert from "node:assert";
import { spawn } from "node:child_process";
import test from "node:test";
import { fileURLToPath, URL } from "node:url";

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

test("block-restricted-commands.js deterministic guardrail rules", async (t) => {
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
      for (const cmd of [
        "rm -rf ./tmp",
        "rmdir ./empty_dir",
        "shred -u secret.txt",
        "truncate -s 0 file.log",
        "(rm -rf ./tmp)",
        "{ rm -rf ./tmp; }",
        "{rm -rf ./tmp;}",
        "chmod +x script.sh",
        "(chmod +x script.sh)",
        "{ chmod +x script.sh; }",
        "chown root:root file",
        "mv old new",
      ]) {
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
      for (const cmd of [
        'git commit -m "fix"',
        "git push origin main",
        "git reset --hard HEAD~1",
        "git checkout main",
        "git switch main",
        "git rebase main",
        "git apply patch.diff",
        "git revert HEAD",
        "git bisect start",
        "git worktree add ../wt",
        "git gc",
        "git repack",
        "git branch new-feature",
        "git branch new-feature main",
        "git branch -m old-branch new-branch",
        "git branch -M old-branch new-branch",
        "git branch --move old-branch new-branch",
        "git branch -c old-branch copy-branch",
        "git branch -d old-branch",
        "git branch -D old-branch",
        "git branch --delete old-branch",
        "git branch -u upstream/main",
        "git branch -f feature-branch main",
        "git tag -d v1.0.0",
        "git tag v1.0.0",
        "gh pr merge 123",
        "gh pr close 123",
        "gh issue close 45",
        "gh release delete v1.0.0",
        "gh pr create --title 'test'",
        "gh issue edit 45",
        "gh -R owner/repo pr merge 123",
        "gh pr checkout 123",
        "gh pr comment 123 --body 'test'",
        "gh pr review 123 --approve",
      ]) {
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
      for (const cmd of [
        "git status",
        "git diff HEAD",
        "git log -n 5",
        "git rev-parse --show-toplevel",
        "git branch",
        "git branch -a",
        "git branch -r",
        "git branch --show-current",
        "git branch -v",
        "git branch -l",
        "git branch --list 'release/*'",
        "gh pr view 123",
        "gh pr list",
        "gh pr status",
        "gh pr diff 123",
        "gh issue list",
        "gh run list",
      ]) {
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
    const substitutions = [
      "echo $(rm -rf /)",
      "echo `rm -rf /`",
      'echo "$(chmod +x script)"',
      'echo "`chown root:root file`"',
      'echo "user\'s directory: $(rm -rf tmp)"',
      'echo "don\'t `rm -rf tmp`"',
    ];
    for (const cmd of substitutions) {
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
    const escapes = [
      "\\rm -rf /",
      "r\\m -rf /",
      "c\\at .gemini/settings.jso\\n",
    ];
    for (const cmd of escapes) {
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
      for (const file of [
        ".gemini/hooks/block-restricted-commands.js",
        ".gemini/settings.json",
        "./.gemini/hooks/block-restricted-commands.js",
        "./.gemini/settings.json",
        ".gemini/hooks/../settings.json",
        ".gemini/hooks/../../.gemini/settings.json",
        ".gemini/hooks/../hooks/block-restricted-commands.js",
        "/absolute/path/.gemini/settings.json",
        "hooks/block-restricted-commands.js",
        "./hooks/block-restricted-commands.js",
        "block-restricted-commands.js",
        "./block-restricted-commands.js",
      ]) {
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
      for (const file of [
        ".vscode/settings.json",
        "./.vscode/settings.json",
        "config/settings.json",
      ]) {
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
    const interpreterCmds = [
      "bash -c 'rm -rf ./tmp'",
      'bash -c"rm -rf ./tmp"',
      "bash -lc 'rm -rf ./tmp'",
      'bash -lc"rm -rf ./tmp"',
      "bash -xec 'rm -rf ./tmp'",
      'sh -c "git push origin main"',
      'sh -c"git push origin main"',
      'sh -ec "git push origin main"',
      'zsh -c "chmod +x script.sh"',
      'zsh -c"chmod +x script.sh"',
      'zsh -lc "chmod +x script.sh"',
      "node -e \"require('fs').rmSync('./tmp', {recursive: true})\"",
      "python -c 'import os; os.system(\"rm -rf ./tmp\")'",
      "python2 -c 'import os; os.system(\"rm -rf ./tmp\")'",
      "python3 -c 'import os; os.system(\"rm -rf ./tmp\")'",
      "ruby -c 'system(\"git commit -m fix\")'",
      "perl -e 'system(\"mv old new\")'",
    ];
    for (const cmd of interpreterCmds) {
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
      const destructive = [
        "/bin/rm -rf ./tmp",
        "/usr/bin/rm -rf ./tmp",
        "command rm -rf ./tmp",
        "command -p rm -rf ./tmp",
        "sudo rm -rf ./tmp",
        "sudo /bin/rm -rf ./tmp",
        "env rm -rf ./tmp",
        "env -i rm -rf ./tmp",
        "VAR=1 /bin/rm -rf ./tmp",
        "/bin/mv old new",
        "/usr/bin/chmod +x script.sh",
        "/usr/sbin/chown root:root file",
        "echo hello && /bin/rm -rf ./tmp",
        "timeout 10s rm -rf ./tmp",
        "timeout -k 5s 10s rm -rf ./tmp",
      ];
      for (const cmd of destructive) {
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
    "blocks wrapped and path-qualified mutating git commands",
    async () => {
      const mutating = [
        '/usr/bin/git commit -m "fix"',
        "/usr/bin/git push origin main",
        "command git checkout main",
        "sudo git reset --hard HEAD~1",
        "env git rebase main",
        'git -C ./repo commit -m "fix"',
        "git --no-pager push origin main",
        "git apply patch.diff",
        "git revert HEAD",
      ];
      for (const cmd of mutating) {
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
      const wrappedBlocked = [
        "watch rm -rf /",
        'watch "rm -rf /"',
        'watch -n 2 "rm -rf /"',
        "watch git push origin main",
        'watch -n 5 "git commit -m evil"',
        'su -c "rm -rf /"',
        'su root -c "rm -rf /"',
        'su - root -c "rm -rf /"',
        'su -c"rm -rf /"',
        'su --command="rm -rf /"',
        'su root --command "rm -rf /"',
        "doas rm -rf /",
        "doas -u root rm -rf /",
        'doas "rm -rf /"',
        "doas git push origin main",
        "chroot /mnt rm -rf /",
        "chroot --userspec=0:0 /mnt rm -rf /",
        "unshare -r rm -rf /",
        "unshare -m -u -r rm -rf /",
        "screen rm -rf /",
        "screen -d -m rm -rf /",
        'screen -S sess "rm -rf /"',
        'tmux new-session "rm -rf /"',
        'tmux new-session -d "rm -rf /"',
        'tmux run-shell "rm -rf /"',
        'tmux new-window -n test "git push origin main"',
        'tmux split-window -h "rm -rf /"',
      ];
      for (const cmd of wrappedBlocked) {
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

      const wrappedAllowed = [
        "watch ls -la",
        "watch -n 10 'git status'",
        "doas ls -l",
        "chroot /mnt ls",
        "unshare -r id",
        "screen -ls",
        "tmux list-sessions",
        "su root",
      ];
      for (const cmd of wrappedAllowed) {
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
      const xargsBlocked = [
        "xargs -I {} rm -rf {}",
        "xargs -I {} /bin/rm -rf {}",
        "find . | xargs -0 -I % rm -f %",
        "xargs -n 1 rm -f",
        "xargs -R 1 -I {} rm -rf {}",
        "xargs -S 100 -I {} rm -rf {}",
        "xargs -P 2 rm -rf",
        "xargs --process-slot-var=VAR rm -rf",
        "xargs --process-slot-var VAR rm -rf",
        "xargs --replace={} rm -rf {}",
        'xargs -d "\\n" -I {} rm -rf {}',
        "xargs -I {} git push origin main",
        'xargs -I {} sh -c "rm -rf {}"',
      ];
      for (const cmd of xargsBlocked) {
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

      const xargsAllowed = [
        "xargs -I {} echo {}",
        "cat list.txt | xargs ls -l",
        "xargs -n 1 cat",
      ];
      for (const cmd of xargsAllowed) {
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
      const safe = [
        "/usr/bin/git status",
        "command git diff HEAD",
        "git -C ./repo log -n 5",
        "/usr/bin/git rev-parse --show-toplevel",
      ];
      for (const cmd of safe) {
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
      for (const envPath of [".env", "./.env", "config/.env", "/path/to/.env"]) {
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

      for (const envrcPath of [".envrc", "./.envrc", "config/.envrc", "/path/to/.envrc"]) {
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
      const findBlocked = [
        "find . -delete",
        "find . -name '*.tmp' -delete",
        "find . -exec rm -rf {} +",
        "find . -exec rm -f {} \\;",
        "find . -execdir chmod 777 {} \\;",
        "find . -exec chown root:root {} \\;",
        "find . -exec mv {} {}.bak \\;",
        "find . -exec git push origin main \\;",
        "find . -exec git commit -m 'evil' \\;",
      ];
      for (const cmd of findBlocked) {
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

      const findAllowed = [
        "find . -name '*.go'",
        "find . -type f -print",
        "find . -exec echo {} \\;",
      ];
      for (const cmd of findAllowed) {
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
      const blockedPayloads = [
        {
          tool_name: "web_fetch",
          tool_input: {
            url: "https://github.com/rancher/terraform-provider-rancher2/pull/123",
          },
        },
        {
          tool_name: "web_fetch",
          tool_input: {
            url: "https://github.com/rancher/terraform-provider-rancher2/issues/456",
          },
        },
        {
          tool_name: "web_fetch",
          tool_input: { prompt: "Fetch https://github.com/owner/repo/pull/10" },
        },
        {
          tool_name: "web_fetch",
          tool_input: {
            url: "https://github.com/something",
            prompt: "check /pull/20",
          },
        },
      ];

      for (const payload of blockedPayloads) {
        const res = await runHook(payload);
        const parsed = parseJSON(res.stdout);
        assert.strictEqual(
          parsed.decision,
          "deny",
          `Expected web_fetch payload ${JSON.stringify(payload)} to be denied`,
        );
      }

      const allowedPayloads = [
        {
          tool_name: "web_fetch",
          tool_input: {
            url: "https://registry.terraform.io/providers/rancher/rancher2",
          },
        },
        {
          tool_name: "web_fetch",
          tool_input: { prompt: "Read documentation on terraform.io" },
        },
      ];

      for (const payload of allowedPayloads) {
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
      const blockedFlagCmds = [
        "curl --output=.gemini/settings.json https://example.com",
        'curl --output=".gemini/settings.json" https://example.com',
        "curl -o.gemini/settings.json https://example.com",
        "wget -O.env https://example.com",
        "wget -O=.env https://example.com",
        "wget -Oid_rsa https://example.com",
        "cat --file=/etc/passwd",
        "tar -C=/etc/ -xvf archive.tar",
      ];

      for (const cmd of blockedFlagCmds) {
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

      const allowedFlagCmds = [
        "curl --output=output.txt https://example.com",
        "git log --format='%H %s'",
        "tar -C /tmp -xvf archive.tar",
      ];

      for (const cmd of allowedFlagCmds) {
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
});
