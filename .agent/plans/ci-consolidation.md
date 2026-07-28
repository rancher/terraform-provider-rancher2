# CI Workflow Simplification, Consolidation, and Standardization

- **Executed Date:** pending
- **Purpose:** Simplify, consolidate, and standardize all CI workflows and workflow scripts on a pre-built Nix container image, using a unified `nix-run.sh` wrapper, a consolidated `lint.sh` script, domain-focused consolidated JS scripts, and a standardized CSpell configuration to streamline workflow execution and remove boilerplate.

---

## Background & Rationale

Currently, all GitHub Actions workflows in this repository run on raw `ubuntu-latest` environments and repeatedly download and install Nix from scratch using a custom inline bash step (`install-nix` with curl). This approach introduces several inefficiencies:
1. **Redundant Download Time:** Every single job (e.g., `build`, `test`, `terraform`, `actionlint`, `shellcheck`, etc.) spends ~1–2 minutes downloading and installing Nix on every run.
2. **Boilerplate Bloat:** There is massive duplication of the installation step across multiple workflow files.
3. **Complex Inline Shell Invocation:** Invoking `nix develop` via custom `shell` configurations is verbose, hard to maintain, and prone to credential-leaking or environment mismatch issues.

Additionally, the repository contains a large number of single-purpose JavaScript helper scripts run by `actions/github-script` under `.github/workflows/scripts/`. Keeping them scattered across 14+ different files leads to directory clutter and fragmented logic.

Furthermore, the spell-checking utility `cspell` is run inside `validate-commit-message.sh` without a standard configuration file, relying on an ad-hoc custom word list file `aspell_custom.txt`. To align with standard practices, we need to introduce a proper `cspell.json` configuration, standard dictionary mapping to a renamed `custom_words.txt`, and resolve any potential early-exit bugs in the validation loop.

By adopting a consolidated, domain-driven script architecture, we can significantly simplify our pipeline. The proposed architecture introduces:
* **`nix-run.sh`**: A robust wrapper that sets up CA certificates, registers the workspace as a safe directory for Git inside the container, sets proper directory permissions for the non-root container user (`suse`), and runs arbitrary commands inside the Nix development shell (`nix develop`) with explicit variable forwarding.
* **`lint.sh`**: A centralized, mode-driven bash script that implements all lint and check routines under a unified entrypoint.
* **Consolidated JavaScript Modules**: Merging 14 separate JavaScript helper scripts into 3 cohesive, domain-specific modules (`backports.js`, `releases.js`, and `e2e.js`) which use a clean dispatcher pattern driven by a `SCRIPT_MODE` environment variable.
* **Full CSpell Implementation**: Standardizing spelling checks with a proper root-level `cspell.json` configuration and a `custom_words.txt` dictionary, while fixing the early-exit loop bug in `validate-commit-message.sh`.
* **Streamlined Yaml Workflows**: Jobs remain separate in the GitHub Actions UI so contributors get clear, granular feedback, but their definitions are reduced to simple, standardized steps. Every single step in every workflow is updated to ensure the `name` attribute is the first field and has a unique, descriptive `id` attribute.

---

## Goals

1. **Eliminate Nix Setup Boilerplate:** Replace all `install-nix` steps with container-based jobs.
2. **Centralize Lint Logic:** Move all lint commands (Terraform format, tflint, actionlint, eslint, node syntax check, shellcheck, gitleaks) into a single bash script `.github/workflows/scripts/lint.sh`.
3. **Consolidate JavaScript Scripts:** Group the 14 scattered GitHub-Script `.js` files into 3 domain-focused files:
   - **`backports.js`**: Merges `backport-issues.js`, `backport-pr.js`, `wait-for-settle.js`, and `merge-label.js`.
   - **`releases.js`**: Merges `check-maintainer.js`, `rc-notify.js`, `publish-release.js`, and `tracking-issue.js`.
   - **`e2e.js`**: Merges `check-ip.js`, `check-lock.js`, `check-run.js`, `clear-runner.js`, `wait-for-e2e.js`, and `report-e2e-status.js`.
4. **Implement Standard CSpell Config:** Rename `aspell_custom.txt` to `custom_words.txt` and introduce a proper root-level `cspell.json`.
5. **Standardize Nix Execution:** Wrap all Nix command executions inside the Nix development environment via `.github/workflows/scripts/nix-run.sh`.
6. **Clean up Pull Request Workflows:** Apply these changes to `.github/workflows/pull_request.yaml`.
7. **Clean up Other Workflows:** Propagate the same standard wrapper and consolidated scripts to other workflows (`release.yml`, `manual-rc-release.yml`, and `manual-release.yml`) to simplify their maintenance.

---

## Detailed Implementation Plan

### Step 1: Create `.github/workflows/scripts/nix-run.sh`
This script ensures SSL/TLS certificate bundles are set up and launches commands via `nix develop` as the preconfigured non-root `suse` user.

*(See earlier plan content for standard `nix-run.sh` code.)*

### Step 2: Create Centralized `.github/workflows/scripts/lint.sh`
This script will centralize all linters and check procedures.

*(See earlier plan content for standard `lint.sh` code.)*

### Step 3: Implement Standard CSpell Configuration & Fix Commit Linting
1. **Rename Dictionary File:** Rename `aspell_custom.txt` to `custom_words.txt`.
2. **Create `cspell.json`:**
   ```json
   {
       "dictionaries": ["custom-words"],
       "dictionaryDefinitions": [
           {
               "name": "custom-words",
               "path": "custom_words.txt",
               "addWords": true
           }
       ]
   }
   ```
3. **Update `.github/workflows/scripts/validate-commit-message.sh`:** Use `cspell stdin` and fix the early-exit bug in `spell_check` (change `exit 0` to `return 0` on `Merge ` match):
   ```bash
   spell_check() {
     message="$1"
     if grep -q -e '^Merge ' <<<"$message"; then
       return 0
     fi
     WORDS="$(cspell stdin <<<"$message")"
     if [ "" != "$WORDS" ]; then
       echo "...Commit message contains spelling errors on: ^$WORDS\$"
       echo "...Also try updating the PR title."
       echo "...If this is a mistake, add your word to the custom_words.txt file."
       exit 1
     else
       echo "...Commit message doesn't contain spelling errors."
     fi
   }
   ```

### Step 4: Consolidate JavaScript Workflow Scripts

To make our scripts clean, maintainable, and unified, we will consolidate them into three domain-specific files under `.github/workflows/scripts/`:

* **`backports.js`**
* **`releases.js`**
* **`e2e.js`**

*(See earlier plan content for ES module structure.)*

### Step 5: Update Workflow YAML Definitions
Rewrite the workflow steps to call the unified `.js` files using the `actions/github-script` action, explicitly setting the target subroutine using the `SCRIPT_MODE` environment variable.

### Step 6: Clean up Obsolete Files
Once the consolidated files are tested and verified, completely remove the 14 individual Javascript files to leave a beautifully clean `.github/workflows/scripts/` directory.

### Step 7: Address Copilot Review Comments
To resolve issues raised in pull request review comments by the Copilot bot, we will perform the following targeted refactorings:

1. **Argument boundary preservation in `nix-run.sh`:**
   Instead of writing commands using `printf "%s\n" "$*"`, which merges arguments into a single space-separated string and loses quoting, we will write them utilizing `printf "%q "` to shell-escape and preserve original argument boundaries and quoting when executing multi-argument CLI payloads (e.g. `nix-run.sh bash -c "..."`).
2. **Conditional Lock Artifact Upload in `release.yml`:**
   We will gate the `Create Lock Artifact` step inside the `test` job in `.github/workflows/release.yml` on the `if: steps.proceed.outputs.run == 'true'` condition to ensure locks are only uploaded/created when the test attempt actually claims execution.
3. **Enforce Strict Parameter Contracts (Fail-Fast):**
   Instead of using soft fallback defaults (such as `globalThis.process` or empty fallback objects) that mask configuration bugs, we will maintain strict parameter destructuring (`{ github, context, core, process }`) in `.github/workflows/scripts/e2e.js`, `backports.js`, and `releases.js` so that any missing parameters fail fast and loud.
4. **Comprehensive Workflow Call Audit:**
   We will audit all workflow files (`release.yml`, `backport-merge-label.yml`, etc.) and ensure they always pass all required context parameters (specifically `process`, `github`, `context`, and `core`) when invoking the consolidated JavaScript scripts (e.g., `await script({ github, context, core, process });`).
5. **Verify Conflict Markers:**
   Confirm that all temporary conflict markers are completely resolved in `.github/workflows/scripts/validate-commit-message.sh`.

---

## Verification and Safety Strategy

1. **Pre-commit Checks:** Validate `lint.sh` and `nix-run.sh` scripts locally using `shellcheck` to make sure they are error-free.
2. **GHA Action Validation:** Run `actionlint` locally (via the Nix shell) to verify that the refactored workflow files have valid syntax and meet GitHub Actions schema rules.
3. **Script Syntax & Import Checks:** Run syntax/compile validation check on the new JavaScript scripts to ensure there are no parser errors.
4. **CSpell Validation:** Test the updated commit validation script with both passing and intentionally misspelled commit messages to ensure it correctly blocks typos while passing valid messages.
