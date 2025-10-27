package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigRotateCertificatesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"generation": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Desired certificate rotation generation.",
		},
		"services": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Service certificates to rotate with this generation.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	return s
}
