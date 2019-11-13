package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	cloudProviderCustomName = "custom"
)

var (
	cloudProviderList = []string{cloudProviderAwsName, cloudProviderAzureName, cloudProviderCustomName, cloudProviderOpenstackName, cloudProviderVsphereName}
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
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(cloudProviderList, true),
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
