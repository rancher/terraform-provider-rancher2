#!/usr/bin/env bash
#
# Skill: generate-otp.sh
# Description: Generates a secure, cryptographically random One-Time Pad (OTP) token.
#              Writes the token safely to disk with strict 0600 permissions and outputs the token.
# Conforms to shell-scripts.instructions.md guidelines.
# Usage: 
#   OTP_TOKEN=$(bash .agent/skills/generate-otp.sh)

set -euo pipefail

# Helper to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Display script help usage instructions
show_help() {
  cat <<EOF
Usage: generate-otp.sh [options]

Generates a secure, cryptographically random 16-byte (32-character) One-Time Pad (OTP) token.
The token is securely saved to ~/.gemini/tmp/terraform-provider-rancher2/active-otp.token with 0600 permissions.
The generated token hex string is printed directly to standard output.

Options:
  -h, --help           Show this message and exit.

Examples:
  OTP_TOKEN=\$(bash .agent/skills/generate-otp.sh)
EOF
}

main() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      *)
        echo "Error: Unknown argument '$1'" >&2
        show_help >&2
        exit 1
        ;;
    esac
  done

  if [[ -z "${HOME:-}" ]]; then
    echo "Error: HOME environment variable is not set or is empty." >&2
    exit 1
  fi

  if ! command_exists openssl; then
    echo "Error: 'openssl' utility is required for secure random generation. Please install openssl." >&2
    exit 1
  fi

  local target_dir="$HOME/.gemini/tmp/terraform-provider-rancher2"
  local token_file="${target_dir}/active-otp.token"

  # Generate 16 bytes of cryptographically secure randomness
  local otp_token
  otp_token=$(openssl rand -hex 16)

  # Create directories safely
  mkdir -p "$target_dir"

  # Write to a freshly-created temp file, then atomically rename it into place, rather than
  # rm -f followed by a redirect. The rm-then-write pattern leaves a TOCTOU window: a symlink
  # could be recreated at $token_file between the two steps, and the write would follow it.
  # mv -f on the same filesystem uses rename(2), which replaces whatever is at the destination
  # atomically without ever dereferencing an existing symlink there.
  local old_umask
  old_umask="$(umask)"
  umask 077
  local tmp_file
  tmp_file="$(mktemp "${target_dir}/.otp.XXXXXX")"
  umask "$old_umask"

  echo "$otp_token" > "$tmp_file"
  mv -f "$tmp_file" "$token_file"

  # Output the token string natively to stdout so calling scripts can bind it to a variable
  echo "$otp_token"
}

main "$@"
