package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2AppImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, appID, err := splitAppID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(appID)
	d.Set("project_id", projectID)

	diag := resourceRancher2AppRead(ctx, d, meta)
	if diag.HasError() {
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
