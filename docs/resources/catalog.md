---
page_title: "Rancher2: rancher2_catalog Resource"
---

# rancher2\_catalog Resource

Provides a Rancher v2 Catalog resource. This can be used to create cluster, global and/or project catalogs for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Global Catalog
resource "rancher2_catalog" "foo-global" {
  name = "foo-global"
  url = "https://<CATALOG_URL>"
}
# Create a new Rancher2 Cluster Catalog
resource "rancher2_catalog" "foo-cluster" {
  name = "foo-cluster"
  url = "https://<CATALOG_URL>"
  scope = "cluster"
}
# Create a new Rancher2 Project Catalog
resource "rancher2_catalog" "foo-project" {
  name = "foo-project"
  url = "https://<CATALOG_URL>"
  scope = "project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the catalog (string)
* `url` - (Required) The url of the catalog repo (string)
* `branch` - (Optional) The branch of the catalog repo to use. Default `master` (string)
* `cluster_id` - (Optional/ForceNew) The cluster id of the catalog. Mandatory if `scope = cluster` (string)
* `description` - (Optional) A catalog description (string)
* `kind` - (Optional) The kind of the catalog. Just helm by the moment (string)
* `password` - (Optional/Sensitive) The password to access the catalog if needed (string)
* `project_id` - (Optional/ForceNew) The project id of the catalog. Mandatory if `scope = project` (string)
* `refresh` - (Optional) Catalog will wait for refresh after tf creation and on every tf read. Default `false` (bool)
* `scope` - (Optional) The scope of the catalog. `cluster`, `global`, and `project` are supported. Default `global` (string)
* `username` - (Optional/Sensitive) The username to access the catalog if needed (string)
* `version` - (Optional/ForceNew/Computed) Helm version for the catalog. Available options: `helm_v2` and `helm_v3` (string)
* `annotations` - (Optional/Computed) Annotations for the catalog (map)
* `labels` - (Optional/Computed) Labels for the catalog (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_catalog` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating catalogs.
- `update` - (Default `10 minutes`) Used for catalog modifications.
- `delete` - (Default `10 minutes`) Used for deleting catalogs.

## Import

Catalogs can be imported using the Rancher Catalog ID and its scope.

```
$ terraform import rancher2_catalog.foo &lt;SCOPE&gt;.&lt;CATALOG_ID_ID&gt;
```
