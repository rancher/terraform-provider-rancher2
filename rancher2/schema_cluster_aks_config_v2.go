package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterAKSV2Kind   = "aksV2"
	clusterDriverAKSV2 = "AKS"
)

//Schemas

func clusterAKSConfigV2NodePoolsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The AKS node group name",
		},
		"availability_zones": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The AKS node pool availability zones",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The AKS node pool count",
		},
		"enable_auto_scaling": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is AKS node pool auto scaling enabled?",
		},
		"max_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The AKS node pool max count",
		},
		"max_pods": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     110,
			Description: "The AKS node pool max pods",
		},
		"min_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The AKS node pool min count",
		},
		"mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "System",
			Description: "The AKS node pool mode",
		},
		"orchestrator_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AKS node pool orchestrator version",
		},
		"os_disk_size_gb": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     128,
			Description: "The AKS node pool os disk size gb",
		},
		"os_disk_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Managed",
			Description: "The AKS node pool os disk type",
		},
		"os_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Linux",
			Description: "Enable AKS node pool os type",
		},
		"vm_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS node pool vm size",
		},
	}

	return s
}

func clusterAKSConfigV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_credential_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The AKS Cloud Credential ID to use",
		},
		"resource_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The AKS resource group",
		},
		"resource_location": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The AKS resource location",
		},
		"auth_base_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AKS auth base url",
		},
		"authorized_ip_ranges": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The AKS authorized ip ranges",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"base_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AKS base url",
		},
		"dns_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The AKS dns prefix. Required if `import=false`",
		},
		"http_application_routing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable AKS http application routing?",
		},
		"imported": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is AKS cluster imported?",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The kubernetes master version. Required if `import=false`",
		},
		"linux_admin_username": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS linux admin username",
		},
		"linux_ssh_public_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS linux ssh public key",
		},
		"load_balancer_sku": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS load balancer sku",
		},
		"log_analytics_workspace_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS log analytics workspace group",
		},
		"log_analytics_workspace_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS log analytics workspace name",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Is AKS cluster monitoring enabled?",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS cluster name",
		},
		"network_dns_service_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS network dns service ip",
		},
		"network_docker_bridge_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS network docker bridge cidr",
		},
		"network_plugin": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AKS network plugin. Required if `import=false`",
		},
		"network_pod_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS network pod cidr",
		},
		"network_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS network policy",
		},
		"network_service_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS network service cidr",
		},
		"node_pools": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The AKS node pools to use. Required if `import=false`",
			Elem: &schema.Resource{
				Schema: clusterAKSConfigV2NodePoolsFields(),
			},
		},
		"private_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Is AKS cluster private?",
		},
		"subnet": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS subnet",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The AKS cluster tags",
		},
		"virtual_network": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS virtual network",
		},
		"virtual_network_resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS virtual network resource group",
		},
	}

	return s
}
