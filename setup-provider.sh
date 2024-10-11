#!/bin/bash

# QA has a new process to test Terraform using RCs. We will not publish RCs
# onto the Terraform registry because Hashicorp could potentially block our
# testing by being slow to publish.

# Instead, we will test using a downloaded binary from the RC. This script
# sets up the correct binary using a defined <provider> <version> to test
# updates locally.

# ./setup-provider.sh <provider> <version>

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

# Check for unzip
if ! command -v "unzip" &> /dev/null; then
  echo "Missing required tool unzip. unzip can be installed via apt on Linux"
  exit 1
fi

# Download binary

OS_PLATFORM=$(uname -sp | tr '[:upper:] ' '[:lower:]_' | sed 's/x86_64/amd64/' | sed 's/i386/amd64/' | sed 's/arm/arm64/')

DIR=~/.terraform.d/plugins/terraform.local/local/${PROVIDER}/${VERSION_TAG}/${OS_PLATFORM}
mkdir -p $DIR

# unzip can't handle files from stdin, so we create a temporary file
ZIP_FILE=terraform-provider-${PROVIDER}_${VERSION_TAG}_${OS_PLATFORM}.zip
curl -sfL https://github.com/rancher/terraform-provider-${PROVIDER}/releases/download/${VERSION}/${ZIP_FILE} -o ${ZIP_FILE}

unzip $ZIP_FILE -d ${DIR}
rm ${ZIP_FILE}

#Note the required 'v' in front of VERSION_TAG
mv ${DIR}/terraform-provider-${PROVIDER}_v${VERSION_TAG} ${DIR}/terraform-provider-${PROVIDER}

# Mod binary
chmod +x ${DIR}/terraform-provider-${PROVIDER}

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
