package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func registryCredentialFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}

	return s
}

func registryFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add docker registry",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the docker registry",
		},
		"registries": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: registryCredentialFields(),
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the docker registry",
		},
		"namespace_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Namespace ID to add docker registry",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
