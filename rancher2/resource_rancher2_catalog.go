package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "helm",
				ValidateFunc: validation.StringInSlice([]string{"helm"}, true),
			},
			"branch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
		},
	}
}

func resourceRancher2CatalogCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	log.Printf("[INFO] Creating Catalog %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	catalog := &managementClient.Catalog{
		Name:        name,
		URL:         d.Get("url").(string),
		Description: d.Get("description").(string),
		Kind:        d.Get("kind").(string),
		Branch:      d.Get("branch").(string),
	}

	newCatalog, err := client.Catalog.Create(catalog)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    CatalogStateRefreshFunc(client, newCatalog.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for catalog (%s) to be created: %s", newCatalog.ID, waitErr)
	}

	d.SetId(newCatalog.ID)

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

	d.Set("name", catalog.Name)
	d.Set("url", catalog.URL)
	d.Set("description", catalog.Description)
	d.Set("kind", catalog.Kind)
	d.Set("branch", catalog.Branch)

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

	update := map[string]string{
		"url":         d.Get("url").(string),
		"description": d.Get("description").(string),
		"kind":        d.Get("kind").(string),
		"branch":      d.Get("branch").(string),
	}

	newCatalog, err := client.Catalog.Update(catalog, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"refreshed"},
		Target:     []string{"active"},
		Refresh:    CatalogStateRefreshFunc(client, newCatalog.ID),
		Timeout:    10 * time.Minute,
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
		Refresh:    CatalogStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
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

func resourceRancher2CatalogImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	catalog, err := client.Catalog.ByID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(catalog.ID)
	d.Set("name", catalog.Name)
	d.Set("url", catalog.URL)
	d.Set("description", catalog.Description)
	d.Set("kind", catalog.Kind)
	d.Set("branch", catalog.Branch)

	return []*schema.ResourceData{d}, nil
}

// CatalogStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Catalog.
func CatalogStateRefreshFunc(client *managementClient.Client, catalogID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cat, err := client.Catalog.ByID(catalogID)
		if err != nil {
			if IsNotFound(err) {
				return cat, "removed", nil
			}
			return nil, "", err
		}

		return cat, cat.State, nil
	}
}
