package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedFlexVolumes(in []policyv1.AllowedFlexVolume) []interface{} {

	out := make([]interface{}, len(in))

	for i, v := range in {
        obj := make(map[string]interface{})

		obj["driver"] = string(v.Driver)
		
		out[i] = obj
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedFlexVolumes(in []interface{}) []policyv1.AllowedFlexVolume {

	obj := make([]policyv1.AllowedFlexVolume, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = policyv1.AllowedFlexVolume{
				Driver: m["driver"].(string),
			}
		}
	}

	return obj
}