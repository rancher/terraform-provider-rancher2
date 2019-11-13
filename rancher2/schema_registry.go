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
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Project ID to add docker registry",
		},
		"name": &schema.Schema{
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
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the docker registry",
		},
		"namespace_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Namespace ID to add docker registry",
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the docker registry",
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the docker registry",
		},
	}

	return s
}
