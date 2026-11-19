# GitHub Workflows & Automation Architecture

This document provides a comprehensive architectural guide to the GitHub Actions workflows, automation scripts, operational actor journeys, and design decisions in `terraform-provider-rancher2`.

---

## Table of Contents

- [1. High-Level Overview & Value Stream Maps](#1-high-level-overview--value-stream-maps)
  - [1.1 Core Principles](#11-core-principles)
  - [1.2 Observer & Project Manager (PjM) Value Stream Map](#12-observer--project-manager-pjm-value-stream-map)
  - [1.3 Workflow Maintainer Technical Value Stream Map](#13-workflow-maintainer-technical-value-stream-map)
  - [1.4 Scheduled & Auxiliary Maintenance Workflows](#14-scheduled--auxiliary-maintenance-workflows)
- [2. Actor Perspectives & Operational Journeys](#2-actor-perspectives--operational-journeys)
  - [2.1 Issue Creator's Perspective](#21-issue-creators-perspective)
  - [2.2 Project Manager's Perspective (PjM)](#22-project-managers-perspective-pjm)
  - [2.3 PR Author's Perspective (Developer)](#23-pr-authors-perspective-developer)
  - [2.4 QA Engineer's Perspective](#24-qa-engineers-perspective)
  - [2.5 Release Operator's Perspective](#25-release-operators-perspective)
- [3. Workflow Reference Guide](#3-workflow-reference-guide)
  - [3.1 Pull Request Validation & Testing](#31-pull-request-validation--testing)
  - [3.2 Merge & Queue Automation](#32-merge--queue-automation)
  - [3.3 Issue Tracking & Backport Orchestration](#33-issue-tracking--backport-orchestration)
  - [3.4 Release Branch Lifecycle & Release Candidate (RC) Automation](#34-release-branch-lifecycle--release-candidate-rc-automation)
  - [3.5 Branch Setup & Synchronization](#35-branch-setup--synchronization)
  - [3.6 Compliance & Infrastructure Maintenance](#36-compliance--infrastructure-maintenance)
- [4. Workflow Scripts Reference Guide](#4-workflow-scripts-reference-guide)
- [5. Rancher2 Terraform Project Lifecycle & State Machine](#5-rancher2-terraform-project-lifecycle--state-machine)
- [6. Architectural Decisions & Engineering Standards](#6-architectural-decisions--engineering-standards)
  - [6.1 Architectural Decisions](#61-architectural-decisions)
  - [6.2 Security & Authentication Model](#62-security--authentication-model)
  - [6.3 Script Execution & Local Safety](#63-script-execution--local-safety)

---

## 1. High-Level Overview & Value Stream Maps

### 1.1 Core Principles

The CI/CD pipeline in this repository is an event-driven automation system built on four architectural pillars:

1. **Issue-First Value Flow**: Every code modification targeting `main` must link to an open GitHub issue in the repository. Pull requests without an associated issue fail automated validation and cannot merge.
2. **Trunk-Based with Automated Backports**: Developers develop and merge fixes strictly against `main`. Once merged, automation generates cherry-pick pull requests against all targeted maintenance branches (`release/v*`).
3. **Automated Release Candidate (RC) Delivery**: Merging an automated backport into a release branch triggers automated Git tagging (e.g. `v15.1.0-rc.1`) and builds an RC release with cross-platform provider binaries.
4. **Project Board Lifecycle Automation**: Tracking issues and sub-issues advance across columns on Rancher2 Terraform Project (`Ready for QA` -> `Ready For Release` -> `Done`) through automated event handlers.
5. **Modular JavaScript Engines with Full Unit Test Coverage**: Workflow YAML definitions remain thin declarative orchestrators. All non-trivial business and automation logic is encapsulated in testable Node.js modules in `.github/workflows/scripts/`, validated by an extensive unit test harness.

---

### 1.2 Observer & Project Manager (PjM) Value Stream Map

This map outlines what happens across each phase of value delivery without workflow implementation details. It allows Project Managers and cross-functional observers to track issues, pull requests, and releases across their lifecycle:

```
[User Discovers Defect]
          │
          ▼
   1. Issue Ingestion & Triage
      - User submits issue describing problem
      - Repository maintainer assigns issue to an engineering team
      - Issue is tracked on Rancher2 Terraform Project Board
          │
          ▼
   2. Engineering Assignment & Fix
      - Team Project Manager (PjM) adds issue to milestone
      - PjM assigns issue to a team developer
      - Developer develops fix on personal fork
      - Developer opens PR targeting main linking the issue ("- Addresses: #...")
      - Developer adds release labels to issue as necessary
      - CoPilot Review is automatically triggered
          │
          ▼
   3. Teammate Peer Review & Main Merge
      - Issue link to PR is automatically validated
      - Continuous integration runs static analysis and unit test suite
      - Developer responds to and resolves all review comment threads (including CoPilot comments)
      - Teammate reviews and approves the PR (collaborator triage role)
      - PR is automatically squash-merged into main
          │
          ▼
   4. Automated Backport Distribution
      - Automation identifies targeted Rancher release branches from sub-issues on issue
      - Automation cherry-picks the fix and opens backport PRs against each release branch
      - Release branch CI verifies the cherry-picked changes, including e2e tests
      - Backport PRs merge into their respective release branches
          │
          ▼
   5. Release Candidate (RC) Staging & QA Notification
      - Merging into release branch automatically tags next Release Candidate (e.g. v15.2.0-rc.1)
      - RC release is created as "pre-release" with full binary build artifacts
      - Backport issue status transitions to "Ready for QA"
          │
          ▼
   6. QA Verification
      - QA engineer picks up backport issue from "Ready for QA"
      - QA validates fix against the published Release Candidate binary
      - When tests pass, QA posts results and closes the backport issue
          │
          ▼
   7. Release Readiness & Milestone Gate
      - When all backport issues for an issue are closed, parent tracking issue is automatically closed if it has a milestone
      - Issue status transitions to "Ready For Release"
          │
          ▼
   8. Production Release
      - Release manager audits "Ready For Release" items and verifies milestones
      - Release manager merges the staged Release PR on the release branch
      - Production release is compiled, signed, and published to Terraform Registry
      - Issue status on Rancher2 Terraform Project board automatically transitions to "Done"
```

---

### 1.3 Workflow Maintainer Technical Value Stream Map

This map provides the detailed execution flow for CI/CD workflow maintainers, showing exact workflow file names, triggers, script delegation, artifact signaling, and state mutations:

```
[User Discovers Defect]
          │
          ▼
   Issue Created (e.g. #2500)
          │
          ▼
   Maintainer assigns team label and project issue "team" field
          │
          ├──► Developer/Maintainer adds release/v* label(s) to #2500
          │         │
          │         ▼
          │    [backport-issues.yml] (on issues: labeled)
          │         │
          │         ├──► calls backports.js (mode: backport-issues)
          │         └──► creates unassigned backport sub-issue [v<major>]
          │
          ▼
   Developer opens PR targeting main ("- Addresses: #2500")
          │
          ├──► [validate-pr.yaml]
          │         │
          │         ├──► calls validate-pr-description.js
          │         └──► emits artifact: signal/pr-signal.json
          │                   │
          │                   ▼
          │              [tracking-issue.yml] (on workflow_run completed)
          │                   │
          │                   ├──► calls releases.js (mode: tracking-issue)
          │                   ├──► labels #2500 with internal/tracking
          │                   └──► dispatches [backport-issues.yml] (for any newly discovered release/v* labels)
          │                             │
          │                             └──► calls backports.js (mode: backport-issues)
          │                                  creates sub-issue for each release/v* label
          │
          ├──► [pull_request.yaml] (11-job hermetic Nix container matrix)
          │
          ├──► [fossa.yml] (license compliance and security vulnerability scan)
          │
          ├──► [dependabot-approve.yml] (submits automated approval for qualifying Dependabot PRs)
          │
          └──► [auto-merge-listener.yml] (on pull_request_review: submitted)
                    │
                    └──► emits artifact: approval-signal.json
                              │
                              ▼
                         [auto-merge.yml] (on workflow_run completed)
                              │
                              ├──► calls auto-merge.js (mode: execute)
                              ├──► verifies reviewer has triage, write, or admin role
                              ├──► asserts no active CHANGES_REQUESTED
                              ├──► asserts no changes to .github/**
                              └──► executes GraphQL mutation: enablePullRequestAutoMerge
                                        │
                                        ▼
                                [Squash-Merge into main]
                                        │
                                        ▼
                                 [backport-prs.yml] (on push to main)
                                        │
                                        ├──► calls backports.js (mode: backport-pr)
                                        └──► opens backport PRs targeting release/v*
                                                  │
                         ┌────────────────────────┴────────────────────────┐
                         ▼                                                 ▼
             [Backport PR -> release/v15]                      [Backport PR -> release/v16]
                         │                                                 │
                         ▼                                                 ▼
             [pull_request.yaml] (Release CI)                  [pull_request.yaml] (Release CI)
                         │                                                 │
                         ▼                                                 ▼
             [Merge into release/v15]                          [Merge into release/v16]
                         │                                                 │
                         ├──► [release-branch-merge.yml]                   ├──► [release-branch-merge.yml]
                         │         │                                       │         │
                         │         ├──► calls backports.js (merge-label)   │         ├──► calls backports.js (merge-label)
                         │         └──► Project 73: "Ready for QA"         │         └──► Project 73: "Ready for QA"
                         │                                                 │
                         └──► [release-please.yml] (on push)               └──► [release-please.yml] (on push)
                                   │                                                 │
                                   ├──► release-please-pr (PR Only)                  ├──► release-please-pr (PR Only)
                                   │         │                                       │         │
                                   │         ├──► updates Staging Release PR         │         ├──► updates Staging Release PR
                                   │         │                                       │         │
                                   │         └──► [sync-manifest-pr] (Event-Driven)  │         └──► [sync-manifest-pr] (Event-Driven)
                                   │                   │                             │                   │
                                   │                   └──► calls sync-manifest.js   │                   └──► calls sync-manifest.js
                                   │                                                 │
                                   └──► release-please-rc (if no full release)       └──► release-please-rc (if no full release)
                                             │                                                 │
                                             ├──► calls releases.js (prepare-rc)               ├──► calls releases.js (prepare-rc)
                                             ├──► tags v15.1.0-rc.1                            ├──► tags v16.0.0-rc.1
                                             └──► dispatches [release.yml]                     └──► dispatches [release.yml]
                                                       │                                                 │
                                                       └──► calls goreleaser.sh                          └──► calls goreleaser.sh
                                                                                               │
                                                                                               ▼
                                                                                 [QA tests RC & closes sub-issues]
                                                                                               │
                                                                                               ▼
                                                                                  [close-tracking-issue.yml] (on issue closed)
                                                                                               │
                                                                                               ├──► calls backports.js (close-tracking-issue)
                                                                                               ├──► validates all sub-issues closed
                                                                                               ├──► verifies/inherits milestone
                                                                                               ├──► closes tracking issue #2500
                                                                                               └──► sets issue status on project 73: "Ready For Release"
                                                                                               │
                                                                                               ▼
                                                                                 [Maintainer Merges Release PR]
                                                                                               │
                                                                                               ▼
                                                                                  [release-please.yml] (on push)
                                                                                               │
                                                                                               ├──► release_created == 'true'
                                                                                               └──► dispatches [release.yml]
                                                                                                         │
                                                                                                         ├──► GoReleaser builds binaries
                                                                                                         ├──► GPG signs assets
                                                                                                         └──► publishes GitHub release
                                                                                                                   │
                                                                                                                   ▼
                                                                                                      [release-complete.yml] (on workflow_run)
                                                                                                                   │
                                                                                                                   ├──► calls backports.js (release-done)
                                                                                                                   ├──► verifies non-draft, non-RC tag
                                                                                                                   └──► Updates "Ready For Release" Issues on Project 73 to "Done"
                                                                                                                            │
                                                                                                                            ▼
                                                                                                               [Infrastructure Teardown]
                                                                                                                            │
                                                                                                                            ▼
                                                                                                                 [cleanup-e2e.yml] (Schedule: Daily 13:00 UTC)
                                                                                                                            │
                                                                                                                            ├──► calls e2e.js (find-completed-e2e-runs)
                                                                                                                            ├──► identifies completed e2e test runs in window
                                                                                                                            └──► executes parallel matrix cleanups via cleanup.sh
```

---

### 1.4 Scheduled & Auxiliary Maintenance Workflows

In addition to the primary event-driven PR-to-Release pipeline, the following workflows run on schedules or manual triggers to maintain repository health, infrastructure hygiene, and release branch state:

```
[Daily Schedule: 13:00 UTC / Dispatch]
       │
       ▼
[cleanup-e2e.yml]
       │
       ├──► Job: find-completed-runs ──► calls e2e.js (find-completed-e2e-runs)
       │                                  identifies all completed e2e runs within 24h window
       │
       └──► Job: cleanup-matrix (parallel)
                 │
                 ├──► assumes AWS CI role (arn:aws:iam::270074865685:role/terraform-module-ci-test)
                 └──► runs cleanup.sh in hermetic Nix container to delete orphaned AWS resources

[Manual Dispatch]
       │
       ▼
[sync-release-manifest.yml]
       │
       ├──► calls sync-manifest.js
       └──► manually synchronizes manifests and configs for specified branch out-of-band

[Manual Dispatch]
       │
       ▼
[setup-release-branch.yml]
       │
       ├──► calls setup-release-branch.js
       └──► initializes new release/v<N> OR idempotently synchronizes caller stubs on existing branches
```

---

## 2. Actor Perspectives & Operational Journeys

### 2.1 Issue Creator's Perspective

The Issue Creator is typically an external consumer, customer, or internal developer reporting a bug or requesting an enhancement.

1. **Submission**:
   - The creator opens an issue on GitHub describing the unexpected behavior, Terraform configuration, provider version, and Rancher version.
2. **Triage & Team Assignment**:
   - Repository maintainers review incoming issues.
   - Maintainers manually assign the issue to a designated functional engineering team (e.g., Cloud, Provisioning).
3. **Visibility**:
   - The creator observes team assignment on the issue.
   - When engineering begins work, a linked pull request appears on the issue timeline.
   - If the fix impacts multiple supported Rancher releases, backport sub-issues appear for each target release branch (e.g. `release/v14`, `release/v15`).
   - When the fix is verified and ready for release, the issue is associated with an official milestone (e.g. `v15.2.0`) and closed.

---

### 2.2 Project Manager's Perspective (PjM)

The Project Manager (PjM) oversees team execution, release alignment, and issue progression for issues routed to their functional team.

1. **Intake & Team Review**:
   - The repository maintainer routes the issue to the PjM's engineering team by assigning the team label and setting the issue `team` field.
   - The PjM reviews the issue for scope, priority, and release targets.
   - The PjM assigns the issue to a milestone (e.g. `TFR2P November 30 2026`).
2. **Developer Assignment**:
   - The PjM assigns a developer from their team to implement the fix.
3. **Lifecycle & Status Oversight**:
   - The PjM tracks the issue and its associated backport sub-issues across the [Rancher2 Terraform Project](https://github.com/orgs/rancher/projects/73):
     - `Backlog` / `In Progress`: Engineering actively working on PR targeting `main`.
     - `Ready for QA`: Backport PR has merged into the release branch, and an automated RC pre-release binary is available for testing.
     - `Ready For Release`: QA has validated and closed all backport sub-issues; parent tracking issue automatically closes with the assigned milestone.
     - `Done`: Official release published to the Terraform Registry.

---

### 2.3 PR Author's Perspective (Developer)

The PR Author is an engineering team member assigned to resolve the issue.

1. **Assignment**:
   - After the repository maintainer assigns the issue to a team, the team's Project Manager (PjM) sets the target release milestone and assigns the issue to the developer.
2. **Development & Verification**:
   - The developer branches from the latest upstream `main`:
     ```bash
     git checkout main && git pull upstream main
     git checkout -b fix-cluster-sync
     ```
   - Implements the fix and validates locally:
     ```bash
     make build
     make test
     make fmt
     make lint
     ```
3. **Opening the Pull Request**:
   - The developer opens a pull request targeting `main`.
   - **Mandatory Linkage**: The pull request body must explicitly reference the parent issue:
     ```markdown
     - Addresses: #2500
     ```
   - **Conventional Commits**: Commits must follow Conventional Commits (e.g., `fix:`, `feat:`, `test:`) with subject lines under 70 characters.
4. **Automated CI & Teammate Auto-Merge**:
   - `validate-pr.yaml` parses the issue reference and emits `pr-signal.json`.
   - `tracking-issue.yml` automatically converts `#2500` into an `internal/tracking` issue and copies release branch labels.
   - `pull_request.yaml` runs the 11-job hermetic Nix CI matrix.
   - **Teammate Peer Review**: Maintainer review is not required to merge. Any team member or collaborator with `triage`, `write`, or `admin` permissions on the repository can review and approve the PR.
   - Once a teammate approves and all CI checks pass, `auto-merge.yml` automatically validates the review, submits bot approval for branch protection, and squash-merges into `main`.
5. **Hands-Off Backports**:
   - The developer does not create cherry-pick PRs manually. `backport-prs.yml` runs immediately upon merge into `main`, cherry-picking the squash commit and opening backport pull requests against each release branch specified on the issue (e.g. `release/v15`, `release/v16`).

---

### 2.4 QA Engineer's Perspective

The QA Engineer validates bug fixes and new features on specific Rancher version release candidates.

1. **Initial State of Backport Sub-Issues**:
   - Automation creates a dedicated sub-issue for each target release branch (e.g. `[v15] Fix cluster creation timeout`).
   - The parent tracking issue remains assigned to the functional team.
   - The backport sub-issue is initially unassigned.
2. **Transition to "Ready for QA"**:
   - When the developer's backport PR merges into `release/v15`, `release-branch-merge.yml` automatically sets the backport sub-issue status on [Rancher2 Terraform Project](https://github.com/orgs/rancher/projects/73) to **`Ready for QA`**.
   - `release-please.yml` automatically detects the merge, tags the next Release Candidate (e.g. `v15.2.0-rc.1`), and triggers `release.yml` to publish compiled provider binaries.
3. **Execution**:
   - The QA engineer filters [Rancher2 Terraform Project](https://github.com/orgs/rancher/projects/73) for items in **`Ready for QA`**.
   - Picks up the unassigned backport sub-issue and assigns themselves.
   - Downloads the published RC binary for `v15.2.0-rc.1` (or points Terraform to local `bin/` via `TF_CLI_CONFIG_FILE`).
   - Executes validation test scenarios.
4. **Resolution**:
   - **Validation Passes**: QA comments on the sub-issue with test evidence and closes the sub-issue.
   - **Validation Fails**: QA comments on the sub-issue with reproduction steps and coordinates with the engineering team.
5. **Automatic Parent Progression**:
   - When QA closes all sub-issues linked to a tracking issue, `close-tracking-issue.yml` automatically closes the parent tracking issue and moves it to **`Ready For Release`** (provided an active milestone is assigned).

---

### 2.5 Release Operator's Perspective

The Release Operator is a repository maintainer or release manager cutting an official production release.

1. **Pre-Release Audit on Rancher2 Terraform Project**:
   - The operator navigates to [Rancher2 Terraform Project](https://github.com/orgs/rancher/projects/73).
   - Verifies all tracking issues scheduled for the release target (e.g. milestone `v15.2.0`) are in the **`Ready For Release`** column.
   - **Milestone Check**: Ensures all tracking issues have an assigned milestone. (The automation blocks closure and status transition if an issue lacks a milestone).
2. **Reviewing the Staging Release PR**:
   - Release Please continuously maintains an open staging PR on each `release/v*` branch (e.g. `chore(release/v15): release 15.2.0`).
   - The operator verifies:
     - Accumulated changelog entries.
     - Version bumps in `package.json` and `.release-please-manifest.json`.
     - End-to-end test results on the release branch.
3. **Triggering the Release**:
   - The operator approves and merges the staging Release PR into `release/v15`.
4. **Automated Publishing & Board Finalization**:
   - `release-please.yml` detects `release_created == 'true'` and dispatches `release.yml`.
   - `release.yml` compiles multi-architecture provider binaries via GoReleaser, computes checksums, signs with GPG, generates `terraform-registry-manifest.json`, and publishes the release to GitHub and the Terraform Registry.
   - `release-complete.yml` verifies the release tag (non-draft, non-RC) and automatically advances all issues in **`Ready For Release`** on Rancher2 Terraform Project to **`Done`**.

---

## 3. Workflow Reference Guide

Below is the technical reference for each workflow residing in `.github/workflows/`.

### 3.1 Pull Request Validation & Testing

#### `validate-pr.yaml`

- **Goal**: Enforce the issue-first policy on all PRs targeting `main` and publish a verified pre-flight signal.
- **Triggers**: `pull_request` targeting `main`, `merge_group`.
- **Permissions**: `contents: read`, `pull-requests: read`.
- **Script**: `.github/workflows/scripts/validate-pr-description.js`.
- **Key Behaviors**:
  - Scans PR description for valid `- Addresses: #<issue>` syntax.
  - Verifies referenced issue exists, is currently open, and is not a pull request.
  - Emits `signal/pr-signal.json` and uploads it as artifact `validate-pr-signal` (retention: 1 day).

#### `pull_request.yaml`

- **Goal**: Run static analysis, Go unit tests, code formatting checks, and commit validations across a hermetic Nix environment.
- **Triggers**: `pull_request`, `push` (`main`, `release/v*`), `merge_group`.
- **Permissions**: `contents: read`.
- **Matrix Jobs (11 Hermetic Nix Jobs)**:
  1. `build`: Compiles provider binary (`make build`).
  2. `test`: Runs Go unit test suite (`make test`).
  3. `lint`: Executes `golangci-lint` (`make lint`).
  4. `fmt`: Checks Go formatting with `gofmt` (`make fmt`).
  5. `codegen`: Runs `make generate` and asserts clean git status.
  6. `shellcheck`: Lints all repository shell scripts.
  7. `actionlint`: Lints GitHub workflow definitions.
  8. `terraform-fmt`: Lints Terraform files (`terraform fmt -check -recursive`).
  9. `tflint`: Analyzes example Terraform configurations.
  10. `validate-commit-message`: Validates Conventional Commits and product directory impacts.
  11. `sensitive-checks`: Scans Go schema files for unmasked sensitive keywords.

#### `dependabot-approve.yml`

- **Goal**: Automatically approve and enable auto-merge for qualifying Dependabot PRs.
- **Triggers**: `pull_request_target` (`opened`, `reopened`).
- **Permissions**: `contents: write`, `pull-requests: write`.
- **Behaviors**: Evaluates Dependabot actor and updates PR with automated approval when criteria are met.

---

### 3.2 Merge & Queue Automation

#### `auto-merge-listener.yml`

- **Goal**: Capture PR review events and forward approval signals to the auto-merge engine.
- **Triggers**: `pull_request_review` (`submitted`), `pull_request` (`auto_merge_enabled`).
- **Permissions**: `contents: read`.
- **Behaviors**: Emits an `approval-signal.json` artifact containing review ID, state, commit SHA, and author login. Dispatches downstream evaluation.

#### `auto-merge.yml`

- **Goal**: Evaluate PR eligibility and enable GitHub auto-merge via GraphQL.
- **Triggers**: `workflow_run` (after `auto-merge-listener`), `workflow_dispatch`.
- **Permissions**: `contents: write`, `pull-requests: write`.
- **Script**: `.github/workflows/scripts/auto-merge.js`.
- **Key Behaviors**:
  - Authenticates using Vault token `pull_requests--write`.
  - Verifies approving reviewer has qualifying repository permission (`triage`, `write`, or `admin`).
  - Asserts no active `CHANGES_REQUESTED` reviews block the PR.
  - Rejects PRs that modify `.github/**` (security boundary).
  - Executes GraphQL mutation `enablePullRequestAutoMerge`.

---

### 3.3 Issue Tracking & Backport Orchestration

#### `tracking-issue.yml`

- **Goal**: Adopt the issue referenced in a PR targeting `main` as an official tracking issue.
- **Triggers**: `workflow_run` (after `Validate PR Description`), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/releases.js` (`SCRIPT_MODE: tracking-issue`).
- **Key Behaviors**:
  - Downloads `validate-pr-signal` artifact to resolve target PR.
  - Labels referenced issue with `internal/tracking`.
  - Resolves target release branch labels (e.g. `release/v15`, `release/v16`) from PR labels or defaults to latest release branch.
  - Automatically dispatches `backport-issues.yml` for each target release branch.

#### `backport-issues.yml`

- **Goal**: Create linked backport sub-issues for each release branch label on a tracking issue.
- **Triggers**: `issues` (`labeled`), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: backport-issues`).
- **Key Behaviors**:
  - Triggers when a `release/v*` label is added to an issue.
  - Creates a child sub-issue (e.g. `[v15] <title>`) linked to the parent tracking issue.
  - Omits initial assignees so sub-issues remain available for QA pickup.
  - Idempotent: Skips creation if a sub-issue for that release branch already exists.

#### `backport-prs.yml`

- **Goal**: Cherry-pick merged commits on `main` into targeted release branches.
- **Triggers**: `push` (`main`), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: backport-pr`).
- **Key Behaviors**:
  - Uses Vault token `pull_requests--write` to push cherry-pick branches and open PRs.
  - Identifies target release branches from the linked tracking issue.
  - Creates branches named `backport-<PR>-v<major>`.
  - Automatically opens PR against target `release/v<major>`.

#### `close-tracking-issue.yml`

- **Goal**: Evaluate closed sub-issues, enforce milestone presence, close tracking issue, and transition status to `Ready For Release`.
- **Triggers**: `issues` (`closed`), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: close-tracking-issue`).
- **Key Behaviors**:
  - When an issue is closed, checks if it is a sub-issue of an `internal/tracking` parent.
  - Queries GitHub GraphQL API for all linked sub-issues.
  - **Milestone Guardrail**: Verifies the tracking issue has a milestone. If missing, attempts to inherit from closed sub-issues. If none found, skips closure and posts a notification comment.
  - If all sub-issues are closed and milestone is present, closes the tracking issue and sets status to `Ready For Release`.

---

### 3.4 Release Branch Lifecycle & Release Candidate (RC) Automation

#### `release-branch-merge.yml`

- **Goal**: Detect merges on `release/v*`, label addressed issues with `internal/merged`, and transition issue status to `Ready for QA`.
- **Triggers**: `workflow_run` (`Release`, `Release Please` on `release/v*`), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: merge-label`).
- **Key Behaviors**:
  - Resolves PR associated with the merge commit.
  - Finds referenced backport issues.
  - Adds label `internal/merged`.
  - Updates issue status to `Ready for QA`.

#### `release-please.yml`

- **Goal**: Maintain release PR staging, orchestrate event-driven manifest synchronization, generate automated RC releases on backport merge, and dispatch full releases.
- **Triggers**: `push` (`release/v*`), `workflow_call`.
- **Permissions**: `contents: write`, `pull-requests: write`.
- **Scripts**: `.github/workflows/scripts/releases.js` (`prepare-rc`, `rc-notify`), `.github/workflows/scripts/sync-manifest.js`.
- **Key Behaviors**:
  - Runs Google Release Please action to evaluate releases and manage staging PRs.
  - **Event-Driven Manifest Sync** (`sync-manifest-pr`):
    - Triggers automatically upon completion of `release-please-pr`.
    - Invokes `sync-manifest.js` to ensure `.release-please-manifest.json` and `release-please-config.json` anchors remain synchronized with published release tags.
  - **RC Path** (`release_created != 'true'`):
    - Computes next RC tag (`v<major>.<minor>.<patch>-rc.<N>`).
    - Skips branch setup commits, manifest syncs, and release chores.
    - Pushes RC tag and calls `release.yml`.
    - Comments on matching issues with RC release URL.
  - **Full Release Path** (`release_created == 'true'`):
    - Calls `release.yml` with the official release tag.
    - Dispatches post-release completion.

#### `release.yml`

- **Goal**: Compile provider binaries, generate Terraform Registry manifests, GPG-sign assets, and publish GitHub releases.
- **Triggers**: `workflow_call`, `workflow_dispatch`.
- **Permissions**: `contents: write`.
- **Script**: `.github/workflows/scripts/releases.js` (`dispatch-tag`, `prepare-manifest`, `verify-tag`), `goreleaser.sh`.
- **Key Behaviors**:
  - Validates release tag format (`validateReleaseTag`).
  - Prepares `terraform-registry-manifest.json`.
  - Executes GoReleaser via `goreleaser.sh` with GPG signing.
  - Uploads compiled binaries and SHA256 checksums to GitHub release.

#### `release-complete.yml`

- **Goal**: Finalize Rancher2 Terraform Project issues after a full production release.
- **Triggers**: `workflow_run` (`Release Please` completed), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: release-done`).
- **Key Behaviors**:
  - Verifies release tag represents a stable release (non-draft, non-RC).
  - Queries Rancher2 Terraform Project for issues in `Ready For Release`.
  - Updates status of released issues to `Done`.

---

### 3.5 Branch Setup & Synchronization

#### `setup-release-branch.yml`

- **Goal**: Initialize a new canonical `release/v<N>` maintenance branch, or idempotently update existing release branches with new or modified caller stubs.
- **Triggers**: `workflow_dispatch` (input: `branch`).
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/setup-release-branch.js`.
- **Key Behaviors**:
  - Enforces canonical naming (`release/v<N>`).
  - **New Branches**: Branches strictly from `origin/main`, cleans up non-stub workflows, writes all standard caller stubs pointing to `@main`, initializes `rancher2/release.md`, `release-please-config.json`, and `.release-please-manifest.json`, and opens an initialization PR.
  - **Existing Branches (Idempotent Update)**: Fetches the branch, checks out a chore branch, preserves active manifest versions and config anchors, synchronizes caller stubs (`pull_request.yaml`, `release.yml`, `e2e.yml`, `sync-release-manifest.yml`), and cleans up non-stub workflows. If changes are detected, opens an update PR labeled `internal/pr-release-chore`; if already up-to-date, exits cleanly with no-op.

#### `sync-release-manifest.yml`

- **Goal**: Manually synchronize `.release-please-manifest.json` and `release-please-config.json` versions across active release branches and `main` out-of-band.
- **Triggers**: `workflow_dispatch`, `workflow_call`.
- **Permissions**: `contents: write`, `pull-requests: write`, `id-token: write`.
- **Script**: `.github/workflows/scripts/sync-manifest.js`.
- **Key Behaviors**:
  - Aligns version records across branches to prevent semver downgrades during backports.
  - Allows maintainers to provide manual overrides for `branch`, `version`, `release_as`, and `last_release_sha`.

---

### 3.6 Compliance & Infrastructure Maintenance

#### `cleanup-e2e.yml`

- **Goal**: Clean up orphaned AWS infrastructure from completed acceptance test runs.
- **Triggers**: `schedule` (daily 13:00 UTC), `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/e2e.js` (`find-completed-e2e-runs`, `execute-cleanup`).
- **Key Behaviors**:
  - Queries completed `e2e.yml` runs.
  - Assumes AWS CI role (`arn:aws:iam::270074865685:role/terraform-module-ci-test`).
  - Runs `cleanup.sh` in parallel matrix inside hermetic Nix container.

#### `e2e.yml`

- **Goal**: Run live acceptance tests on AWS infrastructure.
- **Triggers**: `workflow_dispatch`.
- **Permissions**: `contents: read`.
- **Script**: `.github/workflows/scripts/e2e.js`, `run_tests.sh`.

#### `fossa.yml`

- **Goal**: Scan dependencies for license compliance and security vulnerabilities.
- **Triggers**: `push` (`main`), `pull_request`.
- **Permissions**: `contents: read`.

---

## 4. Workflow Scripts Reference Guide

All automation logic resides in `.github/workflows/scripts/`.

| Script File                      | Core Responsibilities & Handlers                                                                                                      | Underlying Goals                                                                                                                                                                                                                      |
| :------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **`auto-merge.js`**              | `runPreflight`, `runAutoMerge`, `evaluateReviews`, `checkReviewerPermission`, `checkWorkflowModifications`                            | Evaluates pull request approval signals, verifies reviewer has qualifying collaborator role (`triage`, `write`, or `admin`), ensures no change requests exist, guards `.github/` modifications, and activates auto-merge via GraphQL. |
| **`backports.js`**               | `runBackportPr`, `runBackportIssues`, `runMergeLabel`, `runCloseTrackingIssue`, `runReleaseDone`, `updateIssueProjectStatus`          | Multi-mode handler responsible for backport PR creation, sub-issue generation, issue merge labeling, issue status updates, tracking issue closure with milestone inheritance, and final release completion.                           |
| **`releases.js`**                | `runTrackingIssue`, `runPrepareRc`, `computeNextRcTag`, `runRcNotify`, `runVerifyTag`, `runDispatchTag`, `runPrepareRegistryManifest` | Manages tracking issue adoption, RC semver calculations, RC tagging/notifications, release tag verification, and Terraform registry manifest preparation.                                                                             |
| **`setup-release-branch.js`**    | `runSetupReleaseBranch`, `parseReleaseBranch`, `buildReleasePleaseConfig`, `buildWorkflowStubs`, `cleanupNonStubWorkflows`            | Initializes canonical release branches (`release/v<N>`) or idempotently updates existing release branches, ensuring caller stubs match `@main`, non-stub workflows are cleaned up, and manifests are preserved.                       |
| **`sync-manifest.js`**           | `runSyncManifest`, `calculateNextMinorVersion`, `updateManifestContent`, `updateConfigContent`                                        | Aligns `.release-please-manifest.json` versions across active release branches to prevent backport version regressions.                                                                                                               |
| **`validate-pr-description.js`** | `validatePrDescription`, `validateIssueExists`, `run`                                                                                 | Validates PR descriptions for valid `- Addresses: #<issue>` references, verifies issue status via Octokit, and writes `signal/pr-signal.json`.                                                                                        |
| **`validate-commit-message.js`** | `validateMessage`, `validatePrefix`, `validateProductImpact`, `isProductFile`                                                         | Enforces Conventional Commits, validates subject length under 70 characters, checks spelling via cSpell, and restricts `feat:` commits to product directories.                                                                        |
| **`lint.js`**                    | `runTerraform`, `runActionlint`, `runShellcheck`, `runNodeCheck`, `runEslint`, `runGitleaks`, `runManifestCheck`                      | Unified linting wrapper orchestrating terraform fmt, actionlint, shellcheck, node syntax check, eslint, gitleaks secret scanning, and manifest downgrade checks.                                                                      |
| **`e2e.js`**                     | `runFilterPr`, `runE2EStatus`, `runFindCompletedE2eRuns`, `runExecuteCleanup`                                                         | Filters E2E test triggers, posts commit statuses, discovers completed test runs, and orchestrates AWS resource cleanup via `cleanup.sh`.                                                                                              |
| **`goreleaser.sh`**              | Bash script wrapping GoReleaser                                                                                                       | Validates GPG signing keys, generates multi-arch release archives, creates checksums, and signs release manifests.                                                                                                                    |
| **`nix-run.sh`**                 | Bash script wrapper for Nix execution                                                                                                 | Loads dependencies from `flake.nix` and executes commands within a hermetic Nix environment.                                                                                                                                          |
| **`check_sensitive.go`**         | Go utility scanning schema files                                                                                                      | Parses Go AST schema definitions in `rancher2/` to ensure sensitive parameters have `Sensitive: true`.                                                                                                                                |

---

## 5. Rancher2 Terraform Project Lifecycle & State Machine

All tracking and backport issues are tracked on [**Rancher2 Terraform Project**](https://github.com/orgs/rancher/projects/73).

```
┌──────────────────┐       Backport PR Merged to release/v*       ┌──────────────────┐
│  Backlog / Open  │ ───────────────────────────────────────────► │   Ready for QA   │
└──────────────────┘                                              └──────────────────┘
                                                                            │
                                                                            │ All sub-issues closed
                                                                            │ Milestone assigned
                                                                            ▼
┌──────────────────┐            Full Release Completed            ┌──────────────────┐
│       Done       │ ◄─────────────────────────────────────────── │ Ready For Release│
└──────────────────┘                                              └──────────────────┘
```

### Column States & Automated Transitions

1. **`Ready for QA`**:
   - **Trigger**: Merge of an automated backport PR into `release/v*`.
   - **Handler**: `release-branch-merge.yml` -> `backports.js` (`runMergeLabel`).
   - **Logic**: Adds issue to project if not already present, updates status field to `Ready for QA`.
2. **`Ready For Release`**:
   - **Trigger**: Closure of a backport sub-issue by QA.
   - **Handler**: `close-tracking-issue.yml` -> `backports.js` (`runCloseTrackingIssue`).
   - **Milestone Guardrail**: Verifies all sub-issues are closed and parent tracking issue has an active milestone. If unassigned, attempts milestone inheritance from closed sub-issues. If missing, blocks closure and comments on the issue.
3. **`Done`**:
   - **Trigger**: Successful completion of a full production release on `release/v*`.
   - **Handler**: `release-complete.yml` -> `backports.js` (`runReleaseDone`).
   - **Logic**: Filters Rancher2 Terraform Project for items in `Ready For Release` and advances them to `Done`. Ignores Release Candidates and Draft releases.

---

## 6. Architectural Decisions & Engineering Standards

### 6.1 Architectural Decisions

- **`D1`: Issue-First Decoupling**: PR validation requires an open issue, but developers do not create tracking or backport issues manually. Automation (`tracking-issue.yml`) converts the referenced issue into a tracking issue and manages sub-issues.
- **`D2`: Trunk-Based Automated Backports**: Developers write and review code exclusively against `main`. Automated cherry-pick PRs eliminate manual cherry-picking and divergence across branches.
- **`D3`: Event-Driven RC Generation**: RC tags and binary releases are created automatically when backport PRs merge into release branches, eliminating manual RC dispatching.
- **`D4`: Direct GraphQL Project Integration**: Automation uses GitHub's GraphQL API (`projectV2` mutations) to query and update Rancher2 Terraform Project directly, avoiding fragile third-party actions.
- **`D5`: Milestone Guardrail**: Tracking issues cannot close or transition to `Ready For Release` without an assigned milestone, ensuring every released fix is tied to a release version.
- **`D6`: Caller Stubs on Release Branches**: Release branches (`release/v*`) contain lightweight caller stubs pointing to workflow implementations on `@main`, ensuring centralized workflow management without backporting workflow modifications.
- **`D7`: Hermetic Nix Environment**: All CI checks and tests run inside containers built from `flake.nix`, guaranteeing reproducible builds across development machines and CI runners.
- **`D8`: Fail-Closed Security Handlers**: All automation scripts fail closed if required environment variables, commit SHAs, or Vault tokens cannot be resolved.
- **`D9`: Idempotent Execution**: All scripts check existing state (existing labels, sub-issues, issue project fields, remote Git tags) before applying mutations, allowing safe re-runs.
- **`D10`: Modular JavaScript Architecture & Unit Testability**: Rather than embedding complex shell scripts or multi-line inline JavaScript directly into workflow YAML definitions, all business logic and orchestration routines are externalized into discrete ES modules in `.github/workflows/scripts/`. Each module is designed around dependency injection (accepting `github`, `context`, `core`, `process`, and execution helpers as parameters). This enables comprehensive mocking and rapid local unit test execution via Node's native test runner (`node --test`), catching edge cases, regressions, and API payload discrepancies before automation ever runs in production GitHub Actions.

---

### 6.2 Security & Authentication Model

1. **Zero Root Permissions**: Every workflow declares `permissions: {}` at root level. Jobs declare only minimal required read/write scopes.
2. **HashiCorp Vault OIDC Exchange**: Privileged operations exchange GitHub OIDC tokens for short-lived Vault access tokens via `rancher-eio/read-vault-secrets`:
   - `github/token/rancher--terraform-provider-rancher2--pull_requests--write`: Auto-merging, pushing cherry-pick branches, opening backport PRs, tagging RCs.
   - `github/token/rancher--terraform-provider-rancher2--issues--write`: Creating issues, managing milestones, mutating issue project fields.
3. **AWS OIDC Role**:
   - `arn:aws:iam::270074865685:role/terraform-module-ci-test`: Used exclusively by acceptance test runs and scheduled infrastructure cleanups.

---

### 6.3 Script Execution & Local Safety

1. **Hermetic Test Execution**:
   All workflow scripts are covered by extensive unit test suites under `.github/workflows/scripts/*.test.js`. Tests mock Octokit REST/GraphQL calls, GitHub Actions context payloads, and child processes, executing hundreds of assertions across edge cases in under one second without making external network calls.
2. **Local Environment Protection**:
   All workflow scripts in `.github/workflows/scripts/` guard live Git and network mutations behind:
   ```javascript
   if (process.env.GITHUB_ACTIONS === "true") {
     // execute remote git or api mutation
   }
   ```
   This ensures that running unit test suites locally (e.g. `node --test` or `make test`) never alters local developer `.git/config`, remotes, or global configurations.
