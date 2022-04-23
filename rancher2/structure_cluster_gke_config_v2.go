package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	defaultClusterGKEConfigVClusterIpv4CidrBlock = "10.0.0.0/16"
	defaultClusterGKEConfigV2LoggingService      = "logging.googleapis.com/kubernetes"
	defaultClusterGKEConfigV2MaintenanceWindow   = "00:00"
	defaultClusterGKEConfigV2MonitoringService   = "monitoring.googleapis.com/kubernetes"
)

// Flatteners

func flattenClusterGKEConfigV2ClusterAddons(in *managementClient.GKEClusterAddons) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	obj["http_load_balancing"] = in.HTTPLoadBalancing
	obj["horizontal_pod_autoscaling"] = in.HorizontalPodAutoscaling
	obj["network_policy_config"] = in.NetworkPolicyConfig

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2IPAllocationPolicy(in *managementClient.GKEIPAllocationPolicy) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	if len(in.ClusterIpv4CidrBlock) > 0 {
		obj["cluster_ipv4_cidr_block"] = in.ClusterIpv4CidrBlock
	}
	if len(in.ClusterSecondaryRangeName) > 0 {
		obj["cluster_secondary_range_name"] = in.ClusterSecondaryRangeName
	}
	obj["create_subnetwork"] = in.CreateSubnetwork
	if len(in.NodeIpv4CidrBlock) > 0 {
		obj["node_ipv4_cidr_block"] = in.NodeIpv4CidrBlock
	}
	if len(in.ServicesIpv4CidrBlock) > 0 {
		obj["services_ipv4_cidr_block"] = in.ServicesIpv4CidrBlock
	}
	if len(in.ServicesSecondaryRangeName) > 0 {
		obj["services_secondary_range_name"] = in.ServicesSecondaryRangeName
	}
	if len(in.SubnetworkName) > 0 {
		obj["subnetwork_name"] = in.SubnetworkName
	}
	obj["use_ip_aliases"] = in.UseIPAliases

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2CidrBlocks(input []managementClient.GKECidrBlock) []interface{} {
	if input == nil {
		return nil
	}
	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}
		if len(in.CidrBlock) > 0 {
			obj["cidr_block"] = in.CidrBlock
		}
		if len(in.DisplayName) > 0 {
			obj["display_name"] = in.DisplayName
		}
		out[i] = obj
	}

	return out
}

func flattenClusterGKEConfigV2MasterAuthorizedNetworksConfig(in *managementClient.GKEMasterAuthorizedNetworksConfig) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	if len(in.CidrBlocks) > 0 {
		obj["cidr_blocks"] = flattenClusterGKEConfigV2CidrBlocks(in.CidrBlocks)
	}
	obj["enabled"] = in.Enabled

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2NodePoolConfigAutoscaling(in *managementClient.GKENodePoolAutoscaling) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	obj["enabled"] = in.Enabled
	if in.MaxNodeCount > 0 {
		obj["max_node_count"] = int(in.MaxNodeCount)
	}
	if in.MinNodeCount > 0 {
		obj["min_node_count"] = int(in.MinNodeCount)
	}

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2NodeTaintsConfig(p []managementClient.GKENodeTaintConfig) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}
		if len(in.Value) > 0 {
			obj["value"] = in.Value
		}
		if len(in.Effect) > 0 {
			obj["effect"] = in.Effect
		}
		out[i] = obj
	}

	return out
}

func flattenClusterGKEConfigV2NodeConfig(in *managementClient.GKENodeConfig) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	if in.DiskSizeGb > 0 {
		obj["disk_size_gb"] = int(in.DiskSizeGb)
	}
	if len(in.DiskType) > 0 {
		obj["disk_type"] = in.DiskType
	}
	if len(in.ImageType) > 0 {
		obj["image_type"] = in.ImageType
	}
	if len(in.Labels) > 0 {
		obj["labels"] = toMapInterface(in.Labels)
	}
	if in.LocalSsdCount > 0 {
		obj["local_ssd_count"] = int(in.LocalSsdCount)
	}
	if len(in.MachineType) > 0 {
		obj["machine_type"] = in.MachineType
	}
	if len(in.OauthScopes) > 0 {
		obj["oauth_scopes"] = toArrayInterfaceSorted(in.OauthScopes)
	}
	if len(in.Tags) > 0 {
		obj["tags"] = toArrayInterfaceSorted(in.Tags)
	}
	obj["preemptible"] = in.Preemptible
	if len(in.Taints) > 0 {
		obj["taints"] = flattenClusterGKEConfigV2NodeTaintsConfig(in.Taints)
	}

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2NodePoolManagement(in *managementClient.GKENodePoolManagement) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	obj["auto_repair"] = in.AutoRepair
	obj["auto_upgrade"] = in.AutoUpgrade

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2NodePoolsConfig(input []managementClient.GKENodePoolConfig) []interface{} {
	if input == nil {
		return nil
	}
	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}

		if in.Autoscaling != nil {
			obj["autoscaling"] = flattenClusterGKEConfigV2NodePoolConfigAutoscaling(in.Autoscaling)
		}
		if in.Config != nil {
			obj["config"] = flattenClusterGKEConfigV2NodeConfig(in.Config)
		}
		if in.InitialNodeCount != nil {
			obj["initial_node_count"] = int(*in.InitialNodeCount)
		}
		if in.Management != nil {
			obj["management"] = flattenClusterGKEConfigV2NodePoolManagement(in.Management)
		}
		if in.MaxPodsConstraint != nil {
			obj["max_pods_constraint"] = int(*in.MaxPodsConstraint)
		}
		if in.Name != nil && len(*in.Name) > 0 {
			obj["name"] = *in.Name
		}
		if in.Version != nil && len(*in.Version) > 0 {
			obj["version"] = *in.Version
		}
		out[i] = obj
	}

	return out
}

func flattenClusterGKEConfigV2PrivateClusterConfig(in *managementClient.GKEPrivateClusterConfig) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}

	obj["enable_private_endpoint"] = in.EnablePrivateEndpoint
	obj["enable_private_nodes"] = in.EnablePrivateNodes
	if len(in.MasterIpv4CidrBlock) > 0 {
		obj["master_ipv4_cidr_block"] = in.MasterIpv4CidrBlock
	}

	return []interface{}{obj}
}

func flattenClusterGKEConfigV2(in *managementClient.GKEClusterConfigSpec, p []interface{}) []interface{} {
	if in == nil {
		return nil
	}

	obj := map[string]interface{}{}
	if len(p) != 0 && p[0] != nil {
		obj = p[0].(map[string]interface{})
	}

	if len(in.ClusterName) > 0 {
		obj["name"] = in.ClusterName
	}
	if len(in.GoogleCredentialSecret) > 0 {
		obj["google_credential_secret"] = in.GoogleCredentialSecret
	}
	if len(in.ProjectID) > 0 {
		obj["project_id"] = in.ProjectID
	}
	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}
	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}
	obj["imported"] = in.Imported
	// Returning here for imported clusters
	if in.Imported {
		return []interface{}{obj}
	}

	if in.ClusterAddons != nil {
		obj["cluster_addons"] = flattenClusterGKEConfigV2ClusterAddons(in.ClusterAddons)
	}
	if in.ClusterIpv4CidrBlock != nil && len(*in.ClusterIpv4CidrBlock) > 0 {
		obj["cluster_ipv4_cidr_block"] = *in.ClusterIpv4CidrBlock
	}
	if len(in.Description) > 0 {
		obj["description"] = in.Description
	}
	if in.EnableKubernetesAlpha != nil {
		obj["enable_kubernetes_alpha"] = *in.EnableKubernetesAlpha
	}
	if in.IPAllocationPolicy != nil {
		obj["ip_allocation_policy"] = flattenClusterGKEConfigV2IPAllocationPolicy(in.IPAllocationPolicy)
	}
	if in.KubernetesVersion != nil && len(*in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = *in.KubernetesVersion
	}
	if in.Labels != nil && len(in.Labels) > 0 {
		obj["labels"] = toMapInterface(in.Labels)
	}
	if in.Locations != nil && len(in.Locations) > 0 {
		obj["locations"] = toArrayInterface(in.Locations)
	}
	if in.LoggingService != nil && len(*in.LoggingService) > 0 {
		obj["logging_service"] = *in.LoggingService
	}
	if in.MaintenanceWindow != nil && len(*in.MaintenanceWindow) > 0 {
		obj["maintenance_window"] = *in.MaintenanceWindow
	}
	if in.MasterAuthorizedNetworksConfig != nil {
		obj["master_authorized_networks_config"] = flattenClusterGKEConfigV2MasterAuthorizedNetworksConfig(in.MasterAuthorizedNetworksConfig)
	}
	if in.MonitoringService != nil && len(*in.MonitoringService) > 0 {
		obj["monitoring_service"] = *in.MonitoringService
	}
	if in.Network != nil && len(*in.Network) > 0 {
		obj["network"] = *in.Network
	}
	if in.NetworkPolicyEnabled != nil {
		obj["network_policy_enabled"] = *in.NetworkPolicyEnabled
	}
	if in.NodePools != nil {
		obj["node_pools"] = flattenClusterGKEConfigV2NodePoolsConfig(in.NodePools)
	}
	if in.PrivateClusterConfig != nil {
		obj["private_cluster_config"] = flattenClusterGKEConfigV2PrivateClusterConfig(in.PrivateClusterConfig)
	}
	if in.Subnetwork != nil && len(*in.Subnetwork) > 0 {
		obj["subnetwork"] = *in.Subnetwork
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterGKEConfigV2ClusterAddons(p []interface{}) *managementClient.GKEClusterAddons {
	obj := &managementClient.GKEClusterAddons{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["http_load_balancing"].(bool); ok {
		obj.HTTPLoadBalancing = v
	}
	if v, ok := in["horizontal_pod_autoscaling"].(bool); ok {
		obj.HorizontalPodAutoscaling = v
	}
	if v, ok := in["network_policy_config"].(bool); ok {
		obj.NetworkPolicyConfig = v
	}

	return obj
}

func expandClusterGKEConfigV2IPAllocationPolicy(p []interface{}) *managementClient.GKEIPAllocationPolicy {
	obj := &managementClient.GKEIPAllocationPolicy{}
	if len(p) == 0 || p[0] == nil {
		return &managementClient.GKEIPAllocationPolicy{
			UseIPAliases: true,
		}
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.ClusterIpv4CidrBlock = v
	}
	if v, ok := in["cluster_secondary_range_name"].(string); ok && len(v) > 0 {
		obj.ClusterSecondaryRangeName = v
	}
	if v, ok := in["create_subnetwork"].(bool); ok {
		obj.CreateSubnetwork = v
	}
	if v, ok := in["node_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.NodeIpv4CidrBlock = v
	}
	if v, ok := in["services_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.ServicesIpv4CidrBlock = v
	}
	if v, ok := in["services_secondary_range_name"].(string); ok && len(v) > 0 {
		obj.ServicesSecondaryRangeName = v
	}
	if v, ok := in["subnetwork_name"].(string); ok && len(v) > 0 {
		obj.SubnetworkName = v
	}
	if v, ok := in["use_ip_aliases"].(bool); ok {
		obj.UseIPAliases = v
	}

	return obj
}

func expandClusterGKEConfigV2CidrBlocks(p []interface{}) []managementClient.GKECidrBlock {
	if p == nil || p[0] == nil {
		return []managementClient.GKECidrBlock{}
	}
	out := make([]managementClient.GKECidrBlock, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.GKECidrBlock{}

		if v, ok := in["cidr_block"].(string); ok {
			obj.CidrBlock = v
		}
		if v, ok := in["display_name"].(string); ok {
			obj.DisplayName = v
		}
		out[i] = obj
	}
	return out
}

func expandClusterGKEConfigV2MasterAuthorizedNetworksConfig(p []interface{}) *managementClient.GKEMasterAuthorizedNetworksConfig {
	obj := &managementClient.GKEMasterAuthorizedNetworksConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}
	if v, ok := in["cidr_blocks"].([]interface{}); ok && len(v) > 0 {
		obj.CidrBlocks = expandClusterGKEConfigV2CidrBlocks(v)
	}

	return obj
}

func expandClusterGKEConfigV2NodePoolConfigAutoscaling(p []interface{}) *managementClient.GKENodePoolAutoscaling {
	obj := &managementClient.GKENodePoolAutoscaling{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}
	if v, ok := in["max_node_count"].(int); ok {
		obj.MaxNodeCount = int64(v)
	}
	if v, ok := in["min_node_count"].(int); ok {
		obj.MinNodeCount = int64(v)
	}
	return obj
}

func expandClusterGKEConfigV2NodeTaintsConfig(p []interface{}) []managementClient.GKENodeTaintConfig {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.GKENodeTaintConfig{}
	}

	obj := make([]managementClient.GKENodeTaintConfig, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj[i].Key = v
		}
		if v, ok := in["value"].(string); ok && len(v) > 0 {
			obj[i].Value = v
		}
		if v, ok := in["effect"].(string); ok && len(v) > 0 {
			obj[i].Effect = v
		}
	}

	return obj
}

func expandClusterGKEConfigV2NodeConfig(p []interface{}) *managementClient.GKENodeConfig {
	obj := &managementClient.GKENodeConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["disk_size_gb"].(int); ok {
		obj.DiskSizeGb = int64(v)
	}
	if v, ok := in["disk_type"].(string); ok && len(v) > 0 {
		obj.DiskType = v
	}
	if v, ok := in["image_type"].(string); ok && len(v) > 0 {
		obj.ImageType = v
	}
	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}
	if v, ok := in["local_ssd_count"].(int); ok {
		obj.LocalSsdCount = int64(v)
	}
	if v, ok := in["machine_type"].(string); ok && len(v) > 0 {
		obj.MachineType = v
	}
	if v, ok := in["oauth_scopes"].([]interface{}); ok {
		obj.OauthScopes = toArrayStringSorted(v)
	}
	if v, ok := in["preemptible"].(bool); ok {
		obj.Preemptible = v
	}
	if v, ok := in["tags"].([]interface{}); ok {
		obj.Tags = toArrayStringSorted(v)
	}
	if v, ok := in["taints"].([]interface{}); ok && len(v) > 0 {
		obj.Taints = expandClusterGKEConfigV2NodeTaintsConfig(v)
	}
	return obj
}

func expandClusterGKEConfigV2NodePoolManagement(p []interface{}) *managementClient.GKENodePoolManagement {
	obj := &managementClient.GKENodePoolManagement{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["auto_repair"].(bool); ok {
		obj.AutoRepair = v
	}
	if v, ok := in["auto_upgrade"].(bool); ok {
		obj.AutoUpgrade = v
	}
	return obj
}

func expandClusterGKEConfigV2NodePoolsConfig(p []interface{}) []managementClient.GKENodePoolConfig {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}
	out := make([]managementClient.GKENodePoolConfig, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.GKENodePoolConfig{}

		if v, ok := in["autoscaling"].([]interface{}); ok {
			obj.Autoscaling = expandClusterGKEConfigV2NodePoolConfigAutoscaling(v)
		}
		if v, ok := in["config"].([]interface{}); ok {
			obj.Config = expandClusterGKEConfigV2NodeConfig(v)
		}
		if v, ok := in["initial_node_count"].(int); ok {
			count := int64(v)
			obj.InitialNodeCount = &count
		}
		if v, ok := in["management"].([]interface{}); ok {
			obj.Management = expandClusterGKEConfigV2NodePoolManagement(v)
		}
		if v, ok := in["max_pods_constraint"].(int); ok {
			max := int64(v)
			obj.MaxPodsConstraint = &max
		}
		if v, ok := in["name"].(string); ok {
			obj.Name = &v
		}
		if v, ok := in["version"].(string); ok {
			obj.Version = &v
		}

		out[i] = obj
	}

	return out
}

func expandClusterGKEConfigV2PrivateClusterConfig(p []interface{}) *managementClient.GKEPrivateClusterConfig {
	obj := &managementClient.GKEPrivateClusterConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enable_private_endpoint"].(bool); ok {
		obj.EnablePrivateEndpoint = v
	}
	if v, ok := in["enable_private_nodes"].(bool); ok {
		obj.EnablePrivateNodes = v
	}
	if v, ok := in["master_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.MasterIpv4CidrBlock = v
	}
	return obj
}

func expandClusterGKEConfigV2(p []interface{}) *managementClient.GKEClusterConfigSpec {
	obj := &managementClient.GKEClusterConfigSpec{
		Labels: map[string]string{},
	}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	obj.ClusterName = in["name"].(string)
	obj.GoogleCredentialSecret = in["google_credential_secret"].(string)
	obj.ProjectID = in["project_id"].(string)
	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}
	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	if v, ok := in["imported"].(bool); ok {
		obj.Imported = v
		// Returning here for imported clusters
		if obj.Imported {
			return obj
		}
	}

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.KubernetesVersion = &v
	}

	if v, ok := in["cluster_addons"].([]interface{}); ok {
		obj.ClusterAddons = expandClusterGKEConfigV2ClusterAddons(v)
	}
	if v, ok := in["cluster_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.ClusterIpv4CidrBlock = &v
	}

	if v, ok := in["description"].(string); ok && len(v) > 0 {
		obj.Description = v
	}
	if v, ok := in["enable_kubernetes_alpha"].(bool); ok {
		obj.EnableKubernetesAlpha = &v
	}
	if v, ok := in["ip_allocation_policy"].([]interface{}); ok {
		obj.IPAllocationPolicy = expandClusterGKEConfigV2IPAllocationPolicy(v)
	}
	if v, ok := in["labels"].(map[string]interface{}); ok {
		labels := toMapString(v)
		obj.Labels = labels
	}
	if v, ok := in["locations"].([]interface{}); ok {
		locations := toArrayString(v)
		obj.Locations = locations
	}
	if v, ok := in["logging_service"].(string); ok && len(v) > 0 {
		obj.LoggingService = &v
	}
	if v, ok := in["maintenance_window"].(string); ok && len(v) > 0 {
		obj.MaintenanceWindow = &v
	}
	if v, ok := in["master_authorized_networks_config"].([]interface{}); ok {
		obj.MasterAuthorizedNetworksConfig = expandClusterGKEConfigV2MasterAuthorizedNetworksConfig(v)
	}
	if v, ok := in["monitoring_service"].(string); ok && len(v) > 0 {
		obj.MonitoringService = &v
	}
	if v, ok := in["network"].(string); ok && len(v) > 0 {
		obj.Network = &v
	}
	if v, ok := in["network_policy_enabled"].(bool); ok {
		obj.NetworkPolicyEnabled = &v
	}
	if v, ok := in["node_pools"].([]interface{}); ok {
		obj.NodePools = expandClusterGKEConfigV2NodePoolsConfig(v)
	}
	if v, ok := in["private_cluster_config"].([]interface{}); ok {
		obj.PrivateClusterConfig = expandClusterGKEConfigV2PrivateClusterConfig(v)
	}
	if v, ok := in["subnetwork"].(string); ok && len(v) > 0 {
		obj.Subnetwork = &v
	}

	return obj
}

// This fix is required due to some fields doesn't have proper type at Rancher go cli
func fixClusterGKEConfigV2(values map[string]interface{}) map[string]interface{} {
	if values == nil {
		return nil
	}

	affectedFields := map[string]interface{}{
		"enableKubernetesAlpha": newFalse(),
		"clusterIpv4Cidr":       newEmptyString(),
		"loggingService":        newEmptyString(),
		"monitoringService":     newEmptyString(),
		"network":               newEmptyString(),
		"subnetwork":            newEmptyString(),
		"networkPolicyEnabled":  newFalse(),
		"locations":             []string{},
		"maintenanceWindow":     newEmptyString(),
	}

	for k, v := range affectedFields {
		if values[k] == nil {
			values[k] = v
		}
	}

	return values
}
