---
page_title: "rancher2_feature Resource"
---

# rancher2\_feature Resource

Provides a Rancher v2 Feature resource. This can be used to enable/disable [experimental features](https://rancher.com/docs/rancher/v2.x/en/installation/resources/feature-flags/) for Rancher v2 environments.

Experimental features already exist at Rancher v2.5.x systems, so they can just be updated: 
* On create, provider will read Feature from Rancher and update its value. It will return an error if feature doesn't exist
* On destroy, provider will not delete feature from Rancher, just from tfstate

## Example Usage

```hcl
# Create a new rancher2 Feature
resource "rancher2_feature" "foo" {
  name = "foo"
  value = "<VALUE>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the feature (string)
* `value` - (Optional) The value of the feature. Default: `false` (bool)
* `annotations` - (Optional/Computed) Annotations for feature object (map)
* `labels` - (Optional/Computed) Labels for feature object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
