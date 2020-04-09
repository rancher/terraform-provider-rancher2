package rancher2

// Flatteners

func flattenCloudCredentialLinode(in *linodeCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialLinode(p []interface{}) *linodeCredentialConfig {
	obj := &linodeCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	return obj
}
