package rancher2

// Flatteners

func flattenClusterAKSConfig(in *AzureKubernetesServiceConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AADClientAppID) > 0 {
		obj["add_client_app_id"] = in.AADClientAppID
	}

	if len(in.AADServerAppID) > 0 {
		obj["add_server_app_id"] = in.AADServerAppID
	}

	if len(in.AADServerAppSecret) > 0 {
		obj["aad_server_app_secret"] = in.AADServerAppSecret
	}

	if len(in.AADTenantID) > 0 {
		obj["aad_tenant_id"] = in.AADTenantID
	}

	if len(in.AdminUsername) > 0 {
		obj["admin_username"] = in.AdminUsername
	}

	if len(in.AgentDNSPrefix) > 0 {
		obj["agent_dns_prefix"] = in.AgentDNSPrefix
	}

	if in.AgentOsdiskSizeGB > 0 {
		obj["agent_os_disk_size"] = int(in.AgentOsdiskSizeGB)
	}

	if len(in.AgentPoolName) > 0 {
		obj["agent_pool_name"] = in.AgentPoolName
	}

	if len(in.AgentStorageProfile) > 0 {
		obj["agent_storage_profile"] = in.AgentStorageProfile
	}

	if len(in.AgentVMSize) > 0 {
		obj["agent_vm_size"] = in.AgentVMSize
	}

	if len(in.AuthBaseURL) > 0 {
		obj["auth_base_url"] = in.AuthBaseURL
	}

	if len(in.BaseURL) > 0 {
		obj["base_url"] = in.BaseURL
	}

	if len(in.ClientID) > 0 {
		obj["client_id"] = in.ClientID
	}

	if len(in.ClientSecret) > 0 {
		obj["client_secret"] = in.ClientSecret
	}

	if in.Count > 0 {
		obj["count"] = int(in.Count)
	}

	if len(in.DNSServiceIP) > 0 {
		obj["dns_service_ip"] = in.DNSServiceIP
	}

	if len(in.DockerBridgeCIDR) > 0 {
		obj["docker_bridge_cidr"] = in.DockerBridgeCIDR
	}

	obj["enable_http_application_routing"] = in.EnableHTTPApplicationRouting
	obj["enable_monitoring"] = *in.EnableMonitoring

	if len(in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = in.KubernetesVersion
	}

	if len(in.LoadBalancerSku) > 0 {
		obj["load_balancer_sku"] = in.LoadBalancerSku
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	if len(in.LogAnalyticsWorkspace) > 0 {
		obj["log_analytics_workspace"] = in.LogAnalyticsWorkspace
	}

	if len(in.LogAnalyticsWorkspaceResourceGroup) > 0 {
		obj["log_analytics_workspace_resource_group"] = in.LogAnalyticsWorkspaceResourceGroup
	}

	if len(in.MasterDNSPrefix) > 0 {
		obj["master_dns_prefix"] = in.MasterDNSPrefix
	}

	if in.MaxPods > 0 {
		obj["max_pods"] = int(in.MaxPods)
	}

	if len(in.NetworkPlugin) > 0 {
		obj["network_plugin"] = in.NetworkPlugin
	}

	if len(in.NetworkPolicy) > 0 {
		obj["network_policy"] = in.NetworkPolicy
	}

	if len(in.PodCIDR) > 0 {
		obj["pod_cidr"] = in.PodCIDR
	}

	if len(in.ResourceGroup) > 0 {
		obj["resource_group"] = in.ResourceGroup
	}

	if len(in.SSHPublicKeyContents) > 0 {
		obj["ssh_public_key_contents"] = in.SSHPublicKeyContents
	}

	if len(in.ServiceCIDR) > 0 {
		obj["service_cidr"] = in.ServiceCIDR
	}

	if len(in.Subnet) > 0 {
		obj["subnet"] = in.Subnet
	}

	if len(in.SubscriptionID) > 0 {
		obj["subscription_id"] = in.SubscriptionID
	}

	if len(in.Tags) > 0 {
		obj["tags"] = toArrayInterface(in.Tags)
	}

	if len(in.TenantID) > 0 {
		obj["tenant_id"] = in.TenantID
	}

	if len(in.VirtualNetwork) > 0 {
		obj["virtual_network"] = in.VirtualNetwork
	}

	if len(in.VirtualNetworkResourceGroup) > 0 {
		obj["virtual_network_resource_group"] = in.VirtualNetworkResourceGroup
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterAKSConfig(p []interface{}, name string) (*AzureKubernetesServiceConfig, error) {
	obj := &AzureKubernetesServiceConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	obj.Name = name
	obj.DisplayName = name
	obj.DriverName = clusterDriverAKS

	if v, ok := in["add_client_app_id"].(string); ok && len(v) > 0 {
		obj.AADClientAppID = v
	}

	if v, ok := in["add_server_app_id"].(string); ok && len(v) > 0 {
		obj.AADServerAppID = v
	}

	if v, ok := in["aad_server_app_secret"].(string); ok && len(v) > 0 {
		obj.AADServerAppSecret = v
	}

	if v, ok := in["aad_tenant_id"].(string); ok && len(v) > 0 {
		obj.AADTenantID = v
	}

	if v, ok := in["admin_username"].(string); ok && len(v) > 0 {
		obj.AdminUsername = v
	}

	if v, ok := in["agent_dns_prefix"].(string); ok && len(v) > 0 {
		obj.AgentDNSPrefix = v
	}

	if v, ok := in["agent_os_disk_size"].(int); ok && v > 0 {
		obj.AgentOsdiskSizeGB = int64(v)
	}

	if v, ok := in["agent_pool_name"].(string); ok && len(v) > 0 {
		obj.AgentPoolName = v
	}

	if v, ok := in["agent_storage_profile"].(string); ok && len(v) > 0 {
		obj.AgentStorageProfile = v
	}

	if v, ok := in["agent_vm_size"].(string); ok && len(v) > 0 {
		obj.AgentVMSize = v
	}

	if v, ok := in["auth_base_url"].(string); ok && len(v) > 0 {
		obj.AuthBaseURL = v
	}

	if v, ok := in["base_url"].(string); ok && len(v) > 0 {
		obj.BaseURL = v
	}

	if v, ok := in["client_id"].(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in["client_secret"].(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in["count"].(int); ok && v > 0 {
		obj.Count = int64(v)
	}

	if v, ok := in["dns_service_ip"].(string); ok && len(v) > 0 {
		obj.DNSServiceIP = v
	}

	if v, ok := in["docker_bridge_cidr"].(string); ok && len(v) > 0 {
		obj.DockerBridgeCIDR = v
	}

	if v, ok := in["enable_http_application_routing"].(bool); ok {
		obj.EnableHTTPApplicationRouting = v
	}

	if v, ok := in["enable_monitoring"].(bool); ok {
		obj.EnableMonitoring = &v
	}

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.KubernetesVersion = v
	}

	if v, ok := in["load_balancer_sku"].(string); ok && len(v) > 0 {
		obj.LoadBalancerSku = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["log_analytics_workspace"].(string); ok && len(v) > 0 {
		obj.LogAnalyticsWorkspace = v
	}

	if v, ok := in["log_analytics_workspace_resource_group"].(string); ok && len(v) > 0 {
		obj.LogAnalyticsWorkspaceResourceGroup = v
	}

	if v, ok := in["master_dns_prefix"].(string); ok && len(v) > 0 {
		obj.MasterDNSPrefix = v
	}

	if v, ok := in["max_pods"].(int); ok && v > 0 {
		obj.MaxPods = int64(v)
	}

	if v, ok := in["network_plugin"].(string); ok && len(v) > 0 {
		obj.NetworkPlugin = v
	}

	if v, ok := in["network_policy"].(string); ok && len(v) > 0 {
		obj.NetworkPolicy = v
	}

	if v, ok := in["pod_cidr"].(string); ok && len(v) > 0 {
		obj.PodCIDR = v
	}

	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = v
	}

	if v, ok := in["ssh_public_key_contents"].(string); ok && len(v) > 0 {
		obj.SSHPublicKeyContents = v
	}

	if v, ok := in["service_cidr"].(string); ok && len(v) > 0 {
		obj.ServiceCIDR = v
	}

	if v, ok := in["subnet"].(string); ok && len(v) > 0 {
		obj.Subnet = v
	}

	if v, ok := in["subscription_id"].(string); ok && len(v) > 0 {
		obj.SubscriptionID = v
	}

	if v, ok := in["tag"].(map[string]interface{}); ok && len(v) > 0 {
		tagMap := toMapString(v)
		for k, value := range tagMap {
			obj.Tags = append(obj.Tags, k+"="+value)
		}
	}

	if v, ok := in["tags"].([]interface{}); ok && len(v) > 0 {
		obj.Tags = toArrayString(v)
	}

	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in["virtual_network"].(string); ok && len(v) > 0 {
		obj.VirtualNetwork = v
	}

	if v, ok := in["virtual_network_resource_group"].(string); ok && len(v) > 0 {
		obj.VirtualNetworkResourceGroup = v
	}

	return obj, nil
}
