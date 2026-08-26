# Testing & Quality Verification

---

## Abstract

This topic details the linting, formatting, and static analysis configurations, local execution commands, and pipeline integrations that enforce the repository's code quality and execution standards.

---

## Technical Specifications

The codebase maintains a strict, zero-finding policy for all code, shell scripts, formatting, and markdown assets. This is verified locally and in CI/CD pipelines through native, multi-tiered engines.

### 1. Unified Linter and Formatter Configurations

- **Nix Tooling**: All linting and formatting dependencies are loaded natively through `flake.nix`.
- **Prettier (`.prettierrc`)**: Enforces uniform trailing commas, quote marks, and tab-widths across JS, JSON, and MD files.
- **MarkdownLint (`.markdownlint.yaml`)**: Validates semantic structure while intentionally ignoring inline HTML or line-length limits required by agent prompt formatting.
- **GoLint (`.golangci.yml`)**: Enables `gosec` (security), `errcheck` (safety), `revive`/`stylecheck` (idioms), and `gocyclo` (complexity).

### 2. Go Unit & Integration Testing (Future Expansion)

_(Placeholder for future standardizations around `make test`, coverage thresholds, and parallelization)._

### 3. Acceptance & End-to-End Testing (Future Expansion)

_(Placeholder for future standardizations around Terraform Acceptance Testing `make testacc`)._
