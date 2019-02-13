---
layout: "rancher2"
page_title: "Rancher2: rancher2_auth_config_activedirectory"
sidebar_current: "docs-rancher2-auth-config-activedirectory"
description: |-
  Provides a Rancher v2 Auth Config ActiveDirectory resource. This can be used to configure and enable Auth Config ActiveDirectory for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_auth\_config\_activedirectory

Provides a Rancher v2 Auth Config ActiveDirectory resource. This can be used to configure and enable Auth Config ActiveDirectory for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Auth Config ActiveDirectory
resource "rancher2_auth_config_activedirectory" "activedirectory" {
  servers = ["<ACTIVEDIRECTORY_SERVER>"]
  service_account_username = "<SERVICE_DN>"
  service_account_password = "<SERVICE_PASSWORD>"
  user_search_base = "<SEARCH_BASE>"
  username = "<TEST_USER>"
  password = "<TEST_USER_PASSWORD>"
  port = <ACTIVEDIRECTORY_PORT>
}
```

## Argument Reference

The following arguments are supported:

* `servers` - (Required) ActiveDirectory servers list ([]string).
* `service_account_username` - (Required) Service account DN for access ActiveDirectory service (string).
* `service_account_password` - (Required) Service account password for access ActiveDirectory service (string).
* `user_search_base` - (Required) User search base DN (string).
* `username` - (Required) User name to test ActiveDirectory access (string).
* `password` - (Required) User password to test ActiveDirectory access (string).
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted`
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `activedirectory_user://<DN>`  `activedirectory_group://<DN>`
* `certificate` - (Optional) CA certificate for TLS if selfsigned (string).
* `connection_timeout` - (Optional) ActiveDirectory connection timeout. Default `5000`
* `default_login_domain` - (Optional) ActiveDirectory defult lgoin domain (string).
* `enabled` - (Optional) Enable auth config provider. Default `true`.
* `group_dn_attribute` - (Optional/Computed) Group DN attribute. Default `distinguishedName`.
* `group_member_mapping_attribute` - (Optional/Computed) Group member mapping attribute. Default `member`.
* `group_member_user_attribute` - (Optional/Computed) Group member user attribute. Default `distinguishedName`.
* `group_name_attribute` - (Optional/Computed) Group name attribute. Default `name`.
* `group_object_class` - (Optional/Computed) Group object class. Default `group`.
* `group_search_attribute` - (Optional/Computed) Group search attribute. Default `sAMAccountName`.
* `group_search_base` - (Optional/Computed) Group search base (string).
* `group_search_filter` - (Optional/Computed) Group search filter (string).
* `nested_group_membership_enabled` - (Optional/Computed) Nested group membership enable. Default `false`.
* `port` - (Optional) ActiveDirectory port. Default `389`.
* `user_disabled_bit_mask` - (Optional) User disabled bit mask. Default `2`.
* `user_enabled_attribute` - (Optional/Computed) User enable attribute (string)
* `user_login_attribute` - (Optional/Computed) User login attribute. Default `sAMAccountName`.
* `user_name_attribute` - (Optional/Computed) User name attribute. Default `name`.
* `user_object_class` - (Optional/Computed) User object class. Default `person`.
* `user_search_attribute` - (Optional/Computed) User search attribute. Default `sAMAccountName|sn|givenName`.
* `user_search_filter` - (Optional/Computed) User search filter (string)
* `tls` - (Optional/Computed) Enable TLS connection (bool).
* `annotations` - (Optional/Computed) Annotations of the resource (map).
* `labels` - (Optional/Computed) Labels of the resource (map).
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `name` - (Computed) The name of the resource.
* `type` - (Computed) The type of the resource.

