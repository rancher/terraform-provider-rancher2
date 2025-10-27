package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenTolerations(p []managementClient.Toleration) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}

		if len(in.Effect) > 0 {
			obj["effect"] = in.Effect
		}

		if len(in.Operator) > 0 {
			obj["operator"] = in.Operator
		}

		if in.TolerationSeconds != nil {
			obj["seconds"] = int(*in.TolerationSeconds)
		}

		if len(in.Value) > 0 {
			obj["value"] = in.Value
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandTolerations(p []interface{}) []managementClient.Toleration {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.Toleration{}
	}

	obj := make([]managementClient.Toleration, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj[i].Key = v
		}

		if v, ok := in["effect"].(string); ok && len(v) > 0 {
			obj[i].Effect = v
		}

		if v, ok := in["operator"].(string); ok && len(v) > 0 {
			obj[i].Operator = v
		}

		if v, ok := in["seconds"].(int); ok && v > 0 {
			sec := int64(v)
			obj[i].TolerationSeconds = &sec
		}

		if v, ok := in["value"].(string); ok && len(v) > 0 {
			obj[i].Value = v
		}
	}

	return obj
}
