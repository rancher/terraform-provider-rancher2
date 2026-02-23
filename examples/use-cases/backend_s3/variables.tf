variable "identifier" {
  type        = string
  description = <<-EOT
    Unique ID for the resource, a tag will be added to the resource.
    This helps with identifying and cleaning up resources.
  EOT
}
variable "owner" {
  type        = string
  description = <<-EOT
    Owner tag to be added to the resource, helps when identifying and cleaning up resources.
    Often this is an email address, so that someone can see and contact the person who generated the object.
  EOT
}
