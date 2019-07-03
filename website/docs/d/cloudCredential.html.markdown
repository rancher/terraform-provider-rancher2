---
layout: "rancher2"
page_title: "Rancher2: rancher2_cloud_credential"
sidebar_current: "docs-rancher2-datasource-cloud_credential"
description: |-
  Get information on a Rancher v2 Cloud Credentials.
---

# rancher2\_cloud\_credential

Use this data source to retrieve information about a Rancher v2 Cloud Credential.

## Example Usage

```
data "rancher2_cloud_credential" "test" {
    name = "test"
}
```

## Argument Reference

 * `name` - (Required) The Cloud Credential name.

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `annotations` - (Computed) Annotations for the Cloud Credential (map)
* `labels` - (Computed) Labels for the Cloud Credential (map)
