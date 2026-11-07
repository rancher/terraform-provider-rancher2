---
page_title: "rancher2_auth_config_githubapp Resource"
---

# rancher2\_auth\_config\_githubapp Resource

Provides a Rancher v2 Auth Config GitHub App resource. This can be used to configure and enable a GitHub App as an authentication provider for Rancher v2.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time, unless the `multiple-auth-providers` feature is enabled in Rancher.

Destroying this resource disables the auth config in Rancher, it does not delete it.

## Example Usage

```hcl
resource "rancher2_auth_config_githubapp" "githubapp" {
  client_id     = "<GITHUB_APP_CLIENT_ID>"
  client_secret = "<GITHUB_APP_CLIENT_SECRET>"
  app_id        = "<GITHUB_APP_ID>"
  private_key   = file("<PATH_TO_PRIVATE_KEY_PEM>")
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required/Sensitive) GitHub App Client ID (string)
* `client_secret` - (Required/Sensitive) GitHub App Client secret (string)
* `app_id` - (Required) The GitHub App ID, shown on the GitHub App settings page. Must be an integer (string)
* `private_key` - (Required/Sensitive) PEM-encoded RSA private key for the GitHub App, used for signing requests (string)
* `installation_id` - (Optional) The GitHub App installation ID. If not provided, all installations for the App are queried. Must be an integer (string)
* `hostname` - (Optional) GitHub hostname to connect to. Default `github.com` (string)
* `tls` - (Optional) Enable TLS connection. Default `true` (bool)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal IDs for auth. Required if `access_mode` is `required` or `restricted`. Ex: `githubapp_user://<USER_ID>` `githubapp_team://<TEAM_ID>` `githubapp_org://<ORG_ID>` (list)
* `enabled` - (Optional) Enable the auth config provider. Default `true` (bool)
* `name` - (Optional/ForceNew) The name of the auth config in Rancher. Changing this replaces the resource. Default `githubapp` (string)
* `display_name` - (Optional) Display name for the auth config in the Rancher UI. If unset, any display name set in Rancher is cleared (string)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `type` - (Computed) The type of the resource (string)
