output "cluster_id" {
  description = "The Rancher cluster ID."
  value       = rancher2_cluster_v2.cluster.id
}

output "cluster_v1_id" {
  description = "The Rancher v1 cluster ID (used for rancher2_cluster_sync)."
  value       = rancher2_cluster_v2.cluster.cluster_v1_id
}

output "server_machine_config_name" {
  description = "Generated name of the server machine config."
  value       = rancher2_machine_config_v2.server.name
}

output "worker_machine_config_name" {
  description = "Generated name of the worker machine config."
  value       = rancher2_machine_config_v2.worker.name
}
