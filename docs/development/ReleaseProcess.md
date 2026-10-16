# Release Process & Automation

This topic overview details the standard repository release process, tracing how codebase modifications on the `main` branch are systematically packaged, GPG-signed, and published to production.

---

## Abstract

The release process establishes a highly automated, branch-based deployment pipeline that translates Conventional Commits into formal, GPG-signed product releases. By utilizing strict manifest-driven versioning on release branches, we ensure that every release is secure, auditable, and deterministic.

Acceptance tests execute against qualifying automated Pull Requests prior to merging onto release branches. When backports land on release branches, automated workflows generate and publish signed Release Candidates (RCs) and notify tracking issues. When maintainers merge the pending `release-please` PR on a release branch, the pipeline cryptographically compiles, signs, and publishes the official release.

---

## 🧭 How Our Components Work Together

Our release process is designed around two core architectural components that work in tandem to orchestrate the software delivery lifecycle:

### 1. Branch-Based Release Strategy

We utilize a branch-based development release strategy. Code is primarily developed and merged into the `main` branch, but all releases (both RCs and stable versions) are generated directly from their respective `release/v<major>` branches (e.g., `release/v15`, `release/v14`, `release/v13`).

- More details on this strategy, GPG-signing configuration, and release candidate lifecycle can be found in **[Branch-Based Releases](./ReleaseProcess/BranchReleases.md)**.
- Details on the end-to-end pipeline can be found in **[End-to-End Value Stream Lifecycle](./GitHubWorkflows/Overview.md)**.

### 2. Manifest-Driven Automation

To automate versioning and changelog generation, we leverage Google's `release-please` running in manifest mode (`release-please-config.json` and `.release-please-manifest.json`). This tool scans Conventional Commits on the `release/v<major>` branches, computes the correct SemVer increment, and maintains an active draft "Release PR." Once this PR is merged, the system automatically tags the release and initiates the GoReleaser compilation and signing pipeline.

- More details on the action parameters, CLI usage, and manifest configurations can be found in **[Release Please](./ReleaseProcess/ReleasePlease.md)**.

---

## 📊 Actor-vs-Automation Interaction Swimlanes

This swimlane diagram traces the actual event triggers, security boundaries, and data flows between development roles and automated GitHub Actions runners:

```text
|                [ Actor ]                    |                  [ Automation (GHA / Nix / Vault) ]            |
├─────────────────────────────────────────────┼────────────────────────────────────────────────────────────────┤
|                                             |                                                                |
| === PART 1: PULL REQUEST TO SQUASH-MERGE ===|                                                                |
|                                             |                                                                |
|  1. PR Opened or Code Updated ------------> | ──► Trigger: pull_request (main)                               |
|     (Contributor submits change)            |    ├─► validate-pr.yaml: Checks "- Addresses: #issue" format   |
|                                             |    └─► pull_request.yaml (Nix Container):                      |
|                                             |         - Cross-platform builds (make build)                   |
|                                             |         - Go unit tests (make test)                            |
|                                             |         - Linters: tflint, actionlint, eslint, shellcheck      |
|                                             |         - Conventional Commits, cspell, gitleaks, schemas      |
|                                             |                                                                |
|  2. Maintainer Reviews & Approves --------> | ──► Collaborator review & approval                             |
|     (TERRAFORM_MAINTAINERS)                 |                                                                |
|                                             |                                                                |
|  3. Maintainer Squash-Merges PR ----------> | ──► Merges PR into main (Single atomic commit on trunk)        |
|                                             |                                                                |
| === PART 2: AUTOMATED BACKPORT & PRE-MERGE =|                                                                |
|                                             |                                                                |
|  4. Merge Commit Lands on main -----------> | ──► Trigger: push to main (backport-prs.yml)                   |
|                                             |    ├─► wait-for-settle: Waits for commit indexing              |
|                                             |    ├─► Vault: Retrieves GITHUB_MERGE_TOKEN                     |
|                                             |    └─► backports.js: Discovers tracking sub-issues             |
|                                             |         - Cherry-picks commit: git cherry-pick -X theirs       |
|                                             |         - Opens Backport PR targeting release/v*               |
|                                             |         - Enables native GitHub squash auto-merge              |
|                                             |                                                                |
|  5. Backport PR Triggered on release/v* --> | ──► Trigger: e2e-pr-listener.yml (Lightweight event emitter)   |
|                                             |         │                                                      |
|                                             |         ▼                                                      |
|                                             | ──► Trigger: e2e-pr.yml (workflow_run from main)               |
|                                             |    ├─► Filter: Validates bot author & release branch           |
|                                             |    ├─► Distributed Lock & IP collision prevention              |
|                                             |    ├─► AWS OIDC: Assumes test IAM role                         |
|                                             |    ├─► Nix: Runs live AWS acceptance tests (run_tests.sh)      |
|                                             |    ├─► Cleanup: Guarantees teardown of AWS cloud resources    |
|                                             |    └─► Status Gate: Submits commit status check (E2E Tests / executor) |
|                                             |                                                                |
| === PART 3: RELEASE CANDIDATE (RC) STAGING =|                                                                |
|                                             |                                                                |
|  6. Backport PR Merged to release/v* -----> | ──► Trigger: push to release/v* (release.yml)                   |
|     (Auto-merges upon passing CI & E2E)     |    ├─► release-please-release: Skipped (not a release PR)      |
|                                             |    └─► release-please-pr: Updates draft Release PR via Vault   |
|                                             |         │                                                      |
|                                             |         ▼                                                      |
|                                             | ──► Trigger: release-branch-merge.yml (workflow_run)           |
|                                             |    ├─► merge-label: Marks backport issue internal/merged       |
|                                             |    └─► trigger-rc-release: Computes next RC tag                |
|                                             |         │                                                      |
|                                             |         ▼                                                      |
|                                             | ──► Trigger: rc-release.yml (Dispatched with tag & branch)     |
|                                             |    ├─► Tags commit and pushes to origin                        |
|                                             |    ├─► Vault: Retrieves GPG signing keys                       |
|                                             |    ├─► Nix: GoReleaser builds & signs (.goreleaser_rc.yml)     |
|                                             |    └─► rc-notify: Comments on issues with validation link      |
|                                             |                                                                |
| === PART 4: FULL PRODUCTION RELEASE ========|                                                                |
|                                             |                                                                |
|  7. Maintainer Merges Release PR ---------> | ──► Trigger: push to release/v* (Release PR merged)            |
|     (After RC validation by QA team)        |    ├─► release-please-release: Outputs release_created = true  |
|                                             |    ├─► Full Release Step:                                      |
|                                             |    │    - Clears runner disk space                             |
|                                             |    │    - Vault: Retrieves GPG signing credentials             |
|                                             |    │    - Nix: GoReleaser compiles & signs (.goreleaser.yml)   |
|                                             |    │    - publish-release: Promotes draft to published         |
|                                             |    └─► release-please-pr: Creates next release PR cycle        |
|                                             |                                                                |
```

---

## 🛠️ Modular Architectural Blueprints Map

- **[End-to-End Value Stream Lifecycle](./GitHubWorkflows/Overview.md)**: Detailed specification of the complete 7-stage delivery value stream.
- **[Pull Request & Merge Lifecycle on Main](./GitHubWorkflows/PullRequestLifecycle.md)**: PR validation, Nix CI checks, and maintainer squash-merge.
- **[Backport Orchestration & Issue Tracking](./GitHubWorkflows/BackportOrchestration.md)**: Tracking issue adoption, GitHub Sub-Issues hierarchy, and automated cherry-picking.
- **[Privileged AWS End-to-End Testing](./GitHubWorkflows/EndToEndTesting.md)**: Listener-executor security model, AWS OIDC, matrix locking, and status gate.
- **[Branch-Based Releases](./ReleaseProcess/BranchReleases.md)**: Deep dive into release branches, RC calculation, and GoReleaser signing.
- **[Release Please Automation](./ReleaseProcess/ReleasePlease.md)**: Manifest configuration schema, conventional commit mappings, and CLI usage.
