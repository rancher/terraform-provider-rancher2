package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyAllowedHostPaths(in []managementClient.AllowedHostPath) []interface{} {

	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))

	for i, v := range in {
		obj := make(map[string]interface{})

		obj["path_prefix"] = v.PathPrefix
		obj["read_only"] = v.ReadOnly

		out[i] = obj
	}

	return out
}

// Expanders

func expandPodSecurityPolicyAllowedHostPaths(in []interface{}) []managementClient.AllowedHostPath {

	if len(in) == 0 || in[0] == nil {
		return []managementClient.AllowedHostPath{}
	}

	obj := make([]managementClient.AllowedHostPath, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			hp := managementClient.AllowedHostPath{
				PathPrefix: m["path_prefix"].(string),
			}

			if ro, ok := m["read_only"].(bool); ok {
				hp.ReadOnly = ro
			}

			obj[i] = hp
		}
	}

	return obj
}
