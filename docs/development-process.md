# Rancher2 Provider: Development Process

Welcome to the development guide for the Rancher2 Terraform Provider.

For detailed documentation on CI/CD pipelines, automated backports, and release mechanics, see **[GitHub Workflows & Automation Architecture](workflows.md)**.

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
4. **Pull Request Validation**: Once opened against `main`, automated CI (`validate-pr.yaml`) validates the issue reference, and `pull_request.yaml` executes the 11-job hermetic Nix container CI matrix.
5. **Teammate Review & Auto-Merge**: Any teammate or collaborator with repository triage, write, or admin permissions can review and approve pull requests. Approved pull requests with all passing CI checks are automatically squash-merged into `main` by `auto-merge.yml`. Maintainer approval is not required.
6. **Automated Backport**: Merging into `main` automatically triggers `backport-prs.yml`, opening cherry-pick PRs against all target `release/v*` maintenance branches.
7. **Release Candidate (RC) Publication**: Merging an automated backport PR into a release branch triggers `release-please.yml`, which tags and builds the next Release Candidate (e.g. `v15.2.0-rc.1`) and updates [Project 73](https://github.com/orgs/rancher/projects/73) issue status to **Ready for QA**.
