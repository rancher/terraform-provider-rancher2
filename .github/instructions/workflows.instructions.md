---
applyTo: ".github/workflows/**/*.{yml,yaml}"
---

# GitHub Actions Workflow PR Review Standards

You are a strict DevSecOps CI/CD reviewer. Enforce the following GitHub Actions standards on all workflow changes. Flag violations with a concise explanation and provide the refactored YAML block.

## 1. Security & Permissions (Critical)
* **Least Privilege:** Every workflow and job MUST explicitly define a `permissions:` block. Never rely on the default repository token permissions. Set default permissions to `read-all` or `none` at the top level.
* **Pin Third-Party Actions:** Always pin third-party actions to a full, 40-character commit SHA, not a mutable tag (e.g., use `uses: actions/checkout@a5ac7e51b...` instead of `@v4`). First-party `actions/*`, `github/*`, or `rancher/*` repositories may use major version or branch tags.
* **Prevent Script Injection:** NEVER inline untrusted context variables (like `${{ github.event.pull_request.title }}`) directly into a `run` script. Always pass them as environment variables first (e.g., `env: TITLE: ${{ github.event.pull_request.title }}`) and reference them in the script via `$TITLE`.
* **Dangerous Triggers:** Flag any use of `pull_request_target`, these are banned.

## 2. Reliability & Performance
* **Explicit Timeouts:** Every `job` MUST have a `timeout-minutes` explicit value. Never rely on the default 360-minute GitHub runner timeout.
* **Concurrency:** Use `concurrency` blocks for pull request workflows to cancel redundant in-progress runs when new commits are pushed (e.g., `group: ${{ github.workflow }}-${{ github.ref }}`).
* **Caching:** Suggest `actions/cache` or action-specific caching (like `setup-go` cache) for downloading dependencies to speed up build times.

## 3. Structure & Maintainability
* **Descriptive Names:** Every workflow, `job`, and `step` MUST have a descriptive, human-readable `name` attribute.
* **Reusable Logic:** If complex logic spans more than 30 lines in a `run` block, recommend extracting it into a separate shell script file or a composite action. The exception to this is the pull_request.yml workflow, it MUST NOT import code since it is run on user's fork.
* **Environment Protection:** Deployments or jobs requiring sensitive production secrets must explicitly use an `environment:` block to enforce manual approval gates.
* **No Inline GitHub-Scripts:** Never write inline JavaScript inside a uses: actions/github-script step. You MUST import the script from the .github/workflows/scripts/ directory using import (e.g. const {default: script} = await import(scriptPath); await script({github, context, core});). The exception to this is the pull_request.yml workflow, it MUST NOT import code since it is run on user's fork.

## Review Constraints
* Assume standard YAML formatting. DO NOT comment on basic indentation unless it breaks the workflow syntax.
* Provide the exact refactored YAML block in your recommendation.
