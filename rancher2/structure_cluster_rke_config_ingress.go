package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigIngress(in *managementClient.IngressConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.DefaultBackend != nil {
		obj["default_backend"] = *in.DefaultBackend
	}

	if len(in.DNSPolicy) > 0 {
		obj["dns_policy"] = in.DNSPolicy
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if in.HTTPPort > 0 {
		obj["http_port"] = int(in.HTTPPort)
	}

	if in.HTTPSPort > 0 {
		obj["https_port"] = int(in.HTTPSPort)
	}

	if len(in.NetworkMode) > 0 {
		obj["network_mode"] = in.NetworkMode
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	if in.Tolerations != nil && len(in.Tolerations) > 0 {
		obj["tolerations"] = flattenTolerations(in.Tolerations)
	}

	if in.UpdateStrategy != nil {
		obj["update_strategy"] = flattenDaemonSetStrategy(in.UpdateStrategy)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigIngress(p []interface{}) (*managementClient.IngressConfig, error) {
	obj := &managementClient.IngressConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["default_backend"].(bool); ok {
		obj.DefaultBackend = &v
	}

	if v, ok := in["dns_policy"].(string); ok && len(v) > 0 {
		obj.DNSPolicy = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["http_port"].(int); ok {
		obj.HTTPPort = int64(v)
	}

	if v, ok := in["https_port"].(int); ok {
		obj.HTTPSPort = int64(v)
	}

	if v, ok := in["network_mode"].(string); ok && len(v) > 0 {
		obj.NetworkMode = v
	}

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	if v, ok := in["tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.Tolerations = expandTolerations(v)
	}

	if v, ok := in["update_strategy"].([]interface{}); ok && v != nil {
		obj.UpdateStrategy = expandDaemonSetStrategy(v)
	}

	return obj, nil
}
