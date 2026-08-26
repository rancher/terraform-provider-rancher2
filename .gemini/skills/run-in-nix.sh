#!/usr/bin/env bash
#
# Skill: run-in-nix.sh
# Description: Executes a given command inside the standardized Nix development environment, or lists installed tools.
# Usage: .gemini/skills/run-in-nix.sh "<command>"

set -euo pipefail

show_help() {
  cat <<EOF
Usage: run-in-nix.sh [options] "<command>"

Executes a given command inside the standardized Nix development environment, or queries installed tools.

Options:
  --list-tools         List all executable tools installed in the Nix development shell and exit.
  -h, --help           Show this help message and exit.

Examples:
  .gemini/skills/run-in-nix.sh --list-tools
  .gemini/skills/run-in-nix.sh "terraform validate"
  .gemini/skills/run-in-nix.sh "go test ./..."
EOF
}

query_nix_tools() {
  echo "Querying standard Nix environment for installed tools..." >&2

  # Execute a inline script inside Nix shell to find the dev-shell-package/bin directory on the PATH
  local bin_dir
  if [[ -n "${IN_NIX_SHELL:-}" ]]; then
    bin_dir=$(echo "$PATH" | tr ":" "\n" | grep "dev-shell-package/bin" | head -n 1 || echo "")
  else
    # shellcheck disable=SC2016
    bin_dir=$(nix develop \
      --extra-experimental-features nix-command \
      --extra-experimental-features flakes \
      --command bash -c 'echo "$PATH" | tr ":" "\n" | grep "dev-shell-package/bin" | head -n 1' 2>/dev/null || echo "")
  fi

  if [[ -z "$bin_dir" || ! -d "$bin_dir" ]]; then
    echo "Error: Could not locate dev-shell-package/bin in the Nix environment." >&2
    exit 1
  fi

  echo "Installed Nix Tools:"
  echo "===================="
  # List all files in the symlink directory, filtering out internal dot-prefixed executables
  find "$bin_dir" -maxdepth 1 -type l -not -name ".*" -exec basename {} \; | sort
}

execute_in_nix() {
  local command="$1"

  if [[ -n "${IN_NIX_SHELL:-}" ]]; then
    echo "Already in a Nix shell environment (IN_NIX_SHELL=${IN_NIX_SHELL}). Executing command directly..." >&2
    bash -c "${command}"
    return
  fi

  echo "Running command in Nix environment: ${command}" >&2

  nix develop \
    --extra-experimental-features nix-command \
    --extra-experimental-features flakes \
    --command bash -c "${command}"
}

main() {
  if [[ $# -eq 0 ]]; then
    echo "Error: Command required." >&2
    show_help
    exit 1
  fi

  if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
    exit 0
  fi

  if [[ "$1" == "--list-tools" ]]; then
    query_nix_tools
    exit 0
  fi

  execute_in_nix "$1"
}

main "$@"
