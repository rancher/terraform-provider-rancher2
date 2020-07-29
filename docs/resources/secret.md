---
page_title: "rancher2_secret Resource"
---

# rancher2\_secret Resource

Provides a Rancher v2 Secret resource. This can be used to create secrets for Rancher v2 environments and retrieve their information.

Depending of the availability, there are 2 types of Rancher v2 secrets:
- Project secret: Available to all namespaces in the `project_id`
- Namespaced secret: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Create a new rancher2 Project Secret
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret foo"
  project_id = "<project_id>"
  data = {
    address = base64encode("test.io")
    username = base64encode("user2")
    password = base64encode("pass")
  }
}
```

```hcl
# Create a new rancher2 Namespaced Secret
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret foo"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
  data = {
    address = base64encode("test.io")
    username = base64encode("user2")
    password = base64encode("pass")
  }
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Required/Sensitive) Secret key/value data. Base64 encoding required for values (map)
* `project_id` - (Required/ForceNew) The project id where to assign the secret (string)
* `description` - (Optional) A secret description (string)
* `name` - (Optional/ForceNew) The name of the secret (string)
* `namespace_id` - (Optional/ForceNew) The namespace id where to assign the namespaced secret (string)
* `annotations` - (Optional/Computed) Annotations for secret object (map)
* `labels` - (Optional/Computed) Labels for secret object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_secret` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating registries.
- `update` - (Default `10 minutes`) Used for secret modifications.
- `delete` - (Default `10 minutes`) Used for deleting registries.

## Import

Secrets can be imported using the secret ID in the format `<namespace_id>.<project_id>.<secret_id>`

```
$ terraform import rancher2_secret.foo &lt;namespace_id&gt;.&lt;project_id&gt;.&lt;secret_id&gt;
```

`<namespace_id>` is optional, just needed for namespaced secret.
