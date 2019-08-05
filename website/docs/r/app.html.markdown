---
layout: "rancher2"
page_title: "Rancher2: rancher2_app"
sidebar_current: "docs-rancher2-resource-app"
description: |-
  Provides a Rancher v2 app resource. This can be used to install apps on rancher v2 namespaces.
---

# rancher2\_app

Provides a Rancher v2 App resource. This can be used to install apps on rancher v2 namespaces.

## Example Usage

```hcl
# Create a new rancher2 App
resource "rancher2_app" "foo" {
  name = "foo"
  description = "Terraform app foo"
  project_id = "<project_id>"
  target_namespace = "<namespace_name>"
  external_id = "catalog://?catalog=<catalog_name>&template=<template_name>&version=<template_version>"
  answers = {
    something = "Anything"
    "web.env[0].name" = "Rancher2 is cool"
    "ingress.annotations.nginx.ingress.kubernetes.io/force-ssl-redirect" = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the app (string)
* `project_id` - (Required/ForceNew) The project id where the app will be installed (string)
* `target_namespace` - (Required/ForceNew) The namespace name where the app will be installed (string)
* `external_id` - (Required) The url of the app template on a catalog (string)
* `answers` - (Optional) Answers for the app template (map)
* `description` - (Optional) An app description (string)
* `annotations` - (Optional/Computed) Annotations for App object (map)
* `labels` - (Optional/Computed) Labels for App object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_app` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating apps.
- `update` - (Default `10 minutes`) Used for app modifications.
- `delete` - (Default `10 minutes`) Used for deleting apps.

## Import

Apps can be imported using the app ID in the format `<<project_id>.<app_id>`

```
$ terraform import rancher2_app.foo <project_id>.<app_id>
```