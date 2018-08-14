---
layout: "rancher2"
page_title: "Rancher2: rancher2_project_role_template_binding"
sidebar_current: "docs-rancher2-resource-project_role_template_binding"
description: |-
  Provides a Rancher v2 Project Role Template Binding resource. This can be used to create Project Role Template Bindings for rancher v2 environments and retrieve their information.
---

# rancher2\_project_role_template_binding

Provides a Rancher v2 Project Role Template Binding resource. This can be used to create Project Role Template Bindings for rancher v2 environments and retrieve their information.

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

* `project_id` - (Required) The project id where bind project role template.
* `role_template_id` - (Required) The role template id from create project role template binding.
* `name` - (Required) The name of the project role template binding.
* `group_id` - (Optional) The group ID to assign project role template binding.
* `group_principal_id` - (Optional) The group_principal ID to assign project role template binding
* `user_id` - (Optional) The user ID to assign project role template binding
* `user_principal_id` - (Optional) The user_principal ID to assign project role template binding
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Project Role Template Bindings can be imported using the rancher Project Role Template Binding ID

```
$ terraform import rancher2_project_role_template_binding.foo <project_role_template_binding_id>
```

