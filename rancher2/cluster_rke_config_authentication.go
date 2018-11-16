package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func authenticationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"sans": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"strategy": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenAuthentication(in *managementClient.AuthnConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.SANs) > 0 {
		obj["sans"] = toArrayInterface(in.SANs)
	}

	if len(in.Strategy) > 0 {
		obj["strategy"] = in.Strategy
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandAuthentication(p []interface{}) (*managementClient.AuthnConfig, error) {
	obj := &managementClient.AuthnConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["sans"].([]interface{}); ok && len(v) > 0 {
		obj.SANs = toArrayString(v)
	}

	if v, ok := in["strategy"].(string); ok && len(v) > 0 {
		obj.Strategy = v
	}

	return obj, nil
}
