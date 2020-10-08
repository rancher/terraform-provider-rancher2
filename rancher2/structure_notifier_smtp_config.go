package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNotifierSMTPConfig(in *managementClient.SMTPConfig, p []interface{}) []interface{} {
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
	obj["host"] = in.Host
	obj["port"] = int(in.Port)
	obj["sender"] = in.Sender

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	obj["tls"] = *in.TLS

	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}

	return []interface{}{obj}

}

// Expanders

func expandNotifierSMTPConfig(p []interface{}) *managementClient.SMTPConfig {
	obj := &managementClient.SMTPConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.DefaultRecipient = in["default_recipient"].(string)
	obj.Host = in["host"].(string)
	obj.Port = int64(in["port"].(int))
	obj.Sender = in["sender"].(string)

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["tls"].(bool); ok {
		obj.TLS = &v
	}

	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}

	return obj
}
