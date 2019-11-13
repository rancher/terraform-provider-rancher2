---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster_sync"
sidebar_current: "docs-rancher2-resource-cluster_sync"
description: |-
  Provides a Rancher v2 Cluster Sync dummy resource. This can be used to create a Cluster Sync to wait for a Rancher v2 Cluster resource `active` state.
---

# rancher2\_cluster\_sync

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
# Create a new rancher2 Cluster Sync
resource "rancher2_cluster_sync" "foo-custom" {
  cluster_id =  "${rancher2_cluster.foo-custom.id}"
  node_pool_ids = ["${rancher2_node_pool.foo.id}"]
}
# Create a new rancher2 Project
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "${rancher2_cluster_sync.foo-custom.id}"
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


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource. Same as `cluster_id` (string)
* `default_project_id` - (Computed) Default project ID for the cluster sync (string)
* `kube_config` - (Computed) Kube Config generated for the cluster sync (string)
* `system_project_id` - (Computed) System project ID for the cluster sync (string)

## Timeouts

`rancher2_cluster_sync` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating cluster sync.
- `update` - (Default `30 minutes`) Used for cluster sync modifications.
- `delete` - (Default `30 minutes`) Used for deleting cluster sync.
