package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterAKSConfigV2NodePools(input []managementClient.AKSNodePool, p []interface{}) []interface{} {
	if input == nil {
		return nil
	}
	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}
		if i < len(p) && p[i] != nil {
			obj = p[i].(map[string]interface{})
		}

		if in.Name != nil {
			obj["name"] = *in.Name
		}
		if in.AvailabilityZones != nil && len(*in.AvailabilityZones) > 0 {
			obj["availability_zones"] = toArrayInterface(*in.AvailabilityZones)
		}
		if in.Count != nil {
			obj["count"] = int(*in.Count)
		}
		if in.EnableAutoScaling != nil {
			obj["enable_auto_scaling"] = *in.EnableAutoScaling
			// Assigning max_count and min_count just if *obj.EnableAutoScaling is true
			if *in.EnableAutoScaling {
				if in.MaxCount != nil {
					obj["max_count"] = int(*in.MaxCount)
				}
				if in.MinCount != nil {
					obj["min_count"] = int(*in.MinCount)
				}
			}
		}
		if in.MaxPods != nil {
			obj["max_pods"] = int(*in.MaxPods)
		}
		if len(in.Mode) > 0 {
			obj["mode"] = in.Mode
		}
		if in.OrchestratorVersion != nil {
			obj["orchestrator_version"] = *in.OrchestratorVersion
		}
		if in.OsDiskSizeGB != nil {
			obj["os_disk_size_gb"] = int(*in.OsDiskSizeGB)
		}
		if len(in.OsDiskType) > 0 {
			obj["os_disk_type"] = in.OsDiskType
		}
		if len(in.OsType) > 0 {
			obj["os_type"] = in.OsType
		}
		if len(in.VMSize) > 0 {
			obj["vm_size"] = in.VMSize
		}
		out[i] = obj
	}

	return out
}

func flattenClusterAKSConfigV2(in *managementClient.AKSClusterConfigSpec, p []interface{}) []interface{} {
	if in == nil {
		return nil
	}

	obj := map[string]interface{}{}
	if len(p) != 0 && p[0] != nil {
		obj = p[0].(map[string]interface{})
	}

	if len(in.AzureCredentialSecret) > 0 {
		obj["cloud_credential_id"] = in.AzureCredentialSecret
	}
	if len(in.ClusterName) > 0 {
		obj["name"] = in.ClusterName
	}
	if len(in.ResourceGroup) > 0 {
		obj["resource_group"] = in.ResourceGroup
	}
	if len(in.ResourceLocation) > 0 {
		obj["resource_location"] = in.ResourceLocation
	}
	obj["imported"] = in.Imported
	if in.Imported {
		// Return if the cluster is imported
		return []interface{}{obj}
	}

	if in.AuthBaseURL != nil && len(*in.AuthBaseURL) > 0 {
		obj["auth_base_url"] = *in.AuthBaseURL
	}
	if in.AuthorizedIPRanges != nil && len(*in.AuthorizedIPRanges) > 0 {
		obj["authorized_ip_ranges"] = toArrayInterface(*in.AuthorizedIPRanges)
	}
	if in.BaseURL != nil && len(*in.BaseURL) > 0 {
		obj["base_url"] = *in.BaseURL
	}
	if in.DNSPrefix != nil && len(*in.DNSPrefix) > 0 {
		obj["dns_prefix"] = *in.DNSPrefix
	}
	if in.HTTPApplicationRouting != nil {
		obj["http_application_routing"] = *in.HTTPApplicationRouting
	}
	if in.KubernetesVersion != nil && len(*in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = *in.KubernetesVersion
	}
	if in.LinuxAdminUsername != nil && len(*in.LinuxAdminUsername) > 0 {
		obj["linux_admin_username"] = *in.LinuxAdminUsername
	}
	if in.LinuxSSHPublicKey != nil && len(*in.LinuxSSHPublicKey) > 0 {
		obj["linux_ssh_public_key"] = *in.LinuxSSHPublicKey
	}
	if in.LoadBalancerSKU != nil && len(*in.LoadBalancerSKU) > 0 {
		obj["load_balancer_sku"] = *in.LoadBalancerSKU
	}
	if in.LogAnalyticsWorkspaceGroup != nil && len(*in.LogAnalyticsWorkspaceGroup) > 0 {
		obj["log_analytics_workspace_group"] = *in.LogAnalyticsWorkspaceGroup
	}
	if in.LogAnalyticsWorkspaceName != nil && len(*in.LogAnalyticsWorkspaceName) > 0 {
		obj["log_analytics_workspace_name"] = *in.LogAnalyticsWorkspaceName
	}
	if in.Monitoring != nil {
		obj["monitoring"] = *in.Monitoring
	}
	if in.NetworkDNSServiceIP != nil && len(*in.NetworkDNSServiceIP) > 0 {
		obj["network_dns_service_ip"] = *in.NetworkDNSServiceIP
	}
	if in.NetworkDockerBridgeCIDR != nil && len(*in.NetworkDockerBridgeCIDR) > 0 {
		obj["network_docker_bridge_cidr"] = *in.NetworkDockerBridgeCIDR
	}
	if in.NetworkPlugin != nil && len(*in.NetworkPlugin) > 0 {
		obj["network_plugin"] = *in.NetworkPlugin
	}
	if in.NetworkPodCIDR != nil && len(*in.NetworkPodCIDR) > 0 {
		obj["network_pod_cidr"] = *in.NetworkPodCIDR
	}
	if in.NetworkPolicy != nil && len(*in.NetworkPolicy) > 0 {
		obj["network_policy"] = *in.NetworkPolicy
	}
	if in.NetworkServiceCIDR != nil && len(*in.NetworkServiceCIDR) > 0 {
		obj["network_service_cidr"] = *in.NetworkServiceCIDR
	}
	if in.NodePools != nil && len(in.NodePools) > 0 {
		v, ok := obj["node_pools"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["node_pools"] = flattenClusterAKSConfigV2NodePools(in.NodePools, v)
	}
	if in.PrivateCluster != nil {
		obj["private_cluster"] = *in.PrivateCluster
	}
	if in.Subnet != nil && len(*in.Subnet) > 0 {
		obj["subnet"] = *in.Subnet
	}
	if in.Tags != nil && len(in.Tags) > 0 {
		obj["tags"] = toMapInterface(in.Tags)
	}
	if in.VirtualNetwork != nil && len(*in.VirtualNetwork) > 0 {
		obj["virtual_network"] = *in.VirtualNetwork
	}
	if in.VirtualNetworkResourceGroup != nil && len(*in.VirtualNetworkResourceGroup) > 0 {
		obj["virtual_network_resource_group"] = *in.VirtualNetworkResourceGroup
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterAKSConfigV2NodePools(p []interface{}) []managementClient.AKSNodePool {
	if p == nil || len(p) == 0 {
		return []managementClient.AKSNodePool{}
	}
	out := make([]managementClient.AKSNodePool, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.AKSNodePool{}

		if v, ok := in["name"].(string); ok {
			obj.Name = &v
		}
		if v, ok := in["availability_zones"].([]interface{}); ok {
			availabilityZones := toArrayString(v)
			obj.AvailabilityZones = &availabilityZones
		}
		if v, ok := in["count"].(int); ok {
			size := int64(v)
			obj.Count = &size
		}
		if v, ok := in["enable_auto_scaling"].(bool); ok {
			obj.EnableAutoScaling = &v
			// Assigning max_count and min_count just if *obj.EnableAutoScaling is true
			if *obj.EnableAutoScaling {
				if v, ok := in["max_count"].(int); ok {
					size := int64(v)
					obj.MaxCount = &size
				}
				if v, ok := in["min_count"].(int); ok {
					size := int64(v)
					obj.MinCount = &size
				}
			}
		}
		if v, ok := in["max_pods"].(int); ok {
			size := int64(v)
			obj.MaxPods = &size
		}
		if v, ok := in["mode"].(string); ok {
			obj.Mode = v
		}
		if v, ok := in["orchestrator_version"].(string); ok {
			obj.OrchestratorVersion = &v
		}
		if v, ok := in["os_disk_size_gb"].(int); ok {
			size := int64(v)
			obj.OsDiskSizeGB = &size
		}
		if v, ok := in["os_disk_type"].(string); ok {
			obj.OsDiskType = v
		}
		if v, ok := in["os_type"].(string); ok {
			obj.OsType = v
		}
		if v, ok := in["vm_size"].(string); ok {
			obj.VMSize = v
		}
		out[i] = obj
	}

	return out
}

func expandClusterAKSConfigV2(p []interface{}) *managementClient.AKSClusterConfigSpec {
	obj := &managementClient.AKSClusterConfigSpec{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.AzureCredentialSecret = in["cloud_credential_id"].(string)
	obj.ClusterName = in["name"].(string)
	if v, ok := in["resource_group"].(string); ok {
		obj.ResourceGroup = v
	}
	if v, ok := in["resource_location"].(string); ok {
		obj.ResourceLocation = v
	}
	if v, ok := in["imported"].(bool); ok {
		obj.Imported = v
		if obj.Imported {
			// Return if the cluster is imported
			return obj
		}
	}

	// These fields should be assigned just when import=false, no matter if empty
	if v, ok := in["dns_prefix"].(string); ok {
		obj.DNSPrefix = &v
	}
	if v, ok := in["kubernetes_version"].(string); ok {
		obj.KubernetesVersion = &v
	}
	if v, ok := in["network_plugin"].(string); ok {
		obj.NetworkPlugin = &v
	}
	if v, ok := in["http_application_routing"].(bool); ok {
		obj.HTTPApplicationRouting = &v
	}
	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = &v
	}
	if v, ok := in["private_cluster"].(bool); ok {
		obj.PrivateCluster = &v
	}

	// These fields should be assigned just when import=false and no empty
	if v, ok := in["auth_base_url"].(string); ok && len(v) > 0 {
		obj.AuthBaseURL = &v
	}
	if v, ok := in["authorized_ip_ranges"].([]interface{}); ok && len(v) > 0 {
		authorizedIPRanges := toArrayString(v)
		obj.AuthorizedIPRanges = &authorizedIPRanges
	}
	if v, ok := in["base_url"].(string); ok && len(v) > 0 {
		obj.BaseURL = &v
	}
	if v, ok := in["linux_admin_username"].(string); ok && len(v) > 0 {
		obj.LinuxAdminUsername = &v
	}
	if v, ok := in["linux_ssh_public_key"].(string); ok && len(v) > 0 {
		obj.LinuxSSHPublicKey = &v
	}
	if v, ok := in["load_balancer_sku"].(string); ok && len(v) > 0 {
		obj.LoadBalancerSKU = &v
	}
	if v, ok := in["log_analytics_workspace_group"].(string); ok && len(v) > 0 {
		obj.LogAnalyticsWorkspaceGroup = &v
	}
	if v, ok := in["log_analytics_workspace_name"].(string); ok && len(v) > 0 {
		obj.LogAnalyticsWorkspaceName = &v
	}
	if v, ok := in["network_dns_service_ip"].(string); ok && len(v) > 0 {
		obj.NetworkDNSServiceIP = &v
	}
	if v, ok := in["network_docker_bridge_cidr"].(string); ok && len(v) > 0 {
		obj.NetworkDockerBridgeCIDR = &v
	}
	if v, ok := in["network_pod_cidr"].(string); ok && len(v) > 0 {
		obj.NetworkPodCIDR = &v
	}
	if v, ok := in["network_policy"].(string); ok && len(v) > 0 {
		obj.NetworkPolicy = &v
	}
	if v, ok := in["network_service_cidr"].(string); ok && len(v) > 0 {
		obj.NetworkServiceCIDR = &v
	}
	if v, ok := in["node_pools"].([]interface{}); ok && len(v) > 0 {
		nodePools := expandClusterAKSConfigV2NodePools(v)
		obj.NodePools = nodePools
	}
	if v, ok := in["subnet"].(string); ok && len(v) > 0 {
		obj.Subnet = &v
	}
	if v, ok := in["tags"].(map[string]interface{}); ok && len(v) > 0 {
		tags := toMapString(v)
		obj.Tags = tags
	}
	if v, ok := in["virtual_network"].(string); ok && len(v) > 0 {
		obj.VirtualNetwork = &v
	}
	if v, ok := in["virtual_network_resource_group"].(string); ok && len(v) > 0 {
		obj.VirtualNetworkResourceGroup = &v
	}

	return obj
}
