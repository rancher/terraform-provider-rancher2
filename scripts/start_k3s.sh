#!/usr/bin/env bash

set -x

source $(dirname $0)/gotestacc_vars.sh

# Setting default vars
## k3s
K3S_DOCKER_IMAGE=${K3S_DOCKER_IMAGE:-"rancher/k3s"}
K3S_KUBECONFIG_NAME=${K3S_KUBECONFIG_NAME:-"k3s_kubeconfig.yaml"}
K3S_KUBECONFIG=${TEMP_DIR}"/"${K3S_KUBECONFIG_NAME}
K3S_PORT=${K3S_PORT:-6443}
K3S_INGRESS_PORT=${K3S_INGRESS_PORT:-8080}
K3S_INGRESS_PORT_TLS=${K3S_INGRESS_PORT_TLS:-8443}
K3S_SECRET=${K3S_SECRET:-"somethingtotallyrandom"}
K3S_VERSION=${K3S_VERSION:-"v1.19.10-k3s1"}

# Setting exposed ports
if [ "${EXPOSE_HOST_PORTS}" == "true" ]; then
  k3s_exposed_port="-p ${K3S_PORT}:${K3S_PORT}"
  k3s_ingress_port="-p ${K3S_INGRESS_PORT}:80 -p ${K3S_INGRESS_PORT_TLS}:443"
fi

# Starting k3s 
## server
K3S_SERVER=$(${DOCKER_BIN} run -d \
  ${k3s_exposed_port} \
  ${k3s_ingress_port} \
  --privileged \
  --tmpfs /run \
  --tmpfs /var/run \
  -e K3S_CLUSTER_SECRET=${K3S_SECRET} \
  -e K3S_KUBECONFIG_OUTPUT=/tmp/${K3S_KUBECONFIG_NAME} \
  -e K3S_KUBECONFIG_MODE=666 \
  ${K3S_DOCKER_IMAGE}:${K3S_VERSION} server --https-listen-port ${K3S_PORT})
echo ${K3S_SERVER} >> ${DOCKER_LIST}
K3S_SERVER_IP=$(${DOCKER_BIN} inspect ${K3S_SERVER} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
K3S_URL="https://${K3S_SERVER_IP}:${K3S_PORT}"

## agent
k3s_node=$(${DOCKER_BIN} run -d \
  --privileged \
  --tmpfs /run \
  --tmpfs /var/run \
  -e K3S_URL=${K3S_URL} \
  -e K3S_CLUSTER_SECRET=${K3S_SECRET} \
  ${K3S_DOCKER_IMAGE}:${K3S_VERSION} agent)
echo ${k3s_node} >> ${DOCKER_LIST}
k3s_node_ip=$(${DOCKER_BIN} inspect ${k3s_node} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')

## Waiting for start
while [ "${k3s_ready}" != "pong" ]; do
  sleep 5
  k3s_ready=$(${CURL_BIN} -sk ${K3S_URL}/ping || echo starting rancher k3s cluster)
done
sleep 2

${DOCKER_BIN} cp ${K3S_SERVER}:/tmp/${K3S_KUBECONFIG_NAME} ${K3S_KUBECONFIG}

export K3S_SERVER=${K3S_SERVER}
export K3S_SERVER_IP=${K3S_SERVER_IP}
export K3S_VERSION=${K3S_VERSION}
