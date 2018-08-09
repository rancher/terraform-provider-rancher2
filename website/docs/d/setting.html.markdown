---
layout: "cattle"
page_title: "Cattle: cattle_setting"
sidebar_current: "docs-cattle-datasource-setting"
description: |-
  Get information on a Rancher v2 setting.
---

# cattle\_setting

Use this data source to retrieve information about a Rancher v2 setting.

## Example Usage

```
data "cattle_setting" "server-image" {
    name = "server-image"
}
```

## Argument Reference

 * `name` - (Required) The setting name.

## Attributes Reference

 * `value` - the settting's value.
