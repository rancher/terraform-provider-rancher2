package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const authorizationDefaultMode = "rbac"

var (
	authorizationModes = []string{"rbac", "none"}
)

//Schemas

func clusterRKEConfigAuthorizationFields() map[string]*schema.Schema {
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
