---
page_title: "Rancher2: rancher2_config_map_v2 Resource"
---

# rancher2\_config\_map\_v2 Resource

Provides a Rancher ConfigMap v2 resource. This can be used to create k8s configMaps for Rancher v2 environments and retrieve their information. ConfigMap v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
# Create a new Rancher2 ConfigMap V2
resource "rancher2_config_map_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = "foo"
  data = {
  	mydata1 = "<data1>"
  	mydata2 = "<data2>"
  	mydata3 = "<data3>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster id of the configMap V2 (string)
* `data` - (Required) The data of the configMap v2 (map)
* `name` - (Required/ForceNew) The name of the configMap v2 (string)
* `namespace` - (Optional/ForceNew) The namespaces of the configMap v2. Default: `default` (string)
* `immutable` - (Optional) If set to true, any configMap update will remove and recreate the configMap. This is a beta field enabled by k8s `ImmutableEphemeralVolumes` feature gate. Default: `false` (bool)
* `annotations` - (Optional/Computed) Annotations for the configMap v2 (map)
* `labels` - (Optional/Computed) Labels for the configMap v2 (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)

## Timeouts

`rancher2_configMap` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating v2 configMaps.
- `update` - (Default `10 minutes`) Used for v2 configMap modifications.
- `delete` - (Default `10 minutes`) Used for deleting v2 configMaps.

## Import

V2 configMaps can be imported using the Rancher cluster ID, ConfigMap V2 namespace and name.

```
$ terraform import rancher2_config_map_v2.foo &lt;CLUSTER_ID&gt;.&lt;SECRET_V2_NAMESPACE&gt;/&lt;SECRET_V2_NAME&gt;
```
