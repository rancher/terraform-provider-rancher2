#!/usr/bin/env bash
#
# Description: Verifies staging file limits and overrides.
#

# Verify staging files and file-count limits
verify_staging_limits() {
  local max_allowed=5
  if [[ -n "${COMMIT_LIMIT_OVERRIDE:-}" ]]; then
    if [[ ! "${COMMIT_LIMIT_OVERRIDE}" =~ ^[0-9]+$ ]]; then
      echo "Error: COMMIT_LIMIT_OVERRIDE must be a positive integer, got: '${COMMIT_LIMIT_OVERRIDE}'" >&2
      exit 1
    fi
    max_allowed="${COMMIT_LIMIT_OVERRIDE}"
    echo "--> [OVERRIDE] Using custom staged file limit from COMMIT_LIMIT_OVERRIDE: ${max_allowed}" >&2
  fi

  local staged_count
  staged_count=$(git diff --cached --name-only | wc -l | tr -d ' ')

  if [[ "$staged_count" -eq 0 ]]; then
    if [[ -f .git/MERGE_HEAD || -f .git/CHERRY_PICK_HEAD || -f .git/REBASE_HEAD ]]; then
      echo "--> [MERGE STATE] Active merge/rebase/cherry-pick in progress. Allowing 0 staged files to create merge commit."
      return 0
    fi
    echo "Error: No changes are currently staged for commit." >&2
    echo "       Please stage your changes first using 'git add <files>...'." >&2
    exit 1
  fi

  if [[ "$staged_count" -gt "$max_allowed" ]]; then
    echo "Error: Committing too much code at once is prohibited ($staged_count files staged; max allowed is $max_allowed)." >&2
    echo "       In accordance with Phase 5, Step 11 of 'docs/development/AgenticFramework/DevelopmentProcess.md', please split your commit into smaller, surgical layers." >&2
    exit 1
  fi
}
