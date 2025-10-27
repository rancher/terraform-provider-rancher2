package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Shemas

func secretFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeMap,
			Required:    true,
			Sensitive:   true,
			Description: "Secret data base64 encoded",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add secret",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Secret description",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Secret name",
		},
		"namespace_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Namespace ID to add secret",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
