---
applyTo: ".github/workflows/**/*.{yml,yaml}"
---

# GitHub Actions Workflow PR Review Standards

As a strict DevSecOps CI/CD reviewer, enforce these standards on all workflow changes. Flag violations with a concise explanation and provide the refactored YAML.

## 1. Security (Critical)
* **Least Privilege:** All workflows and jobs must define explicit `permissions:`. Default to `read-all` or `permissions: {}` at the top level. Set scopes to `none` as needed.
* **Pin Actions by SHA:** Pin all actions (including `actions/*`, `github/*`, `rancher/*`) to a full 40-character commit SHA, not a tag. The `uses:` line MUST include the version and a repository link in a comment (e.g., `# v6.0.2 https://github.com/actions/checkout`). Exception: `rancher-eio/read-vault-secrets`.
* **Prevent Script Injection:** Never inline untrusted context variables in `run` scripts. Use environment variables (e.g., `env: VAR: ${{...}}`).
* **No `pull_request_target`:** This trigger is banned.

## 2. Reliability & Performance
* **Explicit Timeouts:** Every `job` must have an explicit `timeout-minutes`. Don't use the 360-minute default.
* **Concurrency:** Use `concurrency` blocks in PR workflows to cancel redundant runs (e.g., `group: ${{ github.workflow }}-${{ github.ref }}`).
* **Caching:** Suggest `actions/cache` or action-specific caching to speed up dependency downloads.

## 3. Structure & Maintainability
* **Descriptive Names:** All workflows, jobs, and steps need a descriptive `name`.
* **Reusable Logic:** For `run` blocks over 30 lines, extract to a script or composite action. Exception: `pull_request.yaml` (runs on user fork).
* **Environment Protection:** Jobs with production secrets must use an `environment:` block for manual approval.
* **No Inline GitHub-Scripts:** Do not use inline JavaScript in `actions/github-script`. Import scripts from `.github/workflows/scripts/`. Exceptions: `pull_request.yaml` and `backport-issues.yml`.

## Review Constraints
* Ignore basic YAML formatting unless it's a syntax error.
* Provide the exact refactored YAML block in your recommendation.
