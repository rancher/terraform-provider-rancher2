package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Shemas

func globalRoleBindingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"global_role_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"user_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"group_principal_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
