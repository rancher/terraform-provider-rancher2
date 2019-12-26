package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyRunAsGroup(in *policyv1.RunAsGroupStrategyOptions) []interface{} {

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

func expandPodSecurityPolicyRunAsGroup(in []interface{}) *policyv1.RunAsGroupStrategyOptions {

	obj := &policyv1.RunAsGroupStrategyOptions{}

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = policyv1.RunAsGroupStrategy(v)
	}

	if v, ok := m["ranges"].([]interface{}); ok && len(v) > 0 {
		obj.Ranges = expandPodSecurityPolicyIDRanges(v)
	}

	return obj
}
