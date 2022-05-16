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

```hcl
# Create a new rancher2 User with Token
resource "rancher2_user" "foo" {
  name = "Foo user"
  username = "foo"
  password = "changeme"
  enabled = true
  token_config {
    description = "Token for user Foo"
    ttl = 18000
  }
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required/ForceNew) The user username (string)
* `password` - (Required/ForceNew) The user password (string)
* `name` - (Optional) The user full name (string)
* `annotations` - (Optional/Computed) Annotations for global role binding (map)
* `labels` - (Optional/Computed) Labels for global role binding (map)
* `token_config` - (Optional) The config for a token to generate. This automatically adds a `user-base` global role binding. (list maxitems:1)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `principal_ids` - (Computed) The user principal IDs (list)
* `login_role_binding_id` - (Computed) The ID of the login role binding resource (string)
* `access_key` - (Computed) Token access key part (string)
* `auth_token` - (Computed/Sensitive) Token value (string)
* `secret_key` - (Computed/Sensitive) Token secret key part (string)
* `token_enabled` - (Computed) Token is enabled (bool)
* `token_expired` - (Computed) Token is expired (bool)
* `token_id` - (Computed) The ID of the token resource (string)
* `token_name` - (Computed) Token name (string)
* `temp_token` - (Computed/Sensitive) Generated API temporary token as helper. Should be empty (string)
* `temp_token_id` - (Computed) Generated API temporary token id as helper. Should be empty (string)

## Nested blocks

### `token_config`

#### Arguments

* `cluster_id` - (Optional) Cluster ID for scoped token (string)
* `description` - (Optional) Token description (string)
* `renew` - (Optional) Renew token if expired or disabled. If `true`, a terraform diff would be generated to renew the token if it's disabled or expired. If `false`, the token will not be renewed. Default `true` (bool)
* `ttl` - (Optional) Token time to live in seconds. Default `0` (int) 

## Timeouts

`rancher2_user` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating users.
- `update` - (Default `5 minutes`) Used for user modifications.
- `delete` - (Default `5 minutes`) Used for deleting users.

## Import

Users can be imported using the Rancher User ID (note: a token cannot be imported)

```
$ terraform import rancher2_user.foo &lt;user_id&gt;
```
