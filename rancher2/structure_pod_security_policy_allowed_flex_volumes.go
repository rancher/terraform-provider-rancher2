package rancher2

import (
    policyv1 "k8s.io/api/policy/v1beta1"
)

// Flatteners

func flattenPodSecurityPolicyAllowedFlexVolumes(fv []policyv1.AllowedFlexVolume) []interface{} {
	if len(fv) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(fv))
	for i, in := range fv {
        obj := make(map[string]interface{})

		obj["driver"] = string(in.Driver)
		
		out[i] = obj
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedFlexVolumes(fv []interface{}) []policyv1.AllowedFlexVolume {
	if len(fv) == 0 || fv[0] == nil {
		return []policyv1.AllowedFlexVolume{}
	}

	obj := make([]policyv1.AllowedFlexVolume, len(fv))

	for i := range fv {
		in := fv[i].(map[string]interface{})

		obj[i].Driver = in["driver"].(string)
	}

	return obj
}