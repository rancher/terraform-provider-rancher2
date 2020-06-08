package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func userFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"principal_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"annotations": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
