---
page_title: "rancher2_bootstrap Resource"
---

# rancher2\_bootstrap Resource

Provides a Rancher v2 bootstrap resource. This can be used to bootstrap Rancher v2 environments and output information. It just works if `bootstrap` provider config is added to the .tf file. More info at [rancher2 provider](../index.html)

This resource bootstraps a Rancher system by performing the following tasks:
- Updates the default admin password, provided by setting `password` or generating a random one.
- Sets `server-url` setting, based on `api_url`.
- Sets `telemetry-opt` setting.
- Creates a token for admin user with concrete TTL.

**Note** Starting from Rancher v2.6.0, the Rancher2 installation is setting a random initial admin password by default. To specify the initial password during rancher2 installation, helm chart [`bootstrapPassword`](https://github.com/rancher/rancher/blob/release/v2.6/chart/values.yaml#L157) value for HA installation or docker env variable [`CATTLE_BOOTSTRAP_PASSWORD`](https://github.com/rancher/rancher/blob/release/v2.6/chart/templates/deployment.yaml#L135) for single node installation can be used. To properly use this resource for Rancher v2.6.0 and above, set the `initial_password` argument to the password generated or set during installation.

Rancher2 admin password can be updated after the initial run of terraform by setting `password` field and applying this resource again.

Rancher2 admin `token` can also be regenerated if `token_update` is set to true. Refresh resource function will check if token is expired. If it is expired, `token_update` will be set to true to force token regeneration on next `terraform apply`.

To login Rancher2, the provider tries until success using `token`, then `current_password` and then `initial_password`. If the admin password has been changed outside of terraform and the `token` is expired, the login will fails and the resource will be regenerated. To recover the bootstrap resource, set `initial_password` argument to the proper password and apply.

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
# Provider bootstrap config
provider "rancher2" {
  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}

# Create a new rancher2_bootstrap for Rancher v2.6.0 and above
resource "rancher2_bootstrap" "admin" {
  initial_password = "<INSTALL_PASSWORD>"
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

* `initial_password` - (Optional/Computed/Sensitive) Initial password for Admin user. Default: `admin` (string)
* `password` - (Optional/Computed/Sensitive) Password for Admin user or random generated if empty (string)
* `telemetry` - (Optional) Send telemetry anonymous data. Default: `false` (bool)
* `token_ttl` - (Optional) TTL in seconds for generated admin token. Default: `0`  (int)
* `token_update` - (Optional) Regenerate admin token. Default: `false` (bool)
* `ui_default_landing` - (Optional) Default UI landing for k8s clusters. Available options: `ember` (cluster manager ui)  and `vue` (cluster explorer ui). Default: `ember` (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `current_password` - (Computed/Sensitive) Current password for Admin user (string)
* `token` - (Computed) Generated API token for Admin User (string)
* `token_id` - (Computed) Generated API token id for Admin User (string)
* `url` - (Computed) URL set as server-url (string)
* `user` - (Computed) Admin username (string)
* `temp_token` - (Computed) Generated API temporary token as helper. Should be empty (string)
* `temp_token_id` - (Computed) Generated API temporary token id as helper. Should be empty (string)
