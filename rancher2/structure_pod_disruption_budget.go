package rancher2

import (
	v3 "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodDisruptionBudget(in *v3.PodDisruptionBudgetSpec) []interface{} {
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

func expandPodDisruptionBudget(p []interface{}) *v3.PodDisruptionBudgetSpec {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &v3.PodDisruptionBudgetSpec{
		MaxUnavailable: "",
		MinAvailable:   "",
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["min_available"].(string); ok && v != "" {
		obj.MinAvailable = v
	}

	if v, ok := in["max_unavailable"].(string); ok && v != "" {
		obj.MaxUnavailable = v
	}

	return obj
}
