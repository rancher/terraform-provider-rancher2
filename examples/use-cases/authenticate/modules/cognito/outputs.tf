output "client_id" {
  value = aws_cognito_user_pool_client.client.id
}

output "client_secret" {
  value     = aws_cognito_user_pool_client.client.client_secret
  sensitive = true
}

output "issuer_url" {
  value = "https://${aws_cognito_user_pool.main.endpoint}"
}
