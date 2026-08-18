# Plan: Release Process

* **Executed Date:** pending
* **Purpose:** Establish and enable the newly designed "Standard Repository Release Process: Architectural Blueprint & Tooling Specification" as our repository's standard. This involves documenting the standard in `RELEASING.md`, communicating it in `README.md`, and making the necessary codebase updates to convert our GHA and script logic to this new process.
* **Goals & Code Snippets:**

---

# Standard Repository Release Process: Architectural Blueprint & Tooling Specification

## **1. Actor-vs-Automation Interaction Swimlanes**

This swimlane diagram traces the detailed event triggers and data flow between development roles and automated GHA runners:

```
|                [ Actor ]                    |                  [ Automation (GHA / Nix / Vault) ]            |
├─────────────────────────────────────────────┼────────────────────────────────────────────────────────────────┤
|                                             |                                                                |
| === PART 1: PULL REQUEST TO SQUASH-MERGE ===|                                                                |
|                                             |                                                                |
|  1. PR Opened or Code Updated ------------> | ──► Trigger: pull_request (opened / synchronize)               |
|     (Contributor submits change)            |    ├─► Native Copilot Review triggers automatically            |
|                                             |    └─► Nix: Runs static tests, linting, and unit-checks        |
|                                             |         │                                                      |
|                                             |         ▼ (Checks Completed successfully)                      |
|                                             |       Trigger: workflow_run (completed)                        |
|                                             |       ──► Coordinator executes (Dry-Run Mode)                  |
|                                             |           - Scans approval state                               |
|                                             |           - Finds: No trusted reviews or approvals yet         |
|                                             |           - Action: Posts comment: "PR Needs Collaborator Review"|
|                                             |                                                                |
|  2. Collaborator Reviews PR & Comments ---->| ──► Trigger: pull_request_review (submitted)                   |
|     (Collaborator leaves feedback/questions)|       ──► Coordinator executes (Merge Mode)                    |
|                                             |           - Finds: No approvals, unresolved comment threads    |
|                                             |           - Action: Posts status comment: "Unresolved Comments"|
|                                             |                                                                |
|  3. Contributor Resolves Comments           |                                                                |
|     (Option A: Fixes code & pushes commit)  |                                                                |
|     └─► Pushes Commit ────────────────────> | ──► Re-runs Step 1 (Nix Unit Tests -> Dry-Run check)           |
|                                             |                                                                |
|     (Option B: Resolves via browser)        |                                                                |
|     └─► Resolves conversation on GitHub     |      [ No GHA trigger fired for "clicking resolve" ]           |
|                                             |                                                                |
|     (Option C: "Pokes" via comment)         |                                                                |
|     └─► Types top-level or thread comment -> | ──► Trigger: issue_comment / pull_request_review_comment       |
|         (e.g., "resolved", "ready")         |       ──► Coordinator executes (Merge Mode)                    |
|                                             |           - GraphQL: Queries all review threads                |
|                                             |           - Finds: All threads are marked resolved             |
|                                             |           - Action: Updates status comment, waits for Approval │
|                                             |                                                                |
|  4. Collaborator Submits Final Approval ──> | ──► Trigger: pull_request_review (submitted)                   |
|     (Trusted Collaborator clicks "Approve") |       ──► Coordinator executes (Merge Mode)                    |
|                                             |           - GQL Check: Confirms 100% of comment threads resolved│
|                                             |           - Review Check: Validates ≥ 1 Collaborator Approval  |
|                                             |           - Action: Fires Proxy Approval (vouching for review) │
|                                             |           - Action: Fires SemVer Guard (scopes file boundary)  |
|                                             |           - Action: Sanitizes commit message (product-safe)    |
|                                             |           - Action: Executes NATIVE AUTO-MERGE into 'main'     |
|                                             |           - Action: Automatically deletes status comments      |
|                                             |                                                                |
| === PART 2: SQUASH-MERGE TO PRODUCTION RELEASE =================─────────────────────────────────────────────┤
|                                             |                                                                |
|  5. Squash Merge Lands on main ───────────> | ──► Trigger: push to main                                      |
|                                             |       ──► Release Please Action runs:                          |
|                                             |           - Scans Conventional Commit squash titles            |
|                                             |           - Calculates next version increment                  |
|                                             |           - Action: Updates/creates draft "Release PR"         |
|                                             |             (e.g., "chore: release v1.2.3")                    |
|                                             |                                                                |
|                                             | ──► Trigger: pull_request (targeting Release Please branch)    |
|                                             |       ──► Release PR CI runs:                                  |
|                                             |           - Nix: Runs Unit Tests                               |
|                                             |           - AWS: Assumes OIDC IAM Role                         |
|                                             |           - Nix: Executes FULL ACCEPTANCE TEST SUITE           |
|                                             |                  (acc-relay: Real AWS resource deployment)      |
|                                             |         │                                                      |
|                                             |         ▼ (Heavy Integration Tests Pass)                       |
|                                             |       ──► Release Candidate (RC) Step:                         |
|                                             |           - Calculates next RC tag (e.g., v1.2.3-rc.0) via API |
|                                             |           - Vault: Securely extracts GPG keys                  |
|                                             |           - Nix: GoReleaser compiles & signs RC binaries       |
|                                             |           - Action: Publishes GPG-Signed RC Release            |
|                                             |                                                                |
|  6. Maintainer Merges Release PR ─────────> | ──► Trigger: push to main (Release PR merged)                  |
|     (Maintainer approves/merges Release PR) |       ──► Release Please Action runs:                          |
|                                             |           - Detects Release PR merge                           |
|                                             |           - Action: Outputs: release_created = true            |
|                                             |         │                                                      |
|                                             |         ▼                                                      |
|                                             |       ──► Full Release Step:                                   |
|                                             |           - Action: Automatically tags version (v1.2.3)        |
|                                             |           - Vault: Securely extracts GPG credentials           |
|                                             |           - Keyring Workaround: Dynamically parses primary ID  |
|                                             |           - Nix: GoReleaser cross-compiles stable binaries     |
|                                             |           - Action: Publishes Final GPG-Signed Release         |
|                                             |           - Action: Reconciles Release PR labels in GitHub     |
|                                             |                                                                |
```

---

## **2. Core Automation & Security Tooling Specifications**

To enforce a zero-trust, reproducible release pipeline, the repository utilizes four foundational tools: **Nix**, **The CI-Image**, **Release Please**, and **GoReleaser**.

### **A. Nix: Zero-Trust Hermetic Reproducibility**
* **Purpose:** Nix acts as a declarative package manager that defines the **exact build and test environment** down to the cryptographic hash. It locks the compiler, linter, runtime, and CLI utilities (Go, Node.js, Terraform, actionlint, etc.) in `flake.nix`.
* **Security & Automation Contribution:**
  - **Eliminates Environment Drift:** Standardizes the toolchain so that a developer running a test locally uses the *exact same binary byte-code* as the CI runner, eliminating "works on my machine" failures.
  - **Zero-Dependency Host Runners:** The GitHub Actions runner does not need pre-installed tools. Nix fetches and isolates everything inside a sandbox, securing the pipeline against malicious or outdated runner environments.

### **B. The CI-Image: Pre-Built Dependency Caching**
* **Purpose:** A pre-built Docker image (`ghcr.io/rancher/ci-image/nix`) containing a base Nix environment and pre-cached tool dependencies.
* **Security & Automation Contribution:**
  - **Time Optimization:** Bootstrapping Nix and compiling developer tools on every GHA workflow run can take several minutes. The pre-built CI image slashes initialization overhead to under **15 seconds**.
  - **Immutable Runtime Environment:** By freezing the CI-Image version (e.g. `nix:20260603-18`), the project secures its pipeline against supply-chain updates and runtime image modifications.

### **C. Release Please: Declarative Versioning & Changelog Automation**
* **Purpose:** An automated release management engine that parses Conventional Commits (`feat:`, `fix:`, `chore:`) to calculate Semantic Versioning (SemVer) jumps.
* **Security & Automation Contribution:**
  - **Manual Versioning Eradication:** Completely automates version calculations and generates high-fidelity changelogs.
  - **The Release PR Pattern:** Instead of tagging immediately on merge, it maintains a long-lived "Release PR" that acts as a staging queue. This allows the team to inspect version jumps and provides a physical gateway where final integration tests are executed.

### **D. GoReleaser: Automated Compiling, Packaging & GPG Signing**
* **Purpose:** A release automation engine designed to build, package, sign, and publish compiled binaries (such as Terraform providers) for multiple CPU architectures and Operating Systems.
* **Security & Automation Contribution:**
  - **Cryptographic Signing (GPG):** GoReleaser integrates with local GPG keys (securely pulled from Vault in memory) to cryptographically sign provider binaries and generate SHA256 checksums. This guarantees to the Terraform Registry that the binary has not been tampered with since compilation.
  - **Standardized Multi-Platform Matrixing:** Automatically cross-compiles for `linux`, `darwin`, and `windows` across `amd64` and `arm64` in a single, atomic step.
  
#### **🔧 Security Engineering Tip: GPG Key ID Extraction Workaround**
Static configuration of a GPG Key ID in secrets managers often leads to breaking releases. For example, if a key is rotated, or if Vault is accidentally configured with an encryption subkey ID rather than the primary signing key ID, GoReleaser will abort with a signing failure.

**The Workaround:** 
To handle this key-matching weirdness, the pipeline imports the raw secret key and then **dynamically inspects the GPG key-ring in real-time** to extract the true primary signing key ID (`sec`). It overrides any static `GPG_KEY_ID` configuration with the dynamically detected ID:
```bash
# 1. Strip whitespace/spaces from static GPG_KEY_ID config
export GPG_KEY_ID=$(echo -n "${GPG_KEY_ID}" | tr -d '[:space:]')

# 2. Import raw key block into local keyring
echo "${GPG_KEY}" | gpg --import --batch

# 3. Query the GPG keyring to extract the actual imported primary secret key ID (sec)
SEC_LINE=$(gpg --batch --list-secret-keys --keyid-format LONG | grep -E '^sec' | head -n1 || true)
if [[ -n "${SEC_LINE}" ]]; then
  # Parse out the key ID after the '/' separator
  DETECTED_KEY_ID=$(echo "${SEC_LINE}" | awk '{print $2}' | cut -d'/' -f2)
  if [[ -n "${DETECTED_KEY_ID}" ]]; then
    # Overwrite the static variable with the actual imported ID
    GPG_KEY_ID="${DETECTED_KEY_ID}"
  fi
fi
```
This guarantees that GPG signing never fails due to subkey mismatches, whitespace issues, or misconfigured key identifiers.

---

## **3. Standard Phase Specifications**

### **Phase 1: Zero-Trust Pull Request Checking (CHECK)**
* **PR Opened:** A contributor submits a Pull Request targeting `main`.
* **Checked:** The standard PR checkers run inside the hermetic **Nix** shell on GHA. This executes static code linters (`golangci-lint`, `actionlint`, `shellcheck`, `gitleaks`) and runs localized unit tests with zero-trust permissions.
* **Copilot Review:** Natively triggered repository integration initiates automated AI review comments.

---

### **Phase 2: Secure, Event-Driven Merge Coordination (COORD)**
* **Event Triggered:** Completion of the checks or reviews initiates the event coordinator. This executes on `workflow_run` in the secure default branch context (`main`), protecting secrets while enabling write-level access.
* **Validated Reviews/GQL:** The coordinator checks that the PR requirements are satisfied:
  * **Standard Pull Requests:** Requires **at least 1 approval** from a trusted role (Collaborator, Member, Owner, or Triage permission) and runs GraphQL queries to guarantee **100% of all review comments are marked resolved** (whether left by humans or AI).
  * **Dependabot Pull Requests:** Bypasses human reviewer constraints. Allows auto-merging with **at least 1 AI review approval/comment** (e.g. from Copilot) once all other functional check runs have completed successfully.
* **Proxy Approval:** If the requirements are met, but the PR lacks a Write-level approval (e.g., the approving reviewer has Triage-level access, or it is a Dependabot PR approved by AI), the GHA bot automatically submits an `APPROVE` review on the PR. Since the bot has Write access, its approval satisfies GitHub's branch protection requirements, serving as a proxy to allow the merge.

---

### **Phase 3: Automated SemVer Guard & Native Auto-Merge (MERGE)**
* **Scoped Boundary Check:** Evaluates modified files. If changes are exclusively non-product (e.g. docs, tests, CI files outside the core `internal/` directory), the SemVer Guard is activated.
* **Title Sanitized:** If SemVer Guard is active, conventional commit types like `feat` or `refactor` and breaking indicators (`!`) are dynamically stripped or downgraded to `chore` or `fix` to prevent unintentional Minor or Major version bumps.
* **Native Auto-Merge:** The PR is merged using GitHub's native Auto-Merge backend (`gh pr merge --auto --squash`) with custom, AI-sanitized commit messages. This cleanly bypasses the REST API GITHUB_TOKEN merge limitation for fork-authored PRs while respecting all branch protection settings.
* **Concurrency & Exploitation Protection:** If a contributor pushes a new commit *after* native auto-merge is enabled, GitHub's native system instantly pauses/cancels the auto-merge because previous approvals are automatically dismissed and new checks are triggered. The event coordinator (`pr-executor.yml`) must re-evaluate and re-verify the new commit SHA before auto-merge can be re-enabled.

---

### **Phase 3b: Direct Fork PR Auto-Merge (The Ruleset Bypass)**
* **Fork PR Evaluation:** When `pr-executor.yml` (running under `workflow_run`) evaluates a fork PR, it has full write/merge privileges on the base repository. If all requirements are met, it executes the native squash-merge using the GitHub CLI (`gh pr merge --auto --squash`).
* **The Recursion Bypass:** Normally, a GITHUB_TOKEN push (or merge) event prevents downstream workflows (such as Release Please) from triggering due to recursion prevention rules. However, the repository utilizes **GitHub Workflow Execution Protections (Rulesets)** to explicitly allow `github-actions[bot]` to trigger push-based workflows on `main`, ensuring downstream release pipelines run seamlessly.
* **🔧 Ruleset Configuration Prerequisites & Guardrails:**
  * **Location:** Repository settings -> *Actions > Policies*.
  * **Restrict Events:** Permit the **`push`** event to trigger actions (so the squash-merge commit on `main` initiates downstream release workflows). Keep other high-risk events (like `pull_request_target`) strictly restricted.
  * **Restrict Actors:** Add the system-level **`GitHub Actions` App** (representing the GITHUB_TOKEN) and your team's authorized roles (like `Write`, `Maintain`, `Admin` or the `k3s` group) as allowed actors to trigger the `push` event.
  * **Scope:** Ensure this ruleset is active on the repository (Actions Policies apply globally and do not use branch selectors).

---

### **Phase 4: Release Management & The Integration Test Gate (GATES)**
* **Release PR Maintained:** Merging into `main` triggers `Release Please`. It calculates the next version and automatically updates a draft "Release PR" containing updated version coordinates and changelogs.
* **Integration Tests Gate:** The Release PR acts as a staging queue where a dedicated CI workflow executes the **full integration and acceptance test suite** (using real cloud resources/relays). Merging is blocked until this suite passes.
* **Release PR Merged:** Merging this PR triggers the final release process.

---

### **Phase 5: Release Tagging & Cryptographic Signing (PUBLISH)**
* **Tagged vX.Y.Z:** `Release Please` registers the release merge and automatically creates and pushes the official semantic git tag.
* **Key Extracted:** The workflow imports the signing GPG block from Vault and dynamically queries the keyring to extract the actual primary secret key ID to work around subkey/mismatch weirdness.
* **Signed & Published:** GoReleaser compiles binaries inside the reproducibly locked Nix environment, signs them with the GPG key, and publishes the signed provider assets directly to the GitHub Release.
* **Label Reconciliation:** Since `skip-github-release: true` is configured in `release-please-config.json` to let GoReleaser manage the release, `release-please` skips post-merge tagging/labeling actions. To prevent release PRs from being left with outdated/pending labels, the publish-release script automatically reconciles labels on the merged Release PR (removing `autorelease: pending` and `ready-to-merge`, and adding `autorelease: tagged`).

---

### **Phase 6: Manual Release Escape Hatches (MANUAL)**
The `manual-release.yml` and `manual-rc-release.yml` workflows let a maintainer tag and publish out-of-band, bypassing Release Please. They mirror Phase 5 (tag → Vault GPG → GoReleaser → publish), but the maintainer supplies the tag by hand via `workflow_dispatch`.

#### **🔧 Engineering Rule: Container Jobs Must Never Consume `${{ github.workspace }}` in `run:` Steps**
Every release job now executes inside the pinned CI-Image (`container:` block). This creates two distinct, non-interchangeable views of the same directory:

| Context | Value | Valid where |
| --- | --- | --- |
| `${{ github.workspace }}` (expression) | `/home/runner/work/<repo>/<repo>` | Host runner only |
| `$GITHUB_WORKSPACE` (env var) | `/__w/<repo>/<repo>` | Inside the container |

The expression is evaluated by the runner on the **host** before the command is handed to the container, so it always yields the host path. The host path is not mounted at that location inside the container, so any `cd`/`test`/redirect against it fails.

* **Rule:** Inside a containerized job, `run:` steps MUST reference `$GITHUB_WORKSPACE`. `${{ github.workspace }}` MUST NOT be injected into `run:` steps via `env:`.
* **Actions are exempt:** `actions/checkout` translates host paths to container paths internally, which is why a `path:` input built from `${{ github.workspace }}` silently keeps working and masks the bug. Prefer a **relative** `path:` (relative to `$GITHUB_WORKSPACE`) so both views agree.

#### **🔧 Engineering Rule: `grep` Filters Under `set -e` Require `|| true`**
Container `run:` steps default to `sh -e {0}` (the runner falls back from `bash` when it is absent from the image). Under `-e`, an assignment whose command substitution fails aborts the step — so `var=$(cmd | grep -v ...)` is fatal whenever `grep` matches nothing and exits 1.

This is latent in the tag-pruning logic. `actions/checkout` fetches a tag ref using the refspec glob `refs/tags/<tag>*`, so the fresh checkout usually contains only the target tag and `grep -v "^$TAG$"` legitimately returns empty. It survived historically by coincidence: for `v13.2.0-rc.1` the glob also matched `v13.2.0-rc.10`, `rc.11`, …, leaving residue for `grep -v` to print. A tag with no numeric siblings (e.g. `v15.0.0-rc.13`) produces an empty result and kills the step.

* **Rule:** Any `grep`/`comm`/`diff` filter feeding a command substitution in a `sh -e` step MUST be suffixed with `|| true`, with the empty case handled explicitly by the following `if [ -n "$var" ]` guard.

#### **🔧 Engineering Rule: Crossing the `nix-run.sh` Boundary**
`.github/workflows/scripts/nix-run.sh` launches `nix develop --ignore-environment` with an explicit `--keep` allowlist (`NIX_SSL_CERT_FILE`, `HOME`, `GITHUB_TOKEN`, `GPG_*`, `AWS_*`, …). Anything absent from that allowlist is **empty inside the Nix shell**. The script also serialises its arguments with `printf %q`, so a deferred `\$VAR` reaches the inner shell as a literal `$VAR` and expands there, against the sanitized environment.

* **Rule (preferred):** A shell variable consumed by a command passed *through* `nix-run.sh` MUST be added to the `--keep` allowlist in `nix-run.sh`, and its expansion deliberately **deferred** to the inner shell (`\"\$VAR\"`). `GITHUB_WORKSPACE` and `TAG` are on the allowlist for exactly this reason. Deferring is both the clearest expression of intent and the safer default — see the security rule below.
* **Rule (security, non-negotiable):** An outer-shell-expanded value MUST NEVER be interpolated into a string that the inner `bash -c` will parse. `bash -c "cd \"$PATH\" && cmd"` re-parses the value as shell source, so an attacker-controlled `inputs.tag` such as `v1.2.3-rc.1";echo PWNED;"` executes arbitrary code with `GITHUB_TOKEN` and the Vault GPG secrets in scope — `git check-ref-format` permits `;`, `"` and `$(…)`, and the `Verify Tag` gate only greps for the substring `rc`. Deferred expansion is immune: parameter expansion inside double quotes substitutes the value, it does not re-parse it. If a value genuinely cannot be passed via the environment, pass it as a **positional parameter** (`bash -c 'cd "$1" || exit 1; …' _ "$path"`), never by interpolation.
* **Failure mode if the allowlist regresses:** dropping a variable from `--keep` degrades the expansion to an empty string — e.g. `cd "$GITHUB_WORKSPACE/tags/$TAG"` becomes `cd /tags/`. Verified to fail loudly (`cd` returns non-zero, `&&` short-circuits, the step exits 1) rather than silently releasing from the wrong directory. Always chain the real work behind `&&` or `|| exit 1` so this remains true: the inner `bash -c` does **not** inherit the outer `set -e`.

---

## Implementation Checklist

### Phase 3: Sequential Implementation (Act)
- [x] Create `RELEASING.md` in the root of the repository containing the co-designed release standard documentation.
- [x] Add a prominent section/link in `README.md` pointing to `RELEASING.md`.
- [x] Update `.github/workflows/review-trigger.yml` to trigger on comment events (`pull_request_review_comment` and `issue_comment` types: `created, edited, deleted`) to support the comment-based re-evaluation ("poke") mechanism.
- [x] Refactor `.github/workflows/scripts/verify-pr-requirements.js` to change the review threshold from requiring `2 humans OR 1 human + 1 AI` to **at least 1 human approval** (with trusted access), since AI review is natively triggered at the repository level now.
- [x] Refactor `.github/workflows/scripts/verify-pr-requirements.js` to fail the verification check run if the PR is in draft mode in non-auto-merge mode.
- [x] Rename `verify-pr-requirements.js` to `verify-pr-requirements.mjs` and update `pull_request.yaml` and `pr-executor.yml` references to ESM `.mjs` extension to resolve runner syntax error.
- [x] Resolve syntax error in `verify-pr-requirements.mjs` caused by hybrid code replacement.
- [x] Remove `triggerAIReviewIfNeeded` function and references from `verify-pr-requirements.mjs` since Copilot review is natively repository-triggered.

### Phase 4: Testing, Verification & Proactive Review (Quality Gate 1)
- [x] Run `actionlint` locally to verify that `.github/workflows/review-trigger.yml`, `pull_request.yaml` and `pr-executor.yml` syntax remains valid.
- [x] Run `node --check` to verify that `.github/workflows/scripts/verify-pr-requirements.mjs` syntax remains valid.
- [x] Verify that markdown formatting and hyperlinks are valid in both `RELEASING.md` and `README.md`.
- [x] Enter Proactive Review Mode on the written code diff to proactively resolve any issues human or Copilot automated reviews might flag.

### Phase 5: Chunking & Staging Isolation (Quality Gate 2)
- [x] Present the unstaged diff to the developer in the chat.
- [x] Solicit manual developer review and obtain explicit approval in the chat.

### Phase 6: Authorized Commit & PR Generation (Quality Gate 3)
- [x] Stage only the specific modified/created files (no `git add .` or `-A`).
- [x] Commit changes locally with a conventional prefix (e.g., `ci: implement standard release process and documentation`) using the secure `commit-push.sh` skill (or signed commit fallback).
- [x] Push the branch `feature/document-standard-release` to the user's origin fork.
- [x] Generate a Draft Pull Request targeting upstream `main` in draft mode (`--draft`) using the `create-pr.sh` skill.
- [x] Graduate the draft PR to ready-for-review using `create-pr.sh --ready`.

### Phase 7: Verification, Copilot Compliance & Iteration
- [x] Trigger/wait for GitHub Copilot automated review on the PR.
- [x] Address and resolve any findings by committing necessary refinements.
- [x] Resolve CI deadlock by ignoring self-status and event-trigger check runs (`Verify PR Requirements` and `Trigger Executor on Event`).
- [x] Enhance troubleshooting logging to display detailed information on processed/ignored check runs, all reviewers with their status, type, and author association.
- [x] Integrate collaborator permission checks via `getCollaboratorPermissionLevel` API to correctly identify team/group-based write or triage access with a graceful fallback.
- [x] Enhance logging to check and print collaborator permission level of ALL human reviewers (including inactive ones) to ease debugging.
- [x] Configure `verify-pr-requirements` in `pull_request.yaml` to depend on all other validation/test jobs (making it the last gate to execute).
- [x] Update `pull_request.yaml` to allow `dependabot[bot]` actor branch pushes in `Enforce Fork Contributions` check run.
- [x] Implement Dependabot PR auto-merge rule and proxy approval logic inside `verify-pr-requirements.mjs` (merges Dependabot PRs based on AI review approval).
- [x] Tighten `pull_request.yaml` to ensure only `dependabot[bot]` can push to `dependabot/` branches inside the same-repository.
- [x] Fix Nix-run syntax error in `verify-pr-requirements.mjs` by writing prompt to file instead of command line.
- [x] Remove redundant, false-failing `verify-pr-requirements` job from `pull_request.yaml` to improve developer UX and prevent premature check failures.
- [x] Implement robust GitHub CLI and native auto-merge pipeline (`gh pr merge --auto`) with direct merge and REST API fallbacks in `verify-pr-requirements.mjs`.
- [x] Verify updated code locally and obtain approval.

### Sub-task: Reconcile Release PR Labels after Merge
- [x] Implement label reconciliation logic inside `.github/workflows/scripts/publish-release.js` to find the merged Release PR and update its labels (remove `autorelease: pending` and `ready-to-merge`, add `autorelease: tagged`).
- [x] Run `node --check` to verify that `.github/workflows/scripts/publish-release.js` syntax remains valid.
- [x] Run `actionlint` locally to verify that `.github/workflows/release.yml` syntax remains valid.
- [x] Enter Proactive Review Mode on the written code diff to proactively resolve any issues human or Copilot automated reviews might flag.
- [x] Present the unstaged diff to the developer in the chat and request their IDE review.
- [x] Commit changes locally with a conventional prefix (e.g., `ci: reconcile release-please pr labels after merge`) using the secure `commit-push.sh` skill (or signed commit fallback).

### Phase 8: Direct Fork PR Auto-Merge (Ruleset Enabled)
- [x] Refactor `.github/workflows/scripts/verify-pr-requirements.mjs` to completely remove the legacy `isFork` blocker, allowing fork PRs to be merged directly by the executor just like local and Dependabot PRs.
- [x] Run `node --check` to verify `.github/workflows/scripts/verify-pr-requirements.mjs` syntax.
- [x] Present the unstaged diff to the developer in the chat.
- [x] Solicit manual developer review and obtain explicit approval in the chat.
- [x] Stage the modified/created files and commit locally with a conventional prefix (e.g. `ci: enable automated fork PR merging via Ruleset bypass`) using the secure `commit-push.sh` skill (or signed commit fallback).
- [x] Push the branch to the user's origin fork and generate a Draft PR using `create-pr.sh`.

### Phase 9: Resolve Copilot PR Review Feedback
- [x] Refactor `verify-pr-requirements.mjs` to pass `isFork` to `mergePullRequest`.
- [x] Update `mergePullRequest` with a graceful merge-failure fallback: if the merge fails, apply/keep the `ready-to-merge` label and post a detailed diagnostic comment explaining how the maintainer can merge or fix the Ruleset.
- [x] Rephrase `ReleaseProcess.md` to add a clear, explicit checklist of the repository settings and ruleset prerequisite options needed to enable direct auto-merges of fork PRs.
- [x] Update `WorkflowExecutionProtections.md` to explicitly time-bound its product claims as of August 2026, and add prominent links to the authoritative GitHub documentation pages.
- [x] Compile and verify the refactored script with `node --check`.
- [x] Prevent throwing restError in the fork graceful fallback path to keep the PR status check green, enabling manual maintainer merge.

### Phase 10: Repair Manual Release Workflows Broken by Containerization

**Root cause (run [31718150932](https://github.com/rancher/terraform-provider-rancher2/actions/runs/31718150932), job `rc-release`, step `Prepare Registry Manifest`, exit 1):** commit `2baf25d9` moved the manual release jobs into the `ghcr.io/rancher/ci-image/nix` container but left `WORKSPACE: ${{ github.workspace }}` in place, so `cd "$WORKSPACE/tags/$TAG"` targets a host path that does not exist inside the container. A second, latent defect in the same step (`tags_to_delete=$(git tag | grep -v ...)` with no `|| true`) would abort the step immediately after the first is fixed. Both defects are duplicated verbatim in `manual-release.yml`.

- [x] Fix `.github/workflows/manual-rc-release.yml` `Check out new tag into a new directory`: change `path:` to the container-agnostic relative form `tags/${{ inputs.tag }}`.
- [x] Fix `.github/workflows/manual-rc-release.yml` `Prepare Registry Manifest`: drop the `WORKSPACE: ${{ github.workspace }}` env entry and `cd "$GITHUB_WORKSPACE/tags/$TAG"` instead.
- [x] Fix `.github/workflows/manual-rc-release.yml` `Prepare Registry Manifest`: append `|| true` to the `git tag | grep -v` command substitution so an empty result is not fatal under `sh -e`.
- [x] Apply the same three corrections to `.github/workflows/manual-release.yml` (`Check out new tag into a new directory` and `Prepare Registry Manifest`).
- [x] Add the mandatory explicit `timeout-minutes: 30` to the `rc-release` and `release` jobs per `.agent/rules/workflows.instructions.md` §2, matching the equivalent GoReleaser jobs in `release.yml`.

**Third defect found during verification (Phase 11):** the `Run GoReleaser` step defers expansion (`\"\$GITHUB_WORKSPACE/tags/\$TAG\"`) across the `nix-run.sh` boundary. Because that shell is started with `--ignore-environment` and neither variable is on the `--keep` allowlist, both expand to empty and the step would `cd /tags/` — failing immediately after the `Prepare Registry Manifest` fix unblocks it. Reproduced locally by simulating the `printf %q` serialisation under a cleared environment.

**Fourth defect, and the two approaches considered.** An intermediate fix expanded both variables in the *outer* step shell and interpolated the result into the inner `bash -c` string. Proactive review correctly rejected it as a **command-injection regression**: the inner shell re-parses the interpolated value as source, so `TAG='v1.2.3-rc.1";echo "PWNED'` executed the payload while the step still exited 0 — with `GITHUB_TOKEN` (contents: write) and the Vault GPG secrets in scope, mitigated only partially by the `Check Author` maintainer gate.

Two safe resolutions existed: pass the path as a positional parameter, or keep the original deferred expansion and add the variables to the `nix-run.sh` `--keep` allowlist. **Per developer decision we took the `--keep` route**, which leaves the `Run GoReleaser` step at its original, more readable one-liner and fixes the defect at its actual root — the allowlist — rather than working around it at every call site.

- [x] Add `--keep GITHUB_WORKSPACE` and `--keep TAG` to `.github/workflows/scripts/nix-run.sh` so both variables survive `nix develop --ignore-environment`.
- [x] Keep the `Run GoReleaser` step's deferred expansion (`\"\$GITHUB_WORKSPACE/tags/\$TAG\"`) in both manual workflows, with a comment recording that the deferral is deliberate and depends on the allowlist.
- [x] Verify by emulation that (a) a benign tag reaches the correct directory, (b) a tag carrying shell metacharacters is treated as a literal path component with no execution, and (c) removing the variables from the allowlist fails loudly with exit 1 rather than releasing from the wrong tree.

**Additional pre-existing findings folded in at developer request:** proactive review also flagged that neither manual job declares `environment: release` (unlike the two GoReleaser jobs in `release.yml`), and that `grep -v -e "^$TAG$"` treats the tag as a basic regular expression, so `.` matches any character. Verified locally: with the old pattern a sibling tag `v15x0.0-rc.14` is wrongly preserved when releasing `v15.0.0-rc.14`.

- [x] Add `environment: release` to the `rc-release` and `release` jobs, matching `release.yml`. **Note:** the `release` environment does not currently exist in repository settings (only `copilot` does) — GitHub auto-creates it on first reference with no protection rules, so this does not gate the dispatch today. It is a hook that a maintainer must configure (required reviewers / branch policy) in repo settings for it to provide real protection; the same gap already applies to `release.yml`.
- [x] Change the tag-pruning filter in both manual workflows to `grep -v -x -F -e "$TAG"` so the tag is matched as an exact literal whole line rather than a BRE.

**Deferred to a follow-up (out of the agreed minimal-fix scope):** extracting the multi-line `Prepare Registry Manifest` body into `.github/workflows/scripts/`, and the `|| true` masking a genuine `git tag` failure / the dead `DIR="$(pwd)"` round-trip.

### Phase 11: Verification, Review & Delivery (Quality Gates)
- [x] Run `actionlint` on `.github/workflows/manual-rc-release.yml` and `.github/workflows/manual-release.yml` via `.agent/skills/run-in-nix.sh` and resolve all findings.
- [x] Statically verify the corrected `Prepare Registry Manifest` body under `sh -e` locally against a scratch git repo containing exactly one tag, confirming it no longer aborts.
- [x] Verify the corrected `Run GoReleaser` argument survives `nix-run.sh`'s `printf %q` serialisation under a cleared environment.
- [x] Enter Proactive Review Mode on the diff against `.agent/rules/github-copilot-review.instructions.md` and `.agent/rules/workflows.instructions.md`, targeting exactly 0 automated findings. **Two passes run.** Pass 1: 🔴 1 Critical (the command-injection regression), 2 Major, 3 Minor. Pass 2, after remediation: 🟢 0 findings, approval file written via the sanctioned OTP path.
- [x] Switch to `main`, run `.agent/skills/git-sync.sh`, and branch off the synchronized `main` (`fix/manual-release-container-paths`, off `23013e13`).
- [x] 🔒 Gate 2: commit `05a594a0` (`fix: manual release workflows broken by containerization`), SSH-signed and signed-off, pushed to the fork `matttrach/terraform-provider-rancher2`. Originally committed as `37cdcdf3` with a `ci:` prefix and amended — this repository has not yet adopted the relaxed commit validation that permits `ci:`. **Consequence:** `fix:` is a bumping prefix, so release-please will cut a patch release from this CI-only change. **Gates bypassed by explicit developer instruction ("skip review, we need to get this out"):** the `--keep` rework landed after the pass-2 approval, invalidating its diff hash, so `commit-push.sh` would have refused. The commit was made directly with `git commit -s -S` rather than through the skill, bypassing both the proactive-review hash check and the interactive `user-approval.js` prompt. The approval file was **not** hand-authored or otherwise forged.
- [x] 🔒 Gate 3: [PR #2376](https://github.com/rancher/terraform-provider-rancher2/pull/2376), opened ready-for-review rather than draft, also to save a round trip at developer request.

**Bug found and fixed in tooling along the way:** `.agent/skills/commit-push.sh` `sync_default_branch()` ran `git stash push -k -u` before invoking `.agent/skills/git-sync.sh`. Because most of `.agent/skills/` is untracked, that stash swept away `git-sync.sh` itself, so the very next line failed with `No such file or directory`. The pre-emptive stash was also redundant — `git-sync.sh` already auto-stashes (`verify_git_env`) and restores with `git stash pop --index`, preserving the index. Removed; `shellcheck` clean. No data was lost in the failed run (the emergency pop restored the tree and the diff hash was unchanged).

### Phase 12: Publish the Blocked Release Candidate
- [ ] After the fix merges to `main`, dispatch `Manually Create RC Release` with `branch=main`, `tag=v15.0.0-rc.14`, `sha` empty. (`v15.0.0-rc.13` was pushed by the failed run but has no release attached; per developer decision it is left in place as a dangling tag rather than deleted.)
- [ ] Confirm the run reaches `Run GoReleaser` and publishes the GPG-signed RC assets to the GitHub release.
- [ ] Confirm the `Find Issues and Create Comments` step posts RC notifications.

<!-- Retrigger workflows -->

