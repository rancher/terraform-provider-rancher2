package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//Schemas

func appFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"catalog_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Catalog name of the app",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the app",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add app",
		},
		"target_namespace": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Namespace name to add app",
		},
		"template_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Template name of the app",
		},
		"answers": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Answers of the app",
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"external_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "External ID of the app",
		},
		"force_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force app upgrade",
		},
		"revision_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "App revision id",
		},
		"template_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Template version of the app",
		},
		"values_yaml": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "values.yaml base64 encoded file content of the app",
			ValidateFunc: validation.StringIsBase64,
			StateFunc: func(val interface{}) string {
				s, _ := Base64Decode(val.(string))
				return Base64Encode(TrimSpace(s))
			},
		},
		"wait": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Wait until app is deployed and active",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
