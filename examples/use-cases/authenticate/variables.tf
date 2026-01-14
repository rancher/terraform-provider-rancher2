variable "identifier" {
  type        = string
  description = <<-EOT
    A unique string used to discern resources between tests in remote APIs.
  EOT
}
variable "owner" {
  type        = string
  description = <<-EOT
    The owner to tag on all of the AWS resources, must be a valid email address.
  EOT
}
variable "rancher_url" {
  type        = string
  description = <<-EOT
    The Rancher URL to configure.
    Format should be "https://<hostname>.<domain>.<tld>"
  EOT
}
