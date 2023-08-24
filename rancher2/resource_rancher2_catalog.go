package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2Catalog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2CatalogCreate,
		ReadContext:   resourceRancher2CatalogRead,
		UpdateContext: resourceRancher2CatalogUpdate,
		DeleteContext: resourceRancher2CatalogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2CatalogImport,
		},

		Schema: catalogFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CatalogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope := d.Get("scope").(string)
	name := d.Get("name").(string)
	catalog := expandCatalog(d)

	log.Printf("[INFO] Creating %s Catalog %s", scope, name)

	newCatalog, err := meta.(*Config).CreateCatalog(scope, catalog)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ""
	switch scope {
	case catalogScopeCluster:
		id = newCatalog.(*managementClient.ClusterCatalog).ID
	case catalogScopeGlobal:
		id = newCatalog.(*managementClient.Catalog).ID
	case catalogScopeProject:
		id = newCatalog.(*managementClient.ProjectCatalog).ID
	}

	d.SetId(id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    catalogStateRefreshFunc(meta, id, scope),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for catalog (%s) to be created: %s", id, waitErr)
	}

	return resourceRancher2CatalogRead(ctx, d, meta)
}

func resourceRancher2CatalogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope := d.Get("scope").(string)
	id := d.Id()
	log.Printf("[INFO] Refreshing %s Catalog ID %s", scope, id)

	catalog, err := meta.(*Config).GetCatalog(id, scope)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] %s Catalog ID %s not found.", scope, id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if d.Get("refresh").(bool) {
		_, err := meta.(*Config).RefreshCatalog(scope, catalog)
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"refreshed"},
			Target:     []string{"active"},
			Refresh:    catalogStateRefreshFunc(meta, id, scope),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr := stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			return diag.Errorf(
				"[ERROR] waiting for catalog (%s) to be refreshed: %s", id, waitErr)
		}
	}

	return diag.FromErr(flattenCatalog(d, catalog))
}

func resourceRancher2CatalogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope := d.Get("scope").(string)
	id := d.Id()
	log.Printf("[INFO] Updating %s Catalog ID %s", scope, id)

	catalog, err := meta.(*Config).GetCatalog(id, scope)
	if err != nil {
		return diag.FromErr(err)
	}

	update := map[string]interface{}{
		"url":         d.Get("url").(string),
		"branch":      d.Get("branch").(string),
		"description": d.Get("description").(string),
		"kind":        d.Get("kind").(string),
		"password":    d.Get("password").(string),
		"username":    d.Get("username").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	_, err = meta.(*Config).UpdateCatalog(scope, catalog, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    catalogStateRefreshFunc(meta, id, scope),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for %s catalog (%s) to be updated: %s", scope, id, waitErr)
	}

	return resourceRancher2CatalogRead(ctx, d, meta)
}

func resourceRancher2CatalogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope := d.Get("scope").(string)
	id := d.Id()
	log.Printf("[INFO] Deleting %s catalog ID %s", scope, id)
	catalog, err := meta.(*Config).GetCatalog(id, scope)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] %s Catalog ID %s not found.", scope, id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = meta.(*Config).DeleteCatalog(scope, catalog)
	if err != nil {
		return diag.Errorf("Error removing %s Catalog: %s", scope, err)
	}

	log.Printf("[DEBUG] Waiting for %s catalog (%s) to be removed", scope, id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    catalogStateRefreshFunc(meta, id, scope),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for %s catalog (%s) to be removed: %s", scope, id, waitErr)
	}

	d.SetId("")
	return nil
}

// catalogStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Catalog.
func catalogStateRefreshFunc(meta interface{}, catalogID, scope string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := meta.(*Config).GetCatalog(catalogID, scope)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		var state string

		switch scope {
		case catalogScopeCluster:
			state = obj.(*managementClient.ClusterCatalog).State
		case catalogScopeGlobal:
			state = obj.(*managementClient.Catalog).State
		case catalogScopeProject:
			state = obj.(*managementClient.ProjectCatalog).State
		}

		return obj, state, nil
	}
}
