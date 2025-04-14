package rancher2

import (
	v3 "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func flattenAgentSchedulingCustomization(in *v3.AgentSchedulingCustomization) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.PriorityClass != nil {
		obj["priority_class"] = flattenClusterAgentPriorityClass(in.PriorityClass)
	}

	if in.PodDisruptionBudget != nil {
		obj["pod_disruption_budget"] = flattenPodDisruptionBudget(in.PodDisruptionBudget)
	}

	return []interface{}{obj}
}

func expandAgentSchedulingCustomization(p []interface{}) *v3.AgentSchedulingCustomization {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v3.AgentSchedulingCustomization{}

	in := p[0].(map[string]interface{})

	if v, ok := in["priority_class"].([]interface{}); ok && v != nil {
		obj.PriorityClass = expandClusterAgentPriorityClass(v)
	}

	if v, ok := in["pod_disruption_budget"].([]interface{}); ok && v != nil {
		obj.PodDisruptionBudget = expandPodDisruptionBudget(v)
	}

	return obj
}
