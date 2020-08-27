---
page_title: "rancher2_token Resource"
---

# rancher2\_token Resource

Provides a Rancher v2 Token resource. This can be used to create Tokens for Rancher v2 provider user and retrieve their information.

There are 2 kind of tokens:
- no scoped: valid for global system.
- scoped: valid for just a specific cluster (`cluster_id` should be provided).

Tokens can't be updated once created. Any diff in token data will recreate the token. If any token expire, Rancher2 provider will generate a diff to regenerate it.

## Example Usage

```hcl
# Create a new rancher2 Token
resource "rancher2_token" "foo" {
  description = "foo token"
  ttl = 1200
}
# Create a new rancher2 Token scoped
resource "rancher2_token" "foo" {
  cluster_id = "<cluster-id>"
  description = "foo token"
  ttl = 1200
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional/ForceNew) Cluster ID for scoped token (string)
* `description` - (Optional/ForceNew) Token description (string)
* `renew` - (Optional/ForceNew) Renew token if expired or disabled. If `true`, a terraform diff would be generated to renew the token if it's disabled or expired. If `false`, the token will not be renewed. Default `true` (bool)
* `ttl` - (Optional/ForceNew) Token time to live in seconds. Default `0` (int) 

From Rancher v2.4.6 `ttl` is readed in minutes at Rancher API. To avoid breaking change on the provider, we still read in seconds but rounding up division if required.

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

## Timeouts

`rancher2_token` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating tokens.
- `update` - (Default `5 minutes`) Used for token modifications.
- `delete` - (Default `5 minutes`) Used for deleting tokens.
