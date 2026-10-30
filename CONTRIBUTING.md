# Contributing to the Rancher2 Terraform Provider

Thank you for your interest in contributing to the Rancher2 Terraform Provider (`terraform-provider-rancher2`)!

This guide provides the necessary information for community members, developers, and maintainers looking to report bugs, propose enhancements, write code, run tests, and submit pull requests.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Participate](#how-to-participate)
  - [Reporting Issues & Bugs](#reporting-issues--bugs)
  - [Proposing Features & Enhancements](#proposing-features--enhancements)
  - [Scope Boundaries](#scope-boundaries)
- [Development Environment Setup](#development-environment-setup)
  - [Prerequisites](#prerequisites)
  - [Hermetic Environment with Nix (Recommended)](#hermetic-environment-with-nix-recommended)
  - [Forking and Setting Up Remotes](#forking-and-setting-up-remotes)
- [Branching Strategy & Value Stream](#branching-strategy--value-stream)
    - [Value Flow](#value-flow)
    - [Developer Flow](#developer-flow)
    - [Trunk-Based Development (`main`)](#trunk-based-development-main)
    - [Release Branches & Rancher Version Alignment](#release-branches--rancher-version-alignment)
    - [Automated Backport Orchestration](#automated-backport-orchestration)
- [The Issue-First Requirement](#the-issue-first-requirement)
  - [Linking Issues in PR Descriptions](#linking-issues-in-pr-descriptions)
  - [Target Release Labels & Sub-Issues](#target-release-labels--sub-issues)
- [Building and Local Development](#building-and-local-development)
  - [Building the Provider Binary](#building-the-provider-binary)
  - [Code Generation](#code-generation)
  - [Formatting and Linting](#formatting-and-linting)
- [Testing Guidelines](#testing-guidelines)
  - [Unit Tests](#unit-tests)
  - [Acceptance Tests (End-to-End)](#acceptance-tests-end-to-end)
  - [Cleaning Up Acceptance Test Leftovers](#cleaning-up-acceptance-test-leftovers)
- [Coding Standards & Conventions](#coding-standards--conventions)
  - [Go Guidelines](#go-guidelines)
  - [Terraform Schema Guidelines](#terraform-schema-guidelines)
  - [Security & Sensitive Data Standards](#security--sensitive-data-standards)
- [Commit Message Guidelines (Conventional Commits)](#commit-message-guidelines-conventional-commits)
- [Submitting a Pull Request](#submitting-a-pull-request)
  - [Pre-Flight Checklist](#pre-flight-checklist)
  - [Hermetic Nix CI Matrix (11 Jobs)](#hermetic-nix-ci-matrix-11-jobs)
  - [Code Review and Squash-Merge](#code-review-and-squash-merge)
- [Additional Architectural References](#additional-architectural-references)

---

## Code of Conduct

We are committed to providing a welcoming, inclusive, and harassment-free experience for everyone. Please be respectful, considerate, and collaborative in all communications, issue threads, and pull request reviews.

---

## How to Participate

### Reporting Issues & Bugs

Before filing a new bug report, search the [open issues](https://github.com/rancher/terraform-provider-rancher2/issues) to verify whether the problem has already been reported.

When creating a new bug report, provide a detailed description including:
- **Terraform Provider Version**: The exact version of `terraform-provider-rancher2` in use.
- **Terraform Version**: The output of `terraform version`.
- **Rancher Version**: The specific Rancher Server version (e.g., `v2.14.0`).
- **Minimal Reproduction Configuration**: A minimal, sanitized Terraform configuration (`.tf`) reproducing the issue.
- **Execution Output**: Relevant debug logs (`TF_LOG=DEBUG` output) with any sensitive tokens, secrets, or passwords redacted.
- **Expected vs. Actual Behavior**: Clear description of what should have occurred versus what happened.

### Proposing Features & Enhancements

If you wish to add a new resource, data source, or field to the provider:
1. Open an enhancement issue outlining the proposed change and its use case.
2. Reference the corresponding Rancher API endpoints or documentation.
3. Discuss the architectural approach with repository maintainers before writing substantial code.

### Scope Boundaries

- **Node drivers:** Native provider support is limited to node drivers bundled and shipped by default with Rancher. Third-party drivers versioned separately from Rancher (such as `docker-machine-driver-pve`) are outside the scope of built-in provider support because they introduce an unbundled dependency that the project cannot maintain and support. See [Rancher's machine-driver definitions](https://github.com/rancher/rancher/blob/main/pkg/data/management/machinedriver_data.go).

---

## Development Environment Setup

### Prerequisites

To build and test the provider locally, your machine requires:
- **Go**: Version 1.26.5 or later.
- **Terraform**: Version 1.5.0 or later.
- **Git**: Configured with a GPG signing key for commits.

### Hermetic Environment with Nix (Recommended)

This repository includes a hermetic Nix configuration (`flake.nix`). If you have Nix installed with flakes enabled, you can load the exact toolchains, linters, and compilers used by our CI runners without installing dependencies globally:

```bash
# Enter the hermetic Nix development shell
source .envrc
```

### Forking and Setting Up Remotes

All external contributions must be submitted from your own fork:

1. Fork `rancher/terraform-provider-rancher2` on GitHub.
2. Clone your fork locally:
   ```bash
   git clone git@github.com:<your-username>/terraform-provider-rancher2.git
   cd terraform-provider-rancher2
   ```
3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/rancher/terraform-provider-rancher2.git
   ```
4. Keep your local `main` branch synchronized with upstream:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```
5. Create a branch from (your fork's) main to work out of:
   ```bash
   git checkout -b <your-branch-name>
   ```
   - this step helps if you need to make more than one PR at a time to the repo

*Note: Direct pushes to the `upstream` repository are restricted. Always push your changes to branches on your personal fork.*

---

## Branching Strategy & Value Stream

### Value Flow

1. Issue is created describing a problem or feature, this is required
2. PR targeting `main` is created to address issue, linked with `- Addresses: #<issue number>`
3. Proper SUSE team is added to the PR to review
4. The issue is labeled with the latest release branch and a sub-issue is created to track release and testing
5. Additional labels may be added to release to multiple branches, each label will generate a sub-issue
6. Copilot Review is automatically triggered for every PR and push
7. All CI tests pass, PR review approved from a SUSE teammate, and all review conversations are resolved (including Copilot review)
8. The PR merges into `main`
9. CI automatically generates `cherry-pick` PRs to release branches for each sub-issue attached to the Issue
10. CI tests all pass on automated cherry-pick PRs and they are automatically merged into the release branch
11. QA Engineers pick up sub-issues and close them when they are approved for release
12. After release is cut the issue (created in step 1) is closed as resolved

### Developer Flow

1. Create issue for feature or problem if it doesn't already exist.
2. Create PR
   1. fork repo
   2. sync main
   3. generate branch from main
   4. update codebase
   5. use `make build` and `make test` to validate
   6. validate against real modules, or take a look at run_tests.sh for how we will run our pre-release e2e tests
3. Respond and address review comments
   1. AI comments are meant to be helpful, not blocking
   2. Respond to AI comments with why something isn't valid or is out or scope before resolving the thread
   3. Respond and address all human reviewer concerns before resolving them
4. Wait and Watch
   1. Maintainers may ask for help if the backport to release branch has merge conflicts
   2. QA Engineers may find problems that need to be addressed before release
   3. QA Engineers will close sub-issues when they are approved for release
   4. Issue (from step 1) will close when the code is released


### Trunk-Based Development (`main`)

All enhancements, bug fixes, and maintenance work **must target the `main` branch**. Do not open pull requests directly against release maintenance branches (such as `release/v14` or `release/v15`) unless explicitly requested by maintainers to resolve cherry-pick merge conflicts.

### Release Branches & Rancher Version Alignment

The Rancher2 Terraform Provider major versions correspond to Rancher minor release lines:

| Terraform Provider Branch | Provider Major Version | Target Rancher Minor Version |
| :--- | :--- | :--- |
| `release/v13` | `v13.x` | Rancher `2.13.x` |
| `release/v14` | `v14.x` | Rancher `2.14.x` |
| `release/v15` | `v15.x` | Rancher `2.15.x` |
| `release/v16` | `v16.x` | Rancher `2.16.x` |

See [docs/compatibility-matrix.md](docs/compatibility-matrix.md) for full version alignment details.

### Automated Backport Orchestration

Our goal is for all backports to happen automatically:
1. When your PR merges into `main`, automation automatically cherry-picks the commit into target `release/v<major>` branches.
2. Automated backport pull requests are opened and verified through acceptance test suites on each target release branch.
3. Once tests pass, the backport PRs merge automatically.

---

## The Issue-First Requirement

Every pull request targeting `main` **must be associated with an open GitHub issue**. Pull requests submitted without an associated issue will fail automated PR validation.

### Linking Issues in PR Descriptions

In your pull request description, you must explicitly link the tracking issue using one of the following formats:

```markdown
- Addresses: #1234
```
or
```markdown
* Addresses: #1234
```

Automated PR validation (`validate-pr.yaml`) validates that:
- The referenced issue exists in the repository.
- The reference points to a GitHub issue (not another pull request).
- The referenced issue is currently **open**.

### Target Release Labels & Sub-Issues

Before merging a pull request on `main`, assign release labels to the linked issue (e.g., `release/v15`, `release/v14`) to identify which release branches to port the changes to:
- Automation adds the latest release branch label to the linked issue.
- Automation creates child sub-issues for each targeted release line.
- Once merged, the automation cherry-picks the change into the respective branches.

---

## Building and Local Development

The repository includes a `GNUmakefile` with standard targets for compiling, linting, and formatting.

### Building the Provider Binary

Compile the provider binary locally:

```bash
make build
```

This builds the binary and outputs it to `./bin/terraform-provider-rancher2`.

### Formatting and Linting

Run local formatters and static analysis before submitting code:

```bash
# Format Go code
make fmt

# Run Go linters (golangci-lint)
make lint

# Validate that test packages compile without errors
make lint-tests

# Lint GitHub Actions workflows and automation scripts
make lint-workflows

# Run the complete automation linting suite (Terraform, ShellCheck, ESLint, Gitleaks)
node .github/workflows/scripts/lint.js all
```

Individual linting suites can also be run via `lint.js`:
- `node .github/workflows/scripts/lint.js terraform`: Runs `terraform fmt` and `tflint`.
- `node .github/workflows/scripts/lint.js shellcheck`: Scans shell scripts with `shellcheck -x`.
- `node .github/workflows/scripts/lint.js eslint`: Runs ESLint on JavaScript files.
- `node .github/workflows/scripts/lint.js gitleaks`: Scans for secrets.

---

## Testing Guidelines

Thorough testing is essential for all contributions. Pull requests introducing new features or bug fixes must include corresponding tests.

### Unit Tests

Unit tests execute quickly and do not require live infrastructure:

```bash
# Run all unit tests with gotestsum
make test

# Run a specific unit test by name
make dt -- t='-run=TestName ./...'
```

### Acceptance Tests (End-to-End)

Acceptance tests execute full infrastructure lifecycle scenarios (creation, update, refresh, deletion) against real Rancher and cloud infrastructure.

> **WARNING**: Acceptance tests provision live cloud resources (typically in AWS) and incur cloud infrastructure costs. You must configure cloud credentials in your local environment.

```bash
# Run acceptance tests
./run_tests.sh
# or
make testacc

# Run a specific acceptance test
./run_tests.sh -t TestOneBasic
# or
make et -- t=TestOneBasic

# Run acceptance tests sequentially
./run_tests.sh -s
```

### Cleaning Up Acceptance Test Leftovers

If an acceptance test run is interrupted, orphan resources may remain in your cloud account. Clean them up using the cleanup script:

```bash
# Clean up test leftovers by test identifier
./run_tests.sh -c <identifier>
# or
make clean -- i=<identifier>
```

---

## Coding Standards & Conventions

### Go Guidelines

Adhere strictly to standard Go practices and project conventions:
- **Error Handling**:
  - Always check errors immediately. Never ignore errors with `_` unless explicitly documented.
  - Wrap errors to preserve context using `fmt.Errorf("...: %w", err)`.
  - Avoid nested `if` blocks for error handling; keep the happy path left-aligned.
  - Never use `panic()` for control flow.
- **Context & Concurrency**:
  - Pass `ctx context.Context` as the first argument in any function making network calls or blocking operations.
  - Never store `context.Context` inside struct definitions.
  - Ensure every goroutine has an explicit exit path to prevent leaks.
- **HTTP Clients**:
  - Never use `http.DefaultClient`. Always provide explicit timeouts (e.g., `Timeout: 10 * time.Second`).
- **Memory & Slices**:
  - When the final length or capacity is known, pre-allocate slices and maps (`make([]T, 0, capacity)`).
- **Naming**:
  - Use `PascalCase` for exported identifiers and `camelCase` for unexported identifiers.
  - Avoid getter prefixes (use `User()` instead of `GetUser()`).
- **Test Module Structure**:
  - The `./test` directory is a standalone Go module. Tests must maintain clean fixtures and independent scopes.

Detailed Go rules can be found in [docs/development/CodingStandards/Go.md](docs/development/CodingStandards/Go.md).

### Terraform Schema Guidelines

- **Attribute Naming**: Follow snake_case naming consistent with the upstream Rancher API.
- **Validations**:
  - Use `variable` validations for simple regex matches on single attributes.
  - Use `check` blocks for multi-variable validations.
  - Use `terraform_data` preconditions for complex validations evaluating output attributes of dependent resources.
- Detailed Terraform rules can be found in [docs/development/CodingStandards/Terraform.md](docs/development/CodingStandards/Terraform.md).

### Security & Sensitive Data Standards

- **Zero Hardcoded Secrets**: Never commit passwords, tokens, API keys, certificates, or credentials.
- **Sensitive Schema Flagging**: Ensure sensitive Terraform schema fields (such as passwords, tokens, private keys) are marked with `Sensitive: true`.
- **Automated Secret Audits**: Every pull request is verified against `gitleaks` and schema sensitivity audits (`check_sensitive.go`).

---

## Commit Message Guidelines (Conventional Commits)

This repository enforces Conventional Commits to support automated changelog generation and SemVer release management via Release Please.

### Format Requirements

Commit subject lines must:
1. Start with an approved prefix:
   - `fix: ` — Bug fixes (generates a patch release).
   - `feat: ` or `feature: ` — New features (generates a minor release).
   - `refactor!: `, `feat!: `, or `feature!: ` — Breaking changes (generates a major release).
   - `chore(main): ` — Maintenance tasks targeting `main`.
2. Be concise: **72 characters or fewer** is recommended (automated CI enforces a maximum of 100 characters).
3. Use correct spelling in standard technical English (validated via CSpell).

### Example Commit Messages

```text
feat: add support for cluster proxy configuration v2
fix: resolve nil pointer when parsing custom user tokens
```

---

## Submitting a Pull Request

### Pre-Flight Checklist

Before opening your pull request, verify that:
1. [ ] Your changes are based on the latest `upstream/main`.
2. [ ] An open GitHub issue exists and is referenced in your PR description (`- Addresses: #<issue>`).
3. [ ] Code is formatted (`make fmt`).
4. [ ] All linters pass without errors (`make lint`, `make lint-workflows`).
5. [ ] Relevant unit and acceptance tests pass (`make test`).
6. [ ] Commits follow Conventional Commits formatting and length rules.
7. [ ] No sensitive credentials, keys, or tokens are committed (`gitleaks detect`).

### Hermetic Nix CI Matrix (11 Jobs)

When you open a PR targeting `main`, GitHub Actions runs an 11-job CI matrix inside hermetic Nix containers:
1. **`build`**: Cross-platform binary compilation for Linux, Darwin, and Windows (AMD64 & ARM64).
2. **`test`**: Go unit tests via `gotestsum`.
3. **`terraform`**: Terraform format and tflint checks.
4. **`actionlint`**: Workflow schema validation.
5. **`node-check`**: Node.js syntax checks on automation scripts.
6. **`node-test`**: Node.js unit tests for scripts.
7. **`eslint`**: JavaScript code quality checks.
8. **`shellcheck`**: Shell script static analysis.
9. **`validate-commit-message`**: Conventional commit prefix, length, and spelling validation.
10. **`gitleaks`**: Repository secret scanning.
11. **`sensitive-checks`**: Schema sensitivity analysis.

### Code Review and Squash-Merge

- **Reviewers**: A SUSE employee will review your pull request for architectural soundness, test coverage, and code quality.
  - if no-one responds to your PR within a week, feel free to message in [rancher-users.slack](https://rancher-users.slack.com) discuss-terraform channel
- **Addressing Feedback**: Push additional commits to your branch to address review comments.
  - AI doesn't need to approve your PR, but please comment on open threads before resolving them
- **Merge**: Once approved and all CI checks pass, automation will squash-merge your PR into `main`.
- **Backport**: Following merge, backport automation automatically handles cherry-picking into targeted release branches.

---

## Additional Architectural References

For deeper insights into the provider's architecture and automation systems, consult the developer documentation:

- **[Project Management & Delivery Guide](docs/development/ProjectManagementGuide.md)**: Overview of how work flows through issues, PRs, backports, QA, and releases.
- **[Pull Request & Merge Lifecycle](docs/development/GitHubWorkflows/PullRequestLifecycle.md)**: Technical details on PR validation and CI workflows.
- **[Backport Orchestration](docs/development/GitHubWorkflows/BackportOrchestration.md)**: Deep dive into automated cherry-picking and sub-issue trees.
- **[Privileged AWS End-to-End Testing](docs/development/GitHubWorkflows/EndToEndTesting.md)**: Acceptance test execution, AWS OIDC, and distributed locks.
- **[Release Process & Automation](docs/development/ReleaseProcess.md)**: Release candidate tagging, Vault GPG signing, and Release Please manifests.
- **[Coding Standards Index](docs/development/CodingStandards.md)**: Central directory for Go, Terraform, Shell, and workflow standards.
- **[Version Compatibility Matrix](docs/compatibility-matrix.md)**: Compatibility between Rancher versions and provider releases.
