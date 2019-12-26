package rancher2

import (
    managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyAllowedFlexVolumes(in []managementClient.AllowedFlexVolume) []interface{} {

	out := make([]interface{}, len(in))

	for i, v := range in {
        obj := make(map[string]interface{})

		obj["driver"] = v.Driver
		
		out[i] = obj
	}
	
	return out
}

// Expanders

func expandPodSecurityPolicyAllowedFlexVolumes(in []interface{}) []managementClient.AllowedFlexVolume {

	obj := make([]managementClient.AllowedFlexVolume, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = managementClient.AllowedFlexVolume{
				Driver: m["driver"].(string),
			}
		}
	}

	return obj
}