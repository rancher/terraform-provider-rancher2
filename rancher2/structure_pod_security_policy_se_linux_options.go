package rancher2

import (
    v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenPodSecurityPolicySELinuxOptions(in *v1.SELinuxOptions) []interface{} {
	
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

func expandPodSecurityPolicySELinuxOptions(selinuxOptions []interface{}) *v1.SELinuxOptions {

	obj := &v1.SELinuxOptions{}

	if len(selinuxOptions) == 0 || selinuxOptions[0] == nil {
		return obj
	}

	in := selinuxOptions[0].(map[string]interface{})

	if v, ok := in["user"].(string); ok {
		obj.User = v
	}

	if v, ok := in["role"].(string); ok {
		obj.Role = v
	}

    if v, ok := in["type"].(string); ok {
		obj.Type = v
	}

    if v, ok := in["level"].(string); ok {
		obj.Level = v
	}

	return obj
}