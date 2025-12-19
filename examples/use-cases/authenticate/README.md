# Authenticate

This is the first thing most users will be concerned with after installing Rancher.

The Rancher Helm chart only includes the built-in local password auth provider, and only the admin user.
While this is a logical starting point it leaves much to be desired from a security aspect.

Most users will want to configure an auth provider, generate a user,
authenticate as that user and start using that token for further calls to the provider.
This example shows that process.

## Beginning user

This example assumes you have configured Rancher via the Helm chart,
so the only user present is the admin user setup during that process.
We therefore use the admin user to configure the new auth process
and we expect the username and password ("admin" + whatever password you set up in the Helm chart)
to be configured in the provider block.

## Authentication Takeover

After the user has been generated, we have a special resource "rancher_authenticate"
whose goal is to authenticate as the new user and reconfigure the provider to use the new token that provides.
This resource attempts to authenticate with the generated token during read and if it fails it will request recreate.
The recreate logic will reauthenticate as the user and retrieve a new token.
When this resource is present, the authentication information in the provider block will be overridden and no longer used.

## AWS Cognito

We are already authenticated to AWS for end to end tests,
so it makes the most sense for us to configure the AWS Cognito auth provider for this example.
We generate a new AWS Cognito client config because our callback URL is random between runs for our testing.

## Bootstrap

This example can be seen somewhat as the replacement for the "rancher_bootstrap" resource, which is no longer available.
Bootstrap did other things which this doesn't do, but those things aren't necessary anymore or should be handled separately.
The goal is the same, allow users to configure Rancher from beginning to end using the same Terraform config.
