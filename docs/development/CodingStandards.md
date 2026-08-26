# Coding Standards & Guidelines

This topic overview establishes and indexes the mandatory coding standards, security baselines, and quality criteria enforced across all repository languages, configuration formats, and automation components.

---

## Abstract

Coding standards are a central pillar of our engineering quality. By establishing strict, clear, and language-specific instructions for every file type in our project, we ensure that both human developers and autonomous agents write highly readable, secure, and idiomatic code. These guidelines are coupled with automated local linters, secret scanning, and automated pre-commit review gates to guarantee 100% compliant codebase changes.

---

## 🧭 How Our Coding Standards Work Together

Our coding standards are organized into language-specific and tool-specific component specifications. While they address different ecosystems, they work together as a cohesive compliance framework:

### 1. Core Development Languages

- **Go Standards**: Enforce clean Go styles, explicit error handling/wrapping, context propagation, and network timeouts. More details can be found in **[Go Standards](./CodingStandards/Go.md)**.
- **Terraform Standards**: Prescribe resource naming conventions, variable validations, and pre/postconditions. More details can be found in **[Terraform Standards](./CodingStandards/Terraform.md)**.

### 2. CI/CD & Automation Ecosystems

- **GitHub Actions Workflows**: Enforce least-privilege permission blocks, execution timeouts, and action pinning by 40-character SHA hash. More details can be found in **[GitHub Actions Workflows](./CodingStandards/Workflows.md)**.
- **GitHub Script Rules**: Enforce standards for JavaScript-based GHA automation, including octokit pagination and exception handling. More details can be found in **[GitHub Scripts](./CodingStandards/GitHubScript.md)**.
- **Shell Script Rules**: Enforce strict Bash rules, such as `set -euo pipefail` fail-fasts, double brackets `[[ ]]`, and main execution blocks. More details can be found in **[Shell Scripts](./CodingStandards/ShellScripts.md)**.

### 3. Documentation & Architectural Blueprints

- **Blueprints & Planning**: Dictates the lifecycle of drafting, expanding, and executing checklists in Topic Overviews and Component Specifications under `docs/development/`. More details can be found in **[Architectural Blueprints](./CodingStandards/Blueprints.md)**.
- **Documentation Styles**: Enforces Diátaxis-compliant technical writing, topics, components, and layout structures. More details can be found in **[Documentation Standards](./CodingStandards/Documentation.md)**.

### 4. Code Review Criteria

- **Pre-Commit Review Guidelines**: Establishes the exact validation checklist executed by our automated review subagent (`review_agent`) during our proactive quality gate. More details can be found in **[GitHub Copilot Review Guidelines](./CodingStandards/GitHubCopilotReview.md)**.

---

## 🔒 Global Baselines & Security Standards

In addition to language-specific specs, every file in the repository must strictly comply with three global quality and security baselines, which are verified locally and in CI/CD:

### A. Code Formatting (Prettier & shfmt)

All files (Markdown, JSON, JavaScript, Bash) must be formatted automatically before commit:

- Prettier is used for general layout formatting (checked via `prettier --check .` and fixed via `prettier --write .`).
- `shfmt` is used for shell script formatting (fixed via `shfmt -w .`).

### B. Secret Scanning (Gitleaks)

To prevent accidental exposure of sensitive keys, API tokens, GPG credentials, or certificates:

- Every single commit, staged difference, and active branch history is scanned via `gitleaks` (executed via `gitleaks detect` locally).
- No hardcoded secrets are allowed in any file, including `.envrc` or configuration settings.

### C. Spelling Validation (CSpell)

All source code, documentation, comments, and commit messages must conform to standard technical English:

- Spellchecking is automatically run using `cspell` (configured via `cspell.json` and customized via `custom_words.txt`).
- No typo or unrecognized abbreviations are allowed in committed files.
