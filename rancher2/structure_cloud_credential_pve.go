package rancher2

// Flatteners

func flattenCloudCredentialPve(in map[string]interface{}, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if v, ok := in["pveUrl"].(string); ok && len(v) > 0 {
		obj["pve_url"] = v
	}
	if v, ok := in["pveTokenId"].(string); ok && len(v) > 0 {
		obj["pve_token_id"] = v
	}
	if v, ok := in["pveTokenSecret"].(string); ok && len(v) > 0 {
		obj["pve_token_secret"] = v
	}
	if v, ok := in["pveInsecureTls"].(bool); ok {
		obj["pve_insecure_tls"] = v
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialPve(p []interface{}) map[string]interface{} {
	obj := make(map[string]interface{})
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["pve_url"].(string); ok && len(v) > 0 {
		obj["pveUrl"] = v
	}
	if v, ok := in["pve_token_id"].(string); ok && len(v) > 0 {
		obj["pveTokenId"] = v
	}
	if v, ok := in["pve_token_secret"].(string); ok && len(v) > 0 {
		obj["pveTokenSecret"] = v
	}
	if v, ok := in["pve_insecure_tls"].(bool); ok {
		obj["pveInsecureTls"] = v
	}

	return obj
}
