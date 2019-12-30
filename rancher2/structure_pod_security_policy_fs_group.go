package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyFSGroup(in *managementClient.FSGroupStrategyOptions) []interface{} {

	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.Rule) > 0 {
		obj["rule"] = in.Rule
	}
	if len(in.Ranges) > 0 {
		obj["ranges"] = flattenPodSecurityPolicyIDRanges(in.Ranges)
	}

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicyFSGroup(in []interface{}) *managementClient.FSGroupStrategyOptions {

	obj := &managementClient.FSGroupStrategyOptions{}

	if len(in) == 0 || in[0] == nil {
		return obj
	}

	m := in[0].(map[string]interface{})

	if v, ok := m["rule"].(string); ok {
		obj.Rule = v
	}

	if v, ok := m["ranges"].([]interface{}); ok && len(v) > 0 {
		obj.Ranges = expandPodSecurityPolicyIDRanges(v)
	}

	return obj
}
