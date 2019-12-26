package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicyCapabilities(v []v1.Capability) []interface{} {
	
	out := make([]interface{}, len(v))

	for i, in := range v {
		out[i] = string(in)
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyCapabilities(v []interface{}) []v1.Capability {

	obj := make([]v1.Capability, len(v))

    for i, in := range v {
		if s, ok := in.(string); ok {
			obj[i] = v1.Capability(s)
		}
	}
	
	return obj
}
