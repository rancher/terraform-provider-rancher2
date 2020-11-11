package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
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

	newCatalog, err := meta.(*Config).CreateCatalogV2(clusterID, catalog)
	if err != nil {
		return err
	}
	d.SetId(clusterID + catalogV2ClusterIDsep + newCatalog.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
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

	_, rancherID := splitID(d.Id())
	catalog, err := meta.(*Config).GetCatalogV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Catalog V2 %s not found", name)
			d.SetId("")
			return nil
		}
		return err
	}
	return flattenCatalogV2(d, catalog)
}

func resourceRancher2CatalogV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	catalog := expandCatalogV2(d)
	log.Printf("[INFO] Updating Catalog V2 %s", name)

	_, rancherID := splitID(d.Id())
	newCatalog, err := meta.(*Config).UpdateCatalogV2(clusterID, rancherID, catalog)
	if err != nil {
		return err
	}
	d.SetId(clusterID + catalogV2ClusterIDsep + newCatalog.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
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
	catalog, err := meta.(*Config).GetCatalogV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = meta.(*Config).DeleteCatalogV2(clusterID, catalog)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    catalogV2StateRefreshFunc(meta, clusterID, catalog.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for catalog (%s) to be active: %s", catalog.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// catalogV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Catalog v2.
func catalogV2StateRefreshFunc(meta interface{}, clusterID, catalogID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := meta.(*Config).GetCatalogV2ByID(clusterID, catalogID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		for i := range obj.Status.Conditions {
			if obj.Status.Conditions[i].Type == string(v1.RepoDownloaded) {
				if obj.Status.Conditions[i].Status == "True" {
					return obj, "active", nil
				}
				return nil, "error", fmt.Errorf("%s", obj.Status.Conditions[i].Message)
			}
		}
		return obj, "transitioning", nil
	}
}
