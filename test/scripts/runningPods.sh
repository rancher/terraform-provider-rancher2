#!/bin/bash
set -x

JSONPATH='{range .items[*]}{.metadata.name}{"\t"}{.metadata.namespace}{"\t"}{.status.phase}{"\n"}{end}'

ATTEMPTS="${ATTEMPTS:-0}"
MAX="${1:-${MAX:-30}}"
INTERVAL="${2:-${INTERVAL:-10}}"
EXITCODE=1

notReady() {
  PODS=$(kubectl get pods -A -o jsonpath="$JSONPATH")
  if [ -z "$PODS" ]; then
    # No pods returned yet or API unreachable; continue waiting
    return 0
  fi
  NOT_READY=$(echo "$PODS" | awk '$NF != "Running" && $NF != "Succeeded"' | tr -d " \t\n'" || true)
  if [ -n "$NOT_READY" ]; then
    # Some pods aren't running
    return 0
  else
    # All pods are running
    return 1
  fi
}

readyWait() {
  ATTEMPTS="${ATTEMPTS:-0}"
  MAX="${1:-${MAX:-30}}"
  INTERVAL="${2:-${INTERVAL:-10}}"

  while notReady; do
    if [ "$ATTEMPTS" -lt "$MAX" ]; then
      ATTEMPTS=$((ATTEMPTS + 1))
      sleep "$INTERVAL"
    else
      return 1
    fi
  done
  return 0
}

SUCCESSES=0
SUCCESSES_NEEDED=3 # require three successes to make sure everything is settled

while [ "$SUCCESSES" -lt "$SUCCESSES_NEEDED" ] && readyWait "$@"; do
  ATTEMPTS=0
  SUCCESSES=$((SUCCESSES + 1))
  echo "succeeded $SUCCESSES times..."
  if [ "$SUCCESSES" -lt "$SUCCESSES_NEEDED" ]; then
    sleep 30
  fi
done

echo "nodes..."
kubectl get nodes || true

echo "all..."
kubectl get all -A || true

echo "pods..."
kubectl get pods -A || true

if [ "$SUCCESSES" -lt "$SUCCESSES_NEEDED" ]; then
  echo "Error: Timeout reached waiting for pods to settle" >&2
  exit 1
fi

echo "$SUCCESSES_NEEDED successes reached, passed.."
EXITCODE=0

exit "${EXITCODE:-0}"
