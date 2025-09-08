---
page_title: "rancher2_auth_config_generic_oidc Resource"
---

# rancher2\_auth\_config\_generic\_oidc Resource

Provides a Rancher v2 Auth Config Generic OIDC resource. This can be used to configure and enable the Generic OIDC authentication provider for Rancher v2.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time.

## Example Usage

This example configures Rancher to use GitLab as a Generic OIDC provider.

```hcl
resource "rancher2_auth_config_generic_oidc" "generic_oidc" {
  name          = "genericoidc"
  client_id     = "<GITLAB_APPLICATION_ID>"
  client_secret = "<GITLAB_CLIENT_SECRET>"
  issuer        = "https://gitlab.com"
  rancher_url   = "https://<RANCHER_URL>/verify-auth"

  # OIDC claim mapping
  scopes               = "openid profile email read_api"
  groups_field         = "groups"
  user_name_field      = "preferred_username"
  uid_field            = "sub"
  display_name_field   = "name"
  
  # For the 'genericoidc' provider, group processing must be explicitly enabled.
  group_search_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required/Sensitive) The OIDC Client ID.
* `client_secret` - (Required/Sensitive) The OIDC Client Secret.
* `issuer` - (Required) The OIDC issuer URL.
* `rancher_url` - (Required) The URL of the Rancher server. This is used as the redirect URI for the OIDC provider.
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal IDs for auth. Required if `access_mode` is `required` or `restricted`. Ex: `genericoidc_user://<USER_ID>` `genericoidc_group://<GROUP_ID>` (list)
* `auth_endpoint` - (Optional/Computed) The OIDC Auth Endpoint URL.
* `certificate` - (Optional/Sensitive) A PEM-encoded CA certificate for the OIDC provider.
* `display_name_field` - (Optional/Computed) The name of the OIDC claim to use for the user's display name. Default `name` (string)
* `enabled` - (Optional) Enable the auth config provider. Default `true` (bool)
* `groups_field` - (Optional/Computed) The name of the OIDC claim to use for the user's group memberships. Default `groups` (string)
* `group_search_enabled` - (Optional) Enable group search. Default `false` (bool)
* `jwks_url` - (Optional/Computed) The OIDC JWKS URL.
* `private_key` - (Optional/Sensitive) A PEM-encoded private key for the OIDC provider.
* `scopes` - (Optional/Computed) The OIDC scopes to request. Defaults to `openid profile email` (string)
* `token_endpoint` - (Optional/Computed) The OIDC Token Endpoint URL.
* `uid_field` - (Optional/Computed) The name of the OIDC claim to use for the user's unique ID. Default `sub` (string)
* `user_name_field` - (Optional/Computed) The name of the OIDC claim to use for the user's username. Default `preferred_username` (string)
* `userinfo_endpoint` - (Optional/Computed) The OIDC User Info Endpoint URL.
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)

## Import

Generic OIDC auth config can be imported using its name.

```
$ terraform import rancher2_auth_config_generic_oidc.generic_oidc genericoidc
```
