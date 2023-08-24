package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2ClusterAlertGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diag := resourceRancher2ClusterAlertGroupRead(ctx, d, meta)
	if diag.HasError() {
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
