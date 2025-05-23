output "kubeconfig" {
  value       = module.this.kubeconfig
  description = <<-EOT
    The kubeconfig for the server.
  EOT
  sensitive   = true
}
output "address" {
  value = module.this.address
}
output "admin_token" {
  value     = module.this.admin_token
  sensitive = true
}
output "admin_password" {
  value     = module.this.admin_password
  sensitive = true
}
output "cluster_data" {
  value     = jsonencode(data.rancher2_cluster.local)
  sensitive = true
}
