---
layout: "rancher2"
page_title: "Rancher2: rancher2_project_role_template_binding"
sidebar_current: "docs-rancher2-resource-project_role_template_binding"
description: |-
  Provides a Rancher v2 Project Role Template Binding resource. This can be used to create Project Role Template Bindings for Rancher v2 environments and retrieve their information.
---

# rancher2\_project\_role\_template\_binding

Provides a Rancher v2 Project Role Template Binding resource. This can be used to create Project Role Template Bindings for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Project Role Template Binding
resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "<project_id>"
  role_template_id = "<role_template_id>"
  user_id = "<user_id>"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required/ForceNew) The project id where bind project role template (string)
* `role_template_id` - (Required/ForceNew) The role template id from create project role template binding (string)
* `name` - (Required) The name of the project role template binding (string)
* `group_id` - (Optional) The group ID to assign project role template binding (string)
* `group_principal_id` - (Optional/Computed) The group_principal ID to assign project role template binding (string)
* `user_id` - (Optional) The user ID to assign project role template binding (string)
* `user_principal_id` - (Optional/Computed) The user_principal ID to assign project role template binding (string)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_project_role_template_binding` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating project role template bindings.
- `update` - (Default `10 minutes`) Used for project role template binding modifications.
- `delete` - (Default `10 minutes`) Used for deleting project role template bindings.

## Import

Project Role Template Bindings can be imported using the Rancher Project Role Template Binding ID

```
$ terraform import rancher2_project_role_template_binding.foo <project_role_template_binding_id>
```

