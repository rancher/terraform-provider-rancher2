package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifierDingtalkConfig(in *managementClient.DingtalkConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["url"] = in.URL

	if len(in.ProxyURL) > 0 {
		obj["proxy_url"] = in.ProxyURL
	}

	if len(in.Secret) > 0 {
		obj["secret"] = in.Secret
	}

	return []interface{}{obj}

}

// Expanders

func expandNotifierDingtalkConfig(p []interface{}) *managementClient.DingtalkConfig {
	obj := &managementClient.DingtalkConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.URL = in["url"].(string)

	if v, ok := in["proxy_url"].(string); ok && len(v) > 0 {
		obj.ProxyURL = v
	}

	if v, ok := in["secret"].(string); ok && len(v) > 0 {
		obj.Secret = v
	}

	return obj
}
