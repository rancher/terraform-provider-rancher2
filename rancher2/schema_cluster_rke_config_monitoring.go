package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigMonitoringFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}
