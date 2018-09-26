---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_pool"
sidebar_current: "docs-rancher2-resource-node_pool"
description: |-
  Provides a Rancher v2 Node Pool resource. This can be used to create Node pool for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_node\_pool

Provides a Rancher v2 Node Pool resource. This can be used to create Node Pool for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Node Pool
resource "rancher2_node_pool" "foo" {
  cluster_id =  "foo_clusterID"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  node_template_id = "foo_templateID"
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The rke cluster id to use Node Pool.
* `name` - (Required) The name of the Node Pool.
* `hostname_prefix` - (Required) The prefix for created nodes of the Node Pool.
* `node_template_id` - (Required) The Node Template ID to use for node creation.
* `quantity` - (Required) The number of nodes to create on Node Pool.
* `control_plane` - (Optional) RKE control plane role for created nodes.
* `etcd` - (Optional) RKE etcd role for created nodes.
* `worker` - (Optional) RKE role role for created nodes.
* `annotations` - (Optional/Computed) Annotations for Node Pool object.
* `labels` - (Optional/Computed) Labels for Node Pool object.
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Node Pool can be imported using the rancher Node Pool ID

```
$ terraform import rancher2_node_pool.foo <node_pool_id>
```

