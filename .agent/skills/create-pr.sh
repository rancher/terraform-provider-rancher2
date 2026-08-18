#!/usr/bin/env bash
#
# Skill: create-pr.sh
# Description: Safely creates a pull request from your local fork branch to the upstream main branch,
#              bypassing ambient GITHUB_TOKEN environment overrides to use the gh CLI's authenticated credential.
# Usage: .agent/skills/create-pr.sh --title "<Title>" --body "<Body>" [--base "<Base>"] [--draft]

set -euo pipefail

show_help() {
  cat <<EOF
Usage: create-pr.sh [options]

Safely creates or graduates a pull request from your local fork branch to the upstream repository.

Options:
  --title TITLE        The title of the pull request (Required for creation).
  --body BODY          The markdown description body of the pull request (Required for creation).
  --base BASE          The target upstream branch (default: main).
  --draft              Create the pull request as a draft.
  --ready [TARGET]     Graduate a draft pull request to ready-for-review (accepts PR number, branch, or URL; defaults to current branch).
  -h, --help           Show this help message and exit.

Examples:
  .agent/skills/create-pr.sh --title "fix: logic error" --body "Fixes the loop bounds"
  .agent/skills/create-pr.sh --title "feat: new helper" --body "Adds skill" --draft
  .agent/skills/create-pr.sh --ready
  .agent/skills/create-pr.sh --ready 280
EOF
}

verify_git_env() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: This command must be run inside a Git repository." >&2
    exit 1
  fi
}

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

get_upstream_repo() {
  local upstream_repo="rancher/terraform-provider-rancher2"
  local origin_url
  origin_url=$(git remote get-url origin 2>/dev/null || true)
  if [[ "$origin_url" =~ github\.com[:/][^/]+/([^/]+)\.git ]]; then
    upstream_repo="rancher/${BASH_REMATCH[1]}"
  elif [[ "$origin_url" =~ github\.com[:/][^/]+/([^/]+) ]]; then
    upstream_repo="rancher/${BASH_REMATCH[1]}"
  fi
  echo "$upstream_repo"
}

create_pull_request() {
  local title="$1"
  local body="$2"
  local base="$3"
  local draft_flag="$4"
  local branch
  local fork_owner
  local upstream_repo

  branch=$(git branch --show-current)
  fork_owner=$(get_fork_owner)
  upstream_repo=$(get_upstream_repo)

  echo "Creating Pull Request for branch '${branch}' on fork '${fork_owner}' to upstream '${upstream_repo}'..."

  local extra_flags=()
  if [[ -n "$draft_flag" ]]; then
    extra_flags+=("$draft_flag")
  fi

  # Create the PR non-interactively, bypassing token overrides
  GITHUB_TOKEN="" GH_TOKEN="" gh pr create \
    --repo "$upstream_repo" \
    --base "$base" \
    --head "${fork_owner}:${branch}" \
    --title "$title" \
    --body "$body" \
    "${extra_flags[@]}"
}

graduate_pull_request() {
  local target="$1"
  local branch
  branch=$(git branch --show-current)
  local upstream_repo
  upstream_repo=$(get_upstream_repo)

  if [[ -z "$target" ]]; then
    echo "Graduating draft pull request for the current branch '${branch}'..."
    GITHUB_TOKEN="" GH_TOKEN="" gh pr ready --repo "$upstream_repo"
  else
    echo "Graduating draft pull request for target '${target}'..."
    GITHUB_TOKEN="" GH_TOKEN="" gh pr ready "$target" --repo "$upstream_repo"
  fi
  echo "✅ Pull Request successfully graduated to ready for review!"
}

main() {
  local title=""
  local body=""
  local base="main"
  local draft_flag=""
  local action="create"
  local ready_target=""

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --title)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --title requires an argument." >&2
          exit 1
        fi
        title="$2"
        shift 2
        ;;
      --body)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --body requires an argument." >&2
          exit 1
        fi
        body="$2"
        shift 2
        ;;
      --base)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --base requires an argument." >&2
          exit 1
        fi
        base="$2"
        shift 2
        ;;
      --draft)
        draft_flag="--draft"
        shift
        ;;
      --ready)
        action="ready"
        if [[ $# -gt 1 && ! "$2" =~ ^- ]]; then
          ready_target="$2"
          shift 2
        else
          ready_target=""
          shift
        fi
        ;;
      *)
        echo "Unknown parameter: $1" >&2
        show_help
        exit 1
        ;;
    esac
  done

  if [[ "$action" == "ready" ]]; then
    verify_git_env
    graduate_pull_request "$ready_target"
    exit 0
  fi

  if [[ -z "$title" || -z "$body" ]]; then
    echo "Error: Both --title and --body parameters are required." >&2
    show_help
    exit 1
  fi

  verify_git_env
  create_pull_request "$title" "$body" "$base" "$draft_flag"
}

main "$@"
