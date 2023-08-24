package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2CatalogImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	scope, id := splitID(d.Id())
	if len(scope) == 0 {
		scope = catalogScopeGlobal
	}

	catalog, err := meta.(*Config).GetCatalog(id, scope)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.Set("scope", scope)

	err = flattenCatalog(d, catalog)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
