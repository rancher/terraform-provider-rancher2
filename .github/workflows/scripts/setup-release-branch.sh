#!/usr/bin/env bash
#
# setup-release-branch.sh
# Automates preparing a release branch with thin caller stubs that delegate to main,
# and configures release-please-config.json and .release-please-manifest.json.
#
# Usage:
#   ./.github/workflows/scripts/setup-release-branch.sh [branch-name]
# Example:
#   ./.github/workflows/scripts/setup-release-branch.sh release/v15
#   ./.github/workflows/scripts/setup-release-branch.sh v16
#
set -euo pipefail

# 1. Determine target branch name
TARGET_BRANCH="${1:-}"
if [[ -z "$TARGET_BRANCH" ]]; then
  TARGET_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)
fi

if [[ -z "$TARGET_BRANCH" || "$TARGET_BRANCH" == "HEAD" ]]; then
  echo "Error: Branch name must be provided as an argument or checked out in git." >&2
  echo "Usage: $0 <release/vX | vX>" >&2
  exit 1
fi

# Normalize branch version
# e.g. "release/v15" -> "15", "v15" -> "15", "15" -> "15"
RAW_VERSION=$(echo "$TARGET_BRANCH" | sed -E 's#^(refs/heads/)?(release/)?v?##')
if [[ -z "$RAW_VERSION" || ! "$RAW_VERSION" =~ ^[0-9]+$ ]]; then
  echo "Error: Target branch '$TARGET_BRANCH' does not match a recognized release branch format (e.g., release/v15 or v15)." >&2
  exit 1
fi

RELEASE_VERSION="v${RAW_VERSION}.0.0"
MANIFEST_VERSION="${RAW_VERSION}.0.0"

echo "=== Setting up release branch for '$TARGET_BRANCH' (Version: $RELEASE_VERSION) ==="

# 2. Update release-please-config.json with release-as: <branch name>.0.0
CONFIG_FILE="release-please-config.json"
if [[ -f "$CONFIG_FILE" ]]; then
  echo "Updating $CONFIG_FILE (release-as: $RELEASE_VERSION)..."
  node -e '
    import("fs").then(fs => {
      const file = process.argv[1];
      const version = process.argv[2];
      const config = JSON.parse(fs.readFileSync(file, "utf8"));
      config["release-as"] = version;
      fs.writeFileSync(file, JSON.stringify(config, null, 2) + "\n", "utf8");
    });
  ' "$CONFIG_FILE" "$RELEASE_VERSION"
else
  echo "Creating $CONFIG_FILE (release-as: $RELEASE_VERSION)..."
  cat <<EOF > "$CONFIG_FILE"
{
  "release-as": "$RELEASE_VERSION",
  "always-update": true,
  "commit-search-depth": 100,
  "prerelease": true,
  "draft": true,
  "force-tag-creation": true,
  "include-v-in-tag": true,
  "include-component-in-tag": false,
  "packages": {
    ".": {
      "release-type": "go"
    }
  },
  "\$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json"
}
EOF
fi

# 3. Update .release-please-manifest.json
MANIFEST_FILE=".release-please-manifest.json"
echo "Updating $MANIFEST_FILE ($MANIFEST_VERSION)..."
cat <<EOF > "$MANIFEST_FILE"
{
  ".": "$MANIFEST_VERSION"
}
EOF

# 4. Clean up any existing workflows in .github/workflows except the standard stubs and listener
echo "Cleaning up non-stub workflows in .github/workflows..."
mkdir -p .github/workflows
for file in .github/workflows/*.yml .github/workflows/*.yaml; do
  [[ -f "$file" ]] || continue
  filename=$(basename "$file")
  if [[ "$filename" != "pull_request.yaml" && "$filename" != "release.yml" && "$filename" != "e2e-pr-listener.yml" ]]; then
    rm -f "$file"
  fi
done

# 5. Create Pull Request caller stub referencing @main
echo "Creating .github/workflows/pull_request.yaml stub calling @main..."
cat << 'EOF' > .github/workflows/pull_request.yaml
name: Pull Request

on:
  pull_request:
  merge_group:

permissions: {}

jobs:
  pull-request:
    permissions:
      contents: read
      pull-requests: read
    uses: rancher/terraform-provider-rancher2/.github/workflows/pull_request.yaml@main
    secrets: inherit
EOF

# 6. Create Release caller stub referencing @main
echo "Creating .github/workflows/release.yml stub calling @main..."
cat << 'EOF' > .github/workflows/release.yml
name: Release Please

on:
  push:
    branches:
      - 'release/v*'

permissions: {}

jobs:
  release:
    permissions:
      contents: write
      id-token: write
      issues: write
      pull-requests: write
      actions: read
    uses: rancher/terraform-provider-rancher2/.github/workflows/release-please.yml@main
    secrets: inherit
EOF

# 7. Ensure e2e-pr-listener.yml is present to trigger privileged E2E tests from main
LISTENER_FILE=".github/workflows/e2e-pr-listener.yml"
if [[ ! -f "$LISTENER_FILE" ]]; then
  echo "Creating $LISTENER_FILE to preserve E2E test triggering..."
  cat << 'EOF' > "$LISTENER_FILE"
name: PR Executor Listener
# This unprivileged workflow acts as a middleman to securely trigger e2e tests.

on:
  pull_request:
    branches:
      - 'release/v*'
    types: [opened, synchronize, reopened]

permissions: {}

jobs:
  emit-event:
    name: 'Emit E2E Workflow Run Event'
    runs-on: ubuntu-latest
    timeout-minutes: 5
    permissions:
      contents: read
    steps:
      - name: 'Execute PR Listener'
        run: echo "PR listener executed. The completion of this workflow will trigger the privileged E2E Tests (Release PRs) workflow safely from the base branch via the workflow_run event."
EOF
fi

echo "=== Release branch setup complete for '$TARGET_BRANCH' ==="
echo "Files ready to stage:"
echo "  - .github/workflows/pull_request.yaml (caller stub -> @main)"
echo "  - .github/workflows/release.yml (caller stub -> @main)"
echo "  - $LISTENER_FILE (E2E test listener)"
echo "  - $CONFIG_FILE (release-as: $RELEASE_VERSION)"
echo "  - $MANIFEST_FILE (version: $MANIFEST_VERSION)"
