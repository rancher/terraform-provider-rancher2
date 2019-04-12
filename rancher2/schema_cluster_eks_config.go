package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	clusterEKSKind = "eks"
)

//Schemas

func clusterEKSConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Access key for EKS",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Required Secret key for EKS",
		},
		"ami": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AMI image for EKS worker nodes",
		},
		"associate_worker_node_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Associate public ip EKS worker nodes",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required Instance type for EKS",
		},
		"maximum_nodes": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required maximum number of node",
		},
		"minimum_nodes": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Required minimum number of nodes",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required EKS region",
		},
		"security_groups": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Required EKS Security groups",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"service_role": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required EKS region",
		},
		"subnets": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Required EKS subnets",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"virtual_network": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Required EKS region",
		},
	}

	return s
}
