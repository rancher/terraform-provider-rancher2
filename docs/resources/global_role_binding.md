---
page_title: "rancher2_global_role_binding Resource"
---

# rancher2\_global\_role\_binding Resource

Provides a Rancher v2 Global Role Binding resource. This can be used to create Global Role Bindings for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Global Role Binding using user_id
resource "rancher2_global_role_binding" "foo" {
  name = "foo"
  global_role_id = "admin"
  user_id = "user-XXXXX"
}
# Create a new rancher2 Global Role Binding using group_principal_id
resource "rancher2_global_role_binding" "foo2" {
  name = "foo2"
  global_role_id = "admin"
  group_principal_id = "local://g-XXXXX"
}
```

## Argument Reference

The following arguments are supported:

* `global_role_id` - (Required/ForceNew) The role id from create global role binding (string)
* `group_principal_id` - (Optional/Computed/ForceNew) The group principal ID to assign global role binding (only works with external auth providers that support groups). Rancher v2.4.0 or higher is required (string)
* `user_id` - (Optional/Computed/ForceNew) The user ID to assign global role binding (string)
* `name` - (Optional/Computed/ForceNew) The name of the global role binding (string)
* `annotations` - (Optional/Computed) Annotations for global role binding (map)
* `labels` - (Optional/Computed) Labels for global role binding (map)

**Note** user `user_id` OR group `group_principal_id` must be defined

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
$ terraform import rancher2_global_role_binding.foo &lt;GLOBAL_ROLE_BINDING_ID&gt;
```

