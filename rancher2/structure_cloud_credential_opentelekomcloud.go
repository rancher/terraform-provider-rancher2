package rancher2

import "fmt"

func flattenCloudCredentialOpentelekomcloud(in *opentelekomcloudCredentialConfig, prior []interface{}) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	out := map[string]interface{}{
		"access_key":   in.AccessKey,
		"auth_url":     in.AuthURL,
		"domain_name":  in.DomainName,
		"password":     in.Password,
		"project_id":   in.ProjectID,
		"project_name": in.ProjectName,
		"region":       in.Region,
		"secret_key":   in.SecretKey,
		"username":     in.Username,
	}
	if len(prior) > 0 && prior[0] != nil {
		previous := prior[0].(map[string]interface{})
		// Rancher omits these secrets when reading a cloud credential.
		for _, key := range []string{"password", "secret_key"} {
			if out[key] == "" {
				if value, ok := previous[key].(string); ok {
					out[key] = value
				}
			}
		}
	}
	return []interface{}{out}
}

func expandCloudCredentialOpentelekomcloud(in []interface{}) *opentelekomcloudCredentialConfig {
	out := &opentelekomcloudCredentialConfig{}
	if len(in) == 0 || in[0] == nil {
		return out
	}
	values := in[0].(map[string]interface{})
	out.AccessKey, _ = values["access_key"].(string)
	out.AuthURL, _ = values["auth_url"].(string)
	out.DomainName, _ = values["domain_name"].(string)
	out.Password, _ = values["password"].(string)
	out.ProjectID, _ = values["project_id"].(string)
	out.ProjectName, _ = values["project_name"].(string)
	out.Region, _ = values["region"].(string)
	out.SecretKey, _ = values["secret_key"].(string)
	out.Username, _ = values["username"].(string)
	return out
}

func validateCloudCredentialOpentelekomcloud(in *opentelekomcloudCredentialConfig) error {
	if in != nil {
		if in.AccessKey != "" || in.SecretKey != "" {
			if in.AccessKey != "" && in.SecretKey != "" {
				return nil
			}
			return fmt.Errorf("tcloud_public_credential_config requires access_key and secret_key together; the driver prioritizes access-key authentication")
		}
		if in.Username != "" && in.Password != "" && in.DomainName != "" {
			return nil
		}
	}
	return fmt.Errorf("tcloud_public_credential_config requires access_key and secret_key, or username, password and domain_name")
}
