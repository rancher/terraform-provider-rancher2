---
page_title: "rancher2_cluster_v2 Resource"
---

# rancher2\_cluster\_v2 Resource

Provides a Rancher v2 Cluster v2 resource. This can be used to create RKE2 and K3S Clusters for Rancher v2 environments and retrieve their information. This resource is supported as tech preview from Rancher v2.6.0 and above.

## Example Usage


### Creating Rancher v2 custom cluster v2

```hcl
# Create a new rancher v2 RKE2 custom Cluster v2
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  fleet_namespace = "fleet-ns"
  kubernetes_version = "v1.21.4+rke2r2"
  enable_network_policy = false
  default_cluster_role_for_project_members = "user"
}

# Create a new rancher v2 K3S custom Cluster v2
resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  fleet_namespace = "fleet-ns"
  kubernetes_version = "v1.21.4+k3s1"
  enable_network_policy = false
  default_cluster_role_for_project_members = "user"
}
```

**Note** Once created, get the node command from `rancher2_cluster_v2.foo.cluster_registration_token`

### Creating Rancher v2 amazonec2 cluster v2

```hcl
# Create amazonec2 cloud credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  amazonec2_credential_config {
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"
  }
}

# Create amazonec2 machine config v2
resource "rancher2_machine_config_v2" "foo" {
  generate_name = "test-foo"
  amazonec2_config {
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = [<AWS_SG>]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}

# Create a new rancher v2 amazonec2 RKE2 Cluster v2
resource "rancher2_cluster_v2" "foo-rke2" {
  name = "foo-rke2"
  kubernetes_version = "v1.21.4+rke2r2"
  enable_network_policy = false
  default_cluster_role_for_project_members = "user"
  rke_config {
    machine_pools {
      name = "pool1"
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

# Create a new rancher v2 amazonec2 K3S Cluster v2
resource "rancher2_cluster_v2" "foo-k3s" {
  name = "foo-k3s"
  kubernetes_version = "v1.21.4+k3s1"
  enable_network_policy = false
  default_cluster_role_for_project_members = "user"
  rke_config {
    machine_pools {
      name = "pool1"
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

```hcl
# Create amazonec2 cloud credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  amazonec2_credential_config {
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"
  }
}

# Create amazonec2 machine config v2
resource "rancher2_machine_config_v2" "foo" {
  generate_name = "test-foo"
  amazonec2_config {
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = [<AWS_SG>]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}

resource "rancher2_cluster_v2" "foo" {
  name = "foo"
  kubernetes_version = "v1.21.4+k3s1"
  enable_network_policy = false
  rke_config {
    machine_pools {
      name = "pool1"
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
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required/ForceNew) The name of the Cluster v2 (string)
* `fleet_namespace` - (Optional/ForceNew) The fleet namespace of the Cluster v2. Default: `\"fleet-default\"` (string)
* `kubernetes_version` - (Required) The kubernetes version of the Cluster v2 (list maxitems:1)
* `agent_env_vars` - (Optional) Optional Agent Env Vars for Rancher agent (list)
* `rke_config` - (Optional/Computed) The RKE configuration for `k3s` and `rke2` Clusters v2. (list maxitems:1)
* `cloud_credential_secret_name` - (Optional) Cluster V2 cloud credential secret name (string)
* `default_pod_security_policy_template_name` - (Optional) Cluster V2 default pod security policy template name (string)
* `default_cluster_role_for_project_members` - (Optional) Cluster V2 default cluster role for project members (string)
* `enable_network_policy` - (Optional) Enable k8s network policy at Cluster V2 (bool)
* `annotations` - (Optional/Computed) Annotations for the Cluster V2 (map)
* `labels` - (Optional/Computed) Labels for the Cluster V2 (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cluster_registration_token` - (Computed/Sensitive) Cluster Registration Token generated for the cluster v2 (list maxitems:1)
* `kube_config` - (Computed/Sensitive) Kube Config generated for the cluster v2 (string)
* `cluster_v1_id` - (Computed) Cluster v1 id for cluster v2. (e.g to be used with `rancher2_sync`) (string)
* `resource_version` - (Computed) Cluster v2 k8s resource version (string)

## Nested blocks

### `agent_env_vars`

#### Arguments

* `name` - (Required) Rancher agent env var name (string)
* `value` - (Required) Rancher agent env var value (string)

### `rke_config`

#### Arguments

* `additional_manifest` - (Optional) Cluster V2 additional manifest (string)
* `local_auth_endpoint` - (Optional) Cluster V2 local auth endpoint (list maxitems:1)
* `upgrade_strategy` - (Optional) Cluster V2 upgrade strategy (list maxitems:1)
* `chart_values` - (Optional) Cluster V2 chart values. Must be in YAML format (string)
* `machine_global_config` - (Optional) Cluster V2 machine global config. Must be in YAML format (string)
* `machine_pools` - (Optional/Computed) Cluster V2 machine pools (list)
* `machine_selector_config` - (Optional/Computed) Cluster V2 machine selector config (list)
* `registries` - (Optional) Cluster V2 docker registries (list maxitems:1)
* `etcd` - (Optional) Cluster V2 etcd (list maxitems:1)

#### `local_auth_endpoint`

##### Arguments

* `ca_certs` - (Optional) CA certs for the authorized cluster endpoint (string)
* `enabled` - (Optional) Enable the authorized cluster endpoint. Default `false` (bool)
* `fqdn` - (Optional) FQDN for the authorized cluster endpoint (string)

#### `upgrade_strategy`

##### Arguments

* `control_plane_concurrency` - (Optional) How many controlplane nodes should be upgrade at time, 0 is infinite. Percentages are also accepted (string)
* `control_plane_drain_options` - (Optional) Controlplane nodes drain options (list maxitems:1)
* `worker_concurrency` - (Optional) How many worker nodes should be upgrade at time. Percentages are also accepted (string)
* `worker_drain_options` - (Optional) Worker nodes drain options (list maxitems:1)

##### `control_plane_drain_options` `worker_drain_options`

###### Arguments

* `enabled` - (Optional) Drain options enabled? Default `true` (bool)
* `force` - (Optional) Drain options force. Default `false` (bool)
* `ignore_daemon_sets` - (Optional) Drain options ignore daemon sets. Default `true` (bool)
* `ignore_errors` - (Optional) Drain options ignore errors. Default `false` (bool)
* `delete_empty_dir_data` - (Optional) Drain options delete empty dir data. Default `false` (bool)
* `disable_eviction` - (Optional) Drain options disable eviction. Default `false` (bool)
* `grace_period` - (Optional/Computed) Drain options grace period (int)
* `timeout` - (Optional/Computed) Drain options timeout (int)
* `skip_wait_for_delete_timeout_seconds` - (Optional/Computed) Drain options skip wait for delete timeout seconds (int)

#### `machine_pools`

##### Arguments

* `name` - (Required) Machine pool name (string)
* `cloud_credential_name` - (Required) Machine pool cloud credential secret name (string)
* `machine_config` - (Required) Machine pool node config (list)
* `control_plane_role` - (Optional) Machine pool control plane role? (bool)
* `etcd_role` - (Optional) Machine pool etcd role? (bool)
* `paused` - (Optional) Machine pool paused? (bool)
* `quantity` - (Optional) Machine pool quantity (int)
* `rolling_update` - (Optional) Machine pool rolling update (List maxitems:1)
* `taints` - (Optional) Machine pool taints (list)
* `worker_role` - (Optional) Machine pool worker role? (bool)

##### `machine_config`

###### Arguments

* `kind` - (Required) Machine config kind (string)
* `name` - (Required) Machine config name (string)

##### `rolling_update`

###### Arguments

* `max_unavailable` - (Optional) Rolling update max unavailable (string)
* `max_surge` - (Optional) Rolling update max surge (string)

##### `taints`

###### Arguments

* `key` - (Required) The taint key (string)
* `value` - (Required) The taint value (string)
* `effect` - (Optional) The taint effect. Default: `\"NoExecute\"` (string)

#### `machine_selector_config`

##### Arguments

* `machine_label_selector` - (Optional) Machine selector label (list maxitems:1)
* `config` - (Optional) Machine selector config (map)

##### `machine_label_selector`

###### Arguments

* `match_expressions` - (Optional) Machine selector label match expressions (list)
* `match_labels` - (Optional) Machine selector label match labels (map)

###### `match_expressions`

###### Arguments

* `key` - (Optional) Machine selector label match expressions key (string)
* `operator` - (Optional) Machine selector label match expressions operator (string)
* `values` - (Optional) Machine selector label match expressions values (List string)

#### `registries`

##### Arguments

* `configs` - (Optional) Cluster V2 docker registries config (list)
* `mirrors` - (Optional) Cluster V2 docker registries mirror (list)

##### `configs`

###### Arguments

* `hostname` - (Required) Registry hostname (string)
* `auth_config_secret_name` - (Optional) Registry auth config secret name (string)
* `tls_secret_name` - (Optional/Sensitive) Registry TLS secret name. TLS is a pair of Cert/Key (string)
* `ca_bundle` - (Optional) Registry CA bundle (string)
* `insecure` - (Optional) Registry insecure connectivity (bool)

##### `mirrors`

###### Arguments

* `hostname` - (Required) Registry hostname (string)
* `endpoints` - (Optional) Registry mirror endpoints (List)
* `rewrites` - (Optional) Registry mirror rewrites (map)

#### `etcd`

##### Arguments

* `disable_snapshots` - (Optional) Disable ETCD snapshots. Default: `false` (bool)
* `snapshot_schedule_cron` - (Optional) ETCD snapshot schedule cron (e.g `\"0 */5 * * *\"`) (string)
* `snapshot_retention` - (Optional) ETCD snapshot retention (int)
* `s3_config` - (Optional) Creation option for etcd service (list maxitems:1)

##### `s3_config`

###### Arguments

* `bucket` - (Required) Bucket name for S3 service (string)
* `endpoint` - (Required) ETCD snapshot S3 endpoint (string)
* `cloud_credential_name` - (Optional) ETCD snapshot S3 cloud credential name (string)
* `endpoint_ca` - (Optional) ETCD snapshot S3 endpoint CA (string)
* `folder` - (Optional) ETCD snapshot S3 folder (string)
* `region` - (Optional) ETCD snapshot S3 region (string)
* `skip_ssl_verify` - (Optional) Disable ETCD skip ssl verify. Default: `false` (bool)

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

`rancher2_cluster_v2` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters v2.
- `update` - (Default `30 minutes`) Used for cluster v2 modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters v2.

## Import

Clusters v2 can be imported using the Rancher Cluster v2 ID, that is in the form &lt;FLEET_NAMESPACE&gt;/&lt;CLUSTER_NAME&gt;

```
$ terraform import rancher2_cluster_v2.foo &lt;FLEET_NAMESPACE&gt;/&lt;CLUSTER_NAME&gt;
```
