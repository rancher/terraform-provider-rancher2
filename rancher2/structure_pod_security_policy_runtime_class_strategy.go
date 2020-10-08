package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyRuntimeClassStrategy(in *managementClient.RuntimeClassStrategyOptions) []interface{} {

	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.AllowedRuntimeClassNames) > 0 {
		obj["allowed_runtime_class_names"] = toArrayInterface(in.AllowedRuntimeClassNames)
	}
	if len(in.DefaultRuntimeClassName) > 0 {
		obj["default_runtime_class_name"] = in.DefaultRuntimeClassName
	}

	return []interface{}{obj}
}

// Expanders

func expandPodSecurityPolicyRuntimeClassStrategy(in []interface{}) *managementClient.RuntimeClassStrategyOptions {

	obj := &managementClient.RuntimeClassStrategyOptions{}

	if len(in) == 0 || in[0] == nil {
		return obj
	}

	m := in[0].(map[string]interface{})

	if v, ok := m["allowed_runtime_class_names"].([]interface{}); ok {
		obj.AllowedRuntimeClassNames = toArrayString(v)
	}

	if v, ok := m["default_runtime_class_name"].(string); ok {
		obj.DefaultRuntimeClassName = v
	}

	return obj
}
