#!/bin/bash

# QA has a new process to test Terraform using rcs. We will not publish rcs
# onto the Terraform registry because Hashicorp could potentially block our
# testing by being slow to publish.

# Instead, we will test using a downloaded binary from the rc. This script
# sets up the correct binary using a defined <provider> <version> to test
# updates locally.

# Example
# ./setup-provider.sh rancher2 v3.0.0-rc1

set -e 

# Validate user input

if [ $# -ne 2 ]; then
  echo "Usage: $0 <provider> <version>"
  exit 1
fi

# Set global vars

PROVIDER=$1
VERSION=$2
VERSION_TAG=$(echo $2 | cut -c 2-)

# Install gzip

if ! command -v "gzip" &> /dev/null; then
  echo "Missing gzip. Installing..."
  brew install gzip
fi

# Download binary

OS_PLATFORM=$(uname -sp | tr '[:upper:] ' '[:lower:]_' | sed 's/x86_64/amd64/')

if [[ $OS_PLATFORM == "darwin"* ]]
then
  OS_PLATFORM="darwin_amd64"
elif [[ $OS_PLATFORM == "linux"* ]]
then
  OS_PLATFORM="linux_amd64"
fi

DIR=~/.terraform.d/plugins/terraform.local/local/${PROVIDER}/${VERSION_TAG}/${OS_PLATFORM}
mkdir -p $DIR
curl -sfL https://github.com/rancher/terraform-provider-${PROVIDER}/releases/download/${VERSION}/terraform-provider-${PROVIDER}_${VERSION_TAG}_${OS_PLATFORM}.zip | gunzip -c - > ${DIR}/terraform-provider-${PROVIDER}

echo -e "Terraform provider ${PROVIDER} ${VERSION} is ready to test!
Please update the required_providers block in your Terraform config file

terraform {
  required_providers {
    rancher2 = {
      source = "terraform.local/local/${PROVIDER}"
      version = "${VERSION_TAG}"
    }
  }
}"
