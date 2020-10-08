package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyAllowedCSIDrivers(in []managementClient.AllowedCSIDriver) []interface{} {

	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))

	for i, v := range in {
		obj := make(map[string]interface{})

		obj["name"] = v.Name

		out[i] = obj
	}

	return out
}

// Expanders

func expandPodSecurityPolicyAllowedCSIDrivers(in []interface{}) []managementClient.AllowedCSIDriver {

	if len(in) == 0 || in[0] == nil {
		return []managementClient.AllowedCSIDriver{}
	}

	obj := make([]managementClient.AllowedCSIDriver, len(in))

	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			obj[i] = managementClient.AllowedCSIDriver{
				Name: m["name"].(string),
			}
		}
	}

	return obj
}
