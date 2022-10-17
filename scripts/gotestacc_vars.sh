#!/usr/bin/env bash

set -e

# Setting temporary directory
export TEMP_DIR=${TEMP_DIR:-$(dirname $0)"/tmp"}
## Converting to absolute path if needed
case $TEMP_DIR in
     /*) ;;
     *) TEMP_DIR=$(pwd)/${TEMP_DIR} ;;
esac

if [ ! -d ${TEMP_DIR} ]; then
  mkdir ${TEMP_DIR}
fi

EXPOSE_HOST_PORTS=${EXPOSE_HOST_PORTS:-"false"}

# Setting required software
DOCKER_NAME=docker
DOCKER_URL="https://download.docker.com/linux/static/stable/x86_64/docker-20.10.17.tgz"
DOCKER_BIN=$(which ${DOCKER_NAME} || echo none)
if [ "${DOCKER_BIN}" == "none" ] ; then
  export DOCKER_BIN=${TEMP_DIR}/${DOCKER_NAME}
  if [ ! -x "${TEMP_DIR}/${DOCKER_NAME}" ]; then
    curl -sL ${DOCKER_URL} | tar -xzf - 
    mv docker/docker ${DOCKER_BIN} && rm -rf docker
    chmod 755 ${DOCKER_BIN}
  fi
fi
DOCKER_LIST_NAME=${DOCKER_LIST_NAME:-"docker_ids"}
DOCKER_LIST=${TEMP_DIR}"/"${DOCKER_LIST_NAME}
## curl 
CURL_BIN="${DOCKER_BIN} run -i --rm curlimages/curl"
