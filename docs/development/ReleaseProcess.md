# Release Process & Automation

This topic overview details the standard repository release process, tracing how codebase modifications on the `main` branch are systematically packaged, GPG-signed, and published to production.

---

## Abstract

The release process establishes a highly automated, branch-based deployment pipeline that translates Conventional Commits into formal, GPG-signed product releases. By utilizing strict manifest-driven versioning on release branches, we ensure that every release is secure, auditable, and deterministic.

---

## 🧭 How Our Components Work Together

Our release process is designed around two core architectural components that work in tandem to orchestrate the software delivery lifecycle:

### 1. Branch-Based Release Strategy

We utilize a branch-based development release strategy. Code is primarily developed and merged into the `main` branch, but all releases (both RCs and stable versions) are generated directly from their respective `release/v<major>` branches.

- More details on this strategy, GPG-signing configuration, and release candidate lifecycle can be found in **[Branch-Based Releases](./ReleaseProcess/BranchReleases.md)**.

### 2. Manifest-Driven Automation

To automate versioning and changelog generation, we leverage `release-please` in manifest mode. This tool scans Conventional Commits on the `release/v<major>` branches, computes the correct SemVer increment, and maintains a running "Release PR." Once this PR is merged, the system automatically tags the release and initiates the compilation pipeline.

- More details on the action parameters, CLI usage, and manifest configurations can be found in **[Release Please](./ReleaseProcess/ReleasePlease.md)**.

### 🔄 The Combined Execution Lifecycle

When a developer's contribution lands on the `main` branch:

1. **Development & Backporting**: Code is merged to `main`. An automated workflow triggers to cherry-pick those changes into the targeted release branch (e.g., `release/v15.2.0`), creating a "backport PR".
2. **Release Candidate (RC) Generation**: When the backport PR is approved and merged into the release branch (`release/v*`), it triggers the unified RC release workflow. This computes the next RC tag, compiles binaries via GoReleaser, cryptographically signs them, and pushes the Release Candidate. Concurrently, `release-please` updates a pending "Release PR" on the release branch.
3. **Verification & Stable Release**: Once the Release PR is approved and merged by a maintainer on the release branch, the pipeline executes the full acceptance test suite. Upon success, it automatically tags the release, extracts GPG credentials, compiles stable binaries, and publishes them to the GitHub Release Registry.

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
| === PART 2: AUTOMATED BACKPORT TO RELEASE ==|                                                                |
|                                             |                                                                |
|  5. Squash Merge Lands on main ───────────> | ──► Trigger: push to main                                      |
|                                             |       ──► Automated Backport Action runs:                      |
|                                             |           - Cherry-picks commits to target release/v* branch   |
|                                             |           - Action: Creates a "backport PR" against release/v* |
|                                             |                                                                |
|  6. Backport PR Merged to Release Branch ─> | ──► Trigger: push to release/v*                                |
|     (Maintainer approves/merges Backport)   |       ──► Release Candidate (RC) Step:                         |
|                                             |           - Calculates next RC tag (e.g., v1.2.3-rc.0)         |
|                                             |           - Vault: Securely extracts GPG keys                  |
|                                             |           - Nix: GoReleaser compiles & signs RC binaries       |
|                                             |           - Action: Publishes GPG-Signed RC Release            |
|                                             |                                                                |
|                                             |       ──► Release Please Action runs:                          |
|                                             |           - Scans Conventional Commits on release branch       |
|                                             |           - Action: Updates/creates draft "Release PR"         |
|                                             |                                                                |
| === PART 3: STABLE RELEASE GENERATION ======|                                                                |
|                                             |                                                                |
|  7. Maintainer Merges Release PR ─────────> | ──► Trigger: push to release/v* (Release PR merged)            |
|     (Maintainer approves/merges Release PR) |       ──► Release Please Action runs:                          |
|                                             |           - Detects Release PR merge                           |
|                                             |           - Action: Outputs: release_created = true            |
|                                             |         │                                                      |
|                                             |         ▼                                                      |
|                                             |       ──► Full Release Step:                                   |
|                                             |           - Action: Automatically tags version (v1.2.3)        |
|                                             |           - AWS: Assumes OIDC IAM Role                         |
|                                             |           - Nix: Executes FULL ACCEPTANCE TEST SUITE           |
|                                             |           - Vault: Securely extracts GPG credentials           |
|                                             |           - Keyring Workaround: Dynamically parses primary ID  |
|                                             |           - Nix: GoReleaser cross-compiles stable binaries     |
|                                             |           - Action: Publishes Final GPG-Signed Release         |
|                                             |                                                                |
```
