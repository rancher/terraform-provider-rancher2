---
page_title: "rancher2_cluster_template Data Source"
---

# rancher2\_cluster\_template Data Source

Use this data source to retrieve information about a Rancher v2 cluster template.

Cluster Templates are available from Rancher v2.3.x and above.

## Example Usage

```
data "rancher2_cluster_template" "foo" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The cluster template name (string)
* `decription` - (Optional/Computed) The cluster template description (string)

## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `default_revision_id` - (Computed) Default cluster template revision ID (string)
* `members` - (Computed) Cluster template members (list)
* `template_revisions` - (Computed) Cluster template revisions (list)
* `annotations` - (Computed) Annotations for the cluster template (map)
* `labels` - (Computed) Labels for the cluster template (map)
