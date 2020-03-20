package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyAllowedCSIDriverFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name is the registered name of the CSI driver",
			Required:    true,
		},
	}

	return s
}
