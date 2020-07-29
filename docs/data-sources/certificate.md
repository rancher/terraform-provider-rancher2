---
page_title: "rancher2_certificate Data Source"
---

# rancher2\_certificate

Use this data source to retrieve information about a Rancher v2 certificate.

Depending of the availability, there are 2 types of Rancher v2 certificates:
- Project certificate: Available to all namespaces in the `project_id`
- Namespaced certificate: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Retrieve a rancher2 Project Certificate
data "rancher2_certificate" "foo" {
  name = "<name>"
  project_id = "<project_id>"
}
```

```hcl
# Retrieve a rancher2 Namespaced Certificate
data "rancher2_certificate" "foo" {
  name = "<name>"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the certificate (string)
* `project_id` - (Required) The project id where to assign the certificate (string)
* `namespace_id` - (Optional) The namespace id where to assign the namespaced certificate (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `certs` - (Computed) Base64 encoded certs (string)
* `description` - (Computed) A certificate description (string)
* `annotations` - (Computed) Annotations for certificate object (map)
* `labels` - (Computed) Labels for certificate object (map)
