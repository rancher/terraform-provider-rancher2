package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterAKSConfig(in *managementClient.AzureKubernetesServiceConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AdminUsername) > 0 {
		obj["admin_username"] = in.AdminUsername
	}

	if len(in.AgentDNSPrefix) > 0 {
		obj["agent_dns_prefix"] = in.AgentDNSPrefix
	}

	if len(in.AgentPoolName) > 0 {
		obj["agent_pool_name"] = in.AgentPoolName
	}

	if len(in.AgentVMSize) > 0 {
		obj["agent_vm_size"] = in.AgentVMSize
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

	if len(in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = in.KubernetesVersion
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	if len(in.MasterDNSPrefix) > 0 {
		obj["master_dns_prefix"] = in.MasterDNSPrefix
	}

	if in.OsDiskSizeGB > 0 {
		obj["os_disk_size_gb"] = int(in.OsDiskSizeGB)
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

	if len(in.Tag) > 0 {
		obj["tag"] = toMapInterface(in.Tag)
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

func expandClusterAKSConfig(p []interface{}) (*managementClient.AzureKubernetesServiceConfig, error) {
	obj := &managementClient.AzureKubernetesServiceConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["admin_username"].(string); ok && len(v) > 0 {
		obj.AdminUsername = v
	}

	if v, ok := in["agent_dns_prefix"].(string); ok && len(v) > 0 {
		obj.AgentDNSPrefix = v
	}

	if v, ok := in["agent_pool_name"].(string); ok && len(v) > 0 {
		obj.AgentPoolName = v
	}

	if v, ok := in["agent_vm_size"].(string); ok && len(v) > 0 {
		obj.AgentVMSize = v
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

	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		obj.KubernetesVersion = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["master_dns_prefix"].(string); ok && len(v) > 0 {
		obj.MasterDNSPrefix = v
	}

	if v, ok := in["os_disk_size_gb"].(int); ok && v > 0 {
		obj.OsDiskSizeGB = int64(v)
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
		obj.Tag = toMapString(v)
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
