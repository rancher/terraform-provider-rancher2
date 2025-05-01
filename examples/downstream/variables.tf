variable "key_name" {
  type        = string
  description = <<-EOT
    The name of an AWS key pair to use for SSH access to the instance.
    This key should already be added to your ssh agent for server authentication.
  EOT
}
variable "key" {
  type        = string
  description = <<-EOT
    The contents of an AWS key pair to use for SSH access to the instance.
    This is necessary for installing rke2 on the nodes and will be removed after installation.
  EOT
}
variable "identifier" {
  type        = string
  description = <<-EOT
    A unique identifier for the project, this helps when generating names for infrastructure items."
  EOT
}
variable "owner" {
  type        = string
  description = <<-EOT
    The owner of the project, this helps when generating names for infrastructure items."
  EOT
}
variable "zone" {
  type        = string
  description = <<-EOT
    The Route53 DNS zone to deploy the cluster into.
    This is used to generate the DNS name for the cluster.
    The zone must already exist.
  EOT
}
variable "rke2_version" {
  type        = string
  description = <<-EOT
    The version of rke2 to install on the nodes.
    eg. v1.30.2+rke2r1
  EOT
}
variable "rancher_version" {
  type        = string
  description = <<-EOT
    The version of rancher to install on the rke2 cluster.
  EOT
  default     = "2.9.2"
}
variable "file_path" {
  type        = string
  description = <<-EOT
    The path to the file containing the rke2 install script.
  EOT
  default     = "./rke2"
}
variable "aws_access_key_id" {
  type        = string
  description = <<-EOT
    AWS access key ID.
  EOT
  sensitive   = true
}
variable "aws_secret_access_key" {
  type        = string
  description = <<-EOT
    AWS secret key for EC2 services.
  EOT
  sensitive   = true
}
variable "aws_session_token" {
  type        = string
  description = <<-EOT
    AWS session token for EC2 services.
    If left empty the AWS provider will assume you are using permanent AWS credentials.
  EOT
  sensitive   = true
  default     = ""
}
variable "aws_region" {
  type        = string
  description = <<-EOT
    AWS region EC2 services.
  EOT
  sensitive   = true
}
variable "email" {
  type        = string
  description = <<-EOT
    Email used for TLS certification registration.
    If left blank this will be <identifier>@<zone>.
  EOT
  default     = ""
}
