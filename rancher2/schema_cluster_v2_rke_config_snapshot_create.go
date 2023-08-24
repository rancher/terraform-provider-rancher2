package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

//Types

func clusterV2RKEConfigETCDSnapshotCreateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"generation": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ETCD generation to initiate a snapshot",
		},
	}

	return s
}
