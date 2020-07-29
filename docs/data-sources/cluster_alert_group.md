---
page_title: "rancher2_cluster_alert_group Data Source"
---

# rancher2\_cluster\_alert\_group Data Source

Use this data source to retrieve information about a Rancher v2 cluster alert group.

## Example Usage

```
data "rancher2_cluster_alert_group" "foo" {
  cluster_id = "<cluster_id>"
  name = "<cluster_alert_group_name>"
}
```

## Argument Reference

* `cluster_id` - (Required) The cluster id where create cluster alert group (string)
* `name` - (Required) The cluster alert group name (string)

## Attributes Reference

* `description` - (Computed) The cluster alert group description (string)
* `group_interval_seconds` - (Computed) The cluster alert group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Computed) The cluster alert group wait seconds. Default: `180` (int)
* `recipients` - (Computed) The cluster alert group recipients (list)
* `repeat_interval_seconds` - (Computed) The cluster alert group wait seconds. Default: `3600` (int)
* `annotations` - (Computed) The cluster alert group annotations (map)
* `labels` - (Computed) The cluster alert group labels (map)

