package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2Catalog() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2CatalogCreate,
		Read:   resourceRancher2CatalogRead,
		Update: resourceRancher2CatalogUpdate,
		Delete: resourceRancher2CatalogDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2CatalogImport,
		},

		Schema: catalogFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CatalogCreate(d *schema.ResourceData, meta interface{}) error {
	catalog := expandCatalog(d)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Catalog %s", catalog.Name)

	newCatalog, err := client.Catalog.Create(catalog)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    catalogStateRefreshFunc(client, newCatalog.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for catalog (%s) to be created: %s", newCatalog.ID, waitErr)
	}

	err = flattenCatalog(d, newCatalog)
	if err != nil {
		return err
	}

	return resourceRancher2CatalogRead(d, meta)
}

func resourceRancher2CatalogRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Catalog ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	catalog, err := client.Catalog.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Catalog ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenCatalog(d, catalog)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2CatalogUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Catalog ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	catalog, err := client.Catalog.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"url":         d.Get("url").(string),
		"description": d.Get("description").(string),
		"kind":        d.Get("kind").(string),
		"branch":      d.Get("branch").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newCatalog, err := client.Catalog.Update(catalog, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    catalogStateRefreshFunc(client, newCatalog.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for catalog (%s) to be updated: %s", newCatalog.ID, waitErr)
	}

	return resourceRancher2CatalogRead(d, meta)
}

func resourceRancher2CatalogDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting catalog ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	catalog, err := client.Catalog.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Catalog ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Catalog.Delete(catalog)
	if err != nil {
		return fmt.Errorf("Error removing Catalog: %s", err)
	}

	log.Printf("[DEBUG] Waiting for catalog (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    catalogStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for catalog (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// catalogStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Catalog.
func catalogStateRefreshFunc(client *managementClient.Client, catalogID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Catalog.ByID(catalogID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.State, nil
	}
}
