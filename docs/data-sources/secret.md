---
page_title: "rancher2_secret Data Source"
---

# rancher2\_secret Data Source

Use this data source to retrieve information about a Rancher v2 secret.

Depending of the availability, there are 2 types of Rancher v2 secrets:
- Project secret: Available to all namespaces in the `project_id`
- Namespaced secret: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Retrieve a rancher2 Project Secret
data "rancher2_secret" "foo" {
  name = "<name>"
  project_id = "<project_id>"
}
```

```hcl
# Retrieve a rancher2 Namespaced Secret
data "rancher2_secret" "foo" {
  name = "<name>"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the secret (string)
* `project_id` - (Required) The project id where to assign the secret (string)
* `namespace_id` - (Optional) The namespace id where to assign the namespaced secret (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `data` - (Computed) Secret key/value data. Base64 encoding required for values (map)
* `description` - (Computed) A secret description (string)
* `annotations` - (Computed) Annotations for secret object (map)
* `labels` - (Computed) Labels for secret object (map)
