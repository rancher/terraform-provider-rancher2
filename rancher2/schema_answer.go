package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func answerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Cluster ID for answer",
		},
		"project_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Project ID for answer",
		},
		"values": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Key/values for answer",
		},
	}

	return s
}
