package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigPrivateRegistriesECRCredentials(in *managementClient.ECRCredentialPlugin, p []interface{}) []interface{} {
	if in == nil {
		return nil
	}

	var obj map[string]interface{}
	if p == nil || len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if len(in.AwsAccessKeyID) > 0 {
		obj["aws_access_key_id"] = in.AwsAccessKeyID
	}
	if len(in.AwsSecretAccessKey) > 0 {
		obj["aws_secret_access_key"] = in.AwsSecretAccessKey
	}
	if len(in.AwsSessionToken) > 0 {
		obj["aws_session_token"] = in.AwsSessionToken
	}

	return []interface{}{obj}
}

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

		if in.ECRCredentialPlugin != nil {
			ecrCredential := []interface{}{}
			if objECR, ok := obj["ecr_credential_plugin"].([]interface{}); ok {
				ecrCredential = objECR
			}
			obj["ecr_credential_plugin"] = flattenClusterRKEConfigPrivateRegistriesECRCredentials(in.ECRCredentialPlugin, ecrCredential)
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

func expandClusterRKEConfigPrivateRegistriesECRCredentials(p []interface{}) *managementClient.ECRCredentialPlugin {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	obj := &managementClient.ECRCredentialPlugin{}
	in := p[0].(map[string]interface{})

	if v, ok := in["aws_access_key_id"].(string); ok && len(v) > 0 {
		obj.AwsAccessKeyID = v
	}
	if v, ok := in["aws_secret_access_key"].(string); ok && len(v) > 0 {
		obj.AwsSecretAccessKey = v
	}
	if v, ok := in["aws_session_token"].(string); ok && len(v) > 0 {
		obj.AwsSessionToken = v
	}
	return obj
}

func expandClusterRKEConfigPrivateRegistries(p []interface{}) ([]managementClient.PrivateRegistry, error) {
	out := make([]managementClient.PrivateRegistry, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.PrivateRegistry{}

		if v, ok := in["ecr_credential_plugin"].([]interface{}); ok && len(v) > 0 {
			obj.ECRCredentialPlugin = expandClusterRKEConfigPrivateRegistriesECRCredentials(v)
		}

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
