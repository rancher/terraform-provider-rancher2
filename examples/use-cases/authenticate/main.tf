
provider "rancher2" {}

provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
}

terraform {
  backend "s3" {
    # This needs to be set in the backend configs on the command line or somewhere that your identifier can be set.
    # terraform init -reconfigure -backend-config="bucket=<identifier>"
    # https://developer.hashicorp.com/terraform/language/backend/s3
    # https://developer.hashicorp.com/terraform/language/backend#partial-configuration
    key = "tfstate"
  }
}

locals {
  identifier  = var.identifier
  owner       = var.owner
  rancher_url = var.rancher_url
  # client_id     = var.client_id     # "4cus7h997mtfb3fqoo5b3v4amj"
  # client_secret = var.client_secret # "plaintextpassword"
  # issuer        = var.issuer        # "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_YtjZFSwyl"
}

module "cognito" {
  source      = "./modules/cognito"
  identifier  = local.identifier
  rancher_url = local.rancher_url
  owner       = local.owner
}

# It is much better to configure the client using environment variables,
#  this keeps these secure values out of your state file.
resource "rancher2_client_basic" "admin" {
  id      = "admin"
  api_url = local.rancher_url # RANCHER_API_URL
  # access_key        = "" RANCHER_ACCESS_KEY
  # secret_key        = "" RANCHER_SECRET_KEY
  # token_key         = "" RANCHER_TOKEN_KEY
  # ca_certs          = "" RANCHER_CA_CERTS
  # ignore_system_ca  = "" RANCHER_IGNORE_SYSTEM_CA
  # insecure          = "" RANCHER_INSECURE
  # max_redirects     = "" RANCHER_MAX_REDIRECTS
  # timeout           = "" RANCHER_TIMEOUT // this is http client connection timeout not token TTL
}

# at this point the provider is configured

# The name of the resource reflects the API endpoint
# https://<rancher url>/<API path>/authconfigs/cognito
# https://rancher.example.com/apis/management.cattle.io/v3/authconfigs/cognito
# the resource pkg will determine API path, based on which API it exposes (norman/steve/kubernetesExt)
# resource "rancher_authconfig_cognito" {
# The attributes reflect the payload
# the schema for the payload isn't found in the CRD, which has a dynamic mapping
# since the CRD is dynamic you need to create an object using the UI first to get the actual fields
# since the payload is essentially generating an object from the CRD it should reflect the object's fields
# rancher_client_id = rancher2_client_basic.admin.id

# rancher_url       = join("/", [local.rancher_url, "auth-verify"]) #<<-- required
# enabled     = true #<<-- optional with default (true)

# client_id      = local.client_id #<<-- required
# client_secret  = local.client_secret #<<-- required
# issuer         = local.issuer #<<-- required
# RANCHER_AUTHCONFIG_COGNITO_CLIENT_SECRET keeps plaintext password out of state
# however, if state doesn't hold password then it must exist in the environment to facilitate create and update
# on read we resolve the secret in kubernetes using the kubernetes API, it can then be compared to the environment or config
# if both env variable and config are nil then we ignore in read, but error in create/update
#   the idea behind this is that the environment variable only needs to be present when a create/update is necessary
#   diff will only trigger an update when the config doesn't match state, so if the password is nil in both then no update will trigger
# clientSecretReference = "cattle-global-data:cognitoconfig-clientsecret" #<<-- calculated, read only attribute, this is what the object will output

# name        = "cognito" #<<-- calculated, derived from the resource name
# type        = "cognitoConfig" #<<-- calculated, derived from the resource name
# kind        = "AuthConfig" #<<-- calculated, derived from the resource name
# apiVersion  = "management.cattle.io/v3" #<<-- calculated, determined by the resource

# accessMode          = "unrestricted" # <<-- optional with default?
# scope               = "openid email" # <<-- optional with default, must minimally have "openid email"
# logoutAllSupported  = true  # <<-- mutually exclusive with other logout options, optional with default, calculated
# logoutAllForced     = false # <<-- mutually exclusive with other logout options, optional with default, calculated
# logoutAllEnabled    = false # <<-- mutually exclusive with other logout options, optional with default, calculated
# groupSearchEnabled  = false # <<-- optional with default
# endSessionEndpoint  = null # <<-- optional

# groupsClaim         = "cognito:groups" #<<-- calculated by Rancher, read only attribute
# allowedPrincipalIds = [
#   "cognito_user://d82153b0-8061-70b8-6abd-07e56ff4c646" #<<-- calculated by Rancher
# ]
# metadata = {
#   annotations = { # <<-- optional
#     "key" = "value"
#   }
#   all_annotations = { # <<- calculated
#     "auth.cattle.io/unused-secrets-cleaned"      = "true"
#     "management.cattle.io/auth-provider-cleanup" = "unlocked"
#   }
#   labels = { # <<-- optional
#     "key" = "value"
#   }
#   all_labels = { # <<-- calculated
#     "cattle.io/creator" = "norman"
#   }
#   # maybe we should do like the AWS provider does tags, we have "labels" which is configurable, and "all_labels" which is an attribute and includes everything
#   # read won't know the difference, but if a label is missing that is in the config then a diff will show and the object can be updated
#   creationTimestamp = "2025-12-11T19:33:21Z" # <<-- calculated
#   generation        = 4 # <<-- calculated
#   resourceVersion   = "10722" # <<-- calculated
#   name              = "cognito" # <<-- calculated
#   uid               = "669910a4-1d15-40a8-9540-8f0584763605" # <<-- calculated
# }
# status = { # <<-- calculated
#   conditions = [
#     {
#       status  = "True"
#       type    = "SecretsMigrated"
#     }
#   ]
# }
# }

# resource "rancher2_client_aws_cognito" "admin" {
#   id               = "admin"
#   api_url          = local.rancher_url # RANCHER_API_URL
#   access_key       = "" RANCHER_ACCESS_KEY
#   secret_key       = "" RANCHER_SECRET_KEY
#   token_key        = "" RANCHER_TOKEN_KEY
#   ca_certs         = "" RANCHER_CA_CERTS
#   ignore_system_ca = "" RANCHER_IGNORE_SYSTEM_CA
#   insecure         = "" RANCHER_INSECURE
#   max_redirects    = "" RANCHER_MAX_REDIRECTS
#   timeout          = "" RANCHER_TIMEOUT // this is http client connection timeout not token TTL
# }
