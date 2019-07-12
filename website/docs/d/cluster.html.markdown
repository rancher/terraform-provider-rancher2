---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster"
sidebar_current: "docs-rancher2-resource-cluster"
description: |-
  Get information on a Rancher v2 cluster.
---

# rancher2\_cluster

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
* `cluster_registration_token` - (Computed) Cluster Registration Token generated for the cluster (list maxitems:1)
* `default_project_id` - (Computed) Default project ID for the cluster (string)
* `driver` - (Computed) The driver used for the Cluster. `imported`, `azurekubernetesservice`, `amazonelasticcontainerservice`, `googlekubernetesengine` and `rancherKubernetesEngine` are supported (string)
* `kube_config` - (Computed) Kube Config generated for the cluster (string)
* `system_project_id` - (Computed) System project ID for the cluster (string)
* `rke_config` - (Computed) The rke configuration for `rke` Clusters. Conflicts with `aks_config`, `eks_config` and `gke_config` (list maxitems:1)
* `aks_config` - (Computed) The Azure aks configuration for `aks` Clusters. Conflicts with `eks_config`, `gke_config` and `rke_config` (list maxitems:1)
* `eks_config` - (Computed) The Amazon eks configuration for `eks` Clusters. Conflicts with `aks_config`, `gke_config` and `rke_config` (list maxitems:1)
* `gke_config` - (Computed) The Google gke configuration for `gke` Clusters. Conflicts with `aks_config`, `eks_config` and `rke_config` (list maxitems:1)
* `description` - (Computed) The description for Cluster (string)
* `cluster_auth_endpoint` - (Computed) Enabling the [local cluster authorized endpoint](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#local-cluster-auth-endpoint) allows direct communication with the cluster, bypassing the Rancher API proxy. (list maxitems:1)
* `default_pod_security_policy_template_id` - (Computed) [Default pod security policy template id](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#pod-security-policy-support). `restricted` and `unrestricted` are supported (string)
* `enable_network_policy` - (Computed) Enable project network isolation. Default `false` (bool)
* `annotations` - (Computed) Annotations for Node Pool object (map)
* `labels` - (Computed) Labels for Node Pool object (map)
