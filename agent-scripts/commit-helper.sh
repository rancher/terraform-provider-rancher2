#!/usr/bin/env bash
#
# Description: Modular Git commit signature and sign-off helper.
#

execute_commit() {
  local commit_msg="$1"
  local branch="$2"

  echo "Staging changes..." >&2
  # Stage the changes cleanly
  git add -A

  echo "Creating conventional signed commit on branch '$branch'..." >&2

  # Perform GPG-signed or SSH-signed commit natively with mandatory sign-off (-s)
  if ! git commit -S -s -m "$commit_msg"; then
    echo "======================================================================" >&2
    echo "❌ GPG/SSH COMMIT SIGNATURE FAILURE!" >&2
    echo "   The git commit signature operation failed or was cancelled." >&2
    echo "   Please verify that your GPG key is loaded, and that you have" >&2
    echo "   properly configured your Touch ID biometric hardware binding." >&2
    echo "======================================================================" >&2
    exit 1
  fi

  echo "✅ Conventional GPG/SSH-signed commit successfully created!" >&2
}
