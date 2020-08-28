package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenMembers(p []managementClient.Member) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.AccessType) > 0 {
			obj["access_type"] = in.AccessType
		}

		if len(in.GroupPrincipalID) > 0 {
			obj["group_principal_id"] = in.GroupPrincipalID
		}

		if len(in.UserPrincipalID) > 0 {
			obj["user_principal_id"] = in.UserPrincipalID
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandMembers(p []interface{}) []managementClient.Member {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Member{}
	}

	obj := make([]managementClient.Member, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["access_type"].(string); ok && len(v) > 0 {
			obj[i].AccessType = v
		}

		if v, ok := in["group_principal_id"].(string); ok && len(v) > 0 {
			obj[i].GroupPrincipalID = v
		}

		if v, ok := in["user_principal_id"].(string); ok && len(v) > 0 {
			obj[i].UserPrincipalID = v
		}
	}

	return obj
}
