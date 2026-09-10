#!/usr/bin/env bash
#
# Imports the release signing key and runs GoReleaser.
#
# This MUST run inside nix-run.sh. nix-run.sh drops privileges to the 'suse' user,
# so a GPG key imported by an earlier workflow step (which runs as root) lands in
# root's keyring and is invisible here. Importing in the same shell that invokes
# goreleaser is what keeps the key and the signer in the same GnuPG home.
#
# Environment (all must be on the --keep allowlist in nix-run.sh):
#   GPG_KEY          - ASCII-armored private key                    (required)
#   GPG_KEY_ID       - key id, used by goreleaser as --local-user   (required)
#   GPG_PASSPHRASE   - passphrase for the private key               (required)
#   WORKING_DIR      - directory to release from, relative or absolute (optional)
#   GORELEASER_CONFIG- config path, relative to WORKING_DIR (default .goreleaser.yml)
#   SKIP_VALIDATE    - "true" to pass --skip=validate               (optional)

set -euo pipefail

cleanup() {
  # clear history just in case
  history -c || true
}
trap cleanup EXIT TERM

# Support releasing from a tag checked out into a subdirectory (manual releases)
if [[ -n "${WORKING_DIR:-}" ]]; then
  cd "${WORKING_DIR}"
fi

# sanitize variables
if [[ -z "${GPG_PASSPHRASE:-}" ]]; then echo "Error: gpg passphrase empty" >&2; exit 1; fi
if [[ -z "${GPG_KEY_ID:-}" ]]; then echo "Error: key id empty" >&2; exit 1; fi
if [[ -z "${GPG_KEY:-}" ]]; then echo "Error: key contents empty" >&2; exit 1; fi

# Trim whitespace/newlines so the key id matches what GPG reports
export GPG_KEY_ID
GPG_KEY_ID=$(echo -n "${GPG_KEY_ID}" | tr -d '[:space:]')

echo "Importing gpg key"
echo "${GPG_KEY}" | gpg --import --batch > /dev/null || { echo "Error: Failed to import GPG key" >&2; exit 1; }

# Prefer the secret primary key id actually present in the keyring. Vault may hold
# the id of an encryption subkey, which cannot sign. '|| true' avoids tripping
# pipefail when grep matches nothing.
SEC_LINE=$(gpg --batch --list-secret-keys --keyid-format LONG | grep -E '^sec' | head -n1 || true)
if [[ -n "${SEC_LINE}" ]]; then
  DETECTED_KEY_ID=$(echo "${SEC_LINE}" | awk '{print $2}' | cut -d'/' -f2)
  if [[ -n "${DETECTED_KEY_ID}" ]]; then
    echo "Detected gpg key id from imported key: ${DETECTED_KEY_ID}"
    GPG_KEY_ID="${DETECTED_KEY_ID}"
  fi
fi

# https://www.gnupg.org/documentation/manuals/gnupg24/gpg.1.html
# https://goreleaser.com/customization/sign/sign/
# troubleshooting information
gpg --version
# this only lists UIDs, no secret material
gpg --batch --list-secret-keys --keyid-format LONG
# this fails loudly here rather than deep inside a 17 minute goreleaser run
gpg --batch --list-secret-keys --keyid-format LONG "${GPG_KEY_ID}" >/dev/null
# troubleshooting information
goreleaser --version

CONFIG_FILE="${GORELEASER_CONFIG:-.goreleaser.yml}"
if [[ ! -f "${CONFIG_FILE}" ]]; then
  echo "Error: ${CONFIG_FILE} not found (working directory: $(pwd))" >&2
  exit 1
fi

extra_args=()
if [[ "${SKIP_VALIDATE:-false}" == "true" ]]; then
  extra_args+=("--skip=validate")
fi

goreleaser release --clean --config "${CONFIG_FILE}" "${extra_args[@]+"${extra_args[@]}"}"
