---
page_title: "rancher2_auth_config_cognito Resource"
---

# rancher2\_auth\_config\_cognito Resource

Provides a Rancher v2 Auth Config Cognito resource. This can be used to configure and enable Amazon Cognito as an authentication provider for Rancher v2.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time, unless the `multiple-auth-providers` feature is enabled in Rancher.

Destroying this resource disables the auth config in Rancher, it does not delete it.

## Example Usage

```hcl
resource "rancher2_auth_config_cognito" "cognito" {
  client_id     = "<COGNITO_CLIENT_ID>"
  client_secret = "<COGNITO_CLIENT_SECRET>"
  issuer        = "https://cognito-idp.<REGION>.amazonaws.com/<USER_POOL_ID>"
  rancher_url   = "https://<RANCHER_URL>/verify-auth"
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required/Sensitive) The OIDC Client ID (string)
* `client_secret` - (Required/Sensitive) The OIDC Client Secret (string)
* `issuer` - (Required) The OIDC issuer URL (string)
* `rancher_url` - (Required) The URL of the Rancher server. This is used as the redirect URI for the OIDC provider (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal IDs for auth. Required if `access_mode` is `required` or `restricted`. Ex: `cognito_user://<USER_ID>` `cognito_group://<GROUP_ID>` (list)
* `enabled` - (Optional) Enable the auth config provider. Default `true` (bool)
* `name` - (Optional/ForceNew) The name of the auth config in Rancher. Changing this replaces the resource. Default `cognito` (string)
* `display_name` - (Optional) Display name for the auth config in the Rancher UI. If unset, any display name set in Rancher is cleared (string)
* `auth_endpoint` - (Optional/Computed) The OIDC Auth Endpoint URL (string)
* `token_endpoint` - (Optional/Computed) The OIDC Token Endpoint URL (string)
* `userinfo_endpoint` - (Optional/Computed) The OIDC User Info Endpoint URL (string)
* `jwks_url` - (Optional/Computed) The OIDC JWKS URL (string)
* `scopes` - (Optional/Computed) The OIDC scopes to request. Defaults to `openid profile email` (string)
* `certificate` - (Optional/Sensitive) A PEM-encoded CA certificate for the OIDC provider (string)
* `private_key` - (Optional/Sensitive) A PEM-encoded private key for the OIDC provider (string)
* `group_search_enabled` - (Optional) Enable group search. Default `false` (bool)
* `groups_field` - (Optional) The name of the OIDC claim to use for the user's group memberships. Default `cognito:groups` (string)
* `name_claim` - (Optional) The OIDC claim to use for the user name. Default `cognito:username` (string)
* `email_claim` - (Optional) The OIDC claim to use for the user email (string)
* `logout_all_enabled` - (Optional) Allow the user to choose whether or not to log out of their session with the IdP (bool)
* `logout_all_forced` - (Optional) Force the user to log out of their session with the IdP (bool)
* `end_session_endpoint` - (Optional) The provider specific HTTPS URL used for logging a user out of their session (string)
* `pkce_method` - (Optional/Computed) PKCE verification method. Set to `S256` to enable PKCE verification, or empty to disable it (string)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `type` - (Computed) The type of the resource (string)
