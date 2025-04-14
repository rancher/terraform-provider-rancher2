package rancher2

import (
	"github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
)

// Flatteners

func flattenPodDisruptionBudgetV2(in *v1.PodDisruptionBudgetSpec) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.MinAvailable != "" {
		obj["min_available"] = in.MinAvailable
	}

	if in.MaxUnavailable != "" {
		obj["max_unavailable"] = in.MaxUnavailable
	}

	return []interface{}{obj}
}

// Expanders

func expandPodDisruptionBudgetV2(p []interface{}) *v1.PodDisruptionBudgetSpec {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v1.PodDisruptionBudgetSpec{}
	in := p[0].(map[string]interface{})

	if v, ok := in["min_available"].(string); ok {
		obj.MinAvailable = v
	}

	if v, ok := in["max_unavailable"].(string); ok {
		obj.MaxUnavailable = v
	}

	return obj
}
