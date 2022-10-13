#!/usr/bin/env bash

set -e

echo "==> Running dockerized acceptance testing..."

# Setting docker
DOCKER_NAME=docker
DOCKER_URL="https://download.docker.com/linux/static/stable/x86_64/docker-20.10.17.tgz"
DOCKER_BIN=$(which ${DOCKER_NAME} || echo none)
if [ "${DOCKER_BIN}" == "none" ] ; then
  export DOCKER_BIN=${TESTACC_TEMP_DIR}/${DOCKER_NAME}
  if [ ! -x "${DOCKER_BIN}" ]; then
    curl -sL ${DOCKER_URL} | tar -xzf - 
    mv docker/docker ${DOCKER_BIN} && rm -rf docker
    chmod 755 ${DOCKER_BIN}
  fi
fi

BUILDER_TAG=${BUILDER_TAG:-"terraform-provider-rancher2_builder"}

${DOCKER_BIN} build -t ${BUILDER_TAG} -f $(dirname $0)/Dockerfile.builder .

${DOCKER_BIN} run -i --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD:/go/src/github.com/terraform-providers/terraform-provider-rancher2 \
  -e DOCKERIZED=true \
  -e EXPOSE_HOST_PORTS=${EXPOSE_HOST_PORTS} \
  ${BUILDER_TAG} make testacc
