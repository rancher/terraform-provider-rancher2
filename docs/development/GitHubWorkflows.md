# GitHub Workflows & Automation

This topic overview details the architecture, standards, and runtime security policies governing GitHub Actions workflows and automation scripts within the repository.

---

## Abstract

GitHub Actions workflows are the central orchestration engine of our repository's automation pipeline. By strictly separating workflow orchestration (declaring the execution model) from the actual implementation scripts, and applying platform-level execution policies, we ensure that our automations are highly modular, testable, and secure.

---

## 🧭 How Our Workflow Components Work Together

Our workflow and automation ecosystem is built on top of a highly secure, dual-layered architecture:

### 1. Unified Script Orchestration

Workflows in this repository act as orchestrators rather than executors. Any non-trivial logic—such as Git tagging, manifest scans, or PR validation—is extracted into dedicated, testable JavaScript or Bash scripts under `.github/workflows/scripts/`. This separation allows us to run standard unit tests locally and verify our scripts dynamically without relying on the GHA runner environment.

### 2. Platform-Level Security Policies & Protections

To prevent malicious script execution and secure our privilege boundary against unreviewed fork contributions (such as poisoned pipeline attacks), we enforce centralized GitHub Rulesets and execution allowlists.

- More details on configuring these actor and event restrictions, setup parameters, and threat mitigations can be found in **[Workflow Execution Protections](./GitHubWorkflows/WorkflowExecutionProtections.md)**.

---

# Architectural Blueprint & Specification

## 1. Workflow Architecture & Best Practices

To ensure high maintainability, security, and velocity, all GitHub Workflows in this repository MUST comply with the following architectural rules:

### A. Orchestrate, Don't Execute

Workflows should act as orchestrators, not execution scripts.

- All non-trivial logic (e.g., git tagging, GPG lookup, maintainer checks, PR verification) must live in external scripts inside the `.github/workflows/scripts/` directory.
- Workflows must call these external scripts using the Nix environment (`.github/workflows/scripts/nix-run.sh`) or native `actions/github-script` runners.

### B. Security & Least Privilege

- **Top-Level Scopes:** Every workflow MUST have `permissions: {}` defined at the top-level.
- **Job-Level Permissions:** Each job MUST define explicit, minimal `permissions` required for its operation (e.g., `contents: read`, `pull-requests: write`).
- **Pin Actions by SHA:** All GitHub Actions must be pinned to a full 40-character commit SHA (not a tag). Include a release URL comment on the line before `uses:` (e.g., `# https://github.com/actions/checkout/releases`).

### C. Execution Safety

- **Timeouts:** Every job MUST have `timeout-minutes` explicitly set (default: `5`).
- **Graceful Fallbacks:** Scripts that communicate with the GitHub API must incorporate robust retry logic and handle API rate limit exceptions gracefully.

---

## 2. Shared Automated Script Standards

All helper and runner scripts in `.github/workflows/scripts/` must strictly comply with language-specific rules defined under our **Coding Standards** topic:

- **JavaScript Scripts**: Must use ESM modules, handle paginated lists via `github.paginate()`, and be fully unit-tested with the native Node.js test runner.
- **Bash Scripts**: Must use double brackets `[[ ]]`, fail-fast with `set -euo pipefail`, and handle error logging to stderr.
