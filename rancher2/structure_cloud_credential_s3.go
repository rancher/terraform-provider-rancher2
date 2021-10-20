package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	s3ConfigDriver = "s3"
)

// Flatteners

func flattenCloudCredentialS3(in *managementClient.S3CredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}
	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}
	if len(in.DefaultBucket) > 0 {
		obj["default_bucket"] = in.DefaultBucket
	}
	if len(in.DefaultEndpoint) > 0 {
		obj["default_endpoint"] = in.DefaultEndpoint
	}
	if len(in.DefaultEndpointCA) > 0 {
		obj["default_endpoint_ca"] = in.DefaultEndpointCA
	}
	if len(in.DefaultFolder) > 0 {
		obj["default_folder"] = in.DefaultFolder
	}
	if len(in.DefaultRegion) > 0 {
		obj["default_region"] = in.DefaultRegion
	}
	if len(in.DefaultSkipSSLVerify) > 0 {
		if ToLower(in.DefaultSkipSSLVerify) == "true" {
			obj["default_skip_ssl_verify"] = true
		} else {
			obj["default_skip_ssl_verify"] = false
		}
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialS3(p []interface{}) *managementClient.S3CredentialConfig {
	obj := &managementClient.S3CredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}
	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}
	if v, ok := in["default_bucket"].(string); ok && len(v) > 0 {
		obj.DefaultBucket = v
	}
	if v, ok := in["default_endpoint"].(string); ok && len(v) > 0 {
		obj.DefaultEndpoint = v
	}
	if v, ok := in["default_endpoint_ca"].(string); ok && len(v) > 0 {
		obj.DefaultEndpointCA = v
	}
	if v, ok := in["default_folder"].(string); ok && len(v) > 0 {
		obj.DefaultFolder = v
	}
	if v, ok := in["default_region"].(string); ok && len(v) > 0 {
		obj.DefaultRegion = v
	}
	if v, ok := in["default_skip_ssl_verify"].(bool); ok {
		if v {
			obj.DefaultSkipSSLVerify = "true"
		} else {
			obj.DefaultSkipSSLVerify = "false"
		}
	}

	return obj
}
