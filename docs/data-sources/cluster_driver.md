---
page_title: "rancher2_cluster_driver Data Source"
---

# rancher2\_cluster\_driver Data Source

Use this data source to retrieve information about a Rancher v2 Cluster Driver resource.

## Example Usage

```hcl
data "rancher2_cluster_driver" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) Name of the cluster driver (string)
* `url` - (Optional/Computed) The URL to download the machine driver binary for 64-bit Linux (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `active` - (Computed) Specify if the cluster driver state (bool)
* `builtin` - (Computed) Specify whether the cluster driver is an internal cluster driver or not (bool)
* `actual_url` - (Computed) Actual url of the cluster driver (string)
* `checksum` - (Computed) Verify that the downloaded driver matches the expected checksum (string)
* `ui_url` - (Computed) The URL to load for customized Add Clusters screen for this driver (string)
* `whitelist_domains` - (Computed) Domains to whitelist for the ui (list)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)
