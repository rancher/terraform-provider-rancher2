package rancher2

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const authDefaultAccessMode = "unrestricted"

var (
	authAccessModes = []string{"required", "restricted", "unrestricted"}
)

//Schemas

func authConfigFields(defaultName string) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// The name identifies the AuthConfig in Rancher, so renaming requires a new one.
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  defaultName,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"access_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(authAccessModes, true),
			Default:      authDefaultAccessMode,
		},
		"allowed_principal_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	maps.Copy(s, commonAnnotationLabelFields())

	return s
}
