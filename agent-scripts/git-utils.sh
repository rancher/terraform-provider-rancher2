#!/usr/bin/env bash
#
# Description: Modular Git Utility shell library.
#

# Helper to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Retrieve file owner uid in a cross-platform manner
get_file_owner_uid() {
  local file="$1"
  stat -c %u "$file" 2>/dev/null || stat -f %u "$file" 2>/dev/null || echo ""
}

# Calculate SHA-256 hash of standard input
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
