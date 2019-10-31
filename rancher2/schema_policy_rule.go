package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	policyRuleAll        = "*"
	policyRuleVerbCreate = "create"
	policyRuleVerbDelete = "delete"
	policyRuleVerbGet    = "get"
	policyRuleVerbList   = "list"
	policyRuleVerbPatch  = "patch"
	policyRuleVerbUpdate = "update"
	policyRuleVerbWatch  = "watch"
)

var (
	policyRuleVerbs = []string{
		policyRuleAll,
		policyRuleVerbCreate,
		policyRuleVerbDelete,
		policyRuleVerbGet,
		policyRuleVerbList,
		policyRuleVerbPatch,
		policyRuleVerbUpdate,
		policyRuleVerbWatch,
	}
)

//Schemas

func policyRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Policy rule api groups",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"non_resource_urls": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Policy rule non resource urls",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"resource_names": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Policy rule resource names",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"resources": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Policy rule resources",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"verbs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Policy rule verbs",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(policyRuleVerbs, true),
			},
		},
	}

	return s
}
