---
page_title: "rancher2_project_alert_group Resource"
---

# rancher2\_project\_alert\_group Resource

Provides a Rancher v2 Project Alert Group resource. This can be used to create Project Alert Group for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Project Alert Group
resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert group"
  project_id = "<project_id>"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The project alert group name (string)
* `project_id` - (Required) The project id where create project alert group (string)
* `description` - (Optional) The project alert group description (string)
* `group_interval_seconds` - (Optional) The project alert group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Optional) The project alert group wait seconds. Default: `180` (int)
* `recipients` - (Optional) The project alert group recipients (list)
* `repeat_interval_seconds` - (Optional) The project alert group wait seconds. Default: `3600` (int)
* `annotations` - (Optional/Computed) The project alert group annotations (map)
* `labels` - (Optional/Computed) The project alert group labels (map)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `recipients`

#### Arguments

* `notifier_id` - (Required) Recipient notifier ID (string)
* `recipient` - (Optional/Computed) Recipient (string)

#### Attributes

* `notifier_type` - (Computed) Recipient notifier ID. Supported values : `"pagerduty" | "slack" | "email" | "webhook" | "wechat"` (string)

## Timeouts

`rancher2_project_alert_group` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating project alert groups.
- `update` - (Default `10 minutes`) Used for project alert group modifications.
- `delete` - (Default `10 minutes`) Used for deleting project alert groups.

## Import

Project Alert Group can be imported using the Rancher project alert group ID

```
$ terraform import rancher2_project_alert_group.foo &lt;project_alert_group_id&gt;
```
