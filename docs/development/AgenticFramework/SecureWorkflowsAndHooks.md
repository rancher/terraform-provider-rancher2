# Agentic Framework: Secure Workflows & Hook Enforcements

## Abstract

To establish a zero-trust, compliant development workspace, the Agentic Framework implements a **secure workflow and git hook enforcement system**. This system enforces a zero-bypass interception pattern that blocks direct, unvetted git commands (such as manual `git commit` or `git push`), mandates planning validation before modifying source code, and provides highly resilient branch synchronization utilities to guarantee zero developer data loss.

---

## Technical Specification

### 1. Zero-Bypass Git Interception Pattern

Direct execution of native git commit or push commands is highly prone to bypasses of pre-commit quality gates. To eliminate this risk, the framework blocks all manual commit/push attempts and intercepts them at the shell/tool invocation layer.

```text
Developer Command [git commit] or [git push]
                    │
                    ▼
   [BeforeTool Hook: block-restricted-commands.js]
                    │
                    ▼
        [Unconditional Rejection]
                    │
                    ▼
  [Redirect to .gemini/skills/commit-push.sh]
```

- **Direct Git Block Hook (`block-restricted-commands.js`)**: Configured on the `run_shell_command` matcher in `.gemini/settings.json`. It intercepts all shell commands. If a developer or agent tries to execute `git commit` or `git push` directly, the hook unconditionally rejects the command and redirects them to use the secure `.gemini/skills/commit-push.sh` skill.
- **Upstream Push Block**: It additionally parses remote destinations. Any attempt to push directly to the upstream "rancher" organization repositories is instantly blocked. Push operations can only target the developer's authorized personal fork.

### 2. Live Planning Enforcement Hook (`02-plan-phase.js`)

To guarantee compliance with the **Planning Protocol (Gate 1)**, the framework prevents files from being modified without an approved blueprint in `docs/development/`.

- **Interception**: Configured on the `write_file` and `replace` matchers.
- **Verification**: Before allowing any file-writing or text-replacement tool to execute on source files (Go, Terraform, workflows, scripts), the hook reads `git status --porcelain`.
- **Rule**: It verifies that at least one added, modified, or untracked file exists under the `docs/development/` directory, confirming that a planning blueprint is actively in place. If no active blueprint is found, the tool execution is immediately denied.

### 3. High-Resilience Auto-Stashing Sync (`git-sync.sh`)

When syncing the local fork with the upstream default branch (`main`), uncommitted changes in the working tree often trigger merge conflicts, leading to potential code loss or corrupted state. To solve this, `.gemini/skills/git-sync.sh` implements an atomic auto-stashing state machine:

```text
Active Working Tree (Dirty)
           │
           ▼
[1. git stash push -k -u -m "temp-auto-stash"] ──► Working Tree Cleaned
           │
           ▼
[2. Switch to main -> Fetch Upstream -> Fast-Forward] ──► Sync Successful
           │
           ▼
[3. git stash pop (via Trap 'cleanup' EXIT)] ──► Working Tree Restored
```

- **Atomic Cleanup Trap**: On execution, the skill immediately registers an exit trap (`trap 'cleanup' EXIT`).
- **Isolation**: It pushes all active changes, untracked files, and staged changes safely onto a temporary git stash using `git stash push --keep-index --include-untracked`. This leaves the working tree completely clean.
- **Fast-Forward Sync**: It checks out `main`, fetches upstream updates, and fast-forwards cleanly.
- **Guaranteed Restoration**: On script completion (or if an unexpected error occurs midway), the exit trap automatically executes `git stash pop`, restoring the developer's working tree to its exact original state.

---

## Standing Implementation Decisions

1. **No-Bypass Hook Policy**: Bypass environment flags (such as `BYPASS_COMMIT_HOOK=1`) are strictly removed and prohibited inside block scripts. Hook execution is absolute.
2. **Path Whitelisting**: To allow the agent to create and update planning blueprints, hook configurations, and session parameters, writes/edits targeting files inside the `.gemini/` directory or `GEMINI.md` are unconditionally allowed, bypassing the planning check.
3. **Stash Safety**: Stash pushing must always use `--include-untracked` (`-u`) to ensure newly created files are safely isolated and not left behind to cause merge collisions during checkout.
