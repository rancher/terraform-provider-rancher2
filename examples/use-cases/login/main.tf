# Rancher2 Login Resource
# These resources are special because they modify the running provider's client.
# The provider instantiates a client at config time (pre-read),
#   these resources add their tokens to the client so that it can use them for further communications.
# Further resources should depend on these to use the proper authentication.
# Optionally, further resources can set the "user_token" argument to set their own user token.
# Each time 
provider "rancher2" {
  url = "https://rancher.example.com" # url is required, but can be set with environment variables
}

resource "rancher2_login" "initial_admin" {
  # username = "" provided by the RANCHER_USERNAME environment variable, this won't be saved in state
  # password = "" provided by the RANCHER_PASSWORD environment variable, this won't be saved in state
  refresh_token_ttl               = "1h"
  refresh_token_expiration_ignore = true # don't automatically refresh the refresh token if it expires
  user_token_ttl                  = "1h"
  user_token_expiration_ignore    = true # don't automatically refresh the user token if it expires
  # By allowing the refresh and user tokens to expire we essentially make this a one time operation
  # This says not to login as the admin user again unless the resource is tainted.
}

# all further resources should depend on the login resource
# this resource will be accessed as the initial admin user
resource "rancher2_user" "kate" {
  depends_on = [
    rancher2_login.initial_admin,
  ]
  user_token = rancher2_login.initial_admin.user_token # this is optional
  username   = ""
  password   = ""
  user_data  = "" # ...
}

resource "rancher2_login" "kate" {
  depends_on = [
    rancher2_user.kate,
  ]
  username              = rancher2_user.kate.username
  password              = rancher2_password.kate.password
  refresh_token_ttl     = "90d"
  refresh_token_timeout = "10d" # the amount of time before expiration to refresh
  user_token_ttl        = "16h"
}

# all further resources should depend on the login user token resource
# this resource will be accessed as the kate user
resource "rancher2_dev_resource" "test" {
  depends_on = [
    rancher2_login.kate,
  ]
  user_token = rancher2_login.kate.user_token # this is technically optional, but recommended when using multiple login resources
  id         = "test"
}
