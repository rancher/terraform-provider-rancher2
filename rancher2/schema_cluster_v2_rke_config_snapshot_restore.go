package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

//Types

func clusterV2RKEConfigETCDSnapshotRestoreFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ETCD snapshot name to restore",
		},
		"generation": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ETCD snapshot desired generation",
		},
		"restore_rke_config": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ETCD restore RKE config (set to none, all, or kubernetesVersion)",
		},
	}

	return s
}
