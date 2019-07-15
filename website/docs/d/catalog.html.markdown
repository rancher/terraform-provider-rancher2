---
layout: "rancher2"
page_title: "Rancher2: rancher2_catalog"
sidebar_current: "docs-rancher2-datasource-catalog"
description: |-
  Get information on a Rancher v2 catalog.
---

# rancher2\_catalog

Use this data source to retrieve information about a Rancher v2 catalog.

## Example Usage

```
data "rancher2_catalog" "library" {
    name = "catalog"
}
```

## Argument Reference

* `name` - (Required) The catalog name.
* `scope` - (Optional) The scope of the catalog. `cluster`, `global`, and `project` are supported. Default `global` (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `branch` - (Computed) The branch of the catalog repo to use (string)
* `cluster_id` - (Computed) The cluster id of the catalog (string)
* `description` - (Computed) A catalog description (string)
* `kind` - (Computed) The kind of the catalog. Just helm by the moment (string)
* `password` - (Computed/Sensitive) The password to access the catalog if needed (string)
* `project_id` - (Computed) The project id of the catalog (string)
* `username` - (Computed/Sensitive) The username to access the catalog if needed (string)
* `url` - (Computed) The url of the catalog repo (string)
* `annotations` - (Computed) Annotations for the catalog (map)
* `labels` - (Computed) Labels for the catalog (map)

