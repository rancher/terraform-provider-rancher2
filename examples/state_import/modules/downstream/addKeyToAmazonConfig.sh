#!/bin/sh

ACCESS_KEY="$1"
SECRET_KEY="$2"

if [ -z "$ACCESS_KEY" ]; then echo "need access key to proceed..."; exit 1; fi
if [ -z "$SECRET_KEY" ]; then echo "need secret key to proceed..."; exit 1; fi

NAMESPACE="fleet-default"
NAMES=$(kubectl get Amazonec2Config -n "$NAMESPACE" -o jsonpath='{.items[*].metadata.name}')

if [ -z "$NAMES" ]; then echo "names not found..."; exit 1; fi

cat <<EOT > patch.yaml
accessKey: '$ACCESS_KEY'
secretKey: '$SECRET_KEY'
EOT

for name in $NAMES; do

  kubectl patch Amazonec2Config "$name" -n "$NAMESPACE" --type merge --patch-file patch.yaml
  RESULT=$?
  if [ $RESULT -ne 0 ]; then
    echo "error: $RESULT"
    exit $RESULT
  fi

  KEY="$(kubectl get Amazonec2Config "$name" -n "$NAMESPACE" -o json | jq -r '.accessKey')"

  if [ -z "$KEY" ] || [ "null" = "$KEY" ]; then
    echo "error: key not found on object"
    exit 1
  else
    if [ "$KEY" != "$ACCESS_KEY" ]; then
      echo "error: key not replaced properly"
      exit 1
    fi
    echo "key replaced properly"
  fi

  SECRET="$(kubectl get Amazonec2Config "$name" -n "$NAMESPACE" -o json | jq -r '.secretKey')"

  if [ -z "$SECRET" ] || [ "null" = "$SECRET" ]; then
    echo "error: secret not found on object"
    exit 1
  else
    if [ "$SECRET" != "$SECRET_KEY" ]; then
      echo "error: secret not replaced properly"
      exit 1
    fi
    echo "secret replaced properly"
  fi

  echo "Amazonec2Config $name in namespace $NAMESPACE updated."

done

rm -f patch.yaml
