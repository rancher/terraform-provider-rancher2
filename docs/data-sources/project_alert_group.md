---
page_title: "rancher2_project_alert_group Data Source"
---

# rancher2\_project\_alert\_group Data Source

Use this data source to retrieve information about a Rancher v2 project alert group.

## Example Usage

```
data "rancher2_project_alert_group" "foo" {
  project_id = "<project_id>"
  name = "<project_alert_group_name>"
}
```

## Argument Reference

* `project_id` - (Required) The project id where create project alert group (string)
* `name` - (Required) The project alert group name (string)

## Attributes Reference

* `description` - (Computed) The project alert group description (string)
* `group_interval_seconds` - (Computed) The project alert group interval seconds. Default: `180` (int)
* `group_wait_seconds` - (Computed) The project alert group wait seconds. Default: `180` (int)
* `recipients` - (Computed) The project alert group recipients (list)
* `repeat_interval_seconds` - (Computed) The project alert group wait seconds. Default: `3600` (int)
* `annotations` - (Computed) The project alert group annotations (map)
* `labels` - (Computed) The project alert group labels (map)

