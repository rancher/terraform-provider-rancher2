package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func ingressFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"node_selector": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
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

func flattenIngress(in *managementClient.IngressConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
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

func expandIngress(p []interface{}) (*managementClient.IngressConfig, error) {
	obj := &managementClient.IngressConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	return obj, nil
}
