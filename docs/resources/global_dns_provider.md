---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_dns_provider"
sidebar_current: "docs-rancher2-resource-global_dns_provider"
description: |-
  Provides a Rancher V2 Global DNS Provider resource. This can be used to create Global DNS Providers for Rancher V2.
---

# rancher2\_global\_dns\_provider

Provides a Rancher V2 Global DNS Provider resource. This can be used to create Global DNS Providers for Rancher V2. Supported Global DNS Providers: `alidns, cloudflare, route53`

## Example Usage

```hcl
# Create a new rancher2 Global DNS Provider - alidns
resource "rancher2_global_dns_provider" "foo" {
  name = "foo"
  root_domain = "example.com"
  alidns_config {
    access_key = "YYYYYYYYYYYYYYYYYYYY"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
```

```hcl
# Create a new rancher2 Global DNS Provider - cloudflare
resource "rancher2_global_dns_provider" "foo" {
  name = "foo"
  root_domain = "example.com"
  cloudflare_config {
    api_email = "test@test.local"
    api_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    proxy_setting = true
  }
}
```

```hcl
# Create a new rancher2 Global DNS Provider - route53
resource "rancher2_global_dns_provider" "foo" {
  name = "foo"
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
* `root_domain` - (Required) The user ID to assign Global DNS Provider (string)
* `alidns_domain` - (Optional) The AliDNS provider config (list maxitems:1)
* `cloudflare_domain` - (Optional) The ClodFlare provider config (list maxitems:1)
* `route53_domain` - (Optional) The Route53 provider config (list maxitems:1)
* `annotations` - (Optional/Computed) Annotations for Global DNS Provider (map)
* `labels` - (Optional/Computed) Labels for Global DNS Provider (map)

## Nested blocks

### `alidns_config`

#### Arguments
* `access_key` - (Required) The AliDNS Access key (string)
* `secret_key` - (Required) The AliDNS Secret key (string)

### `cloudflare_config`

#### Arguments
* `api_email` - (Required) The CloudFlare API Email (string)
* `api_key` - (Required) The CloudFlare API Key (string)
* `proxy_setting` - (Optional) CloudFlare Proxy Setting. Default: `false` (bool)

### `route_53`

#### Arguments

* `access_key` - (Required) The AWS Access key (string)
* `secret_key` - (Required) The AWS Secret key (string)
* `credentials_path` - (Optional) The AWS credentials path. Default: `"/.aws"` (string)
* `region` - (Optional) The AWS Region. Default: `"us-west-2"` (string)
* `role_arn` - (Optional) The AWS Role ARN (string)
* `zone_type` - (Optional) The Route53 zone type `public, private`. Default: `"public"` (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `dns_provider` - (Computed) The Global DNS Provider `alidns, cloudflare, route53` (string)

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

