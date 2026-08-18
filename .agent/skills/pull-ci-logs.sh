#!/usr/bin/env bash
#
# Skill: pull-ci-logs.sh
# Description: Retrieve GitHub CI workflow logs and save to a temporary file, supporting listing failed runs and jobs.
#              Conforms to shell-scripts.instructions.md guidelines.
# Usage: .agent/skills/pull-ci-logs.sh [run-id] [options]

set -euo pipefail

show_help() {
  cat <<EOF
Usage: pull-ci-logs.sh [run-id] [options]

Downloads workflow logs from a GitHub Actions run or specific job to a temporary file.

Arguments:
  run-id                  Optional. The database ID of the run. If omitted,
                          the single most recent run matching the criteria is used.

Options:
  -r, --repo OWNER/REPO   The GitHub repository (default: rancher/terraform-provider-rancher2).
  -w, --workflow NAME     Filter runs by workflow name (e.g., "Release", "pull_request").
  -s, --status STATUS     Filter runs by status (e.g., "completed", "failure").
  -f, --failed-only       Only fetch logs for failed steps.
  -o, --output FILE       The file path to save logs (defaults to a securely generated temp file).
  --list-failed           List recently failed workflow runs and exit.
  --list-jobs RUN_ID      List failed jobs inside a specific workflow run and exit.
  -j, --job JOB_ID        Download/view logs for a specific job ID instead of the full run.
  -h, --help              Show this help message and exit.

Examples:
  # List recently failed workflow runs
  $ .agent/skills/pull-ci-logs.sh --list-failed

  # List failed jobs within workflow run 123456789
  $ .agent/skills/pull-ci-logs.sh --list-jobs 123456789

  # Download logs for the most recent run
  $ .agent/skills/pull-ci-logs.sh

  # Download logs from a specific failed job ID
  $ .agent/skills/pull-ci-logs.sh --job 987654321
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

list_failed_runs() {
  local repo="$1"
  echo "Fetching recently failed runs for $repo..." >&2

  local run_data
  run_data=$(run_with_retry gh run list -R "$repo" -s failure --limit 10 --json databaseId,workflowName,headBranch,createdAt,displayTitle 2>/dev/null)

  if [[ -z "$run_data" || "$run_data" == "[]" ]]; then
    echo "No recently failed workflow runs found for $repo." >&2
    return 0
  fi

  echo -e "RUN ID\t\tWORKFLOW / DETAILS"
  echo "================================================================================"
  echo "$run_data" | jq -r '.[] | "\(.databaseId)\t\(.workflowName) (\(.headBranch)) - \(.displayTitle) (\(.createdAt))"'
}

list_failed_jobs() {
  local repo="$1"
  local run_id="$2"
  echo "Fetching failed jobs for run $run_id in $repo..." >&2

  local jobs_data
  jobs_data=$(run_with_retry gh api "repos/${repo}/actions/runs/${run_id}/jobs" 2>/dev/null)

  if [[ -z "$jobs_data" ]]; then
    echo "Error: Could not retrieve jobs for run $run_id." >&2
    exit 1
  fi

  local failed_jobs
  failed_jobs=$(echo "$jobs_data" | jq -r '.jobs[] | select(.conclusion == "failure") | "\(.id)\t\(.name)"' || true)

  if [[ -z "$failed_jobs" ]]; then
    echo "No failed jobs found for run $run_id (all jobs passed or are pending)." >&2
    return 0
  fi

  echo -e "JOB ID\t\tFAILED JOB NAME"
  echo "================================================================================"
  echo "$failed_jobs"
}

get_latest_run_id() {
  local repo="$1"
  local workflow="$2"
  local status="$3"

  local msg="Fetching the latest run ID"
  if [[ -n "$workflow" ]]; then
    msg="$msg for workflow '$workflow'"
  fi
  if [[ -n "$status" ]]; then
    msg="$msg with status '$status'"
  fi
  msg="$msg from $repo..."
  echo "$msg" >&2

  local extra_args=()
  if [[ -n "$workflow" ]]; then
    extra_args+=("-w" "$workflow")
  fi
  if [[ -n "$status" ]]; then
    extra_args+=("-s" "$status")
  fi

  local run_id
  run_id=$(run_with_retry gh run list -R "$repo" --limit 1 "${extra_args[@]}" --json databaseId --jq '.[0].databaseId' 2>/dev/null || echo "")

  if [[ -z "$run_id" || "$run_id" == "null" ]]; then
    echo "Error: No recent workflow runs found matching the criteria for repository '$repo'." >&2
    exit 1
  fi

  echo "$run_id"
}

main() {
  local repo="rancher/terraform-provider-rancher2"
  local failed_only=false
  local output_file=""
  local run_id=""
  local workflow=""
  local status=""
  local job_id=""
  
  local action="download"
  local target_run_id=""

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      -r|--repo)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --repo requires an argument." >&2
          exit 1
        fi
        repo="$2"
        shift 2
        ;;
      -w|--workflow)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --workflow requires an argument." >&2
          exit 1
        fi
        workflow="$2"
        shift 2
        ;;
      -s|--status)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --status requires an argument." >&2
          exit 1
        fi
        status="$2"
        shift 2
        ;;
      -f|--failed-only)
        failed_only=true
        shift
        ;;
      -o|--output)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --output requires an argument." >&2
          exit 1
        fi
        output_file="$2"
        shift 2
        ;;
      -j|--job)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --job requires a job ID." >&2
          exit 1
        fi
        job_id="$2"
        shift 2
        ;;
      --list-failed)
        action="list-failed"
        shift
        ;;
      --list-jobs)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --list-jobs requires a run ID argument." >&2
          exit 1
        fi
        action="list-jobs"
        target_run_id="$2"
        shift 2
        ;;
      -*)
        echo "Error: Unknown option: $1" >&2
        show_help
        exit 1
        ;;
      *)
        if [[ -n "$run_id" ]]; then
          echo "Error: Only one run-id can be specified." >&2
          exit 1
        fi
        if ! [[ "$1" =~ ^[0-9]+$ ]]; then
          echo "Error: Invalid run-id '$1'. Run-id must be a number." >&2
          exit 1
        fi
        run_id="$1"
        shift
        ;;
    esac
  done

  # Ensure gh CLI is installed
  if ! command -v gh &>/dev/null; then
    echo "Error: The GitHub CLI (gh) is not installed or not in PATH." >&2
    exit 1
  fi

  # Ensure jq is installed
  if ! command -v jq &>/dev/null; then
    echo "Error: jq is required but not installed or not in PATH." >&2
    exit 1
  fi

  # Execute requested action
  if [[ "$action" == "list-failed" ]]; then
    list_failed_runs "$repo"
    exit 0
  fi

  if [[ "$action" == "list-jobs" ]]; then
    list_failed_jobs "$repo" "$target_run_id"
    exit 0
  fi

  # Retrieve logs path (for downloading logs)
  local view_flags=()
  if [[ "$failed_only" == "true" ]]; then
    view_flags+=("--log-failed")
  else
    view_flags+=("--log")
  fi

  if [[ -n "$job_id" ]]; then
    if [[ -z "$output_file" ]]; then
      output_file="$(mktemp "/tmp/gh-job-${job_id}.XXXXXX.log")"
    fi
    mkdir -p "$(dirname "$output_file")"
    echo "Downloading logs for job $job_id from $repo..."
    if ! run_with_retry gh run view --job "$job_id" -R "$repo" "${view_flags[@]}" > "$output_file"; then
      echo "Error: Failed to fetch logs for job $job_id." >&2
      exit 1
    fi
  else
    if [[ -z "$run_id" ]]; then
      run_id=$(get_latest_run_id "$repo" "$workflow" "$status")
    fi
    if [[ -z "$output_file" ]]; then
      output_file="$(mktemp "/tmp/gh-run-${run_id}.XXXXXX.log")"
    fi
    mkdir -p "$(dirname "$output_file")"
    echo "Downloading logs for run $run_id from $repo..."
    if ! run_with_retry gh run view "$run_id" -R "$repo" "${view_flags[@]}" > "$output_file"; then
      echo "Error: Failed to fetch logs for run $run_id." >&2
      exit 1
    fi
  fi

  echo "Logs successfully written to: $output_file"
  echo "You can view them using: less -R \"$output_file\" or code \"$output_file\""
}

main "$@"
