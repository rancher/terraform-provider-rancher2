---
page_title: "rancher2_cluster_v2 Resource"
---

# rancher2\_cluster\_v2 Resource

Provides a Rancher v2 Cluster v2 resource. This can be used to create RKE2 and K3s Clusters for Rancher v2 environments and retrieve their information. 

This resource is available in Rancher v2.6.0 and above.

## Example Usage

You can create a Rancher v2 cluster v2 that runs either RKE2 or K3s.

There are some distribution-specific arguments, especially the ones under the `rke_config` section, that you can utilize to configure your RKE2 or K3s cluster.
More details and examples can be found on this page.

You can create two types of clusters depending on how nodes are managed:

- a custom cluster to which your existing VM(s) can be registered
- a node-driver cluster in which Rancher provisions and manages the VM(s) on the specified infrastructure provider

The cluster will be created as a custom cluster if there are no `machine_pools` in the configuration;
otherwise, it will be created as a node-driver cluster.

All arguments, except some distribution-specific ones, are applied to both custom and node-driver clusters of both distributions.

### Create a custom cluster

Below is the minimum configuration for creating a custom cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2-/k3s-version"
}
```

Once the cluster is created, you get the node registration command from `rancher2_cluster_v2.foo.cluster_registration_token`.

### Create a node-driver cluster

Before creating a node-driver cluster, you need to create a `rancher2_machine_config_v2` resource which will be referred to in the machine pool(s) of the cluster.

The example below demonstrates how to create a `rancher2_machine_config_v2` resource with `AmazonEC2` as the infrastructure provider:

```hcl
# Create AmazonEC2 cloud credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  amazonec2_credential_config {
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"
  }
}

# Create AmazonEC2 machine config v2
resource "rancher2_machine_config_v2" "foo" {
  generate_name = "test-foo"
  amazonec2_config {
    ami =  "ami-id"
    region = "region"
    security_group = ["security-group"]
    subnet_id = "subnet-id"
    vpc_id = "vpc-id"
    zone = "zone"
  }
}
```

For the full list of supported infrastructure providers and their arguments, please refer to the page for the `rancher2_machine_config_v2` resource.

Now, you can create an RKE2 or K3s node-driver cluster with one or more machine pools:

```hcl
# Create a cluster with multiple machine pools
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  enable_network_policy = false
  rke_config {
    # Nodes in this pool have control plane role and etcd roles
    machine_pools {
      name = "pool1"
      cloud_credential_secret_name = rancher2_cloud_credential.foo.id
      control_plane_role = true
      etcd_role = true
      worker_role = false
      quantity = 1
      drain_before_delete = true
      machine_config {
        kind = rancher2_machine_config_v2.foo.kind
        name = rancher2_machine_config_v2.foo.name
      }
    }
    # Each machine pool must be passed separately
    # Nodes in this pool have only the worker role
    machine_pools {
      name = "pool2"
      cloud_credential_secret_name = rancher2_cloud_credential.foo.id
      control_plane_role = false
      etcd_role = false
      worker_role = true
      quantity = 2
      drain_before_delete = true
      machine_config {
        kind = rancher2_machine_config_v2.foo.kind
        name = rancher2_machine_config_v2.foo.name
      }
    }
  }
}

# Create a cluster with a single machine pool
resource "rancher2_cluster_v2" "foo-k3s" {
  name = "foo-k3s"
  kubernetes_version = "rke2/k3s-version"
  enable_network_policy = false
  rke_config {
    # nodes in this pool have all roles 
    machine_pools {
      name = "pool"
      cloud_credential_secret_name = rancher2_cloud_credential.foo.id
      control_plane_role = true
      etcd_role = true
      worker_role = true
      quantity = 1
      machine_config {
        kind = rancher2_machine_config_v2.foo.kind
        name = rancher2_machine_config_v2.foo.name
      }
    }
  }
}
```

### Create a node-driver cluster with Harvester as the infrastructure provider

```hcl
# Get imported harvester cluster info
data "rancher2_cluster_v2" "foo-harvester" {
  name = "foo-harvester"
}

# Create a new cloud credential for an imported Harvester cluster
resource "rancher2_cloud_credential" "foo-harvester" {
  name = "foo-harvester"
  harvester_credential_config {
    cluster_id = data.rancher2_cluster_v2.foo-harvester.cluster_v1_id
    cluster_type = "imported"
    kubeconfig_content = data.rancher2_cluster_v2.foo-harvester.kube_config
  }
}

# Create a new rancher2 machine config v2 using harvester node_driver
resource "rancher2_machine_config_v2" "foo-harvester-v2" {
  generate_name = "foo-harvester-v2"
  harvester_config {
    vm_namespace = "default"
    cpu_count = "2"
    memory_size = "4"
    disk_info = <<EOF
    {
        "disks": [{
            "imageName": "harvester-public/image-57hzg",
            "size": 40,
            "bootOrder": 1
        }]
    }
    EOF
    network_info = <<EOF
    {
        "interfaces": [{
            "networkName": "harvester-public/vlan1"
        }]
    }
    EOF
    ssh_user = "ubuntu"
    user_data = <<EOF
    package_update: true
    packages:
      - qemu-guest-agent
      - iptables
    runcmd:
      - - systemctl
        - enable
        - '--now'
        - qemu-guest-agent.service
    EOF
  }
}

# Create a new cluster 
resource "rancher2_cluster_v2" "foo-harvester-v2" {
  name = "foo-harvester-v2"
  kubernetes_version = "<rke2/k3s-version>"
  rke_config {
    machine_pools {
      name = "pool1"
      cloud_credential_secret_name = rancher2_cloud_credential.foo-harvester.id
      control_plane_role = true
      etcd_role = true
      worker_role = true
      quantity = 1
      machine_config {
        kind = rancher2_machine_config_v2.foo-harvester-v2.kind
        name = rancher2_machine_config_v2.foo-harvester-v2.name
      }
    }
    machine_selector_config {
      config = {
        cloud-provider-name = ""
      }
    }
    machine_global_config = <<EOF
cni: "calico"
disable-kube-proxy: false
etcd-expose-metrics: false
EOF
    upgrade_strategy {
      control_plane_concurrency = "10%"
      worker_concurrency = "10%"
    }
    etcd {
      snapshot_schedule_cron = "0 */5 * * *"
      snapshot_retention = 5
    }
  }
}
```

### Create a node-driver cluster with Harvester as both the infrastructure provider and cloud provider 

The example below utilizes the arguments such as `machine_selector_config`, `machine_global_config`, and `chart_values`.
More explanations and examples for those arguments can be found on this page.


```hcl
# Get imported harvester cluster info
data "rancher2_cluster_v2" "foo-harvester" {
  name = "foo-harvester"
}
```

You need the kubeconfig file of the Harvester cluster to use it as the cloud provider for your cluster.

```bash
# Generate harvester cloud provider kubeconfig
RANCHER_SERVER_URL="<RANCHER_SERVER_URL>"
RANCHER_ACCESS_KEY="<RANCHER_ACCESS_KEY>"
RANCHER_SECRET_KEY="<RANCHER_SECRET_KEY>"
HARVESTER_CLUSTER_ID="<HARVESTER_CLUSTER_ID>"
CLUSTER_NAME="foo-harvester-v2-cloud-provider"
curl -k -X POST ${RANCHER_SERVER_URL}/k8s/clusters/${HARVESTER_CLUSTER_ID}/v1/harvester/kubeconfig \
   -H 'Content-Type: application/json' \
   -u ${RANCHER_ACCESS_KEY}:${RANCHER_SECRET_KEY} \
   -d '{"clusterRoleName": "harvesterhci.io:cloudprovider", "namespace": "default", "serviceAccountName": "'${CLUSTER_NAME}'"}' | xargs | sed 's/\\n/\n/g' > ${CLUSTER_NAME}-kubeconfig
```


```hcl
# Create a new Cloud Credential for an imported Harvester cluster
resource "rancher2_cloud_credential" "foo-harvester" {
  name = "foo-harvester"
  harvester_credential_config {
    cluster_id = data.rancher2_cluster_v2.foo-harvester.cluster_v1_id
    cluster_type = "imported"
    kubeconfig_content = data.rancher2_cluster_v2.foo-harvester.kube_config
  }
}

# Create a new rancher2 machine config v2 using harvester node_driver
resource "rancher2_machine_config_v2" "foo-harvester-v2-cloud-provider" {
  generate_name = "foo-harvester-v2-cloud-provider"
  harvester_config {
    vm_namespace = "default"
    cpu_count = "2"
    memory_size = "4"
    disk_info = <<EOF
    {
        "disks": [{
            "imageName": "harvester-public/image-57hzg",
            "size": 40,
            "bootOrder": 1
        }]
    }
    EOF
    network_info = <<EOF
    {
        "interfaces": [{
            "networkName": "harvester-public/vlan1"
        }]
    }
    EOF
    ssh_user = "ubuntu"
    user_data = <<EOF
    package_update: true
    packages:
      - qemu-guest-agent
      - iptables
    runcmd:
      - - systemctl
        - enable
        - '--now'
        - qemu-guest-agent.service
    EOF
  }
}

# Create a new cluster with harvester as the cloud provider
resource "rancher2_cluster_v2" "foo-harvester-v2-cloud-provider" {
  name = "foo-harvester-v2-cloud-provider"
  kubernetes_version = "<rke2/k3s-version>"
  rke_config {
    machine_pools {
      name = "pool1"
      cloud_credential_secret_name = rancher2_cloud_credential.foo-harvester.id
      control_plane_role = true
      etcd_role = true
      worker_role = true
      quantity = 1
      machine_config {
        kind = rancher2_machine_config_v2.foo-harvester-v2-cloud-provider.kind
        name = rancher2_machine_config_v2.foo-harvester-v2-cloud-provider.name
      }
    }
    # The kubeconfig file of the Harvester cluster is sent to all nodes via the machine_selector_config
    machine_selector_config {
      config = {
        cloud-provider-config = file("${path.module}/foo-harvester-v2-cloud-provider-kubeconfig")
        cloud-provider-name = "harvester"
      }
    }
    machine_global_config = <<EOF
cni: "calico"
disable-kube-proxy: false
etcd-expose-metrics: false
EOF
    upgrade_strategy {
      control_plane_concurrency = "10%"
      worker_concurrency = "10%"
    }
    etcd {
      snapshot_schedule_cron = "0 */5 * * *"
      snapshot_retention = 5
    }
    chart_values = <<EOF
harvester-cloud-provider:
  clusterName: foo-harvester-v2-cloud-provider
  cloudConfigPath: /var/lib/rancher/rke2/etc/config-files/cloud-provider-config
EOF
  }
}
```


### Customize the agent environment variables

The example below demonstrates how to set agent environment variables on a cluster.

```hcl
resource "rancher2_cluster_v2" "foo" {
  name               = "cluster-with-agent-env-vars"
  kubernetes_version = "rke2/k3s-version"
  agent_env_vars {
    name  = "foo1"
    value = "boo1"
  }
  # Each agent environment variable must be passed separately
  agent_env_vars {
    name  = "foo2"
    value = "boo2"
  }
  rke_config {
    # ...
  }
}
```


### Customize the cluster agent and the fleet agent

This argument is available in Rancher v2.7.5 and above.

You can configure the tolerations, affinity rules, and resource requirements for the `cattle-cluster-agent` and `fleet-agent` deployments.

The example below demonstrates how to set `cluster_agent_deployment_customization` and `fleet_agent_deployment_customization`:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  cluster_agent_deployment_customization {
    append_tolerations {
      key    = "tolerate/control-plane"
      effect = "NoSchedule"
      value  = "true"
    }
    # Each toleration must be passed separately
    append_tolerations {
      key    = "tolerate/etcd"
      effect = "NoSchedule"
      value  = "true"
    }
    # Override_affinity accepts a string whose value is an Affinity, a k8s native object, in JSON format
    override_affinity = <<EOF
{
  "nodeAffinity": {
    "requiredDuringSchedulingIgnoredDuringExecution": {
      "nodeSelectorTerms": [{
        "matchExpressions": [{
          "key": "not.this/nodepool",
          "operator": "In",
          "values": [
            "true"
          ]
        }]
      }]
    }
  }
}
EOF
    # The resource units follow Kubernetes' standard
    override_resource_requirements {
      cpu_limit      = "800m"
      cpu_request    = "500m"
      memory_limit   = "800Mi"
      memory_request = "500Mi"
    }
  }
  
  fleet_agent_deployment_customization {
    # The same format and requirement for values as the cluster_agent_deployment_customization argument
    # ...
  }
  rke_config {
    # In the case of a node-driver cluster
    machine_pools {
      # ...
      }
    } 
}
```

#### Customize scheduling for the cluster agent

This argument is available in Rancher 2.11.0 and above.

You can configure a Priority Class and or Pod Disruption Budget to be automatically deployed for the cattle cluster agent when provisioning or updating downstream clusters. 

In order to use this field, you must ensure that the `cluster-agent-scheduling-customization` feature is enabled in the Rancher server. 


The example below demonstrates how to set the `scheduling_customization` field to deploy a Priority Class and Pod Disruption Budget. Currently, this field is only supported for the cluster agent. 

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  cluster_agent_deployment_customization {
    scheduling_customization {
      priority_class {
        # The preemption_policy must be set to 'Never', 'PreemptLowerPriority', or omitted. 
        # If omitted, the default of 'PreemptLowerPriority' is used.
        preemption_policy = "PreemptLowerPriority"
        # The value cannot be less than negative 1 billion, or greater than 1 billion
        value = 1000000000
      }
      pod_disruption_budget {
        # min_available and max_unavailable must either be non-negative whole integers, 
        # or whole number percentages greater than 0 and less than or equal to 100 (e.g. "50%").
        # You cannot set both min_available and max_unavailable at the same time.
        min_available = "1"
        
        # max_unavailable = "1"
      }
    }
    
  }

  rke_config {
    # In the case of a node-driver cluster
    machine_pools {
      # ...
      }
    } 
}
```

### Enable Pod Security Policy Admission Configuration Template (PSACT) on the cluster

This argument is available in Rancher v2.7.2 and above.

**Note:** When PSACT is enabled in an RKE2 or K3s cluster, Rancher webhook sets the `kube-apiserver-arg` to the Pod Security Admission mount path. To suppress Terraform constantly trying to reconcile that arg, it must be set locally.

```hcl
locals {
  version = "rke2"  // "k3s" for K3s clusters
  rancher_psact_mount_path = "/etc/rancher/${local.version}/config/rancher-psact.yaml"
  kube_apiserver_arg = var.default_psa_template != null && var.default_psa_template != "" ? ["admission-control-config-file=${local.rancher_psact_mount_path}"] : []
}

resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  enable_network_policy = false
  // There are two builtin PSACT: rancher-privileged and rancher-restricted. You can also create new ones.
  default_pod_security_admission_configuration_template_name = "rancher-restricted"
  rke_config {
    machine_global_config = yamlencode({
      cni = "calico"
      disable-kube-proxy = false
      etcd-expose-metrics = false
      kube-apiserver-arg = local.kube_apiserver_arg
    })
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

### Create a cluster that uses a cluster-level authenticated `system-default-registry`

The `<auth-config-secret-name>` represents a generic Kubernetes secret that contains two keys with base64 encoded values: the `username` and `password` for the specified custom registry. If the `system-default-registry` is not authenticated, no secret is required and the section within the `rke_config` can be omitted if not otherwise needed.

Many registries may be specified in the `rke_config`s `registries` section, however, the `system-default-registry` from which core system images are pulled is always denoted via the `system-default-registry` key of the `machine_selector_config` or the `machine_global_config`. For more information on private registries, please refer to [the Rancher documentation](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/authentication-permissions-and-global-configuration/global-default-private-registry#setting-a-private-registry-with-credentials-when-deploying-a-cluster).

```hcl
resource "rancher2_cluster_v2" "foo_cluster_v2" {
  name = "cluster-with-custom-registry"
  kubernetes_version = "rke2/k3s-version"
  rke_config {
    machine_selector_config {
      config = {
        system-default-registry: "custom-registry-hostname"
      }
    }
    registries {
      configs {
        hostname = "custom-registry-hostname"
        auth_config_secret_name = "<auth-config-secret-name>"
        insecure = "<tls-insecure-bool>"
        tls_secret_name = ""
        ca_bundle = ""
      }
    }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

### Creating Rancher V2 Cluster with Machine Selector Files 

This argument is available in Rancher v2.7.2 and above.

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  enable_network_policy = false
  rke_config {
    machine_selector_files {
      machine_label_selector {
        match_labels = {
          # You can specify multiple labels 
          "rke.cattle.io/control-plane-role" = "true",
          "rke.cattle.io/etcd-role" = "true",
        }
        # You can also specify one or more match expressions
        match_expressions  {
          key = "name"
          values = ["a", "b"]
          operator = "In"
        }
        match_expressions  {
          key = "department"
          operator = "In"
          values = ["a", "b"]
        }
      }
      file_sources {
        # The same format for configmap
        secret {
          name = "config-file-v1"
          default_permissions = "644"
          items {
            key = "audit-policy"
            path ="/etc/rancher/rke2/custom/policy-v1.yaml"
            permissions = "666"
          }
        }
      }
    }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

### Create a cluster with machine global config or machine selector config

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  # An RKE2 version and REK2-specific server configuration are used in this example
  kubernetes_version = "rke2-version"
  enable_network_policy = false
  rke_config {
    machine_selector_config {
      machine_label_selector {
        match_labels = {
          # You can specify multiple labels 
          "rke.cattle.io/control-plane-role" = "true",
          "rke.cattle.io/etcd-role" = "true",
        }
        # You can also specify one or more match expressions
        match_expressions  {
          key = "name"
          values = ["a", "b"]
          operator = "In"
        }
        match_expressions  {
          key = "department"
          operator = "In"
          values = ["a", "b"]
        }
      }
      config = <<EOF
        kubelet-arg:
          - cloud-provider-name=external
      EOF
    }
    machine_selector_config {
      # The config will be applied to all nodes if there is no machine label selector
      config = <<EOF
        kube-proxy-arg:
          - log_file_max_size=1800
      EOF
    }
    machine_global_config = <<EOF
disable-kube-proxy: false
etcd-expose-metrics: false
kubelet-arg:
  - xxx=xxx
kube-proxy-arg:
  - xxx=xxx
kube-apiserver-arg:
  - xxx=xxx
kube-scheduler-arg:
  - xxx=xxx
kube-cloud-controller-manager-arg:
  - xxx=xxx
EOF
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

### Create a cluster with additional manifest

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  rke_config {
    additional_manifest = <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: testing-namespace-1
---
apiVersion: v1
kind: Namespace
metadata:
  name: testing-namespace-2
EOF
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

### Customize the ETCD snapshot feature on the cluster 

```hcl
resource "rancher2_cloud_credential" "credentials" {
  name = "rancher-creds"
  s3_credential_config {
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"
  }
}

resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  rke_config {
    etcd {
      snapshot_schedule_cron = "0 */12 * * *"
      snapshot_retention = 10
      s3_config {
        bucket = "backups"
        endpoint = "https://minio.host:9000"
        cloud_credential_name = rancher2_cloud_credential.credentials.id
      }
    }
  }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
}
```

### Customize distribution-specified server configurations in a cluster

You can customize all server configurations on the cluster by utilizing the `machine_global_config` argument. 

For the full list of server configurations, please refer to [RKE2 server configuration](https://docs.rke2.io/reference/server_config) and [K3s server configuration](https://docs.k3s.io/cli/server).

The example below demonstrates how to disable the system services in a K3s cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "k3s-version"
  rke_config {
    machine_global_config= <<EOF
disable:
  - coredns
  - servicelb
  - traefik
  - local-storage
  - metrics-server
EOF
  }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
}
```

The example below demonstrates how to disable the system services in an RKE2 cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2-version"
  rke_config {
    machine_global_config= <<EOF
disable:
  - rke2-coredns
  - rke2-ingress-nginx
  - rke2-metrics-server
  - metrics-server
EOF
  }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
}
```

The example below demonstrates how to add additional hostnames or IPv4/IPv6 addresses as Subject Alternative Names on the server TLS cert in an RKE2/K3s cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  rke_config {
    machine_global_config= <<EOF
tls-san: ["example-website.com", "100.100.100.100", "2002:db8:3333:4444:5555:6666:7777:8888"]
EOF
  }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
}
```

The example below demonstrates how to configure the IPv4/IPv6 network CIDRs to use for pod IPs and service IPs in an RKE2/K3s cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2/k3s-version"
  rke_config {
    machine_global_config= <<EOF
cluster-cidr: "0.42.0.0/16"
service-cidr: "0.42.0.0/16"
EOF
  }
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
}
```

### Customize chart values in a cluster

You can specify the values for the system charts installed by RKE2 or K3s. 

For more information about how RKE2 or K3s manage packaged components, please refer to [RKE2 documentation](https://docs.rke2.io/helm) or [K3s documentation](https://docs.k3s.io/installation/packaged-components).

The example below demonstrates how to customize chart values in an RKE2 cluster:

```hcl
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "rke2-version"
  enable_network_policy = false
  rke_config {
    chart_values = <<EOF
rke2-calico:
  calicoctl:
    image: rancher/mirrored-calico-ctl
    tag: v3.19.2
  certs:
    node:
      cert: null
      commonName: null
      key: null
    typha:
      caBundle: null
      cert: null
      commonName: null
      key: null
  felixConfiguration:
    featureDetectOverride: ChecksumOffloadBroken=true
  global:
    systemDefaultRegistry: ""
  imagePullSecrets: {}
  installation:
    calicoNetwork:
      bgp: Disabled
      ipPools:
      - blockSize: 24
        cidr: 10.42.0.0/16
        encapsulation: VXLAN
        natOutgoing: Enabled
    controlPlaneTolerations:
    - effect: NoSchedule
      key: node-role.kubernetes.io/control-plane
      operator: Exists
    - effect: NoExecute
      key: node-role.kubernetes.io/etcd
      operator: Exists
    enabled: true
    imagePath: rancher
    imagePrefix: mirrored-calico-
    kubernetesProvider: ""
  ipamConfig:
    autoAllocateBlocks: true
    strictAffinity: true
  tigeraOperator:
    image: rancher/mirrored-calico-operator
    registry: docker.io
    version: v1.17.6
EOF
    # In the case of a node-driver cluster
    machine_pools {
      # ...
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, forceNew, string) The name of the cluster.
* `fleet_namespace` - (Optional, ForceNew, string, default: fleet-default) Fleet namespace is the namespace where the cluster is to create in the local cluster. It is recommended to leave it as the default value. 
* `kubernetes_version` - (Required, string) The RKE2 or K3s version for the cluster.
* `agent_env_vars` - (Optional, list) Agent env vars is a list of additional environment variables to be appended to the `cattle-cluster-agent` and `fleet-agent` deployment, and the plan for the [system upgrade controller](https://github.com/rancher/system-upgrade-controller) to upgrade nodes.
* `cluster_agent_deployment_customization` - (Optional, list) Cluster agent deployment customization specifies the additional tolerations, new affinity rules, and new resource requirements on the `cattle-cluster-agent` deployment. This argument is available in Rancher v2.7.5 and above.
* `fleet_agent_deployment_customization` - (Optional, list) Fleet agent deployment customization specifies the additional tolerations, new affinity rules, and new resource requirements on the `fleet-agent` deployment. The argument is available in Rancher v2.7.5 and above.
* `rke_config` - (Optional/computed, list, max length: 1) The RKE configuration for the cluster.
* `local_auth_endpoint` - (Optional, list, max length: 1) Local auth endpoint configures the Authorized Cluster Endpoint (ACE) which can be used to directly access the Kubernetes API server, without requiring communication through Rancher. For more information, please refer to [Rancher Documentation](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/kubernetes-clusters-in-rancher-setup/register-existing-clusters#authorized-cluster-endpoint-support-for-rke2-and-k3s-clusters).
* `cloud_credential_secret_name` - (Optional, string) Cloud credential secret name is the secret to be used when a cloud credential secret name is not specified at the machine pool level. 
* `default_pod_security_admission_configuration_template_name` - (Optional, string) The name of the pre-defined pod security admission configuration template to be applied to the cluster. Rancher admins (or those with the right permissions) can create, manage, and edit those templates. For more information, please refer to [Rancher Documentation](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/authentication-permissions-and-global-configuration/psa-config-templates). The argument is available in Rancher v2.7.2 and above.
* `default_cluster_role_for_project_members` - (Optional, string) Default cluster role for project members.
* `enable_network_policy` - (Optional, bool, default: false) Enable k8s network policy on the cluster.
* `annotations` - (Optional/computed, map) Annotations for the Cluster.
* `labels` - (Optional/computed, map) Labels for the Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed, string) The ID of the resource.
* `cluster_registration_token` - (Computed, sensitive, list, max length: 1) Cluster Registration Token generated for the cluster.
* `kube_config` - (Computed/Sensitive) Kube Config generated for the cluster. Note: When the cluster has `local_auth_endpoint` enabled, the kube_config will not be available until the cluster is `connected`.
* `cluster_v1_id` - (Computed, string) Cluster v1 id for cluster v2. (e.g. to be used with `rancher2_sync`).
* `resource_version` - (Computed, string) Cluster's k8s resource version.

**Note:** For Rancher 2.6.0 and above: if setting `kubeconfig-generate-token=false` then the generated `kube_config` will not contain any user token. `kubectl` will generate the user token executing the [rancher cli](https://github.com/rancher/cli/releases/tag/v2.6.0), so it should be installed previously.

## Nested blocks

### `agent_env_vars`

#### Arguments

* `name` - (Required, string) Rancher agent env var name.
* `value` - (Required, string) Rancher agent env var value.

### `cluster_agent_deployment_customization` and `fleet_agent_deployment_customization`

These arguments are available in Rancher v2.7.5 and above. The `scheduling_customization` argument is only available in Rancher 2.11 and above, may only be set within `cluster_agent_deployment_customization`, and requires that the `cattle-cluster-agent-scheduling-customization` feature be enabled.

#### Arguments

* `append_tolerations` - (Optional, list) A list of tolerations to be appended to the default tolerations.
* `override_affinity` - (Optional, string, JSON format) Override affinity overrides the global default affinity setting.
* `override_resource_requirements` - (Optional, list) Override resource requirements overrides the default value for requests and/or limits. 
+ `scheduling_customization` - (Optional, list) Supported in Rancher 2.11.0 and above. Defines the configuration of a Priority Class and or Pod Disruption Budget. Currently only supported in the `cluster_agent_deployment_customization` field, and requires the `cattle_cluster_agent_scheduling_customization` feature to be enabled.

### `append_tolerations`

#### Arguments

* `key` - (Required, string) The toleration key.
* `effect` - (Optional, string) The toleration effect. Default: `\"NoSchedule\"`.
* `operator` - (Optional, string) The toleration operator.
* `seconds` - (Optional, int) The number of seconds a pod will stay bound to a node with a matching taint.
* `value` - (Optional, string) The toleration value.

### `override_resource_requirements`

#### Arguments

The resource units follow Kubernetes' standard, 
see more information on [Resource Management for Pods and Containers](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#resource-units-in-kubernetes).

* `cpu_limit` - (Optional, string) The maximum CPU limit for agent. 
* `cpu_request` - (Optional, string) The minimum CPU required for agent.
* `memory_limit` - (Optional, string) The maximum memory limit for agent.
* `memory_request` - (Optional, string) The minimum memory required for agent.

### `scheduling_customization`

#### Arguments

* `pod_disruption_budget` - (Optional, list) The definition of a Pod Disruption Budget deployed for the cluster agent
* `priority_class` - (Optional, list) The definition of a Priority Class deployed for the cluster agent

### `pod_disruption_budget`

#### Arguments

* `min_available` - (Optional, string) The minimum number of agent replicas that must be running at a given time. This can be a non-negative whole number or a whole number percentage (e.g. "1", "50%"). This field cannot be used at the same time as `max_unavailable`.
* `max_unavailable` - (Optional, string) The maximum number of agent replicas that can be unavailable at a given time. This can be a non-negative whole number or a whole number percentage (e.g. "1", "50%"). This field cannot be used at the same time as `min_available`. 

### `priority_class`

#### Arguments

* `value` - (Optional, int) The priority value set for the Priority Class. Must be greater than or equal to negative 1 billion, and less than or equal to 1 billion. 
* `preemption_policy` (Optional, string) The preemption policy set for the Priority Class. Must be set to either 'Never', 'PreemptLowerPriority', or omitted.

### `rke_config`

#### Arguments

* `additional_manifest` - (Optional, string, must be in YAML format) The value of the additional manifest is delivered to the path `/var/lib/rancher/rke2/server/manifests/rancher/addons.yaml` or `/var/lib/rancher/k3s/server/manifests/rancher/addons.yaml` on the control plane nodes.
* `local_auth_endpoint` - (Deprecated) Use rancher2_cluster_v2.local_auth_endpoint instead.
* `upgrade_strategy` - (Optional, list, max length: 1) Cluster upgrade strategy.
* `chart_values` - (Optional, string, must be in YAML format) The value for the system charts installed by the distribution. For more information about how RKE2 or K3s manage packaged components, please refer to [RKE2 documentation](https://docs.rke2.io/helm) or [K3s documentation](https://docs.k3s.io/installation/packaged-components).
* `machine_global_config` - (Optional, string, must be in YAML format) Machine global config specifies the distribution-specified server configuration applied to all nodes. For the full list of server configurations, please refer to [RKE2 server configuration](https://docs.rke2.io/reference/server_config) or [K3s server configuration](https://docs.k3s.io/cli/server).
* `machine_pools` - (Optional/computed, list) Cluster V2 machine pools.
* `machine_selector_config` - (Optional/computed, list) Machine selector config is the same as machine_global_config except that a label selector can be specified with the configuration. The configuration will only be applied to nodes that match the provided label selector. The configuration from machine_selector_config takes precedence over the one from machine_global_config. This argument is available in Rancher v2.7.2 and later.
* `machine_selector_files` - (Optional/computed, list) Machine selector files provide a means to deliver files to nodes so that the files can be in place before initiating RKE2/K3s server or agent processes. Please refer to Rancher documentation for [RKE2 Cluster Configuration Reference](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/rke2-cluster-configuration#machineselectorfiles) and [K3s Cluster Configuration Reference](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/k3s-cluster-configuration#machineselectorfiles). This argument is available in Rancher v2.7.2 and later.
* `registries` - (Optional, list, max length: 1) Docker registries from which the cluster pulls images. 
* `etcd` - (Optional/computed, list, max length: 1) Etcd configures the behavior of the automatic etcd snapshot feature.
* `rotate_certificates` (Optional, list, max length: 1) Cluster V2 certificate rotation.
* `etcd_snapshot_create` (Optional, list, max length: 1) Cluster V2 etcd snapshot create.
* `etcd_snapshot_restore` (Optional, list, max length: 1) Cluster V2 etcd snapshot restore.

#### `local_auth_endpoint`

##### Arguments

* `enabled` - (Optional, bool, default: false) Enable the authorized cluster endpoint.
* `fqdn` - (Optional, string) FQDN for the authorized cluster endpoint. If one is entered, it should point to the downstream cluster.
* `ca_certs` - (Optional, string) CA certs for the authorized cluster endpoint. It is only needed if there is a load balancer in front of the downstream cluster that is using an untrusted certificate. If you have a valid certificate, then nothing needs to be added to the CA Certificates field.

#### `upgrade_strategy`

##### Arguments

* `control_plane_concurrency` - (Optional, string) How many control plane nodes should be upgraded at a time, 0 is infinite. Percentages are also accepted.
* `control_plane_drain_options` - (Optional, list, max length: 1) Controlplane nodes drain options.
* `worker_concurrency` - (Optional, string) How many worker nodes should be upgraded at a time. Percentages are also accepted.
* `worker_drain_options` - (Optional, list, max length: 1) Worker nodes drain options.

##### `control_plane_drain_options`, `worker_drain_options`

###### Arguments

* `enabled` - (Optional, bool, default: true) If `enabled` is set to true, nodes will be drained before upgrade.
* `force` - (Optional, bool, default: false) If `force` is set to true, drain nodes even if there are standalone pods that are not managed by a ReplicationController, Job, or DaemonSet. Drain will not proceed without `force` set to true if there are such pods.
* `ignore_daemon_sets` - (Optional, bool, default: true) If `ignore_daemon_sets` is set to false, drain will not proceed if there are DaemonSet-managed pods.
* `ignore_errors` - (Optional, bool, default: false) If `ignore_errors` is set to true,  errors that occurred between drain nodes in group are ignored.
* `delete_empty_dir_data` - (Optional, bool, default: false) if `delete_empty_dir_data` is set to true, continue draining even if there are pods using emptyDir (local storage).
* `disable_eviction` - (Optional, bool, default: false) If `disable_eviction` is set to true, force drain to use delete rather than evict.
* `grace_period` - (Optional/computed, int) Time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used.
* `timeout` - (Optional/computed, int) Time to wait (in seconds) before giving up for one try.
* `skip_wait_for_delete_timeout_seconds` - (Optional/computed, int) Skip waiting for the pods that have a DeletionTimeStamp > N seconds to be deleted. Seconds must be greater than 0 to skip. Such pods will be force deleted. 

#### `machine_pools`

##### Arguments

* `name` - (Required, string) Machine pool name.
* `cloud_credential_secret_name` - (Optional, string) Machine pool cloud credential secret name.
* `machine_config` - (Required, list) Machine pool node config.
* `control_plane_role` - (Optional, bool) Machine pool control plane role?
* `etcd_role` - (Optional, bool) Machine pool etcd role?
* `drain_before_delete` - (Optional, bool) Machine Pool Drain Before Delete?
* `node_drain_timeout` - (Optional, int) Seconds a machine has to drain before deletion.
* `paused` - (Optional, bool) Machine pool paused?
* `quantity` - (Optional, int) Machine pool quantity.
* `rolling_update` - (Optional, list, max length: 1) Machine pool rolling update.
* `taints` - (Optional, list) Machine pool taints.
* `worker_role` - (Optional, bool) Machine pool worker role?
* `node_startup_timeout_seconds` - (Optional, int) Seconds a new node has to become active before it is replaced.
* `unhealthy_node_timeout_seconds` - (Optional, int) Seconds an unhealthy node has to become active before it is replaced.
* `max_unhealthy` - (Optional, string) Max unhealthy nodes for automated replacement to be allowed.
* `unhealthy_range` - (Optional, string) Range of unhealthy nodes for automated replacement to be allowed.
* `machine_labels` - (Optional, map) Labels for Machine pool nodes.
* `machine_os` - (Optional) OS Type in machine pool. Default `linux`(string)
* `labels` - (Optional, map) Labels for Machine Deployment Resource.
* `annotations` - (Optional, map) Annotations for Machine Deployment Resource. 

##### `machine_config`

###### Arguments

* `kind` - (Required, string) Machine config kind.
* `name` - (Required, string) Machine config name.
* `api_version` - (Optional, string) Api version of the machine_config.

##### `rolling_update`

###### Arguments

* `max_unavailable` - (Optional, string) Rolling update max unavailable.
* `max_surge` - (Optional, string) Rolling update max surge.

##### `taints`

###### Arguments

* `key` - (Required, string) The taint key.
* `value` - (Required, string) The taint value.
* `effect` - (Optional, string) The taint effect. Default: `\"NoExecute\"`.

#### `machine_selector_config`

This argument is available in Rancher v2.7.2 and later.

##### Arguments

* `machine_label_selector` - (Optional, list, max length: 1) Machine selector label is a label query over a set of resources. The result of match_labels and match_expressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
* `config` - (Optional, string, must be in YAML format) Config is the distribution-specify configuration to be applied to nodes that match the provided label selector. For more information, please refer to Rancher's documentation for [RKE2 Cluster Configuration](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/rke2-cluster-configuration#machineselectorconfig) or [K3s Cluster Configuration](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/k3s-cluster-configuration#machineselectorconfig)

##### `machine_label_selector`

###### Arguments

* `match_expressions` - (Optional, list) Match expressions is a list of label selector requirements. The requirements are ANDed.
* `match_labels` - (Optional, map) Machine selector label is a map of {key,value} pairs, the requirements are ANDed.

###### `match_expressions`

###### Arguments

* `key` - (Optional, string) Key is the label key that the selector applies to.
* `operator` - (Optional, string)  Operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
* `values` - (Optional, string list) Values is a list of string values.

#### `machine_selector_files`

This argument is available in Rancher v2.7.2 and later.

##### Arguments

* `machine_label_selector` - (Optional, list, max length: 1) Machine selector label is a label query over a set of resources. The result of match_labels and match_expressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
* `file_sources` - (Optional, list) File sources represents the source of the files. Multiple files can be delivered to nodes that match the provided label selector.

#### `file_sources`

##### Arguments

* `secret` - (Optional, list, max length: 1) Secret represents a K8s secret which is the source of files. It is mutually exclusive with configmap.
* `configmap` - (Optional, list, max length: 1) Configmap represents a K8s configmap which is the source of files. It is mutually exclusive with secret.

#### `secret`

##### Arguments

* `name` - (Required, string) The name of the secret. 
* `default_permissions` - (Optional, string) The numeric representation of the default file permissions for all files defined under the items. 
* `items` - (Optional, list) Items is a list of configurations for files, such as where to retrieve the content from the source, where to put the file on nodes, and etc. 

#### `configmap`

##### Arguments

* `name` - (Required, string) The name of the configmap.
* `default_permissions` - (Optional, string) The numeric representation of the default file permissions for all files defined under the items.
* `items` - (Optional, list) Items is a list of configurations for files, such as where to retrieve the content from the source, where to put the file on nodes, etc.

#### `items`

##### Arguments

* `key` - (Required, string) Key is the name of the key of the item to retrieve.
* `path` - (Required, string) Path is the absolute path to put the file in the target node.
* `dynamic` - (Optional, boolean, default: true) If true, the file is ignored when determining whether the node should be drained before updating the node plan.
* `permissions` - (Optional, string) Permissions is the numeric representation of the file permission. It takes precedence over the default permissions at the outer level.
* `hash` - (Optional, string) Hash is the base64 encoded value of the SHA256 checksum of the file's content. If specified, it is used to validate the integrity of the file content.


#### `registries`

##### Arguments

* `configs` - (Optional, list) Cluster V2 docker registries config.
* `mirrors` - (Optional, list) Cluster V2 docker registries mirror.

##### `configs`

###### Arguments

* `hostname` - (Required, string) Registry hostname.
* `auth_config_secret_name` - (Optional, string) Name of the secret that contains two keys with base64 encoded values: the username and password for the specified custom registry. No secret is required if the system-default-registry is not authenticated.
* `tls_secret_name` - (Optional, sensitive, string) Registry TLS secret name. TLS is a pair of Cert/Key.
* `ca_bundle` - (Optional, string) Registry CA bundle.
* `insecure` - (Optional, bool) Registry insecure connectivity.

##### `mirrors`

###### Arguments

* `hostname` - (Required, string) Registry hostname.
* `endpoints` - (Optional, list) Registry mirror endpoints.
* `rewrites` - (Optional, map) Registry mirror rewrites.

#### `etcd`

##### Arguments

* `disable_snapshots` - (Optional, bool, default: false) Disable ETCD snapshots.
* `snapshot_schedule_cron` - (Optional, string) ETCD snapshot schedule cron (e.g `\"0 */5 * * *\"`).
* `snapshot_retention` - (Optional, int) ETCD snapshot retention.
* `s3_config` - (Optional, list, max length: 1) Creation option for etcd service.

##### `s3_config`

###### Arguments

* `bucket` - (Required, string) Bucket name for S3 service.
* `endpoint` - (Required, string) ETCD snapshot S3 endpoint.
* `cloud_credential_name` - (Optional, string) ETCD snapshot S3 cloud credential name.
* `endpoint_ca` - (Optional, string) ETCD snapshot S3 endpoint CA.
* `folder` - (Optional, string) ETCD snapshot S3 folder.
* `region` - (Optional, string) ETCD snapshot S3 region.
* `skip_ssl_verify` - (Optional, bool, default: false) Disable ETCD skip ssl verify.

##### `rotate_certificates`

###### Arguments

* `generation` - (Required, int) Desired certificate rotation generation.
* `services` - (Optional, list of string) Service certificates to rotate with this generation.

##### `etcd_snapshot_create`

###### Arguments

* `generation` - (Required, int) ETCD generation to initiate a snapshot.

##### `etcd_snapshot_restore`

###### Arguments

* `name` - (Required, string) ETCD snapshot name to restore.
* `generation` (Required, int) ETCD snapshot desired generation.
* `restore_rke_config` (Optional, string) ETCD restore RKE config (set to none, all, or kubernetesVersion).

### `cluster_registration_token`

#### Attributes

* `cluster_id` - (Computed, string) Cluster ID.
* `name` - (Computed, string) Name of cluster registration token.
* `command` - (Computed, string) Command to execute in an imported k8s cluster.
* `insecure_command` - (Computed, string) Insecure command to execute in an imported k8s cluster.
* `insecure_node_command` - (Computed, string) Insecure node command to execute in an imported k8s cluster.
* `insecure_windows_node_command` - (Computed, string) Insecure windows command to execute in an imported k8s cluster.
* `manifest_url` - (Computed, string) K8s manifest url to execute with `kubectl` to import an existing k8s cluster.
* `node_command` - (Computed, string) Node command to execute in Linux nodes for custom k8s cluster.
* `token` - (Computed, string) Token for cluster registration token object.
* `windows_node_command` - (Computed, string) Node command to execute in Windows nodes for custom k8s cluster.
* `annotations` - (Computed, map) Annotations for cluster registration token object.
* `labels` - (Computed, map) Labels for cluster registration token object.

## Timeouts

`rancher2_cluster_v2` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters v2.
- `update` - (Default `30 minutes`) Used for cluster v2 modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters v2.

## Import

Clusters v2 can be imported using the Rancher Cluster v2 ID, that is in the form &lt;FLEET_NAMESPACE&gt;/&lt;CLUSTER_NAME&gt;

```
$ terraform import rancher2_cluster_v2.foo <FLEET_NAMESPACE>/<CLUSTER_NAME>
```
