locals {
  identifier   = var.identifier
  rancher_url  = var.rancher_url
  rancher_host = split(".", split("://", local.rancher_url)[1])[0]
  owner        = var.owner
  callback_url = "${local.rancher_url}/verify-auth"
}

resource "aws_cognito_user_pool" "main" {
  name                     = local.identifier
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]
}

resource "aws_cognito_user_pool_domain" "main" {
  depends_on = [
    aws_cognito_user_pool.main,
  ]
  domain       = lower(local.rancher_host)
  user_pool_id = aws_cognito_user_pool.main.id
}

resource "aws_cognito_user_pool_client" "client" {
  depends_on = [
    aws_cognito_user_pool.main,
  ]
  name            = local.identifier
  user_pool_id    = aws_cognito_user_pool.main.id
  generate_secret = true

  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes                 = ["email", "openid", "profile"]

  callback_urls = [local.callback_url]
  logout_urls   = ["${local.rancher_url}/logout"]

  supported_identity_providers = ["COGNITO"]
}

resource "random_password" "admin_user" {
  length           = 16
  min_special      = 1
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
  override_special = "!#$-_=+"
}

resource "aws_cognito_user" "admin" {
  depends_on = [
    aws_cognito_user_pool.main,
    random_password.admin_user,
  ]

  user_pool_id = aws_cognito_user_pool.main.id
  username     = local.owner

  attributes = {
    email          = local.owner # CHANGE THIS to your real email
    email_verified = "true"      # Required so Cognito doesn't ask for a code
  }

  temporary_password = random_password.admin_user.result
  message_action     = "SUPPRESS"
}
