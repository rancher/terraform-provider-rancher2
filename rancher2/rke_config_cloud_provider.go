package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	cloudProviderList = []string{"aws", "azure", "custom", "openstack", "vsphere"}
)

//Schemas

func cloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"azure_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: azureCloudProviderFields(),
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
				Schema: openstackCloudProviderFields(),
			},
		},
		"vsphere_cloud_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: vsphereCloudProviderFields(),
			},
		},
	}
	return s
}

// Flatteners

func flattenCloudProvider(in *managementClient.CloudProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.AzureCloudProvider != nil {
		azureProvider, err := flattenAzureCloudProvider(in.AzureCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["azure_cloud_provider"] = azureProvider
	}

	if len(in.CustomCloudProvider) > 0 {
		obj["custom_cloud_provider"] = in.CustomCloudProvider
	}

	if len(in.Name) > 0 {
		obj["name"] = in.Name
	}

	if in.OpenstackCloudProvider != nil {
		openstackProvider, err := flattenOpenstackCloudProvider(in.OpenstackCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["openstack_cloud_provider"] = openstackProvider
	}

	if in.VsphereCloudProvider != nil {
		vsphereProvider, err := flattenVsphereCloudProvider(in.VsphereCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["vsphere_cloud_provider"] = vsphereProvider
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandCloudProvider(p []interface{}) (*managementClient.CloudProvider, error) {
	obj := &managementClient.CloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["azure_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		azureProvider, err := expandAzureCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.AzureCloudProvider = azureProvider
	}

	if v, ok := in["custom_cloud_provider"].(string); ok && len(v) > 0 {
		obj.CustomCloudProvider = v
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["openstack_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		openstackProvider, err := expandOpenstackCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.OpenstackCloudProvider = openstackProvider
	}

	if v, ok := in["vsphere_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		vsphereProvider, err := expandVsphereCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.VsphereCloudProvider = vsphereProvider
	}

	return obj, nil
}
