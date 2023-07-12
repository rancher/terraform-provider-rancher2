package rancher2

import (
	corev1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenTolerationsV2(p []corev1.Toleration) []interface{} {
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
			obj["effect"] = string(in.Effect)
		}

		if len(in.Operator) > 0 {
			obj["operator"] = string(in.Operator)
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

func expandTolerationsV2(p []interface{}) []corev1.Toleration {
	if len(p) == 0 || p[0] == nil {
		return []corev1.Toleration{}
	}

	obj := make([]corev1.Toleration, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj[i].Key = v
		}

		if v, ok := in["effect"].(string); ok && len(v) > 0 {
			obj[i].Effect = corev1.TaintEffect(v)
		}

		if v, ok := in["operator"].(string); ok && len(v) > 0 {
			obj[i].Operator = corev1.TolerationOperator(v)
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
