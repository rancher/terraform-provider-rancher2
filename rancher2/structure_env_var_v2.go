package rancher2

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

// Flatteners

func flattenEnvVarsV2(p []rkev1.EnvVar) []interface{} {
	if p == nil || len(p) == 0 {
		return nil
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.Name) > 0 {
			obj["name"] = in.Name
		}

		if len(in.Value) > 0 {
			obj["value"] = in.Value
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandEnvVarsV2(p []interface{}) []rkev1.EnvVar {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := make([]rkev1.EnvVar, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["name"].(string); ok && len(v) > 0 {
			obj[i].Name = v
		}

		if v, ok := in["value"].(string); ok && len(v) > 0 {
			obj[i].Value = v
		}
	}

	return obj
}
