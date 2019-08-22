---
layout: "rancher2"
page_title: "Rancher2: rancher2_cluster_logging"
sidebar_current: "docs-rancher2-datasource-cluster_logging"
description: |-
  Get information on a Rancher v2 Cluster Logging resource. 
---

# rancher2\_cluster\_logging

Use this data source to retrieve information about a Rancher v2 Cluster Logging.

## Example Usage

```hcl
data "rancher2_cluster_logging" "foo" {
  cluster_id = "<cluster_id>"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The cluster id to configure logging (string)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `kind` - (Computed) The kind of the Cluster Logging. `elasticsearch`, `fluentd`, `kafka`, `splunk` and `syslog` are supported (string)
* `elasticsearch_config` - (Computed) The elasticsearch config for Cluster Logging. For `kind = elasticsearch`  (list maxitems:1)
* `fluentd_config` - (Computed) The fluentd config for Cluster Logging. For `kind = fluentd` (list maxitems:1)
* `kafka_config` - (Computed) The kafka config for Cluster Logging. For `kind = kafka` (list maxitems:1)
* `name` - (Computed) The name of the cluster logging config (string)
* `namespace_id` - (Computed) The namespace id from cluster logging (string)
* `output_flush_interval` - (Computed) How often buffered logs would be flushed. Default: `3` seconds (int)
* `output_tags` - (computed) The output tags for Cluster Logging (map)
* `splunk_config` - (Computed) The splunk config for Cluster Logging. For `kind = splunk` (list maxitems:1)
* `syslog_config` - (Computed) The syslog config for Cluster Logging. For `kind = syslog` (list maxitems:1)
* `annotations` - (Computed) Annotations for Cluster Logging object (map)
* `labels` - (Computed) Labels for Cluster Logging object (map)

