package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	clusterAKSKind = "aks"
)

//Schemas

func clusterAKSConfigFields() map[string]*schema.Schema {
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
