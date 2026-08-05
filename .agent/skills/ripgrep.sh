#!/usr/bin/env bash
#
# Skill: ripgrep.sh
# Description: Executes ripgrep (rg) to search for a pattern in a file or directory.
# Usage: ./ripgrep.sh options

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "Error: Pattern required."
  echo "Usage: $0 ripgrep-options"
  exit 1
fi

# shellcheck disable=SC2145
echo "Executing Ripgrep with options: $@"

# use ripgrep in Nix environment
# shellcheck disable=SC2145
nix develop \
  --extra-experimental-features nix-command \
  --extra-experimental-features flakes \
  --command bash -c "rg $@"
