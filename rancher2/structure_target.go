package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenTargets(p []managementClient.Target) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.ProjectID) > 0 {
			obj["project_id"] = in.ProjectID
		}

		if len(in.AppID) > 0 {
			obj["app_id"] = in.AppID
		}

		if len(in.Healthstate) > 0 {
			obj["health_state"] = in.Healthstate
		}

		if len(in.State) > 0 {
			obj["state"] = in.State
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandTargets(p []interface{}) []managementClient.Target {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Target{}
	}

	obj := make([]managementClient.Target, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["project_id"].(string); ok && len(v) > 0 {
			obj[i].ProjectID = v
		}

		if v, ok := in["app_id"].(string); ok && len(v) > 0 {
			obj[i].AppID = v
		}

		if v, ok := in["health_state"].(string); ok && len(v) > 0 {
			obj[i].Healthstate = v
		}

		if v, ok := in["state"].(string); ok && len(v) > 0 {
			obj[i].State = v
		}
	}

	return obj
}
