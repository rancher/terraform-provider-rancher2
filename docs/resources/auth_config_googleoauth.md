---
page_title: "rancher2_auth_config_googleoauth Resource"
---

# rancher2\_auth\_config\_googleoauth

Provides a Rancher v2 Auth Config GoogleOauth resource. This can be used to configure and enable Auth Config GoogleOauth for Rancher v2 RKE clusters and retrieve their information.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time. 

## Example Usage

```hcl
# Create a new rancher2 Auth Config GoogleOauth
resource "rancher2_auth_config_googleoauth" "googleoauth" {
  admin_email = "<GOOGLEOAUTH_ADMIN_EMAIL>"
  hostname = "<GOOGLEOAUTH_HOSTNAME>"
  oauth_credential = "<GOOGLEOAUTH_OAUTH_CREDENTIAL>"
  service_account_credential = "<GOOGLEOAUTH_SERVICE_ACCOUNT_CREDENTIAL>"
  nested_group_membership_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `admin_email` - (Required) GoogleOauth admin email (string)
* `hostname` - (Required) GoogleOauth hostname to connect (string)
* `oauth_credential` - (Required/Sensitive) GoogleOauth credential (string)
* `service_account_credential` - (Required/Sensitive) GoogleOauth service account credential (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `googleoauth_user://<USER_ID>` `googleoauth_group://<GROUP_ID>` (list)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `nested_group_membership_enabled` - (Optional) Nested group membership enable. Default `false` (bool)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)
