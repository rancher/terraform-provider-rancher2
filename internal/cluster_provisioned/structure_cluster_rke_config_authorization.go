package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigAuthorization(in *managementClient.AuthzConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Mode) > 0 {
		obj["mode"] = in.Mode
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigAuthorization(p []interface{}) (*managementClient.AuthzConfig, error) {
	obj := &managementClient.AuthzConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["mode"].(string); ok && len(v) > 0 {
		obj.Mode = v
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	return obj, nil
}
