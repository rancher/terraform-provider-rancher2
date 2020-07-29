---
page_title: "rancher2_setting Data Source"
---

# rancher2\_setting Data Source

Use this data source to retrieve information about a Rancher v2 setting.

## Example Usage

```
data "rancher2_setting" "server-image" {
    name = "server-image"
}
```

## Argument Reference

 * `name` - (Required) The setting name.

## Attributes Reference

 * `value` - the settting's value.
