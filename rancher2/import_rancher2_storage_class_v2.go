package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2StorageClassV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, name := splitID(d.Id())
	d.Set("cluster_id", clusterID)
	d.Set("name", name)

	diag := resourceRancher2StorageClassV2Read(ctx, d, meta)
	if diag.HasError() || d.Id() == "" { // TODO - Is the same as the old one, but can break if diag == nil FELIPE - ANDY
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
