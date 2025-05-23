provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
  region = local.aws_region
}

provider "acme" {
  server_url = "${local.acme_server_url}/directory"
}

provider "github" {}
provider "kubernetes" {} # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)
provider "helm" {}       # make sure you set the env variable KUBE_CONFIG_PATH to local_file_path (file_path variable)

provider "rancher2" {
  api_url   = "https://${local.domain}.${local.zone}"
  token_key = module.rancher.admin_token
  timeout   = "3000s"
}

locals {
  identifier            = var.identifier
  example               = "basic"
  project_name          = "tf-${substr(md5(join("-", [local.example, local.identifier])), 0, 5)}"
  username              = local.project_name
  domain                = local.project_name
  zone                  = var.zone
  key_name              = var.key_name
  key                   = var.key
  owner                 = var.owner
  rke2_version          = var.rke2_version
  local_file_path       = var.file_path
  runner_ip             = chomp(data.http.myip.response_body) # "runner" is the server running Terraform
  rancher_version       = var.rancher_version
  cert_manager_version  = "1.16.3" #"1.13.1"
  os                    = "sle-micro-61"
  acme_server_url       = "https://acme-v02.api.letsencrypt.org"
  aws_access_key_id     = var.aws_access_key_id
  aws_secret_access_key = var.aws_secret_access_key
  aws_region            = var.aws_region
  aws_session_token     = var.aws_session_token
  email                 = (var.email != "" ? var.email : "${local.identifier}@${local.zone}")
}

data "http" "myip" {
  url = "https://ipinfo.io/ip"
}

module "rancher" {
  source  = "rancher/aws/rancher2"
  version = "1.2.2"
  # project
  identifier                   = local.identifier
  owner                        = local.owner
  project_name                 = local.project_name
  domain                       = local.domain
  zone                         = local.zone
  skip_project_cert_generation = true
  # access
  key_name = local.key_name
  key      = local.key
  username = local.username
  admin_ip = local.runner_ip
  # rke2
  rke2_version    = local.rke2_version
  local_file_path = local.local_file_path
  install_method  = "rpm" # rpm only for now, need to figure out local helm chart installs otherwise
  cni             = "canal"
  node_configuration = {
    "rancher" = {
      type            = "all-in-one"
      size            = "large"
      os              = local.os
      indirect_access = true
      initial         = true
    }
  }
  # rancher
  rancher_version        = local.rancher_version
  cert_manager_version   = local.cert_manager_version
  configure_cert_manager = true
  cert_manager_configuration = {
    aws_access_key_id     = local.aws_access_key_id
    aws_secret_access_key = local.aws_secret_access_key
    aws_region            = local.aws_region
    aws_session_token     = local.aws_session_token
    acme_email            = local.email
    acme_server_url       = local.acme_server_url
  }
}

module "rke2_image" {
  source              = "rancher/server/aws"
  version             = "v1.4.0"
  server_use_strategy = "skip"
  image_use_strategy  = "find"
  image_type          = local.os
}

# you can add this one multiple times, or use a loop to deploy multiple clusters
module "downstream" {
  depends_on = [
    module.rancher,
    module.rke2_image,
  ]
  source = "./modules/downstream"
  # general
  name       = "${local.project_name}-mp" # multi-pool
  identifier = local.identifier
  owner      = local.owner
  # aws access
  aws_access_key_id     = local.aws_access_key_id
  aws_secret_access_key = local.aws_secret_access_key
  aws_session_token     = trimspace(chomp(local.aws_session_token))
  aws_region            = local.aws_region
  aws_region_letter = replace(
    module.rancher.subnets[keys(module.rancher.subnets)[0]].availability_zone,
    local.aws_region,
    ""
  )
  # aws project info
  vpc_id                        = module.rancher.vpc.id
  security_group_id             = module.rancher.security_group.id
  load_balancer_security_groups = module.rancher.load_balancer_security_groups
  subnet_id                     = module.rancher.subnets[keys(module.rancher.subnets)[0]].id
  # node info
  aws_instance_type        = "m5.large"
  ami_id                   = module.rke2_image.image.id
  ami_ssh_user             = module.rke2_image.image.user
  ami_admin_group          = module.rke2_image.image.admin_group
  worker_node_count        = 2
  control_plane_node_count = 2
  direct_node_access = {
    runner_ip       = local.runner_ip
    ssh_access_key  = local.key
    ssh_access_user = local.project_name
  }
  # rke2 info
  rke2_version = local.rke2_version
}
