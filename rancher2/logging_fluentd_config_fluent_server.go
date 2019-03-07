package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func fluentServerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"shared_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"standby": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"username": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	return s
}

// Flatteners

func flattenFluentServer(data []managementClient.FluentServer) ([]interface{}, error) {
	result := []interface{}{}

	for _, in := range data {
		obj := make(map[string]interface{})

		if len(in.Endpoint) > 0 {
			obj["endpoint"] = in.Endpoint
		}

		if len(in.Hostname) > 0 {
			obj["hostname"] = in.Hostname
		}

		if len(in.Password) > 0 {
			obj["password"] = in.Password
		}

		if len(in.SharedKey) > 0 {
			obj["shared_key"] = in.SharedKey
		}

		obj["standby"] = in.Standby

		if len(in.Username) > 0 {
			obj["username"] = in.Username
		}

		if in.Weight > 0 {
			obj["weight"] = int(in.Weight)
		}

		result = append(result, obj)
	}

	return result, nil
}

// Expanders

func expandFluentServer(p []interface{}) ([]managementClient.FluentServer, error) {
	result := []managementClient.FluentServer{}

	if len(p) == 0 || p[0] == nil {
		return result, nil
	}

	for i := range p {
		obj := managementClient.FluentServer{}
		in := p[i].(map[string]interface{})

		if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
			obj.Endpoint = v
		}

		if v, ok := in["hostname"].(string); ok && len(v) > 0 {
			obj.Hostname = v
		}

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			obj.Password = v
		}

		if v, ok := in["shared_key"].(string); ok && len(v) > 0 {
			obj.SharedKey = v
		}

		if v, ok := in["standby"].(bool); ok {
			obj.Standby = v
		}

		if v, ok := in["username"].(string); ok && len(v) > 0 {
			obj.Username = v
		}

		if v, ok := in["weight"].(int); ok {
			obj.Weight = int64(v)
		}

		result = append(result, obj)
	}
	return result, nil
}
