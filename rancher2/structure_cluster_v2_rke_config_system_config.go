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
			values, _ := interfaceToGhodssyaml(in.Config.Data)
			obj["config"] = values
		}
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigMachineSelectorFiles(p []rkev1.RKEProvisioningFiles) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}

		if in.MachineLabelSelector != nil {
			obj["machine_label_selector"] = flattenClusterV2RKEConfigSystemConfigLabelSelector(in.MachineLabelSelector)
		}
		if len(in.FileSources) > 0 {
			obj["file_sources"] = flattenClusterV2RKEConfigFileSources(in.FileSources)
		}
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigFileSources(p []rkev1.ProvisioningFileSource) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}
		if results := flattenClusterV2RKEConfigK8sObjectFileSource(in.Secret); len(results) > 0 {
			obj["secret"] = results
		}
		if results := flattenClusterV2RKEConfigK8sObjectFileSource(in.ConfigMap); len(results) > 0 {
			obj["configmap"] = results
		}
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigK8sObjectFileSource(p rkev1.K8sObjectFileSource) []interface{} {
	if p.Name == "" {
		return []interface{}{}
	}

	obj := make(map[string]interface{})
	obj["name"] = p.Name
	obj["default_permissions"] = p.DefaultPermissions
	if len(p.Items) > 0 {
		obj["items"] = flattenClusterV2RKEConfigKeyToPath(p.Items)
	}

	return []interface{}{obj}
}

func flattenClusterV2RKEConfigKeyToPath(p []rkev1.KeyToPath) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}
		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}
		if len(in.Path) > 0 {
			obj["path"] = in.Path
		}
		obj["dynamic"] = in.Dynamic
		if len(in.Permissions) > 0 {
			obj["permissions"] = in.Permissions
		}
		if len(in.Hash) > 0 {
			obj["hash"] = in.Hash
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
		if v, ok := in["config"].(string); ok && len(v) > 0 {
			values, _ := ghodssyamlToMapInterface(v)
			obj.Config.Data = values
		}
		out[i] = obj
	}

	return out
}

func expandClusterV2RKEConfigProvisioningFiles(p []interface{}) []rkev1.RKEProvisioningFiles {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]rkev1.RKEProvisioningFiles, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.RKEProvisioningFiles{}

		if v, ok := in["machine_label_selector"].([]interface{}); ok && len(v) > 0 {
			obj.MachineLabelSelector = expandClusterV2RKEConfigSystemConfigLabelSelector(v)
		}
		if v, ok := in["file_sources"].([]interface{}); ok && len(v) > 0 {
			obj.FileSources = expandClusterV2RKEConfigFileSources(v)
		}
		out[i] = obj
	}

	return out
}

func expandClusterV2RKEConfigFileSources(p []interface{}) []rkev1.ProvisioningFileSource {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]rkev1.ProvisioningFileSource, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.ProvisioningFileSource{}

		if v, ok := in["secret"].([]interface{}); ok && len(v) > 0 {
			obj.Secret = expandClusterV2RKEConfigK8sObjectFileSource(v)
		}
		if v, ok := in["configmap"].([]interface{}); ok && len(v) > 0 {
			obj.ConfigMap = expandClusterV2RKEConfigK8sObjectFileSource(v)
		}
		out[i] = obj
	}

	return out
}

func expandClusterV2RKEConfigK8sObjectFileSource(p []interface{}) rkev1.K8sObjectFileSource {
	obj := rkev1.K8sObjectFileSource{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["default_permissions"].(string); ok && len(v) > 0 {
		obj.DefaultPermissions = v
	}
	if v, ok := in["items"].([]interface{}); ok && len(v) > 0 {
		obj.Items = expandClusterV2RKEConfigKeyToPath(v)
	}

	return obj
}

func expandClusterV2RKEConfigKeyToPath(p []interface{}) []rkev1.KeyToPath {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]rkev1.KeyToPath, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.KeyToPath{}

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj.Key = v
		}
		if v, ok := in["path"].(string); ok && len(v) > 0 {
			obj.Path = v
		}
		if v, ok := in["dynamic"].(bool); ok {
			obj.Dynamic = v
		}
		if v, ok := in["permissions"].(string); ok && len(v) > 0 {
			obj.Permissions = v
		}
		if v, ok := in["hash"].(string); ok && len(v) > 0 {
			obj.Hash = v
		}
		out[i] = obj
	}

	return out
}
