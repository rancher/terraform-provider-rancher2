#!/usr/bin/env bash

set -x

TESTACC_TEMP_DIR=${TESTACC_TEMP_DIR:-"$(dirname $0)/tmp"}
TESTACC_DOCKER_LIST_NAME=${TESTACC_DOCKER_LIST_NAME:-"testacc_docker_ids"}
TESTACC_DOCKER_LIST=${TESTACC_TEMP_DIR}"/"${TESTACC_DOCKER_LIST_NAME}
TESTACC_K3S_KUBECONFIG_NAME=${TESTACC_K3S_KUBECONFIG_NAME:-"testacc_kubeconfig.yaml"}
TESTACC_K3S_KUBECONFIG=${TESTACC_TEMP_DIR}"/"${TESTACC_K3S_KUBECONFIG_NAME}

DOCKER_BIN=${DOCKER_BIN:-$(which docker)}


if [ -f ${TESTACC_DOCKER_LIST} ]; then
	echo Cleaning up testacc docker list ${TESTACC_DOCKER_LIST}
	${DOCKER_BIN} rm -fv $(cat ${TESTACC_DOCKER_LIST})
	rm ${TESTACC_DOCKER_LIST}
fi

if [ -f ${TESTACC_K3S_KUBECONFIG} ]; then
	echo Cleaning up testacc k3s kubeconfig ${TESTACC_K3S_KUBECONFIG}
	rm ${TESTACC_K3S_KUBECONFIG}
fi

if [ -d ${TESTACC_TEMP_DIR} ] && [ "${TESTACC_TEMP_DIR}" != "/tmp" ]; then
	echo Cleaning up testacc temporary dir ${TESTACC_TEMP_DIR}
  	rm -r ${TESTACC_TEMP_DIR}
fi
