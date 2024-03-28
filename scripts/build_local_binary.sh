#!/bin/bash

<<comment
When you run a `make bin` in the terraform-provider-rancher2 repo, it
downloads a binary into your local cache from the specified source. To
tell Terraform to use a local binary, you need to overwrite the one in
the cache.

Add the following to your Terraform config file. This tells Terraform
to use the local cache https://terraform.io/cli/config/config-file.

terraform {
required_providers {
    rancher2 = {
      source = "terraform.local/local/rancher2"
      version = "1.0.0"
    }
  }
}
comment

# Compile the local binary
make bin

opsys=windows
if [[ "$OSTYPE" == linux* ]]; then
  opsys=linux
elif [[ "$OSTYPE" == darwin* ]]; then
  opsys=darwin
fi

# Supported values of 'arch': amd64, arm64, ppc64le, s390x
case $(uname -m) in
x86_64)
    arch=amd64
    ;;
arm64)
    arch=arm64
    ;;
ppc64le)
    arch=ppc64le
    ;;
s390x)
    arch=s390x
    ;;
*)
    arch=amd64
    ;;
esac

LOCAL_BINARY_PATH=~/.terraform.d/plugins/terraform.local/local/rancher2/1.0.0/${opsys}_${arch}

if [ ! -d $LOCAL_BINARY_PATH ]
then
  echo "Directory ${LOCAL_BINARY_PATH} does not exist. Creating..."
  mkdir -p $LOCAL_BINARY_PATH
fi

# Overwrite the cached binary with your local one
cp terraform-provider-rancher2 ${LOCAL_BINARY_PATH}/terraform-provider-rancher2

# Check if user specified options for init or apply
while getopts 'ia' OPTION; do
  case "$OPTION" in
    i)
      rm -r .terraform.lock.hcl
      terraform init
      ;;
    a)
      rm -r .terraform.lock.hcl
      terraform init
      terraform apply
      ;;
    ?)
      echo "script usage: $(basename \$0) [-i] [-a]" >&2
      exit 1
  esac
done
