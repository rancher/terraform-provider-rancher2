#!/usr/bin/env bash
#
# Skill: run-in-nix.sh
# Description: Executes a given command inside the standardized Nix development environment.
# Usage: ./run-in-nix.sh "<command>"
# - quotes around the command are critical

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "Error: Command required."
  echo "Usage: $0 \"<command>\""
  exit 1
fi

COMMAND="$1"

echo "Running command in Nix environment: ${COMMAND}"

nix develop \
  --extra-experimental-features nix-command \
  --extra-experimental-features flakes \
  --command bash -c "${COMMAND}"
