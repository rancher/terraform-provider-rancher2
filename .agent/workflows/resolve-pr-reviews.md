# Workflow: PR Review Comment Resolution

This workflow defines a repeatable, high-standard engineering process for analyzing, planning, and executing resolutions for Pull Request review comments (such as from Copilot or automated review bots).

## Purpose
Automated review bots and scanners are excellent at detecting edge cases, safety risks, and logic flaws, but their actual code recommendations are often sub-optimal, non-idiomatic, or overly verbose. 

This workflow enforces a **"Discernment-First"** protocol: AI agents must analyze and validate the *underlying concerns* raised by reviewers, design custom, idiomatic, high-quality solutions, document them in active project plans, and execute them iteratively under user supervision.

---

## Detailed Step-by-Step Procedure

### 1. Retrieve PR Review Comments
Before modifying any files or assuming a solution, gather the latest feedback directly from the PR context:
* Execute the repository skill script to fetch a unified timeline of all top-level and inline reviews:
  ```bash
  .agent/skills/get-pr-comments.sh
  ```
* Review the output timeline, taking note of the file paths, lines, and the specific feedback left by reviewers.

### 2. Identify the Active Project Plans
Locate the project plans associated with the current branch/feature development to keep documentation in sync:
* **Persistent Plan:** Located under `.agent/plans/<feature-name>.md`
* **Temporary Plan:** Located under `.agent/agent-memory/<feature-name>-temporary.md`

### 3. Separation of Concerns (Discernment Phase)
For each comment retrieved, perform a strict architectural assessment:
* **The Concern (Valid):** Is there an actual logic flaw, unhandled edge case, syntax warning, or security vulnerability being pointed out? If yes, this concern is valid and **must** be resolved.
* **The Suggestion (Ignore/Improve):** Analyze the actual code block or command suggested by the bot. **Do not blindly adopt it.** Reject non-idiomatic workarounds, hacks that disable warnings, or "soft-failure" defaults that mask configuration errors.
* **The Idiomatic Design:** Formulate a clean, native, and high-quality solution that addresses the *concern* perfectly while aligning with the workspace standards and language paradigms (e.g., using Bash's native `printf %q` instead of custom string escape hacks; enforcing strict contracts and failing fast instead of introducing silent defaults).

### 4. Adjust the Feature Plans
Before implementing any changes, record the planned actions in the feature plans to ensure traceability:
* **Update the Persistent Plan (`.agent/plans/<feature-name>.md`):**
  * Add a dedicated section (e.g., `Address Review Comments` or `Step 7: Address Copilot Review Comments`).
  * Explicitly document the core *concerns* and outline the designed *idiomatic solutions* (contrasting them against the sub-optimal bot recommendations if necessary).
* **Update the Temporary Plan (`.agent/agent-memory/<feature-name>-temporary.md`):**
  * Append a new task section to the checklist tracking each of the designed resolutions.
  * Keep tasks unchecked `[ ]` to guide the execution phase.

### 5. Obtain User Approval
Present your formulated strategy and plan updates to the user:
* Explain the underlying concerns clearly.
* Provide the technical rationale behind your custom, idiomatic solutions.
* Wait for the user's feedback or explicit directive to proceed before making any code modifications.

### 6. Execute and Validate Iteratively
Once authorized by the user, proceed to implement and verify the changes:
* **Act surgically:** Apply changes step-by-step, touching only the necessary files.
* **Validate thoroughly:** Run local linters, compilers, and test suites (e.g., `shellcheck`, `actionlint`, local unit/acceptance tests) to verify correctness.
* **Track progress:** Update the temporary plan checklist as each task is successfully completed and verified.
