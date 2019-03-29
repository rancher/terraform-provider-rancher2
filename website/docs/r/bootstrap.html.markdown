---
layout: "rancher2"
page_title: "Rancher2: rancher2_bootstrap"
sidebar_current: "docs-rancher2-resource-bootstrap"
description: |-
  Provides a Rancher v2 bootstrap resource. This can be used to bootstrap rancher v2 environments and output information.
---

# rancher2\_bootstrap

Provides a Rancher v2 bootstrap resource. This can be used to bootstrap rancher v2 environments and output information.

This resource is indeed to bootstrap a rancher system doing following tasks:
- Update default admin password, provided by `password` or generating a random one.
- Set `server-url` setting, based on `api_url`.
- Set `telemetry-opt` setting.
- Create a token for admin user with concrete TTL.

It just works if `bootstrap = true` is added to the provider configuration or exporting env variable `RANCHER_BOOTSTRAP=true`. In this mode, `token_key` or `access_key` and `secret_key` can not be provided.

Rancher2 admin password could be updated setting `password` field and applying this resource again. 

Rancher2 admin `token` could also be updated if `token_update` is set to true. Refresh resource function will check if token is expired. If it's expired, `token_update` will be set to true to force token regeneration on next `terraform apply`.

Login to rancher2 is done using `token` first and if fails, using admin `current_password`. If admin password has been changed from other methods and terraform token is expired, `current_password` field could be especified to recover terraform configuration and reset admin password and token.

## Example Usage

```hcl
# Provider config
provider "rancher2" {
  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}

# Create a new rancher2 Bootstrap
resource "rancher2_bootstrap" "admin" {
  password = "blahblah"
  telemetry = true
}
```

## Argument Reference

The following arguments are supported:

* `current_password` - (Optional/computed/sensitive) Current password for Admin user. Just needed for recover if admin password has been changed from other resources and token is expired (string)
* `password` - (Optional/computed/sensitive) Password for Admin user or random generated if empty (string)
* `telemetry` - (Optional) Send telemetry anonymous data. Default: `false` (bool)
* `token_ttl` - (Optional) TTL in seconds for generated admin token. Default: `0`  (int)
* `token_update` - (Optional) Regenerate admin token. Default: `false` (bool)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `token` - (Computed) Generated API token for Admin User (string)
* `token_id` - (Computed) Generated API token id for Admin User (string)
* `url` - (Computed) URL set as server-url (string)
* `user` - (Computed) Admin username (string)
* `temp_token` - (Computed) Generated API temporary token as helper. Should be empty (string)
* `temp_token_id` - (Computed) Generated API temporary token id as helper. Should be empty (string)
