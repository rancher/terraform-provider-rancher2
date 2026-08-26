#!/usr/bin/env bash
#
# Description: Verifies branch status and remote ancestry safety.
#

# Check if the current branch has an already merged PR on GitHub (Branch Defunct Protection)
check_defunct_branch() {
  local branch="$1"
  if [[ "$branch" != "main" ]]; then
    local pr_status
    if pr_status=$(GITHUB_TOKEN="" gh pr view "$branch" --json state,number --template '{{.state}} {{.number}}' 2>/dev/null); then
      local pr_state
      pr_state=$(echo "$pr_status" | cut -d' ' -f1)
      local pr_number
      pr_number=$(echo "$pr_status" | cut -d' ' -f2)

      if [[ "$pr_state" == "MERGED" ]]; then
        echo "Error: The current branch '$branch' already has a merged Pull Request (#$pr_number) on GitHub." >&2
        echo "       This branch is defunct. In accordance with 'docs/development/AgenticFramework/DevelopmentProcess.md' Phase 5, Step 12, you MUST:" >&2
        echo "       1. Switch to 'main': git checkout main" >&2
        echo "       2. Synchronize with upstream default branch: bash .gemini/skills/git-sync.sh" >&2
        echo "       3. Check out a clean, new branch off updated main: git checkout -b feature/workflows-new-branch" >&2
        exit 1
      fi
    fi
  fi
}

# Verify ancestry check to fail fast if we are behind remote
verify_remote_ancestry() {
  local branch="$1"
  local remote_name="origin"

  echo "Verifying local branch ancestry is fully up-to-date with remote fork..." >&2

  # Fetch the remote tracking branch safely
  if ! git fetch -q "$remote_name" "$branch" >/dev/null 2>&1; then
    echo "--> [FETCH SKIPPED] Tracking reference on remote '$remote_name/$branch' does not exist yet. Safe to proceed." >&2
    return 0
  fi

  local local_sha
  local_sha=$(git rev-parse HEAD)
  local remote_sha
  remote_sha=$(git rev-parse "$remote_name/$branch" 2>/dev/null || echo "")

  if [[ -z "$remote_sha" ]]; then
    echo "--> [RESOLVE SKIPPED] No remote tracking branch SHA found. Safe to proceed." >&2
    return 0
  fi

  if [[ "$local_sha" == "$remote_sha" ]]; then
    echo "✅ Local branch is identical to remote fork tracking reference." >&2
    return 0
  fi

  # Check if remote reference is an ancestor of local HEAD
  if git merge-base --is-ancestor "$remote_sha" "$local_sha" >/dev/null 2>&1; then
    echo "✅ Local branch contains all remote tracking changes (Fast-Forward ancestry confirmed)." >&2
    return 0
  fi

  echo "======================================================================" >&2
  echo "❌ CRITICAL FAILURE: LOCAL BRANCH IS OUT-OF-SYNC WITH REMOTE FORK!" >&2
  echo "   Remote tracking reference '$remote_name/$branch' ($remote_sha) has changes" >&2
  echo "   that are not present in your local branch. Direct pushing is unsafe." >&2
  echo "   To resolve, please execute our synchronized branch rebase:" >&2
  echo "   1. Stash any unstaged changes: git stash" >&2
  echo "   2. Rebase HEAD onto remote tracking ref: git rebase $remote_name/$branch" >&2
  echo "   3. Restore stashed changes: git stash pop" >&2
  echo "======================================================================" >&2
  exit 1
}
