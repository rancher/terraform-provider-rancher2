---
layout: "rancher"
page_title: "Provider: Cattle"
sidebar_current: "docs-cattle-index"
description: |-
  The Cattle provider is used to interact with Rancher v2 container platforms.
---

# Cattle Provider

The Cattle provider is used to interact with the
resources supported by Rancher v2. The provider needs to be configured
with the URL of the Rancher server at minimum and API credentials if
access control is enabled on the server.

## Example Usage

```hcl
# Configure the Cattle provider
provider "cattle" {
  api_url    = "https://rancher.my-domain.com"
  access_key = "${var.cattle_access_key}"
  secret_key = "${var.rancher_secret_key}"
}
```

## Argument Reference

The following arguments are supported:

* `api_url` - (Required) Rancher API url. It must be provided, but it can also be sourced from the `CATTLE_URL` environment variable.
* `access_key` - (Optional) Rancher API access key. It can also be sourced from the `CATTLE_ACCESS_KEY` environment variable.
* `secret_key` - (Optional) Rancher API secret key. It can also be sourced from the `CATTLE_SECRET_KEY` environment variable.
* `token_key` - (Optional) Rancher API toke key. It can also be sourced from the `CATTLE_TOKEN_KEY` environment variable. Could be used instead `access_key` and `secret_key`.
* `cacerts` - CA certificates used to sign rancher server tls certificates. Mandatory if self signed.
* `config` - Path to the Rancher client cli.json config file.
