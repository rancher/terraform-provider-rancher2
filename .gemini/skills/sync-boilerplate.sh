#!/usr/bin/env bash
#
# Skill: sync-boilerplate.sh
# Description: Lightweight, manifest-driven utility to compare and synchronize configuration files
#              and boilerplate files with a centralized master template Git repository.
# Conforms to shell-scripts.instructions.md guidelines.

set -euo pipefail

# Global workspace directories
MANIFEST_FILE=".boilerplate-sync.json"
TMP_WORKSPACE=""
MODE="help"

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
Usage: sync-boilerplate.sh [options]

Lightweight utility to compare and synchronize repository configuration and boilerplate files.

Options:
  -h, --help            Show this help message and exit.
  -d, --diff            Compare local files to remote templates and print visual differences.
  -p, --pull            Pull remote template files to overwrite/update local configurations.
  -u, --push            Push local file changes back to the centralized template repository.
  -s, --status          Summarize the synchronization status of all manifest-tracked files.

Examples:
  .gemini/skills/sync-boilerplate.sh --diff
  .gemini/skills/sync-boilerplate.sh --pull
  .gemini/skills/sync-boilerplate.sh --push
  .gemini/skills/sync-boilerplate.sh --status
EOF
}

# Cleanup trap mathematically guaranteeing zero temporary file residue on exit
cleanup() {
  if [[ -n "${TMP_WORKSPACE:-}" && -d "${TMP_WORKSPACE}" ]]; then
    echo "Cleaning up temporary clone workspace..." >&2
    rm -rf "${TMP_WORKSPACE}"
  fi
}

# Validate that the system has required binaries and the local manifest exists
validate_environment() {
  if [[ ! -f "${MANIFEST_FILE}" ]]; then
    echo "Error: Boilerplate sync manifest file '${MANIFEST_FILE}' not found in root directory." >&2
    echo "       Please create '.boilerplate-sync.json' defining 'template_repo' and 'files' mapping array." >&2
    exit 1
  fi

  if ! command_exists jq; then
    echo "Error: 'jq' JSON parser utility is required but not found in current PATH." >&2
    exit 1
  fi

  if ! command_exists git; then
    echo "Error: 'git' is required but not found in current PATH." >&2
    exit 1
  fi

  # Basic JSON validation
  if ! jq empty "${MANIFEST_FILE}" 2>/dev/null; then
    echo "Error: Manifest file '${MANIFEST_FILE}' is not valid JSON." >&2
    exit 1
  fi

  # Validate template repository source (unless overridden by CENTRAL_FILE_REPO)
  if [[ -z "${CENTRAL_FILE_REPO:-}" ]]; then
    local repo
    repo=$(jq -r '.template_repo' "${MANIFEST_FILE}" 2>/dev/null || true)
    if [[ -z "$repo" || "$repo" == "null" ]]; then
      echo "Error: Manifest must define a non-empty '.template_repo' string (or CENTRAL_FILE_REPO environment variable)." >&2
      exit 1
    fi
  fi

  local files_type
  files_type=$(jq -r '.files | type' "${MANIFEST_FILE}" 2>/dev/null || true)
  if [[ "$files_type" != "array" ]]; then
    echo "Error: Manifest must define a '.files' array of mappings." >&2
    exit 1
  fi
}

# Create sandbox workspace and cleanly clone template repo
clone_template_repo() {
  local template_repo
  if [[ -n "${CENTRAL_FILE_REPO:-}" ]]; then
    template_repo="${CENTRAL_FILE_REPO}"
    echo "--> [OVERRIDE] Using central template repository from CENTRAL_FILE_REPO: ${template_repo}" >&2
  else
    template_repo=$(jq -r '.template_repo' "${MANIFEST_FILE}")
  fi

  echo "Preparing secure sandbox workspace..." >&2
  TMP_WORKSPACE=$(mktemp -d -t boilerplate-sync-XXXXXX)
  trap cleanup EXIT

  if [[ "${MODE}" == "push" ]]; then
    echo "Cloning remote template repository (depth 1, full checkout for push)..." >&2
    if ! git clone --depth 1 "${template_repo}" "${TMP_WORKSPACE}" >/dev/null 2>&1; then
      echo "Error: Failed to clone template repository at '${template_repo}'." >&2
      echo "       Verify the repository URL and SSH/agent access." >&2
      exit 1
    fi
  else
    echo "Cloning remote template repository recursively (depth 1)..." >&2
    # Fetch without checking out files immediately to keep local checkout sparse/lightweight
    if ! git clone --depth 1 --no-checkout "${template_repo}" "${TMP_WORKSPACE}" >/dev/null 2>&1; then
      echo "Error: Failed to clone template repository at '${template_repo}'." >&2
      echo "       Verify the repository URL and SSH/agent access." >&2
      exit 1
    fi

    # Identify files we need to checkout
    local files_count i remote_path
    files_count=$(jq '.files | length' "${MANIFEST_FILE}")

    # Pre-resolve remote paths while we are still in the local root directory
    local remote_paths=()
    for ((i = 0; i < files_count; i++)); do
      remote_paths+=("$(jq -r ".files[$i].remote" "${MANIFEST_FILE}")")
    done

    echo "Checking out tracked files sparsely..." >&2
    cd "${TMP_WORKSPACE}"
    for remote_path in "${remote_paths[@]}"; do
      # Force Git to sparse checkout the specific target remote file path
      git checkout HEAD -- "${remote_path}" >/dev/null 2>&1 || true
    done
    cd - >/dev/null
  fi
}

# ==============================================================================
# OPERATIONAL COMMANDS
# ==============================================================================

# Compare local file to remote template file and output line differences
run_diff() {
  local files_count i local_path remote_path full_remote_path
  files_count=$(jq '.files | length' "${MANIFEST_FILE}")

  echo "=============================================================="
  echo "🔍 COMPARATIVE BLUEPRINT DIFF (Local vs Remote Template)"
  echo "=============================================================="

  local exit_code=0
  for ((i = 0; i < files_count; i++)); do
    local_path=$(jq -r ".files[$i].local" "${MANIFEST_FILE}")
    remote_path=$(jq -r ".files[$i].remote" "${MANIFEST_FILE}")
    full_remote_path="${TMP_WORKSPACE}/${remote_path}"

    if [[ ! -f "${full_remote_path}" ]]; then
      echo "⚠️  [NOT FOUND IN REMOTE] Remote source '${remote_path}' missing in template repo for '${local_path}'."
      exit_code=1
      continue
    fi

    if [[ ! -f "${local_path}" ]]; then
      echo "❌ [MISSING LOCALLY] Local file '${local_path}' does not exist."
      echo "   ---> To retrieve: Run sync-boilerplate.sh --pull"
      exit_code=1
      continue
    fi

    if diff -u "${local_path}" "${full_remote_path}" >/dev/null; then
      echo "✅ [IN SYNC] '${local_path}' is identical to remote boilerplate."
    else
      echo "⚠️  [OUT OF SYNC] '${local_path}' has drifted from template:"
      diff -u "${local_path}" "${full_remote_path}" || true
      exit_code=1
    fi
    echo "--------------------------------------------------------------"
  done

  if [[ $exit_code -eq 0 ]]; then
    echo "🟢 SUCCESS: All configuration and boilerplate files are fully in-sync!"
  else
    echo "🔴 DRIFT DETECTED: Review differences above and run with '--pull' to synchronize."
  fi

  return $exit_code
}

# Pull remote templates to overwrite local workspace files
run_pull() {
  local files_count i local_path remote_path full_remote_path local_dir
  files_count=$(jq '.files | length' "${MANIFEST_FILE}")

  echo "=============================================================="
  echo "📥 PULLING REMOTES (Synchronizing Boilerplate Configuration)"
  echo "=============================================================="

  for ((i = 0; i < files_count; i++)); do
    local_path=$(jq -r ".files[$i].local" "${MANIFEST_FILE}")
    remote_path=$(jq -r ".files[$i].remote" "${MANIFEST_FILE}")
    full_remote_path="${TMP_WORKSPACE}/${remote_path}"

    if [[ ! -f "${full_remote_path}" ]]; then
      echo "⚠️  [SKIPPED] Remote template file '${remote_path}' not found in cloned source."
      continue
    fi

    local_dir=$(dirname "${local_path}")
    if [[ ! -d "${local_dir}" ]]; then
      mkdir -p "${local_dir}"
    fi

    if [[ -f "${local_path}" ]]; then
      if diff -u "${local_path}" "${full_remote_path}" >/dev/null; then
        echo "✅ [UP TO DATE] '${local_path}' already matches template."
        continue
      fi
      echo "🔄 [OVERWRITING] '${local_path}' with remote template..."
    else
      echo "➕ [CREATING] '${local_path}' from remote template..."
    fi

    cp "${full_remote_path}" "${local_path}"
  done

  echo "🟢 SUCCESS: Workspace boilerplate sync pull operation completed successfully!"
}

# Push local changes back to update the centralized template repo
run_push() {
  local files_count i local_path remote_path full_remote_path remote_dir
  files_count=$(jq '.files | length' "${MANIFEST_FILE}")

  # Resolve the current local repo name up front while we are still in our local repository
  local local_repo_name
  local_repo_name=$(basename "$(git rev-parse --show-toplevel 2>/dev/null || pwd)")

  echo "=============================================================="
  echo "📤 PUSHING LOCAL CHANGES (Updating Central Template Repo)"
  echo "=============================================================="

  local copied_count=0
  for ((i = 0; i < files_count; i++)); do
    local_path=$(jq -r ".files[$i].local" "${MANIFEST_FILE}")
    remote_path=$(jq -r ".files[$i].remote" "${MANIFEST_FILE}")
    full_remote_path="${TMP_WORKSPACE}/${remote_path}"

    if [[ ! -f "${local_path}" ]]; then
      echo "⚠️  [SKIPPED] Local file '${local_path}' does not exist."
      continue
    fi

    # Create directory inside the clone if it doesn't exist
    remote_dir=$(dirname "${full_remote_path}")
    if [[ ! -d "${remote_dir}" ]]; then
      mkdir -p "${remote_dir}"
    fi

    # Only copy if the file is different (to avoid touch/staging unchanged files)
    if [[ -f "${full_remote_path}" ]] && diff -u "${local_path}" "${full_remote_path}" >/dev/null; then
      echo "✅ [UP TO DATE] '${local_path}' is identical to remote boilerplate."
      continue
    fi

    echo "🔄 [COPYING] '${local_path}' into template remote at '${remote_path}'..."
    cp "${local_path}" "${full_remote_path}"
    copied_count=$((copied_count + 1))
  done

  if [[ "${copied_count}" -eq 0 ]]; then
    echo "🟢 All files are already up-to-date in the central repository. Nothing to push."
    return 0
  fi

  # Committing and pushing changes inside the sandbox workspace
  echo "Committing and pushing changes back to central repository..." >&2
  cd "${TMP_WORKSPACE}"

  git add -A

  if git diff --cached --quiet; then
    echo "🟢 No diff staged. Nothing to push."
    cd - >/dev/null
    return 0
  fi

  git commit -m "sync: update boilerplate from ${local_repo_name}" -s

  echo "Pushing changes back securely using your Git credentials..." >&2
  if git push origin HEAD; then
    echo "🟢 SUCCESS: Successfully committed and pushed local changes to central repository!"
  else
    echo "❌ ERROR: Git push failed. Verify write permissions to the central repository."
    cd - >/dev/null
    exit 1
  fi
  cd - >/dev/null
}

# Display compact status of files tracked
run_status() {
  local files_count i local_path remote_path full_remote_path
  files_count=$(jq '.files | length' "${MANIFEST_FILE}")

  printf "%-40s %-40s %s\n" "LOCAL WORKSPACE FILE" "REMOTE TEMPLATE FILE" "SYNC STATUS"
  printf "%-40s %-40s %s\n" "--------------------" "--------------------" "-----------"

  for ((i = 0; i < files_count; i++)); do
    local_path=$(jq -r ".files[$i].local" "${MANIFEST_FILE}")
    remote_path=$(jq -r ".files[$i].remote" "${MANIFEST_FILE}")
    full_remote_path="${TMP_WORKSPACE}/${remote_path}"

    local status="UNKNOWN"
    if [[ ! -f "${full_remote_path}" ]]; then
      status="MISSING IN REMOTE"
    elif [[ ! -f "${local_path}" ]]; then
      status="MISSING LOCALLY"
    elif diff -u "${local_path}" "${full_remote_path}" >/dev/null; then
      status="IN SYNC"
    else
      status="OUT OF SYNC"
    fi

    printf "%-40s %-40s %s\n" "${local_path}" "${remote_path}" "${status}"
  done
}

# ==============================================================================
# MAIN ENTRY POINT
# ==============================================================================

main() {
  # Parse arguments
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        MODE="help"
        shift
        ;;
      -d|--diff)
        MODE="diff"
        shift
        ;;
      -p|--pull)
        MODE="pull"
        shift
        ;;
      -u|--push)
        MODE="push"
        shift
        ;;
      -s|--status)
        MODE="status"
        shift
        ;;
      *)
        echo "Error: Unknown option '$1'" >&2
        show_help >&2
        exit 1
        ;;
    esac
  done

  if [[ "${MODE}" == "help" ]]; then
    show_help
    exit 0
  fi

  # Execute operation
  validate_environment
  clone_template_repo

  case "${MODE}" in
    diff)
      run_diff
      ;;
    pull)
      run_pull
      ;;
    push)
      run_push
      ;;
    status)
      run_status
      ;;
  esac
}

main "$@"
