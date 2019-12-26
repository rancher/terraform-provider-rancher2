package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyIDRanges(in []policyv1.IDRange) []interface{} {

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

func expandPodSecurityPolicyIDRanges(in []interface{}) []policyv1.IDRange {

	obj := make([]policyv1.IDRange, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = policyv1.IDRange{
				Min: int64(m["min"].(int)),
				Max: int64(m["max"].(int)),
			}
		}
	}

	return obj
}
