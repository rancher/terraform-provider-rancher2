#!/usr/bin/env bash

set -e

echo "==> Running acceptance testing..."

source $(dirname $0)/gotestacc_vars.sh

cleanup()
{
    $(dirname $0)/cleanup_testacc.sh
}
trap cleanup EXIT TERM

source $(dirname $0)/start_k3s.sh

K3S_SERVER=${K3S_SERVER:-""}
K3S_SERVER_IP=${K3S_SERVER_IP:-"$(${DOCKER_BIN} inspect ${K3S_SERVER} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')"}
RANCHER_VERSION=${RANCHER_VERSION:-"v2.3.6"}

source $(dirname $0)/start_rancher.sh

if [ "${DOCKERIZED}" == "true" ]; then
  echo "${RANCHER_IP} ${RANCHER_HOSTNAME}" >> /etc/hosts
fi

RANCHER_URL=${RANCHER_URL:-""}
RANCHER_TOKEN_KEY=${RANCHER_TOKEN_KEY:-""}
RANCHER_INSECURE=${RANCHER_INSECURE:-true}
RANCHER_ACC_CLUSTER_NAME=${RANCHER_ACC_CLUSTER_NAME:-"local"}
RANCHER_ADMIN_PASS=${RANCHER_ADMIN_PASS:-"admin"}
RANCHER_BOOTSTRAP=${RANCHER_BOOTSTRAP:-false}

echo Running acceptance tests

#PACKAGES="$(find . -name '*.go' | xargs -I{} dirname {} |  cut -f2 -d/ | sort -u | grep -Ev '(^\.$|.git|vendor|bin)' | sed -e 's!^!./!' -e 's!$!/...!')"
TF_ACC=1 go test -cover -tags=test ./rancher2/... -v -timeout 120m
