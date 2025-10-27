package rancher2

// Flatteners

func flattenCloudCredentialDigitalocean(in *digitaloceanCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialDigitalocean(p []interface{}) *digitaloceanCredentialConfig {
	obj := &digitaloceanCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	return obj
}
