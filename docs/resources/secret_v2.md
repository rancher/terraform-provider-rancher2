---
page_title: "Rancher2: rancher2_secret_v2 Resource"
---

# rancher2\_secret\_v2 Resource

Provides a Rancher Secret v2 resource. This can be used to create k8s secrets for Rancher v2 environments and retrieve their information. Secret v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
# Create a new Rancher2 Secret V2
resource "rancher2_secret_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = "foo"
  data = {
  	mydata1 = "<data1>"
  	mydata2 = "<data2>"
  	mydata3 = "<data3>"
  }
}
# Create a new Rancher2 Secret V2 basic-auth
resource "rancher2_secret_v2" "foo" {
  cluster_id = <CLUSTER_ID>
  name = "foo"
  namespace = "<mynamespace>"
  type = "kubernetes.io/basic-auth"
  data = {
  	password = "<mysecret>"
  	username = "<myuser>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster id of the secret V2 (string)
* `data` - (Required/Sensitive) The data of the secret v2 (map)
* `name` - (Required) The name of the secret v2 (string)
* `namespace` - (Optional/ForceNew) The namespaces of the secret v2. Default: `default` (string)
* `type` - (Optional) The type of the k8s secret, used to facilitate programmatic handling of secret data, [More info](https://github.com/kubernetes/api/blob/release-1.20/core/v1/types.go#L5772) about k8s secret types and expected format.  Default: `Opaque` (string)
* `immutable` - (Optional) If set to true, any secret update will remove and recreate the secret. This is a beta field enabled by k8s `ImmutableEphemeralVolumes` feature gate. Default: `false` (bool)
* `annotations` - (Optional/Computed) Annotations for the secret v2 (map)
* `labels` - (Optional/Computed) Labels for the secret v2 (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `resource_version` - (Computed) The k8s resource version (string)

## Timeouts

`rancher2_secret` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating v2 secrets.
- `update` - (Default `10 minutes`) Used for v2 secret modifications.
- `delete` - (Default `10 minutes`) Used for deleting v2 secrets.

## Import

V2 secrets can be imported using the Rancher cluster ID, Secret V2 namespace and name.

```
$ terraform import rancher2_secret_v2.foo &lt;CLUSTER_ID&gt;.&lt;SECRET_V2_NAMESPACE&gt;/&lt;SECRET_V2_NAME&gt;
```
