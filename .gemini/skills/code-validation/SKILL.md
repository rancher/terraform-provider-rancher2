---
name: code-validation
description: Execute static analysis, unit/acceptance testing, automated code formatting, and end-to-end validation. Use when validating code changes, running linters (eslint, golangci-lint, cspell, shellcheck, markdownlint), running tests, or formatting code across the repository.
---

# Code Validation and Testing Harness

## Overview

The `code-validation` skill provides standardized, asynchronous Node.js utilities to execute static analysis, test suites, automated formatting, and complete validation pipelines. It encapsulates workspace verification scripts without triggering interactive security confirmation prompts or terminal execution warnings.

**Rule:** You should always prefer using the scripts in this skill over running raw shell commands (`bash lint.sh`, `bash test.sh`, `gofmt`, `prettier`) directly to validate workspace changes.

## Universal Options & Script Discovery

- **Help Option (`-h`, `--help`)**: Every script in this skill supports `-h` and `--help`. Run any script with `--help` to inspect arguments and flags without opening the file:
  ```bash
  node .gemini/skills/code-validation/scripts/<script_name>.js --help
  ```

## Core Tasks

### 1. Static Analysis & Linting

To run the complete static analysis suite across all tools:

```bash
node .gemini/skills/code-validation/scripts/lint.js
```

To run a specific linting suite:

```bash
node .gemini/skills/code-validation/scripts/lint.js [mode]
```

Available modes: `all`, `terraform`, `actionlint`, `eslint`, `shellcheck`, `shfmt`, `prettier`, `markdownlint`, `golangci-lint`, `tests`, `gitleaks`, `cspell`.

To automatically format code where supported during static analysis:

```bash
node .gemini/skills/code-validation/scripts/lint.js --fix
```

### 2. Test Execution

To run unit tests across all Go packages (default):

```bash
node .gemini/skills/code-validation/scripts/test.js
```

To run all local tests (compile checks, unit tests, workflow script tests, agent script tests):

```bash
node .gemini/skills/code-validation/scripts/test.js all
```

To target a specific test suite:

```bash
node .gemini/skills/code-validation/scripts/test.js [mode]
```

Available modes: `unit`, `compile`, `acc`, `acc-relay`, `workflow-scripts`, `agent-scripts`, `all`.

### 3. Automated Code Formatting

To execute all code formatters (`gofmt`, `prettier`, `shfmt`, `golangci-lint --fix`) across the codebase:

```bash
node .gemini/skills/code-validation/scripts/format.js
```

To format a specific language or component:

```bash
node .gemini/skills/code-validation/scripts/format.js [target]
```

Available targets: `all`, `go`, `gofmt`, `prettier`, `shfmt`, `golangci-lint`.

To verify code formatting without writing changes (dry-run):

```bash
node .gemini/skills/code-validation/scripts/format.js --check
```

### 4. Unified Validation Pipeline

To run a full end-to-end verification pipeline (linting followed by unit tests):

```bash
node .gemini/skills/code-validation/scripts/validate.js
```

To auto-format code before running validation suites:

```bash
node .gemini/skills/code-validation/scripts/validate.js --fix
```

To run all static analysis and all test suites:

```bash
node .gemini/skills/code-validation/scripts/validate.js --all
```

For advanced options, custom modes, or skipping steps:

```bash
node .gemini/skills/code-validation/scripts/validate.js --lint-mode eslint --skip-test
```

## Additional Documentation

- For in-depth descriptions of each linting mode, test suite, and troubleshooting guidance, consult [validation-modes.md](references/validation-modes.md).
