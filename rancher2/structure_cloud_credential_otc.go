package rancher2

// Flatteners

func flattenCloudCredentialOpenTelekomCloud(in *openTelekomCloudCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.UserName) > 0 {
		obj["user_name"] = in.UserName
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialOpenTelekomCloud(p []interface{}) *openTelekomCloudCredentialConfig {
	obj := &openTelekomCloudCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["user_name"].(string); ok && len(v) > 0 {
		obj.UserName = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	return obj
}
