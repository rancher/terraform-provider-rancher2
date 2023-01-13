---
page_title: "rancher2_cluster_template Resource"
---

# rancher2\_cluster\_template Resource

Provides a Rancher v2 Cluster Template resource. This can be used to create Cluster Templates for Rancher v2 RKE clusters and retrieve their information. 

Cluster Templates are available from Rancher v2.3.x and above.

## Example Usage

```hcl
# Create a new rancher2 Cluster Template
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
  description = "Terraform cluster template foo"
}
```

Creating Rancher v2 RKE cluster template with upgrade strategy. For Rancher v2.4.x or above.

```hcl
# Create a new rancher2 Cluster Template
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
        upgrade_strategy {
          drain = true
          max_unavailable_worker = "20%"
        }
      }
    }
    default = true
  }
  description = "Terraform cluster template foo"
}
```


## Argument Reference

* `name` - (Required) The cluster template name (string)
* `decription` - (Optional) The cluster template description (string)
* `members` - (Optional) Cluster template members (list)
* `template_revisions` - (Optional/Computed) Cluster template revisions (list)
* `annotations` - (Optional/Computed) Annotations for the cluster template (map)
* `labels` - (Optional/Computed) Labels for the cluster template (map)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `default_revision_id` - (Computed) Default cluster template revision ID (string)

## Nested blocks

### `members`

#### Arguments

* `access_type` - (Optional) Member access type. Valid values: `["read-only" | "owner"]` (string)
* `group_principal_id` - (Optional) Member group principal id (string)
* `user_principal_id` - (Optional) Member user principal id (string)

### `template_revisions`

#### Arguments

* `name` - (Required) The cluster template revision name (string)
* `cluster_config` - (Required) Cluster configuration (list maxitem: 1)
* `default` - (Optional) Default cluster template revision. Default `false` (bool)
* `enabled` - (Optional) Enable cluster template revision. Default `true` (bool)
* `questions` - (Optional) Cluster template questions (list)
* `annotations` - (Optional/Computed) Annotations for the cluster template revision (map)
* `labels` - (Optional/Computed) Labels for the cluster template revision (map)

#### Attributes

* `id` - (Computed) The cluster template revision ID (string)
* `cluster_template_id` - (Computed) Cluster template ID (string)

#### `cluster_config`

##### Arguments

* `cluster_auth_endpoint` - (Optional/Computed) Local cluster auth endpoint (list maxitems: 1)
* `default_cluster_role_for_project_members` - (Optional/Computed) Default cluster role for project members (string)
* `default_pod_security_policy_template_id` - (Optional/Computed) Default pod security policy template ID (string)
* `desired_agent_image` - (Optional/Computed) Desired agent image (string)
* `desired_auth_image` - (Optional/Computed) Desired auth image (string)
* `docker_root_dir` - (Optional/Computed) Desired auth image (string)
* `enable_cluster_alerting` - (Optional) Enable built-in cluster alerting. Default: `false` (bool)
* `enable_cluster_monitoring` - (Optional) Enable built-in cluster monitoring. Default: `false` (bool)
* `enable_network_policy` - (Optional) Enable project network isolation. Default: `false` (bool)
* `rke_config` - (Required) Rancher Kubernetes Engine Config (list maxitems: 1)
* `windows_prefered_cluster` - (Optional) Windows prefered cluster. Default: `false` (bool)

#### `questions`

##### Arguments

* `default` - (Required) Default variable value (string)
* `required` - (Optional) Required variable. Default `false` (bool)
* `type` - (Optional) Variable type. `boolean`, `int` and `string` are allowed. Default `string` (string)
* `variable` - (Optional) Variable name (string)

## Timeouts

`rancher2_cluster_template` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cluster templates.
- `update` - (Default `10 minutes`) Used for cluster template modifications.
- `delete` - (Default `10 minutes`) Used for deleting cluster templates.

## Import

Cluster Template can be imported using the rancher Cluster Template ID

```
$ terraform import rancher2_cluster_template.foo &lt;CLUSTER_TEMPLATE_ID&gt;
```
