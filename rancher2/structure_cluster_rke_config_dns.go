package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigDNSNodelocal(in *managementClient.Nodelocal) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return nil
	}

	if len(in.IPAddress) > 0 {
		obj["ip_address"] = in.IPAddress
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigDNSLinearAutoscalerParams(in *managementClient.LinearAutoscalerParams) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return nil
	}

	if in.CoresPerReplica > 0 {
		obj["cores_per_replica"] = in.CoresPerReplica
	}

	if in.NodesPerReplica > 0 {
		obj["nodes_per_replica"] = in.NodesPerReplica
	}

	if in.Max >= 0 {
		obj["max"] = int(in.Max)
	}

	if in.Min > 0 {
		obj["min"] = int(in.Min)
	}

	obj["prevent_single_point_failure"] = in.PreventSinglePointFailure

	return []interface{}{obj}
}

func flattenClusterRKEConfigDNS(in *managementClient.DNSConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.NodeSelector != nil && len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	if in.Nodelocal != nil {
		obj["nodelocal"] = flattenClusterRKEConfigDNSNodelocal(in.Nodelocal)
	}

	if in.LinearAutoscalerParams != nil {
		obj["linear_autoscaler_params"] = flattenClusterRKEConfigDNSLinearAutoscalerParams(in.LinearAutoscalerParams)
	}

	if in.Options != nil && len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	if len(in.ReverseCIDRs) > 0 {
		obj["reverse_cidrs"] = toArrayInterface(in.ReverseCIDRs)
	}

	if in.Tolerations != nil && len(in.Tolerations) > 0 {
		obj["tolerations"] = flattenTolerations(in.Tolerations)
	}

	if len(in.UpstreamNameservers) > 0 {
		obj["upstream_nameservers"] = toArrayInterface(in.UpstreamNameservers)
	}

	if in.UpdateStrategy != nil {
		obj["update_strategy"] = flattenDeploymentStrategy(in.UpdateStrategy)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigDNSNodelocal(p []interface{}) *managementClient.Nodelocal {
	obj := &managementClient.Nodelocal{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["ip_address"].(string); ok && len(v) > 0 {
		obj.IPAddress = v
	}

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	return obj
}

func expandClusterRKEConfigDNSLinearAutoscalerParams(p []interface{}) *managementClient.LinearAutoscalerParams {
	obj := &managementClient.LinearAutoscalerParams{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cores_per_replica"].(float64); ok && v > 0 {
		obj.CoresPerReplica = v
	}

	if v, ok := in["nodes_per_replica"].(float64); ok && v > 0 {
		obj.NodesPerReplica = v
	}

	if v, ok := in["max"].(int); ok && v >= 0 {
		obj.Max = int64(v)
	}

	if v, ok := in["min"].(int); ok && v > 0 {
		obj.Min = int64(v)
	}

	if v, ok := in["prevent_single_point_failure"].(bool); ok {
		obj.PreventSinglePointFailure = v
	}

	return obj
}

func expandClusterRKEConfigDNS(p []interface{}) (*managementClient.DNSConfig, error) {
	obj := &managementClient.DNSConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["nodelocal"].([]interface{}); ok && len(v) > 0 {
		obj.Nodelocal = expandClusterRKEConfigDNSNodelocal(v)
	}

	if v, ok := in["linear_autoscaler_params"].([]interface{}); ok && len(v) > 0 {
		obj.LinearAutoscalerParams = expandClusterRKEConfigDNSLinearAutoscalerParams(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	if v, ok := in["reverse_cidrs"].([]interface{}); ok && len(v) > 0 {
		obj.ReverseCIDRs = toArrayString(v)
	}

	if v, ok := in["tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.Tolerations = expandTolerations(v)
	}

	if v, ok := in["upstream_nameservers"].([]interface{}); ok && len(v) > 0 {
		obj.UpstreamNameservers = toArrayString(v)
	}

	if v, ok := in["update_strategy"].([]interface{}); ok && v != nil {
		obj.UpdateStrategy = expandDeploymentStrategy(v)
	}

	return obj, nil
}
