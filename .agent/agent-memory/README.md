# AI Agent Memory

This directory contains persistent context and learnings to retain across AI sessions.

Agents can read from and write to this directory to remember past decisions, project-specific quirks, and user preferences, improving future interactions.

Here is an example to use as a template:

A file in this directory named "WorkflowStandards-TemporaryPlan.md" which tracks progress of a plan named "WorkflowStandards.md" in the plans directory.
```
# Temporary Plan: Workflow Standards

**Status:** Completed
**Purpose:** Track progress of applying `WorkflowStandards.md` rules to `.github/workflows/*`

## Global / Scripts Setup
- [x] Create `.github/workflows/scripts` directory.
- [x] Ensure `.github/workflows/validate.yaml`'s `shellcheck` job includes validation for scripts in `.github/workflows/scripts/`.

## `fossa.yml`
- [x] Change top-level `permissions` to `permissions: {}`.
- [x] Move `contents: read` and `id-token: write` to the `fossa-scanning` job's explicit permissions block.
- [x] Ensure all steps contain both a `name` and an `id`.

## `release.yaml`
- [x] Change top-level `permissions: write-all` to `permissions: {}`.
- [x] Assign explicit, granular `permissions:` at the job level for `release`, `test`, `cleanup`, and `report`.
- [x] Ensure all steps contain both a `name` and an `id`.
- [x] Extract inline `github-script` comment script in the `release` job to `.github/workflows/scripts/release-pr-comment.js`.
- [x] Extract inline `run` block in the `test` job to `.github/workflows/scripts/run-tests.sh`.
- [x] Extract inline `run` block in the `cleanup` job to `.github/workflows/scripts/cleanup-tests.sh`.
- [x] Extract inline `github-script` block in the `report` job to `.github/workflows/scripts/report-success.js`.

## `validate.yaml`
- [x] Add a `concurrency:` block to cancel redundant PR runs.
- [x] Ensure all steps contain both a `name` and an `id`.
- [x] Fix context injection in `validate-commit-message`: use an environment variable for `${{github.event.number}}` rather than inlining it into the script.
- [x] Extract `lint terraform` inline run block to `.github/workflows/scripts/lint-terraform.sh`.
- [x] Extract `action lint` inline run block to `.github/workflows/scripts/action-lint.sh`.
- [x] Extract `shell check` inline run block to `.github/workflows/scripts/shell-check.sh`.
- [x] Extract `Check commit message` inline run block to `.github/workflows/scripts/check-commit-message.sh`.
- [x] Extract `Check for secrets` inline run block to `.github/workflows/scripts/gitleaks-scan.sh`.
- [x] Extract `test-compile-check` inline run block to `.github/workflows/scripts/test-compile-check.sh`.
- [x] Extract `lint-tests` inline run block to `.github/workflows/scripts/lint-tests.sh`.
```
