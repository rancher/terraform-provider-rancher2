---
page_title: "rancher2_cluster Resource"
---

# rancher2\_cluster Resource

Provides a Rancher v2 Cluster resource. This can be used to create Clusters for Rancher v2 environments and retrieve their information.

## Example Usage

**Note optional/computed arguments** If any `optional/computed` argument of this resource is defined by the user, removing it from tf file will NOT reset its value. To reset it, let its definition at tf file as empty/false object. Ex: `enable_cluster_monitoring = false`, `cloud_provider {}`, `name = ""`

### Creating Rancher v2 imported cluster

```hcl
# Create a new rancher2 imported Cluster
resource "rancher2_cluster" "foo-imported" {
  name = "foo-imported"
  description = "Foo rancher2 imported cluster"
}
```

Creating Rancher v2 RKE cluster

```hcl
# Create auditlog policy yaml file
auditlog_policy.yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
  - level: RequestResponse
    resources:
    - group: ""
      resources: ["pods"]

# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      kube_api {
        audit_log {
          enabled = true
          configuration {
            max_age = 5
            max_backup = 5
            max_size = 100
            path = "-"
            format = "json"
            policy = file("auditlog_policy.yaml")
          }
        }
      }
    }
  }
}
```

### Creating Rancher v2 RKE cluster enabling and customizing monitoring

**Note** Cluster monitoring version `0.2.0` or above, can't be enabled until cluster is fully deployed as [`kubeVersion`](https://github.com/rancher/system-charts/blob/52be656700468904b9bf15c3f39cd7112e1f8c9b/charts/rancher-monitoring/v0.2.0/Chart.yaml#L12) requirement has been introduced to helm chart 

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
  }
  enable_cluster_monitoring = true
  cluster_monitoring_input {
    answers = {
      "exporter-kubelets.https" = true
      "exporter-node.enabled" = true
      "exporter-node.ports.metrics.port" = 9796
      "exporter-node.resources.limits.cpu" = "200m"
      "exporter-node.resources.limits.memory" = "200Mi"
      "grafana.persistence.enabled" = false
      "grafana.persistence.size" = "10Gi"
      "grafana.persistence.storageClass" = "default"
      "operator.resources.limits.memory" = "500Mi"
      "prometheus.persistence.enabled" = "false"
      "prometheus.persistence.size" = "50Gi"
      "prometheus.persistence.storageClass" = "default"
      "prometheus.persistent.useReleaseName" = "true"
      "prometheus.resources.core.limits.cpu" = "1000m",
      "prometheus.resources.core.limits.memory" = "1500Mi"
      "prometheus.resources.core.requests.cpu" = "750m"
      "prometheus.resources.core.requests.memory" = "750Mi"
      "prometheus.retention" = "12h"
    }
    version = "0.1.0"
  }
}
```

### Creating Rancher v2 RKE cluster enabling/customizing monitoring and istio

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
  }
  enable_cluster_monitoring = true
  cluster_monitoring_input {
    answers = {
      "exporter-kubelets.https" = true
      "exporter-node.enabled" = true
      "exporter-node.ports.metrics.port" = 9796
      "exporter-node.resources.limits.cpu" = "200m"
      "exporter-node.resources.limits.memory" = "200Mi"
      "grafana.persistence.enabled" = false
      "grafana.persistence.size" = "10Gi"
      "grafana.persistence.storageClass" = "default"
      "operator.resources.limits.memory" = "500Mi"
      "prometheus.persistence.enabled" = "false"
      "prometheus.persistence.size" = "50Gi"
      "prometheus.persistence.storageClass" = "default"
      "prometheus.persistent.useReleaseName" = "true"
      "prometheus.resources.core.limits.cpu" = "1000m",
      "prometheus.resources.core.limits.memory" = "1500Mi"
      "prometheus.resources.core.requests.cpu" = "750m"
      "prometheus.resources.core.requests.memory" = "750Mi"
      "prometheus.retention" = "12h"
    }
    version = "0.1.0"
  }
}
# Create a new rancher2 Cluster Sync for foo-custom cluster
resource "rancher2_cluster_sync" "foo-custom" {
  cluster_id =  rancher2_cluster.foo-custom.id
  wait_monitoring = rancher2_cluster.foo-custom.enable_cluster_monitoring
}
# Create a new rancher2 Namespace
resource "rancher2_namespace" "foo-istio" {
  name = "istio-system"
  project_id = rancher2_cluster_sync.foo-custom.system_project_id
  description = "istio namespace"
}
# Create a new rancher2 App deploying istio (should wait until monitoring is up and running)
resource "rancher2_app" "istio" {
  catalog_name = "system-library"
  name = "cluster-istio"
  description = "Terraform app acceptance test"
  project_id = rancher2_namespace.foo-istio.project_id
  template_name = "rancher-istio"
  template_version = "0.1.1"
  target_namespace = rancher2_namespace.foo-istio.id
  answers = {
    "certmanager.enabled" = false
    "enableCRDs" = true
    "galley.enabled" = true
    "gateways.enabled" = false
    "gateways.istio-ingressgateway.resources.limits.cpu" = "2000m"
    "gateways.istio-ingressgateway.resources.limits.memory" = "1024Mi"
    "gateways.istio-ingressgateway.resources.requests.cpu" = "100m"
    "gateways.istio-ingressgateway.resources.requests.memory" = "128Mi"
    "gateways.istio-ingressgateway.type" = "NodePort"
    "global.monitoring.type" = "cluster-monitoring"
    "global.rancher.clusterId" = rancher2_cluster_sync.foo-custom.cluster_id
    "istio_cni.enabled" = "false"
    "istiocoredns.enabled" = "false"
    "kiali.enabled" = "true"
    "mixer.enabled" = "true"
    "mixer.policy.enabled" = "true"
    "mixer.policy.resources.limits.cpu" = "4800m"
    "mixer.policy.resources.limits.memory" = "4096Mi"
    "mixer.policy.resources.requests.cpu" = "1000m"
    "mixer.policy.resources.requests.memory" = "1024Mi"
    "mixer.telemetry.resources.limits.cpu" = "4800m",
    "mixer.telemetry.resources.limits.memory" = "4096Mi"
    "mixer.telemetry.resources.requests.cpu" = "1000m"
    "mixer.telemetry.resources.requests.memory" = "1024Mi"
    "mtls.enabled" = false
    "nodeagent.enabled" = false
    "pilot.enabled" = true
    "pilot.resources.limits.cpu" = "1000m"
    "pilot.resources.limits.memory" = "4096Mi"
    "pilot.resources.requests.cpu" = "500m"
    "pilot.resources.requests.memory" = "2048Mi"
    "pilot.traceSampling" = "1"
    "security.enabled" = true
    "sidecarInjectorWebhook.enabled" = true
    "tracing.enabled" = true
    "tracing.jaeger.resources.limits.cpu" = "500m"
    "tracing.jaeger.resources.limits.memory" = "1024Mi"
    "tracing.jaeger.resources.requests.cpu" = "100m"
    "tracing.jaeger.resources.requests.memory" = "100Mi"
  }
}
```

### Creating Rancher v2 RKE cluster assigning a node pool (overlapped planes)

```hcl
# Create a new rancher2 RKE Cluster
resource "rancher2_cluster" "foo-custom" {
  name = "foo-custom"
  description = "Foo rancher2 custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}
# Create a new rancher2 Node Template
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = ["<AWS_SECURITY_GROUP>"]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}
# Create a new rancher2 Node Pool
resource "rancher2_node_pool" "foo" {
  cluster_id =  rancher2_cluster.foo-custom.id
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  node_template_id = rancher2_node_template.foo.id
  quantity = 3
  control_plane = true
  etcd = true
  worker = true
}
```

### Creating Rancher v2 RKE cluster from template. For Rancher v2.3.x or above.

```hcl
# Create a new rancher2 cluster template
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "owner"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V1"
    cluster_config {
      rke_config {
        network {
          plugin = "canal"
        }
        services {
          etcd {
            creation = "6h"
            retention = "24h"
          }
        }
      }
    }
    default = true
  }
  description = "Test cluster template v2"
}
# Create a new rancher2 RKE Cluster from template
resource "rancher2_cluster" "foo" {
  name = "foo"
  cluster_template_id = rancher2_cluster_template.foo.id
  cluster_template_revision_id = rancher2_cluster_template.foo.template_revisions.0.id
}
```

### Creating Rancher v2 RKE cluster with upgrade strategy. For Rancher v2.4.x or above.

```hcl
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
      kube_api {
        audit_log {
          enabled = true
          configuration {
            max_age = 5
            max_backup = 5
            max_size = 100
            path = "-"
            format = "json"
            policy = "apiVersion: audit.k8s.io/v1\nkind: Policy\nmetadata:\n  creationTimestamp: null\nomitStages:\n- RequestReceived\nrules:\n- level: RequestResponse\n  resources:\n  - resources:\n    - pods\n"
          }
        }
      }
    }
    upgrade_strategy {
      drain = true
      max_unavailable_worker = "20%"
    }
  }
}
```

### Importing EKS cluster to Rancher v2, using `eks_config_v2`. For Rancher v2.5.x or above.

```hcl
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform EKS cluster"
  eks_config_v2 {
    cloud_credential_id = rancher2_cloud_credential.foo.id
    name = "<CLUSTER_NAME>"
    region = "<EKS_REGION>"
    imported = true
  }
}
```

### Creating EKS cluster from Rancher v2, using `eks_config_v2`. For Rancher v2.5.x or above.

```hcl
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform EKS cluster"
  eks_config_v2 {
    cloud_credential_id = rancher2_cloud_credential.foo.id
    region = "<EKS_REGION>"
    kubernetes_version = "1.24"
    logging_types = ["audit", "api"]
    node_groups {
      name = "node_group1"
      instance_type = "t3.medium"
      desired_size = 3
      max_size = 5
    }
    node_groups {
      name = "node_group2"
      instance_type = "m5.xlarge"
      desired_size = 2
      max_size = 3
      node_role = "arn:aws:iam::role/test-NodeInstanceRole"
    }
    private_access = true
    public_access = false
  }
}
```

### Creating EKS cluster from Rancher v2, using `eks_config_v2` and launch template. For Rancher v2.5.6 or above.

```hcl
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform EKS cluster"
  eks_config_v2 {
    cloud_credential_id = rancher2_cloud_credential.foo.id
    region = "<EKS_REGION>"
    kubernetes_version = "1.24"
    logging_types = ["audit", "api"]
    node_groups {
      desired_size = 3
      max_size = 5
      name = "node_group1"
      launch_template {
        id = "<EC2_LAUNCH_TEMPLATE_ID>"
        version = 1
      }
    }
    private_access = true
    public_access = true
  }
}
```

### Importing GKE cluster from Rancher v2, using `gke_config_v2`. For Rancher v2.5.8 or above.

```hcl
resource "rancher2_cloud_credential" "foo-google" {
  name = "foo-google"
  description= "Terraform cloudCredential acceptance test"
  google_credential_config {
    auth_encoded_json = file(<GOOGLE_AUTH_ENCODED_JSON>)
  }
}

resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Foo rancher2 imported GKE cluster"
  gke_config_v2 {
    name = "foo"
    google_credential_secret = rancher2_cloud_credential.foo-google.id
    region = <REGION> # Zone argument could also be used instead of region
    project_id = <PROJECT_ID>
    imported = true
  }
}
```

### Creating GKE cluster from Rancher v2, using `gke_config_v2`. For Rancher v2.5.8 or above.

**Note** At the moment, routed-based GKE clusters are not supported due to [rancher/issues/32585]](https://github.com/rancher/rancher/issues/32585)

```hcl
resource "rancher2_cloud_credential" "foo-google" {
  name = "foo-google"
  description= "Terraform cloudCredential acceptance test"
  google_credential_config {
    auth_encoded_json = file(<GOOGLE_AUTH_ENCODED_JSON>)
  }
}

resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform GKE cluster"
  gke_config_v2 {
    name = "foo"
    google_credential_secret = rancher2_cloud_credential.foo-google.id
    region = <REGION> # Zone argument could also be used instead of region
    project_id = <PROJECT_ID>
    kubernetes_version = <K8S_VERSION>
    network = <NETWORK>
    subnetwork = <SUBNET>
    node_pools {
      initial_node_count = 1
      max_pods_constraint = 110
      name = <NODE_POOL_NAME>
      version = <VERSION>
    }
  }
}
```


### Importing AKS cluster from Rancher v2, using `aks_config_v2`. For Rancher v2.6.0 or above.

```hcl
resource "rancher2_cloud_credential" "foo-aks" {
  name = "foo-aks"
  azure_credential_config {
    client_id = "<CLIENT_ID>"
    client_secret = "<CLIENT_SECRET>"
    subscription_id = "<SUBSCRIPTION_ID>"
  }
}
# For imported AKS clusters, don't add any other aks_config_v2 field
resource "rancher2_cluster" "foo" {
  name = <CLUSTER_NAME>
  description = "Terraform AKS cluster"
  aks_config_v2 {
    cloud_credential_id = rancher2_cloud_credential.foo-aks.id
    resource_group = "<RESOURCE_GROUP>"
    resource_location = "<RESOURCE_LOCATION>"
    imported = true
  }
}
```

### Creating AKS cluster from Rancher v2, using `aks_config_v2`. For Rancher v2.6.0 or above.

```hcl
resource "rancher2_cloud_credential" "foo-aks" {
  name = "foo-aks"
  azure_credential_config {
    client_id = "<CLIENT_ID>"
    client_secret = "<CLIENT_SECRET>"
    subscription_id = "<SUBSCRIPTION_ID>"
  }
}

resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform AKS cluster"
  aks_config_v2 {
    cloud_credential_id = rancher2_cloud_credential.foo-aks.id
    resource_group = "<RESOURCE_GROUP>"
    resource_location = "<RESOURCE_LOCATION>"
    dns_prefix = "<DNS_PREFIX>"
    kubernetes_version = "1.24.6"
    network_plugin = "<NETWORK_PLUGIN>"
    node_pools {
      availability_zones = ["1", "2", "3"]
      name = "<NODEPOOL_NAME>"
      count = 1
      orchestrator_version = "1.21.2"
      os_disk_size_gb = 128
      vm_size = "Standard_DS2_v2"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cluster (string)
* `agent_env_vars` - (Optional) Optional Agent Env Vars for Rancher agent. Just for Rancher v2.5.6 and above (list)
* `rke_config` - (Optional/Computed) The RKE configuration for `rke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` and `k3s_config` (list maxitems:1)
* `rke2_config` - (Optional/Computed) The RKE2 configuration for `rke2` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `gke_config`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* `k3s_config` - (Optional/Computed) The K3S configuration for `k3s` imported Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` and `rke_config` (list maxitems:1)
* `aks_config` - (Optional) The Azure AKS configuration for `aks` Clusters. Conflicts with `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config` (list maxitems:1)
* `aks_config_v2` - (Optional) The Azure AKS v2 configuration for creating/import `aks` Clusters. Conflicts with `aks_config`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config` (list maxitems:1)
* `eks_config` - (Optional) The Amazon EKS configuration for `eks` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config` (list maxitems:1)
* `eks_config_v2` - (Optional/Computed) The Amazon EKS V2 configuration to create or import `eks` Clusters. Conflicts with `aks_config`, `eks_config`, `gke_config`, `gke_config_v2`, `oke_config` `k3s_config` and `rke_config`. For Rancher v2.5.x or above (list maxitems:1)
* `gke_config` - (Optional) The Google GKE configuration for `gke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config_v2`, `oke_config`, `k3s_config` and `rke_config` (list maxitems:1)
* `gke_config_v2` - (Optional) The Google GKE V2 configuration for `gke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `oke_config`, `k3s_config` and `rke_config`. For Rancher v2.5.8 or above (list maxitems:1)
* `oke_config` - (Optional) The Oracle OKE configuration for `oke` Clusters. Conflicts with `aks_config`, `aks_config_v2`, `eks_config`, `eks_config_v2`, `gke_config`, `gke_config_v2`, `k3s_config` and `rke_config` (list maxitems:1)
* `description` - (Optional) The description for Cluster (string)
* `cluster_auth_endpoint` - (Optional/Computed) Enabling the [local cluster authorized endpoint](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#local-cluster-auth-endpoint) allows direct communication with the cluster, bypassing the Rancher API proxy. (list maxitems:1)
* `cluster_monitoring_input` - (Optional) Cluster monitoring config. Any parameter defined in [rancher-monitoring charts](https://github.com/rancher/system-charts/tree/dev/charts/rancher-monitoring) could be configured  (list maxitems:1)
* `cluster_template_answers` - (Optional/Computed) Cluster template answers. Just for Rancher v2.3.x and above (list maxitems:1)
* `cluster_template_id` - (Optional) Cluster template ID. Just for Rancher v2.3.x and above (string)
* `cluster_template_questions` - (Optional/Computed) Cluster template questions. Just for Rancher v2.3.x and above (list)
* `cluster_template_revision_id` - (Optional) Cluster template revision ID. Just for Rancher v2.3.x and above (string)
* `default_pod_security_policy_template_id` - (Optional/Computed) [Default pod security policy template id](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/options/#pod-security-policy-support) (string)
* `desired_agent_image` - (Optional/Computed) Desired agent image. Just for Rancher v2.3.x and above (string)
* `desired_auth_image` - (Optional/Computed) Desired auth image. Just for Rancher v2.3.x and above (string)
* `docker_root_dir` - (Optional/Computed) Desired auth image. Just for Rancher v2.3.x and above (string)
* `enable_cluster_alerting` - (Optional/Computed) Enable built-in cluster alerting (bool)
* `enable_cluster_monitoring` - (Optional/Computed) Enable built-in cluster monitoring (bool)
* `enable_cluster_istio` - (Deprecated) Deploy istio on `system` project and `istio-system` namespace, using rancher2_app resource instead. See above example.
* `enable_network_policy` - (Optional/Computed) Enable project network isolation (bool)
* `fleet_workspace_name` - (Optional/Computed) Fleet workspace name (string)
* `annotations` - (Optional/Computed) Annotations for the Cluster (map)
* `labels` - (Optional/Computed) Labels for the Cluster (map)
* `windows_prefered_cluster` - (Optional) Windows preferred cluster. Default: `false` (bool)


## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cluster_registration_token` - (Computed) Cluster Registration Token generated for the cluster (list maxitems:1)
* `default_project_id` - (Computed) Default project ID for the cluster (string)
* `driver` - (Computed) The driver used for the Cluster. `imported`, `azurekubernetesservice`, `amazonelasticcontainerservice`, `googlekubernetesengine` and `rancherKubernetesEngine` are supported (string)
* `istio_enabled` - (Computed) Is istio enabled at cluster? Just for Rancher v2.3.x and above (bool)
* `kube_config` - (Computed/Sensitive) Kube Config generated for the cluster. Note: For Rancher 2.6.0 and above, when the cluster has `cluster_auth_endpoint` enabled, the kube_config will not be available until the cluster is `connected` (string)
* `ca_cert` - (Computed/Sensitive) K8s cluster ca cert (string)
* `system_project_id` - (Computed) System project ID for the cluster (string)

**Note** For Rancher 2.6.0 and above: if setting `kubeconfig-generate-token=false` then the generated `kube_config` will not contain any user token. `kubectl` will generate the user token executing the [rancher cli](https://github.com/rancher/cli/releases/tag/v2.6.0), so it should be installed previously.

## Nested blocks

### `agent_env_vars`

#### Arguments

* `name` - (Required) Rancher agent env var name (string)
* `value` - (Required) Rancher agent env var value (string)

### `rke_config`

**Note** `rke_config` works the same as within Rancher GUI; it will not _provision_ hosts when not using `node_pool` nor `node_driver`. It is expected that nodes are registered by having the `node_command` run on each node. Running the `node_command` is outside the scope of this provider.

#### Arguments

* `addon_job_timeout` - (Optional/Computed) Duration in seconds of addon job (int)
* `addons` - (Optional) Addons descripton to deploy on RKE cluster.
* `addons_include` - (Optional) Addons yaml manifests to deploy on RKE cluster (list)
* `authentication` - (Optional/Computed) Kubernetes cluster authentication (list maxitems:1)
* `authorization` - (Optional/Computed) Kubernetes cluster authorization (list maxitems:1)
* `bastion_host` - (Optional/Computed) RKE bastion host (list maxitems:1)
* `cloud_provider` - (Optional/Computed) RKE cloud provider [rke-cloud-providers](https://rancher.com/docs/rke/v0.1.x/en/config-options/cloud-providers/) (list maxitems:1)
* `dns` - (Optional/Computed) RKE dns add-on. Just for Rancher v2.2.x (list maxitems:1)
* `enable_cri_dockerd` - (Optional) Enable/disable using cri-dockerd. Deafult: `false` [enable_cri_dockerd](https://rancher.com/docs/rke/latest/en/config-options/#cri-dockerd) (bool)
* `ignore_docker_version` - (Optional) Ignore docker version. Default `true` (bool)
* `ingress` - (Optional/Computed) Kubernetes ingress configuration (list maxitems:1)
* `kubernetes_version` - (Optional/Computed) K8s version to deploy. Default: `Rancher default` (string) (Note - if rke_config is set at cluster_template, kubernetes_version must be set to the active cluster version so Rancher can clone the RKE template)
* `monitoring` - (Optional/Computed) Kubernetes cluster monitoring (list maxitems:1)
* `network` - (Optional/Computed) Kubernetes cluster networking (list maxitems:1)
* `nodes` - (Optional) RKE cluster nodes (list)
* `prefix_path` - (Optional/Computed) Prefix to customize Kubernetes path (string)
* `win_prefix_path` - (Optional/Computed) Prefix to customize Kubernetes path for windows (string)
* `private_registries` - (Optional) private registries for docker images (list)
* `services` - (Optional/Computed) Kubernetes cluster services (list maxitems:1)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false`
* `ssh_cert_path` - (Optional/Computed) Cluster level SSH certificate path (string)
* `ssh_key_path` - (Optional/Computed) Cluster level SSH private key path (string)
* `upgrade_strategy` - (Optional/Computed) RKE upgrade strategy (list maxitems:1)

#### `authentication`

##### Arguments

* `sans` - (Optional/Computed) RKE sans for authentication ([]string)
* `strategy` - (Optional/Computed) RKE strategy for authentication (string)

#### `authorization`

##### Arguments

* `mode` - (Optional) RKE mode for authorization. `rbac` and `none` modes are available. Default `rbac` (string)
* `options` - (Optional/Computed) RKE options for authorization (map)

#### `bastion_host`

##### Arguments

* `address` - (Required) Address ip for the bastion host (string)
* `user` - (Required) User to connect bastion host (string)
* `port` - (Optional) Port for bastion host. Default `22` (string)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false` (bool)
* `ssh_key` - (Optional/Computed/Sensitive) Bastion host SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Bastion host SSH private key path (string)

#### `cloud_provider`

##### Arguments

* `aws_cloud_provider` - (Optional/Computed) RKE AWS Cloud Provider config for Cloud Provider [rke-aws-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/aws/) (list maxitems:1)
* `azure_cloud_provider` - (Optional/Computed) RKE Azure Cloud Provider config for Cloud Provider [rke-azure-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/azure/) (list maxitems:1)
* `custom_cloud_provider` - (Optional/Computed) RKE Custom Cloud Provider config for Cloud Provider (string)
* `name` - (Optional) RKE Cloud Provider name (string)
* `openstack_cloud_provider` - (Optional/Computed) RKE Openstack Cloud Provider config for Cloud Provider [rke-openstack-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/openstack/) (list maxitems:1)
* `vsphere_cloud_provider` - (Optional/Computed) RKE Vsphere Cloud Provider config for Cloud Provider [rke-vsphere-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/vsphere/) Extra argument `name` is required on `virtual_center` configuration. (list maxitems:1)

##### `aws_cloud_provider`

###### Arguments

* `global` - (Optional/Computed) (list maxitems:1)
* `service_override` - (Optional) (list)

###### `global`

###### Arguments

* `disable_security_group_ingress` - (Optional) Default `false` (bool)
* `disable_strict_zone_check` - (Optional) Default `false` (bool)
* `elb_security_group` - (Optional/Computed) (string)
* `kubernetes_cluster_id` - (Optional/Computed) (string)
* `kubernetes_cluster_tag` - (Optional/Computed) (string)
* `role_arn` - (Optional/Computed) (string)
* `route_table_id` - (Optional/Computed/Sensitive) (string)
* `subnet_id` - (Optional/Computed) (string)
* `vpc` - (Optional/Computed) (string)
* `zone` - (Optional/Computed) (string)

###### `service_override`

###### Arguments

* `service` - (Required) (string)
* `region` - (Optional/Computed) (string)
* `signing_method` - (Optional/Computed) (string)
* `signing_name` - (Optional/Computed) (string)
* `signing_region` - (Optional/Computed) (string)
* `url` - (Optional/Computed) (string)

##### `azure_cloud_provider`

###### Arguments

* `aad_client_id` - (Required/Sensitive) (string)
* `aad_client_secret` - (Required/Sensitive) (string)
* `subscription_id` - (Required/Sensitive) (string)
* `tenant_id` - (Required/Sensitive) (string)
* `aad_client_cert_password` - (Optional/Computed/Sensitive) (string)
* `aad_client_cert_path` - (Optional/Computed) (string)
* `cloud` - (Optional/Computed) (string)
* `cloud_provider_backoff` - (Optional/Computed) (bool)
* `cloud_provider_backoff_duration` - (Optional/Computed) (int)
* `cloud_provider_backoff_exponent` - (Optional/Computed) (int)
* `cloud_provider_backoff_jitter` - (Optional/Computed) (int)
* `cloud_provider_backoff_retries` - (Optional/Computed) (int)
* `cloud_provider_rate_limit` - (Optional/Computed) (bool)
* `cloud_provider_rate_limit_bucket` - (Optional/Computed) (int)
* `cloud_provider_rate_limit_qps` - (Optional/Computed) (int)
* `load_balancer_sku` - (Optional) Allowed values: `basic` (default) `standard` (string)
* `location` - (Optional/Computed) (string)
* `maximum_load_balancer_rule_count` - (Optional/Computed) (int)
* `primary_availability_set_name` - (Optional/Computed) (string)
* `primary_scale_set_name` - (Optional/Computed) (string)
* `resource_group` - (Optional/Computed) (string)
* `route_table_name` - (Optional/Computed) (string)
* `security_group_name` - (Optional/Computed) (string)
* `subnet_name` - (Optional/Computed) (string)
* `use_instance_metadata` - (Optional/Computed) (bool)
* `use_managed_identity_extension` - (Optional/Computed) (bool)
* `vm_type` - (Optional/Computed) (string)
* `vnet_name` - (Optional/Computed) (string)
* `vnet_resource_group` - (Optional/Computed) (string)

##### `openstack_cloud_provider`

###### Arguments

* `global` - (Required) (list maxitems:1)
* `block_storage` - (Optional/Computed) (list maxitems:1)
* `load_balancer` - (Optional/Computed) (list maxitems:1)
* `metadata` - (Optional/Computed) (list maxitems:1)
* `route` - (Optional/Computed) (list maxitems:1)

###### `global`

###### Arguments

* `auth_url` - (Required) (string)
* `password` - (Required/Sensitive) (string)
* `username` - (Required/Sensitive) (string)
* `ca_file` - (Optional/Computed) (string)
* `domain_id` - (Optional/Computed/Sensitive) Required if `domain_name` not provided. (string)
* `domain_name` - (Optional/Computed) Required if `domain_id` not provided. (string)
* `region` - (Optional/Computed) (string)
* `tenant_id` - (Optional/Computed/Sensitive) Required if `tenant_name` not provided. (string)
* `tenant_name` - (Optional/Computed) Required if `tenant_id` not provided. (string)
* `trust_id` - (Optional/Computed/Sensitive) (string)

###### `block_storage`

###### Arguments

* `bs_version` - (Optional/Computed) (string)
* `ignore_volume_az` - (Optional/Computed) (string)
* `trust_device_path` - (Optional/Computed) (string)

###### `load_balancer`

###### Arguments

* `create_monitor` - (Optional/Computed) (bool)
* `floating_network_id` - (Optional/Computed) (string)
* `lb_method` - (Optional/Computed) (string)
* `lb_provider` - (Optional/Computed) (string)
* `lb_version` - (Optional/Computed) (string)
* `manage_security_groups` - (Optional/Computed) (bool)
* `monitor_delay` - (Optional/Computed) Default `60s` (string)
* `monitor_max_retries` - (Optional/Computed) Default 5 (int)
* `monitor_timeout` - (Optional/Computed) Default `30s` (string)
* `subnet_id` - (Optional/Computed) (string)
* `use_octavia` - (Optional/Computed) (bool)

###### `metadata`

###### Arguments

* `request_timeout` - (Optional/Computed) (int)
* `search_order` - (Optional/Computed) (string)

###### `route`

###### Arguments

* `router_id` - (Optional/Computed) (string)

##### `vsphere_cloud_provider`

###### Arguments

* `virtual_center` - (Required) (List)
* `workspace` - (Required) (list maxitems:1)
* `disk` - (Optional/Computed) (list maxitems:1)
* `global` - (Optional/Computed) (list maxitems:1)
* `network` - (Optional/Computed) (list maxitems:1)

###### `virtual_center`

###### Arguments

* `datacenters` - (Required) (string)
* `name` - (Required) Name of virtualcenter config for Vsphere Cloud Provider config (string)
* `password` - (Required/Sensitive) (string)
* `user` - (Required/Sensitive) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

###### `workspace`

###### Arguments

* `datacenter` - (Required) (string)
* `folder` - (Required) (string)
* `server` - (Required) (string)
* `default_datastore` - (Optional/Computed) (string)
* `resourcepool_path` - (Optional/Computed) (string)

###### `disk`

###### Arguments

* `scsi_controller_type` - (Optional/Computed) (string)

###### `global`

###### Arguments

* `datacenters` - (Optional/Computed) (string)
* `insecure_flag` - (Optional/Computed) (bool)
* `password` - (Optional/Computed) (string)
* `user` - (Optional/Computed) (string)
* `port` - (Optional/Computed) (string)
* `soap_roundtrip_count` - (Optional/Computed) (int)

###### `network`

###### Arguments

* `public_network` - (Optional/Computed) (string)

#### `dns`

##### Arguments

* `nodelocal` - (Optional) Nodelocal dns config  (list Maxitem: 1)
* `linear_autoscaler_params` - (Optional) LinearAutoScalerParams dns config (list Maxitem: 1)
* `node_selector` - (Optional/Computed) DNS add-on node selector (map)
* `options` - (Optional/Computed) DNS add-on options (map)
* `provider` - (Optional) DNS add-on provider. `kube-dns`, `coredns` (default), and `none` are supported (string)
* `tolerations` - (Optional) DNS add-on tolerations (list)
* `reverse_cidrs` - (Optional/Computed) DNS add-on reverse cidr  (list)
* `upstream_nameservers` - (Optional/Computed) DNS add-on upstream nameservers  (list)
* `update_strategy` - (Optional) DNS update strategy (list Maxitems: 1)

##### `nodelocal`

###### Arguments

* `ip_address` - (required) Nodelocal dns ip address (string)
* `node_selector` - (Optional) Node selector key pair (map)

##### `linear_autoscaler_params`

###### Arguments

* `cores_per_replica` - (Optional) number of replicas per cluster cores (float64)
* `nodes_per_replica` - (Optional) number of replica per cluster nodes (float64)
* `max` - (Optional) maximum number of replicas (int64)
* `min` - (Optional) minimum number of replicas (int64)
* `prevent_single_point_failure` - (Optional) prevent single point of failure

##### `tolerations`

###### Arguments

* `key` - (Required) The toleration key (string)
* `effect` - (Optional) The toleration effect. `NoExecute`, `NoSchedule`, and `PreferNoSchedule` are supported. Default: `NoExecute` (string)
* `operator` - (Optional) The toleration operator. `Equal`, and `Exists` are supported. Default: `Equal` (string)
* `seconds` - (Optional) The toleration seconds (int)
* `value` - (Optional) The toleration value (string)

#### `ingress`

##### Arguments

* `default_backend` - (Optional) Enable ingress default backend. Default: `true` (bool)
* `dns_policy` - (Optional/Computed) Ingress controller DNS policy. `ClusterFirstWithHostNet`, `ClusterFirst`, `Default`, and `None` are supported. [K8S dns Policy](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy) (string)
* `extra_args` - (Optional/Computed) Extra arguments for RKE Ingress (map)
* `http_port` - (Optional/Computed) HTTP port for RKE Ingress (int)
* `https_port` - (Optional/Computed) HTTPS port for RKE Ingress (int)
* `network_mode` - (Optional/Computed) Network mode for RKE Ingress (string)
* `node_selector` - (Optional/Computed) Node selector for RKE Ingress (map)
* `options` - (Optional/Computed) RKE options for Ingress (map)
* `provider` - (Optional/Computed) Provider for RKE Ingress (string)
* `tolerations` - (Optional) Ingress add-on tolerations (list)
* `update_strategy` - (Optional) RKE ingress update strategy (list Maxitems: 1)

##### `tolerations`

###### Arguments

* `key` - (Required) The toleration key (string)
* `effect` - (Optional) The toleration effect. `NoExecute`, `NoSchedule`, and `PreferNoSchedule` are supported. Default: `NoExecute` (string)
* `operator` - (Optional) The toleration operator. `Equal`, and `Exists` are supported. Default: `Equal` (string)
* `seconds` - (Optional) The toleration seconds (int)
* `value` - (Optional) The toleration value (string)

##### `update_strategy`

###### Arguments

* `rolling_update` - (Optional) Monitoring daemon set rolling update (list Maxitems: 1)
* `strategy` - (Optional) Monitoring daemon set update strategy (string)

###### `rolling_update`

###### Arguments

* `max_unavailable` - (Optional) Monitoring deployment rolling update max unavailable. Default: `1` (int)

#### `monitoring`

##### Arguments

* `node_selector` - (Optional) RKE monitoring node selector (map)
* `options` - (Optional/Computed) RKE options for monitoring (map)
* `provider` - (Optional/Computed) RKE monitoring provider (string)
* `replicas` - (Optional/Computed) RKE monitoring replicas (int)
* `tolerations` - (Optional) RKE monitoring tolerations (list)
* `update_strategy` - (Optional) RKE monitoring update strategy (list Maxitems: 1)

##### `tolerations`

###### Arguments

* `key` - (Required) The toleration key (string)
* `effect` - (Optional) The toleration effect. `NoExecute`, `NoSchedule`, and `PreferNoSchedule` are supported. Default: `NoExecute` (string)
* `operator` - (Optional) The toleration operator. `Equal`, and `Exists` are supported. Default: `Equal` (string)
* `seconds` - (Optional) The toleration seconds (int)
* `value` - (Optional) The toleration value (string)

##### `update_strategy`

###### Arguments

* `rolling_update` - (Optional) Monitoring deployment rolling update (list Maxitems: 1)
* `strategy` - (Optional) Monitoring deployment update strategy (string)

###### `rolling_update`

###### Arguments

* `max_surge` - (Optional) Monitoring deployment rolling update max surge. Default: `1` (int)
* `max_unavailable` - (Optional) Monitoring deployment rolling update max unavailable. Default: `1` (int)

#### `network`

##### Arguments

* `aci_network_provider` - (Optional/Computed) ACI provider config for RKE network (list maxitems:63)
* `calico_network_provider` - (Optional/Computed) Calico provider config for RKE network (list maxitems:1)
* `canal_network_provider` - (Optional/Computed) Canal provider config for RKE network (list maxitems:1)
* `flannel_network_provider` - (Optional/Computed) Flannel provider config for RKE network (list maxitems:1)
* `weave_network_provider` - (Optional/Computed) Weave provider config for RKE network (list maxitems:1)
* `mtu` - (Optional) Network provider MTU. Default `0` (int)
* `options` - (Optional/Computed) RKE options for network (map)
* `plugin` - (Optional/Computed) Plugin for RKE network. `canal` (default), `flannel`, `calico`, `none` and `weave` are supported. (string)
* `tolerations` - (Optional) Network add-on tolerations (list)

##### `aci_network_provider`

###### Arguments

* `aep` - (Required) Attachable entity profile (string)
* `apic_hosts` - (Required) List of APIC hosts to connect for APIC API (list)
* `apic_refresh_ticker_adjust` - (Optional) APIC refresh ticker adjust amount (string)
* `apic_refresh_time` - (Optional) APIC refresh time in seconds (string)
* `apic_subscription_delay` - (Optional) APIC subscription delay amount (string)
* `apic_user_crt` - (Required/Sensitive) APIC user certificate (string)
* `apic_user_key` - (Required/Sensitive) APIC user key (string)
* `apic_user_name` - (Required) APIC user name (string)
* `capic` - (Optional) cAPIC cloud (string)
* `controller_log_level` - (Optional) Log level for ACI controller (string)
* `disable_periodic_snat_global_info_sync` - (Optional) Whether to disable periodic SNAT global info sync (string)
* `disable_wait_for_network` - (Optional) Whether to disable waiting for network (string)
* `drop_log_enable` - (Optional) Whether to enable drop log (string)
* `duration_wait_for_network` - (Optional) The duration to wait for network (string)
* `extern_dynamic` - (Required) Subnet to use for dynamic external IPs (string)
* `enable_endpoint_slice` - (Optional) Whether to enable endpoint slices (string)
* `encap_type` - (Required) Encap type: vxlan or vlan (string)
* `ep_registry` - (Optional) EP registry (string)
* `gbp_pod_subnet` - (Optional) GBH pod subnet (string)
* `host_agent_log_level` - (Optional) Log level for ACI host agent (string)
* `image_pull_policy` - (Optional) Image pull policy (string)
* `image_pull_secret` - (Optional) Image pull policy (string)
* `infra_vlan` - (Optional) The VLAN used by ACI infra (string)
* `install_istio` - (Optional) Whether to install Istio (string)
* `istio_profile` - (Optional) Istio profile name (string)
* `kafka_brokers` - (Optional) List of Kafka broker hosts (list)
* `kafka_client_crt` - (Optional) Kafka client certificate (string)
* `kafka_client_key` - (Optional) Kafka client key (string)
* `kube_api_vlan` - (Required) The VLAN used by the physdom for nodes (string)
* `l3out` - (Required) L3out (string)
* `l3out_external_networks` - (Required) L3out external networks (list)
* `max_nodes_svc_graph` - (Optional) Max nodes in service graph (string)
* `mcast_range_end` - (Required) End of mcast range (string)
* `mcast_range_start` - (Required) Start of mcast range (string)
* `mtu_head_room` - (Optional) MTU head room amount (string)
* `multus_disable` - (Optional) Whether to disable Multus (string)
* `no_priority_class` - (Optional) Whether to use priority class (string)
* `node_pod_if_enable` - (Optional) Whether to enable node pod interface (string)
* `node_subnet` - (Required) Subnet to use for nodes (string)
* `ovs_memory_limit` - (Optional) OVS memory limit (string)
* `opflex_log_level` - (Optional) Log level for ACI opflex (string)
* `opflex_client_ssl` - (Optional) Whether to use client SSL for Opflex (string)
* `opflex_device_delete_timeout` - (Optional) Opflex device delete timeout (string)
* `opflex_mode` - (Optional) Opflex mode (string)
* `opflex_server_port` - (Optional) Opflex server port (string)
* `overlay_vrf_name` - (Optional) Overlay VRF name (string)
* `pbr_tracking_non_snat` - (Optional) Policy-based routing tracking non snat (string)
* `pod_subnet_chunk_size` - (Optional) Pod subnet chunk size (string)
* `run_gbp_container` - (Optional) Whether to run GBP container (string)
* `run_opflex_server_container` - (Optional) Whether to run Opflex server container (string)
* `node_svc_subnet` - (Required) Subnet to use for service graph (string)
* `service_monitor_interval` - (Optional) Service monitor interval (string)
* `service_vlan` - (Required) The VLAN used by LoadBalancer services (string)
* `snat_contract_scope` - (Optional) Snat contract scope (string)
* `snat_namespace` - (Optional) Snat namespace (string)
* `snat_port_range_end` - (Optional) End of snat port range (string)
* `snat_port_range_start` - (Optional) End of snat port range (string)
* `snat_ports_per_node` - (Optional) Snat ports per node (string)
* `sriov_enable` - (Optional) Whether to enable SR-IOV (string)
* `extern_static` - (Required)  Subnet to use for static external IPs (string)
* `subnet_domain_name` - (Optional) Subnet domain name (string)
* `system_id` - (Required) ACI system ID (string)
* `tenant` - (Optional) ACI tenant (string)
* `token` - (Required/Sensitive) ACI token (string)
* `use_aci_anywhere_crd` - (Optional) Whether to use ACI anywhere CRD (string)
* `use_aci_cni_priority_class` - (Optional) Whether to use ACI CNI priority class (string)
* `use_cluster_role` - (Optional) Whether to use cluster role (string)
* `use_host_netns_volume` - (Optional) Whether to use host netns volume (string)
* `use_opflex_server_volume` - (Optional) Whether use Opflex server volume (string)
* `use_privileged_container` - (Optional) Whether ACI containers should run as privileged (string)
* `vrf_name` - (Required) VRF name (string)
* `vrf_tenant` - (Required) VRF tenant (string)
* `vmm_controller` - (Optional) VMM controller configuration (string)
* `vmm_domain` - (Optional) VMM domain configuration (string)

##### `calico_network_provider`

###### Arguments

* `cloud_provider` - (Optional/Computed) RKE options for Calico network provider (string)

##### `canal_network_provider`

###### Arguments

* `iface` - (Optional/Computed) Iface config Canal network provider (string)

##### `flannel_network_provider`

###### Arguments

* `iface` - (Optional/Computed) Iface config Flannel network provider (string)

##### `weave_network_provider`

###### Arguments

* `password` - (Optional/Computed) Password config Weave network provider (string)

##### `tolerations`

###### Arguments

* `key` - (Required) The toleration key (string)
* `effect` - (Optional) The toleration effect. `NoExecute`, `NoSchedule`, and `PreferNoSchedule` are supported. Default: `NoExecute` (string)
* `operator` - (Optional) The toleration operator. `Equal`, and `Exists` are supported. Default: `Equal` (string)
* `seconds` - (Optional) The toleration seconds (int)
* `value` - (Optional) The toleration value (string)

#### `nodes`

##### Arguments

* `address` - (Required) Address ip for node (string)
* `role` - (Requires) Roles for the node. `controlplane`, `etcd` and `worker` are supported. (list)
* `user` - (Required/Sensitive) User to connect node (string)
* `docker_socket` - (Optional/Computed) Docker socket for node (string)
* `hostname_override` - (Optional) Hostname override for node (string)
* `internal_address` - (Optional) Internal ip for node (string)
* `labels` - (Optional) Labels for the node (map)
* `node_id` - (Optional) Id for the node (string)
* `port` - (Optional) Port for node. Default `22` (string)
* `ssh_agent_auth` - (Optional) Use ssh agent auth. Default `false` (bool)
* `ssh_key` - (Optional/Computed/Sensitive) Node SSH private key (string)
* `ssh_key_path` - (Optional/Computed) Node SSH private key path (string)

#### `private_registries`

##### Arguments

* `url` - (Required) Registry URL (string)
* `ecr_credential_plugin` - (Optional) ECR credential plugin config (list maxitems:1)
* `is_default` - (Optional) Set as default registry. Default `false` (bool)
* `password` - (Optional/Sensitive) Registry password (string)
* `user` - (Optional/Sensitive) Registry user (string)

##### `ecr_credential_plugin`

###### Arguments

* `aws_access_key_id` - (Optional) AWS access key ID (string)
* `aws_secret_access_key` - (Optional/Sensitive) AWS secret access key (string)
* `aws_session_token` - (Optional/Sensitive) AWS session token (string)

#### `services`

##### Arguments

* `etcd` - (Optional/Computed) Etcd options for RKE services (list maxitems:1)
* `kube_api` - (Optional/Computed) Kube API options for RKE services (list maxitems:1)
* `kube_controller` - (Optional/Computed) Kube Controller options for RKE services (list maxitems:1)
* `kubelet` - (Optional/Computed) Kubelet options for RKE services (list maxitems:1)
* `kubeproxy` - (Optional/Computed) Kubeproxy options for RKE services (list maxitems:1)
* `scheduler` - (Optional/Computed) Scheduler options for RKE services (list maxitems:1)

##### `etcd`

###### Arguments

* `backup_config` - (Optional/Computed) Backup options for etcd service. Just for Rancher v2.2.x (list maxitems:1)
* `ca_cert` - (Optional/Computed) TLS CA certificate for etcd service (string)
* `cert` - (Optional/Computed/Sensitive) TLS certificate for etcd service (string)
* `creation` - (Optional/Computed) Creation option for etcd service (string)
* `external_urls` - (Optional) External urls for etcd service (list)
* `extra_args` - (Optional/Computed) Extra arguments for etcd service (map)
* `extra_binds` - (Optional) Extra binds for etcd service (list)
* `extra_env` - (Optional) Extra environment for etcd service (list)
* `gid` - (Optional) Etcd service GID. Default: `0`. For Rancher v2.3.x or above (int)
* `image` - (Optional/Computed) Docker image for etcd service (string)
* `key` - (Optional/Computed/Sensitive) TLS key for etcd service (string)
* `path` - (Optional/Computed) Path for etcd service (string)
* `retention` - (Optional/Computed) Retention option for etcd service (string)
* `snapshot` - (Optional/Computed) Snapshot option for etcd service (bool)
* `uid` - (Optional) Etcd service UID. Default: `0`. For Rancher v2.3.x or above (int)

###### `backup_config`

###### Arguments

* `enabled` - (Optional) Enable etcd backup (bool)
* `interval_hours` - (Optional) Interval hours for etcd backup. Default `12` (int)
* `retention` - (Optional) Retention for etcd backup. Default `6` (int)
* `s3_backup_config` - (Optional) S3 config options for etcd backup (list maxitems:1)
* `safe_timestamp` - (Optional) Safe timestamp for etcd backup. Default: `false` (bool)
* `timeout` - (Optional/Computed) Timeout in seconds for etcd backup. Default: `300`. Just for Rancher v2.5.6 and above (int)

###### `s3_backup_config`

###### Arguments

* `access_key` - (Optional/Sensitive) Access key for S3 service (string)
* `bucket_name` - (Required) Bucket name for S3 service (string)
* `custom_ca` - (Optional) Base64 encoded custom CA for S3 service. Use filebase64(<FILE>) for encoding file. Available from Rancher v2.2.5 (string)
* `endpoint` - (Required) Endpoint for S3 service (string)
* `folder` - (Optional) Folder for S3 service. Available from Rancher v2.2.7 (string)
* `region` - (Optional) Region for S3 service (string)
* `secret_key` - (Optional/Sensitive) Secret key for S3 service (string)

##### `kube_api`

###### Arguments

* `admission_configuration` - (Optional) Admission configuration (map)
* `always_pull_images` - (Optional) Enable [AlwaysPullImages](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#alwayspullimages) Admission controller plugin. [Rancher docs](https://rancher.com/docs/rke/latest/en/config-options/services/#kubernetes-api-server-options) Default: `false` (bool)
* `audit_log` - (Optional) K8s audit log configuration. (list maxitems: 1)
* `event_rate_limit` - (Optional) K8s event rate limit configuration. (list maxitems: 1)
* `extra_args` - (Optional/Computed) Extra arguments for kube API service (map)
* `extra_binds` - (Optional) Extra binds for kube API service (list)
* `extra_env` - (Optional) Extra environment for kube API service (list)
* `image` - (Optional/Computed) Docker image for kube API service (string)
* `pod_security_policy` - (Optional) Pod Security Policy option for kube API service. Default `false` (bool)
* `secrets_encryption_config` - (Optional) [Encrypt k8s secret data configration](https://rancher.com/docs/rke/latest/en/config-options/secrets-encryption/). (list maxitem: 1)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster IP Range option for kube API service (string)
* `service_node_port_range` - (Optional/Computed) Service Node Port Range option for kube API service (string)

###### `audit_log`

###### Arguments

* `configuration` - (Optional) Audit log configuration. (list maxitems: 1)
* `enabled` - (Optional) Enable audit log. Default: `false` (bool)

###### `configuration`

###### Arguments

* `format` - (Optional) Audit log format. Default: 'json' (string)
* `max_age` - (Optional) Audit log max age. Default: `30` (int)
* `max_backup` - (Optional) Audit log max backup. Default: `10` (int)
* `max_size` - (Optional) Audit log max size. Default: `100` (int)
* `path` - (Optional) (Optional) Audit log path. Default: `/var/log/kube-audit/audit-log.json` (string)
* `policy` - (Optional/Computed) Audit policy yaml encoded definition. `apiVersion` and `kind: Policy\nrules:"` fields are required in the yaml. [More info](https://rancher.com/docs/rke/latest/en/config-options/audit-log/) (string) Ex:

```
policy = <<EOF
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: RequestResponse
  resources:
  - resources:
    - pods
EOF
```

###### `event_rate_limit`

###### Arguments

* `configuration` - (Optional) Event rate limit configuration yaml encoded definition. `apiVersion` and `kind: Configuration"` fields are required in the yaml. [More info](https://rancher.com/docs/rke/latest/en/config-options/rate-limiting/) (string) Ex:

```
configuration = <<EOF
apiVersion: eventratelimit.admission.k8s.io/v1alpha1
kind: Configuration
limits:
- type: Server
  burst: 35000
  qps: 6000
EOF
```

* `enabled` - (Optional) Enable event rate limit. Default: `false` (bool)

###### `secrets_encryption_config`

###### Arguments

* `enabled` - (Optional) Enable secrets encryption. Default: `false` (bool)
* `custom_config` - (Optional) Secrets encryption yaml encoded custom configuration. `"apiVersion"` and `"kind":"EncryptionConfiguration"` fields are required in the yaml. [More info](https://rancher.com/docs/rke/latest/en/config-options/secrets-encryption/) (string) Ex:

```
custom_config = <<EOF
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
- resources:
  - secrets
  providers:
  - aescbc:
      keys:
      - name: k-fw5hn
        secret: RTczRjFDODMwQzAyMDVBREU4NDJBMUZFNDhCNzM5N0I=
    identity: {}
EOF

```

##### `kube_controller`

###### Arguments

* `cluster_cidr` - (Optional/Computed) Cluster CIDR option for kube controller service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kube controller service (map)
* `extra_binds` - (Optional) Extra binds for kube controller service (list)
* `extra_env` - (Optional) Extra environment for kube controller service (list)
* `image` - (Optional/Computed) Docker image for kube controller service (string)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster ip Range option for kube controller service (string)

##### `kubelet`

###### Arguments

* `cluster_dns_server` - (Optional/Computed) Cluster DNS Server option for kubelet service (string)
* `cluster_domain` - (Optional/Computed) Cluster Domain option for kubelet service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kubelet service (map)
* `extra_binds` - (Optional) Extra binds for kubelet service (list)
* `extra_env` - (Optional) Extra environment for kubelet service (list)
* `fail_swap_on` - (Optional/Computed) Enable or disable failing when swap on is not supported (bool)
* `generate_serving_certificate` [Generate a certificate signed by the kube-ca](https://rancher.com/docs/rke/latest/en/config-options/services/#kubelet-serving-certificate-requirements). Default `false` (bool)
* `image` - (Optional/Computed) Docker image for kubelet service (string)
* `infra_container_image` - (Optional/Computed) Infra container image for kubelet service (string)

##### `kubeproxy`

###### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for kubeproxy service (map)
* `extra_binds` - (Optional) Extra binds for kubeproxy service (list)
* `extra_env` - (Optional) Extra environment for kubeproxy service (list)
* `image` - (Optional/Computed) Docker image for kubeproxy service (string)

##### `scheduler`

###### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for scheduler service (map)
* `extra_binds` - (Optional) Extra binds for scheduler service (list)
* `extra_env` - (Optional) Extra environment for scheduler service (list)
* `image` - (Optional/Computed) Docker image for scheduler service (string)

#### `upgrade_strategy`

##### Arguments

* `drain` - (Optional) RKE drain nodes. Default: `false` (bool)
* `drain_input` - (Optional/Computed) RKE drain node input (list Maxitems: 1)
* `max_unavailable_controlplane` - (Optional) RKE max unavailable controlplane nodes. Default: `1` (string)
* `max_unavailable_worker` - (Optional) RKE max unavailable worker nodes. Default: `10%` (string)

##### `drain_input`

###### Arguments

* `delete_local_data` - Delete RKE node local data. Default: `false` (bool)
* `force` - Force RKE node drain. Default: `false` (bool)
* `grace_period` - RKE node drain grace period. Default: `-1` (int)
* `ignore_daemon_sets` - Ignore RKE daemon sets. Default: `true` (bool)
* `timeout` - RKE node drain timeout. Default: `60` (int)

### `rke2_config`

#### Arguments

The following arguments are supported:

* `upgrade_strategy` - (Optional/Computed) RKE2 upgrade strategy (List maxitems: 1)
* `version` - (Optional/Computed) RKE2 kubernetes version (string)

#### `upgrade_strategy`

##### Arguments

* `drain_server_nodes` - (Optional) Drain server nodes. Default: `false` (bool)
* `drain_worker_nodes` - (Optional) Drain worker nodes. Default: `false` (bool)
* `server_concurrency` - (Optional) Server concurrency. Default: `1` (int)
* `worker_concurrency` - (Optional) Worker concurrency. Default: `1` (int)

### `k3s_config`

#### Arguments

The following arguments are supported:

* `upgrade_strategy` - (Optional/Computed) K3S upgrade strategy (List maxitems: 1)
* `version` - (Optional/Computed) K3S kubernetes version (string)

#### `upgrade_strategy`

##### Arguments

* `drain_server_nodes` - (Optional) Drain server nodes. Default: `false` (bool)
* `drain_worker_nodes` - (Optional) Drain worker nodes. Default: `false` (bool)
* `server_concurrency` - (Optional) Server concurrency. Default: `1` (int)
* `worker_concurrency` - (Optional) Worker concurrency. Default: `1` (int)

### `aks_config`

#### Arguments

The following arguments are supported:

* `agent_dns_prefix` - (Required) DNS prefix to be used to create the FQDN for the agent pool (string)
* `client_id` - (Required/Sensitive) Azure client ID to use (string)
* `client_secret` - (Required/Sensitive) Azure client secret associated with the \"client id\" (string)
* `kubernetes_version` - (Required) Specify the version of Kubernetes. To check available versions exec `az aks get-versions -l eastus -o table` (string)
* `master_dns_prefix` - (Required) DNS prefix to use the Kubernetes cluster control pane (string)
* `resource_group` - (Required) The name of the Cluster resource group (string)
* `ssh_public_key_contents` - (Required) Contents of the SSH public key used to authenticate with Linux hosts (string)
* `subnet` - (Required) The name of an existing Azure Virtual Subnet. Composite of agent virtual network subnet ID (string)
* `subscription_id` - (Required) Subscription credentials which uniquely identify Microsoft Azure subscription (string)
* `tenant_id` - (Required) Azure tenant ID to use (string)
* `virtual_network` - (Required) The name of an existing Azure Virtual Network. Composite of agent virtual network subnet ID (string)
* `virtual_network_resource_group` - (Required) The resource group of an existing Azure Virtual Network. Composite of agent virtual network subnet ID (string)
* `add_client_app_id` - (Optional/Sensitive) The ID of an Azure Active Directory client application of type \"Native\". This application is for user login via kubectl (string)
* `add_server_app_id` - (Optional/Sensitive) The ID of an Azure Active Directory server application of type \"Web app/API\". This application represents the managed cluster's apiserver (Server application) (string)
* `aad_server_app_secret` - (Optional/Sensitive) The secret of an Azure Active Directory server application (string)
* `aad_tenant_id` - (Optional/Sensitive) The ID of an Azure Active Directory tenant (string)
* `admin_username` - (Optional) The administrator username to use for Linux hosts. Default `azureuser` (string)
* `agent_os_disk_size` - (Optional) GB size to be used to specify the disk for every machine in the agent pool. If you specify 0, it will apply the default according to the \"agent vm size\" specified. Default `0` (int)
* `agent_pool_name` - (Optional) Name for the agent pool, upto 12 alphanumeric characters. Default `agentpool0` (string)
* `agent_storage_profile` - (Optional) Storage profile specifies what kind of storage used on machine in the agent pool. Chooses from [ManagedDisks StorageAccount]. Default `ManagedDisks` (string)
* `agent_vm_size` - (Optional) Size of machine in the agent pool. Default `Standard_D1_v2` (string)
* `auth_base_url` - (Optional) Different authentication API url to use. Default `https://login.microsoftonline.com/` (string)
* `base_url` - (Optional) Different resource management API url to use. Default `https://management.azure.com/` (string)
* `count` - (Optional) Number of machines (VMs) in the agent pool. Allowed values must be in the range of 1 to 100 (inclusive). Default `1` (int)
* `dns_service_ip` - (Optional) An IP address assigned to the Kubernetes DNS service. It must be within the Kubernetes Service address range specified in \"service cidr\". Default `10.0.0.10` (string)
* `docker_bridge_cidr` - (Required) A CIDR notation IP range assigned to the Docker bridge network. It must not overlap with any Subnet IP ranges or the Kubernetes Service address range specified in \"service cidr\". Default `172.17.0.1/16` (string)
* `enable_http_application_routing` - (Optional) Enable the Kubernetes ingress with automatic public DNS name creation. Default `false` (bool)
* `enable_monitoring` - (Optional) Turn on Azure Log Analytics monitoring. Uses the Log Analytics \"Default\" workspace if it exists, else creates one. if using an existing workspace, specifies \"log analytics workspace resource id\". Default `true` (bool)
* `load_balancer_sku` - (Optional/Computed) Load balancer type (basic | standard). Must be standard for auto-scaling
* `location` - (Optional) Azure Kubernetes cluster location. Default `eastus` (string)
* `log_analytics_workspace` - (Optional) The name of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses '{resource group}-{subscription id}-{location code}' (string)
* `log_analytics_workspace_resource_group` - (Optional) The resource group of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses the 'Cluster' resource group (string)
* `max_pods` - (Optional) Maximum number of pods that can run on a node. Default `110` (int)
* `network_plugin` - (Optional) Network plugin used for building Kubernetes network. Chooses from `azure` or `kubenet`. Default `azure` (string)
* `network_policy` - (Optional) Network policy used for building Kubernetes network. Chooses from `calico` (string)
* `pod_cidr` - (Optional) A CIDR notation IP range from which to assign Kubernetes Pod IPs when \"network plugin\" is specified in \"kubenet\". Default `172.244.0.0/16` (string)
* `service_cidr` - (Optional) A CIDR notation IP range from which to assign Kubernetes Service cluster IPs. It must not overlap with any Subnet IP ranges. Default `10.0.0.0/16` (string)
* `tag` - (Deprecated) Use `tags` argument instead as []string
* `tags` - (Optional/Computed) Tags for Kubernetes cluster. For example, `["foo=bar","bar=foo"]` (list)

### `aks_config_v2`

#### Arguments

The following arguments are supported:

* `cloud_credential_id` - (Required) The AKS Cloud Credential ID to use (string)
* `resource_group` - (Required) The AKS resource group (string)
* `resource_location` - (Required) The AKS resource location (string)
* `imported` - (Optional) Is AKS cluster imported? Defaul: `false` (bool)

The following arguments are supported just for creating new AKS clusters (`imported=false`): 

* `node_pools` - (Optional) The AKS nnode pools. Required if `imported=false` (list)
* `dns_prefix` - (Optional/ForceNew) The AKS dns prefix. Required if `imported=false` (string)
* `kubernetes_version` - (Optional) The kubernetes master version. Required if `imported=false` (string)
* `network_plugin` - (Optional) The AKS network plugin. Required if `imported=false` (string)
* `name` - (Optional/Computed) The AKS cluster name (string)
* `auth_base_url` - (Optional) The AKS auth base url (string)
* `authorized_ip_ranges` - (Optional) The AKS authorized ip ranges (list)
* `base_url` - (Optional) The AKS base url (string)
* `http_application_routing` - (Optional/Computed) Enable AKS http application routing? (bool)
* `linux_admin_username` - (Optional/Computed) The AKS linux admin username (string)
* `linux_ssh_public_key` - (Optional/Computed) The AKS linux ssh public key (string)
* `load_balancer_sku` - (Optional/Computed) The AKS load balancer sku (string)
* `log_analytics_workspace_group` - (Optional/Computed) The AKS log analytics workspace group (string)
* `log_analytics_workspace_name` - (Optional/Computed) The AKS log analytics workspace name (string)
* `monitoring` - (Optional/Computed) Is AKS cluster monitoring enabled? (bool)
* `network_dns_service_ip` - (Optional/Computed) The AKS network dns service ip (string)
* `network_docker_bridge_cidr` - (Optional/Computed) The AKS network docker bridge cidr (string)
* `network_pod_cidr` - (Optional/Computed) The AKS network pod cidr (string)
* `network_policy` - (Optional/Computed) The AKS network policy (string)
* `network_service_cidr` - (Optional/Computed) The AKS network service cidr (string)
* `private_cluster` - (Optional/Computed) Is AKS cluster private? (bool)
* `subnet` - (Optional/Computed) The AKS subnet (string)
* `tags` - (Optional/Computed) The AKS cluster tags" (map)
* `virtual_network` - (Optional/Computed) The AKS virtual network (string)
* `virtual_network_resource_group` - (Optional/Computed) The AKS virtual network resource group (string)

#### `node_pools`

##### Arguments

* `name` - (Required) The AKS node group name (string)
* `availability_zones` - (Optional) The AKS node pool availability zones (list)
* `count` - (Optional) The AKS node pool count. Default: `1` (int)
* `enable_auto_scaling` - (Optional) Is AKS node pool auto scaling enabled? Default: `false` (bool)
* `max_count` - (Optional) The AKS node pool max count. Required if `enable_auto_scaling=true` (int)
* `max_pods` - (Optional) The AKS node pool max pods. Default: `110` (int)
* `min_count` - (Optional) The AKS node pool min count. Required if `enable_auto_scaling=true` (int)
* `mode` - (Optional) The AKS node group mode. Default: `System` (string)
* `orchestrator_version` - (Optional) The AKS node pool orchestrator version (string)
* `os_disk_size_gb` - (Optional) The AKS node pool os disk size gb. Default: `128` (int)
* `os_disk_type` - (Optional) The AKS node pool os disk type. Default: `Managed` (string)
* `os_type` - (Optional) The AKS node pool os type. Default: `Linux` (string)
* `vm_size` - (Optional/computed) The AKS node pool orchestrator version (string)

### `eks_config`

#### Arguments

The following arguments are supported:

* `access_key` - (Required/Sensitive) The AWS Client ID to use (string)
* `kubernetes_version` - (Required) The Kubernetes master version (string)
* `secret_key` - (Required/Sensitive) The AWS Client Secret associated with the Client ID (string)
* `ami` - (Optional) AMI ID to use for the worker nodes instead of the default (string)
* `associate_worker_node_public_ip` - (Optional) Associate public ip EKS worker nodes. Default `true` (bool)
* `desired_nodes` - (Optional) The desired number of worker nodes. Just for Rancher v2.3.x and above. Default `3` (int)
* `instance_type` - (Optional) The type of machine to use for worker nodes. Default `t2.medium` (string)
* `key_pair_name` - (Optional) Allow user to specify key name to use. Just for Rancher v2.2.7 and above (string)
* `maximum_nodes` - (Optional) The maximum number of worker nodes. Default `3` (int)
* `minimum_nodes` - (Optional) The minimum number of worker nodes. Default `1` (int)
* `node_volume_size` - (Optional) The volume size for each node. Default `20` (int)
* `region` - (Optional) The AWS Region to create the EKS cluster in. Default `us-west-2` (string)
* `security_groups` - (Optional) List of security groups to use for the cluster. If it's not specified Rancher will create a new security group (list)
* `service_role` - (Optional) The service role to use to perform the cluster operations in AWS. If it's not specified Rancher will create a new service role (string)
* `session_token` - (Optional/Sensitive) A session token to use with the client key and secret if applicable (string)
* `subnets` - (Optional) List of subnets in the virtual network to use. If it's not specified Rancher will create 3 news subnets (list)
* `user_data` - (Optional/Computed) Pass user-data to the nodes to perform automated configuration tasks (string)
* `virtual_network` - (Optional) The name of the virtual network to use. If it's not specified Rancher will create a new VPC (string)

### `eks_config_v2`

#### Arguments

The following arguments are supported:

* `cloud_credential_id` - (Required) The EKS cloud_credential id (string)
* `imported` - (Optional) Set to `true` to import EKS cluster. Default: `false` (bool)
* `name` - (Optional/Computed) The EKS cluster name to import. Required to import a cluster (string)
* `kms_key` - (Optional) The AWS kms key to use (string)
* `kubernetes_version` - (Optional/Computed) The EKS cluster kubernetes version. Required to create a new cluster (string)
* `logging_types` - (Optional) The AWS cloudwatch logging types. `audit`, `api`, `scheduler`, `controllerManager` and `authenticator` values are allowed (list)
* `node_groups` - (Optional/Computed) The EKS cluster name to import. Required to create a new cluster (list)
* `private_access` - (Optional/Computed) The EKS cluster has private access (bool)
* `public_access` - (Optional/Computed) The EKS cluster has public access (bool)
* `public_access_sources` - (Optional/Computed) The EKS cluster public access sources (map)
* `region` - (Optional) The EKS cluster region. Default: `us-west-2` (string)
* `secrets_encryption` - (Optional/Computed) Enable EKS cluster secret encryption (bool)
* `security_groups` - (Optional/Computed) List of security groups to use for the cluster (list)
* `service_role` - (Optional) The AWS service role to use (string)
* `subnets` - (Optional) List of subnets in the virtual network to use (list)
* `tags` - (Optional) The EKS cluster tags (map)

#### `node_groups`

##### Arguments

The following arguments are supported:

* `name` - (Required) The EKS node group name to import (string)
* `desired_size` - (Optional) The EKS node group desired size. Default: `2` (int)
* `disk_size` - (Optional) The EKS node group disk size (Gb). Default: `20` (int)
* `ec2_ssh_key` - (Optional) The EKS node group ssh key (string)
* `gpu` - (Optional) Set true to EKS use gpu. Default: `false` (bool)
* `image_id` - (Optional) The EKS node group image ID (string)
* `instance_type` - (Optional) The EKS node group instance type. Default: `t3.medium` (string)
* `labels` - (Optional) The EKS cluster labels (map)
* `launch_template` - (Optional) The EKS node groups launch template (list Maxitem: 1)
* `max_size` - (Optional) The EKS node group maximum size. Default `2` (int)
* `min_size` - (Optional) The EKS node group maximum size. Default `2` (int)
* `node_role` - (Optional) The EKS node group node role ARN. Default `""` (string)
* `request_spot_instances` - (Optional) Enable EKS node group request spot instances (bool)
* `resource_tags` - (Optional) The EKS node group resource tags (map)
* `spot_instance_types` - (Optional) The EKS node group sport instace types (list string)
* `subnets` - (Optional) The EKS node group subnets (list string)
* `tags` - (Optional) The EKS cluster tags (map)
* `user_data` - (Optional) The EKS node group user data (string)
* `version` - (Computed) The EKS node group version (string)

##### `launch_template`

###### Arguments

* `id` - (Required) The EKS node group launch template ID (string)
* `name` - (Optional/Computed) The EKS node group launch template name (string)
* `version` - (Optional) The EKS node group launch template version. Default: `1` (int)

### `gke_config`

#### Arguments

The following arguments are supported:

* `cluster_ipv4_cidr` - (Required) The IP address range of the container pods (string)
* `credential` - (Required/Sensitive) The contents of the GC credential file (string)
* `disk_type` - (Required) Type of the disk attached to each node (string)
* `image_type` - (Required) The image to use for the worker nodes (string)
* `ip_policy_cluster_ipv4_cidr_block` - (Required) The IP address range for the cluster pod IPs (string)
* `ip_policy_cluster_secondary_range_name` - (Required) The name of the secondary range to be used for the cluster CIDR block (string)
* `ip_policy_node_ipv4_cidr_block` - (Required) The IP address range of the instance IPs in this cluster (string)
* `ip_policy_services_ipv4_cidr_block` - (Required) The IP address range of the services IPs in this cluster (string)
* `ip_policy_services_secondary_range_name` - (Required) The name of the secondary range to be used for the services CIDR block (string)
* `ip_policy_subnetwork_name` - (Required) A custom subnetwork name to be used if createSubnetwork is true (string)
* `locations` - (Required) Locations for GKE cluster (list)
* `machine_type` - (Required) Machine type for GKE cluster (string)
* `maintenance_window` - (Required) Maintenance window for GKE cluster (string)
* `master_ipv4_cidr_block` - (Required) The IP range in CIDR notation to use for the hosted master network (string)
* `master_version` - (Required) Master version for GKE cluster (string)
* `network` - (Required) Network for GKE cluster (string)
* `node_pool` - (Required) The ID of the cluster node pool (string)
* `node_version` - (Required) Node version for GKE cluster (string)
* `oauth_scopes` - (Required) The set of Google API scopes to be made available on all of the node VMs under the default service account (list)
* `project_id` - (Required) Project ID for GKE cluster (string)
* `service_account` - (Required) The Google Cloud Platform Service Account to be used by the node VMs (string)
* `sub_network` - (Required) Subnetwork for GKE cluster (string)
* `description` - (Optional) An optional description of this cluster (string)
* `disk_size_gb` - (Optional) Size of the disk attached to each node. Default `100` (int)
* `enable_alpha_feature` - (Optional) To enable Kubernetes alpha feature. Default `true` (bool)
* `enable_auto_repair` - (Optional) Specifies whether the node auto-repair is enabled for the node pool. Default `false` (bool)
* `enable_auto_upgrade` - (Optional) Specifies whether node auto-upgrade is enabled for the node pool. Default `false` (bool)
* `enable_horizontal_pod_autoscaling` - (Optional) Enable horizontal pod autoscaling for the cluster. Default `true` (bool)
* `enable_http_load_balancing` - (Optional) Enable HTTP load balancing on GKE cluster. Default `true` (bool)
* `enable_kubernetes_dashboard` - (Optional) Whether to enable the Kubernetes dashboard. Default `false` (bool)
* `enable_legacy_abac` - (Optional) Whether to enable legacy abac on the cluster. Default `false` (bool)
* `enable_master_authorized_network` - (Optional) Enable master authorized network. Set to `true` if `master_authorized_network_cidr_blocks` is set. Default `false` (bool)
* `enable_network_policy_config` - (Optional) Enable network policy config for the cluster. Default `true` (bool)
* `enable_nodepool_autoscaling` - (Optional) Enable nodepool autoscaling. Default `false` (bool)
* `enable_private_endpoint` - (Optional) Whether the master's internal IP address is used as the cluster endpoint. Default `false` (bool)
* `enable_private_nodes` - (Optional) Whether nodes have internal IP address only. Default `false` (bool)
* `enable_stackdriver_logging` - (Optional) Enable stackdriver monitoring. Default `true` (bool)
* `enable_stackdriver_monitoring` - (Optional) Enable stackdriver monitoring on GKE cluster (bool)
* `ip_policy_create_subnetwork` - (Optional) Whether a new subnetwork will be created automatically for the cluster. Default `false` (bool)
* `issue_client_certificate` - (Optional) Issue a client certificate. Default `false` (bool)
* `kubernetes_dashboard` - (Optional) Enable the Kubernetes dashboard. Default `false` (bool)
* `labels` - (Optional/Computed) The map of Kubernetes labels to be applied to each node (map)
* `local_ssd_count` - (Optional) The number of local SSD disks to be attached to the node. Default `0` (int)
* `master_authorized_network_cidr_blocks` - (Optional) Define up to 10 external networks that could access Kubernetes master through HTTPS (list)
* `max_node_count` - (Optional) Maximum number of nodes in the NodePool. Must be >= minNodeCount. There has to enough quota to scale up the cluster. Default `0` (int)
* `min_node_count` - (Optional) Minimmum number of nodes in the NodePool. Must be >= 1 and <= maxNodeCount. Default `0` (int)
* `node_count` - (Optional) Node count for GKE cluster. Default `3` (int)
* `preemptible` - (Optional) Whether the nodes are created as preemptible VM instances. Default `false` (bool)
* `region` - (Optional) GKE cluster region. Conflicts with `zone` (string)
* `resource_labels` - (Optional/Computed) The map of Kubernetes labels to be applied to each cluster (map)
* `use_ip_aliases` - (Optional) Whether alias IPs will be used for pod IPs in the cluster. Default `false` (bool)
* `taints` - (Required) List of Kubernetes taints to be applied to each node (list)
* `zone` - (Optional) GKE cluster zone. Conflicts with `region` (string)

### `gke_config_v2`

#### Arguments

* `name` - (Required/ForceNew) The GKE ip v4 cidr block (string)
* `google_credential_secret` - (Required/Sensitive) Google credential secret (string)
* `project_id` - (Required/ForceNew) The GKE cluster project id (string)
* `cluster_ipv4_cidr_block` - (Optional/Computed/ForceNew) The GKE ip v4 cidr block (string)
* `cluster_addons` - (Optional/Computed) The GKE cluster addons (List maxitems:1)
* `description` - (Optional/Computed/ForceNew) The GKE cluster addons (string)
* `enable_kubernetes_alpha` - (Optional/Computed/ForceNew) Enable Kubernetes alpha. Default: `false` (bool)
* `ip_allocation_policy` - (Optional/Computed/ForceNew) The GKE ip allocation policy (List maxitems:1)
* `imported` - (Optional/ForceNew) Is GKE cluster imported? Default: `false` (bool)
* `kubernetes_version` - (Optional/Computed) The kubernetes master version. Required for create new cluster (string)
* `labels` - (Optional/Computed) The GKE cluster labels (map)
* `locations` - (Optional/Computed) The GKE cluster locations (List)
* `logging_service` - (Optional/Computed) The GKE cluster logging service (string)
* `maintenance_window` - (Optional/Computed) The GKE cluster maintenance window (string)
* `master_authorized_networks_config` - (Optional/Computed/ForceNew) The GKE cluster master authorized networks config (List maxitems:1)
* `monitoring_service` - (Optional/Computed) The GKE cluster monitoring service (string)
* `network` - (Optional/Computed/ForceNew) The GKE cluster network. Required for create new cluster (string)
* `network_policy_enabled` - (Optional/Computed) Is GKE cluster network policy enabled? Default: `false` (bool)
* `node_pools` - (Optional/Computed) The GKE cluster node pools. Required for create new cluster (List)
* `private_cluster_config` - (Optional/Computed/ForceNew) The GKE private cluster config (List maxitems:1)
* `region` - (Optional/Computed/ForceNew) The GKE cluster region. Required if `zone` not set (string)
* `subnetwork` - (Optional/Computed/ForceNew) The GKE cluster subnetwork. Required for create new cluster (string)
* `zone` - (Optional/Computed/ForceNew) The GKE cluster zone. Required if `region` not set (string)

#### `cluster_addons`

##### Arguments

* `http_load_balancing` - (Optional/Computed) Enable GKE HTTP load balancing. Default: `false` (bool)
* `horizontal_pod_autoscaling` - (Optional/Computed) Enable GKE horizontal pod autoscaling. Default: `false` (bool)
* `network_policy_config` - (Optional/Computed) Enable GKE network policy config. Default: `false` (bool)

#### `ip_allocation_policy`

##### Arguments

* `cluster_ipv4_cidr_block` - (Optional/Computed) The GKE cluster ip v4 allocation cidr block (string)
* `cluster_secondary_range_name` - (Optional/Computed) The GKE cluster ip v4 allocation secondary range name(string)
* `create_subnetwork` - (Optional/Computed) Create GKE subnetwork? Default: `false` (bool)
* `node_ipv4_cidr_block` - (Optional/Computed) The GKE node ip v4 allocation cidr block (string)
* `services_ipv4_cidr_block` - (Optional/Computed) The GKE services ip v4 allocation cidr block (string)
* `services_secondary_range_name` - (Optional/Computed) The GKE services ip v4 allocation secondary range name (string)
* `subnetwork_name` - (Optional/Computed) The GKE cluster subnetwork name (string)
* `use_ip_aliases` - (Optional/Computed) Use GKE ip aliases? Default: `true` (bool)

#### `master_authorized_networks_config`

##### Arguments

* `cidr_blocks` - (Required) The GKE master authorized network config cidr blocks (List)
* `enabled` - (Optional) Enable GKE master authorized network config Default: `false` (bool)

##### `cidr_blocks`

###### Arguments

* `cidr_block` - (Required) The GKE master authorized network config cidr block (string)
* `display_name` - (Optional) The GKE master authorized network config cidr block dispaly name (string)

#### `node_pools`

##### Arguments

* `name` - (Required) The GKE node pool config name (string)
* `initial_node_count` - (Required) The GKE node pool config initial node count (int)
* `version` - (Required) The GKE node pool config version. Required for create new cluster (string)
* `autoscaling` - (Optional/computed) The GKE node pool config autoscaling (List maxitems:1)
* `config` - (Optional/Computed/ForceNew) The GKE node pool node config (List maxitems:1)
* `management` - (Optional/Computed) The GKE node pool config management (List maxitems:1)
* `max_pods_constraint` - (Optional/Computed) The GKE node pool config max pods constraint. Required for create new cluster if `ip_allocation_policy.use_ip_aliases = true` (int)

##### `autoscaling`

###### Arguments

* `enabled` - (Optional) Enable GKE node pool config autoscaling. Default: `false` (bool)
* `max_node_count` - (Optional/Computed) The GKE node pool config max node count (int)
* `min_node_count` - (Optional/Computed) The GKE node pool config min node count (int)

##### `config`

###### Arguments

* `disk_size_gb` - (Optional/Computed) The GKE node config disk size Gb (int)
* `disk_type` - (Optional/Computed) The GKE node config disk type (string)
* `image_type` - (Optional/Computed) The GKE node config image type (string)
* `labels` - (Optional/Computed) The GKE node config labels (map)
* `local_ssd_count` - (Optional/Computed) The GKE node config local ssd count (int)
* `machine_type` - (Optional/Computed) The GKE node config machine type (string)
* `oauth_scopes` - (Optional/Computed) The GKE node config oauth scopes (List)
* `preemptible` - (Optional) Enable GKE node config preemptible. Default: `false` (bool)
* `tags` - (Optional/Computed) The GKE node config tags (List)
* `taints` - (Optional) The GKE node config taints (List)

###### `taints`

####### Arguments

* `key` - (Required) The GKE taint key (string)
* `value` - (Required) The GKE taint value (string)
* `effect` - (Required) The GKE taint effect (string)

##### `management`

###### Arguments

* `auto_repair` - (Optional/Computed) Enable GKE node pool config management auto repair. Default: `false` (bool)
* `auto_upgrade` - (Optional/Computed) Enable GKE node pool config management auto upgrade. Default: `false` (bool)

#### `private_cluster_config`

##### Arguments

* `master_ipv4_cidr_block` - (Required) The GKE cluster private master ip v4 cidr block (string)
* `enable_private_endpoint` - (Optional) Enable GKE cluster private endpoint. Default: `false` (bool)
* `enable_private_nodes` - (Optional) Enable GKE cluster private endpoint. Default: `false` (bool)

### `oke_config`

#### Arguments

The following arguments are supported:

* `compartment_id` - (Required) The OCID of the compartment in which to create resources OKE cluster and related resources (string)
* `custom_boot_volume_size` - (Optional) Optional custom boot volume size (GB) for all nodes. If you specify 0, it will apply the default according to the `node_image` specified. Default `0` (int)
* `description` - (Optional) An optional description of this cluster (string)
* `enable_private_control_plane` - (Optional) Specifies whether Kubernetes API endpoint is a private IP only accessible from within the VCN. Default `false` Just for Rancher v2.5.10 or above (bool)
* `enable_kubernetes_dashboard` - (Optional) Specifies whether to enable the Kubernetes dashboard. Default `false` (bool)
* `enable_private_nodes` - (Optional) Specifies whether worker nodes will be deployed into a new, private, subnet. Default `false` (bool)
* `fingerprint` - (Required) The fingerprint corresponding to the specified user's private API Key (string)
* `flex_ocpus` - (Optional) Specifies number of OCPUs for nodes (requires flexible shape specified with `node_shape`) (int)
* `kms_key_id` - (Optional) The OCID of a KMS vault master key used to encrypt secrets at rest. See [here](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengencryptingdata.htm) for help creating a vault and master encryption key. Just for Rancher v2.5.9 or above (string)
* `kubernetes_version` - (Required) The Kubernetes version that will be used for your master *and* OKE worker nodes (string)
* `limit_node_count` - (Optional) The maximum number of worker nodes. Can limit `quantity_per_subnet`. Default `0` (no limit) (int)
* `load_balancer_subnet_name_1` - (Optional) The name of the first existing subnet to use for Kubernetes services / LB. `vcn_name` is also required when specifying an existing subnet. (string)
* `load_balancer_subnet_name_2` - (Optional) The name of a second existing subnet to use for Kubernetes services / LB. A second subnet is only required when it is AD-specific (non-regional) (string)
* `node_image` - (Required) The Oracle Linux OS image name to use for the OKE node(s). See [here](https://docs.cloud.oracle.com/en-us/iaas/images/) for a list of images. (string)
* `node_pool_dns_domain_name` - (Optional) Name for DNS domain of node pool subnet. Default `nodedns` (string)
* `node_pool_subnet_name` - (Optional) Name for node pool subnet. Default `nodedns` (string)
* `node_public_key_contents` - (Optional) The contents of the SSH public key file to use for the nodes (string)
* `node_shape` - (Required) The shape of the node (determines number of CPUs and  amount of memory on each OKE node) (string)
* `pod_cidr` - (Optional) A CIDR IP range from which to assign Kubernetes Pod IPs (string)
* `private_key_contents` - (Required/Sensitive) The private API key file contents for the specified user, in PEM format (string)
* `private_key_passphrase` - (Optional/Sensitive) The passphrase (if any) of the private key for the OKE cluster (string)
* `quantity_of_node_subnets` - (Optional) Number of node subnets. Default `1` (int)
* `quantity_per_subnet` - (Optional) Number of OKE worker nodes in each subnet / availability domain. Default `1` (int)
* `region` - (Required) The availability domain within the region to host the cluster. See [here](https://docs.cloud.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm) for a list of region names. (string)
* `service_cidr` - (Optional) A CIDR IP range from which to assign Kubernetes Service IPs (string)
* `service_dns_domain_name` - (Optional) Name for DNS domain of service subnet. Default `svcdns` (string)
* `skip_vcn_delete` - (Optional) Specifies whether to skip deleting the virtual cloud network (VCN) on destroy. Default `false` (bool)
* `tenancy_id` - (Required) The OCID of the tenancy in which to create resources (string)
* `user_ocid` - (Required) The OCID of a user who has access to the tenancy/compartment (string)
* `vcn_compartment_id` - (Optional) The OCID of the compartment (if different from `compartment_id`) in which to find the pre-existing virtual network set with `vcn_name`. (string)
* `vcn_name` - (Optional) The name of an existing virtual network to use for the cluster creation. If set, you must also set `load_balancer_subnet_name_1`. A VCN and subnets will be created if none are specified. (string)
* `worker_node_ingress_cidr` - (Optional) Additional CIDR from which to allow ingress to worker nodes (string)

### `cluster_auth_endpoint`

#### Arguments

* `ca_certs` - (Optional) CA certs for the authorized cluster endpoint (string)
* `enabled` - (Optional) Enable the authorized cluster endpoint. Default `true` (bool)
* `fqdn` - (Optional) FQDN for the authorized cluster endpoint (string)

### `cluster_monitoring_input`

#### Arguments

* `answers` - (Optional/Computed) Key/value answers for monitor input (map)
* `version` - (Optional) rancher-monitoring chart version (string)

### `cluster_template_answers`

#### Arguments

* `cluster_id` - (Optional) Cluster ID to apply answer (string)
* `project_id` - (Optional) Project ID to apply answer (string)
* `values` - (Optional) Key/values for answer (map)

### `cluster_template_questions`

#### Arguments

* `default` - (Required) Default variable value (string)
* `required` - (Optional) Required variable. Default `false` (bool)
* `type` - (Optional) Variable type. `boolean`, `int`, `password`, and `string` are allowed. Default `string` (string)
* `variable` - (Optional) Variable name (string)

### `cluster_registration_token`

#### Attributes

* `cluster_id` - (Computed) Cluster ID (string)
* `name` - (Computed) Name of cluster registration token (string)
* `command` - (Computed) Command to execute in a imported k8s cluster (string)
* `insecure_command` - (Computed) Insecure command to execute in a imported k8s cluster (string)
* `insecure_node_command` - (Computed) Insecure node command to execute in a imported k8s cluster (string)
* `insecure_windows_node_command` - (Computed) Insecure windows command to execute in a imported k8s cluster (string)
* `manifest_url` - (Computed) K8s manifest url to execute with `kubectl` to import an existing k8s cluster (string)
* `node_command` - (Computed) Node command to execute in linux nodes for custom k8s cluster (string)
* `token` - (Computed) Token for cluster registration token object (string)
* `windows_node_command` - (Computed) Node command to execute in windows nodes for custom k8s cluster (string)
* `annotations` - (Computed) Annotations for cluster registration token object (map)
* `labels` - (Computed) Labels for cluster registration token object (map)


## Timeouts

`rancher2_cluster` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters.
- `update` - (Default `30 minutes`) Used for cluster modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters.

## Import

Clusters can be imported using the Rancher Cluster ID

```
$ terraform import rancher2_cluster.foo &lt;CLUSTER_ID&gt;
```
