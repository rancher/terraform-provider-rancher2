
locals {
  # general
  identifier   = var.identifier
  owner        = var.owner
  cluster_name = var.name
  # aws access
  aws_access_key_id     = var.aws_access_key_id
  aws_secret_access_key = var.aws_secret_access_key
  aws_session_token     = var.aws_session_token
  aws_region            = var.aws_region
  aws_region_letter     = var.aws_region_letter
  # networking info
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  lbsg              = sort(var.load_balancer_security_groups)
  load_balancer_security_group_id = [
    for i in range(length(local.lbsg)) :
    local.lbsg[i] if local.lbsg[i] != local.security_group_id
    # load balancers only have 2 security groups, the project and its own
    # this eliminates the project security group to just return the load balancer's security group
  ][0]
  downstream_security_group_name = "${local.cluster_name}-sgroup"
  # node info
  aws_instance_type = var.aws_instance_type
  ami_id            = var.ami_id
  ami_ssh_user      = var.ami_ssh_user
  node_count        = var.node_count
  node_ips          = { for i in range(local.node_count) : tostring(i) => data.aws_instances.rke2_instance_nodes.public_ips[i] }
  node_id           = "${local.cluster_name}-nodes"
  ami_admin_group   = (var.ami_admin_group != "" ? var.ami_admin_group : "tty")
  runner_ip         = (var.direct_node_access != null ? var.direct_node_access.runner_ip : "10.1.1.1") # the IP running Terraform
  ssh_access_key    = (var.direct_node_access != null ? var.direct_node_access.ssh_access_key : "fake123abc")
  ssh_access_user   = (var.direct_node_access != null ? var.direct_node_access.ssh_access_user : "fake")
  # rke2 info
  rke2_version = var.rke2_version
}

resource "aws_security_group" "downstream_cluster" {
  description = "Access to downstream cluster"
  name        = local.downstream_security_group_name
  vpc_id      = local.vpc_id
  tags = {
    Name = local.downstream_security_group_name
  }
  lifecycle {
    ignore_changes = [
      ingress,
      egress,
    ]
  }
}
# this allows servers attached to the project security group to accept connections initiated by the downstream cluster
resource "aws_vpc_security_group_ingress_rule" "downstream_ingress_rancher" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  referenced_security_group_id = aws_security_group.downstream_cluster.id
  security_group_id            = local.security_group_id
  ip_protocol                  = "-1"
}
# this allows the load balancer to accept connections initiated by the downstream cluster
resource "aws_vpc_security_group_ingress_rule" "downstream_ingress_loadbalancer" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  referenced_security_group_id = aws_security_group.downstream_cluster.id
  security_group_id            = local.load_balancer_security_group_id
  ip_protocol                  = "-1"
}

# this allows the downstream cluster to reach out to any public ipv4 address
resource "aws_vpc_security_group_egress_rule" "downstream_egress_ipv4" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
  security_group_id = aws_security_group.downstream_cluster.id
}
# this allows the downstream cluster to reach out to any public ipv6 address
resource "aws_vpc_security_group_egress_rule" "downstream_egress_ipv6" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  ip_protocol       = "-1"
  cidr_ipv6         = "::/0"
  security_group_id = aws_security_group.downstream_cluster.id
}
# this allows the downstream cluster to reach out to any server attached to the project security group
resource "aws_vpc_security_group_egress_rule" "downstream_egress_project_link" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  referenced_security_group_id = local.security_group_id
  security_group_id            = aws_security_group.downstream_cluster.id
  ip_protocol                  = "-1"
}
# this allows nodes to talk to each other
resource "aws_vpc_security_group_ingress_rule" "downstream_ingress_internal_ipv4" {
  depends_on = [
    aws_security_group.downstream_cluster,
  ]
  ip_protocol       = "-1"
  cidr_ipv4         = "10.0.0.0/16"
  security_group_id = aws_security_group.downstream_cluster.id
}
resource "rancher2_machine_config_v2" "all_in_one" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
  ]
  generate_name = local.cluster_name
  amazonec2_config {
    ami            = local.ami_id
    region         = local.aws_region
    security_group = [local.downstream_security_group_name]
    subnet_id      = local.subnet_id
    vpc_id         = local.vpc_id
    zone           = local.aws_region_letter
    session_token  = local.aws_session_token
    instance_type  = local.aws_instance_type
    ssh_user       = local.ami_ssh_user
    tags           = join(",", ["Id", local.identifier, "Owner", local.owner, "NodeId", local.node_id])
    userdata       = <<-EOT
      #cloud-config

      merge_how:
       - name: list
         settings: [replace]
       - name: dict
         settings: [replace]

      users:
        - name: ${local.ssh_access_user}
          gecos: ${local.ssh_access_user}
          sudo: ALL=(ALL) NOPASSWD:ALL
          groups: users, ${local.ami_admin_group}
          lock_passwd: true
          ssh_authorized_keys:
            - ${local.ssh_access_key}
          homedir: /home/${local.ssh_access_user}
    EOT
  }
}
resource "terraform_data" "patch_machine_configs" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
  ]
  triggers_replace = {
    core_config = rancher2_machine_config_v2.all_in_one.id
  }
  provisioner "local-exec" {
    command = <<-EOT
      ${path.module}/addKeyToAmazonConfig.sh "${local.aws_access_key_id}" "${local.aws_secret_access_key}"
    EOT
  }
}

resource "rancher2_cluster_v2" "rke2_cluster" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
  ]
  name                  = local.cluster_name
  kubernetes_version    = local.rke2_version
  enable_network_policy = true
  rke_config {
    machine_pools {
      name               = local.cluster_name
      control_plane_role = true
      etcd_role          = true
      worker_role        = true
      quantity           = local.node_count
      machine_config {
        kind = rancher2_machine_config_v2.all_in_one.kind
        name = rancher2_machine_config_v2.all_in_one.name
      }
    }
  }
  timeouts {
    create = "120m"
  }
}

resource "time_sleep" "wait_for_nodes" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
  ]
  create_duration = "120s"
}

data "aws_instances" "rke2_instance_nodes" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
    time_sleep.wait_for_nodes,
  ]
  filter {
    name   = "tag:NodeId"
    values = [local.node_id]
  }
}

# this allows the load balancer to accept connections initiated by the downstream cluster's public ip addresses
# this weird in-flight grab of the nodes and manipulating the security groups is not good,
#  but the only way to allow ingress when the downstream cluster has public IPs
# FYI: security group references only work with private IPs
resource "aws_vpc_security_group_ingress_rule" "downstream_public_ingress_loadbalancer" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
    time_sleep.wait_for_nodes,
    data.aws_instances.rke2_instance_nodes,
  ]
  for_each          = local.node_ips
  security_group_id = local.load_balancer_security_group_id
  ip_protocol       = "-1"
  cidr_ipv4         = "${each.value}/32"
}

resource "aws_vpc_security_group_ingress_rule" "downstream_public_ingress_runner" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
    time_sleep.wait_for_nodes,
    data.aws_instances.rke2_instance_nodes,
  ]
  security_group_id = aws_security_group.downstream_cluster.id
  ip_protocol       = "tcp"
  from_port         = 22
  to_port           = 22
  cidr_ipv4         = "${local.runner_ip}/32"
}

resource "rancher2_cluster_sync" "sync" {
  depends_on = [
    aws_security_group.downstream_cluster,
    aws_vpc_security_group_ingress_rule.downstream_ingress_rancher,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv4,
    aws_vpc_security_group_egress_rule.downstream_egress_ipv6,
    aws_vpc_security_group_egress_rule.downstream_egress_project_link,
    rancher2_machine_config_v2.all_in_one,
    terraform_data.patch_machine_configs,
    rancher2_cluster_v2.rke2_cluster,
  ]
  cluster_id = rancher2_cluster_v2.rke2_cluster.cluster_v1_id
}
