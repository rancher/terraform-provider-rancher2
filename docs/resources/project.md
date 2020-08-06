---
page_title: "rancher2_project Resource"
---

# rancher2\_project Resource

Provides a Rancher v2 Project resource. This can be used to create projects for Rancher v2 environments and retrieve their information.

## Example Usage

```hcl
# Create a new rancher2 Project
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "<CLUSTER_ID>"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "2000m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
```

```hcl
# Create a new rancher2 Project enabling and customizing monitoring
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "<CLUSTER_ID>"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "2000m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
  enable_project_monitoring = true
  project_monitoring_input {
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
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project (string)
* `cluster_id` - (Required) The cluster id where create project (string)
* `container_resource_limit` - (Optional) Default containers resource limits on project (List maxitem:1)
* `description` - (Optional) A project description (string)
* `enable_project_monitoring` - (Optional) Enable built-in project monitoring. Default `false` (bool)
* `pod_security_policy_template_id` - (Optional) Default Pod Security Policy ID for the project (string)
* `project_monitoring_input` - (Optional) Project monitoring config. Any parameter defined in [rancher-monitoring charts](https://github.com/rancher/system-charts/tree/dev/charts/rancher-monitoring) could be configured (list maxitems:1)
* `resource_quota` - (Optional) Resource quota for project. Rancher v2.1.x or higher (list maxitems:1)
* `wait_for_cluster` - (Optional) Wait for cluster becomes active. Default `false` (bool)
* `annotations` - (Optional/Computed) Annotations for Node Pool object (map)
* `labels` - (Optional/Computed) Labels for Node Pool object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `container_resource_limit`

#### Arguments

* `limits_cpu` - (Optional) CPU limit for containers (string)
* `limits_memory` - (Optional) Memory limit for containers (string)
* `requests_cpu` - (Optional) CPU reservation for containers (string)
* `requests_memory` - (Optional) Memory reservation for containers (string)

### `project_monitoring_input`

#### Arguments

* `answers` - (Optional/Computed) Key/value answers for monitor input (map)
* `version` - (Optional) rancher-monitoring chart version (string)

### `resource_quota`

#### Arguments

* `project_limit` - (Required) Resource quota limit for project (list maxitems:1)
* `namespace_default_limit` - (Required) Default resource quota limit for  namespaces in project (list maxitems:1)

#### `project_limit` and `namespace_default_limit`

##### Arguments

The following arguments are supported:

* `config_maps` - (Optional) Limit for config maps in project (string)
* `limits_cpu` - (Optional) Limit for limits cpu in project (string)
* `limits_memory` - (Optional) Limit for limits memory in project (string)
* `persistent_volume_claims` - (Optional) Limit for persistent volume claims in project (string)
* `pods` - (Optional) Limit for pods in project (string)
* `replication_controllers` - (Optional) Limit for replication controllers in project (string)
* `requests_cpu` - (Optional) Limit for requests cpu in project (string)
* `requests_memory` - (Optional) Limit for requests memory in project (string)
* `requests_storage` - (Optional) Limit for requests storage in project (string)
* `secrets` - (Optional) Limit for secrets in project (string)
* `services_load_balancers` - (Optional) Limit for services load balancers in project (string)
* `services_node_ports` - (Optional) Limit for services node ports in project (string)

More info at [resource-quotas](https://rancher.com/docs/rancher/v2.x/en/k8s-in-rancher/projects-and-namespaces/resource-quotas/)

## Timeouts

`rancher2_project` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating projects.
- `update` - (Default `10 minutes`) Used for project modifications.
- `delete` - (Default `10 minutes`) Used for deleting projects.

## Import

Projects can be imported using the Rancher Project ID

```
$ terraform import rancher2_project.foo &lt;project_id&gt;
```

