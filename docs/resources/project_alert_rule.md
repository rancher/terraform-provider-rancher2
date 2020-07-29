---
page_title: "rancher2_project_alert_rule Resource"
---

# rancher2\_project\_alert\_rule Resource

Provides a Rancher v2 Project Alert Rule resource. This can be used to create Project Alert Rule for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Project
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "<cluster_id>"
  description = "Terraform project "
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
# Create a new Rancher2 Project Alert Group
resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group"
  project_id = rancher2_project.foo.id
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
# Create a new Rancher2 Project Alert Rule
resource "rancher2_project_alert_rule" "foo" {
  project_id = rancher2_project_alert_group.foo.project_id
  group_id = rancher2_project_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project id where create project alert rule (string)
* `group_id` - (Required) The project alert rule alert group ID (string)
* `name` - (Required) The project alert rule name (string)
* `group_interval_seconds` - (Optional) The project alert rule group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Optional) The project alert rule group wait seconds. Default: `180` (int)
* `inherited` - (Optional) The project alert rule inherited. Default: `true` (bool)
* `metric_rule` - (Optional) The project alert rule metric rule. ConflictsWith: `"pod_rule", "workload_rule"`` (list Maxitems:1)
* `pod_rule` - (Optional) The project alert rule pod rule. ConflictsWith: `"metric_rule", "workload_rule"`` (list Maxitems:1)
* `repeat_interval_seconds` - (Optional) The project alert rule wait seconds. Default: `3600` (int)
* `severity` - (Optional) The project alert rule severity. Supported values : `"critical" | "info" | "warning"`. Default: `critical` (string)
* `workload_rule` - (Optional) The project alert rule workload rule. ConflictsWith: `"metric_rule", "pod_rule"`` (list Maxitems:1)
* `annotations` - (Optional/Computed) The project alert rule annotations (map)
* `labels` - (Optional/Computed) The project alert rule labels (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `metric_rule`

#### Arguments

* `duration` - (Required) Metric rule duration (string)
* `expression` - (Required) Metric rule expression (string)
* `threshold_value` - (Required) Metric rule threshold value (float64)
* `comparison` - (Optional) Metric rule comparison. Supported values : `"equal" | "greater-or-equal" | "greater-than" | "less-or-equal" | "less-than" | "not-equal" | "has-value"`. Default: `equal`  (string)
* `description` - (Optional) Metric rule description (string)

### `pod_rule`

#### Arguments

* `pod_id` - (Required) Pod ID (string)
* `condition` - (Optional) Pod rule condition. Supported values : `"notrunning" | "notscheduled" | "restarts"`. Default: `notrunning` (string)
* `restart_interval_seconds` - (Optional) Pod rule restart interval seconds. Default: `300` (int)
* `restart_times` - (Optional) Pod rule restart times. Default: `3`  (int)

### `workload_rule`

#### Arguments

* `available_percentage` - (Optional) Workload rule available percentage. Default: `70` (int)
* `selector` - (Optional) Workload rule selector (map)
* `workload_id` - (Optional) Workload ID (string)

## Timeouts

`rancher2_project_alert_rule` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating project alert rules.
- `update` - (Default `10 minutes`) Used for project alert rule modifications.
- `delete` - (Default `10 minutes`) Used for deleting project alert rules.

## Import

Project Alert Rule can be imported using the Rancher project alert rule ID

```
$ terraform import rancher2_project_alert_rule.foo &lt;project_alert_rule_id&gt;
```
