package rancher2

import (
	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenClusterAgentPriorityClassV2(in *v1.PriorityClassSpec) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.PreemptionPolicy != nil && *in.PreemptionPolicy != "" {
		obj["preemption_policy"] = *in.PreemptionPolicy
	}

	obj["value"] = in.Value
	return []interface{}{obj}
}

// Expanders

func expandClusterAgentPriorityClassV2(p []interface{}) *v1.PriorityClassSpec {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v1.PriorityClassSpec{}

	in := p[0].(map[string]interface{})

	if v, ok := in["value"].(int); ok {
		obj.Value = v
	}

	if v, ok := in["preemption_policy"].(string); ok && v != "" {
		preemption := corev1.PreemptionPolicy(v)
		obj.PreemptionPolicy = &preemption
	}

	return obj
}
