package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterAlertGroupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert group Cluster ID",
		},
	}

	for k, v := range alertGroupFields() {
		s[k] = v
	}

	return s
}
