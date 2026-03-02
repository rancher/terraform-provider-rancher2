// Copyright 2020 Oracle and/or its affiliates.

package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterOKEKind   = "oke"
	clusterDriverOKE = "oraclecontainerengine"
)

//Types

type OracleKubernetesEngineConfig struct {
	ClusterType                   string `json:"clusterType,omitempty" yaml:"clusterType,omitempty"`
	CompartmentID                 string `json:"compartmentId,omitempty" yaml:"compartmentId,omitempty"`
	ControlPlaneSubnetName        string `json:"controlPlaneSubnetName,omitempty" yaml:"controlPlaneSubnetName,omitempty"`
	CustomBootVolumeSize          int64  `json:"customBootVolumeSize,omitempty" yaml:"customBootVolumeSize,omitempty"`
	Description                   string `json:"description,omitempty" yaml:"description,omitempty"`
	DisplayName                   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DriverName                    string `json:"driverName,omitempty" yaml:"driverName,omitempty"`
	EnableKubernetesDashboard     bool   `json:"enableKubernetesDashboard,omitempty" yaml:"enableKubernetesDashboard,omitempty"`
	EvictionGraceDuration         string `json:"evictionGraceDuration,omitempty" yaml:"evictionGraceDuration,omitempty"`
	Fingerprint                   string `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	FlexMemoryInGBs               int64  `json:"flexMemoryInGbs,omitempty" yaml:"flexMemoryInGbs,omitempty"`
	FlexOCPUs                     int64  `json:"flexOcpus,omitempty" yaml:"flexOcpus,omitempty"`
	ForceDeleteAfterGraceDuration bool   `json:"forceDeleteAfterGraceDuration,omitempty" yaml:"forceDeleteAfterGraceDuration,omitempty"`
	ImageVerificationKmsKeyID     string `json:"imageVerificationKmsKeyId,omitempty" yaml:"imageVerificationKmsKeyId,omitempty"`
	KMSKeyID                      string `json:"kmsKeyId" yaml:"kmsKeyId"`
	KubernetesVersion             string `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	LimitNodeCount                int64  `json:"limitNodeCount,omitempty" yaml:"limitNodeCount,omitempty"`
	Name                          string `json:"name,omitempty" yaml:"name,omitempty"`
	NodeImage                     string `json:"nodeImage,omitempty" yaml:"nodeImage,omitempty"`
	NodePoolSubnetDNSDomainName   string `json:"nodePoolDnsDomainName,omitempty" yaml:"nodePoolDnsDomainName,omitempty"`
	NodePoolSubnetName            string `json:"nodePoolSubnetName,omitempty" yaml:"nodePoolSubnetName,omitempty"`
	NodePublicSSHKeyContents      string `json:"nodePublicKeyContents,omitempty" yaml:"nodePublicKeyContents,omitempty"`
	NodeShape                     string `json:"nodeShape,omitempty" yaml:"nodeShape,omitempty"`
	NodeUserDataContents          string `json:"nodeUserDataContents,omitempty" yaml:"nodeUserDataContents,omitempty"`
	PodNetwork                    string `json:"podNetwork,omitempty" yaml:"podNetwork,omitempty"`
	PodSubnetName                 string `json:"podSubnetName,omitempty" yaml:"podSubnetName,omitempty"`
	PodCidr                       string `json:"podCidr,omitempty" yaml:"podCidr,omitempty"`
	PrivateControlPlane           bool   `json:"enablePrivateControlPlane,omitempty" yaml:"enablePrivateControlPlane,omitempty"`
	PrivateKeyContents            string `json:"privateKeyContents,omitempty" yaml:"privateKeyContents,omitempty"`
	PrivateKeyPassphrase          string `json:"privateKeyPassphrase,omitempty" yaml:"privateKeyPassphrase,omitempty"`
	PrivateNodes                  bool   `json:"enablePrivateNodes,omitempty" yaml:"enablePrivateNodes,omitempty"`
	QuantityOfSubnets             int64  `json:"quantityOfNodeSubnets,omitempty" yaml:"quantityOfNodeSubnets,omitempty"`
	QuantityPerSubnet             int64  `json:"quantityPerSubnet,omitempty" yaml:"quantityPerSubnet,omitempty"`
	Region                        string `json:"region,omitempty" yaml:"region,omitempty"`
	ServiceCidr                   string `json:"serviceCidr,omitempty" yaml:"serviceCidr,omitempty"`
	ServiceLBSubnet1Name          string `json:"loadBalancerSubnetName1,omitempty" yaml:"loadBalancerSubnetName1,omitempty"`
	ServiceLBSubnet2Name          string `json:"loadBalancerSubnetName2,omitempty" yaml:"loadBalancerSubnetName2,omitempty"`
	ServiceSubnetDNSDomainName    string `json:"serviceDnsDomainName,omitempty" yaml:"serviceDnsDomainName,omitempty"`
	SkipVCNDelete                 bool   `json:"skipVcnDelete,omitempty" yaml:"skipVcnDelete,omitempty"`
	TenancyID                     string `json:"tenancyId,omitempty" yaml:"tenancyId,omitempty"`
	UserOCID                      string `json:"userOcid,omitempty" yaml:"userOcid,omitempty"`
	VcnCompartmentID              string `json:"vcnCompartmentId,omitempty" yaml:"vcnCompartmentId,omitempty"`
	VCNName                       string `json:"vcnName,omitempty" yaml:"vcnName,omitempty"`
	WorkerNodeIngressCidr         string `json:"workerNodeIngressCidr,omitempty" yaml:"workerNodeIngressCidr,omitempty"`
}

//Schemas

func clusterOKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{

		"cluster_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optionally specify a cluster type of basic or enhanced",
		},
		"compartment_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OCID of the compartment in which to create resources (VCN, worker nodes, etc.)",
		},
		"custom_boot_volume_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "An optional custom boot volume size (in GB) for the nodes",
		},
		"control_plane_subnet_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The (optional) name of a pre-existing subnet (public or private) for the Kubernetes API endpoint",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional description of this cluster",
		},
		"enable_kubernetes_dashboard": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable the kubernetes dashboard",
		},
		"enable_private_control_plane": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether Kubernetes API endpoint is a private IP only accessible from within the VCN",
		},
		"enable_private_nodes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether worker nodes are deployed into a new private subnet",
		},
		"eviction_grace_duration": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The optional grace period in minutes to allow cordon and drain to complete successfuly",
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fingerprint corresponding to the specified user's private API Key",
		},
		"flex_memory_in_gbs": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Optional amount of memory in GB for nodes (requires flexible node_shape)",
		},
		"flex_ocpus": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Optional number of OCPUs for nodes (requires flexible node_shape)",
		},
		"force_delete_after_grace_duration": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to send a SIGKILL signal if a pod does not terminate within the specified grace period",
		},
		"image_verification_kms_key_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional specify a comma separated list of master encryption key OCID(s) to verify images",
		},
		"kms_key_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional specify the OCID of the KMS Vault master key",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Kubernetes version that will be used for your master *and* worker nodes e.g. v1.33.1",
		},
		"limit_node_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Optional limit on the total number of nodes in the pool",
		},
		"load_balancer_subnet_name_1": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the first existing subnet to use for Kubernetes services / LB",
		},
		"load_balancer_subnet_name_2": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The (optional) name of a second existing subnet to use for Kubernetes services / LB",
		},
		"node_image": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OS for the node image",
		},
		"node_pool_dns_domain_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "nodedns",
			Description: "Optional name for DNS domain of node pool subnet",
		},
		"node_pool_subnet_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "nodedns",
			Description: "Optional pre-existing subnet (public or private) for nodes",
		},
		"node_public_key_contents": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The contents of the SSH public key file to use for the nodes",
		},
		"node_shape": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The shape of the node (determines number of CPUs and  amount of memory on each node)",
		},
		"node_user_data_contents": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The contents of custom cloud-init / user_data for the nodes - will be base64 encoded internally if it is not already",
		},
		"pod_network": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional Pod Network plugin. Choose flannel or native. Defaults to flannel",
		},
		"pod_subnet_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The (optional) name of a pre-existing subnet that pods will be assigned IPs from when using native pod networking",
		},
		"pod_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional specify the pod CIDR, defaults to 10.244.0.0/16",
		},
		"private_key_contents": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The private API key file contents for the specified user, in PEM format",
		},
		"private_key_passphrase": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The passphrase of the private key for the OKE cluster",
		},
		"quantity_of_node_subnets": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of node subnets (defaults to creating 1 regional subnet)",
		},
		"quantity_per_subnet": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Number of worker nodes in each subnet / availability domain",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The availability domain within the region to host the OKE cluster",
		},
		"service_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional specify the service CIDR, defaults to 10.96.0.0/16",
		},
		"service_dns_domain_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "svcdns",
			Description: "Optional name for DNS domain of service subnet",
		},
		"skip_vcn_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to skip deleting VCN",
		},
		"tenancy_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OCID of the tenancy in which to create resources",
		},
		"user_ocid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The OCID of a user who has access to the tenancy/compartment",
		},
		"vcn_compartment_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The OCID of the compartment (if different from compartment_id) in which to find the pre-existing virtual network set with vcn_name.",
		},
		"vcn_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The optional name of an existing virtual network to use for the cluster creation. A new VCN will be created if not specified.",
		},
		"worker_node_ingress_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional CIDR from which to allow ingress to worker nodes",
		},
	}

	return s
}
