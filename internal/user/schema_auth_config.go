package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const authDefaultAccessMode = "unrestricted"

var (
	authAccessModes = []string{"required", "restricted", "unrestricted"}
)

//Schemas

func authConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
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
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}
	return s
}
