package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedHostPaths(in []policyv1.AllowedHostPath) []interface{} {

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

func expandPodSecurityPolicyAllowedHostPaths(in []interface{}) []policyv1.AllowedHostPath {

	obj := make([]policyv1.AllowedHostPath, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			hp := policyv1.AllowedHostPath{
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
