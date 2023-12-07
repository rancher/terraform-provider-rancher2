---
page_title: "rancher2_project_role_template_binding Resource"
---

# rancher2\_project\_role\_template\_binding Resource

Provides a Rancher v2 Project Role Template Binding resource. This can be used to create Project Role Template Bindings for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Project Role Template Binding
resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "<project_id>"
  role_template_id = "<role_template_id>"
  user_id = "<user_id>"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required/ForceNew) The project id where bind project role template, in format `cluster_id:project_id`, e.g. `c-abc12:p-def34` (string)
* `role_template_id` - (Required/ForceNew) The role template id from create project role template binding (string)
* `name` - (Required/ForceNew) The name of the project role template binding (string)
* `group_id` - (Optional/Computed/ForceNew) The group ID to assign project role template binding (string)
* `group_principal_id` - (Optional/Computed/ForceNew) The group_principal ID to assign project role template binding (string)
* `user_id` - (Optional/Computed/ForceNew) The user ID to assign project role template binding (string)
* `user_principal_id` - (Optional/Computed/ForceNew) The user_principal ID to assign project role template binding (string)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

**Note:** user `user_id | user_principal_id` OR group `group_id | group_principal_id` must be defined

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource, e.g. `p-def34-blabla` (string)

## Timeouts

`rancher2_project_role_template_binding` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating project role template bindings.
- `update` - (Default `10 minutes`) Used for project role template binding modifications.
- `delete` - (Default `10 minutes`) Used for deleting project role template bindings.

## Import

Project Role Template Bindings can be imported using the Rancher Project Role Template Binding ID

```
$ terraform import rancher2_project_role_template_binding.foo &lt;project_role_template_binding_id&gt;
```

Retrieve the ID from the `projectroletemplatebindings` API URL: look at `id` field of the object with `baseType": "projectRoleTemplateBinding"` and matching name and project ID. The URL is provided by the error message you get if you attempt import with an incorrect ID, for example:
```
Error: Bad response statusCode [404]. Status [404 Not Found]. Body: [message=failed to find resource by id, baseType=error, code=NotFound] from [https://<domain>/v3/projectroletemplatebindings]
```
