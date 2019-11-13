---
layout: "rancher2"
page_title: "Rancher2: rancher2_bootstrap"
sidebar_current: "docs-rancher2-resource-bootstrap"
description: |-
  Provides a Rancher v2 bootstrap resource. This can be used to bootstrap Rancher v2 environments and output information.
---

# rancher2\_bootstrap

Provides a Rancher v2 bootstrap resource. This can be used to bootstrap Rancher v2 environments and output information. It just works if `bootstrap` provider config is added to the .tf file. More info at [rancher2 provider](../index.html)

This resource bootstraps a Rancher system by performing the following tasks:
- Updates the default admin password, provided by setting `password` or generating a random one.
- Sets `server-url` setting, based on `api_url`.
- Sets `telemetry-opt` setting.
- Creates a token for admin user with concrete TTL.

Rancher2 admin password can be updated after the initial run of terraform by setting `password` field and applying this resource again.

Rancher2 admin `token` can also be regenerated if `token_update` is set to true. Refresh resource function will check if token is expired. If it is expired, `token_update` will be set to true to force token regeneration on next `terraform apply`.

Login to Rancher2 is done by trying to use `token` first. If it fails, it uses admin `current_password`. If admin password has been changed outside of terraform and the terraform `token` is expired, `current_password` field can be specified to allow terraform to manage admin password and token again.

## Example Usage

```hcl
# Provider bootstrap config
provider "rancher2" {
  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}

# Create a new rancher2_bootstrap
resource "rancher2_bootstrap" "admin" {
  password = "blahblah"
  telemetry = true
}
```

```hcl
# Provider bootstrap config with alias
provider "rancher2" {
  alias = "bootstrap"

  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}

# Create a new rancher2_bootstrap using bootstrap provider config
resource "rancher2_bootstrap" "admin" {
  provider = "rancher2.bootstrap"

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
