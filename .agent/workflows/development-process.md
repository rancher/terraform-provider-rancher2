# Workflow: Standard Development Process

This is the standard development process. All AI agents MUST strictly read, understand, and follow this unified process for any session, feature implementation, refactoring, maintenance, or bug fix.

---

## 🔒 The Three Authoritative Approval Gates

To ensure high-signal coordination and eliminate redundant or disjointed prompts, this workflow is strictly consolidated around three mandatory **Approval Gates**. Outside of these gates, the agent is granted full autonomous authorization to execute.

### **Gate 1: Planning Gate (Initial Strategy Approval)**
* **Location**: Phase 2 (Blueprint & Planning) - Step 5.
* **Protocol**: The agent MUST NOT modify any source files or run any mutating development commands before presenting the planning blueprint and implementation checklist (inside `.agent/plans/AgenticFramework.md` or other plan file) in the chat and receiving explicit developer approval.
* **Autonomous Phase**: Once Gate 1 is approved, the agent operates with **full autonomous authorization** through Phase 3 (Surgical Implementation) and Phase 4 (Proactive Quality Gate). The agent does *not* need to stop and ask for "interim" permissions to compile, run tests, lint, or invoke the review agent.

### **Gate 2: IDE & Commit Gate (Implementation Approval)**
* **Location**: Phase 5 (IDE Review & Secure Commit-Push) - Steps 13 & 14.
* **Protocol**: The agent presents the active unstaged git diff in the chat for the developer's visual IDE review and proposes a conventional commit message. Once agreed, the agent stages the files and executes the custom secure commit-and-push skill in the chat:
  ```bash
  .agent/skills/commit-push.sh -m "your commit message"
  ```
  This script acts as the definitive Gate 2 execution engine. It strictly requires the developer to use `@user-approval` (or the internal `.agent/skills/user-approval.js`) to cryptographically sign the active diff hash *prior* to execution. If a valid, current developer signature is found, the skill verifies it cleanly, synchronizes with upstream, and autonomously signs/commits/pushes the changes with **zero interactive prompts**.
  *(Fallback: If helper skills are not present in your active branch, stage changes and perform standard Conventional Commits with a signed and signed-off commit and push: `git commit -s -S -m "fix(scope): description" && git push origin <branch>` after receiving explicit approval in the chat).*

### **Gate 3: Draft PR Review Gate (Ready-for-Review Approval)**
* **Location**: Phase 6 (Draft PR & Ready Conversion) - Steps 16 & 17.
* **Protocol**: The agent programmatically generates a Draft Pull Request on GitHub using `.agent/skills/create-pr.sh --draft` and halts execution.
  *(Fallback: If helper skills are not present, create the draft PR natively using GitHub CLI: `gh pr create --draft --repo rancher/terraform-provider-rancher2`).*
  The developer inspects the draft PR on GitHub. Upon receiving the developer's explicit approval in the chat, the agent converts the PR to "ready for review" (`gh pr ready <pr-number>`), presents the standard PR link, and cleanly **closes the development session**.
  * **Asynchronous Review Cycle**: The developer waits asynchronously for team and AI reviews. If changes or comments are received on GitHub, the developer starts a **brand new, separate development session** executing the specialized `resolve-pr-reviews.md` workflow to resolve comments and merge.

---

## Core Mandates

1. **Zero Data Loss Guarantee:** Never run destructive git commands (such as `git reset --hard`, `git checkout .`, or `git clean -fd`) on modified workspace files unless explicitly requested by the developer, or after backing up work to a temporary branch/stash or the standard backup folder.
2. **IDE Review Priority:** The developer prefers to review code changes directly in their IDE while they are **unstaged** in the Git working directory to maintain color-coded diff visibility. Never execute a `git commit` without presenting the exact unstaged diff and receiving explicit approval in the chat.
3. **No Upstream Pushes:** Never push directly to upstream "rancher" remotes. All remote operations must target the user's fork.
4. **Strict Release-Please & SemVer Rules (Product-Centric):** All draft commit messages must strictly adhere to Release-Please rules from the end-user product's perspective:
  - **`feat`** or **`feature`** (bumping SemVer Minor) and **`refactor`** / **`!`** (bumping SemVer Major) MUST ONLY be used if the change directly modifies the Go source files or schemas defining the published provider itself (e.g. inside `rancher2/` folder or `main.go`, `go.mod`, etc.).
  - **Internal Dev Changes:** Changes to helper scripts, CI/CD configuration, linters, internal hooks, or test suites DO NOT affect the published product. They MUST NOT use `feat`, `feature`, `refactor`, or `!` types. Instead, use non-bumping conventional prefixes such as `fix:` or `chore(main):` to satisfy the repository's commit-message validation constraints (which block prefixes like `build`, `ci`, `test`, `docs`, or plain `chore`).
5. **Secure Local Backup & Isolation (~/.gemini/tmp):** To isolate staged commits for pristine IDE review with zero clutter, the agent MUST temporarily backup all non-layer modified and untracked files to the standard `~/.gemini/tmp/<repo-name>/backup_changes` directory.

---

## Step-by-Step Procedure

### Phase 1: Research & Reproduce
1. **Understand Goal & Hurdle:** Map the user's high-level goal and hurdle. If an existing workflow matches (e.g. CI failure matches `troubleshoot-workflows.md`), declare it explicitly.
2. **Codebase Exploration:** Search the codebase for existing patterns, conventions, and affected source/test files.
3. **Empirical Bug Reproduction:** For bug fixes, write a reproduction script or local test that demonstrates the failure, and run it to confirm the bug state.

### Phase 2: Planning, Strategy & Blueprint Synchronization (Gate 1)
4. **Acquire, Edit, or Create Architectural Plan:** Following `.agent/rules/plans.instructions.md`, verify if an existing plan covers the target domain.
   - **If an existing plan covers the domain:** You MUST NOT create a new plan file. Instead, *edit* and *adapt* the existing plan, modifying its top-half blueprint/specification and expanding/re-opening the bottom-half implementation checklist.
   - **If no existing plan matches the domain:** Create a brand new unified plan file in `.agent/plans/<PlanName>.md`.
   - **Checklist Construction:** Build and expand the sequential implementation checklist to detail the specific sub-tasks. You MUST explicitly incorporate all standard quality gates (local build/test verification, static linters, proactive review, upstream sync, unstaged IDE review, authorized commit, and draft PR generation) as checkbox-tracked items.
5. **🔒 Solicit Plan Approval (Gate 1):** Present the updated blueprint and implementation checklist to the developer in the chat for explicit approval. **The agent is strictly prohibited from modifying any source files or running mutating development commands until Gate 1 approval is received.**

### Phase 3: Surgical Implementation (Autonomous Action)
6. **Execute Plan & Track State (No Stage/Commit):** Implement the plan sequentially, updating checkboxes in the plan file in place. Keep edits simple, precise, and idiomatic. Do NOT stage (`git add`) or commit (`git commit`).
7. **Build & Test Verification:** Compile, build, and run tests locally.
   * **Full Test Suite Context Warning:** The full test suite can take over an hour and generate massive logs. Redirect output (e.g., `./run_tests.sh [options] > /tmp/run_tests.log 2>&1`) and run `.agent/skills/parse-test-logs.sh` to prevent context window flooding.
   * **Fast Verification Option:** Validate changes quickly on a single fixture:
     ```bash
     ./run_tests.sh -f sle-micro-61-canal-stable-one-rpm-ipv4
     ```
8. **Static Analysis & Linters:** Run ecosystem linters (e.g., `golangci-lint`, `shellcheck`, `tflint --recursive`, `actionlint`) and resolve all warnings.

### Phase 4: Proactive Quality Gate (Autonomous Action)
9. **Proactive Code Review:** Delegate a proactive code review of your active local git diff directly to the custom review subagent by running `@review_agent` in the chat. The agent will rigorously verify your modifications against `.agent/rules/github-copilot-review.instructions.md` and all repository standards, generating a pre-commit review report and the secure cryptographic SHA-256 approval signature.
10. **Resolve Findings:** Refactor and fix any concerns discovered by the review agent, ensuring exactly 0 automated Copilot or linter findings.

### Phase 5: Chunking & IDE Review (Gate 2)
11. **Logical Partitioning:** If there is a large volume of changes, group files into focused, independent **subsystem boundaries** (layers).
12. **Upstream Synchronization:** Before checkout, switch to `main` and execute `.agent/skills/git-sync.sh` to ensure our branch off point is completely up-to-date with upstream.
13. **Isolate First Layer (Keep Unstaged):** Create a dedicated branch directly off the updated `main`. To keep the workspace clean, backup all other non-layer files to the standard `~/.gemini/tmp/<repo-name>/backup_changes` directory. Clean other files from the working directory, leaving **exclusively** the target layer's changes unstaged.
14. **🔒 Solicit IDE & Commit Approval (Gate 2):** Present the unstaged diff to the developer in the chat for visual IDE review. Propose and agree on a conventional commit message.
15. **🔒 Execute Secure Commit & Push**: Once Gate 2 approval is confirmed in the chat, stage your changes and run the secure skill `.agent/skills/commit-push.sh -m "commit message"`. **Direct manual git commit and push commands are strictly prohibited.** The developer's confirmation inside the interactive TTY prompt of the script concludes Gate 2.

### Phase 6: Draft PR & Ready Conversion (Gate 3)
16. **Generate Draft PR:** Create a Draft Pull Request using `.agent/skills/create-pr.sh --draft`.
17. **🔒 Solicit Draft PR Approval (Gate 3):** Present the draft PR link in the chat and halt execution. Wait for the developer to inspect the draft PR on GitHub.
18. **Convert PR to Ready & Conclude Session:** Once the developer explicitly approves the draft PR in the chat, convert the PR from Draft to "Ready for Review" using the GitHub CLI: `gh pr ready <pr-number>`. Provide a completion summary, present the final PR link, and cleanly **close the current development session**.

### Phase 7: Asynchronous PR Iteration & Next Layer Restoration
19. **Asynchronous Review Wait-State:** The developer waits asynchronously for team and AI reviews. If comments or requested changes are submitted on GitHub, the developer starts a **brand new development session** running the dedicated `.agent/workflows/resolve-pr-reviews.md` workflow to resolve comments.
20. **Proceed to Next Layer:** Switch back to the synchronized `main`, restore the remaining files from the backup directory `~/.gemini/tmp/<repo-name>/backup_changes` back into the active workspace, and return to Step 11 for the next layer.
21. **Completion Summary:** Once all layers are successfully complete and merged, provide a concise summary with links to all Pull Requests.
