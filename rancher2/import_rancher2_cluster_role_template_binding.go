package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2ClusterRoleTemplateBindingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	clusterRole, err := client.ClusterRoleTemplateBinding.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenClusterRoleTemplateBinding(d, clusterRole)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
