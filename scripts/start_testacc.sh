#!/usr/bin/env bash

set -x

# Setting temporary directory
export TESTACC_TEMP_DIR=${TESTACC_TEMP_DIR:-$(dirname $0)"/tmp"}
## Converting to absolute path if needed
case $TESTACC_TEMP_DIR in
     /*) ;;
     *) TESTACC_TEMP_DIR=$(pwd)/${TESTACC_TEMP_DIR} ;;
esac

if [ ! -d ${TESTACC_TEMP_DIR} ]; then
  mkdir ${TESTACC_TEMP_DIR}
fi

# Setting default vars
export TESTACC_DOCKER_LIST_NAME=${TESTACC_DOCKER_LIST_NAME:-"testacc_docker_ids"}
TESTACC_DOCKER_LIST=${TESTACC_TEMP_DIR}"/"${TESTACC_DOCKER_LIST_NAME}

TESTACC_EXPOSE_HOST_PORTS=${TESTACC_EXPOSE_HOST_PORTS:-false}

export TESTACC_K3S_KUBECONFIG_NAME=${TESTACC_K3S_KUBECONFIG_NAME:-"testacc_kubeconfig.yaml"}
TESTACC_K3S_KUBECONFIG=${TESTACC_TEMP_DIR}"/"${TESTACC_K3S_KUBECONFIG_NAME}
TESTACC_K3S_PORT=${TESTACC_K3S_PORT:-6443}
TESTACC_K3S_SECRET=${TESTACC_K3S_SECRET:-"somethingtotallyrandom"}
TESTACC_K3S_VERSION=${TESTACC_K3S_VERSION:-"v1.18.8-k3s1"}

TESTACC_RANCHER_PORT=${TESTACC_RANCHER_PORT:-44443}
TESTACC_RANCHER_VERSION=${TESTACC_RANCHER_VERSION:-"v2.4.8"}

# Download required software if not available
## jq
JQ_NAME=jq
JQ_URL="https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64"
JQ_BIN=$(which ${JQ_NAME} || echo none)
if [ "${JQ_BIN}" == "none" ] ; then
  echo Downloading ${JQ_NAME}
  JQ_BIN=${TESTACC_TEMP_DIR}/${JQ_NAME}
  curl -sL ${JQ_URL} -o ${JQ_BIN}
  chmod 755 ${JQ_BIN}
fi
## kubectl
KUBECTL_NAME=kubectl
KUBECTL_URL="https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
KUBECTL_BIN=$(which ${KUBECTL_NAME} || echo none)
if [ "${KUBECTL_BIN}" == "none" ] ; then
  echo Downloading ${KUBECTL_NAME}
  KUBECTL_BIN=${TESTACC_TEMP_DIR}/${KUBECTL_NAME}
  curl -sL ${KUBECTL_URL} -o ${KUBECTL_BIN}
  chmod 755 ${KUBECTL_BIN}
fi
## docker
DOCKER_NAME=docker
DOCKER_URL="https://download.docker.com/linux/static/stable/x86_64/docker-20.10.17.tgz"
DOCKER_BIN=$(which ${DOCKER_NAME} || echo none)
if [ "${DOCKER_BIN}" == "none" ] ; then
  export DOCKER_BIN=${TESTACC_TEMP_DIR}/${DOCKER_NAME}
  curl -sL ${DOCKER_URL} | tar -xzf - 
  mv docker/docker ${DOCKER_BIN} && rm -rf docker
  chmod 755 ${DOCKER_BIN}
fi

# Setting exposed ports
if [ "${TESTACC_EXPOSE_HOST_PORTS}" == "true" ]; then
  rancher_exposed_port="-p ${TESTACC_RANCHER_PORT}:${TESTACC_RANCHER_PORT}"
  k3s_exposed_port="-p ${TESTACC_K3S_PORT}:${TESTACC_K3S_PORT}"
fi

# Starting rancher server
rancher_server=$(${DOCKER_BIN} run -d \
  ${rancher_exposed_port} \
  rancher/rancher:${TESTACC_RANCHER_VERSION} --https-listen-port=${TESTACC_RANCHER_PORT})
rancher_server_ip=$(${DOCKER_BIN} inspect ${rancher_server} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
echo ${rancher_server} >> ${TESTACC_DOCKER_LIST}

export RANCHER_URL=https://${rancher_server_ip}:${TESTACC_RANCHER_PORT}
export RANCHER_INSECURE=true

# Starting k3s cluster
k3s_imported_server=$(${DOCKER_BIN} run -d \
  ${k3s_exposed_port} \
  --privileged \
  --tmpfs /run \
  --tmpfs /var/run \
  -e K3S_CLUSTER_SECRET=${TESTACC_K3S_SECRET} \
  -e K3S_KUBECONFIG_OUTPUT=/tmp/${TESTACC_K3S_KUBECONFIG_NAME} \
  -e K3S_KUBECONFIG_MODE=666 \
  rancher/k3s:${TESTACC_K3S_VERSION} server --https-listen-port ${TESTACC_K3S_PORT})
echo ${k3s_server} >> ${TESTACC_DOCKER_LIST}
k3s_imported_server_ip=$(${DOCKER_BIN} inspect ${k3s_server} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
k3s_imported_node=$(${DOCKER_BIN} run -d \
  --privileged \
  --tmpfs /run \
  --tmpfs /var/run \
  -e K3S_URL=https://${k3s_server_ip}:${TESTACC_K3S_PORT} \
  -e K3S_CLUSTER_SECRET=${TESTACC_K3S_SECRET} \
  rancher/k3s:${TESTACC_K3S_VERSION})
echo ${k3s_imported_node} >> ${TESTACC_DOCKER_LIST}

export RANCHER_ACC_CLUSTER_NAME=bootstrap-imported-k3s-cluster

# Show running dockers
${DOCKER_BIN} ps

# Rancher bootstrap
## Waiting for rancher server start
while [ "${rancher_ready}" != "pong" ]; do
  sleep 5
  rancher_ready=$(${DOCKER_BIN} exec -i ${rancher_server} curl -sk ${RANCHER_URL}/ping || echo starting)
done

## Resetting rancher admin password
rancher_password=$(${DOCKER_BIN} exec -i ${rancher_server} reset-password | grep -v '^New'|tr -d '\r')
export RANCHER_ADMIN_PASS=${rancher_password}

## Admin login in rancher server
login_token=$(${DOCKER_BIN} exec -i ${rancher_server} curl -X POST -sk "${RANCHER_URL}/v3-public/localProviders/local?action=login" \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    --data-binary "{\"password\":\"${rancher_password}\",\"username\":\"admin\"}" | ${JQ_BIN} -r .token)

## Getting admin token 
rancher_token=$(${DOCKER_BIN} exec -i ${rancher_server} curl -sSk \
  "${RANCHER_URL}/v3/token" \
  -H 'content-type: application/json' \
  -H "Authorization: Bearer ${login_token}" \
  --data-binary '{"type":"token","description":"automation","name":""}' | ${JQ_BIN} -r '.token')

export RANCHER_TOKEN_KEY=${rancher_token}

## Setting rancher server-url
${DOCKER_BIN} exec -i ${rancher_server} curl -X PUT -sk "${RANCHER_URL}/v3/settings/server-url" \
    -H 'Accept: application/json' \
    -H "Authorization: Bearer ${rancher_token}" \
    -H 'Content-Type: application/json' \
    --data-binary "{\"name\": \"server-url\", \"value\":\"${RANCHER_URL}\"}"

# Creating imported cluster
cluster_json=$(cat <<-EOM
{
  "type": "cluster",
  "dockerRootDir": "/var/lib/docker",
  "enableNetworkPolicy": "false",
  "name": "${RANCHER_ACC_CLUSTER_NAME}"
}
EOM
)

cluster_obj=$(${DOCKER_BIN} exec -i ${rancher_server} curl -sk -X POST -H "Authorization: Bearer ${rancher_token}" \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    ${RANCHER_URL}/v3/clusters \
    -d "${cluster_json}")

cluster_id=$(echo $cluster_obj | ${JQ_BIN} -r '.id')
registration_link=$(echo $cluster_obj | ${JQ_BIN} -r '.links.clusterRegistrationTokens')

# Creating cluster registration token
${DOCKER_BIN} exec -i ${rancher_server} curl -sk -X POST -H "Authorization: Bearer ${rancher_token}" \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    ${RANCHER_URL}/v3/clusterregistrationtoken \
    -d "{\"clusterId\": \"${cluster_id}\", \"type\":\"clusterRegistrationToken\"}"
sleep 1

# Copying kubeconfig
${DOCKER_BIN} cp ${k3s_server}:/tmp/${TESTACC_K3S_KUBECONFIG_NAME} ${TESTACC_K3S_KUBECONFIG}

# Setting kubeconfig and rancher_url if exposed host ports
if [ "${TESTACC_EXPOSE_HOST_PORTS}" == "true" ]; then
  export RANCHER_URL=https://localhost:${TESTACC_RANCHER_PORT}
else
  sed -i -e 's/localhost/'"${k3s_server_ip}"'/g' ${TESTACC_K3S_KUBECONFIG}
fi

# Registering k3s cluster
## Getting manifest
manifest_url=$(${DOCKER_BIN} exec -i ${rancher_server} curl -sk -X GET -H "Authorization: Bearer ${rancher_token}" \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    ${registration_link} | ${JQ_BIN} -r '.data[0].manifestUrl')

## Applying manifest on k3s cluster
${DOCKER_BIN} exec -i ${rancher_server} curl -sfLk ${manifest_url} | KUBECONFIG=${TESTACC_K3S_KUBECONFIG} ${KUBECTL_BIN} apply -f -

## Waiting for k3s cluster becomes active
while [ "${cluster_ready}" != "active" ]; do
  sleep 10
  cluster_ready=$(${DOCKER_BIN} exec -i ${rancher_server} curl -sk -H "Authorization: Bearer ${rancher_token}" \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    ${RANCHER_URL}/v3/clusters/${cluster_id} | ${JQ_BIN} -r '.state')
done

