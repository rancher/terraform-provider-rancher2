package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigCloudProviderAwsGlobalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"disable_security_group_ingress": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"disable_strict_zone_check": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"elb_security_group": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"kubernetes_cluster_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"kubernetes_cluster_tag": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"role_arn": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"route_table_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"zone": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderAwsServiceOverrideFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"service": {
			Type:     schema.TypeString,
			Required: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"signing_method": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"signing_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"signing_region": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderAwsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderAwsGlobalFields(),
			},
		},
		"service_override": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderAwsServiceOverrideFields(),
			},
		},
	}
	return s
}
