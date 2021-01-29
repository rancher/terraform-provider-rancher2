---
page_title: "rancher2_role_template Data Source"
---

# rancher2\_role\_template Data Source

Use this data source to retrieve information about a Rancher v2 role template resource.

## Example Usage

```hcl
data "rancher2_role_template" "foo" {
  name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Role Template (string)
* `context` - (Optional/Computed) Role template context. `cluster` and `project` values are supported (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `builtin` - (Computed) Builtin role template (string)
* `administrative` - (Computed) Administrative role template (bool)
* `default_role` - (Computed) Default role template for new created cluster or project (bool)
* `description` - (Computed) Role template description (string)
* `external` - (Computed) External role template (bool)
* `hidden` - (Computed) Hidden role template (bool)
* `locked` - (Computed) Locked role template (bool)
* `role_template_ids` - (Computed) Inherit role template IDs (list)
* `rules` - (Computed) Role template policy rules (list)
* `annotations` - (Computed) Annotations for role template object (map)
* `labels` - (Computed) Labels for role template object (map)
