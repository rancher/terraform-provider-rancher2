# Rancher2 Provider: Development Process

Welcome to the development guide for the Rancher2 Terraform Provider.

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
