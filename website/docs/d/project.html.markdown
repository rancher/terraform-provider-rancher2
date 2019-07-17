---
layout: "rancher2"
page_title: "Rancher2: rancher2_project"
sidebar_current: "docs-rancher2-datasource-project"
description: |-
  Get information on a Rancher v2 project.
---

# rancher2\_project

Use this data source to retrieve information about a Rancher v2 project.
This data source can be used in conjunction with the Terraform
[Kubernetes provider](https://www.terraform.io/docs/providers/kubernetes/)
to associate Namespaces with projects.

## Example Usage

```
data "rancher2_project" "system" {
    cluster_id = "${var.my_cluster_id}"
    name = "System"
}

resource "kubernetes_namespace" "my_namespace" {
  metadata {
    annotations {
      "field.cattle.io/projectId" = "${data.rancher2_project.system.id}"
    }
    name = "my-namespace"
  }
}
```

## Argument Reference

 * `cluster_id` - (Required) ID of the Rancher 2 cluster.
 * `name` - (Required) The project name.

## Attributes Reference

 * `id` - Cluster-wide unique ID of the Rancher 2 project.
 * `uuid` - UUID of the project as stored by Rancher 2.
 * `description` - The project's description.
 * `enable_project_monitoring` - (Computed) Enable built-in project monitoring. Default `false` (bool)
 * `annotations` - Annotations of the rancher2 project (map).
 * `labels` - Labels of the rancher2 project (map).
