#!/usr/bin/env bash
#
# Description: Modular Git push with lease helper.
#

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

execute_push() {
  local remote_name="$1"
  local branch="$2"
  local force_push="$3"

  verify_push_safety "$remote_name"

  if [[ "$force_push" == "true" ]]; then
    echo "Safely force-pushing branch '$branch' to '$remote_name' with lease..." >&2
    if ! git push -u "$remote_name" "$branch" --force-with-lease; then
      echo "Error: Remote force-push with lease failed." >&2
      exit 1
    fi
  else
    echo "Pushing branch '$branch' to '$remote_name'..." >&2
    if ! git push -u "$remote_name" "$branch"; then
      echo "Error: Remote push failed." >&2
      exit 1
    fi
  fi

  echo "✅ Changes successfully pushed to remote '$remote_name/$branch'!" >&2
}
