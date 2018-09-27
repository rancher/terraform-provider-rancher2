---
layout: "rancher2"
page_title: "Rancher2: rancher2_auth_config_github"
sidebar_current: "docs-rancher2-auth-config-github"
description: |-
  Provides a Rancher v2 Auth Config Github resource. This can be used to configure and enable Auth Config Github for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_auth\_config\_github

Provides a Rancher v2 Auth Config Github resource. This can be used to configure and enable Auth Config Github for rancher v2 rke clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Auth Config Github
resource "rancher2_auth_config_github" "github" {
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

* `client_id` - (Required) Github auth Client ID (string).
* `client_secret` - (Required) Github auth Client secret (string).
* `code` - (Required) Github auth code. Generated from `https://github.com/login/oauth/authorize?client_id=<CLIENT_ID>` (string).
* `hostname` - (Optional) Gtihub hostname to connect. Defaulf `github.com`.
* `access_mode` - (Optional) Access mode for ActiveDirectory auth. `required`, `restricted`, `unrestricted` are supported. Default `restricted`
* `allowed_principal_ids` - (Optional/Computed) Allowed principal ids for auth (string).
* `enabled` - (Optional) Enable auth config for ActiveDirectory backend. Default `true`.
* `tls` - (Optional/Computed) Enable TLS connection (bool).
* `annotations` - (Optional/Computed) Annotations of the resource (map).
* `labels` - (Optional/Computed) Labels of the resource (map).
                

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `name` - (Computed) The name of the resource.
* `type` - (Computed) The type of the resource.

