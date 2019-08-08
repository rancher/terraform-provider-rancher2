---
layout: "rancher2"
page_title: "Rancher2: rancher2_app"
sidebar_current: "docs-rancher2-datasource-app"
description: |-
  Get information on a Rancher v2 app.
---

# rancher2\_app

Use this data source to retrieve information about a Rancher v2 app.

## Example Usage

```
data "rancher2_app" "rancher2" {
    name = "backoffice"
    project_id = "c-1a2bc:p-3d4fg"
    target_namespace = "webservers"
}
```

## Argument Reference

* `name` - (Required) The app name.
* `project_id` - (Required) The id of the project where the app is installed.
* `target_namespace` - (Optional/Computed) The namespace where the app is installed.

## Attributes Reference

* `external_id` - (Computed) The URL of the helm catalog app (string)
* `annotations` - (Computed) Annotations for the catalog (map)
* `answers` - (Computed) Answers for the app (map)
* `labels` - (Computed) Labels for the catalog (map)
* `values_yaml` - (Computed) values.yaml file content for the app (string)
