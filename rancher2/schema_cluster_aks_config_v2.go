package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterAKSV2Kind   = "aksV2"
	clusterDriverAKSV2 = "AKS"
)

var (
	clusterAKSOutboundType = []string{"loadbalancer", "managednatgateway", "userassignednatgateway", "userdefinedrouting"}
)

//Types

type AzureKubernetesServiceConfig struct {
	AADClientAppID                     string   `json:"addClientAppId,omitempty" yaml:"addClientAppId,omitempty"`
	AADServerAppID                     string   `json:"addServerAppId,omitempty" yaml:"addServerAppId,omitempty"`
	AADServerAppSecret                 string   `json:"addServerAppSecret,omitempty" yaml:"addServerAppSecret,omitempty"`
	AADTenantID                        string   `json:"addTenantId,omitempty" yaml:"addTenantId,omitempty"`
	AdminUsername                      string   `json:"adminUsername,omitempty" yaml:"adminUsername,omitempty"`
	AgentDNSPrefix                     string   `json:"agentDnsPrefix,omitempty" yaml:"agentDnsPrefix,omitempty"`
	AgentOsdiskSizeGB                  int64    `json:"agentOsdiskSize,omitempty" yaml:"agentOsdiskSize,omitempty"`
	AgentPoolName                      string   `json:"agentPoolName,omitempty" yaml:"agentPoolName,omitempty"`
	AgentStorageProfile                string   `json:"agentStorageProfile,omitempty" yaml:"agentStorageProfile,omitempty"`
	AgentVMSize                        string   `json:"agentVmSize,omitempty" yaml:"agentVmSize,omitempty"`
	AuthBaseURL                        string   `json:"authBaseUrl" yaml:"authBaseUrl"`
	BaseURL                            string   `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
	ClientID                           string   `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret                       string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	Count                              int64    `json:"count,omitempty" yaml:"count,omitempty"`
	DisplayName                        string   `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DriverName                         string   `json:"driverName,omitempty" yaml:"driverName,omitempty"`
	DNSServiceIP                       string   `json:"dnsServiceIp,omitempty" yaml:"dnsServiceIp,omitempty"`
	DockerBridgeCIDR                   string   `json:"dockerBridgeCidr,omitempty" yaml:"dockerBridgeCidr,omitempty"`
	EnableHTTPApplicationRouting       bool     `json:"enableHttpApplicationRouting,omitempty" yaml:"enableHttpApplicationRouting,omitempty"`
	EnableMonitoring                   *bool    `json:"enableMonitoring,omitempty" yaml:"enableMonitoring,omitempty"`
	KubernetesVersion                  string   `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	LoadBalancerSku                    string   `json:"loadBalancerSku,omitempty" yaml:"loadBalancerSku,omitempty"`
	Location                           string   `json:"location,omitempty" yaml:"location,omitempty"`
	LogAnalyticsWorkspace              string   `json:"logAnalyticsWorkspace,omitempty" yaml:"logAnalyticsWorkspace,omitempty"`
	LogAnalyticsWorkspaceResourceGroup string   `json:"logAnalyticsWorkspaceResourceGroup,omitempty" yaml:"logAnalyticsWorkspaceResourceGroup,omitempty"`
	MasterDNSPrefix                    string   `json:"masterDnsPrefix,omitempty" yaml:"masterDnsPrefix,omitempty"`
	MaxPods                            int64    `json:"maxPods,omitempty" yaml:"maxPods,omitempty"`
	Name                               string   `json:"name,omitempty" yaml:"name,omitempty"`
	NetworkPlugin                      string   `json:"networkPlugin,omitempty" yaml:"networkPlugin,omitempty"`
	NetworkPolicy                      string   `json:"networkPolicy,omitempty" yaml:"networkPolicy,omitempty"`
	PodCIDR                            string   `json:"podCidr,omitempty" yaml:"podCidr,omitempty"`
	ResourceGroup                      string   `json:"resourceGroup,omitempty" yaml:"resourceGroup,omitempty"`
	SSHPublicKeyContents               string   `json:"sshPublicKeyContents,omitempty" yaml:"sshPublicKeyContents,omitempty"`
	ServiceCIDR                        string   `json:"serviceCidr,omitempty" yaml:"serviceCidr,omitempty"`
	Subnet                             string   `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	SubscriptionID                     string   `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
	Tags                               []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	TenantID                           string   `json:"tenantId,omitempty" yaml:"tenantId,omitempty"`
	VirtualNetwork                     string   `json:"virtualNetwork,omitempty" yaml:"virtualNetwork,omitempty"`
	VirtualNetworkResourceGroup        string   `json:"virtualNetworkResourceGroup,omitempty" yaml:"virtualNetworkResourceGroup,omitempty"`
}

// Schemas

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
		"max_surge": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AKS node pool max surge",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The AKS node pool labels",
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The AKS node pool taints",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
		"node_resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The AKS node resource group name",
		},
		"outbound_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "loadBalancer",
			Description:  "The AKS outbound type for the egress traffic",
			ValidateFunc: validation.StringInSlice(clusterAKSOutboundType, true),
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
