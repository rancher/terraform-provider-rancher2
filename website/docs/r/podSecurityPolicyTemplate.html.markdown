---
layout: "rancher2"
page_title: "Rancher2: rancher2_pod_security_policy_template"
sidebar_current: "docs-rancher2-resource-pod-security-policy-template"
description: |-
  Provides a Rancher v2 PodSecurityPolicyTemplate resource. This can be used to create PodSecurityPolicyTemplates for Rancher v2 environments and retrieve their information.
---

# rancher2\_pod\_security\_policy\_template

Provides a Rancher v2 PodSecurityPolicyTemplate resource. This can be used to create PodSecurityPolicyTemplates for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl-terraform
# Create a new rancher2 PodSecurityPolicyTemplate
resource "rancher2_pod_security_policy_template" "foo" {
  name = "foo"
  description = "Terraform PodSecurityPolicyTemplate acceptance test - update"
  allow_privilege_escalation = false
  allowed_csi_driver {
    name = "something"
  }
  allowed_csi_driver {
    name = "something-else"
  }
  allowed_flex_volume {
    driver = "something"
  }
  allowed_flex_volume {
    driver = "something-else"
  }
  allowed_host_path {
    path_prefix = "/"
    read_only = true
  }
  allowed_host_path {
    path_prefix = "//"
    read_only = false
  }
  allowed_proc_mount_types = ["Default"]
  default_allow_privilege_escalation = false
  fs_group {
    rule = "MustRunAs"
    range {
      min = 0
      max = 100
    }
    range {
      min = 0
      max = 100
    }
  }
  host_ipc = false
  host_network = false
  host_pid = false
  host_port {
    min = 0
    max = 65535
  }
  host_port {
    min = 1024
    max = 8080
  }
  privileged = false
  read_only_root_filesystem = false
  required_drop_capabilities = ["something"]

  run_as_user {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  run_as_group {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  runtime_class {
    default_runtime_class_name = "something"
    allowed_runtime_class_names  = ["something"]
  }
  se_linux {
    rule = "RunAsAny"
  }
  supplemental_group {
    rule = "RunAsAny"
  }
  volumes = ["azureFile"]
}
```

## Argument Reference

The following arguments are supported:

Any field without a description is taken from the PodSecurityPolicy object definition in Kubernetes: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#podsecuritypolicy-v1beta1-extensions](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#podsecuritypolicy-v1beta1-extensions)

* `name` - (Required) The name of the PodSecurityPolicyTemplate (string)
* `description` - (Optional) The PodSecurityPolicyTemplate description (string)
* `annotations` - (Optional/Computed) Annotations for PodSecurityPolicyTemplate object (map)
* `labels` - (Optional/Computed) Labels for PodSecurityPolicyTemplate object (map)
* `allow_privilege_escalation` = (Optional)
* `allowed_capabilities` - (Optional) (list)
* `allowed_csi_driver` - (Optional) (list)
* `allowed_flex_volume` - (Optional) (list)
* `allowed_host_path` - (Optional) (list)
* `allowed_proc_mount_types` - (Optional) (list)
* `allowed_unsafe_sysctls` - (Optional) (list)
* `default_add_capabilities` - (Optional) (list)
* `default_allow_privilege_escalation` - (Optional) (list)
* `forbidden_sysctls` - (Optional) (list)
* `fs_group` - (Optional) (list maxitems:1)
* `host_ipc` - (Optional) (bool)
* `host_pid` - (Optional) (bool)
* `host_port` - (Optional) (list)
* `privileged` - (Optional) (bool)
* `read_only_root_filesystem` - (Optional) (bool)
* `required_drop_capabilities` - (Optional) (list)
* `run_as_user` - (Optional) (list maxitems:1)
* `run_as_group` - (Optional) (list maxitems:1)
* `runtime_class` - (Optional) (list maxitems:1)
* `se_linux` - (Optional) (list maxitems:1)
* `supplemental_group` - (Optional) (list maxitems:1)
* `volumes` - (Optional) (list)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

Any field without a description is taken from the PodSecurityPolicy object definition in Kubernetes: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#podsecuritypolicy-v1beta1-extensions](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#podsecuritypolicy-v1beta1-extensions)

### `allowed_host_path`

#### Arguments

* `path_prefix` - (Required) (string)
* `read_only` - (Optional) (string)

### `fs_group`

#### Arguments

* `range` - (Optional) (list)
* `rule` - (Optional) (string)

### `range`

#### Arguments

* `min` - (Required) (int)
* `max` - (Required) (int)

### `host_port`

#### Arguments

* `min` - (Required) (int)
* `max` - (Required) (int)

### `run_as_user`

#### Arguments

* `rule` - (Required) (string)
* `range` - (Optional) (list)

### `run_as_group`

#### Arguments

* `rule` - (Required) (string)
* `range` - (Optional) (list)

### `runtime_class`

#### Arguments

* `allowed_runtime_class_names` - (Required) (list)
* `default_runtime_class_name` - (Optional) (string)

### `se_linux`

#### Arguments

* `rule` - (Required) (string)
* `se_linux_option` - (Optional) (list maxitems:1)

### `se_linux_option`

#### Arguments

* `level` - (Optional) (string)
* `role` - (Optional) (string)
* `type` - (Optional) (string)
* `user` - (Optional) (string)

### `supplemental_group`

#### Arguments

* `rule` - (Required) (string)
* `range` - (Optional) (list)


## Timeouts

Because the `rancher2_pod_security_policy_template` API endpoint does not have a `state` field defined, this provider does
not implement timeouts, only direct API calls (without retry on failure).

## Import

PodSecurityPolicyTemplate can be imported using the Rancher PodSecurityPolicyTemplate Name

```
$ terraform import rancher2_pod_security_policy_template.foo &lt;pod_security_policy_name&gt;
```

