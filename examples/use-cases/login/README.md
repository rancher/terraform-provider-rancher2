# Login

This example shows how to login to Rancher within an apply operation.
In testing we use this to allow a single apply operation to deploy VMs on AWS,
install Rancher using Helm, and configure Rancher with the provider.
In at least one of the tests we configure Rancher to provision downstream clusters.

## Multiple Ways To Login

There are multiple ways to login to Rancher with the provider.

1. Configure the provider with a token.
2. Use the rancher2_login resource to inject a token into the client.
3. Each resource has an optional "user_token" argument that can inject a token into the client.

### Provider Config

You can configure the provider with a token either by assigning it to the arguments or by setting the RANCHER_TOKEN environment variable, provider configurations are never saved in state.

When you assign a token in the provider config every request will authenticate with this token.
This forces the user to maintain the token, and is useful when the user has other methods of managing token generation.
This method is generally not recommended unless you have a well established process for obtaining and rotating tokens properly.

### Login Resource

The login resource is our recommended approach to maintaining a token for Rancher authentication, 
it requires a username and password that can be set in the configuration or passed in the environment variables RANCHER_USERNAME and RANCHER_PASSWORD. When set with environment variables the contents aren't saved to state, but each time the resource is recreated it will need those variables present.

The rancher2_login resource overwrites the provider's client's auth_token attribute, all further resources will then use that token.
The provider first logs into Rancher with the username and password, retrieving a "session token" which only lasts until the end of the current execution. Then it uses that token to create a "refresh token" which only has access to create other tokens, but has a long lifetime. Then it uses the refresh token to get a "user token" which has a shorter lifespan but has the full authorization of the user. The user token is what will be added to the client for future requests.
After the initial Terraform apply, the login resource will attempt to use the user token it has in state, if that is expired it will attempt to get a new user token using the refresh token in state, only if that token is invalid will it login with the username and password again to restart the process.

### Resource's Token Argument

Each resource has an optional attribute "user_token" which can be supplied at apply time.
This allows each resource to override a previous authentication mechanism.
For instance, if you wanted to supply a different token for a resource you could do that in the config of that resource.
You could use multiple login resources and supply each resource the token from the correlating login to build a robust authorization management system for your Rancher configuration.

