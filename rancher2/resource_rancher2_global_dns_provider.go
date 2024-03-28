package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2GlobalDNSProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2GlobalDNSProviderCreate,
		Read:   resourceRancher2GlobalDNSProviderRead,
		Update: resourceRancher2GlobalDNSProviderUpdate,
		Delete: resourceRancher2GlobalDNSProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2GlobalDNSProviderImport,
		},

		Schema: globalDNSProviderFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2GlobalDNSProviderCreate(d *schema.ResourceData, meta interface{}) error {
	globalDNSProvider := expandGlobalDNSProvider(d)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Global DNS Provider %s", globalDNSProvider.Name)

	newGlobalDNSProvider, err := client.GlobalDnsProvider.Create(globalDNSProvider)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, newGlobalDNSProvider.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be created: %s", newGlobalDNSProvider.ID, waitErr)
	}

	err = flattenGlobalDNSProvider(d, newGlobalDNSProvider)
	if err != nil {
		return err
	}

	return resourceRancher2GlobalDNSProviderRead(d, meta)
}

func resourceRancher2GlobalDNSProviderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Global DNS Provider ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		globalDNSProvider, err := client.GlobalDnsProvider.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) {
				log.Printf("[INFO] Global DNS Provider ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenGlobalDNSProvider(d, globalDNSProvider); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2GlobalDNSProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Global DNS Provider ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNSProvider, err := client.GlobalDnsProvider.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newGlobalDNSProvider, err := client.GlobalDnsProvider.Update(globalDNSProvider, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, newGlobalDNSProvider.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be updated: %s", newGlobalDNSProvider.ID, waitErr)
	}

	return resourceRancher2GlobalDNSProviderRead(d, meta)
}

func resourceRancher2GlobalDNSProviderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Global DNS Provider ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNSProvider, err := client.GlobalDnsProvider.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Global DNS Provider ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.GlobalDnsProvider.Delete(globalDNSProvider)
	if err != nil {
		return fmt.Errorf("Error removing Global DNS Provider: %s", err)
	}

	log.Printf("[DEBUG] Waiting for global dns provider (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// globalDNSProviderStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Global DNS Provider.
func globalDNSProviderStateRefreshFunc(client *managementClient.Client, globalDNSProviderID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalDnsProvider.ByID(globalDNSProviderID)
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
