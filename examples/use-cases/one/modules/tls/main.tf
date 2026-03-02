
locals {
  domain = var.domain
}

resource "tls_private_key" "ca_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "ca_cert" {
  private_key_pem = tls_private_key.ca_key.private_key_pem

  subject {
    common_name  = "Example CA"
    organization = "Example"
  }

  validity_period_hours = 8760
  is_ca_certificate     = true

  allowed_uses = [
    "cert_signing",
    "crl_signing",
  ]
}

// TLS Certificate
resource "tls_private_key" "tls_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "tls_csr" {
  private_key_pem = tls_private_key.tls_key.private_key_pem

  subject {
    common_name  = local.domain
    organization = "Example"
  }

  dns_names = [local.domain]
}

resource "tls_locally_signed_cert" "tls_cert" {
  cert_request_pem   = tls_cert_request.tls_csr.cert_request_pem
  ca_private_key_pem = tls_private_key.ca_key.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca_cert.cert_pem

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}
