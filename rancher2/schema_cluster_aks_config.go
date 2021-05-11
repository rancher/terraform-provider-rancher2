package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterAKSKind                    = "aks"
	clusterDriverAKS                  = "azurekubernetesservice"
	clusterAKSLoadBalancerSkuBasic    = "basic"
	clusterAKSLoadBalancerSkuStandard = "standard"
)

var (
	clusterAKSAgentStorageProfile = []string{"ManagedDisks", "StorageAccount"}
	clusterAKSNetworkPlugin       = []string{"azure", "kubenet"}
	clusterAKSNetworkPolicy       = []string{"calico"}
	clusterAKSLoadBalancerSkuList = []string{
		clusterAKSLoadBalancerSkuBasic,
		clusterAKSLoadBalancerSkuStandard,
	}
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

//Schemas

func clusterAKSConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"agent_dns_prefix": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "DNS prefix to be used to create the FQDN for the agent pool",
		},
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure client ID to use",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure client secret associated with the \"client id\"",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specify the version of Kubernetes",
		},
		"master_dns_prefix": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "DNS prefix to use the Kubernetes cluster control pane",
		},
		"resource_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Cluster resource group",
		},
		"ssh_public_key_contents": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Contents of the SSH public key used to authenticate with Linux hosts",
		},
		"subnet": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of an existing Azure Virtual Subnet. Composite of agent virtual network subnet ID",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subscription credentials which uniquely identify Microsoft Azure subscription",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Azure tenant ID to use",
		},
		"virtual_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of an existing Azure Virtual Network. Composite of agent virtual network subnet ID",
		},
		"virtual_network_resource_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The resource group of an existing Azure Virtual Network. Composite of agent virtual network subnet ID",
		},
		"add_client_app_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The ID of an Azure Active Directory client application of type \"Native\". This application is for user login via kubectl",
		},
		"add_server_app_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The ID of an Azure Active Directory server application of type \"Web app/API\". This application represents the managed cluster's apiserver (Server application)",
		},
		"aad_server_app_secret": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The secret of an Azure Active Directory server application",
		},
		"aad_tenant_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The ID of an Azure Active Directory tenant",
		},
		"admin_username": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "azureuser",
			Description: "The administrator username to use for Linux hosts",
		},
		"agent_os_disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "GB size to be used to specify the disk for every machine in the agent pool. If you specify 0, it will apply the default according to the \"agent vm size\" specified",
		},
		"agent_pool_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "agentpool0",
			Description: "Name for the agent pool, upto 12 alphanumeric characters",
		},
		"agent_storage_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ManagedDisks",
			Description: "Storage profile specifies what kind of storage used on machine in the agent pool. Chooses from [ManagedDisks StorageAccount]",
		},
		"agent_vm_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Standard_D1_v2",
			Description: "Size of machine in the agent pool",
		},
		"auth_base_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "https://login.microsoftonline.com/",
			Description: "Different authentication API url to use",
		},
		"base_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "https://management.azure.com/",
			Description: "Different resource management API url to use",
		},
		"count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Number of machines (VMs) in the agent pool. Allowed values must be in the range of 1 to 100 (inclusive)",
		},
		"dns_service_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "10.0.0.10",
			Description: "An IP address assigned to the Kubernetes DNS service. It must be within the Kubernetes Service address range specified in \"service cidr\"",
		},
		"docker_bridge_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "172.17.0.1/16",
			Description: "A CIDR notation IP range assigned to the Docker bridge network. It must not overlap with any Subnet IP ranges or the Kubernetes Service address range specified in \"service cidr\"",
		},
		"enable_http_application_routing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable the Kubernetes ingress with automatic public DNS name creation",
		},
		"enable_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Turn on Azure Log Analytics monitoring. Uses the Log Analytics \"Default\" workspace if it exists, else creates one. if using an existing workspace, specifies \"log analytics workspace resource id\"",
		},
		"load_balancer_sku": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "Load balancer type (basic | standard). Must be standard for auto-scaling",
			ValidateFunc: validation.StringInSlice(clusterAKSLoadBalancerSkuList, true),
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "eastus",
			Description: "Azure Kubernetes cluster location",
		},
		"log_analytics_workspace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses '{resource group}-{subscription id}-{location code}'",
		},
		"log_analytics_workspace_resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The resource group of an existing Azure Log Analytics Workspace to use for storing monitoring data. If not specified, uses the 'Cluster' resource group",
		},
		"max_pods": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     110,
			Description: "Maximum number of pods that can run on a node",
		},
		"network_plugin": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "azure",
			Description:  "Network plugin used for building Kubernetes network. Chooses from [azure kubenet]",
			ValidateFunc: validation.StringInSlice(clusterAKSNetworkPlugin, true),
		},
		"network_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Network policy used for building Kubernetes network. Chooses from [calico]",
			ValidateFunc: validation.StringInSlice(clusterAKSNetworkPolicy, true),
		},
		"pod_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "172.244.0.0/16",
			Description: "A CIDR notation IP range from which to assign Kubernetes Pod IPs when \"network plugin\" is specified in \"kubenet\".",
		},
		"service_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "10.0.0.0/16",
			Description: "A CIDR notation IP range from which to assign Kubernetes Service cluster IPs. It must not overlap with any Subnet IP ranges",
		},
		"tag": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Tags for Kubernetes cluster. For example, foo=bar",
			Deprecated:  "Use tags argument instead as []string",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Tags for Kubernetes cluster. For example, `[\"foo=bar\",\"bar=foo\"]`",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	return s
}
