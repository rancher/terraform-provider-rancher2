package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterImportedConfigFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"private_registry_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Private registry URL",
		},
		"private_registry_pull_secrets": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Private registry image pull secrets",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
