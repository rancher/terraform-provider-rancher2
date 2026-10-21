# Rancher2 Provider: Development Process

Welcome to the development guide for the Rancher2 Terraform Provider.

For the comprehensive architectural specifications, end-to-end delivery lifecycle, and automation blueprints, please consult the primary developer documentation under **[docs/development/](development/README.md)**:

- **[Project Management & Delivery Guide](development/ProjectManagementGuide.md)**: High-level value stream navigation for PMs, release coordinators, and QA to follow work moving through the system.
- **[End-to-End Value Stream Lifecycle](development/GitHubWorkflows/Overview.md)**: Complete guide to how changes move from pull requests through testing, backporting, RC staging, and production releases.
- **[Pull Request & Merge Lifecycle](development/GitHubWorkflows/PullRequestLifecycle.md)**: PR description validation rules, Conventional Commits format, and the 11-job Nix container CI matrix.
- **[Backport Orchestration & Issue Tracking](development/GitHubWorkflows/BackportOrchestration.md)**: Automated issue tracking adoption, GitHub Sub-Issues hierarchy, and cherry-picking to release branches.
- **[Privileged AWS End-to-End Testing](development/GitHubWorkflows/EndToEndTesting.md)**: Pre-merge acceptance test matrix, AWS OIDC authentication, and distributed locking.
- **[Release Process & Automation](development/ReleaseProcess.md)**: Release candidate generation, Vault GPG signing, and Google Release Please manifest automation.

---

### Quick Start for Contributors

1. **Issue Association**: Every contribution targeting `main` must be associated with an open GitHub issue. In your PR description, include:
   ```markdown
   - Addresses: #<issue-number>
   ```
2. **Conventional Commits**: Ensure commit messages follow Conventional Commits (e.g. `fix:`, `feat:`, or `chore(main):`) and subject lines stay under 100 characters.
3. **Local Testing & Linting**: Run tests and linters locally before submitting:
   ```bash
   make build
   make test
   make fmt
   make lint
   ```
4. **Pull Request Validation**: Once opened against `main`, automated CI will validate the issue reference and execute the hermetic Nix container CI matrix.
5. **Maintainer Review & Backport**: Repository maintainers (`TERRAFORM_MAINTAINERS`) will review and squash-merge approved changes into `main`. Merging triggers automated cherry-picking into target `release/v<major>` maintenance branches.
