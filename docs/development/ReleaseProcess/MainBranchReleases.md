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

### 2. Automated Release Candidate (RC) Tagging

When the `release-please` workflow triggers, it calculates the target version and generates intermediate **Release Candidate (RC) tags** for binary pre-releases.

- **Dynamic RC Tag Calculation**: `.github/workflows/scripts/create-push-rc-tag.js` executes inside an `actions/github-script` environment. It fetches existing tags via the GitHub API, identifies the largest existing RC suffix for the target base version (e.g. `v1.2.3-rc.x`), increments the count (`v1.2.3-rc.0` -> `v1.2.3-rc.1`), and pushes the new tag directly to the active `context.sha`.
- **Biometric / Bot Attribution**: Tag creation is handled programmatically via `github.rest.git.createRef` utilizing the runner's authenticated token, ensuring precise, bot-attributed GitHub releases without local SSH or GPG keychain dependencies on the GHA runner.

---

## Standing Implementation Decisions

1. **No Outdated Branch Releases**: Releases must always be cut directly from `main`. Branches prefixed with `release/v*` are obsolete and blocked from triggering stable releases.
2. **Attribution Principle**: All release candidate tags are created using the native GitHub REST API and authenticated GITHUB_TOKEN rather than direct command-line git commands, ensuring cryptographically secure, traceable tag histories.
3. **No Issue Notification Overhead**: Outdated manual release notification scripts (such as `rc-notify.js`) are permanently removed, streamlining releases of core Terraform binaries without tracking issue spam.
