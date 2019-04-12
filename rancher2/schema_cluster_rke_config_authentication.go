package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func clusterRKEConfigAuthenticationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"sans": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"strategy": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}
