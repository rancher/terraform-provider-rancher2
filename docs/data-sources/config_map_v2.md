---
page_title: "rancher2_config_map_v2 Datasource"
---

# rancher2\_config\_map\_v2 Datasource

Use this data source to retrieve information about a Rancher2 configMap v2. ConfigMap v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
data "rancher2_config_map_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = <CONFIG_MAP_V2_NAME>
  namespace = <CONFIG_MAP_V2_NAMESPACE>
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster id of the configMap V2 (string)
* `name` - (Required) The name of the configMap v2 (string)
* `namespace` - (Optional) The namespaces of the configMap v2. Default: `default` (string)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)
* `data` - (Computed) The data of the configMap v2 (map)
* `immutable` - (Computed) If set to true, any configMap update will remove and recreate the configMap. This is a beta field enabled by k8s `ImmutableEphemeralVolumes` feature gate (bool)
* `annotations` - (Computed) Annotations for the configMap v2 (map)
* `labels` - (Computed) Labels for the configMap v2 (map)
