package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyIDRanges(in []managementClient.IDRange) []interface{} {

	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))

	for i, v := range in {
		out[i] = map[string]interface{}{
			"min": int(v.Min),
			"max": int(v.Max),
		}
	}

	return out
}

// Expanders

func expandPodSecurityPolicyIDRanges(in []interface{}) []managementClient.IDRange {

	if len(in) == 0 || in[0] == nil {
		return []managementClient.IDRange{}
	}

	obj := make([]managementClient.IDRange, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = managementClient.IDRange{
				Min: int64(m["min"].(int)),
				Max: int64(m["max"].(int)),
			}
		}
	}

	return obj
}
