package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2NodePoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	nodePool, err := client.NodePool.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNodePool(d, nodePool)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
