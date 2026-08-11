#!/usr/bin/env bash
#
# Skill: run-acc-test.sh
# Description: Safely runs a specific Terraform provider acceptance test within the Nix environment.
#              It explicitly keeps necessary AWS and Terraform environment variables so tests don't fail.
# Usage: ./run-acc-test.sh <TestName>

set -euo pipefail

if [ -z "${1:-}" ]; then
  echo "Error: Must provide a test name."
  echo "Usage: $0 <TestName>"
  exit 1
fi

TEST_NAME="$1"

echo "Executing acceptance test: ${TEST_NAME} inside Nix environment..."

nix develop --ignore-environment \
  --extra-experimental-features nix-command \
  --extra-experimental-features flakes \
  --keep HOME --keep SSH_AUTH_SOCK \
  --keep GPG_SIGNING_KEY \
  --keep NIX_SSL_CERT_FILE \
  --keep NIX_ENV_LOADED \
  --keep TERM \
  --keep AWS_ROLE \
  --keep AWS_REGION \
  --keep AWS_DEFAULT_REGION \
  --keep AWS_ACCESS_KEY_ID \
  --keep AWS_SECRET_ACCESS_KEY \
  --keep AWS_SESSION_TOKEN \
  --keep TF_VAR_aws_access_key_id \
  --keep TF_VAR_aws_secret_access_key \
  --keep TF_VAR_aws_session_token \
  --keep TF_VAR_aws_region \
  --command bash -c "./run_tests.sh -t ${TEST_NAME}"
