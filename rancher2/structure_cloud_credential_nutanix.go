package rancher2

// Flatteners

func flattenCloudCredentialNutanix(in *nutanixCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}

	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.Port) > 0 {
		obj["port"] = in.Port
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialNutanix(p []interface{}) *nutanixCredentialConfig {
	obj := &nutanixCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.Port = v
	}

	return obj
}
