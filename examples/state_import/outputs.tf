output "kubeconfig" {
  value       = module.rancher.kubeconfig
  description = <<-EOT
    The kubeconfig for the rancher cluster.
  EOT
  sensitive   = true
}
output "address" {
  value = module.rancher.address
  description = <<-EOT
    The url for the rancher cluster.
  EOT
}
output "admin_password" {
  value     = module.rancher.admin_password
  description = <<-EOT
    The admin password for the rancher cluster.
  EOT
  sensitive = true
}
output "downstream_kubeconfig" {
  value     = data.rancher2_cluster.downstream_cluster.kube_config
  description = <<-EOT
    The kubeconfig for the downstream cluster.
  EOT
  sensitive = true
}
