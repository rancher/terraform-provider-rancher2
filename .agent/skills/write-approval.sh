#!/usr/bin/env bash
#
# Skill: write-approval.sh
# Description: Securely writes and verifies the programmatic proactive review approval JSON file.
# Conforms to shell-scripts.instructions.md guidelines.
# Usage: 
#   Write Mode:  .agent/skills/write-approval.sh -t TOKEN -d DIFF_HASH -m "MESSAGE"
#   Verify Mode: .agent/skills/write-approval.sh --verify

set -euo pipefail

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
Usage: write-approval.sh [options] -t TOKEN -d DIFF_HASH -m "MESSAGE"
   or  write-approval.sh --verify

Securely writes and verifies the programmatic proactive review approval JSON file.

Options:
  -h, --help           Show this message and exit.
  --verify             Verify the active proactive review approval file.
  -t TOKEN             The secure One-Time Pad (OTP) token (Required for write).
  -d DIFF_HASH         The cryptographic SHA-256 diff hash of active changes (Required for write).
  -m MESSAGE           The review approval message (Required for write).

Examples:
  .agent/skills/write-approval.sh -t "token" -d "sha256" -m "PR Review status: 🟢 PERFECT - 0 findings."
  .agent/skills/write-approval.sh --verify
EOF
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

# ==============================================================================
# OPERATION STAGES
# ==============================================================================

verify_approval() {
  local approval_file="$1"
  echo "Verifying proactive review approval status..." >&2

  if [[ ! -f "$approval_file" ]]; then
    echo "Error: Proactive review approval not found!" >&2
    echo "       In accordance with 'development-process.md', you MUST delegate" >&2
    echo "       a proactive review of your changes to our specialized subagent" >&2
    echo "       before committing: @review_agent" >&2
    exit 1
  fi

  if [[ -L "$approval_file" ]]; then
    echo "Error: Proactive review approval file at '$approval_file' is a symbolic link." >&2
    echo "       Symlink-based approval files are prohibited for security." >&2
    exit 1
  fi

  local owner_uid
  owner_uid=$(get_file_owner_uid "$approval_file")
  if [[ "$owner_uid" != "$UID" ]]; then
    echo "Error: Proactive review approval file at '$approval_file' is not owned by the current user (UID: $UID, Owner: $owner_uid)." >&2
    echo "       This is a security violation. Please ensure the file is generated securely." >&2
    exit 1
  fi

  if ! command_exists jq; then
    echo "Error: 'jq' utility is required to parse review approval file. Please install jq." >&2
    exit 1
  fi

  local status
  status=$(jq -r '.status' "$approval_file" 2>/dev/null || true)
  local approval_hash
  approval_hash=$(jq -r '.diff_hash' "$approval_file" 2>/dev/null || true)

  if [[ "$status" != "approved" ]]; then
    echo "Error: The proactive review approval status is '$status' (not approved)." >&2
    exit 1
  fi

  # Recalculate the active local diff hash securely using SHA-256 (staged + unstaged combined)
  local active_hash
  active_hash=$(git diff HEAD | calculate_sha256)

  if [[ "$approval_hash" != "$active_hash" ]]; then
    echo "Error: Local changes have been modified since your last proactive review approval!" >&2
    echo "       Approved SHA-256 hash: ${approval_hash}" >&2
    echo "       Current active SHA-256 hash: ${active_hash}" >&2
    echo "       Please re-run the review agent on your latest changes: @review_agent" >&2
    exit 1
  fi

  echo "✅ Proactive review approval verified! (SHA-256 Hash: $active_hash)"
  exit 0
}

write_approval() {
  local token="$1"
  local diff_hash="$2"
  local message="$3"
  local target_dir="$4"
  local approval_file="$5"
  local token_file="$6"

  if [[ -z "$token" || -z "$diff_hash" || -z "$message" ]]; then
    echo "Error: Missing required arguments for write mode." >&2
    show_help >&2
    exit 1
  fi

  # ----------------------------------------------------------------------------
  # SECURE TOKEN VALIDATION (Anti-Spoofing Symlink/Owner checks)
  # ----------------------------------------------------------------------------
  if [[ ! -f "$token_file" ]]; then
    echo "Error: No active verification token found on disk. Proactive review is unauthorized." >&2
    exit 1
  fi

  if [[ -L "$token_file" ]]; then
    echo "Error: Token file at '$token_file' is a symbolic link. Prohibited for security." >&2
    exit 1
  fi

  local token_uid
  token_uid=$(get_file_owner_uid "$token_file")
  if [[ "$token_uid" != "$UID" ]]; then
    echo "Error: Token file is not owned by current user (UID: $UID, Owner: $token_uid). Security violation." >&2
    exit 1
  fi

  local active_token
  active_token=$(tr -d ' \n' < "$token_file")
  if [[ "$token" != "$active_token" || ${#active_token} -lt 16 ]]; then
    echo "Error: Invalid verification token. Proactive review signature rejected." >&2
    exit 1
  fi

  # ----------------------------------------------------------------------------
  # SECURE APPROVAL GENERATION
  # ----------------------------------------------------------------------------
  
  # Ensure target diff_hash exactly matches active diff (prevents hash mismatch spoofing)
  local active_hash
  active_hash=$(git diff HEAD | calculate_sha256)
  if [[ "$diff_hash" != "$active_hash" ]]; then
    echo "Error: Requested hash '$diff_hash' does not match current active hash '$active_hash'." >&2
    exit 1
  fi

  if ! command_exists jq; then
    echo "Error: 'jq' utility is required for secure JSON generation. Please install jq." >&2
    exit 1
  fi

  # Success! Construct the approval JSON securely with umask 077 (0600 permissions)
  mkdir -p "$target_dir"

  # Write to a freshly-created temp file, then atomically rename it into place, rather than
  # `rm -f` followed by a plain `> file` redirect. The rm-then-write pattern leaves a TOCTOU
  # window: a symlink could be recreated at $approval_file between the two steps, and the
  # write would silently follow it. `mv` on the same filesystem uses rename(2), which replaces
  # whatever is at the destination atomically without ever dereferencing an existing symlink
  # there, so a pre-existing symlink can't redirect the write.
  local old_umask
  old_umask="$(umask)"
  umask 077
  local tmp_file
  tmp_file="$(mktemp "${target_dir}/.approval.XXXXXX")"
  umask "$old_umask"

  jq -n \
    --arg status "approved" \
    --arg msg "${message}" \
    --arg sha "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" \
    --arg hash "${diff_hash}" \
    '{status: $status, message: $msg, commit_sha: $sha, diff_hash: $hash}' > "$tmp_file"
  mv -f "$tmp_file" "$approval_file"

  # Single-use token: Immediately destroy the active token to prevent replay
  rm -f "$token_file"

  echo "✅ Proactive review approval file successfully generated at: $approval_file"
}

# ==============================================================================
# MAIN ENTRY POINT
# ==============================================================================

main() {
  if [[ -z "${HOME:-}" ]]; then
    echo "Error: HOME environment variable is not set or is empty." >&2
    exit 1
  fi

  local token=""
  local diff_hash=""
  local message=""
  local verify_only=false

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      --verify)
        verify_only=true
        shift
        ;;
      -t)
        if [[ $# -lt 2 ]]; then
          echo "Error: Option -t requires an argument." >&2
          show_help >&2
          exit 1
        fi
        token="$2"
        shift 2
        ;;
      -d)
        if [[ $# -lt 2 ]]; then
          echo "Error: Option -d requires an argument." >&2
          show_help >&2
          exit 1
        fi
        diff_hash="$2"
        shift 2
        ;;
      -m)
        if [[ $# -lt 2 ]]; then
          echo "Error: Option -m requires an argument." >&2
          show_help >&2
          exit 1
        fi
        message="$2"
        shift 2
        ;;
      *)
        echo "Error: Unknown argument '$1'" >&2
        show_help >&2
        exit 1
        ;;
    esac
  done

  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: This command must be run inside a Git repository." >&2
    exit 1
  fi

  local target_dir="$HOME/.gemini/tmp/terraform-provider-rancher2"
  local approval_file="${target_dir}/review-approval.json"
  local token_file="${target_dir}/active-otp.token"

  if [[ "$verify_only" == "true" ]]; then
    verify_approval "$approval_file"
  else
    write_approval "$token" "$diff_hash" "$message" "$target_dir" "$approval_file" "$token_file"
  fi
}

main "$@"
