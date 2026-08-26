#!/usr/bin/env bash
# shellcheck disable=SC2016
#
# Skill: resolve-pr-reviews.sh
# Description: Programmatically list and resolve review comment threads on a GitHub Pull Request using GraphQL and the GitHub CLI.
# Conforms to shell-scripts.instructions.md guidelines.

set -euo pipefail

# Helper to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Display script help usage instructions
show_help() {
  cat <<EOF
Usage: resolve-pr-reviews.sh [PR_ID] [options/file_pattern]

Programmatically list and resolve review comment threads on a GitHub Pull Request.

Arguments:
  PR_ID                 The numeric ID of the Pull Request (optional if on a branch with an open PR).
  OPTIONS/PATTERN       Filter or resolution command options.

Options:
  -h, --help            Show this message and exit.
  --all                 Resolve ALL unresolved comment threads.
  --bypass-token        Bypass ambient GITHUB_TOKEN environment variables and force local keyring auth.
  <pattern>             Resolve threads where the file path contains the given literal pattern.

Examples:
  .gemini/skills/resolve-pr-reviews.sh 390
  .gemini/skills/resolve-pr-reviews.sh 390 --all
  .gemini/skills/resolve-pr-reviews.sh 390 --bypass-token --all
  .gemini/skills/resolve-pr-reviews.sh 390 publish-release.test.js
EOF
}

verify_environment() {
  if ! command_exists gh; then
    echo "Error: GitHub CLI (gh) is required but not installed." >&2
    exit 1
  fi

  if ! command_exists jq; then
    echo "Error: jq (JSON processor) is required but not installed." >&2
    exit 1
  fi

  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: Must be run inside a Git repository." >&2
    exit 1
  fi
}

resolve_thread() {
  local thread_id="$1"
  local path="$2"
  local author="$3"
  local owner="$4"
  local repo="$5"
  
  echo "Resolving thread $thread_id on '$path' by @$author..."
  local mutation='
  mutation($threadId: ID!) {
    resolveReviewThread(input: {threadId: $threadId}) {
      thread {
        id
        isResolved
      }
    }
  }
  '
  gh api graphql -F threadId="$thread_id" -f query="$mutation" >/dev/null
  echo "✅ Thread successfully resolved!"
}

main() {
  local pr_id=""
  local mode="list"
  local filter=""
  local bypass_token=false

  # Parse Help / Arguments
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --bypass-token)
        bypass_token=true
        shift
        ;;
      --all)
        mode="all"
        shift
        ;;
      [0-9]*)
        pr_id="$1"
        shift
        ;;
      *)
        # If it is not a recognized flag, treat it as the file path filter
        mode="filter"
        filter="$1"
        shift
        ;;
    esac
  done

  # Perform conditional token bypass if explicitly requested
  if [[ "$bypass_token" == "true" ]]; then
    echo "Bypassing GITHUB_TOKEN environment variable to use native keychain credentials..."
    unset GITHUB_TOKEN
  fi

  verify_environment

  # Autodetect PR ID if missing
  if [[ -z "$pr_id" ]]; then
    local current_branch
    current_branch=$(git branch --show-current)
    if [[ -n "$current_branch" ]]; then
      echo "Autodetecting open PR for branch '$current_branch' on origin/upstream..."
      pr_id=$(gh pr list --head "$current_branch" --json number --jq '.[0].number' 2>/dev/null || true)
    fi
  fi

  if [[ -z "$pr_id" ]]; then
    echo "Error: No Pull Request number provided and could not autodetect an open PR for the current branch." >&2
    show_help >&2
    exit 1
  fi

  # Report Mode
  if [[ "$mode" == "filter" ]]; then
    echo "Filtering threads containing file path pattern: '$filter'"
  fi

  # Query the PR URL to parse the exact target upstream Owner and Repo
  local pr_url
  pr_url=$(gh pr view "$pr_id" --json url --jq '.url' 2>/dev/null || true)
  
  local owner=""
  local repo=""
  if [[ -n "$pr_url" && "$pr_url" =~ github\.com/([^/]+)/([^/]+)/pull/([0-9]+) ]]; then
    owner="${BASH_REMATCH[1]}"
    repo="${BASH_REMATCH[2]}"
    pr_id="${BASH_REMATCH[3]}"
  else
    # Local fallback if gh pr view fails
    local origin_url
    origin_url=$(git remote get-url origin 2>/dev/null || true)
    if [[ -n "$origin_url" && "$origin_url" =~ github\.com[:/]([^/]+)/([^/.]+)(\.git)?$ ]]; then
      owner="${BASH_REMATCH[1]}"
      repo="${BASH_REMATCH[2]}"
    else
      echo "Error: Could not determine target Owner and Repo." >&2
      exit 1
    fi
  fi

  echo "Connected to PR #$pr_id on $owner/$repo"

  # Fetch Unresolved Threads via GraphQL
  echo "Fetching unresolved review comment threads..."
  local query='
  query($owner: String!, $repo: String!, $pullNumber: Int!) {
    repository(owner: $owner, name: $repo) {
      pullRequest(number: $pullNumber) {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 1) {
              nodes {
                body
                path
                author {
                  login
                }
              }
            }
          }
        }
      }
    }
  }
  '

  local threads_json
  threads_json=$(gh api graphql -F owner="$owner" -F repo="$repo" -F pullNumber="$pr_id" -f query="$query")

  # Detect and report GraphQL errors
  if echo "$threads_json" | jq -e '.errors' >/dev/null 2>&1; then
    local err_msg
    err_msg=$(echo "$threads_json" | jq -r '.errors[0].message')
    echo "Error: GitHub GraphQL API returned errors: $err_msg" >&2
    exit 1
  fi

  local unresolved_threads
  unresolved_threads=$(echo "$threads_json" | jq '.data.repository.pullRequest.reviewThreads.nodes | map(select(.isResolved == false))')
  
  local thread_count
  thread_count=$(echo "$unresolved_threads" | jq '. | length')

  if [[ "$thread_count" -eq 0 ]]; then
    echo "🎉 No unresolved comment threads found on PR #$pr_id!"
    exit 0
  fi

  echo "Found $thread_count unresolved comment thread(s) on PR #$pr_id:"

  # Loop and display/resolve
  local thread_id=""
  local file_path=""
  local author=""
  local body=""
  local thread=""

  for i in $(seq 0 $((thread_count - 1))); do
    thread=$(echo "$unresolved_threads" | jq ".[$i]")
    thread_id=$(echo "$thread" | jq -r '.id')
    file_path=$(echo "$thread" | jq -r '.comments.nodes[0].path')
    author=$(echo "$thread" | jq -r '.comments.nodes[0].author.login')
    body=$(echo "$thread" | jq -r '.comments.nodes[0].body' | tr '\r\n' ' ' | cut -c1-80)

    echo "------------------------------------------------------------"
    echo "Thread ID : $thread_id"
    echo "File Path : $file_path"
    echo "Author    : @$author"
    echo "Comment   : $body..."

    if [[ "$mode" == "all" ]]; then
      resolve_thread "$thread_id" "$file_path" "$author" "$owner" "$repo"
    elif [[ "$mode" == "filter" ]]; then
      # Safe literal substring matching to avoid regex crash and metacharacter matches
      if [[ "$file_path" == *"$filter"* ]]; then
        resolve_thread "$thread_id" "$file_path" "$author" "$owner" "$repo"
      else
        echo "  -> Skipping (does not match filter)"
      fi
    else
      echo "  -> Running in list mode. Run with '$0 $pr_id --all' or '$0 $pr_id $file_path' to resolve."
    fi
  done
  echo "------------------------------------------------------------"
}

main "$@"
