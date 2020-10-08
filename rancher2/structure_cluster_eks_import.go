package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterEKSImport(in *managementClient.EKSClusterConfigSpec, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AmazonCredentialSecret) > 0 {
		obj["cloud_credential"] = in.AmazonCredentialSecret
	}

	if len(in.DisplayName) > 0 {
		obj["name"] = in.DisplayName
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterEKSImport(p []interface{}) *managementClient.EKSClusterConfigSpec {
	obj := &managementClient.EKSClusterConfigSpec{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cloud_credential"].(string); ok && len(v) > 0 {
		obj.AmazonCredentialSecret = v
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.DisplayName = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	obj.Imported = true

	return obj
}
