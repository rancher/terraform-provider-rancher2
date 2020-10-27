---
page_title: "rancher2_registry Resource"
---

# rancher2\_registry Resource

Provides a Rancher v2 Registry resource. This can be used to create docker registries for Rancher v2 environments and retrieve their information.

Depending of the availability, there are 2 types of Rancher v2 docker registries:
- Project registry: Available to all namespaces in the `project_id`
- Namespaced registry: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Create a new rancher2 Project Registry
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry foo"
  project_id = "<project_id>"
  registries {
    address = "test.io"
    username = "user"
    password = "pass"
  }
}
```

```hcl
# Create a new rancher2 Namespaced Registry
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry foo"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
  registries {
    address = "test.io"
    username = "user2"
    password = "pass"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required/ForceNew) The name of the registry (string)
* `project_id` - (Required/ForceNew) The project id where to assign the registry (string)
* `registries` - (Required) Registries data for registry (list)
* `description` - (Optional) A registry description (string)
* `namespace_id` - (Optional) The namespace id where to assign the namespaced registry (string)
* `annotations` - (Optional/Computed) Annotations for Registry object (map)
* `labels` - (Optional/Computed) Labels for Registry object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `registries`

#### Arguments

* `address` - (Required) Address for registry.
* `password` - (Optional) Password for the registry (string)
* `username` - (Optional) Username for the registry (string)

## Timeouts

`rancher2_registry` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating registries.
- `update` - (Default `10 minutes`) Used for registry modifications.
- `delete` - (Default `10 minutes`) Used for deleting registries.

## Import

Registries can be imported using the registry ID in the format `<namespace_id>.<project_id>.<registry_id>`

```
$ terraform import rancher2_registry.foo &lt;namespace_id&gt;.&lt;project_id&gt;.&lt;registry_id&gt;
```

`<namespace_id>` is optional, just needed for namespaced registry.
