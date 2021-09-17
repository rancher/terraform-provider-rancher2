package rancher2

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Flatteners

func flattenClusterV2RKEConfigSystemConfigLabelSelectorExpression(p []metav1.LabelSelectorRequirement) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}

		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}
		if len(in.Operator) > 0 {
			obj["operator"] = string(in.Operator)
		}
		if len(in.Values) > 0 {
			obj["values"] = toArrayInterfaceSorted(in.Values)
		}
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigSystemConfigLabelSelector(in *metav1.LabelSelector) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.MatchLabels) > 0 {
		obj["match_labels"] = toMapInterface(in.MatchLabels)
	}
	if len(in.MatchExpressions) > 0 {
		obj["match_expressions"] = flattenClusterV2RKEConfigSystemConfigLabelSelectorExpression(in.MatchExpressions)
	}

	return []interface{}{obj}
}

func flattenClusterV2RKEConfigSystemConfig(p []rkev1.RKESystemConfig) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}

		if in.MachineLabelSelector != nil {
			obj["machine_label_selector"] = flattenClusterV2RKEConfigSystemConfigLabelSelector(in.MachineLabelSelector)
		}
		if len(in.Config.Data) > 0 {
			obj["config"] = in.Config.Data
		}
		out[i] = obj
	}

	return out
}

// Expanders

func expandClusterV2RKEConfigSystemConfigLabelSelectorExpression(p []interface{}) []metav1.LabelSelectorRequirement {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]metav1.LabelSelectorRequirement, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := metav1.LabelSelectorRequirement{}

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj.Key = v
		}
		if v, ok := in["operator"].(string); ok && len(v) > 0 {
			obj.Operator = metav1.LabelSelectorOperator(v)
		}
		if v, ok := in["values"].([]interface{}); ok && len(v) > 0 {
			obj.Values = toArrayStringSorted(v)
		}
		out[i] = obj
	}

	return out
}

func expandClusterV2RKEConfigSystemConfigLabelSelector(p []interface{}) *metav1.LabelSelector {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &metav1.LabelSelector{}

	in := p[0].(map[string]interface{})

	if v, ok := in["match_labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.MatchLabels = toMapString(v)
	}
	if v, ok := in["match_expressions"].([]interface{}); ok && len(v) > 0 {
		obj.MatchExpressions = expandClusterV2RKEConfigSystemConfigLabelSelectorExpression(v)
	}

	return obj
}

func expandClusterV2RKEConfigSystemConfig(p []interface{}) []rkev1.RKESystemConfig {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]rkev1.RKESystemConfig, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.RKESystemConfig{}

		if v, ok := in["machine_label_selector"].([]interface{}); ok && len(v) > 0 {
			obj.MachineLabelSelector = expandClusterV2RKEConfigSystemConfigLabelSelector(v)
		}
		if v, ok := in["config"].(map[string]interface{}); ok && len(v) > 0 {
			obj.Config.Data = v
		}
		out[i] = obj
	}

	return out
}
