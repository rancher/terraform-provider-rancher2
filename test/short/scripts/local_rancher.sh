#!/bin/bash

set -e


: "${RANCHER_VERSION:=v2.13.2}"
: "${CLUSTER_NAME:=local-rancher}"
: "${K3S_IMAGE:=rancher/k3s:v1.33.8-k3s1}"

# --- Helper functions ---
info() {
  echo '[INFO] ' "$@"
}

fatal() {
  echo '[ERROR] ' "$@" >&2
  exit 1
}
# Function to clean up resources
cleanup() {
  info "Cleaning up resources..."
  rm -f traefik-values.yaml || true
  info "Deleting k3d cluster..."
  k3d cluster delete "$CLUSTER_NAME" || true
  info "Pruning Docker..."
  DOCKERPS=$(docker ps | grep -v -e '^CONTAINER' | awk '{print $1}')
  if [ -n "$DOCKERPS" ]; then
    info "Stopping Docker containers..."
    for c in $DOCKERPS; do
      docker container stop "$c" || true
    done
  fi
  docker system prune -af || true
  docker volume prune -af || true
  if [ "$CI" != "true" ]; then
    info "Stopping colima..."
    colima stop || true
    colima delete -f || true
  fi
}

if [ "$1" == "cleanup" ]; then
  cleanup
  exit 0
fi

# --- Prerequisite checks ---
for cmd in k3d kubectl helm ssh openssl; do
  if ! command -v "$cmd" &> /dev/null; then
    fatal "$cmd is not installed. Please install it before running this script."
  fi
done
if [ "$CI" != "true" ]; then
  if ! command -v "colima" &> /dev/null; then
    fatal "colima is not installed. Please install it before running this script."
  fi
fi

info "Cleaning up any previous environments..."
cleanup

# --- Main script ---
info "Starting local Rancher setup..."
info "Rancher version: $RANCHER_VERSION"
info "Cluster name: $CLUSTER_NAME"

if [ "$CI" != "true" ]; then
  info "Configuring and starting Colima with 4 CPUs and 4GB RAM..."
  colima start --cpu 4 --memory 4
  colima status
fi

info "Creating new k3d cluster..."
# Map host's 443 to the load balancer's 443, making Rancher accessible on the standard HTTPS port.
k3d cluster create "$CLUSTER_NAME" --api-port 6443 -p "80:80@loadbalancer" -p "443:443@loadbalancer" --image "$K3S_IMAGE" --k3s-arg "--disable=traefik@server:0" --wait

info "Installing custom Traefik..."
helm repo add traefik https://traefik.github.io/charts
helm repo update

info "Creating Traefik values file..."
cat <<EOF > traefik-values.yaml
additionalArguments:
  - "--entrypoints.web.forwardedHeaders.insecure=true"
  - "--entrypoints.websecure.forwardedHeaders.insecure=true"
EOF

helm install traefik traefik/traefik --namespace kube-system --values traefik-values.yaml --wait

kubectl cluster-info

info "Getting k3d load balancer IP address..."
RETRY_COUNT=0
MAX_RETRIES=10
RETRY_INTERVAL=3
LOAD_BALANCER_IP=""
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  LOAD_BALANCER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "k3d-${CLUSTER_NAME}-serverlb")
  if [ -n "$LOAD_BALANCER_IP" ]; then
    break
  fi
  info "Waiting for k3d load balancer to get an IP... ($((RETRY_COUNT+1))/$MAX_RETRIES)"
  sleep $RETRY_INTERVAL
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ -z "$LOAD_BALANCER_IP" ]; then
  fatal "Could not determine k3d load balancer IP address after trying for $((MAX_RETRIES * RETRY_INTERVAL)) seconds."
fi
info "k3d load balancer IP address is: $LOAD_BALANCER_IP"

HOSTNAME="127.0.0.1.nip.io"
info "Hostname: $HOSTNAME"


info "Adding Rancher helm repository..."
helm repo add rancher-latest https://releases.rancher.com/server-charts/latest
helm repo update

info "Pre-flight checks and configuration..."

info "Generating self-signed certificate..."
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes -subj "/CN=${HOSTNAME}" -addext "subjectAltName = DNS:${HOSTNAME}"

info "Creating TLS secret..."
kubectl create namespace cattle-system || true
kubectl -n cattle-system create secret tls tls-rancher-ingress \
   --cert=cert.pem \
   --key=key.pem

info "Deploying Rancher..."
helm install rancher rancher-latest/rancher \
  --namespace cattle-system \
  --create-namespace \
  --set hostname="$HOSTNAME" \
  --set bootstrapPassword=admin \
  --set replicas=1 \
  --set rancherImageTag="$RANCHER_VERSION" \
  --set telemetry.enabled=false \
  --set ingress.tls.source=secret \
  --set ingress.tls.secretName=tls-rancher-ingress \
  --set ingress.ingressClassName=traefik \
  --wait \
  --debug

info "Rancher is deployed!"
info "Access it at: https://$HOSTNAME"
info "Bootstrap password is: admin"

rm cert.pem key.pem
