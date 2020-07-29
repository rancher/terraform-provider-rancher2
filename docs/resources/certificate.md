---
page_title: "rancher2_certificate Resource"
---

# rancher2\_certificate Resource

Provides a Rancher v2 certificate resource. This can be used to create certificates for Rancher v2 environments and retrieve their information.

There are 2 types of Rancher v2 certificates:
- Project certificate: Available to all namespaces in the `project_id`
- Namespaced certificate: Available to just `namespace_id` in the `project_id`

## Example Usage

```hcl
# Create a new rancher2 Project Certificate
resource "rancher2_certificate" "foo" {
  certs = base64encode(<PUBLIC_CERTS>)
  key = base64encode(<PRIVATE_KEY>)
  name = "foo"
  description = "Terraform certificate foo"
  project_id = "<project_id>"
}
```

```hcl
# Create a new rancher2 Namespaced Certificate
resource "rancher2_certificate" "foo" {
  certs = base64encode(<PUBLIC_CERTS>)
  key = base64encode(<PRIVATE_KEY>)
  name = "foo"
  description = "Terraform certificate foo"
  project_id = "<project_id>"
  namespace_id = "<namespace_id>"
}
```

## Argument Reference

The following arguments are supported:

* `certs` - (Required) Base64 encoded public certs (string)
* `key` - (Required/Sensitive) Base64 encoded private key (string)
* `project_id` - (Required/ForceNew) The project id where the certificate should be created  (string)
* `description` - (Optional) A certificate description (string)
* `name` - (Optional/ForceNew) The name of the certificate (string)
* `namespace_id` - (Optional/ForceNew) The namespace id where the namespaced certificate should be created (string)
* `annotations` - (Optional/Computed) Annotations for certificate object (map)
* `labels` - (Optional/Computed) Labels for certificate object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_certificate` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating registries.
- `update` - (Default `10 minutes`) Used for certificate modifications.
- `delete` - (Default `10 minutes`) Used for deleting registries.
