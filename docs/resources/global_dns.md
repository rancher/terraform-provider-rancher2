---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_dns"
sidebar_current: "docs-rancher2-resource-global_dns"
description: |-
  Provides a Rancher V2 Global DNS resource. This can be used to create Global DNS records for Rancher V2.
---

# rancher2\_global\_dns

Provides a Rancher V2 Global DNS resource. This can be used to create Global DNS records for Rancher V2.

## Example Usage

```hcl
# Create a new rancher2 Global DNS Provider
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
# Create a new rancher2 Global DNS using project IDs
resource "rancher2_global_dns" "foo" {
	name = "foo"
	fqdn = "foo.example.com"
	provider_id = rancher2_global_dns_provider.foo.id
	project_ids = ["project1", "project2"]
}
```

```hcl
# Create a new rancher2 Global DNS Provider
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
# Create a new rancher2 Global DNS using MultiClusterApp ID
resource "rancher2_global_dns" "foo" {
	name = "foo"
	fqdn = "foo.example.com"
	provider_id = rancher2_global_dns_provider.foo.id
	multi_cluster_app_id = "<MCA_ID>"
}
```

## Argument Reference

The following arguments are supported:

* `fqdn` - (Required) The Global DNS record (string)
* `provider_id` - (Required) The Global DNS provider ID to use (string)
* `multi_cluster_app_id` - (Optional) The MultiCluster App ID to assign to the Global DNS. Conflicts with `project_ids` (string)
* `name` - (Optional/Computed/ForceNew) The name of the Global DNS (string)
* `project_ids` - (Optional) A list of project_ids to assign to the Global DNS. Conflicts with `multi_cluster_app_id` (list(string))
* `ttl` - (Optional) TTL in seconds for DNS record. Default: `300` (int)
* `annotations` - (Optional/Computed) Annotations for Global DNS (map)
* `labels` - (Optional/Computed) Labels for Global DNS (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_global_dns_entry` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating Global DNS
- `update` - (Default `5 minutes`) Used for Global DNS modifications
- `delete` - (Default `5 minutes`) Used for deleting Global DNS

## Import

Global DNS Entry can be imported using the Rancher Global DNS ID

```
$ terraform import rancher2_global_dns_entry.foo <global_dns_id>
```
