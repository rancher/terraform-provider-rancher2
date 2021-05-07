package rancher2

// Flatteners

func flattenCloudCredentialGoogle(in *googleCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AuthEncodedJSON) > 0 {
		obj["auth_encoded_json"] = in.AuthEncodedJSON
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialGoogle(p []interface{}) *googleCredentialConfig {
	obj := &googleCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["auth_encoded_json"].(string); ok && len(v) > 0 {
		obj.AuthEncodedJSON = v
	}

	return obj
}
