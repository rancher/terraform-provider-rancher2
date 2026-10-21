# Project Management & Delivery Guide: Value Stream Navigation

## Abstract

This guide provides a high-level view of how work moves from initial issue creation to customer release in the Rancher2 Terraform Provider project. 

Rather than focusing on technical GitHub Actions workflow triggers, runner configurations, or script implementations, this document explains the **lifecycle of work**, how **developers, QA, and automation** create and move data through the system, and how to observe state transitions and answer key status questions at each stage of delivery.

---

## 💡 Key Points of Interest

The following nuances address common points of interest and areas of confusion encountered when tracking work through the pipeline:

- **The "Main" Issue is the Public Anchor**: Created by community members, users, developers, or project managers, the main issue is the primary public point of contact for anyone searching for known bugs, reproduction details, or enhancement requests.
- **Release Labels Designate Target Branches**: Target release branch labels (e.g., `release/v15`, `release/v14`) indicate which maintenance branches require the change. While tracking automation actively supports placing release labels directly on the PR (extracting them to select the target branch and syncing them to the main issue), placing them on the **main issue** serves as the public source of truth.
- **Sub-Issues Must Exist Before the PR Merges**: Release branch child sub-issues must be created **before** the pull request on `main` is merged in order for automated backport cherry-picking to trigger.
- **Backports Auto-Merge on Passing Tests (Best-Effort)**: When an automated backport pull request passes acceptance and integration testing, automation attempts to automatically merge it into the target release branch. If auto-merge cannot be enabled (e.g., due to branch protection or configuration restrictions), maintainers perform a manual merge fallback once tests pass.
- **Backport Issues Become QA Issues**: Once a backport merges to its release branch and is packaged into a Release Candidate (RC), that backport sub-issue transitions into the active QA verification issue for that release line.
- **QA Engineers Close QA Sub-Issues**: QA engineers validate functionality on the RC build and close their respective backport/QA sub-issue once verification is complete.
- **Main Issues Close When the Release Happens**: The main parent issue stays open across all release branch testing and closes automatically only when the official production release is published.

---

## 👥 Key System Actors

The following actors create, review, verify, and package changes through the value stream:

| Actor | Primary Activities in the System | Key System Artifact |
| :--- | :--- | :--- |
| **User / Reporter** | Submits bug reports, feature requests, and reproduction steps. | The original **Main Issue**. |
| **Developer / Contributor** | Authors code and automated tests, associates PRs with the main issue, and resolves review feedback. | Pull Request opened against `main` referencing `- Addresses: #<main-issue>`. |
| **Repository Maintainer** | Conducts code reviews, ensures quality and safety, squash-merges PRs to `main`, and reviews final release PRs. | PR approved and squash-merged into trunk. |
| **Automation ("The Bot")** | Adopts main issues, generates backport branches/PRs, runs automated integration tests, auto-merges passing backports, cuts RCs, and drafts changelogs. | Sub-issue tree, backport PRs, RC builds, and QA notification comments. |
| **QA Engineer** | Validates fixes and features on live test environments using Release Candidate builds, closing the backport/QA sub-issue when finished. | QA verification sign-off and sub-issue closure. |

---

## 🔄 The 7-Stage Delivery Value Stream

The diagram below illustrates how work items progress through the delivery lifecycle:

```text
[ Intake: Main Issue Created ] ──► (Primary Public Point of Contact)
              │
              ▼
[ 1. Development & PR on Main ] ──► (- Addresses: #<main-issue>)
              │
              ▼
[ 2. Release Labels & Sub-Issues Added ] ──► (MUST exist BEFORE PR merge)
              │
              ▼
[ 3. Maintainer Review & Squash-Merge ] ──► (Trunk Commit on main)
              │
              ▼
[ 4. Automated Backport PRs Generated ] ──► (Cherry-picks to release/v*)
              │
              ▼
[ 5. Acceptance Testing & Auto-Merge ] ──► (Best-effort auto-merge when tests pass)
              │
              ▼
[ 6. RC Staging & QA Validation ] ──► (Backport issue becomes QA issue;
              │                          QA closes sub-issue upon completion)
              ▼
[ 7. Production Release ] ──► (Main issue closes upon final release)
```

---

### Intake: The "Main" Issue (Primary Public Anchor)

All work begins with a **Main Issue** filed on GitHub.

- **Purpose & Ownership**: 
  - The main issue can be created by external users, community contributors, customers, developers, or project managers.
  - It acts as the single primary point of contact for anyone searching the repository for known bugs, ongoing discussions, reproduction steps, or planned feature requests.
  - The main issue remains open throughout development and intermediate release branch testing.

---

### Stage 1: Development & PR on `main`

All feature development and bug fixes are implemented against `main` (trunk) first. Work is never developed directly on maintenance or release branches.

- **What Happens:** 
  - A developer creates a feature branch and opens a Pull Request targeting `main`.
  - The PR description must explicitly link the main issue:
    ```markdown
    - Addresses: #1234
    ```
  - Automated continuous integration runs hermetic builds, unit tests, linters, security audits, and commit format checks.

---

### Stage 2: Release Branch Labeling & Sub-Issue Preparation

Before the PR on `main` is merged, the target release hierarchy is established.

- **Release Line Assignment**:
  - Target release branch labels (such as `release/v15`, `release/v14`) identify which maintenance branches require backports.
  - Labels may be added directly to the **main issue** or placed on the **pull request** (where automation detects them, selects the target release branch, and synchronizes the label to the main issue).
  - Automation inspects these labels to identify which maintenance branches require backports.
- **Sub-Issue Creation Prerequisite**:
  - Release branch child sub-issues (e.g., `[release/v15] Fix cluster proxy config`) **must be created before the PR is merged into `main`**.
  - When a PR merges to `main`, automation queries the main issue for existing child sub-issues to determine which branches to cherry-pick into. If sub-issues are not present at the time of merge, automated backports will not trigger.
  - Sub-issues are created automatically when the main issue is adopted by tracking automation, or manually added under GitHub's **Sub-issues** panel on the main issue.

---

### Stage 3: Maintainer Review & Squash-Merge

Code cannot enter `main` without explicit human review and approval.

- **What Happens:**
  - Maintainers review the architecture, tests, and documentation.
  - Once approved and all CI checks pass, maintainers perform a **Squash-Merge** into `main`.
  - Merging produces a single atomic commit on `main` following Conventional Commits format (e.g., `feat: ...` or `fix: ...`).

---

### Stage 4: Automated Backport Generation

Once the PR merges to `main`, automation processes the backport tree.

- **What Happens:**
  - Automation reads the sub-issues linked under the main issue.
  - For each target branch identified in the sub-issues (e.g., `release/v15`), automation cherry-picks the squash commit from `main` into a dedicated backport branch.
  - A backport pull request targeting the release branch is opened automatically.
  - If a merge conflict occurs during cherry-picking, automation notifies maintainers to resolve the conflict manually.

---

### Stage 5: Acceptance Testing & Auto-Merge

Backports undergo automated testing before entering maintenance lines.

- **What Happens:**
  - Backport PRs run automated end-to-end (E2E) acceptance tests on real cloud infrastructure.
  - **Automated Merge (Best-Effort)**: When all acceptance tests pass, automation enables **auto-merge** to merge the backport PR into the target release branch. If auto-merge fails to engage (e.g., due to branch rules or permissions), repository maintainers manually review and merge the passing backport.

---

### Stage 6: Release Candidate (RC) Staging & QA Validation

Merging into a release branch triggers packaging and hands off work to QA.

- **What Happens:**
  - Automation calculates the next RC version tag (e.g., `v15.2.0-rc.1`), builds and cryptographically signs binaries, and creates a GitHub pre-release.
  - Automation posts a notification comment on the main issue with the pre-release link for testing.
  - **Backport Issues Become QA Issues**: Once the backport merges, its corresponding child sub-issue transitions into the active QA verification issue for that release branch.
  - **QA Sub-Issue Closure**: QA engineers perform testing against the RC build in test environments. When validation is successful, the **QA engineer closes the child sub-issue**.

---

### Stage 7: Production Release & Main Issue Closure

When all planned items for a release line have completed QA testing, the production release is cut.

- **What Happens:**
  - Maintainers merge the pending Release PR on `release/v*`.
  - Official release tags (e.g., `v15.2.0`) are generated, final binaries are signed, and the provider is published to the **Terraform Registry**.
  - **Main Issue Closes**: The main parent issue, having remained open while sub-issues were validated across release branches, **automatically closes upon the final production release**.

---

## 🏷️ Visual Status Cheat Sheet

The following indicators mark work as it moves through GitHub Issues and Pull Requests:

| Indicator / Label | Where to Find It | What It Means | Next State in Pipeline |
| :--- | :--- | :--- | :--- |
| **Main Issue** | Issues tab | Primary public point of contact for bug reports and enhancements. | PR opened on `main`; `release/v*` labels added. |
| **PR on `main`** | Pull Requests tab | Development of change targeting trunk (`- Addresses: #<issue>`). | CI matrix checks pass; maintainer review. |
| **`release/v<major>`** | **Main Issue** / PR (Labels) | Designates target release branch (synced to main issue by automation if placed on PR). | Sub-issues created for backport generation. |
| **Sub-issues Widget** | Main Issue (Sidebar) | Displays child sub-issues per release branch. **Must exist before PR merge.** | Directs automated cherry-pick targets upon merge to `main`. |
| **Backport PR** | Pull Requests tab | Automated cherry-pick PR targeting a release branch. | Auto-merges once tests pass (with manual fallback if auto-merge fails). |
| **Child Sub-Issue** | Sub-issues list | Transitions into the **QA verification issue** once backport merges. | Closed by QA engineer after testing the RC build. |
| **RC Notification Comment** | Main Issue (Comments) | Automation posts the `vX.Y.Z-rc.N` pre-release link for testing. | QA validates the RC build. |
| **Main Issue Closed** | Main Issue | Final production release published to Terraform Registry. | Work is complete and officially released. |

---

## ❓ Answering Common Status Questions

### 1. Has this bug fix or feature landed in trunk?
- Open the main issue.
- Inspect the linked pull request. If the PR status is **Merged** into `main`, the change has successfully landed in trunk.

### 2. Which release branches will receive this change?
- Check the **Sub-issues** list on the main issue.
- Each child issue is prefixed with the target branch name (e.g., `[release/v15]`, `[release/v14]`). 
- **Important**: If a release branch is missing, ensure the release label is on the main issue or PR, and ensure the sub-issue exists **before** merging the PR on `main`.

### 3. Will the backport PR require a manual merge?
- In most cases, no. Backport PRs have auto-merge enabled automatically and will merge once acceptance tests pass. However, auto-merge is best-effort; if auto-merge cannot be engaged (or is blocked by repository rules), maintainers will manually merge the PR once tests are green.

### 4. Is the fix ready for QA testing?
- Look for an automated comment from `github-actions[bot]` on the main issue stating that a Release Candidate (e.g., `v...-rc...`) is available.
- When this comment appears, the backport sub-issue is active as a QA issue ready for validation.

### 5. Who closes the child sub-issue?
- The QA engineer closes the child sub-issue after validating the fix on the RC build.

### 6. When does the main issue close?
- The main issue stays open across all sub-issue testing cycles and closes automatically when the official production release is published to the Terraform Registry.
