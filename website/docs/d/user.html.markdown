---
layout: "rancher2"
page_title: "Rancher2: rancher2_user"
sidebar_current: "docs-rancher2-datasource-user"
description: |-
  Get information on a Rancher v2 user.
---

# rancher2\_user

Use this data source to retrieve information about a Rancher v2 user

## Example Usage

```
data "rancher2_user" "foo" {
    username = "foo"
    global_role_id = "foo_id"
}
```

## Argument Reference

* `username` - (Required) The name of the user (string)
* `global_role_id` - (Optional/Computed) The global role id (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `user_id` - (Computed) The user ID to assign global role binding (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)
