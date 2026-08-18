#!/usr/bin/env bash
#
# Skill: commit-push.sh
# Description: Programmatically commit and push local changes with GPG/SSH signature, sign-off, and fork synchronization.
# Conforms to shell-scripts.instructions.md guidelines.

set -euo pipefail

# Determine repository root directory
REPO_ROOT=""
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  REPO_ROOT="$(git rev-parse --show-toplevel)"
fi

# ==============================================================================
# HELPER FUNCTIONS
# ==============================================================================

# Helper to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Display script help usage instructions
show_help() {
  cat <<EOF
Usage: commit-push.sh [options] -m "COMMIT_MESSAGE"

Programmatically commit and push local changes with GPG/SSH signature, sign-off, and fork synchronization.

Options:
  -h, --help            Show this message and exit.
  -m MESSAGE            The conventional commit message (Required).
  -f, --force           Bypass remote ancestry check and perform safe force-push with lease.

Examples:
  .agent/skills/commit-push.sh -m "ci(workflows): add new automated checks"
  .agent/skills/commit-push.sh -f -m "refactor(hooks): force push after rebase"
EOF
}

verify_environment() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: Must be run inside a Git repository." >&2
    exit 1
  fi

  if ! command_exists gh; then
    echo "Error: 'gh' (GitHub CLI) is required to run this script. Please install gh." >&2
    exit 1
  fi

  if ! gh auth status >/dev/null 2>&1; then
    echo "Error: GitHub CLI is not authenticated. Please run 'gh auth login' to authenticate." >&2
    exit 1
  fi
}

get_file_owner_uid() {
  local file="$1"
  stat -c %u "$file" 2>/dev/null || stat -f %u "$file" 2>/dev/null || echo ""
}

calculate_sha256() {
  if command_exists shasum; then
    shasum -a 256 | cut -d' ' -f1
  elif command_exists sha256sum; then
    sha256sum | cut -d' ' -f1
  else
    echo "Error: No SHA-256 utility (shasum or sha256sum) found on this system." >&2
    exit 1
  fi
}

verify_push_safety() {
  local remote_name="$1"
  local url
  url=$(git remote get-url "$remote_name" 2>/dev/null || true)
  if [[ -z "$url" ]]; then
    echo "Error: Remote '$remote_name' has no configured URL." >&2
    exit 1
  fi
  # Compare against a lowercased copy — bash's [[ =~ ]] is case-sensitive by default, and
  # GitHub org names aren't (github.com/Rancher/... would otherwise bypass this check).
  local url_lower="${url,,}"
  if [[ "$url_lower" =~ [/:](rancher|rancherlabs)/ ]]; then
    echo "======================================================================" >&2
    echo "❌ CRITICAL SECURITY ERROR: UNSAFE PUSH PREVENTED!" >&2
    echo "   The remote '$remote_name' points to a Rancher-owned repository:" >&2
    echo "   $url" >&2
    echo "   Pushing directly to upstream Rancher repositories is strictly forbidden." >&2
    echo "======================================================================" >&2
    exit 1
  fi
}

# ==============================================================================
# OPERATION STAGES
# ==============================================================================

# Parse options and arguments
parse_args() {
  # Initialize variables with global defaults
  COMMIT_MSG=""
  FORCE_PUSH=false

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

# Check if the current branch has an already merged PR on GitHub (Branch Defunct Protection)
check_defunct_branch() {
  local branch="$1"
  if [[ "$branch" != "main" ]]; then
    local pr_status
    if pr_status=$(gh pr view "$branch" --json state,number --template '{{.state}} {{.number}}' 2>/dev/null); then
      local pr_state
      pr_state=$(echo "$pr_status" | cut -d' ' -f1)
      local pr_number
      pr_number=$(echo "$pr_status" | cut -d' ' -f2)
      
      if [[ "$pr_state" == "MERGED" ]]; then
        echo "Error: The current branch '$branch' already has a merged Pull Request (#$pr_number) on GitHub." >&2
        echo "       This branch is defunct. In accordance with 'development-process.md' Phase 5, Step 12, you MUST:" >&2
        echo "       1. Switch to 'main': git checkout main" >&2
        echo "       2. Synchronize with upstream default branch: bash .agent/skills/git-sync.sh" >&2
        echo "       3. Check out a clean, new branch off updated main: git checkout -b feature/workflows-new-branch" >&2
        exit 1
      fi
    fi
  fi
}

# Verify staging files and file-count limits
verify_staging_limits() {
  local max_allowed=5
  STAGED_COUNT=$(git diff --cached --name-only | wc -l | tr -d ' ')

  if [[ "$STAGED_COUNT" -eq 0 ]]; then
    echo "Error: No changes are currently staged for commit." >&2
    echo "       Please stage your changes first using 'git add <files>...'." >&2
    exit 1
  fi

  if [[ "$STAGED_COUNT" -gt "$max_allowed" ]]; then
    echo "Error: Committing too much code at once is prohibited ($STAGED_COUNT files staged; max allowed is $max_allowed)." >&2
    echo "       In accordance with Phase 5, Step 11 of 'development-process.md', please split your commit into smaller, surgical layers." >&2
    exit 1
  fi
}

# Enforce secure proactive review validation
verify_proactive_review() {
  # Delegate verification cleanly and securely to write-approval.sh skill
  if ! bash "$REPO_ROOT/.agent/skills/write-approval.sh" --verify; then
    exit 1
  fi
}

# Sync with Upstream parent repository
sync_default_branch() {
  local branch="$1"
  if [[ "$branch" != "main" ]]; then
    echo "Synchronizing local 'main' branch and tags with upstream parent repository..."
    # NOTE: Do NOT stash before calling git-sync.sh. That skill already auto-stashes
    # uncommitted AND untracked files itself (see its verify_git_env) and restores them
    # with 'git stash pop --index', so the staged index survives. Stashing here as well
    # is redundant and actively harmful: much of '.agent/skills/' is untracked, so a
    # 'git stash push -u' sweeps away git-sync.sh itself before this line can invoke it,
    # failing with "No such file or directory".
    if ! bash "$REPO_ROOT/.agent/skills/git-sync.sh"; then
      echo "Error: Upstream synchronization failed." >&2
      exit 1
    fi

    # Switch back to the active feature branch in case git-sync.sh left the checkout on main
    echo "Switching back to branch '$branch'..."
    if ! git checkout "$branch" >/dev/null 2>&1; then
      echo "Error: Failed to switch back to branch '$branch' after sync." >&2
      exit 1
    fi
  fi
}

# Verify ancestry check to fail fast if we are behind remote
verify_remote_ancestry() {
  local branch="$1"
  if [[ "$FORCE_PUSH" == "true" ]]; then
    echo "Force-push option specified. Skipping ancestry check."
  else
    echo "Checking remote branch status on origin..."
    # Check existence separately from fetching it — treating a fetch failure as "branch
    # doesn't exist" would also silently swallow real failures (network/auth issues) and let
    # the push proceed without ever checking whether we're behind.
    if git ls-remote --exit-code origin "$branch" >/dev/null 2>&1; then
      if ! git fetch origin "$branch" >/dev/null 2>&1; then
        echo "Error: Remote branch 'origin/$branch' exists but fetching it failed (network or authentication issue?)." >&2
        echo "       Please resolve connectivity/authentication and try again." >&2
        exit 1
      fi
      local behind_count
      behind_count=$(git rev-list --count "HEAD..origin/$branch" 2>/dev/null || echo "0")
      if [[ "$behind_count" -gt 0 ]]; then
        echo "Error: Your local branch is behind 'origin/$branch' by $behind_count commit(s)." >&2
        echo "       Please pull and integrate the remote changes before pushing." >&2
        exit 1
      fi
      echo "  -> Local branch is up to date with remote."
    else
      echo "  -> Remote branch 'origin/$branch' does not exist yet. Safe to proceed."
    fi
  fi
}

# Verify developer manual IDE review approval
verify_developer_approval() {
  # If a valid user-approval signature is already present, verify it cleanly
  if node "$REPO_ROOT/.agent/skills/user-approval.js" --verify >/dev/null 2>&1; then
    echo "✅ Developer visual IDE review approval verified!" >&2
    return 0
  fi

  # Otherwise, programmatically prompt the developer for approval now
  echo "No prior developer approval signature found on disk." >&2
  if ! node "$REPO_ROOT/.agent/skills/user-approval.js" "Do you approve GPG-signing, committing, and pushing these changes?"; then
    echo "❌ Commit and push aborted by developer." >&2
    exit 1
  fi
}

# Execute signed and signed-off git commit
execute_signed_commit() {
  echo "Committing $STAGED_COUNT staged file(s) with signature (-S) and sign-off (-s)..."
  if ! git commit -s -S -m "$COMMIT_MSG"; then
    echo "Error: Git commit failed. Ensure GPG/SSH signing is configured." >&2
    exit 1
  fi
}

# Perform secure push to fork remote
secure_push() {
  local branch="$1"
  if [[ "$FORCE_PUSH" == "true" ]]; then
    echo "Pushing signed commit with FORCE securely to fork remote 'origin/$branch'..."
    if ! git push origin "$branch" --force-with-lease; then
      echo "Error: Git push failed." >&2
      exit 1
    fi
  else
    echo "Pushing signed commit securely to fork remote 'origin/$branch'..."
    if ! git push origin "$branch"; then
      echo "Error: Git push failed." >&2
      exit 1
    fi
  fi
}

# ==============================================================================
# MAIN ENTRY POINT
# ==============================================================================
main() {
  # Global state variables
  STAGED_COUNT=0

  parse_args "$@"
  verify_environment
  verify_push_safety origin

  local current_branch
  current_branch=$(git branch --show-current)
  if [[ -z "$current_branch" ]]; then
    echo "Error: Could not determine current branch name." >&2
    exit 1
  fi

  check_defunct_branch "$current_branch"
  verify_staging_limits
  verify_proactive_review
  sync_default_branch "$current_branch"
  verify_remote_ancestry "$current_branch"
  verify_developer_approval
  execute_signed_commit
  secure_push "$current_branch"

  echo "✅ Changes programmatically committed and pushed successfully!"
}

main "$@"