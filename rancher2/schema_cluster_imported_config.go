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
	}
}
