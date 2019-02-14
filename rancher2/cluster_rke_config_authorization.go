package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

const authorizationDefaultMode = "rbac"

var (
	authorizationModes = []string{"rbac", "none"}
)

//Schemas

func authorizationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      authorizationDefaultMode,
			ValidateFunc: validation.StringInSlice(authorizationModes, true),
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenAuthorization(in *managementClient.AuthzConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Mode) > 0 {
		obj["mode"] = in.Mode
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandAuthorization(p []interface{}) (*managementClient.AuthzConfig, error) {
	obj := &managementClient.AuthzConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["mode"].(string); ok && len(v) > 0 {
		obj.Mode = v
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	return obj, nil
}
