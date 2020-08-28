package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifierSlackConfig(in *managementClient.SlackConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["default_recipient"] = in.DefaultRecipient
	obj["url"] = in.URL

	if len(in.ProxyURL) > 0 {
		obj["proxy_url"] = in.ProxyURL
	}

	return []interface{}{obj}

}

// Expanders

func expandNotifierSlackConfig(p []interface{}) *managementClient.SlackConfig {
	obj := &managementClient.SlackConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.DefaultRecipient = in["default_recipient"].(string)
	obj.URL = in["url"].(string)

	if v, ok := in["proxy_url"].(string); ok && len(v) > 0 {
		obj.ProxyURL = v
	}

	return obj
}
