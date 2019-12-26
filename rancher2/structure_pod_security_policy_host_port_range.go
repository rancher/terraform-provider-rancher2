package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyHostPortRanges(hp []policyv1.HostPortRange) []interface{} {
	if len(hp) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(hp))
	for i, in := range hp {
        obj := make(map[string]interface{})

        obj["min"] = int(in.Min)
		obj["max"] = int(in.Max)

		out[i] = obj
	}
	
	return out

}

// Expanders

func expandPodSecurityPolicyHostPortRanges(hp []interface{}) []policyv1.HostPortRange {

	if len(hp) == 0 || hp[0] == nil {
		return []policyv1.HostPortRange{}
	}

	obj := make([]policyv1.HostPortRange, len(hp))

	for i := range hp {
		in := hp[i].(map[string]interface{})

		obj[i].Min = int32(in["min"].(int))
		obj[i].Max = int32(in["max"].(int))

	}

	return obj

}
