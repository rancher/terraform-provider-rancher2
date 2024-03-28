package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ClusterDriver() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterDriverCreate,
		Read:   resourceRancher2ClusterDriverRead,
		Update: resourceRancher2ClusterDriverUpdate,
		Delete: resourceRancher2ClusterDriverDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterDriverImport,
		},
		Schema: clusterDriverFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterDriverCreate(d *schema.ResourceData, meta interface{}) error {
	clusterDriver := expandClusterDriver(d)

	log.Printf("[INFO] Creating Cluster Driver %s", clusterDriver.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newClusterDriver, err := client.KontainerDriver.Create(clusterDriver)
	if err != nil {
		return err
	}

	d.SetId(newClusterDriver.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "activating"},
		Target:     []string{"active", "inactive"},
		Refresh:    clusterDriverStateRefreshFunc(client, newClusterDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster driver (%s) to be created: %s", newClusterDriver.ID, waitErr)
	}

	return resourceRancher2ClusterDriverRead(d, meta)
}

func resourceRancher2ClusterDriverRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		clusterDriver, err := client.KontainerDriver.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Driver ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}

			resource.NonRetryableError(err)
		}

		if err = flattenClusterDriver(d, clusterDriver); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterDriverUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterDriver, err := client.KontainerDriver.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"active":           d.Get("active").(bool),
		"actualUrl":        d.Get("actual_url").(string),
		"builtin":          d.Get("builtin").(bool),
		"checksum":         d.Get("checksum").(string),
		"name":             d.Get("name").(string),
		"uiUrl":            d.Get("ui_url").(string),
		"url":              d.Get("url").(string),
		"whitelistDomains": toArrayString(d.Get("whitelist_domains").([]interface{})),
		"annotations":      toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":           toMapString(d.Get("labels").(map[string]interface{})),
	}

	newClusterDriver, err := client.KontainerDriver.Update(clusterDriver, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "inactive", "downloading", "activating", "deactivating"},
		Target:     []string{"active", "inactive"},
		Refresh:    clusterDriverStateRefreshFunc(client, newClusterDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster driver (%s) to be updated: %s", newClusterDriver.ID, waitErr)
	}

	return resourceRancher2ClusterDriverRead(d, meta)
}

func resourceRancher2ClusterDriverDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Driver ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterDriver, err := client.KontainerDriver.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster Driver ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	if !clusterDriver.BuiltIn {
		err = client.KontainerDriver.Delete(clusterDriver)
		if err != nil {
			return fmt.Errorf("Error removing Cluster Driver: %s", err)
		}

		log.Printf("[DEBUG] Waiting for cluster driver (%s) to be removed", id)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"removing"},
			Target:     []string{"removed"},
			Refresh:    clusterDriverStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf(
				"[ERROR] waiting for cluster driver (%s) to be removed: %s", id, waitErr)
		}
	}

	d.SetId("")
	return nil
}

// clusterDriverStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ClusterDriver.
func clusterDriverStateRefreshFunc(client *managementClient.Client, clusterDriverID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.KontainerDriver.ByID(clusterDriverID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
