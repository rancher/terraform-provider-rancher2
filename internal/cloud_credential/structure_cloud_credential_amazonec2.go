package rancher2

// Flatteners

func flattenCloudCredentialAmazonec2(in *amazonec2CredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.DefaultRegion) > 0 {
		obj["default_region"] = in.DefaultRegion
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialAmazonec2(p []interface{}) *amazonec2CredentialConfig {
	obj := &amazonec2CredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["default_region"].(string); ok && len(v) > 0 {
		obj.DefaultRegion = v
	}

	return obj
}
