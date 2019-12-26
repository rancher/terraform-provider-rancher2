package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyVolumes(in []policyv1.FSType) []interface{} {
	
	out := make([]interface{}, len(in))
	
	for i, v := range in {
		out[i] = string(v)
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyVolumes(in []interface{}) []policyv1.FSType {

	obj := make([]policyv1.FSType, len(in))

    for i, v := range in {
		if s, ok := v.(string); ok {
			obj[i] = policyv1.FSType(s)
		}
	}

	return obj
}
