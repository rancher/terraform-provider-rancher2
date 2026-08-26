# Component Specification: PR Executor Target PR Resolution

- **Related Topic:** [GitHub Workflows & Automation](../GitHubWorkflows.md)
- **Target Component:** `.github/workflows/pr-executor.yml`

---

## Abstract

This component details the architectural modifications to the PR Executor and Auto-Merge workflow (`.github/workflows/pr-executor.yml`) to ensure reliable pull request resolution. It resolves a platform-level limitation where the GHA `workflow_run` event payload's `pull_requests` array is left empty when triggered by a pull request originating from a fork.

In accordance with our workflow controller standards, the logic is extracted from the GHA workflow file into a standalone, testable Node.js module at `.github/workflows/scripts/get-target-pr.js`. The repository checkout step is moved to the top of the job to facilitate executing this module.

---

## 1. Architectural Strategy & Context

The `PR Executor and Auto-Merge` workflow is triggered upon the completion of a parent `pull_request` or `pull_request_review_trigger` workflow. To verify and merge the correct pull request, it must resolve the triggering pull request number.

The workflow extracts this number directly from `github.event.workflow_run.pull_requests[0].number` or queries `getWorkflowRun`. However, due to GitHub Actions security and API constraints, the `pull_requests` array is empty when the parent workflow is triggered by a pull request originating from a fork repository. Since all human contributions in this repository are required to originate from forks, this blocks the automated merge execution pipeline.

### Robust Fallback Strategy

When both the payload and the direct `getWorkflowRun` query fail to provide a `pull_requests` entry, the executor executes a multi-layered fallback strategy:

1. **Associated Commit Lookup**: Fetches any pull requests associated with the head commit SHA (`parentRun.head_sha`) using the GitHub REST API (`listPullRequestsAssociatedWithCommit`).
2. **Open PR Search**: Fetches all active open pull requests for the repository and matches the head commit SHA (`parentRun.head_sha`) or branch name (`parentRun.head_branch`) against the active head ref (`pr.head.sha` or `pr.head.ref`).

To keep workflows clean and adhere to controller-only design standards, this logic is modularized into `.github/workflows/scripts/get-target-pr.js` and thoroughly unit-tested.

### `pull_request_review_trigger` Workflow Checkout Requirement

The `pull_request_review_trigger` (`review-trigger.yml`) workflow is a parent/trigger workflow for the PR Executor. It executes the script `.github/workflows/scripts/log-trigger.sh` upon detecting valid comments.

- **Checkout Step**: Like other controller workflows, it checks out the codebase prior to executing any script located under `.github/workflows/scripts/`. Without the checkout step, the runner environment lacks access to the script.
- **Implementation**: The workflow uses `contents: read` permissions and runs the standard `actions/checkout@v7.0.1` step (targeting the `main` branch) within the `trigger` job prior to running the script.

---

## 2. Security Analysis & Threat Mitigations

Because the PR Executor runs with elevated permissions (`contents: write`, `pull-requests: write`), security and tamper-proofing are paramount:

### A. Fork Branch Name Hijacking Protection (Collision Prevention)

If the workflow fell back to matching a branch name like `patch-1` or `update` alone, a malicious actor could name a branch on their fork to collide with a trusted branch/PR, potentially hijacking the target `prNumber` and triggering a merge run.

- **Mitigation**: The fallback branch matching logic strictly validates the fork repository owner:

  ```javascript
  p.head.ref === parentRun.head_branch && p.head.repo.owner.login === parentRun.head_repository.owner.login;
  ```

  This ensures branch-name matching only succeeds if the PR originates from the exact same fork that triggered the workflow run.

### B. Event & Code Isolation (Base Ref Guarantee)

A bad actor might attempt to modify the `pr-executor.yml` workflow file or the validation/merge scripts inside their PR branch to bypass checks or run arbitrary code.

- **Mitigation**: GitHub Actions executes `workflow_run` workflows **exclusively using the workflow YAML file from the base repository's default branch ref**. Furthermore, the `actions/checkout` step checks out the default branch ref of the base repository by default. As a result, all executed YAML definitions and imported scripts (`verify-pr-requirements.mjs`, `merge-pr.js`) are read strictly from the base repository's trusted default branch commit, making PR branch tampering impossible.

### C. Immutable Gated Verification Checks

Even if a bad actor manages to resolve their PR number in the executor, they cannot bypass the security requirements. The `Verify PR Requirements` job enforces:

1. **Verified GPG Signatures**: Every single commit in the PR must be signed and verified by GitHub.
2. **Trusted Role Approval**: The PR must have at least one approval from a trusted role (Owner, Member, Collaborator) with write permissions.
3. **AI Review Gate**: The PR must have a valid AI review from Copilot or our agent.
4. **Resolved Review Conversations**: 100% of review comment threads must be resolved.
5. **Trusted /merge Trigger**: For human PRs, a `/merge` comment must be explicitly posted in the conversation thread by a trusted repository member.

---

## 3. Technical Blueprint

The inline `actions/github-script` step in `.github/workflows/pr-executor.yml` imports and executes the external module:

```javascript
const scriptPath = `${process.env.GITHUB_WORKSPACE}/.github/workflows/scripts/get-target-pr.js`;
const { default: script } = await import(scriptPath);
return await script({ github, context, core });
```

The external module `.github/workflows/scripts/get-target-pr.js` exports the robust resolution logic:

```javascript
export default async ({ github, context, core }) => {
  let prNumber;
  const parentRun = context.payload.workflow_run;

  // ... (PR resolution & fallback lookup logic) ...

  return prNumber;
};
```
