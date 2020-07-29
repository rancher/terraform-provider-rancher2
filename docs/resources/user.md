---
page_title: "rancher2_user Resource"
---

# rancher2\_user Resource

Provides a Rancher v2 User resource. This can be used to create Users for Rancher v2 environments and retrieve their information.

When a Rancher User is created, it doesn't have a global role binding. At least, `user-base` global role binding in needed in order to enable user login.

## Example Usage

```hcl
# Create a new rancher2 User
resource "rancher2_user" "foo" {
  name = "Foo user"
  username = "foo"
  password = "changeme"
  enabled = true
}
# Create a new rancher2 global_role_binding for User
resource "rancher2_global_role_binding" "foo" {
  name = "foo"
  global_role_id = "user-base"
  user_id = rancher2_user.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required/ForceNew) The user username (string)
* `password` - (Required/ForceNew) The user password (string)
* `name` - (Optional) The user full name (string)
* `annotations` - (Optional/Computed) Annotations for global role binding (map)
* `labels` - (Optional/Computed) Labels for global role binding (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `principal_ids` - (Computed) The user principal IDs (list)

## Timeouts

`rancher2_user` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating users.
- `update` - (Default `5 minutes`) Used for user modifications.
- `delete` - (Default `5 minutes`) Used for deleting users.

## Import

Users can be imported using the Rancher User ID

```
$ terraform import rancher2_user.foo &lt;user_id&gt;
```
