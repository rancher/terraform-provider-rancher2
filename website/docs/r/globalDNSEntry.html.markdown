---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_dns_entry"
sidebar_current: "docs-rancher2-resource-global_dns_entry"
description: |-
  Provides a Rancher v2 Global DNS Entry resource. This can be used to create Global DNS Entry for Rancher v2.
---

# rancher2\_global\_dns\_entry

Provides a Rancher v2 Global DNS Entry resource. This can be used to create Global DNS Entry for Rancher v2.

## Example Usage

```hcl
# Create a new rancher2 Global DNS Entry
resource "rancher2_global_dns_entry" "entry" {
	name = "test-entry"
	fqdn = "test-entry.example.com"
	provider_id = rancher2_global_dns_provider.provider.id
	project_ids = ["${data.rancher2_project.proj.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional/Computed/ForceNew) The name of the Global DNS Entry (string)
* `fqdn` - (Required/ForceNew) The Global DNS Entry `route53, alidns, cloudflare` (string)
* `provider_id` - (Required/ForceNew) The provider ID to assign Global DNS Entry (string)
* `project_ids` - (Optional/ForceNew) A list of project_ids to assign to the Global DNS Entry, conflicts with `multi_cluster_app_id` (list(string))
* `multi_cluster_app_id` - (Optional/ForceNew) The MultiCluster App ID to assign to the Global DNS Entry, conflicts with `project_ids` (string)
* `annotations` - (Optional/Computed) Annotations for Global DNS Entry (map)
* `labels` - (Optional/Computed) Labels for Global DNS Entry (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_global_dns_entry` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating Global DNS Entry.
- `update` - (Default `5 minutes`) Used for Global DNS Entry modifications.
- `delete` - (Default `5 minutes`) Used for deleting Global DNS Entry.

## Import

Global DNS Entry can be imported using the Rancher Global DNS Entry ID

```
$ terraform import rancher2_global_dns_entry.foo <global_dns_entry_id>
```

