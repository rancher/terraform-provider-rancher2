package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ProjectRoleTemplateBinding() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ProjectRoleTemplateBindingRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ProjectRoleTemplateBindingRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	roleTemplateID := d.Get("role_template_id").(string)

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}
	if len(roleTemplateID) > 0 {
		filters["roleTemplateId"] = roleTemplateID
	}
	listOpts := NewListOpts(filters)

	projectRoleTemplateBindings, err := client.ProjectRoleTemplateBinding.List(listOpts)
	if err != nil {
		return err
	}

	count := len(projectRoleTemplateBindings.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] project role template binding with name \"%s\" on project ID \"%s\" not found", name, projectID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d project role template binding with name \"%s\" on project ID \"%s\"", count, name, projectID)
	}

	return flattenProjectRoleTemplateBinding(d, &projectRoleTemplateBindings.Data[0])
}
