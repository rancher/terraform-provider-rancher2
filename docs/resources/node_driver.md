---
page_title: "rancher2_node_driver Resource"
---

# rancher2\_node\_driver Resource

Provides a Rancher v2 Node Driver resource. This can be used to create Node Driver for Rancher v2 RKE clusters and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Node Driver
resource "rancher2_node_driver" "foo" {
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

* `active` - (Required) Specify if the node driver state (bool)
* `builtin` - (Required) Specify wheter the node driver is an internal node driver or not (bool)
* `name` - (Required) Name of the node driver (string)
* `url` - (Required) The URL to download the machine driver binary for 64-bit Linux (string)
* `checksum` - (Optional) Verify that the downloaded driver matches the expected checksum (string)
* `description` - (Optional) Description of the node driver (string)
* `external_id` - (Optional) External ID (string)
* `ui_url` - (Optional) The URL to load for customized Add Nodes screen for this driver (string)
* `whitelist_domains` - (Optional) Domains to whitelist for the ui (list)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_node_driver` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node drivers.
- `update` - (Default `10 minutes`) Used for node driver modifications.
- `delete` - (Default `10 minutes`) Used for deleting node drivers.

## Import

Node Driver can be imported using the Rancher Node Driver ID

```
$ terraform import rancher2_node_driver.foo &lt;node_driver_id&gt;
```
