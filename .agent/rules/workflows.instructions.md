---
applyTo: ".github/workflows/**/*.{yml,yaml}"
---

# GitHub Actions Workflow PR Review Standards

As a strict DevSecOps CI/CD reviewer, enforce these standards on all workflow changes. Flag violations with a concise explanation and provide the refactored YAML.

## 1. Security (Critical)
* **Least Privilege:** All jobs must define explicit `permissions:`. All workflows should have `permissions: {}` at the top level. Set scopes to `none` as needed. Permissions should implement least privilege necessary access.
* **Pin Actions by SHA:** Pin all actions (including `actions/*`, `github/*`, `rancher/*`) to a full 40-character commit SHA, not a tag. The `uses:` line MUST include the version (e.g., `# v6.0.2`). On the line before the `uses:` there should be a comment with a link to the releases page for the action (e.g. `# https://github.com/actions/github-script/releases`).
* **Approved Action Namespaces:** Only pre-approved action namespaces are allowed. Approved namespaces are documented at: https://github.com/rancher/security-team/blob/main/docs/standards/rancher-gha-standards.md#allowed-github-actions. Important ones include: `actions/*`, `aquasecurity/*`, `aws-actions/*`, `dependabot/*`, `fossas/*`, `golang/*`, `golangci/*`, `google-github-actions/*`, `google/*`, `googleapis/*`, `goreleaser/*`, `hashicorp/*`, `rancher-eio/*`, `renovatebot/*`, and `updatecli/*`.
* **Prevent Script Injection:** Never inline untrusted context variables in `run` scripts. Use environment variables (e.g., `env: VAR: ${{...}}`).
* **No `pull_request_target`:** This trigger is banned.

## 2. Reliability & Performance
* **Explicit Timeouts:** Every `job` must have an explicit `timeout-minutes`. Don't use the 360-minute default.
* **Concurrency:** Use `concurrency` blocks in PR workflows to cancel redundant runs (e.g., `group: ${{ github.workflow }}-${{ github.ref }}`).
* **Caching:** Suggest `actions/cache` or action-specific caching to speed up dependency downloads.

## 3. Structure & Maintainability
* **Orchestrate, Don't Execute:** Workflows should orchestrate, not execute. They may call out to external actions or internal scripts, but must not execute full steps by themselves.
* **External Scripts:** All `run` or `github-script` scripts should be placed in the `.github/workflows/scripts` directory. Prefer using external scripts over inline JavaScript in `actions/github-script` for complex logic. Minimal inline bootstrapping (such as importing and calling an external script or loading dependencies) is allowed where necessary.
* **Script Validation:** All scripts should be validated in the `pull_request.yaml` workflow.
* **Descriptive Names:** All workflows, jobs, and steps need a descriptive `name`.
* **Environment Protection:** Jobs with production secrets must use an `environment:` block for manual approval.

## Review Constraints
* Ignore basic YAML formatting unless it's a syntax error.
* Provide the exact refactored YAML block in your recommendation.
