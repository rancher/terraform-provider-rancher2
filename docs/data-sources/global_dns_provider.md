---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_dns_provider"
sidebar_current: "docs-rancher2-datasource-global_dns_provider"
description: |-
  Get information on a Rancher v2 global DNS provider.
---

# rancher2\_global\_dns\_provider

Use this data source to retrieve information about a Rancher v2 global DNS provider.

## Example Usage

```
data "rancher2_global_dns_provider" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The name of the global DNS provider (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the global DNS provider (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)