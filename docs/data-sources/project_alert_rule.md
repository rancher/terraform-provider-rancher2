---
page_title: "rancher2_project_alert_rule Data Source"
---

# rancher2\_project\_alert\_rule Data Source

Use this data source to retrieve information about a Rancher v2 project alert rule.

## Example Usage

```
data "rancher2_project_alert_rule" "foo" {
  project_id = "<project_id>"
  name = "<project_alert_rule_name>"
}
```

## Argument Reference

* `project_id` - (Required) The project id where create project alert rule (string)
* `name` - (Required) The project alert rule name (string)

## Attributes Reference

* `group_id` - (Computed) The project alert rule alert group ID (string)
* `group_interval_seconds` - (Computed) The project alert rule group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Computed) The project alert rule group wait seconds. Default: `180` (int)
* `inherited` - (Computed) The project alert rule inherited. Default: `true` (bool)
* `metric_rule` - (Computed) The project alert rule metric rule. ConflictsWith: `"pod_rule", "workload_rule"`` (list Maxitems:1)
* `pod_rule` - (Computed) The project alert rule pod rule. ConflictsWith: `"metric_rule", "workload_rule"`` (list Maxitems:1)
* `repeat_interval_seconds` - (Computed) The project alert rule wait seconds. Default: `3600` (int)
* `severity` - (Computed) The project alert rule severity. Supported values : `"critical" | "info" | "warning"`. Default: `critical` (string)
* `workload_rule` - (Computed) The project alert rule workload rule. ConflictsWith: `"metric_rule", "pod_rule"`` (list Maxitems:1)
* `annotations` - (Computed) The project alert rule annotations (map)
* `labels` - (Computed) The project alert rule labels (map)

