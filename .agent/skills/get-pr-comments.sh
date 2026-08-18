#!/usr/bin/env bash
#
# Skill: get-pr-comments.sh
# Description: Retrieves and sorts general and inline PR review comments chronologically for an open PR.
# Usage: .agent/skills/get-pr-comments.sh [PR_ID] [FORMAT]

# Fail fast on any errors, unbound variables, or pipe failures
set -euo pipefail

# Display usage instructions for the script
show_usage() {
    cat <<EOF
Usage: $0 [PR_ID] [FORMAT]

Retrieves both top-level (general conversation) and inline (diff review) comments 
for a given GitHub Pull Request ID, merging and sorting them chronologically.

If no PR_ID is provided, the script will automatically detect the open pull request
for the current Git branch.

Arguments:
  PR_ID       The numeric ID of the Pull Request (optional if on a branch with an open PR)
  FORMAT      Output format: "markdown" (default) or "json" (optional)

Dependencies:
  - gh (GitHub CLI), configured and authenticated
  - jq (JSON processor)
  - git

Options:
  -h, --help  Show this help message and exit
EOF
}

# Verify that all necessary CLI dependencies are installed and authenticated
check_dependencies() {
    if ! command -v gh &> /dev/null; then
        echo "Error: 'gh' (GitHub CLI) is not installed." >&2
        echo "Please install it from https://cli.github.com/ and authenticate with 'gh auth login'." >&2
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        echo "Error: 'jq' is not installed." >&2
        echo "Please install 'jq' to run this script." >&2
        exit 1
    fi

    if ! command -v git &> /dev/null; then
        echo "Error: 'git' is not installed." >&2
        exit 1
    fi

    # Check if gh CLI is authenticated (either via active login or GITHUB_TOKEN env var)
    if [[ -z "${GITHUB_TOKEN:-}" ]] && ! gh auth status &> /dev/null; then
        echo "Error: 'gh' CLI is not authenticated." >&2
        echo "Please run 'gh auth login' or set the GITHUB_TOKEN environment variable." >&2
        exit 1
    fi
}

# Safely executes a command with retry and exponential backoff
# Usage: run_with_retry cmd args...
run_with_retry() {
    local max_attempts=5
    local base_delay=2
    local attempt=1
    local exit_code=0

    while true; do
        if "$@"; then
            return 0
        else
            exit_code=$?
        fi

        if [[ ${attempt} -ge ${max_attempts} ]]; then
            echo "Error: Command '$*' failed after ${max_attempts} attempts." >&2
            return ${exit_code}
        fi

        local delay
        delay=$(( base_delay * (2 ** (attempt - 1)) ))
        echo "Warning: Command failed (exit code ${exit_code}). Retrying in ${delay} seconds (attempt ${attempt}/${max_attempts})..." >&2
        sleep "${delay}"
        attempt=$((attempt + 1))
    done
}

# Parse the owner/repo from the git remote url, defaulting to canonical upstream
get_repo_context() {
    # Allow overriding via environment variable
    if [[ -n "${REPO_CONTEXT:-}" ]]; then
        echo "${REPO_CONTEXT}"
        return
    fi

    # Try upstream remote first
    local url
    url=$(git remote get-url upstream 2>/dev/null || git config --get remote.upstream.url 2>/dev/null || echo "")

    # If no upstream, default to canonical rancher repository
    if [[ -z "${url}" ]]; then
      echo "rancher/terraform-provider-rancher2"
      return
    fi

    local clean_url="${url}"
    clean_url="${clean_url#*github.com[:/]}"
    clean_url="${clean_url%.git}"
    echo "${clean_url}"
}

# Determine the fork owner from the git remote origin url (to support cross-repository PRs)
get_fork_owner() {
    local url
    url=$(git remote get-url origin 2>/dev/null || git config --get remote.origin.url 2>/dev/null || echo "")
    if [[ -z "${url}" ]]; then
        echo ""
        return
    fi
    local clean_url="${url}"
    clean_url="${clean_url#*github.com[:/]}"
    clean_url="${clean_url%.git}"
    echo "${clean_url%%/*}"
}

# Find the PR ID associated with the current branch
get_current_branch_pr_id() {
    local target_repo="$1"
    local branch
    branch=$(git branch --show-current 2>/dev/null || git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    
    if [[ -z "${branch}" ]]; then
        echo "Error: Could not determine current git branch." >&2
        exit 1
    fi

    local pr_id=""

    # 1. Try querying as a cross-repository PR (fork-owner:branch)
    local fork_owner
    fork_owner=$(get_fork_owner)
    if [[ -n "${fork_owner}" && "${fork_owner}" != "rancher" ]]; then
        pr_id=$(run_with_retry gh pr list --repo "${target_repo}" --state open --head "${fork_owner}:${branch}" --json number --jq '.[0].number' 2>/dev/null || echo "")
    fi

    # 2. Try querying as a same-repository PR (branch) if cross-repo query yielded nothing
    if [[ -z "${pr_id}" ]]; then
        pr_id=$(run_with_retry gh pr list --repo "${target_repo}" --state open --head "${branch}" --json number --jq '.[0].number' 2>/dev/null || echo "")
    fi
    
    if [[ -z "${pr_id}" ]]; then
        echo "Error: No open pull request found in '${target_repo}' for current branch '${branch}'." >&2
        exit 1
    fi

    echo "${pr_id}"
}

# Fetch comments from GitHub API and format the output
fetch_comments() {
    local repo_context="$1"
    local pr_id="$2"
    local format="${3:-markdown}"

    # Verify that the PR exists in the repository
    if ! run_with_retry gh pr view "${pr_id}" --repo "${repo_context}" &> /dev/null; then
        echo "Error: Pull Request #${pr_id} was not found in repo '${repo_context}'." >&2
        exit 1
    fi

    # Fetch top-level and inline/review comments using parsed repository
    local general_comments
    local review_comments

    if ! general_comments=$(run_with_retry gh api "repos/${repo_context}/issues/${pr_id}/comments" --paginate 2>/dev/null); then
        echo "Error: Failed to fetch general comments from the GitHub API." >&2
        exit 1
    fi

    if ! review_comments=$(run_with_retry gh api "repos/${repo_context}/pulls/${pr_id}/comments" --paginate 2>/dev/null); then
        echo "Error: Failed to fetch review comments from the GitHub API." >&2
        exit 1
    fi

    # Ensure we got valid JSON arrays
    if ! echo "${general_comments}" | jq -e 'type == "array"' &>/dev/null || \
       ! echo "${review_comments}" | jq -e 'type == "array"' &>/dev/null; then
        echo "Error: Received invalid JSON response from GitHub API." >&2
        exit 1
    fi

    # Format the combined comments based on user preference
    if [[ "${format}" == "json" ]]; then
        # Combine general and review comments, tag each with its type, and sort them chronologically
        jq -n \
            --argjson gen "${general_comments}" \
            --argjson rev "${review_comments}" \
            '
                # 1. Map general comments and append type: "general"
                ($gen | map(. + {type: "general"})) + 
                # 2. Map review comments and append type: "review"
                ($rev | map(. + {type: "review"})) | 
                # 3. Combine both lists and sort by creation timestamp
                sort_by(.created_at)
            '
    else
        local markdown_output
        # Compile and format comments chronologically as a beautiful, aligned Markdown document
        markdown_output=$(jq -n -r \
            --argjson gen "${general_comments}" \
            --argjson rev "${review_comments}" \
            '
                # 1. Combine and tag both comment streams
                (($gen | map(. + {type: "general"})) + 
                 ($rev | map(. + {type: "review"}))) | 
                # 2. Sort the combined list chronologically
                sort_by(.created_at) |
                # 3. Check if the list is empty
                if length == 0 then
                    "No comments found."
                else
                    # 4. Iterate over each comment and generate its Markdown block
                    .[] | (
                        # A. Header line with type icon (💬 for general, 📝 for review) and username
                        "### " + (if .type == "general" then "💬" else "📝" end) + " @" + .user.login + 
                        
                        # B. Context details (such as file path and line number for inline reviews)
                        " (" + (if .type == "general" then "General Comment" else "Inline Review on `" + .path + ":" + (.line // .original_line // "unknown" | tostring) + "`" end) + ") - " + 
                        
                        # C. Timestamp parsed to readable UTC format
                        (.created_at | sub("T"; " ") | sub("Z"; " UTC")) + "\n\n" + 
                        
                        # D. The markdown body of the comment itself
                        .body + "\n\n" + 
                        
                        # E. Standard separator line
                        "---"
                    )
                end
            ')

        if [[ "${markdown_output}" == "No comments found." ]]; then
            printf "No comments found for PR #%s.\n" "${pr_id}"
        else
            printf "# Comments for PR #%s\n\n%s\n" "${pr_id}" "${markdown_output}"
        fi
    fi
}

main() {
    # Check for help options
    if [[ "${1:-}" == "-h" ]] || [[ "${1:-}" == "--help" ]]; then
        show_usage
        exit 0
    fi

    check_dependencies

    # Extract owner/repo context (defaults to upstream canonical repo)
    local repo_context
    repo_context=$(get_repo_context)

    local pr_id=""
    local format="markdown"

    # Argument parsing
    if [[ $# -eq 0 ]]; then
        # No arguments given -> find PR for the current branch
        pr_id=$(get_current_branch_pr_id "${repo_context}")
    elif [[ $# -eq 1 ]]; then
        # One argument given: could be PR ID or format
        local arg="$1"
        if [[ "${arg}" =~ ^[0-9]+$ ]]; then
            pr_id="${arg}"
        else
            # Try to interpret it as format, and detect PR ID
            format="${arg}"
            pr_id=$(get_current_branch_pr_id "${repo_context}")
        fi
    else
        # Two or more arguments given
        pr_id="$1"
        format="$2"
    fi

    # Normalize format casing
    format=$(echo "${format}" | tr '[:upper:]' '[:lower:]')

    if [[ "${format}" != "markdown" ]] && [[ "${format}" != "json" ]]; then
        echo "Error: Invalid output format '${format}'. Supported formats: markdown, json" >&2
        exit 1
    fi

    fetch_comments "${repo_context}" "${pr_id}" "${format}"
}

# Execute main function with all arguments
main "$@"
