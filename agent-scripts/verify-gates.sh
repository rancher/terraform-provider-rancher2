#!/usr/bin/env bash
#
# Description: Verifies Gate 3 proactive review status and hashes.
#

# Enforce secure proactive review validation
verify_proactive_review() {
  local target_dir
  local repo_name
  repo_name=$(basename "$(git rev-parse --show-toplevel 2>/dev/null || pwd)")
  target_dir="${AGENT_STATE_DIR:-${HOME}/.gemini/tmp/${repo_name}}"
  local review_file="${target_dir}/review-approval.json"

  echo "Verifying proactive review approval status..." >&2

  if [[ ! -f "$review_file" ]]; then
    echo "Error: Proactive review approval file not found!" >&2
    echo "       In accordance with Gate 3 (Review Gate) of 'docs/development/AgenticFramework/DevelopmentProcess.md'," >&2
    echo "       you MUST run the review agent first: @review_agent" >&2
    exit 1
  fi

  # Reject symbolic links
  if [[ -L "$review_file" ]]; then
    echo "Error: Proactive review approval file is a symbolic link (Prohibited)." >&2
    exit 1
  fi

  # Verify ownership natively
  local file_uid=""
  file_uid=$(get_file_owner_uid "$review_file")

  if [[ -z "$file_uid" ]]; then
    echo "Error: Could not determine owner UID for proactive review approval file." >&2
    exit 1
  fi

  if [[ "$file_uid" -ne "$(id -u)" ]]; then
    echo "Error: Proactive review approval file is not owned by the current user (UID: $(id -u), Owner: $file_uid)." >&2
    exit 1
  fi

  # Ensure jq utility is present
  if ! command_exists jq; then
    echo "Error: jq utility is required but not found on this system." >&2
    exit 1
  fi

  # Check diff_hash inside the JSON using jq and grep (since main on feature branches)
  local active_hash=""
  local current_branch=""
  current_branch=$(git branch --show-current 2>/dev/null || echo "")

  if [[ "$current_branch" != "main" && -n "$current_branch" ]]; then
    active_hash=$(git diff main | calculate_sha256)
  else
    active_hash=$(git diff HEAD | calculate_sha256)
  fi

  # Check if status is approved and diff_hash matches
  local status=""
  status=$(jq -r '.status // empty' "$review_file" 2>/dev/null || echo "")
  local diff_hash=""
  diff_hash=$(jq -r '.diff_hash // empty' "$review_file" 2>/dev/null || echo "")

  if [[ "$status" != "approved" ]]; then
    echo "Error: Proactive review approval status is '$status' (not approved)." >&2
    exit 1
  fi

  if [[ "$diff_hash" != "$active_hash" ]]; then
    echo "Error: Local changes have been modified since your last proactive review!" >&2
    echo "       Approved SHA-256 hash: $diff_hash" >&2
    echo "       Current active SHA-256 hash: $active_hash" >&2
    echo "       Please run the review agent again on your latest changes." >&2
    exit 1
  fi

  echo "✅ Proactive review approval verified! (SHA-256 Hash: $active_hash)" >&2
}
