package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenAnswers(p []managementClient.Answer) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.ClusterID) > 0 {
			obj["cluster_id"] = in.ClusterID
		}

		if len(in.ProjectID) > 0 {
			obj["project_id"] = in.ProjectID
		}

		if len(in.Values) > 0 {
			obj["values"] = toMapInterface(in.Values)
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandAnswers(p []interface{}) []managementClient.Answer {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Answer{}
	}

	obj := make([]managementClient.Answer, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["cluster_id"].(string); ok && len(v) > 0 {
			obj[i].ClusterID = v
		}

		if v, ok := in["project_id"].(string); ok && len(v) > 0 {
			obj[i].ProjectID = v
		}

		if v, ok := in["values"].(map[string]interface{}); ok && len(v) > 0 {
			obj[i].Values = toMapString(v)
		}
	}

	return obj
}
