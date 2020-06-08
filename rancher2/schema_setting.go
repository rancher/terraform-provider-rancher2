package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func settingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
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
