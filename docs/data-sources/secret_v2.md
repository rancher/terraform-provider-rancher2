---
page_title: "rancher2_secret_v2 Datasource"
---

# rancher2\_secret\_v2 Datasource

Use this data source to retrieve information about a Rancher2 secret v2. Secret v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
data "rancher2_secret_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = <SECRET_V2_NAME>
  namespace = <SECRET_V2_NAMESPACE>
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster id of the secret V2 (string)
* `name` - (Required) The name of the secret v2 (string)
* `namespace` - (Optional) The namespaces of the secret v2. Default: `default` (string)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)
* `data` - (Computed/Sensitive) The data of the secret v2 (map)
* `type` - (Computed) The type of the k8s secret, used to facilitate programmatic handling of secret data, [More info](https://github.com/kubernetes/api/blob/release-1.20/core/v1/types.go#L5772) about k8s secret types and expected format (string)
* `immutable` - (Computed) If set to true, any secret update will remove and recreate the secret. This is a beta field enabled by k8s `ImmutableEphemeralVolumes` feature gate (bool)
* `annotations` - (Computed) Annotations for the secret v2 (map)
* `labels` - (Computed) Labels for the secret v2 (map)
