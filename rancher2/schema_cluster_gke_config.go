package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterGKEKind   = "gke"
	clusterDriverGKE = "googlekubernetesengine"
)

//Types

type GoogleKubernetesEngineConfig struct {
	ClusterIpv4Cidr                    string            `json:"clusterIpv4Cidr,omitempty" yaml:"clusterIpv4Cidr,omitempty"`
	Credential                         string            `json:"credential,omitempty" yaml:"credential,omitempty"`
	Description                        string            `json:"description,omitempty" yaml:"description,omitempty"`
	DiskSizeGb                         int64             `json:"diskSizeGb,omitempty" yaml:"diskSizeGb,omitempty"`
	DiskType                           string            `json:"diskType,omitempty" yaml:"diskType,omitempty"`
	DisplayName                        string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DriverName                         string            `json:"driverName,omitempty" yaml:"driverName,omitempty"`
	EnableAlphaFeature                 bool              `json:"enableAlphaFeature,omitempty" yaml:"enableAlphaFeature,omitempty"`
	EnableAutoRepair                   bool              `json:"enableAutoRepair,omitempty" yaml:"enableAutoRepair,omitempty"`
	EnableAutoUpgrade                  bool              `json:"enableAutoUpgrade,omitempty" yaml:"enableAutoUpgrade,omitempty"`
	EnableHorizontalPodAutoscaling     *bool             `json:"enableHorizontalPodAutoscaling,omitempty" yaml:"enableHorizontalPodAutoscaling,omitempty"`
	EnableHTTPLoadBalancing            *bool             `json:"enableHttpLoadBalancing,omitempty" yaml:"enableHttpLoadBalancing,omitempty"`
	EnableKubernetesDashboard          bool              `json:"enableKubernetesDashboard,omitempty" yaml:"enableKubernetesDashboard,omitempty"`
	EnableLegacyAbac                   bool              `json:"enableLegacyAbac,omitempty" yaml:"enableLegacyAbac,omitempty"`
	EnableMasterAuthorizedNetwork      bool              `json:"enableMasterAuthorizedNetwork,omitempty" yaml:"enableMasterAuthorizedNetwork,omitempty"`
	EnableNetworkPolicyConfig          *bool             `json:"enableNetworkPolicyConfig,omitempty" yaml:"enableNetworkPolicyConfig,omitempty"`
	EnableNodepoolAutoscaling          bool              `json:"enableNodepoolAutoscaling,omitempty" yaml:"enableNodepoolAutoscaling,omitempty"`
	EnablePrivateEndpoint              bool              `json:"enablePrivateEndpoint,omitempty" yaml:"enablePrivateEndpoint,omitempty"`
	EnablePrivateNodes                 bool              `json:"enablePrivateNodes,omitempty" yaml:"enablePrivateNodes,omitempty"`
	EnableStackdriverLogging           *bool             `json:"enableStackdriverLogging,omitempty" yaml:"enableStackdriverLogging,omitempty"`
	EnableStackdriverMonitoring        *bool             `json:"enableStackdriverMonitoring,omitempty" yaml:"enableStackdriverMonitoring,omitempty"`
	ImageType                          string            `json:"imageType,omitempty" yaml:"imageType,omitempty"`
	IPPolicyClusterIpv4CidrBlock       string            `json:"ipPolicyClusterIpv4CidrBlock,omitempty" yaml:"ipPolicyClusterIpv4CidrBlock,omitempty"`
	IPPolicyClusterSecondaryRangeName  string            `json:"ipPolicyClusterSecondaryRangeName,omitempty" yaml:"ipPolicyClusterSecondaryRangeName,omitempty"`
	IPPolicyCreateSubnetwork           bool              `json:"ipPolicyCreateSubnetwork,omitempty" yaml:"ipPolicyCreateSubnetwork,omitempty"`
	IPPolicyNodeIpv4CidrBlock          string            `json:"ipPolicyNodeIpv4CidrBlock,omitempty" yaml:"ipPolicyNodeIpv4CidrBlock,omitempty"`
	IPPolicyServicesIpv4CidrBlock      string            `json:"ipPolicyServicesIpv4CidrBlock,omitempty" yaml:"ipPolicyServicesIpv4CidrBlock,omitempty"`
	IPPolicyServicesSecondaryRangeName string            `json:"ipPolicyServicesSecondaryRangeName,omitempty" yaml:"ipPolicyServicesSecondaryRangeName,omitempty"`
	IPPolicySubnetworkName             string            `json:"ipPolicySubnetworkName,omitempty" yaml:"ipPolicySubnetworkName,omitempty"`
	IssueClientCertificate             bool              `json:"issueClientCertificate,omitempty" yaml:"issueClientCertificate,omitempty"`
	KubernetesDashboard                bool              `json:"kubernetesDashboard,omitempty" yaml:"kubernetesDashboard,omitempty"`
	Labels                             map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	LocalSsdCount                      int64             `json:"localSsdCount,omitempty" yaml:"localSsdCount,omitempty"`
	Locations                          []string          `json:"locations,omitempty" yaml:"locations,omitempty"`
	MachineType                        string            `json:"machineType,omitempty" yaml:"machineType,omitempty"`
	MaintenanceWindow                  string            `json:"maintenanceWindow,omitempty" yaml:"maintenanceWindow,omitempty"`
	MasterAuthorizedNetworkCidrBlocks  []string          `json:"masterAuthorizedNetworkCidrBlocks,omitempty" yaml:"masterAuthorizedNetworkCidrBlocks,omitempty"`
	MasterIpv4CidrBlock                string            `json:"masterIpv4CidrBlock,omitempty" yaml:"masterIpv4CidrBlock,omitempty"`
	MasterVersion                      string            `json:"masterVersion,omitempty" yaml:"masterVersion,omitempty"`
	MaxNodeCount                       int64             `json:"maxNodeCount,omitempty" yaml:"maxNodeCount,omitempty"`
	MinNodeCount                       int64             `json:"minNodeCount,omitempty" yaml:"minNodeCount,omitempty"`
	Name                               string            `json:"name,omitempty" yaml:"name,omitempty"`
	Network                            string            `json:"network,omitempty" yaml:"network,omitempty"`
	NodeCount                          int64             `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
	NodePool                           string            `json:"nodePool,omitempty" yaml:"nodePool,omitempty"`
	NodeVersion                        string            `json:"nodeVersion,omitempty" yaml:"nodeVersion,omitempty"`
	OauthScopes                        []string          `json:"oauthScopes,omitempty" yaml:"oauthScopes,omitempty"`
	Preemptible                        bool              `json:"preemptible,omitempty" yaml:"preemptible,omitempty"`
	ProjectID                          string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Region                             string            `json:"region,omitempty" yaml:"region,omitempty"`
	ResourceLabels                     map[string]string `json:"resourceLabels,omitempty" yaml:"resourceLabels,omitempty"`
	ServiceAccount                     string            `json:"serviceAccount,omitempty" yaml:"serviceAccount,omitempty"`
	SubNetwork                         string            `json:"subNetwork,omitempty" yaml:"subNetwork,omitempty"`
	UseIPAliases                       bool              `json:"useIpAliases,omitempty" yaml:"useIpAliases,omitempty"`
	Taints                             []string          `json:"taints,omitempty" yaml:"taints,omitempty"`
	Zone                               string            `json:"zone,omitempty" yaml:"zone,omitempty"`
}

//Schemas

func clusterGKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_ipv4_cidr": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP address range of the container pods",
		},
		"credential": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The contents of the GC credential file",
		},
		"disk_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Type of the disk attached to each node",
		},
		"image_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The image to use for the worker nodes",
		},
		"ip_policy_cluster_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP address range for the cluster pod IPs",
		},
		"ip_policy_cluster_secondary_range_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the secondary range to be used for the cluster CIDR block",
		},
		"ip_policy_node_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP address range of the instance IPs in this cluster",
		},
		"ip_policy_services_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP address range of the services IPs in this cluster",
		},
		"ip_policy_services_secondary_range_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the secondary range to be used for the services CIDR block",
		},
		"ip_policy_subnetwork_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A custom subnetwork name to be used if createSubnetwork is true",
		},
		"locations": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Locations to use for the cluster",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"machine_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The machine type to use for the worker nodes",
		},
		"maintenance_window": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "When to performance updates on the nodes, in 24-hour time",
		},
		"master_ipv4_cidr_block": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP range in CIDR notation to use for the hosted master network",
		},
		"master_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The kubernetes master version",
		},
		"network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The network to use for the cluster",
		},
		"node_pool": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the cluster node pool",
		},
		"node_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The version of kubernetes to use on the nodes",
		},
		"oauth_scopes": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The set of Google API scopes to be made available on all of the node VMs under the default service account",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of your project to use when creating a cluster",
		},
		"service_account": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Google Cloud Platform Service Account to be used by the node VMs",
		},
		"sub_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The sub-network to use for the cluster",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional description of this cluster",
		},
		"disk_size_gb": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     100,
			Description: "Size of the disk attached to each node",
		},
		"enable_alpha_feature": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "To enable kubernetes alpha feature",
		},
		"enable_auto_repair": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Specifies whether the node auto-repair is enabled for the node pool",
		},
		"enable_auto_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Specifies whether node auto-upgrade is enabled for the node pool",
		},
		"enable_horizontal_pod_autoscaling": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable horizontal pod autoscaling for the cluster",
		},
		"enable_http_load_balancing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable http load balancing for the cluster",
		},
		"enable_kubernetes_dashboard": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to enable the kubernetes dashboard",
		},
		"enable_legacy_abac": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to enable legacy abac on the cluster",
		},
		"enable_master_authorized_network": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether or not master authorized network is enabled",
		},
		"enable_network_policy_config": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable network policy config for the cluster",
		},
		"enable_nodepool_autoscaling": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable nodepool autoscaling",
		},
		"enable_private_endpoint": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the master's internal IP address is used as the cluster endpoint",
		},
		"enable_private_nodes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether nodes have internal IP address only",
		},
		"enable_stackdriver_logging": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable stackdriver logging",
		},
		"enable_stackdriver_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable stackdriver monitoring",
		},
		"ip_policy_create_subnetwork": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether a new subnetwork will be created automatically for the cluster",
		},
		"issue_client_certificate": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Issue a client certificate",
		},
		"kubernetes_dashboard": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable the kubernetes dashboard",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The map of Kubernetes labels (key/value pairs) to be applied to each node",
		},
		"local_ssd_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The number of local SSD disks to be attached to the node",
		},
		"master_authorized_network_cidr_blocks": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Define up to 10 external networks that could access Kubernetes master through HTTPS",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"max_node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Maximum number of nodes in the NodePool. Must be >= minNodeCount. There has to enough quota to scale up the cluster",
		},
		"min_node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Minimmum number of nodes in the NodePool. Must be >= 1 and <= maxNodeCount",
		},
		"node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "The number of nodes to create in this cluster",
		},
		"preemptible": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the nodes are created as preemptible VM instances",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The region to launch the cluster. Region or zone should be used",
		},
		"resource_labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "The map of Kubernetes labels (key/value pairs) to be applied to each cluster",
		},
		"use_ip_aliases": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether alias IPs will be used for pod IPs in the cluster",
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of kubernetes taints to be applied to each node",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The zone to launch the cluster. Zone or region should be used",
		},
	}

	return s
}
