---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster"
sidebar_current: "docs-rancher2-resource-cluster"
description: |-
  Provides a Rancher v2 Cluster resource. This can be used to create Clusters for rancher v2 environments and retrieve their information.
---

# rancher2\_cluster

Provides a Rancher v2 Cluster resource. This can be used to create Clusters for rancher v2 environments and retrieve their information.

## Example Usage

Creating Rancher v2 imported cluster
```hcl
# Create a new rancher2 imported Cluster 
resource "rancher2_cluster" "foo-imported" {
  name = "foo-imported"
  description = "Foo rancher2 imported cluster"
  kind = "imported"
}
```

Creating Rancher v2 rke cluster
```hcl
# Create a new rancher2 rke Cluster 
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  kind = "rke"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}
```

Creating Rancher v2 rke cluster assingning a node pool (overlapped planes)
```hcl
# Create a new rancher2 rke Cluster 
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  kind = "rke"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}
# Create a new rancher2 Node Template
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_config {
    access_key = "AWS_ACCESS_KEY"
    secret_key = "<AWS_SECRET_KEY>"
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = ["<AWS_SECURITY_GROUP>"]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}
# Create a new rancher2 Node Pool
resource "rancher2_node_pool" "foo" {
  cluster_id =  "${rancher2_cluster.foo-custom.id}"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  node_template_id = "${rancher2_node_template.foo.id}"
  quantity = 3
  control_plane = true
  etcd = true
  worker = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cluster (string)
* `kind` - (Required) The kind of the Cluster. `imported`, `eks`, `aks`, `gke` and `rke` are supported
* `rke_config` - (Optional) The rke configuration for `rke` kind Clusters. Conflicts with `aks_config`, `eks_config` and `gke_config`
* `eks_config` - (Optional) The Amazon eks configuration for `eks` kind Clusters
* `aks_config` - (Optional) The Azure aks configuration for `aks` kind Clusters. Conflicts with `eks_config`, `gke_config` and `rke_config`
* `gke_config` - (Optional) The Google gke configuration for `gke` kind Clusters. Conflicts with `aks_config`, `eks_config` and `rke_config`
* `description` - (Optional) The description for Cluster (string)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

### Rancher `rke_config`

The following arguments are supported:

* `addon_job_timeout` - (Optional/Computed) Duration in seconds of addon job (int)
* `addons` - (Optional) Addons descripton to deploy on rke cluster.
* `addons_include` - (Optional) Addons yaml manisfests to deploy on rke cluster (list)
* `authentication` - (Optional/Computed) Kubernetes cluster authentication
* `authorization` - (Optional/Computed) Kubernetes cluster authorization
* `bastion_host` - (Optional/Computed) RKE bastion host
* `cloud_provider` - (Optional/Computed) Kubernetes cluster authentication [rke-cloud-providers](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/)
* `ignore_docker_version` - (Optional/Computed) Ignore docker version (bool)
* `ingress` - (Optional/Computed) Kubernetes ingress configuration
* `kubernetes_version` - (Optional/Computed) Kubernetes version to deploy (string)
* `monitoring` - (Optional/Computed) Kubernetes cluster monitoring
* `network` - (Optional/Computed) Kubernetes cluster networking
* `nodes` - (Optional) RKE cluster nodes (list)
* `prefix_path` - (Optional/Computed) Prefix to customize kubernetes path (string)
* `private_registries` - (Optional) private registries for docker images (list)
* `services` - (Optional/Computed) Kubernetes cluster services
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false`
* `ssh_key_path` - (Optional/Computed) Cluster level SSH private key path (string)

The following arguments are supported for `authentication`:

* `options` - (Optional/Computed) RKE options for authentication (map)
* `sans` - (Optional/Computed) RKE sans for authentication ([]string)
* `strategy` - (Optional/Computed) RKE strategy for authentication (string)

The following arguments are supported for `authorization`:

* `mode` - (Optional) RKE mode for authorization. `rbac` and `none`  are available. Default `rbac`
* `options` - (Optional/Computed) RKE options for authorization (map)

The following arguments are supported for `bastion_host`:

* `user` - (Required) User to connect bastion host (string)
* `address` - (Required) Address ip for the bastion host
* `port` - (Optional) Port for bastion host. Default `22`
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false`
* `ssh_key` - (Optional/Computed) Bastion host SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Bastion host SSH private key path (string)

The following arguments are supported for `cloud_provider`:

* `azure_cloud_provider` - (Optional/Computed) RKE Azure Cloud Provider config for Cloud Provider [rke-azure-cloud-provider](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/azure)
* `custom_cloud_provider` - (Optional/Computed) RKE Custom Cloud Provider config for Cloud Provider (string)
* `name` - (Optional/Computed) RKE sans for Cloud Provider. `aws`, `azure`, `custom`, `openstack`, `vsphere` are supported.
* `openstack_cloud_provider` - (Optional/Computed) RKE Openstack Cloud Provider config for Cloud Provider [rke-openstack-cloud-provider](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/openstack) Extra argument `name` is required on `virtual_center` configuration.
* `vsphere_cloud_provider` - (Optional/Computed) RKE Vsphere Cloud Provider config for Cloud Provider [rke-vsphere-cloud-provider](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/vsphere)

The following arguments are supported for `azure_cloud_provider`:

* `aad_client_id` - (Required) (string)
* `aad_client_secret` - (Required) (string)
* `subscription_id` - (Required) (string)
* `tenant_id` - (Required) (string)
* `aad_client_cert_password` - (Optional/Computed) (string)
* `aad_client_cert_path` - (Optional/Computed) (string)
* `cloud` - (Optional/Computed) (string)
* `cloud_provider_backoff` - (Optional/Computed) (bool)
* `cloud_provider_backoff_duration` - (Optional/Computed) (int)
* `cloud_provider_backoff_exponent` - (Optional/Computed) (int)
* `cloud_provider_backoff_jitter` - (Optional/Computed) (int)
* `cloud_provider_backoff_retries` - (Optional/Computed) (int)
* `cloud_provider_rate_limit` - (Optional/Computed) (bool)
* `cloud_provider_rate_limit_bucket` - (Optional/Computed) (int)
* `cloud_provider_rate_limit_qps` - (Optional/Computed) (int)
* `location` - (Optional/Computed) (string)
* `maximum_load_balancer_rule_count` - (Optional/Computed) (int)
* `primary_availability_set_name"` - (Optional/Computed) (string)
* `primary_scale_set_name` - (Optional/Computed) (string)
* `resource_group` - (Optional/Computed) (string)
* `route_table_name` - (Optional/Computed) (string)
* `security_group_name` - (Optional/Computed) (string)
* `subnet_name` - (Optional/Computed) (string)
* `use_instance_metadata` - (Optional/Computed) (bool)
* `use_managed_identity_extension` - (Optional/Computed) (bool)
* `vm_type` - (Optional/Computed) (string)
* `vnet_name` - (Optional/Computed) (string)
* `vnet_resource_group` - (Optional/Computed) (string)

The following arguments are supported for `openstack_cloud_provider`:

* `global` - (Required)
* `block_storage` - (Optional)
* `load_balancer` - (Optional)
* `metadata` - (Optional)
* `route` - (Optional)

The following arguments are supported for `global`:

* `auth_url` - (Required) (string)
* `password` - (Required) (string)
* `tenant_id` - (Required) (string)
* `user_id` - (Required) (string)
* `username` - (Required) (string)
* `ca_file` - (Optional/Computed) (string)
* `domain_id` - (Optional/Computed) (string)
* `domain_name` - (Optional/Computed) (string)
* `region` - (Optional/Computed) (string)
* `tenant_name` - (Optional/Computed) (string)
* `trust_id` - (Optional/Computed) (string)

The following arguments are supported for `block_storage`:

* `bs_version` - (Optional/Computed) (string)
* `ignore_volume_az` - (Optional/Computed) (string)
* `trust_device_path` - (Optional/Computed) (string)

The following arguments are supported for `load_balancer`:

* `create_monitor` - (Optional/Computed) (bool)
* `floating_network_id` - (Optional/Computed) (string)
* `lb_method` - (Optional/Computed) (string)
* `lb_provider` - (Optional/Computed) (string)
* `lb_version` - (Optional/Computed) (string)
* `manage_security_groups` - (Optional/Computed) (bool)
* `monitor_delay` - (Optional/Computed) Default 60 (int)
* `monitor_max_retries` - (Optional/Computed) Default 5 (int)
* `monitor_timeout` - (Optional/Computed) Default 30 (int)
* `subnet_id` - (Optional/Computed) (string)
* `use_octavia` - (Optional/Computed) (bool)

The following arguments are supported for `metadata`:

* `request_timeout` - (Optional/Computed) (int)
* `search_order` - (Optional/Computed) (string)

The following arguments are supported for `route`:

* `router_id` - (Optional/Computed) (string)

The following arguments are supported for `vsphere_cloud_provider`:

* `virtual_center` - (Required) (List)
* `workspace` - (Required)
* `disk` - (Optional/Computed)
* `global` - (Optional/Computed)
* `network` - (Optional/Computed)

The following arguments are supported for `virtual_center`:

* `datacenters` - (Required) (string)
* `name` - (Required) Name of virtualcenter config for Vsphere Cloud Provider config (string)
* `password` - (Required) (string)
* `user` - (Required) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

The following arguments are supported for `workspace`:

* `datacenter` - (Required) (string)
* `folder` - (Required) (string)
* `server` - (Required) (string)
* `default_datastore` - (Optional/Computed) (string)
* `resourcepool_path` - (Optional/Computed) (string)

The following arguments are supported for `disk`:

* `scsi_controller_type` - (Optional/Computed) (string)

The following arguments are supported for `global`:

* `datacenters` - (Optional/Computed) (string)
* `insecure_flag` - (Optional/Computed) (bool)
* `password` - (Optional/Computed) (string)
* `user` - (Optional/Computed) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

The following arguments are supported for `network`:

* `public_network` - (Optional/Computed) (string)

The following arguments are supported for `ingress`:

* `extra_args` - (Optional/Computed) Extra arguments for RKE Ingress (map)
* `node_selector` - (Optional/Computed) Node selector for RKE Ingress (map)
* `options` - (Optional/Computed) RKE options for Ingress (map)
* `provider` - (Optional/Computed) Provider for RKE Ingress (string)

The following arguments are supported for `monitoring`:

* `options` - (Optional/Computed) RKE options for monitoring (map)
* `provider` - (Optional/Computed) Provider for RKE monitoring (string)

The following arguments are supported for `network`:

* `calico_network_provider` - (Optional/Computed) Calico provider config for RKE network
* `canal_network_provider` - (Optional/Computed) Canal provider config for RKE network
* `flannel_network_provider` - (Optional/Computed) Flannel provider config for RKE network
* `options` - (Optional/Computed) RKE options for network (map)
* `plugin` - (Optional/Computed) Plugin for RKE network. `canal` (default), `flannel` and `calico` are supported.

The following arguments are supported for `calico_network_provider`:

* `cloud_provider` - (Optional/Computed) RKE options for Calico network provider (string)

The following arguments are supported for `canal_network_provider`:

* `iface` - (Optional/Computed) Iface config Canal network provider (string)

The following arguments are supported for `flannel_network_provider`:

* `iface` - (Optional/Computed) Iface config Flannel network provider (string)

The following arguments are supported for `nodes`:

* `address` - (Required) Address ip for node (string)
* `role` - (Requires) Roles for the node. `controlplane`, `etcd` and `worker` are supported. (list)
* `user` - (Required) User to connect node (string)
* `docker_socket` - (Optional/Computed) Docker scojer for node (string)
* `hostname_override` - (Optional) Hostname override for node (string)
* `internal_address` - (Optional) Internal ip for node (string)
* `labels` - (Optional) Labels for the node (map)
* `node_id` - (Optional) Id for the node (string)
* `port` - (Optional) Port for node. Default `22`
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false`
* `ssh_key` - (Optional/Computed) Node SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Node SSH private key path (string)

The following arguments are supported for `private_registries`:

* `is_default` - (Optional) Set as default registry. Default `false`
* `password` - (Optional) Registry password (string)
* `url` - (Required) Registry URL (string)
* `user` - (Optional) Registry user (string)


The following arguments are supported for `services`:

* `etcd` - (Optional/Computed) Etcd options for RKE services
* `kube_api` - (Optional/Computed) Kube API options for RKE services
* `kube_controller` - (Optional/Computed) Kube Controller options for RKE services
* `kubelet` - (Optional/Computed) Kubelet options for RKE services
* `kubeproxy` - (Optional/Computed) Kubeproxy options for RKE services

The following arguments are supported for `etcd`:

* `ca_cert` - (Optional/Computed) Tls CA certificate for etcd service (string)
* `cert` - (Optional/Computed) Tls certificate for etcd service (string)
* `creation` - (Optional/Computed) Creation option for etcd service (string)
* `external_urls` - (Optional) External urls for etcd service (list)
* `extra_args` - (Optional/Computed) Extra arguments for etcd service (map)
* `extra_binds` - (Optional) Extra binds for etcd service (list)
* `extra_env` - (Optional) Extra environment for etcd service (list)
* `image` - (Optional/Computed) Docker image for etcd service (string)
* `key` - (Optional/Computed) Tls key for etcd service (string)
* `path` - (Optional/Computed) Path for etcd service (string)
* `retention` - (Optional/Computed) Retention option for etcd service (string)
* `snapshot` - (Optional/Computed) Snapshot option for etcd service (bool)

The following arguments are supported for `kube_api`:

* `extra_args` - (Optional/Computed) Extra arguments for kube API service (map)
* `extra_binds` - (Optional) Extra binds for kube API service (list)
* `extra_env` - (Optional) Extra environment for kube API service (list)
* `image` - (Optional/Computed) Docker image for kube API service (string)
* `pod_security_policy` - (Optional/Computed) Pod Security Policy option for kube API service (bool)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster IP Range option for kube API service (string)
* `service_node_port_range` - (Optional/Computed) Service Node Port Range option for kube API service (string)

The following arguments are supported for `kube_controller`:

* `cluster_cidr` - (Optional/Computed) Cluster CIDR option for kube controller service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kube controller service (map)
* `extra_binds` - (Optional) Extra binds for kube controller service (list)
* `extra_env` - (Optional) Extra environment for kube controller service (list)
* `image` - (Optional/Computed) Docker image for kube controller service (string)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster ip Range option for kube controller service (string)

The following arguments are supported for `kubelet`:

* `cluster_dns_server` - (Optional/Computed) Cluster DNS Server option for kubelet service (string)
* `cluster_domain` - (Optional/Computed) Cluster Domain option for kubelet service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kubelet service (map)
* `extra_binds` - (Optional) Extra binds for kubelet service (list)
* `extra_env` - (Optional) Extra environment for kubelet service (list)
* `fail_swap_on` - (Optional/Computed) Enable or disable failing when swap on is not supported (bool)
* `image` - (Optional/Computed) Docker image for kubelet service (string)
* `infra_container_image` - (Optional/Computed) Infre container image for kubelet service (string)

The following arguments are supported for `kubeproxy`:

* `extra_args` - (Optional/Computed) Extra arguments for kubeproxy service (map)
* `extra_binds` - (Optional) Extra binds for kubeproxy service (list)
* `extra_env` - (Optional) Extra environment for kubeproxy service (list)
* `image` - (Optional/Computed) Docker image for kubeproxy service (string)

### Amazon `eks_config`

The following arguments are supported:

* `access_key` - (Required) Access key for EKS (string)
* `secret_key` - (Required) Secret key for EKS (string)
* `ami` - (Optional) AMI image for EKS worker nodes (string)
* `associate_worker_node_public_ip` - (Optional) Associate public ip EKS worker nodes (bool). Default `true`
* `instance_type` - (Required) Intance type for EKS cluster (string)
* `maximum_nodes` - (Required) Maximum instaces for EKS cluster (int)
* `minimum_nodes` - (Required) Minimum instaces for EKS cluster (int)
* `region` - (Required) Region for EKS cluster (string)
* `security_groups` - (Required) Security groups for EKS cluster ([]string)
* `service_role` - (Required) Service role for EKS cluster (string)
* `subnets` - (Required) Subnets for EKS cluster ([]string)
* `virtual_network` - (Required) Virtual network for EKS cluster (string)

### Azure `aks_config`

The following arguments are supported:

* `admin_username` - (Required) Admin username for AKS (string)
* `agent_dns_prefix` - (Required) Agent dns prefix for AKS (string)
* `agent_pool_name` - (Required) Agent pool name for AKS cluster (string)
* `agent_vm_size` - (Required) Agent vm size for AKS cluster (string)
* `base_url` - (Required) Base URL for AKS cluster (int)
* `client_id` - (Required) Client ID for AKS (string)
* `client_secret` - (Required) Client secret for AKS (string)
* `count` - (Required) Number of agents for AKS cluster (int)
* `location` - (Required) Location for AKS cluster (string)
* `dns_service_ip` - (Required) DNS service ip for AKS cluster (string)
* `docker_bridge_cidr` - (Required) Docker birdge CIDR for AKS cluster (string)
* `kubernetes_version` - (Required) Kubernetes version for AKS cluster (string)
* `master_dns_prefix` - (Required) Master dns prefix for AKS cluster (string)
* `os_disk_size_gb` - (Required) OS disk size for agents for AKS cluster (int)
* `resource_group` - (Required) Resource group for AKS (string)
* `ssh_public_key_contents` - (Required) SSH public key for AKS cluster (string)
* `service_cidr` - (Required) Services CIDR for AKS cluster (string)
* `subnet` - (Required) Subnet for AKS (string)
* `subscription_id` - (Required) Subscription ID for AKS (string)
* `tag` - (Required) Tags for AKS cluster (map)
* `tenant_id` - (Required) Tenant ID for AKS (string)
* `virtual_network` - (Required) Virtual Network for AKS (string)
* `virtual_network_resource_group` - (Required) Virtual Network resource group for AKS (string)


### Google `gke_config`

The following arguments are supported:

* `cluster_ipv4_cidr` - (Required) Cluster ipv4 CIDR for GKE (string)
* `credential` - (Required) Credential for GKE (string)
* `description` - (Optional) Description for GKE cluster (string)
* `disk_size_gb` - (Required) Disk size for agents for GKE cluster (int)
* `enable_alpha_feature` - (Required) Enable alpha features on GKE cluster (bool)
* `enable_http_load_balancing` - (Required) Enable HTTP load balancing on GKE cluster (bool)
* `enable_horizontal_pod_autoscaling` - (Required) Enable Horitzontal Pod Autoscaling on GKE cluster (bool)
* `enable_kubernetes_dashboard` - (Required) Enable kubernetes dashboard on GKE cluster (bool)
* `enable_legacy_abac` - (Required) Enable legacy abac on GKE cluster (bool)
* `enable_network_policy_config` - (Required) Enable network policy config on GKE cluster (bool)
* `enable_stackdriver_logging` - (Required) Enable stackdriver logging on GKE cluster (bool)
* `enable_stackdriver_monitoring` - (Required) Enable stackdriver monitoring on GKE cluster (bool)
* `image_type` - (Required) Image type for GKE cluster (string)
* `labels` - (Optional/Computed) Labels for GKE cluster (map)
* `locations` - (Required) Locations for GKE cluster ([]string)
* `machine_type` - (Required) Machine type for GKE cluster (string)
* `maintenance_window` - (Required) Maintenance window for GKE cluster (string)
* `master_version` - (Required) Master version for GKE cluster (string)
* `network` - (Required) Network for GKE cluster (string)
* `node_count` - (Required) Node count for GKE cluster (int)
* `node_version` - (Required) Node version for GKE cluster (string)
* `project_id` - (Required) Project ID for GKE cluster (string)
* `sub_network` - (Required) Subnetwork for GKE cluster (string)
* `zone` - (Required) Zone GKE cluster (string)

### Timeouts

`rancher2_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters.
- `update` - (Default `30 minutes`) Used for cluster modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `cluster_registration_token` - (Computed) Cluster Registration Token generated for the cluster.
* `kube_config` - (Computed) Kube Config generated for the cluster.

The following attributes are exported for `cluster_registration_token`:

* `cluster_id` - (Computed) Cluster ID
* `name` - (Computed) Name of cluster registration token
* `command` - (Computed) Command to execute in a imported k8s cluster
* `insecure_command` - (Computed) Insecure command to execute in a imported k8s cluster
* `manifest_url` - (Computed) K8s mnifest url to execute kubectl in a imported k8s cluster
* `node_command` - (Computed) Node command to execute in linux nodes for custom k8s cluster
* `windows_node_command` - (Computed) Node command to execute in windows nodes for custom k8s cluster
* `annotations` - (Computed) Annotations for cluster registration token object
* `labels` - (Computed) Labels for cluster registration token object

## Import

Clusters can be imported using the rancher Cluster ID

```
$ terraform import rancher2_cluster.foo <cluster>
```

