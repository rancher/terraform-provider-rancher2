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

func expandPodSecurityPolicySELinuxStrategy(selinux []interface{}) policyv1.SELinuxStrategyOptions {

	obj := policyv1.SELinuxStrategyOptions{}

	if len(selinux) == 0 || selinux[0] == nil {
		return obj
	}

	in := selinux[0].(map[string]interface{})

	if v, ok := in["rule"].(string); ok {
		obj.Rule = policyv1.SELinuxStrategy(v)
	}

	if v, ok := in["se_linux_options"].([]interface{}); ok && len(v) > 0 {
		obj.SELinuxOptions = expandPodSecurityPolicySELinuxOptions(v)
	}

	return obj
}
