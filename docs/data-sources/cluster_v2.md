---
page_title: "rancher2_cluster_v2 Data Source"
---

# rancher2\_cluster\_v2 Data Source

Use this data source to retrieve information about a Rancher v2 cluster.

## Example Usage

```hcl
data "rancher2_cluster_v2" "foo" {
  name = "foo"
  fleet_namespace = "fleet-ns"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cluster v2 (string)
* `fleet_namespace` - (Optional) The fleet namespace of the Cluster v2. Default: `\"fleet-default\"` (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cluster_registration_token` - (Computed/Sensitive) Cluster Registration Token generated for the cluster v2 (list maxitems:1)
* `kube_config` - (Computed/Sensitive) Kube Config generated for the cluster v2 (string)
* `cluster_v1_id` - (Computed) Cluster v1 id for cluster v2 (string)
* `resource_version` - (Computed) Cluster v2 k8s resource version (string)
* `kubernetes_version` - (Computed) The kubernetes version of the Cluster v2 (list maxitems:1)
* `agent_env_vars` - (Computed) Optional Agent Env Vars for Rancher agent (list)
* `rke_config` - (Computed) The RKE configuration for `k3s` and `rke2` Clusters v2. (list maxitems:1)
* `cloud_credential_secret_name` - (Computed) Cluster V2 cloud credential secret name (string)
* `default_pod_security_policy_template_name` - (Computed) Cluster V2 default pod security policy template name (string)
* `default_cluster_role_for_project_members` - (Computed) Cluster V2 default cluster role for project members (string)
* `enable_network_policy` - (Computed) Enable k8s network policy at Cluster V2 (bool)
