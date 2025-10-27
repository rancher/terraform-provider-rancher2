package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigETCDSnapshotS3Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"bucket": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ETCD snapshot S3 bucket",
		},
		"cloud_credential_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD snapshot S3 cloud credential name",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ETCD snapshot S3 endpoint",
		},
		"endpoint_ca": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD snapshot S3 endpoint CA",
		},
		"folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD snapshot S3 folder",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD snapshot S3 region",
		},
		"skip_ssl_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable ETCD skip ssl verify",
		},
	}

	return s
}

func clusterV2RKEConfigETCDFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"disable_snapshots": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable ETCD snapshots",
		},
		"snapshot_schedule_cron": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD snapshot schedule cron (e.g `\"0 */5 * * *\"`)",
		},
		"snapshot_retention": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ETCD snapshot retention",
		},
		"s3_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "ETCD snapshot S3 config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigETCDSnapshotS3Fields(),
			},
		},
	}

	return s
}
