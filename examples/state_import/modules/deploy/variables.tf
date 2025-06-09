variable "inputs" {
  type        = string
  description = <<-EOT
    Contents of an inputs.tfvars file to save in the deployment path.
  EOT
  default     = ""
}
variable "template_path" {
  type        = string
  description = <<-EOT
    Path to the module to deploy.
    These files will be copied to the deploy path, not used directly.
    This is optional, but one of template_path or template_files must be specified.
    Only one of template_path or template_files can be specified.
  EOT
  default     = null
}
variable "template_files" {
  type        = list(any)
  description = <<-EOT
    List of file paths that will be copied to the deploy path.
    This is optional, but one of template_path or template_files must be specified.
    Only one of template_path or template_files can be specified.
  EOT
  default     = null
}
variable "deploy_path" {
  type        = string
  description = <<-EOT
    Path to preform deployment in, this will be Terraform's working directory.
  EOT
}
variable "data_path" {
  type        = string
  description = <<-EOT
    Should match your TF_DATA_DIR environment variable.
    This directory is used to stage all of the various files for your implementation.
    If left null, this will match "path.root".
    This should be a full path, not relative.
  EOT
  default     = null
}
variable "environment_variables" {
  type        = map(any)
  description = <<-EOT
    Map of environment variables to set before running Terraform.
    Key is the name and Value is the value of the variable.
    We export this before running Terraform, eg. "export KEY_1=VARIABLE_1;export KEY_2=VARIABLE_2".
  EOT
  default     = null
}
variable "attempts" {
  type        = number
  description = <<-EOT
    Number of attempts to deploy module.
    Each time Terraform apply is run we check for a successful exit code,
     if the exit code !=0 then we try again, up to the value set in this argument.
  EOT
  default     = 3
}
variable "interval" {
  type        = number
  description = <<-EOT
    A number of seconds to sleep between Terraform apply or destroy attempts.
  EOT
  default     = 30
}
variable "timeout" {
  type        = string
  description = <<-EOT
    A (linux coreutils) timeout DURATION string.
    This will be used to kill the Terraform run in case there is an endless loop.
    If this DURATION is reached a single TERM will be sent, then KILL 1 minute later.
  EOT
  default     = "45m"
}
variable "init" {
  type        = bool
  description = <<-EOT
    Set to false to prevent running Terraform init.
    This is helpful when testing a local bin version of the provider.
  EOT
  default     = true
}
variable "skip_destroy" {
  type        = bool
  description = <<-EOT
    Set to true to ignore calls to destroy the deployed substate.
    State and deploy path will still exist, this essentially divorces the parent from the child.
    This only effects specifically calls to destroy the deploy module, not taint or recreate.
    Be careful as this can leave objects in your API unmanaged by IAC.
  EOT
  default     = false
}
