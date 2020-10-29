#!/usr/bin/env bash

set -x

source $(dirname $0)/gotestacc_vars.sh

if [ -f ${DOCKER_LIST} ]; then
	echo Cleaning up k3s docker list ${DOCKER_LIST}
	${DOCKER_BIN} rm -fv $(cat ${DOCKER_LIST})
	rm ${DOCKER_LIST}
fi

if [ -d ${TEMP_DIR} ] && [ "${TEMP_DIR}" != "/tmp" ]; then
	echo Cleaning up testacc temporary dir ${TEMP_DIR}
  	rm -rf ${TEMP_DIR}
fi
