---
layout: "rancher2"
page_title: "Rancher2: rancher2_global_role_binding"
sidebar_current: "docs-rancher2-datasource-global_role_binding"
description: |-
  Get information on a Rancher v2 global role binding.
---

# rancher2\_global\_role\_binding

Use this data source to retrieve information about a Rancher v2 global role binding.

## Example Usage

```
data "rancher2_global_role_binding" "foo" {
    name = "foo"
    global_role_id = "foo_id"
}
```

## Argument Reference

* `name` - (Required) The name of the global role binding (string)
* `global_role_id` - (Optional/Computed) The global role id (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `user_id` - (Computed) The user ID to assign global role binding (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)