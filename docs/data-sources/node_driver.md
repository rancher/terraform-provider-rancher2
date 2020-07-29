---
page_title: "rancher2_node_driver Data Source"
---

# rancher2\_node\_driver Data Source

Use this data source to retrieve information about a Rancher v2 Node Driver resource. 

## Example Usage

```hcl
data "rancher2_node_driver" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) Name of the node driver (string)
* `url` - (Optional/Computed) The URL to download the machine driver binary for 64-bit Linux (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `active` - (Computed) Specify if the node driver state (bool)
* `builtin` - (Computed) Specify wheter the node driver is an internal cluster driver or not (bool)
* `checksum` - (Computed) Verify that the downloaded driver matches the expected checksum (string)
* `description` - (Computed) Description of the node driver (string)
* `external_id` - (Computed) External ID (string)
* `ui_url` - (Computed) The URL to load for customized Add Node screen for this driver (string)
* `whitelist_domains` - (Computed) Domains to whitelist for the ui (list)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)

