#!/usr/bin/env bash
#
# Skill: commit-push.sh
# Description: Programmatically commit and push local changes with GPG/SSH signature, sign-off, and fork synchronization.
#              This acts as a lightweight controller, importing and coordinating modular operations under agent-scripts/.
# Conforms to shell-scripts.instructions.md guidelines.

set -euo pipefail

# Locate absolute directory path of the workspace root dynamically
WORKSPACE_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

# 1. Source modular utilities from agent-scripts/
# shellcheck source=agent-scripts/git-utils.sh
source "${WORKSPACE_ROOT}/agent-scripts/git-utils.sh"
# shellcheck source=agent-scripts/check-branch.sh
source "${WORKSPACE_ROOT}/agent-scripts/check-branch.sh"
# shellcheck source=agent-scripts/verify-limits.sh
source "${WORKSPACE_ROOT}/agent-scripts/verify-limits.sh"
# shellcheck source=agent-scripts/verify-gates.sh
source "${WORKSPACE_ROOT}/agent-scripts/verify-gates.sh"
# shellcheck source=agent-scripts/commit-helper.sh
source "${WORKSPACE_ROOT}/agent-scripts/commit-helper.sh"
# shellcheck source=agent-scripts/push-helper.sh
source "${WORKSPACE_ROOT}/agent-scripts/push-helper.sh"

# ==============================================================================
# GLOBAL VARIABLES
# ==============================================================================
COMMIT_MSG=""
FORCE_PUSH=false

show_help() {
  cat <<EOF
Usage: commit-push.sh [options] -m "COMMIT_MESSAGE"

Programmatically commit and push local changes with GPG/SSH signature, sign-off, and fork synchronization.

Options:
  -h, --help            Show this message and exit.
  -m MESSAGE            The conventional commit message (Required).
  -f, --force           Bypass remote ancestry check and perform safe force-push with lease.
EOF
}

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      -f|--force)
        FORCE_PUSH=true
        shift
        ;;
      -m)
        if [[ -z "${2:-}" ]]; then
          echo "Error: -m option requires a non-empty commit message argument." >&2
          exit 1
        fi
        COMMIT_MSG="$2"
        shift 2
        ;;
      *)
        echo "Error: Unknown argument '$1'" >&2
        show_help >&2
        exit 1
        ;;
    esac
  done

  if [[ -z "$COMMIT_MSG" ]]; then
    echo "Error: Commit message is required. Specify using -m \"message\"." >&2
    show_help >&2
    exit 1
  fi
}

# Sync with Upstream parent repository
sync_default_branch() {
  local branch="$1"
  if [[ "$branch" != "main" ]]; then
    if [[ -f .git/MERGE_HEAD || -f .git/CHERRY_PICK_HEAD || -f .git/REBASE_HEAD ]]; then
      echo "--> [MERGE STATE] Active merge/rebase/cherry-pick in progress. Skipping sync_default_branch to preserve merge state."
      return 0
    fi
    echo "Synchronizing local 'main' branch and tags with upstream parent repository..."
    local stash_created=false
    if git status --porcelain | grep -v '^[A-Z]' >/dev/null; then
      echo "  -> Temporarily stashing unstaged/untracked files..."
      git stash push -k -u -m "temp-commit-push-stash" >/dev/null
      stash_created=true
    fi

    # Run sync skill
    if ! bash "${WORKSPACE_ROOT}/.gemini/skills/git-sync.sh"; then
      echo "Error: Upstream synchronization failed." >&2
      if [[ "$stash_created" == "true" ]]; then
        if ! git stash pop --index >/dev/null 2>&1; then
          echo "Warning: Stash pop failed during emergency exit. Your stashed changes remain preserved in Git stash." >&2
        fi
      fi
      exit 1
    fi

    echo "Switching back to branch '$branch'..."
    if ! git checkout "$branch" >/dev/null 2>&1; then
      echo "Error: Failed to switch back to branch '$branch' after sync." >&2
      if [[ "$stash_created" == "true" ]]; then
        if ! git stash pop --index >/dev/null 2>&1; then
          echo "Warning: Stash pop failed during emergency exit. Your stashed changes remain preserved in Git stash." >&2
        fi
      fi
      exit 1
    fi

    if [[ "$stash_created" == "true" ]]; then
      echo "  -> Restoring stashed unstaged/untracked files..."
      if ! git stash pop --index >/dev/null 2>&1; then
        echo "Warning: Re-applying stashed changes resulted in merge conflicts." >&2
        echo "         Your stashed changes have been PRESERVED in the Git stash list." >&2
      fi
    fi
  fi
}

main() {
  parse_args "$@"

  local active_branch
  active_branch=$(git branch --show-current)

  # 1. Defunct branch protection check
  check_defunct_branch "$active_branch"

  # 2. Gate 3 (Proactive Review) validation on disk
  verify_proactive_review

  # 3. Stage & limit checks
  # Stage files first so we can verify limits accurately
  git add -A
  verify_staging_limits

  # 4. Sync up with parent repository if we are not on main
  sync_default_branch "$active_branch"

  # 5. Remote ancestry validation (unless force is requested)
  if [[ "$FORCE_PUSH" == "false" ]]; then
    verify_remote_ancestry "$active_branch"
  fi

  # 6. Conventional commit execution (adds staging, sign-off and signature)
  execute_commit "$COMMIT_MSG" "$active_branch"

  # 7. Safe Push with lease
  execute_push "origin" "$active_branch" "$FORCE_PUSH"
}

main "$@"
