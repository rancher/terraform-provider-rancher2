---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_dns_provider"
sidebar_current: "docs-rancher2-resource-global_dns_provider"
description: |-
  Provides a Rancher v2 Global DNS Provider resource. This can be used to create Global DNS Providers for Rancher v2.
---

# rancher2\_global\_dns\_provider

Provides a Rancher v2 Global DNS Provider resource. This can be used to create Global DNS Providers for Rancher v2.

## Example Usage

```hcl
# Create a new rancher2 Global DNS Provider
resource "rancher2_global_dns_provider" "dns" {
  name = "foo-test2"
  dns_provider = "route53"
  root_domain = "example.com"

  route53_config {
    access_key = "YYYYYYYYYYYYYYYYYYYY"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    zone_type = "private"
    region = "us-east-1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional/Computed/ForceNew) The name of the Global DNS Provider (string)
* `dns_provider` - (Required/ForceNew) The Global DNS Provider `route53, alidns, cloudflare` (string)
* `root_domain` - (Required/ForceNew) The user ID to assign Global DNS Provider (string)
* `route53_domain` - (Optional/Computed) The Route53 provider config (set)
* `alidns_domain` - (Optional/Computed) The AliDNS provider config (set)
* `cloudflare_domain` - (Optional/Computed) The ClodFlare provider config (set)
* `annotations` - (Optional/Computed) Annotations for Global DNS Provider (map)
* `labels` - (Optional/Computed) Labels for Global DNS Provider (map)

## Nested blocks

### `route_53`

#### Arguments

* `access_key` - (Required) The AWS Access key (string)
* `secret_key` - (Required) The AWS Secret key (string)
* `zone_type` - (Required/ForceNew) The Route53 zone type `public, private` (string)
* `region` - (Required/ForceNew) The AWS Region (string)
* `credentials_path` - (Optional) The AWS credentials path (string)
* `role_arn` - (Optional) The AWS Role ARN (string)

### `cloudflare_config`

#### Arguments
* `api_email` - (Required) The CloudFlare API Email (string)
* `api_key` - (Required) The CloudFlare API Key (string)
* `proxy_setting` - (Optional) CloudFlare Proxy Setting (bool)

### `alidns_config`

#### Arguments
* `access_key` - (Required) The AliDNS Access key (string)
* `secret_key` - (Required) The AliDNS Secret key (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_global_dns_provider` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating Global DNS Providers.
- `update` - (Default `5 minutes`) Used for Global DNS Provider modifications.
- `delete` - (Default `5 minutes`) Used for deleting Global DNS Providers.

## Import

Global DNS Providers can be imported using the Rancher Global DNS Provider ID

```
$ terraform import rancher2_global_dns_provider.foo <global_dns_provider_id>
```

