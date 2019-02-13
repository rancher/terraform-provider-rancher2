---
layout: "rancher2"
page_title: "Rancher2: rancher2_auth_config_freeipa"
sidebar_current: "docs-rancher2-auth-config-freeipa"
description: |-
  Provides a Rancher v2 Auth Config FreeIpa resource. This can be used to configure and enable Auth Config FreeIpa for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_auth\_config\_freeipa

Provides a Rancher v2 Auth Config FreeIpa resource. This can be used to configure and enable Auth Config FreeIpa for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Auth Config FreeIpa
resource "rancher2_auth_config_freeipa" "freeipa" {
  servers = ["<FREEIPA_SERVER>"]
  service_account_distinguished_name = "<SERVICE_DN>"
  service_account_password = "<SERVICE_PASSWORD>"
  user_search_base = "<SEARCH_BASE>"
  username = "<TEST_USER>"
  password = "<TEST_USER_PASSWORD>"
  port = <FREEIPA_PORT>
}
```

## Argument Reference

The following arguments are supported:

* `servers` - (Required) FreeIpa servers list ([]string).
* `service_account_distinguished_name` - (Required) Service account DN for access FreeIpa service (string).
* `service_account_password` - (Required) Service account password for access FreeIpa service (string).
* `user_search_base` - (Required) User search base DN (string).
* `username` - (Required) User name to test FreeIpa access (string).
* `password` - (Required) User password to test FreeIpa access (string).
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted`
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `freeipa_user://<DN>`  `freeipa_group://<DN>`
* `certificate` - (Optional) CA certificate for TLS if selfsigned (string).
* `connection_timeout` - (Optional) FreeIpa connection timeout. Default `5000`
* `enabled` - (Optional) Enable auth config provider. Default `true`.
* `group_dn_attribute` - (Optional/Computed) Group DN attribute. Default `entryDN`.
* `group_member_mapping_attribute` - (Optional/Computed) Group member mapping attribute. Default `member`.
* `group_member_user_attribute` - (Optional/Computed) Group member user attribute. Default `entryDN`.
* `group_name_attribute` - (Optional/Computed) Group name attribute. Default `cn`.
* `group_object_class` - (Optional/Computed) Group object class. Default `groupOfNames`.
* `group_search_attribute` - (Optional/Computed) Group search attribute. Default `cn`.
* `group_search_base` - (Optional/Computed) Group search base (string).
* `nested_group_membership_enabled` - (Optional/Computed) Nested group membership enable. Default `false`.
* `port` - (Optional) FreeIpa port. Default `389`.
* `user_disabled_bit_mask` - (Optional/Computed) User disabled bit mask (int).
* `user_enabled_attribute` - (Optional/Computed) User enable attribute (string)
* `user_login_attribute` - (Optional/Computed) User login attribute. Default `uid`.
* `user_member_attribute` - (Optional/Computed) User member attribute. Default `memberOf`.
* `user_name_attribute` - (Optional/Computed) User name attribute. Default `givenName`.
* `user_object_class` - (Optional/Computed) User object class. Default `inetorgperson`.
* `user_search_attribute` - (Optional/Computed) User search attribute. Default `uid|sn|givenName`.
* `tls` - (Optional/Computed) Enable TLS connection (bool).
* `annotations` - (Optional/Computed) Annotations of the resource (map).
* `labels` - (Optional/Computed) Labels of the resource (map).
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `name` - (Computed) The name of the resource.
* `type` - (Computed) The type of the resource.

