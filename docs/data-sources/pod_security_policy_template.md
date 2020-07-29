---
page_title: "rancher2_pod_security_policy_template Data Source"
---

# rancher2\_pod\_security\_policy\_template Data Source

Use this data source to retrieve information about a Rancher v2 PodSecurityPolicyTemplate.

## Example Usage

```hcl-terraform
data "rancher2_pod_security_policy_template" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The name of the PodSecurityPolicyTemplate (string)

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
