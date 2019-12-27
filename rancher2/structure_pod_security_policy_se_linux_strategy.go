package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenPodSecurityPolicySELinuxStrategy(in *managementClient.SELinuxStrategyOptions) []interface{} {

	obj := make(map[string]interface{})

	if len(in.Rule) > 0 {
		obj["rule"] = in.Rule
	}
	if in.SELinuxOptions != nil {
		obj["se_linux_options"] = flattenPodSecurityPolicySELinuxOptions(in.SELinuxOptions)
	}

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicySELinuxStrategy(in []interface{}) *managementClient.SELinuxStrategyOptions {

	obj := &managementClient.SELinuxStrategyOptions{}

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = v
	}

	if v, ok := m["se_linux_options"].([]interface{}); ok && len(v) > 0 {
		obj.SELinuxOptions = expandPodSecurityPolicySELinuxOptions(v)
	}

	return obj
}
