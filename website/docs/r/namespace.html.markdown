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

* `name` - (Required) The name of the namespace.
* `project_id` - (Required) The project id where assign namespace.
* `description` - (Optional) A namespace description.
* `resource_quota` - (Optional) Resource quota for namespace. Rancher v2.1.x or higher 
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

### Namespace resource quota `resource_quota`

The following arguments are supported:

* `limit` - (Required) Resource quota limit for namespace.

The following arguments are supported for `limit`:

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

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Projects can be imported using the namespace ID in the format `<cluster_id>:<namespace_id>`

```
$ terraform import rancher2_namespace.foo <cluster_id>:<namespace_id>
```

