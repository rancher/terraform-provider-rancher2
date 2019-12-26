package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicySELinuxOptions(in *v1.SELinuxOptions) []interface{} {
	
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

func expandPodSecurityPolicySELinuxOptions(in []interface{}) *v1.SELinuxOptions {

	obj := &v1.SELinuxOptions{}

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