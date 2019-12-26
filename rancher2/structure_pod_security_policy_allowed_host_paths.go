package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedHostPaths(hp []policyv1.AllowedHostPath) []interface{} {
	if len(hp) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(hp))
	for i, in := range hp {
        obj := make(map[string]interface{})

        obj["path_prefix"] = in.PathPrefix
		obj["read_only"] = in.ReadOnly

		out[i] = obj
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedHostPaths(hp []interface{}) []policyv1.AllowedHostPath {

	if len(hp) == 0 || hp[0] == nil {
		return []policyv1.AllowedHostPath{}
	}

	obj := make([]policyv1.AllowedHostPath, len(hp))

	for i := range hp {
		in := hp[i].(map[string]interface{})

		obj[i].PathPrefix = in["path_prefix"].(string)

        if ro, ok := in["read_only"].(bool); ok {
				obj[i].ReadOnly = ro
		}

	}

	return obj
}
