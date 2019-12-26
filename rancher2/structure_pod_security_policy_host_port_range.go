package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyHostPortRanges(in []policyv1.HostPortRange) []interface{} {

	out := make([]interface{}, len(in))

	for i, v := range in {
        out[i] = map[string]interface{}{
			"min": int(v.Min),
			"max": int(v.Max),
		}
	}
	
	return out

}

// Expanders

func expandPodSecurityPolicyHostPortRanges(in []interface{}) []policyv1.HostPortRange {

	obj := make([]policyv1.HostPortRange, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = policyv1.HostPortRange{
				Min: int32(m["min"].(int)),
				Max: int32(m["max"].(int)),
			}
		}
	}

	return obj

}
