---
page_title: "rancher2_principal Data Source"
---

# rancher2\_principal Data Source

Use this data source to retrieve information about a Rancher v2 Principal resource.

## Example Usage

```hcl
data "rancher2_principal" "foo" {
  name = "user@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The full name of the principal (string)
* `type` - (Optional) The type of the identity (string). Defaults to `user`. Only `user` and `group` values are supported (string)
* `exact_match` - (Optional) If set to `true`, only the exactly matched result is returned. Defaults to `false`, which means a partially matched result can be returned (for example: `foo2` also matches for `foo` search input) (bool)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
