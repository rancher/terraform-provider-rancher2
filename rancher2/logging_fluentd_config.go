package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	loggingFluentdKind = "fluentd"
)

//Schemas

func fluentdConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"fluent_servers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: fluentServerFields(),
			},
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"compress": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enable_tls": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	return s
}

// Flatteners

func flattenFluentdConfig(in *managementClient.FluentForwarderConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.FluentServers != nil {
		servers, err := flattenFluentServer(in.FluentServers)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["fluent_servers"] = servers
	}

	if len(in.Certificate) > 0 {
		obj["certificate"] = in.Certificate
	}

	obj["compress"] = in.Compress

	obj["enable_tls"] = in.EnableTLS

	return []interface{}{obj}, nil
}

// Expanders

func expandFluentdConfig(p []interface{}) (*managementClient.FluentForwarderConfig, error) {
	obj := &managementClient.FluentForwarderConfig{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["fluent_servers"].([]interface{}); ok && len(v) > 0 {
		servers, err := expandFluentServer(v)
		if err != nil {
			return obj, err
		}
		obj.FluentServers = servers
	}

	if v, ok := in["certificate"].(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in["compress"].(bool); ok {
		obj.Compress = v
	}

	if v, ok := in["enable_tls"].(bool); ok {
		obj.EnableTLS = v
	}

	return obj, nil
}
