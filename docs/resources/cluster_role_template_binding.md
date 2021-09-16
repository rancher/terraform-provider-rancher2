---
page_title: "rancher2_cluster_role_template_binding Resource"
---

# rancher2\_cluster\_role\_template\_binding Resource

Provides a Rancher v2 Cluster Role Template Binding resource. This can be used to create Cluster Role Template Bindings for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher2 Cluster Role Template Binding
resource "rancher2_cluster_role_template_binding" "foo" {
  name = "foo"
  cluster_id = "<cluster_id>"
  role_template_id = "<role_template_id>"
  user_id = "<user_id>"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster id where bind cluster role template binding (string)
* `role_template_id` - (Required/ForceNew) The role template id from create cluster role template binding (string)
* `name` - (Required/ForceNew) The name of the cluster role template binding (string)
* `group_id` - (Optional/Computed/ForceNew) The group ID to assign cluster role template binding (string)
* `group_principal_id` - (Optional/Computed/ForceNew) The group_principal ID to assign cluster role template binding (string)
* `user_id` - (Optional/Computed/ForceNew) The user ID to assign cluster role template binding (string)
* `user_principal_id` - (Optional/Computed/ForceNew) The user_principal ID to assign cluster role template binding (string)
* `annotations` - (Optional/Computed) Annotations for cluster role template binding (map)
* `labels` - (Optional/Computed) Labels for cluster role template binding (map)

**Note** user `user_id | user_principal_id` OR group `group_id | group_principal_id` must be defined

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_cluster_role_template_binding` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cluster role template bindings.
- `update` - (Default `10 minutes`) Used for cluster role template binding modifications.
- `delete` - (Default `10 minutes`) Used for deleting cluster role template bindings.

## Import

Cluster Role Template Bindings can be imported using the Rancher cluster Role Template Binding ID

```
$ terraform import rancher2_cluster_role_template_binding.foo &lt;CLUSTER_ROLE_TEMPLATE_BINDING_ID&gt;
```
