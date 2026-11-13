#!/usr/bin/env bash
#
# setup-release-branch.sh
#
# DEPRECATED: This script is deprecated and has been superseded by the GitHub Actions workflow
# '.github/workflows/setup-release-branch.yml' and '.github/workflows/scripts/setup-release-branch.js'.
#
# To create and configure a release branch, trigger the workflow manually:
#   gh workflow run setup-release-branch.yml -f branch="release/v<version>"
#
set -euo pipefail

echo "Error: setup-release-branch.sh is deprecated." >&2
echo "Release branches are now created and configured via the 'Setup Release Branch' GitHub Actions workflow." >&2
echo "Run: gh workflow run setup-release-branch.yml -f branch=\"release/v<version>\"" >&2
exit 1
