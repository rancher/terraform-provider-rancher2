package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigCloudProvider(in *managementClient.CloudProvider, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.AWSCloudProvider != nil {
		awsProvider, err := flattenClusterRKEConfigCloudProviderAws(in.AWSCloudProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["aws_cloud_provider"] = awsProvider
	}

	if in.AzureCloudProvider != nil {
		v, ok := obj["azure_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		azureProvider, err := flattenClusterRKEConfigCloudProviderAzure(in.AzureCloudProvider, v)
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
		v, ok := obj["openstack_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		openstackProvider, err := flattenClusterRKEConfigCloudProviderOpenstack(in.OpenstackCloudProvider, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["openstack_cloud_provider"] = openstackProvider
	}

	if in.VsphereCloudProvider != nil {
		v, ok := obj["vsphere_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		vsphereProvider, err := flattenClusterRKEConfigCloudProviderVsphere(in.VsphereCloudProvider, v)
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

	if v, ok := in["aws_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		awsProvider, err := expandClusterRKEConfigCloudProviderAws(v)
		if err != nil {
			return obj, err
		}
		obj.AWSCloudProvider = awsProvider
	}

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

	if v, ok := in["name"].(string); ok {
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
