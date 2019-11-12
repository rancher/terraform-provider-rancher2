package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigPrivateRegistriesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"is_default": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}
