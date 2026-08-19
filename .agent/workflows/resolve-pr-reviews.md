# Workflow: PR Review Comment Resolution

This workflow defines a repeatable, high-standard engineering process for analyzing, planning, and executing resolutions for Pull Request review comments in a dedicated, fresh development session.

## Purpose
Review comments from human maintainers and automated bots provide valuable insights. This workflow enforces a **"Discernment-First"** protocol: AI agents must analyze the *underlying concerns*, critically evaluate if they are valid, design high-quality, standard-compliant solutions, post clear responses to each thread on GitHub, and programmatically resolve them after refactoring.

---

## Detailed Step-by-Step Procedure

### 1. Initiate Session & Retrieve Comments
This workflow is initiated in a **brand new development session** specifically started to resolve PR comments.
* First, retrieve a chronological timeline of all general and inline review comments:
  ```bash
  .agent/skills/get-pr-comments.sh [PR_ID]
  ```
* Analyze the timeline to map file paths, lines, authors, and feedback.

### 2. Separation of Concerns & Evaluation (Discernment Phase)
For each comment retrieved, perform a critical architectural assessment:
* **Evaluate Validity**: Is there an actual logic flaw, security vulnerability, syntax error, or style deviation? 
  - **If Valid**: Acknowledge the concern and design a custom, idiomatic fix conforming to `.agent/rules/`.
  - **If Invalid**: Prepare a clear, professional, and technical explanation why the current implementation is correct.
* **Reject Bot Hacks**: Never blindly copy sub-optimal recommendations, hacks that disable warnings, or "soft-failure" defaults that mask configuration errors.

### 3. Update active plan
Before making any file edits, adapt the active plan file under `.agent/plans/` to document the comments:
* Add a dedicated comment resolution section mapping out the evaluated concerns and custom solutions.
* Append corresponding task checkboxes to the plan's checklist.

### 4. Respond to Comment Threads on GitHub
To maintain high collaboration standards, post an explicit response to each comment thread on GitHub before or during the fix:
* Explain your technical evaluation and solution, or provide your counter-rationale if the concern is invalid.
* *(Note: Response comments can be posted using native GitHub CLI or discussion APIs).*

### 5. Surgical Refactoring & Verification
Implement and verify changes autonomously off the approved plan:
* **Act surgically**: Touch only the necessary files.
* **Validate thoroughly**: Run local linters (`eslint`, `shellcheck`), compilers, and test suites (`make test`) to ensure exactly 0 findings and 0 regressions.

### 6. Secure Commit & Push
Once verified, commit and push your changes using our custom secure commit-push skill:
```bash
.agent/skills/commit-push.sh -m "fix(hooks): resolve review findings on PR #<id>"
```
*(Fallback: If helper skills are not present in your active branch, perform standard Conventional Commits with a signed and signed-off commit and push: `git commit -s -S -m "fix(scope): description" && git push origin <branch>`).*

### 7. Programmatic Thread Resolution
Once changes are pushed and verified on GitHub, programmatically resolve all comment threads on GitHub:
```bash
.agent/skills/resolve-pr-reviews.sh [PR_ID] --bypass-token --all
```
*(Fallback: If helper skills are not present, resolve threads manually in the GitHub UI, or use GitHub GraphQL API `resolveReviewThread` mutation).*
Verify that all threads are fully closed on GitHub, concluding the resolution session.
