package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func multiClusterAppFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"catalog_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Multi cluster app catalog name",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Multi cluster app name",
		},
		"roles": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Multi cluster app roles",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"targets": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Multi cluster app targets",
			Elem: &schema.Resource{
				Schema: targetFields(),
			},
		},
		"template_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Multi cluster app template name",
		},
		"answers": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Multi cluster app answers",
			Elem: &schema.Resource{
				Schema: answerFields(),
			},
		},
		"members": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Multi cluster app members",
			Elem: &schema.Resource{
				Schema: memberFields(),
			},
		},
		"revision_history_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10,
			Description: "Multi cluster app revision history limit",
		},
		"revision_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Multi cluster app revision name",
		},
		"template_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Multi cluster app template version",
		},
		"template_version_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Multi cluster app template version ID",
		},
		"upgrade_strategy": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Multi cluster app upgrade strategy",
			Elem: &schema.Resource{
				Schema: upgradeStrategyFields(),
			},
		},
		"wait": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Wait until multi cluster app is active",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
