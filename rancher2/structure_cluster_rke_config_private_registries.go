package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigPrivateRegistries(p []managementClient.PrivateRegistry, v []interface{}) ([]interface{}, error) {
	out := make([]interface{}, len(p))
	lenV := len(v)
	for i, in := range p {
		var obj map[string]interface{}
		if lenV <= i {
			obj = make(map[string]interface{})
		} else {
			obj = v[i].(map[string]interface{})
		}
		obj["is_default"] = in.IsDefault

		if len(in.Password) > 0 {
			obj["password"] = in.Password
		}

		if len(in.URL) > 0 {
			obj["url"] = in.URL
		}

		if len(in.User) > 0 {
			obj["user"] = in.User
		}

		out[i] = obj
	}

	return out, nil
}

// Expanders

func expandClusterRKEConfigPrivateRegistries(p []interface{}) ([]managementClient.PrivateRegistry, error) {
	out := make([]managementClient.PrivateRegistry, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.PrivateRegistry{}

		if v, ok := in["is_default"].(bool); ok {
			obj.IsDefault = v
		}

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			obj.Password = v
		}

		if v, ok := in["url"].(string); ok && len(v) > 0 {
			obj.URL = v
		}

		if v, ok := in["user"].(string); ok && len(v) > 0 {
			obj.User = v
		}
		out[i] = obj
	}

	return out, nil
}
