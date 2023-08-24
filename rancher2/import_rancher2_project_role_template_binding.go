package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2ProjectRoleTemplateBindingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	projectRole, err := client.ProjectRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenProjectRoleTemplateBinding(d, projectRole)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
