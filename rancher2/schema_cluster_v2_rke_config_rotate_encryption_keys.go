package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigRotateEncryptionKeysFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"generation": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Desired encryption keys rotation generation.",
		},
	}

	return s
}
