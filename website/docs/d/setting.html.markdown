---
layout: "rancher2"
page_title: "Rancher2: rancher2_setting"
sidebar_current: "docs-rancher2-datasource-setting"
description: |-
  Get information on a Rancher v2 setting.
---

# rancher2\_setting

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
