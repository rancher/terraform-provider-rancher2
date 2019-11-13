package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func targetFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Project ID for target",
		},
		"app_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "App ID for target",
		},
		"health_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "App health state for target",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "App state for target",
		},
	}

	return s
}
