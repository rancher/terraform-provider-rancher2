package rancher2

import (
	v3 "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterAgentPriorityClass(in *v3.PriorityClassSpec) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.PreemptionPolicy != "" {
		obj["preemption_policy"] = in.PreemptionPolicy
	}

	obj["value"] = in.Value
	return []interface{}{obj}
}

// Expanders

func expandClusterAgentPriorityClass(p []interface{}) *v3.PriorityClassSpec {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v3.PriorityClassSpec{}

	in := p[0].(map[string]interface{})

	if v, ok := in["value"].(int); ok {
		obj.Value = int64(v)
	}

	if v, ok := in["preemption_policy"].(string); ok && v != "" {
		obj.PreemptionPolicy = v
	}

	return obj
}
