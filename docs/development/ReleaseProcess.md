# Release Process & Automation

This topic overview details the standard repository release process, tracing how codebase modifications on the `main` branch are systematically packaged, GPG-signed, and published to production.

---

## Abstract

The release process establishes a highly automated, trunk-based deployment pipeline that translates Conventional Commits into formal, GPG-signed product releases. By eliminating manual tagging and utilizing strict manifest-driven versioning, we ensure that every release is secure, auditable, and deterministic.

---

## 🧭 How Our Components Work Together

Our release process is designed around two core architectural components that work in tandem to orchestrate the software delivery lifecycle:

### 1. Trunk-Based Release Strategy

We release all product versions directly from our single source of truth—the `main` branch. This eliminates the complexity and drift of maintaining parallel release branches.

- More details on this strategy, GPG-signing configuration, and release candidate lifecycle can be found in **[Release From Main](./ReleaseProcess/MainBranchReleases.md)**.

### 2. Manifest-Driven Automation

To automate versioning and changelog generation, we leverage `release-please` in manifest mode. This tool scans Conventional Commit squash-merge titles on the `main` branch, computes the correct SemVer increment, and maintains a running "Release PR." Once this PR is merged, the system automatically tags the release and initiates the compilation pipeline.

- More details on the action parameters, CLI usage, and manifest configurations can be found in **[Release Please](./ReleaseProcess/ReleasePlease.md)**.

### 🔄 The Combined Execution Lifecycle

When a developer's contribution lands on the `main` branch:

1. **Trigger & Version Calculation**: The `release-please` action is invoked. It scans the squash-merge commit title and updates the running Release PR (or opens a new one if none exists), generating an automated changelog.
2. **Integration & Acceptance Testing**: To guarantee stability, any push to a Release Please branch triggers our comprehensive OpenID Connect (OIDC) acceptance tests inside a Nix shell, deploying real AWS infrastructure to verify binary correctness.
3. **Verification & Signing**: Upon maintainer merge of the Release PR, the pipeline securely extracts our GPG signing credentials, compiles the binaries via GoReleaser inside a hardened container, cryptographically signs the assets, and publishes them natively to the GitHub Release Registry.

---

## 📊 Actor-vs-Automation Interaction Swimlanes

This swimlane diagram traces the detailed event triggers and data flow between development roles and automated GHA runners:

```text
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
