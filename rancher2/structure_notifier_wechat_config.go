package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifierWechatConfig(in *managementClient.WechatConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["agent"] = in.Agent
	obj["corp"] = in.Corp
	obj["default_recipient"] = in.DefaultRecipient

	if len(in.Secret) > 0 {
		obj["secret"] = in.Secret
	}

	if len(in.ProxyURL) > 0 {
		obj["proxy_url"] = in.ProxyURL
	}

	obj["recipient_type"] = in.RecipientType

	return []interface{}{obj}

}

// Expanders

func expandNotifierWechatConfig(p []interface{}) *managementClient.WechatConfig {
	obj := &managementClient.WechatConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.Agent = in["agent"].(string)
	obj.Corp = in["corp"].(string)
	obj.DefaultRecipient = in["default_recipient"].(string)
	obj.Secret = in["secret"].(string)

	if v, ok := in["proxy_url"].(string); ok && len(v) > 0 {
		obj.ProxyURL = v
	}

	obj.RecipientType = in["recipient_type"].(string)

	return obj
}
