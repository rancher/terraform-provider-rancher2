package rancher2

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2SecretImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	namespaceID, projectID, resourceID := splitRegistryID(d.Id())

	d.SetId(resourceID)
	d.Set("project_id", projectID)
	d.Set("namespace_id", namespaceID)

	diag := resourceRancher2SecretRead(ctx, d, meta)
	if diag.HasError() {
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
