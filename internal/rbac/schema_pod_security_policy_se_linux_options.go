package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicySELinuxOptionsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"level": {
			Type:        schema.TypeString,
			Description: "Level is SELinux level label that applies to the container.",
			Optional:    true,
		},
		"role": {
			Type:        schema.TypeString,
			Description: "Role is a SELinux role label that applies to the container.",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type is a SELinux type label that applies to the container.",
			Optional:    true,
		},
		"user": {
			Type:        schema.TypeString,
			Description: "User is a SELinux user label that applies to the container.",
			Optional:    true,
		},
	}

	return s
}
