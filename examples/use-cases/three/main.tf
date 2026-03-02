provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
}

provider "acme" {
  server_url = local.acme_server_url
}

provider "github" {}
provider "kubernetes" {} # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)
provider "helm" {}       # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)


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
  identifier           = var.identifier
  example              = "three"
  project_name         = "tf-${substr(md5(join("-", [local.example, local.identifier])), 0, 5)}"
  username             = local.project_name
  domain               = local.project_name
  zone                 = var.zone
  key_name             = var.key_name
  key                  = var.key
  acme_server_url      = "https://acme-staging-v02.api.letsencrypt.org/directory"
  owner                = var.owner
  rke2_version         = var.rke2_version
  rancher_helm_repo    = "https://releases.rancher.com/server-charts"
  rancher_helm_channel = "stable"
  helm_chart_strategy  = "provide"
  # These options use the Let's Encrypt cert that the module generates for you when you deploy the VPC and Domain.
  # WARNING! "hostname" must be an fqdn
  helm_chart_values = {
    "hostname"               = "${local.domain}.${local.zone}"
    "replicas"               = "1"
    "bootstrapPassword"      = "admin"
    "tls"                    = "ingress"
    "ingress.enabled"        = "true"
    "ingress.tls.source"     = "secret"
    "ingress.tls.secretName" = "tls-rancher-ingress"
    "certmanager.version"    = local.cert_manager_version
    "agentTLSMode"           = "strict"
    "privateCA"              = "true"
    "additionalTrustedCAs"   = "true"
  }
  node_configuration = {
    "rancherA" = {
      type            = "all-in-one"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = true
    }
    "rancherB" = {
      type            = "all-in-one"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "rancherC" = {
      type            = "all-in-one"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
  }
  local_file_path      = var.file_path
  runner_ip            = chomp(data.http.myip.response_body) # "runner" is the server running Terraform
  rancher_version      = var.rancher_version
  cert_manager_version = "1.18.3"
  os                   = "sle-micro-61"
}

data "http" "myip" {
  url = "https://ipinfo.io/ip"
}

# you shouldn't do this in production, I am trying to show/prove self-signed certificates working with the Rancher configuration
# this could easily be replaced by some secret resource from Vault or if you are using Terraform 1.11+ you should use the ephemeral resources
module "tls" {
  source = "./modules/tls"
  domain = "${local.domain}.${local.zone}"
}

module "rancher" {
  depends_on = [
    module.tls,
  ]
  source  = "rancher/aws/rancher2"
  version = "3.1.1"
  # project
  identifier   = local.identifier
  owner        = local.owner
  project_name = local.project_name
  domain       = local.domain
  zone         = local.zone
  # access
  key_name = local.key_name
  key      = local.key
  username = local.username
  admin_ip = local.runner_ip
  # rke2
  rke2_version       = local.rke2_version
  local_file_path    = local.local_file_path
  install_method     = "rpm" # rpm only for now, need to figure out local helm chart installs otherwise
  cni                = "canal"
  node_configuration = local.node_configuration
  # rancher
  cert_manager_version            = local.cert_manager_version
  cert_use_strategy               = "supply"
  tls_public_cert                 = module.tls.tls_public_certificate # just the cert, not any CA
  tls_private_key                 = module.tls.tls_private_key
  tls_public_chain                = module.tls.certificate_chain # just the chain, it should not include the cert itself
  rancher_version                 = local.rancher_version
  rancher_helm_repo               = local.rancher_helm_repo
  rancher_helm_channel            = local.rancher_helm_channel
  rancher_helm_chart_use_strategy = local.helm_chart_strategy
  rancher_helm_chart_values       = local.helm_chart_values
  acme_server_url                 = local.acme_server_url
}

provider "rancher2" {
  alias     = "authenticate"
  bootstrap = true
  api_url   = "https://${local.domain}.${local.zone}"
  timeout   = "300s"
  ca_certs  = module.rancher.tls_certificate_chain
}

resource "rancher2_bootstrap" "authenticate" {
  depends_on = [
    module.tls,
    module.rancher,
  ]
  provider         = rancher2.authenticate
  initial_password = module.rancher.admin_password
  password         = module.rancher.admin_password
  token_update     = true
  token_ttl        = 7200 # 2 hours
}

provider "rancher2" {
  alias     = "default"
  api_url   = "https://${local.domain}.${local.zone}"
  token_key = rancher2_bootstrap.authenticate.token
  timeout   = "300s"
  ca_certs  = module.rancher.tls_certificate_chain
}

data "rancher2_cluster" "local" {
  depends_on = [
    module.tls,
    module.rancher,
    rancher2_bootstrap.authenticate,
  ]
  provider = rancher2.default
  name     = "local"
}
