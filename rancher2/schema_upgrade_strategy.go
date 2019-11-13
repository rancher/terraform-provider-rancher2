package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rollingUpdateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"batch_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Rolling update batch size",
		},
		"interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Rolling update interval",
		},
	}

	return s
}

func upgradeStrategyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"rolling_update": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Rolling update for upgrade strategy",
			Elem: &schema.Resource{
				Schema: rollingUpdateFields(),
			},
		},
	}

	return s
}
