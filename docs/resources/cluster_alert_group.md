---
page_title: "rancher2_cluster_alert_group Resource"
---

# rancher2\_cluster\_alert\_group Resource

Provides a Rancher v2 Cluster Alert Group resource. This can be used to create Cluster Alert Group for Rancher v2 environments and retrieve their information.

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
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster id where create cluster alert group (string)
* `name` - (Required) The cluster alert group name (string)
* `description` - (Optional) The cluster alert group description (string)
* `group_interval_seconds` - (Optional) The cluster alert group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Optional) The cluster alert group wait seconds. Default: `180` (int)
* `recipients` - (Optional) The cluster alert group recipients (list)
* `repeat_interval_seconds` - (Optional) The cluster alert group wait seconds. Default: `3600` (int)
* `annotations` - (Optional/Computed) The cluster alert group annotations (map)
* `labels` - (Optional/Computed) The cluster alert group labels (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `recipients`

#### Arguments

* `notifier_id` - (Required) Recipient notifier ID (string)
* `recipient` - (Optional/Computed) Recipient (string)
* `default_recipient` - (Optional) Use notifier default recipient, overriding `recipient` argument if set.  Default: `false` (bool)

#### Attributes

* `notifier_type` - (Computed) Recipient notifier ID. Supported values : `"dingtalk" | "msteams" | "pagerduty" | "slack" | "email" | "webhook" | "wechat"` (string)

## Timeouts

`rancher2_cluster_alert_group` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cluster alert groups.
- `update` - (Default `10 minutes`) Used for cluster alert group modifications.
- `delete` - (Default `10 minutes`) Used for deleting cluster alert groups.

## Import

Cluster Alert Group can be imported using the Rancher cluster alert group ID

```
$ terraform import rancher2_cluster_alert_group.foo &lt;CLUSTER_ALERT_GROUP_ID&gt;
```
