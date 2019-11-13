package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterAlertGroupFields() map[string]*schema.Schema {
	r := alertGroupFields()
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert group Cluster ID",
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}
