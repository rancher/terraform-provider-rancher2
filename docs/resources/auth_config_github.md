---
page_title: "rancher2_auth_config_github Resource"
---

# rancher2\_auth\_config\_github Resource

Provides a Rancher v2 Auth Config Github resource. This can be used to configure and enable Auth Config Github for Rancher v2 RKE clusters and retrieve their information.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time. 

## Example Usage

```hcl
# Create a new rancher2 Auth Config Github
resource "rancher2_auth_config_github" "github" {
  client_id = "<GITHUB_CLIENT_ID>"
  client_secret = "<GITHUB_CLIENT_SECRET>"
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required/Sensitive) Github auth Client ID (string)
* `client_secret` - (Required/Sensitive) Github auth Client secret (string)
* `hostname` - (Optional) Github hostname to connect. Default `github.com` (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `github_user://<USER_ID>`  `github_team://<GROUP_ID>` `github_org://<ORG_ID>` (list)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `tls` - (Optional) Enable TLS connection. Default `true` (bool)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)
