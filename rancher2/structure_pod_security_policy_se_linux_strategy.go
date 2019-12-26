package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicySELinuxStrategy(in policyv1.SELinuxStrategyOptions) []interface{} {
	
	obj := make(map[string]interface{})

	if len(in.Rule) > 0 {
		obj["rule"] = string(in.Rule)
	}
    if in.SELinuxOptions != nil {
	    obj["se_linux_options"] = flattenPodSecurityPolicySELinuxOptions(in.SELinuxOptions)
    }

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicySELinuxStrategy(in []interface{}) policyv1.SELinuxStrategyOptions {

	obj := policyv1.SELinuxStrategyOptions{}

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = policyv1.SELinuxStrategy(v)
	}

	if v, ok := m["se_linux_options"].([]interface{}); ok && len(v) > 0 {
		obj.SELinuxOptions = expandPodSecurityPolicySELinuxOptions(v)
	}

	return obj
}
