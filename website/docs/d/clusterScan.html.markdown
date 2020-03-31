---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster_scan"
sidebar_current: "docs-rancher2-datasource-cluster_scan"
description: |-
  Get information on a Rancher v2 Cluster CIS Scan resource.
---

# rancher2\_cluster\_scan

Use this data source to retrieve information about a Rancher v2 Cluster CIS Scan resource.

## Example Usage

```hcl
data "rancher2_cluster_scan" "foo" {
    cluster_id = <clusterID>
    name = "foo"
}
```

## Argument Reference

* `cluster_id` - (Required) Cluster ID for CIS Scan (string)
* `name` - (Optional/Computed) Name of the cluster Scan (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `run_type` - (Computed) Cluster Scan run type (string)
* `scan_config` - (Computed) Cluster Scan config (bool)
* `scan_type` - (Computed) Cluster Scan type (string)
* `status` - (Computed) Cluster Scan status (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)

## Nested blocks

### `scan_config`

#### Arguments

* `cis_scan_config` - (Computed) Cluster Cis Scan config (List maxitems:1)

#### `cis_scan_config`

##### Arguments

* `debug_master` - (Computed) Debug master. Default: `false` (bool)
* `debug_worker` - (Computed) Debug worker. Default: `false` (bool)
* `override_benchmark_version` - (Computed) Override benchmark version (string)
* `override_skip` - (Computed) Override skip (string)
* `profile` - (Computed) Cis scan profile. Allowed values: `"permissive" (default) || "hardened"` (string)
