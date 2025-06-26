---
page_title: "rancher2_pod_security_admission_configuration_template Resource"
---

# rancher2\_pod\_security\_admission\_configuration\_template Resource

Provides a rancher v2 pod security admission configration template resource. 
This can be used to create pod security admission configration templates and retrieve their information.

For more information, please refer to [Rancher Documentation](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/authentication-permissions-and-global-configuration/psa-config-templates)

## Example Usage

```hcl
# Create a Pod Security Admission Configuration Template resource
resource "rancher2_pod_security_admission_configuration_template" "foo" {
  name = "custom-psact"
  description = "This is my custom Pod Security Admission Configuration Template"
  defaults {
    audit = "restricted"
    audit_version = "latest"
    enforce = "restricted"
    enforce_version = "latest"
    warn = "restricted"
    warn_version = "latest"
  }
  exemptions {
    usernames = ["testuser"]
    runtime_classes = ["testclass"]
    namespaces = ["ingress-nginx","kube-system"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required/ForceNew) The name of the pod security admission configuration template (string)
* `description` - (Optional) The description of the pod security admission configuration template (string) 
* `defaults` - (Required) The default level labels and version labels to be applied when labels for a mode is not set (list maxitems:1)
* `exemptions`- (Optional) The authenticated usernames, runtime class names, and namespaces to exempt (list maxitems:1)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `defaults`

#### Arguments

Level label values must be one of the following: `privileged` (default), `baseline`, `restricted`.

Version label values must be one of the following: `latest` (default), specific version like `v1.33`.

* `audit` - (Optional) The audit-mode's level label. Default `privileged` (string)
* `audit_version` - (Optional) The audit-mode's version label. Default `latest` (string)
* `enforce` - (Optional) The enforce-mode's level label. Default `privileged` (string)
* `enforce_version` - (Optional) The enforce-mode's version label. Default `latest` (string)
* `warn` - (Optional) The warn-mode's level label. Default `privileged` (string)
* `warn_version` - (Optional) The warn-mode's version label. Default `latest` (string)

### `exemptions`

#### Arguments

* `usernames` - (Optional) The list of authenticated usernames to exempt (List)
* `runtime_classes` - (Optional) The list of runtime class names to exempt (List)
* `namespaces` - (Optional) The list of namespaces to exempt (List)

## Timeouts

`rancher2_pod_security_admission_configuration_template` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating the resource.
- `update` - (Default `10 minutes`) Used for updating the resource.
- `delete` - (Default `10 minutes`) Used for deleting the resource.

## Import

Pod Security Admission Configration Templates can be imported using its ID:

```
$ terraform import rancher2_pod_security_admission_configuration_template.foo &lt;resource_id&gt;
```

