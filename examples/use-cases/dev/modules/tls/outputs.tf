output "tls_public_certificate" {
  value = tls_locally_signed_cert.tls_cert.cert_pem
}

output "tls_private_key" {
  value     = tls_private_key.tls_key.private_key_pem
  sensitive = true
}

output "certificate_chain" {
  value = tls_self_signed_cert.ca_cert.cert_pem
}
