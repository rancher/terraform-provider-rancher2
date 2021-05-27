---
page_title: "rancher2_auth_config_keycloak Resource"
---

# rancher2\_auth\_config\_keycloak Resource

Provides a Rancher v2 Auth Config KeyCloak resource. This can be used to configure and enable Auth Config KeyCloak for Rancher v2 RKE clusters and retrieve their information.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time.

## Example Usage

```hcl
# Create a new rancher2 Auth Config KeyCloak
resource "rancher2_auth_config_keycloak" "keycloak" {
  display_name_field = "<DISPLAY_NAME_FIELD>"
  groups_field = "<GROUPS_FIELD>"
  idp_metadata_content = "<IDP_METADATA_CONTENT>"
  rancher_api_host = "https://<RANCHER_API_HOST>"
  sp_cert = "<SP_CERT>"
  sp_key = "<SP_KEY>"
  uid_field = "<UID_FIELD>"
  user_name_field = "<USER_NAME_FIELD>"
}
```

## Argument Reference

The following arguments are supported:

* `display_name_field` - (Required) KeyCloak display name field (string)
* `groups_field` - (Required) KeyCloak group field (string)
* `idp_metadata_content` - (Required/Sensitive) KeyCloak IDP metadata content (string)
* `rancher_api_host` - (Required) Rancher URL. URL scheme needs to be specified, `https://<RANCHER_API_HOST>` (string)
* `sp_cert` - (Required/Sensitive) KeyCloak SP cert (string)
* `sp_key` - (Required/Sensitive) KeyCloak SP key (string)
* `uid_field` - (Required) KeyCloak UID field (string)
* `user_name_field` - (Required) KeyCloak user name field (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `keycloak_user://<USER_ID>`  `keycloak_group://<GROUP_ID>` (list)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)
* `entity_id` - (Optional/Computed) KeyCloak Client ID field (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)
