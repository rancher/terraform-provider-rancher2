# Agentic Framework: Workflow Optimization & Design

## Abstract

To maximize efficiency and eliminate friction in a human-agent collaborative environment, the Agentic Framework streamlines processes, eliminates mechanical checkpoints from human loops, and structures developer interactions around three authoritative checkpoints. Simultaneously, it prunes redundant styling instructions from agent profiles to reduce context overhead and processing costs.

---

## Technical Specification

### 1. The Gated 4-Phase Lifecycle & Three Gates

To optimize collaboration and ensure zero unvetted changes, the framework coordinates work across 4 distinct phases (`Plan`, `Implement`, `Review`, `Commit`) and enforces three strict, sequential gates:

```text
       [ Plan Phase ]
             │
             ▼
 🔒 Gate 1: Planning Gate (User-Facing / GPG Touch ID)
             │
             ▼
       [ Implement Phase (Autonomous) ]
             │
             ▼
       [ Review Phase ]
             │
             ▼
 🔒 Gate 2: Quality Gate (Programmatic / Test & Subagent)
             │
             ▼
       [ Commit Phase ]
             │
             ▼
 🔒 Gate 3: Commit Gate (User-Facing / GPG Touch ID)
```

1. **Planning Gate (Gate 1 - User-Facing)**:
   - **Phase Transition**: Plan $\rightarrow$ Implement.
   - **Security**: Prompts macOS Touch ID to GPG-sign the active strategy checklist on disk, writing `plan-approval.json`.
   - **Authorization**: Unlocks autonomous file modification, compilation, and testing capabilities.
2. **Quality Gate (Gate 2 - Programmatic)**:
   - **Phase Transition**: Implement $\rightarrow$ Review.
   - **Security**: Natively verifies that local unit/integration tests pass successfully and delegates an automated code review to our sandboxed `review_agent` to secure the review signature (`review-approval.json`).
3. **Commit Gate (Gate 3 - User-Facing)**:
   - **Phase Transition**: Review $\rightarrow$ Commit.
   - **Security**: Displays the live unstaged Git diff in chat, requesting Conventional Commit message approval. It triggers macOS Touch ID to verify the developer's physical sign-off and write `user-approval.json`.
   - **Automation**: Upon biometric verification, the hook automatically stages files, commits with the signature, pushes, and programmatically opens a Draft PR on GitHub.

---

## Asynchronous Review Iterations

To prevent clogging active workspace contexts with long-lived PR review wait states:

- Once a PR is opened, the active session is cleanly **closed**.
- If external maintainers or automated reviewers leave requested changes on GitHub, the developer starts a **brand new development session** running a dedicated `.gemini/workflows/resolve-pr-reviews.md` workflow.
- In accordance with our PR iteration standards, comments are resolved by updating the review agent's rules first, reproducing findings, implementing fixes, re-verifying Gate 2, and committing via Gate 3.
- This keeps individual sessions extremely short-lived, fast, and completely free of state contamination.

---

## ✂️ Prompt Pruning & Tooling Synergies

Mechanical linter validation (such as scanning for trailing whitespace, checking bracket indentation, or verifying formatting) is highly repetitive and computationally expensive to delegate to LLM reasoning.

By implementing strict, deterministic, and hermetic formatting tools (Prettier, shfmt, gofmt) in our local environment:

- All formatting enforcements are offloaded to local compiler binaries.
- We **prune all mechanical style rules** from the instructions of our AI agents (e.g., `review_agent.md`).
- This dramatically reduces prompt sizes, minimizing context window footprint and cloud-processing API costs.
- The `review_agent` can focus 100% of its cognitive window on high-signal architectural logic, security vectors, and structural compliance.
