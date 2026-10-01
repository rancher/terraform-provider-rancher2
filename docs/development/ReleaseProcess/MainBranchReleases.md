# Release Process: Trunk-Based Release from Main

## Abstract

To establish a highly standardized, predictable, and simple release lifecycle, the repository adopts a **trunk-based development release strategy**. This system coordinates automated production releases directly from the `main` branch, completely eliminating long-lived, error-prone `release/v*` branch structures and aligning our SemVer delivery directly with squash-merged conventional commits on `main`.

---

## Technical Specification

### 1. Trunk-Based Release Please Trigger

Our pipeline coordinates automatically with `release-please` to manage stable SemVer versions and automate draft release manifests.

- **Trigger**: Opening/merging pull requests targeting `main`.
- **Workflow**: `.github/workflows/release.yml` triggers upon direct pushes to `main`.
- **Release Please Integration**: It scans conventional commits squash titles to calculate version bumps, then generates and maintains a pending "Release PR" (e.g., `chore: release v1.2.3`).
- **Heavy Acceptance Testing**: Merging the Release PR triggers our exhaustive integration and end-to-end acceptance test suites (`test.sh acc-relay`) to guarantee full production compatibility before binary compilation.

### 2. Automated Release Candidate (RC) Workflow

When pull requests are merged into release branches (`release/v*`), the `.github/workflows/release-branch-merge.yml` workflow triggers and dispatches the unified `rc-release.yml` workflow to cut and build a new Release Candidate.

- **Dynamic RC Tag Calculation**: `computeNextRcTag` in `.github/workflows/scripts/releases.js` inspects existing tags, `release-please-config.json` (`release-as`), and `.release-please-manifest.json` on the target release branch to compute the next sequential RC tag (e.g. `v15.2.0-rc.1` -> `v15.2.0-rc.2`).
- **Unified Pipeline & Attribution**: The release candidate is built and signed via `.github/workflows/rc-release.yml` executing on `main`, ensuring standardized GoReleaser configuration, Nix tooling, and Vault GPG signing. Git tag creation and push are authenticated with maintainer or GitHub App credentials for complete traceability.
- **Automated Validation Notifications**: Upon successful RC generation, `rc-notify` automatically identifies and comments on open backport issues (`internal/backport` and `internal/merged`) to notify developers that the candidate binary is ready for testing.

---

## Standing Implementation Decisions

1. **No Outdated Branch Releases**: Releases must always be cut directly from `main`. Branches prefixed with `release/v*` are obsolete and blocked from triggering stable releases.
2. **Attribution Principle**: All release candidate tags are created and pushed in the unified `rc-release.yml` workflow using authenticated git credentials attributed to the triggering maintainer or GitHub App bot.
3. **Targeted RC Notifications**: Release candidate notifications are integrated directly into `rc-release.yml` via the `rc-notify` mode in `releases.js`. Upon completion of an RC release, it automatically discovers and comments on open issues associated with the release branch that have been marked `internal/backport` and `internal/merged`, ensuring immediate visibility for validation without issue spam.
4. **Manual Full Release Publishing**: The manual full release workflow (`manual-release.yml`) builds draft releases with GoReleaser (configured with `draft: true` in `.goreleaser.yml`). To finalize these, the workflow automatically invokes the `publish-release` SCRIPT_MODE in `releases.js` at the end of the run to promote the draft release to published. To ensure extreme robustness, `releases.js` uses a bidirectional fallback mechanism that checks for both the exact tag and its `v` prefix variant, while the manual workflows programmatically enforce a `v*` prefix validation on triggering.
5. **Unified RC Release Workflow (`rc-release.yml`)**: Release candidate generation is unified into `.github/workflows/rc-release.yml` ("Create RC Release"), which always executes from the `main` branch. This guarantees consistent builds, GoReleaser configuration, Nix tooling, and Vault GPG signing regardless of target branch age. The workflow is invoked manually via `workflow_dispatch` or automatically upon PR merges into release branches via `release-branch-merge.yml` using the Vault GitHub App token.
