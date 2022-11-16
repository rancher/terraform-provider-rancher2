package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"golang.org/x/sync/errgroup"
)

func resourceRancher2CatalogV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2CatalogV2Create,
		Read:   resourceRancher2CatalogV2Read,
		Update: resourceRancher2CatalogV2Update,
		Delete: resourceRancher2CatalogV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2CatalogV2Import,
		},
		Schema: catalogV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CatalogV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	catalog := expandCatalogV2(d)

	log.Printf("[INFO] Creating Catalog V2 %s", name)

	newCatalog, err := createCatalogV2(meta.(*Config), clusterID, catalog)
	if err != nil {
		return err
	}
	d.SetId(clusterID + catalogV2ClusterIDsep + newCatalog.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"downloaded"},
		Refresh:    catalogV2StateRefreshFunc(meta, clusterID, newCatalog.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for catalog (%s) to be active: %s", newCatalog.ID, waitErr)
	}
	return resourceRancher2CatalogV2Read(d, meta)
}

func resourceRancher2CatalogV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Catalog V2 %s", name)

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		_, rancherID := splitID(d.Id())
		catalog, err := getCatalogV2ByID(meta.(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Catalog V2 %s not found", name)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenCatalogV2(d, catalog); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2CatalogV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	catalog := expandCatalogV2(d)
	log.Printf("[INFO] Updating Catalog V2 %s", name)

	_, rancherID := splitID(d.Id())
	newCatalog, err := updateCatalogV2(meta.(*Config), clusterID, rancherID, catalog)
	if err != nil {
		return err
	}
	d.SetId(clusterID + catalogV2ClusterIDsep + newCatalog.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"downloaded"},
		Refresh:    catalogV2StateRefreshFunc(meta, clusterID, newCatalog.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for catalog (%s) to be active: %s", newCatalog.ID, waitErr)
	}
	return resourceRancher2CatalogV2Read(d, meta)
}

func resourceRancher2CatalogV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Catalog V2 %s", name)

	_, rancherID := splitID(d.Id())
	catalog, err := getCatalogV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteCatalogV2(meta.(*Config), clusterID, catalog)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    catalogV2StateRefreshFunc(meta, clusterID, catalog.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for catalog (%s) to be deleted: %s", catalog.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// catalogV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Catalog v2.
func catalogV2StateRefreshFunc(meta interface{}, clusterID, catalogID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getCatalogV2ByID(meta.(*Config), clusterID, catalogID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		for i := range obj.Status.Conditions {
			if obj.Status.Conditions[i].Type == string(v1.RepoDownloaded) {
				if obj.Status.Conditions[i].Status == "Unknown" {
					return obj, "transitioning", nil
				}
				if obj.Status.Conditions[i].Status == "True" {
					return obj, "downloaded", nil
				}
				return nil, "error", fmt.Errorf("%s", obj.Status.Conditions[i].Message)
			}
		}
		return obj, "transitioning", nil
	}
}

// Rancher2 Catalog V2 API CRUD functions
func createCatalogV2(c *Config, clusterID string, repo *ClusterRepo) (*ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating catalog V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("Creating catalog V2: Cluster ID is nil")
	}
	if repo == nil {
		return nil, fmt.Errorf("Creating catalog V2: object is nil")
	}
	resp := &ClusterRepo{}
	err := c.createObjectV2(clusterID, catalogV2APIType, repo, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating Catalog V2: %s", err)
	}
	return resp, nil
}

func deleteCatalogV2(c *Config, clusterID string, obj *ClusterRepo) error {
	if c == nil {
		return fmt.Errorf("Deleting catalog V2: Provider config is nil")
	}
	if clusterID == "" {
		return fmt.Errorf("Deleting catalog V2: Cluster ID is nil")
	}
	if obj == nil {
		return fmt.Errorf("Deleting catalog V2: Catalog V2 is nil")
	}

	resource := &types.Resource{
		ID:      obj.ID,
		Type:    obj.Type,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(clusterID, resource)
}

func getCatalogV2ByID(c *Config, clusterID, id string) (*ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting catalog V2: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Getting catalog V2: Cluster ID and/or Catalog V2 ID is nil")
	}
	resp := &ClusterRepo{}
	err := c.getObjectV2ByID(clusterID, id, catalogV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting Catalog V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateCatalogV2(c *Config, clusterID, id string, obj *ClusterRepo) (*ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating catalog V2: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Updating catalog V2: Cluster ID and/or Catalog V2 ID is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating catalog V2: Cluster V2 is nil")
	}
	resp := &ClusterRepo{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(clusterID, id, catalogV2APIType, obj, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read clusterRepo again and update ObjectMeta.ResourceVersion before retry
			newObj := &ClusterRepo{}
			err = c.getObjectV2ByID(clusterID, id, catalogV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating catalog V2 ID %s: %v", id, err)
		}
	}
}

func getCatalogV2List(c *Config, clusterID string) ([]ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting catalog V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("Cluster v2 ID is nil")
	}
	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	listOpts := NewListOpts(nil)
	resp := &ClusterRepoCollection{}
	for {
		err = client.List(catalogV2APIType, listOpts, resp)
		if err == nil {
			return resp.Data, nil
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsNotFound(err) {
			return nil, err
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout getting catalog V2 list at cluster ID %s: %v", clusterID, err)
		}
	}
}

func waitCatalogV2Downloaded(c *Config, clusterID, catalogID string) (*ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Waiting for catalog V2: Provider config is nil")
	}
	if clusterID == "" || catalogID == "" {
		return nil, fmt.Errorf("Waiting for catalog V2: Cluster ID and/or Catalog V2 ID is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		obj, err := getCatalogV2ByID(c, clusterID, catalogID)
		if err != nil {
			return nil, fmt.Errorf("Getting catalog V2 ID (%s): %v", catalogID, err)
		}
		for i := range obj.Status.Conditions {
			if obj.Status.Conditions[i].Type == string(v1.RepoDownloaded) {
				// Status of the condition, one of True, False, Unknown.
				if obj.Status.Conditions[i].Status == "Unknown" {
					break
				}
				if obj.Status.Conditions[i].Status == "True" {
					return obj, nil
				}
				return nil, fmt.Errorf("Catalog V2 ID %s: %s", catalogID, obj.Status.Conditions[i].Message)
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("[ERROR] Timeout waiting for catalog V2 ID %s at cluster ID (%s): %v", catalogID, clusterID, err)
		}
	}
}

func waitAllCatalogV2Downloaded(c *Config, clusterID string) ([]ClusterRepo, error) {
	if c == nil {
		return nil, fmt.Errorf("Waiting for all catalogs V2: Provider config is nil")
	}
	clusterRepos, err := getCatalogV2List(c, clusterID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] getting catalog V2 list at cluster ID (%s): %s", clusterID, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, clusterRepo := range clusterRepos {
		repoID := clusterRepo.ID
		g.Go(func() error {
			repo := repoID
			_, err := waitCatalogV2Downloaded(c, clusterID, repo)
			if err != nil {
				return err
			}
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return clusterRepos, fmt.Errorf("[ERROR] waiting for all catalogs V2 to be active at cluster ID (%s): %s", clusterID, err)
	}

	return clusterRepos, nil
}
