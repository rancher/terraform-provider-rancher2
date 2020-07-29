---
page_title: "Rancher2: rancher2_cluster_driver Resource"
---

# rancher2\_cluster\_driver Resource

Provides a Rancher v2 Cluster Driver resource. This can be used to create Cluster Driver for Rancher v2.2.x Kontainer Engine clusters and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Cluster Driver
resource "rancher2_cluster_driver" "foo" {
    active = true
    builtin = false
    checksum = "0x0"
    description = "Foo description"
    external_id = "foo_external"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
    whitelist_domains = ["*.foo.com"]
}
```

## Argument Reference

The following arguments are supported:

* `active` - (Required) Specify the cluster driver state (bool)
* `builtin` - (Required) Specify whether the cluster driver is an internal cluster driver or not (bool)
* `name` - (Required) Name of the cluster driver (string)
* `url` - (Required) The URL to download the machine driver binary for 64-bit Linux (string)
* `actual_url` - (Optional) Actual url of the cluster driver (string)
* `checksum` - (Optional) Verify that the downloaded driver matches the expected checksum (string)
* `ui_url` - (Optional) The URL to load for customized Add Clusters screen for this driver (string)
* `whitelist_domains` - (Optional) Domains to whitelist for the ui (list)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_cluster_driver` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cluster drivers.
- `update` - (Default `10 minutes`) Used for cluster driver modifications.
- `delete` - (Default `10 minutes`) Used for deleting cluster drivers.

## Import

Cluster Driver can be imported using the Rancher Cluster Driver ID

```
$ terraform import rancher2_cluster_driver.foo &lt;CLUSTER_DRIVER_ID&gt;
```
