---
layout: "cattle"
page_title: "Cattle: cattle_project"
sidebar_current: "docs-cattle-resource-project"
description: |-
  Provides a Cattle Project resource. This can be used to create projects for rancher v2 environments and retrieve their information.
---

# cattle\_project

Provides a Cattle Project resource. This can be used to create projects for rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new Cattle Project
resource "cattle_project" "foo" {
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
$ terraform import cattle_project.foo <project_id>
```

