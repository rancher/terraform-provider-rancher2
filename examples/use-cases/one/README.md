# Single Server All in One Rancher Deployment

This module deploys a single Rancher server with all kubernetes roles.

This shows the most basic use case for the provider, is functions as a good start for configuring your Rancher deployment.

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

## Expandable

This module doesn't limit you to only development use.
Many times a development deployment becomes the only working deployment and then somehow the production deployment, this isn't ideal, but does happen.
With that in mind, this example is fully expandable to a production ready cluster.

I know "production ready" means a lot of things to a lot of people, so to break that down:

- you can increase the node count and it will automatically deploy to separate availability zones
- it includes a network load balancer so that it will automatically balance loads as more nodes are added
- you can add nodes with specific roles, then remove the all-in-one nodes to have dedicated roles
- you can change the cert from self signed to externally valid by replacing the tls internal module
  - you then supply the cert to Rancher in the same way as the self-signed
  - you may need to add new nodes and destroy the old ones afterwords

After these steps are complete your cluster should match the "production" example.
The "production" example is what the RKE2 team considers the most scalable approach for cluster infrastructure.

Due to the expandable nature of this example, it will deploy expensive than is absolutely necessary for a single node.
This can lead to additional charges from AWS, please review the deployment carefully to see whether it works for you.
The network load balancer that is included is often a cost consideration.
