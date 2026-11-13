#!/usr/bin/env bash
# Wrapper script to exercise models in the correct environment

set -euo pipefail

show_help() {
  cat <<EOF
Usage: exercise-cron.sh [options]

Wrapper script to exercise models in the correct environment for headless cron execution.

Options:
  -h, --help    Show this help message and exit
EOF
}

# parse_args handles the command-line argument contract.
# Note: This is intentionally simplified to support a single-parameter contract
# (specifically checking for -h/--help) and does not support multi-argument shifting.
parse_args() {
  case "${1:-}" in
    -h | --help)
      show_help
      exit 0
      ;;
    "")
      ;;
    *)
      echo "Error: Unknown argument: $1" >&2
      show_help
      exit 1
      ;;
  esac
}

run_exercise() {
  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
  local project_dir
  project_dir="$(cd "${script_dir}/../.." >/dev/null 2>&1 && pwd)"

  cd "${project_dir}"

  # Source environment variables for headless cron execution
  if [[ -f ".variables" ]]; then
    # shellcheck disable=SC1091,SC1090
    source .variables
  fi
  if [[ -f "${HOME}/.config/gemini/default" ]]; then
    # shellcheck disable=SC1091,SC1090
    source "${HOME}/.config/gemini/default"
  fi

  # Execute exercise-agents.js directly
  node "${script_dir}/exercise-agents.js"
}

main() {
  parse_args "$@"
  run_exercise
}

main "$@"
