---
page_title: "rancher2_global_role Resource"
---

# rancher2\_global\_role Resource

Provides a Rancher v2 Global Role resource. This can be used to create Global Role for Rancher v2 and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Global Role
resource "rancher2_global_role" "foo" {
  name             = "foo"
  new_user_default = true
  description      = "Terraform global role acceptance test"

  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["create"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Global role name (string)
* `description` - (Optional/Computed) Global role description (string)
* `new_user_default` - (Optional) Whether or not this role should be added to new users. Default `false` (bool)
* `rules` - (Optional/Computed) Global role policy rules (list)
* `annotations` - (Optional/Computed) Annotations for global role object (map)
* `labels` - (Optional/Computed) Labels for global role object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `builtin` - (Computed) Builtin global role (bool)

## Nested blocks

### `rules`

#### Arguments

* `api_groups` - (Optional) Policy rule api groups (list)
* `non_resource_urls` - (Optional) Policy rule non resource urls (list)
* `resource_names` - (Optional) Policy rule resource names (list)
* `resources` - (Optional) Policy rule resources (list)
* `verbs` - (Optional) Policy rule verbs. `bind`, `create`, `delete`, `deletecollection`, `escalate`, `get`, `impersonate`, `list`, `patch`, `update`, `use`, `view`, `watch`, `own` and `*` values are supported (list)

## Timeouts

`rancher2_global_role` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating global role.
- `update` - (Default `10 minutes`) Used for global role modifications.
- `delete` - (Default `10 minutes`) Used for deleting global role.

## Import

Global Role can be imported using the Rancher Global Role ID

```
$ terraform import rancher2_global_role.foo &lt;global_role_id&gt;
```
