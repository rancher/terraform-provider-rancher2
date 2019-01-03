---
layout: "rancher2"
page_title: "Rancher2: rancher2_catalog"
sidebar_current: "docs-rancher2-resource-catalog"
description: |-
  Provides a Rancher v2 Catalog resource. This can be used to create catalogs for rancher v2 environments and retrieve their information.
---

# rancher2\_catalog

Provides a Rancher v2 Catalog resource. This can be used to create catalogs for rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Catalog
resource "rancher2_catalog" "foo" {
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

### Timeouts

`rancher2_catalog` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating catalogs.
- `update` - (Default `10 minutes`) Used for catalog modifications.
- `delete` - (Default `10 minutes`) Used for deleting catalogs.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Catalogs can be imported using the rancher Catalog ID.

```
$ terraform import rancher2_catalog.foo <catalog_id>
```

