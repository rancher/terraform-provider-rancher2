output "cluster_id" {
  value = rancher2_cluster_v2.rke2_cluster.id
}

output "cluster_v1_id" {
  value = rancher2_cluster_v2.rke2_cluster.cluster_v1_id
}

output "machine_config_kind" {
  value = rancher2_machine_config_v2.all_in_one.kind
}

output "machine_config_name" {
  value = rancher2_machine_config_v2.all_in_one.name
}

output "downstream_security_group" {
  value = local.downstream_security_group_name
}

output "node_id" {
  value = local.node_id
}
