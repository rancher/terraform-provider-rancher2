package rancher2

import v1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"

func flattenAgentSchedulingCustomizationV2(in *v1.AgentSchedulingCustomization) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.PriorityClass != nil {
		obj["priority_class"] = flattenClusterAgentPriorityClassV2(in.PriorityClass)
	}

	if in.PodDisruptionBudget != nil {
		obj["pod_disruption_budget"] = flattenPodDisruptionBudgetV2(in.PodDisruptionBudget)
	}

	return []interface{}{obj}
}

func expandAgentSchedulingCustomizationV2(p []interface{}) *v1.AgentSchedulingCustomization {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v1.AgentSchedulingCustomization{}

	in := p[0].(map[string]interface{})

	if v, ok := in["priority_class"].([]interface{}); ok {
		obj.PriorityClass = expandClusterAgentPriorityClassV2(v)
	}

	if v, ok := in["pod_disruption_budget"].([]interface{}); ok {
		obj.PodDisruptionBudget = expandPodDisruptionBudgetV2(v)
	}

	return obj
}
