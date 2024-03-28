package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2NodeDriver() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NodeDriverCreate,
		Read:   resourceRancher2NodeDriverRead,
		Update: resourceRancher2NodeDriverUpdate,
		Delete: resourceRancher2NodeDriverDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NodeDriverImport,
		},
		Schema: nodeDriverFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NodeDriverCreate(d *schema.ResourceData, meta interface{}) error {
	nodeDriver := expandNodeDriver(d)

	log.Printf("[INFO] Creating Node Driver %s", nodeDriver.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newNodeDriver, err := client.NodeDriver.Create(nodeDriver)
	if err != nil {
		return err
	}

	d.SetId(newNodeDriver.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "activating"},
		Target:     []string{"active", "inactive"},
		Refresh:    nodeDriverStateRefreshFunc(client, newNodeDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for node driver (%s) to be created: %s", newNodeDriver.ID, waitErr)
	}

	return resourceRancher2NodeDriverRead(d, meta)
}

func resourceRancher2NodeDriverRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		nodeDriver, err := client.NodeDriver.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Node Driver ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenNodeDriver(d, nodeDriver); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2NodeDriverUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Node Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver, err := client.NodeDriver.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"active":           d.Get("active").(bool),
		"builtin":          d.Get("builtin").(bool),
		"checksum":         d.Get("checksum").(string),
		"description":      d.Get("description").(string),
		"externalId":       d.Get("external_id").(string),
		"name":             d.Get("name").(string),
		"uiUrl":            d.Get("ui_url").(string),
		"url":              d.Get("url").(string),
		"whitelistDomains": toArrayString(d.Get("whitelist_domains").([]interface{})),
		"annotations":      toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":           toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNodeDriver, err := client.NodeDriver.Update(nodeDriver, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "inactive", "downloading", "activating", "deactivating"},
		Target:     []string{"active", "inactive"},
		Refresh:    nodeDriverStateRefreshFunc(client, newNodeDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node driver (%s) to be updated: %s", newNodeDriver.ID, waitErr)
	}

	return resourceRancher2NodeDriverRead(d, meta)
}

func resourceRancher2NodeDriverDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Node Driver ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver, err := client.NodeDriver.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Node Driver ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	if !nodeDriver.Builtin {
		err = client.NodeDriver.Delete(nodeDriver)
		if err != nil {
			return fmt.Errorf("Error removing Node Driver: %s", err)
		}

		log.Printf("[DEBUG] Waiting for node driver (%s) to be removed", id)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"removing"},
			Target:     []string{"removed"},
			Refresh:    nodeDriverStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf(
				"[ERROR] waiting for node driver (%s) to be removed: %s", id, waitErr)
		}
	}

	d.SetId("")
	return nil
}

// nodeDriverStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher NodeDriver.
func nodeDriverStateRefreshFunc(client *managementClient.Client, nodeDriverID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.NodeDriver.ByID(nodeDriverID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
