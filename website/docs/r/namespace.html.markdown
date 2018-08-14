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
  cluster_id = "<CLUSTER_ID>"
  project_name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the namespace.
* `cluster_id` - (Required) The cluster id where create namespace.
* `description` - (Optional) A namespace description.
* `project_name` - (Optional) Rancher Project name where assign namespace.
* `resource_quota_template_id` - (Optional) Rancher v2 resource quota template id to apply to namespace.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `project_id` - (Computed) The project ID of the associated project.

## Import

Projects can be imported using the namespace ID in the format `<cluster_id>:<namespace_id>`

```
$ terraform import rancher2_namespace.foo <cluster_id>:<namespace_id>
```

