# Standard Gated Development Process

---

## Abstract

This component outlines the repository's standard, step-by-step developer and agent development process. It is structured around a **Gated 4-Phase Lifecycle** (Plan, Implement, Review, Commit) and **Three Authoritative Approval Gates** (two of which are user-facing) to enforce strict planning, automated quality, and cryptographic biometric commits.

---

## 🔒 The Three Authoritative Approval Gates

To maintain absolute system integrity and prevent unvetted code modifications, the development lifecycle is anchored around three sequential approval gates. The agent operates with full autonomous authorization between gates, but is strictly blocked from advancing phases or committing code until the respective gate is satisfied.

### **Gate 1: Planning Gate (User-Facing)**

- **Phase Transition**: Plan $\rightarrow$ Implement.
- **Verification File**: `plan-approval.json` and `plan-approval.json.sig`.
- **Protocol**: Before any source or configuration files can be modified, the developer must review and cryptographically sign the dynamic implementation plan checklist using SSH Keys (via `ssh-keygen -Y sign`).
- **Authorization**: Once Gate 1 is signed, the agent is granted full autonomous authorization to modify files, compile, and run tests.

### **Gate 2: Quality Gate (Programmatic)**

- **Phase Transition**: Implement $\rightarrow$ Review.
- **Verification File**: `review-approval.json`.
- **Protocol**: This gate is programmatically validated by the enforcer hooks. It requires that:
  1. The automated unit and linter checks pass successfully, storing the tested diff_hash in `phase-state.json`.
  2. A proactive code review is executed by the isolated, sandboxed **Review Subagent** (`review_agent`), which writes `review-approval.json` upon reporting a clean review.
- **Enforcement**: If any workspace files are modified after Gate 2 is signed, the enforcer hooks automatically delete the signatures, revoking approval and requiring re-testing and re-review.

### **Gate 3: Commit Gate (User-Facing)**

- **Phase Transition**: Review $\rightarrow$ Commit.
- **Verification File**: `user-approval.json`.
- **Protocol**: The agent presents the active unstaged Git diff for visual IDE review. The developer explicitly approves the proposed Conventional Commit message via `ask_user` (format: `Commit Message: "docs: <description>"`).
- **Execution**: The hook (`04-commit-phase.js --after-ask`) intercepts the approval and triggers SSH Key signing (via `ssh-keygen -Y sign`) to write `user-approval.json` and its signature. It then securely stages the files, commits, pushes the branch, and programmatically generates a Draft Pull Request on GitHub.

---

## Core Mandates

1. **Zero Data Loss Guarantee:** Destructive Git commands (`git reset --hard`, `git checkout .`, `git clean -fd`) must never be run on uncommitted workspace files unless explicitly requested by the developer, or after backing up work to a temporary branch or the standard backup directory (`~/.gemini/tmp/<repo-name>/backup_changes`).
2. **IDE Review Priority:** Developers review code changes directly in their IDE while they are **unstaged** to maintain color-coded diff visibility. No commits can occur without presenting the unstaged diff and receiving explicit SSH Key-signed commit approval via the `ask_user` tool.
3. **No Upstream Pushes:** All remote pushes must target the developer's fork, never the upstream "rancher" remote.
4. **Strict Release-Please & SemVer Rules:** All draft commit messages must strictly adhere to Conventional Commits from the end-user product's perspective:
   - **`feat`** (bumping SemVer Minor) and **`refactor`/`!`** (bumping SemVer Major) are strictly reserved for changes directly modifying the Terraform definition files (`main.tf`, `variables.tf`, `versions.tf`, or `outputs.tf`).
   - **Internal Dev Changes:** Changes to helper scripts, CI/CD configuration, linters, internal hooks, or test suites must use non-bumping conventional prefixes such as `build`, `ci`, `test`, `docs`, `fix`, or `chore`.

---

## Step-by-Step Procedure

### Phase 1: Plan Phase (Gate 1)

1. **Research & Explore**: Map the goal and hurdle. Search the codebase for existing patterns and affected source or test files.
2. **Empirical Bug Reproduction**: For bug fixes, write a reproduction script or local test that demonstrates the failure, and run it to confirm the bug state.
3. **Draft Plan**: Draft a step-by-step imperative plan checklist under `plans/` in the session workspace. Explain the implementation details and testing strategy.
4. **Solicit Plan Approval (Gate 1)**: Present the plan in the chat and request cryptographic approval via `ask_user` (using a choice option labeled "Approve Plan"). The backend hooks sign the plan using the developer's SSH key material and write `plan-approval.json` along with its signature file.
5. **Phase Transition**: Call `exit_plan_mode` to transition the session from Plan to Implement.

### Phase 2: Implement Phase (Autonomous)

1. **Surgical Refactoring**: Sequentially implement the tasks from the approved plan, updating the checkboxes in the plan file. Keep edits focused and surgical.
2. **Verification Tests**: Compile and execute tests locally to verify correctness.
3. **Linter & Static Analysis Compliance**: Run `make lint` to verify compliance with repository style guides and static analysis rules.

### Phase 3: Review Phase (Gate 2)

1. **Testing Sign-Off**: Run `make test` to verify full codebase integration, which stores the tested diff_hash in `phase-state.json`.
2. **Delegate Proactive Review**: Delegate a proactive code review of the active local Git diff to the sandboxed `review_agent` using `invoke_agent`.
3. **Resolve Findings**: The review agent's primary goal is to be a critical, adversarial peer reviewer. If the subagent flags any architectural gaps or documentation inconsistencies under the `Findings & Comments` section of the report, surgically resolve them and re-run the review until all 4 passes are checked (`- [x]`) and exactly `0 comments/findings` are reported, which allows the enforcer hook to programmatically sign and write `review-approval.json`.

### Phase 4: Commit Phase (Gate 3)

1. **Isolate Changes**: Create a dedicated feature branch off the updated `main`. Keep the changes unstaged in the working directory.
2. **Solicit Commit Approval (Gate 3)**: Present the unstaged diff and request final commit approval via `ask_user` using the proposed conventional commit message format (`Commit Message: "docs: <description>"`).
3. **Automated Commit & Push**: The developer's SSH key material signs `user-approval.json` via `ssh-keygen -Y sign`. The enforcer hook intercepts the approval, stages the changes, commits, pushes the branch, and programmatically opens a Draft Pull Request on GitHub.
4. **Graduation & Conclude**: Convert the Draft PR to "Ready for Review" via `gh pr ready <pr-number>` once finalized, and cleanly close the session.

---

## Pull Request Iteration & Comment Resolution Protocol

When resolving comments or feedback on an open Pull Request, developers and agents must adhere to the following systematic quality iteration loop:

1. **Update Review Agent Guidelines First**: Translate the review comments into strict checking rules and append them to `.gemini/agents/review_agent.md` under Section 3 ("Strict Quality Gates, Refactoring, & Safety Verification").
2. **Run Review Agent on Unmodified Codebase**: Run the `review_agent` on the unmodified codebase files. The subagent must successfully reproduce the findings in its report and conclude with `PR Review status: 🔴 FINDINGS - Violations detected.`.
3. **Implement Surgical Code Fixes**: Only after the subagent successfully reproduces the findings in its local report are you authorized to modify files to address the comments.
4. **Local Verification & Final Review**: Run local linters and tests. Then, execute the `review_agent` one final time to confirm it approves the workspace changes with a clean `PR Review status: 🟢 PERFECT - 0 findings`.
5. **Push Updates & Resolve**: Stage the changes, commit using the Commit Gate, and push to the remote branch. Execute `.gemini/skills/resolve-pr-reviews.sh` to programmatically resolve the comment threads on GitHub.
