package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyHostPortRange(in *policyv1.HostPortRange) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["min"] = int(in.Min)
	obj["max"] = int(in.Max)

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicyHostPortRange(hpRange []interface{}) *policyv1.HostPortRange {

	obj := &policyv1.HostPortRange{}

	if len(hpRange) == 0 || hpRange[0] == nil {
		return obj
	}

	in := hpRange[0].(map[string]interface{})

	obj.Min = int32(in["min"].(int))
	obj.Max = int32(in["max"].(int))

	return obj
}
