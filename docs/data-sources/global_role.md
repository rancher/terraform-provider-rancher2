---
page_title: "rancher2_global_role Data Source"
---

# rancher2\_global\_role Data Source

Use this data source to retrieve information about a Rancher v2 global role resource.

## Example Usage

```hcl
data "rancher2_global_role" "foo" {
  name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Global Role (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `builtin` - (Computed) Builtin global role (bool)
* `description` - (Computed) Global role description (string)
* `new_user_default` - (Computed) Whether or not this role should be added to new users (bool)
* `rules` - (Computed) Global role policy rules (list)
* `annotations` - (Computed) Annotations for global role object (map)
* `labels` - (Computed) Labels for global role object (map)
