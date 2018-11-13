package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func monitoringFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenMonitoring(in *managementClient.MonitoringConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandMonitoring(p []interface{}) (*managementClient.MonitoringConfig, error) {
	obj := &managementClient.MonitoringConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	return obj, nil
}
