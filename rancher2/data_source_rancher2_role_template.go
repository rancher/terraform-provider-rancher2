package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceRancher2RoleTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2RoleTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role template policy name",
			},
			"context": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(roleTemplateContexts, true),
				Description:  "Context role template",
			},
			"administrative": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Administrative role template",
			},
			"builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Builtin role template",
			},
			"default_role": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Default role template for new created cluster or project",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role template policy description",
			},
			"external": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "External role template",
			},
			"hidden": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Hidden role template",
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Locked role template",
			},
			"role_template_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Inherit role template IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Role template policy rules",
				Elem: &schema.Resource{
					Schema: policyRuleFields(),
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the role template",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the role template",
			},
		},
	}
}

func dataSourceRancher2RoleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	context := d.Get("context").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	if len(context) > 0 {
		filters["context"] = context
	}
	listOpts := NewListOpts(filters)

	roleTemplates, err := client.RoleTemplate.List(listOpts)
	if err != nil {
		return err
	}

	count := len(roleTemplates.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] role template with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d role template with name \"%s\"", count, name)
	}

	return flattenRoleTemplate(d, &roleTemplates.Data[0])
}
