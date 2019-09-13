---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_role_binding"
sidebar_current: "docs-rancher2-resource-global_role_binding"
description: |-
  Provides a Rancher v2 Global Role Binding resource. This can be used to create Global Role Bindings for Rancher v2 environments and retrieve their information.
---

# rancher2\_global\_role\_binding

Provides a Rancher v2 Global Role Binding resource. This can be used to create Global Role Bindings for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Global Role Binding
resource "rancher2_global_role_binding" "foo" {
  name = "foo"
  global_role_id = "<global_role_id>"
  user_id = "<user_id>"
}
```

## Argument Reference

The following arguments are supported:

* `global_role_id` - (Required/ForceNew) The role id from create global role binding (string)
* `user_id` - (Required/ForceNew) The user ID to assign global role binding (string)
* `name` - (Optional/Computed/ForceNew) The name of the global role binding (string)
* `annotations` - (Optional/Computed) Annotations for global role binding (map)
* `labels` - (Optional/Computed) Labels for global role binding (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_global_role_binding` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating global role bindings.
- `update` - (Default `5 minutes`) Used for global role binding modifications.
- `delete` - (Default `5 minutes`) Used for deleting global role bindings.

## Import

Global Role Bindings can be imported using the Rancher Global Role Binding ID

```
$ terraform import rancher2_global_role_binding.foo <global_role_binding_id>
```

