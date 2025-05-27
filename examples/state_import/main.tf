provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
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
  timeout   = "300s"
}

locals {
  identifier            = var.identifier
  example               = "state_import_downstream"
  project_name          = "tf-${substr(md5(join("-", [local.example, local.identifier])), 0, 5)}"
  tf_data_dir           = abspath(var.data_path != null ? var.data_path : path.root)
  username              = local.project_name
  domain                = local.project_name
  zone                  = var.zone
  key_name              = var.key_name
  key                   = var.key
  owner                 = var.owner
  rke2_version          = var.rke2_version
  local_file_path       = abspath(var.file_path)
  runner_ip             = chomp(data.http.myip.response_body) # "runner" is the server running Terraform
  rancher_version       = var.rancher_version
  cert_manager_version  = "1.16.3" #"1.13.1"
  os                    = "sle-micro-61"
  aws_access_key_id     = var.aws_access_key_id
  aws_secret_access_key = var.aws_secret_access_key
  aws_region            = var.aws_region
  aws_session_token     = var.aws_session_token
  aws_instance_type     = "m5.large"
  node_count            = 3
  email                 = (var.email != "" ? var.email : "${local.identifier}@${local.zone}")
  acme_server_url       = "https://acme-v02.api.letsencrypt.org"
  cluster_name          = "tf-all-in-one-config"
  project_mismatch      = var.project_mismatch # if this is true, then the import should fail
  project_id            = (local.project_mismatch ? rancher2_project.test.id : data.rancher2_cluster.downstream_cluster.default_project_id)
  # tflint-ignore: terraform_unused_declarations
  fail_project_id = (strcontains(local.project_id, ":") != true ? one([local.project_id, "project_id_malformed"]) : false)
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
    aws_session_token     = local.aws_session_token
    aws_region            = local.aws_region
    acme_email            = local.email
    acme_server_url       = local.acme_server_url
  }
}

module "rke2_image" {
  source              = "rancher/server/aws"
  version             = "v1.4.0"
  server_use_strategy = "skip"
  image_use_strategy  = "find"
  image_type          = local.os # this is not required to match Rancher, it just seemed easier in this example
}

module "downstream_cluster" {
  source = "./modules/downstream"
  # general
  name            = local.cluster_name
  identifier      = local.identifier
  owner           = local.owner
  domain          = local.domain
  zone            = local.zone
  kubeconfig_path = "${local.local_file_path}/kubeconfig"
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
  aws_instance_type = local.aws_instance_type
  ami_id            = module.rke2_image.image.id
  ami_ssh_user      = module.rke2_image.image.user
  ami_admin_group   = module.rke2_image.image.admin_group
  node_count        = local.node_count
  direct_node_access = {
    runner_ip       = local.runner_ip
    ssh_access_key  = chomp(local.key)
    ssh_access_user = local.project_name
  }
  # rke2 info
  rke2_version = local.rke2_version
  # rancher info
  rancher_token = module.rancher.admin_token
}

data "rancher2_cluster" "downstream_cluster" {
  depends_on = [
    module.rancher,
    module.rke2_image,
    module.downstream_cluster,
  ]
  name = local.cluster_name
}

resource "rancher2_namespace" "test" {
  depends_on = [
    module.rancher,
    module.rke2_image,
    module.downstream_cluster,
    data.rancher2_cluster.downstream_cluster,
  ]
  name        = "test"
  project_id  = data.rancher2_cluster.downstream_cluster.default_project_id
  description = "testing namespace"
  resource_quota {
    limit {
      limits_cpu       = "100m"
      limits_memory    = "100Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu      = "20m"
    limits_memory   = "20Mi"
    requests_cpu    = "1m"
    requests_memory = "1Mi"
  }
}

resource "rancher2_project" "test" {
  depends_on = [
    module.rancher,
    module.rke2_image,
    module.downstream_cluster,
    data.rancher2_cluster.downstream_cluster,
    rancher2_namespace.test,
  ]
  name       = "test"
  cluster_id = data.rancher2_cluster.downstream_cluster.id
}

resource "local_file" "import_main" {
  depends_on = [
    module.rancher,
    module.rke2_image,
    module.downstream_cluster,
    data.rancher2_cluster.downstream_cluster,
    rancher2_namespace.test,
    rancher2_project.test,
  ]
  filename = "${local.tf_data_dir}/tf-rancher-imported/main.tf"
  content = templatefile("${path.module}/modules/import/main.tftpl", {
    cluster_id   = module.downstream_cluster.cluster_id
    namespace_id = join(".", [local.project_id, rancher2_namespace.test.id])
  })
}

module "import" {
  source = "./modules/deploy"
  depends_on = [
    module.downstream_cluster,
    local_file.import_main,
    data.rancher2_cluster.downstream_cluster,
    rancher2_namespace.test,
  ]
  deploy_path = "${local.tf_data_dir}/tf-rancher-imported"
  data_path   = local.tf_data_dir
  template_files = [
    "${abspath(path.module)}/modules/import/cloud-config.tftpl",
    "${abspath(path.module)}/modules/import/variables.tf",
    "${abspath(path.module)}/modules/import/versions.tf",
  ]
  inputs       = <<-EOT
    cluster_name        = "${local.cluster_name}"
    rke2_version        = "${local.rke2_version}"
    node_count          = "${local.node_count}"
    rancher_key         = "${module.rancher.admin_token}"
    domain              = "${local.domain}"
    zone                = "${local.zone}"
    machine_config_kind = "${module.downstream_cluster.machine_config_kind}"
    machine_config_name = "${module.downstream_cluster.machine_config_name}"
    project_mismatch    = "${local.project_mismatch}"
  EOT
  skip_destroy = true  // this is for testing purposes, it prevents an issue where the imported resources destroy the API objects and the main resources error out on destroy (not found)
  init         = false // this is for testing purposes, it allow us to use dev overrides in the terraformrc to use the locally built binary rather than the registry provider
}
