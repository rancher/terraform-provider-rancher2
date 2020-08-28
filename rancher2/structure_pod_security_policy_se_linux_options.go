package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicySELinuxOptions(in *managementClient.SELinuxOptions) []interface{} {

	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.User) > 0 {
		obj["user"] = in.User
	}

	if len(in.Role) > 0 {
		obj["role"] = in.Role
	}

	if len(in.Type) > 0 {
		obj["type"] = in.Type
	}

	if len(in.Level) > 0 {
		obj["level"] = in.Level
	}

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicySELinuxOptions(in []interface{}) *managementClient.SELinuxOptions {

	obj := &managementClient.SELinuxOptions{}

	if len(in) == 0 || in[0] == nil {
		return obj
	}

	m := in[0].(map[string]interface{})

	if v, ok := m["user"].(string); ok {
		obj.User = v
	}

	if v, ok := m["role"].(string); ok {
		obj.Role = v
	}

	if v, ok := m["type"].(string); ok {
		obj.Type = v
	}

	if v, ok := m["level"].(string); ok {
		obj.Level = v
	}

	return obj
}
