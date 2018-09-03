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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the namespace.
* `project_id` - (Required) The project id where create namespace.
* `description` - (Optional) A namespace description.
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `project_id` - (Computed) The project ID of the associated project.

## Import

Projects can be imported using the namespace ID in the format `<cluster_id>:<namespace_id>`

```
$ terraform import rancher2_namespace.foo <cluster_id>:<namespace_id>
```

