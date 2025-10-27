package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"aws_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderAwsFields(),
			},
		},
		"azure_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderAzureFields(),
			},
		},
		"custom_cloud_provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"openstack_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackFields(),
			},
		},
		"vsphere_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereFields(),
			},
		},
	}
	return s
}
