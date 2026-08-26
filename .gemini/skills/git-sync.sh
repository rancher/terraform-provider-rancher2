#!/usr/bin/env bash
#
# Skill: git-sync.sh
# Description: Safely syncs the default branch, tags, and optionally the current branch of a local fork with the upstream parent repository.
#              Enforces strict security checks to ensure we never push to any Rancher-owned remote/repository.
# Usage: .gemini/skills/git-sync.sh [stay]
#        If 'stay' (or any non-empty value) is provided, it will also sync and checkout the current branch.

set -euo pipefail

# ==============================================================================
# CORE FUNCTIONS
# ==============================================================================

# Global variables to track state
AUTO_STASH_CREATED=false
STAY_STASH_CREATED=false
ORIGINAL_BRANCH=""
ORIGINAL_BRANCH_COMMIT=""

verify_git_env() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: This command must be run inside a Git repository." >&2
    exit 1
  fi

  # Securely auto-stash uncommitted or untracked changes rather than erroring
  if [[ "$(git status --porcelain=v1 2>/dev/null | wc -l)" -gt 0 ]]; then
    echo "  -> Uncommitted or untracked changes detected in local tree."
    echo "  -> Securely stashing changes to guarantee zero data loss..."
    if ! git stash push -u -m "git-sync-auto-stash" >/dev/null; then
      echo "Error: Failed to stash changes cleanly. Aborting sync." >&2
      exit 1
    fi
    AUTO_STASH_CREATED=true
  fi
}

get_origin_url() {
  local url
  url=$(git remote get-url origin 2>/dev/null || true)
  if [[ -z "$url" ]]; then
    echo "Error: Could not retrieve configured origin remote URL." >&2
    echo "       Please ensure 'origin' is set up." >&2
    exit 1
  fi
  echo "$url"
}

verify_push_safety() {
  local remote_name="$1"
  local url
  url=$(git remote get-url "$remote_name" 2>/dev/null || true)
  if [[ -z "$url" ]]; then
    echo "Error: Remote '$remote_name' has no configured URL." >&2
    exit 1
  fi
  if [[ "$url" =~ [/:](rancher|rancherlabs)/ ]]; then
    echo "======================================================================" >&2
    echo "❌ CRITICAL SECURITY ERROR: UNSAFE PUSH PREVENTED!" >&2
    echo "   The remote '$remote_name' points to a Rancher-owned repository:" >&2
    echo "   $url" >&2
    echo "   Pushing directly to upstream Rancher repositories is strictly forbidden." >&2
    echo "======================================================================" >&2
    exit 1
  fi
}

safe_git_push() {
  local remote=""
  local arg
  for arg in "$@"; do
    if [[ ! "$arg" =~ ^- ]]; then
      remote="$arg"
      break
    fi
  done

  if [[ -z "$remote" ]]; then
    echo "Error: Could not determine remote from git push arguments." >&2
    exit 1
  fi

  verify_push_safety "$remote"
  git push "$@"
}

get_origin_owner() {
  local origin_url="$1"
  if [[ "$origin_url" =~ github\.com[:/]([^/]+)/ ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    echo ""
  fi
}

get_upstream_owner() {
  # The parent repository owner is always "rancher"
  echo "rancher"
}

get_upstream_repo() {
  # Parse repository name from origin remote url
  local origin_url
  origin_url=$(git remote get-url origin 2>/dev/null || true)
  
  if [[ "$origin_url" =~ github\.com[:/][^/]+/([^/]+)\.git ]]; then
    echo "${BASH_REMATCH[1]}"
  elif [[ "$origin_url" =~ github\.com[:/][^/]+/([^/]+) ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    # Fallback to the local directory name
    basename "$(git rev-parse --show-toplevel)"
  fi
}

get_default_branch() {
  local default_branch=""
  default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || true)
  if [[ -z "$default_branch" ]]; then
    default_branch=$(git remote show origin 2>/dev/null | grep 'HEAD branch' | cut -d' ' -f5 || true)
  fi
  if [[ -z "$default_branch" ]]; then
    echo "Error: Could not determine the default branch for the origin remote." >&2
    exit 1
  fi
  echo "$default_branch"
}

get_upstream_url() {
  local origin_url="$1"
  local owner="$2"
  local repo="$3"
  if [[ "$origin_url" =~ ^git@ ]]; then
    echo "git@github.com:${owner}/${repo}.git"
  else
    echo "https://github.com/${owner}/${repo}.git"
  fi
}

sync_branch() {
  local branch="$1"
  local upstream_url="$2"
  local is_default="${3:-false}"

  echo "Syncing branch '$branch' with upstream..."
  git checkout "$branch"
  git remote rm upstream 2>/dev/null || true
  git remote add upstream "$upstream_url"
  git fetch --all
  git pull --rebase upstream "$branch"
  
  if [[ "$is_default" == "true" ]]; then
    safe_git_push --tags origin
  fi
  
  git remote rm upstream
  safe_git_push -f origin "$branch"
}

# ==============================================================================
# MAIN ORCHESTRATION
# ==============================================================================

cleanup() {
  local exit_code=$?
  if [[ $exit_code -ne 0 ]]; then
    echo "⚠️ Script interrupted or failed. Running cleanup..." >&2
    git remote rm upstream 2>/dev/null || true
    
    if [[ -n "$ORIGINAL_BRANCH" ]]; then
      echo "Restoring original branch '$ORIGINAL_BRANCH'..." >&2
      git checkout "$ORIGINAL_BRANCH" 2>/dev/null || true
      
      if [[ -n "$ORIGINAL_BRANCH_COMMIT" ]]; then
        echo "Resetting '$ORIGINAL_BRANCH' back to its original commit: $ORIGINAL_BRANCH_COMMIT..." >&2
        git reset --hard "$ORIGINAL_BRANCH_COMMIT" >/dev/null 2>&1 || true
      fi
    fi

    if [[ "$STAY_STASH_CREATED" == "true" ]]; then
      echo "⚠️ WARNING: A temporary stay stash ('git-sync-stay-stash') was created containing your in-progress work!" >&2
      echo "   The stash has been PRESERVED in your Git stash list to prevent data loss." >&2
      echo "   Please inspect/restore it manually using 'git stash list' and 'git stash pop' after recovery." >&2
    fi
  fi

  # Always restore the auto-stash on exit if it was successfully created
  if [[ "$AUTO_STASH_CREATED" == "true" ]]; then
    if [[ -n "${ORIGINAL_BRANCH:-}" && "$(git branch --show-current)" != "$ORIGINAL_BRANCH" ]]; then
      echo "Switching back to original branch '$ORIGINAL_BRANCH' before restoring stash..." >&2
      git checkout "$ORIGINAL_BRANCH" >/dev/null 2>&1 || true
    fi
    echo "  -> Restoring stashed changes from auto-stash..." >&2
    if ! git stash pop --index >/dev/null 2>&1; then
      echo "Warning: Re-applying stashed changes resulted in merge conflicts." >&2
      echo "         Your changes have been safely preserved in the Git stash." >&2
    fi
    AUTO_STASH_CREATED=false
  fi

  exit "$exit_code"
}

show_help() {
  cat <<EOF
Usage: git-sync.sh [options] [stay]

Safely syncs the default branch, tags, and optionally the current branch of a local fork with the upstream parent repository.

Options:
  stay                 Sync and checkout the current branch in addition to the default branch.
  -h, --help           Show this help message and exit.

Examples:
  .gemini/skills/git-sync.sh
  .gemini/skills/git-sync.sh stay
EOF
}

main() {
  trap cleanup EXIT

  local stay=""

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      stay)
        stay="true"
        shift
        ;;
      -*)
        echo "Error: Unknown option: $1" >&2
        show_help
        exit 1
        ;;
      *)
        echo "Error: Unexpected argument: $1" >&2
        show_help
        exit 1
        ;;
    esac
  done

  verify_git_env

  ORIGINAL_BRANCH="$(git branch --show-current)"
  ORIGINAL_BRANCH_COMMIT="$(git rev-parse HEAD 2>/dev/null || true)"
  echo "Current branch is $ORIGINAL_BRANCH..."

  local origin_url
  origin_url=$(get_origin_url)

  echo "Running push safety verification on origin..."
  verify_push_safety origin

  local origin_owner
  origin_owner=$(get_origin_owner "$origin_url")
  echo "Found origin owner: \"$origin_owner\"..."

  local upstream_owner
  upstream_owner=$(get_upstream_owner)
  echo "Found upstream owner: $upstream_owner..."

  local upstream_repo
  upstream_repo=$(get_upstream_repo)
  echo "Found upstream repo name: $upstream_repo..."

  if [[ "$origin_owner" == "$upstream_owner" ]]; then
    echo "Origin is already the upstream repository ($upstream_owner), nothing to sync."
    exit 0
  fi

  local default_branch
  default_branch=$(get_default_branch)
  echo "Found default branch as: $default_branch..."

  local upstream_url
  upstream_url=$(get_upstream_url "$origin_url" "$upstream_owner" "$upstream_repo")
  echo "Configuring upstream URL: $upstream_url"

  # Clear any legacy upstream remotes
  git remote rm upstream 2>/dev/null || true

  # Expert 6-Step Stay-Sync Workflow
  local has_changes="false"
  if [[ -n "$stay" ]]; then
    echo "User requested to stay on current branch. Implementing expert 6-step sync workflow..."
    
    # 1. set history back to when the branch was first merged from main (merge-base)
    echo "  -> Step 1: Finding merge base between '$ORIGINAL_BRANCH' and '$default_branch'..."
    local merge_base
    merge_base=$(git merge-base "$default_branch" "$ORIGINAL_BRANCH" 2>/dev/null || true)
    if [[ -z "$merge_base" ]]; then
      echo "Error: Could not determine merge base between '$ORIGINAL_BRANCH' and '$default_branch'." >&2
      exit 1
    fi

    if ! git diff --quiet "$merge_base" "$ORIGINAL_BRANCH"; then
      has_changes="true"
      echo "  -> Resetting branch history soft-wise to merge base: $merge_base..."
      git checkout "$ORIGINAL_BRANCH"
      git reset --soft "$merge_base"
      
      # 2. add and stash the changes
      echo "  -> Step 2: Adding and stashing changes securely..."
      git add -A
      git stash push -u -m "git-sync-stay-stash"
      STAY_STASH_CREATED="true"
    else
      echo "  -> No unique changes detected on '$ORIGINAL_BRANCH' compared to merge base."
    fi
  fi

  # 3. checkout main and sync it using the git-sync skill (default branch sync)
  echo "  -> Step 3: Checking out and syncing default branch '$default_branch'..."
  git checkout "$default_branch"
  git reset --hard
  sync_branch "$default_branch" "$upstream_url" true

  # 4 & 5 & 6. switch back to branch, merge main, pop stash, and force push
  if [[ -n "$stay" ]]; then
    # 4. switch back to the branch and merge main into it
    echo "  -> Step 4: Switching back to original branch '$ORIGINAL_BRANCH' and merging '$default_branch'..."
    git checkout "$ORIGINAL_BRANCH"
    git merge "$default_branch"
    
    # 5. pop the stash and resolve any conflicts
    if [[ "$has_changes" == "true" ]]; then
      echo "  -> Step 5: Restoring stashed branch changes..."
      STAY_STASH_CREATED="false" # Popping now, so don't drop on failure
      if ! git stash pop; then
        echo "⚠️ WARNING: Merge conflicts detected during stash pop!" >&2
        echo "   Please resolve conflicts manually in your IDE, then run:" >&2
        echo "   git push -f origin $ORIGINAL_BRANCH" >&2
        ORIGINAL_BRANCH_COMMIT="" # Disable destructive rollback
        exit 1
      fi
    else
      echo "  -> Step 5: Skipping stash pop (no changes were stashed)."
    fi
    
    # 6. force push the change to origin
    echo "  -> Step 6: Force-pushing finalized branch changes to origin..."
    safe_git_push -f origin "$ORIGINAL_BRANCH"
  fi

  echo "✅ Git sync completed successfully!"
}

main "$@"
