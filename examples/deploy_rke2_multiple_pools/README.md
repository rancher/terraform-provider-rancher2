# Deploy RKE2 with Multiple Machine Pools

This example module configures Rancher to deploy a downstream RKE2 cluster with two different roles: worker and control plane.

## Deploy Rancher

This module starts by using the rancher-aws module to deploy Rancher on AWS.
The Rancher cluster is a single node.

## Downstream

This module has a local module that provides a logical separation for deploying a downstream cluster using the rancher2_cluster_v2 resource.

## Machine Config Patch

There is a local exec that runs kubectl to patch the Amazonec2Config objects.
The AWS access key id and secret access key attributes are unable to be set directly with the provider in favor of the Amazonec2Credential object.
The Amazonec2Credential object doesn't support a session tokens, making it impossible to use temporary credentials.
Our CI must use temporary AWS credentials supplied via OIDC, so this may be something that you eliminate from the example in your implementation.
We bypass the Amazonec2Credential object by manually patching the Amazonec2Config objects with the aws_access_key_id and aws_secret_access_key.
The AWS session token isn't obfuscated and is available as an argument when creating a rancher2_machine_config_v2 resource.

## Dependencies

The `flake.nix` file in the root of the module explains all of the dependencies for the development of the module, it also includes the dependencies to run it.
You can see the list on lines 50-80, but a more specific list is below (with explanations).

- bash -> born again shell with linux core utils facilitates CLI actions
- tfswitch -> handy for installing Terraform at specific verisons
- git -> required by Terraform
- curl -> required by Terraform as well as dependent modules (when downloading RKE2 for install)
- openssh -> required by Terraform and used in dependent modules to connect to servers for initial configuration
- openssl -> required by Terraform and used in dependent modules to verify TLS certificates
- ssh-agent -> used for connecting to remote server for initial configuration, you need to have the key you send into the module loaded in your agent
- gh -> the github cli tool, used to find releases when downloading RKE2 for install
- jq -> json parsing tool, used in dependent modules to parse submodule outputs
- kubectl -> used in local exec to patch kubernetes objects
- awscli2 -> the aws cli tool, used in some dependent modules in some use cases (dualstack)
- yq -> yaml parsing tool, used in dependent modules to parse kubectl outputs
- go -> necessary to run tests

## Environment Variables

I like to set my AWS credentials in environment variables:

- AWS_ROLE -> role to assume when using OIDC
- AWS_REGION -> AWS region to deploy to, make sure there are multiple availability zones when needing HA
- AWS_DEFAULT_REGION -> same as region
- AWS_ACCESS_KEY_ID -> access key, this will make it into the state, please secure it properly
- AWS_SECRET_ACCESS_KEY -> secret key, this will make it into the state, please secure it properly
- AWS_SESSION_TOKEN -> used with temporary AWS credentials, this will make it into the state, please secure it properly
- TF_VAR_aws_access_key_id -> access key, this will make it into the state, please secure it properly
- TF_VAR_aws_secret_access_key -> secret key, this will make it into the state, please secure it properly
- TF_VAR_aws_session_token -> used with temporary AWS credentials, this will make it into the state, please secure it properly
- TF_VAR_aws_region -> AWS region to deploy to, make sure there are multiple availability zones when needing HA

