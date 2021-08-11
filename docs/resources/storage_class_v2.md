---
page_title: "Rancher2: rancher2_storage_class_v2 Resource"
---

# rancher2\_storageClass\_v2 Resource

Provides a Rancher Storage Class v2 resource. This can be used to manage k8s storage classes for Rancher v2 clusters and retrieve their information. Storage Class v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
# Create a new Rancher2 Storage Class V2
resource "rancher2_storage_class_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = "foo"
  parameters = {
    "gidAllocate" = "true"
    "gidMax" = "50000"
    "gidMin" = "40000"
  }
  k8s_provisioner: "example.com/aws-efs"
  reclaim_policy: "Delete"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster id of the storageClass V2 (string)
* `name` - (Required/ForceNew) The name of the storageClass v2 (string)
* `k8s_provisioner` - (Required/ForceNew) The provisioner of the storageClass v2 (string)
* `allow_volume_expansion` - (Optional/Computed) Is the provisioner of the storageClass v2 allowing volume expansion? (bool)
* `mount_options` - (Optional/Computed) The mount options for storageClass v2 (list)
* `parameters` - (Optional) The parameters for storageClass v2 (string)
* `reclaim_policy` - (Optional/ForceNew) The reclaim policy for storageClass v2. `Delete`, `Recycle` and `Retain` values are allowed. Default: `Delete` (string)
* `volume_binding_mode` - (Optional/ForceNew) The volume binding mode for storageClass v2. `Immediate` and `WaitForFirstConsumer` values are allowed. Default: `Immediate` (string)
* `annotations` - (Optional/Computed) Annotations for the storageClass v2 (map)
* `labels` - (Optional/Computed) Labels for the storageClass v2 (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)

## Timeouts

`rancher2_storage_class_v2` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating v2 storageClasss.
- `update` - (Default `10 minutes`) Used for v2 storageClass modifications.
- `delete` - (Default `10 minutes`) Used for deleting v2 storageClasss.

## Import

V2 storage classs can be imported using the Rancher cluster ID and StorageClass V2 name.

```
$ terraform import rancher2_storage_class_v2.foo &lt;CLUSTER_ID&gt;.&lt;STORAGE_CLASS_V2_NAME&gt;
```
