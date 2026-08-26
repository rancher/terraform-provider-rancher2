# Review Agent Component Specification

This component specification details the persona, checking checklists, security boundaries, and validation criteria governing our automated pre-commit Review Subagent (`review_agent`).

---

## Abstract

The Review Subagent is a hardened, read-only AI agent tasked with performing exhaustive pre-commit code audits on our active workspace diff. By running line-by-line checks against our coding standards and security safeguards, the agent guarantees exactly zero styling, formatting, or security violations, programmatically writing the `review-approval.json` signature to unblock GPG commits.

---

## 🛠️ Architectural Design & Constraints

### 1. Hardened Subagent Sandbox

To guarantee zero-trust execution and prevent untrusted code changes from executing unsafe operations during reviews:

- The Review Subagent's custom model definition is located at `.gemini/agents/review_agent.md`.
- It operates under a strictly sandboxed permission model: it does NOT have access to file-modifying tools like `replace` or `write_file`, nor shell-execution tools like `run_shell_command`. It is purely read-only, interacting with the workspace exclusively via `read_file` and `glob`.

### 2. Line-by-Line Checking Protocols

The agent evaluates our active `git diff HEAD` against our language-specific instructions located under `docs/development/CodingStandards/`:

- **Go Components**: Checked against **[Go Standards](../CodingStandards/Go.md)** (checks for immediate error wrapping, early returns, context propagation, and banned panics).
- **Terraform Components**: Checked against **[Terraform Standards](../CodingStandards/Terraform.md)** (checks variable validation and check/precondition blocks).
- **Workflows & Scripts**: Checked against **[Workflow Standards](../CodingStandards/Workflows.md)**, **[GitHub Script Standards](../CodingStandards/GitHubScript.md)**, and **[Shell Script Standards](../CodingStandards/ShellScripts.md)**.
- **Documentation**: Checked against **[Documentation Standards](../CodingStandards/Documentation.md)**.

### 3. Gating & Chaining Integration

The Review Subagent is executed during Phase 4 of our **Development Process**:

- It requires Gate 1 (Planning Gate) and Gate 2 (Testing Gate) to be successfully signed and valid before execution is authorized by `.gemini/hooks/03-review-phase.js`.
- Upon successful execution with zero findings, the system's `03-review-phase.js --after-invoke` hook automatically intercepts, cryptographically signs, and writes `review-approval.json` to disk, chaining it to our active plan and diff hashes.
