---
page_title: "rancher2_principal Data Source"
---

# rancher2\_principal Data Source

Use this data source to retrieve information about a Rancher v2 Principal resource.

## Example Usage

```hcl
data "rancher2_principal" "foo" {
  email = "group@example.com"
  type  = "group"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email address of the identity (string)
* `type` - (Required) The type of the identity (string). Only `group` and `user` values are supported (string)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
