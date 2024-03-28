package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenResourceRequirements(in *managementClient.ResourceRequirements) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if in.Limits != nil {
		if len(in.Limits["cpu"]) > 0 {
			obj["cpu_limit"] = in.Limits["cpu"]
		}
		if len(in.Limits["memory"]) > 0 {
			obj["memory_limit"] = in.Limits["memory"]
		}
	}
	if in.Requests != nil {
		if len(in.Requests["cpu"]) > 0 {
			obj["cpu_request"] = in.Requests["cpu"]
		}
		if len(in.Requests["memory"]) > 0 {
			obj["memory_request"] = in.Requests["memory"]
		}
	}

	return []interface{}{obj}
}

// Expanders

func expandResourceRequirements(p []interface{}) *managementClient.ResourceRequirements {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &managementClient.ResourceRequirements{}

	for i := range p {
		in := p[i].(map[string]interface{})

		obj.Limits = make(map[string]string)
		obj.Requests = make(map[string]string)

		if v, ok := in["cpu_limit"].(string); ok && len(v) > 0 {
			obj.Limits["cpu"] = v
		}

		if v, ok := in["cpu_request"].(string); ok && len(v) > 0 {
			obj.Requests["cpu"] = v
		}

		if v, ok := in["memory_limit"].(string); ok && len(v) > 0 {
			obj.Limits["memory"] = v
		}

		if v, ok := in["memory_request"].(string); ok && len(v) > 0 {
			obj.Requests["memory"] = v
		}
	}

	return obj
}
