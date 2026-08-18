#!/usr/bin/env bash
#
# Skill: update-action-versions.sh
# Description: Automatically audits and updates the commit SHAs and version tags of GitHub Actions workflows to the latest available releases.
# Usage: .agent/skills/update-action-versions.sh

set -euo pipefail

show_help() {
  cat <<EOF
Usage: update-action-versions.sh [options]

Automatically audits and updates the commit SHAs and version tags of GitHub Actions workflows
to their latest available releases.

Options:
  --list-actions       List all actions used in workflow files and their current versions.
  -h, --help           Show this help message and exit.

Examples:
  .agent/skills/update-action-versions.sh --list-actions
  .agent/skills/update-action-versions.sh
EOF
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

list_current_actions() {
  local search_dir=".github/workflows"
  echo "Scanning GitHub Actions workflow files to list action usages and versions..." >&2

  local actions_data
  # shellcheck disable=SC2016
  actions_data=$(find "${search_dir}" -type f \( -name "*.yml" -o -name "*.yaml" \) -print0 2>/dev/null \
    | xargs -0 awk '
      # Match lines declaring action usages, e.g. "uses: actions/checkout@v4"
      /[[:space:]]+uses:[[:space:]]+/ {
        str = $0
        sub(/^[[:space:]]*uses:[[:space:]]*/, "", str)
        comment = ""
        
        # Check if there is a trailing version comment (e.g. # v4.2.2)
        if (str ~ /#/) {
          idx = index(str, "#")
          comment = substr(str, idx + 1)
          gsub(/^[[:space:]]+|[[:space:]]+$/, "", comment)
          str = substr(str, 1, idx - 1)
          gsub(/[[:space:]]+$/, "", str)
        }
        
        # Split action name and reference SHA/tag on "@"
        split(str, parts, "@")
        action = parts[1]
        ver = parts[2]
        
        # Only print valid action declarations (must contain slash and have version)
        if (action ~ /\// && ver != "") {
          # Strip path prefix to make filename output clean
          filename_clean = FILENAME
          sub(/^.*\.github\/workflows\//, "", filename_clean)
          
          # Output tab-separated records: action \t ver \t comment \t filename
          print action "\t" ver "\t" (comment != "" ? comment : "N/A") "\t" filename_clean
        }
      }
    ' 2>/dev/null || echo "")

  if [[ -z "${actions_data}" ]]; then
    echo "No GitHub Actions references found in workflow files." >&2
    return 0
  fi

  # Format and print the tab-separated results as an aligned table
  echo "${actions_data}" | sort | awk -F'\t' '
    # Print the table headers and divider line
    BEGIN {
      printf "%-38s %-43s %-15s %s\n", "ACTION", "COMMIT SHA / REF", "RELEASE TAG", "WORKFLOW FILE"
      printf "========================================================================================================================\n"
    }
    # Print each aligned row
    {
      printf "%-38s %-43s %-15s %s\n", $1, $2, $3, $4
    }
  '
}

update_workflow_releases() {
  local workflows
  workflows=$(find .github/workflows -type f \( -name '*.yml' -o -name '*.yaml' \) 2>/dev/null || echo "")

  if [[ -z "${workflows}" ]]; then
    echo "No workflow files found." >&2
    return 0
  fi

  local workflow
  for workflow in ${workflows}; do
    echo "Processing ${workflow}..." >&2

    # Extract unique repos from the workflow file
    # Format expected: # https://github.com/owner/repo/releases
    local repos
    repos=$(awk '/^[[:space:]]*#[[:space:]]*https:\/\/github\.com\/[^\/]+\/[^\/]+\/releases/ {
      str = $0
      sub(/^[[:space:]]*#[[:space:]]*https:\/\/github\.com\//, "", str)
      sub(/\/releases.*/, "", str)
      print str
    }' "${workflow}" | sort -u || true)

    local repo
    for repo in ${repos}; do
      echo "Found release link for ${repo}" >&2

      # Get latest release tag
      local tag
      tag=$(run_with_retry gh api "/repos/${repo}/releases/latest" --jq '.tag_name' 2>/dev/null || echo "")

      if [[ -z "${tag}" || "${tag}" == "null" ]]; then
        echo "Failed to get latest tag for ${repo}" >&2
        continue
      fi

      # Get commit sha for tag
      local sha
      sha=$(run_with_retry gh api "/repos/${repo}/commits/${tag}" --jq '.sha' 2>/dev/null || echo "")

      if [[ -z "${sha}" ]]; then
        echo "Failed to get commit SHA for ${repo} tag ${tag}" >&2
        continue
      fi

      echo "Latest version for ${repo} is ${tag} (${sha})" >&2

      # Use a temporary file for replacement
      local tmp_file
      tmp_file=$(mktemp)

      awk -v repo="${repo}" -v new_sha="${sha}" -v new_tag="${tag}" '
      {
        # Check if the current line is a release link comment acting as an anchor
        # e.g., "# https://github.com/actions/checkout/releases"
        if ($0 ~ "^[[:space:]]*#[[:space:]]*https://github\\.com/" repo "/releases") {
          print $0 # Print the anchor comment line unchanged
          
          # Read the immediately following line (the action declaration)
          getline
          
          # Check if the next line is the expected "- uses: owner/repo@<sha>" format
          if ($0 ~ "^[[:space:]]*- uses: " repo "@") {
            # Find where "- uses: " starts to preserve the original indentation
            idx = index($0, "- uses: ")
            indent = substr($0, 1, idx - 1)
            
            # Print the line updated with the new commit SHA and trailing tag comment
            # e.g., "      - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2"
            print indent "- uses: " repo "@" new_sha " # " new_tag
          } else {
            # If the next line is not a "- uses:" line, print it unchanged
            print $0
          }
          next
        }
        # Print all other non-anchor lines unchanged
        print $0
      }
      ' "${workflow}" > "${tmp_file}"

      mv "${tmp_file}" "${workflow}"
    done
  done
}

main() {
  local action="update"

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --list-actions)
        action="list"
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

  if [[ "$action" == "list" ]]; then
    list_current_actions
    exit 0
  fi

  update_workflow_releases
  echo "Done updating action versions."
}

main "$@"
