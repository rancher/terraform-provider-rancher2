---
page_title: "rancher2_cluster_alert_rule Data Source"
---

# rancher2\_cluster\_alert\_rule Data Source

Use this data source to retrieve information about a Rancher v2 cluster alert rule.

## Example Usage

```
data "rancher2_cluster_alert_rule" "foo" {
  cluster_id = "<cluster_id>"
  name = "<cluster_alert_rule_name>"
}
```

## Argument Reference

* `cluster_id` - (Required) The cluster id where create cluster alert rule (string)
* `name` - (Required) The cluster alert rule name (string)

## Attributes Reference

* `group_id` - (Computed) The cluster alert rule alert group ID (string)
* `event_rule` - (Computed) The cluster alert rule event rule. ConflictsWith: `"metric_rule", "node_rule", "system_service_rule"` (list Maxitems:1)
* `group_interval_seconds` - (Computed) The cluster alert rule group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Computed) The cluster alert rule group wait seconds. Default: `180` (int)
* `inherited` - (Computed) The cluster alert rule inherited. Default: `true` (bool)
* `metric_rule` - (Computed) The cluster alert rule metric rule. ConflictsWith: `"event_rule", "node_rule", "system_service_rule"`` (list Maxitems:1)
* `node_rule` - (Computed) The cluster alert rule node rule. ConflictsWith: `"event_rule", "metric_rule", "system_service_rule"`` (list Maxitems:1)
* `repeat_interval_seconds` - (Optional) The cluster alert rule wait seconds. Default: `3600` (int)
* `severity` - (Computed) The cluster alert rule severity. Supported values : `"critical" | "info" | "warning"`. Default: `critical` (string)
* `system_service_rule` - (Computed) The cluster alert rule system service rule. ConflictsWith: `"event_rule", "metric_rule", "node_rule"`` (list Maxitems:1)
* `annotations` - (Computed) The cluster alert rule annotations (map)
* `labels` - (Computed) The cluster alert rule labels (map)

