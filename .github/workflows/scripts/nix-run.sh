#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${NIX_SSL_CERT_FILE:-}" ]]; then
  for cert in /etc/ssl/certs/ca-certificates.crt \
              /etc/ssl/certs/ca-bundle.crt \
              /etc/pki/tls/certs/ca-bundle.crt \
              /etc/ssl/ca-bundle.pem \
              /var/lib/ca-certificates/ca-bundle.pem; do
    if [[ -f "${cert}" ]]; then
      export NIX_SSL_CERT_FILE="${cert}"
      break
    fi
  done
fi

export SSL_CERT_FILE="${NIX_SSL_CERT_FILE:-}"
export CURL_CA_BUNDLE="${NIX_SSL_CERT_FILE:-}"

{
  echo "git config --global --add safe.directory \"$PWD\""
  printf "exec "
  printf "%q " "$@"
  printf "\n"
} > .nix-script.sh

trap 'rm -f .nix-script.sh' EXIT

# Ensure the suse user can read/write the script and current directory
chown -R suse:suse . || true

# Ensure parent directories are traversable by the suse user
p="$PWD"
while [[ "${p}" != "/" && -n "${p}" ]]; do
  chmod a+rx "${p}" 2>/dev/null || true
  p="$(dirname "${p}")"
done

sudo -E -u suse /home/suse/.nix-profile/bin/nix develop \
  --ignore-environment \
  --extra-experimental-features nix-command \
  --extra-experimental-features flakes \
  --keep NIX_SSL_CERT_FILE \
  --keep SSL_CERT_FILE \
  --keep CURL_CA_BUNDLE \
  --keep NIX_ENV_LOADED \
  --keep TERM \
  --keep HOME \
  --keep SSH_AUTH_SOCK \
  --keep GITHUB_TOKEN \
  --keep GITHUB_OWNER \
  --keep GITHUB_WORKSPACE \
  --keep AWS_ACCESS_KEY_ID \
  --keep AWS_SECRET_ACCESS_KEY \
  --keep AWS_SESSION_TOKEN \
  --keep AWS_ROLE \
  --keep AWS_REGION \
  --keep AWS_DEFAULT_REGION \
  --keep IDENTIFIER \
  --keep ZONE \
  --keep ACME_SERVER_URL \
  --keep PR_NUMBER \
  --keep TAG \
  --keep GPG_PASSPHRASE \
  --keep GPG_KEY_ID \
  --keep GPG_KEY \
  --command bash -e .nix-script.sh
