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

func rollingUpdateDaemonSetFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"max_unavailable": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Rolling update max unavailable",
		},
	}

	return s
}

func rollingUpdateDeploymentFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"max_surge": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Rolling update max surge",
		},
		"max_unavailable": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Rolling update max unavailable",
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

func daemonSetStrategyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"rolling_update": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Rolling update for update strategy",
			Elem: &schema.Resource{
				Schema: rollingUpdateDaemonSetFields(),
			},
		},
		"strategy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Strategy",
		},
	}

	return s
}

func deploymentStrategyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"rolling_update": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Rolling update for update strategy",
			Elem: &schema.Resource{
				Schema: rollingUpdateDeploymentFields(),
			},
		},
		"strategy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Strategy",
		},
	}

	return s
}
