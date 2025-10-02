# There are many ways to orchestrate Terraform configurations with the goal of breaking it down
# I am using Terraform resources to orchestrate Terraform
# I felt this was the best way to accomplish the goal without incurring additional dependencies

locals {
  inputs         = var.inputs
  inputs_hash    = md5(local.inputs)
  template_path  = var.template_path
  template_files = var.template_files
  # tflint-ignore: terraform_unused_declarations
  fail_no_template = ((local.template_path == null && local.template_files == null) ? one([local.template_path, "missing_template"]) : false)
  # tflint-ignore: terraform_unused_declarations
  fail_too_much_template = ((local.template_path != null && local.template_files != null) ? one([local.template_path, "template_path_or_template_files"]) : false)
  # template_file_list should be a list of absolute file paths (not relative)
  template_file_list = (
    local.template_path != null ?
    [
      for i in range(length(fileset(local.template_path, "**"))) :
      join("/", [local.template_path, tolist(fileset(local.template_path, "**"))[i]])
    ]
    : local.template_files
  )
  template_files_hash          = md5(join("-", [for f in local.template_file_list : filemd5(f)]))
  template_file_list_basenames = [for f in local.template_file_list : basename(f)]
  outputs                      = contains(local.template_file_list_basenames, "outputs.tf")

  deploy_path           = chomp(var.deploy_path)
  environment_variables = var.environment_variables
  export_contents = (
    local.environment_variables != null ?
    join(";", [for k, v in local.environment_variables : "export ${k}=${v}"])
    : ""
  )
  export_hash  = md5(local.export_contents)
  attempts     = var.attempts
  interval     = var.interval
  timeout      = var.timeout
  init         = var.init
  init_script  = (local.init ? "terraform init -upgrade" : "")
  tf_data_dir  = abspath(var.data_path != null ? var.data_path : path.root)
  skip_destroy = (var.skip_destroy ? "true" : "")
}

resource "terraform_data" "path" {
  triggers_replace = {
    inputs = local.inputs_hash
    files  = local.template_files_hash
    env    = local.export_hash
  }
  provisioner "local-exec" {
    command = <<-EOT
      install -d ${local.deploy_path}
    EOT
  }
}

# keep the deploy config up to date with the template
resource "local_file" "copy_template" {
  depends_on = [
    terraform_data.path
  ]
  for_each = toset(local.template_file_list)
  content  = file(each.key)
  filename = "${local.deploy_path}/${basename(each.key)}"
}

resource "local_file" "inputs" {
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
  ]
  content  = local.inputs
  filename = "${local.deploy_path}/inputs.tfvars"
}

resource "terraform_data" "destroy" {
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
  ]
  triggers_replace = {
    inputs = local.inputs_hash
    files  = local.template_files_hash
    env    = local.export_hash
    ec     = local.export_contents
    dp     = local.deploy_path
    to     = local.timeout
    dd     = local.tf_data_dir
    sd     = local.skip_destroy
  }
  provisioner "local-exec" {
    when    = destroy
    command = <<-EOT
      ${self.triggers_replace.ec}
      cd ${self.triggers_replace.dp}
      export TF_DATA_DIR="${self.triggers_replace.dd}"
      if [ -z "${self.triggers_replace.sd}" ]; then
        timeout -k 1m ${self.triggers_replace.to} terraform destroy -var-file="${self.triggers_replace.dp}/inputs.tfvars" -auto-approve -state="${self.triggers_replace.dp}/tfstate" || true
      else
        echo "Not destroying deployed module, it will no longer be managed here."
      fi
    EOT
  }
}

resource "terraform_data" "create" {
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
    terraform_data.destroy,
  ]
  triggers_replace = {
    inputs = local.inputs_hash
    files  = local.template_files_hash
    env    = local.export_hash
  }
  provisioner "local-exec" {
    command = <<-EOT
      ${local.export_contents}
      cd ${local.deploy_path}
      export TF_DATA_DIR="${local.tf_data_dir}"

      ${local.init_script}

      MAX=${local.attempts}
      EXITCODE=1
      ATTEMPTS=0
      E=1
      E1=0
      while [ $EXITCODE -gt 0 ] && [ $ATTEMPTS -lt $MAX ]; do
        A=0
        while [ $E -gt 0 ] && [ $A -lt $MAX ]; do
          timeout -k 1m ${local.timeout} terraform apply -var-file="${local.deploy_path}/inputs.tfvars" -auto-approve -state="${local.deploy_path}/tfstate"
          E=$?
          if [ $E -eq 124 ]; then echo "Apply timed out after ${local.timeout}"; fi
          A=$((A+1))
        done
        # don't destroy if the last attempt fails
        if [ $E -gt 0 ] && [ $ATTEMPTS != $((MAX-1)) ]; then
          A1=0
          while [ $E1 -gt 0 ] && [ $A1 -lt $MAX ]; do
            timeout -k 1m ${local.timeout} terraform destroy -var-file="${local.deploy_path}/inputs.tfvars" -auto-approve -state="${local.deploy_path}/tfstate"
            E1=$?
            if [ $E1 -eq 124 ]; then echo "Apply timed out after ${local.timeout}"; fi
            A1=$((A1+1))
          done
        fi
        if [ $E -gt 0 ]; then
          echo "apply failed..."
        fi
        if [ $E1 -gt 0 ]; then
          echo "destroy failed..."
        fi
        if [ $E -gt 0 ] || [ $E1 -gt 0 ]; then
          EXITCODE=1
        else
          EXITCODE=0
        fi
        ATTEMPTS=$((ATTEMPTS+1))
        if [ $EXITCODE -gt 0 ] && [ $ATTEMPTS -lt $MAX ]; then
          echo "wait ${local.interval} seconds between attempts..."
          sleep ${local.interval}
        fi
      done
      if [ $ATTEMPTS -eq $MAX ]; then echo "max attempts reached..."; fi
      if [ $EXITCODE -ne 0 ]; then echo "failure, exit code $EXITCODE..."; fi
      if [ $EXITCODE -eq 0 ]; then 
        echo "success...";
        terraform output -json -state="${local.deploy_path}/tfstate" > ${local.deploy_path}/outputs.json
      fi
      exit $EXITCODE
    EOT
  }
}

data "local_file" "outputs" {
  count = (local.outputs ? 1 : 0)
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
    terraform_data.create,
  ]
  filename = "${local.deploy_path}/outputs.json"
}

data "local_file" "state" {
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
    terraform_data.create,
  ]
  filename = "${local.deploy_path}/tfstate"
}

# this ensures that the parent state can hold the deployed modules outputs
# when you plan on a new server you will see the need to deploy the outputs file if it doesn't exist
resource "local_file" "outputs" {
  count = (local.outputs ? 1 : 0)
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
    terraform_data.create,
  ]
  content  = data.local_file.outputs[0].content
  filename = "${local.deploy_path}/outputs.json"
}

# this ensures that the parent state can hold the deployed modules state
# when you plan on a new server you will see the need to deploy the state file if it doesn't exist
resource "local_file" "state_backup" {
  depends_on = [
    terraform_data.path,
    local_file.copy_template,
    local_file.inputs,
    terraform_data.create,
  ]
  content  = data.local_file.state.content
  filename = "${local.deploy_path}/tfstate.backup"
}
