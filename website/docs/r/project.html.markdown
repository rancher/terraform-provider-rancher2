---
layout: "rancher2"
page_title: "Rancher2: rancher2_project"
sidebar_current: "docs-rancher2-resource-project"
description: |-
  Provides a Rancher v2 Project resource. This can be used to create projects for rancher v2 environments and retrieve their information.
---

# rancher2\_project

Provides a Rancher v2 Project resource. This can be used to create projects for rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Project
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "<CLUSTER_ID>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project.
* `cluster_id` - (Required) The cluster id where create project.
* `description` - (Optional) A project description.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.

## Import

Projects can be imported using the rancher Project ID

```
$ terraform import rancher2_project.foo <project_id>
```

