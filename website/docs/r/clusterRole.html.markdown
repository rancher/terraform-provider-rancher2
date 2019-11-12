---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster_role_template_binding"
sidebar_current: "docs-rancher2-resource-cluster_role_template_binding"
description: |-
  Provides a Rancher v2 Cluster Role Template Binding resource. This can be used to create Cluster Role Template Bindings for Rancher v2 environments and retrieve their information.
---

# rancher2\_cluster\_role\_template\_binding

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

* `cluster_id` - (Required) The cluster id where bind cluster role template binding (string)
* `role_template_id` - (Required) The role template id from create cluster role template binding (string)
* `name` - (Required) The name of the cluster role template binding (string)
* `group_id` - (Optional) The group ID to assign cluster role template binding (string)
* `group_principal_id` - (Optional) The group_principal ID to assign cluster role template binding (string)
* `user_id` - (Optional) The user ID to assign cluster role template binding (string)
* `user_principal_id` - (Optional) The user_principal ID to assign cluster role template binding (string)
* `annotations` - (Optional/Computed) Annotations for cluster role template binding (map)
* `labels` - (Optional/Computed) Labels for cluster role template binding (map)


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
$ terraform import rancher2_cluster_role_template_binding.foo <cluster_role_template_binding_id>
```
