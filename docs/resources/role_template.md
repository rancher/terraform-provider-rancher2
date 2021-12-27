---
page_title: "rancher2_role_template Resource"
---

# rancher2\_role\_template Resource

Provides a Rancher v2 Role Template resource. This can be used to create Role Template for Rancher v2 and retrieve their information. 

`cluster` and `project` scopes are supported for role templates.

## Example Usage

```hcl
# Create a new rancher2 cluster Role Template
resource "rancher2_role_template" "foo" {
  name = "foo"
  context = "cluster"
  default_role = true
  description = "Terraform role template acceptance test"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["create"]
  }
}
```

```hcl
# Create a new rancher2 project Role Template
resource "rancher2_role_template" "foo" {
  name = "foo"
  context = "project"
  default_role = true
  description = "Terraform role template acceptance test"
  rules {
    api_groups = ["*"]
    resources = ["secrets"]
    verbs = ["create"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Role template name (string)
* `administrative` - (Optional) Administrative role template. Default `false` (bool)
* `context` - (Optional) Role template context. `cluster` and `project` values are supported. Default: `cluster` (string)
* `default_role` - (Optional) Default role template for new created cluster or project. Default `false` (bool)
* `description` - (Optional/Computed) Role template description (string)
* `external` - (Optional) External role template. Default `false` (bool)
* `hidden` - (Optional) Hidden role template. Default `false` (bool)
* `locked` - (Optional) Locked role template. Default `false` (bool)
* `role_template_ids` - (Optional/Computed) Inherit role template IDs (list)
* `rules` - (Optional/Computed) Role template policy rules (list)
* `annotations` - (Optional/Computed) Annotations for role template object (map)
* `labels` - (Optional/Computed) Labels for role template object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `builtin` - (Computed) Builtin role template (string)

## Nested blocks

### `rules`

#### Arguments

* `api_groups` - (Optional) Policy rule api groups (list)
* `non_resource_urls` - (Optional) Policy rule non resource urls (list)
* `resource_names` - (Optional) Policy rule resource names (list)
* `resources` - (Optional) Policy rule resources (list)
* `verbs` - (Optional) Policy rule verbs. `bind`, `create`, `delete`, `deletecollection`, `escalate`, `get`, `impersonate`, `list`, `patch`, `update`, `use`, `view`, `watch`, `own` and `*` values are supported (list)

## Timeouts

`rancher2_role_template` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating role templates.
- `update` - (Default `10 minutes`) Used for role template modifications.
- `delete` - (Default `10 minutes`) Used for deleting role templates.

## Import

Role Template can be imported using the Rancher Role Template ID

```
$ terraform import rancher2_role_template.foo &lt;role_template_id&gt;
```
