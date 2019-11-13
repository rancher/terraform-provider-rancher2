package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func appFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"catalog_name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Catalog name of the app",
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the app",
		},
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add app",
		},
		"target_namespace": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Namespace name to add app",
		},
		"template_name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Template name of the app",
		},
		"answers": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Answers of the app",
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"external_id": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "External ID of the app",
		},
		"force_upgrade": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force app upgrade",
		},
		"revision_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "App revision id",
		},
		"template_version": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Template version of the app",
		},
		"values_yaml": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "values.yaml base64 encoded file content of the app",
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the app",
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
