package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	roleTemplateContextCluster = "cluster"
	roleTemplateContextProject = "project"
)

var (
	roleTemplateContexts = []string{roleTemplateContextCluster, roleTemplateContextProject}
)

//Schemas

func roleTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Role template policy name",
		},
		"administrative": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Administrative role template",
		},
		"builtin": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Builtin role template",
		},
		"context": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      roleTemplateContextCluster,
			ValidateFunc: validation.StringInSlice(roleTemplateContexts, true),
			Description:  "Context role template",
		},
		"default_role": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Default role template for new created cluster or project",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Role template policy description",
		},
		"external": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "External role template",
		},
		"hidden": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Hidden role template",
		},
		"locked": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Locked role template",
		},
		"role_template_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Inherit role template IDs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"rules": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Role template policy rules",
			Elem: &schema.Resource{
				Schema: policyRuleFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
