#!/usr/bin/env bash
#
# Skill: run-in-nix.sh
# Description: Executes a given command inside the standardized Nix development environment, or lists installed tools.
# Usage: .agent/skills/run-in-nix.sh "<command>"

set -euo pipefail

show_help() {
  cat <<EOF
Usage: run-in-nix.sh [options] [command...]

Executes a given command inside the standardized Nix development environment, or queries installed tools.

Options:
  --list-tools         List all executable tools installed in the Nix development shell and exit.
  -h, --help           Show this help message and exit.

Examples:
  .agent/skills/run-in-nix.sh --list-tools
  .agent/skills/run-in-nix.sh terraform validate
  .agent/skills/run-in-nix.sh "go test ./..."
EOF
}

query_nix_tools() {
  echo "Querying standard Nix environment for installed tools..." >&2

  # Execute an inline script inside Nix shell to find the dev-shell-package/bin directory on the PATH
  local bin_dir
  if [[ -n "${IN_NIX_SHELL:-}" ]]; then
    bin_dir=$(echo "$PATH" | tr ":" "\n" | grep "dev-shell-package/bin" | head -n 1 || echo "")
  else
    # Check if nix is installed
    if ! command -v nix >/dev/null 2>&1; then
      echo "Error: 'nix' executable is required but was not found on the system." >&2
      exit 1
    fi
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

  # Check if nix is installed
  if ! command -v nix >/dev/null 2>&1; then
    echo "Error: 'nix' executable is required but was not found on the system." >&2
    exit 1
  fi

  echo "Running command in Nix environment: ${command}" >&2

  nix develop \
    --extra-experimental-features nix-command \
    --extra-experimental-features flakes \
    --command bash -c "${command}"
}

main() {
  local list_tools=false

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --list-tools)
        list_tools=true
        shift
        ;;
      *)
        break
        ;;
    esac
  done

  if [[ "$list_tools" == "true" ]]; then
    query_nix_tools
    exit 0
  fi

  if [[ $# -eq 0 ]]; then
    echo "Error: Command required." >&2
    show_help
    exit 1
  fi

  # Join arguments with spaces to support unquoted commands
  local command="$*"
  execute_in_nix "${command}"
}

main "$@"
