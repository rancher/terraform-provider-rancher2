output "kubeconfig" {
  value       = module.rancher.kubeconfig
  description = <<-EOT
    The kubeconfig for the server.
  EOT
  sensitive   = true
}
output "address" {
  value = module.rancher.address
}
output "admin_password" {
  value     = module.rancher.admin_password
  sensitive = true
}
output "downstream_kubeconfig" {
  value     = data.rancher2_cluster_v2.downstream.kube_config
  sensitive = true
}
