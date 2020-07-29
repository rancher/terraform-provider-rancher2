---
page_title: "rancher2_etcd_backup Data Source"
---

# rancher2\_etcd\_backup Data Source

Use this data source to retrieve information about a Rancher v2 etcd backup.

## Example Usage

```hcl
data "rancher2_etcd_backup" "foo" {
  cluster_id = "<CLUSTER_ID>"
  name = "foo"
}
```

## Argument Reference

* `cluster_id` - (Required) Cluster ID to config Etcd Backup (string)
* `name` - (Required) The name of the Etcd Backup (string)


## Attributes Reference

* `id` - (Computed) The ID of the resource (string)
* `backup_config` - (Computed) Backup config for etcd backup (list maxitems:1)
* `filename` - (Computed) Filename of the Etcd Backup (string)
* `manual` - (Computed) Manual execution of the Etcd Backup. Default `false` (bool)
* `namespace_id` - (Computed) Description for the Etcd Backup (string)
* `annotations` - (Computed) Annotations for Etcd Backup object (map)
* `labels` - (Computed) Labels for Etcd Backup object (map)
