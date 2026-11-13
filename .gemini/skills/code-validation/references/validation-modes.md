# Code Validation Modes & Reference Guide

This document describes all static analysis modes, test suites, formatting tools, and troubleshooting patterns available within the repository.

## Static Analysis Modes (`lint.js`)

The static analysis suite wraps `.github/workflows/scripts/lint.sh` and supports the following modes:

| Mode            | Tooling                   | Description                                                    | Fix Supported |
| :-------------- | :------------------------ | :------------------------------------------------------------- | :------------ |
| `all`           | All tools                 | Executes all static analysis suites sequentially               | Yes (`--fix`) |
| `terraform`     | `terraform fmt`, `tflint` | Verifies Terraform configuration formatting and best practices | No            |
| `actionlint`    | `actionlint`              | Validates GitHub Actions workflow YAML files                   | No            |
| `eslint`        | `eslint`                  | Checks JavaScript and MJS codebase files and catch statements  | Yes (`--fix`) |
| `shellcheck`    | `shellcheck`              | Lints POSIX and Bash shell scripts for syntax and safety       | No            |
| `shfmt`         | `shfmt`                   | Validates shell script indentation and layout                  | Yes (`--fix`) |
| `prettier`      | `prettier`                | Formats Markdown, JSON, and JavaScript files                   | Yes (`--fix`) |
| `markdownlint`  | `markdownlint`            | Enforces Markdown formatting rules and link structure          | No            |
| `golangci-lint` | `golangci-lint`           | Comprehensive Go static analysis (root and `test` modules)     | Yes (`--fix`) |
| `tests`         | `golangci-lint`, `gofmt`  | Runs Go test-specific legacy checks in `test/` directory       | No            |
| `gitleaks`      | `gitleaks`                | Scans workspace and git history for hardcoded credentials      | No            |
| `cspell`        | `cspell`                  | Verifies code and documentation spelling against dictionary    | No            |

## Test Suites (`test.js`)

The test suite wraps `.github/workflows/scripts/test.sh` and supports the following modes:

| Mode               | Description                                                                | Command Executed                                                                         |
| :----------------- | :------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
| `unit`             | Runs Go unit tests across root modules (default)                           | `gotestsum` with coverage                                                                |
| `compile`          | Verifies test binaries compile cleanly                                     | `go test -c` in `test/`                                                                  |
| `workflow-scripts` | Runs Node unit tests for workflow automation                               | `node --test .github/workflows/scripts/tests/**/*.test.js`                               |
| `agent-scripts`    | Runs Node unit tests for agents, hooks, and skills                         | `node --test .gemini/orchestrator/lib/**/*.test.js .gemini/skills/**/tests/**/*.test.js` |
| `release-scripts`  | Runs Node unit tests for release automation scripts                        | `node --test scripts/tests/**/*.test.js`                                                 |
| `acc`              | Runs standard acceptance tests (requires credentials)                      | `make testacc`                                                                           |
| `acc-relay`        | Runs AWS Test Relay acceptance tests                                       | `make testaccrelay`                                                                      |
| `all`              | Executes compile check, unit tests, workflow tests, and agent script tests | Full local verification pipeline                                                         |

## Automatic Code Formatters (`format.js`)

Automated code formatting applies style rules across multiple languages:

| Target          | Executable Command                             | Target Files                                   |
| :-------------- | :--------------------------------------------- | :--------------------------------------------- |
| `go` / `gofmt`  | `gofmt -s -w -e .`                             | All `.go` source files                         |
| `prettier`      | `npx prettier --write .`                       | `.js`, `.mjs`, `.json`, `.md`, `.yml`, `.yaml` |
| `shfmt`         | `shfmt -i 2 -ci -bn -w <files>`                | All shell scripts (`*.sh`)                     |
| `golangci-lint` | `golangci-lint run --timeout=5m --fix`         | Go source files (root and `test/`)             |
| `all`           | All above formatters in deterministic sequence | Entire repository                              |

Dry-run verification is supported using the `--check` flag on `format.js`.

## Common Validation Failures & Remediation

1. **Spell Check (`cspell`)**: If valid domain terms, acronyms, or proper names fail spell check, append them to `custom_words.txt` in lowercase.
2. **Markdown Lint (`markdownlint`)**: Ensure clean line breaks, code fence language tags (e.g. ` ```bash `), and single top-level headings.
3. **Go Linter (`golangci-lint`)**: Address unchecked errors, unused variables, and formatting. Auto-fixable issues can be resolved with `node .gemini/skills/code-validation/scripts/format.js golangci-lint`.
4. **Shell Script Formatting (`shfmt`)**: Indentation issues can be resolved with `node .gemini/skills/code-validation/scripts/format.js shfmt`.
