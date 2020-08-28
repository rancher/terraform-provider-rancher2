package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifierPagerdutyConfig(in *managementClient.PagerdutyConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["service_key"] = in.ServiceKey

	if len(in.ProxyURL) > 0 {
		obj["proxy_url"] = in.ProxyURL
	}

	return []interface{}{obj}

}

// Expanders

func expandNotifierPagerdutyConfig(p []interface{}) *managementClient.PagerdutyConfig {
	obj := &managementClient.PagerdutyConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.ServiceKey = in["service_key"].(string)

	if v, ok := in["proxy_url"].(string); ok && len(v) > 0 {
		obj.ProxyURL = v
	}

	return obj
}
