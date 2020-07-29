---
page_title: "rancher2_project Data Source"
---

# rancher2\_project Data Source

Use this data source to retrieve information about a Rancher v2 project.
This data source can be used in conjunction with the Terraform
[Kubernetes provider](https://www.terraform.io/docs/providers/kubernetes/)
to associate Namespaces with projects.

## Example Usage

```
data "rancher2_project" "system" {
    cluster_id = var.my_cluster_id
    name = "System"
}

resource "kubernetes_namespace" "my_namespace" {
  metadata {
    annotations {
      "field.cattle.io/projectId" = data.rancher2_project.system.id
    }
    name = "my-namespace"
  }
}
```

## Argument Reference

 * `cluster_id` - (Required) ID of the Rancher 2 cluster (string)
 * `name` - (Required) The project name (string)

## Attributes Reference

 * `id` - (Computed) Cluster-wide unique ID of the Rancher 2 project (string)
 * `container_resource_limit` - (Computed) Default containers resource limits on project (List maxitem:1)
 * `enable_project_monitoring` - (Computed) Enable built-in project monitoring. Default `false` (bool)
 * `pod_security_policy_template_id` - (Computed) Default Pod Security Policy ID for the project (string)
 * `resource_quota` - (Computed) Resource quota for project. Rancher v2.1.x or higher (list maxitems:1)
 * `uuid` - (Computed) UUID of the project as stored by Rancher 2 (string)
 * `description` - (Computed) The project's description (string)
 * `annotations` - (Computed) Annotations of the rancher2 project (map)
 * `labels` - (Computed) Labels of the rancher2 project (map)
