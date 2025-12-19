output "client_id" {
  value = module.cognito.client_id
}

output "client_secret" {
  value     = module.cognito.client_secret
  sensitive = true
}

output "issuer_url" {
  value = module.cognito.issuer_url
}
