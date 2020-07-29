---
page_title: "rancher2_multi_cluster_app Resource"
---

# rancher2\_multi_cluster_app Resource

Provides a Rancher v2 multi_cluster_app resource. This can be used to deploy multi_cluster_app on Rancher v2.

This resource can also modify Rancher v2 multi cluster apps in 3 ways:
- `Add/Remove targets`: If `targets` arguments is modified, the multi cluster app targets will be updated.
- `Rollback`: If `revision_id` argument is provided or modified the app will be rolled back accordingly. A new `revision_id` will be generated in Rancher. It will also generate a non-empty terraform plan that will require manual .tf file intervention. Use carefully.
- `Update`: If any other argument is modified the app will be upgraded.

Note: In case of multiple resource modification in a row, `rollback` has preference.

## Example Usage

```hcl
# Create a new rancher2 Multi Cluster App
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "<catalog_name>"
  name = "foo"
  targets {
    project_id = "<project_id>"
  }
  template_name = "<template_name>"
  template_version = "<template_version>"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  roles = ["project-member"]
}
```

```hcl
# Create a new rancher2 Multi Cluster App overriding answers
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "<catalog_name>"
  name = "foo"
  targets {
    project_id = "<project_id1>"
  }
  targets {
    project_id = "<project_id2>"
  }
  template_name = "<template_name>"
  template_version = "<template_version>"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  answers {
    project_id = "<project_id2>"
    values = {
      "ingress_host" = "test2.xip.io"
    }
  }
  roles = ["project-member"]
}
```

## Argument Reference

The following arguments are supported:

* `catalog_name` - (Required) The multi cluster app catalog name (string)
* `name` - (Required/ForceNew) The multi cluster app name (string)
* `roles` - (Required) The multi cluster app roles (list)
* `targets` - (Required) The multi cluster app target projects (list)
* `template_name` - (Required) The multi cluster app template name (string)
* `answers` - (Optional/Computed) The multi cluster app answers (list)
* `members` - (Optional) The multi cluster app answers (list)
* `revision_history_limit` - (Computed) The multi cluster app revision history limit. Default `10` (int)
* `revision_id` - (Optional/Computed) Current revision id for the multi cluster app (string)
* `template_version` - (Optional/Computed) The multi cluster app template version. Default: `latest` (string)
* `upgrade_strategy` - (Optional/Computed) The multi cluster app upgrade strategy (list MaxItems:1)
* `wait` - (Optional) Wait until the multi cluster app is active. Default `true` (bool)
* `annotations` - (Optional/Computed) Annotations for multi cluster app object (map)
* `labels` - (Optional/Computed) Labels for multi cluster app object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `template_version_id` - (Computed) The multi cluster app template version ID (string)

## Nested blocks

### `targets`

#### Arguments

* `project_id` - (Required) Project ID for target (string)
* `app_id` - (Computed) App ID for target (string)
* `health_state` - (Computed) App health state for target (string)
* `state` - (Computed) App state for target (string)

### `answers`

#### Arguments

* `cluster_id` - (Optional) Cluster ID for answer (string)
* `project_id` - (Optional) Project ID for target (string)
* `values` - (Optional) Key/values for answer (map)

### `members`

#### Arguments

* `access_type` - (Optional) Member access type. Valid values: `["member" | "owner" | "read-only"]` (string)
* `group_principal_id` - (Optional) Member group principal id (string)
* `user_principal_id` - (Optional) Member user principal id (string)

### `upgrade_strategy`

#### Arguments

* `rolling_update` - (Optional) Upgrade strategy rolling update (list MaxItems:1)

#### `rolling_update`

##### Arguments

* `batch_size` - (Optional) Rolling update batch size. Default `1` (int)
* `interval` - (Optional) Rolling update interval. Default `1` (int)

## Timeouts

`rancher2_app` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating apps.
- `update` - (Default `10 minutes`) Used for app modifications.
- `delete` - (Default `10 minutes`) Used for deleting apps.

## Import

Multi cluster app can be imported using the multi cluster app ID in the format `<multi_cluster_app_name>`

```
$ terraform import rancher2_multi_cluster_app.foo &lt;MULTI_CLUSTER_APP_ID&gt;
```
