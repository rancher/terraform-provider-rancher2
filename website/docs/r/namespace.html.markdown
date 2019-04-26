---
layout: "rancher2"
page_title: "Rancher2: rancher2_namespace"
sidebar_current: "docs-rancher2-resource-namespace"
description: |-
  Provides a Rancher v2 Namespace resource. This can be used to create namespaces for rancher v2 environments and retrieve their information.
---

# rancher2\_namespace

Provides a Rancher v2 Namespace resource. This can be used to create namespaces for rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Namespace
resource "rancher2_namespace" "foo" {
  name = "foo"
  project_id = "<PROJECT_ID>"
  description = "foo namespace"
  resource_quota {
    limit {
      limits_cpu = "100m"
      limits_memory = "100Mi"
      requests_storage = "1Gi"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the namespace (string)
* `project_id` - (Required) The project id where assign namespace. It's on the form `project_id=<cluster_id>:<id>`. Updating `<id>` part on same `<cluster_id>` namespace will be moved between projects (string)
* `description` - (Optional) A namespace description (string)
* `resource_quota` - (Optional/Computed) Resource quota for namespace. Rancher v2.1.x or higher (list maxitems:1)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `resource_quota`

#### Arguments

* `limit` - (Required) Resource quota limit for namespace (list maxitems:1)

#### `limit`

##### Arguments

* `config_maps` - (Optional) Limit for config maps in namespace (string)
* `limits_cpu` - (Optional) Limit for limits cpu in namespace (string)
* `limits_memory` - (Optional) Limit for limits memory in namespace (string)
* `persistent_volume_claims` - (Optional) Limit for persistent volume claims in namespace (string)
* `pods` - (Optional) Limit for pods in namespace (string)
* `replication_controllers` - (Optional) Limit for replication controllers in namespace (string)
* `requests_cpu` - (Optional) Limit for requests cpu in namespace (string)
* `requests_memory` - (Optional) Limit for requests memory in namespace (string)
* `requests_storage` - (Optional) Limit for requests storage in namespace (string)
* `secrets` - (Optional) Limit for secrets in namespace (string)
* `services_load_balancers` - (Optional) Limit for services load balancers in namespace (string)
* `services_node_ports` - (Optional) Limit for services node ports in namespace (string)

More info at [resource-quotas](https://rancher.com/docs/rancher/v2.x/en/k8s-in-rancher/projects-and-namespaces/resource-quotas/)

## Timeouts

`rancher2_namespace` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating namespaces.
- `update` - (Default `10 minutes`) Used for namespace modifications.
- `delete` - (Default `10 minutes`) Used for deleting namespaces.

## Import

Projects can be imported using the namespace ID in the format `<cluster_id>:<namespace_id>`

```
$ terraform import rancher2_namespace.foo <cluster_id>:<namespace_id>
```

When you import a raw k8s namespace, `project_id=<cluster_id>`. It'll not be assigned to any project. To move it into a project, update `project_id=<cluster_id>:<id>`. Namespace move is only supported inside same `cluster_id`.
