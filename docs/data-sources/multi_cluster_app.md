---
page_title: "rancher2_multi_cluster_app Data Source"
---

# rancher2\_multi\_cluster\_app Data Source

Use this data source to retrieve information about a Rancher v2 multi cluster app.

## Example Usage

```
data "rancher2_multi_cluster_app" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The multi cluster app name (string)

## Attributes Reference

* `catalog_name` - (Computed) The multi cluster app catalog name (string)
* `id` - (Computed) The ID of the resource (string)
* `roles` - (Computed) The multi cluster app roles (list)
* `targets` - (Computed) The multi cluster app target projects (list)
* `template_name` - (Computed) The multi cluster app template name (string)
* `template_version` - (Computed) The multi cluster app template version (string)
* `template_version_id` - (Computed) The multi cluster app template version ID (string)
* `answers` - (Computed) The multi cluster app answers (list)
* `members` - (Computed) The multi cluster app members (list)
* `revision_history_limit` - (Computed) The multi cluster app revision history limit (int)
* `revision_id` - (Computed) Current revision id for the multi cluster app (string)
* `upgrade_strategy` - (Computed) The multi cluster app upgrade strategy (list)
* `annotations` - (Computed) Annotations for multi cluster app object (map)
* `labels` - (Computed) Labels for multi cluster app object (map)
