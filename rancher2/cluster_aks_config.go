package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	clusterAksKind = "aks"
)

//Schemas

func aksConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"admin_username": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Admin username for AKS",
		},
		"agent_dns_prefix": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Agent dns prefix for AKS",
		},
		"agent_pool_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Agent pool name for AKS",
		},
		"agent_vm_size": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Agent vm size for AKS",
		},
		"base_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Base URL for AKS",
		},
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Client ID for AKS",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Client secret for AKS",
		},
		"count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required Number of agents for AKS cluster",
		},
		"location": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Location for AKS cluster",
		},
		"dns_service_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required DNS service ip for AKS cluster",
		},
		"docker_bridge_cidr": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Docker birdge CIDR for AKS cluster",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Kubernetes version for AKS cluster",
		},
		"master_dns_prefix": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Master dns prefix for AKS cluster",
		},
		"os_disk_size_gb": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required OS disk size for agents for AKS cluster",
		},
		"resource_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Resource group for AKS",
		},
		"ssh_public_key_contents": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required SSH public key for AKS",
		},
		"service_cidr": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Services CIDR for AKS",
		},
		"subnet": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Subnet for AKS",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Subscription ID for AKS",
		},
		"tag": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Optional Tags for AKS",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Tenant ID for AKS",
		},
		"virtual_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Virtual Network for AKS",
		},
		"virtual_network_resource_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Virtual Network resource group for AKS",
		},
	}

	return s
}

// Flatteners

func flattenAksConfig(in *managementClient.AzureKubernetesServiceConfig) ([]interface{}, error) {
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

func expandAksConfig(p []interface{}) (*managementClient.AzureKubernetesServiceConfig, error) {
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
