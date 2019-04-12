package rancher2

import (
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterRKEConfigCloudProvider(in *managementClient.CloudProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.AzureCloudProvider != nil {
		azureProvider, err := flattenClusterRKEConfigCloudProviderAzure(in.AzureCloudProvider)
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
		openstackProvider, err := flattenClusterRKEConfigCloudProviderOpenstack(in.OpenstackCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["openstack_cloud_provider"] = openstackProvider
	}

	if in.VsphereCloudProvider != nil {
		vsphereProvider, err := flattenClusterRKEConfigCloudProviderVsphere(in.VsphereCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["vsphere_cloud_provider"] = vsphereProvider
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigCloudProvider(p []interface{}) (*managementClient.CloudProvider, error) {
	obj := &managementClient.CloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["azure_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		azureProvider, err := expandClusterRKEConfigCloudProviderAzure(v)
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
		openstackProvider, err := expandClusterRKEConfigCloudProviderOpenstack(v)
		if err != nil {
			return obj, err
		}
		obj.OpenstackCloudProvider = openstackProvider
	}

	if v, ok := in["vsphere_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		vsphereProvider, err := expandClusterRKEConfigCloudProviderVsphere(v)
		if err != nil {
			return obj, err
		}
		obj.VsphereCloudProvider = vsphereProvider
	}

	return obj, nil
}
