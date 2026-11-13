# Portable `.gemini` Agent Platform & Repository Integration Contract

This document defines the operating contract between the `.gemini/` directory and any host repository adopting it.

---

## 1. Overview & Architecture

The `.gemini/` directory is a self-contained, language-agnostic agent orchestration platform. It houses:

- The programmatic developer orchestrator and multi-agent Ralph review loop (`.gemini/orchestrator/`).
- Deterministic lifecycle hooks and security guardrails (`.gemini/hooks/`).
- Specialized skills for git, GitHub PRs, GitHub Actions CI, and code validation (`.gemini/skills/`).
- Autonomous planning agents (`.gemini/agents/`).
- Self-contained dependency definitions (`.gemini/package.json`).

The platform interacts with the host repository exclusively through standard POSIX and `GNUmakefile` interfaces.

---

## 2. Consolidation & Deletion Instructions for Migrating Repositories

When copying this `.gemini/` directory into a repository that previously used legacy agent scripts or legacy process documentation:

### Safe to Delete from Host Repository

1. `agent-scripts/`: Delete the entire directory. All universal orchestrator tools and libraries now reside under `.gemini/orchestrator/`.
2. `docs/development/`: Delete the entire directory. All legacy development process blueprints, Diátaxis TOML rulebooks, and documentation created to box agents into multi-step constraints are obsolete. The orchestrator drives agents on single scoped tasks with direct test and linter validation backpressure.
3. Documentation around hooks: Delete all legacy documentation or guides attempting to box agents into manual hook workflows. Enforcer hooks are replaced by single-task orchestrator pipelines.
4. `compile-docs.js` and `docs/docs-compiled.json`: Delete. The orchestrator now injects context and backpressure directly via the Ralph loop without compiling documentation to JSON indexes.
5. Obsolete helper copies: Delete root or `agent-scripts` copies of `update-action-versions.sh`, `update-modules.sh`, `run-in-nix.sh`, and `cleanup-data.js`.
6. Obsolete agent definitions: Delete `.gemini/agents/quality_assurance.toml`. Single-agent QA has been superseded by `runReviewRalphLoop`.
7. Root `package.json` test commands: Remove any `compile-docs.js` invocation from the `"test"` script.

---

## 3. Host Repository Contract (What the Repository Must Provide)

To function with `.gemini/`, the host repository must satisfy four standard contracts:

### C1: `GNUmakefile` Standard Interface

The host repository must provide a `GNUmakefile` (or `Makefile`) implementing these core targets:

- `make build`: Compiles workspace binaries or packages.
- `make unit`: Executes unit test suites.
- `make acc`: Executes acceptance or integration test suites.
- `make test`: Executes the default test harness (e.g. unit tests or full test suite).
- `make lint`: Executes all workspace static analysis tools across all codebase languages.
- `make fmt`: Auto-formats code across all supported formatters.

Optional extended targets (with automated fallback to core targets in `test.js` if omitted):

- `make compile`: Compiles both product and test binaries (falls back to `make build`).
- `make test-all`: Executes full test suites across all languages (falls back to `make test`).
- `make testaccrelay`: Executes relay/cloud acceptance tests (falls back to `make acc`).

### C2: Tooling & Environment Contract

- `flake.nix` & `flake.lock`: Must supply `nodejs_24`, `gemini-cli`, `gh`, `git`, `make`, and repository-specific toolchains.
- `.envrc`: Must activate the Nix shell and export `GEMINI_SYSTEM_MD=.gemini/system.md`.

### C3: Git and Plan Management Contract

- `plans/current.md`: Must be gitignored. The orchestrator tracks all execution plans and Ralph loop iterations in this file.
- `.gemini/tmp/`: Must be gitignored. All orchestrator logs (`orchestrator.log`) and transient review diffs reside here.
- Canonical Upstream Remote: The canonical upstream repository remote must be named `upstream`. The developer's fork remote must be named `origin`.

### C4: Dependency Setup Contract

To install the orchestrator SDK and dependencies, execute:

```bash
npm --prefix .gemini ci && node .gemini/orchestrator/install-sdk.js
```

The host repository does NOT need to be a Node.js project. Dependencies can be isolated entirely inside `.gemini/node_modules/`.

---

## 4. What `.gemini/` Does NOT Need from the Repository

- No root `package.json` requirement: Non-Node repositories (Go, Rust, Python, C) do not need root Node artifacts.
- No `.github/workflows/scripts/test.sh` or `lint.sh` dependency: The `code-validation` skill calls `make` targets directly.
- No hardcoded product directory assumptions: Product detection scans standard roots (`internal/`, `src/`, `lib/`, `pkg/`, `cmd/`) or custom paths in `.gemini/settings.json`.
- Flexible test boundary detection: Orchestrator test boundary enforcement accepts standard conventions across languages (`*_test.go`, `*.test.*`, `*.spec.*`, `test_*.*`, and test directories), while host test execution delegates to standard make targets (`make test`, `make unit`, `make compile`, `make test-all`).

---

## 5. Execution Reference

### Starting the Developer Orchestrator

```bash
node .gemini/orchestrator/orchestrator.js
```

### Running Code Validation Skill

```bash
# Full validation (make lint + make test)
node .gemini/skills/code-validation/scripts/validate.js

# Linting only (make lint)
node .gemini/skills/code-validation/scripts/lint.js all

# Unit tests only (make unit)
node .gemini/skills/code-validation/scripts/test.js unit

# Formatting (make fmt)
node .gemini/skills/code-validation/scripts/format.js all
```

### Running Git Operations

```bash
# Branch creation synced with canonical upstream
node .gemini/skills/git/scripts/create-branch.js <branch-name>

# Conventional commit with spell-checking and boundary enforcement
node .gemini/skills/git/scripts/commit.js -m "feat: descriptive title"
```

---

## 6. Migration Checklist for Autonomous Agents

When an agent enters a target repository to complete `.gemini/` consolidation:

- [ ] Step 1: Ensure `GNUmakefile` defines `build`, `unit`, `acc`, `test`, `lint`, and `fmt`.
- [ ] Step 2: Delete `agent-scripts/`, `docs/development/`, and legacy doc compilation scripts.
- [ ] Step 3: Run dependency setup: `npm --prefix .gemini ci && node .gemini/orchestrator/install-sdk.js`.
- [ ] Step 4: Verify test execution via `node .gemini/skills/code-validation/scripts/test.js unit`.
- [ ] Step 5: Verify static analysis via `node .gemini/skills/code-validation/scripts/lint.js all`.
- [ ] Step 6: Verify orchestrator dry-run via `node .gemini/orchestrator/inspect-prompts.js`.
