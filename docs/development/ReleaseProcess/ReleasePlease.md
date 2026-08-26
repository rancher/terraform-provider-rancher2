# Release Please Automation

---

## Abstract

This component details the configuration, inputs, outputs, and CLI usage of the Google Release Please tool, which automates our changelog generation and version bumping on merge.

---

## Release Please Action

`release-please-action` automates CHANGELOG generation, GitHub release creation, and version bumps by parsing git history for Conventional Commits (`fix:`, `feat:`, `feat!:`). It maintains a running Release PR containing the updated changelog and version bumps. Upon merge, the action tags the commit and creates a GitHub Release.

### Action Inputs

Key inputs configured under the `with:` block:

- `token`: GitHub secret token (default: `secrets.GITHUB_TOKEN`). Recommended to use a custom PAT to allow triggers for subsequent workflows.
- `release-type`: Release strategy (e.g., `node`, `python`, `go`, `terraform-module`, `simple`). If omitted, defaults to manifest-based config.
- `path`: Create a release from a subdirectory.
- `config-file` / `manifest-file`: Paths for manifest-based monorepo configuration (defaults to `release-please-config.json` / `.release-please-manifest.json`).

### Action Outputs

Outputs available to orchestrate downstream steps:

- `releases_created`: `true` if a release was created/tagged.
- `release_created`: `true` if root release was created.
- `tag_name`: Generated git tag name (e.g., `v1.2.3`).
- `version`, `major`, `minor`, `patch`: Version information.
- `upload_url` / `html_url`: URLs to the GitHub release.
- _Monorepo outputs_ are prefixed with the path, e.g., `<path>--release_created`.

### Usage Examples

**Basic Configuration (Single Project):**

```yaml
steps:
  - uses: googleapis/release-please-action@v4
    with:
      token: ${{ secrets.MY_RELEASE_PLEASE_TOKEN }}
      release-type: go
```

**Manifest Mode (Advanced/Monorepos):**

```yaml
steps:
  - uses: googleapis/release-please-action@v4
    with:
      token: ${{ secrets.MY_RELEASE_PLEASE_TOKEN }}
      config-file: release-please-config.json
      manifest-file: .release-please-manifest.json
```

**Downstream Publication (e.g. running goreleaser after release created):**

```yaml
steps:
  - uses: googleapis/release-please-action@v4
    id: release
  - uses: actions/checkout@v4
    if: ${{ steps.release.outputs.release_created }}
  - run: make publish
    if: ${{ steps.release.outputs.release_created }}
```

## Release Please CLI & Core Concepts

`release-please` relies on the **Conventional Commits** specification to determine how to bump versions and generate changelogs.

### Conventional Commits Handling

- **`fix:`**: Bug fixes. Bumps **Patch** version (`1.0.0` -> `1.0.1`).
- **`feat:`**: New feature. Bumps **Minor** version (`1.0.0` -> `1.1.0`).
- **`!` (Breaking Change)**: Appending `!` to the type (e.g., `feat!:`) or `BREAKING CHANGE:` in footer bumps **Major** version (`1.0.0` -> `2.0.0`).
- **Force Version**: Add `Release-As: x.y.z` to the commit body to force a specific version.

### CLI Commands

The `release-please` CLI is executed to automate release PRs and GitHub Releases.

1. **`release-pr`**: Scans git history, determines next version, and opens/updates a Release PR.
2. **`github-release`**: Executed after the Release PR is merged. Tags the merge commit and creates the GitHub Release.

**Example:**

```bash
# Update release PR
npx release-please release-pr --token=$GITHUB_TOKEN --repo-url=owner/repo --release-type=node

# Create GitHub release on merge
npx release-please github-release --token=$GITHUB_TOKEN --repo-url=owner/repo
```

### Manifest-Based Configuration (Monorepos)

Uses two root files:

- `.release-please-manifest.json`: Tracks the current version of each path.

  ```json
  {
    ".": "1.0.0",
    "packages/core": "2.3.1"
  }
  ```

- `release-please-config.json`: Configures the release strategy per path.

  ```json
  {
    "packages": {
      ".": { "release-type": "node" },
      "packages/core": { "release-type": "typescript", "extra-files": ["src/version.ts"] }
    }
  }
  ```

### Manifest Configuration Options (`release-please-config.json`)

The configuration schema handles global behavior and per-package options. Global options sit at the root level, while package-specific options sit under the `packages` map (where the keys are the directory paths).

**Global Options:**

- `packages`: **(Required)** Map of per-path component configurations.
- `plugins`: Array of plugins to apply for extra PR processing (e.g., `node-workspace`, `cargo-workspace`, `linked-versions`).
- `label` / `release-label`: Comma-separated labels applied to the open release PR, or after it has been tagged.
- `sequential-calls`: `boolean` to open PRs sequentially (avoids GitHub secondary rate limits).

**Package-Level Options (Can be set globally as defaults):**

- **Strategy & Versioning:**
  - `release-type`: The release strategy/language to use (`node`, `go`, `rust`, `simple`, etc.).
  - `bump-minor-pre-major` / `bump-patch-for-minor-pre-major`: Booleans to restrict semver bumps when version is `< 1.0.0`.
- **Changelog:**
  - `changelog-path`: Path to track release notes (default: `CHANGELOG.md`).
  - `skip-changelog`: `boolean` to disable changelog generation.
  - `changelog-sections`: Array of objects (`type`, `section`, `hidden`) to customize headers (e.g., hiding `chore` or renaming `feat` to "Features").
- **GitHub Release & Tagging:**
  - `skip-github-release`: `boolean` to skip tagging GitHub releases.
  - `include-component-in-tag`: `boolean` to prefix tag with the package path (default `true` for monorepos).
  - `draft` / `prerelease`: `boolean` flags for the created GitHub release.
- **Pull Request:**
  - `draft-pull-request`: `boolean` to open the Release PR in draft mode.
  - `pull-request-title-pattern`: Customize the PR title.
- **Files:**
  - `extra-files`: Array of additional files to update alongside the primary language files. Supports paths (using generic `# x-release-please-version` annotations) or objects specifying parsers (`jsonpath`, `xpath`).
