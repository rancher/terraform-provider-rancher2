package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterGKEConfig(in *managementClient.GoogleKubernetesEngineConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterIpv4Cidr) > 0 {
		obj["cluster_ipv4_cidr"] = in.ClusterIpv4Cidr
	}

	if len(in.Credential) > 0 {
		obj["credential"] = in.Credential
	}

	if len(in.Description) > 0 {
		obj["description"] = in.Description
	}

	if in.DiskSizeGb > 0 {
		obj["disk_size_gb"] = int(in.DiskSizeGb)
	}

	obj["enable_alpha_feature"] = in.EnableAlphaFeature

	if in.EnableHTTPLoadBalancing != nil {
		obj["enable_http_load_balancing"] = *in.EnableHTTPLoadBalancing
	}

	if in.EnableHorizontalPodAutoscaling != nil {
		obj["enable_horizontal_pod_autoscaling"] = *in.EnableHorizontalPodAutoscaling
	}

	obj["enable_kubernetes_dashboard"] = in.EnableKubernetesDashboard

	obj["enable_legacy_abac"] = in.EnableLegacyAbac

	if in.EnableNetworkPolicyConfig != nil {
		obj["enable_network_policy_config"] = *in.EnableNetworkPolicyConfig
	}

	if in.EnableStackdriverLogging != nil {
		obj["enable_stackdriver_logging"] = *in.EnableStackdriverLogging
	}

	if in.EnableHorizontalPodAutoscaling != nil {
		obj["enable_stackdriver_monitoring"] = *in.EnableStackdriverMonitoring
	}

	if len(in.ImageType) > 0 {
		obj["image_type"] = in.ImageType
	}

	if len(in.Labels) > 0 {
		obj["labels"] = toMapInterface(in.Labels)
	}

	if len(in.Locations) > 0 {
		obj["locations"] = toArrayInterface(in.Locations)
	}

	if len(in.MachineType) > 0 {
		obj["machine_type"] = in.MachineType
	}

	if len(in.MaintenanceWindow) > 0 {
		obj["maintenance_window"] = in.MaintenanceWindow
	}

	if len(in.MasterVersion) > 0 {
		obj["master_version"] = in.MasterVersion
	}

	if len(in.Network) > 0 {
		obj["network"] = in.Network
	}

	if in.NodeCount > 0 {
		obj["node_count"] = int(in.NodeCount)
	}

	if len(in.NodeVersion) > 0 {
		obj["node_version"] = in.NodeVersion
	}

	if len(in.ProjectID) > 0 {
		obj["project_id"] = in.ProjectID
	}

	if len(in.SubNetwork) > 0 {
		obj["sub_network"] = in.SubNetwork
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterGKEConfig(p []interface{}) (*managementClient.GoogleKubernetesEngineConfig, error) {
	obj := &managementClient.GoogleKubernetesEngineConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_ipv4_cidr"].(string); ok && len(v) > 0 {
		obj.ClusterIpv4Cidr = v
	}

	if v, ok := in["credential"].(string); ok && len(v) > 0 {
		obj.Credential = v
	}

	if v, ok := in["description"].(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in["disk_size_gb"].(int); ok && v > 0 {
		obj.DiskSizeGb = int64(v)
	}

	if v, ok := in["enable_alpha_feature"].(bool); ok {
		obj.EnableAlphaFeature = v
	}

	if v, ok := in["enable_http_load_balancing"].(bool); ok {
		obj.EnableHTTPLoadBalancing = &v
	}

	if v, ok := in["enable_horizontal_pod_autoscaling"].(bool); ok {
		obj.EnableHorizontalPodAutoscaling = &v
	}

	if v, ok := in["enable_kubernetes_dashboard"].(bool); ok {
		obj.EnableKubernetesDashboard = v
	}

	if v, ok := in["enable_legacy_abac"].(bool); ok {
		obj.EnableLegacyAbac = v
	}

	if v, ok := in["enable_network_policy_config"].(bool); ok {
		obj.EnableNetworkPolicyConfig = &v
	}

	if v, ok := in["enable_stackdriver_logging"].(bool); ok {
		obj.EnableStackdriverLogging = &v
	}

	if v, ok := in["enable_stackdriver_monitoring"].(bool); ok {
		obj.EnableStackdriverMonitoring = &v
	}

	if v, ok := in["image_type"].(string); ok && len(v) > 0 {
		obj.ImageType = v
	}

	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in["locations"].([]interface{}); ok && len(v) > 0 {
		obj.Locations = toArrayString(v)
	}

	if v, ok := in["machine_type"].(string); ok && len(v) > 0 {
		obj.MachineType = v
	}

	if v, ok := in["maintenance_window"].(string); ok && len(v) > 0 {
		obj.MaintenanceWindow = v
	}

	if v, ok := in["master_version"].(string); ok && len(v) > 0 {
		obj.MasterVersion = v
	}

	if v, ok := in["network"].(string); ok && len(v) > 0 {
		obj.Network = v
	}

	if v, ok := in["node_count"].(int); ok && v > 0 {
		obj.NodeCount = int64(v)
	}

	if v, ok := in["node_version"].(string); ok && len(v) > 0 {
		obj.NodeVersion = v
	}

	if v, ok := in["project_id"].(string); ok && len(v) > 0 {
		obj.ProjectID = v
	}

	if v, ok := in["sub_network"].(string); ok && len(v) > 0 {
		obj.SubNetwork = v
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj, nil
}
