package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedProcMountTypes(p []v1.ProcMountType) []interface{} {

	out := make([]interface{}, len(p))

	for i, in := range p {
		out[i] = string(in)
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedProcMountTypes(p []interface{}) []v1.ProcMountType {

	obj := make([]v1.ProcMountType, len(p))

    for i, in := range p {
		if s, ok := in.(string); ok {
			obj[i] = v1.ProcMountType(s)
		}
	}

	return obj
}