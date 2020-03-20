package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyAllowedFlexVolumesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"driver": {
			Type:        schema.TypeString,
			Description: "driver is the name of the Flexvolume driver.",
			Required:    true,
		},
	}

	return s
}
