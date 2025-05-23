provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
}

provider "acme" {
  server_url = "https://acme-v02.api.letsencrypt.org/directory"
}

provider "github" {}
provider "kubernetes" {} # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)
provider "helm" {}       # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)

provider "rancher2" {
  alias     = "default"
  api_url   = "https://${local.domain}.${local.zone}"
  token_key = module.this.admin_token
  timeout   = "300s"
}

provider "rancher2" {
  alias     = "authenticate"
  bootstrap = true
  api_url   = "https://${local.domain}.${local.zone}"
  timeout   = "300s"
}

resource "rancher2_bootstrap" "authenticate" {
  depends_on = [
    module.this,
  ]
  provider         = rancher2.authenticate
  initial_password = module.this.admin_password
  password         = module.this.admin_password
  token_update     = true
  token_ttl        = 7200 # 2 hours
}


locals {
  identifier           = var.identifier
  example              = "prod"
  project_name         = "tf-${substr(md5(join("-", [local.example, local.identifier])), 0, 5)}"
  username             = local.project_name
  domain               = local.project_name
  zone                 = var.zone
  key_name             = var.key_name
  key                  = var.key
  owner                = var.owner
  rke2_version         = var.rke2_version
  local_file_path      = var.file_path
  os                   = "sle-micro-61"
  runner_ip            = chomp(data.http.myip.response_body) # "runner" is the server running Terraform
  rancher_version      = var.rancher_version
  cert_manager_version = "1.16.3" #"1.13.1"
  node_configuration = {
    "initial" = {
      type            = "database"
      size            = "xl"
      os              = local.os
      indirect_access = false
      initial         = true # this will set the first server as the inital node
    }
    "db2" = {
      type            = "database"
      size            = "xl"
      os              = local.os
      indirect_access = false
      initial         = false
    }
    "db3" = {
      type            = "database"
      size            = "xl"
      os              = local.os
      indirect_access = false
      initial         = false
    }
    "api1" = {
      type            = "api"
      size            = "xl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "api2" = {
      type            = "api"
      size            = "xl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "api3" = {
      type            = "api"
      size            = "xl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "wrk1" = {
      type            = "worker"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "wrk2" = {
      type            = "worker"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = false
    }
    "wrk3" = {
      type            = "worker"
      size            = "xxl"
      os              = local.os
      indirect_access = true
      initial         = false
    },
  }
}

data "http" "myip" {
  url = "https://ipinfo.io/ip"
}

module "this" {
  source               = "rancher/aws/rancher2"
  version              = "1.2.2"
  identifier           = local.identifier
  owner                = local.owner
  project_name         = local.project_name
  domain               = local.domain
  zone                 = local.zone
  key_name             = local.key_name
  key                  = local.key
  username             = local.username
  admin_ip             = local.runner_ip
  rke2_version         = local.rke2_version
  local_file_path      = local.local_file_path
  install_method       = "rpm" # rpm only for now, need to figure out local helm chart installs otherwise
  cni                  = "canal"
  node_configuration   = local.node_configuration
  cert_manager_version = local.cert_manager_version
  rancher_version      = local.rancher_version
}

# this will fail if the default self signed cert is found
resource "terraform_data" "get_cert_info" {
  depends_on = [
    module.this,
  ]
  provisioner "local-exec" {
    command = <<-EOT
      CERT="$(echo | openssl s_client -showcerts -servername ${local.domain}.${local.zone} -connect ${local.domain}.${local.zone}:443 2>/dev/null | openssl x509 -inform pem -noout -text)"
      echo "$CERT"
      FAKE="$(echo "$CERT" | grep 'Kubernetes Ingress Controller Fake Certificate')"
      if [ -z "$FAKE" ]; then
        echo "cert is not fake"
        exit 0
      else
        echo "cert is fake"
        exit 1
      fi
    EOT
  }
}

data "rancher2_cluster" "local" {
  depends_on = [
    module.this,
    rancher2_bootstrap.authenticate,
  ]
  provider = rancher2.default
  name     = "local"
}
