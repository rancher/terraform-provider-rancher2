#!/usr/bin/env bash

set -e

echo "==> Running acceptance testing..."

cleanup()
{
    $(dirname $0)/cleanup_testacc.sh
}
trap cleanup EXIT TERM

source $(dirname $0)/start_testacc.sh

RANCHER_URL=${RANCHER_URL:-""}
RANCHER_TOKEN_KEY=${RANCHER_TOKEN_KEY:-""}
RANCHER_ACCESS_KEY=${RANCHER_ACCESS_KEY:-""}
RANCHER_SECRET_KEY=${RANCHER_SECRET_KEY:-""}
RANCHER_INSECURE=${RANCHER_INSECURE:-true}
RANCHER_ACC_CLUSTER_NAME=${RANCHER_ACC_CLUSTER_NAME:-"local"}
RANCHER_ADMIN_PASS=${RANCHER_ADMIN_PASS:-""}
RANCHER_BOOTSTRAP=${RANCHER_BOOTSTRAP:-false}

echo Running acceptance tests

PACKAGES="$(find . -name '*.go' | xargs -I{} dirname {} |  cut -f2 -d/ | sort -u | grep -Ev '(^\.$|.git|vendor|bin)' | sed -e 's!^!./!' -e 's!$!/...!')"
TF_ACC=1 go test -cover -tags=test ${PACKAGES} -v -timeout 120m
