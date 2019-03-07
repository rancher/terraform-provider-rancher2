package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	loggingSplunkKind = "splunk"
)

//Schemas

func splunkConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"token": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_cert": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"client_key_pass": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"index": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"source": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssl_verify": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

// Flatteners

func flattenSplunkConfig(in *managementClient.SplunkConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["endpoint"] = in.Endpoint
	obj["token"] = in.Token

	if len(in.Certificate) > 0 {
		obj["certificate"] = in.Certificate
	}

	if len(in.ClientCert) > 0 {
		obj["client_cert"] = in.ClientCert
	}

	if len(in.ClientKey) > 0 {
		obj["client_key"] = in.ClientKey
	}

	if len(in.ClientKeyPass) > 0 {
		obj["client_key_pass"] = in.ClientKeyPass
	}

	if len(in.Index) > 0 {
		obj["index"] = in.Index
	}

	if len(in.Source) > 0 {
		obj["source"] = in.Source
	}

	obj["ssl_verify"] = in.SSLVerify

	return []interface{}{obj}, nil
}

const (
	splunkKind = "splunk"
)

// Expanders

func expandSplunkConfig(p []interface{}) (*managementClient.SplunkConfig, error) {
	obj := &managementClient.SplunkConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	if v, ok := in["certificate"].(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in["client_cert"].(string); ok && len(v) > 0 {
		obj.ClientCert = v
	}

	if v, ok := in["client_key"].(string); ok && len(v) > 0 {
		obj.ClientKey = v
	}

	if v, ok := in["client_key_pass"].(string); ok && len(v) > 0 {
		obj.ClientKeyPass = v
	}

	if v, ok := in["index"].(string); ok && len(v) > 0 {
		obj.Index = v
	}

	if v, ok := in["source"].(string); ok && len(v) > 0 {
		obj.Source = v
	}

	if v, ok := in["ssl_verify"].(bool); ok {
		obj.SSLVerify = v
	}

	return obj, nil
}
