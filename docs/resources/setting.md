---
page_title: "rancher2_setting Resource"
---

# rancher2\_setting Resource

Provides a Rancher v2 Setting resource. This can be used to create settings for Rancher v2 environments and retrieve their information.

On create, if setting already exists, provider will import it and update its value.

On destroy, if setting is a system setting like `server-url`, provider'll not delete it from Rancher, it'll just update setting value to default and remove it from tfstate. 

## Example Usage

```hcl
# Create a new rancher2 Setting
resource "rancher2_setting" "foo" {
  name = "foo"
  value = "<VALUE>"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the setting (string)
* `value` - (Required) The value of the setting (string)
* `annotations` - (Optional/Computed) Annotations for setting object (map)
* `labels` - (Optional/Computed) Labels for setting object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Import

Setting can be imported using the Rancher setting ID.

```
$ terraform import rancher2_setting.foo &lt;setting_id&gt;
```

