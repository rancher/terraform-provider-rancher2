package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Pod Security Policy template policy name",
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Pod Security Policy template policy description",
		},
		"spec": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Pod Security Policy template spec",
			Elem: &schema.Resource{
				Schema: podSecurityPolicySpecFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the Pod Security Policy template",
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the Pod Security Policy template",
		},
	}

	return s
}