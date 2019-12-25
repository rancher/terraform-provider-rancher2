package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicySupplementalGroup(in *policyv1.SupplementalGroupsStrategyOptions) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.Rule) > 0 {
		obj["rule"] = string(in.Rule)
	}
    if len(in.Ranges) > 0 {
	    obj["ranges"] = flattenPodSecurityPolicyIDRanges(in.Ranges)
    }

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicySupplementalGroup(supplementalGroups []interface{}) *policyv1.SupplementalGroupsStrategyOptions {

	obj := &policyv1.SupplementalGroupsStrategyOptions{}

	if len(supplementalGroups) == 0 || supplementalGroups[0] == nil {
		return obj
	}

	in := supplementalGroups[0].(map[string]interface{})

	if v, ok := in["rule"].(string); ok {
		obj.Rule = policyv1.SupplementalGroupsStrategyType(v)
	}

	if v, ok := in["ranges"].([]interface{}); ok && len(v) > 0 {
		obj.Ranges = expandPodSecurityPolicyIDRanges(v)
	}

	return obj
}
