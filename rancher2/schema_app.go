package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func appFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add app",
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    false,
			Description: "Name of the app",
		},
		"target_namespace": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Namespace name to add app",
		},
		"external_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    false,
			Description: "External ID of the app",
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the app",
		},
		"answers": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    false,
			Description: "Answers of the app",
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the app",
		},
	}

	return s
}
