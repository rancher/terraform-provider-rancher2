---
page_title: "rancher2_pod_security_admission_configuration_template Data Source"
---

# rancher2\_pod\_security\_admission\_configuration\_template Resource

Use this data source to retrieve information about a rancher v2 pod security admission configration template.

## Example Usage

```hcl
data "rancher2_pod_security_admission_configuration_template" "foo" {
    name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the pod security admission configuration template (string)

## Attributes Reference

The following attributes are exported:

* `description` - (Computed) The description of the pod security admission configuration template (string)
* `defaults` - (Computed) The default level labels and version labels to be applied when labels for a mode is not set (list maxitems:1)
* `exemptions`- (Computed) The authenticated usernames, runtime class names, and namespaces to exempt (list maxitems:1)
* `annotations` - (Computed) Annotations of the resource (map)
* `labels` - (Computed) Labels of the resource (map)

