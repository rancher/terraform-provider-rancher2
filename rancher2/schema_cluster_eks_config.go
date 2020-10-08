package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterEKSKind   = "eks"
	clusterDriverEKS = "amazonelasticcontainerservice"
)

//Types

type AmazonElasticContainerServiceConfig struct {
	AMI                         string   `json:"ami,omitempty" yaml:"ami,omitempty"`
	AccessKey                   string   `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	AssociateWorkerNodePublicIP *bool    `json:"associateWorkerNodePublicIp,omitempty" yaml:"associateWorkerNodePublicIp,omitempty"`
	DesiredNodes                int64    `json:"desiredNodes,omitempty" yaml:"desiredNodes,omitempty"`
	DisplayName                 string   `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DriverName                  string   `json:"driverName,omitempty" yaml:"driverName,omitempty"`
	EBSEncryption               bool     `json:"ebsEncryption,omitempty" yaml:"ebsEncryption,omitempty"`
	InstanceType                string   `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	KeyPairName                 string   `json:"keyPairName,omitempty" yaml:"keyPairName,omitempty"`
	KubernetesVersion           string   `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	MaximumNodes                int64    `json:"maximumNodes,omitempty" yaml:"maximumNodes,omitempty"`
	MinimumNodes                int64    `json:"minimumNodes,omitempty" yaml:"minimumNodes,omitempty"`
	NodeVolumeSize              int64    `json:"nodeVolumeSize,omitempty" yaml:"nodeVolumeSize,omitempty"`
	Region                      string   `json:"region,omitempty" yaml:"region,omitempty"`
	SecretKey                   string   `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	SecurityGroups              []string `json:"securityGroups,omitempty" yaml:"securityGroups,omitempty"`
	ServiceRole                 string   `json:"serviceRole,omitempty" yaml:"serviceRole,omitempty"`
	SessionToken                string   `json:"sessionToken,omitempty" yaml:"sessionToken,omitempty"`
	Subnets                     []string `json:"subnets,omitempty" yaml:"subnets,omitempty"`
	UserData                    string   `json:"userData,omitempty" yaml:"userData,omitempty"`
	VirtualNetwork              string   `json:"virtualNetwork,omitempty" yaml:"virtualNetwork,omitempty"`
}

//Schemas

func clusterEKSConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The AWS Client ID to use",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The kubernetes master version",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The AWS Client Secret associated with the Client ID",
		},
		"ami": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A custom AMI ID to use for the worker nodes instead of the default",
		},
		"associate_worker_node_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Associate public ip EKS worker nodes",
		},
		"desired_nodes": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "The desired number of worker nodes",
		},
		"ebs_encryption": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enables EBS encryption of worker nodes",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "t2.medium",
			Description: "The type of machine to use for worker nodes",
		},
		"key_pair_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Allow user to specify key name to use",
		},
		"maximum_nodes": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "The maximum number of worker nodes",
		},
		"minimum_nodes": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The minimum number of worker nodes",
		},
		"node_volume_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     20,
			Description: "The volume size for each node",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "us-west-2",
			Description: "The AWS Region to create the EKS cluster in",
		},
		"security_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of security groups to use for the cluster",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"service_role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The service role to use to perform the cluster operations in AWS",
		},
		"session_token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "A session token to use with the client key and secret if applicable",
		},
		"subnets": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of subnets in the virtual network to use",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Pass user-data to the nodes to perform automated configuration tasks",
		},
		"virtual_network": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the virtual network to use",
		},
	}

	return s
}
