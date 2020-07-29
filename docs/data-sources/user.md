---
page_title: "rancher2_user Data Source"
---

# rancher2\_user

Use this data source to retrieve information about a Rancher v2 user

## Example Usage

```
data "rancher2_user" "foo" {
    username = "foo"
}
```

## Argument Reference

* `is_external` - (Optional) Set is the user if the user is external. Default: `false` (bool)
* `name` - (Optional) The name of the user (string)
* `username` - (Optional) The username of the user (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The user common name (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `enabled` - (Computed) The user is enabled (bool)
* `principal_ids` - (Computed) The user principal IDs (list)
* `labels` - (Computed) Labels of the resource (map)
