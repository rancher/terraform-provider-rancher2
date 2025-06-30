package rancher2

import managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"

func expandClusterImportedConfig(p []interface{}) *managementClient.ImportedConfig {
	obj := &managementClient.ImportedConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["private_registry_url"].(string); ok && len(v) > 0 {
		obj.PrivateRegistryURL = v
	}

	return obj
}

func flattenClusterImportedConfig(in *managementClient.ImportedConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.PrivateRegistryURL) > 0 {
		obj["private_registry_url"] = in.PrivateRegistryURL
	}

	return []interface{}{obj}, nil
}
