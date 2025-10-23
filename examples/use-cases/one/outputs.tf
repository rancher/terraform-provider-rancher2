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
output "admin_token" {
  value     = module.rancher.admin_token
  sensitive = true
}
output "admin_password" {
  value     = module.rancher.admin_password
  sensitive = true
}
output "cluster_data" {
  value     = jsonencode(data.rancher2_cluster.local)
  sensitive = true
}
