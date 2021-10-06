---
page_title: "Rancher2: rancher2_catalog_v2 Resource"
---

# rancher2\_app\_v2 Resource

Provides a Rancher App v2 resource. This can be used to manage helm charts for Rancher v2 environments and retrieve their information. App v2 resource is available at Rancher v2.5.x and above.

## Example Usage

```hcl
# Create a new Rancher2 App V2 using
resource "rancher2_app_v2" "foo" {
  cluster_id = "<CLUSTER_ID>"
  name = "rancher-monitoring"
  namespace = "cattle-monitoring-system"
  repo_name = "rancher-charts"
  chart_name = "rancher-monitoring"
  chart_version = "9.4.200"
  values = file("values.yaml")
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required/ForceNew) The cluster id of the app (string)
* `name` - (Required/ForceNew) The name of the app v2 (string)
* `namespace` - (Required/ForceNew) The namespace of the app v2 (string)
* `repo_name` - (Required) Repo name (string)
* `chart_name` - (Required) The app v2 chart name (string)
* `chart_version` - (Optional/Computed) The app v2 chart version (string)
* `project_id` - (Optional) Deploy the app v2 within project ID (string)
* `values` - (Optional/Sensitive) The app v2 values yaml. Yaml format is required (string)
* `cleanup_on_fail` - (Optional) Cleanup app v2 on failed chart upgrade. Default: `false` (bool)
* `disable_hooks` - (Optional) Disable app v2 chart hooks. Default: `false` (bool)
* `disable_open_api_validation` - (Optional) Disable app V2 Open API Validation. Default: `false` (bool)
* `force_upgrade` - (Optional) Force app V2 chart upgrade. Default: `false` (bool)
* `wait` - (Optional) Wait until app is deployed. Default: `true` (bool)
* `annotations` - (Optional/Computed) Annotations for the app v2 (map)
* `labels` - (Optional/Computed) Labels for the app v2 (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cluster_name` - (Computed) The cluster name of the app (string)
* `system_default_registry` - (Computed) The system default registry of the app (string)

## Timeouts

`rancher2_app_v2` provides the following
[Timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating Rancher v2 apps.
- `update` - (Default `10 minutes`) Used for Rancher v2 app modifications.
- `delete` - (Default `10 minutes`) Used for deleting Rancher v2 apps.

## Import

V2 apps can be imported using the Rancher cluster ID and App V2 name, which is composed of `<namespace>/<application_name>`.

```
$ terraform import rancher2_app_v2.foo &lt;CLUSTER_ID&gt;.&lt;APP_V2_NAME&gt;
```
