---
page_title: "rancher2_cluster Data Source"
---

# rancher2\_cluster Data Source

Use this data source to retrieve information about a Rancher v2 cluster.

## Example Usage

```hcl
data "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cluster (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `agent_env_vars` - (Computed) Optional Agent Env Vars for Rancher agent. Just for Rancher v2.5.6 and above (list)
* `cluster_registration_token` - (Computed) Cluster Registration Token generated for the cluster (list maxitems:1)
* `default_project_id` - (Computed) Default project ID for the cluster (string)
* `driver` - (Computed) The driver used for the Cluster. `imported`, `azurekubernetesservice`, `amazonelasticcontainerservice`, `googlekubernetesengine` and `rancherKubernetesEngine` are supported (string)
* `kube_config` - (Computed) Kube Config generated for the cluster (string)
* `ca_cert` - (Computed) K8s cluster ca cert (string)
* `system_project_id` - (Computed) System project ID for the cluster (string)
* `rke_config` - (Computed) The RKE configuration for `rke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` and `k3s_config` (list maxitems:1)
* `rke2_config` - (Computed) The RKE2 configuration for `rke2` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `gke_config`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* `k3s_config` - (Computed) The K3S configuration for `k3s` imported Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` and `rke_config` (list maxitems:1)
* `aks_config` - (Computed) The Azure aks configuration for `aks` Clusters. Conflicts with `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* `aks_config_v2` - (Optional) The Azure AKS v2 configuration for creating/import `aks` Clusters. Conflicts with `aks_config`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config` (list maxitems:1)
* `eks_config` - (Computed) The Amazon eks configuration for `eks` Conflicts with `aks_config`, `aks_config_v2`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* `eks_config_v2` - (Computed) The Amazon EKS V2 configuration to create or import `eks` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `gke_config`, `gke_config_v2`, `oke_config`, `k3s_config` and `rke_config`. For Rancher v2.5.x or above (list maxitems:1)
* `gke_config` - (Computed) The Google gke configuration for `gke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config_v2`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1) (list maxitems:1)
* `gke_config_v2` - (Computed) The Google GKE V2 configuration for `gke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `oke_config`, `k3s_config` and `rke_config`. For Rancher v2.5.8 or above (list maxitems:1)
* `oke_config` - (Computed) The Oracle OKE configuration for `oke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `k3s_config` and `rke_config` (list maxitems:1)
* `description` - (Computed) The description for Cluster (string)
* `cluster_auth_endpoint` - (Computed) Enabling the [local cluster authorized endpoint](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#local-cluster-auth-endpoint) allows direct communication with the cluster, bypassing the Rancher API proxy. (list maxitems:1)
* `cluster_monitoring_input` - (Computed) Cluster monitoring config (list maxitems:1)
* `cluster_template_answers` - (Computed) Cluster template answers (list maxitems:1)
* `cluster_template_id` - (Computed) Cluster template ID (string)
* `cluster_template_questions` - (Computed) Cluster template questions (list)
* `cluster_template_revision_id` - (Computed) Cluster template revision ID (string)
* `default_pod_security_policy_template_id` - (Optional/Computed) [Default pod security policy template id](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#pod-security-policy-support) (string)
* `enable_cluster_monitoring` - (Computed) Enable built-in cluster monitoring. Default `false` (bool)
* `enable_network_policy` - (Computed) Enable project network isolation. Default `false` (bool)
* `enable_cluster_istio` - (Computed) Enable built-in cluster istio. Default `false` (bool)
* `fleet_workspace_name` - (Computed) Fleet workspace name (string)
* `annotations` - (Computed) Annotations for Node Pool object (map)
* `labels` - (Computed) Labels for Node Pool object (map)
