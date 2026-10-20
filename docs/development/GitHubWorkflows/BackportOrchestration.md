# Component Specification: Backport Orchestration & Issue Tracking

- **Related Topic:** [GitHub Workflows & Automation](../GitHubWorkflows.md)
- **Related Topic:** [Release Process & Automation](../ReleaseProcess.md)
- **Overview:** [End-to-End Value Stream Lifecycle](./Overview.md)
- **Upstream:** [Pull Request & Merge Lifecycle on Main](./PullRequestLifecycle.md)
- **Downstream:** [Privileged AWS End-to-End Testing](./EndToEndTesting.md)

---

## Abstract

This component specification details the automated backport orchestration architecture of `terraform-provider-rancher2`. It coordinates tracking issue adoption, GitHub Sub-Issues hierarchy creation, and automated Git cherry-picking into targeted release branches (`release/v<major>`).

The backport system bridges human-approved pull requests on trunk (`main`) with long-lived release maintenance branches, automating issue tracking, branch creation, cherry-picking, auto-merge configuration, and collaborator assignment.

---

## Architectural Design & Two-Tier Lifecycle

The backport system operates across two distinct phases:

1. **Pre-Merge Phase (Tracking & Sub-Issue Creation)**: Discovers open pull requests on `main`, attaches tracking metadata to the referenced GitHub issue, and creates linked backport sub-issues for each targeted release branch.
2. **Post-Merge Phase (Cherry-Pick & PR Automation)**: Detects merge commits on `main`, discovers the associated tracking issue and sub-issues, cherry-picks the commit to each release branch, and opens backport PRs with auto-merge enabled.

```mermaid
flowchart TD
    subgraph PreMerge["Phase 1: Issue Tracking & Sub-Issue Tree"]
        PR_MAIN[PR opened on main with - Addresses: #ISSUE]
        CRON[tracking-issue.yml<br/>Runs every 4h or manual dispatch]
        PARSE[releases.js: tracking-issue<br/>Detects highest release/v* branch]
        LABEL_PARENT[Label Issue: internal/tracking, internal/user, release/v*]
        COMMENT_LINK[Comment: tracking-pr: #PR and instructions]
        DISPATCH_BP[Dispatch backport-issues.yml]
        CREATE_SUB[backports.js: backport-issues<br/>Create sub-issue: [release/v*] Title]
        LINK_SUB[Link via GitHub Sub-Issues API]

        PR_MAIN --> CRON --> PARSE --> LABEL_PARENT --> COMMENT_LINK --> DISPATCH_BP --> CREATE_SUB --> LINK_SUB
    end

    subgraph PostMerge["Phase 2: Cherry-Pick & Backport PR Generation"]
        MERGE[Maintainer merges PR into main]
        PUSH[backport-prs.yml on push to main]
        SETTLE[backports.js: wait-for-settle<br/>Wait for API commit indexing]
        VAULT[rancher-eio/read-vault-secrets<br/>Retrieve GITHUB_MERGE_TOKEN]
        GET_TREE[Query Sub-Issues API for linked release branches]
        CHERRY_PICK[git cherry-pick --allow-empty -x SHA -X theirs]
        OPEN_BPR[Open Backport PR targeting release/v*]
        AUTO_MERGE[GraphQL: enablePullRequestAutoMerge SQUASH]

        MERGE --> PUSH --> SETTLE --> VAULT --> GET_TREE --> CHERRY_PICK --> OPEN_BPR --> AUTO_MERGE
    end
```

---

## Technical Specifications: Workflow & Script Mechanics

### 1. Tracking Issue Adoption (`.github/workflows/tracking-issue.yml`)

- **Trigger**: Cron schedule every 4 hours (`0 */4 * * *`) or manual `workflow_dispatch`.
- **Permissions**: `issues: write`, `pull-requests: read`, `contents: read`, `id-token: write`, `actions: write`.
- **Vault Token**: Retrieves `GITHUB_MERGE_TOKEN` from HashiCorp Vault path:
  `github/token/rancher--terraform-provider-rancher2--pull_requests--write token | GITHUB_MERGE_TOKEN`
- **Execution Script**: `.github/workflows/scripts/releases.js` (`SCRIPT_MODE: tracking-issue`).
- **Logic**:
  1. **Release Branch Discovery**: Fetches all repository branches via `github.rest.repos.listBranches` and extracts branches matching `/^release\/v\d+$/`. Sorts them descending by major version (e.g. `release/v15`, `release/v14`, `release/v13`) to identify `latestReleaseBranch`.
  2. **Candidate Search**: Queries GitHub Search API for open, non-draft PRs targeting `main` that do not have `internal/ignore`, `internal/pr-backport`, or `autorelease` labels.
  3. **Issue Parsing**: Validates the PR body using `validatePrDescription(pr.body)` to extract referenced issue numbers.
  4. **Adoption & Labeling**:
     - Ensures the issue is open and not a PR.
     - Adds `internal/tracking` and `internal/user` labels.
     - Adds target release branch label (derived from PR labels or falling back to `latestReleaseBranch`).
     - Adds linking comment: `<!-- tracking-pr: #<pr_number> -->\nThis is the tracking issue for PR #<pr_number>`.
     - Adds instruction comment prompting maintainers to add any additional release branch labels required for multi-branch backports without skipping versions.
  5. **Sub-Issue Dispatch**: Uses `GITHUB_MERGE_TOKEN` to call `actions.createWorkflowDispatch` triggering `backport-issues.yml` on `main`.

### 2. Backport Sub-Issue Creation (`.github/workflows/backport-issues.yml`)

- **Trigger**:
  - Automatically fired when an issue is labeled (`issues: [labeled]`) where `startsWith(label.name, 'release/v')` and the issue is not a pull request.
  - Or manually via `workflow_dispatch` with inputs `issue_number`, `release_label`, and optional `pr_number`.
- **Execution Script**: `.github/workflows/scripts/backports.js` (`SCRIPT_MODE: backport-issues`).
  *(Note: To avoid arbitrary checkout execution, the workflow fetches `backports.js` via the GitHub REST API using the default branch SHA and imports it dynamically).*
- **Logic**:
  1. **PR Number Resolution**: If not supplied directly, parses the issue's comments in reverse chronological order for the tracking comment pattern `/<!--\s*tracking-pr:\s*#?(\d+)\s*-->/i`.
  2. **Idempotency Gate**: Queries `GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues` to verify whether a sub-issue for this release label already exists. If found, skips creation.
  3. **Issue Creation**: Creates a new issue:
     - **Title**: `[<release_label>] <parent_issue_title>`
     - **Body**: References the backport PR and parent tracking issue.
     - **Labels**: `[release_label, "internal/backport"]`
     - **Assignees**: Parsed from `TERRAFORM_MAINTAINERS`.
  4. **Sub-Issue Association**: Sends a `POST` request to `/repos/{owner}/{repo}/issues/{parent_issue_number}/sub_issues` with `{ sub_issue_id: newIssue.id }` using GitHub API version `2022-11-28`, establishing the native Sub-Issues relationship.

### 3. Automated Cherry-Pick & Backport PR Generation (`.github/workflows/backport-prs.yml`)

- **Trigger**: Runs on `push` to `main` or manually via `workflow_dispatch` taking `merge_commit_sha`.
- **Timeout**: 15 minutes.
- **Permissions**: `contents: write`, `pull-requests: write`, `id-token: write`, `issues: write`, `actions: write`.
- **Step 1: Settle (`wait-for-settle`)**:
  - Merged PRs may take several seconds to appear in GitHub's commit association index.
  - Calls `runWaitForSettle` in `backports.js`, retrying up to 5 times with exponential backoff (3s, 6s, 12s, 24s) until `listPullRequestsAssociatedWithCommit` returns a merged PR targeting `main`.
- **Step 2: Vault Token Retrieval**:
  - Retrieves privileged `GITHUB_MERGE_TOKEN` from HashiCorp Vault (`rancher-eio/read-vault-secrets`).
- **Step 3: Cherry-Pick & PR Dispatch (`backport-pr`)**:
  - Calls `runBackportPr` in `backports.js`.
  - Searches for the `internal/tracking` issue linked to the PR.
  - Queries `GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues` to retrieve all child sub-issues.
  - For each sub-issue having a label matching `/^release\/v\d{1,2}$/`:
    1. **Branch Checkout**: Checks out `origin/<target_branch>`.
    2. **Branch Creation**: Creates a fresh local branch `backport-<pr_number>-<target_branch>`.
    3. **Cherry-Pick**: Executes `git cherry-pick --allow-empty -x <merge_commit_sha> -X theirs`.
       - `-x`: Appends `(cherry picked from commit ...)` to the commit log for cryptographic traceability.
       - `-X theirs`: Automatically favors incoming changes in trivial merge conflicts. If a hard conflict occurs, the script safely aborts the cherry-pick, reports the error, and continues with remaining branches.
    4. **Push**: Pushes the backport branch to the repository.
    5. **PR Creation**: Creates a PR targeting `<target_branch>` with warning notices regarding the cherry-pick merge strategy.
    6. **Auto-Merge**: Executes a GraphQL mutation (`enablePullRequestAutoMerge`) with merge method `SQUASH`.
    7. **Assignment & Labels**: Assigns maintainers from `TERRAFORM_MAINTAINERS` and applies labels `internal/pr-backport` and `<target_branch>`.

---

## Standing Implementation Decisions

1. **Sub-Issues Graph Architecture**: Backports are modeled as explicit sub-issues under the main tracking issue. This gives engineering managers and release teams unified visibility into all release branches affected by a single bug or feature.
2. **Deterministic Cherry-Pick Flagging**: The cherry-pick command always uses `-x` to preserve the original trunk commit SHA in the Git log.
3. **Automated Auto-Merge with Squash**: Backport PRs have GitHub auto-merge pre-enabled. Once the required CI and E2E tests pass on the backport branch, GitHub merges the PR without requiring additional manual clicks from maintainers.
4. **Token Isolation via Injected Factory**: The runner uses the injected `getOctokit(mergeToken)` factory specifically for operations requiring elevated PR write and auto-merge permissions, while standard issues queries retain the default token, enforcing least-privilege separation.
5. **Multi-Branch Continuity Rule**: When backporting across multiple major versions (e.g. `release/v13` through `release/v15`), maintainers must not skip intermediate branches, ensuring schema and state migrations remain backward compatible.
