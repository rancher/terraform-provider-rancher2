---
page_title: "rancher2_etcd_backup Resource"
---

# rancher2\_etcd\_backup Resource

Provides a Rancher v2 Etcd Backup resource. This can be used to create an Etcd Backup for Rancher v2.2.x and above, and to retrieve their information. 

The `rancher2_etcd_backup` resource is used to define extra etcd backups for a `rancher2_cluster`, which will be created as a local or S3 backup in accordance with the etcd backup config for the cluster. The main etcd backup config for the cluster should be set on the [cluster config](https://www.terraform.io/docs/providers/rancher2/r/cluster.html#backup_config-1)

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
      folder = "/folder"
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

* `cluster_id` - (Required) Cluster ID to config Etcd Backup (string)
* `backup_config` - (Optional/Computed) Backup config for etcd backup (list maxitems:1)
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

* `access_key` - (Optional/Sensitive) Access key for S3 service (string)
* `bucket_name` - (Required) Bucket name for S3 service (string)
* `custom_ca` - (Optional) Base64 encoded custom CA for S3 service. Use filebase64(<FILE>) for encoding file. Available from Rancher v2.2.5 (string)
* `endpoint` - (Required) Endpoint for S3 service (string)
* `folder` - (Optional) Folder for S3 service. Available from Rancher v2.2.7 (string)
* `region` - (Optional) Region for S3 service (string)
* `secret_key` - (Optional/Sensitive) Secret key for S3 service (string)

## Timeouts

`rancher2_etcd_backup` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cloud credentials.
- `update` - (Default `10 minutes`) Used for cloud credential modifications.
- `delete` - (Default `10 minutes`) Used for deleting cloud credentials.

## Import

Etcd Backup can be imported using the Rancher etcd backup ID

```
$ terraform import rancher2_etcd_backup.foo &lt;ETCD_BACKUP_ID&gt;
```
