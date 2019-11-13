---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_template"
sidebar_current: "docs-rancher2-datasource-node_template"
description: |-
  Get information on a Rancher v2 Node Template resource.
---

# rancher2\_node\_template

Use this data source to retrieve information about a Rancher v2 Node Template resource.

## Example Usage

```hcl
data "rancher2_node_template" "foo" {
  name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Template (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `cloud_credential_id` - (Computed) Cloud credential ID for the Node Template. Required from Rancher v2.2.x (string)
* `description` - (Computed) Description for the Node Template (string)
* `driver` - (Computed) The driver of the node template (string)
* `engine_env` - (Computed) Engine environment for the node template (string)
* `engine_insecure_registry` - (Computed) Insecure registry for the node template (list)
* `engine_install_url` - (Computed) Docker engine install URL for the node template (string)
* `engine_label` - (Computed) Engine label for the node template (string)
* `engine_opt` - (Computed) Engine options for the node template (map)
* `engine_registry_mirror` - (Computed) Engine registry mirror for the node template (list)
* `engine_storage_driver` - (Computed) Engine storage driver for the node template (string)
* `use_internal_ip_address` - (Computed) Engine storage driver for the node template (bool)
* `annotations` - (Computed) Annotations for Node Template object (map)
* `labels` - (Computed) Labels for Node Template object (map)

