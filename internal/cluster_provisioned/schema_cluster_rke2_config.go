package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterRKE2Kind   = "rke2"
	clusterDriverRKE2 = "rke2"
)

//Schemas

func clusterRKE2ConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"upgrade_strategy": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The RKE2 upgrade strategy",
			Elem: &schema.Resource{
				Schema: clusterK3SUpgradeStrategyConfigFields(),
			},
		},
		"version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The RKE2 kubernetes version",
		},
	}

	return s
}
