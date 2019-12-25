package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyIDRanges(idRange []policyv1.IDRange) []interface{} {
	if len(idRange) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(idRange))
	
	for i, in := range idRange {
		obj := make(map[string]interface{})

		obj["min"] = int(in.Min)
		obj["max"] = int(in.Max)

		out[i] = obj
	}

	return out
}

// Expanders

func expandPodSecurityPolicyIDRanges(idRange []interface{}) []policyv1.IDRange {

	if len(idRange) == 0 || idRange[0] == nil {
		return []policyv1.IDRange{}
	}

	obj := make([]policyv1.IDRange, len(idRange))

	for i := range idRange {
		in := idRange[i].(map[string]interface{})

		obj[i].Min = int64(in["min"].(int))
		obj[i].Max = int64(in["max"].(int))

	}

	return obj
}
