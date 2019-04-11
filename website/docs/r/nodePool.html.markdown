---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_pool"
sidebar_current: "docs-rancher2-resource-node_pool"
description: |-
  Provides a Rancher v2 Node Pool resource. This can be used to create Node pool, using Node template for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_node\_pool

Provides a Rancher v2 Node Pool resource. This can be used to create Node Pool, using Node template for rancher v2 rke clusters and retrieve their information.

## Example Usage

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
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The rke cluster id to use Node Pool (string)
* `name` - (Required) The name of the Node Pool (string)
* `hostname_prefix` - (Required) The prefix for created nodes of the Node Pool (string)
* `node_template_id` - (Required) The Node Template ID to use for node creation (string)
* `quantity` - (Required) The number of nodes to create on Node Pool (int)
* `control_plane` - (Optional) RKE control plane role for created nodes (bool)
* `etcd` - (Optional) RKE etcd role for created nodes (bool)
* `worker` - (Optional) RKE role role for created nodes (bool)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_node_pool` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node pools.
- `update` - (Default `10 minutes`) Used for node pool modifications.
- `delete` - (Default `10 minutes`) Used for deleting node pools.

## Import

Node Pool can be imported using the rancher Node Pool ID

```
$ terraform import rancher2_node_pool.foo <node_pool_id>
```

