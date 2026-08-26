#!/usr/bin/env bash
#
# Description: Modular helper functions for GitHub PR operations.
#

get_fork_owner() {
  local origin_url
  origin_url=$(git remote get-url origin 2>/dev/null || true)
  if [[ -z "$origin_url" ]]; then
    echo "Error: Could not retrieve configured origin remote URL." >&2
    exit 1
  fi

  if [[ "$origin_url" =~ github\.com[:/]([^/]+)/ ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    echo "Error: Could not parse fork owner from origin URL $origin_url" >&2
    exit 1
  fi
}

create_pull_request() {
  local title="$1"
  local body="$2"
  local base="$3"
  local draft_flag="$4"
  local branch
  local fork_owner

  branch=$(git branch --show-current)
  fork_owner=$(get_fork_owner)

  # Check if a pull request already exists for this branch
  local existing_pr
  existing_pr=$(GITHUB_TOKEN="" gh pr list --head "${branch}" --json url --jq '.[0].url' 2>/dev/null || true)
  if [[ -n "${existing_pr}" ]]; then
    echo "✅ Pull Request already exists for branch '${branch}': ${existing_pr}"
    return 0
  fi

  # Extract upstream repo name dynamically (defaulting to the parsed repo name)
  local upstream_repo=""
  local origin_url
  origin_url=$(git remote get-url origin 2>/dev/null || true)
  if [[ "$origin_url" =~ github\.com[:/]([^/]+)/([^/]+)\.git ]]; then
    upstream_repo="rancher/${BASH_REMATCH[2]}"
  elif [[ "$origin_url" =~ github\.com[:/]([^/]+)/([^/]+) ]]; then
    upstream_repo="rancher/${BASH_REMATCH[2]}"
  else
    local repo_name
    repo_name=$(basename "$(git rev-parse --show-toplevel 2>/dev/null || pwd)")
    upstream_repo="rancher/${repo_name}"
  fi

  echo "Creating Pull Request for branch '${branch}' on fork '${fork_owner}' to upstream '${upstream_repo}'..."

  # Set the default repository to ensure gh targets upstream
  GITHUB_TOKEN="" gh repo set-default "$upstream_repo"

  # Create the PR non-interactively, bypassing GITHUB_TOKEN override
  # shellcheck disable=SC2086
  GITHUB_TOKEN="" gh pr create \
    --repo "$upstream_repo" \
    --base "$base" \
    --head "${fork_owner}:${branch}" \
    --title "$title" \
    --body "$body" \
    $draft_flag
}

graduate_pull_request() {
  local target="$1"
  local branch
  branch=$(git branch --show-current)

  if [[ -z "$target" ]]; then
    echo "Graduating draft pull request for the current branch '${branch}'..."
    GITHUB_TOKEN="" gh pr ready
  else
    echo "Graduating draft pull request for target '${target}'..."
    GITHUB_TOKEN="" gh pr ready "$target"
  fi
  echo "✅ Pull Request successfully graduated to ready for review!"
}
