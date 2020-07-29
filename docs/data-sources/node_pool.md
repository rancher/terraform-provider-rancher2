---
page_title: "rancher2_node_pool Data Source"
---

# rancher2\_node\_pool Data Source

Use this data source to retrieve information about a Rancher v2 Node Pool resource.

## Example Usage

```hcl
data "rancher2_node_pool" "foo" {
  cluster_id =  rancher2_cluster.foo-custom.id
  name = "foo"
}
```

## Argument Reference

* `cluster_id` - (Required) The RKE cluster id to use Node Pool (string)
* `name` - (Required) The name of the Node Pool (string)
* `node_template_id` - (Optional/Computed) The Node Template ID to use for node creation (string)


## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `hostname_prefix` - (Computed) The prefix for created nodes of the Node Pool (string)
* `delete_not_ready_after_secs` - (Computed) Delete not ready node after secs. Default `0` (int)
* `node_taints` - (Computed) Node taints (List)
* `quantity` - (Computed) The number of nodes to create on Node Pool (int)
* `control_plane` - (Computed) RKE control plane role for created nodes (bool)
* `etcd` - (Computed) RKE etcd role for created nodes (bool)
* `worker` - (Computed) RKE role role for created nodes (bool)
* `annotations` - (Computed) Annotations for Node Pool object (map)
* `labels` - (Computed) Labels for Node Pool object (map)
