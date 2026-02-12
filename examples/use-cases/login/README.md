# Login

This example shows how to login to Rancher within an apply operation.
In testing we use this to allow a single apply operation to deploy VMs on AWS,
install Rancher using Helm, and configure Rancher with the provider.
In at least one of the tests we configure Rancher to provision downstream clusters.

## Multiple Ways To Login

There are multiple ways to login to Rancher with the provider.

1. Configure the provider with a token.
2. Use the rancher2_login resource to inject a token into the client.
3. Each resource has an optional "user_token" argument that can inject a token into the client when operating that resource.

### Provider Config

You can configure the provider with a token either by assigning it to the arguments 
 or by setting the RANCHER_TOKEN environment variable, provider configurations are never saved in state.

When you assign a token in the provider config every request will authenticate with this token.
This forces the user to maintain the token, and is useful when the user has other methods of managing token generation.
This method is generally not recommended unless you have a well established process for obtaining and rotating tokens properly.

### Login Resource

The login resource is our recommended approach to maintaining a token for Rancher authentication, 
it requires a username and password that can be set in the configuration 
 or passed in the environment variables RANCHER_USERNAME and RANCHER_PASSWORD.
When set with environment variables the contents aren't saved to state,
 but each time the resource is recreated it will need those variables present.

The rancher2_login resource overwrites the provider's client's auth_token attribute, all further resources will then use that token.
The provider first logs into Rancher with the username and password,
 retrieving a "session token" which only lasts until the end of the current execution.
Then it uses that token to get a "user token" with a longer lifespan, 90 days by default.
The user token is what will be added to the client for future requests.
After the initial Terraform apply, the login resource will attempt to use the user token it has in state,
 if that is expired it will get a new session token and refresh the user token.
On creation a date will be set for the user_token to be refreshed,
 after that date update will be triggered which will attempt to get a new token using the current token,
 it will then save the new token to state and update the refresh and expiration dates.

### Resource's Token Argument

Each resource has an optional attribute "user_token" which can be supplied at apply time.
This allows each resource to override a previous authentication mechanism.
For instance, if you wanted to supply a different token for a resource you could do that in the config of that resource.
You could use multiple login resources and supply each resource the token from the correlating login
 to build a robust authorization management system for your Rancher configuration.

## How this resource operates

This resource is special because it modifies the running provider's client.
The provider instantiates a client at config time (pre-read),
  this resource adds its token to the client so that it can use it for further communications.
Further resources should depend on this resoruce to use the proper authentication.
Optionally, further resources can set the "user_token" argument to set their own user token.

### CREATE

On create, the resource attempts to login to Rancher using the username and password supplied.
The successful response will return what is normally considered a "session token",
 however in API terms this is just a user token with a short lifespan (16 hours by default).
We use this token to get a "user token" with a longer lifespan (90 days by default).
We then save the new token to state with dates for when it was created, when it expires, and when we will refresh it.

### READ

On read, the resource attempts to get the token resource using the user_token in state to authenticate.
- If this request fails with a 401 then it is assumed the token is expired and read will clear the token from state, prompting an update.
- If this request passes, then the token is valid and the refresh date is checked.
- If the refresh date has passed then the refresh date is cleared, prompting an update.
- If the refresh date has not passed then the read completes and no changes are necessary.

### UPDATE

On update, the resource attempts to refresh the user token.
- If the current token is null or if the current token fails to authenticate
 then the resource falls back on the username and password to get a new token.
- If the refresh passes then the new token is added to state and the refresh, create, and expiration dates are updated.

### DELETE

On delete, the resource attempts to delete the user token.
- If this request fails with a 401 then the token isn't valid and delete succeeds.
- If the request succeeds then the delete succeeds.
- Anything else results in an error.
