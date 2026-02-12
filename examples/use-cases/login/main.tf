# Rancher2 Login Resource
provider "rancher2" {
  url = local.rancher_url
}
locals {
  rancher_url = var.rancher_url
}

resource "rancher2_login" "initial_admin" {
  # username_environment_variable = "RANCHER_USERNAME" # optional, tells resource where to look for username, defaults to RANCHER_USERNAME
  # password_environment_variable = "RANCHER_PASSWORD" # optional, tells resource where to look for password, defaults to RANCHER_PASSWORD
  # username = "" provided by the RANCHER_USERNAME environment variable, this won't be saved in state
  # password = "" provided by the RANCHER_PASSWORD environment variable, this won't be saved in state
  token_ttl    = "90d" # this token will expire in 90 days from the time it is created or refreshed
  refresh_at   = "10d" # this token will be refreshed 10 days before it expires, resetting the 90 day ttl
  ignore_token = false # if set to true this won't save the resulting token to state and will always recreate on plan/apply
}

# resource "random_password" "kates_password" {
#   length           = 16
#   special          = true
#   override_special = "!#$&*()-_=+[]{}<>?"
# }

# # all further resources should depend on the login resource
# # this resource will be accessed as the initial admin user
# resource "rancher2_user" "kate" {
#   depends_on = [
#     rancher2_login.initial_admin,
#   ]
#   user_token = rancher2_login.initial_admin.user_token # this is optional, but recommended when using multiple login resources
#   username   = "kate"
#   password   = random_password.kates_password.result
#   user_data  = "" # ...
# }

# resource "rancher2_login" "kate" {
#   depends_on = [
#     rancher2_user.kate,
#   ]
#   username       = rancher2_user.kate.username
#   password       = rancher2_user.kate.password
#   user_token_ttl = "1d"
#   ignore_token   = true # this user will always login
# }

# # all further resources should depend on the login user token resource
# # this resource will be accessed as the kate user
# resource "rancher2_dev_resource" "test" {
#   depends_on = [
#     rancher2_login.kate, # depending on kate's login will ensure that this uses her token, unless some other login takes precedence
#   ]
#   user_token = rancher2_login.kate.user_token # this is technically optional, but recommended when using multiple login resources
#   id         = "test"
# }
