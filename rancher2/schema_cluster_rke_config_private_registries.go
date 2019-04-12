package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func clusterRKEConfigPrivateRegistriesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
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
		"url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"user": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}
