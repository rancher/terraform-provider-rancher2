output "outputs_json" {
  value = try(local_file.outputs[0].content, "")
}

output "output" {
  value = try({
    for i in range(length(keys(jsondecode(local_file.outputs[0].content)))) :
    keys(jsondecode(local_file.outputs[0].content))[i] => jsondecode(local_file.outputs[0].content)[keys(jsondecode(local_file.outputs[0].content))[i]].value
  }, "")
}

output "state" {
  value = data.local_file.state
}
