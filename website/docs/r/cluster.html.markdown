---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster"
sidebar_current: "docs-rancher2-resource-cluster"
description: |-
  Provides a Rancher v2 Cluster resource. This can be used to create Clusters for Rancher v2 environments and retrieve their information.
---

# rancher2\_cluster

Provides a Rancher v2 Cluster resource. This can be used to create Clusters for Rancher v2 environments and retrieve their information.

## Example Usage

Creating Rancher v2 imported cluster

```hcl
# Create a new rancher2 imported Cluster
resource "rancher2_cluster" "foo-imported" {
  name = "foo-imported"
  description = "Foo rancher2 imported cluster"
}
```

Creating Rancher v2 RKE cluster

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}
```

Creating Rancher v2 RKE cluster enabling and customizing monitoring

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
  }
  enable_cluster_monitoring = true
  cluster_monitoring_input {
    answers = {
      "exporter-kubelets.https" = true
      "exporter-node.enabled" = true
      "exporter-node.ports.metrics.port" = 9796
      "exporter-node.resources.limits.cpu" = "200m"
      "exporter-node.resources.limits.memory" = "200Mi"
      "grafana.persistence.enabled" = false
      "grafana.persistence.size" = "10Gi"
      "grafana.persistence.storageClass" = "default"
      "operator.resources.limits.memory" = "500Mi"
      "prometheus.persistence.enabled" = "false"
      "prometheus.persistence.size" = "50Gi"
      "prometheus.persistence.storageClass" = "default"
      "prometheus.persistent.useReleaseName" = "true"
      "prometheus.resources.core.limits.cpu" = "1000m",
      "prometheus.resources.core.limits.memory" = "1500Mi"
      "prometheus.resources.core.requests.cpu" = "750m"
      "prometheus.resources.core.requests.memory" = "750Mi"
      "prometheus.retention" = "12h"
    }
  }
}
```


Creating Rancher v2 RKE cluster assigning a node pool (overlapped planes)

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
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

Creating Rancher v2 RKE cluster from template. For Rancher v2.3.x or above.

```hcl
# Create a new rancher2 cluster template
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "owner"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V1"
    cluster_config {
      rke_config {
        network {
          plugin = "canal"
        }
        services {
          etcd {
            creation = "6h"
            retention = "24h"
          }
        }
      }
    }
    default = true
  }
  description = "Test cluster template v2"
}
# Create a new rancher2 RKE Cluster from template
resource "rancher2_cluster" "foo" {
  name = "foo"
  cluster_template_id = "${rancher2_cluster_template.foo.id}"
  cluster_template_revision_id = "${rancher2_cluster_template.foo.default_revision_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cluster (string)
* `rke_config` - (Optional/Computed) The RKE configuration for `rke` Clusters. Conflicts with `aks_config`, `eks_config` and `gke_config` (list maxitems:1)
* `aks_config` - (Optional) The Azure AKS configuration for `aks` Clusters. Conflicts with `eks_config`, `gke_config` and `rke_config` (list maxitems:1)
* `eks_config` - (Optional) The Amazon EKS configuration for `eks` Clusters. Conflicts with `aks_config`, `gke_config` and `rke_config` (list maxitems:1)
* `gke_config` - (Optional) The Google GKE configuration for `gke` Clusters. Conflicts with `aks_config`, `eks_config` and `rke_config` (list maxitems:1)
* `description` - (Optional) The description for Cluster (string)
* `cluster_auth_endpoint` - (Optional/Computed) Enabling the [local cluster authorized endpoint](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#local-cluster-auth-endpoint) allows direct communication with the cluster, bypassing the Rancher API proxy. (list maxitems:1)
* `cluster_monitoring_input` - (Optional/Computed) Cluster monitoring config. Any parameter defined in [rancher-monitoring charts](https://github.com/rancher/system-charts/tree/dev/charts/rancher-monitoring) could be configured  (list maxitems:1)
* `cluster_template_answers` - (Optional) Cluster template answers. Just for Rancher v2.3.x and above (list maxitems:1)
* `cluster_template_id` - (Optional) Cluster template ID. Just for Rancher v2.3.x and above (string)
* `cluster_template_questions` - (Optional) Cluster template questions. Just for Rancher v2.3.x and above (list)
* `cluster_template_revision_id` - (Optional) Cluster template revision ID. Just for Rancher v2.3.x and above (string)
* `default_pod_security_policy_template_id` - (Optional/Computed) [Default pod security policy template id](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#pod-security-policy-support) (string)
* `desired_agent_image` - (Optional/Computed) Desired agent image. Just for Rancher v2.3.x and above (string)
* `desired_auth_image` - (Optional/Computed) Desired auth image. Just for Rancher v2.3.x and above (string)
* `docker_root_dir` - (Optional/Computed) Desired auth image. Just for Rancher v2.3.x and above (string)
* `enable_cluster_alerting` - (Optional) Enable built-in cluster alerting. Default `false` (bool)
* `enable_cluster_monitoring` - (Optional) Enable built-in cluster monitoring. Default `false` (bool)
* `enable_cluster_istio` - (Optional) Enable built-in cluster istio. Default `false`. Just for Rancher v2.3.x and above (bool)
* `enable_network_policy` - (Optional) Enable project network isolation. Default `false` (bool)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cluster_registration_token` - (Computed) Cluster Registration Token generated for the cluster (list maxitems:1)
* `default_project_id` - (Computed) Default project ID for the cluster (string)
* `driver` - (Computed) The driver used for the Cluster. `imported`, `azurekubernetesservice`, `amazonelasticcontainerservice`, `googlekubernetesengine` and `rancherKubernetesEngine` are supported (string)
* `kube_config` - (Computed) Kube Config generated for the cluster (string)
* `system_project_id` - (Computed) System project ID for the cluster (string)

## Nested blocks

### `rke_config`

#### Arguments

* `addon_job_timeout` - (Optional/Computed) Duration in seconds of addon job (int)
* `addons` - (Optional) Addons descripton to deploy on RKE cluster.
* `addons_include` - (Optional) Addons yaml manifests to deploy on RKE cluster (list)
* `authentication` - (Optional/Computed) Kubernetes cluster authentication (list maxitems:1)
* `authorization` - (Optional/Computed) Kubernetes cluster authorization (list maxitems:1)
* `bastion_host` - (Optional/Computed) RKE bastion host (list maxitems:1)
* `cloud_provider` - (Optional/Computed) RKE cloud provider [rke-cloud-providers](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/) (list maxitems:1)
* `dns` - (Optional/Computed) RKE dns add-on. Just for Rancher v2.2.x (list maxitems:1)
* `ignore_docker_version` - (Optional) Ignore docker version. Default `true` (bool)
* `ingress` - (Optional/Computed) Kubernetes ingress configuration (list maxitems:1)
* `kubernetes_version` - (Optional/Computed) Kubernetes version to deploy (string)
* `monitoring` - (Optional/Computed) Kubernetes cluster monitoring (list maxitems:1)
* `network` - (Optional/Computed) Kubernetes cluster networking (list maxitems:1)
* `nodes` - (Optional) RKE cluster nodes (list)
* `prefix_path` - (Optional/Computed) Prefix to customize Kubernetes path (string)
* `private_registries` - (Optional) private registries for docker images (list)
* `services` - (Optional/Computed) Kubernetes cluster services (list maxitems:1)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false`
* `ssh_cert_path` - (Optional/Computed) Cluster level SSH certificate path (string)
* `ssh_key_path` - (Optional/Computed) Cluster level SSH private key path (string)

#### `authentication`

##### Arguments

* `sans` - (Optional/Computed) RKE sans for authentication ([]string)
* `strategy` - (Optional/Computed) RKE strategy for authentication (string)

#### `authorization`

##### Arguments

* `mode` - (Optional) RKE mode for authorization. `rbac` and `none` modes are available. Default `rbac` (string)
* `options` - (Optional/Computed) RKE options for authorization (map)

#### `bastion_host`

##### Arguments

* `address` - (Required) Address ip for the bastion host (string)
* `user` - (Required) User to connect bastion host (string)
* `port` - (Optional) Port for bastion host. Default `22` (string)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false` (bool)
* `ssh_key` - (Optional/Computed/Sensitive) Bastion host SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Bastion host SSH private key path (string)

#### `cloud_provider`

##### Arguments

* `aws_cloud_provider` - (Optional/Computed) RKE AWS Cloud Provider config for Cloud Provider [rke-aws-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/aws/) (list maxitems:1)
* `azure_cloud_provider` - (Optional/Computed) RKE Azure Cloud Provider config for Cloud Provider [rke-azure-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/azure/) (list maxitems:1)
* `custom_cloud_provider` - (Optional/Computed) RKE Custom Cloud Provider config for Cloud Provider (string) (string)
* `name` - (Optional/Computed) RKE sans for Cloud Provider. `aws`, `azure`, `custom`, `openstack`, `vsphere` are supported. (string)
* `openstack_cloud_provider` - (Optional/Computed) RKE Openstack Cloud Provider config for Cloud Provider [rke-openstack-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/openstack/) (list maxitems:1)
* `vsphere_cloud_provider` - (Optional/Computed) RKE Vsphere Cloud Provider config for Cloud Provider [rke-vsphere-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/vsphere/) Extra argument `name` is required on `virtual_center` configuration. (list maxitems:1)

##### `aws_cloud_provider`

###### Arguments

* `global` - (Optional) (list maxitems:1)
* `service_override` - (Optional) (list)

###### `global`

###### Arguments

* `disable_security_group_ingress` - (Optional) Default `false` (bool)
* `disable_strict_zone_check` - (Optional) Default `false` (bool)
* `elb_security_group` - (Optional/Computed) (string)
* `kubernetes_cluster_id` - (Optional/Computed) (string)
* `kubernetes_cluster_tag` - (Optional/Computed) (string)
* `role_arn` - (Optional/Computed) (string)
* `route_table_id` - (Optional/Computed/Sensitive) (string)
* `subnet_id` - (Optional/Computed) (string)
* `vpc` - (Optional/Computed) (string)
* `zone` - (Optional/Computed) (string)

###### `service_override`

###### Arguments

* `service` - (Required) (string)
* `region` - (Optional/Computed) (string)
* `signing_method` - (Optional/Computed) (string)
* `signing_name` - (Optional/Computed) (string)
* `signing_region` - (Optional/Computed) (string)
* `url` - (Optional/Computed) (string)

##### `azure_cloud_provider`

###### Arguments

* `aad_client_id` - (Required/Sensitive) (string)
* `aad_client_secret` - (Required/Sensitive) (string)
* `subscription_id` - (Required/Sensitive) (string)
* `tenant_id` - (Required/Sensitive) (string)
* `aad_client_cert_password` - (Optional/Computed/Sensitive) (string)
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
* `primary_availability_set_name` - (Optional/Computed) (string)
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

##### `openstack_cloud_provider`

###### Arguments

* `global` - (Required) (list maxitems:1)
* `block_storage` - (Optional/Computed) (list maxitems:1)
* `load_balancer` - (Optional/Computed) (list maxitems:1)
* `metadata` - (Optional/Computed) (list maxitems:1)
* `route` - (Optional/Computed) (list maxitems:1)

###### `global`

###### Arguments

* `auth_url` - (Required) (string)
* `password` - (Required/Sensitive) (string)
* `username` - (Required/Sensitive) (string)
* `ca_file` - (Optional/Computed) (string)
* `domain_id` - (Optional/Computed/Sensitive) Required if `domain_name` not provided. (string)
* `domain_name` - (Optional/Computed) Required if `domain_id` not provided. (string)
* `region` - (Optional/Computed) (string)
* `tenant_id` - (Optional/Computed/Sensitive) Required if `tenant_name` not provided. (string)
* `tenant_name` - (Optional/Computed) Required if `tenant_id` not provided. (string)
* `trust_id` - (Optional/Computed/Sensitive) (string)

###### `block_storage`

###### Arguments

* `bs_version` - (Optional/Computed) (string)
* `ignore_volume_az` - (Optional/Computed) (string)
* `trust_device_path` - (Optional/Computed) (string)

###### `load_balancer`

###### Arguments

* `create_monitor` - (Optional/Computed) (bool)
* `floating_network_id` - (Optional/Computed) (string)
* `lb_method` - (Optional/Computed) (string)
* `lb_provider` - (Optional/Computed) (string)
* `lb_version` - (Optional/Computed) (string)
* `manage_security_groups` - (Optional/Computed) (bool)
* `monitor_delay` - (Optional/Computed) Default `60s` (string)
* `monitor_max_retries` - (Optional/Computed) Default 5 (int)
* `monitor_timeout` - (Optional/Computed) Default `30s` (string)
* `subnet_id` - (Optional/Computed) (string)
* `use_octavia` - (Optional/Computed) (bool)

###### `metadata`

###### Arguments

* `request_timeout` - (Optional/Computed) (int)
* `search_order` - (Optional/Computed) (string)

###### `route`

###### Arguments

* `router_id` - (Optional/Computed) (string)

##### `vsphere_cloud_provider`

###### Arguments

* `virtual_center` - (Required) (List)
* `workspace` - (Required) (list maxitems:1)
* `disk` - (Optional/Computed) (list maxitems:1)
* `global` - (Optional/Computed) (list maxitems:1)
* `network` - (Optional/Computed) (list maxitems:1)

###### `virtual_center`

###### Arguments

* `datacenters` - (Required) (string)
* `name` - (Required) Name of virtualcenter config for Vsphere Cloud Provider config (string)
* `password` - (Required/Sensitive) (string)
* `user` - (Required/Sensitive) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

###### `workspace`

###### Arguments

* `datacenter` - (Required) (string)
* `folder` - (Required) (string)
* `server` - (Required) (string)
* `default_datastore` - (Optional/Computed) (string)
* `resourcepool_path` - (Optional/Computed) (string)

###### `disk`

###### Arguments

* `scsi_controller_type` - (Optional/Computed) (string)

###### `global`

###### Arguments

* `datacenters` - (Optional/Computed) (string)
* `insecure_flag` - (Optional/Computed) (bool)
* `password` - (Optional/Computed) (string)
* `user` - (Optional/Computed) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

###### `network`

###### Arguments

* `public_network` - (Optional/Computed) (string)

#### `dns`

##### Arguments

* `node_selector` - (Optional/Computed) DNS add-on node selector (map)
* `provider` - (Optional) DNS add-on provider. `kube-dns`, `coredns` (default), and `none` are supported (string)
* `reverse_cidrs` - (Optional/Computed) DNS add-on reverse cidr  (list)
* `upstream_nameservers` - (Optional/Computed) DNS add-on upstream nameservers  (list)

#### `ingress`

##### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for RKE Ingress (map)
* `node_selector` - (Optional/Computed) Node selector for RKE Ingress (map)
* `options` - (Optional/Computed) RKE options for Ingress (map)
* `provider` - (Optional/Computed) Provider for RKE Ingress (string)

#### `monitoring`

##### Arguments

* `options` - (Optional/Computed) RKE options for monitoring (map)
* `provider` - (Optional/Computed) Provider for RKE monitoring (string)

#### `network`

##### Arguments

* `calico_network_provider` - (Optional/Computed) Calico provider config for RKE network (list maxitems:1)
* `canal_network_provider` - (Optional/Computed) Canal provider config for RKE network (list maxitems:1)
* `flannel_network_provider` - (Optional/Computed) Flannel provider config for RKE network (list maxitems:1)
* `weave_network_provider` - (Optional/Computed) Weave provider config for RKE network (list maxitems:1)
* `options` - (Optional/Computed) RKE options for network (map)
* `plugin` - (Optional/Computed) Plugin for RKE network. `canal` (default), `flannel`, `calico` and `weave` are supported. (string)

##### `calico_network_provider`

###### Arguments

* `cloud_provider` - (Optional/Computed) RKE options for Calico network provider (string)

##### `canal_network_provider`

###### Arguments

* `iface` - (Optional/Computed) Iface config Canal network provider (string)

##### `flannel_network_provider`

###### Arguments

* `iface` - (Optional/Computed) Iface config Flannel network provider (string)

##### `weave_network_provider`

###### Arguments

* `password` - (Optional/Computed) Password config Weave network provider (string)

#### `nodes`

##### Arguments

* `address` - (Required) Address ip for node (string)
* `role` - (Requires) Roles for the node. `controlplane`, `etcd` and `worker` are supported. (list)
* `user` - (Required/Sensitive) User to connect node (string)
* `docker_socket` - (Optional/Computed) Docker socket for node (string)
* `hostname_override` - (Optional) Hostname override for node (string)
* `internal_address` - (Optional) Internal ip for node (string)
* `labels` - (Optional) Labels for the node (map)
* `node_id` - (Optional) Id for the node (string)
* `port` - (Optional) Port for node. Default `22` (string)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false` (bool)
* `ssh_key` - (Optional/Computed/Sensitive) Node SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Node SSH private key path (string)

#### `private_registries`

##### Arguments

* `url` - (Required) Registry URL (string)
* `is_default` - (Optional) Set as default registry. Default `false` (bool)
* `password` - (Optional/Sensitive) Registry password (string)
* `user` - (Optional/Sensitive) Registry user (string)


#### `services`

##### Arguments

* `etcd` - (Optional/Computed) Etcd options for RKE services (list maxitems:1)
* `kube_api` - (Optional/Computed) Kube API options for RKE services (list maxitems:1)
* `kube_controller` - (Optional/Computed) Kube Controller options for RKE services (list maxitems:1)
* `kubelet` - (Optional/Computed) Kubelet options for RKE services (list maxitems:1)
* `kubeproxy` - (Optional/Computed) Kubeproxy options for RKE services (list maxitems:1)
* `scheduler` - (Optional/Computed) Scheduler options for RKE services (list maxitems:1)

##### `etcd`

###### Arguments

* `backup_config` - (Optional/Computed) Backup options for etcd service. Just for Rancher v2.2.x (list maxitems:1)
* `ca_cert` - (Optional/Computed) TLS CA certificate for etcd service (string)
* `cert` - (Optional/Computed/Sensitive) TLS certificate for etcd service (string)
* `creation` - (Optional/Computed) Creation option for etcd service (string)
* `external_urls` - (Optional) External urls for etcd service (list)
* `extra_args` - (Optional/Computed) Extra arguments for etcd service (map)
* `extra_binds` - (Optional) Extra binds for etcd service (list)
* `extra_env` - (Optional) Extra environment for etcd service (list)
* `gid` - (Optional) Etcd service GID. Default: `0`. For Rancher v2.3.x or above (int)
* `image` - (Optional/Computed) Docker image for etcd service (string)
* `key` - (Optional/Computed/Sensitive) TLS key for etcd service (string)
* `path` - (Optional/Computed) Path for etcd service (string)
* `retention` - (Optional/Computed) Retention option for etcd service (string)
* `snapshot` - (Optional/Computed) Snapshot option for etcd service (bool)
* `uid` - (Optional) Etcd service UID. Default: `0`. For Rancher v2.3.x or above (int)

###### `backup_config`

###### Arguments

* `enabled` - (Optional) Enable etcd backup (bool)
* `interval_hours` - (Optional) Interval hours for etcd backup. Default `12` (int)
* `retention` - (Optional) Retention for etcd backup. Default `6` (int)
* `s3_backup_config` - (Optional) S3 config options for etcd backup (list maxitems:1)
* `safe_timestamp` - (Optional) Safe timestamp for etcd backup. Default: `false` (bool)

###### `s3_backup_config`

###### Arguments

* `access_key` - (Optional/Sensitive) Access key for S3 service (string)
* `bucket_name` - (Required) Bucket name for S3 service (string)
* `custom_ca` - (Optional) Base64 encoded custom CA for S3 service. Use filebase64(<FILE>) for encoding file. Available from Rancher v2.2.5 (string)
* `endpoint` - (Required) Endpoint for S3 service (string)
* `folder` - (Optional) Folder for S3 service. Available from Rancher v2.2.7 (string)
* `region` - (Optional) Region for S3 service (string)
* `secret_key` - (Optional/Sensitive) Secret key for S3 service (string)

##### `kube_api`

###### Arguments

* `admission_configuration` - (Optional) Admission configuration (map)
* `always_pull_images` - (Optional) Enable [AlwaysPullImages](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#alwayspullimages) Admission controller plugin. [Rancher docs](https://rancher.com/docs/rke/latest/en/config-options/services/#kubernetes-api-server-options) Default: `false` (bool)
* `audit_log` - (Optional) K8s audit log configuration. (list maxitem: 1)
* `event_rate_limit` - (Optional) K8s event rate limit configuration. (list maxitem: 1)
* `extra_args` - (Optional/Computed) Extra arguments for kube API service (map)
* `extra_binds` - (Optional) Extra binds for kube API service (list)
* `extra_env` - (Optional) Extra environment for kube API service (list)
* `image` - (Optional/Computed) Docker image for kube API service (string)
* `pod_security_policy` - (Optional) Pod Security Policy option for kube API service. Default `false` (bool)
* `secrets_encryption_config` - (Optional) [Encrypt k8s secret data configration](https://rancher.com/docs/rke/latest/en/config-options/secrets-encryption/). (list maxitem: 1)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster IP Range option for kube API service (string)
* `service_node_port_range` - (Optional/Computed) Service Node Port Range option for kube API service (string)

###### `audit_log`

###### Arguments

* `configuration` - (Optional) Audit log configuration. (list maxtiem: 1)
* `enabled` - (Optional) Enable audit log. Default: `false` (bool)

###### `configuration`

###### Arguments

* `format` - (Optional) Audit log format. Default: 'json' (string)
* `max_age` - (Optional) Audit log max age. Default: `30` (int)
* `max_backup` - (Optional) Audit log max backup. Default: `10` (int)
* `max_size` - (Optional) Audit log max size. Default: `100` (int)
* `path` - (Optional) (Optional) Audit log path. Default: `/var/log/kube-audit/audit-log.json` (string)
* `policy` - (Optional) Audit log policy (map)

###### `event_rate_limit`

###### Arguments

* `configuration` - (Optional) Event rate limit configuration. (map)
* `enabled` - (Optional) Enable event rate limit. Default: `false` (bool)

###### `secrets_encryption_config`

###### Arguments

* `custom_config` - (Optional) Secrets encryption configuration. (map)
* `enabled` - (Optional) Enable secrets encryption. Default: `false` (bool)

##### `kube_controller`

###### Arguments

* `cluster_cidr` - (Optional/Computed) Cluster CIDR option for kube controller service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kube controller service (map)
* `extra_binds` - (Optional) Extra binds for kube controller service (list)
* `extra_env` - (Optional) Extra environment for kube controller service (list)
* `image` - (Optional/Computed) Docker image for kube controller service (string)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster ip Range option for kube controller service (string)

##### `kubelet`

###### Arguments

* `cluster_dns_server` - (Optional/Computed) Cluster DNS Server option for kubelet service (string)
* `cluster_domain` - (Optional/Computed) Cluster Domain option for kubelet service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kubelet service (map)
* `extra_binds` - (Optional) Extra binds for kubelet service (list)
* `extra_env` - (Optional) Extra environment for kubelet service (list)
* `fail_swap_on` - (Optional/Computed) Enable or disable failing when swap on is not supported (bool)
* `generate_serving_certificate` [Generate a certificate signed by the kube-ca](https://rancher.com/docs/rke/latest/en/config-options/services/#kubelet-serving-certificate-requirements). Default `false` (bool)
* `image` - (Optional/Computed) Docker image for kubelet service (string)
* `infra_container_image` - (Optional/Computed) Infra container image for kubelet service (string)

##### `kubeproxy`

###### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for kubeproxy service (map)
* `extra_binds` - (Optional) Extra binds for kubeproxy service (list)
* `extra_env` - (Optional) Extra environment for kubeproxy service (list)
* `image` - (Optional/Computed) Docker image for kubeproxy service (string)

##### `scheduler`

###### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for scheduler service (map)
* `extra_binds` - (Optional) Extra binds for scheduler service (list)
* `extra_env` - (Optional) Extra environment for scheduler service (list)
* `image` - (Optional/Computed) Docker image for scheduler service (string)

### `aks_config`

#### Arguments

The following arguments are supported:

* `agent_dns_prefix` - (Required) DNS prefix to be used to create the FQDN for the agent pool (string)
* `client_id` - (Required/Sensitive) Azure client ID to use (string)
* `client_secret` - (Required/Sensitive) Azure client secret associated with the \"client id\" (string)
* `kubernetes_version` - (Required) Specify the version of Kubernetes. To check available versions exec `az aks get-versions -l eastus -o table` (string)
* `master_dns_prefix` - (Required) DNS prefix to use the Kubernetes cluster control pane (string)
* `resource_group` - (Required) The name of the Cluster resource group (string)
* `ssh_public_key_contents` - (Required) Contents of the SSH public key used to authenticate with Linux hosts (string)
* `subnet` - (Required) The name of an existing Azure Virtual Subnet. Composite of agent virtual network subnet ID (string)
* `subscription_id` - (Required) Subscription credentials which uniquely identify Microsoft Azure subscription (string)
* `tenant_id` - (Required) Azure tenant ID to use (string)
* `virtual_network` - (Required) The name of an existing Azure Virtual Network. Composite of agent virtual network subnet ID (string)
* `virtual_network_resource_group` - (Required) The resource group of an existing Azure Virtual Network. Composite of agent virtual network subnet ID (string)
* `add_client_app_id` - (Optional/Sensitive) The ID of an Azure Active Directory client application of type \"Native\". This application is for user login via kubectl (string)
* `add_server_app_id` - (Optional/Sensitive) The ID of an Azure Active Directory server application of type \"Web app/API\". This application represents the managed cluster's apiserver (Server application) (string)
* `aad_server_app_secret` - (Optional/Sensitive) The secret of an Azure Active Directory server application (string)
* `aad_tenant_id` - (Optional/Sensitive) The ID of an Azure Active Directory tenant (string)
* `admin_username` - (Optional) The administrator username to use for Linux hosts. Default `azureuser` (string)
* `agent_os_disk_size` - (Optional) GB size to be used to specify the disk for every machine in the agent pool. If you specify 0, it will apply the default according to the \"agent vm size\" specified. Default `0` (int)
* `agent_pool_name` - (Optional) Name for the agent pool, upto 12 alphanumeric characters. Default `agentpool0` (string)
* `agent_storage_profile` - (Optional) Storage profile specifies what kind of storage used on machine in the agent pool. Chooses from [ManagedDisks StorageAccount]. Default `ManagedDisks` (string)
* `agent_vm_size` - (Optional) Size of machine in the agent pool. Default `Standard_D1_v2` (string)
* `auth_base_url` - (Optional) Different authentication API url to use. Default `https://login.microsoftonline.com/` (string)
* `base_url` - (Optional) Different resource management API url to use. Default `https://management.azure.com/` (string)
* `count` - (Optional) Number of machines (VMs) in the agent pool. Allowed values must be in the range of 1 to 100 (inclusive). Default `1` (int)
* `dns_service_ip` - (Optional) An IP address assigned to the Kubernetes DNS service. It must be within the Kubernetes Service address range specified in \"service cidr\". Default `10.0.0.10` (string)
* `docker_bridge_cidr` - (Required) A CIDR notation IP range assigned to the Docker bridge network. It must not overlap with any Subnet IP ranges or the Kubernetes Service address range specified in \"service cidr\". Default `172.17.0.1/16` (string)
* `enable_http_application_routing` - (Optional) Enable the Kubernetes ingress with automatic public DNS name creation. Default `false` (bool)
* `enable_monitoring` - (Optional) Turn on Azure Log Analytics monitoring. Uses the Log Analytics \"Default\" workspace if it exists, else creates one. if using an existing workspace, specifies \"log analytics workspace resource id\". Default `true` (bool)
* `location` - (Optional) Azure Kubernetes cluster location. Default `eastus` (string)
* `log_analytics_workspace` - (Optional) The name of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses '{resource group}-{subscription id}-{location code}' (string)
* `log_analytics_workspace_resource_group` - (Optional) The resource group of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses the 'Cluster' resource group (string)
* `max_pods` - (Optional) Maximum number of pods that can run on a node. Default `110` (int)
* `network_plugin` - (Optional) Network plugin used for building Kubernetes network. Chooses from `azure` or `kubenet`. Default `azure` (string)
* `network_policy` - (Optional) Network policy used for building Kubernetes network. Chooses from `calico` (string)
* `pod_cidr` - (Optional) A CIDR notation IP range from which to assign Kubernetes Pod IPs when \"network plugin\" is specified in \"kubenet\". Default `172.244.0.0/16` (string)
* `service_cidr` - (Optional) A CIDR notation IP range from which to assign Kubernetes Service cluster IPs. It must not overlap with any Subnet IP ranges. Default `10.0.0.0/16` (string)
* `tag` - (Optional/Computed) Tags for Kubernetes cluster. For example, foo=bar (map)

### `eks_config`

#### Arguments

The following arguments are supported:

* `access_key` - (Required/Sensitive) The AWS Client ID to use (string)
* `kubernetes_version` - (Required) The Kubernetes master version (string)
* `secret_key` - (Required/Sensitive) The AWS Client Secret associated with the Client ID (string)
* `ami` - (Optional) AMI ID to use for the worker nodes instead of the default (string)
* `associate_worker_node_public_ip` - (Optional) Associate public ip EKS worker nodes. Default `true` (bool)
* `desired_nodes` - (Optional) The desired number of worker nodes. Just for Rancher v2.3.x and above. Default `3` (int)
* `instance_type` - (Optional) The type of machine to use for worker nodes. Default `t2.medium` (string)
* `key_pair_name` - (Optional) Allow user to specify key name to use. Just for Rancher v2.2.7 and above (string)
* `maximum_nodes` - (Optional) The maximum number of worker nodes. Default `3` (int)
* `minimum_nodes` - (Optional) The minimum number of worker nodes. Default `1` (int)
* `node_volume_size` - (Optional) The volume size for each node. Default `20` (int)
* `region` - (Optional) The AWS Region to create the EKS cluster in. Default `us-west-2` (string)
* `security_groups` - (Optional) List of security groups to use for the cluster. If it's not specified Rancher will create a new security group (list)
* `service_role` - (Optional) The service role to use to perform the cluster operations in AWS. If it's not specified Rancher will create a new service role (string)
* `session_token` - (Optional/Sensitive) A session token to use with the client key and secret if applicable (string)
* `subnets` - (Optional) List of subnets in the virtual network to use. If it's not specified Rancher will create 3 news subnets (list)
* `user_data` - (Optional/Computed) Pass user-data to the nodes to perform automated configuration tasks (string)
* `virtual_network` - (Optional) The name of the virtual network to use. If it's not specified Rancher will create a new VPC (string)

### `gke_config`

#### Arguments

The following arguments are supported:

* `cluster_ipv4_cidr` - (Required) The IP address range of the container pods (string)
* `credential` - (Required/Sensitive) The contents of the GC credential file (string)
* `disk_type` - (Required) Type of the disk attached to each node (string)
* `image_type` - (Required) The image to use for the worker nodes (string)
* `ip_policy_cluster_ipv4_cidr_block` - (Required) The IP address range for the cluster pod IPs (string)
* `ip_policy_cluster_secondary_range_name` - (Required) The name of the secondary range to be used for the cluster CIDR block (string)
* `ip_policy_node_ipv4_cidr_block` - (Required) The IP address range of the instance IPs in this cluster (string)
* `ip_policy_services_ipv4_cidr_block` - (Required) The IP address range of the services IPs in this cluster (string)
* `ip_policy_services_secondary_range_name` - (Required) The name of the secondary range to be used for the services CIDR block (string)
* `ip_policy_subnetwork_name` - (Required) A custom subnetwork name to be used if createSubnetwork is true (string)
* `locations` - (Required) Locations for GKE cluster (list)
* `machine_type` - (Required) Machine type for GKE cluster (string)
* `maintenance_window` - (Required) Maintenance window for GKE cluster (string)
* `master_ipv4_cidr_block` - (Required) The IP range in CIDR notation to use for the hosted master network (string)
* `master_version` - (Required) Master version for GKE cluster (string)
* `network` - (Required) Network for GKE cluster (string)
* `node_pool` - (Required) The ID of the cluster node pool (string)
* `node_version` - (Required) Node version for GKE cluster (string)
* `oauth_scopes` - (Required) The set of Google API scopes to be made available on all of the node VMs under the default service account (list)
* `project_id` - (Required) Project ID for GKE cluster (string)
* `service_account` - (Required) The Google Cloud Platform Service Account to be used by the node VMs (string)
* `sub_network` - (Required) Subnetwork for GKE cluster (string)
* `description` - (Optional) An optional description of this cluster (string)
* `disk_size_gb` - (Optional) Size of the disk attached to each node. Default `100` (int)
* `enable_alpha_feature` - (Optional) To enable Kubernetes alpha feature. Default `true` (bool)
* `enable_auto_repair` - (Optional) Specifies whether the node auto-repair is enabled for the node pool. Default `false` (bool)
* `enable_auto_upgrade` - (Optional) Specifies whether node auto-upgrade is enabled for the node pool. Default `false` (bool)
* `enable_horizontal_pod_autoscaling` - (Optional) Enable horizontal pod autoscaling for the cluster. Default `true` (bool)
* `enable_http_load_balancing` - (Optional) Enable HTTP load balancing on GKE cluster. Default `true` (bool)
* `enable_kubernetes_dashboard` - (Optional) Whether to enable the Kubernetes dashboard. Default `false` (bool)
* `enable_legacy_abac` - (Optional) Whether to enable legacy abac on the cluster. Default `false` (bool)
* `enable_network_policy_config` - (Optional) Enable stackdriver logging. Default `true` (bool)
* `enable_nodepool_autoscaling` - (Optional) Enable nodepool autoscaling. Default `false` (bool)
* `enable_private_endpoint` - (Optional) Whether the master's internal IP address is used as the cluster endpoint. Default `false` (bool)
* `enable_private_nodes` - (Optional) Whether nodes have internal IP address only. Default `false` (bool)
* `enable_stackdriver_logging` - (Optional) Enable stackdriver monitoring. Default `true` (bool)
* `enable_stackdriver_monitoring` - (Optional) Enable stackdriver monitoring on GKE cluster (bool)
* `ip_policy_create_subnetwork` - (Optional) Whether a new subnetwork will be created automatically for the cluster. Default `false` (bool)
* `issue_client_certificate` - (Optional) Issue a client certificate. Default `false` (bool)
* `kubernetes_dashboard` - (Optional) Enable the Kubernetes dashboard. Default `false` (bool)
* `labels` - (Optional/Computed) The map of Kubernetes labels to be applied to each node (map)
* `local_ssd_count` - (Optional) The number of local SSD disks to be attached to the node. Default `0` (int)
* `master_authorized_network_cidr_blocks` - (Optional) Define up to 10 external networks that could access Kubernetes master through HTTPS (list)
* `max_node_count` - (Optional) Maximum number of nodes in the NodePool. Must be >= minNodeCount. There has to enough quota to scale up the cluster. Default `0` (int)
* `min_node_count` - (Optional) Minimmum number of nodes in the NodePool. Must be >= 1 and <= maxNodeCount. Default `0` (int)
* `node_count` - (Optional) Node count for GKE cluster. Default `3` (int)
* `preemptible` - (Optional) Whether the nodes are created as preemptible VM instances. Default `false` (bool)
* `resource_labels` - (Optional/Computed) The map of Kubernetes labels to be applied to each cluster (map)
* `use_ip_aliases` - (Optional) Whether alias IPs will be used for pod IPs in the cluster. Default `false` (bool)
* `taints` - (Required) List of Kubernetes taints to be applied to each node (list)
* `zone` - (Required) Zone GKE cluster (string)

### `cluster_auth_endpoint`

#### Arguments

* `ca_certs` - (Optional) CA certs for the authorized cluster endpoint (string)
* `enabled` - (Optional) Enable the authorized cluster endpoint. Default `true` (bool)
* `fqdn` - (Optional) FQDN for the authorized cluster endpoint (string)

### `cluster_monitoring_input`

#### Arguments

* `answers` - (Optional/Computed) Key/value answers for monitor input (map)
=======
### `cluster_template_answers`

#### Arguments

* `cluster_id` - (Optional) Cluster ID to apply answer (string)
* `project_id` - (Optional) Project ID to apply answer (string)
* `values` - (Optional) Key/values for answer (map)

### `cluster_template_questions`

#### Arguments

* `default` - (Required) Default variable value (string)
* `required` - (Optional) Required variable. Default `false` (bool)
* `type` - (Optional) Variable type. `boolean`, `int` and `string` are allowed. Default `string` (string)
* `variable` - (Optional) Variable name (string)
>>>>>>> c6a2cbc... Feat: added rancher2_cluster_template datasource and resource. For rancher V2.3.x. Doc files

### `cluster_registration_token`

#### Attributes

* `cluster_id` - (Computed) Cluster ID (string)
* `name` - (Computed) Name of cluster registration token (string)
* `command` - (Computed) Command to execute in a imported k8s cluster (string)
* `insecure_command` - (Computed) Insecure command to execute in a imported k8s cluster (string)
* `manifest_url` - (Computed) K8s manifest url to execute with `kubectl` to import an existing k8s cluster (string)
* `node_command` - (Computed) Node command to execute in linux nodes for custom k8s cluster (string)
* `token` - (Computed) Token for cluster registration token object (string)
* `windows_node_command` - (Computed) Node command to execute in windows nodes for custom k8s cluster (string)
* `annotations` - (Computed) Annotations for cluster registration token object (map)
* `labels` - (Computed) Labels for cluster registration token object (map)

## Timeouts

`rancher2_cluster` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters.
- `update` - (Default `30 minutes`) Used for cluster modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters.

## Import

Clusters can be imported using the Rancher Cluster ID

```
$ terraform import rancher2_cluster.foo <cluster>
```
