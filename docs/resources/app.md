---
page_title: "rancher2_app Resource"
---

# rancher2\_app Resource

Provides a Rancher v2 app resource. This can be used to deploy apps within Rancher v2 projects.

This resource can also modify Rancher v2 apps in 3 ways:
- `Update`: If `description`, `annotations` or `labels` arguments are modified the app will be updated. No new `revision_id` will be generated in Rancher.
- `Upgrade`: If `answers`, `catalog_name`, `template_name`, `template_version` or `values_yaml` arguments are modified, the app will be upgraded. A new `revision_id` will be generated in Rancher.
- `Rollback`: If `revision_id` argument is provided or modified the app will be rolled back accordingly. A new `revision_id` will be generated in Rancher. It will also generate a non-empty terraform plan that will require manual .tf file intervention. Use carefully.

Note: In case of multiple resource modifications in a row, `rollback` has preference over `upgrade`.

## Example Usage

```hcl
# Create a new rancher2 App
resource "rancher2_app" "foo" {
  catalog_name = "<catalog_name>"
  name = "foo"
  description = "Foo app"
  project_id = "<project_id>"
  template_name = "<template_name>"
  template_version = "<template_version>"
  target_namespace = "<namespace_name>"
  answers = {
    "ingress_host" = "test.xip.io"
    "foo" = "bar"
    "ingress.annotations.nginx.ingress.kubernetes.io/force-ssl-redirect" = true
  }
}
```

```hcl
# Create a new rancher2 App in a new namespace
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace"
  project_id = "<project_id>"
  resource_quota {
    limit {
      limits_cpu = "100m"
      limits_memory = "100Mi"
      requests_storage = "1Gi"
    }
  }
}

resource "rancher2_app" "foo" {
  catalog_name = "<catalog_name>"
  name = "foo"
  description = "Foo app"
  project_id = "<project_id>"
  template_name = "<template_name>"
  template_version = "<template_version>"
  target_namespace = rancher2_namespace.foo.id
  answers = {
    "ingress_host" = "test.xip.io"
    "foo" = "bar"
    "ingress.annotations.nginx.ingress.kubernetes.io/force-ssl-redirect" = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `catalog_name` - (Required) Catalog name of the app. If modified, app will be upgraded. For use scoped catalogs:
  * add cluster ID before name, `local:<name>` or `c-XXXXX:<name>`
  * add project ID before name, `p-XXXXX:<name>`
* `name` - (Required/ForceNew) The name of the app (string)
* `project_id` - (Required/ForceNew) The project id where the app will be installed (string)
* `target_namespace` - (Required/ForceNew) The namespace id where the app will be installed (string)
* `template_name` - (Required) Template name of the app. If modified, app will be upgraded (string)
* `answers` - (Optional) Answers for the app template. If modified, app will be upgraded (map)
* `description` - (Optional/Computed) Description for the app (string)
* `force_upgrade` - (Optional) Force app upgrade (string)
* `revision_id` - (Optional/Computed) Current revision id for the app. If modified, If this argument is provided or modified, app will be rollbacked to `revision_id` (string)
* `template_version` - (Optional/Computed) Template version of the app. If modified, app will be upgraded. Default: `latest` (string)
* `values_yaml` - (Optional) values.yaml base64 encoded file content for the app template. If modified, app will be upgraded (string)
* `wait` - (Optional) Wait until app is deployed and active. Default: `true` (bool)
* `annotations` - (Optional/Computed) Annotations for App object (map)
* `labels` - (Optional/Computed) Labels for App object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `external_id` - (Computed) The url of the app template on a catalog (string)

## Timeouts

`rancher2_app` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating apps.
- `update` - (Default `10 minutes`) Used for app modifications.
- `delete` - (Default `10 minutes`) Used for deleting apps.

## Import

Apps can be imported using the app ID in the format `<project_id>:<app_name>`

```
$ terraform import rancher2_app.foo &lt;PROJECT_ID_ID&gt;:&lt;APP_NAME&gt;
```
