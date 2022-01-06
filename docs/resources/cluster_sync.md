---
page_title: "rancher2_cluster_sync Resource"
---

# rancher2\_cluster\_sync Resource

Provides a Rancher v2 Cluster Sync dummy resource. This can be used to create a Cluster Sync to wait for a Rancher v2 Cluster resource `active` state.

This dummy resource doesn't create anything at Rancher side. It's used to sync terraform resources that depends of Rancher v2 Cluster resource in `active` state. This resource will wait until `cluster_id` is `active` on `terraform apply`. It also helps to sync `terraform destroy` dependencies, specially useful if cluster is using node pools.

This resource will also compute attributes with useful cluster related data (see Attributes Reference section). 

## Example Usage

```hcl
# Create a new rancher2 rke Cluster 
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
    access_key = "<AWS_ACCESS_KEY>"
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
  cluster_id =  rancher2_cluster.foo-custom.id
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  node_template_id = rancher2_node_template.foo.id
  quantity = 3
  control_plane = true
  etcd = true
  worker = true
}
# Create a new rancher2 Cluster Sync
resource "rancher2_cluster_sync" "foo-custom" {
  cluster_id =  rancher2_cluster.foo-custom.id
  node_pool_ids = [rancher2_node_pool.foo.id]
}
# Create a new rancher2 Project
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = rancher2_cluster_sync.foo-custom.id
  description = "Terraform namespace acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster ID that is syncing (string)
* `node_pool_ids` - (Optional) The node pool IDs used by the cluster id (list)
* `wait_alerting` - (Optional) Wait until alerting is up and running. Default: `false` (bool)
* `wait_catalogs` - (Optional) Wait until all catalogs are downloaded and active. Default: `false` (bool)
* `wait_monitoring` - (Optional) Wait until monitoring is up and running. Default: `false` (bool)
* `state_confirm` - (Optional) Wait until active status is confirmed a number of times (wait interval of 5s). Default: `1` means no confirmation (int)

**Note** `state_confirm` would be useful, if you have troubles for creating/updating custom clusters that eventually are reaching `active` state before they are fully installed. For example: setting `state_confirm = 2` will assure that the cluster has been in `active` state for at least 5 seconds, `state_confirm = 3` assure at least 10 seconds, etc

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource. Same as `cluster_id` (string)
* `default_project_id` - (Computed) Default project ID for the cluster sync (string)
* `kube_config` - (Computed/Sensitive) Kube Config generated for the cluster sync (string)
* `nodes` - (Computed) The cluster nodes (list).
* `system_project_id` - (Computed) System project ID for the cluster sync (string)

**Note** For Rancher 2.6.0 and above: if setting `kubeconfig-generate-token=false` then the generated `kube_config` will not contain any user token. `kubectl` will generate the user token executing the [rancher cli](https://github.com/rancher/cli/releases/tag/v2.6.0), so it should be installed previously.

## Nested blocks

### `nodes`

#### Arguments

* `annotations` - (Computed) Annotations of the node (map).
* `capacity` - (Computed) The total resources of a node (map).
* `cluster_id` - (Computed) The Cluster ID of the node (string).
* `external_ip_address` - (Computed)  The external IP address of the node (string).
* `hostname` - (Computed) The hostname of the node (string).
* `id` - (Computed) The ID of the node (string)
* `ip_address` - (Computed) The private IP address of the node (string).
* `labels` - (Computed) Labels of the node (map).
* `name` - (Computed) The name of the node (string).
* `node_pool_id` - (Computed) The Node Pool ID of the node (string).
* `node_template_id` - (Computed) The Node Template ID of the node (string).
* `provider_id` - (Computed) The Provider ID of the node (string).
* `requested_hostname` - (Computed) The requested hostname (string).
* `roles` - (Computed) Roles of the node. `controlplane`, `etcd` and `worker`. (list)
* `ssh_user` - (Computed/Sensitive) The user to connect to the node (string).
* `system_info` - (Computed) General information about the node, such as kernel version, kubelet and kube-proxy version, Docker version (if used), and OS name.

### `system_info`

* `container_runtime_version` - (Computed) ContainerRuntime Version reported by the node through runtime remote API (e.g. docker://1.5.0).
* `kernel_version` - (Computed) Kernel Version reported by the node from 'uname -r' (e.g. 3.16.0-0.bpo.4-amd64).
* `kube_proxy_version` - (Computed) KubeProxy Version reported by the node.
* `kubelet_version` - (Computed) Kubelet Version reported by the node.
* `operating_system` - (Computed) The Operating System reported by the node.

## Timeouts

`rancher2_cluster_sync` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating cluster sync.
- `update` - (Default `30 minutes`) Used for cluster sync modifications.
- `delete` - (Default `30 minutes`) Used for deleting cluster sync.
