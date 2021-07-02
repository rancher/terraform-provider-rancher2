---
page_title: "rancher2_node_pool Resource"
---

# rancher2\_node\_pool Resource

Provides a Rancher v2 Node Pool resource. This can be used to create Node Pool, using Node template for Rancher v2 RKE clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 RKE Cluster 
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
# Create a new rancher2 Cloud Credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description= "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
    access_key = "XXXXXXXXXXXXXXXXXXXX"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
# Create a new rancher2 Node Template
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "foo test"
  cloud_credential_id = rancher2_cloud_credential.foo.id
  amazonec2_config {
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
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The RKE cluster id to use Node Pool (string)
* `name` - (Required/ForceNew) The name of the Node Pool (string)
* `hostname_prefix` - (Required) The prefix for created nodes of the Node Pool (string)
* `node_template_id` - (Required) The Node Template ID to use for node creation (string)
* `delete_not_ready_after_secs` - (Optional) Delete not ready node after secs. For Rancher v2.3.3 or above. Default `0` (int)
* `drain_before_delete` - (Optional) Drain nodes before delete. Default: `false` (bool)
* `node_taints` - (Required) Node taints. For Rancher v2.3.3 or above (List)
* `control_plane` - (Optional) RKE control plane role for created nodes (bool)
* `etcd` - (Optional) RKE etcd role for created nodes (bool)
* `quantity` - (Optional) The number of nodes to create on Node Pool. Default `1`. Only values >= 1 allowed (int)
* `worker` - (Optional) RKE role role for created nodes (bool)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Nested blocks

### `node_taints`

#### Arguments

* `key` - (Required) Taint key (string)
* `value` - (Required) Taint value (string)
* `effect` - (Optional) Taint effect. Supported values : `"NoExecute" | "NoSchedule" | "PreferNoSchedule"` (string)
* `time_added` - (Optional) Taint time added (string)

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

Node Pool can be imported using the Rancher Node Pool ID

```
$ terraform import rancher2_node_pool.foo &lt;node_pool_id&gt;
```

