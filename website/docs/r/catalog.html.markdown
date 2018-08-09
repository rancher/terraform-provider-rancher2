---
layout: "cattle"
page_title: "Cattle: cattle_catalog"
sidebar_current: "docs-cattle-resource-catalog"
description: |-
  Provides a Cattle Catalog resource. This can be used to create catalogs for rancher v2 environments and retrieve their information.
---

# cattle\_certificate

Provides a Cattle Catalog resource. This can be used to create catalogs for rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Cattle Catalog
resource "cattle_catalog" "foo" {
  name = "foo"
  url = "https://<CATALOG_URL>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the catalog.
* `url` - (Required) The url of the catalog repo.
* `description` - (Optional) A catalog description.
* `kind` - (Optional) The kind of the catalog. Just helm by the moment.
* `branch` - (Optional) The branch of the catalog repo to use.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Catalogs can be imported using the rancher Catalog ID.

```
$ terraform import cattle_catalog.foo <catalog_id>
```

