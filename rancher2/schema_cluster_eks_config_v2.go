package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterEKSV2Kind                     = "eksV2"
	clusterDriverEKSV2                   = "EKS"
	clusterEKSV2LoggingAudit             = "audit"
	clusterEKSV2LoggingAPI               = "api"
	clusterEKSV2LoggingScheduler         = "scheduler"
	clusterEKSV2LoggingcontrollerManager = "controllerManager"
	clusterEKSV2LoggingAuthenticator     = "authenticator"
)

//Schemas

func clusterEKSConfigV2NodeGroupsLaunchTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EKS node group launch template ID",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The EKS node group launch template name",
		},
		"version": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The EKS node group launch template version",
		},
	}

	return s
}

func clusterEKSConfigV2NodeGroupsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EKS node group name",
		},
		"desired_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     2,
			Description: "The EKS node group desired size",
		},
		"disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     20,
			Description: "The EKS node group disk size",
		},
		"ec2_ssh_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The EKS node group ssh key",
		},
		"gpu": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is EKS cluster using gpu?",
		},
		"image_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The EKS node group image ID",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "t3.medium",
			Description: "The EKS node group instance type",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The EKS node group tags",
		},
		"launch_template": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The EKS node groups launch template",
			Elem: &schema.Resource{
				Schema: clusterEKSConfigV2NodeGroupsLaunchTemplateFields(),
			},
		},
		"max_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     2,
			Description: "The EKS node group maximum size",
		},
		"min_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     2,
			Description: "The EKS node group minimum size",
		},
		"node_role": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The EKS node group node role ARN",
		},
		"request_spot_instances": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable EKS node group request spot instances",
		},
		"resource_tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The EKS node group resource tags",
		},
		"spot_instance_types": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The EKS node group spot instance types",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"subnets": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The EKS node group subnets",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The EKS node group tags",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The EKS node group user data",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The EKS node group k8s version",
		},
	}

	return s
}

func clusterEKSConfigV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_credential_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The AWS Cloud Credential ID to use",
		},
		"imported": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is EKS cluster imported?",
		},
		"kms_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AWS kms key to use",
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The kubernetes master version",
		},
		"logging_types": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The AWS logging types",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The EKS cluster name",
		},
		"node_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The AWS node groups to use",
			Elem: &schema.Resource{
				Schema: clusterEKSConfigV2NodeGroupsFields(),
			},
		},
		"private_access": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "The EKS cluster has private access",
		},
		"public_access": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "The EKS cluster has public access",
		},
		"public_access_sources": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "The EKS cluster public access sources",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "us-west-2",
			Description: "The AWS Region to create the EKS cluster in",
		},
		"secrets_encryption": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable EKS cluster secret encryption",
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
			Description: "The AWS service role to use",
		},
		"subnets": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of subnets in the virtual network to use",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The EKS cluster tags",
		},
	}

	return s
}
