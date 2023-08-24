package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2SecretV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, _ := splitID(d.Id())
	d.Set("cluster_id", clusterID)

	diag := resourceRancher2SecretV2Read(ctx, d, meta)
	if diag.HasError() {
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}
	if d.Id() == "" {
		return []*schema.ResourceData{}, nil
	}

	return []*schema.ResourceData{d}, nil
}
