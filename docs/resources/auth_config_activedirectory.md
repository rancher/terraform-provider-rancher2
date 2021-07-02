---
page_title: "rancher2_auth_config_activedirectory Resource"
---

# rancher2\_auth\_config\_activedirectory Resource

Provides a Rancher v2 Auth Config ActiveDirectory resource. This can be used to configure and enable Auth Config ActiveDirectory for Rancher v2 RKE clusters and retrieve their information.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time.

## Example Usage

```hcl
# Create a new rancher2 Auth Config ActiveDirectory
resource "rancher2_auth_config_activedirectory" "activedirectory" {
  servers = ["<ACTIVEDIRECTORY_SERVER>"]
  service_account_username = "<SERVICE_DN>"
  service_account_password = "<SERVICE_PASSWORD>"
  user_search_base = "<SEARCH_BASE>"
  port = <ACTIVEDIRECTORY_PORT>
  test_username = "<USER_NAME>"
  test_password = "<USER_PASSWORD>"
}
```

## Argument Reference

The following arguments are supported:

* `servers` - (Required) ActiveDirectory servers list (list)
* `service_account_username` - (Required/Sensitive) Service account DN for access ActiveDirectory service (string)
* `service_account_password` - (Required/Sensitive) Service account password for access ActiveDirectory service (string)
* `test_username` - (Required) Username for test access to ActiveDirectory service (string)
* `test_password` - (Required/Sensitive) Password for test access to ActiveDirectory service (string)
* `user_search_base` - (Required) User search base DN (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `activedirectory_user://<DN>`  `activedirectory_group://<DN>`. The local admin (`local://<admin id>`) and the `test_username` must be added too. (list)
* `certificate` - (Optional/Sensitive) CA certificate for TLS if selfsigned (string)
* `connection_timeout` - (Optional) ActiveDirectory connection timeout. Default `5000` (int)
* `default_login_domain` - (Optional) ActiveDirectory defult login domain (string)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `group_dn_attribute` - (Optional/Computed) Group DN attribute. Default `distinguishedName` (string)
* `group_member_mapping_attribute` - (Optional/Computed) Group member mapping attribute. Default `member` (string)
* `group_member_user_attribute` - (Optional/Computed) Group member user attribute. Default `distinguishedName` (string)
* `group_name_attribute` - (Optional/Computed) Group name attribute. Default `name` (string)
* `group_object_class` - (Optional/Computed) Group object class. Default `group` (string)
* `group_search_attribute` - (Optional/Computed) Group search attribute. Default `sAMAccountName` (string)
* `group_search_base` - (Optional/Computed) Group search base (string)
* `group_search_filter` - (Optional/Computed) Group search filter (string)
* `nested_group_membership_enabled` - (Optional/Computed) Nested group membership enable. Default `false` (bool)
* `port` - (Optional) ActiveDirectory port. Default `389` (int)
* `start_tls` - (Optional/Computed) Enable start TLS connection (bool)
* `tls` - (Optional/Computed) Enable TLS connection (bool)
* `user_disabled_bit_mask` - (Optional) User disabled bit mask. Default `2` (int)
* `user_enabled_attribute` - (Optional/Computed) User enable attribute (string)
* `user_login_attribute` - (Optional/Computed) User login attribute. Default `sAMAccountName` (string)
* `user_name_attribute` - (Optional/Computed) User name attribute. Default `name` (string)
* `user_object_class` - (Optional/Computed) User object class. Default `person` (string)
* `user_search_attribute` - (Optional/Computed) User search attribute. Default `sAMAccountName|sn|givenName` (string)
* `user_search_filter` - (Optional/Computed) User search filter (string)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)
