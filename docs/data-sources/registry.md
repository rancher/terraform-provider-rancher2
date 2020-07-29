---
page_title: "rancher2_registry Data Source"
---

# rancher2\_registry Data Source

Use this data source to retrieve information about a Rancher v2 docker registry.

Depending of the availability, there are 2 types of Rancher v2 docker registries:
- Project registry: Available to all namespaces in the `project_id`
- Namespaced registry: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Retrieve a rancher2 Project Registry
data "rancher2_registry" "foo" {
  name = "<name>"
  project_id = "<project_id>"
}
```

```hcl
# Retrieve a rancher2 Namespaced Registry
data "rancher2_registry" "foo" {
  name = "<name>"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the registry (string)
* `project_id` - (Required) The project id where to assign the registry (string)
* `namespace_id` - (Optional) The namespace id where to assign the namespaced registry (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `registries` - (Computed) Registries data for registry (list)
* `description` - (Computed) A registry description (string)
* `annotations` - (Computed) Annotations for Registry object (map)
* `labels` - (Computed) Labels for Registry object (map)
