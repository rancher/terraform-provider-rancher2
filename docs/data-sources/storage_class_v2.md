---
page_title: "rancher2_storage_class_v2 Data Source"
---

# rancher2\_storage_class\_v2 Data Source

Use this data source to retrieve information about a Rancher2 Storage Class v2. Storage Class v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```
data "rancher2_storage_class_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = "foo"
}
```

## Argument Reference

* `cluster_id` - (Required) The cluster id of the storageClass V2 (string)
* `name` - (Required) The name of the storageClass v2 (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)
* `k8s_provisioner` - (Computed) The provisioner of the storageClass v2 (string)
* `allow_volume_expansion` - (Computed) Is the provisioner of the storageClass v2 allowing volume expansion? (bool)
* `mount_options` - (Computed) The mount options for storageClass v2 (list)
* `parameters` - (Computed) The parameters for storageClass v2 (string)
* `reclaim_policy` - (Computed) The reclaim policy for storageClass v2 (string)
* `volume_binding_mode` - (Computed) The volume binding mode for storageClass v2 (string)
* `annotations` - (Computed) Annotations for the storageClass v2 (map)
* `labels` - (Computed) Labels for the storageClass v2 (map)