#!/usr/bin/env bash

set -x

source $(dirname $0)/gotestacc_vars.sh

# Setting vars
K3S_SERVER=${K3S_SERVER:-""}
if [ ${K3S_SERVER} == "" ]; then
  echo "[ERROR] K3S_SERVER env var should be provided"
  exit 1
fi
K3S_SERVER_IP=${K3S_SERVER_IP:-"$(${DOCKER_BIN} inspect ${K3S_SERVER} -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')"}
K3S_INGRESS_PORT_TLS=${K3S_INGRESS_PORT_TLS:-8443}
KUBECTL_BIN=${KUBECTL_BIN:-"${DOCKER_BIN} exec -i ${K3S_SERVER} kubectl"}

## cert-manager
CERTMANAGER_VERSION=${CERTMANAGER_VERSION:-"v1.5.2"}
CERTMANAGER_CRD=${CERTMANAGER_CRD:-"https://github.com/jetstack/cert-manager/releases/download/${CERTMANAGER_VERSION}/cert-manager.crds.yaml"}
CERTMANAGER_NS=${CERTMANAGER_NS:-"cert-manager"}

## rancher
RANCHER_VERSION=${RANCHER_VERSION:-"v2.6.1"}
RANCHER_NS=${RANCHER_NS:-"cattle-system"}
RANCHER_DNS_DOMAIN="nip.io"
export RANCHER_HOSTNAME="rancher.${K3S_SERVER_IP}.${RANCHER_DNS_DOMAIN}"
export RANCHER_IP=${K3S_SERVER_IP}

# Installing helm charts 
## cert-manager
cat << EOF > ${TEMP_DIR}"/cert-manager.yaml"
---
kind: Namespace
apiVersion: v1
metadata:
  name: ${CERTMANAGER_NS}
  labels:
    app: ${CERTMANAGER_NS}
---
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: cert-manager
  namespace: kube-system
spec:
  chart: cert-manager
  repo: https://charts.jetstack.io
  targetNamespace: ${CERTMANAGER_NS}
  version: ${CERTMANAGER_VERSION}
EOF

${KUBECTL_BIN} apply -f ${CERTMANAGER_CRD}
${DOCKER_BIN} cp ${TEMP_DIR}"/cert-manager.yaml" ${K3S_SERVER}:/var/lib/rancher/k3s/server/manifests/
## waiting for HelmChart cert-manager
while [[ $(${KUBECTL_BIN} -n kube-system get helmchart cert-manager -o 'jsonpath={..spec.version}') != ${CERTMANAGER_VERSION} ]] ; 
do echo "Waiting for HelmChart rancher" && sleep 2;
done
## waiting for helm-install-cert-manager
helm_job=$(${KUBECTL_BIN} -n kube-system get helmchart cert-manager -o 'jsonpath={..status.jobName}')
while [[ $(${KUBECTL_BIN} -n kube-system get jobs ${helm_job} -o 'jsonpath={..status.conditions[?(@.type=="Complete")].status}') != "True" ]] ; 
do echo "waiting for ${helm_job} job" && sleep 10; 
done
## waiting for cert-manager
while [[ $(for i in $(${KUBECTL_BIN} -n ${CERTMANAGER_NS} get pods -l app.kubernetes.io/instance=cert-manager -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}'); do if [ $i != "True" ]; then break; fi;  done && echo $i) != "True" ]] ; 
do echo "waiting for cert-manager pods" && sleep 10; 
done

## rancher
cat << EOF > ${TEMP_DIR}"/rancher.yaml"
---
kind: Namespace
apiVersion: v1
metadata:
  name: ${RANCHER_NS}
---
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: rancher
  namespace: kube-system
spec:
  chart: rancher
  repo: https://releases.rancher.com/server-charts/latest
  targetNamespace: ${RANCHER_NS}
  version: ${RANCHER_VERSION}
  set:
    hostname: ${RANCHER_HOSTNAME}
    replicas: 1
    rancherImageTag: ${RANCHER_VERSION}
    bootstrapPassword: "admin"
EOF

${DOCKER_BIN} cp ${TEMP_DIR}"/rancher.yaml" ${K3S_SERVER}:/var/lib/rancher/k3s/server/manifests/
## waiting for HelmChart rancher
while [[ $(${KUBECTL_BIN} -n kube-system get helmchart rancher -o 'jsonpath={..spec.set.rancherImageTag}') != ${RANCHER_VERSION} ]] ; 
do echo "Waiting for HelmChart rancher" && sleep 2;
done
## waiting for helm-install-rancher
helm_job=$(${KUBECTL_BIN} -n kube-system get helmchart rancher -o 'jsonpath={..status.jobName}')
while [[ $(${KUBECTL_BIN} -n kube-system get jobs ${helm_job} -o 'jsonpath={..status.conditions[?(@.type=="Complete")].status}') != "True" ]] ; 
do echo "waiting for ${helm_job} job" && sleep 10; 
done
## waiting for rancher
while [[ $(for i in $(${KUBECTL_BIN} -n ${RANCHER_NS} get pods -l app=rancher -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}'); do if [ $i != "True" ]; then break; fi;  done && echo $i) != "True" ]] ; 
do echo "waiting for rancher pods" && sleep 10; 
done
## Waiting for rancher server start
while [ "$(${CURL_BIN} -sk https://${RANCHER_HOSTNAME}/ping || echo starting)" != "pong" ]; 
do echo "waiting for rancher service" && sleep 5;
done
sleep 5

export RANCHER_INSECURE=true
export RANCHER_URL="https://${RANCHER_HOSTNAME}"
export RANCHER_VERSION=${RANCHER_VERSION}
if [ ${EXPOSE_HOST_PORTS} == "true" ]; then 
  export RANCHER_EXPOSED_URL="https://rancher.127.0.0.1.${RANCHER_DNS_DOMAIN}:${K3S_INGRESS_PORT_TLS}"
fi
