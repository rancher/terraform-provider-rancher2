# Temporary Plan: CI Workflow Simplification, Consolidation, and Standardization

This is a progress-tracking temporary plan to guide the iterative execution phase.

## Tasks & Checklist

- [x] **Step 1: Set up central execution and linting scripts**
  - [x] Create `.github/workflows/scripts/nix-run.sh` with correct environment variables and permissions.
  - [x] Create `.github/workflows/scripts/lint.sh` backing `terraform`, `actionlint`, `shellcheck`, `node-check`, `eslint`, and `gitleaks` options.
  - [x] Make both scripts executable.

- [x] **Step 2: Implement Standard CSpell Configuration**
  - [x] Rename `aspell_custom.txt` to `custom_words.txt`.
  - [x] Create root-level `cspell.json` referencing `custom_words.txt`.
  - [x] Refactor `.github/workflows/scripts/validate-commit-message.sh` to use `cspell stdin` and fix the early-exit loop bug (using `return 0` instead of `exit 0` on `Merge ` commits).

- [x] **Step 3: Consolidate JavaScript Scripts**
  - [x] Build `.github/workflows/scripts/backports.js` containing:
    - [x] `wait-for-settle` (from `wait-for-settle.js`)
    - [x] `backport-pr` (from `backport-pr.js`)
    - [x] `backport-issues` (from `backport-issues.js`)
    - [x] `merge-label` (from `merge-label.js`)
  - [x] Build `.github/workflows/scripts/releases.js` containing:
    - [x] `check-maintainer` (from `check-maintainer.js`)
    - [x] `rc-notify` (from `rc-notify.js`)
    - [x] `publish-release` (from `publish-release.js`)
    - [x] `tracking-issue` (from `tracking-issue.js`)
  - [x] Build `.github/workflows/scripts/e2e.js` containing:
    - [x] `check-ip` (from `check-ip.js`)
    - [x] `check-lock` (from `check-lock.js`)
    - [x] `check-run` (from `check-run.js`)
    - [x] `clear-runner` (from `clear-runner.js`)
    - [x] `wait-for-e2e` (from `wait-for-e2e.js`)
    - [x] `report-e2e-status` (from `report-e2e-status.js`)

- [x] **Step 4: Simplify and standardize `pull_request.yaml`**
  - [x] Update top-level configuration (remove redundant env/NIX variables).
  - [x] Update `build` job to run in container and use `nix-run.sh`.
  - [x] Update `test` job to run in container and use `nix-run.sh`.
  - [x] Update `terraform` job to use `nix-run.sh` and `lint.sh terraform`.
  - [x] Update `actionlint` job to use `nix-run.sh` and `lint.sh actionlint`.
  - [x] Update `node-check` job to use `nix-run.sh` and `lint.sh node-check`.
  - [x] Update `eslint` job to use `nix-run.sh` and `lint.sh eslint`.
  - [x] Update `shellcheck` job to use `nix-run.sh` and `lint.sh shellcheck`.
  - [x] Update `validate-commit-message` job to use `nix-run.sh` and `.github/workflows/scripts/validate-commit-message.sh`.
  - [x] Update `gitleaks` job to use `nix-run.sh` and `lint.sh gitleaks`.

- [x] **Step 5: Update other Workflows with Nix and JavaScript Consolidation**
  - [x] Update `release.yml` jobs to use pre-built Nix container, `nix-run.sh`, and consolidated JavaScript files (`releases.js` and `e2e.js`).
  - [x] Update `manual-rc-release.yml` to use pre-built Nix container, `nix-run.sh`, and consolidated JavaScript files.
  - [x] Update `manual-release.yml` to use pre-built Nix container, `nix-run.sh`, and consolidated JavaScript files.
  - [x] Update `backport-issues.yml`, `backport-merge-label.yml`, `backport-pr-manual.yml`, `backport-prs.yml`, and `rc-notifications.yml` to use the consolidated JavaScript files.

- [x] **Step 6: Clean Up Obsolete Files**
  - [x] Delete the 14 individual JS scripts under `.github/workflows/scripts/`.

- [x] **Step 7: Verify and Validate**
  - [x] Run `shellcheck` on `nix-run.sh` and `lint.sh` locally.
  - [x] Run `actionlint` locally on all updated workflows to verify their schema validity.
  - [x] Verify that there are no remaining instances of manual Nix install steps (`install-nix`) in the codebase.

- [x] **Step 8: Address Copilot Review Comments**
  - [x] Fix argument boundaries and quoting in `nix-run.sh` using shell-escaped `printf '%q '`
  - [x] Add conditional execution check `if: steps.proceed.outputs.run == 'true'` to `Create Lock Artifact` in `release.yml`
  - [x] Enforce strict parameter contracts (fail-fast on missing parameters) in `e2e.js`, `backports.js`, and `releases.js`
  - [x] Audit and update all workflow calls in `release.yml` (lines 144, 336, 451, 537) to ensure all required parameters (`process`, `github`, `context`, `core`) are explicitly forwarded
  - [x] Audit and update `backport-merge-label.yml` (line 38) to explicitly forward all required parameters (`process`, `github`, `context`, `core`)
  - [x] Verify that no merge conflict markers exist in `validate-commit-message.sh` (already verified clean)
