# PR Review Comment Resolution

---

## Abstract

This component defines our high-standard engineering process for analyzing, planning, and executing resolutions for Pull Request review comments in a dedicated development session.

---

## Purpose

Review comments from human maintainers and automated bots provide valuable insights. This workflow enforces a **"Discernment-First"** protocol: AI agents must analyze the _underlying concerns_, critically evaluate if they are valid, design high-quality, standard-compliant solutions, post clear responses to each thread on GitHub, and programmatically resolve them after refactoring.

---

## Detailed Step-by-Step Procedure

### 1. Retrieve Comments

- First, retrieve a chronological timeline of all general and inline review comments:

  ```bash
  .gemini/skills/get-pr-comments.sh [PR_ID]
  ```

- Analyze the timeline to map file paths, lines, authors, and feedback.

### 2. Separation of Concerns & Evaluation (Discernment Phase)

For each comment retrieved, perform a critical architectural assessment:

- **Evaluate Validity**: Is there an actual logic flaw, security vulnerability, syntax error, or style deviation?
  - **If Valid**: Acknowledge the concern and design a custom, idiomatic fix conforming to `docs/development/CodingStandards/`.
  - **If Invalid**: Prepare a clear, professional, and technical explanation why the current implementation is correct.
- **Reject Bot Hacks**: Never blindly copy sub-optimal recommendations, hacks that disable warnings, or "soft-failure" defaults that mask configuration errors.

### 3. Update active plan

Before making any file edits, adapt the active plan under `plans/` in the session workspace to document the comments:

- Add a dedicated comment resolution section mapping out the evaluated concerns and custom solutions.
- Append corresponding task checkboxes to the plan's checklist.

### 4. Respond to Comment Threads on GitHub

To maintain high collaboration standards, post an explicit response to each comment thread on GitHub before or during the fix:

- Explain your technical evaluation and solution, or provide your counter-rational if the concern is invalid.
- _(Note: Response comments can be posted using native GitHub CLI or discussion APIs)._

### 5. Surgical Refactoring & Verification

Implement and verify changes autonomously off the approved plan:

- **Act surgically**: Touch only the necessary files.
- **Validate thoroughly**: Run local linters (`eslint`, `shellcheck`), compilers, and test suites (`make test`) to ensure exactly 0 findings and 0 regressions.

### 6. Secure Commit & Push

Once verified, present the unstaged diff and request approval via `ask_user` (format: `Commit Message: "fix(hooks): resolve review findings on PR #<id>"`). The hook will automatically write the signature, commit, push, and update the PR on GitHub.

### 7. Programmatic Thread Resolution

Once changes are pushed and verified on GitHub, programmatically resolve all comment threads on GitHub:

```bash
.gemini/skills/resolve-pr-reviews.sh [PR_ID] --bypass-token --all
```

Verify that all threads are fully closed on GitHub, concluding the resolution session.
