package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigAuthentication(in *managementClient.AuthnConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.SANs) > 0 {
		obj["sans"] = toArrayInterface(in.SANs)
	}

	if len(in.Strategy) > 0 {
		obj["strategy"] = in.Strategy
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigAuthentication(p []interface{}) (*managementClient.AuthnConfig, error) {
	obj := &managementClient.AuthnConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["sans"].([]interface{}); ok && len(v) > 0 {
		obj.SANs = toArrayString(v)
	}

	if v, ok := in["strategy"].(string); ok && len(v) > 0 {
		obj.Strategy = v
	}

	return obj, nil
}
