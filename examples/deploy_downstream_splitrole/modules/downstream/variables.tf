
variable "name" {
  type        = string
  description = <<-EOT
    The name to give the new downstream cluster.
  EOT
}
variable "identifier" {
  type        = string
  description = <<-EOT
    A unique string that identifies the project as a whole.
  EOT
}
variable "owner" {
  type        = string
  description = <<-EOT
    A string to identify the "owner" of the resource, this is added as a tag to all resources.
    This is important both for finding resources in the future and removing them in case of problems.
  EOT
}
variable "aws_access_key_id" {
  type        = string
  description = <<-EOT
    The AWS access key to use to deploy the machines.
    We suggest using temporary credentials.
  EOT
}
variable "aws_secret_access_key" {
  type        = string
  description = <<-EOT
    The AWS secret key to use to deploy the machines.
    We suggest using temporary credentials.
  EOT
}
variable "aws_session_token" {
  type        = string
  description = <<-EOT
    The AWS session token to use to deploy the machines.
    While this is not required we suggest using it.
  EOT
  default     = ""
}
variable "aws_region" {
  type        = string
  description = <<-EOT
    The AWS region name to deploy to.
  EOT
}
variable "aws_region_letter" {
  type        = string
  description = <<-EOT
    The letter of the availability zone from the full region name.
    Eg. in AZ "us-west-2a", this would be "a", in "us-west-2b" this would be "b".
  EOT
}
variable "vpc_id" {
  type        = string
  description = <<-EOT
    The id of the VPC where the cluster will be deployed.
    Currently this expects the same VPC as the Rancher cluster.
  EOT
}
variable "security_group_id" {
  type        = string
  description = <<-EOT
    The id of the security group to add the downstream cluster's security group to.
  EOT
}

variable "load_balancer_security_groups" {
  type        = list(any)
  description = <<-EOT
    The security groups for the load balancer.
  EOT
}
variable "subnet_id" {
  type        = string
  description = <<-EOT
    The id of the subnet to deploy the cluster to.
    Currently we only support deploying to the same subnet as the Rancher cluster.
  EOT
}
variable "aws_instance_type" {
  type        = string
  description = <<-EOT
    The AWS instance type to deploy.
  EOT
}
variable "ami_id" {
  type        = string
  description = <<-EOT
    The AWS id for the OS image to deploy the cluster on.
  EOT
}
variable "ami_ssh_user" {
  type        = string
  description = <<-EOT
    The ssh user that is default for the image, the deployment will replace this user.
  EOT
}
variable "ami_admin_group" {
  type        = string
  description = <<-EOT
    The group name for the root/admin user.
    Usually "wheel" or "admin".
  EOT
}
variable "control_plane_node_count" {
  type        = number
  description = <<-EOT
    The number of control plane nodes to deploy.
  EOT
}
variable "worker_node_count" {
  type        = number
  description = <<-EOT
    The number of worker nodes to deploy.
  EOT
}
variable "direct_node_access" {
  type = object({
    runner_ip       = string
    ssh_access_key  = string
    ssh_access_user = string
  })
  description = <<-EOT
    This object has the information to give the computer running Terraform (runner) the ability to directly access the nodes provisioned by Rancher.
    When left in the default state no access will be given.
    This requires the downstream cluster to have public IP addresses.
  EOT
  default     = null
}
variable "rke2_version" {
  type        = string
  description = <<-EOT
    The RKE2 version to deploy.
    This module only supports RKE2.
  EOT
}
