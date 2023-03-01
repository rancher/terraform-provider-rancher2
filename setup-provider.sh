#!/bin/bash

# QA has a new process to test Terraform using rcs. We will not publish rcs
# onto the Terraform registry because Hashicorp could potentially block our
# testing by being slow to publish.

# Instead, we will test using a downloaded binary from the rc. This script
# sets up the correct binary using a defined <provider> <version> to test
# updates locally.

# Example
# ./setup-provider.sh rancher2 v3.0.0-rc1

# Validate user input

if [ -z "$1" ] || [ -z "$2" ]; then
    echo -e "Please specify a Terraform provider and version tag."
    echo -e "./setup-provider.sh <provider> <version>"
    exit 1
fi

# Set global vars

PROVIDER=$1
VERSION=$2
VERSION_TAG=$(echo $2 | cut -c 2-)
BASE=~/.terraform.d/plugins/terraform.example.com/local/${PROVIDER}/${VERSION_TAG}

# Install gzip

brew install gzip

# Pull and copy the binary

case $OSTYPE in

  "linux-gnu"*)
    curl -sfL https://github.com/rancher/terraform-provider-${PROVIDER}/releases/download/${VERSION}/terraform-provider-${PROVIDER}_${VERSION_TAG}_linux_amd64.zip | gunzip -c - > terraform-provider-${PROVIDER}_${VERSION}
    if [ ! $? -eq 0 ]; then
      echo "Error: incomplete download. Did you specify the correct provider and version?"
      rm terraform-provider-${PROVIDER}_${VERSION}
      exit 1
    fi
    DIR=${BASE}/linux_amd64
    mkdir -p $DIR && cp terraform-provider-${PROVIDER}_${VERSION} $DIR/terraform-provider-${PROVIDER}
    ;;

  "darwin"*)
    curl -sfL https://github.com/rancher/terraform-provider-${PROVIDER}/releases/download/${VERSION}/terraform-provider-${PROVIDER}_${VERSION_TAG}_darwin_amd64.zip | gunzip -c - > terraform-provider-${PROVIDER}_${VERSION}
    if [ ! $? -eq 0 ]; then
      echo "Error: incomplete download. Did you specify the correct provider and version?"
      rm terraform-provider-${PROVIDER}_${VERSION}
      exit 1
    fi
    DIR=${BASE}/darwin_amd64
    mkdir -p $DIR && cp terraform-provider-${PROVIDER}_${VERSION} $DIR/terraform-provider-${PROVIDER}
    ;;
esac

# Clean up

rm terraform-provider-${PROVIDER}_${VERSION}

echo -e "Terraform provider ${PROVIDER} ${VERSION} is ready to test!
Please update the required_providers block in your Terraform config file. Example:

terraform {
  required_providers {
    rancher2 = {
      source = "terraform.example.com/local/rancher2"
      version = "3.0.0-rc1"
    }
  }
}"
