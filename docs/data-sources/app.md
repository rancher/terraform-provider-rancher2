---
page_title: "rancher2_app Data Source"
---

# rancher2\_app Data Source

Use this data source to retrieve information about a Rancher v2 app.

## Example Usage

```
data "rancher2_app" "rancher2" {
    name = "foo"
    project_id = "<project_id>"
    target_namespace = "<namespace_name>"
}
```

## Argument Reference

* `name` - (Required) The app name (string)
* `project_id` - (Required) The id of the project where the app is deployed (string)
* `target_namespace` - (Optional/Computed) The namespace name where the app is deployed (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `catalog_name` - (Computed) Catalog name of the app (string)
* `answers` - (Computed) Answers for the app (map)
* `description` - (Computed) Description for the app (string)
* `external_id` - (Computed) The URL of the helm catalog app (string)
* `revision_id` - (Computed) Current revision id for the app (string)
* `template_name` - (Computed) Template name of the app (string)
* `template_version` - (Computed) Template version of the app (string)
* `values_yaml` - (Computed) values.yaml base64 encoded file content for the app (string)
* `annotations` - (Computed) Annotations for the catalog (map)
* `labels` - (Computed) Labels for the catalog (map)
