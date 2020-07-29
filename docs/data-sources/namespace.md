---
page_title: "rancher2_namespace Data Source"
---

# rancher2\_namespace Data Source

Use this data source to retrieve information about a Rancher v2 namespace.

## Example Usage

```hcl
data "rancher2_namespace" "foo" {
  name = "foo"
  project_id = rancher2_cluster.foo-custom.default_project_id
}
```

## Argument Reference

* `name` - (Required) The name of the namespace (string)
* `project_id` - (Required) The project id where namespace is assigned (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `container_resource_limit` - (Computed) Default containers resource limits on namespace (List maxitem:1)
* `description` - (Computed) A namespace description (string)
* `resource_quota` - (Computed) Resource quota for namespace. Rancher v2.1.x or higher (list maxitems:1)
* `annotations` - (Computed) Annotations for Node Pool object (map)
* `labels` - (Computed) Labels for Node Pool object (map)
