package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func alertGroupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert group name",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Alert group description",
		},
		"group_interval_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     180,
			Description: "Alert group interval seconds",
		},
		"group_wait_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     180,
			Description: "Alert group wait seconds",
		},
		"recipients": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Alert group recipients",
			Elem: &schema.Resource{
				Schema: recipientFields(),
			},
		},
		"repeat_interval_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3600,
			Description: "Alert group repeat interval seconds",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
