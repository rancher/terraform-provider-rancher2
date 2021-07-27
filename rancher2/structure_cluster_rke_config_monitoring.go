package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigMonitoring(in *managementClient.MonitoringConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
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

	if in.Replicas != nil {
		obj["replicas"] = int(*in.Replicas)
	}

	if in.Tolerations != nil && len(in.Tolerations) > 0 {
		obj["tolerations"] = flattenTolerations(in.Tolerations)
	}

	if in.UpdateStrategy != nil {
		obj["update_strategy"] = flattenDeploymentStrategy(in.UpdateStrategy)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigMonitoring(p []interface{}) (*managementClient.MonitoringConfig, error) {
	obj := &managementClient.MonitoringConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}
	if obj.Provider != clusterRKEConfigMonitoringProviderDisabled {
		// Setting Replicas to 1 if monitoring enabled
		value := int64(1)
		obj.Replicas = &value
	}

	if v, ok := in["replicas"].(int); ok && v > 0 {
		value := int64(v)
		obj.Replicas = &value
	}

	if v, ok := in["tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.Tolerations = expandTolerations(v)
	}

	if v, ok := in["update_strategy"].([]interface{}); ok && v != nil {
		obj.UpdateStrategy = expandDeploymentStrategy(v)
	}

	return obj, nil
}
