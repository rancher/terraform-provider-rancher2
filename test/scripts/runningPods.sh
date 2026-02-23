#!/bin/bash
set -x

JSONPATH="'{range .items[*]}
  {.metadata.name}{\"\\t\"} \
  {.metadata.namespace}{\"\\t\"} \
  {.status.phase}{\"\\n\"} \
{end}'"

notReady() {
  PODS=$(kubectl get pods -A -o jsonpath="$JSONPATH")
  # shellcheck disable=SC2060,SC2140
  NOT_READY=$(echo "$PODS" | grep -v "Running" | grep -v "Succeeded"  | tr -d ["\t","\n"," ","'"] || true)
  if [ -n "$NOT_READY" ]; then
    # Some pods aren't running
    return 0
  else
    # All pods are running
    return 1
  fi
}

readyWait() {
  TIMEOUT=10 # 10 minutes
  TIMEOUT_MINUTES=$((TIMEOUT * 60))
  INTERVAL=30 # 30 seconds
  MAX=$((TIMEOUT_MINUTES / INTERVAL))
  ATTEMPTS=0

  while notReady; do
    if [ "$ATTEMPTS" -lt "$MAX" ]; then
      ATTEMPTS=$((ATTEMPTS + 1))
      sleep "$INTERVAL";
    else
      return 1
    fi
  done
  return 0
}

SUCCESSES=0
SUCCESSES_NEEDED=3 # require three successes to make sure everything is settled

while readyWait && [ "$SUCCESSES" -lt "$SUCCESSES_NEEDED" ]; do
  SUCCESSES=$((SUCCESSES + 1))
  echo "succeeeded $SUCCESSES times..."
  sleep 30
done

if [ "$SUCCESSES" -eq "$SUCCESSES_NEEDED" ]; then
  echo "$SUCCESSES_NEEDED successes reached, passed.."
  EXITCODE=0
else
  echo "$SUCCESSES_NEEDED successes not reached, failed.."
  EXITCODE=1
fi

echo "nodes..."
kubectl get nodes || true

echo "all..."
kubectl get all -A || true

echo "pods..."
kubectl get pods -A || true

exit $EXITCODE
