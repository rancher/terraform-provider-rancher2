---
page_title: "rancher2_project_role_template_binding Data Source"
---

# rancher2\_project\_role\_template\_binding Data Source

Use this data source to retrieve information about a Rancher v2 project role template binding.

## Example Usage

```
data "rancher2_project_role_template_binding" "foo" {
    name = "foo"
    project_id = "foo_id"
}
```

## Argument Reference

* `name` - (Required) The name of the project role template binding (string)
* `project_id` - (Required) The project id where bind project role template (string)
* `role_template_id` - (Optional/Computed) The role template id from create project role template binding (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `group_id` - (Computed) The group ID to assign project role template binding (string)
* `group_principal_id` - (Computed) The group_principal ID to assign project role template binding (string)
* `user_id` - (Computed) The user ID to assign project role template binding (string)
* `user_principal_id` - (Computed) The user_principal ID to assign project role template binding (string)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)

