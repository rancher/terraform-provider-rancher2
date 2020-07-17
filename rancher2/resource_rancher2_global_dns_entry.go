package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2GlobalDNSEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2GlobalDNSEntryCreate,
		Read:   resourceRancher2GlobalDNSEntryRead,
		Update: resourceRancher2GlobalDNSEntryUpdate,
		Delete: resourceRancher2GlobalDNSEntryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2GlobalDNSEntryImport,
		},

		Schema: GlobalDNSEntryFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2GlobalDNSEntryCreate(d *schema.ResourceData, meta interface{}) error {
	globalDNSEntry, err := expandGlobalDNSEntry(d)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Global Role Binding %s", globalDNSEntry.Name)

	newglobalDNSEntry, err := client.GlobalDNS.Create(globalDNSEntry)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    GlobalDNSEntryStateRefreshFunc(client, newglobalDNSEntry.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be created: %s", newglobalDNSEntry.ID, waitErr)
	}

	err = flattenGlobalDNSEntry(d, newglobalDNSEntry)
	if err != nil {
		return err
	}

	return resourceRancher2GlobalDNSEntryRead(d, meta)
}

func resourceRancher2GlobalDNSEntryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Global Role Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNSEntry, err := client.GlobalDNS.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Global Role Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenGlobalDNSEntry(d, globalDNSEntry)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2GlobalDNSEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Global Role Binding ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNSEntry, err := client.GlobalDNS.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newglobalDNSEntry, err := client.GlobalDNS.Update(globalDNSEntry, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    GlobalDNSEntryStateRefreshFunc(client, newglobalDNSEntry.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be updated: %s", newglobalDNSEntry.ID, waitErr)
	}

	return resourceRancher2GlobalDNSEntryRead(d, meta)
}

func resourceRancher2GlobalDNSEntryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Global Role Binding ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNSEntry, err := client.GlobalDNS.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Global Role Binding ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.GlobalDNS.Delete(globalDNSEntry)
	if err != nil {
		return fmt.Errorf("Error removing Global Role Binding: %s", err)
	}

	log.Printf("[DEBUG] Waiting for global role binding (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    GlobalDNSEntryStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global role binding (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// GlobalDNSEntryStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Global Role Binding.
func GlobalDNSEntryStateRefreshFunc(client *managementClient.Client, globalDNSEntryID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalDNS.ByID(globalDNSEntryID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, "active", nil
	}
}
