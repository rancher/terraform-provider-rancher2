package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func globalRoleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"builtin": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Builtin global role",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Global role policy description",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Global role policy name",
		},
		"new_user_default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether or not this role should be added to new users",
		},
		"rules": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Global role policy rules",
			Elem: &schema.Resource{
				Schema: policyRuleFields(),
			},
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global role uuid",
		},
		"inherited_cluster_roles": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Names of role templates whose permissions are granted by this global role in every cluster besides the local cluster",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
