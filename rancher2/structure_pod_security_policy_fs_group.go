package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyFSGroup(in policyv1.FSGroupStrategyOptions) []interface{} {

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

func expandPodSecurityPolicyFSGroup(in []interface{}) policyv1.FSGroupStrategyOptions {

	obj := policyv1.FSGroupStrategyOptions{}	

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = policyv1.FSGroupStrategyType(v)
	}

	if v, ok := m["ranges"].([]interface{}); ok && len(v) > 0 {
		obj.Ranges = expandPodSecurityPolicyIDRanges(v)
	}

	return obj
}
