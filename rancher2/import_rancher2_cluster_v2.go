package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2ClusterV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diag := resourceRancher2ClusterV2Read(ctx, d, meta)
	if diag.HasError() || d.Id() == "" { // TODO - Check this if, it can break, didn't change the logic,.
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
