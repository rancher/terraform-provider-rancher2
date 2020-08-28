package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicySELinuxStrategy(in *managementClient.SELinuxStrategyOptions) []interface{} {

	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.Rule) > 0 {
		obj["rule"] = in.Rule
	}
	if in.SELinuxOptions != nil {
		obj["se_linux_option"] = flattenPodSecurityPolicySELinuxOptions(in.SELinuxOptions)
	}

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicySELinuxStrategy(in []interface{}) *managementClient.SELinuxStrategyOptions {

	obj := &managementClient.SELinuxStrategyOptions{}

	if len(in) == 0 || in[0] == nil {
		return obj
	}

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = v
	}

	if v, ok := m["se_linux_option"].([]interface{}); ok && len(v) > 0 {
		obj.SELinuxOptions = expandPodSecurityPolicySELinuxOptions(v)
	}

	return obj
}
