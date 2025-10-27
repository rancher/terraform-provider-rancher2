package rancher2

// Flatteners

func flattenCloudCredentialHarvester(in *harvesterCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.ClusterID) > 0 {
		obj["cluster_id"] = in.ClusterID
	}

	if len(in.ClusterType) > 0 {
		obj["cluster_type"] = in.ClusterType
	}

	if len(in.KubeconfigContent) > 0 {
		obj["kubeconfig_content"] = in.KubeconfigContent
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialHarvester(p []interface{}) *harvesterCredentialConfig {
	obj := &harvesterCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_id"].(string); ok && len(v) > 0 {
		obj.ClusterID = v
	}

	if v, ok := in["cluster_type"].(string); ok && len(v) > 0 {
		obj.ClusterType = v
	}

	if v, ok := in["kubeconfig_content"].(string); ok && len(v) > 0 {
		obj.KubeconfigContent = v
	}

	return obj
}
