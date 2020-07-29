---
page_title: "rancher2_cluster_alert_rule Resource"
---

# rancher2\_cluster\_alert\_rule Resource

Provides a Rancher v2 Cluster Alert Rule resource. This can be used to create Cluster Alert Rule for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Cluster Alert Group
resource "rancher2_cluster_alert_group" "foo" {
  cluster_id = "<cluster_id>"
  name = "foo"
  description = "Terraform cluster alert group"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
# Create a new Rancher2 Cluster Alert Rule
resource "rancher2_cluster_alert_rule" "foo" {
  cluster_id = rancher2_cluster_alert_group.foo.cluster_id
  group_id = rancher2_cluster_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster id where create cluster alert rule (string)
* `group_id` - (Required) The cluster alert rule alert group ID (string)
* `name` - (Required) The cluster alert rule name (string)
* `event_rule` - (Optional) The cluster alert rule event rule. ConflictsWith: `"metric_rule", "node_rule", "system_service_rule"`` (list Maxitems:1)
* `group_interval_seconds` - (Optional) The cluster alert rule group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Optional) The cluster alert rule group wait seconds. Default: `180` (int)
* `inherited` - (Optional) The cluster alert rule inherited. Default: `true` (bool)
* `metric_rule` - (Optional) The cluster alert rule metric rule. ConflictsWith: `"event_rule", "node_rule", "system_service_rule"`` (list Maxitems:1)
* `node_rule` - (Optional) The cluster alert rule node rule. ConflictsWith: `"event_rule", "metric_rule", "system_service_rule"`` (list Maxitems:1)
* `repeat_interval_seconds` - (Optional) The cluster alert rule wait seconds. Default: `3600` (int)
* `severity` - (Optional) The cluster alert rule severity. Supported values : `"critical" | "info" | "warning"`. Default: `critical` (string)
* `system_service_rule` - (Optional) The cluster alert rule system service rule. ConflictsWith: `"event_rule", "metric_rule", "node_rule"` (list Maxitems:1)
* `annotations` - (Optional/Computed) The cluster alert rule annotations (map)
* `labels` - (Optional/Computed) The cluster alert rule labels (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `event_rule`

#### Arguments

* `resource_kind` - (Required) Resource kind. Supported values : `"DaemonSet" | "Deployment" | "Node" | "Pod" | "StatefulSet"` (string)
* `event_type` - (Optional) Event type. Supported values : `"Warning" | "Normal"`. Default: `Warning` (string)

### `metric_rule`

#### Arguments

* `duration` - (Required) Metric rule duration (string)
* `expression` - (Required) Metric rule expression (string)
* `threshold_value` - (Required) Metric rule threshold value (float64)
* `comparison` - (Optional) Metric rule comparison. Supported values : `"equal" | "greater-or-equal" | "greater-than" | "less-or-equal" | "less-than" | "not-equal" | "has-value"`. Default: `equal`  (string)
* `description` - (Optional) Metric rule description (string)

### `node_rule`

#### Arguments

* `cpu_threshold` - (Optional) Node rule cpu threshold. Default: `70` (int)
* `condition` - (Optional) Node rule condition. Supported values : `"cpu" | "mem" | "notready"`. Default: `notready` (string)
* `mem_threshold` - (Optional) Node rule mem threshold. Default: `70` (int)
* `node_id` - (Optional) Node ID (string)
* `selector` - (Optional) Node rule selector (map)

### `system_service_rule`

#### Arguments

* `condition` - (Optional) System service rule condition. Supported values : `"controller-manager" | "etcd" | "scheduler"`. Default: `scheduler` (string)

## Timeouts

`rancher2_cluster_alert_rule` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cluster alert rules.
- `update` - (Default `10 minutes`) Used for cluster alert rule modifications.
- `delete` - (Default `10 minutes`) Used for deleting cluster alert rules.

## Import

Cluster Alert Rule can be imported using the Rancher cluster alert rule ID

```
$ terraform import rancher2_cluster_alert_rule.foo &lt;CLUSTER_ALERT_RULE_ID&gt;
```
