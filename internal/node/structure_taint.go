package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenTaints(p []managementClient.Taint) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}

		if len(in.Value) > 0 {
			obj["value"] = in.Value
		}

		if len(in.Effect) > 0 {
			obj["effect"] = in.Effect
		}

		if len(in.TimeAdded) > 0 {
			obj["time_added"] = in.TimeAdded
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandTaints(p []interface{}) []managementClient.Taint {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Taint{}
	}

	obj := make([]managementClient.Taint, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj[i].Key = v
		}

		if v, ok := in["value"].(string); ok && len(v) > 0 {
			obj[i].Value = v
		}

		if v, ok := in["effect"].(string); ok && len(v) > 0 {
			obj[i].Effect = v
		}

		if v, ok := in["time_added"].(string); ok && len(v) > 0 {
			obj[i].TimeAdded = v
		}
	}

	return obj
}
