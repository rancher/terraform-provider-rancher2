---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_driver"
sidebar_current: "docs-rancher2-resource-node_driver"
description: |-
  Provides a Rancher v2 Node Driver resource. This can be used to create Node Driver for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_node\_driver

Provides a Rancher v2 Node Driver resource. This can be used to create Node Driver for rancher v2 rke clusters and retrieve their information.

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

* `active` - (Required) Specify if the node driver state.
* `builtin` - (Required) Specify wheter the node driver is an internal node driver or not.
* `checksum` - (Optional) Verify that the downloaded driver matches the expected checksum.
* `description` - (Optional) Description of the node driver.
* `external_id` - (Optional) External ID
* `name` - (Required) Name of the node driver.
* `ui_url` - (Optional) The URL to load for customized Add Nodes screen for this driver.
* `url` - (Required) The URL to download the machine driver binary for 64-bit Linux.
* `whitelist_domains` - (Optional) Domains to whitelist for the ui.
* `annotations` - (Optional/Computed) Annotations of the resource (map).
* `labels` - (Optional/Computed) Labels of the resource (map).

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Node Driver can be imported using the rancher Node Driver ID

```
$ terraform import rancher2_node_driver.foo <node_driver_id>
```
