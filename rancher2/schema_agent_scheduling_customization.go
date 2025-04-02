package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func agentSchedulingCustomizationFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"priority_class": {
			Type:        schema.TypeList,
			Description: "The Priority Class created for the cattle cluster agent",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: priorityClassFields(),
			},
		},
		"pod_disruption_budget": {
			Type:        schema.TypeList,
			Description: "The Pod Disruption Budget created for the cattle cluster agent",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podDisruptionBudgetFields(),
			},
		},
	}
}

func priorityClassFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:        schema.TypeInt,
			Description: "The priority value for the cattle cluster agent. Must be between negative 1 billion and 1 billion.",
			Required:    true,
		},
		"preemption_policy": {
			Type:        schema.TypeString,
			Description: "The preemption behavior for the cattle cluster agent. Must be either 'PreemptLowerPriority' or 'Never'",
			Optional:    true,
		},
	}
}

func podDisruptionBudgetFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"min_available": {
			Type:        schema.TypeString,
			Description: "The minimum number of cattle cluster agent replicas that must be running at a given time.",
			Optional:    true,
		},
		"max_unavailable": {
			Type:        schema.TypeString,
			Description: "The maximum number of cattle cluster agent replicas that can be down at a given time.",
			Optional:    true,
		},
	}
}
