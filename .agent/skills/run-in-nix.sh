#!/usr/bin/env bash
#
# Skill: run-in-nix.sh
# Description: Executes a given command inside the standardized Nix development environment.
# Usage: ./run-in-nix.sh "<command>"

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "Error: Command required."
  echo "Usage: $0 \"<command>\""
  exit 1
fi

COMMAND="$1"

echo "Running command in Nix environment: ${COMMAND}"

nix develop --ignore-environment \
  --extra-experimental-features nix-command \
  --extra-experimental-features flakes \
  --keep HOME --keep SSH_AUTH_SOCK --keep GPG_SIGNING_KEY --keep NIX_SSL_CERT_FILE --keep NIX_ENV_LOADED --keep TERM \
  --command bash -c "${COMMAND}"
