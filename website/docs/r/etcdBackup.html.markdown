---
layout: "rancher2"
page_title: "Rancher2: rancher2_etcd_backup"
sidebar_current: "docs-rancher2-resource-etcd_backup"
description: |-
  Provides a Rancher v2.2.x Etcd Backup resource. This can be used to create Etcd Backup for rancher v2.2 node templates and retrieve their information.
---

# rancher2\_etcd\_backup

Provides a Rancher v2 Etcd Backup resource. This can be used to create Etcd Backup for rancher v2.2.x and retrieve their information. 

## Example Usage

```hcl
# Create a new rancher2 Etcd Backup
resource "rancher2_etcd_backup" "foo" {
  backup_config {
    enabled = true
    interval_hours = 20
    retention = 10
    s3_backup_config {
      access_key = "access_key"
      bucket_name = "bucket_name"
      endpoint = "endpoint"
      region = "region"
      secret_key = "secret_key"
    }
  }
  cluster_id = "<CLUSTER_ID>"
  name = "foo"
  filename = "<FILENAME>"
}
```

## Argument Reference

The following arguments are supported:

* `backup_config` - (Optional/Computed) Backup config for etcd backup (list maxitems:1)
* `cluster_id` - (Required) Cluster ID to config Etcd Backup (string)
* `filename` - (Optional/Computed) Filename of the Etcd Backup (string)
* `manual` - (Optional) Manual execution of the Etcd Backup. Default `false` (bool)
* `name` - (Required) The name of the Etcd Backup (string)
* `namespace_id` - (Optional/Computed) Description for the Etcd Backup (string)
* `annotations` - (Optional) Annotations for Etcd Backup object (map)
* `labels` - (Optional/Computed) Labels for Etcd Backup object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Nested blocks

### `backup_config`

#### Arguments

* `enabled` - (Optional) Enable etcd backup (bool)
* `interval_hours` - (Optional) Interval hours for etcd backup. Default `12` (int)
* `retention` - (Optional) Retention for etcd backup. Default `6` (int)
* `s3_backup_config` - (Optional) S3 config options for etcd backup. Valid for `imported` and `rke` clusters. (list maxitems:1)

#### `s3_backup_config`

##### Arguments

* `access_key` - (Required/Sensitive) Access key for S3 service (string)
* `bucket_name` - (Required) Bucket name for S3 service (string)
* `endpoint` - (Required) Endpoint for S3 service (string)
* `region` - (Required) Region for S3 service (string)
* `secret_key` - (Required/Sensitive) Secret key for S3 service (string)

## Timeouts

`rancher2_etcd_backup` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cloud credentials.
- `update` - (Default `10 minutes`) Used for cloud credential modifications.
- `delete` - (Default `10 minutes`) Used for deleting cloud credentials.

## Import

Etcd Backup can be imported using the rancher etcd backup ID

```
$ terraform import rancher2_etcd_backup.foo <etcd_backup_id>
```
