---
page_title: "rancher2_service_account_token Resource"
---

# rancher2\_service\_account\_token Resource

Provides a Rancher v2 Token resource, specifically to create service account tokens. Service accounts tokens are tokens for other users (service accounts) than the Rancher v2 provider user. To create a service account token the username/password for the Rancher User must be known.

There are 2 kind of tokens:
- no scoped: valid for global system.
- scoped: valid for just a specific cluster (`cluster_id` should be provided).

Tokens can only be created for a Rancher User with at least the `user-base` global role binding in order to enable user login.

Tokens can't be updated once created. Any diff in token data will recreate the token. If any token expire, Rancher2 provider will generate a diff to regenerate it.

## Example Usage

```hcl
# Create a rancher2 Token
resource "rancher2_user" "foo" {
  name = "foo"
  username = "foo"
  password = "changeme"
  enabled = true
}

resource "rancher2_global_role_binding" "foo-login" {
  name = "foo-login-binding"
  global_role_id = "user-base"
  user_id = rancher2_user.foo.id
}

resource "rancher2_service_account_token" "foo" {
  username = rancher2_user.foo.username
  password = rancher2_user.foo.password
  description = "foo token"
  ttl = 0

  depends_on = [
    rancher2_global_role_binding.foo-login
  ]
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required/ForceNew) The user username (string)
* `password` - (Required/ForceNew) The user password (string)
* `cluster_id` - (Optional/ForceNew) Cluster ID for scoped token (string)
* `description` - (Optional/ForceNew) Token description (string)
* `renew` - (Optional/ForceNew) Renew token if expired or disabled. If `true`, a terraform diff would be generated to renew the token if it's disabled or expired. If `false`, the token will not be renewed. Default `true` (bool)
* `ttl` - (Optional/ForceNew) Token time to live in seconds. Default `0` (int) 

From Rancher v2.4.6 `ttl` is read in minutes at Rancher API. To avoid breaking change on the provider, we still read in seconds but rounding up division if required.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `access_key` - (Computed) Token access key part (string)
* `enabled` - (Computed) Token is enabled (bool)
* `expired` - (Computed) Token is expired (bool)
* `name` - (Computed) Token name (string)
* `secret_key` - (Computed/Sensitive) Token secret key part (string)
* `token` - (Computed/Sensitive) Token value (string)
* `user_id` - (Computed) Token user ID (string)
* `annotations` - (Computed) Annotations of the token (map)
* `labels` - (Computed) Labels of the token (map)
* `temp_token` - (Computed) Generated API temporary token as helper. Should be empty (string)
* `temp_token_id` - (Computed) Generated API temporary token id as helper. Should be empty (string)

## Timeouts

`rancher2_service_account_token` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating tokens.
- `update` - (Default `5 minutes`) Used for token modifications.
- `delete` - (Default `5 minutes`) Used for deleting tokens.
